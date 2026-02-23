package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type PromoterHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
	brandId    int64
	campaignId int64
}

func (suite *PromoterHandlerIntegrationTestSuite) SetupSuite() {
	baseURL := os.Getenv("DMH_INTEGRATION_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8889"
	}
	suite.baseURL = strings.TrimRight(baseURL, "/")
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
	suite.ensureTestData()
}

func (suite *PromoterHandlerIntegrationTestSuite) loginAsAdmin() {
	username := os.Getenv("DMH_TEST_ADMIN_USERNAME")
	if username == "" {
		username = "admin"
	}
	password := os.Getenv("DMH_TEST_ADMIN_PASSWORD")
	if password == "" {
		password = "123456"
	}

	loginReq := map[string]string{
		"username": username,
		"password": password,
	}
	reqBody, _ := json.Marshal(loginReq)

	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
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
	_ = json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败: %s", string(body))
		return
	}

	suite.adminToken = loginResp.Token
}

func (suite *PromoterHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, token string) (int, []byte) {
	var reqBody io.Reader
	if payload != nil {
		body, _ := json.Marshal(payload)
		reqBody = bytes.NewBuffer(body)
	}

	req, _ := http.NewRequest(method, suite.baseURL+path, reqBody)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

func (suite *PromoterHandlerIntegrationTestSuite) ensureTestData() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/brands?page=1&pageSize=1", nil, suite.adminToken)
	if status == http.StatusOK {
		var listResp struct {
			Brands []struct {
				Id int64 `json:"id"`
			} `json:"brands"`
		}
		if err := json.Unmarshal(body, &listResp); err == nil && len(listResp.Brands) > 0 {
			suite.brandId = listResp.Brands[0].Id
		}
	}

	status, body = suite.doRequest(http.MethodGet, "/api/v1/campaigns?page=1&pageSize=1", nil, suite.adminToken)
	if status == http.StatusOK {
		var listResp struct {
			Campaigns []struct {
				Id int64 `json:"id"`
			} `json:"campaigns"`
		}
		if err := json.Unmarshal(body, &listResp); err == nil && len(listResp.Campaigns) > 0 {
			suite.campaignId = listResp.Campaigns[0].Id
		}
	}
}

func (suite *PromoterHandlerIntegrationTestSuite) Test_1_GetPromoterList() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/promoter/list?page=1&pageSize=10", nil, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Total     int64 `json:"total"`
		Promoters []struct {
			Id     int64  `json:"id"`
			Status string `json:"status"`
		} `json:"promoters"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.T().Logf("✓ 推广员列表查询成功，共 %d 个推广员", resp.Total)
}

func (suite *PromoterHandlerIntegrationTestSuite) Test_2_GetPromoterListWithFilters() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/promoter/list?status=active&page=1&pageSize=5", nil, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Total     int64 `json:"total"`
		Promoters []struct {
			Status string `json:"status"`
		} `json:"promoters"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	for _, p := range resp.Promoters {
		suite.Equal("active", p.Status)
	}
	suite.T().Logf("✓ 带过滤条件的推广员列表查询成功")
}

func (suite *PromoterHandlerIntegrationTestSuite) Test_3_GetPromoterDetailNotFound() {
	status, _ := suite.doRequest(http.MethodGet, "/api/v1/promoter/detail/999999999", nil, "")
	suite.NotEqual(http.StatusOK, status)
	suite.T().Log("✓ 查询不存在的推广员返回错误")
}

func (suite *PromoterHandlerIntegrationTestSuite) Test_4_GeneratePromoterLinkMissingParams() {
	payload := map[string]interface{}{}
	status, body := suite.doRequest(http.MethodPost, "/api/v1/promoter/generate-link", payload, "")
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 缺少参数的生成链接请求被拒绝，状态码: %d，响应: %s", status, string(body))
}

func (suite *PromoterHandlerIntegrationTestSuite) Test_5_GetPromoterRewardsNotFound() {
	status, _ := suite.doRequest(http.MethodGet, "/api/v1/promoter/rewards/999999999?page=1&pageSize=10", nil, "")
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	suite.Contains(allowedStatus, status)
	suite.T().Logf("✓ 查询不存在推广员的奖励记录，状态码: %d", status)
}

func (suite *PromoterHandlerIntegrationTestSuite) Test_6_GetPromoterRewardsWithFilters() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/promoter/rewards/1?type=commission&status=paid&page=1&pageSize=10", nil, "")
	allowedStatus := []int{http.StatusOK, http.StatusBadRequest}
	suite.Contains(allowedStatus, status)

	var resp struct {
		Total   int64 `json:"total"`
		Rewards []struct {
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"rewards"`
	}
	_ = json.Unmarshal(body, &resp)
	suite.T().Logf("✓ 带过滤条件的奖励记录查询，状态码: %d，总数: %d", status, resp.Total)
}

func (suite *PromoterHandlerIntegrationTestSuite) Test_7_InvalidPromoterId() {
	status, _ := suite.doRequest(http.MethodGet, "/api/v1/promoter/detail/abc", nil, "")
	suite.NotEqual(http.StatusOK, status)
	suite.T().Log("✓ 无效的推广员ID被拒绝")
}

func TestPromoterHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PromoterHandlerIntegrationTestSuite))
}
