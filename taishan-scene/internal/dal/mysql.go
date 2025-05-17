package dal

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"scene/internal/conf"
	"scene/internal/query"
	"time"
)

var (
	db *gorm.DB
)

const dsnTemplate = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"

func MustInitMySQL() {
	var err error

	c := conf.Conf
	dsn := fmt.Sprintf(dsnTemplate, c.MySQL.Username, c.MySQL.Password, c.MySQL.Host, c.MySQL.Port, c.MySQL.DBName, c.MySQL.Charset, url.QueryEscape("Asia/Shanghai"))

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("fatal error mysql init: %w", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 开启日志
	//if c.Base.IsDebug {
	//	db = db.Debug()
	//}

	fmt.Println("mysql initialized")
}

func MysqlDB() *gorm.DB {
	return db
}

func GetQuery() *query.Query {
	return query.Use(db)
}
