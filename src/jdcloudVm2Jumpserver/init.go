package jdcloudVm2Jumpserver

import (
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/config"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/jdcloud"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/jumpserver"
	"log"
	"strings"
	"time"
)

func Init() {
	// 读取配置文件
	config, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	for _, da := range config.Tags {
		log.Printf("key %s value %s", da.Key, da.Value)
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

func runTask(config *config.Config) {

	// 获取 Jumpserver token
	token, err := jumpserver.GetToken(config.Jumpserver.URL, config.Jumpserver.User, config.Jumpserver.Password)
	if err != nil {
		log.Printf("Error getting Jumpserver token: %v", err)
		return
	}

	// 获取资产列表
	assetMap, err := jumpserver.FetchAssetObjectListLabels(config, config.Jumpserver.URL, "/api/v1/assets/hosts/", token)
	if err != nil {
		log.Fatalf("Error fetching asset list: %v", err)
	}
	log.Println("Fetched asset list: ", assetMap)
	// 遍历所有地域，获取云主机列表
	var allInstances []models.Instance
	var allInstancesSource []models.Instance

	for _, region := range config.JDCloud.Regions {
		instances, err := jdcloud.GetCloudHostList(config.JDCloud.AccessKey, config.JDCloud.SecretKey, region)
		if err != nil {
			log.Printf("Error getting cloud host list for region %s: %v \n", region, err)
			continue
		}

		matchedInstances := matchJdCloudVmTag(instances, config.Tags)
		allInstances = append(allInstances, matchedInstances...)
		allInstancesSource = append(allInstancesSource, instances...)
	}

	// 打印云主机数量
	log.Printf("Total number of cloud hosts: %d\n", len(allInstances))

	var cloudHostAddresses []string
	for _, host := range allInstances {
		cloudHostAddresses = append(cloudHostAddresses, host.PrivateIpAddress)
	}
	if cloudHostAddresses != nil && len(cloudHostAddresses) > 0 {
		log.Printf("Fetched cloud host list: %s \n", strings.Join(cloudHostAddresses, ", "))
	}

	jumpserver.DeleteByJdCloudNotExistButJumpServerExist(config, token, allInstances)

	jumpserver.DeleteNewJumpServerInstance(config, token, allInstances)

	//创建新的云主机
	jumpserver.CreteNewJumpServerInstance(config, token, allInstances)

	jumpserver.UpdateNewJumpServerInstance(config, token, allInstances)

}
