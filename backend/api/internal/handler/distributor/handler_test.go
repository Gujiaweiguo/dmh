package distributor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dmh/api/internal/handler/testutil"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDistributorHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(
		&model.Brand{},
		&model.Campaign{},
		&model.DistributorApplication{},
		&model.DistributorLink{},
		&model.Distributor{},
		&model.DistributorLevelReward{},
		&model.DistributorReward{},
		&model.User{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	testutil.ClearTables(db, "distributors", "distributor_applications", "distributor_links", "distributor_rewards", "distributor_level_rewards", "brands", "campaigns", "users")

	return db
}

func createTestBrand(t *testing.T, db *gorm.DB, name string) *model.Brand {
	t.Helper()
	brand := &model.Brand{
		Name:   name + fmt.Sprintf("_%d", time.Now().UnixNano()),
		Status: "active",
	}
	if err := db.Create(brand).Error; err != nil {
		t.Fatalf("Failed to create test brand: %v", err)
	}
	return brand
}

func createTestUser(t *testing.T, db *gorm.DB, usernamePrefix string) *model.User {
	t.Helper()
	user := &model.User{
		Username: testutil.GenUniqueUsername(usernamePrefix),
		Password: "hashed_password",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func TestDistributorHandlersConstruct(t *testing.T) {
	assert.NotNil(t, DistributorApplyHandler(nil))
	assert.NotNil(t, GetDistributorApplicationsHandler(nil))
	assert.NotNil(t, GetDistributorApplicationHandler(nil))
	assert.NotNil(t, ApproveDistributorApplicationHandler(nil))
	assert.NotNil(t, GetMyDistributorStatusHandler(nil))
	assert.NotNil(t, GetMyDistributorDashboardHandler(nil))
	assert.NotNil(t, GetDistributorSubordinatesHandler(nil))
	assert.NotNil(t, GetDistributorStatisticsHandler(nil))
	assert.NotNil(t, GetDistributorByCodeHandler(nil))
	assert.NotNil(t, GenerateDistributorLinkHandler(nil))
	assert.NotNil(t, TrackDistributorLinkHandler(nil))
	assert.NotNil(t, GetDistributorLinksHandler(nil))
	assert.NotNil(t, GetDistributorQrcodeHandler(nil))
	assert.NotNil(t, GetDistributorRewardsHandler(nil))
	assert.NotNil(t, GetBrandDistributorsHandler(nil))
	assert.NotNil(t, GetBrandDistributorHandler(nil))
	assert.NotNil(t, UpdateDistributorStatusHandler(nil))
	assert.NotNil(t, GetBrandDistributorApplicationsHandler(nil))
	assert.NotNil(t, GetBrandDistributorApplicationHandler(nil))
	assert.NotNil(t, UpdateDistributorLevelHandler(nil))
	assert.NotNil(t, GetDistributorLevelRewardsHandler(nil))
	assert.NotNil(t, SetDistributorLevelRewardsHandler(nil))
}

func TestDistributorHandler_ParseErrors(t *testing.T) {
	testCases := []struct {
		name        string
		handler     http.HandlerFunc
		method      string
		url         string
		body        string
		contentType string
	}{
		{
			name:        "DistributorApplyHandler",
			handler:     DistributorApplyHandler(nil),
			method:      http.MethodPost,
			url:         "/api/v1/distributors/apply",
			body:        "{invalid}",
			contentType: "application/json",
		},
		{
			name:        "GetDistributorByCodeHandler",
			handler:     GetDistributorByCodeHandler(nil),
			method:      http.MethodGet,
			url:         "/api/v1/distributors/by-code/invalid",
			body:        "",
			contentType: "",
		},
		{
			name:        "GetDistributorApplicationHandler",
			handler:     GetDistributorApplicationHandler(nil),
			method:      http.MethodGet,
			url:         "/api/v1/distributors/applications/invalid",
			body:        "",
			contentType: "",
		},
		{
			name:        "ApproveDistributorApplicationHandler",
			handler:     ApproveDistributorApplicationHandler(nil),
			method:      http.MethodPost,
			url:         "/api/v1/distributors/applications/invalid/approve",
			body:        "{invalid}",
			contentType: "application/json",
		},
		{
			name:        "GetDistributorSubordinatesHandler",
			handler:     GetDistributorSubordinatesHandler(nil),
			method:      http.MethodGet,
			url:         "/api/v1/distributors/invalid/subordinates",
			body:        "",
			contentType: "",
		},
		{
			name:        "GetDistributorRewardsHandler",
			handler:     GetDistributorRewardsHandler(nil),
			method:      http.MethodGet,
			url:         "/api/v1/distributors/invalid/rewards",
			body:        "",
			contentType: "",
		},
		{
			name:        "GetBrandDistributorApplicationHandler",
			handler:     GetBrandDistributorApplicationHandler(nil),
			method:      http.MethodGet,
			url:         "/api/v1/brands/invalid/distributors/applications/invalid",
			body:        "",
			contentType: "",
		},
		{
			name:        "UpdateDistributorStatusHandler",
			handler:     UpdateDistributorStatusHandler(nil),
			method:      http.MethodPost,
			url:         "/api/v1/brands/invalid/distributors/invalid/status",
			body:        "{invalid}",
			contentType: "application/json",
		},
		{
			name:        "UpdateDistributorLevelHandler",
			handler:     UpdateDistributorLevelHandler(nil),
			method:      http.MethodPost,
			url:         "/api/v1/brands/invalid/distributors/invalid/level",
			body:        "{invalid}",
			contentType: "application/json",
		},
		{
			name:        "SetDistributorLevelRewardsHandler",
			handler:     SetDistributorLevelRewardsHandler(nil),
			method:      http.MethodPost,
			url:         "/api/v1/distributors/level-rewards",
			body:        "{invalid}",
			contentType: "application/json",
		},
		{
			name:        "GenerateDistributorLinkHandler",
			handler:     GenerateDistributorLinkHandler(nil),
			method:      http.MethodPost,
			url:         "/api/v1/distributors/generate-link",
			body:        "{invalid}",
			contentType: "application/json",
		},
		{
			name:        "TrackDistributorLinkHandler",
			handler:     TrackDistributorLinkHandler(nil),
			method:      http.MethodPost,
			url:         "/api/v1/distributors/track-link",
			body:        "{invalid}",
			contentType: "application/json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"_ParseError", func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			resp := httptest.NewRecorder()

			tc.handler(resp, req)

			assert.NotEqual(t, http.StatusOK, resp.Code)
		})
	}
}

func TestDistributorApplyHandler_Success(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "testuser")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DistributorApplyHandler(svcCtx)

	reqBody := types.DistributorApplyReq{
		BrandId: brand.Id,
		Reason:  "I want to be a distributor",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.Id))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetMyDistributorStatusHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	user := createTestUser(t, db, "testuser")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMyDistributorStatusHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/distributors/status", nil)
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.Id))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetBrandDistributorsHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetBrandDistributorsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/brands/%d/distributors?page=1&pageSize=10", brand.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetBrandDistributorApplicationsHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetBrandDistributorApplicationsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/brands/%d/distributors/applications?page=1&pageSize=10", brand.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetDistributorStatisticsHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	user := createTestUser(t, db, "testuser")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDistributorStatisticsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/distributors/%d/statistics", user.Id), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestDistributorApplyHandler_WithEmptyReason(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := DistributorApplyHandler(svcCtx)

	reqBody := types.DistributorApplyReq{
		BrandId: brand.Id,
		Reason:  "",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/distributors/apply", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetBrandDistributorHandler_Construct(t *testing.T) {
	handler := GetBrandDistributorHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetDistributorApplicationsHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "testuser")

	application := &model.DistributorApplication{
		UserId:  user.Id,
		BrandId: brand.Id,
		Reason:  "Test reason",
		Status:  "pending",
	}
	db.Create(application)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDistributorApplicationsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/distributors/applications?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetDistributorLevelRewardsHandler_Construct(t *testing.T) {
	handler := GetDistributorLevelRewardsHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetDistributorLinksHandler_Construct(t *testing.T) {
	handler := GetDistributorLinksHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetDistributorQrcodeHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "Test Brand")

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	if err := db.Create(distributor).Error; err != nil {
		t.Fatalf("Failed to create distributor: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDistributorQrcodeHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/distributors/%d/qrcode", user.Id), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", user.Id))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetMyDistributorDashboardHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "Test Brand")

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	if err := db.Create(distributor).Error; err != nil {
		t.Fatalf("Failed to create distributor: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetMyDistributorDashboardHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/distributors/dashboard", nil)
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", user.Id))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestApproveDistributorApplicationHandler_Construct(t *testing.T) {
	handler := ApproveDistributorApplicationHandler(nil)
	assert.NotNil(t, handler)
}

func TestUpdateDistributorStatusHandler_Construct(t *testing.T) {
	handler := UpdateDistributorStatusHandler(nil)
	assert.NotNil(t, handler)
}

func TestUpdateDistributorLevelHandler_Construct(t *testing.T) {
	handler := UpdateDistributorLevelHandler(nil)
	assert.NotNil(t, handler)
}

func TestSetDistributorLevelRewardsHandler_Construct(t *testing.T) {
	handler := SetDistributorLevelRewardsHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetDistributorRewardsHandler_Construct(t *testing.T) {
	handler := GetDistributorRewardsHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetDistributorSubordinatesHandler_Construct(t *testing.T) {
	handler := GetDistributorSubordinatesHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetDistributorQrcodeHandler_Construct(t *testing.T) {
	handler := GetDistributorQrcodeHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetMyDistributorDashboardHandler_Construct(t *testing.T) {
	handler := GetMyDistributorDashboardHandler(nil)
	assert.NotNil(t, handler)
}

func TestTrackDistributorLinkHandler_Construct(t *testing.T) {
	handler := TrackDistributorLinkHandler(nil)
	assert.NotNil(t, handler)
}

func TestGetBrandDistributorHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "testuser")

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	if err := db.Create(distributor).Error; err != nil {
		t.Fatalf("Failed to create distributor: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetBrandDistributorHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/brands/%d/distributors/%d", brand.Id, distributor.Id), nil)
	ctx := context.WithValue(req.Context(), "distributorId", distributor.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetDistributorLevelRewardsHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")

	reward := &model.DistributorLevelReward{
		BrandId:          brand.Id,
		Level:            1,
		RewardPercentage: 10.0,
	}
	if err := db.Create(reward).Error; err != nil {
		t.Fatalf("Failed to create reward config: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDistributorLevelRewardsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/distributors/level-rewards", nil)
	ctx := context.WithValue(req.Context(), "brandId", brand.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUpdateDistributorStatusHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "testuser")

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	if err := db.Create(distributor).Error; err != nil {
		t.Fatalf("Failed to create distributor: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateDistributorStatusHandler(svcCtx)

	reqBody := types.UpdateDistributorStatusReq{
		Status: "inactive",
		Reason: "Test reason",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/distributors/%d/status", distributor.Id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "distributorId", distributor.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestUpdateDistributorLevelHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "testuser")

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	if err := db.Create(distributor).Error; err != nil {
		t.Fatalf("Failed to create distributor: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := UpdateDistributorLevelHandler(svcCtx)

	reqBody := types.UpdateDistributorLevelReq{
		Level: 2,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/distributors/%d/level", distributor.Id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "distributorId", distributor.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestSetDistributorLevelRewardsHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")

	svcCtx := &svc.ServiceContext{DB: db}
	handler := SetDistributorLevelRewardsHandler(svcCtx)

	reqBody := types.SetDistributorLevelRewardsReq{
		Rewards: []types.SetDistributorLevelRewardReq{
			{Level: 1, RewardPercentage: 10.0},
		},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/brands/%d/distributors/level-rewards", brand.Id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "brandId", brand.Id)
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestApproveDistributorApplicationHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "applicant")

	application := &model.DistributorApplication{
		UserId:  user.Id,
		BrandId: brand.Id,
		Reason:  "Test reason",
		Status:  "pending",
	}
	if err := db.Create(application).Error; err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := ApproveDistributorApplicationHandler(svcCtx)

	reqBody := types.ApproveDistributorReq{
		Action: "approve",
	}
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

func TestGetDistributorApplicationHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "applicant")

	application := &model.DistributorApplication{
		UserId:  user.Id,
		BrandId: brand.Id,
		Reason:  "Test reason",
		Status:  "pending",
	}
	if err := db.Create(application).Error; err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetDistributorApplicationHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/distributors/applications/%d", application.Id), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", application.Id))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestGetBrandDistributorApplicationHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "applicant")

	application := &model.DistributorApplication{
		UserId:  user.Id,
		BrandId: brand.Id,
		Reason:  "Test reason",
		Status:  "pending",
	}
	if err := db.Create(application).Error; err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetBrandDistributorApplicationHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/brands/%d/distributors/applications/%d", brand.Id, application.Id), nil)
	req.SetPathValue("brandId", fmt.Sprintf("%d", brand.Id))
	req.SetPathValue("id", fmt.Sprintf("%d", application.Id))
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}

func TestTrackDistributorLinkHandler_WithDB(t *testing.T) {
	db := setupDistributorHandlerTestDB(t)
	brand := createTestBrand(t, db, "Test Brand")
	user := createTestUser(t, db, "testuser")

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	if err := db.Create(distributor).Error; err != nil {
		t.Fatalf("Failed to create distributor: %v", err)
	}
	campaign := &model.Campaign{
		BrandId:   brand.Id,
		Name:      "Track Campaign",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		Status:    "active",
	}
	if err := db.Create(campaign).Error; err != nil {
		t.Fatalf("Failed to create campaign: %v", err)
	}

	link := &model.DistributorLink{
		DistributorId: distributor.Id,
		CampaignId:    campaign.Id,
		LinkCode:      "testcode123",
	}
	if err := db.Create(link).Error; err != nil {
		t.Fatalf("Failed to create distributor link: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := TrackDistributorLinkHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/distributors/links/track?code=testcode123", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusInternalServerError, resp.Code)
}
