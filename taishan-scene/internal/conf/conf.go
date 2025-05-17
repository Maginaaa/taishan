package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Http        Http        `yaml:"http"`
	Url         Url         `yaml:"url"`
	MySQL       MySQL       `yaml:"mysql"`
	Mongo       Mongo       `yaml:"mongo"`
	Redis       Redis       `yaml:"redis"`
	RedisReport RedisReport `yaml:"redisReport"`
	Log         Log         `yaml:"log"`
	OSSConfig   OSSConfig   `yaml:"ossConfig"`
	FsConfig    FsConfig    `yaml:"fsConfig"`
}

type Log struct {
	InfoPath string `yaml:"InfoPath"`
	ErrPath  string `yaml:"ErrPath"`
}

type Http struct {
	Port int `yaml:"port"`
}

type MySQL struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type Mongo struct {
	DSN      string `yaml:"dsn"`
	Database string `yaml:"database"`
	PoolSize uint64 `yaml:"pool_size"`
}

type Kafka struct {
	Host  string `yaml:"host"`
	Topic string `yaml:"topic"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
type RedisReport struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type OSSConfig struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	Endpoint        string `yaml:"Endpoint"`
	BucketName      string `yaml:"BucketName"`
}

type Url struct {
	Machine string `yaml:"machine"`
	Report  string `yaml:"report"`
	Account string `yaml:"account"`
	Task    string `yaml:"task"`
	Data    string `yaml:"data"`
}

type FsConfig struct {
	AppID       string `yaml:"AppID"`
	AppSecret   string `yaml:"AppSecret"`
	ChatGroupID string `yaml:"ChatGroupID"`
}

func MustInitConf(configFile string) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unmarshal error config file: %w", err))
	}

	fmt.Println("config initialized")
}
