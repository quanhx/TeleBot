package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
	"xcheck.info/telebot/pkg/conf"
)

func InitDb() *gorm.DB {
	dborm, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True",
		conf.EnvConfig.MySQL.User, conf.EnvConfig.MySQL.Password,
		conf.EnvConfig.MySQL.Host, conf.EnvConfig.MySQL.Port, conf.EnvConfig.MySQL.DB))
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	dborm.DB().SetConnMaxLifetime(30 * time.Minute)

	if conf.EnvConfig.Environment == conf.EnvironmentLocal || conf.EnvConfig.Environment == conf.EnvironmentDev {
		dborm.LogMode(true)
	}

	// set max idle and max open cons
	dborm.DB().SetMaxIdleConns(conf.EnvConfig.MySQL.MaxIdleConns)
	dborm.DB().SetMaxOpenConns(conf.EnvConfig.MySQL.MaxOpenConns)

	return dborm
}
