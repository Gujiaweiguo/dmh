package integration

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

const (
	orderAuthGuardBaseURL      = "http://localhost:8889"
	verificationSignSecret     = "dmh-verification-secret-2026"
	defaultAdminUsername       = "admin"
	defaultAdminPassword       = "123456"
	participantSyntheticUserID = int64(900001)
)

func TestOrderVerifyRoutesAuthGuard(t *testing.T) {
	baseURL := strings.TrimRight(os.Getenv("DMH_INTEGRATION_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = orderAuthGuardBaseURL
	}

	client := &http.Client{Timeout: 10 * time.Second}
	adminUsername, adminPassword := getTestAdminCredentials()

	adminToken, err := loginAndGetToken(client, baseURL, adminUsername, adminPassword)
	if err != nil {
		t.Skipf("integration server unavailable or login failed: %v", err)
	}

	jwtSecret, err := resolveJWTSecret(adminToken)
	if err != nil {
		t.Skipf("cannot resolve jwt secret for participant token: %v", err)
	}

	participantToken, err := buildRoleToken(jwtSecret, participantSyntheticUserID, "participant-auth-guard", []string{"participant"})
	require.NoError(t, err)

	campaignID, err := createActiveCampaign(client, baseURL, adminToken)
	if err != nil {
		t.Skipf("cannot prepare active campaign: %v", err)
	}

	t.Run("verify auth matrix", func(t *testing.T) {
		orderID, phone, err := createOrder(client, baseURL, campaignID)
		require.NoError(t, err)

		code := generateVerificationCode(orderID, phone, time.Now().Unix())
		payload := map[string]string{"code": code}

		status, body, err := doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders/verify", "", payload)
		require.NoError(t, err)
		require.Equalf(t, http.StatusUnauthorized, status, "no-token verify response: %s", string(body))

		status, body, err = doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders/verify", participantToken, payload)
		require.NoError(t, err)
		require.Equalf(t, http.StatusForbidden, status, "participant verify response: %s", string(body))

		status, body, err = doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders/verify", adminToken, payload)
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, status, "admin verify response: %s", string(body))
	})

	t.Run("unverify auth matrix", func(t *testing.T) {
		orderID, phone, err := createOrder(client, baseURL, campaignID)
		require.NoError(t, err)

		code := generateVerificationCode(orderID, phone, time.Now().Unix())

		verifyPayload := map[string]string{"code": code, "remark": "setup verified for unverify auth test"}
		status, body, err := doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders/verify", adminToken, verifyPayload)
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, status, "verify setup response: %s", string(body))

		unverifyPayload := map[string]string{"code": code, "reason": "auth regression check"}

		status, body, err = doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders/unverify", "", unverifyPayload)
		require.NoError(t, err)
		require.Equalf(t, http.StatusUnauthorized, status, "no-token unverify response: %s", string(body))

		status, body, err = doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders/unverify", participantToken, unverifyPayload)
		require.NoError(t, err)
		require.Equalf(t, http.StatusForbidden, status, "participant unverify response: %s", string(body))

		status, body, err = doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders/unverify", adminToken, unverifyPayload)
		require.NoError(t, err)
		require.Equalf(t, http.StatusOK, status, "admin unverify response: %s", string(body))
	})
}

func TestOrderCreateDuplicateMessage(t *testing.T) {
	baseURL := strings.TrimRight(os.Getenv("DMH_INTEGRATION_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = orderAuthGuardBaseURL
	}

	client := &http.Client{Timeout: 10 * time.Second}
	adminUsername, adminPassword := getTestAdminCredentials()

	adminToken, err := loginAndGetToken(client, baseURL, adminUsername, adminPassword)
	if err != nil {
		t.Skipf("integration server unavailable or login failed: %v", err)
	}

	campaignID, err := createActiveCampaign(client, baseURL, adminToken)
	if err != nil {
		t.Skipf("cannot prepare active campaign: %v", err)
	}

	phone := fmt.Sprintf("139%08d", time.Now().UnixNano()%100000000)
	payload := map[string]interface{}{
		"campaignId": campaignID,
		"phone":      phone,
		"formData": map[string]string{
			"name": "重复报名回归用户",
		},
	}

	status, body, err := doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders", "", payload)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, status, "first create response: %s", string(body))

	status, body, err = doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders", "", payload)
	require.NoError(t, err)
	require.NotEqualf(t, http.StatusOK, status, "second create should fail: %s", string(body))
	require.Contains(t, string(body), "请勿重复报名")
	require.NotContains(t, strings.ToLower(string(body)), "duplicate entry")
}

func loginAndGetToken(client *http.Client, baseURL, username, password string) (string, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	status, body, err := doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/auth/login", "", payload)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("login returned %d: %s", status, string(body))
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parse login response: %w", err)
	}
	if resp.Token == "" {
		return "", fmt.Errorf("empty token in login response: %s", string(body))
	}

	return resp.Token, nil
}

func getTestAdminCredentials() (string, string) {
	username := strings.TrimSpace(os.Getenv("DMH_TEST_ADMIN_USERNAME"))
	if username == "" {
		username = defaultAdminUsername
	}

	password := strings.TrimSpace(os.Getenv("DMH_TEST_ADMIN_PASSWORD"))
	if password == "" {
		password = defaultAdminPassword
	}

	return username, password
}

func createActiveCampaign(client *http.Client, baseURL, adminToken string) (int64, error) {
	now := time.Now()
	payload := map[string]interface{}{
		"brandId":     1,
		"name":        fmt.Sprintf("订单鉴权回归活动-%d", now.UnixNano()),
		"description": "order verify auth guard regression",
		"rewardRule":  10.0,
		"startTime":   now.Add(-1 * time.Hour).Format("2006-01-02T15:04:05"),
		"endTime":     now.Add(24 * time.Hour).Format("2006-01-02T15:04:05"),
		"formFields": []map[string]interface{}{
			{
				"type":     "text",
				"name":     "name",
				"label":    "姓名",
				"required": true,
			},
		},
	}

	status, body, err := doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/campaigns", adminToken, payload)
	if err != nil {
		return 0, err
	}
	if status != http.StatusOK {
		return 0, fmt.Errorf("create campaign returned %d: %s", status, string(body))
	}

	var resp struct {
		Id int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("parse create campaign response: %w", err)
	}
	if resp.Id <= 0 {
		return 0, fmt.Errorf("invalid campaign id in response: %s", string(body))
	}

	return resp.Id, nil
}

func createOrder(client *http.Client, baseURL string, campaignID int64) (int64, string, error) {
	phone := fmt.Sprintf("139%08d", time.Now().UnixNano()%100000000)
	payload := map[string]interface{}{
		"campaignId": campaignID,
		"phone":      phone,
		"formData": map[string]string{
			"name": "鉴权回归用户",
		},
	}

	status, body, err := doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/orders", "", payload)
	if err != nil {
		return 0, "", err
	}
	if status != http.StatusOK {
		return 0, "", fmt.Errorf("create order returned %d: %s", status, string(body))
	}

	var resp struct {
		Id int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, "", fmt.Errorf("parse create order response: %w", err)
	}
	if resp.Id <= 0 {
		return 0, "", fmt.Errorf("invalid order id in response: %s", string(body))
	}

	return resp.Id, phone, nil
}

func doJSONRequest(client *http.Client, method, url, token string, payload interface{}) (int, []byte, error) {
	var reqBody io.Reader
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, err
		}
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

func generateVerificationCode(orderID int64, phone string, timestamp int64) string {
	signatureData := fmt.Sprintf("%d_%s_%d_%s", orderID, phone, timestamp, verificationSignSecret)
	hash := md5.Sum([]byte(signatureData))
	return fmt.Sprintf("%d_%s_%d_%s", orderID, phone, timestamp, hex.EncodeToString(hash[:]))
}

func buildRoleToken(secret string, userID int64, username string, roles []string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"userId":   userID,
		"username": username,
		"roles":    roles,
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
		"exp":      now.Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func resolveJWTSecret(referenceToken string) (string, error) {
	candidates := make([]string, 0)
	if envSecret := strings.TrimSpace(os.Getenv("DMH_ACCESS_SECRET")); envSecret != "" {
		candidates = append(candidates, envSecret)
	}

	configPaths := []string{
		"../../api/etc/dmh-api.yaml",
		"../../api/etc/dmh-api.dev.yaml",
		"../../api/etc/dmh-api.docker.yaml",
		"../../api/etc/dmh-api.prod.yaml",
	}
	for _, path := range configPaths {
		if secret, err := readAccessSecret(path); err == nil && secret != "" {
			candidates = append(candidates, secret)
		}
	}

	seen := make(map[string]struct{})
	for _, candidate := range candidates {
		if _, exists := seen[candidate]; exists {
			continue
		}
		seen[candidate] = struct{}{}

		if tokenLooksValidForSecret(referenceToken, candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no config secret can validate the admin token")
}

func readAccessSecret(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?m)^\s*AccessSecret:\s*['\"]?([^'\"#\s]+)['\"]?\s*$`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return "", fmt.Errorf("AccessSecret not found in %s", path)
	}

	return strings.TrimSpace(matches[1]), nil
}

func tokenLooksValidForSecret(tokenString, secret string) bool {
	parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false
	}

	return parsed.Valid
}
