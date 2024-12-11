package config

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
	"im/model"
	"strings"
)

var (
	MongoDBClient *mongo.Client
	AppMode       string
	HttpPort      string
	Db            string
	DbHost        string
	DbPort        string
	DbUser        string
	DbPassword    string
	DbName        string
	MongoDBName   string
	MongoDBAddr   string
	MongoDBPwd    string
	MongoDBPort   string
)

func Init() {
	// 从本地环境读取
	file, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Printf("Fail to parse 'config/config.ini': %v", err)
	}
	LoadServer(file)
	LoadMysql(file)
	LoadMongoDB(file)
	MongoDB() // 连接mongodb
	//MySQL
	path := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	model.Database(path)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassword = file.Section("mysql").Key("DbPassword").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadMongoDB(file *ini.File) {
	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
}

func MongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoDBPort)
	var err error
	MongoDBClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logrus.Info(err)
		panic(err)
	}
	logrus.Info("Connected to MongoDB")
}
