package internal

import (
	"data/internal/biz/log"
	"data/internal/dal"
)

func InitProjects() {
	// 初始化各种中间件
	dal.MustInitInfluxDB()
	dal.MustInitMongo()
	dal.NewKafkaProducer()
	dal.MustInitOSS()
	dal.MustInitSLS()
	dal.MustInitMySQL()
	dal.MustInitRedis()
	// 初始化logger
	log.InitLogger()

}
