package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"gorm.io/gen"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"scene/internal/model"
	"scene/rao"
	"scene/utils"
)

func CreateScene(ctx *gin.Context, req *rao.Scene) (sceneId int32, err error) {
	var exportInfoStr string
	scenes, err := getPlanScenes(ctx, req.PlanID)
	if len(scenes) >= 100 {
		return 0, errors.New("场景数量超过100上限")
	}
	// 获取scene的sort
	currentSceneSort := scenesMaxSort(scenes) + 1
	if req.SceneType == rao.PreFixScene {
		currentSceneSort = 0
		if checkHasPreFixScene(ctx, 0, scenes) {
			return 0, errors.New("已存在前置场景")
		}
	}
	if &req.ExportDataInfo == nil {
		exportInfoStr = "{}"
	} else {
		exportInfoByte, _ := json.Marshal(req.ExportDataInfo)
		exportInfoStr = string(exportInfoByte)
	}

	tx := dal.GetQuery().Scene
	insertData := &model.Scene{
		PlanID:       req.PlanID,
		SceneType:    req.SceneType,
		ExportInfo:   exportInfoStr,
		SceneName:    req.SceneName,
		Sort:         currentSceneSort,
		Disabled:     false,
		CreateUserID: utils.GetCurrentUserID(ctx),
	}
	err = tx.WithContext(ctx).Create(insertData)
	if err != nil {
		return 0, err
	}
	sceneId = insertData.ID
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceScene,
		SourceID:      sceneId,
		OperationType: rao.CreateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    insertData,
	})
	return
}

func UpdateScene(ctx *gin.Context, req *rao.Scene) (success bool, err error) {
	scenes, err := getPlanScenes(ctx, req.PlanID)
	currentSceneSort := req.Sort
	if req.SceneType == rao.PreFixScene {
		currentSceneSort = 0
		if checkHasPreFixScene(ctx, req.SceneID, scenes) {
			return false, errors.New("仅允许创建1个前置场景")
		}
	}
	var beforeScene *model.Scene
	for _, s := range scenes {
		if s.ID == req.SceneID {
			beforeScene = s
		}
	}
	exportInfoStr := "{}"
	if req.SceneType == rao.PreFixScene {
		exportInfoStr, err = convertor.ToJson(req.ExportDataInfo)
		if err != nil {
			log.Logger.Error("logic.scene.UpdateScene.convertorToJson ，err:", err)
			return false, errors.New("导出信息获取失败")
		}
	}
	tx := dal.GetQuery().Scene
	updateData := &model.Scene{
		ID:           req.SceneID,
		PlanID:       req.PlanID,
		SceneType:    req.SceneType,
		SceneName:    req.SceneName,
		ExportInfo:   exportInfoStr,
		Disabled:     req.Disabled,
		Sort:         currentSceneSort,
		UpdateUserID: utils.GetCurrentUserID(ctx),
	}
	err = tx.WithContext(ctx).Save(updateData)
	if err != nil {
		log.Logger.Error("logic.scene.UpdateScene.Save ，err:", err)
		return false, errors.New("数据库更新失败")
	}

	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceScene,
		SourceID:      req.SceneID,
		OperationType: rao.UpdateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   beforeScene,
		ValueAfter:    updateData,
	})
	return true, nil

}

func GetPlanScenesCase(ctx *gin.Context, planId int32) (sceneList []rao.SceneInfo, err error) {
	scenes, err := getPlanScenes(ctx, planId)
	sceneIds := make([]int32, 0)
	for _, s := range scenes {
		sceneIds = append(sceneIds, s.ID)
	}
	caseMap, err := BatchGetSceneCaseTree(ctx, sceneIds)
	if err != nil {
		log.Logger.Error("logic.plan.GetPlanDetail.BatchGetSceneCaseTree(), err：", err)
		return nil, err
	}
	vrbMap, err := batchGetSceneCommonVariableList(ctx, sceneIds)
	if err != nil {
		log.Logger.Error("logic.plan.GetPlanDetail.batchGetSceneCommonVariableList(), err：", err)
		return nil, err
	}
	sceneList = make([]rao.SceneInfo, 0)
	for _, v := range scenes {
		scene := rao.SceneInfo{
			SceneId:      v.ID,
			Sort:         v.Sort,
			SceneType:    v.SceneType,
			SceneName:    v.SceneName,
			Disabled:     v.Disabled,
			CaseTree:     caseMap[v.ID],
			VariableList: vrbMap[v.ID],
		}
		if v.SceneType == rao.PreFixScene {
			var exportInfo rao.ExportDataInfo
			_ = json.Unmarshal([]byte(v.ExportInfo), &exportInfo)
			scene.ExportDataInfo = &exportInfo
			bucket, _ := dal.GetTaishanBucket()
			isExist, _ := bucket.IsObjectExist(fmt.Sprintf("export/scene/%d.csv", scene.SceneId))
			scene.ExportDataInfo.HasCache = isExist
			if !isExist {
				scene.ExportDataInfo.DisableCache = true
			}
		}
		sceneList = append(sceneList, scene)
	}
	return
}

func getPlanScenes(ctx *gin.Context, planId int32) (sceneList []*model.Scene, err error) {
	tx := dal.GetQuery().Scene
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.IsDelete.Not())
	conditions = append(conditions, tx.PlanID.Eq(planId))
	sceneList, err = tx.WithContext(ctx).Where(conditions...).Order(tx.SceneType.Desc()).Order(tx.Sort).Order(tx.ID).Find()
	if err != nil {
		log.Logger.Error("logic.scene.getPlanSceneIds.Find(), err:", err)
		return
	}
	return
}

func GetSceneDetail(ctx *gin.Context, id int32) (*rao.Scene, error) {
	tx := dal.GetQuery().Scene
	s, err := tx.WithContext(ctx).Where(tx.ID.Eq(id)).First()
	if err != nil {
		log.Logger.Error("场景列表--获取场景详情数据失败，err:", err)
		return nil, err
	}

	// TODO: 获取创建人、修改人信息
	sc := &rao.Scene{
		PlanID:       s.PlanID,
		SceneID:      s.ID,
		Sort:         s.Sort,
		SceneType:    s.SceneType,
		SceneName:    s.SceneName,
		Disabled:     s.Disabled,
		CreateUserID: s.CreateUserID,
		UpdateUserID: s.UpdateUserID,
		UpdateTime:   s.UpdateTime.Format(FullTimeFormat),
	}
	if sc.SceneType == rao.PreFixScene && s.ExportInfo != "" {
		var exportInfo rao.ExportDataInfo
		err = json.Unmarshal([]byte(s.ExportInfo), &exportInfo)
		if err != nil {
			return nil, err
		}
		sc.ExportDataInfo = exportInfo
	}

	return sc, nil
}

func DeleteScene(ctx *gin.Context, id int32) (isDelete bool, err error) {
	// TODO: 场景删除时，需要提前软删除场景内所有case
	tx := dal.GetQuery().Scene
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(id)).UpdateSimple(tx.IsDelete.Value(true))
	if err != nil {
		log.Logger.Error("logic.scene.DeleteScene.UpdateSimple，err:", err)
		return false, err
	}
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceScene,
		SourceID:      id,
		OperationType: rao.DeleteOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    struct{}{},
	})
	return true, nil
}

func CopyScene(ctx *gin.Context, id int32, planId int32) (sceneId int32, err error) {
	//根据场景id获取场景信息
	detail, err := GetSceneDetail(ctx, id)
	detail.SceneName = detail.SceneName + "-COPY"
	detail.SceneID = 0
	if planId != 0 {
		detail.PlanID = planId
	}
	//copy场景
	copySceneId, err := CreateScene(ctx, detail)
	if err != nil {
		return
	}
	//获取场景下的case tree
	tree, err := getSceneCaseTree(ctx, id)
	copyCase(ctx, tree, copySceneId, 0)

	return copySceneId, nil
}

// 复制case tree到新的场景id下
func copyCase(ctx *gin.Context, tree []*rao.SceneCaseTree, copySceneId, parentId int32) (err error) {
	for _, t := range tree {
		t.CaseID = 0
		t.SceneID = copySceneId
		tx := dal.GetQuery().SceneCase
		extStr, _ := convertor.ToJson(t.Extend)
		insertData := &model.SceneCase{
			ParentID: parentId,
			Type:     t.Type,
			SceneID:  t.SceneID,
			Disabled: false,
			Sort:     t.Sort,
			Extend:   extStr,
		}
		err = tx.WithContext(ctx).Create(insertData)
		if err != nil {
			log.Logger.Error("logic.scene_case.UpdateTreeSort.createCase(), err:", err)
			return
		}
		if len(t.Children) > 0 {
			caseId := insertData.ID
			err = copyCase(ctx, t.Children, copySceneId, caseId)
			if err != nil {
				return
			}
		}
	}
	return
}

// SceneDebug 查询和执行
func SceneDebug(ctx *gin.Context, sceneId int32) (report rao.SceneExecutionResult, err error) {
	baseScene, err := GetSceneDetail(ctx, sceneId)
	if err != nil {
		log.Logger.Error("SceneDebug.GetSceneDetail err:", err)
		return
	}
	var caseTree []*rao.SceneCaseTree
	if baseScene.Disabled {
		return report, fmt.Errorf("已禁用场景无法调试")
	}
	caseTree, err = getAvailableCaseTree(ctx, sceneId)
	if err != nil {
		log.Logger.Error("SceneDebug.getSceneCaseTree err:", err)
		return
	}

	list, err := getSceneCommonVariableList(ctx, sceneId)
	if err != nil {
		return rao.SceneExecutionResult{}, err
	}
	// 从参数化文件中获取数据
	fileVariableList, err := GetPlanDebugFileVariable(ctx, baseScene.PlanID)
	if err != nil {
		log.Logger.Error("GetPlanDebugFileVariable.FindAll err:", err)
	}
	mergeVariableList := append(fileVariableList, list...)

	scene := rao.Scene{
		SceneID:      baseScene.SceneID,
		SceneName:    baseScene.SceneName,
		CreateUserID: baseScene.CreateUserID,
		UpdateUserID: baseScene.UpdateUserID,
		Cases:        caseTree,
		VariablePool: rao.VariablePool{
			VariableList: mergeVariableList,
		},
	}
	scene.RunScene()
	return scene.Result, nil

}

func GetSceneVariableList(ctx *gin.Context, sceneId int32) (res []*rao.SceneVariable, err error) {

	tx := dal.GetQuery().SceneVariable
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.IsDelete.Not())
	conditions = append(conditions, tx.SceneID.Eq(sceneId))
	variableList, err := tx.WithContext(ctx).Where(conditions...).Find()
	if err != nil {
		log.Logger.Error("GetSceneVariableList.FindAll err:", err)
		return
	}
	res = make([]*rao.SceneVariable, 0)
	for _, v := range variableList {
		res = append(res, &rao.SceneVariable{
			SceneID: sceneId,
			Variable: rao.Variable{
				VariableID:   v.ID,
				VariableName: v.VariableName,
				VariableVal:  v.VariableVal,
				Remark:       v.Remark,
			},
		})
	}
	return
}

func getSceneCommonVariableList(ctx *gin.Context, sceneId int32) (res []rao.Variable, err error) {
	tx := dal.GetQuery().SceneVariable
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.IsDelete.Not())
	conditions = append(conditions, tx.SceneID.Eq(sceneId))
	variableList, err := tx.WithContext(ctx).Where(conditions...).Find()
	if err != nil {
		log.Logger.Error("logic.scene.getSceneCommonVariableList.FindAll err:", err)
		return
	}
	res = make([]rao.Variable, 0)
	for _, v := range variableList {
		res = append(res, rao.Variable{
			VariableID:   v.ID,
			VariableName: v.VariableName,
			VariableVal:  v.VariableVal,
			Remark:       v.Remark,
		})
	}
	return
}

func batchGetSceneCommonVariableList(ctx *gin.Context, sceneIds []int32) (mp map[int32][]*rao.Variable, err error) {
	tx := dal.GetQuery().SceneVariable
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.IsDelete.Not())
	conditions = append(conditions, tx.SceneID.In(sceneIds...))
	variableList, err := tx.WithContext(ctx).Where(conditions...).Find()
	if err != nil {
		log.Logger.Error("logic.scene.batchGetSceneCommonVariableList.FindAll err:", err)
		return
	}
	mp = make(map[int32][]*rao.Variable)
	for _, id := range sceneIds {
		mp[id] = make([]*rao.Variable, 0)
	}
	for _, v := range variableList {
		mp[v.SceneID] = append(mp[v.SceneID], &rao.Variable{
			VariableID:   v.ID,
			VariableName: v.VariableName,
			VariableVal:  v.VariableVal,
			Remark:       v.Remark,
		})
	}
	return
}

func CreateSceneVariable(ctx *gin.Context, variable rao.SceneVariable) (variableId int32) {
	tx := dal.GetQuery().SceneVariable
	insertData := &model.SceneVariable{
		SceneID:      variable.SceneID,
		VariableName: variable.VariableName,
		VariableVal:  variable.VariableVal,
		Remark:       variable.Remark,
		CreateUserID: utils.GetCurrentUserID(ctx),
	}
	err := tx.WithContext(ctx).Create(insertData)
	if err != nil {
		log.Logger.Error("logic.scene_case.CreateSceneVariable.Create err:", err)
		return 0
	}
	return insertData.ID
}

func UpdateSceneVariable(ctx *gin.Context, variable rao.SceneVariable) (success bool) {
	tx := dal.GetQuery().SceneVariable
	info, err := tx.WithContext(ctx).Where(tx.ID.Eq(variable.VariableID)).UpdateSimple(tx.VariableName.Value(variable.VariableName),
		tx.VariableVal.Value(variable.VariableVal), tx.Remark.Value(variable.Remark))
	if err != nil {
		log.Logger.Error("logic.scene_case.UpdateSceneVariable.UpdateSimple err:", err)
		return false
	}
	return info.RowsAffected > 0
}

func DeleteSceneVariable(ctx *gin.Context, variableId int32) (success bool) {
	tx := dal.GetQuery().SceneVariable
	info, err := tx.WithContext(ctx).Where(tx.ID.Eq(variableId)).UpdateSimple(tx.IsDelete.Value(true))
	if err != nil {
		log.Logger.Error("logic.scene_case.DeleteSceneVariable.UpdateSimple err:", err)
		return false
	}
	return info.RowsAffected > 0
}

// ScenePlanIsRunning 查询场景所在计划是否正在执行
func ScenePlanIsRunning(ctx *gin.Context, sceneId int32) (isRunning bool, err error) {
	planId, err := getScenePlanID(ctx, sceneId)
	isRunning, err = PlanIsRunning(ctx, planId)
	if err != nil {
		return false, err
	}
	return
}

func getScenePlanID(ctx *gin.Context, sceneId int32) (planId int32, err error) {
	tx := dal.GetQuery().Scene
	scene, err := tx.WithContext(ctx).Where(tx.ID.Eq(sceneId)).First()
	if err != nil {
		log.Logger.Errorf("logic.scene.getScenePlanID.queryScene: sceneId: %d, err: %v", sceneId, err)
		return 0, err
	}
	return scene.PlanID, nil
}

func checkHasPreFixScene(ctx *gin.Context, currentSceneID int32, sceneList []*model.Scene) bool {
	for _, scene := range sceneList {
		if scene.SceneType == rao.PreFixScene && scene.ID != currentSceneID {
			return true
		}
	}
	return false
}

// 获取当前所有scene的最大sort值
func scenesMaxSort(sceneList []*model.Scene) (sort int32) {
	sort = 0
	for _, scene := range sceneList {
		if scene.Sort > sort {
			sort = scene.Sort
		}
	}
	return
}

func UpdateSceneSort(ctx *gin.Context, req []rao.Scene) (success bool, err error) {
	tx := dal.GetQuery().Scene
	for _, item := range req {
		_, _ = tx.WithContext(ctx).Where(tx.ID.Eq(item.SceneID)).UpdateSimple(tx.Sort.Value(item.Sort))
	}
	return true, nil
}
