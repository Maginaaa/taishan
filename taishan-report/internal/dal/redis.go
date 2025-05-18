package dal

import (
	"context"
	"github.com/go-redis/redis/v8"
	"report/conf"
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

	RDB = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.Redis.Address,
		Password: conf.Conf.Redis.Password,
		DB:       conf.Conf.Redis.DB,
	})

}

func GetRDB() *redis.Client {
	return RDB
}

func SubscribeMsg(ctx context.Context, topic string) (pubSub *redis.PubSub) {
	pubSub = RDB.Subscribe(ctx, topic)
	return
}
