package handler

import (
	"github.com/gin-gonic/gin"
	"scene/internal/biz/errno"
	"scene/internal/response"
	"scene/logic"
	"scene/rao"
)

func GetSamplingData(ctx *gin.Context) {
	var req rao.SamplingDataReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	info, err := logic.GetReportDebugInfo(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, info)
	return
}

func ModifyReportRps(ctx *gin.Context) {
	var req rao.DoReportRpsModifyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	err := logic.DoReportRpsModify(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	response.Success(ctx)
	return
}
