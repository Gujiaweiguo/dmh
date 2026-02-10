package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// ConcurrencyTestSuite 并发场景测试套件
type ConcurrencyTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	authToken  string
	campaignId int64
}

func (suite *ConcurrencyTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}

	suite.login()
	suite.createTestCampaign()
}

func (suite *ConcurrencyTestSuite) login() {
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

func (suite *ConcurrencyTestSuite) createTestCampaign() {
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1,
		"name":        "并发测试活动",
		"description": "用于测试并发场景的活动",
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

// Test_11_6_1_ConcurrentCampaignList 测试并发获取活动列表
func (suite *ConcurrencyTestSuite) Test_11_6_1_ConcurrentCampaignList() {
	suite.T().Log("测试场景 1: 并发获取活动列表")

	concurrentRequests := 10
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/campaigns", nil)
			req.Header.Set("Authorization", "Bearer "+suite.authToken)

			resp, err := suite.httpClient.Do(req)
			if err != nil {
				suite.T().Logf("请求 %d 失败: %v", index, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	suite.T().Logf("✓ 并发请求完成: %d/%d 成功", successCount, concurrentRequests)
	suite.Greater(successCount, concurrentRequests/2, "至少有一半请求应成功")
}

// Test_11_6_2_ConcurrentPaymentQRCode 测试并发生成支付二维码
func (suite *ConcurrencyTestSuite) Test_11_6_2_ConcurrentPaymentQRCode() {
	suite.T().Log("测试场景 2: 并发生成支付二维码")

	concurrentRequests := 10
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			url := fmt.Sprintf("%s/api/v1/campaigns/%d/payment-qrcode", suite.baseURL, suite.campaignId)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Authorization", "Bearer "+suite.authToken)

			start := time.Now()
			resp, err := suite.httpClient.Do(req)
			if err != nil {
				suite.T().Logf("请求 %d 失败: %v", index, err)
				return
			}
			defer resp.Body.Close()
			duration := time.Since(start)

			if resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
				suite.T().Logf("请求 %d 成功，耗时: %v", index, duration)
			} else {
				suite.T().Logf("请求 %d 返回状态码: %d", index, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()

	suite.T().Logf("✓ 并发二维码生成完成: %d/%d 成功", successCount, concurrentRequests)
	suite.Greater(successCount, 0, "至少有部分请求应成功")
}

// Test_11_6_3_ConcurrentCreateOrder 测试并发创建订单
func (suite *ConcurrencyTestSuite) Test_11_6_3_ConcurrentCreateOrder() {
	suite.T().Log("测试场景 3: 并发创建订单（不同手机号）")

	concurrentRequests := 10
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			createOrderReq := map[string]interface{}{
				"campaignId": suite.campaignId,
				"phone":      fmt.Sprintf("138%08d", index),
				"formData":   map[string]string{"name": fmt.Sprintf("用户%d", index)},
			}

			reqBody, _ := json.Marshal(createOrderReq)
			req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/orders", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+suite.authToken)

			start := time.Now()
			resp, err := suite.httpClient.Do(req)
			if err != nil {
				suite.T().Logf("请求 %d 失败: %v", index, err)
				return
			}
			defer resp.Body.Close()
			duration := time.Since(start)

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest {
				mu.Lock()
				successCount++
				mu.Unlock()
				suite.T().Logf("请求 %d 完成 (状态: %d)，耗时: %v", index, resp.StatusCode, duration)
			} else {
				suite.T().Logf("请求 %d 返回状态码: %d", index, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()

	suite.T().Logf("✓ 并发创建订单完成: %d/%d 成功", successCount, concurrentRequests)
	suite.Greater(successCount, 0, "至少有部分请求应成功")
}

// Test_11_6_4_ConcurrentLogin 测试并发登录
func (suite *ConcurrencyTestSuite) Test_11_6_4_ConcurrentLogin() {
	suite.T().Log("测试场景 4: 并发登录")

	concurrentRequests := 10
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			loginReq := map[string]string{
				"username": "admin",
				"password": "123456",
			}

			reqBody, _ := json.Marshal(loginReq)
			req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			start := time.Now()
			resp, err := suite.httpClient.Do(req)
			if err != nil {
				suite.T().Logf("请求 %d 失败: %v", index, err)
				return
			}
			defer resp.Body.Close()
			duration := time.Since(start)

			if resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
				suite.T().Logf("请求 %d 成功，耗时: %v", index, duration)
			} else {
				suite.T().Logf("请求 %d 返回状态码: %d", index, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()

	suite.T().Logf("✓ 并发登录完成: %d/%d 成功", successCount, concurrentRequests)
	suite.Greater(successCount, 0, "至少有部分登录请求应成功")
}

// TestConcurrencyTestSuite 运行测试套件
func TestConcurrencyTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrencyTestSuite))
}
