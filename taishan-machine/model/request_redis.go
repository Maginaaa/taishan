package model

import (
	"encoding/json"
	"github.com/duke-git/lancet/slice"
	"github.com/go-redis/redis"
	"machine/conf"
	"machine/log"
	"time"
)

var (
	RDB              *redis.Client
	MachineListKey         = "TaishanMachineList"
	MachineInfoKey         = "MachineMonitor:%s"
	MaxHeartbeatDiff int64 = 60 // 最大心跳时间差，单位为秒
)

const (
	UsedEngines = "UsedEngines"
)

func InitRedisClient() (err error) {
	RDB = redis.NewClient(
		&redis.Options{
			Addr:     conf.Conf.Redis.Address,
			Password: conf.Conf.Redis.Password,
			DB:       conf.Conf.Redis.DB,
		})
	_, err = RDB.Ping().Result()
	if err != nil {
		log.Logger.Fatal("model.request_redis.InitRedisClient RDB连接失败 ,err: ", err.Error())
		return
	}
	return err
}

func GetMachineMap() map[string]HeartBeat {
	val, err := RDB.HGetAll(MachineListKey).Result()
	if err != nil {
		log.Logger.Error("model.request_redis.GetMachineMap, err: ", err.Error())
		return nil
	}
	machineMap := make(map[string]HeartBeat)
	for k, v := range val {
		var hb HeartBeat
		err = json.Unmarshal([]byte(v), &hb)
		if err != nil {
			log.Logger.Error("model.request_redis.GetMachineMap, err: ", err.Error())
			return nil
		}
		diff := time.Now().Unix() - hb.CreateTime
		if diff > MaxHeartbeatDiff {
			continue
		}
		machineMap[k] = hb
	}
	return machineMap
}

func GetAvailableMachine() map[string]HeartBeat {
	val, err := RDB.HGetAll(MachineListKey).Result()
	if err != nil {
		log.Logger.Error("model.request_redis.GetMachineMap, err: ", err.Error())
		return nil
	}
	usingList, _ := RDB.SMembers(UsedEngines).Result()
	machineMap := make(map[string]HeartBeat)
	for k, v := range val {
		var hb HeartBeat
		err = json.Unmarshal([]byte(v), &hb)
		if !slice.Contain(usingList, hb.IP) {
			machineMap[k] = hb
		}
	}
	return machineMap
}
