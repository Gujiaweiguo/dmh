package order

import (
	"context"
	"sync"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderIntegration_CreateVerifyScanConsistency(t *testing.T) {
	db := setupTestDB(t)
	campaign := createTestCampaign(t, db)
	svcCtx := &svc.ServiceContext{DB: db}

	createLogic := NewCreateOrderLogic(context.Background(), svcCtx)
	createResp, err := createLogic.CreateOrder(&types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13900000001",
		FormData: map[string]string{
			"name": "集成测试用户",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, createResp)

	var created model.Order
	require.NoError(t, db.First(&created, createResp.Id).Error)
	require.NotEmpty(t, created.VerificationCode)

	verifyCtx := context.WithValue(context.Background(), "roles", []string{"brand_admin"})
	verifyCtx = context.WithValue(verifyCtx, "userId", int64(9001))
	verifyLogic := NewVerifyOrderLogic(verifyCtx, svcCtx)
	verifyResp, err := verifyLogic.VerifyOrder(&types.VerifyOrderReq{Code: created.VerificationCode, Remark: "integration verify"})
	require.NoError(t, err)
	require.NotNil(t, verifyResp)
	assert.Equal(t, "verified", verifyResp.Status)

	scanLogic := NewScanOrderLogic(context.Background(), svcCtx)
	scanResp, err := scanLogic.ScanOrder(&types.ScanOrderReq{Code: created.VerificationCode})
	require.NoError(t, err)
	require.NotNil(t, scanResp)
	assert.Equal(t, created.Id, scanResp.OrderId)
	assert.Equal(t, "verified", scanResp.Status)
	assert.Equal(t, "unpaid", scanResp.PayStatus)
	assert.Equal(t, "集成测试用户", scanResp.FormData["name"])
}

func TestOrderIntegration_ConcurrentDuplicateCreateGuard(t *testing.T) {
	db := setupTestDB(t)
	campaign := createTestCampaign(t, db)
	svcCtx := &svc.ServiceContext{DB: db}

	req := &types.CreateOrderReq{
		CampaignId: campaign.Id,
		Phone:      "13900000003",
		FormData: map[string]string{
			"name": "并发测试用户",
		},
	}

	start := make(chan struct{})
	errCh := make(chan error, 2)
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			logic := NewCreateOrderLogic(context.Background(), svcCtx)
			_, err := logic.CreateOrder(req)
			errCh <- err
		}()
	}

	close(start)
	wg.Wait()
	close(errCh)

	successCount := 0
	errCount := 0
	for err := range errCh {
		if err == nil {
			successCount++
		} else {
			errCount++
		}
	}

	assert.Equal(t, 1, successCount, "expected exactly one success")
	assert.Equal(t, 1, errCount, "expected exactly one error due to unique constraint")

	var total int64
	require.NoError(t, db.Model(&model.Order{}).Where("campaign_id = ? AND phone = ?", campaign.Id, req.Phone).Count(&total).Error)
	assert.Equal(t, int64(1), total, "should have exactly one order after concurrent attempts")
}
