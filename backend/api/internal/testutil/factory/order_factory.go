package factory

import (
	"dmh/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ErrCampaignIdRequired 创建订单时缺少活动ID的错误
var ErrCampaignIdRequired = errors.New("campaign_id is required for order creation")

// OrderFactory 订单测试数据工厂
type OrderFactory struct {
	BaseFactory
}

// NewOrderFactory 创建订单工厂实例
func NewOrderFactory() *OrderFactory {
	return &OrderFactory{}
}

// parseTimePointer 解析时间值，支持 time.Time 和 *time.Time 两种类型
func parseTimePointer(v any) *time.Time {
	if pt, ok := v.(*time.Time); ok {
		return pt
	}
	if t, ok := v.(time.Time); ok {
		return &t
	}
	return nil
}

// Build 创建一个订单实例，使用合理的默认值
// 注意：CampaignId 必须通过 overrides 提供
func (f *OrderFactory) Build() *model.Order {
	return f.BuildWith(nil)
}

// BuildWith 创建一个订单实例，允许覆盖指定字段
func (f *OrderFactory) BuildWith(overrides map[string]any) *model.Order {
	order := &model.Order{
		CampaignId:         0, // 必须通过 overrides 提供
		Phone:              f.RandomPhone(),
		FormData:           `{"name":"测试用户","phone":"` + f.RandomPhone() + `"}`,
		ReferrerId:         0,
		DistributorPath:    "",
		Status:             OrderStatusPending,
		Amount:             0.00,
		PayStatus:          PayStatusUnpaid,
		SyncStatus:         "pending",
		VerificationStatus: "unverified",
	}

	// 应用覆盖值
	if overrides != nil {
		if v, ok := overrides["id"]; ok {
			order.Id = v.(int64)
		}
		if v, ok := overrides["campaign_id"]; ok {
			order.CampaignId = v.(int64)
		}
		if v, ok := overrides["campaignId"]; ok {
			order.CampaignId = v.(int64)
		}
		if v, ok := overrides["member_id"]; ok {
			mid := v.(int64)
			order.MemberID = &mid
		}
		if v, ok := overrides["memberId"]; ok {
			mid := v.(int64)
			order.MemberID = &mid
		}
		if v, ok := overrides["unionid"]; ok {
			order.UnionID = v.(string)
		}
		if v, ok := overrides["phone"]; ok {
			order.Phone = v.(string)
		}
		if v, ok := overrides["form_data"]; ok {
			order.FormData = v.(string)
		}
		if v, ok := overrides["formData"]; ok {
			order.FormData = v.(string)
		}
		if v, ok := overrides["referrer_id"]; ok {
			order.ReferrerId = v.(int64)
		}
		if v, ok := overrides["referrerId"]; ok {
			order.ReferrerId = v.(int64)
		}
		if v, ok := overrides["distributor_path"]; ok {
			order.DistributorPath = v.(string)
		}
		if v, ok := overrides["distributorPath"]; ok {
			order.DistributorPath = v.(string)
		}
		if v, ok := overrides["status"]; ok {
			order.Status = v.(string)
		}
		if v, ok := overrides["amount"]; ok {
			order.Amount = v.(float64)
		}
		if v, ok := overrides["pay_status"]; ok {
			order.PayStatus = v.(string)
		}
		if v, ok := overrides["payStatus"]; ok {
			order.PayStatus = v.(string)
		}
		if v, ok := overrides["trade_no"]; ok {
			order.TradeNo = v.(string)
		}
		if v, ok := overrides["tradeNo"]; ok {
			order.TradeNo = v.(string)
		}
		if v, ok := overrides["paid_at"]; ok {
			order.PaidAt = parseTimePointer(v)
		}
		if v, ok := overrides["paidAt"]; ok {
			order.PaidAt = parseTimePointer(v)
		}
		if v, ok := overrides["sync_status"]; ok {
			order.SyncStatus = v.(string)
		}
		if v, ok := overrides["syncStatus"]; ok {
			order.SyncStatus = v.(string)
		}
		if v, ok := overrides["verification_status"]; ok {
			order.VerificationStatus = v.(string)
		}
		if v, ok := overrides["verificationStatus"]; ok {
			order.VerificationStatus = v.(string)
		}
		if v, ok := overrides["verified_at"]; ok {
			order.VerifiedAt = parseTimePointer(v)
		}
		if v, ok := overrides["verifiedAt"]; ok {
			order.VerifiedAt = parseTimePointer(v)
		}
		if v, ok := overrides["verified_by"]; ok {
			uid := v.(int64)
			order.VerifiedBy = &uid
		}
		if v, ok := overrides["verifiedBy"]; ok {
			uid := v.(int64)
			order.VerifiedBy = &uid
		}
		if v, ok := overrides["verification_code"]; ok {
			order.VerificationCode = v.(string)
		}
		if v, ok := overrides["verificationCode"]; ok {
			order.VerificationCode = v.(string)
		}
	}

	return order
}

// BuildList 创建多个订单实例
func (f *OrderFactory) BuildList(count int) []*model.Order {
	orders := make([]*model.Order, count)
	for i := 0; i < count; i++ {
		orders[i] = f.Build()
	}
	return orders
}

// Create 创建并持久化订单到数据库
func (f *OrderFactory) Create(db *gorm.DB) (*model.Order, error) {
	return f.CreateWith(db, nil)
}

// CreateWith 创建并持久化订单到数据库，允许覆盖指定字段
func (f *OrderFactory) CreateWith(db *gorm.DB, overrides map[string]any) (*model.Order, error) {
	order := f.BuildWith(overrides)

	// 验证必填字段
	if order.CampaignId == 0 {
		return nil, ErrCampaignIdRequired
	}

	if err := db.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// CreateList 创建并持久化多个订单到数据库
func (f *OrderFactory) CreateList(db *gorm.DB, count int, campaignId int64) ([]*model.Order, error) {
	orders := make([]*model.Order, count)
	for i := 0; i < count; i++ {
		order, err := f.CreateWith(db, map[string]any{"campaign_id": campaignId})
		if err != nil {
			return nil, err
		}
		orders[i] = order
	}
	return orders, nil
}

// ========== 便捷方法：创建特定类型订单 ==========

// BuildPending 创建待支付订单（不持久化）
func (f *OrderFactory) BuildPending(campaignId int64) *model.Order {
	return f.BuildWith(map[string]any{
		"campaign_id": campaignId,
		"status":      OrderStatusPending,
		"pay_status":  PayStatusUnpaid,
	})
}

// CreatePending 创建并持久化待支付订单
func (f *OrderFactory) CreatePending(db *gorm.DB, campaignId int64) (*model.Order, error) {
	return f.CreateWith(db, map[string]any{
		"campaign_id": campaignId,
		"status":      OrderStatusPending,
		"pay_status":  PayStatusUnpaid,
	})
}

// BuildPaid 创建已支付订单（不持久化）
func (f *OrderFactory) BuildPaid(campaignId int64) *model.Order {
	now := f.Now()
	return f.BuildWith(map[string]any{
		"campaign_id": campaignId,
		"status":      OrderStatusPaid,
		"pay_status":  PayStatusPaid,
		"amount":      100.00,
		"paid_at":     &now,
	})
}

// CreatePaid 创建并持久化已支付订单
func (f *OrderFactory) CreatePaid(db *gorm.DB, campaignId int64) (*model.Order, error) {
	now := f.Now()
	return f.CreateWith(db, map[string]any{
		"campaign_id": campaignId,
		"status":      OrderStatusPaid,
		"pay_status":  PayStatusPaid,
		"amount":      100.00,
		"paid_at":     &now,
	})
}

// BuildVerified 创建已核销订单（不持久化）
func (f *OrderFactory) BuildVerified(campaignId int64, verifierId int64) *model.Order {
	now := f.Now()
	return f.BuildWith(map[string]any{
		"campaign_id":         campaignId,
		"status":              OrderStatusPaid,
		"pay_status":          PayStatusPaid,
		"amount":              100.00,
		"paid_at":             &now,
		"verification_status": "verified",
		"verified_at":         &now,
		"verified_by":         verifierId,
		"verification_code":   "TEST_CODE_" + f.RandomSuffix(),
	})
}

// CreateVerified 创建并持久化已核销订单
func (f *OrderFactory) CreateVerified(db *gorm.DB, campaignId int64, verifierId int64) (*model.Order, error) {
	now := f.Now()
	return f.CreateWith(db, map[string]any{
		"campaign_id":         campaignId,
		"status":              OrderStatusPaid,
		"pay_status":          PayStatusPaid,
		"amount":              100.00,
		"paid_at":             &now,
		"verification_status": "verified",
		"verified_at":         &now,
		"verified_by":         verifierId,
		"verification_code":   "TEST_CODE_" + f.RandomSuffix(),
	})
}
