package logic

import (
	"github.com/gin-gonic/gin"
	"report/internal/dal"
	"report/log"
)

func BatchGetPlanName(ctx *gin.Context, ids []int32) (mp map[int32]string, err error) {
	tx := dal.GetQuery().Plan
	planList, err := tx.WithContext(ctx).Where(tx.ID.In(ids...)).Find()
	if err != nil {
		log.Logger.Error("logic.plan.BatchGetPlanName.planListQuery(),err:", err)
		return nil, err
	}

	mp = make(map[int32]string)
	for _, p := range planList {
		mp[p.ID] = p.PlanName
	}
	return
}

func GetPlanNameByID(ctx *gin.Context, id int32) (planName string, err error) {
	tx := dal.GetQuery().Plan
	plan, err := tx.WithContext(ctx).Where(tx.ID.Eq(id)).First()
	if err != nil {
		log.Logger.Error("logic.plan.BatchGetPlanName.planListQuery(),err:", err)
		return "", err
	}
	return plan.PlanName, nil
}
