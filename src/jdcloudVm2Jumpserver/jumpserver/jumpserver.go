package jumpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/config"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Asset 结构体表示每个资产的结构
type Asset struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Address   string       `json:"address"`
	Comment   string       `json:"comment"`
	Domain    interface{}  `json:"domain"`
	Platform  PlatformData `json:"platform"`
	Nodes     []NodeData   `json:"nodes"`
	Labels    []string     `json:"labels"`
	Protocols []Protocol   `json:"protocols"`
	IsActive  bool         `json:"is_active"`
}

// 定义响应结构体
type NodeData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlatformData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Protocol struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// 定义响应结构体
type LabelData struct {
	ID string `json:"id"`
}

// FetchAssetList 查询资产列表
func FetchAssetList(baseURL, api, token string) (map[string]string, error) {
	req, err := http.NewRequest("GET", baseURL+api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var assets []Asset
	err = json.NewDecoder(resp.Body).Decode(&assets)
	if err != nil {
		return nil, err
	}

	assetMap := make(map[string]string)
	for _, asset := range assets {
		assetMap[asset.Address] = asset.ID
	}
	log.Printf("%s", assetMap)
	return assetMap, nil
}

// 根据config中的sysLabel过滤jumpserver中的节点
func FetchAssetListLabels(config *config.Config, baseURL, api, token string) (map[string]string, error) {
	req, err := http.NewRequest("GET", baseURL+api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var assets []Asset
	err = json.NewDecoder(resp.Body).Decode(&assets)
	if err != nil {
		return nil, err
	}

	assetMap := make(map[string]string)
	for _, asset := range assets {
		if asset.Labels != nil {
			if len(asset.Labels) > 0 {
				// 过滤出包含 sysLabel 的 labels
				filteredLabels := filterLabels(asset.Labels, config.SysLabel)
				if filteredLabels != nil && len(filteredLabels) > 0 {
					assetMap[asset.Address] = asset.ID
				}
			} else {
				continue
			}
		} else {
			continue
		}

	}
	log.Printf("%s", assetMap)
	return assetMap, nil
}

func FetchAssetObjectListLabels(config *config.Config, baseURL, api, token string) (map[string]Asset, error) {
	req, err := http.NewRequest("GET", baseURL+api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var assets []Asset
	err = json.NewDecoder(resp.Body).Decode(&assets)
	if err != nil {
		return nil, err
	}

	assetMap := make(map[string]Asset)
	for _, asset := range assets {
		if asset.Labels != nil {
			if len(asset.Labels) > 0 {
				// 过滤出包含 sysLabel 的 labels
				filteredLabels := filterLabels(asset.Labels, config.SysLabel)
				if filteredLabels != nil && len(filteredLabels) > 0 {
					_, err := GetConfigTagByValueList(config.Tags, asset.Nodes)
					if err == nil {
						assetMap[asset.Address] = asset
					}
				}
			} else {
				continue
			}
		} else {
			continue
		}

	}
	return assetMap, nil
}

func FetchAssetObjectListLabelsAll(config *config.Config, baseURL, api, token string) (map[string]Asset, error) {
	req, err := http.NewRequest("GET", baseURL+api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var assets []Asset
	err = json.NewDecoder(resp.Body).Decode(&assets)
	if err != nil {
		return nil, err
	}

	assetMap := make(map[string]Asset)
	for _, asset := range assets {
		if asset.Labels != nil {
			if len(asset.Labels) > 0 {
				// 过滤出包含 sysLabel 的 labels
				filteredLabels := filterLabels(asset.Labels, config.SysLabel)
				if filteredLabels != nil && len(filteredLabels) > 0 {
					//_, err := GetConfigTagByValueList(config.Tags, asset.Nodes)
					//if err == nil {
					//	assetMap[asset.Address] = asset
					//}
					assetMap[asset.Address] = asset
				}
			} else {
				continue
			}
		} else {
			continue
		}

	}
	return assetMap, nil
}

func GetConfigTagByValueList(configTags []config.Tag, nodeDataValueList []NodeData) (config.Tag, error) {
	var result config.Tag
	// Search for the key
	var searchValueList []string
	for _, searchValue := range nodeDataValueList {
		searchValueList = append(searchValueList, searchValue.ID)
		for _, tag := range configTags {
			if tag.Value == searchValue.ID {
				return tag, nil
			}
		}
	}
	return result, fmt.Errorf("key not found: %s", strings.Join(searchValueList, ","))
}

// createAsset 创建资产
func CreateAsset(baseURL, api, token string, assetData map[string]interface{}) (string, error) {
	jsonObjectBody, err := json.Marshal(assetData)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", baseURL+api, bytes.NewBuffer(jsonObjectBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析 JSON 响应
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", err
	}

	// 提取 id 字段
	id, ok := responseData["id"].(string)
	if !ok {
		responseDataStr, err := json.MarshalIndent(responseData, "", "  ")
		if err != nil {
			return "", err
		}
		return "", errors.New(string(responseDataStr))
	}

	return id, nil
}

// createAsset 创建资产
func UpdateAsset(baseURL, api, token string, assetData map[string]interface{}) (string, error) {
	jsonObjectBody, err := json.Marshal(assetData)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("PUT", baseURL+api, bytes.NewBuffer(jsonObjectBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析 JSON 响应
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", err
	}

	// 提取 id 字段
	id, ok := responseData["id"].(string)
	if !ok {
		responseDataStr, err := json.MarshalIndent(responseData, "", "  ")
		if err != nil {
			return "", err
		}
		return "", errors.New(string(responseDataStr))
	}

	return id, nil
}

// TokenResponse 定义从 Jumpserver 获取 token 的响应结构
type TokenResponse struct {
	Token string `json:"token"`
}

// GetToken 从 Jumpserver 获取 token
func GetToken(jmsurl, username, password string) (string, error) {
	url := jmsurl + "/api/v1/authentication/auth/"
	log.Printf("Requesting token from URL: %s\n", url) // 添加调试信息
	queryArgs := strings.NewReader(fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, username, password))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, queryArgs)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var response TokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	if response.Token == "" {
		return "", fmt.Errorf("failed to get token from response: %s", string(body))
	}

	return response.Token, nil
}

// DeleteAssetByID 根据ID删除资产
func DeleteAssetByID(baseURL, id, token string) error {
	url := fmt.Sprintf("%s/api/v1/assets/hosts/%s/", baseURL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-CSRFToken", token)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if body != nil {
		return nil
	}

	return nil
}

func DeleteAssetByLabelID(baseURL, id, token string, labelId string) error {
	url := fmt.Sprintf("%s/api/v1/labels/labeled-resources/%s/?label=%s", baseURL, id, labelId)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-CSRFToken", token)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
	}
	if body != nil {
		return nil
	}
	return nil
}

// filterLabels 根据 sysLabel 过滤 labels
func filterLabels(labels, sysLabel []string) []string {
	var filteredLabels []string
	for _, label := range labels {
		for _, sys := range sysLabel {
			if strings.Contains(label, sys) {
				filteredLabels = append(filteredLabels, label)
			}
		}
	}
	return filteredLabels
}

func GetLabelID(config *config.Config, token string) (string, error) {
	labelValue := strings.Split(config.SysLabel[0], ":")[1]
	url := fmt.Sprintf("%s/api/v1/labels/labels/?value=%s", config.Jumpserver.URL, labelValue)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var labelData []LabelData
	err = json.NewDecoder(resp.Body).Decode(&labelData)
	if err != nil {
		return "", err
	}
	id := labelData[0].ID
	return id, nil
}

func GetLabelInstanceID(config *config.Config, token string, name string) (string, error) {
	labelValue := url.QueryEscape(name)
	url := fmt.Sprintf("%s/api/v1/labels/labeled-resources/?search=%s", config.Jumpserver.URL, labelValue)
	log.Printf(url)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Printf("%v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%v", err)

		return "", err
	}
	defer resp.Body.Close()

	// 读取响应的 body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var labelData []LabelData
	err = json.Unmarshal(body, &labelData) // 使用 body 进行解码
	if err != nil {
		return "", fmt.Errorf("error decoding response body: %w", err)
	}

	if len(labelData) == 0 {
		return "", fmt.Errorf("no label data found for value: %s", labelValue)
	}
	id := labelData[0].ID
	return id, nil
}

// 删除jumpserver 资源
func deleteJumpServerInstance(config *config.Config, token string, asset Asset) {
	jpLabelId, jpLabelErr := GetLabelID(config, token)
	jpInstanceId, jpnInstanceErr := GetLabelInstanceID(config, token, asset.Name)
	if jpLabelErr != nil || jpnInstanceErr != nil {
		log.Printf("获取jumpserver label配置id失败")
		return
	}
	deleteAssetByLabelIDError := DeleteAssetByLabelID(config.Jumpserver.URL, jpInstanceId, token, jpLabelId)
	err := DeleteAssetByID(config.Jumpserver.URL, asset.ID, token)
	if err != nil && deleteAssetByLabelIDError != nil {
		log.Printf("Error deleting asset with ID %s: %v \n", asset.ID, err)
	} else {
		log.Printf("Successfully deleted asset with ID %s\n", asset.ID)
	}
}

func UpdateAssetJd2JumpServer(config *config.Config, token string, pks []string, asset Asset) (string, error) {
	//根据Instance获取pk
	assetData := map[string]interface{}{
		"name":    asset.Name,
		"address": asset.Address, // 假设使用第一个私有IP地址
		"platform": map[string]interface{}{
			"pk": asset.Platform.ID, // 根据实际情况设置平台的 pk
		},
		"nodes": []map[string]interface{}{
			//{
			//	"pk": pk, // 根据实际情况设置节点的 pk
			//},
		},
		"protocols": asset.Protocols,
		//"domain":    asset.Domain,
		//"is_active": asset.IsActive,
		//"comment":   asset.Comment,
		//"labels":    asset.Labels,
	}

	// 将 pks 组装进 nodes
	for _, pk := range pks {
		node := map[string]interface{}{
			"pk": pk,
		}
		assetData["nodes"] = append(assetData["nodes"].([]map[string]interface{}), node)
	}
	result, err := UpdateAsset(config.Jumpserver.URL, "/api/v1/assets/hosts/"+asset.ID+"/", token, assetData)
	if err != nil {
		log.Printf("Error creating asset for instance %s: %v\n", asset.Address, err)
		if err != nil {
			return "", err
		}
	}
	return result, err
}
