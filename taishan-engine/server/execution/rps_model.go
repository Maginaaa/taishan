package execution

import (
	"encoding/json"
	"engine/internal/biz/log"
	"engine/middleware"
	"engine/model"
	"engine/server/heartbeat"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	RpsValKey = "report:rps:%d"
)

func RpsRateModel(wg *sync.WaitGroup, sceneAction model.SceneAction, dataMsgCh chan *model.ResultDataMsg) string {

	// 初始并发数
	currentConcurrency := parseConcurrency(sceneAction.ScenePressInfo.Concurrency, sceneAction.EngineCount, sceneAction.EngineSerialNumber)

	// 初始rps
	currentRps := queryTotalRps(sceneAction)

	// 压测时长
	duration := sceneAction.PressInfo.Duration * 60

	// 监听过程中rps和并发数调整
	adjustKey := fmt.Sprintf("ReportStatusChange:%d:%d", sceneAction.ReportID, sceneAction.SceneID)
	pubSub := model.SubscribeMsg(adjustKey)
	statusCh := pubSub.Channel()
	defer func(pubSub *redis.PubSub) {
		err := pubSub.Close()
		if err != nil {
			log.Logger.Error("关闭redis订阅失败：", err)
		}
	}(pubSub)

	sceneWg := &sync.WaitGroup{}
	// 达到性能瓶颈
	reachBottleneck := false
	notReachExecuteTime := true

	// 并发模式根据时间进行压测
	log.Logger.Info(fmt.Sprintf("机器ip:%s, 开始性能测试,持续时间 %d秒", middleware.LocalIp, duration))
	targetTime, startTime := time.Now().Unix(), time.Now().Unix()

	timer := time.NewTicker(time.Duration(1) * time.Second)
	defer timer.Stop()

	// 记录每秒的执行次数
	var executeCount int64

	go func() {
		for {
			select {
			case <-timer.C:
				reachBottleneck = heartbeat.GetMemInfo().UsedPercent > 95
				notReachExecuteTime = startTime+duration >= time.Now().Unix()

				// 更新rps
				currentRps = queryTotalRps(sceneAction)

				// 发送当前并发
				dataMsgCh <- &model.ResultDataMsg{
					MsgType: model.TypeSceneResultData,
					SceneInfo: &model.SceneInfo{
						SceneID:     sceneAction.SceneID,
						Concurrency: currentConcurrency,
					},
				}
				//log.Logger.Infof("当前场景id: %d, 并发数: %d, 目标rps为: %d,当前执行次数为: %d", sceneAction.SceneID, currentConcurrency, currentRps, executeCount)
				atomic.StoreInt64(&executeCount, 0)
			}
		}
	}()

	concurrentMap := new(sync.Map)
	for notReachExecuteTime {
		select {
		case c := <-statusCh:
			var reportStatusChange = new(model.ReportStatusChange)
			_ = json.Unmarshal([]byte(c.Payload), reportStatusChange)
			if reportStatusChange == nil {
				continue
			}
			switch reportStatusChange.Type {
			case model.StopPlan:
				log.Logger.Info(fmt.Sprintf("ReportID: %d 报告，场景ID: %d，手动修改为： %s", sceneAction.ReportID, sceneAction.SceneID, c.Payload))
				if reportStatusChange.StopPlan == "stop" {
					concurrentMap.Range(func(key, value any) bool {
						concurrentMap.Delete(key)
						return true
					})
					sceneWg.Wait()
					return fmt.Sprintf("场景ID: %d, 并发数：%d, 总运行时长%ds, 任务手动结束！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency, time.Now().Unix()-targetTime)
				}
			case model.ReportChange:
				// 如果修改后的并发小于当前并发
				newConcurrency := parseConcurrency(reportStatusChange.ActionChangeInfo.Concurrency, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
				if newConcurrency < currentConcurrency {
					diff := currentConcurrency - newConcurrency
					// 将最后的几个并发从map中去掉
					for i := int64(0); i < diff; i++ {
						concurrentMap.Delete(currentConcurrency - 1 - i)
					}
				}
				currentConcurrency = newConcurrency
			}
		default:
			for i := int64(0); i < currentConcurrency; i++ {
				if _, ok := concurrentMap.Load(i); ok {
					continue
				}
				sceneWg.Add(1)
				concurrentMap.Store(i, true)
				go func(goroutineId int64, action model.SceneAction) {
					defer sceneWg.Done()
					defer concurrentMap.Delete(goroutineId)
					for executeCount < currentRps && notReachExecuteTime {
						// 如果当前并发的id不在map中，那么就停止该goroutine
						if _, isOk := concurrentMap.Load(goroutineId); !isOk {
							break
						}
						if reachBottleneck {
							break
						}
						// 记录测试计划被完整的执行了一遍
						DisposeScene(wg, sceneAction, dataMsgCh, goroutineId, &executeCount, currentRps)
					}
				}(i, sceneAction)
			}
		}
		// 这里设置为1，会导致timer.C清除executeCount后等待1s，从而5s能只能执行4s的情况
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	sceneWg.Wait()
	return fmt.Sprintf("场景ID: %d, 并发数：%d, 总运行时长%ds, 任务结束!", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency, time.Now().Unix()-targetTime)
}

func queryTotalRps(sceneAction model.SceneAction) int64 {
	targetTotalRpsStr, _ := model.ReportRdb.Get(fmt.Sprintf(RpsValKey, sceneAction.ReportID)).Result()
	targetTotalRps, _ := strconv.Atoi(targetTotalRpsStr)
	return parseConcurrency(int64(targetTotalRps)*sceneAction.ScenePressInfo.Rate/100, sceneAction.EngineCount, sceneAction.EngineSerialNumber)

}
