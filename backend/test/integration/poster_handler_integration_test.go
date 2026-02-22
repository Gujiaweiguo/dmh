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

type PosterHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL            string
	httpClient         *http.Client
	adminToken         string
	existingCampaignID int64
	existingDistID     int64
	templateID         int64
}

func (suite *PosterHandlerIntegrationTestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8889"
	suite.httpClient = &http.Client{Timeout: 20 * time.Second}
	suite.loginAsAdmin()
	suite.loadExistingCampaignID()
	suite.loadExistingDistributorID()
}

func (suite *PosterHandlerIntegrationTestSuite) loginAsAdmin() {
	loginReq := map[string]string{
		"username": "admin",
		"password": "123456",
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

func (suite *PosterHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, token string) (int, []byte) {
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

func (suite *PosterHandlerIntegrationTestSuite) loadExistingCampaignID() {
	if suite.adminToken == "" {
		return
	}

	status, body := suite.doRequest(http.MethodGet, "/api/v1/campaigns?page=1&pageSize=1", nil, suite.adminToken)
	if status != http.StatusOK {
		return
	}

	var resp struct {
		Campaigns []struct {
			Id int64 `json:"id"`
		} `json:"campaigns"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return
	}

	if len(resp.Campaigns) > 0 {
		suite.existingCampaignID = resp.Campaigns[0].Id
	}
}

func (suite *PosterHandlerIntegrationTestSuite) loadExistingDistributorID() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/poster/records", nil, "")
	if status != http.StatusOK {
		return
	}

	var resp struct {
		Records []struct {
			RecordType    string `json:"recordType"`
			DistributorId int64  `json:"distributorId"`
		} `json:"records"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return
	}

	for _, record := range resp.Records {
		if record.RecordType == "distributor" && record.DistributorId > 0 {
			suite.existingDistID = record.DistributorId
			return
		}
	}
}

func (suite *PosterHandlerIntegrationTestSuite) Test_1_GetPosterTemplatesPublic() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/poster/templates?page=1&pageSize=10", nil, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Total     int64 `json:"total"`
		Templates []struct {
			Id int64 `json:"id"`
		} `json:"templates"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	if len(resp.Templates) > 0 {
		suite.templateID = resp.Templates[0].Id
	}
	suite.T().Logf("✓ 海报模板列表查询成功，total=%d", resp.Total)
}

func (suite *PosterHandlerIntegrationTestSuite) Test_2_GetPosterRecordsPublic() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/poster/records", nil, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Total int64 `json:"total"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.T().Logf("✓ 海报记录列表查询成功，total=%d", resp.Total)
}

func (suite *PosterHandlerIntegrationTestSuite) Test_3_GenerateCampaignPoster() {
	if suite.existingCampaignID == 0 {
		suite.T().Skip("无可用活动ID，跳过测试")
	}
	if suite.templateID == 0 {
		suite.T().Skip("无可用海报模板ID，跳过测试")
	}

	path := fmt.Sprintf("/api/v1/campaigns/%d/poster", suite.existingCampaignID)
	status, body := suite.doRequest(http.MethodPost, path, map[string]interface{}{"templateId": suite.templateID}, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		PosterUrl string `json:"posterUrl"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.NotEmpty(resp.PosterUrl)
	suite.T().Logf("✓ 活动海报生成成功，campaignId=%d", suite.existingCampaignID)
}

func (suite *PosterHandlerIntegrationTestSuite) Test_4_GenerateCampaignPosterNotFound() {
	status, body := suite.doRequest(http.MethodPost, "/api/v1/campaigns/999999/poster", map[string]interface{}{"templateId": 1}, "")
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 不存在活动海报生成被拒绝，状态码: %d，响应: %s", status, string(body))
}

func (suite *PosterHandlerIntegrationTestSuite) Test_5_GenerateDistributorPoster() {
	if suite.existingDistID == 0 {
		suite.T().Skip("无可用分销商ID，跳过测试")
	}

	path := fmt.Sprintf("/api/v1/distributors/%d/poster", suite.existingDistID)
	status, body := suite.doRequest(http.MethodPost, path, map[string]interface{}{"templateId": 1}, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		PosterUrl string `json:"posterUrl"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.NotEmpty(resp.PosterUrl)
	suite.T().Logf("✓ 分销海报生成成功，distributorId=%d", suite.existingDistID)
}

func (suite *PosterHandlerIntegrationTestSuite) Test_6_GenerateDistributorPosterNotFound() {
	status, body := suite.doRequest(http.MethodPost, "/api/v1/distributors/999999/poster", map[string]interface{}{"templateId": 1}, "")
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 不存在分销商海报生成被拒绝，状态码: %d，响应: %s", status, string(body))
}

func (suite *PosterHandlerIntegrationTestSuite) Test_7_GenerateDistributorPosterInvalidId() {
	status, body := suite.doRequest(http.MethodPost, "/api/v1/distributors/abc/poster", map[string]interface{}{"templateId": 1}, "")
	suite.NotEqual(http.StatusOK, status)
	suite.T().Logf("✓ 非法分销商ID被拒绝，状态码: %d，响应: %s", status, string(body))
}

func TestPosterHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PosterHandlerIntegrationTestSuite))
}
