package router

import (
	"github.com/gin-gonic/gin"
	"im/api"
	"im/service"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	v1 := r.Group("/")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, "pong")
		})

		v1.POST("/user/register", api.UserRegister)

		v1.GET("/ws", service.Handler)
	}

	return r
}
