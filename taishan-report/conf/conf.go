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
	MongoDB     MongoDB     `yaml:"mongodb"`
	InfluxDB    InfluxDB    `yaml:"influxdb"`
	Redis       Redis       `yaml:"redis"`
	RedisReport RedisReport `yaml:"redisReport"`
	Log         Log         `yaml:"log"`
}

type Log struct {
	InfoPath string `yaml:"InfoPath"`
	ErrPath  string `yaml:"ErrPath"`
}

type Http struct {
	Port int `yaml:"port"`
}

type Url struct {
	Machine string `yaml:"machine"`
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

type MongoDB struct {
	DSN      string `yaml:"dsn"`
	Database string `yaml:"database"`
	PoolSize uint64 `mapstructure:"pool_size"`
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
