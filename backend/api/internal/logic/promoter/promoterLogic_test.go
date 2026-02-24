package promoter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PromoterLogicTestSuite struct {
	suite.Suite
	db     *gorm.DB
	svcCtx *svc.ServiceContext
}

func (suite *PromoterLogicTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open("root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local"), &gorm.Config{})
	suite.Require().NoError(err)

	suite.db = db
	suite.svcCtx = &svc.ServiceContext{DB: db}
}

func (suite *PromoterLogicTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *PromoterLogicTestSuite) SetupTest() {
	suite.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	suite.db.Exec("TRUNCATE TABLE promoter_rewards")
	suite.db.Exec("TRUNCATE TABLE promoter_links")
	suite.db.Exec("TRUNCATE TABLE promoters")
	suite.db.Exec("TRUNCATE TABLE campaigns")
	suite.db.Exec("TRUNCATE TABLE brands")
	suite.db.Exec("TRUNCATE TABLE users")
	suite.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

func (suite *PromoterLogicTestSuite) createTestUser(id int64, username string) *model.User {
	user := &model.User{
		Id:       id,
		Username: username,
		Password: "hashedpassword",
		Phone:    fmt.Sprintf("138%08d", id),
		Status:   "active",
	}
	suite.db.Create(user)
	return user
}

func (suite *PromoterLogicTestSuite) createTestBrand(id int64, name string) *model.Brand {
	brand := &model.Brand{
		Id:     id,
		Name:   name,
		Status: "active",
	}
	suite.db.Create(brand)
	return brand
}

func (suite *PromoterLogicTestSuite) createTestCampaign(id, brandId int64, name string) *model.Campaign {
	campaign := &model.Campaign{
		Id:        id,
		BrandId:   brandId,
		Name:      name,
		Status:    "active",
		StartTime: time.Now().Add(-time.Hour),
		EndTime:   time.Now().Add(time.Hour),
	}
	suite.db.Create(campaign)
	return campaign
}

func (suite *PromoterLogicTestSuite) createTestPromoter(id, userId, brandId int64, status string) *model.Promoter {
	promoter := &model.Promoter{
		Id:        id,
		UserId:    userId,
		BrandId:   brandId,
		Status:    status,
		Level:     "normal",
		CreatedAt: time.Now(),
	}
	suite.db.Create(promoter)
	return promoter
}

func (suite *PromoterLogicTestSuite) TestGetPromoterList_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestUser(2, "promoter2")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")
	suite.createTestPromoter(2, 2, 1, "active")

	l := NewGetPromoterListLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterListReq{Page: 1, PageSize: 10}

	resp, err := l.GetPromoterList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2), resp.Total)
	assert.Len(suite.T(), resp.Promoters, 2)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterList_FilterByBrand() {
	ctx := context.Background()
	suite.createTestUser(1, "p1")
	suite.createTestUser(2, "p2")
	suite.createTestBrand(1, "Brand1")
	suite.createTestBrand(2, "Brand2")
	suite.createTestPromoter(1, 1, 1, "active")
	suite.createTestPromoter(2, 2, 2, "active")

	l := NewGetPromoterListLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterListReq{Page: 1, PageSize: 10, BrandId: 1}

	resp, err := l.GetPromoterList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.Total)
	assert.Equal(suite.T(), int64(1), resp.Promoters[0].BrandId)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterList_FilterByStatus() {
	ctx := context.Background()
	suite.createTestUser(1, "p1")
	suite.createTestUser(2, "p2")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")
	suite.createTestPromoter(2, 2, 1, "inactive")

	l := NewGetPromoterListLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterListReq{Page: 1, PageSize: 10, Status: "active"}

	resp, err := l.GetPromoterList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.Total)
	assert.Equal(suite.T(), "active", resp.Promoters[0].Status)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterList_Pagination() {
	ctx := context.Background()
	suite.createTestBrand(1, "Brand1")
	for i := 1; i <= 5; i++ {
		suite.createTestUser(int64(i), fmt.Sprintf("user%d", i))
		suite.createTestPromoter(int64(i), int64(i), 1, "active")
	}

	l := NewGetPromoterListLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterListReq{Page: 1, PageSize: 2}

	resp, err := l.GetPromoterList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(5), resp.Total)
	assert.Len(suite.T(), resp.Promoters, 2)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterList_DefaultPagination() {
	ctx := context.Background()
	suite.createTestBrand(1, "Brand1")
	for i := 1; i <= 25; i++ {
		suite.createTestUser(int64(i), fmt.Sprintf("user%d", i))
		suite.createTestPromoter(int64(i), int64(i), 1, "active")
	}

	l := NewGetPromoterListLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterListReq{}

	resp, err := l.GetPromoterList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(25), resp.Total)
	assert.Len(suite.T(), resp.Promoters, 20)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterList_EmptyResult() {
	ctx := context.Background()

	l := NewGetPromoterListLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterListReq{Page: 1, PageSize: 10}

	resp, err := l.GetPromoterList(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(0), resp.Total)
	assert.Len(suite.T(), resp.Promoters, 0)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterDetail_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")

	l := NewGetPromoterDetailLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterDetailReq{Id: 1}

	resp, err := l.GetPromoterDetail(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.Id)
	assert.Equal(suite.T(), "promoter1", resp.Username)
	assert.Equal(suite.T(), "Brand1", resp.BrandName)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterDetail_NotFound() {
	ctx := context.Background()

	l := NewGetPromoterDetailLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterDetailReq{Id: 99999}

	resp, err := l.GetPromoterDetail(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *PromoterLogicTestSuite) TestGetPromoterDetail_WithLinks() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestCampaign(1, 1, "Campaign1")
	suite.createTestPromoter(1, 1, 1, "active")

	suite.db.Create(&model.PromoterLink{
		PromoterId: 1,
		CampaignId: 1,
		LinkCode:   "abc12345",
		ClickCount: 10,
		CreatedAt:  time.Now(),
	})

	l := NewGetPromoterDetailLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterDetailReq{Id: 1}

	resp, err := l.GetPromoterDetail(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Len(suite.T(), resp.Links, 1)
	assert.Equal(suite.T(), "abc12345", resp.Links[0].LinkCode)
	assert.Equal(suite.T(), int64(10), resp.Links[0].ClickCount)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterDetailById_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")

	l := NewGetPromoterDetailLogic(ctx, suite.svcCtx)

	resp, err := l.GetPromoterDetailById(1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.Id)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterRewards_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")

	suite.db.Create(&model.PromoterReward{
		PromoterId:  1,
		Type:        "commission",
		Status:      "pending",
		Amount:      100.0,
		Description: "Test reward",
		CreatedAt:   time.Now(),
	})

	l := NewGetPromoterRewardsLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterRewardsReq{PromoterId: 1, Page: 1, PageSize: 10}

	resp, err := l.GetPromoterRewards(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.Total)
	assert.Len(suite.T(), resp.Rewards, 1)
	assert.Equal(suite.T(), "commission", resp.Rewards[0].Type)
	assert.Equal(suite.T(), 100.0, resp.Rewards[0].Amount)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterRewards_FilterByType() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")

	suite.db.Create(&model.PromoterReward{
		PromoterId: 1, Type: "commission", Status: "pending", Amount: 100.0, CreatedAt: time.Now(),
	})
	suite.db.Create(&model.PromoterReward{
		PromoterId: 1, Type: "bonus", Status: "pending", Amount: 50.0, CreatedAt: time.Now(),
	})

	l := NewGetPromoterRewardsLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterRewardsReq{PromoterId: 1, Type: "commission", Page: 1, PageSize: 10}

	resp, err := l.GetPromoterRewards(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.Total)
	assert.Equal(suite.T(), "commission", resp.Rewards[0].Type)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterRewards_FilterByStatus() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")

	suite.db.Create(&model.PromoterReward{
		PromoterId: 1, Type: "commission", Status: "pending", Amount: 100.0, CreatedAt: time.Now(),
	})
	suite.db.Create(&model.PromoterReward{
		PromoterId: 1, Type: "commission", Status: "paid", Amount: 50.0, CreatedAt: time.Now(),
	})

	l := NewGetPromoterRewardsLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterRewardsReq{PromoterId: 1, Status: "paid", Page: 1, PageSize: 10}

	resp, err := l.GetPromoterRewards(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(1), resp.Total)
	assert.Equal(suite.T(), "paid", resp.Rewards[0].Status)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterRewards_Pagination() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")

	for i := 1; i <= 5; i++ {
		suite.db.Create(&model.PromoterReward{
			PromoterId: 1, Type: "commission", Status: "pending", Amount: float64(i * 10), CreatedAt: time.Now(),
		})
	}

	l := NewGetPromoterRewardsLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterRewardsReq{PromoterId: 1, Page: 1, PageSize: 2}

	resp, err := l.GetPromoterRewards(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(5), resp.Total)
	assert.Len(suite.T(), resp.Rewards, 2)
}

func (suite *PromoterLogicTestSuite) TestGetPromoterRewards_EmptyResult() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")

	l := NewGetPromoterRewardsLogic(ctx, suite.svcCtx)
	req := &types.GetPromoterRewardsReq{PromoterId: 1, Page: 1, PageSize: 10}

	resp, err := l.GetPromoterRewards(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(0), resp.Total)
	assert.Len(suite.T(), resp.Rewards, 0)
}

func (suite *PromoterLogicTestSuite) TestGeneratePromoterLink_Success() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")
	suite.createTestCampaign(100, 1, "Test Campaign")

	l := NewGeneratePromoterLinkLogic(ctx, suite.svcCtx)
	req := &types.GeneratePromoterLinkReq{
		PromoterId: 1,
		CampaignId: 100,
	}

	resp, err := l.GeneratePromoterLink(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Len(suite.T(), resp.LinkCode, 8)
	assert.Contains(suite.T(), resp.LinkUrl, resp.LinkCode)

	var link model.PromoterLink
	suite.db.Where("promoter_id = ? AND campaign_id = ?", 1, 100).First(&link)
	assert.Equal(suite.T(), resp.LinkCode, link.LinkCode)
}

func (suite *PromoterLogicTestSuite) TestGeneratePromoterLink_PromoterNotFound() {
	ctx := context.Background()
	suite.createTestBrand(1, "Brand1")
	suite.createTestCampaign(100, 1, "Test Campaign")

	l := NewGeneratePromoterLinkLogic(ctx, suite.svcCtx)
	req := &types.GeneratePromoterLinkReq{
		PromoterId: 999,
		CampaignId: 100,
	}

	resp, err := l.GeneratePromoterLink(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "promoter not found")
}

func (suite *PromoterLogicTestSuite) TestGeneratePromoterLink_CampaignNotFound() {
	ctx := context.Background()
	suite.createTestUser(1, "promoter1")
	suite.createTestBrand(1, "Brand1")
	suite.createTestPromoter(1, 1, 1, "active")

	l := NewGeneratePromoterLinkLogic(ctx, suite.svcCtx)
	req := &types.GeneratePromoterLinkReq{
		PromoterId: 1,
		CampaignId: 999,
	}

	resp, err := l.GeneratePromoterLink(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "campaign not found")
}

func TestPromoterLogicTestSuite(t *testing.T) {
	suite.Run(t, new(PromoterLogicTestSuite))
}
