package server

import (
	"engine/internal/biz/log"
	"engine/middleware"
	"engine/model"
	"engine/server/execution"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

func ExecutionPlan(_ *gin.Context, action model.Action) {

	// 设置接收数据缓存
	dataMsgCh := make(chan *model.ResultDataMsg, 500000)

	// kafka写入测试结果
	go model.SendKafkaMsg(dataMsgCh, action.PartitionID)

	// 分解任务
	TaskDecomposition(action, dataMsgCh)

}

// TaskDecomposition 分解任务
func TaskDecomposition(action model.Action, dataMsgCh chan *model.ResultDataMsg) {
	defer close(dataMsgCh)
	// 获取计划的参数化文件
	fileList := getDataFile(action)
	planVariableList := getPlanVariableList(action)
	var sceneDownWg = &sync.WaitGroup{}

	dataMsgCh <- &model.ResultDataMsg{
		MsgType: model.TypePlanAction,
		GlobalInfo: &model.GlobalInfo{
			PartitionID: action.PartitionID,
			ReportID:    action.ReportID,
			PlanID:      action.Plan.PlanID,
			Timestamp:   time.Now().UnixMilli(),
			MachineIP:   middleware.LocalIp,
		},
		Start: true,
		End:   false,
	}

	for _, scene := range action.ScenePressInfoList {
		scene.Scene.FileList = fileList
		scene.Scene.VariablePool.SaveList(planVariableList)
		scene.Scene.DefaultHeader = action.Plan.DefaultHeader
		scene.Scene.CasesMarshal()
		sc := model.SceneAction{
			ReportID:           action.ReportID,
			PlanID:             action.Plan.PlanID,
			SceneID:            scene.Scene.SceneID,
			ScenePressInfo:     scene,
			PartitionID:        action.PartitionID,
			PressInfo:          action.Plan.PressInfo,
			BreakType:          action.Plan.BreakType,
			BreakValue:         action.Plan.BreakValue,
			SamplingInfo:       action.Plan.SamplingInfo,
			EngineCount:        action.EngineCount,
			EngineSerialNumber: action.EngineSerialNumber,
		}
		if scene.Scene.SceneType == model.PreFixScene {
			if !ExecutionPrefixScene(sc, dataMsgCh) {
				break
			}
		} else {
			sceneDownWg.Add(1)
			go ExecutionScene(sceneDownWg, sc, dataMsgCh)
		}
	}
	sceneDownWg.Wait()

	dataMsgCh <- &model.ResultDataMsg{
		MsgType: model.TypePlanAction,
		GlobalInfo: &model.GlobalInfo{
			PartitionID: action.PartitionID,
			ReportID:    action.ReportID,
			PlanID:      action.Plan.PlanID,
			Timestamp:   time.Now().UnixMilli(),
			MachineIP:   middleware.LocalIp,
		},
		Start: false,
		End:   true,
	}
	log.Logger.Info("engine施压完成")
}

func ExecutionScene(sceneDownWg *sync.WaitGroup, action model.SceneAction, dataMsgCh chan *model.ResultDataMsg) {
	defer sceneDownWg.Done()

	var sceneExecuteWg = &sync.WaitGroup{}

	dataMsgCh <- &model.ResultDataMsg{
		MsgType: model.TypeSceneAction,
		SceneInfo: &model.SceneInfo{
			SceneID:   action.SceneID,
			SceneType: model.NormalScene,
		},
		Start: true,
		End:   false,
	}

	var msg string
	switch action.PressInfo.PressType {
	case model.ConcurrentModel, model.RpsModel:
		msg = execution.ConcurrentModel(sceneExecuteWg, action, dataMsgCh)
		break
	case model.StepModel:
		msg = execution.StepModel(sceneExecuteWg, action, dataMsgCh)
		break
	case model.RpsRateModel:
		msg = execution.RpsRateModel(sceneExecuteWg, action, dataMsgCh)
		break
	case model.FixedIterationModel:
		msg = execution.FixedIterationModel(sceneExecuteWg, action, dataMsgCh)
		break
	default:
		// TODO:发消息告诉主机执行失败
	}
	sceneExecuteWg.Wait()
	// end msg
	dataMsgCh <- &model.ResultDataMsg{
		MsgType: model.TypeSceneAction,
		SceneInfo: &model.SceneInfo{
			SceneID:   action.SceneID,
			SceneType: model.NormalScene,
		},
		Start: false,
		End:   true,
	}

	log.Logger.Info("单场景执行完成：", msg)
}

func ExecutionPrefixScene(action model.SceneAction, dataMsgCh chan *model.ResultDataMsg) bool {

	var preSceneExecuteWg = &sync.WaitGroup{}
	dataMsgCh <- &model.ResultDataMsg{
		MsgType: model.TypeSceneAction,
		SceneInfo: &model.SceneInfo{
			SceneID:   action.SceneID,
			SceneType: model.PreFixScene,
		},
		Start: true,
		End:   false,
	}
	resume := execution.DisposePrefixScene(preSceneExecuteWg, action, dataMsgCh)
	preSceneExecuteWg.Wait()

	dataMsgCh <- &model.ResultDataMsg{
		MsgType: model.TypeSceneAction,
		SceneInfo: &model.SceneInfo{
			SceneID:   action.SceneID,
			SceneType: model.PreFixScene,
		},
		Start: false,
		End:   true,
	}
	log.Logger.Infof("前置场景执行完成, 继续执行: %t", resume)

	return resume
}

func getDataFile(action model.Action) *model.FileList {
	// 解析参数化文件数据
	fileList := &model.FileList{
		File: action.FileInfo,
	}
	fileList.ParseFile(action.Plan.PlanID, action.EngineSerialNumber, action.EngineCount)
	return fileList
}

func getPlanVariableList(action model.Action) []*model.Variable {
	planVariableList := make([]*model.Variable, 0)
	for _, vrb := range action.Plan.GlobalVariable {
		if vrb.Enable {
			planVariableList = append(planVariableList, &model.Variable{
				VariableName: vrb.Key,
				VariableVal:  vrb.Value,
			})
		}
	}
	return planVariableList
}
