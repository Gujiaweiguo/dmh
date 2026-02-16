package reward

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

func setupRewardHandlerTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Reward{}, &model.DistributorReward{}, &model.UserBalance{}, &model.User{}, &model.Order{}, &model.Campaign{}, &model.Brand{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestRewardHandlersConstruct(t *testing.T) {
	assert.NotNil(t, GetRewardsHandler(nil))
	assert.NotNil(t, GetBalanceHandler(nil))
}

func TestGetRewardsHandler_Success(t *testing.T) {
	db := setupRewardHandlerTestDB(t)

	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{Name: "Test Campaign", BrandId: brand.Id, Status: "active"}
	db.Create(campaign)

	user := &model.User{Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	order := &model.Order{CampaignId: campaign.Id, Phone: "13800138000", Amount: 100.00, PayStatus: "paid", Status: "active"}
	db.Create(order)

	reward := &model.Reward{UserId: user.Id, OrderId: order.Id, CampaignId: campaign.Id, Amount: 10.00, Status: "pending"}
	db.Create(reward)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetRewardsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rewards?page=1&pageSize=10", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetBalanceHandler_Success(t *testing.T) {
	db := setupRewardHandlerTestDB(t)

	user := &model.User{Username: "testuser", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	balance := &model.UserBalance{UserId: user.Id, Balance: 100.00, TotalReward: 500.00}
	db.Create(balance)

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetBalanceHandler(svcCtx)

	userId := user.Id
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/rewards/balance/%d", userId), nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetRewardsHandler_Error(t *testing.T) {
	db := setupRewardHandlerTestDB(t)
	if err := db.Migrator().DropTable(&model.DistributorReward{}); err != nil {
		t.Fatalf("failed to drop DistributorReward table: %v", err)
	}

	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetRewardsHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rewards?userId=1", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func TestGetBalanceHandler_MissingUserId(t *testing.T) {
	db := setupRewardHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetBalanceHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rewards/balance/", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetBalanceHandler_InvalidUserId(t *testing.T) {
	db := setupRewardHandlerTestDB(t)
	svcCtx := &svc.ServiceContext{DB: db}
	handler := GetBalanceHandler(svcCtx)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rewards/balance/abc", nil)
	resp := httptest.NewRecorder()

	handler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
