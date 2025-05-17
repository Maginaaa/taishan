package api

import (
	"encoding/json"
	"engine/internal/biz/errno"
	"engine/internal/biz/log"
	"engine/internal/response"
	middlewares "engine/middleware"
	"engine/model"
	"engine/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func PlanRun(ctx *gin.Context) {
	var act model.Action
	if err := ctx.ShouldBindJSON(&act); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrActionParam, err.Error())
		return
	}
	requestJson, _ := json.Marshal(act)
	log.Logger.Info(fmt.Sprintf("机器ip:%s, 调试场景", middlewares.LocalIp), string(requestJson))
	go server.ExecutionPlan(ctx, act)
	response.Success(ctx)
	return
}
