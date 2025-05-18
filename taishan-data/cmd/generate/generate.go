package main

import (
	"data/config"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"net/url"
	"os"
	"strings"
)

func MustInitConf() {
	viper.SetConfigFile("./conf.yml")
	viper.SetConfigType("yaml")
	wd, err := os.Getwd()
	fmt.Println("当前工作目录：", wd)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&config.Conf); err != nil {
		panic(fmt.Errorf("unmarshal error config file: %w", err))
	}

	fmt.Println("config initialized")
}

const dsnTemplate = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"

func main() {
	MustInitConf()
	c := config.Conf
	dsn := fmt.Sprintf(dsnTemplate, c.MySQL.Username, c.MySQL.Password, c.MySQL.Host, c.MySQL.Port, c.MySQL.DBName, c.MySQL.Charset, url.QueryEscape("Asia/Shanghai"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/query",
	})

	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			detailType, _ := columnType.ColumnType()
			if strings.HasPrefix(detailType, "tinyint(1)") {
				return "bool"
			}
			return "int8"
		},
	}
	g.WithDataTypeMap(dataMap)

	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("snapshot"),
		g.GenerateModel("plan"),
		g.GenerateModel("report"),
		g.GenerateModel("scene_case"),
		g.GenerateModel("sls_store"),
	)

	g.Execute()
}
