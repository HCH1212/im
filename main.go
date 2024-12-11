package main

import (
	"im/config"
	"im/router"
	"im/service"
)

func main() {
	config.Init()

	go service.Manager.Start()

	r := router.NewRouter()
	_ = r.Run(config.HttpPort)
}
