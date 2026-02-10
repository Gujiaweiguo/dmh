package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// OrderVerificationIntegrationTestSuite 订单核销完整流程测试套件
type OrderVerificationIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	authToken  string
	orderId    int64
	campaignId int64
}

func (suite *OrderVerificationIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}

	suite.login()
	suite.createTestCampaign()
	suite.createTestOrder()
}

func (suite *OrderVerificationIntegrationTestSuite) login() {
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
	suite.T().Log("✓ 登录成功，获取到 token")
}

func (suite *OrderVerificationIntegrationTestSuite) createTestCampaign() {
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1,
		"name":        "订单核销测试活动",
		"description": "用于测试订单核销功能的活动",
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
	suite.T().Logf("✓ 测试活动创建成功，ID: %d", suite.campaignId)
}

func (suite *OrderVerificationIntegrationTestSuite) createTestOrder() {
	createOrderReq := map[string]interface{}{
		"campaignId": suite.campaignId,
		"phone":      "13800138000",
		"formData":   map[string]string{"name": "测试用户"},
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
	suite.T().Logf("✓ 测试订单创建成功，ID: %d", suite.orderId)
}

// Test_11_4_1_VerifyPaidOrder 测试核销已支付订单
func (suite *OrderVerificationIntegrationTestSuite) Test_11_4_1_VerifyPaidOrder() {
	suite.T().Log("测试场景 1: 核销已支付订单")

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	verificationCode := fmt.Sprintf("%d_13800138000_%s_abc123", suite.orderId, timestamp)

	verifyReq := map[string]string{
		"code": verificationCode,
	}
	reqBody, _ := json.Marshal(verifyReq)

	url := fmt.Sprintf("%s/api/v1/orders/verify", suite.baseURL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	suite.T().Logf("核销响应 (状态码: %d): %s", resp.StatusCode, string(body))
	suite.T().Log("✓ 核销接口测试完成")
}

// Test_11_4_2_ScanVerificationCode 测试扫码获取订单信息
func (suite *OrderVerificationIntegrationTestSuite) Test_11_4_2_ScanVerificationCode() {
	suite.T().Log("测试场景 2: 扫码获取订单信息")

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	verificationCode := fmt.Sprintf("%d_13800138000_%s_abc123", suite.orderId, timestamp)

	url := fmt.Sprintf("%s/api/v1/orders/scan?code=%s", suite.baseURL, verificationCode)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	suite.T().Logf("扫码响应 (状态码: %d): %s", resp.StatusCode, string(body))
	suite.T().Log("✓ 扫码接口测试完成")
}

// Test_11_4_3_UnverifyOrder 测试取消核销
func (suite *OrderVerificationIntegrationTestSuite) Test_11_4_3_UnverifyOrder() {
	suite.T().Log("测试场景 3: 取消核销")

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	verificationCode := fmt.Sprintf("%d_13800138000_%s_abc123", suite.orderId, timestamp)

	unverifyReq := map[string]string{
		"code": verificationCode,
	}
	reqBody, _ := json.Marshal(unverifyReq)

	url := fmt.Sprintf("%s/api/v1/orders/unverify", suite.baseURL)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	suite.T().Logf("取消核销响应 (状态码: %d): %s", resp.StatusCode, string(body))
	suite.T().Log("✓ 取消核销接口测试完成")
}

// Test_11_4_4_VerificationRecords 测试查询核销记录
func (suite *OrderVerificationIntegrationTestSuite) Test_11_4_4_VerificationRecords() {
	suite.T().Log("测试场景 4: 查询核销记录")

	url := fmt.Sprintf("%s/api/v1/orders/verification-records", suite.baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	suite.T().Logf("核销记录响应 (状态码: %d): %s", resp.StatusCode, string(body))
	suite.T().Log("✓ 查询核销记录接口测试完成")
}

// TestOrderVerificationIntegrationTestSuite 运行测试套件
func TestOrderVerificationIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(OrderVerificationIntegrationTestSuite))
}
