package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/mathutil"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gen"

	"report/conf"
	"report/internal/dal"
	"report/internal/libs"
	"report/internal/model"
	"report/log"
	"report/rao"
)

const (
	ReportRunningLockKey = "report:lock:%d"
	ReportDataCacheKey   = "report:data:%d"
	RpsValKey            = "report:rps:%d"

	GraphMaxLen = 120
)

func CreateReport(ctx *gin.Context, req rao.CreateReportReq) (reportId int32, err error) {
	tx := dal.GetQuery().Report

	insertData := &model.Report{
		PlanID:       req.PlanID,
		ReportName:   fmt.Sprintf("%s_%s", req.PlanName, time.Now().Format("20060102150405")),
		Status:       true,
		Duration:     req.Duration,
		PressType:    req.PressType,
		EngineList:   req.EngineList,
		CreateUserID: req.CreateUserID,
	}
	err = tx.WithContext(ctx).Omit(tx.EndTime).Create(insertData)
	if err != nil {
		log.Logger.Error("logic.report.CreateReport.Create ，err:", err)
		return 0, err
	}
	reportId = insertData.ID
	ms, _ := json.Marshal(req.EngineList)
	_, err = dal.ReportRdb.Set(ctx, fmt.Sprintf(ReportRunningLockKey, reportId), string(ms), 0).Result()
	if err != nil {
		log.Logger.Error("logic.report.CreateReport.SetRedis, err:", err)
		return
	}
	return
}

// ReportPressDown 压测被动结束，清除压测报告和测试计划状态
func ReportPressDown(ctx *gin.Context, reportId int32) (err error) {
	reportData, err := GetReportData(ctx, reportId)
	if err != nil {
		log.Logger.Error("logic.report.ReportPressDown err")
		return
	}
	totalReqTime := reportData.TotalRequestTime
	pressStartTime := reportData.TotalStartTime
	pressEndTime := reportData.TotalEndTime

	// 修改mysql
	tx := dal.GetQuery().Report
	report, err := tx.WithContext(ctx).Where(tx.ID.Eq(reportId)).First()
	machineNum := int64(len(report.EngineList))
	// vum = (((总请求时间 / (压测时长 * 机器数)) / 500) + 1) * 500 * 机器数 * 压测时长(分)
	vum := int64(0)
	if pressEndTime != 0 && pressStartTime != 0 && pressEndTime-pressStartTime != 0 {
		vum = ((totalReqTime/((pressEndTime-pressStartTime)*machineNum))/500 + 1) * 500 * machineNum * (pressEndTime - pressStartTime) / 1000 / 60
	}

	// 运行中的，先改mysql，再改redis内的plan锁
	_, err = tx.WithContext(ctx).Select(tx.Status, tx.Vum, tx.EndTime).Where(tx.ID.Eq(reportId)).Updates(model.Report{
		Status:  false,
		Vum:     int32(vum),
		EndTime: time.UnixMilli(pressEndTime),
	})
	if err != nil {
		log.Logger.Error("logic.report.reportStopPress.UpdateSimple ，err:", err)
		return
	}
	dal.ReportRdb.Del(ctx, fmt.Sprintf(ReportRunningLockKey, reportId))
	err = PlanRunningDelLock(ctx, report.PlanID)
	if err != nil {
		log.Logger.Error("logic.report.reportStopPress.PlanRunningDelLock ，err:", err)
		return
	}
	log.Logger.Infof("测试报告: %d生成测试报告完成...", reportId)
	return
}

func GetReportList(ctx *gin.Context, req rao.ReportListReq) (res []*rao.Report, count int64, err error) {
	tx := dal.GetQuery().Report
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.IsDelete.Not())
	if req.PlanID != 0 {
		conditions = append(conditions, tx.PlanID.Eq(req.PlanID))
	}
	if req.CreateUserId != 0 {
		conditions = append(conditions, tx.CreateUserID.Eq(req.CreateUserId))
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	if req.StartTime != "" {
		startTime, _ := time.ParseInLocation(FullTimeFormat, req.StartTime, loc)
		conditions = append(conditions, tx.CreateTime.Gte(startTime))
	}
	if req.EndTime != "" {
		endTime, _ := time.ParseInLocation(FullTimeFormat, req.EndTime, loc)
		conditions = append(conditions, tx.CreateTime.Lte(endTime))
	}
	reportList, count, err := tx.WithContext(ctx).Where(conditions...).Order(tx.ID.Desc()).FindByPage(offset, limit)
	if err != nil {
		log.Logger.Error("logic.report.GetReportList.FindByPage ，err:", err)
		return nil, 0, err
	}

	// 获取用户信息列表
	var users []rao.User
	users, err = GetUserList()
	if err != nil {
		return nil, 0, errors.New("获取用户信息失败")
	}

	planIds := make([]int32, 0)
	for _, s := range reportList {
		planIds = append(planIds, s.PlanID)
	}
	// TODO:目前为实时取测试计划名，应改为快照
	planMap, err := BatchGetPlanName(ctx, planIds)
	if err != nil {
		log.Logger.Error("logic.report.GetReportList.BatchGetPlanName ，err:", err)
		return
	}
	res = make([]*rao.Report, 0, len(reportList))
	for _, s := range reportList {
		actualDuration, _ := decimal.NewFromFloat(s.EndTime.Sub(s.CreateTime).Minutes()).Round(2).Float64()
		res = append(res, &rao.Report{
			ReportID:       s.ID,
			ReportName:     s.ReportName,
			Status:         s.Status,
			PlanID:         s.PlanID,
			PlanName:       planMap[s.PlanID],
			Duration:       s.Duration,
			ActualDuration: actualDuration,
			PressType:      s.PressType,
			StartTime:      s.CreateTime.Format(FullTimeFormat),
			CreateUserName: GetNameByID(users, s.CreateUserID),
			UpdateUserName: GetNameByID(users, s.UpdateUserID),
			UpdateTime:     s.UpdateTime.Format(FullTimeFormat),
			EndTime:        s.EndTime.Format(FullTimeFormat),
		})
	}
	return res, count, nil
}

// StopPress 压测被手动关闭
func StopPress(ctx *gin.Context, reportId int32) (err error) {
	msg := rao.ReportStatusChange{
		Type:     rao.StopPlan,
		StopPlan: "stop",
	}
	msgByte, err := json.Marshal(msg)
	// 根据reportId获取场景列表
	report, err := GetReportDetail(ctx, reportId)
	if err != nil {
		log.Logger.Error("logic.report.StopPress.getReportById ，err:", err)
		return
	}
	_, lockInfo, err := getPlanRunningLock(ctx, report.PlanID)
	if err != nil {
		log.Logger.Error("logic.report.StopPress.getPlanRunningLock ，err:", err)
		return
	}
	for _, sceneID := range lockInfo.SceneIDs {
		_, err = dal.RDB.Publish(ctx, fmt.Sprintf("ReportStatusChange:%d:%d", reportId, sceneID), string(msgByte)).Result()
		if err != nil {
			log.Logger.Error("logic.report.StopPress.SetRedis ，err:", err)
			return err
		}
	}

	return nil
}

// UpdatePressCurrency 压测被手动修改
func UpdatePressCurrency(ctx *gin.Context, req rao.ConcurrencyChange) (err error) {
	msg := rao.ReportStatusChange{
		Type: rao.ReportChange,
		ActionChangeInfo: rao.ActionChangeInfo{
			Concurrency: req.Concurrency,
		},
	}
	msgByte, err := json.Marshal(msg)
	_, err = dal.RDB.Publish(ctx, fmt.Sprintf("ReportStatusChange:%d:%d", req.ReportID, req.SceneID), string(msgByte)).Result()
	if err != nil {
		log.Logger.Error("logic.report.UpdatePress.SetRedis ，err:", err)
		return err
	}
	return
}

func ReleasePreScene(ctx *gin.Context, req rao.ConcurrencyChange) (err error) {
	msg := rao.ReportStatusChange{
		Type: rao.SceneRelease,
	}
	msgByte, err := json.Marshal(msg)
	_, err = dal.RDB.Publish(ctx, fmt.Sprintf("ReportStatusChange:%d:%d", req.ReportID, req.SceneID), string(msgByte)).Result()
	if err != nil {
		log.Logger.Error("logic.report.UpdatePress.SetRedis ，err:", err)
		return err
	}
	return
}

func UpdateReportName(ctx *gin.Context, req rao.UpdateReportNameReq) (err error) {
	tx := dal.GetQuery().Report
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(req.ReportID)).Update(tx.ReportName, req.ReportName)
	if err != nil {
		log.Logger.Error("logic.report.UpdateReportName.Update ，err:", err)
		return err
	}
	return
}

// GetReportDetail 获取Report的基础信息
func GetReportDetail(ctx *gin.Context, reportId int32) (report rao.Report, err error) {
	reportInfo, err := getReport(ctx, reportId)
	if err != nil {
		log.Logger.Error("logic.report.GetReportDetail.First ，err:", err)
		return report, errors.New("查询报告信息失败")
	}

	users, err := GetUserList()
	if err != nil {
		log.Logger.Error("logic.report.GetReportDetail.GetUserList ，err:", err)
		return report, errors.New("用户信息获取失败")
	}

	planName, err := GetPlanNameByID(ctx, reportInfo.PlanID)
	if err != nil {
		log.Logger.Error("logic.report.GetReportDetail.GetPlanNameByID, err: ", err)
		return report, errors.New("获取报告对应计划名称信失败")
	}

	return rao.Report{
		ReportID:       reportInfo.ID,
		ReportName:     reportInfo.ReportName,
		Status:         reportInfo.Status,
		PlanID:         reportInfo.PlanID,
		PlanName:       planName,
		PressType:      reportInfo.PressType,
		Concurrency:    int64(reportInfo.Concurrency),
		EngineList:     reportInfo.EngineList,
		Duration:       reportInfo.Duration,
		StartTime:      reportInfo.CreateTime.Format(FullTimeFormat),
		EndTime:        reportInfo.EndTime.Format(FullTimeFormat),
		CreateUserID:   reportInfo.CreateUserID,
		CreateUserName: GetNameByID(users, reportInfo.CreateUserID),
		UpdateUserName: GetNameByID(users, reportInfo.UpdateUserID),
		UpdateTime:     reportInfo.UpdateTime.Format(FullTimeFormat),
	}, nil

}

func getReport(ctx *gin.Context, reportId int32) (report *model.Report, err error) {
	tx := dal.GetQuery().Report
	report, err = tx.WithContext(ctx).Where(tx.ID.Eq(reportId)).First()
	if err != nil {
		log.Logger.Error("logic.report.getReport.First ，err:", err)
	}
	return
}

// GetReportData 获取测试报告数据
func GetReportData(ctx *gin.Context, reportId int32) (planResultData rao.PlanResultData, err error) {
	baseData, err := GetReportBaseData(ctx, reportId)
	if err != nil {
		log.Logger.Errorf("logic.report.GetReportData.GetReportBaseData err: %v", err)
		return
	}
	b, _ := dal.ReportRdb.Exists(ctx, fmt.Sprintf(ReportRunningLockKey, reportId)).Result()
	planResultData.End = b == 0
	stageCount := len(baseData)
	timestampKeys := make([]int64, 0, stageCount)
	for k := range baseData {
		timestampKeys = append(timestampKeys, k)
	}
	sort.Slice(timestampKeys, func(i, j int) bool {
		return timestampKeys[i] < timestampKeys[j] // 升序排序,如果需要降序排序，则改为 return slice[i] > slice[j]
	})
	planResultData.Graph = &rao.ReportDataGraphEntity{
		Rps:         make([][3]any, 0, stageCount),
		RequestTime: make([][2]any, 0, stageCount),
		SuccessRate: make([][2]any, 0, stageCount),
	}
	planResultData.Scenes = make([]*rao.SceneResultData, 0)
	planResultData.ErrorUrlArray = make(rao.ApiDistributionList, 0)
	planResultData.ErrorCodeArray = make([]rao.ApiDistribution, 0)
	planResultData.RequestErrorArray = make([]rao.ApiDistribution, 0)
	planResultData.AssertErrorArray = make([]rao.ApiDistribution, 0)

	idMap := make(map[int32]map[int32]struct{})
	errorUrlMp := make(map[string]*rao.ApiDistribution)

	for stage, tms := range timestampKeys {
		timeStr := time.Unix(tms, 0).Format(HourMinSec)
		var (
			stepStartTime        int64 = 0
			stepEndTime          int64 = 0
			stepTotalConcurrency int64 = 0 // 当前时间戳总并发数 = 每个Scene的并发数之和
			stepTotalRequestNum  int64 = 0
			stepTotalRequestTime int64 = 0
			stepTotalSuccessNum  int64 = 0
			sceneConcurrency           = make(map[int32]int64)
		)
		for _, cs := range baseData[tms] {
			// scene、case初始化
			if _, ok := idMap[cs.SceneId]; !ok {
				idMap[cs.SceneId] = make(map[int32]struct{})
				planResultData.Scenes = append(planResultData.Scenes, &rao.SceneResultData{
					SceneID:   cs.SceneId,
					SceneType: cs.SceneType,
					Cases:     make([]*rao.CaseResultData, 0),
				})
			}
			if _, hasCase := idMap[cs.SceneId][cs.CaseId]; !hasCase {
				idMap[cs.SceneId][cs.CaseId] = struct{}{}
				for i, sInfo := range planResultData.Scenes {
					if sInfo.SceneID == cs.SceneId {
						planResultData.Scenes[i].Cases = append(planResultData.Scenes[i].Cases, &rao.CaseResultData{
							CaseID: cs.CaseId,
						})
						break
					}
				}
			}
			if _, ok := errorUrlMp[cs.Url]; !ok && cs.Url != "" {
				errorUrlMp[cs.Url] = &rao.ApiDistribution{
					Url: cs.Url,
				}
			}

			// total数据填充
			for _, s := range planResultData.Scenes {
				if s.SceneID != cs.SceneId {
					continue
				}
				for _, c := range s.Cases {
					if c.CaseID != cs.CaseId {
						continue
					}
					if cs.Url != "" {
						errorUrlMp[cs.Url].ErrorCount += cs.RequestErrorNum + cs.AssertErrorNum
						errorUrlMp[cs.Url].Count += cs.RequestNum
					}
					c.TotalRequestTime += cs.RequestTime
					c.TotalRequestNum += cs.RequestNum
					c.TotalSuccessNum += cs.SuccessNum
					c.TotalRequestErrorNum += cs.RequestErrorNum
					c.TotalAssertErrorNum += cs.AssertErrorNum
					c.TotalErrorNum += cs.ErrorNum
					c.TotalSuccessRate = libs.CalcRate(c.TotalSuccessNum, c.TotalRequestNum)
					c.TwoXxCodeNum += cs.Normal2xxCodeNum
					c.ThreeXxCodeNum += cs.Normal3xxCodeNum
					c.FourXxCodeNum += cs.Normal4xxCodeNum
					c.FiveXxCodeNum += cs.Normal5xxCodeNum
					c.OtherCodeNum += cs.Normal1xxCodeNum + cs.Normal3xxCodeNum + cs.Normal4xxCodeNum + cs.Normal5xxCodeNum
					c.SendBytes += cs.SendBytes
					c.ReceivedBytes += cs.ReceivedBytes
					if c.TotalStartTime == 0 || c.TotalStartTime > cs.StartTime {
						c.TotalStartTime = cs.StartTime
					}
					if c.TotalEndTime == 0 || c.TotalEndTime < cs.EndTime {
						c.TotalEndTime = cs.EndTime
					}
					c.TotalRps = libs.CalcRpsNew(c.TotalRequestNum, c.TotalEndTime-c.TotalStartTime)
					c.TotalAvgRt = libs.CalcDiv(c.TotalRequestTime, c.TotalRequestNum)

					planResultData.TotalRequestTime += cs.RequestTime
					planResultData.TotalRequestNum += cs.RequestNum
					planResultData.TotalSuccessNum += cs.SuccessNum
					planResultData.TotalErrorNum += cs.ErrorNum
					planResultData.TotalSendBytes += cs.SendBytes
					planResultData.TotalReceivedBytes += cs.ReceivedBytes
					if planResultData.TotalStartTime == 0 || planResultData.TotalStartTime > cs.StartTime {
						planResultData.TotalStartTime = cs.StartTime
					}
					if planResultData.TotalEndTime == 0 || planResultData.TotalEndTime < cs.EndTime {
						planResultData.TotalEndTime = cs.EndTime
					}

					c.FiftyRequestTimeLineValue = cs.FiftyRequestTimeLineValue
					c.NinetyRequestTimeLineValue = cs.NinetyRequestTimeLineValue
					c.NinetyFiveRequestTimeLineValue = cs.NinetyFiveRequestTimeLineValue
					c.NinetyNineRequestTimeLineValue = cs.NinetyNineRequestTimeLineValue

					// stage数据填充
					if stage == len(timestampKeys)-1 {
						c.StageRequestTime = cs.RequestTime
						c.StageRequestNum = cs.RequestNum
						c.StageSuccessNum = cs.SuccessNum
						c.StageErrorNum = cs.ErrorNum

						if c.StageStartTime == 0 || c.StageStartTime > cs.StartTime {
							c.StageStartTime = cs.StartTime
						}
						if c.StageEndTime == 0 || c.StageEndTime < cs.EndTime {
							c.StageEndTime = cs.EndTime
						}
						c.StageRps = libs.CalcRpsNew(c.StageRequestNum, calcDuration(c.StageStartTime, c.StageEndTime))
						c.StageSuccessRate = libs.CalcRate(c.StageSuccessNum, c.StageRequestNum)
						c.StageAvgRt = libs.CalcDiv(c.StageRequestTime, c.StageRequestNum)

						//c.MaxRt = mathutil.Max(c.MaxRt, cs.MaxRequestTime)
						//if c.MinRt != 0 {
						//	c.MinRt = mathutil.Min(c.MinRt, cs.MinRequestTime)
						//} else {
						//	c.MinRt = cs.MinRequestTime
						//}

						s.Concurrency = cs.ActualConcurrency

						planResultData.StageRequestTime += cs.RequestTime
						planResultData.StageRequestNum += cs.RequestNum
						planResultData.StageSuccessNum += cs.SuccessNum
						planResultData.StageErrorNum += cs.ErrorNum
						if planResultData.StageStartTime == 0 || planResultData.StageStartTime > cs.StartTime {
							planResultData.StageStartTime = cs.StartTime
						}
						if planResultData.StageEndTime == 0 || planResultData.StageEndTime < cs.EndTime {
							planResultData.StageEndTime = cs.EndTime
						}
					}

				}
			}
			sceneConcurrency[cs.SceneId] = cs.ActualConcurrency
			stepTotalRequestNum += cs.RequestNum
			stepTotalRequestTime += cs.RequestTime
			stepTotalSuccessNum += cs.SuccessNum

			if stepStartTime == 0 || stepStartTime > cs.StartTime {
				stepStartTime = cs.StartTime
			}
			if stepEndTime == 0 || stepEndTime < cs.EndTime {
				stepEndTime = cs.EndTime
			}

			// 判断前置场景是否结束
			if cs.SceneType != rao.PreFixScene && planResultData.PrefixSceneEndTime == "" {
				planResultData.PrefixSceneEndTime = time.Unix(tms, 0).Format(FullTimeFormat)
			}
		}
		for _, value := range sceneConcurrency {
			stepTotalConcurrency += value
		}

		planResultData.Graph.Rps = append(planResultData.Graph.Rps, [3]any{timeStr, stepTotalConcurrency, libs.CalcRpsNew(stepTotalRequestNum, calcDuration(stepStartTime, stepEndTime))})
		// 计算平均响应时间 = 总响应时长 / 总请求次数
		planResultData.Graph.RequestTime = append(planResultData.Graph.RequestTime, [2]any{timeStr, libs.CalcDiv(stepTotalRequestTime, stepTotalRequestNum)})
		planResultData.Graph.SuccessRate = append(planResultData.Graph.SuccessRate, [2]any{timeStr, libs.CalcRate(stepTotalSuccessNum, stepTotalRequestNum)})
		planResultData.Concurrency = stepTotalConcurrency
	}
	planResultData.StageRps = libs.CalcRpsNew(planResultData.StageRequestNum, calcDuration(planResultData.StageStartTime, planResultData.StageEndTime))
	planResultData.StageRT = libs.CalcDiv(planResultData.StageRequestTime, planResultData.StageRequestNum)
	planResultData.StageSuccessRate = libs.CalcRate(planResultData.StageSuccessNum, planResultData.StageRequestNum)

	planResultData.TotalSuccessRate = libs.CalcRate(planResultData.TotalSuccessNum, planResultData.TotalRequestNum)

	for _, scene := range planResultData.Scenes {
		for _, cs := range scene.Cases {
			if cs.ThreeXxCodeNum != 0 {
				planResultData.ErrorCodeArray = append(planResultData.ErrorCodeArray, rao.ApiDistribution{
					CaseID:    cs.CaseID,
					ErrorCode: "3xx",
					Count:     cs.ThreeXxCodeNum,
				})
			}
			if cs.FourXxCodeNum != 0 {
				planResultData.ErrorCodeArray = append(planResultData.ErrorCodeArray, rao.ApiDistribution{
					CaseID:    cs.CaseID,
					ErrorCode: "4xx",
					Count:     cs.FourXxCodeNum,
				})
			}
			if cs.FiveXxCodeNum != 0 {
				planResultData.ErrorCodeArray = append(planResultData.ErrorCodeArray, rao.ApiDistribution{
					CaseID:    cs.CaseID,
					ErrorCode: "5xx",
					Count:     cs.FiveXxCodeNum,
				})
			}
			if cs.TotalRequestErrorNum != 0 {
				planResultData.RequestErrorArray = append(planResultData.RequestErrorArray, rao.ApiDistribution{
					CaseID: cs.CaseID,
					Count:  cs.TotalRequestErrorNum,
				})
			}
			if cs.TotalAssertErrorNum != 0 {
				planResultData.AssertErrorArray = append(planResultData.AssertErrorArray, rao.ApiDistribution{
					CaseID: cs.CaseID,
					Count:  cs.TotalAssertErrorNum,
				})
			}
		}
	}

	for k, v := range errorUrlMp {
		if k == "" || v.Count == 0 {
			continue
		}
		planResultData.ErrorUrlArray = append(planResultData.ErrorUrlArray, rao.ApiDistribution{
			Url:        k,
			Count:      v.Count,
			ErrorCount: v.ErrorCount,
		})
	}
	sort.Sort(planResultData.ErrorUrlArray)
	sort.Sort(planResultData.ErrorCodeArray)
	sort.Sort(planResultData.RequestErrorArray)
	sort.Sort(planResultData.AssertErrorArray)

	// 报告结尾毛刺处理，从后往前去掉最多3个超时stage，然后将-1至替换为-2值
	if planResultData.End && len(timestampKeys) > 5 {
		for i := 1; i <= 3; i++ {
			if planResultData.Graph.RequestTime[len(timestampKeys)-i][1].(float64) > 4500 && planResultData.Graph.SuccessRate[len(timestampKeys)-i][1].(float64) < 1 {
				planResultData.Graph.Rps = planResultData.Graph.Rps[:len(timestampKeys)-i]
				planResultData.Graph.RequestTime = planResultData.Graph.RequestTime[:len(timestampKeys)-i]
				planResultData.Graph.SuccessRate = planResultData.Graph.SuccessRate[:len(timestampKeys)-i]
			}
		}
		graphLen := len(planResultData.Graph.Rps)
		planResultData.Graph.Rps[graphLen-1][2] = planResultData.Graph.Rps[graphLen-2][2]
		planResultData.Graph.RequestTime[graphLen-1][1] = planResultData.Graph.RequestTime[graphLen-2][1]
		planResultData.Graph.SuccessRate[graphLen-1][1] = planResultData.Graph.SuccessRate[graphLen-2][1]
	}

	planResultData.Graph.Rps = ReduceArray(planResultData.Graph.Rps)
	planResultData.Graph.RequestTime = ReduceArray(planResultData.Graph.RequestTime)
	planResultData.Graph.SuccessRate = ReduceArray(planResultData.Graph.SuccessRate)

	return
}

func GetCaseData(ctx *gin.Context, reportId, caseId int32) (planResultData rao.PlanResultData, err error) {
	baseData, err := GetReportBaseData(ctx, reportId)
	if err != nil {
		log.Logger.Errorf("logic.report.GetCaseData.GetReportBaseData err: %v", err)
		return
	}
	b, _ := dal.ReportRdb.Exists(ctx, fmt.Sprintf(ReportRunningLockKey, reportId)).Result()
	planResultData.End = b == 0
	stageCount := len(baseData)
	timestampKeys := make([]int64, 0, stageCount)
	for k := range baseData {
		timestampKeys = append(timestampKeys, k)
	}
	sort.Slice(timestampKeys, func(i, j int) bool {
		return timestampKeys[i] < timestampKeys[j] // 升序排序,如果需要降序排序，则改为 return slice[i] > slice[j]
	})
	planResultData.Graph = &rao.ReportDataGraphEntity{
		Rps:         make([][3]any, 0, stageCount),
		RequestTime: make([][2]any, 0, stageCount),
		SuccessRate: make([][2]any, 0, stageCount),
	}
	planResultData.Scenes = make([]*rao.SceneResultData, 0)

	idMap := make(map[int32]map[int32]struct{})

	for stage, tms := range timestampKeys {
		timeStr := time.Unix(tms, 0).Format(HourMinSec)
		var (
			stepStartTime        int64 = 0
			stepEndTime          int64 = 0
			stepTotalConcurrency int64 = 0 // 当前时间戳总并发数 = 每个Scene的并发数之和
			stepTotalRequestNum  int64 = 0
			stepTotalRequestTime int64 = 0
			stepTotalSuccessNum  int64 = 0
			sceneConcurrency           = make(map[int32]int64)
		)
		for _, cs := range baseData[tms] {
			if cs.CaseId != caseId {
				continue
			}
			// scene、case初始化
			if _, ok := idMap[cs.SceneId]; !ok {
				idMap[cs.SceneId] = make(map[int32]struct{})
				planResultData.Scenes = append(planResultData.Scenes, &rao.SceneResultData{
					SceneID:   cs.SceneId,
					SceneType: cs.SceneType,
					Cases:     make([]*rao.CaseResultData, 0),
				})
			}
			if _, hasCase := idMap[cs.SceneId][cs.CaseId]; !hasCase {
				idMap[cs.SceneId][cs.CaseId] = struct{}{}
				for i, sInfo := range planResultData.Scenes {
					if sInfo.SceneID == cs.SceneId {
						planResultData.Scenes[i].Cases = append(planResultData.Scenes[i].Cases, &rao.CaseResultData{
							CaseID: cs.CaseId,
						})
						break
					}
				}
			}

			// total数据填充
			for _, s := range planResultData.Scenes {
				if s.SceneID != cs.SceneId {
					continue
				}
				for _, c := range s.Cases {
					if c.CaseID != cs.CaseId {
						continue
					}
					c.TotalRequestTime += cs.RequestTime
					c.TotalRequestNum += cs.RequestNum
					c.TotalSuccessNum += cs.SuccessNum
					c.TotalRequestErrorNum += cs.RequestErrorNum
					c.TotalAssertErrorNum += cs.AssertErrorNum
					c.TotalErrorNum += cs.ErrorNum
					c.TotalSuccessRate = libs.CalcRate(c.TotalSuccessNum, c.TotalRequestNum)
					c.TwoXxCodeNum += cs.Normal2xxCodeNum
					c.ThreeXxCodeNum += cs.Normal3xxCodeNum
					c.FourXxCodeNum += cs.Normal4xxCodeNum
					c.FiveXxCodeNum += cs.Normal5xxCodeNum
					c.OtherCodeNum += cs.Normal1xxCodeNum + cs.Normal3xxCodeNum + cs.Normal4xxCodeNum + cs.Normal5xxCodeNum
					c.SendBytes += cs.SendBytes
					c.ReceivedBytes += cs.ReceivedBytes
					if c.TotalStartTime == 0 || c.TotalStartTime > cs.StartTime {
						c.TotalStartTime = cs.StartTime
					}
					if c.TotalEndTime == 0 || c.TotalEndTime < cs.EndTime {
						c.TotalEndTime = cs.EndTime
					}
					c.TotalRps = libs.CalcRpsNew(c.TotalRequestNum, c.TotalEndTime-c.TotalStartTime)
					c.TotalAvgRt = libs.CalcDiv(c.TotalRequestTime, c.TotalRequestNum)

					planResultData.TotalRequestTime += cs.RequestTime
					planResultData.TotalRequestNum += cs.RequestNum
					planResultData.TotalSuccessNum += cs.SuccessNum
					planResultData.TotalErrorNum += cs.ErrorNum
					planResultData.TotalSendBytes += cs.SendBytes
					planResultData.TotalReceivedBytes += cs.ReceivedBytes
					if planResultData.TotalStartTime == 0 || planResultData.TotalStartTime > cs.StartTime {
						planResultData.TotalStartTime = cs.StartTime
					}
					if planResultData.TotalEndTime == 0 || planResultData.TotalEndTime < cs.EndTime {
						planResultData.TotalEndTime = cs.EndTime
					}

					c.FiftyRequestTimeLineValue = cs.FiftyRequestTimeLineValue
					c.NinetyRequestTimeLineValue = cs.NinetyRequestTimeLineValue
					c.NinetyFiveRequestTimeLineValue = cs.NinetyFiveRequestTimeLineValue
					c.NinetyNineRequestTimeLineValue = cs.NinetyNineRequestTimeLineValue

					// stage数据填充
					if stage == len(timestampKeys)-1 {
						c.StageRequestTime = cs.RequestTime
						c.StageRequestNum = cs.RequestNum
						c.StageSuccessNum = cs.SuccessNum
						c.StageErrorNum = cs.ErrorNum

						if c.StageStartTime == 0 || c.StageStartTime > cs.StartTime {
							c.StageStartTime = cs.StartTime
						}
						if c.StageEndTime == 0 || c.StageEndTime < cs.EndTime {
							c.StageEndTime = cs.EndTime
						}
						c.StageRps = libs.CalcRpsNew(c.StageRequestNum, calcDuration(c.StageStartTime, c.StageEndTime))
						c.StageSuccessRate = libs.CalcRate(c.StageSuccessNum, c.StageRequestNum)
						c.StageAvgRt = libs.CalcDiv(c.StageRequestTime, c.StageRequestNum)

						//c.MaxRt = mathutil.Max(c.MaxRt, cs.MaxRequestTime)
						//if c.MinRt != 0 {
						//	c.MinRt = mathutil.Min(c.MinRt, cs.MinRequestTime)
						//} else {
						//	c.MinRt = cs.MinRequestTime
						//}

						s.Concurrency = cs.ActualConcurrency

						planResultData.StageRequestTime += cs.RequestTime
						planResultData.StageRequestNum += cs.RequestNum
						planResultData.StageSuccessNum += cs.SuccessNum
						planResultData.StageErrorNum += cs.ErrorNum
						if planResultData.StageStartTime == 0 || planResultData.StageStartTime > cs.StartTime {
							planResultData.StageStartTime = cs.StartTime
						}
						if planResultData.StageEndTime == 0 || planResultData.StageEndTime < cs.EndTime {
							planResultData.StageEndTime = cs.EndTime
						}
					}

				}
			}
			sceneConcurrency[cs.SceneId] = cs.ActualConcurrency
			stepTotalRequestNum += cs.RequestNum
			stepTotalRequestTime += cs.RequestTime
			stepTotalSuccessNum += cs.SuccessNum

			if stepStartTime == 0 || stepStartTime > cs.StartTime {
				stepStartTime = cs.StartTime
			}
			if stepEndTime == 0 || stepEndTime < cs.EndTime {
				stepEndTime = cs.EndTime
			}
		}
		for _, value := range sceneConcurrency {
			stepTotalConcurrency += value
		}

		planResultData.Graph.Rps = append(planResultData.Graph.Rps, [3]any{timeStr, stepTotalConcurrency, libs.CalcRpsNew(stepTotalRequestNum, calcDuration(stepStartTime, stepEndTime))})
		// 计算平均响应时间 = 总响应时长 / 总请求次数
		planResultData.Graph.RequestTime = append(planResultData.Graph.RequestTime, [2]any{timeStr, libs.CalcDiv(stepTotalRequestTime, stepTotalRequestNum)})
		planResultData.Graph.SuccessRate = append(planResultData.Graph.SuccessRate, [2]any{timeStr, libs.CalcRate(stepTotalSuccessNum, stepTotalRequestNum)})
		planResultData.Concurrency = stepTotalConcurrency
	}
	planResultData.StageRps = libs.CalcRpsNew(planResultData.StageRequestNum, calcDuration(planResultData.StageStartTime, planResultData.StageEndTime))
	planResultData.StageRT = libs.CalcDiv(planResultData.StageRequestTime, planResultData.StageRequestNum)
	planResultData.StageSuccessRate = libs.CalcRate(planResultData.StageSuccessNum, planResultData.StageRequestNum)
	planResultData.TotalSuccessRate = libs.CalcRate(planResultData.TotalSuccessNum, planResultData.TotalRequestNum)

	planResultData.Graph.Rps = ReduceArray(planResultData.Graph.Rps)
	planResultData.Graph.RequestTime = ReduceArray(planResultData.Graph.RequestTime)
	planResultData.Graph.SuccessRate = ReduceArray(planResultData.Graph.SuccessRate)
	return
}

// GetReportBaseData 查询原始数据
// map key为时间戳，value为当前时间戳的值
func GetReportBaseData(ctx *gin.Context, reportId int32) (recordItems map[int64][]rao.ReportRecordEntity, err error) {
	var (
		cursor         = int64(0) // 时间戳，秒级
		newRecordItems = make(map[int64][]rao.ReportRecordEntity)
	)
	recordItems = make(map[int64][]rao.ReportRecordEntity)
	cache, hasCache := dal.LocalCacheGet(fmt.Sprintf(ReportDataCacheKey, reportId))
	if hasCache {
		recordItems = cache.(map[int64][]rao.ReportRecordEntity)
		for ts := range recordItems {
			cursor = mathutil.Max(cursor, ts)
		}
	}
	if hasCache && time.Now().Unix()-cursor < 5 {
		// 有数据且时差<5s 说明还没有下一个stage的数据，不查询
		return
	}
	// 查db
	query := fmt.Sprintf(`from(bucket: "%s")
			|> range(start: %d, stop:-2s)
			|> filter(fn: (r) => r._measurement == "%d")
			|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`, conf.Conf.InfluxDB.Bucket, cursor+1, reportId)
	result, err := dal.Query(query)
	if err != nil {
		log.Logger.Error("login.report.GetReportBaseData.Query err:", err)
		return
	}
	for result.Next() {
		recordTimestamp := result.Record().Time().Unix()
		cursor = mathutil.Max(cursor, recordTimestamp)

		// 序列化当前时间戳Record信息
		var baseRecord rao.ReportRecordBaseEntity
		jsonBytes, _ := json.Marshal(result.Record().Values())
		_ = json.Unmarshal(jsonBytes, &baseRecord)
		caseId, _ := strconv.Atoi(baseRecord.CaseId)
		sceneId, _ := strconv.Atoi(baseRecord.SceneId)
		sceneTp, _ := strconv.Atoi(baseRecord.SceneType)
		record := rao.ReportRecordEntity{
			CaseId:            int32(caseId),
			SceneId:           int32(sceneId),
			SceneType:         int32(sceneTp),
			CollectorCaseData: baseRecord.CollectorCaseData,
		}
		if _, ok := newRecordItems[recordTimestamp]; !ok {
			newRecordItems[recordTimestamp] = make([]rao.ReportRecordEntity, 0)
		}
		newRecordItems[recordTimestamp] = append(newRecordItems[recordTimestamp], record)
	}
	for k, v := range newRecordItems {
		recordItems[k] = v
	}
	if len(newRecordItems) > 0 {
		dal.LocalCacheSet(fmt.Sprintf(ReportDataCacheKey, reportId), recordItems)
	}
	return
}

// GetReportRps 获取测试报告当前Rps明细
func GetReportRps(ctx *gin.Context, reportId int32) (res *rao.ReportRpsResponse, err error) {
	var (
		scenes   = make([]rao.ReportStatusSceneEntity, 0)
		sceneMap = make(map[int32]rao.ReportStatusSceneEntity)
		itemMap  = make(map[int32]map[string]interface{})
	)

	// 从InfluxDB获取测试报告最新状态数据
	query := fmt.Sprintf(`from(bucket: "%s")
			|> range(start: 0, stop:-2s)
			|> filter(fn: (r) => r._measurement == "%d")
			|> last()`, conf.Conf.InfluxDB.Bucket, reportId)

	result, err := dal.Query(query)
	if err != nil {
		log.Logger.Error("logic.report.GetReportStatus.Query ，err:", err)
		return
	}

	for result.Next() {
		record := result.Record()
		sceneId, _ := strconv.ParseInt(record.ValueByKey("scene_id").(string), 10, 32)
		caseId, _ := strconv.ParseInt(record.ValueByKey("case_id").(string), 10, 32)
		item := itemMap[int32(caseId)]
		if item == nil {
			item = make(map[string]interface{})
		}
		item["scene_id"] = int32(sceneId)
		item["case_id"] = int32(caseId)
		item[record.Field()] = record.Value()
		itemMap[int32(caseId)] = item
	}

	for _, value := range itemMap {
		sceneId := value["scene_id"].(int32)
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			log.Logger.Error("Error marshalling to JSON:", err)
			continue
		}
		var itemCase rao.ReportStatusCaseEntity
		err = json.Unmarshal(jsonBytes, &itemCase)
		if err != nil {
			log.Logger.Error("Error unmarshalling from JSON:", err)
			continue
		}
		// 计算当前RPS = 请求数 / (结束时间 - 开始时间) * 1000
		itemCase.CurrentRps = libs.CalcRps(itemCase.RequestNum, itemCase.StartTime, itemCase.EndTime, false)
		scene := sceneMap[sceneId]
		if _, ok := sceneMap[sceneId]; !ok {
			scene = rao.ReportStatusSceneEntity{SceneId: sceneId, Cases: make([]rao.ReportStatusCaseEntity, 0)}
		}
		scene.Cases = append(scene.Cases, itemCase)
		sceneMap[sceneId] = scene
	}

	for _, value := range sceneMap {
		scenes = append(scenes, value)
	}

	return &rao.ReportRpsResponse{
		ReportID: reportId,
		Scenes:   scenes,
	}, nil
}

func GetReportTargetRps(ctx *gin.Context, reportId int32) (rps int64, err error) {
	targetTotalRpsStr, err := dal.ReportRdb.Get(ctx, fmt.Sprintf(RpsValKey, reportId)).Result()
	if err != nil {
		return 0, err
	}
	targetTotalRps, _ := strconv.Atoi(targetTotalRpsStr)
	return int64(targetTotalRps), nil
}

func calcDuration(startTime, endTime int64) int64 {
	// 解决rps模式，统计时差的问题
	if endTime-startTime > 4000 {
		return 5000
	}
	return endTime - startTime
}

// ReduceArray 裁剪超长的图数据
func ReduceArray[T any](arr []T) []T {
	length := len(arr)

	// 如果数组长度小于等于 120，直接返回原数组
	if length <= GraphMaxLen {
		return arr
	}

	// 目标大小为 120
	targetSize := GraphMaxLen

	// 计算步长（需要裁剪中间的数据）
	step := float64(length-2) / float64(targetSize-2) // -2 是为了保留首位和末位
	var result []T

	// 保留第一位
	result = append(result, arr[0])

	// 按照步长均匀采样中间的数据
	for i := 1; i < targetSize-1; i++ {
		index := int(float64(i) * step)
		result = append(result, arr[index])
	}

	// 保留最后一位
	result = append(result, arr[length-1])

	return result
}
