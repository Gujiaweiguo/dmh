package model

import (
	"time"
)

// Campaign 营销活动模型
type Campaign struct {
	Id                  int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BrandId             int64      `gorm:"column:brand_id;not null;index" json:"brandId"`
	Name                string     `gorm:"column:name;type:varchar(200);not null" json:"name"`
	Description         string     `gorm:"column:description;type:text" json:"description"`
	FormFields          string     `gorm:"column:form_fields;type:json" json:"formFields"` // JSON格式存储
	RewardRule          float64    `gorm:"column:reward_rule;type:decimal(10,2);not null;default:0.00" json:"rewardRule"`
	StartTime           time.Time  `gorm:"column:start_time;not null" json:"startTime"`
	EndTime             time.Time  `gorm:"column:end_time;not null" json:"endTime"`
	Status              string     `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"`        // active, paused, ended
	EnableDistribution  bool       `gorm:"column:enable_distribution;not null;default:false;index" json:"enableDistribution"` // 是否启用分销
	DistributionLevel   int        `gorm:"column:distribution_level;not null;default:1" json:"distributionLevel"`             // 分销层级(1/2/3)
	DistributionRewards *string    `gorm:"column:distribution_rewards;type:json" json:"distributionRewards,omitempty"`        // 各级奖励比例
	CreatedAt           time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt           *time.Time `gorm:"column:deleted_at" json:"deletedAt,omitempty"`
}

// TableName 表名
func (m *Campaign) TableName() string {
	return "campaigns"
}

// Order 订单模型
type Order struct {
	Id              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignId      int64      `gorm:"column:campaign_id;not null;index" json:"campaignId"`
	MemberID        *int64     `gorm:"column:member_id;index" json:"memberId"`                // 关联会员ID（可选）
	UnionID         string     `gorm:"column:unionid;type:varchar(100);index" json:"unionid"` // 微信 unionid
	Phone           string     `gorm:"column:phone;type:varchar(20);not null;index" json:"phone"`
	FormData        string     `gorm:"column:form_data;type:json" json:"formData"` // JSON格式存储
	ReferrerId      int64      `gorm:"column:referrer_id;default:0;index" json:"referrerId"`
	DistributorPath string     `gorm:"column:distributor_path;type:varchar(100);default:'';index:idx_distributor_path" json:"distributorPath"` // 分销链路径 "一级ID,二级ID,三级ID"
	Status          string     `gorm:"column:status;type:varchar(20);not null;default:pending;index" json:"status"`                            // pending, paid, cancelled
	Amount          float64    `gorm:"column:amount;type:decimal(10,2);not null;default:0.00" json:"amount"`
	PayStatus       string     `gorm:"column:pay_status;type:varchar(20);not null;default:unpaid;index" json:"payStatus"` // unpaid, paid, refunded
	TradeNo         string     `gorm:"column:trade_no;type:varchar(100);default:''" json:"tradeNo"`
	PaidAt          *time.Time `gorm:"column:paid_at" json:"paidAt,omitempty"`                                               // 支付时间
	SyncStatus      string     `gorm:"column:sync_status;type:varchar(20);not null;default:pending;index" json:"syncStatus"` // pending, synced, failed
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deletedAt,omitempty"`
}

// TableName 表名
func (m *Order) TableName() string {
	return "orders"
}

// Reward 奖励记录模型
type Reward struct {
	Id         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     int64      `gorm:"column:user_id;not null;index" json:"userId"`
	MemberID   *int64     `gorm:"column:member_id;index" json:"memberId"` // 关联会员ID（可选）
	OrderId    int64      `gorm:"column:order_id;not null;index" json:"orderId"`
	CampaignId int64      `gorm:"column:campaign_id;not null;index" json:"campaignId"`
	Amount     float64    `gorm:"column:amount;type:decimal(10,2);not null;default:0.00" json:"amount"`
	Status     string     `gorm:"column:status;type:varchar(20);not null;default:pending;index" json:"status"` // pending, settled, cancelled
	SettledAt  *time.Time `gorm:"column:settled_at" json:"settledAt,omitempty"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName 表名
func (m *Reward) TableName() string {
	return "rewards"
}

// UserBalance 用户余额模型
type UserBalance struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId      int64     `gorm:"column:user_id;not null;uniqueIndex" json:"userId"`
	Balance     float64   `gorm:"column:balance;type:decimal(10,2);not null;default:0.00" json:"balance"`
	TotalReward float64   `gorm:"column:total_reward;type:decimal(10,2);not null;default:0.00" json:"totalReward"`
	Version     int64     `gorm:"column:version;not null;default:0" json:"version"` // 乐观锁版本号
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName 表名
func (m *UserBalance) TableName() string {
	return "user_balances"
}

// SyncLog 同步日志模型
type SyncLog struct {
	Id         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderId    int64      `gorm:"column:order_id;not null;index" json:"orderId"`
	SyncType   string     `gorm:"column:sync_type;type:varchar(20);not null" json:"syncType"`                           // order, reward
	SyncStatus string     `gorm:"column:sync_status;type:varchar(20);not null;default:pending;index" json:"syncStatus"` // pending, synced, failed
	Attempts   int        `gorm:"column:attempts;not null;default:0" json:"attempts"`
	ErrorMsg   string     `gorm:"column:error_msg;type:text" json:"errorMsg"`
	SyncedAt   *time.Time `gorm:"column:synced_at" json:"syncedAt,omitempty"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName 表名
func (m *SyncLog) TableName() string {
	return "sync_logs"
}
