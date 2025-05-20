package dal

import (
	"collector/config"
	"collector/internal/biz/log"
	"collector/middleware"
	"collector/model"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

var (
	influxClient influxdb2.Client
)

func MustInitInfluxDB() {
	influxClient = influxdb2.NewClient(config.Conf.InfluxDB.Url, config.Conf.InfluxDB.Token)
	_, err := influxClient.Ping(context.Background())
	if err != nil {
		log.Logger.Error("influxdb initialized err")
		return
	}
	fmt.Println("influxdb initialized")
}

func BatchInsertTestData(reportId int32, caseMap map[int32]*model.StageCaseResult) (err error) {
	insertTime := time.Now().In(middleware.Location).Unix()
	org := config.Conf.InfluxDB.Org
	bucket := config.Conf.InfluxDB.Bucket
	writeAPI := influxClient.WriteAPIBlocking(org, bucket)
	points := make([]*write.Point, 0)
	for _, caseResult := range caseMap {
		tags := map[string]string{
			"case_id":    fmt.Sprintf("%d", caseResult.CaseID),
			"scene_id":   fmt.Sprintf("%d", caseResult.SceneID),
			"scene_type": fmt.Sprintf("%d", caseResult.SceneType),
		}
		fields, _ := convertor.StructToMap(&caseResult.BaseData)
		point := write.NewPoint(fmt.Sprintf("%d", reportId), tags, fields, time.Unix(insertTime, 0).In(middleware.Location)) // time.Now() 可以替换为自定义时间
		points = append(points, point)
	}
	if err = writeAPI.WritePoint(context.Background(), points...); err != nil {
		log.Logger.Error("测试数据写入influxdb失败：", err)
	}
	return
}
