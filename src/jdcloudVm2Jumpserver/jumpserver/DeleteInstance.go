package jumpserver

import (
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/config"
	"log"
	"strings"
)

// 京东云没有此资源 但是jumpserver有 根据ip判断 并且符合本程序的标签规则
func DeleteByJdCloudNotExistButJumpServerExist(config *config.Config, token string, allInstances []models.Instance) {
	// 获取资产列表
	assetMap, err := FetchAssetObjectListLabels(config, config.Jumpserver.URL, "/api/v1/assets/hosts/", token)
	if err != nil {
		log.Fatalf("Error fetching asset list: %v", err)
	}
	assetList := JdCloudNotExistButJumpServerExistDifference(assetMap, allInstances)
	var addressList []string
	// 删除差集中的资产
	for _, asset := range assetList {
		othersPk := getOtherPk(asset, config.Tags)
		if othersPk == nil || len(othersPk) == 0 {
			deleteJumpServerInstance(config, token, asset)
		} else {
			UpdateAssetJd2JumpServer(config, token, othersPk, asset)
		}
		addressList = append(addressList, asset.Address)
	}
	if len(addressList) > 0 {
		log.Printf("[DeleteByJdCloudNotExistButJumpServerExist]delete[delete] instance %s size: %d ", strings.Join(addressList, ", "), len(addressList))
	}
}

// 京东云有此资源 但是jumpserver也有 但是jumpserver多pk  并且符合本程序的标签规则
func DeleteNewJumpServerInstance(config *config.Config, token string, allInstances []models.Instance) {
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
			pks := lastPk1(config, instance, asset)
			if pks == nil || len(pks) == 0 {
				continue
			}
			log.Printf("[DeleteNewJumpServerInstance]delete[update] instance %s save %d pk", instance.PrivateIpAddress, len(pks))
			UpdateAssetJd2JumpServer(config, token, pks, asset)
		}
	}

}

// 京东云不存在的ip，但是jumpserver存在
func JdCloudNotExistButJumpServerExistDifference(assetMap map[string]Asset, allInstances []models.Instance) []Asset {
	cloudHostMap := make(map[string]bool)
	for _, instance := range allInstances {
		cloudHostMap[instance.PrivateIpAddress] = true
	}
	var diff []Asset
	for address, asset := range assetMap {
		if !cloudHostMap[address] {
			diff = append(diff, asset)
		}
	}
	return diff
}

/*
*
是否有其他的目录 非本程序京东云标签，若有则更新
*/
func getOtherPk(asset Asset, configTags []config.Tag) []string {
	var pks []string
	if asset.Nodes != nil && len(asset.Nodes) > 0 {
		for _, node := range asset.Nodes {
			status := true
			for _, tag := range configTags {
				if node.ID == tag.Value {
					status = false
				}
			}
			if status {
				pks = append(pks, node.ID)
			}
		}
	}
	return pks
}

func lastPk1(config *config.Config, instance models.Instance, asset Asset) []string {
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

	if isSubset(apks, matchTag) {
		return pks
	}

	last1 := intersection(matchTag, apks)

	last := difference(matchTag, apks, config.Tags)

	if last1 != nil && len(last1) > 0 {
		pks = append(pks, last1...)
	}

	if last != nil && len(last) > 0 {
		pks = append(pks, last...)
	}

	if pks == nil || len(pks) == 0 {
		return pks
	}
	return removeDuplicateStrings(pks)
}

// difference 返回切片 b 中比切片 a 多的元素，并且这些元素不在切片 c 中
func difference(a, b []string, configTags []config.Tag) []string {
	var c []string
	for _, tag := range configTags {
		c = append(c, tag.Value)
	}
	var diff []string
	for _, v := range b {
		if !containsString(a, v) && !containsString(c, v) {
			diff = append(diff, v)
		}
	}
	return diff
}

// intersection 返回切片 b 和 a 中相同的元素
func intersection(a, b []string) []string {
	var inter []string
	for _, v := range b {
		if containsString(a, v) {
			inter = append(inter, v)
		}
	}
	return inter
}
