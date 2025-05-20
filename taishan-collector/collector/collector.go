package collector

import (
	"collector/config"
	"collector/internal/biz/log"
	"collector/internal/dal"
	"collector/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	RunKafkaPartition    = "RunKafkaPartition"
	UsedPartitions       = "UsedPartitions"
	ReportRunningLockKey = "report:lock:%d"
)

func KafkaConsumer() {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	pubSub := model.SubscribeMsg(RunKafkaPartition)
	partitionCh := pubSub.Channel()

	// 使用锁+for循环主要是为了后续动态调整partition
	for {
		select {
		case c := <-partitionCh:
			partitionId, err := strconv.Atoi(c.Payload)
			if err != nil {
				log.Logger.Error("分区转换失败：", err.Error())
				continue
			}
			log.Logger.Infof("collector收到通知，开始对分区%d 进行消费", partitionId)

			consumer, err := sarama.NewConsumer(strings.Split(config.Conf.Kafka.Address, ","), sarama.NewConfig())
			pc, err := consumer.ConsumePartition(config.Conf.Kafka.Topic, int32(partitionId), sarama.OffsetNewest)
			if err != nil {
				log.Logger.Error("consumer.ConsumePartition error:", err)
				break
			}
			pc.IsPaused()
			go ReceiveMessage(pc, int32(partitionId))
		}
	}
}

func ReceiveMessage(pc sarama.PartitionConsumer, partition int32) {
	defer pc.AsyncClose()
	defer func() {
		if err := recover(); err != nil {
			// 将异常信息打印到日志
			log.Logger.Error("panic错误,err:", err)
		}
	}()

	if pc == nil {
		return
	}

	// 保存所有测试结果
	caseTestResultDataMsg := make(map[int32]*model.StageCaseResult)
	frequency := int64(5)               // 采集频率
	machineMp := make(map[string]int64) // key为机器ip，value为最后采集到的数据时间
	planId, reportId := int32(0), int32(0)
	isEnd := false

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 每5s上报一次数据
		timer := time.NewTicker(time.Duration(frequency) * time.Second)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				if caseTestResultDataMsg == nil {
					log.Logger.Errorf("消费失败，caseTestResultDataMsg is nil")
					continue
				}
				mu.Lock()
				for ip, tim := range machineMp {
					if time.Now().UnixMilli()-tim > 60*1000 {
						log.Logger.Errorf("ip:%s 已1分钟没获取到数据，强制释放", ip)
						delete(machineMp, ip)
						model.ReleaseEngine(ip)
						if len(machineMp) == 0 {
							_, err := model.RDB.SRem(context.Background(), UsedPartitions, partition).Result()
							model.ReportRdb.Del(context.Background(), fmt.Sprintf(ReportRunningLockKey, reportId))
							if err != nil {
								log.Logger.Error("usedPartitions删除失败：", err)
							}
							log.Logger.Infof("partitionID: %d被释放", partition)
							isEnd = true
						}
					}
				}

				for caseId, result := range caseTestResultDataMsg {
					caseTestResultDataMsg[caseId].BaseData.FiftyRequestTimeLineValue = parseRequestLineTime(result.CalcData.FiftyRequestTimeMap)
					caseTestResultDataMsg[caseId].BaseData.NinetyRequestTimeLineValue = parseRequestLineTime(result.CalcData.NinetyRequestTimeMap)
					caseTestResultDataMsg[caseId].BaseData.NinetyFiveRequestTimeLineValue = parseRequestLineTime(result.CalcData.NinetyFiveRequestTimeMap)
					caseTestResultDataMsg[caseId].BaseData.NinetyNineRequestTimeLineValue = parseRequestLineTime(result.CalcData.NinetyNineRequestTimeMap)

					totalCon := int64(0)
					for _, i := range result.CalcData.ActualConcurrencyMap {
						totalCon += i
					}
					caseTestResultDataMsg[caseId].BaseData.ActualConcurrency = totalCon
				}
				if err := dal.BatchInsertTestData(reportId, caseTestResultDataMsg); err != nil {
					log.Logger.Error("测试数据写入influxdb失败：", err)
					continue
				}

				caseTestResultDataMsg = make(map[int32]*model.StageCaseResult)
				mu.Unlock()
				if isEnd {
					return
				}
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range pc.Messages() {
			var resultDataMsg model.PlanTestResultDataMsg
			err := json.Unmarshal(msg.Value, &resultDataMsg)
			if err != nil {
				log.Logger.Errorf("pc.Messages().jsonUnmarshal error: %s, \nvalue: %s", err.Error(), string(msg.Value))
			}

			if planId == 0 || reportId == 0 {
				planId = resultDataMsg.PlanID
				reportId = resultDataMsg.ReportID
			}

			mu.Lock()
			// 记录engine实例最新一条数据的时间
			machineMp[resultDataMsg.MachineIP] = time.Now().UnixMilli()

			for sceneId, sceneResult := range resultDataMsg.SceneResults {
				for caseId, caseResult := range sceneResult.CaseResults {
					if _, ok := caseTestResultDataMsg[caseId]; !ok {
						caseTestResultDataMsg[caseId] = &model.StageCaseResult{
							CaseID:    caseId,
							SceneID:   sceneId,
							SceneType: sceneResult.SceneType,
							BaseData: &model.BaseData{
								Url: caseResult.Url,
							},
							CalcData: &model.CalcData{
								FiftyRequestTimeMap:      make(map[string]float64),
								NinetyRequestTimeMap:     make(map[string]float64),
								NinetyFiveRequestTimeMap: make(map[string]float64),
								NinetyNineRequestTimeMap: make(map[string]float64),
								ActualConcurrencyMap:     make(map[string]int64),
							},
						}
					}
					if caseTestResultDataMsg[caseId].BaseData.StartTime == 0 || caseTestResultDataMsg[caseId].BaseData.StartTime > caseResult.StartTime {
						caseTestResultDataMsg[caseId].BaseData.StartTime = caseResult.StartTime
					}
					if caseTestResultDataMsg[caseId].BaseData.EndTime == 0 || caseTestResultDataMsg[caseId].BaseData.EndTime < caseResult.EndTime {
						caseTestResultDataMsg[caseId].BaseData.EndTime = caseResult.EndTime
					}
					if caseTestResultDataMsg[caseId].BaseData.MaxRequestTime == 0 || caseTestResultDataMsg[caseId].BaseData.MaxRequestTime < caseResult.MaxRequestTime {
						caseTestResultDataMsg[caseId].BaseData.MaxRequestTime = caseResult.MaxRequestTime
					}
					if caseTestResultDataMsg[caseId].BaseData.MinRequestTime == 0 || caseTestResultDataMsg[caseId].BaseData.MinRequestTime > caseResult.MinRequestTime {
						caseTestResultDataMsg[caseId].BaseData.MinRequestTime = caseResult.MinRequestTime
					}

					caseTestResultDataMsg[caseId].CalcData.FiftyRequestTimeMap[resultDataMsg.MachineIP] = caseResult.FiftyRequestTimeLineValue
					caseTestResultDataMsg[caseId].CalcData.NinetyRequestTimeMap[resultDataMsg.MachineIP] = caseResult.NinetyRequestTimeLineValue
					caseTestResultDataMsg[caseId].CalcData.NinetyFiveRequestTimeMap[resultDataMsg.MachineIP] = caseResult.NinetyFiveRequestTimeLineValue
					caseTestResultDataMsg[caseId].CalcData.NinetyNineRequestTimeMap[resultDataMsg.MachineIP] = caseResult.NinetyNineRequestTimeLineValue

					caseTestResultDataMsg[caseId].BaseData.Normal1xxCodeNum += caseResult.StatusCodeCounter[0]
					caseTestResultDataMsg[caseId].BaseData.Normal2xxCodeNum += caseResult.StatusCodeCounter[1]
					caseTestResultDataMsg[caseId].BaseData.Normal3xxCodeNum += caseResult.StatusCodeCounter[2]
					caseTestResultDataMsg[caseId].BaseData.Normal4xxCodeNum += caseResult.StatusCodeCounter[3]
					caseTestResultDataMsg[caseId].BaseData.Normal5xxCodeNum += caseResult.StatusCodeCounter[4]

					if caseTestResultDataMsg[caseId].CalcData.ActualConcurrencyMap[resultDataMsg.MachineIP] < caseResult.ActualConcurrency {
						caseTestResultDataMsg[caseId].CalcData.ActualConcurrencyMap[resultDataMsg.MachineIP] = caseResult.ActualConcurrency
					}

					caseTestResultDataMsg[caseId].BaseData.RequestNum += caseResult.RequestNum
					caseTestResultDataMsg[caseId].BaseData.RequestTime += caseResult.RequestTime
					caseTestResultDataMsg[caseId].BaseData.SuccessNum += caseResult.SuccessNum
					caseTestResultDataMsg[caseId].BaseData.RequestErrorNum += caseResult.RequestErrorNum
					caseTestResultDataMsg[caseId].BaseData.AssertErrorNum += caseResult.AssertErrorNum
					caseTestResultDataMsg[caseId].BaseData.ErrorNum += caseResult.ErrorNum
					caseTestResultDataMsg[caseId].BaseData.SendBytes += caseResult.SendBytes
					caseTestResultDataMsg[caseId].BaseData.ReceivedBytes += caseResult.ReceivedBytes
				}
			}
			mu.Unlock()

			if resultDataMsg.End {
				delete(machineMp, resultDataMsg.MachineIP)
				model.ReleaseEngine(resultDataMsg.MachineIP)
				if len(machineMp) == 0 {
					_, err := model.RDB.SRem(context.Background(), UsedPartitions, partition).Result()
					model.ReportRdb.Del(context.Background(), fmt.Sprintf(ReportRunningLockKey, reportId))
					if err != nil {
						log.Logger.Error("usedPartitions删除失败：", err)
					}
					log.Logger.Infof("partitionID: %d被释放", partition)
					isEnd = true
					return
				}
			}
		}
	}()
	wg.Wait()
	log.Logger.Infof("计划:%d, 报告:%d 处理完成", planId, reportId)

}

func parseRequestLineTime(mp map[string]float64) float64 {
	machineCount := 0
	totalTime := float64(0)
	for _, t := range mp {
		machineCount += 1
		totalTime += t
	}
	avgTime, _ := decimal.NewFromFloat(totalTime / float64(machineCount)).Round(0).Float64()
	return avgTime
}
