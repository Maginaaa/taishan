package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"gorm.io/datatypes"
	"gorm.io/gen"
	"net/http"
	"scene/internal/biz/log"
	"scene/internal/conf"
	"scene/internal/dal"
	"scene/internal/model"
	"scene/rao"
	"scene/utils"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	partitionLock sync.Mutex
)

const (
	RpsModelCaseDefaultRT = 50
	RpsValKey             = "report:rps:%d"
	NULL                  = "null"
)

func PlanInfoCheck(ctx *gin.Context, plan *rao.Plan) (err error) {

	if !(plan.PressInfo.Duration > 0) {
		return fmt.Errorf("压测时长需大于 0")
	}
	switch plan.PressInfo.PressType {
	case rao.ConcurrentModel:
		for _, scene := range plan.PressInfo.SceneList {
			if scene.SceneType == rao.PreFixScene {
				continue
			}
			if scene.Concurrency > 0 {
				continue
			}
			return fmt.Errorf("场景id: %d, 并发数需大于0", scene.SceneId)
		}
		break
	case rao.StepModel:
		if !(plan.PressInfo.Concurrency > 0) {
			return fmt.Errorf("并发数需大于 0")
		}
		totalRate := int64(0)
		for _, sc := range plan.PressInfo.SceneList {
			if sc.SceneType == rao.PreFixScene {
				continue
			}
			totalRate += sc.Rate
		}
		//所有场景的流量百分比不等于100时返回错误
		if totalRate != 100 {
			return fmt.Errorf("场景流量百分比合不为100")
		}
		break
	case rao.PlanRpsRateMode:
		totalRate := int64(0)
		for _, sc := range plan.PressInfo.SceneList {
			if sc.SceneType == rao.PreFixScene {
				continue
			}
			totalRate += sc.Rate
		}
		//所有场景的流量百分比不等于100时返回错误
		if totalRate != 100 {
			return fmt.Errorf("场景流量百分比合不为100")
		}
		break
	case rao.FixedFrequency:
		for _, scene := range plan.PressInfo.SceneList {
			if scene.SceneType == rao.PreFixScene {
				continue
			}
			if scene.Concurrency > 0 && scene.Iteration > 0 {
				continue
			}
			return fmt.Errorf("场景id: %d, 并发数和压测次数均需大于0", scene.SceneId)
		}
		break
	}

	//pressInfoList, err := getPressInfo(ctx, plan.PressInfo)
	//// 检查集合点信息
	//err = CheckRallyPoint(plan, pressInfoList)
	//if err != nil {
	//	return err
	//}
	return
}

func CreatePlan(ctx *gin.Context, req *rao.Plan) (planId int32, err error) {
	tx := dal.GetQuery().Plan

	emptyPressInfo := rao.PressInfo{
		SceneList: []rao.PressSceneInfo{},
	}
	pressInfo, err := json.Marshal(emptyPressInfo)
	if err != nil {
		return 0, err
	}
	var emptySamplingInfo rao.SamplingInfo
	samplingInfo, _ := json.Marshal(emptySamplingInfo)

	emptyParam, _ := json.Marshal([]rao.ParamsForm{
		{
			Enable: true,
			Key:    "",
			Value:  "",
			Desc:   "",
		},
	})

	insertData := &model.Plan{
		PlanName:       req.PlanName,
		EngineCount:    1,
		Remark:         req.Remark,
		PressInfo:      string(pressInfo),
		GlobalVariable: string(emptyParam),
		DefaultHeader:  string(emptyParam),
		SamplingInfo:   string(samplingInfo),
		ServerInfo:     "[]",
		Tag:            "[]",
		BreakType:      rao.NotBreak,
		DebugStatus:    false,
		CreateUserID:   utils.GetCurrentUserID(ctx),
	}
	err = tx.WithContext(ctx).Create(insertData)
	if err != nil {
		return 0, err
	}
	planId = insertData.ID
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourcePlan,
		SourceID:      planId,
		OperationType: rao.CreateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    insertData,
	})
	createSceneData := &rao.Scene{
		SceneName: "默认场景",
		PlanID:    planId,
	}
	_, err = CreateScene(ctx, createSceneData)
	if err != nil {
		return 0, err
	}
	return
}

// TODO: 疑似冗余逻辑，大促压测结束后再次检查
func doRpsRateModify(pressInfo rao.PressInfo) rao.PressInfo {
	if pressInfo.PressType == rao.PlanRpsRateMode {
		totalRps := pressInfo.RPS
		for _, scene := range pressInfo.SceneList {
			scene.Concurrency = totalRps * scene.RpsRate
		}
	}
	return pressInfo
}

func UpdatePlan(ctx *gin.Context, req *rao.Plan) (isUpdate bool, err error) {
	planTx := dal.GetQuery().Plan
	beforePlan, _ := planTx.WithContext(ctx).Where(planTx.ID.Eq(req.PlanID)).First()
	req.PressInfo = doRpsRateModify(req.PressInfo)
	pressInfo, err := json.Marshal(req.PressInfo)
	if err != nil {
		return false, err
	}
	taskId, _ := UpdatePlanTask(ctx, req)
	samplingInfo, _ := json.Marshal(req.SamplingInfo)
	variableInfo, _ := json.Marshal(req.GlobalVariable)
	headerInfo, _ := json.Marshal(req.DefaultHeader)
	serverInfo, _ := json.Marshal(req.ServerInfo)
	tagList, _ := json.Marshal(req.TagList)
	updateData := &model.Plan{
		PlanName:       req.PlanName,
		EngineCount:    req.EngineCount,
		Remark:         req.Remark,
		PressInfo:      string(pressInfo),
		SamplingInfo:   string(samplingInfo),
		GlobalVariable: string(variableInfo),
		DefaultHeader:  string(headerInfo),
		ServerInfo:     string(serverInfo),
		Tag:            string(tagList),
		BreakType:      req.BreakType,
		BreakValue:     req.BreakValue,
		TaskID:         taskId,
		UpdateUserID:   utils.GetCurrentUserID(ctx),
		UpdateTime:     time.Now(),
	}
	_, err = planTx.WithContext(ctx).Where(planTx.ID.Eq(req.PlanID)).Updates(updateData)
	if err != nil {
		log.Logger.Error("logic.plan.UpdatePlan.Updates(), err:", err)
		return false, err
	}
	afterPlan, _ := planTx.WithContext(ctx).Where(planTx.ID.Eq(req.PlanID)).First()

	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourcePlan,
		SourceID:      req.PlanID,
		OperationType: rao.UpdateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   beforePlan,
		ValueAfter:    afterPlan,
	})
	return true, nil
}

func UpdatePlanBaseInfo(ctx *gin.Context, req *rao.Plan) (isUpdate bool, err error) {
	planTx := dal.GetQuery().Plan
	beforePlan, _ := planTx.WithContext(ctx).Where(planTx.ID.Eq(req.PlanID)).First()

	updateData := &model.Plan{
		PlanName:     req.PlanName,
		Remark:       req.Remark,
		UpdateUserID: utils.GetCurrentUserID(ctx),
		UpdateTime:   time.Now(),
	}
	_, err = planTx.WithContext(ctx).Where(planTx.ID.Eq(req.PlanID)).Updates(updateData)
	if err != nil {
		log.Logger.Error("logic.plan.UpdatePlan.Save(), err:", err)
		return false, err
	}
	afterPlan, _ := planTx.WithContext(ctx).Where(planTx.ID.Eq(req.PlanID)).First()

	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourcePlan,
		SourceID:      req.PlanID,
		OperationType: rao.UpdateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   beforePlan,
		ValueAfter:    afterPlan,
	})
	return true, nil
}

func DeletePlan(ctx *gin.Context, id int32) (isDelete bool, err error) {
	tx := dal.GetQuery().Plan
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(id)).UpdateSimple(tx.IsDelete.Value(true))
	if err != nil {
		log.Logger.Error("logic.plan.DeletePlan.UpdateSimple，err:", err)
		return false, err
	}
	return true, nil
}

func CopyPlan(ctx *gin.Context, req rao.Plan) (int32, error) {
	tx := dal.GetQuery().Plan
	plan, err := tx.WithContext(ctx).Where(tx.ID.Eq(req.PlanID)).First()
	if err != nil {
		log.Logger.Error("logic.plan.CopyPlan.selectPlan，err:", err)
		return 0, err
	}
	plan.ID = 0
	plan.DebugStatus = false
	plan.PlanName = req.PlanName
	plan.Remark = req.Remark
	plan.CreateTime = time.Now()
	plan.UpdateTime = time.Now()
	plan.CreateUserID = utils.GetCurrentUserID(ctx)
	plan.TaskID = 0
	err = tx.WithContext(ctx).Create(plan)
	if err != nil {
		log.Logger.Error("logic.plan.CopyPlan.createPlan，err:", err)
		return 0, err
	}
	copyPlanFile(ctx, plan.ID, req.PlanID)
	scenes, err := getPlanScenes(ctx, req.PlanID)

	// 复制场景压测策略
	sceneIDMapping := make(map[int32]int32)
	for _, scene := range scenes {
		targetSceneID, _ := CopyScene(ctx, scene.ID, plan.ID)
		sceneIDMapping[scene.ID] = targetSceneID
	}
	var pressInfo rao.PressInfo
	err = json.Unmarshal([]byte(plan.PressInfo), &pressInfo)

	for i := range pressInfo.SceneList {
		pressInfo.SceneList[i].SceneId = sceneIDMapping[pressInfo.SceneList[i].SceneId]
	}
	ms, _ := json.Marshal(pressInfo)
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(plan.ID)).UpdateSimple(tx.PressInfo.Value(string(ms)))
	if err != nil {
		log.Logger.Error("logic.plan.CopyPlan.UpdateSimple，err:", err)
		return plan.ID, err
	}

	return plan.ID, nil
}

func GetPlanList(ctx *gin.Context, req rao.PlanListReq) (res rao.PageResponse[rao.Plan], err error) {

	tx := dal.GetQuery().Plan
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	conditions := make([]gen.Condition, 0)
	if req.CreateUserId != 0 {
		conditions = append(conditions, tx.CreateUserID.Eq(req.CreateUserId))
	}

	if req.CaseInfo != "" {
		type scene struct {
			SceneID int32 `json:"scene_id"`
		}
		var sceneIds []int32
		sql := fmt.Sprintf(`select scene_id from scene_case where is_delete = 0 and ( extend->>'$.url' like '%%%s%%'  or extend->>'$.case_name' like '%%%s%%' )`, req.CaseInfo, req.CaseInfo)
		rows, _ := dal.MysqlDB().Raw(sql).Rows()
		for rows.Next() {
			var s scene
			dal.MysqlDB().ScanRows(rows, &s)
			sceneIds = slice.AppendIfAbsent(sceneIds, s.SceneID)
		}
		sceneTx := dal.GetQuery().Scene
		var planIds []int32
		scenes, _ := sceneTx.WithContext(ctx).Where(sceneTx.ID.In(sceneIds...), sceneTx.IsDelete.Not()).Select(sceneTx.PlanID).Find()
		for _, s := range scenes {
			planIds = slice.AppendIfAbsent(planIds, s.PlanID)
		}
		conditions = append(conditions, tx.ID.In(planIds...))
	}
	if len(req.TagList) > 0 {
		where := tx.WithContext(ctx).Where()
		for _, tag := range req.TagList {
			cond := gen.Cond(datatypes.JSONArrayQuery("tag").Contains(tag))
			where = where.Or(cond...)
		}
		conditions = append(conditions, where)
	}

	if req.PlanInfo != "" {
		where := tx.WithContext(ctx).Where(tx.PlanName.Like("%" + req.PlanInfo + "%"))
		planId, err := strconv.Atoi(req.PlanInfo)
		if err == nil {
			where = where.Or(tx.ID.Eq(int32(planId)))
		}
		conditions = append(conditions, where)

	}

	conditions = append(conditions, tx.IsDelete.Not())
	planList, total, err := tx.WithContext(ctx).Where(conditions...).Order(tx.UpdateTime.Desc()).FindByPage(offset, limit)

	if err != nil {
		log.Logger.Error("logic.plan.GetPlanList.planListQuery(), err:", err)
		return
	}

	// 获取用户信息列表
	users, err := GetUserList()
	if err != nil {
		log.Logger.Error("logic.plan.GetPlanList.GetUserList(), err:", err)
		return res, errors.New("获取用户信息失败")
	}

	pList := make([]rao.Plan, 0, len(planList))
	for _, v := range planList {
		// TODO: 后面需要将获取逻辑调整为mget
		isRunning, lock, err := getPlanRunningLock(ctx, v.ID)
		if err != nil {
			log.Logger.Error("logic.plan.GetPlanList.PlanIsRunning(), err:", err)
			return res, err
		}

		var tagList []int32
		if v.Tag == NULL || v.Tag == "" {
			tagList = make([]int32, 0)
		} else {
			_ = json.Unmarshal([]byte(v.Tag), &tagList)
		}

		plan := rao.Plan{
			PlanID:         v.ID,
			PlanName:       v.PlanName,
			IsRunning:      isRunning,
			CreateUserName: GetNameByID(users, v.CreateUserID),
			UpdateUserName: GetNameByID(users, v.UpdateUserID),
			UpdateTime:     v.UpdateTime.Format(FullTimeFormat),
			Remark:         v.Remark,
			TagList:        tagList,
			DebugSuccess:   v.DebugStatus,
		}

		if isRunning {
			plan.ReportID = lock.ReportID
		}

		pList = append(pList, plan)
	}
	sort.Slice(pList, func(i, j int) bool {
		return pList[i].IsRunning
	})
	res.List = pList
	res.Total = total
	return
}

func GetPlanDetail(ctx *gin.Context, planId int32) (plan *rao.Plan, err error) {

	planTx := dal.GetQuery().Plan

	planInfo, err := planTx.WithContext(ctx).Where(planTx.ID.Eq(planId)).First()
	if err != nil {
		log.Logger.Error("logic.plan.GetPlanDetail.planInfoQuery(),err:", err)
		return nil, err
	}
	var taskInfo rao.TaskInfo
	if planInfo.TaskID != 0 {
		task, _ := GetTaskDetail(ctx, planInfo.TaskID)
		taskInfo = rao.TaskInfo{
			TaskId: planInfo.TaskID,
			Enable: task.Enable,
			Cron:   task.Cron,
		}
	}

	var pressInfo rao.PressInfo
	err = json.Unmarshal([]byte(planInfo.PressInfo), &pressInfo)
	if err != nil {
		log.Logger.Error("logic.plan.GetPlanDetail.jsonUnmarshal(pressInfo), err：", err)
		return
	}
	var sampling rao.SamplingInfo
	err = json.Unmarshal([]byte(planInfo.SamplingInfo), &sampling)
	if err != nil {
		log.Logger.Error("logic.plan.GetPlanDetail.jsonUnmarshal(SamplingInfo), err：", err)
		return
	}
	emptyParam := []*rao.ParamsForm{
		{
			Enable: true,
			Key:    "",
			Value:  "",
			Desc:   "",
		},
	}
	var globalVariable []*rao.ParamsForm
	if planInfo.GlobalVariable == NULL {
		globalVariable = emptyParam
	} else {
		err = json.Unmarshal([]byte(planInfo.GlobalVariable), &globalVariable)
		if err != nil {
			log.Logger.Error("logic.plan.GetPlanDetail.jsonUnmarshal(GlobalVariable), err：", err)
			return
		}
	}

	var defaultHeader []*rao.ParamsForm
	if planInfo.DefaultHeader == NULL {
		defaultHeader = emptyParam
	} else {
		err = json.Unmarshal([]byte(planInfo.DefaultHeader), &defaultHeader)
		if err != nil {
			log.Logger.Error("logic.plan.GetPlanDetail.jsonUnmarshal(DefaultHeader), err：", err)
			return
		}
	}
	var serverInfo [][]string
	if planInfo.ServerInfo == NULL {
		serverInfo = make([][]string, 0)
	} else {
		_ = json.Unmarshal([]byte(planInfo.ServerInfo), &serverInfo)
	}

	var tagList []int32
	if planInfo.Tag == NULL || planInfo.Tag == "" {
		tagList = make([]int32, 0)
	} else {
		_ = json.Unmarshal([]byte(planInfo.Tag), &tagList)
	}

	reportList, err := GetReportListByPlanID(ctx, planId)
	if err != nil {
		return nil, err
	}

	users, err := GetUserList()
	if err != nil {
		log.Logger.Error("logic.plan.GetPlanDetail.GetUserList(), err:", err)
		return nil, errors.New("获取用户信息失败")
	}

	lastPressTime := ""
	if len(reportList) > 0 {
		lastPressTime = reportList[0].StartTime
	}
	plan = &rao.Plan{
		PlanID:         planInfo.ID,
		PlanName:       planInfo.PlanName,
		PressInfo:      pressInfo,
		SamplingInfo:   sampling,
		GlobalVariable: globalVariable,
		DefaultHeader:  defaultHeader,
		BreakType:      planInfo.BreakType,
		BreakValue:     planInfo.BreakValue,
		EngineCount:    planInfo.EngineCount,
		CreateUserName: GetNameByID(users, planInfo.CreateUserID),
		CreateTime:     planInfo.CreateTime.Format(FullTimeFormat),
		UpdateUserName: GetNameByID(users, planInfo.UpdateUserID),
		UpdateTime:     planInfo.UpdateTime.Format(FullTimeFormat),
		TagList:        tagList,
		Remark:         planInfo.Remark,
		PressCount:     len(reportList),
		LastPressTime:  lastPressTime,
		DebugSuccess:   planInfo.DebugStatus,
		TaskInfo:       taskInfo,
		ServerInfo:     serverInfo,
	}

	return plan, nil
}

func PlanDebug(ctx *gin.Context, planId int32) (sceneRes []rao.SceneExecutionResult, err error) {
	planInfo, err := GetPlanDetail(ctx, planId)
	scenesInfo, err := GetPlanScenesCase(ctx, planId)
	if err != nil {
		return nil, err
	}
	fileVariableList, err := GetPlanDebugFileVariable(ctx, planId)
	if err != nil {
		log.Logger.Error("GetPlanDebugFileVariable.FindAll err:", err)
	}
	planInfo.InitGlobalVariable()
	planInfo.AppendVariablePool(fileVariableList)

	sceneRes = make([]rao.SceneExecutionResult, 0)
	debugSuccess := true

	// 前置场景
	normalSceneIndex := 0
	firstScene := scenesInfo[0]
	if firstScene.SceneType == rao.PreFixScene {
		normalSceneIndex = 1
		if !firstScene.Disabled && len(firstScene.CaseTree) != 0 {
			variableList, _ := getSceneCommonVariableList(ctx, firstScene.SceneId)
			scene := rao.Scene{
				SceneID:        firstScene.SceneId,
				SceneType:      firstScene.SceneType,
				SceneName:      firstScene.SceneName,
				ExportDataInfo: *firstScene.ExportDataInfo,
				Cases:          firstScene.CaseTree,
				DefaultHeader:  planInfo.DefaultHeader,
				VariablePool: rao.VariablePool{
					VariableList: variableList,
				},
			}
			scene.LoadPlanVariablePool(planInfo)
			scene.RunScene()
			sceneRes = append(sceneRes, scene.Result)
			debugSuccess = debugSuccess && scene.Result.Passed
			planInfo.AppendVariablePool(scene.ExportDataInfo.ExportList)
		}
	}

	// 普通场景
	ch := make(chan rao.SceneExecutionResult, len(scenesInfo[normalSceneIndex:]))
	defer close(ch)
	for _, sc := range scenesInfo[normalSceneIndex:] {
		if sc.Disabled {
			ch <- rao.SceneExecutionResult{}
			continue
		}
		if len(sc.CaseTree) == 0 {
			ch <- rao.SceneExecutionResult{}
			continue
		}
		s := sc
		go func(sc rao.SceneInfo) {
			variableList, _ := getSceneCommonVariableList(ctx, sc.SceneId)
			scene := rao.Scene{
				SceneID:       sc.SceneId,
				SceneType:     sc.SceneType,
				SceneName:     sc.SceneName,
				Cases:         sc.CaseTree,
				DefaultHeader: planInfo.DefaultHeader,
				VariablePool: rao.VariablePool{
					VariableList: variableList,
				},
			}
			scene.LoadPlanVariablePool(planInfo)
			scene.RunScene()
			ch <- scene.Result
		}(s)
	}

	// 消费channel
	for i := 0; i < len(scenesInfo[normalSceneIndex:]); i++ {
		res := <-ch
		if res.SceneID == 0 {
			continue
		}
		sceneRes = append(sceneRes, res)
		debugSuccess = debugSuccess && res.Passed
	}
	sort.Slice(sceneRes[normalSceneIndex:], func(i, j int) bool {
		return sceneRes[i].SceneID < sceneRes[j].SceneID
	})

	go SetPlanDebugRecord(ctx, planId, sceneRes, debugSuccess)
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourcePlan,
		SourceID:      planId,
		OperationType: rao.DebugOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    struct{}{},
	})
	return sceneRes, nil
}

var planStartLock sync.Mutex

func PlanExecute(ctx *gin.Context, planId int32) (reportId int32, err error) {
	planStartLock.Lock()
	defer planStartLock.Unlock()
	// 获取执行锁
	isRunning, _, err := getPlanRunningLock(ctx, planId)
	if isRunning {
		return 0, fmt.Errorf("计划正在执行中")
	}
	// 获取发布计划详情
	planInfo, err := GetPlanDetail(ctx, planId)
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.GetPlanDetail(), err:", err)
		return 0, err
	}
	err = PlanInfoCheck(ctx, planInfo)
	if err != nil {
		return 0, err
	}
	scenePressInfo, err := getPressInfo(ctx, planInfo)
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.getPressInfo(), err:", err)
		return 0, err
	}
	if len(scenePressInfo) == 0 {
		log.Logger.Error("logic.plan.PlanExecute.getPressInfo(): 有效场景数为0，不执行")
		return 0, errors.New("有效场景数为0，不执行")
	}

	// 获取参数化文件
	fileInfo, err := GetPlanDataSource(ctx, planId)

	// 获取可用的kafka分区
	partitionId, err := getAvailablePartition(ctx)
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.getAvailablePartitions(), err:", err)
		return 0, err
	}

	// 获取可用机器
	availableEngineList, err := GetAvailableEngineList(ctx)
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.GetAvailableEngineList(), err:", err)
		return 0, err
	}
	if len(availableEngineList) < int(planInfo.EngineCount) {
		return 0, fmt.Errorf("可用压测机不足，请稍后再试")
	}

	// 设置集合点信息
	err = SetRallyPointScript(ctx, planId, scenePressInfo)
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.SetRallyPointScript(), err:", err)
		return 0, err
	}

	sceneIds := make([]int32, 0)
	for _, s := range scenePressInfo {
		sceneIds = append(sceneIds, s.Scene.SceneID)
	}

	engineList := availableEngineList[:planInfo.EngineCount]
	_, err = useEngine(ctx, engineList)
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.useEngine(), err:", err)
		return
	}

	// 生成测试报告id
	reportId, err = createReport(ctx, planInfo, engineList)
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.createReport(), err:", err)
		return 0, err
	}
	if reportId == 0 {
		return 0, fmt.Errorf("生成测试报告id失败")
	}

	// RPS模式设置缓存
	if planInfo.PressInfo.PressType == rao.PlanRpsRateMode {
		dal.ReportRdb.Set(ctx, fmt.Sprintf(RpsValKey, reportId), planInfo.PressInfo.RPS, time.Duration(planInfo.PressInfo.Duration+10)*time.Minute)
	}

	// 测试计划锁
	planRunningAddLock(ctx, planId, reportId, sceneIds)
	// 分区锁, TODO: 处理失败，资源释放
	_, err = usePartition(ctx, []int32{partitionId})
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.usePartition(), err:", err)
		return
	}

	// 下发给不同机器
	req := &rao.Action{
		Plan:               planInfo,
		ReportID:           reportId,
		ScenePressInfoList: scenePressInfo,
		FileInfo:           fileInfo,
		EngineCount:        planInfo.EngineCount,
		PartitionID:        partitionId,
	}

	// 延迟2s，为了确保消费者先启动
	time.Sleep(time.Duration(2) * time.Second)
	for index, engineIp := range engineList {
		req.EngineSerialNumber = int32(index)
		requestJson, _ := json.Marshal(req)
		engineUri := fmt.Sprintf("http://%s:8002/plan/run", engineIp)
		log.Logger.Info("请求压力机运行情况，report_id:", reportId, " 压测机器url为：", engineUri, " 请求参数为：\n", string(requestJson))
		runResponse, err := resty.New().R().SetBody(req).Post(engineUri)
		log.Logger.Info("请求压力机返回结果，report_id:", reportId, " 压测机器url为：", engineUri, " 返回body为：\n", string(runResponse.Body()))
		if err != nil {
			log.Logger.Error("logic.plan.PlanExecute.RequestEngine(), err:", err)
			return 0, err
		}
	}

	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourcePlan,
		SourceID:      planId,
		OperationType: rao.RunOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    struct{}{},
	})

	go func() {
		notRelease := true
		for notRelease {
			time.Sleep(5 * time.Second)
			availablePartitionIds, err := GetAvailablePartition(ctx)
			if err != nil {
				log.Logger.Error("logic.plan.PlanExecute.GetAvailablePartition(), err:", err)
				// TODO: 强制停止压测，并释放资源
			}
			if planInfo.PressInfo.PressType == rao.PlanRpsRateMode {
				DoRpsModeOperation(ctx, reportId, *req)
			}
			// TODO：后续做熔断、扩缩容等聚合到同一函数里
			//GetReportBottleneck(ctx, reportId)

			log.Logger.Infof("报告: %d 执行中，持续检查释放状态ing", reportId)
			// 判断已用partition是否被释放，被释放证明压测已执行结束
			notRelease = !slice.Contain(availablePartitionIds, partitionId)
			if !notRelease {
				log.Logger.Info("压测机器资源已释放，报告:", reportId)
			}
		}

		err = ReportPressDown(ctx, reportId)
		if err != nil {
			log.Logger.Error("logic.plan.PlanExecute.ReportPressDown(), err:", err)
		}
		if planInfo.CapacitySwitch {
			//err = CalcCapacityInfo(ctx, reportId)
		}
		if err != nil {
			log.Logger.Error("logic.plan.PlanExecute.CalcCapacityInfo(), err:", err)
		}
		log.Logger.Info("压测结束，report_id:", reportId)
	}()

	return reportId, nil
}

// 检查压测信息是否合规
func checkPressEngineInfo(ctx *gin.Context, planInfo *rao.Plan) (err error) {
	minConcurrency := int64(999999)
	switch planInfo.PressInfo.PressType {
	case rao.ConcurrentModel:
		for _, s := range planInfo.PressInfo.SceneList {
			minConcurrency = mathutil.Min(minConcurrency, s.Concurrency)
		}
	case rao.CaseRpsMode:

	case rao.PlanRpsRateMode:
	}
	return
}

// 将场景信息与场景施压信息进行组装
func getPressInfo(ctx *gin.Context, info *rao.Plan) (pressInfo []*rao.ScenePressInfo, err error) {
	sceneList := info.PressInfo.SceneList
	// 获取场景case信息
	pressInfo = make([]*rao.ScenePressInfo, 0)
	for _, scp := range sceneList {
		//totalRps := int64(0)
		// scene的基础信息，不含case
		baseScene, err := GetSceneDetail(ctx, scp.SceneId)
		if err != nil {
			log.Logger.Error("logic.plan.PlanExecute.GetSceneDetail(), err:", err)
			return nil, fmt.Errorf("获取发布计划详情失败")
		}
		if baseScene.Disabled {
			continue
		}
		// 组装caseTree
		tree, err := getAvailableCaseTree(ctx, scp.SceneId)
		if err != nil {
			log.Logger.Error("logic.plan.PlanExecute.GetSceneCaseTree(), err:", err)
			return nil, fmt.Errorf("获取发布计划详情失败")
		}
		if len(tree) == 0 {
			continue
		}
		baseScene.Cases = tree
		// 获取场景变量
		variableList, err := getSceneCommonVariableList(ctx, scp.SceneId)
		baseScene.VariablePool = rao.VariablePool{
			VariableList: variableList,
		}

		//switch info.PressType {
		//case rao.CaseRpsMode:
		//	{
		//		for _, caseInfo := range scp.CaseTree {
		//			totalRps += caseInfo.TargetRps
		//		}
		//		if totalRps < 5 {
		//			scp.Concurrency = 1
		//		} else {
		//			scp.Concurrency = totalRps / (1000 / RpsModelCaseDefaultRT)
		//		}
		//	}
		//case rao.PlanRpsRateMode:
		//	{
		//		totalRps = int64(scp.Rate) * info.RPS / 100
		//		if totalRps < 5 {
		//			scp.Concurrency = 1
		//		} else {
		//			// 默认为平均RT = 125ms
		//			scp.Concurrency = totalRps / (1000 / RpsModelCaseDefaultRT)
		//		}
		//	}
		//}
		if info.PressInfo.PressType == rao.PlanRpsRateMode {
			scp.Concurrency = mathutil.Max(int64(info.EngineCount), info.PressInfo.RPS/int64((1000/RpsModelCaseDefaultRT)*len(tree)*len(sceneList)))
		}

		pressInfo = append(pressInfo, &rao.ScenePressInfo{
			Scene:       baseScene,
			Rate:        scp.Rate,
			Concurrency: scp.Concurrency,
			Iteration:   scp.Iteration,
		})
	}
	// 前置场景排在最前
	sort.SliceStable(pressInfo, func(i, j int) bool {
		return pressInfo[i].Scene.SceneType > pressInfo[j].Scene.SceneType
	})
	return
}

func CheckRallyPoint(planInfo *rao.Plan, sceneList []*rao.ScenePressInfo) (err error) {
	for _, scene := range sceneList {
		for _, caseTree := range scene.Scene.Cases {
			if len(caseTree.Children) > 0 {
				for _, cs := range caseTree.Children {
					err = checkRallyPoint(planInfo, scene.Scene.SceneID, cs.SceneCase)
					if err != nil {
						return
					}
				}
			} else {
				err = checkRallyPoint(planInfo, scene.Scene.SceneID, caseTree.SceneCase)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func checkRallyPoint(planInfo *rao.Plan, sceneId int32, cs *rao.SceneCase) (err error) {
	if cs.Disabled {
		return
	}
	var httpCase rao.HttpCase
	httpCase.Unmarshal(cs.Extend)
	if httpCase.RallyPoint == nil {
		return
	}
	if !httpCase.RallyPoint.Enable {
		return
	}
	for _, scene := range planInfo.PressInfo.SceneList {
		if scene.SceneId == sceneId && scene.Concurrency < httpCase.RallyPoint.Concurrency {
			return fmt.Errorf("场景id: '%d'并发数为%d, 小于接口:'%s'集合点数%d", scene.SceneId, scene.Concurrency, httpCase.CaseName, httpCase.RallyPoint.Concurrency)
		}
	}
	return nil
}

func SetRallyPointScript(ctx context.Context, planId int32, scenePressInfoList []*rao.ScenePressInfo) (err error) {
	for _, scenePressInfo := range scenePressInfoList {
		for _, caseTree := range scenePressInfo.Scene.Cases {
			if len(caseTree.Children) > 0 {
				for _, cs := range caseTree.Children {
					err = setRallyPointScript(ctx, planId, cs.SceneCase)
					if err != nil {
						return
					}
				}
			} else {
				err = setRallyPointScript(ctx, planId, caseTree.SceneCase)
				if err != nil {
					return
				}
			}
		}
	}
	// TODO:缺失定时清理脚本机制
	return
}

func setRallyPointScript(ctx context.Context, planId int32, cs *rao.SceneCase) (err error) {
	if cs.Type != rao.HttpCaseType {
		return
	}
	var httpCase rao.HttpCase
	httpCase.Unmarshal(cs.Extend)
	if httpCase.RallyPoint == nil {
		return
	}
	if !httpCase.RallyPoint.Enable {
		return
	}
	luaScript := fmt.Sprintf(`
		local plan_id = %d
		local case_id = %d
		local concurrency = %d
		
		local key_count = tonumber(redis.call('get', 'report:count:' .. plan_id .. ':' .. case_id))
		local key_index = tonumber(redis.call('get', 'report:index:' .. plan_id .. ':' .. case_id .. ':' .. key_count))
		
		if key_index >= concurrency then
			local new_key_count = key_count + 1
			redis.call('SETEX', 'report:count:' .. plan_id .. ':' .. case_id, 7200, new_key_count)
			redis.call('SETEX', 'report:index:' .. plan_id .. ':' .. case_id .. ':' .. new_key_count, 7200, 0)
			return new_key_count
		else
			redis.call('INCRBY', 'report:index:' .. plan_id .. ':' .. case_id .. ':' .. key_count, 1)
			return key_count
		end
	`, planId, cs.CaseID, httpCase.RallyPoint.Concurrency)
	script := redis.NewScript(luaScript)
	_, err = dal.ReportRdb.SetNX(ctx, fmt.Sprintf("report:count:%d:%d", planId, cs.CaseID), 0, 2*time.Hour).Result()
	if err != nil {
		return err
	}
	_, err = dal.ReportRdb.SetNX(ctx, fmt.Sprintf("report:index:%d:%d:0", planId, cs.CaseID), 0, 2*time.Hour).Result()
	if err != nil {
		return err
	}
	sha, _ := script.Load(ctx, dal.ReportRdb).Result()
	httpCase.RallyPoint.LuaScriptSHA = sha
	cs.Extend = httpCase
	return
}

func DeleteRallyPointScript(ctx context.Context, planId int32, sceneList []*rao.ScenePressInfo) (err error) {
	for _, scene := range sceneList {
		for _, caseTree := range scene.Scene.Cases {
			if len(caseTree.Children) > 0 {
				for _, cs := range caseTree.Children {
					err = deleteRallyPointScript(ctx, planId, cs.SceneCase)
					if err != nil {
						return
					}
				}
			} else {
				err = deleteRallyPointScript(ctx, planId, caseTree.SceneCase)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func deleteRallyPointScript(ctx context.Context, planId int32, cs *rao.SceneCase) (err error) {
	if cs.Type != rao.HttpCaseType {
		return
	}
	var httpCase rao.HttpCase
	httpCase.Unmarshal(cs.Extend)
	if httpCase.RallyPoint == nil {
		return
	}
	if !httpCase.RallyPoint.Enable {
		return
	}
	countStr, err := dal.ReportRdb.Get(ctx, fmt.Sprintf("report:count:%d:%d", planId, cs.CaseID)).Result()
	countInt, err := strconv.Atoi(countStr)
	for i := 0; i < countInt+1; i++ {
		key := fmt.Sprintf("report:index:%d:%d:%d", planId, cs.CaseID, i)
		val, err := dal.ReportRdb.Del(ctx, key).Result()
		if err != nil || val != 1 {
			log.Logger.Error("logic.plan.deleteRallyPointScript.Del(index), val:", val, " err:", err)
			continue
		}
	}
	val, err := dal.ReportRdb.Del(ctx, fmt.Sprintf("report:count:%d:%d", planId, cs.CaseID)).Result()
	if err != nil || val != 1 {
		log.Logger.Error("logic.plan.deleteRallyPointScript.Del(index), val:", val, " err:", err)
	}
	return
}

func getAvailablePartition(ctx *gin.Context) (partitionId int32, err error) {
	partitionLock.Lock()
	defer partitionLock.Unlock()

	availablePartition, err := GetAvailablePartition(ctx)
	if err != nil {
		log.Logger.Error("logic.plan.PlanExecute.GetAvailablePartition(), err:", err)
		return -1, err
	}
	if availablePartition == nil || len(availablePartition) == 0 {
		return -1, errors.New("当前可用kafka分区不足")
	}
	// 需要使用的partition
	return availablePartition[0], nil
}

func GetPlanDebugVariable(ctx *gin.Context, planId int32) (res []rao.Variable, err error) {
	res = make([]rao.Variable, 0)
	fileVariable, err := GetPlanDebugFileVariable(ctx, planId)
	res = append(res, fileVariable...)
	detail, err := GetPlanDetail(ctx, planId)
	if err != nil {
		return nil, err
	}
	for _, vrb := range detail.GlobalVariable {
		if !vrb.Enable {
			continue
		}
		res = append(res, rao.Variable{
			VariableName: vrb.Key,
			VariableVal:  vrb.Value,
		})
	}
	return
}

type SimplePLan struct {
	PlanName string `json:"plan_name"`
	PlanID   int32  `json:"plan_id"`
}

func GetAllPlanList(ctx *gin.Context) (planList []SimplePLan, err error) {
	planList = make([]SimplePLan, 0)
	tx := dal.GetQuery().Plan
	plans, err := tx.WithContext(ctx).Where(tx.IsDelete.Not()).Order(tx.UpdateTime.Desc()).Find()
	if err != nil {
		log.Logger.Error("logic.plan.GetAllPlanList.planListQuery(),err:", err)
		return
	}
	for _, p := range plans {
		planList = append(planList, SimplePLan{
			PlanName: p.PlanName,
			PlanID:   p.ID,
		})
	}
	return
}

func SetPlanDebugRecord(ctx *gin.Context, planID int32, result []rao.SceneExecutionResult, success bool) {
	go SetPlanDebugStatus(ctx, planID, success)

	resultList, err := convertor.ToJson(result)
	if err != nil {
		log.Logger.Error("logic.scene.SceneDebug.convertorToJSon(), err:", err)
		return
	}
	tx := dal.GetQuery().DebugRecord
	insertData := &model.DebugRecord{
		PlanID:       planID,
		Status:       success,
		ResultInfo:   resultList,
		CreateUserID: utils.GetCurrentUserID(ctx),
	}
	_ = tx.WithContext(ctx).Create(insertData)
}

func SetPlanDebugStatus(ctx *gin.Context, planID int32, success bool) (isUpdate bool, err error) {
	tx := dal.GetQuery().Plan
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(planID)).UpdateSimple(tx.DebugStatus.Value(success))
	if err != nil {
		log.Logger.Error("logic.plan.setPlanDebugStatus.UpdateSimple，err:", err)
		return false, err
	}
	return true, nil
}

func GetDebugRecordList(ctx *gin.Context, req rao.PlanDebugRecordReq) (res rao.PageResponse[rao.PlanDebugRecord], err error) {
	tx := dal.GetQuery().DebugRecord
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.IsDelete.Not())
	conditions = append(conditions, tx.PlanID.Eq(req.PlanID))
	recordList, total, err := tx.WithContext(ctx).Where(conditions...).Order(tx.ID.Desc()).FindByPage(offset, limit)
	if err != nil {
		log.Logger.Error("getDebugRecordList.FindAll err:", err)
		return
	}
	list := make([]rao.PlanDebugRecord, 0)
	for _, v := range recordList {
		var resultInfo []rao.SceneExecutionResult
		err = json.Unmarshal([]byte(v.ResultInfo), &resultInfo)
		if err != nil {
			log.Logger.Error("logic.plan.GetDebugRecordList().Unmarshal Err", v.ResultInfo, "err:", err)
			continue
		}
		list = append(list, rao.PlanDebugRecord{
			Time:   v.CreateTime.Format(FullTimeFormat),
			Passed: v.Status,
			Result: resultInfo,
		})
	}
	return rao.PageResponse[rao.PlanDebugRecord]{
		List:  list,
		Total: total,
	}, nil
}

func UpdatePlanTask(ctx *gin.Context, req *rao.Plan) (taskId int32, err error) {
	if req.PlanID == 0 {
		return
	}
	planTx := dal.GetQuery().Plan
	plan, _ := planTx.WithContext(ctx).Where(planTx.ID.Eq(req.PlanID)).Take()
	if plan.TaskID == 0 && !req.TaskInfo.Enable {
		// 没有task，且不开启
		return
	}

	// 通知task服务
	param := fmt.Sprintf(`{"enable":%t,"plan_id":%d,"cron":"%s"}`, req.TaskInfo.Enable, req.PlanID, req.TaskInfo.Cron)
	request, _ := http.NewRequest(http.MethodPost, conf.Conf.Url.Task+"/plan/update", strings.NewReader(param))
	request.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(request)

	// task数据更新
	if plan.TaskID == 0 && req.TaskInfo.Enable {
		// 创建task
		taskId, err = CreateTask(ctx, rao.TaskParam{
			Type: rao.TypePlanTask,
			Cron: req.TaskInfo.Cron,
			TaskInfo: rao.PlanTask{
				PlanID: req.PlanID,
			},
			Enable: req.TaskInfo.Enable,
			UserID: utils.GetCurrentUserID(ctx),
		})

		return
	}
	err = UpdateTask(ctx, rao.TaskParam{
		PlanID: req.PlanID,
		ID:     plan.TaskID,
		Cron:   req.TaskInfo.Cron,
		TaskInfo: rao.PlanTask{
			PlanID: req.PlanID,
		},
		Enable: req.TaskInfo.Enable,
		UserID: utils.GetCurrentUserID(ctx),
	})
	if err != nil {
		return 0, err
	}

	return plan.TaskID, nil
}
