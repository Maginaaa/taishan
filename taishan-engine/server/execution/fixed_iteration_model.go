package execution

import (
	"encoding/json"
	"engine/internal/biz/log"
	"engine/model"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"sync/atomic"
	"time"
)

func FixedIterationModel(wg *sync.WaitGroup, sceneAction model.SceneAction, dataMsgCh chan *model.ResultDataMsg) string {

	currentConcurrency := parseConcurrency(sceneAction.ScenePressInfo.Concurrency, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
	//if !(currentConcurrency > 0) {
	//	return fmt.Sprintf("场景ID: %d, 并发数：%d, 压测无效！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency)
	//}
	if currentConcurrency > 2000 {
		fmt.Printf("场景ID: %d, 并发数：%d, 自动调整到2K！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency)
		currentConcurrency = 2000
	}
	totalExecuteCount := parseConcurrency(sceneAction.ScenePressInfo.Iteration, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
	if !(totalExecuteCount > 0) {
		return fmt.Sprintf("场景ID: %d, 执行次数：%d, 压测无效！", sceneAction.ScenePressInfo.Scene.SceneID, totalExecuteCount)
	}
	if currentConcurrency > totalExecuteCount {
		currentConcurrency = totalExecuteCount
	}

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
	// 定义一个map，管理并发
	concurrentMap := new(sync.Map)
	var executeCount int64

	timer := time.NewTicker(time.Duration(1) * time.Second)
	defer timer.Stop()

	go func() {
		for {
			select {
			case <-timer.C:
				dataMsgCh <- &model.ResultDataMsg{
					MsgType: model.TypeSceneResultData,
					SceneInfo: &model.SceneInfo{
						SceneID:     sceneAction.SceneID,
						SceneType:   model.PreFixScene,
						Concurrency: currentConcurrency,
					},
				}
			}
		}
	}()

	for executeCount < totalExecuteCount {
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
					return fmt.Sprintf("场景ID: %d, 执行次数：%d, 任务手动结束！", sceneAction.ScenePressInfo.Scene.SceneID, executeCount)
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
				go func(goroutineId int64) {
					defer sceneWg.Done()
					defer concurrentMap.Delete(goroutineId)
					for executeCount < totalExecuteCount {
						if _, isOk := concurrentMap.Load(goroutineId); !isOk {
							break
						}
						atomic.AddInt64(&executeCount, 1)
						DisposeScene(wg, sceneAction, dataMsgCh, goroutineId, nil, 0)
					}
				}(i)
			}
		}
	}
	sceneWg.Wait()

	return fmt.Sprintf("场景ID: %d, 并发数：%d, 总执行次数:%d, 任务结束！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency, executeCount)
}
