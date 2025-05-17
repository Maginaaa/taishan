package initialize

import (
	"engine/middleware"
	"engine/routers"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {

	r := gin.Default()
	// 配置跨域
	r.Use(middleware.Cors())

	routers.InitRouter(r)

	return r
}
