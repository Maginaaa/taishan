package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"scene/handler"
	"scene/internal/biz/log"
	"scene/internal/middleware"
	"time"
)

func RegisterRouter(r *gin.Engine) {
	// cors
	//r.Use(cors.Default())
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.Use(RecoverPanic()) // 恢复因接口内部错误导致的panic

	// 探活接口
	r.Use(RecoverPanic()) // 恢复因接口内部错误导致的panic
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Use(middleware.SessionAuthMiddleWare())

	{
		scene := r.Group("/scene")
		scene.GET("/list/:id", handler.GetSceneList)
		scene.POST("/create", handler.CreateScene)
		scene.POST("/update", handler.UpdateScene)
		scene.GET("/delete/:id", handler.DeleteScene)
		//scene.GET("/detail/:id", handler.GetSceneDetail)
		scene.GET("/debug/:id", handler.SceneDebug)
		scene.GET("/copy/:id", handler.CopyScene)
		scene.GET("/variable/list/:id", handler.GetSceneVariableList)
		scene.POST("/variable/create", handler.CreateSceneVariable)
		scene.POST("/variable/update", handler.UpdateSceneVariable)
		scene.GET("/variable/delete/:id", handler.DeleteSceneVariable)
		scene.GET("/debug_record/:id", handler.SceneDebugRecord)
		scene.POST("/sort", handler.SortScene)

		cs := r.Group("/case")
		cs.POST("/create", handler.CreateSceneCase)
		cs.POST("/import", handler.ImportCase)
		cs.POST("/update", handler.UpdateSceneCase)
		cs.GET("/delete/:id", handler.DeleteSceneCase)
		cs.POST("/sort", handler.SortSceneCase)
		cs.GET("/tree/:id", handler.GetSceneCaseTree)
		cs.POST("/debug", handler.CaseDebug)

		pl := r.Group("/plan")
		pl.POST("/create", handler.CreatePlan)
		pl.GET("/delete/:id", handler.DeletePlan)
		pl.POST("/update", handler.UpdatePlan)
		pl.POST("/update/simple", handler.UpdatePlanBaseInfo)
		pl.POST("/list", handler.GetPlanList)
		pl.GET("/list/all", handler.GetAllPlanList)
		pl.GET("/detail/:id", handler.GetPlanDetail)
		pl.GET("/debug/:id", handler.PlanDebug)
		pl.GET("/execute/:id", handler.PlanExecute)
		pl.GET("/debug/false/:id", handler.SetPlanDebugStatusFalse)
		pl.POST("/copy", handler.CopyPlan)
		pl.POST("/record", handler.GetPlanDebugRecord)

		fi := r.Group("/file")
		fi.POST("/upload", handler.FileUpload)
		fi.GET("/list/:id", handler.GetPlanDataSource)
		fi.GET("/delete/:id", handler.DeleteFile)
		fi.POST("/column/update", handler.ColumnUpdate)
		fi.POST("/download", handler.DownloadPlanFile)

		rp := r.Group("/report")
		rp.POST("/sampling", handler.GetSamplingData)
		rp.POST("/rps/modify", handler.ModifyReportRps)

		ut := r.Group("/utils")
		ut.POST("/function_assistant/debug", handler.FunctionAssistantDebug)

		tag := r.Group("/tag")
		tag.POST("/list", handler.GetTagList)

		r.POST("/log", handler.GetOperationLog)
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
