package integration

import (
	"testing"
	"time"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DistributorIntegrationTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *DistributorIntegrationTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open("root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.UserBrand{},
		&model.Brand{},
		&model.Campaign{},
		&model.Order{},
		&model.Reward{},
		&model.Distributor{},
		&model.DistributorReward{},
		&model.DistributorApplication{},
		&model.UserBalance{},
		&model.DistributorLink{},
	)
	suite.Require().NoError(err)

	suite.db = db
}

func (suite *DistributorIntegrationTestSuite) SetupTest() {
	suite.Require().NoError(suite.db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error)
	suite.db.Exec("DELETE FROM distributor_links")
	suite.db.Exec("DELETE FROM distributor_rewards")
	suite.db.Exec("DELETE FROM distributor_applications")
	suite.db.Exec("DELETE FROM distributors")
	suite.db.Exec("DELETE FROM user_balances")
	suite.db.Exec("DELETE FROM rewards")
	suite.db.Exec("DELETE FROM orders")
	suite.db.Exec("DELETE FROM campaigns")
	suite.db.Exec("DELETE FROM user_brands")
	suite.db.Exec("DELETE FROM user_roles")
	suite.db.Exec("DELETE FROM roles")
	suite.db.Exec("DELETE FROM brands")
	suite.db.Exec("DELETE FROM users")
	suite.Require().NoError(suite.db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error)
	suite.createTestData()
}

func (suite *DistributorIntegrationTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *DistributorIntegrationTestSuite) createTestData() {
	// 使用 GORM 的 Unscoped().Delete 强制删除，确保幂等性
	// 先删除所有依赖数据，再创建

	// 1. 创建用户（先删除再创建）
	users := []model.User{
		{Id: 1, Username: "user1", Phone: "13800000001", Email: "user1@test.com", RealName: "用户1", Role: "participant", Status: "active"},
		{Id: 2, Username: "user2", Phone: "13800000002", Email: "user2@test.com", RealName: "用户2", Role: "participant", Status: "active"},
		{Id: 3, Username: "user3", Phone: "13800000003", Email: "user3@test.com", RealName: "用户3", Role: "participant", Status: "active"},
		{Id: 4, Username: "user4", Phone: "13800000004", Email: "user4@test.com", RealName: "用户4", Role: "participant", Status: "active"},
		{Id: 5, Username: "user5", Phone: "13800000005", Email: "user5@test.com", RealName: "用户5", Role: "participant", Status: "active"},
		{Id: 6, Username: "user6", Phone: "13800000006", Email: "user6@test.com", RealName: "用户6", Role: "participant", Status: "active"},
		{Id: 7, Username: "user7", Phone: "13800000007", Email: "user7@test.com", RealName: "用户7", Role: "participant", Status: "active"},
	}
	suite.db.Exec("DELETE FROM users WHERE id IN (1,2,3,4,5,6,7)")
	for _, user := range users {
		suite.Require().NoError(suite.db.Create(&user).Error)
	}

	// 2. 创建角色
	roles := []model.Role{
		{ID: 1, Name: "平台管理员", Code: "platform_admin"},
		{ID: 2, Name: "品牌管理员", Code: "brand_admin"},
		{ID: 3, Name: "参与者", Code: "participant"},
		{ID: 4, Name: "分销商", Code: "distributor"},
	}
	suite.db.Exec("DELETE FROM roles WHERE id IN (1,2,3,4)")
	for _, role := range roles {
		suite.Require().NoError(suite.db.Create(&role).Error)
	}

	// 3. 创建用户角色关系
	userRoles := []model.UserRole{
		{UserID: 1, RoleID: 3},
		{UserID: 2, RoleID: 3},
		{UserID: 3, RoleID: 3},
		{UserID: 4, RoleID: 3},
		{UserID: 5, RoleID: 3},
		{UserID: 6, RoleID: 3},
		{UserID: 7, RoleID: 3},
	}
	suite.db.Exec("DELETE FROM user_roles WHERE user_id IN (1,2,3,4,5,6,7)")
	for _, ur := range userRoles {
		suite.Require().NoError(suite.db.Create(&ur).Error)
	}

	// 4. 创建品牌
	brands := []model.Brand{
		{Id: 1, Name: "品牌A", Status: "active"},
	}
	suite.db.Exec("DELETE FROM brands WHERE id = 1")
	for _, brand := range brands {
		suite.Require().NoError(suite.db.Create(&brand).Error)
	}

	// 5. 创建用户品牌关系
	userBrands := []model.UserBrand{
		{UserId: 1, BrandId: 1},
		{UserId: 2, BrandId: 1},
		{UserId: 3, BrandId: 1},
		{UserId: 4, BrandId: 1},
		{UserId: 5, BrandId: 1},
		{UserId: 6, BrandId: 1},
		{UserId: 7, BrandId: 1},
	}
	suite.db.Exec("DELETE FROM user_brands WHERE user_id IN (1,2,3,4,5,6,7)")
	for _, ub := range userBrands {
		suite.Require().NoError(suite.db.Create(&ub).Error)
	}

	// 6. 创建活动
	distributionRewards := `{"1": 10.0, "2": 5.0, "3": 3.0}`
	campaigns := []model.Campaign{
		{
			Id:                  1,
			BrandId:             1,
			Name:                "测试活动A",
			Description:         "测试描述",
			RewardRule:          10.0,
			StartTime:           time.Now(),
			EndTime:             time.Now().AddDate(0, 1, 0),
			Status:              "active",
			EnableDistribution:  true,
			DistributionLevel:   3,
			DistributionRewards: &distributionRewards,
		},
	}
	suite.db.Exec("DELETE FROM campaigns WHERE id = 1")
	for _, campaign := range campaigns {
		suite.Require().NoError(suite.db.Create(&campaign).Error)
	}

	// 7. 创建用户余额
	suite.db.Exec("DELETE FROM user_balances WHERE user_id IN (1,2,3,4,5,6,7)")
	for i := 1; i <= 7; i++ {
		balance := &model.UserBalance{
			UserId:      int64(i),
			Balance:     0.0,
			TotalReward: 0.0,
			Version:     0,
		}
		suite.Require().NoError(suite.db.Create(balance).Error)
	}
}

func (suite *DistributorIntegrationTestSuite) TestDistributorRecordCreation() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ?", 1)

	distributor := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}

	err := suite.db.Create(distributor).Error
	suite.Require().NoError(err)

	var savedDistributor model.Distributor
	err = suite.db.Where("id = ?", distributor.Id).First(&savedDistributor).Error
	suite.Require().NoError(err)

	assert.Equal(suite.T(), int64(1), savedDistributor.UserId)
	assert.Equal(suite.T(), int64(1), savedDistributor.BrandId)
	assert.Equal(suite.T(), 1, savedDistributor.Level)
	assert.Equal(suite.T(), "active", savedDistributor.Status)
}

func (suite *DistributorIntegrationTestSuite) TestDistributorHierarchy() {
	d1 := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(d1)

	d2 := &model.Distributor{
		UserId:            2,
		BrandId:           1,
		Level:             2,
		Status:            "active",
		ParentId:          &d1.Id,
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(d2)

	d3 := &model.Distributor{
		UserId:            3,
		BrandId:           1,
		Level:             3,
		Status:            "active",
		ParentId:          &d2.Id,
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(d3)

	var parentDistributor model.Distributor
	suite.db.Where("id = ?", d1.Id).First(&parentDistributor)
	assert.NotNil(suite.T(), parentDistributor.Id)

	var childDistributor model.Distributor
	suite.db.Where("id = ?", d2.Id).Preload("Parent").First(&childDistributor)
	assert.NotNil(suite.T(), childDistributor.Parent)
	assert.Equal(suite.T(), d1.Id, childDistributor.Parent.Id)

	var subordinates []model.Distributor
	suite.db.Where("parent_id = ?", d1.Id).Find(&subordinates)
	assert.Len(suite.T(), subordinates, 1)
	assert.Equal(suite.T(), d2.Id, subordinates[0].Id)
}

func (suite *DistributorIntegrationTestSuite) TestDistributorRewardRecord() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ? AND brand_id = ?", 1, 1)
	distributor := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	order := &model.Order{
		Id:              1,
		CampaignId:      1,
		Phone:           "13800000001",
		FormData:        "{}",
		DistributorPath: "1",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          100.0,
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	now := time.Now()
	reward := &model.DistributorReward{
		DistributorId: distributor.Id,
		UserId:        1,
		OrderId:       order.Id,
		CampaignId:    1,
		Amount:        10.0,
		Level:         1,
		RewardRate:    10.0,
		Status:        "settled",
		SettledAt:     &now,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	suite.db.Create(reward)

	var savedReward model.DistributorReward
	err := suite.db.Where("order_id = ?", order.Id).First(&savedReward).Error
	suite.Require().NoError(err)

	assert.Equal(suite.T(), distributor.Id, savedReward.DistributorId)
	assert.Equal(suite.T(), 10.0, savedReward.Amount)
	assert.Equal(suite.T(), 1, savedReward.Level)
	assert.Equal(suite.T(), 10.0, savedReward.RewardRate)
	assert.Equal(suite.T(), "settled", savedReward.Status)
}

func (suite *DistributorIntegrationTestSuite) TestUserBalanceUpdate() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ?", 1)
	suite.db.Exec("DELETE FROM orders WHERE id = ?", 2)

	distributor := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	err := suite.db.Create(distributor).Error
	suite.Require().NoError(err)

	order := &model.Order{
		Id:              2,
		CampaignId:      1,
		Phone:           "13800000002",
		FormData:        "{}",
		DistributorPath: "1",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          100.0,
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	tx := suite.db.Begin()
	var balance model.UserBalance
	err = tx.Where("user_id = ?", 1).First(&balance).Error
	suite.Require().NoError(err)

	balance.Balance += 10.0
	balance.TotalReward += 10.0
	balance.Version += 1
	err = tx.Save(&balance).Error
	suite.Require().NoError(err)

	err = tx.Model(&model.Distributor{}).Where("id = ?", distributor.Id).
		Update("total_earnings", gorm.Expr("total_earnings + ?", 10.0)).Error
	suite.Require().NoError(err)

	err = tx.Commit().Error
	suite.Require().NoError(err)

	var updatedBalance model.UserBalance
	suite.db.Where("user_id = ?", 1).First(&updatedBalance)
	assert.Equal(suite.T(), 10.0, updatedBalance.Balance)
	assert.Equal(suite.T(), 10.0, updatedBalance.TotalReward)

	var updatedDistributor model.Distributor
	suite.db.Where("id = ?", distributor.Id).First(&updatedDistributor)
	assert.Equal(suite.T(), 10.0, updatedDistributor.TotalEarnings)
}

func (suite *DistributorIntegrationTestSuite) TestDistributorLinkCreation() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ? AND brand_id = ?", 1, 1)
	distributor := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	link := &model.DistributorLink{
		DistributorId: distributor.Id,
		CampaignId:    1,
		LinkCode:      "TEST123",
		ClickCount:    0,
		OrderCount:    0,
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	suite.db.Create(link)

	var savedLink model.DistributorLink
	err := suite.db.Where("link_code = ?", "TEST123").First(&savedLink).Error
	suite.Require().NoError(err)

	assert.Equal(suite.T(), distributor.Id, savedLink.DistributorId)
	assert.Equal(suite.T(), int64(1), savedLink.CampaignId)
	assert.Equal(suite.T(), "TEST123", savedLink.LinkCode)
	assert.Equal(suite.T(), "active", savedLink.Status)
}

func (suite *DistributorIntegrationTestSuite) TestOrderWithDistributorPath() {
	suite.db.Exec("DELETE FROM distributors WHERE brand_id = ?", 1)
	d1 := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(d1)

	d2 := &model.Distributor{
		UserId:            2,
		BrandId:           1,
		Level:             2,
		Status:            "active",
		ParentId:          &d1.Id,
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(d2)

	d3 := &model.Distributor{
		UserId:            3,
		BrandId:           1,
		Level:             3,
		Status:            "active",
		ParentId:          &d2.Id,
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(d3)

	order := &model.Order{
		Id:              3,
		CampaignId:      1,
		Phone:           "13800000004",
		FormData:        "{}",
		DistributorPath: "1,2,3",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          200.0,
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	var savedOrder model.Order
	err := suite.db.Where("id = ?", order.Id).First(&savedOrder).Error
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "1,2,3", savedOrder.DistributorPath)

	rewards := []model.DistributorReward{
		{
			DistributorId: d1.Id,
			UserId:        1,
			OrderId:       order.Id,
			CampaignId:    1,
			Amount:        20.0,
			Level:         1,
			RewardRate:    10.0,
			Status:        "settled",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			DistributorId: d2.Id,
			UserId:        2,
			OrderId:       order.Id,
			CampaignId:    1,
			Amount:        10.0,
			Level:         2,
			RewardRate:    5.0,
			Status:        "settled",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			DistributorId: d3.Id,
			UserId:        3,
			OrderId:       order.Id,
			CampaignId:    1,
			Amount:        6.0,
			Level:         3,
			RewardRate:    3.0,
			Status:        "settled",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
	for _, reward := range rewards {
		suite.db.Create(&reward)
	}

	var savedRewards []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Find(&savedRewards)
	assert.Len(suite.T(), savedRewards, 3)
}

func (suite *DistributorIntegrationTestSuite) TestDistributorStatusRestriction() {
	distributor := &model.Distributor{
		UserId:            4,
		BrandId:           1,
		Level:             1,
		Status:            "suspended",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	order := &model.Order{
		Id:              4,
		CampaignId:      1,
		Phone:           "13800000005",
		FormData:        "{}",
		DistributorPath: "4",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          100.0,
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	var rewards []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Find(&rewards)
	assert.Len(suite.T(), rewards, 0)
}

func (suite *DistributorIntegrationTestSuite) TestDistributionDisabledCampaign() {
	distributionRewardsDisabled := "{}"

	campaignNoDist := &model.Campaign{
		Id:                  2,
		BrandId:             1,
		Name:                "不启用分销的活动",
		Description:         "测试描述",
		RewardRule:          10.0,
		StartTime:           time.Now(),
		EndTime:             time.Now().AddDate(0, 1, 0),
		Status:              "active",
		EnableDistribution:  false,
		DistributionLevel:   1,
		DistributionRewards: &distributionRewardsDisabled,
	}
	suite.db.Create(campaignNoDist)

	distributor := &model.Distributor{
		UserId:            5,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	order := &model.Order{
		Id:              5,
		CampaignId:      2,
		Phone:           "13800000006",
		FormData:        "{}",
		DistributorPath: "5",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          100.0,
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	var rewards []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Find(&rewards)
	assert.Len(suite.T(), rewards, 0)
}

func TestDistributorIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(DistributorIntegrationTestSuite))
}
