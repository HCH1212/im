package dao

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"im/global"
	"im/model"
	"time"
)

var DB *gorm.DB

func InitMysql() {
	if global.Config.Mysql.Host == "" {
		logrus.Warnln("未配置mysql, 取消gorm连接")
		return
	}
	dsn := global.Config.Mysql.Dsn()

	var mysqlLogger logger.Interface
	if global.Config.System.Env == "debug" {
		// 开发环境显示所有sql
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error) // 只打印错误的sql
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		global.Log.Fatalln("mysql连接失败", err)
	}

	// 自动迁移表
	_ = db.AutoMigrate(
		&model.User{},
	)

	//// 默认用户
	//if db.First(&model.User{}, "id = 1").Error != nil {
	//	db.Create(&model.User{
	//		Name:          "12",
	//		Password:      "12",
	//		LoginTime:     time.Now(),
	//		HeartbeatTime: time.Now(),
	//		LogoutTime:    time.Now(),
	//	})
	//}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)               // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              // 最大可容纳
	sqlDB.SetConnMaxLifetime(time.Hour * 4) // 连接最大复用时间， 不能超过mysql的wait_timeout
	global.DB = db
}
