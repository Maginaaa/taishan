package handler

import (
	"github.com/gin-gonic/gin"
	"scene/internal/biz/errno"
	"scene/internal/response"
	"scene/logic"
	"scene/rao"
)

func GetOperationLog(ctx *gin.Context) {
	var req rao.OperationLogReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	if req.SourceName == "" {
		response.ErrorWithMsg(ctx, errno.ErrParam, "日志请求类型错误")
		return
	}
	data, err := logic.GetOperationLog(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, data)
	return
}
