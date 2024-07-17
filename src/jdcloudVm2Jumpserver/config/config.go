package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config 是用于解析 config.yml 文件的结构
type Config struct {
	Tags []Tag `yaml:"tags"`

	SysLabel []string `yaml:"sysLabel"`

	JDCloud struct {
		AccessKey string   `yaml:"accessKey"`
		SecretKey string   `yaml:"secretKey"`
		Regions   []string `yaml:"regions"`
	} `yaml:"jdcloud"`
	Jumpserver struct {
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"jumpserver"`

	Schedule struct {
		Interval int `yaml:"interval"` // 执行间隔时间，以分钟为单位
	} `yaml:"schedule"`
}

// Tag 结构体表示标签的键值对
type Tag struct {
	Key      string   `yaml:"key"`
	Value    string   `yaml:"value"`
	Accounts []string `yaml:"accounts"`
}

// 读取配置文件
func ReadConfig(configPath string) (*Config, error) {
	config := &Config{}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return config, nil
}
