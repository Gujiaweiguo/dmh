package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestBindEmailHandler_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := createTestUserForHandler(t, db, "bindemail_user", "pwd1234")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := BindEmailHandler(svcCtx)

	reqBody := types.BindEmailReq{Email: "newemail@example.com", Code: "1234"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/bind-email", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.UserInfoResp
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if result.Email != reqBody.Email {
		t.Fatalf("expected email %s, got %s", reqBody.Email, result.Email)
	}
}

func TestBindPhoneHandler_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := createTestUserForHandler(t, db, "bindphone_user", "pwd1234")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := BindPhoneHandler(svcCtx)

	reqBody := types.BindPhoneReq{Phone: "13800139000", Code: "1234"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/bind-phone", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.UserInfoResp
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if result.Phone != reqBody.Phone {
		t.Fatalf("expected phone %s, got %s", reqBody.Phone, result.Phone)
	}
}

func TestSendEmailCodeHandler_Success(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	handler := SendEmailCodeHandler(svcCtx)

	reqBody := types.SendCodeReq{Target: "test@example.com", Type: "email"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/send-code", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.CommonResp
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if result.Message == "" {
		t.Fatalf("expected non-empty message in response, got empty")
	}
}

func TestSendPhoneCodeHandler_Success(t *testing.T) {
	svcCtx := &svc.ServiceContext{}
	handler := SendPhoneCodeHandler(svcCtx)

	reqBody := types.SendCodeReq{Target: "13800138000", Type: "phone"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/send-code", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.CommonResp
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if result.Message == "" {
		t.Fatalf("expected non-empty message in response, got empty")
	}
}

func TestUpdateProfileHandler_Success(t *testing.T) {
	db := setupHandlerTestDB(t)
	user := createTestUserForHandler(t, db, "update_user", "pwd1234")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateProfileHandler(svcCtx)

	reqBody := types.UpdateProfileReq{RealName: "NewRealName"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/auth/profile", bytes.NewReader(body))
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.UserInfoResp
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if result.RealName != reqBody.RealName {
		t.Fatalf("expected real name %s, got %s", reqBody.RealName, result.RealName)
	}
}
