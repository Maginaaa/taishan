package logic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"scene/internal/model"
	"scene/rao"
	"scene/utils"
)

func GetTaskDetail(ctx *gin.Context, taskId int32) (taskDetail *rao.TaskDetail[any], err error) {
	tx := dal.GetQuery().TaskInfo
	taskInfo, err := tx.WithContext(ctx).Where(tx.ID.Eq(taskId)).First()
	if taskInfo == nil || err != nil {
		log.Logger.Error("logic.GetTaskDetail.Query() error:", err)
		return nil, err
	}
	switch taskInfo.Type {
	case rao.TypePlanTask:
		var dt rao.PlanTask
		_ = json.Unmarshal([]byte(taskInfo.TaskInfo), &dt)
		taskDetail = &rao.TaskDetail[any]{
			TaskID:   taskInfo.ID,
			Cron:     taskInfo.Cron,
			TaskData: &dt,
			Enable:   taskInfo.Enable,
		}
	default:
		return nil, nil
	}
	return
}

func CreateTask(ctx *gin.Context, param rao.TaskParam) (taskId int32, err error) {
	tx := dal.GetQuery().TaskInfo
	ms, _ := json.Marshal(param.TaskInfo)
	insertData := &model.TaskInfo{
		Type:         param.Type,
		Cron:         param.Cron,
		TaskInfo:     string(ms),
		Enable:       param.Enable,
		CreateUserID: param.UserID,
	}
	err = tx.WithContext(ctx).Create(insertData)
	if err != nil {
		log.Logger.Error("logic.CreateTask.Create() error:", err)
		return
	}
	return insertData.ID, nil
}

func UpdateTask(ctx *gin.Context, param rao.TaskParam) (err error) {
	tx := dal.GetQuery().TaskInfo
	ms, _ := json.Marshal(param.TaskInfo)
	beforeTask, _ := tx.WithContext(ctx).Where(tx.ID.Eq(param.ID)).First()
	_, err = tx.WithContext(ctx).Where(tx.ID.Eq(param.ID)).UpdateSimple(tx.Cron.Value(param.Cron), tx.TaskInfo.Value(string(ms)), tx.Enable.Value(param.Enable), tx.Enable.Value(param.Enable), tx.UpdateUserID.Value(param.UserID))
	afterTask, _ := tx.WithContext(ctx).Where(tx.ID.Eq(param.ID)).First()
	OperationInsert(ctx, rao.OperationLog{
		SourceName:    rao.SourcePlan,
		SourceID:      param.PlanID,
		OperationType: rao.UpdateOperation,
		OperatorID:    utils.GetCurrentUserID(ctx),
		ValueBefore: map[string]any{
			"task_cron":   beforeTask.Cron,
			"task_enable": beforeTask.Enable,
		},
		ValueAfter: map[string]any{
			"task_cron":   afterTask.Cron,
			"task_enable": afterTask.Enable,
		},
	})
	return
}
