package jumpserver

import (
	disk "github.com/jdcloud-api/jdcloud-sdk-go/services/disk/models"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/config"
	"log"
	"strings"
)

func CreteNewJumpServerInstance(config *config.Config, token string, allInstances []models.Instance) {
	// 获取资产列表
	assetMap, err := FetchAssetObjectListLabels(config, config.Jumpserver.URL, "/api/v1/assets/hosts/", token)
	if err != nil {
		log.Fatalf("Error fetching asset list: %v \n", err)
	}
	diff1 := DifferenceFet(assetMap, allInstances)
	var diff1Ip []string
	for _, instance := range diff1 {
		diff1Ip = append(diff1Ip, instance.PrivateIpAddress)
	}
	// 将云主机信息写入 Jumpserver
	for _, instance := range diff1 {
		_, err := FetchAssetJd2JumpServer(config, instance, token)
		if err != nil {
			log.Printf("Error creating asset for instance %s: %v \n", instance.InstanceId, err)
			continue
		}
	}

}

// Difference 返回 GetCloudHostList 中存在但 FetchAssetList 中不存在的元素的 ID
func DifferenceFet(assetMap map[string]Asset, allInstances []models.Instance) []models.Instance {
	var diff []models.Instance
	for _, address := range allInstances {
		if _, exists := assetMap[address.PrivateIpAddress]; !exists {
			diff = append(diff, address)
		}
	}
	return diff
}

func FetchAssetJd2JumpServer(config *config.Config, instance models.Instance, token string) (string, error) {
	if instance.Tags != nil {
		pks := matchTags(instance.Tags, config.Tags)
		log.Printf("[FetchAssetJd2JumpServer] add %s and pks size: %d pks:  %s \n", instance.PrivateIpAddress, len(pks), strings.Join(pks, ","))
		FetchAssetJd2JumpServerHasPK(config, instance, token, pks)
	}
	return "", nil
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

func matchTagsAccount(instanceTags []disk.Tag, configTags []config.Tag) []string {
	var matchedValues []string
	for _, configTag := range configTags {
		for _, instanceTag := range instanceTags {
			if instanceTag.Key != nil && *instanceTag.Value == configTag.Key {
				if instanceTag.Value != nil {
					matchedValues = append(matchedValues, configTag.Accounts...)
				}
			}
		}
	}
	return matchedValues
}

func FetchAssetJd2JumpServerHasPK(config *config.Config, instance models.Instance, token string, pks []string) (string, error) {
	name := "JD-" + instance.InstanceName
	//根据Instance获取pk
	assetData := map[string]interface{}{
		"set_active": true,
		"name":       name,
		"address":    instance.PrivateIpAddress, // 假设使用第一个私有IP地址
		"platform": map[string]interface{}{
			"pk": 1, // 根据实际情况设置平台的 pk
		},
		"nodes": []map[string]interface{}{
			//{
			//	"pk": pk, // 根据实际情况设置节点的 pk
			//},
		},
		"accounts": []map[string]interface{}{
			//{
			//	"template": template,
			//},
		},
		"labels": config.SysLabel,
		"protocols": []map[string]interface{}{
			{
				"name": "ssh",
				"port": "22",
			},
			{
				"name": "sftp",
				"port": "22",
			},
		},
	}

	// 将 pks 组装进 nodes
	for _, pk := range pks {
		node := map[string]interface{}{
			"pk": pk,
		}
		assetData["nodes"] = append(assetData["nodes"].([]map[string]interface{}), node)
	}

	accounts := matchTagsAccount(instance.Tags, config.Tags)

	if accounts != nil {
		lastAccounts := removeDuplicateStrings(accounts)
		for _, account := range lastAccounts {
			template := map[string]interface{}{
				"template": account,
			}
			assetData["accounts"] = append(assetData["accounts"].([]map[string]interface{}), template)
		}
	}

	result, err := CreateAsset(config.Jumpserver.URL, "/api/v1/assets/hosts/", token, assetData)
	if err != nil {
		log.Printf("Error creating asset for instance %s: %v \n", instance.InstanceId, err)
		if err != nil {
			return "", err
		}
	}
	log.Printf("Created asset for instance %s: %s\n", instance.InstanceId, result)
	return result, err
}
