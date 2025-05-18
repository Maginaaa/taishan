package rao

import (
	"github.com/go-redis/redis"
	"report/conf"
	"report/log"
)

var (
	ReportRdb *redis.Client
	RDB       *redis.Client
)

func InitRedisClient() (err error) {
	ReportRdb = redis.NewClient(
		&redis.Options{
			Addr:     conf.Conf.Redis.Address,
			Password: conf.Conf.Redis.Password,
			DB:       conf.Conf.Redis.DB,
		})
	_, err = ReportRdb.Ping().Result()
	if err != nil {
		log.Logger.Error("rao.request_redis.InitRedisClient ReportRdb连接失败", "err", err.Error())
		return
	}

	RDB = redis.NewClient(
		&redis.Options{
			Addr:     conf.Conf.RedisReport.Address,
			Password: conf.Conf.RedisReport.Password,
			DB:       conf.Conf.RedisReport.DB,
		})
	_, err = RDB.Ping().Result()
	if err != nil {
		log.Logger.Error("rao.request_redis.InitRedisClient RDB连接失败", "err", err.Error())
	}
	return err
}

func SubscribeMsg(topic string) (pubSub *redis.PubSub) {
	pubSub = RDB.Subscribe(topic)
	return
}
