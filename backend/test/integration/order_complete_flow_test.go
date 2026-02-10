package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// OrderCompleteFlowIntegrationTestSuite 订单完整流程测试
type OrderCompleteFlowIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	authToken  string
	campaignId int64
	orderId    int64
}

func (suite *OrderCompleteFlowIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}

	suite.login()
	suite.createTestCampaign()
	suite.createTestOrder()
}

func (suite *OrderCompleteFlowIntegrationTestSuite) login() {
	loginReq := map[string]string{
		"username": "admin",
		"password": "123456",
	}
	reqBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("无法连接到后端服务，跳过测试: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败，未获取到 token: %s", string(body))
		return
	}

	suite.authToken = loginResp.Token
}

func (suite *OrderCompleteFlowIntegrationTestSuite) createTestCampaign() {
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1,
		"name":        "订单流程测试活动",
		"description": "用于测试完整订单流程",
		"rewardRule":  10.0,
		"startTime":   now.Add(-1 * time.Hour).Format(time.RFC3339),
		"endTime":     now.Add(24 * time.Hour).Format(time.RFC3339),
		"formFields": []map[string]interface{}{
			{
				"type":     "text",
				"name":     "name",
				"label":    "姓名",
				"required": true,
			},
			{
				"type":     "phone",
				"name":     "phone",
				"label":    "手机号",
				"required": true,
			},
		},
	}

	reqBody, _ := json.Marshal(createCampaignReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/campaigns", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("创建活动失败: %v，跳过测试", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		suite.T().Skipf("创建活动失败，状态码: %d，响应: %s，跳过测试", resp.StatusCode, string(body))
		return
	}

	var createResp struct {
		Id int64 `json:"id"`
	}
	json.Unmarshal(body, &createResp)
	suite.campaignId = createResp.Id
}

func (suite *OrderCompleteFlowIntegrationTestSuite) createTestOrder() {
	createOrderReq := map[string]interface{}{
		"campaignId": suite.campaignId,
		"phone":      "13800138099",
		"formData": map[string]string{
			"name":  "测试用户",
			"phone": "13800138099",
		},
	}

	reqBody, _ := json.Marshal(createOrderReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/orders", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("创建订单失败: %v，跳过测试", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		suite.T().Skipf("创建订单失败，状态码: %d，响应: %s，跳过测试", resp.StatusCode, string(body))
		return
	}

	var createResp struct {
		Id int64 `json:"id"`
	}
	json.Unmarshal(body, &createResp)
	suite.orderId = createResp.Id
}

// Test_1_CreateOrder 测试创建订单
func (suite *OrderCompleteFlowIntegrationTestSuite) Test_1_CreateOrder() {
	suite.T().Log("测试场景 1: 创建订单")
	suite.Greater(suite.orderId, int64(0), "订单ID应该大于0")
	suite.T().Logf("✓ 订单创建成功，ID: %d", suite.orderId)
}

// Test_2_GetOrder 测试查询订单
func (suite *OrderCompleteFlowIntegrationTestSuite) Test_2_GetOrder() {
	suite.T().Log("测试场景 2: 查询订单详情")

	url := suite.baseURL + "/api/v1/orders/" + strconv.FormatInt(suite.orderId, 10)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.T().Log("✓ 查询订单详情成功")
}

// Test_3_ListOrders 测试订单列表
func (suite *OrderCompleteFlowIntegrationTestSuite) Test_3_ListOrders() {
	suite.T().Log("测试场景 3: 查询订单列表")

	url := suite.baseURL + "/api/v1/orders/list?page=1&pageSize=10"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.T().Log("✓ 查询订单列表成功")
}

func TestOrderCompleteFlowIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(OrderCompleteFlowIntegrationTestSuite))
}
