package distributor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createTestCampaign(t *testing.T, db *gorm.DB, brandId int64, name string) *model.Campaign {
	t.Helper()
	campaign := &model.Campaign{
		BrandId:    brandId,
		Name:       name,
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		Status:     "active",
		RewardRule: 10,
	}
	if err := db.Create(campaign).Error; err != nil {
		t.Fatalf("failed to create test campaign: %v", err)
	}
	return campaign
}

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
	campaign := createTestCampaign(t, db, brand.Id, "Test Campaign")
	dist := &model.Distributor{UserId: user.Id, BrandId: brand.Id, Status: "active"}
	if err := db.Create(dist).Error; err != nil {
		t.Fatalf("failed to create distributor: %v", err)
	}
	link := &model.DistributorLink{DistributorId: dist.Id, CampaignId: campaign.Id, LinkCode: "abcd1234"}
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

// Test for GetDistributorSubordinatesHandler with DB is pending full integration wiring

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
	campaign := createTestCampaign(t, db, brand.Id, "Link Campaign")
	dist := &model.Distributor{UserId: user.Id, BrandId: campaign.Id, Status: "active"}
	if err := db.Create(dist).Error; err != nil {
		t.Fatalf("failed to create distributor: %v", err)
	}
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GenerateDistributorLinkHandler(svcCtx)

	reqBody := types.GenerateLinkReq{CampaignId: campaign.Id}
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

func TestApproveDistributorApplicationHandler_WithDB_Success(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := seedBrandAndDistributor(t, db)
	user := createTestUser(t, db, "applicant_user")

	application := &model.DistributorApplication{
		UserId:  user.Id,
		BrandId: brand.Id,
		Reason:  "Need to become distributor",
		Status:  "pending",
	}
	if err := db.Create(application).Error; err != nil {
		t.Fatalf("failed to create distributor application: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApproveDistributorApplicationHandler(svcCtx)

	reqBody := types.ApproveDistributorReq{Action: "approve"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/distributors/applications/%d/approve", application.Id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "applicationId", application.Id)
	ctx = context.WithValue(ctx, "userId", user.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}
