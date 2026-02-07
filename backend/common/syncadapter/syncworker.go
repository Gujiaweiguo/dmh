// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package syncadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// SyncWorker 同步Worker
type SyncWorker struct {
	adapter *SyncAdapter
	queue   *SyncQueue
	db      *gorm.DB
	logger  logx.Logger
	stop    chan struct{}
}

// NewSyncWorker 创建同步Worker
func NewSyncWorker(adapter *SyncAdapter, queue *SyncQueue, db *gorm.DB) *SyncWorker {
	return &SyncWorker{
		adapter: adapter,
		queue:   queue,
		db:      db,
		logger:  logx.WithContext(context.Background()),
		stop:    make(chan struct{}),
	}
}

// Start 启动Worker
func (w *SyncWorker) Start() {
	w.logger.Info("SyncWorker started")

	for {
		select {
		case <-w.stop:
			w.logger.Info("SyncWorker stopping...")
			return
		default:
			task, err := w.queue.Dequeue(5 * time.Second)
			if err != nil {
				w.logger.Errorf("Failed to dequeue task: %v", err)
				continue
			}
			if task == nil {
				continue
			}

			if task.Type == "order" {
				w.syncOrder(task.OrderId)
			} else if task.Type == "reward" {
				w.syncReward(task.RewardId)
			}
		}
	}
}

// Stop 停止Worker
func (w *SyncWorker) Stop() {
	close(w.stop)
}

// syncOrder 同步订单
func (w *SyncWorker) syncOrder(orderId int64) {
	var order model.Order
	if err := w.db.Where("id = ? AND deleted_at IS NULL", orderId).First(&order).Error; err != nil {
		w.logger.Errorf("Failed to query order: %v", err)
		return
	}

	// 2. 解析表单数据
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(order.FormData), &formData); err != nil {
		formData = make(map[string]interface{})
	}

	var memberID int64
	if order.MemberID != nil {
		memberID = *order.MemberID
	}

	// 3. 构建同步数据
	syncData := &SyncOrderData{
		OrderId:    order.Id,
		CampaignId: order.CampaignId,
		MemberId:   memberID,
		UnionID:    order.UnionID,
		Phone:      order.Phone,
		FormData:   formData,
		Amount:     order.Amount,
		PayStatus:  order.PayStatus,
		CreatedAt:  order.CreatedAt,
	}

	// 4. 执行同步
	_ = w.adapter.SyncOrder(context.Background(), syncData)
}

// syncReward 同步奖励
func (w *SyncWorker) syncReward(rewardId int64) error {
	// 1. 查询奖励数据
	var reward model.Reward
	if err := w.db.Where("id = ?", rewardId).First(&reward).Error; err != nil {
		return fmt.Errorf("get reward failed: %w", err)
	}

	var rewardMemberID int64
	if reward.MemberID != nil {
		rewardMemberID = *reward.MemberID
	}

	// 2. 构建同步数据
	syncData := &SyncRewardData{
		RewardId: reward.Id,
		UserId:   reward.UserId,
		MemberId: rewardMemberID,
		OrderId:  reward.OrderId,
		Amount:   reward.Amount,
		Status:   reward.Status,
	}

	if reward.SettledAt != nil {
		syncData.SettledAt = *reward.SettledAt
	}

	// 3. 执行同步
	_ = w.adapter.SyncReward(context.Background(), syncData)
	return nil
}

// updateSyncStatus 更新同步状态
func (w *SyncWorker) updateSyncStatus(orderId int64, status, errorMsg string) {
	var syncedAt *time.Time
	if status == "synced" {
		now := time.Now()
		syncedAt = &now
	}

	// 创建或更新sync_logs记录
	log := model.SyncLog{
		OrderId:    orderId,
		SyncType:   "order",
		SyncStatus: status,
		ErrorMsg:   errorMsg,
		SyncedAt:   syncedAt,
		UpdatedAt:  time.Now(),
	}

	// 尝试查找现有记录
	var existing model.SyncLog
	result := w.db.Where("order_id = ? AND sync_type = ?", orderId, "order").First(&existing)

	if result.Error == nil {
		// 更新现有记录
		log.Id = existing.Id
		log.Attempts = existing.Attempts + 1
		w.db.Save(&log)
	} else {
		// 创建新记录
		w.db.Create(&log)
	}

	w.logger.Infof("Sync status updated: orderId=%d, status=%s, attempts=%d", orderId, status, log.Attempts)
}

// updateRewardSyncStatus 更新奖励同步状态
func (w *SyncWorker) updateRewardSyncStatus(rewardId int64, status, errorMsg string) {
	var syncedAt *time.Time
	if status == "synced" {
		now := time.Now()
		syncedAt = &now
	}

	// 创建或更新sync_logs记录
	log := model.SyncLog{
		OrderId:    rewardId,
		SyncType:   "reward",
		SyncStatus: status,
		ErrorMsg:   errorMsg,
		SyncedAt:   syncedAt,
		UpdatedAt:  time.Now(),
	}

	// 尝试查找现有记录
	var existing model.SyncLog
	result := w.db.Where("order_id = ? AND sync_type = ?", rewardId, "reward").First(&existing)

	if result.Error == nil {
		// 更新现有记录
		log.Id = existing.Id
		log.Attempts = existing.Attempts + 1
		w.db.Save(&log)
	} else {
		// 创建新记录
		w.db.Create(&log)
	}

	w.logger.Infof("Reward sync status updated: rewardId=%d, status=%s, attempts=%d", rewardId, status, log.Attempts)
}
