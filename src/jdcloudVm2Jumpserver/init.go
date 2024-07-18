package jdcloudVm2Jumpserver

import (
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/config"
	"log"
	"time"
)

func Init() {
	// 读取配置文件
	config, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	for _, da := range config.Tags {
		log.Printf("key %s value %s \n", da.Key, da.Value)
	}

	// 获取执行间隔时间
	interval := time.Duration(config.Schedule.Interval) * time.Second
	// 创建一个定时器，根据配置文件中的间隔时间触发
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		// 执行任务
		runTask(config)

		// 等待下一个定时器触发
		<-ticker.C
	}
}
