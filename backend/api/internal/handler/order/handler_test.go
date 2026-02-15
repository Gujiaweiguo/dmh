package order

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupOrderHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Order{}, &model.Campaign{}, &model.Brand{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

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
		Phone:      "13800138000",
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
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
		Phone:      "13800138000",
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
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
		Phone:      "13800138000",
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetOrdersHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders?page=1&pageSize=10&phone=13800138000", nil)
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

	body := `{"campaignId":` + fmt.Sprintf("%d", campaign.Id) + `,"phone":"13800138000","formData":{}}`
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
