package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"report/api"
	"report/log"
	"runtime/debug"
)

func RegisterRouter(r *gin.Engine) {
	// cors
	r.Use(cors.Default())
	r.Use(RecoverPanic()) // 恢复因接口内部错误导致的panic

	{
		r.GET("/ping", api.Ping)

		rp := r.Group("/report")
		rp.POST("/create", api.CreateReport)
		rp.POST("/list", api.GetReportList)
		rp.GET("/detail/:id", api.GetReportDetail)
		rp.GET("/data/:id", api.GetReportData)
		rp.POST("/data/case", api.GetReportCaseData)
		rp.GET("/stop/:id", api.StopPressTest)
		rp.GET("/stop/:id/hard", api.StopPressTestHard)
		rp.POST("/update/currency", api.UpdatePressTest)
		rp.POST("/update/release", api.ReleasePreScene)
		rp.POST("/update/name", api.UpdateReportName)
		r.GET("/rps/target/:id", api.GetReportTargetRps)

		ch := r.Group("/cache")
		ch.POST("/set", api.SetCache)
		ch.POST("/get", api.GetCache)
	}
}

func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 将异常信息打印到日志
				log.Logger.Errorf("接口panic错误,err: %v, stack: %s", err, debug.Stack())
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
