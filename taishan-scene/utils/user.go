package utils

import "github.com/gin-gonic/gin"

func GetCurrentUserID(ctx *gin.Context) int32 {
	value, exists := ctx.Get("userID")
	if exists {
		return value.(int32)
	} else {
		return 0
	}
}

func SetAdminUser(ctx *gin.Context) {
	ctx.Set("userID", int32(100001))
	ctx.Set("userName", "系统")
}

func SetNormalizedUser(ctx *gin.Context) {
	ctx.Set("userID", int32(100003))
	ctx.Set("userName", "常态化压测")
}
