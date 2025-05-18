package logic

import (
	"data/internal/dal"
	"github.com/gin-gonic/gin"
	"math"
)

type DashBoard struct {
	TotalPlan   int64     `json:"total_plan"`
	TotalReport int64     `json:"total_report"`
	TotalCost   float64   `json:"total_cost"`
	Graph       GraphData `json:"graph"`
}

type GraphData struct {
	Column [3]string `json:"column"`
	Data   [][3]any  `json:"data"`
}

const (
	HttpCall  = "http"
	RedisCall = "redis"
)

func GetDashboardData(ctx *gin.Context) DashBoard {
	planCount, _ := dal.GetQuery().Plan.WithContext(ctx).Count()
	var rep []struct {
		Cost  float64
		Count int64
	}
	reportTx := dal.GetQuery().Report
	reportTx.WithContext(ctx).Select(reportTx.Vum.Sum().Mul(3).Div(1000).As("cost"), reportTx.Vum.Count().As("count")).Scan(&rep)
	var respArr []struct {
		Sum   int32   // vum合
		Count int32   // vum个数/压测次数
		Cost  float64 // 费用
		Date  string
	}
	sql := `SELECT DATE_FORMAT(end_time,'%m-%d') AS date, sum(vum) as 'sum', count(vum) as 'count', sum(vum)*0.003 as 'cost' from report WHERE end_time>= DATE_SUB(NOW(), INTERVAL 30 DAY) GROUP BY date`
	dal.MysqlDB().Raw(sql).Scan(&respArr)
	graphData := make([][3]any, 0, len(respArr))
	for _, res := range respArr {
		graphData = append(graphData, [3]any{res.Date, res.Count, math.Floor(res.Cost)})
	}
	return DashBoard{
		TotalPlan:   planCount,
		TotalReport: rep[0].Count,
		TotalCost:   math.Floor(rep[0].Cost),
		Graph: GraphData{
			Column: [3]string{"日期", "压测次数", "费用"},
			Data:   graphData,
		},
	}

}
