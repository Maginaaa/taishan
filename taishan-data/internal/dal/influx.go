package dal

import (
	"context"
	"data/config"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
)

var (
	influxClient influxdb2.Client
)

func MustInitInfluxDB() {
	var err error
	influxClient = influxdb2.NewClient(config.Conf.InfluxDB.Url, config.Conf.InfluxDB.Token)
	_, err = influxClient.Ping(context.Background())
	if err != nil {
		log.Fatal("influxdb initialized err")
	}
	fmt.Println("influxdb initialized")
}
