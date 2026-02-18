// Unit tests for order logic
package order

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

func setupTestDB(t *testing.T) *gorm.DB {
	db, _ := testutil.SetupMySQLTestDB(t)
	return db
}

func createTestCampaign(t *testing.T, db *gorm.DB) *model.Campaign {
	campaign := &model.Campaign{
		Name:        "测试活动",
		Description: "这是一个测试活动",
		FormFields:  `[{"type":"text","name":"name","label":"姓名","required":true,"placeholder":"请输入姓名"}]`,
		RewardRule:  10.00,
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "active",
		BrandId:     1,
	}

	if err := db.Create(campaign).Error; err != nil {
		t.Fatalf("Failed to create test campaign: %v", err)
	}

	return campaign
}

// CreateOrderLogic tests
func TestCreateOrderLogic_CreateOrder_Success(t *testing.T) {
	db := setupTestDB(t)
	campaign := createTestCampaign(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateOrderLogic(ctx, svcCtx)

	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData: map[string]string{
			"name": "张三",
		},
		ReferrerId: 0,
	}

	resp, err := logic.CreateOrder(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, campaign.Id, resp.CampaignId)
	assert.Equal(t, "13800138000", resp.Phone)
	assert.Equal(t, "pending", resp.Status)
	assert.NotZero(t, resp.Id)
}

func TestCreateOrderLogic_InvalidPhone(t *testing.T) {
	db := setupTestDB(t)
	campaign := createTestCampaign(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateOrderLogic(ctx, svcCtx)

	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "invalid",
		FormData: map[string]string{
			"name": "张三",
		},
	}

	resp, err := logic.CreateOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "手机号格式错误")
}

func TestCreateOrderLogic_DuplicateOrder(t *testing.T) {
	db := setupTestDB(t)
	campaign := createTestCampaign(t, db)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewCreateOrderLogic(ctx, svcCtx)

	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData: map[string]string{
			"name": "张三",
		},
	}

	// Create first order
	resp1, err1 := logic.CreateOrder(req)
	assert.NoError(t, err1)
	assert.NotNil(t, resp1)

	// Try to create second order with same phone and campaign
	resp2, err2 := logic.CreateOrder(req)

	assert.Error(t, err2)
	assert.Nil(t, resp2)
	assert.Contains(t, err2.Error(), "重复")
}

// VerifyOrderLogic tests
func TestVerifyOrderLogic_VerifyOrder_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create a test order
	createLogic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	orderId := int64(1)
	timestamp := time.Now().Unix()
	verificationCode := createLogic.generateVerificationCode(orderId, "13800138000", timestamp)

	order := &model.Order{
		Id:                 orderId,
		CampaignId:         1,
		Phone:              "13800138000",
		FormData:           `{"name":"张三"}`,
		Status:             "paid",
		PayStatus:          "paid",
		VerificationStatus: "unverified",
		VerificationCode:   verificationCode,
	}
	db.Create(order)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewVerifyOrderLogic(ctx, svcCtx)

	req := &types.VerifyOrderReq{
		Code: verificationCode,
	}

	resp, err := logic.VerifyOrder(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, order.Id, resp.OrderId)
	assert.Equal(t, "verified", resp.Status)
}

func TestVerifyOrderLogic_VerifyOrder_AlreadyVerified(t *testing.T) {
	db := setupTestDB(t)

	// Create a verified order
	createLogic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	orderId := int64(1)
	timestamp := time.Now().Unix()
	verificationCode := createLogic.generateVerificationCode(orderId, "13800138000", timestamp)

	order := &model.Order{
		Id:                 orderId,
		CampaignId:         1,
		Phone:              "13800138000",
		FormData:           `{"name":"张三"}`,
		Status:             "paid",
		PayStatus:          "paid",
		VerificationStatus: "verified",
		VerificationCode:   verificationCode,
	}
	db.Create(order)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "userId", int64(100))
	ctx = context.WithValue(ctx, "roles", []string{"brand_admin"})
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewVerifyOrderLogic(ctx, svcCtx)

	req := &types.VerifyOrderReq{
		Code: verificationCode,
	}

	resp, err := logic.VerifyOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "已核销")
}

// Test verification code generation
func TestCreateOrderLogic_GenerateVerificationCode(t *testing.T) {
	logic := &CreateOrderLogic{}

	orderId := int64(1)
	timestamp := int64(1234567890)
	code := logic.generateVerificationCode(orderId, "13800138000", timestamp)

	assert.NotEmpty(t, code)
	assert.Contains(t, code, "1_13800138000_")
	assert.Contains(t, code, "_1234567890_")

	// Parse and verify
	parts := len(code)
	assert.Greater(t, parts, 0)
}

func TestPaymentCallbackLogic_Success(t *testing.T) {
	db := setupTestDB(t)

	campaign := &model.Campaign{
		Name:               "测试活动",
		Description:        "这是一个测试活动",
		FormFields:         `[{"type":"text","name":"name","label":"姓名","required":true}]`,
		RewardRule:         10.00,
		StartTime:          time.Now().Add(-1 * time.Hour),
		EndTime:            time.Now().Add(24 * time.Hour),
		Status:             "active",
		BrandId:            1,
		EnableDistribution: true,
	}
	db.Create(campaign)

	distributor := &model.Distributor{
		UserId:        100,
		BrandId:       1,
		Level:         1,
		Status:        "active",
		TotalEarnings: 0,
	}
	db.Create(distributor)

	levelReward := &model.DistributorLevelReward{
		BrandId:          1,
		Level:            1,
		RewardPercentage: 50.00,
	}
	db.Create(levelReward)

	order := &model.Order{
		Id:                 1,
		CampaignId:         campaign.Id,
		Phone:              "13800138000",
		FormData:           `{"name":"张三"}`,
		ReferrerId:         100,
		Status:             "pending",
		PayStatus:          "unpaid",
		Amount:             100.00,
		VerificationStatus: "unverified",
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewPaymentCallbackLogic(ctx, svcCtx)

	req := &types.PaymentCallbackReq{
		OrderId: 1,
		TradeNo: "TRADE_123456789",
		Amount:  100.00,
	}

	err := logic.PaymentCallback(req)

	assert.NoError(t, err)

	var updatedOrder model.Order
	db.First(&updatedOrder, order.Id)
	assert.Equal(t, "paid", updatedOrder.PayStatus)
	assert.Equal(t, "TRADE_123456789", updatedOrder.TradeNo)
	assert.Equal(t, 100.00, updatedOrder.Amount)

	var reward model.DistributorReward
	result := db.Where("order_id = ?", order.Id).First(&reward)
	assert.NoError(t, result.Error)
	assert.Equal(t, distributor.Id, reward.DistributorId)
	assert.Equal(t, 500.00, reward.Amount)
	assert.Equal(t, "settled", reward.Status)

	var updatedDistributor model.Distributor
	db.First(&updatedDistributor, distributor.Id)
	assert.Equal(t, 500.00, updatedDistributor.TotalEarnings)
}

func TestPaymentCallbackLogic_OrderNotFound(t *testing.T) {
	db := setupTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewPaymentCallbackLogic(ctx, svcCtx)

	req := &types.PaymentCallbackReq{
		OrderId: 999,
		TradeNo: "TRADE_123456789",
		Amount:  100.00,
	}

	err := logic.PaymentCallback(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Order not found")
}

func TestPaymentCallbackLogic_AlreadyPaid(t *testing.T) {
	db := setupTestDB(t)

	order := &model.Order{
		Id:         1,
		CampaignId: 1,
		Phone:      "13800138000",
		FormData:   `{"name":"张三"}`,
		Status:     "pending",
		PayStatus:  "paid",
		Amount:     100.00,
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewPaymentCallbackLogic(ctx, svcCtx)

	req := &types.PaymentCallbackReq{
		OrderId: 1,
		TradeNo: "TRADE_123456789",
		Amount:  100.00,
	}

	err := logic.PaymentCallback(req)

	assert.NoError(t, err)
}

func TestPaymentCallbackLogic_NoReferrer(t *testing.T) {
	db := setupTestDB(t)

	campaign := &model.Campaign{
		Name:               "测试活动",
		Description:        "这是一个测试活动",
		FormFields:         `[{"type":"text","name":"name","label":"姓名","required":true}]`,
		RewardRule:         10.00,
		StartTime:          time.Now().Add(-1 * time.Hour),
		EndTime:            time.Now().Add(24 * time.Hour),
		Status:             "active",
		BrandId:            1,
		EnableDistribution: true,
	}
	db.Create(campaign)

	order := &model.Order{
		Id:         1,
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData:   `{"name":"张三"}`,
		ReferrerId: 0,
		Status:     "pending",
		PayStatus:  "unpaid",
		Amount:     100.00,
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewPaymentCallbackLogic(ctx, svcCtx)

	req := &types.PaymentCallbackReq{
		OrderId: 1,
		TradeNo: "TRADE_123456789",
		Amount:  100.00,
	}

	err := logic.PaymentCallback(req)

	assert.NoError(t, err)

	var updatedOrder model.Order
	db.First(&updatedOrder, order.Id)
	assert.Equal(t, "paid", updatedOrder.PayStatus)

	var rewardCount int64
	db.Model(&model.DistributorReward{}).Where("order_id = ?", order.Id).Count(&rewardCount)
	assert.Equal(t, int64(0), rewardCount)
}

func TestPaymentCallbackLogic_DistributionDisabled(t *testing.T) {
	db := setupTestDB(t)

	campaign := &model.Campaign{
		Name:               "测试活动",
		Description:        "这是一个测试活动",
		FormFields:         `[{"type":"text","name","name","label":"姓名","required":true}]`,
		RewardRule:         10.00,
		StartTime:          time.Now().Add(-1 * time.Hour),
		EndTime:            time.Now().Add(24 * time.Hour),
		Status:             "active",
		BrandId:            1,
		EnableDistribution: false,
	}
	db.Create(campaign)

	order := &model.Order{
		Id:         1,
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData:   `{"name":"张三"}`,
		ReferrerId: 100,
		Status:     "pending",
		PayStatus:  "unpaid",
		Amount:     100.00,
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewPaymentCallbackLogic(ctx, svcCtx)

	req := &types.PaymentCallbackReq{
		OrderId: 1,
		TradeNo: "TRADE_123456789",
		Amount:  100.00,
	}

	err := logic.PaymentCallback(req)

	assert.NoError(t, err)

	var updatedOrder model.Order
	db.First(&updatedOrder, order.Id)
	assert.Equal(t, "paid", updatedOrder.PayStatus)

	var rewardCount int64
	db.Model(&model.DistributorReward{}).Where("order_id = ?", order.Id).Count(&rewardCount)
	assert.Equal(t, int64(0), rewardCount)
}

func TestPaymentCallbackLogic_ReferrerNotFound(t *testing.T) {
	db := setupTestDB(t)

	campaign := &model.Campaign{
		Name:               "测试活动",
		Description:        "这是一个测试活动",
		FormFields:         `[{"type":"text","name":"name","label":"姓名","required":true}]`,
		RewardRule:         10.00,
		StartTime:          time.Now().Add(-1 * time.Hour),
		EndTime:            time.Now().Add(24 * time.Hour),
		Status:             "active",
		BrandId:            1,
		EnableDistribution: true,
	}
	db.Create(campaign)

	order := &model.Order{
		Id:         1,
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData:   `{"name":"张三"}`,
		ReferrerId: 100,
		Status:     "pending",
		PayStatus:  "unpaid",
		Amount:     100.00,
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewPaymentCallbackLogic(ctx, svcCtx)

	req := &types.PaymentCallbackReq{
		OrderId: 1,
		TradeNo: "TRADE_123456789",
		Amount:  100.00,
	}

	err := logic.PaymentCallback(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Referrer not found")
}

func TestScanOrderLogic_Success(t *testing.T) {
	db := setupTestDB(t)

	createLogic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	timestamp := time.Now().Unix()
	verificationCode := createLogic.generateVerificationCode(1, "13800138000", timestamp)

	order := &model.Order{
		Id:                 1,
		CampaignId:         1,
		Phone:              "13800138000",
		FormData:           `{"name":"张三","age":"25"}`,
		Status:             "paid",
		PayStatus:          "paid",
		VerificationStatus: "unverified",
		VerificationCode:   verificationCode,
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewScanOrderLogic(ctx, svcCtx)

	req := &types.ScanOrderReq{
		Code: verificationCode,
	}

	resp, err := logic.ScanOrder(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, order.Id, resp.OrderId)
	assert.Equal(t, "unverified", resp.Status)
	assert.Equal(t, "paid", resp.PayStatus)
	assert.Equal(t, "13800138000", resp.Phone)
	assert.Contains(t, resp.FormData, "name")
}

func TestScanOrderLogic_InvalidCode(t *testing.T) {
	db := setupTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewScanOrderLogic(ctx, svcCtx)

	req := &types.ScanOrderReq{
		Code: "invalid_code",
	}

	resp, err := logic.ScanOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "核销码无效")
}

func TestScanOrderLogic_OrderNotFound(t *testing.T) {
	db := setupTestDB(t)

	createLogic := NewCreateOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	timestamp := time.Now().Unix()
	verificationCode := createLogic.generateVerificationCode(999, "13800138000", timestamp)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewScanOrderLogic(ctx, svcCtx)

	req := &types.ScanOrderReq{
		Code: verificationCode,
	}

	resp, err := logic.ScanOrder(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "订单不存在")
}

func TestGetOrderLogic_Success(t *testing.T) {
	db := setupTestDB(t)
	campaign := createTestCampaign(t, db)
	order := &model.Order{
		Id:         1,
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData:   `{"name":"张三","age":25}`,
		ReferrerId: 0,
		Status:     "pending",
		Amount:     100,
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetOrderLogic(ctx, svcCtx)

	resp, err := logic.GetOrder(1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, campaign.Id, resp.CampaignId)
	assert.Equal(t, "13800138000", resp.Phone)
	assert.Equal(t, "pending", resp.Status)
	assert.Equal(t, float64(100), resp.Amount)
	assert.NotEmpty(t, resp.FormData)
}

func TestGetOrderLogic_NotFound(t *testing.T) {
	db := setupTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetOrderLogic(ctx, svcCtx)

	resp, err := logic.GetOrder(999)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetOrdersLogic_Success(t *testing.T) {
	db := setupTestDB(t)
	campaign := createTestCampaign(t, db)

	order1 := &model.Order{
		Id:         1,
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		FormData:   `{"name":"张三"}`,
		Status:     "pending",
		Amount:     100,
	}
	db.Create(order1)

	order2 := &model.Order{
		Id:         2,
		CampaignId: campaign.Id,
		Phone:      "13800138001",
		FormData:   `{"name":"李四"}`,
		Status:     "paid",
		Amount:     200,
	}
	db.Create(order2)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetOrdersLogic(ctx, svcCtx)

	resp, err := logic.GetOrders()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Orders, 2)
}

func TestGetOrdersLogic_Empty(t *testing.T) {
	db := setupTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetOrdersLogic(ctx, svcCtx)

	resp, err := logic.GetOrders()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Orders, 0)
}

func TestGetVerificationRecordsLogic_Success(t *testing.T) {
	db := setupTestDB(t)

	now := time.Now()
	record1 := &model.VerificationRecord{
		ID:                 1,
		OrderID:            1,
		VerificationStatus: "verified",
		VerifiedAt:         &now,
		VerifiedBy:         int64Ptr(100),
		VerificationCode:   "CODE_123",
		VerificationMethod: "manual",
		Remark:             "核销成功",
	}
	db.Create(record1)

	record2 := &model.VerificationRecord{
		ID:                 2,
		OrderID:            2,
		VerificationStatus: "cancelled",
		VerificationCode:   "CODE_456",
		VerificationMethod: "manual",
		Remark:             "取消核销",
	}
	db.Create(record2)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetVerificationRecordsLogic(ctx, svcCtx)

	resp, err := logic.GetVerificationRecords()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Records, 2)
}

func TestGetVerificationRecordsLogic_Empty(t *testing.T) {
	db := setupTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetVerificationRecordsLogic(ctx, svcCtx)

	resp, err := logic.GetVerificationRecords()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Records, 0)
}

func int64Ptr(i int64) *int64 {
	return &i
}
