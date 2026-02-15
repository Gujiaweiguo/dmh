package poster

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

func setupPosterHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.PosterTemplate{}, &model.PosterRecord{}, &model.PosterTemplateConfig{}, &model.Campaign{}, &model.Brand{}, &model.Distributor{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestPosterHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetPosterTemplatesHandler(nil))
	assert.NotNil(t, GetPosterRecordsHandler(nil))
	assert.NotNil(t, GenerateCampaignPosterHandler(nil))
	assert.NotNil(t, GenerateDistributorPosterHandler(nil))
}

func TestGetPosterTemplatesHandler_Success(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	template := &model.PosterTemplate{Type: "campaign", TemplateUrl: "https://example.com/template.png"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterTemplatesHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/templates", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPosterRecordsHandler_Success(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active"}
	db.Create(campaign)

	record := &model.PosterRecord{RecordType: "personal", CampaignID: campaign.Id, TemplateName: "Template", PosterUrl: "https://example.com/poster.png"}
	db.Create(record)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterRecordsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/records?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGenerateCampaignPosterHandler_ParseError(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateCampaignPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posters/campaign", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGenerateDistributorPosterHandler_ParseError(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posters/distributor", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetPosterTemplatesHandler_EmptyList(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterTemplatesHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/templates", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPosterRecordsHandler_EmptyList(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterRecordsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/records?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPosterRecordsHandler_WithFilters(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active"}
	db.Create(campaign)

	record := &model.PosterRecord{RecordType: "personal", CampaignID: campaign.Id, TemplateName: "Template", PosterUrl: "https://example.com/poster.png"}
	db.Create(record)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterRecordsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/records?page=1&pageSize=10&recordType=personal", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
