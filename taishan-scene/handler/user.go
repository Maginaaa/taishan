package handler

import (
	"github.com/gin-gonic/gin"
	"scene/internal/biz/errno"
	"scene/internal/response"
	"scene/logic"
	"scene/rao"
)

// GetUserList 获取用户列表
func GetUserList(ctx *gin.Context) {

	res, err := logic.GetUserList()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrGetUser, err.Error())
		return
	}
	response.SuccessWithData(ctx, rao.Response{
		Users: res,
	})
	return
}

// UserLogout 用户登出
func UserLogout(ctx *gin.Context) {
	err := logic.UserLogout(ctx)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrLogout, err.Error())
		return
	}
	response.SuccessWithData(ctx, map[string]interface{}{})
	return
}
