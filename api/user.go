package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"im/service"
)

func UserRegister(c *gin.Context) {
	var userRegisterService service.UserRegisterService
	if err := c.ShouldBind(&userRegisterService); err == nil {
		res := userRegisterService.Register()
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrResponse(err))
		logrus.Info(err.Error())
	}
}
