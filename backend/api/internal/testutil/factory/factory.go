// Package factory 提供统一的测试数据工厂模式
// 用于减少重复的测试数据创建代码，提高测试可维护性
package factory

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

// Factory 定义测试数据工厂的基础接口
type Factory[T any] interface {
	// Build 创建一个实例，使用默认值
	Build() *T

	// BuildWith 创建一个实例，允许覆盖指定字段
	BuildWith(overrides map[string]any) *T

	// BuildList 创建多个实例
	BuildList(count int) []*T

	// Create 创建并持久化到数据库
	Create(db any) (*T, error)

	// CreateWith 创建并持久化到数据库，允许覆盖指定字段
	CreateWith(db any, overrides map[string]any) (*T, error)

	// CreateList 创建并持久化多个实例
	CreateList(db any, count int) ([]*T, error)
}

// BaseFactory 提供工厂的通用功能
type BaseFactory struct {
	counter uint64
}

// NextID 生成唯一递增 ID（用于测试）
func (f *BaseFactory) NextID() int64 {
	return int64(atomic.AddUint64(&f.counter, 1))
}

// RandomSuffix 生成随机后缀，避免唯一约束冲突
func (f *BaseFactory) RandomSuffix() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// RandomPhone 生成随机手机号
func (f *BaseFactory) RandomPhone() string {
	return fmt.Sprintf("138%08d", time.Now().UnixNano()%100000000)
}

// RandomEmail 生成随机邮箱
func (f *BaseFactory) RandomEmail() string {
	return fmt.Sprintf("test_%s@example.com", f.RandomSuffix())
}

// RandomUsername 生成随机用户名
func (f *BaseFactory) RandomUsername(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, f.RandomSuffix())
}

// Now 返回当前时间（便于测试中 mock）
func (f *BaseFactory) Now() time.Time {
	return time.Now()
}

// FutureTime 返回未来某个时间点
func (f *BaseFactory) FutureTime(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

// PastTime 返回过去某个时间点
func (f *BaseFactory) PastTime(duration time.Duration) time.Time {
	return time.Now().Add(-duration)
}

// ========== 常量定义 ==========

// 用户角色
const (
	RolePlatformAdmin = "platform_admin"
	RoleBrandAdmin    = "brand_admin"
	RoleParticipant   = "participant"
)

// 用户状态
const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusLocked   = "locked"
)

// 活动状态
const (
	CampaignStatusActive = "active"
	CampaignStatusPaused = "paused"
	CampaignStatusEnded  = "ended"
)

// 订单状态
const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusCancelled = "cancelled"
)

// 支付状态
const (
	PayStatusUnpaid   = "unpaid"
	PayStatusPaid     = "paid"
	PayStatusRefunded = "refunded"
)
