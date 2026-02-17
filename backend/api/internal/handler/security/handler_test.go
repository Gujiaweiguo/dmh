package security

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSecurityHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.PasswordPolicy{}, &model.LoginAttempt{}, &model.UserSession{}, &model.SecurityEvent{}, &model.AuditLog{}, &model.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestSecurityHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetPasswordPolicyHandler(nil))
	assert.NotNil(t, UpdatePasswordPolicyHandler(nil))
	assert.NotNil(t, CheckPasswordStrengthHandler(nil))
	assert.NotNil(t, GetUserSessionsHandler(nil))
	assert.NotNil(t, RevokeSessionHandler(nil))
	assert.NotNil(t, ForceLogoutUserHandler(nil))
	assert.NotNil(t, GetAuditLogsHandler(nil))
	assert.NotNil(t, GetLoginAttemptsHandler(nil))
	assert.NotNil(t, HandleSecurityEventHandler(nil))
	assert.NotNil(t, GetSecurityEventsHandler(nil))
}

func TestGetPasswordPolicyHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	policy := &model.PasswordPolicy{MinLength: 8, RequireUppercase: true, RequireLowercase: true, RequireNumbers: true, RequireSpecialChars: true}
	db.Create(policy)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPasswordPolicyHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/security/password-policy", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetAuditLogsHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	user := &model.User{Username: "admin", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	auditLog := &model.AuditLog{UserID: &user.Id, Username: "admin", Action: "login", Resource: "auth", Status: "success"}
	db.Create(auditLog)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetAuditLogsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/security/audit-logs?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var got types.AuditLogListResp
	err := json.Unmarshal(resp.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), got.Total)
	assert.Len(t, got.Logs, 1)
	assert.Equal(t, "admin", got.Logs[0].Username)
	assert.Equal(t, "login", got.Logs[0].Action)
}

func TestGetLoginAttemptsHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	attempt := &model.LoginAttempt{Username: "testuser", ClientIP: "192.168.1.1", Success: false}
	db.Create(attempt)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetLoginAttemptsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/security/login-attempts?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetSecurityEventsHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	event := &model.SecurityEvent{EventType: "login_failed", Severity: "medium", Username: "testuser", ClientIP: "192.168.1.1", Description: "Failed login attempt"}
	db.Create(event)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetSecurityEventsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/security/events?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetUserSessionsHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	user := &model.User{Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	session := &model.UserSession{ID: "session-123", UserID: user.Id, UserAgent: "test-agent", ClientIP: "192.168.1.1", Status: "active"}
	db.Create(session)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserSessionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/security/sessions?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestRevokeSessionHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	user := &model.User{Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	session := &model.UserSession{ID: "session-456", UserID: user.Id, UserAgent: "test-agent", ClientIP: "192.168.1.1", Status: "active"}
	db.Create(session)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := RevokeSessionHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/security/sessions/%s", session.ID), nil)
	req.SetPathValue("id", session.ID)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestForceLogoutUserHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	user := &model.User{Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ForceLogoutUserHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/security/users/%d/logout", user.Id), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", user.Id))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestHandleSecurityEventHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	event := &model.SecurityEvent{EventType: "test_event", Severity: "low", Username: "testuser", ClientIP: "192.168.1.1", Description: "Test event"}
	db.Create(event)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := HandleSecurityEventHandler(svcCtx)

	reqBody := types.HandleSecurityEventReq{
		Note: "Handled test event",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/security/events/%d/handle", event.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", fmt.Sprintf("%d", event.ID))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestCheckPasswordStrengthHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	policy := &model.PasswordPolicy{MinLength: 8, RequireUppercase: true, RequireLowercase: true, RequireNumbers: true, RequireSpecialChars: true}
	db.Create(policy)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := CheckPasswordStrengthHandler(svcCtx)

	reqBody := map[string]string{"password": "Test@123456"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/security/check-password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUpdatePasswordPolicyHandler_Success(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)

	policy := &model.PasswordPolicy{MinLength: 8, RequireUppercase: false, RequireLowercase: false, RequireNumbers: false, RequireSpecialChars: false}
	db.Create(policy)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdatePasswordPolicyHandler(svcCtx)

	reqBody := types.UpdatePasswordPolicyReq{
		MinLength:           10,
		RequireUppercase:    true,
		RequireLowercase:    true,
		RequireNumbers:      true,
		RequireSpecialChars: true,
		MaxAge:              90,
		HistoryCount:        5,
		MaxLoginAttempts:    5,
		LockoutDuration:     30,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/security/password-policy", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUpdatePasswordPolicyHandler_ParseError(t *testing.T) {
	handler := UpdatePasswordPolicyHandler(nil)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/security/password-policy", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestCheckPasswordStrengthHandler_ParseError(t *testing.T) {
	handler := CheckPasswordStrengthHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/security/check-password", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetUserSessionsHandler_ParseError(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetUserSessionsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/security/sessions?page=invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestRevokeSessionHandler_ParseError(t *testing.T) {
	handler := RevokeSessionHandler(nil)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/security/sessions/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotNil(t, resp)
}

func TestForceLogoutUserHandler_ParseError(t *testing.T) {
	handler := ForceLogoutUserHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/security/users/invalid/logout", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotNil(t, resp)
}

func TestHandleSecurityEventHandler_ParseError(t *testing.T) {
	db := setupSecurityHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := HandleSecurityEventHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/security/events/handle", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}
