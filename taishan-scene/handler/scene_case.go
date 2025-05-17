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

// CreateSceneCase 保存场景内的case
func CreateSceneCase(ctx *gin.Context) {
	var req rao.SceneCase
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseParam, err.Error())
		return
	}
	if req.CaseID != 0 {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseID, errors.New("id错误").Error())
		return
	}

	running, err := logic.ScenePlanIsRunning(ctx, req.SceneID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanExecute, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	caseId, err := logic.CreateSceneCase(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, caseId)
	return
}

func ImportCase(ctx *gin.Context) {
	var req rao.ImportCase
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	caseId, err := logic.ImportCase(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, caseId)
	return
}

func UpdateSceneCase(ctx *gin.Context) {
	var req rao.SceneCase
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseParam, err.Error())
		return
	}
	if req.CaseID == 0 {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseID, errors.New("caseId不能为空").Error())
		return
	}

	running, err := logic.ScenePlanIsRunning(ctx, req.SceneID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanExecute, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	_, err = logic.UpdateSceneCase(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

// DeleteSceneCase 删除指定case
func DeleteSceneCase(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseID, err.Error())
		return
	}

	running, err := logic.CasePlanIsRunning(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanExecute, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	b, err := logic.DeleteSceneCase(ctx, int32(id))
	if err != nil || !b {
		response.ErrorWithMsg(ctx, errno.ErrSceneCreateFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

// SortSceneCase 对场景内的case进行排序
func SortSceneCase(ctx *gin.Context) {
	var req rao.CaseSortReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseParam, err.Error())
		return
	}
	running, err := logic.ScenePlanIsRunning(ctx, req.Before.SceneID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanExecute, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}
	b, err := logic.UpdateTreeSort(ctx, &req)
	if err != nil || !b {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func GetSceneCaseTree(ctx *gin.Context) {
	sceneId, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseID, err.Error())
		return
	}

	if sceneId == 0 {
		response.ErrorWithMsg(ctx, errno.ErrSceneID, errors.New("sceneId错误").Error())
		return
	}
	tree, err := logic.GetSceneCaseTree(ctx, int32(sceneId))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, tree)
	return
}

func CaseDebug(ctx *gin.Context) {
	var req rao.HttpCase
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseParam, err.Error())
		return
	}
	//err := urlAssert(ctx, req.URL)
	//if err != nil {
	//	response.ErrorWithMsg(ctx, errno.ErrSceneCaseParam, err.Error())
	//	return
	//}
	res, err := logic.HttpCaseDebug(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, res)
	return
}
