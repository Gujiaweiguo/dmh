package factory

import (
	"dmh/model"

	"gorm.io/gorm"
)

// Fixtures 提供统一的测试数据工厂入口
// 用于测试套件中快速创建和管理测试数据
type Fixtures struct {
	db *gorm.DB

	User     *UserFactory
	Brand    *BrandFactory
	Campaign *CampaignFactory
	Order    *OrderFactory
}

// NewFixtures 创建测试数据工厂集合
func NewFixtures(db *gorm.DB) *Fixtures {
	return &Fixtures{
		db:       db,
		User:     NewUserFactory(),
		Brand:    NewBrandFactory(),
		Campaign: NewCampaignFactory(),
		Order:    NewOrderFactory(),
	}
}

// ========== 快捷方法：创建完整的测试数据链 ==========

// SetupUserWithBrand 创建用户和品牌，并建立关联
func (f *Fixtures) SetupUserWithBrand() (*model.User, *model.Brand, *model.UserBrand, error) {
	// 创建用户
	user, err := f.User.CreateBrandAdmin(f.db)
	if err != nil {
		return nil, nil, nil, err
	}

	// 创建品牌
	brand, err := f.Brand.Create(f.db)
	if err != nil {
		return nil, nil, nil, err
	}

	// 创建关联
	userBrand := &model.UserBrand{
		UserId:  user.Id,
		BrandId: brand.Id,
	}
	if err := f.db.Create(userBrand).Error; err != nil {
		return nil, nil, nil, err
	}

	return user, brand, userBrand, nil
}

// SetupCampaignWithBrand 创建品牌和活动
func (f *Fixtures) SetupCampaignWithBrand() (*model.Brand, *model.Campaign, error) {
	// 创建品牌
	brand, err := f.Brand.Create(f.db)
	if err != nil {
		return nil, nil, err
	}

	// 创建活动
	campaign, err := f.Campaign.CreateActive(f.db, brand.Id)
	if err != nil {
		return nil, nil, err
	}

	return brand, campaign, nil
}

// SetupFullOrderChain 创建完整的订单测试数据链：用户 -> 品牌 -> 活动 -> 订单
func (f *Fixtures) SetupFullOrderChain() (*model.User, *model.Brand, *model.Campaign, *model.Order, error) {
	// 创建用户
	user, err := f.User.Create(f.db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// 创建品牌
	brand, err := f.Brand.Create(f.db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// 创建活动
	campaign, err := f.Campaign.CreateActive(f.db, brand.Id)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// 创建订单
	order, err := f.Order.CreatePaid(f.db, campaign.Id)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return user, brand, campaign, order, nil
}

// SetupVerifiedOrder 创建已核销的订单完整链
func (f *Fixtures) SetupVerifiedOrder(verifierId int64) (*model.Brand, *model.Campaign, *model.Order, error) {
	// 创建品牌
	brand, err := f.Brand.Create(f.db)
	if err != nil {
		return nil, nil, nil, err
	}

	// 创建活动
	campaign, err := f.Campaign.CreateActive(f.db, brand.Id)
	if err != nil {
		return nil, nil, nil, err
	}

	// 创建已核销订单
	order, err := f.Order.CreateVerified(f.db, campaign.Id, verifierId)
	if err != nil {
		return nil, nil, nil, err
	}

	return brand, campaign, order, nil
}

// ========== 批量创建 ==========

// CreateTestUsers 批量创建测试用户
func (f *Fixtures) CreateTestUsers(count int) ([]*model.User, error) {
	return f.User.CreateList(f.db, count)
}

// CreateTestCampaigns 批量创建测试活动（需先有品牌）
func (f *Fixtures) CreateTestCampaigns(count int, brandId int64) ([]*model.Campaign, error) {
	return f.Campaign.CreateList(f.db, count, brandId)
}

// CreateTestOrders 批量创建测试订单（需先有活动）
func (f *Fixtures) CreateTestOrders(count int, campaignId int64) ([]*model.Order, error) {
	return f.Order.CreateList(f.db, count, campaignId)
}

// ========== 清理方法 ==========

// Cleanup 清理所有测试数据（按外键依赖顺序）
func (f *Fixtures) Cleanup() error {
	tables := []string{
		"orders",
		"campaigns",
		"user_brands",
		"brands",
		"users",
	}

	for _, table := range tables {
		if err := f.db.Exec("DELETE FROM " + table).Error; err != nil {
			return err
		}
	}
	return nil
}

// CleanupOrders 清理订单数据
func (f *Fixtures) CleanupOrders() error {
	return f.db.Exec("DELETE FROM orders").Error
}

// CleanupCampaigns 清理活动数据（会级联删除订单）
func (f *Fixtures) CleanupCampaigns() error {
	// 先删除订单
	if err := f.db.Exec("DELETE FROM orders").Error; err != nil {
		return err
	}
	return f.db.Exec("DELETE FROM campaigns").Error
}

// CleanupBrands 清理品牌数据（会级联删除活动和订单）
func (f *Fixtures) CleanupBrands() error {
	// 按依赖顺序删除
	if err := f.db.Exec("DELETE FROM orders").Error; err != nil {
		return err
	}
	if err := f.db.Exec("DELETE FROM campaigns").Error; err != nil {
		return err
	}
	if err := f.db.Exec("DELETE FROM user_brands").Error; err != nil {
		return err
	}
	return f.db.Exec("DELETE FROM brands").Error
}

// CleanupUsers 清理用户数据
func (f *Fixtures) CleanupUsers() error {
	// 先清理关联数据
	if err := f.db.Exec("DELETE FROM user_brands").Error; err != nil {
		return err
	}
	return f.db.Exec("DELETE FROM users").Error
}
