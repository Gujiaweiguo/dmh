package promoter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dmh/api/internal/handler/testutil"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupPromoterHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.SetupGormTestDB(t)

	err := db.AutoMigrate(
		&model.Brand{},
		&model.Campaign{},
		&model.Promoter{},
		&model.PromoterLink{},
		&model.PromoterReward{},
		&model.User{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	testutil.ClearTables(db, "promoters", "promoter_links", "promoter_rewards", "brands", "campaigns", "users")

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
		Username: usernamePrefix + fmt.Sprintf("_%d", time.Now().UnixNano()),
		Password: "test_hash",
		Phone:    fmt.Sprintf("138%08d", time.Now().UnixNano()%100000000),
		Role:     "brand_admin",
		Status:   "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestPromoter(t *testing.T, db *gorm.DB, user *model.User, brand *model.Brand) *model.Promoter {
	t.Helper()
	promoter := &model.Promoter{
		UserId:         user.Id,
		BrandId:        brand.Id,
		Status:         "active",
		Level:          "VIP",
		TotalOrders:    10,
		TotalRewards:   100.50,
		ConversionRate: 25.5,
		CampaignCount:  5,
	}
	if err := db.Create(promoter).Error; err != nil {
		t.Fatalf("Failed to create test promoter: %v", err)
	}
	return promoter
}

func TestGetPromoterListHandler(t *testing.T) {
	db := setupPromoterHandlerTestDB(t)
	brand := createTestBrand(t, db, "test_brand")
	user := createTestUser(t, db, "promoter_user")
	_ = createTestPromoter(t, db, user, brand)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPromoterListHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/promoter/list?page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.PromoterListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Promoters, 1)
	assert.Equal(t, "active", resp.Promoters[0].Status)
	assert.Equal(t, "VIP", resp.Promoters[0].Level)
}

func TestGetPromoterListHandlerWithFilters(t *testing.T) {
	db := setupPromoterHandlerTestDB(t)
	brand := createTestBrand(t, db, "test_brand")
	user1 := createTestUser(t, db, "promoter1")
	user2 := createTestUser(t, db, "promoter2")

	p1 := createTestPromoter(t, db, user1, brand)
	p1.Status = "inactive"
	db.Save(p1)

	_ = createTestPromoter(t, db, user2, brand)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPromoterListHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/promoter/list?status=active", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.PromoterListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Promoters, 1)
	assert.Equal(t, "active", resp.Promoters[0].Status)
}

func TestGetPromoterDetailHandler(t *testing.T) {
	db := setupPromoterHandlerTestDB(t)
	brand := createTestBrand(t, db, "test_brand")
	user := createTestUser(t, db, "promoter_user")
	promoter := createTestPromoter(t, db, user, brand)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPromoterDetailHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/promoter/detail/%d", promoter.Id), nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.PromoterDetailResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, promoter.Id, resp.Id)
	assert.Equal(t, "active", resp.Status)
	assert.Equal(t, "VIP", resp.Level)
	assert.Equal(t, int64(10), resp.TotalOrders)
}

func TestGetPromoterDetailHandlerNotFound(t *testing.T) {
	db := setupPromoterHandlerTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPromoterDetailHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/promoter/detail/99999", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
}

func TestGeneratePromoterLinkHandler(t *testing.T) {
	db := setupPromoterHandlerTestDB(t)
	brand := createTestBrand(t, db, "test_brand")
	user := createTestUser(t, db, "promoter_user")
	promoter := createTestPromoter(t, db, user, brand)

	campaign := &model.Campaign{
		BrandId:   brand.Id,
		Name:      "Test Campaign",
		Status:    "active",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	if err := db.Create(campaign).Error; err != nil {
		t.Fatalf("Failed to create test campaign: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GeneratePromoterLinkHandler(svcCtx)

	body := types.GeneratePromoterLinkReq{
		PromoterId: promoter.Id,
		CampaignId: campaign.Id,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/promoter/generate-link", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.GeneratePromoterLinkResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.NotEmpty(t, resp.LinkCode)
	assert.NotEmpty(t, resp.LinkUrl)
}

func TestGetPromoterRewardsHandler(t *testing.T) {
	db := setupPromoterHandlerTestDB(t)
	brand := createTestBrand(t, db, "test_brand")
	user := createTestUser(t, db, "promoter_user")
	promoter := createTestPromoter(t, db, user, brand)

	reward := &model.PromoterReward{
		PromoterId:  promoter.Id,
		Type:        "commission",
		Status:      "paid",
		Amount:      50.00,
		Description: "Test reward",
	}
	if err := db.Create(reward).Error; err != nil {
		t.Fatalf("Failed to create test reward: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPromoterRewardsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/promoter/rewards/%d?page=1&pageSize=10", promoter.Id), nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Logf("Response body: %s", rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.PromoterRewardsListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Rewards, 1)
	assert.Equal(t, "commission", resp.Rewards[0].Type)
	assert.Equal(t, "paid", resp.Rewards[0].Status)
	assert.Equal(t, 50.00, resp.Rewards[0].Amount)
}

func TestGetPromoterRewardsHandlerWithFilters(t *testing.T) {
	db := setupPromoterHandlerTestDB(t)
	brand := createTestBrand(t, db, "test_brand")
	user := createTestUser(t, db, "promoter_user")
	promoter := createTestPromoter(t, db, user, brand)

	r1 := &model.PromoterReward{
		PromoterId:  promoter.Id,
		Type:        "commission",
		Status:      "paid",
		Amount:      50.00,
		Description: "Commission 1",
	}
	r2 := &model.PromoterReward{
		PromoterId:  promoter.Id,
		Type:        "bonus",
		Status:      "pending",
		Amount:      100.00,
		Description: "Bonus 1",
	}
	if err := db.Create(r1).Error; err != nil {
		t.Fatalf("Failed to create test reward r1: %v", err)
	}
	if err := db.Create(r2).Error; err != nil {
		t.Fatalf("Failed to create test reward r2: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetPromoterRewardsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/promoter/rewards/%d?status=pending", promoter.Id), nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp types.PromoterRewardsListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Rewards, 1)
	assert.Equal(t, "pending", resp.Rewards[0].Status)
	assert.Equal(t, "bonus", resp.Rewards[0].Type)
}
