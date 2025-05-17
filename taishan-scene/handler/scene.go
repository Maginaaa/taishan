package handler

import (
	"errors"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"scene/internal/biz/errno"
	"scene/internal/response"
	"scene/logic"
	"scene/rao"
)

// CreateScene 保存场景
func CreateScene(ctx *gin.Context) {
	var req rao.Scene
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneParam, err.Error())
		return
	}
	if req.SceneID != 0 {
		response.ErrorWithMsg(ctx, errno.ErrSceneID, errors.New("场景ID错误").Error())
		return
	}
	if req.SceneName == "" {
		response.ErrorWithMsg(ctx, errno.ErrSceneParam, errors.New("场景名称不能为空").Error())
		return
	}
	if req.PlanID == 0 {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, errors.New("计划ID错误").Error())
		return
	}

	running, err := logic.PlanIsRunning(ctx, req.PlanID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	sceneId, err := logic.CreateScene(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, sceneId)
	return
}

func UpdateScene(ctx *gin.Context) {
	var req rao.Scene
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneParam, err.Error())
		return
	}
	if req.SceneID == 0 {
		response.ErrorWithMsg(ctx, errno.ErrSceneID, "场景ID错误")
		return
	}
	if req.PlanID == 0 {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, "计划ID错误")
		return
	}

	running, err := logic.PlanIsRunning(ctx, req.PlanID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	success, err := logic.UpdateScene(ctx, &req)
	if err != nil || !success {
		response.ErrorWithMsg(ctx, errno.ErrSceneUpdateFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

// GetSceneList 获取场景列表
func GetSceneList(ctx *gin.Context) {
	planId, err := convertor.ToInt(ctx.Param("id"))

	res, err := logic.GetPlanScenesCase(ctx, int32(planId))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, res)
	return
}

// DeleteScene 删除指定场景
func DeleteScene(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneID, err.Error())
		return
	}

	running, err := logic.ScenePlanIsRunning(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneDeleteFailed, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	b, err := logic.DeleteScene(ctx, int32(id))
	if err != nil || !b {
		response.ErrorWithMsg(ctx, errno.ErrSceneDeleteFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

// GetSceneDetail 获取场景详情
func GetSceneDetail(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneID, err.Error())
		return
	}
	detail, err := logic.GetSceneDetail(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, detail)
	return
}

func SceneDebug(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneID, err.Error())
		return
	}
	report, err := logic.SceneDebug(ctx, int32(id))

	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneDebugFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, report)
	return
}

func CopyScene(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneID, err.Error())
		return
	}
	running, err := logic.ScenePlanIsRunning(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}
	_, err = logic.CopyScene(ctx, int32(id), 0)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func GetSceneVariableList(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneID, err.Error())
		return
	}
	list, err := logic.GetSceneVariableList(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, list)
	return
}

func CreateSceneVariable(ctx *gin.Context) {
	var req rao.SceneVariable
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneParam, err.Error())
		return
	}
	if id := logic.CreateSceneVariable(ctx, req); id != 0 {
		response.SuccessWithData(ctx, id)
	} else {
		response.ErrorWithMsg(ctx, errno.ErrSceneVariableCreate, "创建失败")
	}
	return
}

func UpdateSceneVariable(ctx *gin.Context) {
	var req rao.SceneVariable
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneParam, err.Error())
		return
	}
	if logic.UpdateSceneVariable(ctx, req) {
		response.Success(ctx)
	} else {
		response.ErrorWithMsg(ctx, errno.ErrSceneVariableUpdate, "更新失败")
	}
	return
}

func DeleteSceneVariable(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneVariableID, err.Error())
		return
	}
	logic.DeleteSceneVariable(ctx, int32(id))
	response.Success(ctx)
	return
}

func SceneDebugRecord(ctx *gin.Context) {
	//id, err := convertor.ToInt(ctx.Param("id"))
	//if err != nil {
	//	response.ErrorWithMsg(ctx, errno.ErrSceneID, err.Error())
	//	return
	//}
	//recordList, err := logic.GetDebugRecordList(ctx, int32(id))
	//
	//if err != nil {
	//	log.Logger.Error("SceneDebug.GetDebugRecordList err:", err)
	//	response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
	//	return
	//}
	//
	//response.SuccessWithData(ctx, recordList)
	//return
}

func SortScene(ctx *gin.Context) {
	var req []rao.Scene
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneParam, err.Error())
		return
	}
	if len(req) == 0 {
		response.ErrorWithMsg(ctx, errno.ErrSceneParam, "")
		return
	}
	running, err := logic.ScenePlanIsRunning(ctx, req[0].SceneID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanExecute, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}
	b, err := logic.UpdateSceneSort(ctx, req)
	if err != nil || !b {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseFailed, err.Error())
		return
	}
	response.Success(ctx)

}
