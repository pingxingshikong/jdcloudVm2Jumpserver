package jdcloudVm2Jumpserver

import (
	disk "github.com/jdcloud-api/jdcloud-sdk-go/services/disk/models"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/config"
)

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
