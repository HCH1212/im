package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	System System `yaml:"system"`
	Mysql  Mysql  `yaml:"mysql"`
	Logger Logger `yaml:"logger"`
	JWT    JWT    `yaml:"jwt"`
}

const ConfigFile = "config.yaml"

// InitConf 读取yaml文件的配置
func InitConf() *Config {
	c := &Config{}
	yamlConf, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		logrus.Fatalln(fmt.Errorf("yaml Unmarshal err:%v", err))
	}
	logrus.Println("config yamlFile load init success")
	return c
}
