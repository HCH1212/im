package router

import (
	"context"
	"im/global"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func InitRouter() {
	h := server.Default(server.WithHostPorts(global.Config.System.Port))

	// 用户组
	userRouter := h.Group("/user")
	{
		userRouter.POST("")
	}

	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})
	h.Spin()
}
