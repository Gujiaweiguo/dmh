package syncadapter

import (
	"context"
	"testing"
	"time"

	"dmh/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zeromicro/go-zero/core/logx"
)

// Helper to set up an in-memory sqlite DB for Order/Reward models is already provided in syncworker_test.go

func TestSyncWorker_SyncOrder_ViaAdapter(t *testing.T) {
	// Create a sqlite gorm DB and a corresponding sqlmock for adapter
	gormDB := setupSyncWorkerTestDB(t)

	// Create a test order in GORM DB
	order := &model.Order{Id: 123, CampaignId: 456, MemberID: nil, UnionID: "wx_unionid", Phone: "13800138000", FormData: `{"name":"test"}`, Amount: 100.0, PayStatus: "paid", CreatedAt: time.Now()}
	gormDB.Create(order)

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	adapter := &SyncAdapter{db: sqlDB, mapper: NewFieldMapper(), metrics: NewSyncMetrics(), logger: logx.WithContext(context.Background())}

	// Reuse sqlite DB for lookups, and pass to worker
	worker := NewSyncWorker(adapter, &SyncQueue{}, gormDB)

	mock.ExpectExec("INSERT INTO external_orders").
		WithArgs(int64(123), int64(456), int64(0), "wx_unionid", "13800138000", sqlmock.AnyArg(), 100.0, "paid", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	worker.syncOrder(123)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestSyncWorker_SyncReward_ViaAdapter(t *testing.T) {
	gormDB := setupSyncWorkerTestDB(t)
	reward := &model.Reward{Id: 888, UserId: 999, MemberID: nil, OrderId: 777, Amount: 12.34, Status: "settled", SettledAt: func() *time.Time { t0 := time.Now(); return &t0 }()}
	gormDB.Create(reward)

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	adapter := &SyncAdapter{db: sqlDB, mapper: NewFieldMapper(), metrics: NewSyncMetrics(), logger: logx.WithContext(context.Background())}
	worker := NewSyncWorker(adapter, &SyncQueue{}, gormDB)

	mock.ExpectExec("INSERT INTO external_rewards").
		WithArgs(int64(888), int64(999), int64(0), int64(777), 12.34, "settled", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call private method
	_ = worker.syncReward(888)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}
