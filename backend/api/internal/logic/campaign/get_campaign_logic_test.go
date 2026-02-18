package campaign

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/testutil"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupGetCampaignTestDB(t *testing.T) *gorm.DB {
	db, _ := testutil.SetupMySQLTestDB(t)
	return db
}

func createTestCampaignForGet(t *testing.T, db *gorm.DB, brandId int64, name, status string) *model.Campaign {
	now := time.Now()
	campaign := &model.Campaign{
		BrandId:     brandId,
		Name:        name,
		Description: "Test Description",
		RewardRule:  10.0,
		StartTime:   now,
		EndTime:     now.Add(24 * time.Hour),
		Status:      status,
	}
	if err := db.Create(campaign).Error; err != nil {
		t.Fatalf("Failed to create test campaign: %v", err)
	}
	return campaign
}

func TestGetCampaignLogic_GetCampaign_Success(t *testing.T) {
	db := setupGetCampaignTestDB(t)

	campaign := createTestCampaignForGet(t, db, 1, "Test Campaign", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetCampaignLogic(ctx, svcCtx)

	req := &types.GetCampaignReq{Id: campaign.Id}
	resp, err := logic.GetCampaign(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, campaign.Id, resp.Id)
	assert.Equal(t, campaign.Name, resp.Name)
	assert.Equal(t, campaign.BrandId, resp.BrandId)
}

func TestGetCampaignLogic_CampaignNotFound(t *testing.T) {
	db := setupGetCampaignTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetCampaignLogic(ctx, svcCtx)

	req := &types.GetCampaignReq{Id: 999}
	resp, err := logic.GetCampaign(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetCampaignsLogic_GetCampaigns_Success(t *testing.T) {
	db := setupGetCampaignTestDB(t)

	for i := 0; i < 5; i++ {
		createTestCampaignForGet(t, db, 1, fmt.Sprintf("Campaign %d", i), "active")
	}

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetCampaignsLogic(ctx, svcCtx)

	req := &types.GetCampaignsReq{
		Page:     1,
		PageSize: 10,
	}
	resp, err := logic.GetCampaigns(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(5), resp.Total)
	assert.Len(t, resp.Campaigns, 5)
}

func TestGetCampaignsLogic_WithStatusFilter(t *testing.T) {
	db := setupGetCampaignTestDB(t)

	for i := 0; i < 3; i++ {
		createTestCampaignForGet(t, db, 1, fmt.Sprintf("Active %d", i), "active")
	}
	for i := 0; i < 2; i++ {
		createTestCampaignForGet(t, db, 1, fmt.Sprintf("Paused %d", i), "paused")
	}

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetCampaignsLogic(ctx, svcCtx)

	req := &types.GetCampaignsReq{
		Status:   "active",
		Page:     1,
		PageSize: 10,
	}
	resp, err := logic.GetCampaigns(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(3), resp.Total)
	assert.Len(t, resp.Campaigns, 3)
}

func TestGetCampaignsLogic_EmptyResult(t *testing.T) {
	db := setupGetCampaignTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetCampaignsLogic(ctx, svcCtx)

	req := &types.GetCampaignsReq{
		Page:     1,
		PageSize: 10,
	}
	resp, err := logic.GetCampaigns(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Campaigns, 0)
}
