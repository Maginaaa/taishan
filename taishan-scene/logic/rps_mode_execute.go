package logic

import (
	"fmt"
	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"scene/rao"
	"strconv"
)

const (
	LimitConcurrency           = 2000 // 单节点单场景限制并发数
	LimitStepChangeConcurrency = int64(500)
	LimitStepChangeRate        = 0.05
)

func DoRpsModeOperation(ctx *gin.Context, reportId int32, action rao.Action) {
	reportData := GetReportData(ctx, reportId)
	if reportData == nil {
		return
	}
	sceneRateMp := make(map[int32]int64)
	for _, sce := range action.ScenePressInfoList {
		sceneRateMp[sce.Scene.SceneID] = sce.Rate
	}
	totalRpsStr, err := dal.ReportRdb.Get(ctx, fmt.Sprintf(RpsValKey, reportId)).Result()
	if err != nil {
		return
	}
	targetTotalRps, _ := strconv.Atoi(totalRpsStr)
	for _, scene := range reportData.Data.Scenes {
		if scene.SceneType != rao.NormalScene {
			continue
		}
		totalRT, actualTotalRps, caseNum := float64(0), float64(0), 0
		for _, cs := range scene.Cases {
			totalRT += cs.StageAvgRt
			caseNum += 1
			actualTotalRps += cs.StageRps
		}
		floatNum, _ := decimal.NewFromFloat(actualTotalRps / float64(targetTotalRps)).RoundFloor(3).Float64()
		if floatNum > 0.995 {
			continue
		}
		targetRps := int64(targetTotalRps) * sceneRateMp[scene.SceneID] / 100
		preTargetCon := int64(float64(targetRps)*totalRT)/(1000*int64(caseNum)) + 1
		if len(scene.Cases) > 1 {
			caseLen := len(scene.Cases)
			// 最后一个接口执行次数和第一个接口差距较大，说明并发数分配不合理
			if scene.Cases[caseLen-1].StageRequestNum+scene.Concurrency < scene.Cases[0].StageRequestNum {
				// 调整并发数 / 场景总执行 * 第一个接口的执行次数
				preTargetCon = preTargetCon * targetRps * 5 / (int64(caseLen) * scene.Cases[0].StageRequestNum)
			}
		}
		if preTargetCon == scene.Concurrency {
			continue
		}
		targetCon := int64(0)
		if preTargetCon < scene.Concurrency {
			targetCon = mathutil.Max(int64(action.EngineCount), preTargetCon)
		} else {
			targetCon = mathutil.Min(int64(action.EngineCount)*LimitConcurrency, preTargetCon)
		}
		log.Logger.Info("reportId: ", reportId, " sceneId: ", scene.SceneID, " 目标rps:", targetRps, " 需要调整到并发:", targetCon)
		_ = UpdatePress(ctx, rao.ConcurrencyChange{
			ReportID:    reportId,
			Concurrency: targetCon,
			SceneID:     scene.SceneID,
		})

	}

}
