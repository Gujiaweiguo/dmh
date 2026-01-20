package distributor

import (
	"context"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MultiLevelRewardLogicTestSuite struct {
	suite.Suite
	db     *gorm.DB
	svcCtx *svc.ServiceContext
}

func (suite *MultiLevelRewardLogicTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.UserBrand{},
		&model.Brand{},
		&model.Campaign{},
		&model.Order{},
		&model.Distributor{},
		&model.DistributorReward{},
		&model.UserBalance{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.svcCtx = &svc.ServiceContext{DB: db}
	suite.createTestData()
}

func (suite *MultiLevelRewardLogicTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *MultiLevelRewardLogicTestSuite) createTestData() {
	roles := []model.Role{
		{ID: 1, Name: "分销商", Code: "distributor"},
	}
	for _, role := range roles {
		suite.db.Create(&role)
	}

	users := []model.User{
		{Id: 1, Username: "user1", Phone: "13800000001", Email: "user1@test.com", RealName: "用户1", Role: "distributor", Status: "active"},
		{Id: 2, Username: "user2", Phone: "13800000002", Email: "user2@test.com", RealName: "用户2", Role: "distributor", Status: "active"},
		{Id: 3, Username: "user3", Phone: "13800000003", Email: "user3@test.com", RealName: "用户3", Role: "distributor", Status: "active"},
		{Id: 4, Username: "user4", Phone: "13800000004", Email: "user4@test.com", RealName: "用户4", Role: "distributor", Status: "active"},
		{Id: 5, Username: "user5", Phone: "13800000005", Email: "user5@test.com", RealName: "用户5", Role: "participant", Status: "active"},
	}
	for _, user := range users {
		suite.db.Create(&user)
	}

	userRoles := []model.UserRole{
		{UserID: 1, RoleID: 1},
		{UserID: 2, RoleID: 1},
		{UserID: 3, RoleID: 1},
		{UserID: 4, RoleID: 1},
	}
	for _, ur := range userRoles {
		suite.db.Create(&ur)
	}

	brands := []model.Brand{
		{Id: 1, Name: "品牌A", Status: "active"},
	}
	for _, brand := range brands {
		suite.db.Create(&brand)
	}

	userBrands := []model.UserBrand{
		{UserId: 1, BrandId: 1},
		{UserId: 2, BrandId: 1},
		{UserId: 3, BrandId: 1},
		{UserId: 4, BrandId: 1},
	}
	for _, ub := range userBrands {
		suite.db.Create(&ub)
	}

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
	for _, campaign := range campaigns {
		suite.db.Create(&campaign)
	}

	for i := 1; i <= 5; i++ {
		balance := &model.UserBalance{
			UserId:      int64(i),
			Balance:     0.0,
			TotalReward: 0.0,
			Version:     0,
		}
		suite.db.Create(balance)
	}
}

func (suite *MultiLevelRewardLogicTestSuite) TestCalculateAndDistributeRewards_SingleLevel() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

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
		Phone:           "13800000005",
		DistributorPath: "1",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          100.0,
		PaidAt:          timePtr(time.Now()),
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	err := logic.CalculateAndDistributeRewards(order.Id, order.CampaignId, 5, 0, order.Amount)
	suite.Require().NoError(err)

	var rewards []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Find(&rewards)
	assert.Len(suite.T(), rewards, 1)
	assert.Equal(suite.T(), 10.0, rewards[0].Amount)
	assert.Equal(suite.T(), 1, rewards[0].Level)
	assert.Equal(suite.T(), 10.0, rewards[0].RewardRate)
}

func (suite *MultiLevelRewardLogicTestSuite) TestCalculateAndDistributeRewards_ThreeLevel() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

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
		Id:              2,
		CampaignId:      1,
		Phone:           "13800000005",
		DistributorPath: "1,2,3",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          200.0,
		PaidAt:          timePtr(time.Now()),
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	err := logic.CalculateAndDistributeRewards(order.Id, order.CampaignId, 5, 0, order.Amount)
	suite.Require().NoError(err)

	var rewards []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Order("level").Find(&rewards)
	assert.Len(suite.T(), rewards, 3)

	assert.Equal(suite.T(), 20.0, rewards[0].Amount)
	assert.Equal(suite.T(), 1, rewards[0].Level)
	assert.Equal(suite.T(), 10.0, rewards[0].RewardRate)

	assert.Equal(suite.T(), 10.0, rewards[1].Amount)
	assert.Equal(suite.T(), 2, rewards[1].Level)
	assert.Equal(suite.T(), 5.0, rewards[1].RewardRate)

	assert.Equal(suite.T(), 6.0, rewards[2].Amount)
	assert.Equal(suite.T(), 3, rewards[2].Level)
	assert.Equal(suite.T(), 3.0, rewards[2].RewardRate)
}

func (suite *MultiLevelRewardLogicTestSuite) TestCalculateAndDistributeRewards_Idempotency() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

	distributor := &model.Distributor{
		UserId:            4,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	order := &model.Order{
		Id:              3,
		CampaignId:      1,
		Phone:           "13800000005",
		DistributorPath: "4",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          100.0,
		PaidAt:          timePtr(time.Now()),
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	err := logic.CalculateAndDistributeRewards(order.Id, order.CampaignId, 5, 0, order.Amount)
	suite.Require().NoError(err)

	var rewardsFirst []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Find(&rewardsFirst)
	assert.Len(suite.T(), rewardsFirst, 1)

	var balanceAfterFirst model.UserBalance
	suite.db.Where("user_id = ?", 4).First(&balanceAfterFirst)

	err = logic.CalculateAndDistributeRewards(order.Id, order.CampaignId, 5, 0, order.Amount)
	suite.Require().NoError(err)

	var rewardsSecond []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Find(&rewardsSecond)
	assert.Len(suite.T(), rewardsSecond, 1)

	assert.Equal(suite.T(), rewardsFirst[0].Id, rewardsSecond[0].Id)

	var balanceAfterSecond model.UserBalance
	suite.db.Where("user_id = ?", 4).First(&balanceAfterSecond)
	assert.Equal(suite.T(), balanceAfterFirst.Balance, balanceAfterSecond.Balance)
}

func (suite *MultiLevelRewardLogicTestSuite) TestCalculateAndDistributeRewards_DisabledDistribution() {
	distributionRewards := ""

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
		DistributionRewards: &distributionRewards,
	}
	suite.db.Create(campaignNoDist)

	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

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
		Id:              4,
		CampaignId:      2,
		Phone:           "13800000005",
		DistributorPath: "1",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          100.0,
		PaidAt:          timePtr(time.Now()),
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	err := logic.CalculateAndDistributeRewards(order.Id, order.CampaignId, 5, 0, order.Amount)
	suite.Require().NoError(err)

	var rewards []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Find(&rewards)
	assert.Len(suite.T(), rewards, 0)
}

func (suite *MultiLevelRewardLogicTestSuite) TestCalculateAndDistributeRewards_SuspendedDistributor() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

	distributor := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "suspended",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	order := &model.Order{
		Id:              5,
		CampaignId:      1,
		Phone:           "13800000005",
		DistributorPath: "1",
		Status:          "paid",
		PayStatus:       "paid",
		Amount:          100.0,
		PaidAt:          timePtr(time.Now()),
		CreatedAt:       time.Now(),
	}
	suite.db.Create(order)

	err := logic.CalculateAndDistributeRewards(order.Id, order.CampaignId, 5, 0, order.Amount)
	suite.Require().NoError(err)

	var rewards []model.DistributorReward
	suite.db.Where("order_id = ?", order.Id).Find(&rewards)
	assert.Len(suite.T(), rewards, 0)
}

func (suite *MultiLevelRewardLogicTestSuite) TestParseDistributorPath_JSONFormat() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

	path := "[1,2,3]"
	ids := logic.parseDistributorPath(path)

	assert.Len(suite.T(), ids, 3)
	assert.Equal(suite.T(), int64(1), ids[0])
	assert.Equal(suite.T(), int64(2), ids[1])
	assert.Equal(suite.T(), int64(3), ids[2])
}

func (suite *MultiLevelRewardLogicTestSuite) TestParseDistributorPath_CSVFormat() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

	path := "1,2,3"
	ids := logic.parseDistributorPath(path)

	assert.Len(suite.T(), ids, 3)
	assert.Equal(suite.T(), int64(1), ids[0])
	assert.Equal(suite.T(), int64(2), ids[1])
	assert.Equal(suite.T(), int64(3), ids[2])
}

func (suite *MultiLevelRewardLogicTestSuite) TestParseDistributorPath_Empty() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

	path := ""
	ids := logic.parseDistributorPath(path)

	assert.Len(suite.T(), ids, 0)
}

func (suite *MultiLevelRewardLogicTestSuite) TestUpdateUserBalance_NewBalance() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

	distributor := &model.Distributor{
		UserId:            4,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	order := &model.Order{
		Id:         6,
		CampaignId: 1,
		Phone:      "13800000005",
		Status:     "paid",
		PayStatus:  "paid",
		Amount:     100.0,
		PaidAt:     timePtr(time.Now()),
		CreatedAt:  time.Now(),
	}
	suite.db.Create(order)

	tx := suite.db.Begin()
	err := logic.distributeRewardWithTx(tx, distributor, order, 1, 10.0, 10.0)
	suite.Require().NoError(err)
	suite.Require().NoError(tx.Commit().Error)

	var balance model.UserBalance
	suite.db.Where("user_id = ?", 4).First(&balance)
	assert.Equal(suite.T(), 10.0, balance.Balance)
	assert.Equal(suite.T(), 10.0, balance.TotalReward)
}

func (suite *MultiLevelRewardLogicTestSuite) TestUpdateUserBalance_ExistingBalance() {
	logic := NewMultiLevelRewardLogic(context.Background(), suite.svcCtx)

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
		Id:         7,
		CampaignId: 1,
		Phone:      "13800000005",
		Status:     "paid",
		PayStatus:  "paid",
		Amount:     100.0,
		PaidAt:     timePtr(time.Now()),
		CreatedAt:  time.Now(),
	}
	suite.db.Create(order)

	tx := suite.db.Begin()
	err := logic.distributeRewardWithTx(tx, distributor, order, 1, 10.0, 10.0)
	suite.Require().NoError(err)
	suite.Require().NoError(tx.Commit().Error)

	var balance model.UserBalance
	suite.db.Where("user_id = ?", 1).First(&balance)
	assert.Equal(suite.T(), 10.0, balance.Balance)
	assert.Equal(suite.T(), 10.0, balance.TotalReward)
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestMultiLevelRewardLogicTestSuite(t *testing.T) {
	suite.Run(t, new(MultiLevelRewardLogicTestSuite))
}
