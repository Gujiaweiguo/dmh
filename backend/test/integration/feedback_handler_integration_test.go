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

type FeedbackHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
	feedbackID int64
	faqID      int64
}

func (suite *FeedbackHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
}

func (suite *FeedbackHandlerIntegrationTestSuite) loginAsAdmin() {
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
	_ = json.Unmarshal(body, &loginResp)

	if loginResp.Token == "" {
		suite.T().Skipf("登录失败: %s", string(body))
		return
	}

	suite.adminToken = loginResp.Token
	suite.T().Log("✓ Admin 登录成功")
}

func (suite *FeedbackHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, token string) (int, []byte) {
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

func (suite *FeedbackHandlerIntegrationTestSuite) createPublicFeedback() int64 {
	payload := map[string]interface{}{
		"category": "other",
		"title":    fmt.Sprintf("集成测试反馈_%d", time.Now().UnixNano()),
		"content":  "用于集成测试",
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/feedback", payload, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Id int64 `json:"id"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Greater(resp.Id, int64(0))
	return resp.Id
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_1_CreateFeedbackPublic() {
	status, body := suite.doRequest(http.MethodPost, "/api/v1/feedback", map[string]interface{}{
		"category": "other",
		"title":    fmt.Sprintf("公开反馈_%d", time.Now().UnixNano()),
		"content":  "公开反馈内容",
	}, "")

	suite.Equal(http.StatusOK, status)

	var resp struct {
		Id int64 `json:"id"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Greater(resp.Id, int64(0))
	suite.feedbackID = resp.Id
	suite.T().Logf("✓ 公开反馈创建成功，ID: %d", resp.Id)
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_2_ListFeedbackPublic() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/feedback/list?page=1&pageSize=10", nil, "")
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 反馈列表查询成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_3_ListFAQPublic() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/feedback/faq", nil, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Faqs []struct {
			Id int64 `json:"id"`
		} `json:"faqs"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	if len(resp.Faqs) > 0 {
		suite.faqID = resp.Faqs[0].Id
	}
	suite.T().Logf("✓ FAQ 列表查询成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_4_MarkFAQHelpful() {
	if suite.faqID == 0 {
		status, body := suite.doRequest(http.MethodGet, "/api/v1/feedback/faq", nil, "")
		suite.Equal(http.StatusOK, status)
		var resp struct {
			Faqs []struct {
				Id int64 `json:"id"`
			} `json:"faqs"`
		}
		err := json.Unmarshal(body, &resp)
		suite.NoError(err)
		if len(resp.Faqs) == 0 {
			suite.T().Skip("无可用 FAQ，跳过 helpful 测试")
		}
		suite.faqID = resp.Faqs[0].Id
	}

	status, body := suite.doRequest(http.MethodPost, "/api/v1/feedback/faq/helpful", map[string]interface{}{
		"id":   suite.faqID,
		"type": "helpful",
	}, "")

	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ FAQ helpful 请求完成，状态码: %d，响应: %s", status, string(body))
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_5_SubmitSatisfactionSurvey() {
	status, body := suite.doRequest(http.MethodPost, "/api/v1/feedback/satisfaction-survey", map[string]interface{}{
		"feature":             "poster",
		"overallSatisfaction": 5,
	}, "")

	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 满意度调查提交成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_6_RecordFeatureUsage() {
	status, body := suite.doRequest(http.MethodPost, "/api/v1/feedback/feature-usage", map[string]interface{}{
		"feature": "poster",
		"action":  "generate",
		"success": true,
	}, "")

	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 功能使用记录提交成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_7_GetFeedbackDetailAdmin() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}
	if suite.feedbackID == 0 {
		suite.feedbackID = suite.createPublicFeedback()
	}

	status, body := suite.doRequest(http.MethodGet, fmt.Sprintf("/api/v1/feedback/detail?id=%d", suite.feedbackID), nil, suite.adminToken)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 管理员查询反馈详情成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_8_UpdateFeedbackStatusAdmin() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}
	if suite.feedbackID == 0 {
		suite.feedbackID = suite.createPublicFeedback()
	}

	status, body := suite.doRequest(http.MethodPut, "/api/v1/feedback/status", map[string]interface{}{
		"id":       suite.feedbackID,
		"status":   "resolved",
		"response": "done",
	}, suite.adminToken)

	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 管理员更新反馈状态成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_9_GetFeedbackStatisticsAdmin() {
	if suite.adminToken == "" {
		suite.T().Skip("未登录，跳过测试")
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/feedback/statistics", nil, suite.adminToken)
	suite.Equal(http.StatusOK, status)
	suite.T().Logf("✓ 管理员获取反馈统计成功，状态码: %d，响应: %s", status, string(body))
}

func (suite *FeedbackHandlerIntegrationTestSuite) Test_10_UnauthorizedAdminEndpoints() {
	status1, body1 := suite.doRequest(http.MethodGet, "/api/v1/feedback/detail?id=1", nil, "")
	suite.Equal(http.StatusUnauthorized, status1)
	suite.T().Logf("✓ 未授权访问 feedback/detail 被拒绝，状态码: %d，响应: %s", status1, string(body1))

	status2, body2 := suite.doRequest(http.MethodGet, "/api/v1/feedback/statistics", nil, "")
	suite.Equal(http.StatusUnauthorized, status2)
	suite.T().Logf("✓ 未授权访问 feedback/statistics 被拒绝，状态码: %d，响应: %s", status2, string(body2))
}

func TestFeedbackHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FeedbackHandlerIntegrationTestSuite))
}
