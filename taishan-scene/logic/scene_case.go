package logic

import (
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"gorm.io/gen"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"scene/internal/model"
	"scene/rao"
	"scene/utils"
)

func CreateSceneCase(ctx *gin.Context, req *rao.SceneCase) (caseId int32, err error) {
	tx := dal.GetQuery().SceneCase
	// 先查询当前数据，获取sort
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.ParentID.Eq(req.ParentID))
	conditions = append(conditions, tx.SceneID.Eq(req.SceneID))
	caseInfo, err := tx.WithContext(ctx).Where(conditions...).Order(tx.Sort.Desc()).FirstOrInit()
	if err != nil {
		log.Logger.Error("logic.scene_case.CreateSceneCase.queryCaseInfo(), err:", err)
	}
	extStr, err := convertor.ToJson(req.Extend)
	if err != nil {
		log.Logger.Error("logic.scene_case.CreateSceneCase.convertorToJSon(), err:", err)
		return 0, err
	}
	insertData := &model.SceneCase{
		ParentID: req.ParentID,
		Type:     req.Type,
		SceneID:  req.SceneID,
		Disabled: req.Disabled,
		Sort:     caseInfo.Sort + 1,
		Extend:   extStr,
	}

	err = tx.WithContext(ctx).Create(insertData)
	if err != nil {
		log.Logger.Error("logic.scene_case.UpdateTreeSort.createCase(), err:", err)
		return 0, err
	}
	caseId = insertData.ID
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceSceneCase,
		SourceID:      caseId,
		OperationType: rao.CreateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    insertData,
	})
	return
}

func ImportCase(ctx *gin.Context, req rao.ImportCase) (caseId int32, err error) {
	tx := dal.GetQuery().SceneCase
	// 查询baseCase
	baseCase, _ := tx.WithContext(ctx).Where(tx.ID.Eq(req.CaseID)).First()
	var extMp map[string]interface{}
	err = json.Unmarshal([]byte(baseCase.Extend), &extMp)
	if err != nil {
		log.Logger.Error("logic.scene_case.ImportCase.Unmarshal(), err:", err)
		return
	}
	newCase := &rao.SceneCase{
		Type:    baseCase.Type,
		SceneID: req.SceneID,
		Extend:  extMp,
	}
	caseId, err = CreateSceneCase(ctx, newCase)
	return
}

func UpdateSceneCase(ctx *gin.Context, req *rao.SceneCase) (isUpdate bool, err error) {
	tx := dal.GetQuery().SceneCase
	extStr, err := convertor.ToJson(req.Extend)
	if err != nil {
		log.Logger.Error("logic.scene_case.UpdateSceneCaseExt.convertorToJson(), err:", err)
		return false, err
	}
	beforeInfo, _ := tx.WithContext(ctx).Where(tx.ID.Eq(req.CaseID)).First()
	cs, err := tx.WithContext(ctx).Where(tx.ID.Eq(req.CaseID)).UpdateSimple(tx.Disabled.Value(req.Disabled), tx.Extend.Value(extStr))
	if err != nil {
		log.Logger.Error("logic.scene_case.UpdateSceneCaseExt.Update(), err:", err)
		return false, err
	}
	afterInfo, _ := tx.WithContext(ctx).Where(tx.ID.Eq(req.CaseID)).First()
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceSceneCase,
		SourceID:      req.CaseID,
		OperationType: rao.UpdateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   beforeInfo,
		ValueAfter:    afterInfo,
	})
	return cs.RowsAffected > 0, nil
}

func DeleteSceneCase(ctx *gin.Context, id int32) (isDelete bool, err error) {
	tx := dal.GetQuery().SceneCase
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(id)).UpdateSimple(tx.IsDelete.Value(true))
	if err != nil {
		log.Logger.Error("logic.SceneCase.DeleteCase.UpdateSimple，err:", err)
		return false, err
	}
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceSceneCase,
		SourceID:      id,
		OperationType: rao.DeleteOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    struct{}{},
	})
	return true, nil
}

func HttpCaseDebug(ctx *gin.Context, httpCase *rao.HttpCase) (hdr *rao.HttpResponse, err error) {
	// 使用场景变量
	sceneId := httpCase.SceneID
	if sceneId == 0 {
		log.Logger.Error("logic.sceneCase.HttpCaseDebug.sceneId is 0")
		return nil, errors.New("场景ID缺失")
	}
	sceneVariableList, err := GetSceneVariableList(ctx, sceneId)
	if err != nil {
		return nil, err
	}
	// 获取计划id
	planID, err := getScenePlanID(ctx, sceneId)
	if err != nil {
		log.Logger.Error("logic.sceneCase.HttpCaseDebug.getScenePlanID: err", err)
		return nil, err
	}
	// 从参数化文件中获取数据
	planVariableList, err := GetPlanDebugVariable(ctx, planID)
	if err != nil {
		log.Logger.Error("logic.sceneCase.HttpCaseDebug.GetPlanDebugVariable err:", err)
	}
	sceneVariableListNew := make([]rao.Variable, 0)
	for _, item := range sceneVariableList {
		sceneVariableListNew = append(sceneVariableListNew, rao.Variable{
			VariableName: item.VariableName,
			VariableVal:  item.VariableVal,
		})
	}
	planInfo, err := GetPlanDetail(ctx, planID)
	if err != nil {
		return nil, err
	}
	mergeVariableList := append(planVariableList, sceneVariableListNew...)
	httpCase.LoadVariablePool(rao.VariablePool{VariableList: mergeVariableList})
	httpCase.LoadDefaultHeader(planInfo.DefaultHeader)
	err = httpCase.DoRequest()
	if err != nil {
		return
	}
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourceSceneCase,
		SourceID:      httpCase.CaseID,
		OperationType: rao.DebugOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore:   struct{}{},
		ValueAfter:    struct{}{},
	})
	return httpCase.ResponseData, nil
}

func GetSceneCaseTree(ctx *gin.Context, sceneId int32) (tree []*rao.SceneCaseTree, err error) {
	return getSceneCaseTree(ctx, sceneId)
}

func BatchGetSceneCaseTree(ctx *gin.Context, sceneIds []int32) (mp map[int32][]*rao.SceneCaseTree, err error) {
	return batchGetSceneCaseTree(ctx, sceneIds)
}

func GetAvailableCaseTree(ctx *gin.Context, sceneId int32) (tree []*rao.SceneCaseTree, err error) {
	return getAvailableCaseTree(ctx, sceneId)
}

func UpdateTreeSort(ctx *gin.Context, req *rao.CaseSortReq) (isUpdate bool, err error) {
	//查询出after.sort之后的数据
	tx := dal.GetQuery().SceneCase
	if req.Position != rao.Inner {
		var dropSort int32
		conditions := make([]gen.Condition, 0)
		conditions = append(conditions, tx.SceneID.Eq(req.After.SceneID))
		conditions = append(conditions, tx.ParentID.Eq(req.After.ParentID))
		if req.Position == rao.Before {
			conditions = append(conditions, tx.Sort.Gte(req.After.Sort))
			dropSort = req.After.Sort
		} else if req.Position == rao.After {
			conditions = append(conditions, tx.Sort.Gt(req.After.Sort))
			dropSort = req.After.Sort + 1
		}
		caseInfo, err := tx.WithContext(ctx).Where(conditions...).Or(tx.ID.Eq(req.Before.CaseID)).Order(tx.Sort).Find()
		if err != nil {
			log.Logger.Error("logic.scene_case.UpdateTreeSort.queryCaseInfo(), err:", err)
		}
		for _, item := range caseInfo {
			if item.ID == req.Before.CaseID {
				item.Sort = dropSort
			} else {
				item.Sort += 1
			}
			item.ParentID = req.After.ParentID
			err = tx.WithContext(ctx).Save(item)
			if err != nil {
				return false, err
			}
		}
	} else if req.Position == rao.Inner {
		conditions := make([]gen.Condition, 0)
		conditions = append(conditions, tx.ParentID.Eq(req.After.CaseID))
		innerItem, err := tx.WithContext(ctx).Where(conditions...).Order(tx.Sort.Desc()).FirstOrInit()
		if err != nil {
			log.Logger.Error("Case拖拽[innerItem]--case查询错误，err:", err)
		}
		conditions = make([]gen.Condition, 0)
		conditions = append(conditions, tx.ID.Eq(req.Before.CaseID))
		moveItem, err := tx.WithContext(ctx).Where(conditions...).Order(tx.Sort.Desc()).FirstOrInit()
		if err != nil {
			log.Logger.Error("Case拖拽[inner][moveItem]--case查询错误，err:", err)
		}
		moveItem.Sort = innerItem.Sort + 1
		moveItem.ParentID = req.After.CaseID
		err = tx.WithContext(ctx).Save(moveItem)
		if err != nil {
			return false, err
		}
	}
	return true, nil

}

func getSceneCaseList(ctx *gin.Context, sceneId int32) (list []*model.SceneCase, err error) {
	tx := dal.GetQuery().SceneCase
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.SceneID.Eq(sceneId))
	conditions = append(conditions, tx.IsDelete.Not())
	list, err = tx.WithContext(ctx).Where(conditions...).Order(tx.Sort).Find()
	return
}

func batchGetSceneCaseList(ctx *gin.Context, sceneIds []int32) (mp map[int32][]*model.SceneCase, err error) {
	tx := dal.GetQuery().SceneCase
	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.SceneID.In(sceneIds...))
	conditions = append(conditions, tx.IsDelete.Not())
	list, err := tx.WithContext(ctx).Where(conditions...).Order(tx.Sort).Find()
	mp = make(map[int32][]*model.SceneCase)
	for _, id := range sceneIds {
		mp[id] = make([]*model.SceneCase, 0)
	}
	for _, l := range list {
		mp[l.SceneID] = append(mp[l.SceneID], l)
	}
	return
}

func getSceneCaseTree(ctx *gin.Context, sceneId int32) (tree []*rao.SceneCaseTree, err error) {
	sceneCaseList, err := getSceneCaseList(ctx, sceneId)
	if err != nil {
		log.Logger.Error("场景Case -- 获取场景Case列表数据失败，err:", err)
		return
	}
	if len(sceneCaseList) == 0 {
		tree = make([]*rao.SceneCaseTree, 0)
		return
	}
	tree = BuildTree(sceneCaseList, 0)
	return tree, nil
}

func batchGetSceneCaseTree(ctx *gin.Context, sceneIds []int32) (mp map[int32][]*rao.SceneCaseTree, err error) {
	sceneCaseMp, err := batchGetSceneCaseList(ctx, sceneIds)
	if err != nil {
		log.Logger.Error("logic.scene_case.batchGetSceneCaseTree.batchGetSceneCaseList: err", err)
		return
	}
	mp = make(map[int32][]*rao.SceneCaseTree)
	for sceneId, sceneCaseList := range sceneCaseMp {
		tree := make([]*rao.SceneCaseTree, 0)
		if len(sceneCaseList) != 0 {
			tree = BuildTree(sceneCaseList, 0)
		}
		mp[sceneId] = tree
	}
	return
}

func getAvailableCaseTree(ctx *gin.Context, sceneId int32) (tree []*rao.SceneCaseTree, err error) {
	sceneCaseList, err := getSceneCaseList(ctx, sceneId)
	if err != nil {
		log.Logger.Error("场景Case -- 获取场景Case列表数据失败，err:", err)
		return
	}
	if len(sceneCaseList) == 0 {
		tree = make([]*rao.SceneCaseTree, 0)
		return
	}
	baseTree := BuildTree(sceneCaseList, 0)
	for _, t := range baseTree {
		if !t.Disabled {
			switch t.Type {
			case rao.HttpCaseType:
				tree = append(tree, t)
			case rao.LogicControlType:
				newChildren := make([]*rao.SceneCaseTree, 0)
				for _, c := range t.Children {
					if !c.Disabled {
						newChildren = append(newChildren, c)
					}
				}
				t.Children = newChildren
				tree = append(tree, t)
			}

		}
	}
	return tree, nil
}

func BuildTree(list []*model.SceneCase, parentId int32) (result []*rao.SceneCaseTree) {
	result = make([]*rao.SceneCaseTree, 0)
	for _, item := range list {
		var extMp map[string]interface{}
		err := json.Unmarshal([]byte(item.Extend), &extMp)
		if err != nil {
			log.Logger.Error("logic.scene_case.treeSort(), item.Extend:", item.Extend, "err:", err)
			continue
		}
		// 判断extMp内是否存在case_name
		title := ""
		if _, ok := extMp["case_name"]; ok {
			title = extMp["case_name"].(string)
		}
		if item.ParentID == parentId {
			var treeItem = &rao.SceneCaseTree{
				SceneCase: &rao.SceneCase{
					CaseID:   item.ID,
					Title:    title,
					ParentID: item.ParentID,
					Type:     item.Type,
					SceneID:  item.SceneID,
					Sort:     item.Sort,
					Disabled: item.Disabled,
					Extend:   extMp,
				},
			}
			treeItem.Children = BuildTree(list, treeItem.CaseID)
			result = append(result, treeItem)
		}
	}
	return result
}

func CasePlanIsRunning(ctx *gin.Context, caseId int32) (isRunning bool, err error) {
	tx := dal.GetQuery().SceneCase
	_case, err := tx.WithContext(ctx).Where(tx.ID.Eq(caseId)).First()
	if err != nil {
		log.Logger.Error("logic.sceneCase.CasePlanIsRunning.queryCase: err", err)
		return false, err
	}
	planID, err := getScenePlanID(ctx, _case.SceneID)
	if err != nil {
		log.Logger.Error("logic.sceneCase.CasePlanIsRunning.getScenePlanID: err", err)
		return false, err
	}
	return PlanIsRunning(ctx, planID)
}
