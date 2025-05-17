package handler

import (
	"github.com/gin-gonic/gin"
	"scene/internal/biz/errno"
	"scene/internal/response"
	"scene/rao"
	"sync"
)

func FunctionAssistantDebug(ctx *gin.Context) {
	type FunctionAssistantReq struct {
		Key string `json:"key"`
	}
	var req FunctionAssistantReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrSceneCaseParam, err.Error())
		return
	}
	res := rao.FunctionAssistant(req.Key, new(sync.Map))
	response.SuccessWithData(ctx, res)
	return
}
