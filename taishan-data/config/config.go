package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Http      Http      `yaml:"http"`
	Kafka     Kafka     `yaml:"kafka"`
	Mongo     Mongo     `yaml:"mongo"`
	InfluxDB  InfluxDB  `yaml:"influxDB"`
	Log       Log       `yaml:"log"`
	OSSConfig OSSConfig `yaml:"ossConfig"`
	SLSConfig SLSConfig `yaml:"slsConfig"`
	Url       Url       `yaml:"url"`
	MySQL     MySQL     `yaml:"mysql"`
	Redis     Redis     `yaml:"redis"`
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

type OSSConfig struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	Endpoint        string `yaml:"Endpoint"`
	BucketName      string `yaml:"BucketName"`
}

type SLSConfig struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	EndPoint        string `yaml:"EndPoint"`
}

type Url struct {
	Report  string `yaml:"report"`
	Account string `yaml:"account"`
}

type MySQL struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
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
