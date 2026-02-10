package campaign

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCreateCampaignTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Campaign{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestCreateCampaignLogic_CreateCampaign_Success(t *testing.T) {
	db := setupCreateCampaignTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	campaignLogic := NewCreateCampaignLogic(ctx, svcCtx)

	now := time.Now()
	req := &types.CreateCampaignReq{
		BrandId:             1,
		Name:                "Test Campaign",
		Description:         "Test Description",
		RewardRule:          10.0,
		StartTime:           now.Add(1 * time.Hour).Format("2006-01-02T15:04:05"),
		EndTime:             now.Add(24 * time.Hour).Format("2006-01-02T15:04:05"),
		EnableDistribution:  true,
		DistributionLevel:   2,
		DistributionRewards: `[{"level":1,"rate":0.1},{"level":2,"rate":0.05}]`,
		FormFields: []types.FormField{
			{Type: "text", Name: "name", Label: "姓名", Required: true},
			{Type: "phone", Name: "phone", Label: "手机号", Required: true},
		},
	}

	resp, err := campaignLogic.CreateCampaign(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotZero(t, resp.Id)
	assert.Equal(t, req.Name, resp.Name)
	assert.Equal(t, req.BrandId, resp.BrandId)
	assert.Equal(t, "active", resp.Status)
	assert.Equal(t, 2, resp.DistributionLevel)
}

func TestCreateCampaignLogic_InvalidTimeFormat(t *testing.T) {
	db := setupCreateCampaignTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	campaignLogic := NewCreateCampaignLogic(ctx, svcCtx)

	req := &types.CreateCampaignReq{
		BrandId:     1,
		Name:        "Test Campaign",
		StartTime:   "invalid-time",
		EndTime:     "invalid-time",
		RewardRule:  10.0,
	}

	resp, err := campaignLogic.CreateCampaign(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Time format error")
}

func TestCreateCampaignLogic_InvalidDistributionLevel(t *testing.T) {
	db := setupCreateCampaignTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	campaignLogic := NewCreateCampaignLogic(ctx, svcCtx)

	now := time.Now()
	req := &types.CreateCampaignReq{
		BrandId:           1,
		Name:              "Test Campaign",
		StartTime:         now.Add(1 * time.Hour).Format("2006-01-02T15:04:05"),
		EndTime:           now.Add(24 * time.Hour).Format("2006-01-02T15:04:05"),
		RewardRule:        10.0,
		DistributionLevel: 5,
	}

	resp, err := campaignLogic.CreateCampaign(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Distribution level must be between 1 and 3")
}

func TestCreateCampaignLogic_DefaultValues(t *testing.T) {
	db := setupCreateCampaignTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	campaignLogic := NewCreateCampaignLogic(ctx, svcCtx)

	now := time.Now()
	req := &types.CreateCampaignReq{
		BrandId:     1,
		Name:        "Test Campaign",
		StartTime:   now.Add(1 * time.Hour).Format("2006-01-02T15:04:05"),
		EndTime:     now.Add(24 * time.Hour).Format("2006-01-02T15:04:05"),
		RewardRule:  10.0,
	}

	resp, err := campaignLogic.CreateCampaign(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.DistributionLevel)
	assert.Equal(t, int64(1), resp.PosterTemplateId)
}
