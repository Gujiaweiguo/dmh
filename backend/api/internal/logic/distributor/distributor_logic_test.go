package distributor

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDistributorTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(
		&model.Distributor{},
		&model.DistributorApplication{},
		&model.DistributorLink{},
		&model.DistributorLevelReward{},
		&model.DistributorReward{},
		&model.Order{},
		&model.UserBalance{},
		&model.Campaign{},
		&model.Brand{},
		&model.User{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func createTestDistributor(t *testing.T, db *gorm.DB, userId, brandId int64, level int, status string) *model.Distributor {
	distributor := &model.Distributor{
		UserId:            userId,
		BrandId:           brandId,
		Level:             level,
		Status:            status,
		TotalEarnings:     1000.00,
		SubordinatesCount: 5,
	}
	if err := db.Create(distributor).Error; err != nil {
		t.Fatalf("Failed to create test distributor: %v", err)
	}
	return distributor
}

func createTestBrand(t *testing.T, db *gorm.DB, name string) *model.Brand {
	brand := &model.Brand{
		Name:        name,
		Description: "Test brand",
		Status:      "active",
	}
	if err := db.Create(brand).Error; err != nil {
		t.Fatalf("Failed to create test brand: %v", err)
	}
	return brand
}

func createTestUser(t *testing.T, db *gorm.DB, username string) *model.User {
	phoneNum := 13800000000 + rand.Intn(100000000)
	user := &model.User{
		Username: username,
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17zhU", // bcrypt hash of "password"
		Phone:    fmt.Sprintf("%d", phoneNum),
		Email:    fmt.Sprintf("%s@test.com", username),
		RealName: "Test User",
		Status:   "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestDistributorApplication(t *testing.T, db *gorm.DB, userId, brandId int64) *model.DistributorApplication {
	app := &model.DistributorApplication{
		UserId:  userId,
		BrandId: brandId,
		Status:  "pending",
		Reason:  "I want to become a distributor",
	}
	if err := db.Create(app).Error; err != nil {
		t.Fatalf("Failed to create test distributor application: %v", err)
	}
	return app
}

func TestGetBrandDistributorsLogic_GetDistributors_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	createTestDistributor(t, db, 1, 1, 1, "active")
	createTestDistributor(t, db, 2, 1, 2, "active")
	createTestDistributor(t, db, 3, 1, 1, "suspended")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandDistributorsLogic(ctx, svcCtx)

	req := &types.GetDistributorsReq{
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetBrandDistributors(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Total >= 3)
}

func TestGetBrandDistributorsLogic_WithStatusFilter(t *testing.T) {
	db := setupDistributorTestDB(t)

	createTestDistributor(t, db, 1, 1, 1, "active")
	createTestDistributor(t, db, 2, 1, 2, "active")
	createTestDistributor(t, db, 3, 1, 1, "suspended")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandDistributorsLogic(ctx, svcCtx)

	req := &types.GetDistributorsReq{
		Status:   "active",
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetBrandDistributors(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetBrandDistributorsLogic_Pagination(t *testing.T) {
	db := setupDistributorTestDB(t)

	for i := 0; i < 25; i++ {
		createTestDistributor(t, db, int64(i+1), 1, 1, "active")
	}

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandDistributorsLogic(ctx, svcCtx)

	req := &types.GetDistributorsReq{
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetBrandDistributors(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(25), resp.Total)
	assert.Len(t, resp.Distributors, 10)
}

func TestGenerateDistributorLinkLogic_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "TestBrand")
	createTestDistributor(t, db, user.Id, brand.Id, 1, "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGenerateDistributorLinkLogic(ctx, svcCtx)

	req := &types.GenerateLinkReq{
		CampaignId: brand.Id,
	}

	resp, err := logic.GenerateDistributorLink(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotZero(t, resp.LinkId)
	assert.NotEmpty(t, resp.Link)
	assert.NotEmpty(t, resp.LinkCode)
	assert.Equal(t, brand.Id, resp.CampaignId)
}

func TestGetMyDistributorDashboardLogic_HasDistributor(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand1 := createTestBrand(t, db, "Brand1")
	brand2 := createTestBrand(t, db, "Brand2")
	createTestDistributor(t, db, user.Id, brand1.Id, 1, "active")
	createTestDistributor(t, db, user.Id, brand2.Id, 2, "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMyDistributorDashboardLogic(ctx, svcCtx)

	resp, err := logic.GetMyDistributorDashboard()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.Code)
	assert.True(t, resp.Data.HasDistributor)
	assert.Len(t, resp.Data.Brands, 2)
}

func TestGetMyDistributorDashboardLogic_NoDistributor(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMyDistributorDashboardLogic(ctx, svcCtx)

	resp, err := logic.GetMyDistributorDashboard()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.Code)
	assert.False(t, resp.Data.HasDistributor)
	assert.Len(t, resp.Data.Brands, 0)
}

func TestApproveDistributorApplicationLogic_Approve(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "applicant")
	brand := createTestBrand(t, db, "TestBrand")
	reviewer := createTestUser(t, db, "reviewer")
	app := createTestDistributorApplication(t, db, user.Id, brand.Id)

	ctx := context.WithValue(context.Background(), "applicationId", app.Id)
	ctx = context.WithValue(ctx, "userId", reviewer.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewApproveDistributorApplicationLogic(ctx, svcCtx)

	req := &types.ApproveDistributorReq{
		Action: "approved",
		Level:  1,
		Reason: "Approved",
	}

	resp, err := logic.ApproveDistributorApplication(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "approved", resp.Status)

	var distributor model.Distributor
	err = db.Where("user_id = ? AND brand_id = ?", user.Id, brand.Id).First(&distributor).Error
	assert.NoError(t, err)
	assert.Equal(t, "active", distributor.Status)
	assert.Equal(t, 1, distributor.Level)
	assert.Equal(t, reviewer.Id, *distributor.ApprovedBy)
}

func TestApproveDistributorApplicationLogic_Reject(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "applicant")
	brand := createTestBrand(t, db, "TestBrand")
	reviewer := createTestUser(t, db, "reviewer")
	app := createTestDistributorApplication(t, db, user.Id, brand.Id)

	ctx := context.WithValue(context.Background(), "applicationId", app.Id)
	ctx = context.WithValue(ctx, "userId", reviewer.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewApproveDistributorApplicationLogic(ctx, svcCtx)

	req := &types.ApproveDistributorReq{
		Action: "rejected",
		Reason: "Not qualified",
	}

	resp, err := logic.ApproveDistributorApplication(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "rejected", resp.Status)
}

func TestUpdateDistributorLevelLogic_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "TestBrand")
	dist := createTestDistributor(t, db, user.Id, brand.Id, 1, "active")

	ctx := context.WithValue(context.Background(), "distributorId", dist.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateDistributorLevelLogic(ctx, svcCtx)

	req := &types.UpdateDistributorLevelReq{
		Level: 3,
	}

	resp, err := logic.UpdateDistributorLevel(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Distributor level updated successfully", resp.Message)

	var updatedDist model.Distributor
	err = db.First(&updatedDist, dist.Id).Error
	assert.NoError(t, err)
	assert.Equal(t, 3, updatedDist.Level)
}

func TestUpdateDistributorStatusLogic_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "TestBrand")
	dist := createTestDistributor(t, db, user.Id, brand.Id, 1, "active")

	ctx := context.WithValue(context.Background(), "distributorId", dist.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateDistributorStatusLogic(ctx, svcCtx)

	req := &types.UpdateDistributorStatusReq{
		Status: "suspended",
		Reason: "Violation",
	}

	resp, err := logic.UpdateDistributorStatus(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Distributor status updated successfully", resp.Message)

	var updatedDist model.Distributor
	err = db.First(&updatedDist, dist.Id).Error
	assert.NoError(t, err)
	assert.Equal(t, "suspended", updatedDist.Status)
}

func TestGetDistributorLinksLogic_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "TestBrand")
	createTestDistributor(t, db, user.Id, brand.Id, 1, "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGenerateDistributorLinkLogic(ctx, svcCtx)

	req := &types.GenerateLinkReq{
		CampaignId: brand.Id,
	}

	link1, _ := logic.GenerateDistributorLink(req)
	link2, _ := logic.GenerateDistributorLink(req)

	getLinksLogic := NewGetDistributorLinksLogic(ctx, svcCtx)
	links, err := getLinksLogic.GetDistributorLinks()

	assert.NoError(t, err)
	assert.Len(t, links, 2)
	assert.Contains(t, []int64{link1.LinkId, link2.LinkId}, links[0].LinkId)
}

func TestGetBrandDistributorApplicationsLogic_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	user1 := createTestUser(t, db, "applicant1")
	user2 := createTestUser(t, db, "applicant2")
	brand := createTestBrand(t, db, "TestBrand")
	createTestDistributorApplication(t, db, user1.Id, brand.Id)
	createTestDistributorApplication(t, db, user2.Id, brand.Id)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandDistributorApplicationsLogic(ctx, svcCtx)

	req := &types.GetDistributorApplicationsReq{
		Page:     1,
		PageSize: 10,
		BrandId:  brand.Id,
	}

	resp, err := logic.GetBrandDistributorApplications(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.GreaterOrEqual(t, resp.Total, int64(2))
	assert.Len(t, resp.Applications, 2)
}

func TestGetBrandDistributorApplicationsLogic_WithStatusFilter(t *testing.T) {
	db := setupDistributorTestDB(t)

	user1 := createTestUser(t, db, "applicant1")
	user2 := createTestUser(t, db, "applicant2")
	brand := createTestBrand(t, db, "TestBrand")
	app1 := createTestDistributorApplication(t, db, user1.Id, brand.Id)
	app2 := createTestDistributorApplication(t, db, user2.Id, brand.Id)

	var pendingApp model.DistributorApplication
	db.First(&pendingApp, app1.Id)
	pendingApp.Status = "pending"
	db.Save(&pendingApp)

	var rejectedApp model.DistributorApplication
	db.First(&rejectedApp, app2.Id)
	rejectedApp.Status = "rejected"
	db.Save(&rejectedApp)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandDistributorApplicationsLogic(ctx, svcCtx)

	req := &types.GetDistributorApplicationsReq{
		Page:     1,
		PageSize: 10,
		Status:   "pending",
	}

	resp, err := logic.GetBrandDistributorApplications(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Applications, 1)
}

func TestGenerateDistributorLinkLogic_DistributorNotFound(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "TestBrand")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGenerateDistributorLinkLogic(ctx, svcCtx)

	req := &types.GenerateLinkReq{
		CampaignId: brand.Id,
	}

	resp, err := logic.GenerateDistributorLink(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateDistributorLevelLogic_DistributorNotFound(t *testing.T) {
	db := setupDistributorTestDB(t)

	ctx := context.WithValue(context.Background(), "distributorId", int64(999))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateDistributorLevelLogic(ctx, svcCtx)

	req := &types.UpdateDistributorLevelReq{
		Level: 3,
	}

	resp, err := logic.UpdateDistributorLevel(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateDistributorStatusLogic_DistributorNotFound(t *testing.T) {
	db := setupDistributorTestDB(t)

	ctx := context.WithValue(context.Background(), "distributorId", int64(999))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateDistributorStatusLogic(ctx, svcCtx)

	req := &types.UpdateDistributorStatusReq{
		Status: "suspended",
	}

	resp, err := logic.UpdateDistributorStatus(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetMyDistributorDashboardLogic_MultipleBrands(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand1 := createTestBrand(t, db, "Brand1")
	brand2 := createTestBrand(t, db, "Brand2")
	brand3 := createTestBrand(t, db, "Brand3")
	createTestDistributor(t, db, user.Id, brand1.Id, 1, "active")
	createTestDistributor(t, db, user.Id, brand2.Id, 2, "active")
	createTestDistributor(t, db, user.Id, brand3.Id, 3, "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMyDistributorDashboardLogic(ctx, svcCtx)

	resp, err := logic.GetMyDistributorDashboard()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Data.HasDistributor)
	assert.Len(t, resp.Data.Brands, 3)
}

func TestGetMyDistributorDashboardLogic_InactiveDistributors(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand1 := createTestBrand(t, db, "TestBrand1")
	brand2 := createTestBrand(t, db, "TestBrand2")
	createTestDistributor(t, db, user.Id, brand1.Id, 1, "pending")
	createTestDistributor(t, db, user.Id, brand2.Id, 1, "suspended")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMyDistributorDashboardLogic(ctx, svcCtx)

	resp, err := logic.GetMyDistributorDashboard()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Data.HasDistributor)
}

func TestSetDistributorLevelRewardsLogic_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	brand := createTestBrand(t, db, "TestBrand")

	ctx := context.WithValue(context.Background(), "brandId", brand.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewSetDistributorLevelRewardsLogic(ctx, svcCtx)

	req := &types.SetDistributorLevelRewardsReq{
		Rewards: []types.SetDistributorLevelRewardReq{
			{Level: 1, RewardPercentage: 5.0},
			{Level: 2, RewardPercentage: 8.0},
			{Level: 3, RewardPercentage: 10.0},
		},
	}

	resp, err := logic.SetDistributorLevelRewards(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Success", resp.Message)

	var rewards []model.DistributorLevelReward
	db.Where("brand_id = ?", brand.Id).Find(&rewards)
	assert.Len(t, rewards, 3)
}

func TestGetDistributorLevelRewardsLogic_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	brand := createTestBrand(t, db, "TestBrand")

	var levelReward1 model.DistributorLevelReward
	var levelReward2 model.DistributorLevelReward
	var levelReward3 model.DistributorLevelReward
	levelReward1.BrandId = brand.Id
	levelReward1.Level = 1
	levelReward1.RewardPercentage = 5.0
	levelReward2.BrandId = brand.Id
	levelReward2.Level = 2
	levelReward2.RewardPercentage = 8.0
	levelReward3.BrandId = brand.Id
	levelReward3.Level = 3
	levelReward3.RewardPercentage = 10.0
	db.Create(&levelReward1)
	db.Create(&levelReward2)
	db.Create(&levelReward3)

	ctx := context.WithValue(context.Background(), "brandId", brand.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDistributorLevelRewardsLogic(ctx, svcCtx)

	resp, err := logic.GetDistributorLevelRewards()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Rewards, 3)
}

func TestGetDistributorStatisticsLogic_Success(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "TestBrand")
	createTestDistributor(t, db, user.Id, brand.Id, 1, "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDistributorStatisticsLogic(ctx, svcCtx)

	req := &types.GetDistributorStatisticsReq{
		BrandId: brand.Id,
	}

	resp, err := logic.GetDistributorStatistics(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.Code)
}

func TestGetDistributorStatisticsLogic_DistributorNotFound(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")
	brand := createTestBrand(t, db, "TestBrand")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDistributorStatisticsLogic(ctx, svcCtx)

	req := &types.GetDistributorStatisticsReq{
		BrandId: brand.Id,
	}

	resp, err := logic.GetDistributorStatistics(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.Code)
}

func TestGetBrandDistributorsLogic_EmptyResult(t *testing.T) {
	db := setupDistributorTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandDistributorsLogic(ctx, svcCtx)

	req := &types.GetDistributorsReq{
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetBrandDistributors(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Distributors, 0)
}

func TestGetBrandDistributorApplicationsLogic_Pagination(t *testing.T) {
	db := setupDistributorTestDB(t)

	for i := 0; i < 25; i++ {
		user := createTestUser(t, db, fmt.Sprintf("applicant%d", i))
		brand := createTestBrand(t, db, fmt.Sprintf("Brand%d", i))
		createTestDistributorApplication(t, db, user.Id, brand.Id)
	}

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandDistributorApplicationsLogic(ctx, svcCtx)

	req := &types.GetDistributorApplicationsReq{
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetBrandDistributorApplications(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(25), resp.Total)
	assert.Len(t, resp.Applications, 10)
}

func TestGetDistributorLinksLogic_NoDistributor(t *testing.T) {
	db := setupDistributorTestDB(t)

	user := createTestUser(t, db, "testuser")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDistributorLinksLogic(ctx, svcCtx)

	links, err := logic.GetDistributorLinks()

	assert.Error(t, err)
	assert.Nil(t, links)
}

func TestApproveDistributorApplicationLogic_ApplicationNotFound(t *testing.T) {
	db := setupDistributorTestDB(t)

	reviewer := createTestUser(t, db, "reviewer")

	ctx := context.WithValue(context.Background(), "applicationId", int64(999))
	ctx = context.WithValue(ctx, "userId", reviewer.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewApproveDistributorApplicationLogic(ctx, svcCtx)

	req := &types.ApproveDistributorReq{
		Action: "approved",
		Level:  1,
		Reason: "Approved",
	}

	resp, err := logic.ApproveDistributorApplication(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetBrandDistributorApplicationsLogic_Empty(t *testing.T) {
	db := setupDistributorTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetBrandDistributorApplicationsLogic(ctx, svcCtx)

	req := &types.GetDistributorApplicationsReq{
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetBrandDistributorApplications(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Applications, 0)
}

func TestGetDistributorLevelRewardsLogic_Empty(t *testing.T) {
	db := setupDistributorTestDB(t)

	brand := createTestBrand(t, db, "TestBrand")

	ctx := context.WithValue(context.Background(), "brandId", brand.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDistributorLevelRewardsLogic(ctx, svcCtx)

	resp, err := logic.GetDistributorLevelRewards()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Rewards, 0)
}

func TestSetDistributorLevelRewardsLogic_UpdateExisting(t *testing.T) {
	db := setupDistributorTestDB(t)

	brand := createTestBrand(t, db, "TestBrand")

	var existingReward model.DistributorLevelReward
	existingReward.BrandId = brand.Id
	existingReward.Level = 1
	existingReward.RewardPercentage = 5.0
	db.Create(&existingReward)

	ctx := context.WithValue(context.Background(), "brandId", brand.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewSetDistributorLevelRewardsLogic(ctx, svcCtx)

	req := &types.SetDistributorLevelRewardsReq{
		Rewards: []types.SetDistributorLevelRewardReq{
			{Level: 1, RewardPercentage: 8.0},
		},
	}

	resp, err := logic.SetDistributorLevelRewards(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Success", resp.Message)

	var updatedReward model.DistributorLevelReward
	db.Where("brand_id = ? AND level = ?", brand.Id, 1).First(&updatedReward)
	assert.Equal(t, 8.0, updatedReward.RewardPercentage)
}
