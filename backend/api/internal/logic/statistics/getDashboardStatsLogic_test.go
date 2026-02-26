package statistics

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
)

func setupStatsTestDB(t *testing.T) {
	db, _ := testutil.SetupMySQLTestDB(t)
	return db
}

func TestGetDashboardStats_Success(t *testing.T) {
	db := setupStatsTestDB(t)

	// Create test data
	user := &model.User{Id: 100, Username: "testuser", Phone: "13800138000"}
	db.Create(user)

	brand := &model.Brand{Id: 1, Name: "Test Brand", Status: "active"}
	db.Create(brand)

	campaign := &model.Campaign{
		Id:          1,
		Name:        "Test Campaign",
		Description: "Test Description",
		FormFields:  `[]`,
		RewardRule:  10.00,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}
	db.Create(campaign)

	order := &model.Order{
		Id:         1,
		UserId:     100,
		CampaignId: 1,
		BrandId:    1,
		Amount:     100.00,
		Status:     "paid",
		CreatedAt:  time.Now(),
	}
	db.Create(order)

	// Test
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(context.Background(), svcCtx)

	req := &types.GetDashboardStatsReq{
		StartDate: time.Now().Add(-30 * 24 * time.Hour).Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
	}

	resp, err := logic.GetDashboardStats(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
