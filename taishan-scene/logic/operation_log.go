package logic

import (
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"reflect"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"scene/internal/model"
	"scene/rao"
	"scene/utils"
)

func OperationInsert(ctx *gin.Context, info rao.OperationLog) {
	go func(ctx *gin.Context, info rao.OperationLog) {
		beforeMs, _ := json.Marshal(info.ValueBefore)
		afterMs, _ := json.Marshal(info.ValueAfter)
		beforeMp := make(map[string]interface{})
		afterMp := make(map[string]interface{})
		_ = json.Unmarshal(beforeMs, &beforeMp)
		_ = json.Unmarshal(afterMs, &afterMp)
		// 比较差异
		diff := make(map[string]rao.OperationDetail[interface{}])
		for key, value := range afterMp {
			if slice.Contain([]string{"create_user_id", "create_time", "update_user_id", "update_time", "is_delete"}, key) {
				continue
			}
			if oldValue, _ := beforeMp[key]; !reflect.DeepEqual(value, oldValue) {
				diff[key] = rao.OperationDetail[interface{}]{
					Before: oldValue,
					After:  value,
				}
			}
		}
		valueDiff, _ := json.Marshal(diff)
		valueDiffStr := string(valueDiff)
		if info.OperationType == rao.UpdateOperation && valueDiffStr == "{}" {
			return
		}
		tx := dal.GetQuery().OperationLog
		insertData := &model.OperationLog{
			SourceName:    info.SourceName,
			SourceID:      info.SourceID,
			OperationType: info.OperationType,
			OperatorID:    utils.GetCurrentUserID(ctx),
			ValueBefore:   string(beforeMs),
			ValueAfter:    string(afterMs),
			ValueDiff:     valueDiffStr,
		}
		err := tx.WithContext(ctx).Create(insertData)
		if err != nil {
			log.Logger.Error("logic.common.OperationInput.Create(), err:", err)
		}
	}(ctx, info)

}

func GetOperationLog(ctx *gin.Context, req rao.OperationLogReq) (res rao.PageResponse[rao.OperationLogResp], err error) {
	tx := dal.GetQuery().OperationLog
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize

	list := make([]rao.OperationLogResp, 0)
	var result []*model.OperationLog
	total := int64(0)
	if req.SourceName == rao.SourcePlan && req.SourceID != 0 {
		// 计划内额外处理文件、场景、case
		fileTx := dal.GetQuery().ParameterFile
		fileList, _ := fileTx.WithContext(ctx).Where(fileTx.PlanID.Eq(req.SourceID)).Find()
		fileIds := make([]int32, 0)
		for _, f := range fileList {
			fileIds = append(fileIds, f.ID)
		}
		sceneTx := dal.GetQuery().Scene
		sceneList, _ := sceneTx.WithContext(ctx).Where(sceneTx.PlanID.Eq(req.SourceID)).Find()
		sceneIds := make([]int32, 0)
		for _, f := range sceneList {
			sceneIds = append(sceneIds, f.ID)
		}
		caseTx := dal.GetQuery().SceneCase
		caseList, _ := caseTx.WithContext(ctx).Where(caseTx.SceneID.In(sceneIds...)).Find()
		caseIds := make([]int32, 0)
		for _, f := range caseList {
			caseIds = append(caseIds, f.ID)
		}
		result, total, err = tx.WithContext(ctx).Where(tx.SourceName.Eq(rao.SourcePlan), tx.SourceID.Eq(req.SourceID)).
			Or(tx.SourceName.Eq(rao.SourceScene), tx.SourceID.In(sceneIds...)).
			Or(tx.SourceName.Eq(rao.SourceSceneCase), tx.SourceID.In(caseIds...)).
			Or(tx.SourceName.Eq(rao.SourceFile), tx.SourceID.In(fileIds...)).Order(tx.ID.Desc()).FindByPage(offset, limit)

	} else {
		result, total, err = tx.WithContext(ctx).Where(tx.SourceName.Eq(req.SourceName), tx.SourceID.Eq(req.SourceID)).Order(tx.ID.Desc()).FindByPage(offset, limit)
	}
	// 获取用户信息列表
	users, err := GetUserList()
	if err != nil {
		log.Logger.Error("logic.operation_log.GetOperationLog.GetUserList(), err:", err)
		return res, errors.New("获取用户信息失败")
	}
	for _, r := range result {
		df := parseOperationLog(r)
		list = append(list, rao.OperationLogResp{
			ID:            r.ID,
			SourceId:      r.SourceID,
			SourceName:    r.SourceName,
			OperationType: r.OperationType,
			OperatorName:  GetNameByID(users, r.OperatorID),
			ValueDiff:     df,
			CreateTime:    r.CreatedTime.Format(FullTimeFormat),
		})
	}

	res.Total = total
	res.List = list
	return
}

func parseOperationLog(lg *model.OperationLog) (mp map[string]rao.OperationDetail[any]) {

	if lg.OperationType == rao.DeleteOperation {
		mp = make(map[string]rao.OperationDetail[any])
		mp["id"] = rao.OperationDetail[any]{
			Before: lg.SourceID,
			After:  0,
		}
		return
	}
	if lg.ValueDiff == "{}" {
		return nil
	}
	json.Unmarshal([]byte(lg.ValueDiff), &mp)
	if lg.OperationType == rao.CreateOperation {
		retainNameAndID(&mp)
		return
	}
	switch lg.SourceName {
	case rao.SourcePlan, rao.SourceScene, rao.SourceFile:
		return replacePlanKey(mp)
	case rao.SourceSceneCase:
		return replaceCaseKey(mp)
	}
	return
}

func retainNameAndID(mp *map[string]rao.OperationDetail[any]) {
	if mp == nil || *mp == nil {
		return
	}
	for key := range *mp {
		if key != "id" && key != "name" {
			delete(*mp, key)
		}
	}
}

func replacePlanKey(mp map[string]rao.OperationDetail[any]) (newMap map[string]rao.OperationDetail[any]) {
	if mp == nil {
		return
	}
	newMap = make(map[string]rao.OperationDetail[any])
	for key, diffVal := range mp {
		if slice.Contain(nestKey, key) {
			beforeMp := make(map[string]interface{})
			afterMp := make(map[string]interface{})
			_ = json.Unmarshal([]byte(diffVal.Before.(string)), &beforeMp)
			_ = json.Unmarshal([]byte(diffVal.After.(string)), &afterMp)
			for childK, v := range afterMp {
				if slice.Contain([]string{"create_user_id", "create_time", "update_user_id", "update_time", "is_delete"}, key) {
					continue
				}
				if oldValue, _ := beforeMp[childK]; !reflect.DeepEqual(v, oldValue) {
					newMap[getMappingName(key, planKeyMapperMap)+"-"+getMappingName(childK, planKeyMapperMap)] = rao.OperationDetail[interface{}]{
						Before: oldValue,
						After:  v,
					}
				}
			}
		} else {
			newMap[getMappingName(key, planKeyMapperMap)] = diffVal
		}
	}
	return newMap

}

func replaceCaseKey(mp map[string]rao.OperationDetail[any]) (newMap map[string]rao.OperationDetail[any]) {
	if mp == nil {
		return
	}
	newMap = make(map[string]rao.OperationDetail[any])

	beforeMp := make(map[string]interface{})
	afterMp := make(map[string]interface{})
	if _, ok := mp["extend"]; !ok {
		return
	}
	_ = json.Unmarshal([]byte(mp["extend"].Before.(string)), &beforeMp)
	_ = json.Unmarshal([]byte(mp["extend"].After.(string)), &afterMp)
	for k, v := range afterMp {
		if oldValue, _ := beforeMp[k]; !reflect.DeepEqual(v, oldValue) {
			newMap[getMappingName(k, caseKeyMapperMap)] = rao.OperationDetail[interface{}]{
				Before: oldValue,
				After:  v,
			}
		}
	}
	return
}

func getMappingName(key string, mp map[string]string) string {
	if value, exists := mp[key]; exists {
		return value
	} else {
		return key
	}
}

var nestKey = []string{"press_info", "sampling_info", "scene_list"}

var planKeyMapperMap = map[string]string{
	"press_info":      "压测策略",
	"press_type":      "压测类型",
	"plan_name":       "计划名称",
	"scene_list":      "场景列表",
	"duration":        "时长",
	"task_cron":       "定时任务-表达式",
	"task_enable":     "定时任务-开关",
	"sampling_info":   "采样策略",
	"sampling_type":   "采样模式",
	"sampling_rate":   "采样率",
	"server_info":     "管理服务",
	"global_variable": "全局变量",
	"default_header":  "默认请求头",
	"scene_name":      "场景名",
	"column":          "文件参数",
	"export_info":     "数据导出",
}

var caseKeyMapperMap = map[string]string{
	"case_name":       "Case名称",
	"assert_form":     "断言",
	"body":            "请求体",
	"headers_form":    "请求头",
	"overtime_config": "超时配置",
	"params_form":     "请求url参数",
	"rally_point":     "集合点",
	"variable_form":   "变量",
	"waiting_config":  "等待时间",
	"method_type":     "请求方法",
	"url":             "请求URL",
}
