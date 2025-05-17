package dal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"scene/internal/conf"
)

var (
	ReportRdb *redis.Client
	RDB       *redis.Client
)

func MustInitRedis() {

	ReportRdb = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.RedisReport.Address,
		Password: conf.Conf.RedisReport.Password,
		DB:       conf.Conf.RedisReport.DB,
	})
	_, err := ReportRdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("model.request_redis.InitRedisClient ReportRdb连接失败 ,err: ", err.Error())
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.Redis.Address,
		Password: conf.Conf.Redis.Password,
		DB:       conf.Conf.Redis.DB,
	})
	_, err = RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("model.request_redis.InitRedisClient RDB连接失败 ,err: ", err.Error())
	}
	fmt.Println("redis initialized")

}

func GetRDB() *redis.Client {
	return RDB
}

func SubscribeMsg(ctx context.Context, topic string) (pubSub *redis.PubSub) {
	pubSub = RDB.Subscribe(ctx, topic)
	return
}
