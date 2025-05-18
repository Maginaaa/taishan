package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"report/internal/dal"
	"report/log"
	"strconv"
)

type PlanRunLockVal struct {
	ReportID int32   `json:"report_id"`
	SceneIDs []int32 `json:"scene_ids"`
}

const (
	PlanRunLock = "plan:lock:%d"

	AllPartitions  = "AllPartitions"
	UsedPartitions = "UsedPartitions"

	RunKafkaPartition = "RunKafkaPartition"
)

func PlanRunningDelLock(ctx context.Context, planId int32) error {
	_, err := dal.RDB.Del(ctx, fmt.Sprintf(PlanRunLock, planId)).Result()
	if err != nil {
		log.Logger.Error("logic.report.PlanRunningDelLock.Del ，err:", err)
		return err
	}
	return err
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
