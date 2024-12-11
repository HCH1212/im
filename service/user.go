package service

import (
	"github.com/sirupsen/logrus"
	"im/model"
	"im/serializer"
)

type UserRegisterService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	var count int64 = 0
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).Count(&count)
	if count != 0 {
		return serializer.Response{
			Status: 400,
			Msg:    "已经存在该用户了",
		}
	}
	user = model.User{
		UserName: service.UserName,
	}
	//加密密码
	if err := user.SetPassword(service.Password); err != nil {
		logrus.Info(err)
		return serializer.Response{
			Status: 500,
			Msg:    "加密错误",
		}
	}
	//创建用户
	model.DB.Create(&user)
	return serializer.Response{
		Status: 200,
		Msg:    "创建成功",
	}
}
