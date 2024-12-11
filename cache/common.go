package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"strconv"
)

// 缓存redis

var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func init() {
	file, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Println("Redis 配置文件读取错误，请检查文件路径:", err)
	}
	LoadRedisData(file)
	Redis() // 连接redis
}

func LoadRedisData(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

func Redis() {
	db, err := strconv.ParseUint(RedisDbName, 10, 64) // string to uint64
	if err != nil {
		logrus.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		DB:   int(db),
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		logrus.Info(err.Error())
		panic(err)
	}

	RedisClient = client
}
