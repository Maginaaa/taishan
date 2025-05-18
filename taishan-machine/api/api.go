package api

import (
	"github.com/gin-gonic/gin"
	"machine/model"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func GetAllMachine(ctx *gin.Context) {
	res := getAllMachine(ctx)
	SuccessWithData(ctx, res)
	return
}

func GetAvailableMachine(ctx *gin.Context) {
	res := getAvailableMachine(ctx)
	SuccessWithData(ctx, res)
	return
}

func GetAllMachineInfo(ctx *gin.Context) {
	//var req model.MachineInfoReq
	//if err := ctx.ShouldBindJSON(&req); err != nil {
	//	ErrorWithMsg(ctx, ErrParam, err.Error())
	//	return
	//}
	res := GetAllMachineInfoLogic(ctx)
	SuccessWithData(ctx, res)
	return
}

func GetMachineInfo(ctx *gin.Context) {
	var req model.MachineInfoReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ErrorWithMsg(ctx, ErrParam, err.Error())
		return
	}
	res := GetMachineInfoLogic(ctx, req.IPList)
	SuccessWithData(ctx, res)
	return
}
