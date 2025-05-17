package execution

import (
	"encoding/json"
	"engine/internal/biz/log"
	"engine/middleware"
	"engine/model"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"sync/atomic"
	"time"
)

// DisposePrefixScene 执行前置场景
// 返回值表示是否需要继续执行普通场景
func DisposePrefixScene(wg *sync.WaitGroup, sceneAction model.SceneAction, dataMsgCh chan *model.ResultDataMsg) bool {

	exportDataInfo := sceneAction.ScenePressInfo.Scene.ExportDataInfo
	if exportDataInfo.HasCache && !exportDataInfo.DisableCache {
		sceneAction.ScenePressInfo.Scene.FileList.ParsePreSceneFile(sceneAction.SceneID, exportDataInfo, sceneAction.EngineSerialNumber, sceneAction.EngineCount)
		log.Logger.Infof("场景： %d直接使用缓存，跳过执行", sceneAction.SceneID)
		return true
	}

	totalExecuteCount := parseConcurrency(exportDataInfo.ExportTimes, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
	currentConcurrency := parseConcurrency(exportDataInfo.Concurrency, sceneAction.EngineCount, sceneAction.EngineSerialNumber)
	if !(currentConcurrency > 0) {
		log.Logger.Infof("场景ID: %d, 并发数：%d, 压测无效！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency)
		return false
	}
	if currentConcurrency > 2000 {
		log.Logger.Infof("场景ID: %d, 并发数：%d, 自动调整到2K！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency)
		currentConcurrency = 2000
	}

	variableCh := make(chan map[string]string, currentConcurrency)
	defer close(variableCh)
	var executeCount int64

	// 订阅redis中消息  任务状态：包括：报告停止；debug日志状态；任务配置变更
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
	variableWg := &sync.WaitGroup{}
	// 定义一个map，管理并发
	concurrentMap := new(sync.Map)

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

	param := make(map[string]*model.Param)
	for _, key := range exportDataInfo.VariableList {
		param[key] = &model.Param{
			ReadType: model.OrderedRead,
			Val:      make([]string, 0),
		}
	}

	go func(param map[string]*model.Param) {
		tf := model.TransferInfo[model.ExportDataMap]{
			Type:      model.TypePreSceneExport,
			MachineIP: middleware.LocalIp,
			ReportID:  sceneAction.ReportID,
			SceneID:   sceneAction.SceneID,
			End:       false,
		}
		for variableMp := range variableCh {
			variableWg.Done()
			exportDataMp := make(map[string]string)
			for _, key := range exportDataInfo.VariableList {
				if _, ok := param[key]; ok {
					param[key].Val = append(param[key].Val, variableMp[key])
					exportDataMp[key] = variableMp[key]
				}
			}
			tf.Data = exportDataMp
			model.SendPrefixSceneMsg(tf.ToByte())
		}
		tf.End = true
		model.SendPrefixSceneMsg(tf.ToByte())
	}(param)

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
					sceneAction.ScenePressInfo.Scene.FileList.File = append(sceneAction.ScenePressInfo.Scene.FileList.File, &model.FileInfo{
						DataType: model.ExportDataType,
						DataMap:  param,
					})
					log.Logger.Infof("前置场景ID: %d, 并发数：%d, 任务手动结束！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency)
					return false
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
			case model.SceneRelease:
				concurrentMap.Range(func(key, value any) bool {
					concurrentMap.Delete(key)
					return true
				})
				sceneWg.Wait()
				variableWg.Wait()
				sceneAction.ScenePressInfo.Scene.FileList.File = append(sceneAction.ScenePressInfo.Scene.FileList.File, &model.FileInfo{
					DataType: model.ExportDataType,
					DataMap:  param,
				})
				log.Logger.Infof("前置场景ID: %d, 并发数：%d, 任务释放！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency)
				return true
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
						disposePrefixScene(wg, sceneAction, dataMsgCh, variableCh, goroutineId, variableWg)
					}
				}(i)
			}
		}
	}
	sceneWg.Wait()
	variableWg.Wait()
	sceneAction.ScenePressInfo.Scene.FileList.File = append(sceneAction.ScenePressInfo.Scene.FileList.File, &model.FileInfo{
		DataType: model.ExportDataType,
		DataMap:  param,
	})

	log.Logger.Infof("前置场景ID: %d, 并发数：%d, 总执行次数:%d, 任务结束！", sceneAction.ScenePressInfo.Scene.SceneID, currentConcurrency, executeCount)
	return true

}

func disposePrefixScene(wg *sync.WaitGroup, sceneAction model.SceneAction, dataMsgCh chan *model.ResultDataMsg, variableCh chan map[string]string, goroutineId int64, variableWg *sync.WaitGroup) {

	allPass := true

	s := sceneAction.ScenePressInfo.Scene.DeepCopy()
	s.LoadFileVariable()
	s.LoadSceneVariablePool()

	for _, cs := range s.MarshalCases {
		switch cs.Type {
		case model.HttpCaseType:
			caseData := cs.HttpCase.DeepCopy()
			caseData.LoadVariablePool(s.VariablePool)
			caseData.LoadDefaultHeader(s.DefaultHeader)
			wg.Add(1)
			err := DisposeRequest(wg, &caseData, sceneAction, dataMsgCh, goroutineId)
			if err != nil {
				continue
			}
			s.StoreVariablePool(&caseData)
			if !caseData.ResponseData.RequestSuccess || !caseData.ResponseData.AssertSuccess {
				allPass = false
			}
		default:
			continue
		}
	}

	if allPass {
		mp := make(map[string]string)
		s.VariablePool.VariableMap.Range(func(key, value any) bool {
			mp[key.(string)] = value.(string)
			return true
		})
		variableWg.Add(1)
		variableCh <- mp
	}
}
