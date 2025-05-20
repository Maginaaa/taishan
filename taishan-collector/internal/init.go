package internal

import (
	"collector/internal/biz/log"
	"collector/internal/dal"
	"collector/model"
	"go.uber.org/zap"
)

func InitProjects() {
	// 初始化各种中间件
	dal.MustInitInfluxDB()
	model.InitRedisClient()
	// 初始化logger
	zap.S().Debug("初始化logger")
	log.InitLogger()
}
