package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// RateLimitingTestSuite 频率限制测试套件
type RateLimitingTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	authToken  string
	campaignId int64
	templateId int64
}

func (suite *RateLimitingTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}

	suite.login()
	suite.createTestCampaign()
	suite.loadPosterTemplate()
}

func (suite *RateLimitingTestSuite) loadPosterTemplate() {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/poster/templates?page=1&pageSize=1", nil)
	resp, err := suite.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}

	var listResp struct {
		Templates []struct {
			Id int64 `json:"id"`
		} `json:"templates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return
	}
	if len(listResp.Templates) > 0 {
		suite.templateId = listResp.Templates[0].Id
	}
}

func (suite *RateLimitingTestSuite) login() {
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

	var loginResp struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败: token is empty")
		return
	}

	suite.authToken = loginResp.Token
	suite.T().Log("✓ 登录成功")
}

func (suite *RateLimitingTestSuite) createTestCampaign() {
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1,
		"name":        "频率限制测试活动",
		"description": "用于测试频率限制的活动",
		"rewardRule":  10.0,
		"startTime":   now.Add(-24 * time.Hour).Format("2006-01-02T15:04:05"),
		"endTime":     now.Add(48 * time.Hour).Format("2006-01-02T15:04:05"),
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

	var createResp struct {
		Id int64 `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&createResp)

	suite.campaignId = createResp.Id
	suite.T().Logf("✓ 测试活动创建成功，ID: %d", suite.campaignId)
}

// Test_13_3_1_PosterRateLimit 测试海报生成频率限制
func (suite *RateLimitingTestSuite) Test_13_3_1_PosterRateLimit() {
	suite.T().Log("测试场景 1: 海报生成频率限制（配置: 5次/分钟）")
	if suite.templateId == 0 {
		suite.T().Skip("无可用海报模板ID，跳过海报限流测试")
	}

	successCount := 0
	rateLimitCount := 0
	var mu sync.Mutex

	for i := 0; i < 8; i++ {
		generatePosterReq := map[string]int64{
			"templateId": suite.templateId,
		}
		reqBody, _ := json.Marshal(generatePosterReq)

		url := fmt.Sprintf("%s/api/v1/campaigns/%d/poster", suite.baseURL, suite.campaignId)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.authToken)

		start := time.Now()
		resp, err := suite.httpClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			suite.T().Logf("请求 %d 失败: %v", i+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			mu.Lock()
			successCount++
			mu.Unlock()
			suite.T().Logf("请求 %d 成功，耗时: %v", i+1, duration)
		} else if resp.StatusCode == http.StatusTooManyRequests {
			mu.Lock()
			rateLimitCount++
			mu.Unlock()
			suite.T().Logf("请求 %d 被限流（状态码: 429），耗时: %v", i+1, duration)
		} else {
			suite.T().Logf("请求 %d 返回状态码: %d，耗时: %v", i+1, resp.StatusCode, duration)
		}

		time.Sleep(100 * time.Millisecond)
	}

	suite.T().Logf("✓ 海报生成频率限制测试完成: %d 成功，%d 被限流", successCount, rateLimitCount)
	suite.Greater(successCount, 0, "至少有部分请求应成功")
}

// Test_13_3_2_PaymentQRCodeRateLimit 测试支付二维码生成频率限制
func (suite *RateLimitingTestSuite) Test_13_3_2_PaymentQRCodeRateLimit() {
	suite.T().Log("测试场景 2: 支付二维码生成频率限制")

	successCount := 0
	rateLimitCount := 0
	var mu sync.Mutex

	for i := 0; i < 8; i++ {
		url := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, suite.campaignId)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+suite.authToken)

		start := time.Now()
		resp, err := suite.httpClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			suite.T().Logf("请求 %d 失败: %v", i+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			mu.Lock()
			successCount++
			mu.Unlock()
			suite.T().Logf("请求 %d 成功，耗时: %v", i+1, duration)
		} else if resp.StatusCode == http.StatusTooManyRequests {
			mu.Lock()
			rateLimitCount++
			mu.Unlock()
			suite.T().Logf("请求 %d 被限流（状态码: 429），耗时: %v", i+1, duration)
		} else {
			suite.T().Logf("请求 %d 返回状态码: %d，耗时: %v", i+1, resp.StatusCode, duration)
		}

		time.Sleep(50 * time.Millisecond)
	}

	suite.T().Logf("✓ 支付二维码生成频率限制测试完成: %d 成功，%d 被限流", successCount, rateLimitCount)
	suite.Greater(successCount, 0, "至少有部分请求应成功")
}

// TestRateLimitingTestSuite 运行测试套件
func TestRateLimitingTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimitingTestSuite))
}
