package jumpserver

import (
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/config"
	"log"
)

// 京东云有此资源 但是jumpserver也有 但是jumpserver缺少pk  并且符合本程序的标签规则
func UpdateNewJumpServerInstance(config *config.Config, token string, allInstances []models.Instance) {
	// 获取资产列表
	assetMap, err := FetchAssetObjectListLabelsAll(config, config.Jumpserver.URL, "/api/v1/assets/hosts/", token)
	if err != nil {
		log.Fatalf("Error fetching asset list: %v", err)
	}

	instanceList := InstanceListByAssetList(allInstances, assetMap)

	if instanceList == nil {
		return
	}

	for _, instance := range instanceList {
		jpInstances := getAssetListByIp(assetMap, instance.PrivateIpAddress)
		if jpInstances != nil && len(jpInstances) > 0 {
			asset := jpInstances[0]
			pks := lastPk(config, instance, asset)
			if pks == nil || len(pks) == 0 {
				continue
			}
			log.Printf("[UpdateNewJumpServerInstance]update instance %s save %d pk", instance.PrivateIpAddress, len(pks))
			UpdateAssetJd2JumpServer(config, token, pks, asset)
		}

	}

}

func lastPk(config *config.Config, instance models.Instance, asset Asset) []string {
	var pks []string
	instanceTag := instance.Tags
	if instanceTag == nil || len(instanceTag) == 0 {
		return pks
	}
	matchTag := matchTags(instanceTag, config.Tags)
	if matchTag == nil || len(matchTag) == 0 {
		return pks
	}

	nodes := asset.Nodes
	if nodes == nil || len(nodes) == 0 {
		return matchTag
	}
	var apks []string
	for _, node := range nodes {
		apks = append(apks, node.ID)
	}

	if isSubset(matchTag, apks) {
		return pks
	}
	pks = append(pks, matchTag...)
	pks = append(pks, apks...)

	return removeDuplicateStrings(pks)
}

// 根据京东云主机地址获取jumpserver 配置
func AssetListByInstanceList(allInstances []models.Instance, assets map[string]Asset) []Asset {
	var assetList []Asset
	for _, instance := range allInstances {
		for _, asset := range assets {
			if instance.PrivateIpAddress == asset.Address {
				assetList = append(assetList, asset)
			}
		}
	}
	return assetList
}

// 根据京东云主机地址获取jumpserver 配置
func InstanceListByAssetList(allInstances []models.Instance, assets map[string]Asset) []models.Instance {
	var instanceList []models.Instance
	for _, instance := range allInstances {
		for _, asset := range assets {
			if instance.PrivateIpAddress == asset.Address {
				instanceList = append(instanceList, instance)
			}
		}
	}
	return instanceList
}

// jumpserver.Asset去重
func removeDuplicateStrings(pks []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, pk := range pks {
		if _, ok := seen[pk]; !ok {
			seen[pk] = true
			result = append(result, pk)
		}
	}
	return result
}

// containsString 判断切片 b 是否包含字符串 s
func containsString(b []string, s string) bool {
	for _, v := range b {
		if v == s {
			return true
		}
	}
	return false
}

// isSubset 判断切片 a 是否是切片 b 的子集
func isSubset(a, b []string) bool {
	for _, v := range a {
		if !containsString(b, v) {
			return false
		}
	}
	return true
}

/*
*
根据ip获取jumpserver相同ip的节点
*/
func getAssetListByIp(assetMap map[string]Asset, ip string) []Asset {
	var diff []Asset
	for _, asset := range assetMap {
		if asset.Address == ip {
			diff = append(diff, asset)
		}
	}
	return diff
}
