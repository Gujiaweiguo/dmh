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

func TestGenerateCampaignPosterHandler_InvalidBody(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateCampaignPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posters/campaign", strings.NewReader(`{"id": "invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGenerateDistributorPosterHandler_InvalidPath(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/invalid/poster", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGenerateDistributorPosterHandler_ShortPath(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/poster", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGenerateDistributorPosterHandler_TooShortPath(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/1/poster", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetPosterTemplatesHandler_WithType(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	template1 := &model.PosterTemplate{Type: "campaign", TemplateUrl: "https://example.com/template1.png"}
	template2 := &model.PosterTemplate{Type: "distributor", TemplateUrl: "https://example.com/template2.png"}
	db.Create(template1)
	db.Create(template2)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterTemplatesHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/templates?type=campaign", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGenerateCampaignPosterHandler_Success(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active", PosterTemplateId: 1}
	db.Create(campaign)

	template := &model.PosterTemplateConfig{Name: "Test Template", Status: "active"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateCampaignPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posters/campaign", strings.NewReader(`{"id": 1, "templateId": 1}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGenerateDistributorPosterHandler_Success(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	user := &model.User{Username: "testuser", Password: "hashed", Role: "participant"}
	db.Create(user)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	distributor := &model.Distributor{UserId: user.Id, BrandId: brand.Id, Level: 1, Status: "active"}
	db.Create(distributor)

	template := &model.PosterTemplateConfig{Name: "Test Template", Status: "active"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/distributors/%d/poster", distributor.Id), strings.NewReader(`{"templateId": 1}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGenerateDistributorPosterHandler_ValidIdPath(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	user := &model.User{Username: "testuser2", Password: "hashed", Role: "participant"}
	db.Create(user)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	distributor := &model.Distributor{UserId: user.Id, BrandId: brand.Id, Level: 1, Status: "active"}
	db.Create(distributor)

	template := &model.PosterTemplateConfig{Name: "Test Template", Status: "active"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/distributors/%d/poster", distributor.Id), strings.NewReader(`{"templateId": 1}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGenerateDistributorPosterHandler_ExactPathLength(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors//poster", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGenerateCampaignPosterHandler_WithValidTemplate(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active", PosterTemplateId: 1}
	db.Create(campaign)

	template := &model.PosterTemplateConfig{Name: "Test Template", Status: "active"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateCampaignPosterHandler(svcCtx)

	reqBody := fmt.Sprintf(`{"id": %d, "templateId": %d}`, campaign.Id, template.Id)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posters/campaign", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGenerateDistributorPosterHandler_WithValidDistributor(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	user := &model.User{Username: "distuser", Password: "hashed", Role: "participant"}
	db.Create(user)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	distributor := &model.Distributor{UserId: user.Id, BrandId: brand.Id, Level: 1, Status: "active"}
	db.Create(distributor)

	template := &model.PosterTemplateConfig{Name: "Test Template", Status: "active"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	reqBody := fmt.Sprintf(`{"templateId": %d}`, template.Id)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/distributors/%d/poster", distributor.Id), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetPosterRecordsHandler_WithRecordType(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active"}
	db.Create(campaign)

	record := &model.PosterRecord{RecordType: "campaign", CampaignID: campaign.Id, TemplateName: "Template", PosterUrl: "https://example.com/poster.png"}
	db.Create(record)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterRecordsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/records?page=1&pageSize=10&recordType=campaign", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPosterRecordsHandler_WithCampaignId(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active"}
	db.Create(campaign)

	record := &model.PosterRecord{RecordType: "campaign", CampaignID: campaign.Id, TemplateName: "Template", PosterUrl: "https://example.com/poster.png"}
	db.Create(record)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterRecordsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posters/records?page=1&pageSize=10&campaignId=%d", campaign.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPosterTemplatesHandler_WithStatus(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	template1 := &model.PosterTemplateConfig{Name: "Active Template", Status: "active"}
	template2 := &model.PosterTemplateConfig{Name: "Inactive Template", Status: "inactive"}
	db.Create(template1)
	db.Create(template2)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterTemplatesHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/templates?status=active", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPosterTemplatesHandler_WithKeyword(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	template := &model.PosterTemplateConfig{Name: "Special Campaign Template", Status: "active"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterTemplatesHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/templates?keyword=Special", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGenerateDistributorPosterHandler_LongPath(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	user := &model.User{Username: "testuser3", Password: "hashed", Role: "participant"}
	db.Create(user)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	distributor := &model.Distributor{UserId: user.Id, BrandId: brand.Id, Level: 1, Status: "active"}
	db.Create(distributor)

	template := &model.PosterTemplateConfig{Name: "Test Template", Status: "active"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/distributors/%d/poster", distributor.Id), strings.NewReader(`{"templateId": 1}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGenerateCampaignPosterHandler_WithTemplateId(t *testing.T) {
	db := setupPosterHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active", PosterTemplateId: 1}
	db.Create(campaign)

	template := &model.PosterTemplateConfig{Name: "Test Template", Status: "active"}
	db.Create(template)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateCampaignPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posters/campaign", strings.NewReader(fmt.Sprintf(`{"id": %d, "templateId": 1}`, campaign.Id)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetPosterTemplatesHandler_ParseError(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPosterTemplatesHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posters/templates?page=invalid", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetPosterRecordsHandler_WithMemberId(t *testing.T) {
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

func TestGenerateDistributorPosterHandler_InvalidIdFormat(t *testing.T) {
	db := setupPosterHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorPosterHandler(svcCtx)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/abc123/poster", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}
