package model

import (
	"gorm.io/gorm"
	"im/global"
	"time"
)

type User struct {
	gorm.Model
	Name          string    `gorm:"type:varchar(20);not null" json:"name"`
	Password      string    `gorm:"type:varchar(255);not null" json:"password"`
	Phone         string    `gorm:"type:varchar(20)" json:"phone"`
	Email         string    `gorm:"type:varchar(20)" json:"email"`
	Identity      string    `gorm:"type:varchar(20)" json:"identity"`
	ClientIp      string    `gorm:"type:varchar(20)" json:"clientIp"`
	ClientPort    string    `gorm:"type:varchar(20)" json:"clientPort"`
	LoginTime     time.Time `gorm:"type:datetime" json:"loginTime"`
	HeartbeatTime time.Time `gorm:"type:datetime" json:"heartbeatTime"`
	LogoutTime    time.Time `gorm:"type:datetime" json:"logoutTime"`
	IsLogout      bool      `gorm:"type:boolean" json:"isLogout"`
	DeviceInfo    string    `gorm:"type:varchar(255)" json:"deviceInfo"`
}

func (User) TableName() string {
	return "user"
}

func GetUserList() []*User {
	users := make([]*User, 10)
	global.DB.Find(&users)
	return users
}
