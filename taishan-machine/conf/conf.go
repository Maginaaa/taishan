package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Http  Http  `yaml:"http"`
	Redis Redis `yaml:"redis"`
	Log   Log   `yaml:"log"`
}

type Log struct {
	InfoPath string `yaml:"InfoPath"`
	ErrPath  string `yaml:"ErrPath"`
}

type Http struct {
	Port int `yaml:"port"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
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
