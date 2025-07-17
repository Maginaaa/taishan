package model

import (
	"engine/config"
	"engine/internal/biz/log"
	"engine/middleware"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bytedance/sonic"
	"sort"
	"strings"
	"sync"
	"time"
)

type GoroutineCount struct {
	ID           int32
	GoroutineIDs map[int32]int
}

const (
	SamplingDataPartition = 10
	PrefixDataPartition   = 11 // 采样日志默认写死由partition11进行传递
)

var (
	kafkaProducer sarama.SyncProducer
)

// InitKafkaProducer 构建生产者
func InitKafkaProducer() {
	kafkaConf := sarama.NewConfig()
	kafkaConf.Producer.RequiredAcks = sarama.WaitForAll                                                 // 发送完数据需要leader和follow都确认
	kafkaConf.Producer.Partitioner = sarama.NewManualPartitioner                                        // 设置选择分区的策略为Hash,当设置key时，所有的key的消息都在一个分区Partitioner里
	kafkaConf.Producer.Return.Successes = true                                                          // 成功交付的消息将在success channel返回
	kafkaProducer, _ = sarama.NewSyncProducer(strings.Split(config.Conf.Kafka.Address, ","), kafkaConf) // 生产者客户端
	fmt.Println("kafkaProducer initialized")
}

// SendKafkaMsg 将需要的测试数据写入到kafka中
func SendKafkaMsg(dataMsgCh chan *ResultDataMsg, partitionId int32) {
	index, stageIndex := 0, 0

	planTestResultDataMsg := new(PlanTestResultDataMsg)
	caseRequestTimeListMap := make(map[int32]RequestTimeList)
	snapshotCh := make(chan struct {
		planTestResultDataMsg  *PlanTestResultDataMsg
		caseRequestTimeListMap map[int32]RequestTimeList
	}, 100)

	// requestTimeListMap的读写锁
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		timer := time.NewTicker(time.Duration(1) * time.Second)
		for {
			select {
			case <-timer.C:
				log.Logger.Infof("消费了%d条数据，剩余数据:%d条", index-stageIndex, len(dataMsgCh))
				stageIndex = index
				if planTestResultDataMsg.PlanID == 0 || planTestResultDataMsg.SceneResults == nil {
					break
				}
				if len(planTestResultDataMsg.SceneResults) == 0 {
					continue
				}
				mu.Lock()
				var planCopy *PlanTestResultDataMsg
				_ = deepCopy(planTestResultDataMsg, &planCopy)
				var caseMapCopy map[int32]RequestTimeList
				_ = deepCopy(caseRequestTimeListMap, &caseMapCopy)
				snapshotCh <- struct {
					planTestResultDataMsg  *PlanTestResultDataMsg
					caseRequestTimeListMap map[int32]RequestTimeList
				}{planCopy, caseMapCopy}
				if planTestResultDataMsg.End {
					log.Logger.Info(fmt.Sprintf("机器ip: %s，计划: %d，报告: %d, 测试数据向kafka写入完成！本次任务有： %d 条数据", middleware.LocalIp, planTestResultDataMsg.PlanID, planTestResultDataMsg.ReportID, index))
					timer.Stop()
					close(snapshotCh)
					return
				}
				for _, sceneResult := range planTestResultDataMsg.SceneResults {
					sceneResult.CaseResults = make(map[int32]*ApiTestResultDataMsg)
				}
				mu.Unlock()
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if msg, ok := <-dataMsgCh; ok {
				index++
				switch msg.MsgType {
				case TypePlanAction:
					if msg.Start {
						planTestResultDataMsg.ReportID = msg.GlobalInfo.ReportID
						planTestResultDataMsg.PlanID = msg.GlobalInfo.PlanID
						planTestResultDataMsg.MachineIP = msg.GlobalInfo.MachineIP
						planTestResultDataMsg.SceneResults = make(map[int32]*SceneResultDataMsg)
					}
					if msg.End {
						planTestResultDataMsg.End = true
						return
					}
					break
				case TypeSceneAction:
					if msg.Start {
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID] = new(SceneResultDataMsg)
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].SceneType = msg.SceneInfo.SceneType
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].TargetConcurrency = msg.SceneInfo.Concurrency
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults = make(map[int32]*ApiTestResultDataMsg)
					}
					if msg.End {
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].End = true
					}
					break
				case TypeSceneResultData:
					planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].TargetConcurrency = msg.SceneInfo.Concurrency
				case TypeHttpResultData:
					mu.Lock()
					if msg.GlobalInfo.NeedSampling {
						dataMs, _ := sonic.Marshal(msg)
						go sendMsg(SamplingDataPartition, dataMs)
					}
					if _, hasCase := planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID]; !hasCase {
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID] = &ApiTestResultDataMsg{
							CaseID:            msg.HttpResultDataMsg.CaseID,
							Url:               msg.HttpResultDataMsg.Url,
							RallyPoint:        msg.HttpResultDataMsg.RallyPoint,
							StatusCodeCounter: map[int]int64{},
						}
					}
					if msg.HttpResultDataMsg.StatusCode != 0 {
						firstDigit := msg.HttpResultDataMsg.StatusCode / 100
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].StatusCodeCounter[firstDigit-1]++
					}
					if msg.HttpResultDataMsg.IsSuccess {
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].SuccessNum++
					} else {
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].ErrorNum++
						if msg.HttpResultDataMsg.ErrorType == RequestErr {
							planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].RequestErrorNum++
						} else {
							planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].AssertErrorNum++
						}
					}
					planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].RequestNum++
					planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].RequestTime += msg.HttpResultDataMsg.RequestTime

					planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].SendBytes += msg.HttpResultDataMsg.SendBytes
					planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].ReceivedBytes += msg.HttpResultDataMsg.ReceivedBytes

					planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].RallyPoint = msg.HttpResultDataMsg.RallyPoint

					if planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].StartTime == 0 ||
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].StartTime > msg.HttpResultDataMsg.StartTime {
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].StartTime = msg.HttpResultDataMsg.StartTime
					}

					if planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].EndTime == 0 ||
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].EndTime < msg.HttpResultDataMsg.EndTime {
						planTestResultDataMsg.SceneResults[msg.SceneInfo.SceneID].CaseResults[msg.HttpResultDataMsg.CaseID].EndTime = msg.HttpResultDataMsg.EndTime
					}

					caseRequestTimeListMap[msg.HttpResultDataMsg.CaseID] = append(caseRequestTimeListMap[msg.HttpResultDataMsg.CaseID], msg.HttpResultDataMsg.RequestTime)
					mu.Unlock()
					break
				default:
					log.Logger.Error("未知的消息类型")
				}

			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for snapshot := range snapshotCh {
			for _, sceneResult := range snapshot.planTestResultDataMsg.SceneResults {
				for caseId, caseResult := range sceneResult.CaseResults {
					requestTimeList := snapshot.caseRequestTimeListMap[caseId]
					sort.Sort(requestTimeList)
					caseResult.MaxRequestTime = float64(requestTimeList[len(requestTimeList)-1])
					caseResult.MinRequestTime = float64(requestTimeList[0])
					caseResult.FiftyRequestTimeLineValue = float64(requestTimeList[len(requestTimeList)/2])
					caseResult.NinetyRequestTimeLineValue = TimeLineCalculate(90, requestTimeList)
					caseResult.NinetyFiveRequestTimeLineValue = TimeLineCalculate(95, requestTimeList)
					caseResult.NinetyNineRequestTimeLineValue = TimeLineCalculate(99, requestTimeList)
					caseResult.ActualConcurrency = sceneResult.TargetConcurrency
				}
			}
			go sendMsg(partitionId, snapshot.planTestResultDataMsg.ToByte())
		}
	}()

	wg.Wait()
	return
}

func sendMsg(partitionId int32, msg []byte) (err error) {
	proMsg := &sarama.ProducerMessage{}
	proMsg.Topic = config.Conf.Kafka.TopIc
	proMsg.Partition = partitionId
	proMsg.Value = sarama.ByteEncoder(msg)
	_, _, err = kafkaProducer.SendMessage(proMsg)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("机器ip:%s, 向kafka发送消息失败: %s", middleware.LocalIp, err.Error()))
	}
	return
}

func SendPrefixSceneMsg(msg []byte) (err error) {
	proMsg := &sarama.ProducerMessage{}
	proMsg.Topic = config.Conf.Kafka.TopIc
	proMsg.Partition = PrefixDataPartition
	proMsg.Value = sarama.ByteEncoder(msg)
	_, _, err = kafkaProducer.SendMessage(proMsg)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("机器ip:%s, 向kafka发送消息失败: %s", middleware.LocalIp, err.Error()))
	}
	return
}

// 2次/s的调用频率，可直接使用序列化/反序列化进行deepCopy
func deepCopy(src, dst interface{}) error {
	data, err := sonic.Marshal(src)
	if err != nil {
		return err
	}
	return sonic.Unmarshal(data, dst)
}
