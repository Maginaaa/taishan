package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"scene/internal/biz/errno"
	"scene/internal/response"
	"scene/logic"
	"scene/rao"
)

func GetTagList(ctx *gin.Context) {
	var req rao.TagSearchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	if req.Type == 0 {
		response.ErrorWithMsg(ctx, errno.ErrParam, errors.New("tag类型错误").Error())
		return
	}

	data, err := logic.GetTagLiat(ctx, req.Type)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, data)
	return
}
