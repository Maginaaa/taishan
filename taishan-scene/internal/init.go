package internal

import (
	"scene/internal/biz/log"
	"scene/internal/conf"
	"scene/internal/dal"
	"scene/logic"
)

func InitProjects(configFile string) {
	conf.MustInitConf(configFile)

	// 初始化各种中间件
	dal.MustInitMySQL()
	dal.MustInitMongo()
	dal.MustInitRedis()
	dal.MustInitOSS()
	dal.InitFeishu()
	//dal.MustInitK8sClient()

	logic.PartitionInit()

	//dal.MustInitRedisForReport()
	//dal.MustInitBigCache()
	// 初始化logger
	log.InitLogger()
}
