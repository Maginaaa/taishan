package dal

import (
	"context"
	"data/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	RDB *redis.Client
)

func MustInitRedis() {

	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Address,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("model.request_redis.InitRedisClient ReportRdb连接失败 ,err: ", err.Error())
	}
	fmt.Println("redis initialized")

}
