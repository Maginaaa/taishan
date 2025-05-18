package dal

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"report/conf"
)

var (
	client influxdb2.Client
)

// MustInitInflux 初始化influx
func MustInitInflux() {
	client = influxdb2.NewClient(conf.Conf.InfluxDB.Url, conf.Conf.InfluxDB.Token)
	// 延迟关闭客户端连接
	defer client.Close()
}

// Query 查询
func Query(fluxQuery string) (*api.QueryTableResult, error) {
	queryAPI := client.QueryAPI(conf.Conf.InfluxDB.Org)
	results, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return results, nil
}
