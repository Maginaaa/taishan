package logic

import (
	"report/rao"
)

const (
	LimitConcurrency           = 2000 // 单节点限制并发数
	LimitStepChangeConcurrency = int64(200)
	LimitStepChangeRate        = 0.05
)

func getSceneProcessDetail(scenes []rao.ReportStatusSceneEntity, sceneId int32) rao.ReportStatusSceneEntity {
	for _, sceneInfo := range scenes {
		if sceneInfo.SceneId == sceneId {
			return sceneInfo
		}
	}
	return rao.ReportStatusSceneEntity{}
}

//func DoRpsModeOperation(ctx *gin.Context, reportId int32) {
//	reportDetail, _ := getReport(ctx, reportId)
//	//todo 改成结构体返回
//	var sceneResultList []rao.SceneInformation
//	_ = json.Unmarshal([]byte(reportDetail.PressResult), &sceneResultList)
//	processReport, _ := GetReportRps(ctx, reportId)
//	if processReport != nil {
//		for _, scene := range sceneResultList {
//			//获取缓存
//			targetCon := int64(0)
//			floatNum := float64(0)
//			sceneTotalRps := float64(0)
//			actualCon := int64(0)
//			preCon := int64(0)
//			changeCon := int64(0)
//			processScene := getSceneProcessDetail(processReport.Scenes, scene.SceneId)
//			for _, caseInfo := range processScene.Cases {
//				sceneTotalRps += caseInfo.CurrentRps
//			}
//			floatNum = sceneTotalRps / float64(scene.SceneRps)
//			log.Logger.Info("场景目标rps：", scene.SceneRps, "场景实际rps：", sceneTotalRps)
//			if scene.Cases != nil {
//				actualCon = processScene.Cases[0].ActualConcurrency
//				preCon = int64(float64(scene.SceneRps*actualCon) / sceneTotalRps)
//
//				changeCon = mathutil.Max(mathutil.Abs(preCon-actualCon)/3, LimitStepChangeConcurrency)
//			}
//			// 计算当前rps
//			if preCon < actualCon {
//				targetCon = int64(mathutil.Max(1, int(mathutil.Max(preCon, actualCon-changeCon))))
//			} else {
//				targetCon = mathutil.Min(int64(len(reportDetail.EngineList))*LimitConcurrency, mathutil.Min(preCon, actualCon+LimitStepChangeConcurrency))
//			}
//			log.Logger.Info("reportId: ", reportId, " sceneId: ", scene.SceneId, "需要调整到并发:", targetCon)
//			absNum := mathutil.Abs(floatNum - 1)
//			if absNum > LimitStepChangeRate && targetCon > 0 {
//				//todo 计算出的并发数单节点大于2000，触发扩容操作
//				conCurrency := rao.ConcurrencyChange{
//					ReportID:    reportId,
//					Concurrency: targetCon,
//					SceneID:     scene.SceneId,
//				}
//				err := UpdatePressCurrency(ctx, conCurrency)
//				if err != nil {
//					return
//				}
//			}
//		}
//	} else {
//		log.Logger.Info("报告还未生成")
//	}
//
//}
