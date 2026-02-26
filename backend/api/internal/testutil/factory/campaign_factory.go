package factory

import (
	"dmh/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ErrBrandIdRequired 创建活动时缺少品牌ID的错误
var ErrBrandIdRequired = errors.New("brand_id is required for campaign creation")

// BrandFactory 品牌测试数据工厂
type BrandFactory struct {
	BaseFactory
}

// NewBrandFactory 创建品牌工厂实例
func NewBrandFactory() *BrandFactory {
	return &BrandFactory{}
}

// Build 创建一个品牌实例，使用合理的默认值
func (f *BrandFactory) Build() *model.Brand {
	return f.BuildWith(nil)
}

// BuildWith 创建一个品牌实例，允许覆盖指定字段
func (f *BrandFactory) BuildWith(overrides map[string]any) *model.Brand {
	brand := &model.Brand{
		Name:        "测试品牌_" + f.RandomSuffix(),
		Logo:        "https://example.com/logo.png",
		Description: "这是一个测试品牌",
		Status:      StatusActive,
	}

	// 应用覆盖值
	if overrides != nil {
		if v, ok := overrides["id"]; ok {
			brand.Id = v.(int64)
		}
		if v, ok := overrides["name"]; ok {
			brand.Name = v.(string)
		}
		if v, ok := overrides["logo"]; ok {
			brand.Logo = v.(string)
		}
		if v, ok := overrides["description"]; ok {
			brand.Description = v.(string)
		}
		if v, ok := overrides["status"]; ok {
			brand.Status = v.(string)
		}
	}

	return brand
}

// BuildList 创建多个品牌实例
func (f *BrandFactory) BuildList(count int) []*model.Brand {
	brands := make([]*model.Brand, count)
	for i := 0; i < count; i++ {
		brands[i] = f.Build()
	}
	return brands
}

// Create 创建并持久化品牌到数据库
func (f *BrandFactory) Create(db *gorm.DB) (*model.Brand, error) {
	return f.CreateWith(db, nil)
}

// CreateWith 创建并持久化品牌到数据库，允许覆盖指定字段
func (f *BrandFactory) CreateWith(db *gorm.DB, overrides map[string]any) (*model.Brand, error) {
	brand := f.BuildWith(overrides)
	if err := db.Create(brand).Error; err != nil {
		return nil, err
	}
	return brand, nil
}

// CreateList 创建并持久化多个品牌到数据库
func (f *BrandFactory) CreateList(db *gorm.DB, count int) ([]*model.Brand, error) {
	brands := f.BuildList(count)
	for _, brand := range brands {
		if err := db.Create(brand).Error; err != nil {
			return nil, err
		}
	}
	return brands, nil
}

// ========== Campaign Factory ==========

// CampaignFactory 营销活动测试数据工厂
type CampaignFactory struct {
	BaseFactory
}

// NewCampaignFactory 创建活动工厂实例
func NewCampaignFactory() *CampaignFactory {
	return &CampaignFactory{}
}

// Build 创建一个活动实例，使用合理的默认值
// 注意：BrandId 必须通过 overrides 提供
func (f *CampaignFactory) Build() *model.Campaign {
	return f.BuildWith(nil)
}

// BuildWith 创建一个活动实例，允许覆盖指定字段
func (f *CampaignFactory) BuildWith(overrides map[string]any) *model.Campaign {
	now := f.Now()
	campaign := &model.Campaign{
		BrandId:             0, // 必须通过 overrides 提供
		Name:                "测试活动_" + f.RandomSuffix(),
		Description:         "这是一个测试营销活动",
		FormFields:          `[{"type":"text","name":"name","label":"姓名","required":true},{"type":"phone","name":"phone","label":"手机号","required":true}]`,
		RewardRule:          10.00,
		StartTime:           now,
		EndTime:             now.Add(7 * 24 * time.Hour), // 7 天后结束
		Status:              CampaignStatusActive,
		EnableDistribution:  false,
		DistributionLevel:   1,
		PosterTemplateId:    1,
	}

	// 应用覆盖值
	if overrides != nil {
		if v, ok := overrides["id"]; ok {
			campaign.Id = v.(int64)
		}
		if v, ok := overrides["brand_id"]; ok {
			campaign.BrandId = v.(int64)
		}
		if v, ok := overrides["brandId"]; ok {
			campaign.BrandId = v.(int64)
		}
		if v, ok := overrides["name"]; ok {
			campaign.Name = v.(string)
		}
		if v, ok := overrides["description"]; ok {
			campaign.Description = v.(string)
		}
		if v, ok := overrides["form_fields"]; ok {
			campaign.FormFields = v.(string)
		}
		if v, ok := overrides["formFields"]; ok {
			campaign.FormFields = v.(string)
		}
		if v, ok := overrides["reward_rule"]; ok {
			campaign.RewardRule = v.(float64)
		}
		if v, ok := overrides["rewardRule"]; ok {
			campaign.RewardRule = v.(float64)
		}
		if v, ok := overrides["start_time"]; ok {
			campaign.StartTime = v.(time.Time)
		}
		if v, ok := overrides["startTime"]; ok {
			campaign.StartTime = v.(time.Time)
		}
		if v, ok := overrides["end_time"]; ok {
			campaign.EndTime = v.(time.Time)
		}
		if v, ok := overrides["endTime"]; ok {
			campaign.EndTime = v.(time.Time)
		}
		if v, ok := overrides["status"]; ok {
			campaign.Status = v.(string)
		}
		if v, ok := overrides["enable_distribution"]; ok {
			campaign.EnableDistribution = v.(bool)
		}
		if v, ok := overrides["enableDistribution"]; ok {
			campaign.EnableDistribution = v.(bool)
		}
		if v, ok := overrides["distribution_level"]; ok {
			campaign.DistributionLevel = v.(int)
		}
		if v, ok := overrides["distributionLevel"]; ok {
			campaign.DistributionLevel = v.(int)
		}
		if v, ok := overrides["distribution_rewards"]; ok {
			s := v.(string)
			campaign.DistributionRewards = &s
		}
		if v, ok := overrides["distributionRewards"]; ok {
			s := v.(string)
			campaign.DistributionRewards = &s
		}
		if v, ok := overrides["payment_config"]; ok {
			s := v.(string)
			campaign.PaymentConfig = &s
		}
		if v, ok := overrides["paymentConfig"]; ok {
			s := v.(string)
			campaign.PaymentConfig = &s
		}
		if v, ok := overrides["poster_template_id"]; ok {
			campaign.PosterTemplateId = v.(int64)
		}
		if v, ok := overrides["posterTemplateId"]; ok {
			campaign.PosterTemplateId = v.(int64)
		}
	}

	return campaign
}

// BuildList 创建多个活动实例
func (f *CampaignFactory) BuildList(count int) []*model.Campaign {
	campaigns := make([]*model.Campaign, count)
	for i := 0; i < count; i++ {
		campaigns[i] = f.Build()
	}
	return campaigns
}

// Create 创建并持久化活动到数据库
// 注意：必须提供 brandId（通常先创建 Brand）
func (f *CampaignFactory) Create(db *gorm.DB) (*model.Campaign, error) {
	return f.CreateWith(db, nil)
}

// CreateWith 创建并持久化活动到数据库，允许覆盖指定字段
func (f *CampaignFactory) CreateWith(db *gorm.DB, overrides map[string]any) (*model.Campaign, error) {
	campaign := f.BuildWith(overrides)

	// 验证必填字段
	if campaign.BrandId == 0 {
		return nil, ErrBrandIdRequired
	}

	if err := db.Create(campaign).Error; err != nil {
		return nil, err
	}
	return campaign, nil
}

// CreateList 创建并持久化多个活动到数据库
func (f *CampaignFactory) CreateList(db *gorm.DB, count int, brandId int64) ([]*model.Campaign, error) {
	campaigns := make([]*model.Campaign, count)
	for i := 0; i < count; i++ {
		campaign, err := f.CreateWith(db, map[string]any{"brand_id": brandId})
		if err != nil {
			return nil, err
		}
		campaigns[i] = campaign
	}
	return campaigns, nil
}

// ========== 便捷方法：创建特定类型活动 ==========

// BuildActive 创建活跃状态的活动（不持久化）
func (f *CampaignFactory) BuildActive(brandId int64) *model.Campaign {
	return f.BuildWith(map[string]any{
		"brand_id": brandId,
		"status":   CampaignStatusActive,
	})
}

// CreateActive 创建并持久化活跃状态的活动
func (f *CampaignFactory) CreateActive(db *gorm.DB, brandId int64) (*model.Campaign, error) {
	return f.CreateWith(db, map[string]any{
		"brand_id": brandId,
		"status":   CampaignStatusActive,
	})
}

// BuildPaused 创建暂停状态的活动（不持久化）
func (f *CampaignFactory) BuildPaused(brandId int64) *model.Campaign {
	return f.BuildWith(map[string]any{
		"brand_id": brandId,
		"status":   CampaignStatusPaused,
	})
}

// BuildEnded 创建已结束的活动（不持久化）
func (f *CampaignFactory) BuildEnded(brandId int64) *model.Campaign {
	now := f.Now()
	return f.BuildWith(map[string]any{
		"brand_id":   brandId,
		"status":     CampaignStatusEnded,
		"start_time": now.Add(-14 * 24 * time.Hour),
		"end_time":   now.Add(-7 * 24 * time.Hour),
	})
}

// BuildWithDistribution 创建启用分销的活动（不持久化）
func (f *CampaignFactory) BuildWithDistribution(brandId int64, levels int) *model.Campaign {
	rewards := `{"1":0.1,"2":0.05,"3":0.02}`
	return f.BuildWith(map[string]any{
		"brand_id":             brandId,
		"enable_distribution":  true,
		"distribution_level":   levels,
		"distribution_rewards": rewards,
	})
}

// ========== 错误定义已移至文件顶部 ==========

