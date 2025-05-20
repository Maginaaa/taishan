package model

import (
	"collector/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	ReportRdb    *redis.Client
	RDB          *redis.Client
	timeDuration = 3 * time.Second

	UsedEngines = "UsedEngines"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func InitRedisClient() {
	ReportRdb = redis.NewClient(
		&redis.Options{
			Addr:     config.Conf.ReportRedis.Address,
			Password: config.Conf.ReportRedis.Password,
			DB:       config.Conf.ReportRedis.DB,
		})
	_, err := ReportRdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("model.request_redis.InitRedisClient ReportRdb连接失败", "err", err.Error())
	}

	RDB = redis.NewClient(
		&redis.Options{
			Addr:     config.Conf.Redis.Address,
			Password: config.Conf.Redis.Password,
			DB:       config.Conf.Redis.DB,
		})
	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatal("model.request_redis.InitRedisClient RDB连接失败", "err", err.Error())
	}
	fmt.Println("redis initialized")
}

func SubscribeMsg(topic string) (pubSub *redis.PubSub) {
	pubSub = RDB.Subscribe(ctx, topic)
	return
}

//func InsertTestData(sceneTestResultDataMsg *SceneTestResultDataMsg) (err error) {
//	data := sceneTestResultDataMsg.ToJson()
//	key := fmt.Sprintf("reportData:%d:%d", sceneTestResultDataMsg.ReportID, sceneTestResultDataMsg.SceneID)
//	if sceneTestResultDataMsg.End {
//		//duration := sceneTestResultDataMsg.TimeStamp - runTime
//		//pkg.SendStopStressReport(machineMap, sceneTestResultDataMsg.TeamId, sceneTestResultDataMsg.PlanId, sceneTestResultDataMsg.ReportId, duration)
//	}
//	err = ReportRdb.LPush(context.TODO(), key, data).Err()
//	return
//}

func ReleaseEngine(ip string) (err error) {
	err = RDB.SRem(context.TODO(), UsedEngines, ip).Err()
	return
}
