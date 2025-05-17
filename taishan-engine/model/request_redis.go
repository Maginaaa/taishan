package model

import (
	"engine/config"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

var (
	ReportRdb *redis.Client
	RDB       *redis.Client
)

func InitRedisClient() {

	RDB = redis.NewClient(
		&redis.Options{
			Addr:     config.Conf.Redis.Address,
			Password: config.Conf.Redis.Password,
			DB:       config.Conf.Redis.DB,
		})
	_, err := RDB.Ping().Result()
	if err != nil {
		log.Fatal("model.request_redis.InitRedisClient RDB连接失败", "err", err.Error())
	}

	ReportRdb = redis.NewClient(
		&redis.Options{
			Addr:     config.Conf.ReportRedis.Address,
			Password: config.Conf.ReportRedis.Password,
			DB:       config.Conf.ReportRedis.DB,
		})
	_, err = RDB.Ping().Result()
	if err != nil {
		log.Fatal("model.request_redis.InitRedisClient ReportRedis连接失败", "err", err.Error())
	}
	fmt.Println("redis initialized")
}

func InsertHeartbeat(key string, field string, value interface{}) error {
	_, err := RDB.HSet(key, field, value).Result()
	return err
}

func InsertMachineResources(key string, value interface{}) error {
	_, err := RDB.LPush(key, value).Result()
	return err
}

func DelMachine(key, field string) error {
	_, err := RDB.HDel(key, field).Result()
	return err
}

func DelResources(key string) error {
	_, err := RDB.Del(key).Result()
	return err
}

func SubscribeMsg(topic string) (pubSub *redis.PubSub) {
	pubSub = RDB.Subscribe(topic)
	return
}
