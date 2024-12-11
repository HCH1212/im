package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func Database(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalln("mysql连接失败", err)
	}

	_ = db.AutoMigrate(&User{})

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)               // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              // 最大可容纳
	sqlDB.SetConnMaxLifetime(time.Hour * 4) // 连接最大复用时间， 不能超过mysql的wait_timeout
	DB = db
}
