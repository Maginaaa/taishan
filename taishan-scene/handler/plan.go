package handler

import (
	"errors"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"scene/internal/biz/errno"
	"scene/internal/response"
	"scene/logic"
	"scene/rao"
	"strings"
)

func CreatePlan(ctx *gin.Context) {
	var req rao.Plan
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanParam, err.Error())
		return
	}
	if req.PlanID != 0 {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, errors.New("id错误").Error())
		return
	}
	planId, err := logic.CreatePlan(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, planId)
	return

}

// DeletePlan 删除指定计划
func DeletePlan(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}

	running, err := logic.PlanIsRunning(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanDeleteFailed, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	b, err := logic.DeletePlan(ctx, int32(id))
	if err != nil || !b {
		response.ErrorWithMsg(ctx, errno.ErrPlanDeleteFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

// CopyPlan 复制指定计划
func CopyPlan(ctx *gin.Context) {
	var req rao.Plan
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanParam, err.Error())
		return
	}
	if !(req.PlanID > 0) {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, errors.New("id错误").Error())
		return
	}
	planId, err := logic.CopyPlan(ctx, req)

	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, planId)
	return
}

// UpdatePlan 更新指定计划
func UpdatePlan(ctx *gin.Context) {
	var req rao.Plan
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanParam, err.Error())
		return
	}

	running, err := logic.PlanIsRunning(ctx, req.PlanID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanUpdate, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}
	b, err := logic.UpdatePlan(ctx, &req)
	if err != nil || !b {
		response.ErrorWithMsg(ctx, errno.ErrPlanUpdate, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func UpdatePlanBaseInfo(ctx *gin.Context) {
	var req rao.Plan
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanParam, err.Error())
		return
	}

	running, err := logic.PlanIsRunning(ctx, req.PlanID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanUpdate, err.Error())
		return
	}
	if running {
		response.ErrorWithMsg(ctx, errno.RunningPlanEnableChange, " ")
		return
	}

	_, err = logic.UpdatePlanBaseInfo(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanUpdate, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func GetPlanList(ctx *gin.Context) {
	var req rao.PlanListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}
	req.PlanInfo = strings.TrimSpace(req.PlanInfo)
	req.CaseInfo = strings.TrimSpace(req.CaseInfo)
	data, err := logic.GetPlanList(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanMysqlFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, data)
}

func GetPlanDetail(ctx *gin.Context) {

	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}
	detail, err := logic.GetPlanDetail(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanMysqlFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, detail)
	return
}

// PlanDebug 计划调试
func PlanDebug(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}
	res, err := logic.PlanDebug(ctx, int32(id))
	if len(res) == 0 {
		response.ErrorWithMsg(ctx, errno.ErrEmptySceneTestCase, "场景调试执行失败")
		return
	}
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanDebugFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, res)
	return
}

// PlanExecute 执行指定计划, 开始压测
func PlanExecute(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}

	reportId, err := logic.PlanExecute(ctx, int32(id))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanExecute, err.Error())
		return
	}
	response.SuccessWithData(ctx, reportId)
	return
}

func SetPlanDebugStatusFalse(ctx *gin.Context) {
	id, err := convertor.ToInt(ctx.Param("id"))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, err.Error())
		return
	}
	_, err = logic.SetPlanDebugStatus(ctx, int32(id), false)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanExecute, err.Error())
		return
	}
	response.Success(ctx)
	return
}

func GetAllPlanList(ctx *gin.Context) {
	res, err := logic.GetAllPlanList(ctx)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, res)
	return
}

func GetPlanDebugRecord(ctx *gin.Context) {

	var req rao.PlanDebugRecordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanParam, err.Error())
		return
	}
	if req.PlanID == 0 {
		response.ErrorWithMsg(ctx, errno.ErrPlanID, errors.New("id错误").Error())
		return
	}
	data, err := logic.GetDebugRecordList(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrPlanMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}
