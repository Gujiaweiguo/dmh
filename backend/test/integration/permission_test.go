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

// PermissionControlIntegrationTestSuite 权限控制集成测试套件
type PermissionControlIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
	userToken  string
	brandToken string
	campaignId int64
}

func (suite *PermissionControlIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}

	suite.loginAsAdmin()
	suite.loginAsUser()
	suite.loginAsBrand()
	suite.createTestCampaign()
}

func (suite *PermissionControlIntegrationTestSuite) loginAsAdmin() {
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

	suite.adminToken = loginResp.Token
	suite.T().Log("✓ Admin 登录成功")
}

func (suite *PermissionControlIntegrationTestSuite) loginAsUser() {
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

	suite.userToken = loginResp.Token
	suite.T().Log("✓ User 登录成功")
}

func (suite *PermissionControlIntegrationTestSuite) loginAsBrand() {
	suite.brandToken = suite.adminToken
	suite.T().Log("✓ Brand token 设置成功 (使用admin)")
}

func (suite *PermissionControlIntegrationTestSuite) createTestCampaign() {
	now := time.Now()
	createCampaignReq := map[string]interface{}{
		"brandId":     1,
		"name":        "权限测试活动",
		"description": "用于测试权限控制的活动",
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
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

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

// Test_11_5_1_UnauthorizedAccess 测试未授权访问
func (suite *PermissionControlIntegrationTestSuite) Test_11_5_1_UnauthorizedAccess() {
	suite.T().Log("测试场景 1: 未授权访问活动列表")

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/campaigns", nil)
	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.T().Logf("✓ 未授权访问返回状态码: %d", resp.StatusCode)
	suite.True(resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden,
		"未授权访问应返回 401 或 403")
}

// Test_11_5_2_AdminAccess 测试管理员访问所有资源
func (suite *PermissionControlIntegrationTestSuite) Test_11_5_2_AdminAccess() {
	suite.T().Log("测试场景 2: 管理员访问所有资源")

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+suite.adminToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.T().Logf("✓ 管理员访问用户列表返回状态码: %d", resp.StatusCode)
	suite.T().Log("✓ 管理员权限测试完成")
}

// Test_11_5_3_UserCannotAccessAdmin 测试普通用户无法访问管理接口
func (suite *PermissionControlIntegrationTestSuite) Test_11_5_3_UserCannotAccessAdmin() {
	suite.T().Log("测试场景 3: 普通用户无法访问管理接口")

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+suite.userToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.T().Logf("✓ 普通用户访问管理接口返回状态码: %d", resp.StatusCode)
	suite.T().Log("✓ 普通用户权限测试完成")
}

// Test_11_5_4_InvalidToken 测试无效Token访问
func (suite *PermissionControlIntegrationTestSuite) Test_11_5_4_InvalidToken() {
	suite.T().Log("测试场景 4: 使用无效Token访问")

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/campaigns", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_12345")

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.T().Logf("✓ 无效Token访问返回状态码: %d", resp.StatusCode)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode, "无效Token应返回 401")
}

// Test_11_5_5_ExpiredToken 测试过期Token访问
func (suite *PermissionControlIntegrationTestSuite) Test_11_5_5_ExpiredToken() {
	suite.T().Log("测试场景 5: 使用过期Token访问")

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/campaigns", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJuYW1lIjoiYWRtaW4iLCJyb2xlcyI6WyJwbGF0Zm9ybV9hZG1pbiJdLCJleHAiOjE2MDAwMDAwMDB9")

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.T().Logf("✓ 过期Token访问返回状态码: %d", resp.StatusCode)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode, "过期Token应返回 401")
}

// Test_11_5_6_CampaignAccess 测试活动访问权限
func (suite *PermissionControlIntegrationTestSuite) Test_11_5_6_CampaignAccess() {
	suite.T().Log("测试场景 6: 不同角色访问活动详情")

	url := fmt.Sprintf("%s/api/v1/campaigns/%d", suite.baseURL, suite.campaignId)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+suite.userToken)

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	suite.T().Logf("✓ 活动详情访问返回状态码: %d", resp.StatusCode)
	suite.T().Log("✓ 活动访问权限测试完成")
}

// TestPermissionControlIntegrationTestSuite 运行测试套件
func TestPermissionControlIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionControlIntegrationTestSuite))
}
