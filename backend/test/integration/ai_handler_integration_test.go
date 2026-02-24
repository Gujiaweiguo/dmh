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

type AIHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
}

func (suite *AIHandlerIntegrationTestSuite) SetupSuite() {
	baseURL := os.Getenv("DMH_INTEGRATION_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8889"
	}
	suite.baseURL = strings.TrimRight(baseURL, "/")
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
}

func (suite *AIHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}) (int, []byte) {
	var reqBody io.Reader
	if payload != nil {
		body, _ := json.Marshal(payload)
		reqBody = bytes.NewBuffer(body)
	}

	req, _ := http.NewRequest(method, suite.baseURL+path, reqBody)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

func (suite *AIHandlerIntegrationTestSuite) Test_1_GenerateCopywriting() {
	payload := map[string]interface{}{
		"topic":  "双十一大促",
		"style":  "professional",
		"length": "medium",
	}
	status, body := suite.doRequest(http.MethodPost, "/api/v1/ai/generate-copywriting", payload)
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Content string `json:"content"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.NotEmpty(resp.Content)
	suite.Contains(resp.Content, "双十一大促")
	suite.T().Logf("✓ AI 文案生成成功: %s", resp.Content)
}

func (suite *AIHandlerIntegrationTestSuite) Test_2_GenerateCopywritingCasualStyle() {
	payload := map[string]interface{}{
		"topic":  "新品上市",
		"style":  "casual",
		"length": "short",
	}
	status, body := suite.doRequest(http.MethodPost, "/api/v1/ai/generate-copywriting", payload)
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Content string `json:"content"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Contains(resp.Content, "🎉")
	suite.T().Logf("✓ 轻松风格文案生成成功")
}

func (suite *AIHandlerIntegrationTestSuite) Test_3_GenerateCopywritingUrgentStyle() {
	payload := map[string]interface{}{
		"topic":  "限时抢购",
		"style":  "urgent",
		"length": "medium",
	}
	status, body := suite.doRequest(http.MethodPost, "/api/v1/ai/generate-copywriting", payload)
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Content string `json:"content"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Contains(resp.Content, "⚡")
	suite.T().Logf("✓ 紧急风格文案生成成功")
}

func (suite *AIHandlerIntegrationTestSuite) Test_4_GenerateCopywritingEmotionalStyle() {
	payload := map[string]interface{}{
		"topic":  "感恩回馈",
		"style":  "emotional",
		"length": "long",
	}
	status, body := suite.doRequest(http.MethodPost, "/api/v1/ai/generate-copywriting", payload)
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Content string `json:"content"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Contains(resp.Content, "❤️")
	suite.T().Logf("✓ 情感风格文案生成成功")
}

func (suite *AIHandlerIntegrationTestSuite) Test_5_GenerateCopywritingDefaultStyle() {
	payload := map[string]interface{}{
		"topic": "默认风格测试",
	}
	status, body := suite.doRequest(http.MethodPost, "/api/v1/ai/generate-copywriting", payload)
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Content string `json:"content"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Contains(resp.Content, "📢")
	suite.T().Logf("✓ 默认风格文案生成成功")
}

func (suite *AIHandlerIntegrationTestSuite) Test_6_GenerateCopywritingLongLength() {
	payload := map[string]interface{}{
		"topic":  "长文案测试",
		"style":  "professional",
		"length": "long",
	}
	status, body := suite.doRequest(http.MethodPost, "/api/v1/ai/generate-copywriting", payload)
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Content string `json:"content"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Contains(resp.Content, "卓越")
	suite.T().Logf("✓ 长文案生成成功")
}

func (suite *AIHandlerIntegrationTestSuite) Test_7_GenerateCopywritingEmptyTopic() {
	payload := map[string]interface{}{
		"style": "professional",
	}
	status, _ := suite.doRequest(http.MethodPost, "/api/v1/ai/generate-copywriting", payload)
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 空 topic 文案生成被拒绝，状态码: %d", status)
}

func TestAIHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AIHandlerIntegrationTestSuite))
}
