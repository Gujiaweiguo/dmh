package statistics

import (
	"dmh/api/internal/handler/testutil"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupStatisticsHandlerTestDB(t *testing.T) *gorm.DB {
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(&model.Order{}, &model.Campaign{}, &model.Brand{}, &model.User{}, &model.Distributor{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestStatisticsHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetDashboardStatsHandler(nil))
}

func TestGetDashboardStatsHandler_InvalidBrandId(t *testing.T) {
	db := setupStatisticsHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDashboardStatsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics/dashboard?brandId=0", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetDashboardStatsHandler_MissingBrandId(t *testing.T) {
	db := setupStatisticsHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDashboardStatsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics/dashboard", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetDashboardStatsHandler_Success(t *testing.T) {
	db := setupStatisticsHandlerTestDB(t)

	// Create test data
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active"}
	db.Create(campaign)

	user := &model.User{Username: "testuser", Password: "hashed", Role: "participant"}
	db.Create(user)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		Amount:     100.0,
		PayStatus:  "paid",
		Status:     "confirmed",
		CreatedAt:  time.Now(),
	}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDashboardStatsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/dashboard?brandId=%d", brand.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetDashboardStatsHandler_InvalidJSON(t *testing.T) {
	db := setupStatisticsHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDashboardStatsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/statistics/dashboard", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	// Should return error for invalid request method/body
	assert.NotEqual(t, http.StatusOK, resp.Code)
}
