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
	Heartbeat   Heartbeat   `yaml:"heartbeat"`
	Machine     Machine     `yaml:"machine"`
	Log         Log         `yaml:"log"`
	OSSConfig   OSSConfig   `yaml:"ossConfig"`
	Mongo       Mongo       `yaml:"mongo"`
}

type Machine struct {
	MaxGoroutines int    `yaml:"maxGoroutines"`
	ServerType    int    `yaml:"serverType"`
	NetName       string `yaml:"netName"`
	DiskName      string `yaml:"diskName"`
}

type Heartbeat struct {
	Port      int32  `yaml:"port"`
	Region    string `yaml:"region"`
	Duration  int64  `yaml:"duration"`
	Resources int64  `yaml:"resources"`
}
type Http struct {
	Address string `yaml:"address"`
}

type Kafka struct {
	Address string `yaml:"address"`
	TopIc   string `yaml:"topic"`
}

type ReportRedis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}
type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}

type Log struct {
	Path string `yaml:"path"`
}

type OSSConfig struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	Endpoint        string `yaml:"Endpoint"`
	BucketName      string `yaml:"BucketName"`
}

type Mongo struct {
	DSN      string `yaml:"dsn"`
	Password string `yaml:"password"`
	DataBase string `yaml:"database"`
	PoolSize uint64 `yaml:"pool_size"`
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
