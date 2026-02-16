package sync

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSyncHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.SyncLog{}, &model.Order{}, &model.Campaign{}, &model.Brand{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestSyncHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetSyncHealthHandler(nil))
	assert.NotNil(t, GetSyncStatusHandler(nil))
	assert.NotNil(t, GetSyncStatsHandler(nil))
	assert.NotNil(t, RetrySynHandler(nil))
}

func TestGetSyncHealthHandler_Success(t *testing.T) {
	db := setupSyncHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetSyncHealthHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sync/health", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetSyncStatusHandler_Success(t *testing.T) {
	db := setupSyncHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active"}
	db.Create(campaign)

	order := &model.Order{CampaignId: campaign.Id, Phone: "13800138000", Amount: 100.00, PayStatus: "paid", Status: "active", SyncStatus: "synced"}
	db.Create(order)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetSyncStatusHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sync/status?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetSyncStatsHandler_Success(t *testing.T) {
	db := setupSyncHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetSyncStatsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sync/stats", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Ensure RetrySynHandler is wired and returns 200 OK with default logic (nil response)
func TestRetrySynHandler_Success(t *testing.T) {
	db := setupSyncHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := RetrySynHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sync/retry", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
