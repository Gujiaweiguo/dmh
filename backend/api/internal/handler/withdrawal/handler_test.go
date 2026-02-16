package withdrawal

import (
	"bytes"
	"context"
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

func setupWithdrawalHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Withdrawal{}, &model.User{}, &model.Brand{}, &model.Distributor{}, &model.UserBalance{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func createTestUserForWithdrawal(t *testing.T, db *gorm.DB, username string) *model.User {
	user := &model.User{
		Username: username,
		Password: "hashed_password",
		Phone:    "13800138000",
		Role:     "participant",
		Status:   "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestBrandForWithdrawal(t *testing.T, db *gorm.DB, name string) *model.Brand {
	brand := &model.Brand{
		Name:   name,
		Status: "active",
	}
	if err := db.Create(brand).Error; err != nil {
		t.Fatalf("Failed to create test brand: %v", err)
	}
	return brand
}

func TestWithdrawalHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetWithdrawalsHandler(nil))
	assert.NotNil(t, GetWithdrawalHandler(nil))
	assert.NotNil(t, ApplyWithdrawalHandler(nil))
	assert.NotNil(t, ApproveWithdrawalHandler(nil))
}

func TestGetWithdrawalsHandler_Success(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)

	user := &model.User{Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	withdrawal := &model.Withdrawal{UserID: user.Id, BrandId: brand.Id, DistributorId: 1, Amount: 100.00, Status: "pending", PayType: "wechat", PayAccount: "openid"}
	db.Create(withdrawal)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetWithdrawalsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdrawals?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetWithdrawalsHandler_EmptyList(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetWithdrawalsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdrawals?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetWithdrawalsHandler_WithStatus(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	user := createTestUserForWithdrawal(t, db, "testuser")
	brand := createTestBrandForWithdrawal(t, db, "Test Brand")

	withdrawal := &model.Withdrawal{UserID: user.Id, BrandId: brand.Id, Amount: 100.00, Status: "approved", PayType: "wechat", PayAccount: "openid"}
	db.Create(withdrawal)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetWithdrawalsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdrawals?page=1&pageSize=10&status=approved", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetWithdrawalHandler_ParseError(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetWithdrawalHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdrawals/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetWithdrawalHandler_NotFound(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetWithdrawalHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdrawals/99999", nil)
	req.SetPathValue("id", "99999")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestApplyWithdrawalHandler_ParseError(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApplyWithdrawalHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/withdrawals", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

// New tests to improve coverage for Apply and Approve withdrawal handlers
// These tests exercise the successful execution paths and DB side effects
// using the testify assertion framework.
func TestApplyWithdrawalHandler_Success(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	// create test user
	user := &model.User{Id: 1, Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)
	// initial balance sufficient for withdrawal
	balance := &model.UserBalance{UserId: 1, Balance: 500}
	db.Create(balance)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApplyWithdrawalHandler(svcCtx)

	reqBody := types.WithdrawalApplyReq{Amount: 100, BankName: "ICBC", BankAccount: "6222021234567890", AccountName: "张三"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/withdrawals", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// simulate authenticated user in context
	req = req.WithContext(context.WithValue(req.Context(), "userId", int64(1)))
	resp := httptest.NewRecorder()

	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// verify response payload
	var withdrawalResp types.WithdrawalResp
	err := json.NewDecoder(resp.Body).Decode(&withdrawalResp)
	assert.NoError(t, err)
	assert.Equal(t, float64(reqBody.Amount), withdrawalResp.Amount)
	assert.Equal(t, "pending", withdrawalResp.Status)

	// verify DB changes: withdrawal created and balance deducted
	var w model.Withdrawal
	db.Last(&w)
	assert.Equal(t, int64(1), w.UserID)
	assert.Equal(t, reqBody.Amount, w.Amount)

	var b model.UserBalance
	db.Where("user_id = ?", 1).First(&b)
	assert.Equal(t, float64(400), b.Balance)
}

func TestApproveWithdrawalHandler_Success(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	// prepare test data: user, balance, and a pending withdrawal
	user := &model.User{Id: 1, Username: "adminUser", Phone: "13800138000"}
	balance := &model.UserBalance{UserId: 1, Balance: 500}
	db.Create(user)
	db.Create(balance)
	withdrawal := &model.Withdrawal{ID: 1, UserID: 1, BrandId: 0, DistributorId: 0, Amount: 100, Status: "pending", BankName: "ICBC", BankAccount: "6222021234567890", AccountName: "张三"}
	db.Create(withdrawal)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApproveWithdrawalHandler(svcCtx)

	reqBody := types.WithdrawalApproveReq{Status: "approved"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/withdrawals/%d/approve", withdrawal.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userId", int64(1))) // admin id
	req.SetPathValue("id", fmt.Sprintf("%d", withdrawal.ID))
	resp := httptest.NewRecorder()

	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// verify withdrawal status updated and approved metadata set
	var w model.Withdrawal
	db.First(&w, withdrawal.ID)
	assert.Equal(t, "approved", w.Status)
	assert.NotNil(t, w.ApprovedBy)
}

func TestApproveWithdrawalHandler_RejectedWithRefund(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	// prepare test data: user, balance, and a pending withdrawal
	user := &model.User{Id: 1, Username: "adminUser", Phone: "13800138000"}
	balance := &model.UserBalance{UserId: 1, Balance: 500}
	db.Create(user)
	db.Create(balance)
	withdrawal := &model.Withdrawal{ID: 1, UserID: 1, Amount: 100, Status: "pending", BankName: "ICBC", BankAccount: "6222021234567890", AccountName: "张三"}
	db.Create(withdrawal)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApproveWithdrawalHandler(svcCtx)

	reqBody := types.WithdrawalApproveReq{Status: "rejected", Remark: "Invalid"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/withdrawals/%d/approve", withdrawal.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "userId", int64(1))) // admin id
	req.SetPathValue("id", fmt.Sprintf("%d", withdrawal.ID))
	resp := httptest.NewRecorder()

	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// verify withdrawal status updated and balance refunded
	var w model.Withdrawal
	db.First(&w, withdrawal.ID)
	assert.Equal(t, "rejected", w.Status)
	assert.NotNil(t, w.ApprovedBy)

	var b model.UserBalance
	db.Where("user_id = ?", 1).First(&b)
	assert.Equal(t, float64(600), b.Balance) // 500 + 100 refund
}
func TestApproveWithdrawalHandler_ParseError(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApproveWithdrawalHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/withdrawals/1/approve", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestApproveWithdrawalHandler_InvalidID(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApproveWithdrawalHandler(svcCtx)

	reqBody := types.WithdrawalApproveReq{Status: "approved"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/withdrawals/invalid/approve", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetWithdrawalHandler_Success(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	user := createTestUserForWithdrawal(t, db, "testuser")
	brand := createTestBrandForWithdrawal(t, db, "Test Brand")

	withdrawal := &model.Withdrawal{UserID: user.Id, BrandId: brand.Id, Amount: 100.00, Status: "pending", PayType: "wechat", PayAccount: "openid"}
	db.Create(withdrawal)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetWithdrawalHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/withdrawals/%d", withdrawal.ID), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", withdrawal.ID))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestApplyWithdrawalHandler_InvalidAmount(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApplyWithdrawalHandler(svcCtx)

	reqBody := types.WithdrawalApplyReq{Amount: -100}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/withdrawals", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetWithdrawalsHandler_ParseError(t *testing.T) {
	db := setupWithdrawalHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetWithdrawalsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdrawals?page=abc&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}
