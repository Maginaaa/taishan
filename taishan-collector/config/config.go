package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Http        Http        `yaml:"http"`
	Kafka       Kafka       `yaml:"kafka"`
	ReportRedis ReportRedis `yaml:"reportRedis"`
	Redis       Redis       `yaml:"redis"`
	Mongo       Mongo       `yaml:"mongo"`
	InfluxDB    InfluxDB    `yaml:"influxDB"`
	Log         Log         `yaml:"log"`
}

type Machine struct {
	MaxGoroutines int    `yaml:"maxGoroutines"`
	ServerType    int    `yaml:"serverType"`
	NetName       string `yaml:"netName"`
	DiskName      string `yaml:"diskName"`
}

type Http struct {
	Port int `yaml:"port"`
}

type Kafka struct {
	Address string `yaml:"address"`
	Topic   string `yaml:"topic"`
}

type ReportRedis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Log struct {
	Path string `yaml:"path"`
}

type Mongo struct {
	DSN      string `yaml:"dsn"`
	Password string `yaml:"password"`
	DataBase string `yaml:"database"`
	PoolSize uint64 `mapstructure:"pool_size"`
}

type InfluxDB struct {
	Url    string `yaml:"url"`
	Org    string `yaml:"org"`
	Token  string `yaml:"token"`
	Bucket string `yaml:"bucket"`
}

func InitConfig(conf string) {

	viper.SetConfigFile(conf)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err = viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unmarshal error config file: %w", err))
	}

	fmt.Println("config initialized")

}
