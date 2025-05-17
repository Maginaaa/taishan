package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"strconv"
	"time"
)

const (
	PlanRunLock = "plan:lock:%d"

	AllPartitions  = "AllPartitions"
	UsedPartitions = "UsedPartitions"
	UsedEngines    = "UsedEngines"

	RunKafkaPartition = "RunKafkaPartition"
)

type PlanRunLockVal struct {
	ReportID int32   `json:"report_id"`
	SceneIDs []int32 `json:"scene_ids"`
}

func PartitionInit() {
	partitionIds := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	_, err := dal.RDB.SAdd(context.Background(), AllPartitions, partitionIds...).Result()
	if err != nil {
		fmt.Println("logic.redis.UsePartition.SAdd ，err:", err)
		return
	}
}

// planRunningAddLock 添加计划执行锁
func planRunningAddLock(ctx context.Context, planId, reportId int32, sceneIds []int32) bool {
	val := PlanRunLockVal{
		ReportID: reportId,
		SceneIDs: sceneIds,
	}
	s, _ := json.Marshal(val)
	return dal.RDB.SetNX(ctx, fmt.Sprintf(PlanRunLock, planId), string(s), 150*60*time.Second).Val()
}

// PlanIsRunning 查询计划是否正在执行
func PlanIsRunning(ctx context.Context, planId int32) (isRunning bool, err error) {
	result, err := dal.RDB.Exists(ctx, fmt.Sprintf(PlanRunLock, planId)).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func getPlanRunningLock(ctx context.Context, planId int32) (isRunning bool, val PlanRunLockVal, err error) {
	isExists, err := dal.RDB.Exists(ctx, fmt.Sprintf(PlanRunLock, planId)).Result()
	if isExists == 0 {
		return false, val, nil
	}
	result, err := dal.RDB.Get(ctx, fmt.Sprintf(PlanRunLock, planId)).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(result), &val)
	if err != nil {
		return
	}
	return true, val, nil
}

func PlanRunningDelLock(ctx context.Context, planId int32) error {
	_, err := dal.RDB.Del(ctx, fmt.Sprintf(PlanRunLock, planId)).Result()
	if err != nil {
		log.Logger.Error("logic.report.PlanRunningDelLock.Del ，err:", err)
		return err
	}
	return err
}

// GetAvailablePartition 查询当前可用的分区
func GetAvailablePartition(ctx context.Context) (ids []int32, err error) {
	result, err := dal.RDB.SDiff(ctx, AllPartitions, UsedPartitions).Result()
	if err != nil {
		log.Logger.Error("logic.report.GetAvailablePartition.SMembers ，err:", err)
		return
	}
	ids = make([]int32, 0)
	for _, p := range result {
		id, _ := strconv.Atoi(p)
		ids = append(ids, int32(id))
	}
	return
}

// usePartition 使用分区
func usePartition(ctx context.Context, partition []int32) (success bool, err error) {
	// 将分区信息填入UsedPartitions
	ps := make([]string, 0)
	for _, p := range partition {
		ps = append(ps, strconv.Itoa(int(p)))
	}
	length, err := dal.RDB.SAdd(ctx, UsedPartitions, ps).Result()
	if err != nil {
		log.Logger.Error("logic.redis.UsePartition.SAdd ，err:", err)
		return
	}
	// 消息发布
	for _, p := range ps {
		log.Logger.Infof("scene.usePartition 通知partition %s 开始消费", p)
		err = dal.RDB.Publish(ctx, RunKafkaPartition, p).Err()
		if err != nil {
			log.Logger.Error("logic.redis.UsePartition.Publish ，err:", err)
			continue
		}
	}
	return length == int64(len(partition)), nil
}

func useEngine(ctx context.Context, engineIps []string) (success bool, err error) {
	length, err := dal.RDB.SAdd(ctx, UsedEngines, engineIps).Result()
	if err != nil {
		log.Logger.Error("logic.redis.useEngine.SAdd ，err:", err)
		return
	}
	return length == int64(len(engineIps)), nil
}

func getUsingEngine(ctx context.Context) (ips []string, err error) {
	return dal.RDB.SMembers(ctx, UsedEngines).Result()
}
