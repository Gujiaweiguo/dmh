package campaign

import (
	"context"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/testutil"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupCampaignTestDB(t *testing.T) *gorm.DB {
	db, _ := testutil.SetupMySQLTestDB(t)
	return db
}

func TestUpdateCampaignLogic_UpdateCampaign_Success(t *testing.T) {
	db := setupCampaignTestDB(t)

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

func TestUpdateCampaignLogic_UpdateCampaign_WithTimeFormat(t *testing.T) {
	db := setupCampaignTestDB(t)

	campaign := &model.Campaign{
		Name:        "原始活动",
		Description: "原始描述",
		FormFields:  `[]`,
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

	startTime := "2026-03-01T10:00:00"
	endTime := "2026-03-31T23:59:59"

	req := &types.UpdateCampaignReq{
		Id:        campaign.Id,
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUpdateCampaignLogic_UpdateCampaign_InvalidStartTime(t *testing.T) {
	db := setupCampaignTestDB(t)

	campaign := &model.Campaign{
		Name:        "原始活动",
		Description: "原始描述",
		FormFields:  `[]`,
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

	startTime := "invalid-time"

	req := &types.UpdateCampaignReq{
		Id:        campaign.Id,
		StartTime: &startTime,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateCampaignLogic_UpdateCampaign_InvalidEndTime(t *testing.T) {
	db := setupCampaignTestDB(t)

	campaign := &model.Campaign{
		Name:        "原始活动",
		Description: "原始描述",
		FormFields:  `[]`,
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

	endTime := "invalid-time"

	req := &types.UpdateCampaignReq{
		Id:      campaign.Id,
		EndTime: &endTime,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateCampaignLogic_UpdateCampaign_WithDistributionLevel(t *testing.T) {
	db := setupCampaignTestDB(t)

	campaign := &model.Campaign{
		Name:        "原始活动",
		Description: "原始描述",
		FormFields:  `[]`,
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

	level := 2
	enableDistribution := true

	req := &types.UpdateCampaignReq{
		Id:                 campaign.Id,
		DistributionLevel:  &level,
		EnableDistribution: &enableDistribution,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, resp.DistributionLevel)
	assert.True(t, resp.EnableDistribution)
}

func TestUpdateCampaignLogic_UpdateCampaign_InvalidDistributionLevel(t *testing.T) {
	db := setupCampaignTestDB(t)

	campaign := &model.Campaign{
		Name:        "原始活动",
		Description: "原始描述",
		FormFields:  `[]`,
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

	level := 5

	req := &types.UpdateCampaignReq{
		Id:                campaign.Id,
		DistributionLevel: &level,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateCampaignLogic_UpdateCampaign_WithFormFields(t *testing.T) {
	db := setupCampaignTestDB(t)

	campaign := &model.Campaign{
		Name:        "原始活动",
		Description: "原始描述",
		FormFields:  `[]`,
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

	formFields := []types.FormField{
		{Type: "text", Name: "name", Label: "姓名", Required: true},
		{Type: "phone", Name: "phone", Label: "手机号", Required: true},
	}

	req := &types.UpdateCampaignReq{
		Id:         campaign.Id,
		FormFields: formFields,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUpdateCampaignLogic_UpdateCampaign_WithPaymentConfig(t *testing.T) {
	db := setupCampaignTestDB(t)

	campaign := &model.Campaign{
		Name:        "原始活动",
		Description: "原始描述",
		FormFields:  `[]`,
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

	paymentConfig := `{"amount":100,"type":"wechat"}`

	req := &types.UpdateCampaignReq{
		Id:            campaign.Id,
		PaymentConfig: &paymentConfig,
	}

	resp, err := logic.UpdateCampaign(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteCampaignLogic_DeleteCampaign_Success(t *testing.T) {
	db := setupCampaignTestDB(t)

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

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewDeleteCampaignLogic(ctx, svcCtx)

	resp, err := logic.DeleteCampaign(999)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestSavePageConfigLogic_SavePageConfig_Success(t *testing.T) {
	db := setupCampaignTestDB(t)

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
