package sync

import (
	"context"
	"testing"

	"dmh/api/internal/svc"

	"github.com/stretchr/testify/assert"
)

func TestSyncLogicConstructors(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}

	assert.NotNil(t, NewGetSyncHealthLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetSyncStatusLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetSyncStatsLogic(ctx, svcCtx))
	assert.NotNil(t, NewRetrySynLogic(ctx, svcCtx))
}

func TestGetSyncHealthBehavior(t *testing.T) {
	logic := NewGetSyncHealthLogic(context.Background(), &svc.ServiceContext{})
	resp, err := logic.GetSyncHealth()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "healthy", resp.Status)
	assert.NotNil(t, resp.Database)
	assert.NotNil(t, resp.Queue)
}

func TestSyncTodoMethodsCurrentBehavior(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}

	status := NewGetSyncStatusLogic(ctx, svcCtx)
	statusResp, statusErr := status.GetSyncStatus()
	assert.NoError(t, statusErr)
	assert.Nil(t, statusResp)

	stats := NewGetSyncStatsLogic(ctx, svcCtx)
	statsResp, statsErr := stats.GetSyncStats()
	assert.NoError(t, statsErr)
	assert.Nil(t, statsResp)

	retry := NewRetrySynLogic(ctx, svcCtx)
	retryResp, retryErr := retry.RetrySyn()
	assert.NoError(t, retryErr)
	assert.Nil(t, retryResp)
}
