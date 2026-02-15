package campaign

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

func setupCampaignHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Campaign{}, &model.Brand{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestCampaignHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetCampaignsHandler(nil))
	assert.NotNil(t, GetCampaignHandler(nil))
	assert.NotNil(t, CreateCampaignHandler(nil))
	assert.NotNil(t, UpdateCampaignHandler(nil))
	assert.NotNil(t, DeleteCampaignHandler(nil))
	assert.NotNil(t, SavePageConfigHandler(nil))
	assert.NotNil(t, GetPageConfigHandler(nil))
	assert.NotNil(t, GetPaymentQrcodeHandler(nil))
}

func TestGetCampaignsHandler_Success(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	db.Create(campaign)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetCampaignsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Handler construct tests
func TestGetCampaignHandler_Construct(t *testing.T) {
	handler := GetCampaignHandler(nil)
	assert.NotNil(t, handler)
}

func TestCreateCampaignHandler_Construct(t *testing.T) {
	handler := CreateCampaignHandler(nil)
	assert.NotNil(t, handler)
}

func TestUpdateCampaignHandler_Construct(t *testing.T) {
	handler := UpdateCampaignHandler(nil)
	assert.NotNil(t, handler)
}

func TestDeleteCampaignHandler_Construct(t *testing.T) {
	handler := DeleteCampaignHandler(nil)
	assert.NotNil(t, handler)
}

func TestSavePageConfigHandler_Construct(t *testing.T) {
	handler := SavePageConfigHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetPageConfigHandler_Construct(t *testing.T) {
	handler := GetPageConfigHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetPaymentQrcodeHandler_Construct(t *testing.T) {
	handler := GetPaymentQrcodeHandler(nil)
	assert.NotNil(t, handler)
}

// Error path tests
func TestCreateCampaignHandler_ParseError(t *testing.T) {
	handler := CreateCampaignHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetCampaignHandler_ParseError(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetCampaignHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestUpdateCampaignHandler_ParseError(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateCampaignHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/campaigns/1", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetPageConfigHandler_ParseError(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPageConfigHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/invalid/page-config", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetPaymentQrcodeHandler_ParseError(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPaymentQrcodeHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/invalid/payment-qrcode", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestDeleteCampaignHandler_ParseError(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteCampaignHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/campaigns/invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetCampaignHandler_Success(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	db.Create(campaign)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetCampaignHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/campaigns/%d", campaign.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetCampaignHandler_NotFound(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetCampaignHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/99999", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestCreateCampaignHandler_Success(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := CreateCampaignHandler(svcCtx)

	now := time.Now()
	body := fmt.Sprintf(`{"brandId":%d,"name":"New Campaign","description":"Test","startTime":"%s","endTime":"%s","formFields":[],"rewardRule":0}`,
		brand.Id,
		now.Format(time.RFC3339),
		now.Add(24*time.Hour).Format(time.RFC3339),
	)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestUpdateCampaignHandler_ParseError_InvalidID(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateCampaignHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/campaigns/invalid", strings.NewReader(`{"name":"Updated"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestDeleteCampaignHandler_Success(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	db.Create(campaign)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DeleteCampaignHandler(svcCtx)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/campaigns/%d", campaign.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetCampaignsHandler_EmptyList(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetCampaignsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetCampaignsHandler_WithFilters(t *testing.T) {
	db := setupCampaignHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	db.Create(campaign)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetCampaignsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns?page=1&pageSize=10&status=active", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
