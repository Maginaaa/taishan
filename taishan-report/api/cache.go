package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"report/internal/dal"
	"report/internal/errno"
	"report/internal/response"
	"report/rao"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func SetCache(ctx *gin.Context) {
	var req rao.Cache
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	dal.LocalCacheSet(req.Key, req.Value)
	return
}

func GetCache(ctx *gin.Context) {
	var req rao.Cache
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, err.Error())
		return
	}
	bt, b := dal.LocalCacheGet(req.Key)
	if !b {
		response.ErrorWithMsg(ctx, errno.ErrReportParam, "查询缓存失败")
		return
	}
	response.SuccessWithData(ctx, bt)
	return
}
