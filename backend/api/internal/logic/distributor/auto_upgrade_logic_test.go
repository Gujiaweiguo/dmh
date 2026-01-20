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

type AutoUpgradeLogicTestSuite struct {
	suite.Suite
	db     *gorm.DB
	svcCtx *svc.ServiceContext
}

func (suite *AutoUpgradeLogicTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.UserBrand{},
		&model.Brand{},
		&model.Campaign{},
		&model.Distributor{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.svcCtx = &svc.ServiceContext{DB: db}
	suite.createTestData()
}

func (suite *AutoUpgradeLogicTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *AutoUpgradeLogicTestSuite) createTestData() {
	roles := []model.Role{
		{ID: 1, Name: "平台管理员", Code: "platform_admin"},
		{ID: 2, Name: "品牌管理员", Code: "brand_admin"},
		{ID: 3, Name: "参与者", Code: "participant"},
		{ID: 4, Name: "分销商", Code: "distributor"},
	}
	for _, role := range roles {
		suite.db.Create(&role)
	}

	users := []model.User{
		{Id: 1, Username: "user1", Phone: "13800000001", Email: "user1@test.com", RealName: "用户1", Role: "participant", Status: "active"},
		{Id: 2, Username: "user2", Phone: "13800000002", Email: "user2@test.com", RealName: "用户2", Role: "participant", Status: "active"},
		{Id: 3, Username: "user3", Phone: "13800000003", Email: "user3@test.com", RealName: "用户3", Role: "participant", Status: "active"},
		{Id: 4, Username: "user4", Phone: "13800000004", Email: "user4@test.com", RealName: "用户4", Role: "participant", Status: "active"},
		{Id: 5, Username: "user5", Phone: "13800000005", Email: "user5@test.com", RealName: "用户5", Role: "participant", Status: "active"},
	}
	for _, user := range users {
		suite.db.Create(&user)
	}

	userRoles := []model.UserRole{
		{UserID: 1, RoleID: 3},
		{UserID: 2, RoleID: 3},
		{UserID: 3, RoleID: 3},
		{UserID: 4, RoleID: 3},
		{UserID: 5, RoleID: 3},
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
		{UserId: 5, BrandId: 1},
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
}

func (suite *AutoUpgradeLogicTestSuite) TestCheckAndAutoUpgrade_CreateNewDistributor() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ?", 1)

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

	distributor, err := logic.CheckAndAutoUpgrade(1, 1, 0)
	suite.Require().NoError(err)
	suite.Require().NotNil(distributor)

	assert.Equal(suite.T(), int64(1), distributor.UserId)
	assert.Equal(suite.T(), int64(1), distributor.BrandId)
	assert.Equal(suite.T(), 1, distributor.Level)
	assert.Equal(suite.T(), "active", distributor.Status)
	assert.Equal(suite.T(), 0.0, distributor.TotalEarnings)

	var savedDistributor model.Distributor
	err = suite.db.Where("user_id = ? AND brand_id = ?", 1, 1).First(&savedDistributor).Error
	suite.Require().NoError(err)
	assert.Equal(suite.T(), distributor.Id, savedDistributor.Id)
}

func (suite *AutoUpgradeLogicTestSuite) TestCheckAndAutoUpgrade_AlreadyDistributor() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ?", 1)

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

	distributor1, err := logic.CheckAndAutoUpgrade(1, 1, 0)
	suite.Require().NoError(err)
	suite.Require().NotNil(distributor1)

	distributor2, err := logic.CheckAndAutoUpgrade(1, 1, 0)
	suite.Require().NoError(err)
	suite.Require().NotNil(distributor2)

	assert.Equal(suite.T(), distributor1.Id, distributor2.Id)
	assert.Equal(suite.T(), distributor1.UserId, distributor2.UserId)
	assert.Equal(suite.T(), distributor1.BrandId, distributor2.BrandId)

	var count int64
	suite.db.Model(&model.Distributor{}).Where("user_id = ? AND brand_id = ?", 1, 1).Count(&count)
	assert.Equal(suite.T(), int64(1), count)
}

func (suite *AutoUpgradeLogicTestSuite) TestCheckAndAutoUpgrade_WithReferrer() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id IN (?, ?)", 1, 2)

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

	existingDistributor := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(existingDistributor)

	newDistributor, err := logic.CheckAndAutoUpgrade(2, 1, 1)
	suite.Require().NoError(err)
	suite.Require().NotNil(newDistributor)

	assert.Equal(suite.T(), int64(2), newDistributor.UserId)
	assert.Equal(suite.T(), int64(1), newDistributor.BrandId)
	assert.Equal(suite.T(), existingDistributor.Id, *newDistributor.ParentId)
}

func (suite *AutoUpgradeLogicTestSuite) TestCheckAndAutoUpgradeWithCampaign_Success() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ?", 1)

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

	distributor, err := logic.CheckAndAutoUpgradeWithCampaign(1, 1, 1, 0)
	suite.Require().NoError(err)
	suite.Require().NotNil(distributor)

	assert.Equal(suite.T(), int64(1), distributor.UserId)
	assert.Equal(suite.T(), int64(1), distributor.BrandId)
	assert.Equal(suite.T(), 1, distributor.Level)
	assert.Equal(suite.T(), "active", distributor.Status)
}

func (suite *AutoUpgradeLogicTestSuite) TestCheckAndAutoUpgradeWithCampaign_DisabledDistribution() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ?", 3)

	distributionRewardsDisabled := ""

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

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

	distributor, err := logic.CheckAndAutoUpgradeWithCampaign(3, 1, 2, 0)
	suite.Require().Error(err)
	suite.Require().Nil(distributor)
	assert.Contains(suite.T(), err.Error(), "活动未启用分销")
}

func (suite *AutoUpgradeLogicTestSuite) TestCalculateDistributorPath_SingleLevel() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id = ?", 1)

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

	distributor := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	path := logic.CalculateDistributorPath(distributor.Id, 3)
	assert.Equal(suite.T(), "[1]", path)
}

func (suite *AutoUpgradeLogicTestSuite) TestCalculateDistributorPath_ThreeLevel() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id IN (?, ?, ?)", 1, 2, 3)

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

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

	path := logic.CalculateDistributorPath(d3.Id, 3)
	assert.Equal(suite.T(), "[1 2 3]", path)
}

func (suite *AutoUpgradeLogicTestSuite) TestAddDistributorRole() {
	suite.db.Exec("DELETE FROM user_roles WHERE user_id = ?", 4)

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

	err := logic.addDistributorRole(4)
	suite.Require().NoError(err)

	var userRole model.UserRole
	err = suite.db.Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.code = ?", 4, "distributor").
		First(&userRole).Error
	suite.Require().NoError(err)

	var role model.Role
	suite.db.Where("id = ?", userRole.RoleID).First(&role)
	assert.Equal(suite.T(), "distributor", role.Code)
}

func (suite *AutoUpgradeLogicTestSuite) TestUpdateSubordinatesCount() {
	suite.db.Exec("DELETE FROM distributors WHERE user_id IN (?, ?, ?)", 1, 2, 3)

	logic := NewAutoUpgradeLogic(context.Background(), suite.svcCtx)

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
		Level:             2,
		Status:            "active",
		ParentId:          &d1.Id,
		TotalEarnings:     0.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(d3)

	err := logic.updateSubordinatesCount(d1.Id)
	suite.Require().NoError(err)

	var updatedD1 model.Distributor
	suite.db.Where("id = ?", d1.Id).First(&updatedD1)
	assert.Equal(suite.T(), 2, updatedD1.SubordinatesCount)
}

func TestAutoUpgradeLogicTestSuite(t *testing.T) {
	suite.Run(t, new(AutoUpgradeLogicTestSuite))
}
