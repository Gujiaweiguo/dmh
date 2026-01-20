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

type WithdrawalLogicTestSuite struct {
	suite.Suite
	db     *gorm.DB
	svcCtx *svc.ServiceContext
}

func (suite *WithdrawalLogicTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.UserBrand{},
		&model.Brand{},
		&model.Distributor{},
		&model.UserBalance{},
		&model.Withdrawal{},
	)
	suite.Require().NoError(err)

	suite.db = db
	suite.svcCtx = &svc.ServiceContext{DB: db}
	suite.createTestData()
}

func (suite *WithdrawalLogicTestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *WithdrawalLogicTestSuite) createTestData() {
	roles := []model.Role{
		{ID: 1, Name: "分销商", Code: "distributor"},
		{ID: 2, Name: "平台管理员", Code: "platform_admin"},
	}
	for _, role := range roles {
		suite.db.Create(&role)
	}

	users := []model.User{
		{Id: 1, Username: "user1", Phone: "13800000001", Email: "user1@test.com", RealName: "用户1", Role: "distributor", Status: "active"},
		{Id: 2, Username: "admin1", Phone: "13800000002", Email: "admin1@test.com", RealName: "管理员", Role: "platform_admin", Status: "active"},
	}
	for _, user := range users {
		suite.db.Create(&user)
	}

	userRoles := []model.UserRole{
		{UserID: 1, RoleID: 1},
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
	}
	for _, ub := range userBrands {
		suite.db.Create(&ub)
	}

	distributor := &model.Distributor{
		UserId:            1,
		BrandId:           1,
		Level:             1,
		Status:            "active",
		TotalEarnings:     1000.0,
		SubordinatesCount: 0,
	}
	suite.db.Create(distributor)

	balance := &model.UserBalance{
		UserId:      1,
		Balance:     1000.0,
		TotalReward: 1000.0,
		Version:     0,
	}
	suite.db.Create(balance)
}

func (suite *WithdrawalLogicTestSuite) TestApplyWithdrawal_Success() {
	ctx := context.WithValue(context.Background(), "userId", int64(1))
	logic := NewWithdrawalLogic(ctx, suite.svcCtx)

	req := &ApplyWithdrawalReq{
		Amount:      100.0,
		PayType:     "wechat",
		PayAccount:  "13800000001",
		PayRealName: "张三",
	}

	withdrawal, err := logic.ApplyWithdrawal(req)
	suite.Require().NoError(err)
	suite.Require().NotNil(withdrawal)

	assert.Equal(suite.T(), int64(1), withdrawal.UserID)
	assert.Equal(suite.T(), 100.0, withdrawal.Amount)
	assert.Equal(suite.T(), "pending", withdrawal.Status)
	assert.Equal(suite.T(), "wechat", withdrawal.PayType)

	var balance model.UserBalance
	suite.db.Where("user_id = ?", 1).First(&balance)
	assert.Equal(suite.T(), 900.0, balance.Balance)
}

func (suite *WithdrawalLogicTestSuite) TestApplyWithdrawal_AmountExceedsBalance() {
	ctx := context.WithValue(context.Background(), "userId", int64(1))
	logic := NewWithdrawalLogic(ctx, suite.svcCtx)

	req := &ApplyWithdrawalReq{
		Amount:      2000.0,
		PayType:     "wechat",
		PayAccount:  "13800000001",
		PayRealName: "张三",
	}

	withdrawal, err := logic.ApplyWithdrawal(req)
	suite.Require().Error(err)
	suite.Require().Nil(withdrawal)
	assert.Contains(suite.T(), err.Error(), "提现金额超过可用余额")
}

func (suite *WithdrawalLogicTestSuite) TestApplyWithdrawal_InvalidAmount() {
	ctx := context.WithValue(context.Background(), "userId", int64(1))
	logic := NewWithdrawalLogic(ctx, suite.svcCtx)

	req := &ApplyWithdrawalReq{
		Amount:      -100.0,
		PayType:     "wechat",
		PayAccount:  "13800000001",
		PayRealName: "张三",
	}

	withdrawal, err := logic.ApplyWithdrawal(req)
	suite.Require().Error(err)
	suite.Require().Nil(withdrawal)
	assert.Contains(suite.T(), err.Error(), "提现金额必须大于0")
}

func (suite *WithdrawalLogicTestSuite) TestApproveWithdrawal_Success() {
	suite.db.Exec("UPDATE user_balances SET balance = ? WHERE user_id = ?", 900.0, 1)

	withdrawal := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        100.0,
		Status:        "pending",
		PayType:       "wechat",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
	}
	suite.db.Create(withdrawal)

	ctx := context.WithValue(context.Background(), "userId", int64(2))
	logic := NewWithdrawalLogic(ctx, suite.svcCtx)

	err := logic.ApproveWithdrawal(withdrawal.ID, "审批通过")
	suite.Require().NoError(err)

	var updatedWithdrawal model.Withdrawal
	suite.db.Where("id = ?", withdrawal.ID).First(&updatedWithdrawal)
	assert.Equal(suite.T(), "completed", updatedWithdrawal.Status)
	assert.NotNil(suite.T(), updatedWithdrawal.ApprovedAt)
	assert.Equal(suite.T(), int64(2), *updatedWithdrawal.ApprovedBy)
}

func (suite *WithdrawalLogicTestSuite) TestApproveWithdrawal_AlreadyProcessed() {
	suite.db.Exec("UPDATE user_balances SET balance = ? WHERE user_id = ?", 900.0, 1)

	withdrawal := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        100.0,
		Status:        "approved",
		PayType:       "wechat",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
	}
	suite.db.Create(withdrawal)

	ctx := context.WithValue(context.Background(), "userId", int64(2))
	logic := NewWithdrawalLogic(ctx, suite.svcCtx)

	err := logic.ApproveWithdrawal(withdrawal.ID, "审批通过")
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "提现申请已处理")
}

func (suite *WithdrawalLogicTestSuite) TestRejectWithdrawal_Success() {
	suite.db.Exec("UPDATE user_balances SET balance = ? WHERE user_id = ?", 900.0, 1)

	withdrawal := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        100.0,
		Status:        "pending",
		PayType:       "wechat",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
	}
	suite.db.Create(withdrawal)

	ctx := context.WithValue(context.Background(), "userId", int64(2))
	logic := NewWithdrawalLogic(ctx, suite.svcCtx)

	err := logic.RejectWithdrawal(withdrawal.ID, "账户信息错误")
	suite.Require().NoError(err)

	var updatedWithdrawal model.Withdrawal
	suite.db.Where("id = ?", withdrawal.ID).First(&updatedWithdrawal)
	assert.Equal(suite.T(), "rejected", updatedWithdrawal.Status)
	assert.Equal(suite.T(), "账户信息错误", updatedWithdrawal.RejectedReason)

	var balance model.UserBalance
	suite.db.Where("user_id = ?", 1).First(&balance)
	assert.Equal(suite.T(), 1000.0, balance.Balance)
}

func (suite *WithdrawalLogicTestSuite) TestRejectWithdrawal_AlreadyProcessed() {
	suite.db.Exec("UPDATE user_balances SET balance = ? WHERE user_id = ?", 900.0, 1)

	withdrawal := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        100.0,
		Status:        "rejected",
		PayType:       "wechat",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
	}
	suite.db.Create(withdrawal)

	ctx := context.WithValue(context.Background(), "userId", int64(2))
	logic := NewWithdrawalLogic(ctx, suite.svcCtx)

	err := logic.RejectWithdrawal(withdrawal.ID, "账户信息错误")
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "提现申请已处理")
}

func (suite *WithdrawalLogicTestSuite) TestGetMyWithdrawals_Success() {
	ctx := context.WithValue(context.Background(), "userId", int64(1))
	logic := NewWithdrawalLogic(ctx, suite.svcCtx)

	withdrawal1 := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        100.0,
		Status:        "completed",
		PayType:       "wechat",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
		CreatedAt:     time.Now(),
	}
	suite.db.Create(withdrawal1)

	withdrawal2 := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        200.0,
		Status:        "pending",
		PayType:       "alipay",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
		CreatedAt:     time.Now(),
	}
	suite.db.Create(withdrawal2)

	resp, err := logic.GetMyWithdrawals(1, 10)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	assert.Equal(suite.T(), int64(2), resp.Total)
	assert.Equal(suite.T(), 2, len(resp.Withdrawals))
}

func (suite *WithdrawalLogicTestSuite) TestGetWithdrawals_Success() {
	logic := NewWithdrawalLogic(context.Background(), suite.svcCtx)

	withdrawal1 := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        100.0,
		Status:        "completed",
		PayType:       "wechat",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
		CreatedAt:     time.Now(),
	}
	suite.db.Create(withdrawal1)

	withdrawal2 := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        200.0,
		Status:        "pending",
		PayType:       "alipay",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
		CreatedAt:     time.Now(),
	}
	suite.db.Create(withdrawal2)

	resp, err := logic.GetWithdrawals(1, "", 1, 10)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	assert.Equal(suite.T(), int64(2), resp.Total)
	assert.Equal(suite.T(), 2, len(resp.Withdrawals))
}

func (suite *WithdrawalLogicTestSuite) TestGetWithdrawals_WithStatusFilter() {
	logic := NewWithdrawalLogic(context.Background(), suite.svcCtx)

	withdrawal1 := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        100.0,
		Status:        "completed",
		PayType:       "wechat",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
		CreatedAt:     time.Now(),
	}
	suite.db.Create(withdrawal1)

	withdrawal2 := &model.Withdrawal{
		UserID:        1,
		BrandId:       1,
		DistributorId: 1,
		Amount:        200.0,
		Status:        "pending",
		PayType:       "alipay",
		PayAccount:    "13800000001",
		PayRealName:   "张三",
		CreatedAt:     time.Now(),
	}
	suite.db.Create(withdrawal2)

	resp, err := logic.GetWithdrawals(1, "pending", 1, 10)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	assert.Equal(suite.T(), int64(1), resp.Total)
	assert.Equal(suite.T(), 1, len(resp.Withdrawals))
	assert.Equal(suite.T(), "pending", resp.Withdrawals[0].Status)
}

func TestWithdrawalLogicTestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawalLogicTestSuite))
}
