package main

import (
	"context"
	"engine/config"
	"engine/initialize"
	"engine/internal/biz/log"
	"engine/middleware"
	"engine/model"
	"engine/server/heartbeat"
	"flag"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ginRouter *gin.Engine
)

func main() {
	initService()
}

var configFile string

func initService() {
	flag.StringVar(&configFile, "f", "./conf.yml", "use config file")
	flag.Parse()
	config.InitConfig(configFile)

	zap.S().Debug("初始化logger")
	log.InitLogger()

	// 获取本机地址
	heartbeat.InitLocalIp()

	model.InitRedisClient()
	model.MustInitMongo()
	model.MustInitOSS()
	model.InitKafkaProducer()

	//3. 初始化routers
	log.Logger.Debug(fmt.Sprintf("机器ip:%s, 初始化routers", middleware.LocalIp))
	ginRouter = initialize.Routers()

	pprof.Register(ginRouter)

	// 注册服务
	log.Logger.Debug(fmt.Sprintf("机器ip:%s, 注册服务", middleware.LocalIp))
	taishanService := &http.Server{
		Addr:           config.Conf.Http.Address,
		Handler:        ginRouter,
		MaxHeaderBytes: 1 << 20,
	}

	// 初始化全局函数
	//tools.InitPublicFunc()

	go func() {
		if err := taishanService.ListenAndServe(); err != nil {
			log.Logger.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 注册并发送心跳数据
	go func() {
		heartbeat.SendHeartBeatRedis()
	}()

	// 资源监控数据
	go func() {
		heartbeat.SendMachineResources()
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Logger.Info(fmt.Sprintf("机器ip:%s, 注销成功", middleware.LocalIp))
	heartbeat.Logout()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := taishanService.Shutdown(ctx); err != nil {
		log.Logger.Info(fmt.Sprintf("机器ip:%s, 注销成功", middleware.LocalIp))
	}

}
