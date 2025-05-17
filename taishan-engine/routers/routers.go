package routers

import (
	"engine/api"
	"engine/internal/biz/log"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	{
		r.Use(RecoverPanic())
		r.GET("/ping", api.Ping)
		r.POST("/plan/run", api.PlanRun)
		//Router.POST("/run_plan/", api.RunPlan)

		//Router.POST("/stop/", api.Stop)
		//Router.POST("/stop_scene/", api.StopScene)
	}
}

func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 将异常信息打印到日志
				log.Logger.Error("接口panic错误,err:", err)
				// 返回一个错误响应给客户端
				c.JSON(500, gin.H{
					"error": "服务内部错误",
				})
			}
		}()
		// 继续接受和处理新的请求
		c.Next()
	}
}
