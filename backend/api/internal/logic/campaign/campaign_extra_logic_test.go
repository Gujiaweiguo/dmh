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

func setupCampaignTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Campaign{}, &model.PageConfig{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func cleanupCampaignTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

func TestUpdateCampaignLogic_UpdateCampaign_Success(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	campaign := &model.Campaign{
		Name:        "原始活动",
		Description: "原始描述",
		FormFields:  `[{"type":"text","name":"name","label":"姓名"}]`,
		RewardRule:  10.00,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}
	db.Create(campaign)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateCampaignLogic(ctx, svcCtx)

	name := "更新后的活动"
	description := "更新后的描述"
	rewardRule := 20.00

	req := &types.UpdateCampaignReq{
		Id:          campaign.Id,
		Name:        &name,
		Description: &description,
		RewardRule:  &rewardRule,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "更新后的活动", resp.Name)
	assert.Equal(t, "更新后的描述", resp.Description)
	assert.Equal(t, 20.00, resp.RewardRule)
}

func TestUpdateCampaignLogic_UpdateCampaign_NotFound(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateCampaignLogic(ctx, svcCtx)

	name := "不存在的活动"
	description := "测试描述"
	rewardRule := 10.00

	req := &types.UpdateCampaignReq{
		Id:          999,
		Name:        &name,
		Description: &description,
		RewardRule:  &rewardRule,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeleteCampaignLogic_DeleteCampaign_Success(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	campaign := &model.Campaign{
		Name:        "待删除活动",
		Description: "测试描述",
		RewardRule:  10.00,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}
	db.Create(campaign)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewDeleteCampaignLogic(ctx, svcCtx)

	resp, err := logic.DeleteCampaign(campaign.Id)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	var deletedCampaign model.Campaign
	result := db.Unscoped().First(&deletedCampaign, campaign.Id)
	assert.NoError(t, result.Error)
	assert.NotNil(t, deletedCampaign.DeletedAt)
}

func TestDeleteCampaignLogic_DeleteCampaign_NotFound(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewDeleteCampaignLogic(ctx, svcCtx)

	resp, err := logic.DeleteCampaign(999)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestSavePageConfigLogic_SavePageConfig_Success(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	campaign := &model.Campaign{
		Name:        "测试活动",
		Description: "测试描述",
		RewardRule:  10.00,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}
	db.Create(campaign)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewSavePageConfigLogic(ctx, svcCtx)

	components := []map[string]interface{}{
		{"type": "banner", "content": "测试内容"},
	}
	theme := map[string]interface{}{
		"primaryColor": "#ff0000",
	}

	req := &types.PageConfigReq{
		Id:         campaign.Id,
		Components: components,
		Theme:      theme,
	}

	resp, err := logic.SavePageConfig(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	var pageConfig model.PageConfig
	result := db.Where("campaign_id = ?", campaign.Id).First(&pageConfig)
	assert.NoError(t, result.Error)
	assert.Equal(t, campaign.Id, pageConfig.CampaignId)
}

func TestGetPageConfigLogic_GetPageConfig_Success(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	campaign := &model.Campaign{
		Name:        "测试活动",
		Description: "测试描述",
		RewardRule:  10.00,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}
	db.Create(campaign)

	pageConfig := &model.PageConfig{
		CampaignId: campaign.Id,
		Components: `[{"type":"banner","content":"测试内容"}]`,
		Theme:      `{"primaryColor":"#ff0000"}`,
	}
	db.Create(pageConfig)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetPageConfigLogic(ctx, svcCtx)

	req := &types.GetPageConfigReq{
		Id: campaign.Id,
	}

	resp, err := logic.GetPageConfig(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, campaign.Id, resp.Id)
	assert.Equal(t, campaign.Id, resp.CampaignId)
	assert.Len(t, resp.Components, 1)
}

func TestGetPageConfigLogic_GetPageConfig_NotFound(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetPageConfigLogic(ctx, svcCtx)

	req := &types.GetPageConfigReq{
		Id: 999,
	}

	resp, err := logic.GetPageConfig(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(999), resp.CampaignId)
	assert.Empty(t, resp.Components)
}

func TestGetPaymentQrcodeLogic_GetPaymentQrcode_Success(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	campaign := &model.Campaign{
		Name:        "测试活动",
		Description: "测试描述",
		RewardRule:  10.00,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}
	db.Create(campaign)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetPaymentQrcodeLogic(ctx, svcCtx)

	req := types.GetPaymentQrcodeReq{
		Id: campaign.Id,
	}

	resp, err := logic.GetPaymentQrcode(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 100.00, resp.Amount)
	assert.Equal(t, "测试活动", resp.CampaignName)
	assert.NotEmpty(t, resp.QrcodeUrl)
}

func TestGetPaymentQrcodeLogic_GetPaymentQrcode_NotFound(t *testing.T) {
	db := setupCampaignTestDB(t)
	defer cleanupCampaignTestDB(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetPaymentQrcodeLogic(ctx, svcCtx)

	req := types.GetPaymentQrcodeReq{
		Id: 999,
	}

	resp, err := logic.GetPaymentQrcode(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
