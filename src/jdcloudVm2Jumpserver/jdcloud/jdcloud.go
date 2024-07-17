package jdcloud

import (
	"fmt"
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
)

// 获取云主机列表的函数，支持分页
func GetCloudHostList(accessKey, secretKey, region string) ([]models.Instance, error) {
	// 创建JDCloud SDK的配置
	credentials := core.NewCredentials(accessKey, secretKey)
	config := core.NewConfig()
	config.SetEndpoint("vm.jdcloud-api.com")

	// 创建虚拟机客户端
	vmClient := client.NewVmClient(credentials)
	// 初始化分页参数
	pageNumber := 1
	pageSize := 20
	var allInstances []models.Instance
	for {
		// 创建请求对象
		describeInstancesRequest := apis.NewDescribeInstancesRequestWithAllParams(region, &pageNumber, &pageSize, nil, nil)
		// 发送请求
		response, err := vmClient.DescribeInstances(describeInstancesRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %v", err)
		}

		// 将当前页的实例添加到所有实例列表中
		allInstances = append(allInstances, response.Result.Instances...)

		// 检查是否还有下一页
		if len(response.Result.Instances) < pageSize {
			break
		}
		// 移动到下一页
		pageNumber++
	}

	return allInstances, nil
}
