package handler

import (
	"data/internal/response"
	"data/logic"
	"github.com/gin-gonic/gin"
)

func GetDashboardData(ctx *gin.Context) {
	data := logic.GetDashboardData(ctx)
	response.SuccessWithData(ctx, data)
	return
}
