package main

import (
	"im/config"
	"im/dao"
	"im/global"
	"im/log"
	"im/router"
)

func main() {
	// 读取配置文件
	global.Config = config.InitConf()
	// 初始化日志
	log.InitLogger()
	// 连接数据库
	dao.InitMysql()
	// 路由
	router.InitRouter()
}
