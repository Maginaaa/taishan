package execution

import (
	"encoding/json"
	"engine/internal/biz/log"
	"engine/middleware"
	"engine/model"
	"engine/server/heartbeat"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

// ConcurrentModel 并发模式
func ConcurrentModel(wg *sync.WaitGroup, sceneAction model.SceneAction, dataMsgCh chan *model.ResultDataMsg) string {

	currentConcurrency := parseConcurrency(sceneAction.ScenePressInfo.Concurrency, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
	//if !(currentConcurrency > 0) {
	//	return fmt.Sprintf("场景ID: %d, 并发数：%d, 压测无效！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency)
	//}
	if currentConcurrency > 2000 {
		fmt.Printf("场景ID: %d, 并发数：%d, 自动调整到2K！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency)
		currentConcurrency = 2000
	}

	currentPartitionID := sceneAction.PartitionID
	duration := sceneAction.PressInfo.Duration * 60
	// 订阅redis中消息  任务状态：包括：报告停止；debug日志状态；任务配置变更
	adjustKey := fmt.Sprintf("ReportStatusChange:%d:%d", sceneAction.ReportID, sceneAction.SceneID)
	pubSub := model.SubscribeMsg(adjustKey)
	statusCh := pubSub.ChannelSize(1000)
	defer func(pubSub *redis.PubSub) {
		err := pubSub.Close()
		if err != nil {
			log.Logger.Error("关闭redis订阅失败：", err)
		}
	}(pubSub)

	sceneWg := &sync.WaitGroup{}
	// 定义一个map，管理并发
	concurrentMap := new(sync.Map)
	// 达到压力机性能瓶颈
	reachBottleneck := false
	notReachExecuteTime := true

	// 并发模式根据时间进行压测
	log.Logger.Info(fmt.Sprintf("机器ip:%s, 开始性能测试,持续时间 %d秒", middleware.LocalIp, duration))
	targetTime, startTime := time.Now().Unix(), time.Now().Unix()

	timer := time.NewTicker(time.Duration(1) * time.Second)
	defer timer.Stop()

	go func() {
		for {
			select {
			case <-timer.C:
				reachBottleneck = heartbeat.GetMemInfo().UsedPercent > 95
				notReachExecuteTime = startTime+duration >= time.Now().Unix()

				// 发送当前并发
				dataMsgCh <- &model.ResultDataMsg{
					MsgType: model.TypeSceneResultData,
					SceneInfo: &model.SceneInfo{
						SceneID:     sceneAction.SceneID,
						Concurrency: currentConcurrency,
					},
				}
			}
		}
	}()
	for notReachExecuteTime {
		select {
		case c := <-statusCh:
			var reportStatusChange = new(model.ReportStatusChange)
			_ = json.Unmarshal([]byte(c.Payload), reportStatusChange)
			if reportStatusChange == nil {
				continue
			}
			log.Logger.Info(fmt.Sprintf("ReportID: %d 报告，场景ID: %d，手动修改为： %s", sceneAction.ReportID, sceneAction.SceneID, c.Payload))
			switch reportStatusChange.Type {
			case model.StopPlan:
				if reportStatusChange.StopPlan == "stop" {
					concurrentMap.Range(func(key, value any) bool {
						concurrentMap.Delete(key)
						return true
					})
					sceneWg.Wait()
					return fmt.Sprintf("场景ID: %d, 并发数：%d, 总运行时长%ds, 任务手动结束！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency, time.Now().Unix()-targetTime)
				}
			//case constant.DebugStatus:
			//	debug = subscriptionStressPlanStatusChange.Debug
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
				go func(goroutineId int64, sceneAction model.SceneAction) {
					defer sceneWg.Done()
					defer concurrentMap.Delete(goroutineId)
					for notReachExecuteTime {
						// 如果当前并发的id不在map中，那么就停止该goroutine
						if _, isOk := concurrentMap.Load(goroutineId); !isOk {
							break
						}
						if reachBottleneck {
							break
						}
						// 记录测试计划被完整的执行了一遍
						DisposeScene(wg, sceneAction, dataMsgCh, goroutineId, nil, 0)
					}
				}(i, sceneAction)
			}
		}
		// TODO: 将当前逻辑与timer.C进行聚合
		time.Sleep(time.Duration(1) * time.Second)
	}
	sceneWg.Wait()
	return fmt.Sprintf("场景ID: %d, 并发数：%d, 总运行时长%ds, partitionId: %d, 任务结束!", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency, time.Now().Unix()-targetTime, currentPartitionID)

}

// parseConcurrency 通过总并发和机器总数、机器序号，计算真实并发
func parseConcurrency(totalCon int64, engineCount, engineSerialNumber int32) (concurrency int64) {
	concurrency = totalCon / int64(engineCount)
	if int64(engineSerialNumber) < totalCon%int64(engineCount) {
		concurrency += 1
	}
	return concurrency
}
