package distributor

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Helper to initialise a fresh in-memory DB and seed common entities
func seedBrandAndDistributor(t *testing.T, db *gorm.DB) *model.Brand {
	brand := &model.Brand{Name: "CoverageBrand", Status: "active"}
	if err := db.Create(brand).Error; err != nil {
		t.Fatalf("failed to seed brand: %v", err)
	}
	return brand
}

func TestGetDistributorLinksHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	user := createTestUser(t, db, "link_user")
	brand := seedBrandAndDistributor(t, db)
	dist := &model.Distributor{UserId: user.Id, BrandId: brand.Id, Status: "active"}
	if err := db.Create(dist).Error; err != nil {
		t.Fatalf("failed to create distributor: %v", err)
	}
	link := &model.DistributorLink{DistributorId: dist.Id, CampaignId: brand.Id, LinkCode: "abcd1234"}
	if err := db.Create(link).Error; err != nil {
		t.Fatalf("failed to create distributor link: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDistributorLinksHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/distributors/links", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetDistributorApplicationsHandler_WithDB_New(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	user := createTestUser(t, db, "apps_user")
	brand := seedBrandAndDistributor(t, db)
	// seed a distributor application
	app := &model.DistributorApplication{UserId: user.Id, BrandId: brand.Id, Reason: "need to be distributor", Status: "pending"}
	if err := db.Create(app).Error; err != nil {
		t.Fatalf("failed to create distributor application: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDistributorApplicationsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/distributors/applications?page=1&pageSize=10&brandId=1&status=pending", nil)
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetDistributorByCodeHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	user := createTestUser(t, db, "code_user")
	brand := seedBrandAndDistributor(t, db)
	dist := &model.Distributor{UserId: user.Id, BrandId: brand.Id, Status: "active"}
	if err := db.Create(dist).Error; err != nil {
		t.Fatalf("failed to create distributor: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDistributorByCodeHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/distributors/by-code/"+strconv.FormatInt(dist.Id, 10), nil)
	req.SetPathValue("code", strconv.FormatInt(dist.Id, 10))
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGenerateDistributorLinkHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	user := createTestUser(t, db, "gen_user")
	brand := seedBrandAndDistributor(t, db)
	dist := &model.Distributor{UserId: user.Id, BrandId: brand.Id, Status: "active"}
	if err := db.Create(dist).Error; err != nil {
		t.Fatalf("failed to create distributor: %v", err)
	}
	// Prepare request to generate link with campaignId matching brand id (per logic)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorLinkHandler(svcCtx)

	reqBody := types.GenerateLinkReq{CampaignId: brand.Id}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/generate-link", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()
	handler(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
