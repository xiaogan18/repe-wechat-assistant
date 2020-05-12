package main

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/bll"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type YamlConfig struct {
	Database struct {
		DriverName string `yaml:"DriverName"`
		DataSource string `yaml:"DataSource"`
	} `yaml:"Database"`
	System struct {
		WebListen string `yaml:"WebListen"`
		BotListen string `yaml:"BotListen"`
	} `yaml:"System"`
	CommandConfig bll.ContentConfig `yaml:"CommandConfig"`
}

func initConfig(filename string) (*YamlConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	v := &YamlConfig{}
	err = yaml.Unmarshal(data, v)
	return v, err
}

type SetupArgs struct {
	Log    *int
	Config *string
}
