package task

import (
	"encoding/json"
	"fmt"
	"machine/log"
	"machine/model"
	"time"
)

const (
	MachineListKey         = "TaishanMachineList"
	MachineInfoKey         = "MachineMonitor:%s"
	MaxSaveHour            = 2
	MaxHeartbeatDiff int64 = 60 // 最大心跳时间差，单位为秒

)

func DropMachine() (err error) {
	val, err := model.RDB.HGetAll(MachineListKey).Result()
	if err != nil {
		log.Logger.Error("model.request_redis.GetMachineMap, err: ", err.Error())
		return nil
	}
	for k, v := range val {
		var hb model.HeartBeat
		err = json.Unmarshal([]byte(v), &hb)
		if err != nil {
			log.Logger.Error("model.request_redis.GetMachineMap, err: ", err.Error())
			return nil
		}
		diff := time.Now().Unix() - hb.CreateTime
		if diff > MaxHeartbeatDiff {
			// 如果机器心跳超过 60 秒则从 Redis 中移除该机器
			model.RDB.HDel(MachineListKey, k)
			model.RDB.Del(fmt.Sprintf(MachineInfoKey, hb.IP))
			continue
		}
	}
	return nil
}

func ClearEngineHeartbeatData() {
	machineMap, _ := model.RDB.HGetAll(MachineListKey).Result()
	for ip, _ := range machineMap {
		model.RDB.LTrim(fmt.Sprintf(MachineInfoKey, ip), 0, MaxSaveHour*60*12)
	}
}
