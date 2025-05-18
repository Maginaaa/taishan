package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"machine/api"
	"machine/log"
	"runtime"
)

func RegisterRouter(r *gin.Engine) {
	// cors
	r.Use(cors.Default())

	r.Use(RecoverPanic()) // 恢复因接口内部错误导致的panic

	{
		r.GET("/ping", api.Ping)
		r.GET("/machine/list", api.GetAllMachine)
		r.GET("/machine/available/list", api.GetAvailableMachine)
		r.POST("/machine/info/all", api.GetAllMachineInfo)
		r.POST("/machine/info/list", api.GetMachineInfo)
	}
}

func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				pc, filename, line, _ := runtime.Caller(3)
				f := runtime.FuncForPC(pc)
				// 将异常信息打印到日志
				log.Logger.Error("接口panic错误,err:", err, " func:", f.Name(), " file:", filename, " line:", line)
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
