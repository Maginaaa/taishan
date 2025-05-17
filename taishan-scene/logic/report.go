package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gen"
	"scene/internal/biz/log"
	"scene/internal/conf"
	"scene/internal/dal"
	"scene/internal/model"
	"scene/rao"
	"scene/utils"
	"sync"
	"time"
)

const (
	ReportRunningLockKey = "report:lock:%d"
	ReportDataCacheKey   = "report:data:%d"

	RequestFailType = 1
	AssertFailType  = 2
)

func createReport(ctx *gin.Context, plan *rao.Plan, engineList []string) (reportId int32, err error) {
	req := rao.CreateReportReq{
		PlanID:       plan.PlanID,
		PlanName:     plan.PlanName,
		Duration:     int32(plan.PressInfo.Duration),
		PressType:    int32(plan.PressInfo.PressType),
		EngineList:   engineList,
		Rps:          int32(plan.PressInfo.RPS),
		CreateUserID: utils.GetCurrentUserID(ctx),
	}
	res, err := resty.New().R().SetBody(req).Post(conf.Conf.Url.Report + "/report/create")
	if err != nil {
		return 0, errors.New("创建报告失败")
	}
	var resp rao.CommonResponse[int32]
	if err = json.Unmarshal(res.Body(), &resp); err != nil {
		return 0, errors.New("创建报告解析失败")
	}
	reportId = resp.Data
	return
}

// ReportPressDown 压测被动结束，清除压测报告和测试计划状态
func ReportPressDown(ctx *gin.Context, reportId int32) (err error) {
	reportData := GetReportData(ctx, reportId)
	if reportData == nil {
		log.Logger.Error("获取测试报告data失败")
		return
	}
	totalReqTime := reportData.Data.TotalRequestTime
	pressStartTime := reportData.Data.TotalStartTime
	pressEndTime := reportData.Data.TotalEndTime
	if pressEndTime == 0 {
		log.Logger.Error("pressEndTime is 0, err")
		pressEndTime = time.Now().UnixMilli()
	}

	// 修改mysql
	tx := dal.GetQuery().Report
	report, err := tx.WithContext(ctx).Where(tx.ID.Eq(reportId)).First()
	machineNum := int64(len(report.EngineList))
	// vum = (((总请求时间 / (压测时间 * 机器数)) / 500) + 1) * 500 * 机器数 * 压测时长(分)
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
		return err
	}
	dal.ReportRdb.Expire(ctx, fmt.Sprintf(ReportDataCacheKey, reportId), time.Duration(3)*time.Hour)
	err = PlanRunningDelLock(ctx, report.PlanID)
	if err != nil {
		log.Logger.Error("logic.report.reportStopPress.PlanRunningDelLock ，err:", err)
		return err
	}
	log.Logger.Infof("测试报告: %d生成测试报告完成...", reportId)
	return
}

// UpdatePress 调整并发数
func UpdatePress(ctx *gin.Context, req rao.ConcurrencyChange) (err error) {
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
		return
	}
	return
}

func GetReportDebugInfo(ctx *gin.Context, debugInfo rao.SamplingDataReq) (resp rao.SamplingDataResp, err error) {
	filter := bson.M{}
	if len(debugInfo.CaseIdList) > 0 {
		filter["case_id"] = bson.M{"$in": debugInfo.CaseIdList}
	}

	if debugInfo.SceneID != 0 {
		filter["scene_id"] = bson.M{"$eq": debugInfo.SceneID}
	}

	if len(debugInfo.CodeList) > 0 {
		filter["status_code"] = bson.M{"$in": debugInfo.CodeList}
	}

	if len(debugInfo.StatusList) > 0 {
		orConditions := []bson.M{}
		if slice.Contain(debugInfo.StatusList, RequestFailType) {
			orConditions = append(orConditions, bson.M{"request_success": false})
		}
		if slice.Contain(debugInfo.StatusList, AssertFailType) {
			orConditions = append(orConditions, bson.M{"assert_success": false})
		}
		filter["$or"] = orConditions
	}

	if debugInfo.MinRt != 0 || debugInfo.MaxRt != 0 {
		responseTimeCondition := bson.M{}
		if debugInfo.MinRt != 0 {
			responseTimeCondition["$gte"] = debugInfo.MinRt
		}
		if debugInfo.MaxRt != 0 {
			responseTimeCondition["$lte"] = debugInfo.MaxRt
		}
		filter["response_time"] = responseTimeCondition
	}

	if debugInfo.StartTime != "" || debugInfo.EndTime != "" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		startTimeCondition := bson.M{}
		if debugInfo.StartTime != "" {
			startTime, err := time.ParseInLocation(FullTimeFormat, debugInfo.StartTime, loc)
			if err != nil {
				log.Logger.Error("logic.report.GetReportDebugInfo.timeParseStartTime, err:", err)
			}
			startTimeCondition["$gte"] = startTime.UnixMilli()
		}
		if debugInfo.EndTime != "" {
			endTime, err := time.ParseInLocation(FullTimeFormat, debugInfo.EndTime, loc)
			if err != nil {
				log.Logger.Error("logic.report.GetReportDebugInfo.timeParseEndTime, err:", err)
			}
			startTimeCondition["$lte"] = endTime.UnixMilli()
		}
		filter["start_time"] = startTimeCondition
	}

	findOptions := options.Find().SetSort(bson.D{{"_id", -1}}). // 按照 _id 降序排列
									SetSkip(int64((debugInfo.Page - 1) * debugInfo.PageSize)).
									SetLimit(int64(debugInfo.PageSize))
	cursor, err := dal.GetMongoCollection(fmt.Sprintf("%d", debugInfo.ReportID)).Find(ctx, filter, findOptions)
	// 查询total
	resp.Total, err = dal.GetMongoCollection(fmt.Sprintf("%d", debugInfo.ReportID)).CountDocuments(ctx, filter)
	if err != nil {
		log.Logger.Error("logic.report.GetReportDebugInfo.GetMongoCollection, err:", err)
		return
	}

	var results []rao.HttpResponse
	err = cursor.All(ctx, &results)
	if err != nil {
		log.Logger.Error("logic.report.GetReportDebugInfo.cursorAll, err:", err)
	}
	samplingList := make([]rao.SamplingData, 0)
	for _, res := range results {
		samplingList = append(samplingList, rao.SamplingData{
			StartTime:        time.UnixMilli(res.StartTime).Format(FullTimeFormat),
			Success:          res.RequestSuccess && res.AssertSuccess,
			CaseID:           res.CaseID,
			CaseName:         res.CaseName,
			ResponseTime:     res.ResponseTime,
			StatusCode:       res.StatusCode,
			HttpResponseData: res,
		})
	}
	resp.List = samplingList
	return
}

func GetReportListByPlanID(ctx *gin.Context, planID int32) (res []*rao.Report, err error) {
	tx := dal.GetQuery().Report
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.PlanID.Eq(planID))
	conditions = append(conditions, tx.IsDelete.Not())
	reportList, err := tx.WithContext(ctx).Where(conditions...).Order(tx.ID.Desc()).Find()
	if err != nil {
		log.Logger.Error("logic.report.GetReportListByPlanID.Find ，err:", err)
		return nil, err
	}

	// 获取用户信息列表
	var users []rao.User
	users, err = GetUserList()
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}

	planIds := make([]int32, 0)
	for _, s := range reportList {
		planIds = append(planIds, s.PlanID)
	}
	if err != nil {
		log.Logger.Error("logic.report.GetReportList.BatchGetPlanName ，err:", err)
		return
	}
	res = make([]*rao.Report, 0, len(reportList))
	for _, s := range reportList {
		actualDuration, _ := decimal.NewFromFloat(s.EndTime.Sub(s.CreateTime).Minutes()).Round(2).Float64()
		res = append(res, &rao.Report{
			ReportID:       s.ID,
			Status:         s.Status,
			PlanID:         s.PlanID,
			Duration:       s.Duration,
			ActualDuration: actualDuration,
			PressType:      s.PressType,
			StartTime:      s.CreateTime.Format(FullTimeFormat),
			CreateUserName: GetNameByID(users, s.CreateUserID),
			UpdateUserName: GetNameByID(users, s.UpdateUserID),
			UpdateTime:     s.UpdateTime.Format(FullTimeFormat),
		})
	}
	return res, nil
}

func DoReportRpsModify(ctx *gin.Context, req rao.DoReportRpsModifyReq) (err error) {
	_, err = dal.ReportRdb.Set(ctx, fmt.Sprintf(RpsValKey, req.ReportId), req.NewValue, 120*time.Minute).Result()
	return
}

//func DoReportRpsInfoGet(ctx *gin.Context, reportId int32) map[int32]int64 {
//	result := make(map[int32]int64)
//	reportDetail := getReportDetail(ctx, reportId)
//	var sceneResultList []rao.SceneInformation
//	_ = json.Unmarshal([]byte(reportDetail.PressResult), &sceneResultList)
//	for _, scene := range sceneResultList {
//		for _, caseInfo := range scene.Cases {
//			result[caseInfo.CaseId] = caseInfo.TargetRps
//		}
//	}
//	return result
//}

func getReportDetail(ctx *gin.Context, id int32) *model.Report {

	tx := dal.GetQuery().Report
	reportInfo, err := tx.WithContext(ctx).Where(tx.ID.Eq(id)).First()
	if err != nil {
		log.Logger.Error("logic.report.GetReportDetail.First ，err:", err)
		return nil
	}
	return reportInfo
}

func GetReportData(ctx *gin.Context, reportId int32) *rao.CommonResponse[rao.ReportData] {
	runResponse, err := resty.New().R().Get(fmt.Sprintf("%s/report/data/%d", conf.Conf.Url.Report, reportId))
	if err != nil {
		log.Logger.Error("report服务返回异常")
		return nil
	}
	var resp rao.CommonResponse[rao.ReportData]
	err = json.Unmarshal(runResponse.Body(), &resp)
	if err != nil {
		log.Logger.Error("report服务返回异常")
		return nil
	}
	return &resp
}

func GetReportBottleneck(ctx *gin.Context, reportId int32) {
	reportData := GetReportData(ctx, reportId)
	fmt.Println("查询报告")
	serverNameCh := make(chan string, 1000)
	wg := &sync.WaitGroup{}
	for _, scene := range reportData.Data.Scenes {
		for _, cas := range scene.Cases {
			if cas.StageRequestTime > 40 || cas.StageSuccessRate < 99 {
				getCaseServer(reportId, cas, wg, serverNameCh)
			}
		}
	}

	wg.Wait()
	close(serverNameCh)
	// 打印serverNameCh内的所有内容
	for serverName := range serverNameCh {
		fmt.Println(serverName)
	}

}

func getCaseServer(reportId int32, caseData *rao.CaseResultData, wg *sync.WaitGroup, ch chan string) {
	listResp, _ := resty.New().R().SetBody(rao.TraceListReq{
		CaseId:   caseData.CaseID,
		ReportId: reportId,
	}).Post(fmt.Sprintf("%s/arms/trace/list", conf.Conf.Url.Data))
	var resp rao.CommonResponse[[]rao.TraceListResp]
	_ = json.Unmarshal(listResp.Body(), &resp)
	if resp.Data == nil {
		return
	}
	for idx, traceDetail := range resp.Data {
		if idx > 10 {
			break
		}
		wg.Add(1)
		go func(traceId string, funcWg *sync.WaitGroup) {
			defer funcWg.Done()
			traceResp, _ := resty.New().R().SetBody(rao.TraceDetailReq{
				TraceId:  traceId,
				ReportId: reportId,
			}).Post(fmt.Sprintf("%s/arms/trace/detail", conf.Conf.Url.Data))
			var tResp rao.CommonResponse[[]rao.TraceTreeNode]
			_ = json.Unmarshal(traceResp.Body(), &tResp)
			serverName := make([]string, 0)
			for _, dt := range tResp.Data {
				flattenTree(dt, &serverName)
			}
			for _, sv := range serverName {
				ch <- sv
			}
		}(traceDetail.TraceID, wg)
	}
}

func flattenTree(root rao.TraceTreeNode, result *[]string) {
	// 将当前节点的值添加到结果中
	if root.App != "" {
		*result = append(*result, root.App)
	}
	// 递归处理子节点
	for _, child := range root.Children {
		flattenTree(child, result)
	}
}
