package order

import (
	"context"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verificationAdminCtx(userID int64) context.Context {
	ctx := context.WithValue(context.Background(), "roles", []string{"brand_admin"})
	if userID > 0 {
		ctx = context.WithValue(ctx, "userId", userID)
	}
	return ctx
}

func createVerificationCode(orderID int64, phone string) string {
	createLogic := &CreateOrderLogic{}
	return createLogic.generateVerificationCode(orderID, phone, time.Now().Unix())
}

func insertOrderForVerification(t *testing.T, dbCtx *svc.ServiceContext, status string, verificationStatus string) *model.Order {
	order := &model.Order{
		CampaignId:         1,
		Phone:              "13800138000",
		FormData:           `{"name":"测试"}`,
		Status:             status,
		PayStatus:          "paid",
		VerificationStatus: verificationStatus,
	}
	require.NoError(t, dbCtx.DB.Create(order).Error)

	code := createVerificationCode(order.Id, order.Phone)
	require.NoError(t, dbCtx.DB.Model(order).Update("verification_code", code).Error)
	order.VerificationCode = code
	return order
}

func TestVerifyOrderLogic_InvalidCode(t *testing.T) {
	db := setupTestDB(t)

	logic := NewVerifyOrderLogic(context.Background(), &svc.ServiceContext{DB: db})
	resp, err := logic.VerifyOrder(&types.VerifyOrderReq{Code: "invalid_code"})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "核销码无效")
}

func TestVerifyOrderLogic_OrderNotFound(t *testing.T) {
	db := setupTestDB(t)

	logic := NewVerifyOrderLogic(verificationAdminCtx(0), &svc.ServiceContext{DB: db})
	missingCode := createVerificationCode(9999, "13800138000")

	resp, err := logic.VerifyOrder(&types.VerifyOrderReq{Code: missingCode})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "订单不存在")
}

func TestVerifyOrderLogic_CreatesVerificationRecord(t *testing.T) {
	db := setupTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	order := insertOrderForVerification(t, svcCtx, "paid", "unverified")

	logic := NewVerifyOrderLogic(verificationAdminCtx(777), svcCtx)

	resp, err := logic.VerifyOrder(&types.VerifyOrderReq{Code: order.VerificationCode})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "verified", resp.Status)

	var freshOrder model.Order
	require.NoError(t, db.First(&freshOrder, order.Id).Error)
	assert.Equal(t, "verified", freshOrder.VerificationStatus)
	require.NotNil(t, freshOrder.VerifiedBy)
	assert.Equal(t, int64(777), *freshOrder.VerifiedBy)

	var records []model.VerificationRecord
	require.NoError(t, db.Where("order_id = ?", order.Id).Find(&records).Error)
	require.Len(t, records, 1)
	assert.Equal(t, "verified", records[0].VerificationStatus)
	assert.Equal(t, "品牌管理员核销", records[0].Remark)
}

func TestVerifyOrderLogic_PermissionDenied(t *testing.T) {
	db := setupTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	order := insertOrderForVerification(t, svcCtx, "paid", "unverified")

	ctx := context.WithValue(context.Background(), "roles", []string{"participant"})
	logic := NewVerifyOrderLogic(ctx, svcCtx)

	resp, err := logic.VerifyOrder(&types.VerifyOrderReq{Code: order.VerificationCode})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "权限不足")
}

func TestUnverifyOrderLogic_Success(t *testing.T) {
	db := setupTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	order := insertOrderForVerification(t, svcCtx, "paid", "verified")

	logic := NewUnverifyOrderLogic(verificationAdminCtx(0), svcCtx)
	resp, err := logic.UnverifyOrder(&types.UnverifyOrderReq{Code: order.VerificationCode})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "unverified", resp.Status)

	var freshOrder model.Order
	require.NoError(t, db.First(&freshOrder, order.Id).Error)
	assert.Equal(t, "cancelled", freshOrder.VerificationStatus)
	assert.Nil(t, freshOrder.VerifiedBy)

	var records []model.VerificationRecord
	require.NoError(t, db.Where("order_id = ?", order.Id).Find(&records).Error)
	require.Len(t, records, 1)
	assert.Equal(t, "cancelled", records[0].VerificationStatus)
	assert.Equal(t, "品牌管理员取消核销", records[0].Remark)
}

func TestUnverifyOrderLogic_OrderNotVerified(t *testing.T) {
	db := setupTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	order := insertOrderForVerification(t, svcCtx, "paid", "unverified")

	logic := NewUnverifyOrderLogic(verificationAdminCtx(0), svcCtx)
	resp, err := logic.UnverifyOrder(&types.UnverifyOrderReq{Code: order.VerificationCode})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "尚未核销")
}

func TestUnverifyOrderLogic_PermissionDenied(t *testing.T) {
	db := setupTestDB(t)

	svcCtx := &svc.ServiceContext{DB: db}
	order := insertOrderForVerification(t, svcCtx, "paid", "verified")

	ctx := context.WithValue(context.Background(), "roles", []string{"participant"})
	logic := NewUnverifyOrderLogic(ctx, svcCtx)

	resp, err := logic.UnverifyOrder(&types.UnverifyOrderReq{Code: order.VerificationCode})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "权限不足")
}

func TestUnverifyOrderLogic_OrderNotFound(t *testing.T) {
	db := setupTestDB(t)

	logic := NewUnverifyOrderLogic(verificationAdminCtx(0), &svc.ServiceContext{DB: db})
	missingCode := createVerificationCode(9999, "13800138000")

	resp, err := logic.UnverifyOrder(&types.UnverifyOrderReq{Code: missingCode})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "订单不存在")
}
