package jdcloudVm2Jumpserver

import (
	disk "github.com/jdcloud-api/jdcloud-sdk-go/services/disk/models"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/config"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/jdcloud"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/jumpserver"
	"log"
	"strings"
)

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
		return
	}
	log.Println("Fetched asset list: ", assetMap)
	// 遍历所有地域，获取云主机列表
	var allInstances []models.Instance

	for _, region := range config.JDCloud.Regions {
		instances, err := jdcloud.GetCloudHostList(config.JDCloud.AccessKey, config.JDCloud.SecretKey, region)
		if err != nil {
			log.Printf("Error getting cloud host list for region %s: %v \n", region, err)
			continue
		}
		matchedInstances := matchJdCloudVmTag(instances, config.Tags)
		allInstances = append(allInstances, matchedInstances...)
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

	jumpserver.DeleteJumpServerInstance(config, token, allInstances)

	jumpserver.CreteJumpServerInstance(config, token, allInstances)

	jumpserver.UpdateJumpServerInstance(config, token, allInstances)

}

// matchTags 获取实例中匹配配置文件中标签键的所有标签值
func matchTags(instanceTags []disk.Tag, configTags []config.Tag) []string {
	var matchedValues []string
	for _, configTag := range configTags {
		for _, instanceTag := range instanceTags {
			if instanceTag.Key != nil && *instanceTag.Value == configTag.Key {
				if instanceTag.Value != nil {
					matchedValues = append(matchedValues, configTag.Value)
				}
			}
		}
	}
	return matchedValues
}

func matchJdCloudVmTag(allInstances []models.Instance, configTags []config.Tag) []models.Instance {
	var matchedValues []models.Instance
	for _, instance := range allInstances {
		if instance.Tags != nil && len(instance.Tags) > 0 {
			matched := matchTags(instance.Tags, configTags)
			if matched != nil && len(matched) > 0 {
				matchedValues = append(matchedValues, instance)

			}
		}
	}
	return matchedValues
}
