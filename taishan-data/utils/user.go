package utils

import (
	"github.com/gin-gonic/gin"
)

func GetCurrentUserID(ctx *gin.Context) int32 {
	return ctx.MustGet("userID").(int32)
}
