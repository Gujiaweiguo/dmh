package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MaterialHandlerIntegrationTestSuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	adminToken string
}

func (suite *MaterialHandlerIntegrationTestSuite) SetupSuite() {
	baseURL := os.Getenv("DMH_INTEGRATION_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8889"
	}
	suite.baseURL = strings.TrimRight(baseURL, "/")
	suite.httpClient = &http.Client{Timeout: 10 * time.Second}
	suite.loginAsAdmin()
}

func (suite *MaterialHandlerIntegrationTestSuite) loginAsAdmin() {
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

func (suite *MaterialHandlerIntegrationTestSuite) doRequest(method, path string, payload interface{}, token string) (int, []byte) {
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

func (suite *MaterialHandlerIntegrationTestSuite) doMultipartUpload(path, fieldName, filename string, content []byte) (int, []byte) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	suite.Require().NoError(err)
	_, err = part.Write(content)
	suite.Require().NoError(err)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, respBody
}

func (suite *MaterialHandlerIntegrationTestSuite) Test_1_GetMaterialList() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/material/list?page=1&pageSize=10", nil, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Total     int64 `json:"total"`
		Materials []struct {
			Id   int64  `json:"id"`
			Type string `json:"type"`
		} `json:"materials"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.T().Logf("✓ 素材列表查询成功，共 %d 个素材", resp.Total)
}

func (suite *MaterialHandlerIntegrationTestSuite) Test_2_GetMaterialListWithFilters() {
	status, body := suite.doRequest(http.MethodGet, "/api/v1/material/list?type=image&page=1&pageSize=5", nil, "")
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Total     int64 `json:"total"`
		Materials []struct {
			Type string `json:"type"`
		} `json:"materials"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	for _, m := range resp.Materials {
		suite.Equal("image", m.Type)
	}
	suite.T().Logf("✓ 带过滤条件的素材列表查询成功")
}

func (suite *MaterialHandlerIntegrationTestSuite) Test_3_UploadMaterial() {
	status, body := suite.doMultipartUpload("/api/v1/material/upload?type=image", "file", "test.png", []byte("fake image content"))
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Equal("test.png", resp.Name)
	suite.Equal("image", resp.Type)
	suite.T().Logf("✓ 素材上传成功，ID: %d", resp.Id)
}

func (suite *MaterialHandlerIntegrationTestSuite) Test_4_DeleteMaterialNotFound() {
	status, _ := suite.doRequest(http.MethodDelete, "/api/v1/material/delete/999999999", nil, "")
	suite.NotEqual(http.StatusOK, status)
	suite.T().Log("✓ 删除不存在的素材返回错误")
}

func (suite *MaterialHandlerIntegrationTestSuite) Test_5_InvalidMaterialId() {
	status, _ := suite.doRequest(http.MethodDelete, "/api/v1/material/delete/abc", nil, "")
	suite.NotEqual(http.StatusOK, status)
	suite.T().Log("✓ 无效的素材ID被拒绝")
}

func (suite *MaterialHandlerIntegrationTestSuite) Test_6_UploadTextMaterial() {
	status, body := suite.doMultipartUpload("/api/v1/material/upload?type=text", "file", "test.txt", []byte("sample text content"))
	suite.Equal(http.StatusOK, status)

	var resp struct {
		Id   int64  `json:"id"`
		Type string `json:"type"`
	}
	err := json.Unmarshal(body, &resp)
	suite.NoError(err)
	suite.Equal("text", resp.Type)
	suite.T().Logf("✓ 文本素材上传成功，ID: %d", resp.Id)
}

func TestMaterialHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(MaterialHandlerIntegrationTestSuite))
}
