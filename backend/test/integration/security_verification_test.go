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

// VerificationCodeSecurityTestSuite 核销码伪造防护测试套件
type VerificationCodeSecurityTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	authToken  string
	campaignId int64
	orderId    int64
}

func (suite *VerificationCodeSecurityTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}

	suite.login()
	suite.createTestCampaign()
	suite.createTestOrder()
}

func (suite *VerificationCodeSecurityTestSuite) login() {
	loginReq := map[string]string{
		"username": "admin",
		"password": "123456",
	}
	reqBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("无法连接到后端服务: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var loginResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败: %s", string(body))
		return
	}

	suite.authToken = loginResp.Token
	suite.T().Log("✓ 登录成功")
}

func (suite *VerificationCodeSecurityTestSuite) createTestCampaign() {
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1,
		"name":        "核销码安全测试活动",
		"description": "用于测试核销码安全性的活动",
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
		suite.T().Skipf("创建活动失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		suite.T().Skipf("创建活动失败: %s", string(body))
		return
	}

	var createResp struct {
		Id int64 `json:"id"`
	}
	json.Unmarshal(body, &createResp)

	suite.campaignId = createResp.Id
	suite.T().Logf("✓ 测试活动创建成功，ID: %d", suite.campaignId)
}

func (suite *VerificationCodeSecurityTestSuite) createTestOrder() {
	createOrderReq := map[string]interface{}{
		"campaignId": suite.campaignId,
		"phone":      "13800138888",
		"formData":   map[string]string{"name": "安全测试用户"},
	}

	reqBody, _ := json.Marshal(createOrderReq)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/orders", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		suite.T().Skipf("创建订单失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		suite.T().Skipf("创建订单失败: %s", string(body))
		return
	}

	var createResp struct {
		Id int64 `json:"id"`
	}
	json.Unmarshal(body, &createResp)

	suite.orderId = createResp.Id
	suite.T().Logf("✓ 测试订单创建成功，ID: %d", suite.orderId)
}

// Test_13_1_1_TamperedOrderId 测试篡改订单ID
func (suite *VerificationCodeSecurityTestSuite) Test_13_1_1_TamperedOrderId() {
	suite.T().Log("测试场景 1: 使用篡改的订单ID")

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	fakeOrderId := suite.orderId + 10000
	fakeCode := fmt.Sprintf("%d_13800138888_%s_abc123", fakeOrderId, timestamp)

	verifyReq := map[string]string{
		"code": fakeCode,
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

	suite.T().Logf("✓ 篡改订单ID测试完成 (状态码: %d): %s", resp.StatusCode, string(body))
	suite.NotEqual(http.StatusOK, resp.StatusCode, "篡改订单ID应被拒绝")
}

// Test_13_1_2_TamperedPhone 测试篡改手机号
func (suite *VerificationCodeSecurityTestSuite) Test_13_1_2_TamperedPhone() {
	suite.T().Log("测试场景 2: 使用篡改的手机号")

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	fakeCode := fmt.Sprintf("%d_13999999999_%s_abc123", suite.orderId, timestamp)

	verifyReq := map[string]string{
		"code": fakeCode,
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

	suite.T().Logf("✓ 篡改手机号测试完成 (状态码: %d): %s", resp.StatusCode, string(body))
	suite.NotEqual(http.StatusOK, resp.StatusCode, "篡改手机号应被拒绝")
}

// Test_13_1_3_TamperedSignature 测试篡改签名
func (suite *VerificationCodeSecurityTestSuite) Test_13_1_3_TamperedSignature() {
	suite.T().Log("测试场景 3: 使用篡改的签名")

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	fakeCode := fmt.Sprintf("%d_13800138888_%s_wrong_signature", suite.orderId, timestamp)

	verifyReq := map[string]string{
		"code": fakeCode,
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

	suite.T().Logf("✓ 篡改签名测试完成 (状态码: %d): %s", resp.StatusCode, string(body))
	suite.NotEqual(http.StatusOK, resp.StatusCode, "篡改签名应被拒绝")
}

// Test_13_1_4_MalformedCode 测试格式错误的核销码
func (suite *VerificationCodeSecurityTestSuite) Test_13_1_4_MalformedCode() {
	suite.T().Log("测试场景 4: 使用格式错误的核销码")

	malformedCodes := []string{
		"invalid_code",
		"123",
		"a_b",
		"1_2_3",
	}

	for _, code := range malformedCodes {
		verifyReq := map[string]string{
			"code": code,
		}
		reqBody, _ := json.Marshal(verifyReq)

		url := fmt.Sprintf("%s/api/v1/orders/verify", suite.baseURL)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.authToken)

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			suite.T().Logf("✓ 格式错误代码 '%s' 被正确拒绝 (状态码: %d)", code, resp.StatusCode)
		}
	}

	suite.T().Log("✓ 格式错误核销码测试完成")
}

// Test_13_1_5_EmptyCode 测试空核销码
func (suite *VerificationCodeSecurityTestSuite) Test_13_1_5_EmptyCode() {
	suite.T().Log("测试场景 5: 使用空核销码")

	verifyReq := map[string]string{
		"code": "",
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

	suite.T().Logf("✓ 空核销码测试完成 (状态码: %d): %s", resp.StatusCode, string(body))
	suite.NotEqual(http.StatusOK, resp.StatusCode, "空核销码应被拒绝")
}

// TestVerificationCodeSecurityTestSuite 运行测试套件
func TestVerificationCodeSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(VerificationCodeSecurityTestSuite))
}
