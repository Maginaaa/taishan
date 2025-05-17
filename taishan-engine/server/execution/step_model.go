package execution

import (
	"encoding/json"
	"engine/internal/biz/log"
	"engine/model"
	"engine/server/heartbeat"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

// StepModel 阶梯式加压模式
func StepModel(wg *sync.WaitGroup, sceneAction model.SceneAction, dataMsgCh chan *model.ResultDataMsg) string {

	startConcurrency := parseConcurrency(sceneAction.PressInfo.StartConcurrency, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
	stepSize := parseConcurrency(sceneAction.PressInfo.StepSize, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
	maxConcurrency := parseConcurrency(sceneAction.PressInfo.Concurrency, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
	stepDuration := sceneAction.PressInfo.StepDuration
	totalDuration := sceneAction.PressInfo.Duration * 60

	adjustKey := fmt.Sprintf("ReportStatusChange:%d:%d", sceneAction.ReportID, sceneAction.SceneID)
	pubSub := model.SubscribeMsg(adjustKey)
	statusCh := pubSub.ChannelSize(1000)
	defer func(pubSub *redis.PubSub) {
		err := pubSub.Close()
		if err != nil {
			log.Logger.Error("关闭redis订阅失败：", err)
		}
	}(pubSub)

	currentConcurrency := startConcurrency
	// target记录是否到达峰值，step记录步数
	reachPeak := false
	step := int64(0)
	targetTime, startTime, endTime := time.Now().Unix(), time.Now().Unix(), time.Now().Unix()
	sceneWg := &sync.WaitGroup{}
	// 达到机器性能瓶颈
	reachBottleneck := false
	notReachExecuteTime := true

	timer := time.NewTicker(time.Duration(1) * time.Second)
	defer timer.Stop()

	go func() {
		for {
			select {
			case <-timer.C:
				reachBottleneck = heartbeat.GetMemInfo().UsedPercent > 95
				notReachExecuteTime = startTime+stepDuration > endTime

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

	concurrentMap := new(sync.Map)
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
					return fmt.Sprintf("场景ID: %d, 当前并发数：%d, 总运行时长%ds, 任务手动结束！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency, time.Now().Unix()-targetTime)
				}
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
					for notReachExecuteTime {
						// 如果当前并发的id不在map中，那么就停止该goroutine
						if _, isOk := concurrentMap.Load(goroutineId); !isOk {
							break
						}
						if reachBottleneck {
							break
						}
						// 记录测试计划被完整的执行了一遍
						DisposeScene(wg, action, dataMsgCh, goroutineId, nil, 0)
					}
				}(i, sceneAction)
			}
			endTime = time.Now().Unix()
		}
		if currentConcurrency == maxConcurrency && startTime+stepDuration <= endTime {
			sceneWg.Wait()
			return fmt.Sprintf("最大并发数：%d， 总运行时长%ds, 任务正常结束！", currentConcurrency, endTime-targetTime)
		}
		// 如果当前并发数小于最大并发数，
		if currentConcurrency < maxConcurrency {
			if endTime-startTime >= stepDuration {
				// 从开始时间算起，加上持续时长。如果大于现在的时间，说明已经运行了持续时长的时间，那么就要将开始时间的值，变为现在的时间值
				currentConcurrency = currentConcurrency + stepSize
				step++
				if currentConcurrency > maxConcurrency {
					currentConcurrency = maxConcurrency
				}
				if currentConcurrency <= maxConcurrency {
					startTime = endTime
				}
			}
		}
		if currentConcurrency == maxConcurrency {
			if !reachPeak {
				reachPeak = true
				stepDuration = totalDuration - stepDuration*step
				startTime = endTime
			}

		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	sceneWg.Wait()
	return fmt.Sprintf("最大并发数：%d， 总运行时长%ds, 任务非正常结束！", currentConcurrency, time.Now().Unix()-targetTime)

}
