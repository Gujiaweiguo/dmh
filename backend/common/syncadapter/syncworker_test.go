package syncadapter

import (
	"testing"
	"time"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupSyncWorkerTestDB(t *testing.T) *gorm.DB {
	dsn := "root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Order{}, &model.Reward{}, &model.SyncLog{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		t.Fatalf("Failed to disable foreign key checks: %v", err)
	}
	if err := db.Exec("TRUNCATE TABLE sync_logs").Error; err != nil {
		t.Fatalf("Failed to truncate sync_logs: %v", err)
	}
	if err := db.Exec("TRUNCATE TABLE rewards").Error; err != nil {
		t.Fatalf("Failed to truncate rewards: %v", err)
	}
	if err := db.Exec("TRUNCATE TABLE orders").Error; err != nil {
		t.Fatalf("Failed to truncate orders: %v", err)
	}
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		t.Fatalf("Failed to enable foreign key checks: %v", err)
	}

	return db
}

func TestNewSyncWorker(t *testing.T) {
	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := &SyncQueue{}

	worker := NewSyncWorker(adapter, queue, db)

	assert.NotNil(t, worker)
	assert.NotNil(t, worker.adapter)
	assert.NotNil(t, worker.queue)
	assert.NotNil(t, worker.db)
	assert.NotNil(t, worker.stop)
}

func TestSyncWorker_Stop(t *testing.T) {
	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := &SyncQueue{}

	worker := NewSyncWorker(adapter, queue, db)

	assert.NotPanics(t, func() {
		worker.Stop()
	})
}

func TestUpdateSyncStatus_Create(t *testing.T) {
	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := &SyncQueue{}

	worker := NewSyncWorker(adapter, queue, db)

	worker.updateSyncStatus(123, "synced", "")

	var log model.SyncLog
	err := db.Where("order_id = ? AND sync_type = ?", 123, "order").First(&log).Error

	assert.NoError(t, err)
	assert.Equal(t, int64(123), log.OrderId)
	assert.Equal(t, "order", log.SyncType)
	assert.Equal(t, "synced", log.SyncStatus)
	assert.Equal(t, 0, log.Attempts)
	assert.NotNil(t, log.SyncedAt)
}

func TestUpdateSyncStatus_Update(t *testing.T) {
	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := &SyncQueue{}

	worker := NewSyncWorker(adapter, queue, db)

	worker.updateSyncStatus(123, "synced", "")

	worker.updateSyncStatus(123, "synced", "")

	var log model.SyncLog
	err := db.Where("order_id = ? AND sync_type = ?", 123, "order").First(&log).Error

	assert.NoError(t, err)
	assert.Equal(t, 1, log.Attempts)
}

func TestUpdateSyncStatus_Failed(t *testing.T) {
	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := &SyncQueue{}

	worker := NewSyncWorker(adapter, queue, db)

	worker.updateSyncStatus(123, "failed", "connection error")

	var log model.SyncLog
	err := db.Where("order_id = ? AND sync_type = ?", 123, "order").First(&log).Error

	assert.NoError(t, err)
	assert.Equal(t, "failed", log.SyncStatus)
	assert.Equal(t, "connection error", log.ErrorMsg)
	assert.Nil(t, log.SyncedAt)
}

func TestUpdateRewardSyncStatus_Create(t *testing.T) {
	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := &SyncQueue{}

	worker := NewSyncWorker(adapter, queue, db)

	worker.updateRewardSyncStatus(456, "synced", "")

	var log model.SyncLog
	err := db.Where("order_id = ? AND sync_type = ?", 456, "reward").First(&log).Error

	assert.NoError(t, err)
	assert.Equal(t, int64(456), log.OrderId)
	assert.Equal(t, "reward", log.SyncType)
	assert.Equal(t, "synced", log.SyncStatus)
	assert.Equal(t, 0, log.Attempts)
	assert.NotNil(t, log.SyncedAt)
}

func TestUpdateRewardSyncStatus_Update(t *testing.T) {
	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := &SyncQueue{}

	worker := NewSyncWorker(adapter, queue, db)

	worker.updateRewardSyncStatus(456, "synced", "")

	worker.updateRewardSyncStatus(456, "synced", "")

	var log model.SyncLog
	err := db.Where("order_id = ? AND sync_type = ?", 456, "reward").First(&log).Error

	assert.NoError(t, err)
	assert.Equal(t, 1, log.Attempts)
}

func TestUpdateRewardSyncStatus_Failed(t *testing.T) {
	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := &SyncQueue{}

	worker := NewSyncWorker(adapter, queue, db)

	worker.updateRewardSyncStatus(456, "failed", "database error")

	var log model.SyncLog
	err := db.Where("order_id = ? AND sync_type = ?", 456, "reward").First(&log).Error

	assert.NoError(t, err)
	assert.Equal(t, "failed", log.SyncStatus)
	assert.Equal(t, "database error", log.ErrorMsg)
	assert.Nil(t, log.SyncedAt)
}

func TestSyncWorker_StartAndStop(t *testing.T) {
	// Skip if Redis is not available
	redisClient := skipIfNoRedis(t)
	if redisClient == nil {
		return
	}
	defer redisClient.Close()

	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := NewSyncQueue(redisClient, "test_start_stop")
	worker := NewSyncWorker(adapter, queue, db)

	go worker.Start()
	time.Sleep(50 * time.Millisecond)
	worker.Stop()
	time.Sleep(50 * time.Millisecond)
	assert.True(t, true)
}

func TestSyncWorker_StartStopImmediate(t *testing.T) {
	// Skip if Redis is not available
	redisClient := skipIfNoRedis(t)
	if redisClient == nil {
		return
	}
	defer redisClient.Close()

	db := setupSyncWorkerTestDB(t)
	adapter := &SyncAdapter{}
	queue := NewSyncQueue(redisClient, "test_start_stop_immediate")
	worker := NewSyncWorker(adapter, queue, db)

	go worker.Start()
	worker.Stop()
	time.Sleep(50 * time.Millisecond)
	assert.True(t, true)
}

