package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"machine/model"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	UserInfoQueryMaxMin = 5
	AllMachineKey       = "TaishanMachineList"
	MachineInfoKey      = "MachineMonitor:%s"

	HourMinSec = "15:04:05"
)

func getAllMachine(ctx *gin.Context) map[string]model.HeartBeat {
	return model.GetMachineMap()
}

func getAvailableMachine(ctx *gin.Context) map[string]model.HeartBeat {
	return model.GetAvailableMachine()
}

func GetAllMachineInfoLogic(ctx *gin.Context) any {
	machineMap, _ := model.RDB.HGetAll(AllMachineKey).Result()
	ipSet := make([]string, 0, len(machineMap))
	for ip := range machineMap {
		ipSet = append(ipSet, ip)
	}
	sort.Slice(ipSet, func(i, j int) bool {
		parts1 := strings.Split(ipSet[i], ".")
		parts2 := strings.Split(ipSet[j], ".")
		for i := 0; i < 4; i++ {
			// 将部分转换为整数进行比较
			num1, _ := strconv.Atoi(parts1[i])
			num2, _ := strconv.Atoi(parts2[i])
			if num1 < num2 {
				return true
			} else if num1 > num2 {
				return false
			}
		}
		return false
	})
	return getMachineInfo(ctx, ipSet)
}

func GetMachineInfoLogic(ctx *gin.Context, ipArray []string) any {
	return getMachineInfo(ctx, ipArray)
}

func getMachineInfo(ctx *gin.Context, ipSet []string) any {

	startSec := time.Now().Add(time.Minute * -UserInfoQueryMaxMin).Add(5 * time.Second).Unix()
	timestampSet := make([]string, 0, 12*UserInfoQueryMaxMin)
	for i := 0; i < 12*UserInfoQueryMaxMin; i++ {
		timestampSet = append(timestampSet, time.Unix(startSec, 0).Format(HourMinSec))
		startSec += 5
	}

	cpuArr := make([][]any, len(timestampSet)+1)
	memArr := make([][]any, len(timestampSet)+1)
	networkArr := make([][]any, len(timestampSet)+1)

	for ipIndex, ip := range ipSet {
		heartbeatInfos, _ := model.RDB.LRange(fmt.Sprintf(MachineInfoKey, ip), 0, 12*UserInfoQueryMaxMin-1).Result()
		var h model.HeartBeat
		_ = json.Unmarshal([]byte(heartbeatInfos[0]), &h)
		if cpuArr[0] == nil {
			cpuArr[0] = make([]any, len(ipSet)+1)
			cpuArr[0][0] = "time"
		}
		cpuArr[0][ipIndex+1] = h.Name
		if memArr[0] == nil {
			memArr[0] = make([]any, len(ipSet)+1)
			memArr[0][0] = "time"
		}
		memArr[0][ipIndex+1] = h.Name
		if networkArr[0] == nil {
			networkArr[0] = make([]any, len(ipSet)*2+1)
			networkArr[0][0] = "time"
		}
		networkArr[0][ipIndex*2+1] = h.Name + " Sent"
		networkArr[0][ipIndex*2+2] = h.Name + " Received"
		for timestampIndex, timestamp := range timestampSet {
			if cpuArr[timestampIndex+1] == nil {
				cpuArr[timestampIndex+1] = make([]any, len(ipSet)+1)
				cpuArr[timestampIndex+1][0] = timestamp
			}
			if memArr[timestampIndex+1] == nil {
				memArr[timestampIndex+1] = make([]any, len(ipSet)+1)
				memArr[timestampIndex+1][0] = timestamp
			}
			if networkArr[timestampIndex+1] == nil {
				networkArr[timestampIndex+1] = make([]any, len(ipSet)*2+1)
				networkArr[timestampIndex+1][0] = timestamp
			}
			var hb model.HeartBeat
			if len(heartbeatInfos) > 1 {
				_ = json.Unmarshal([]byte(heartbeatInfos[len(heartbeatInfos)-timestampIndex-1]), &hb)
				cpuPercent, _ := decimal.NewFromFloat(hb.CpuUsage).RoundFloor(2).Float64()
				cpuArr[timestampIndex+1][ipIndex+1] = cpuPercent
				memPercent, _ := decimal.NewFromFloat(hb.MemInfo.UsedPercent).RoundFloor(2).Float64()
				memArr[timestampIndex+1][ipIndex+1] = memPercent
			}

			//preSent := uint64(0)
			//preRec := uint64(0)
			//if timestampIndex != 0 {
			//	var preHb model.HeartBeat
			//	_ = json.Unmarshal([]byte(heartbeatInfos[len(heartbeatInfos)-timestampIndex]), &preHb)
			//	for _, pNet := range preHb.Networks {
			//		preSent += pNet.BytesSent
			//		preRec += pNet.BytesRecv
			//	}
			//}
			//currentSent := uint64(0)
			//currentRec := uint64(0)
			//for _, cNet := range hb.Networks {
			//	currentSent += cNet.BytesSent
			//	currentRec += cNet.BytesRecv
			//}
			//sentMbVal, _ := decimal.NewFromUint64(currentSent - preSent).Div(decimal.NewFromUint64(1024 * 1024)).RoundFloor(2).Float64()
			//recMbVal, _ := decimal.NewFromUint64(currentRec - preRec).Div(decimal.NewFromUint64(1024 * 1024)).RoundFloor(2).Float64()
			//networkArr[timestampIndex+1][ipIndex*2+1] = sentMbVal
			//networkArr[timestampIndex+1][ipIndex*2+2] = recMbVal
		}
	}

	return map[string]any{
		"cpu": cpuArr,
		"mem": memArr,
		//"network": networkArr,
	}
}
