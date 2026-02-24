package order

import (
	"bytes"
	"dmh/api/internal/handler/testutil"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupOrderHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(&model.Order{}, &model.Campaign{}, &model.Brand{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	testutil.ClearTables(db, "orders", "campaigns", "brands")
	return db
}

func TestOrderHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetOrdersHandler(nil))
	assert.NotNil(t, GetOrderHandler(nil))
	assert.NotNil(t, CreateOrderHandler(nil))
	assert.NotNil(t, ScanOrderHandler(nil))
	assert.NotNil(t, VerifyOrderHandler(nil))
	assert.NotNil(t, UnverifyOrderHandler(nil))
	assert.NotNil(t, GetVerificationRecordsHandler(nil))
	assert.NotNil(t, PaymentCallbackHandler(nil))
}

func TestGetOrdersHandler_Success(t *testing.T) {
	db := setupOrderHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		FormFields: "[]",
	}
	db.Create(campaign)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      testutil.GenUniquePhone(),
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
		FormData:   "{}",
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetOrdersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Handler construct tests
func TestGetOrderHandler_Construct(t *testing.T) {
	handler := GetOrderHandler(nil)
	assert.NotNil(t, handler)
}

func TestScanOrderHandler_Construct(t *testing.T) {
	handler := ScanOrderHandler(nil)
	assert.NotNil(t, handler)
}

func TestVerifyOrderHandler_Construct(t *testing.T) {
	handler := VerifyOrderHandler(nil)
	assert.NotNil(t, handler)
}

func TestUnverifyOrderHandler_Construct(t *testing.T) {
	handler := UnverifyOrderHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetVerificationRecordsHandler_Construct(t *testing.T) {
	handler := GetVerificationRecordsHandler(nil)
	assert.NotNil(t, handler)
}

func TestPaymentCallbackHandler_Construct(t *testing.T) {
	handler := PaymentCallbackHandler(nil)
	assert.NotNil(t, handler)
}

// Error path tests
func TestCreateOrderHandler_ParseError(t *testing.T) {
	handler := CreateOrderHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetOrderHandler_ParseError(t *testing.T) {
	db := setupOrderHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetOrderHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestScanOrderHandler_ParseError(t *testing.T) {
	db := setupOrderHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ScanOrderHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/scan", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestVerifyOrderHandler_ParseError(t *testing.T) {
	handler := VerifyOrderHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestPaymentCallbackHandler_ParseError(t *testing.T) {
	handler := PaymentCallbackHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/payment-callback", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetOrderHandler_Success(t *testing.T) {
	db := setupOrderHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		FormFields: "[]",
	}
	db.Create(campaign)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      testutil.GenUniquePhone(),
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
		FormData:   "{}",
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetOrderHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/orders/%d", order.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetOrderHandler_NotFound(t *testing.T) {
	db := setupOrderHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetOrderHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/99999", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetOrdersHandler_EmptyList(t *testing.T) {
	db := setupOrderHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetOrdersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetOrdersHandler_WithFilters(t *testing.T) {
	db := setupOrderHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		FormFields: "[]",
	}
	db.Create(campaign)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      testutil.GenUniquePhone(),
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
		FormData:   "{}",
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetOrdersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/orders?page=1&pageSize=10&phone=%s", order.Phone), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestScanOrderHandler_MissingCode(t *testing.T) {
	db := setupOrderHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := ScanOrderHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/scan", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestCreateOrderHandler_Success(t *testing.T) {
	db := setupOrderHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	now := time.Now()
	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		FormFields: "[]",
		StartTime:  now.Add(-24 * time.Hour),
		EndTime:    now.Add(24 * time.Hour),
	}
	db.Create(campaign)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := CreateOrderHandler(svcCtx)

	body := `{"campaignId":` + fmt.Sprintf("%d", campaign.Id) + `,"phone":"` + testutil.GenUniquePhone() + `","formData":{}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestVerifyOrderHandler_ParseError_Body(t *testing.T) {
	db := setupOrderHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := VerifyOrderHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", strings.NewReader(`{"orderId": "invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUnverifyOrderHandler_ParseError_Body(t *testing.T) {
	db := setupOrderHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UnverifyOrderHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/unverify", strings.NewReader(`{"orderId": "invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetVerificationRecordsHandler_Success(t *testing.T) {
	db := setupOrderHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetVerificationRecordsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/verification-records", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestVerifyOrderHandler_Success(t *testing.T) {
	db := setupOrderHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		FormFields: "[]",
	}
	db.Create(campaign)

	order := &model.Order{
		CampaignId:       campaign.Id,
		Phone:            testutil.GenUniquePhone(),
		Amount:           100.00,
		PayStatus:        "paid",
		Status:           "active",
		VerificationCode: "TESTCODE123",
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := VerifyOrderHandler(svcCtx)

	reqBody := types.VerifyOrderReq{
		Code:   "TESTCODE123",
		Remark: "Test verification",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUnverifyOrderHandler_Success(t *testing.T) {
	db := setupOrderHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		FormFields: "[]",
	}
	db.Create(campaign)

	verifiedAt := time.Now()
	order := &model.Order{
		CampaignId:       campaign.Id,
		Phone:            testutil.GenUniquePhone(),
		Amount:           100.00,
		PayStatus:        "paid",
		Status:           "verified",
		VerificationCode: "TESTCODE456",
		VerifiedAt:       &verifiedAt,
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UnverifyOrderHandler(svcCtx)

	reqBody := types.UnverifyOrderReq{
		Code:   "TESTCODE456",
		Reason: "Test unverification",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/unverify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestPaymentCallbackHandler_Success(t *testing.T) {
	db := setupOrderHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		FormFields: "[]",
	}
	db.Create(campaign)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      testutil.GenUniquePhone(),
		Amount:     100.00,
		PayStatus:  "pending",
		Status:     "active",
		TradeNo:    "TEST123456",
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := PaymentCallbackHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/payment-callback", strings.NewReader(`{"outTradeNo":"TEST123456","transactionId":"TXN123"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestScanOrderHandler_Success(t *testing.T) {
	db := setupOrderHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		Name:       "Test Campaign",
		BrandId:    brand.Id,
		Status:     "active",
		FormFields: "[]",
	}
	db.Create(campaign)

	order := &model.Order{
		CampaignId:       campaign.Id,
		Phone:            testutil.GenUniquePhone(),
		Amount:           100.00,
		PayStatus:        "paid",
		Status:           "active",
		VerificationCode: "SCANTEST123",
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ScanOrderHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/scan?code=SCANTEST123", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}
