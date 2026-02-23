package model

import "time"

// Promoter 推广员表
type Promoter struct {
	Id             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId         int64      `gorm:"column:user_id;not null;index:idx_promoter_user;uniqueIndex:idx_user_brand" json:"userId"`
	BrandId        int64      `gorm:"column:brand_id;not null;index:idx_promoter_brand;uniqueIndex:idx_user_brand" json:"brandId"`
	Status         string     `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"` // active/inactive/blocked
	Level          string     `gorm:"column:level;type:varchar(50);default:''" json:"level"`                      // VIP/普通
	TotalOrders    int64      `gorm:"column:total_orders;not null;default:0" json:"totalOrders"`
	TotalRewards   float64    `gorm:"column:total_rewards;type:decimal(10,2);not null;default:0.00" json:"totalRewards"`
	ConversionRate float64    `gorm:"column:conversion_rate;type:decimal(5,2);not null;default:0.00" json:"conversionRate"`
	CampaignCount  int64      `gorm:"column:campaign_count;not null;default:0" json:"campaignCount"`
	LastActiveAt   *time.Time `gorm:"column:last_active_at" json:"lastActiveAt"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deletedAt,omitempty"`

	// 关联
	User  *User  `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Brand *Brand `gorm:"foreignKey:BrandId" json:"brand,omitempty"`
}

// TableName 表名
func (Promoter) TableName() string {
	return "promoters"
}

// PromoterLink 推广员推广链接表
type PromoterLink struct {
	Id         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PromoterId int64      `gorm:"column:promoter_id;not null;index" json:"promoterId"`
	CampaignId int64      `gorm:"column:campaign_id;not null;index" json:"campaignId"`
	LinkCode   string     `gorm:"column:link_code;type:varchar(50);not null;uniqueIndex" json:"linkCode"`
	ClickCount int64      `gorm:"column:click_count;not null;default:0" json:"clickCount"`
	OrderCount int64      `gorm:"column:order_count;not null;default:0" json:"orderCount"`
	Status     string     `gorm:"column:status;type:varchar(20);not null;default:active" json:"status"`
	ExpiresAt  *time.Time `gorm:"column:expires_at" json:"expiresAt"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`

	// 关联
	Promoter *Promoter `gorm:"foreignKey:PromoterId" json:"promoter,omitempty"`
	Campaign *Campaign `gorm:"foreignKey:CampaignId" json:"campaign,omitempty"`
}

// TableName 表名
func (PromoterLink) TableName() string {
	return "promoter_links"
}

// PromoterReward 推广员奖励记录表
type PromoterReward struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PromoterId  int64     `gorm:"column:promoter_id;not null;index" json:"promoterId"`
	Type        string    `gorm:"column:type;type:varchar(20);not null;index" json:"type"`                     // commission/bonus/penalty
	Status      string    `gorm:"column:status;type:varchar(20);not null;default:pending;index" json:"status"` // pending/paid/cancelled
	Amount      float64   `gorm:"column:amount;type:decimal(10,2);not null" json:"amount"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	CampaignId  *int64    `gorm:"column:campaign_id;index" json:"campaignId"`
	OrderId     *int64    `gorm:"column:order_id;index" json:"orderId"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`

	// 关联
	Promoter *Promoter `gorm:"foreignKey:PromoterId" json:"promoter,omitempty"`
	Campaign *Campaign `gorm:"foreignKey:CampaignId" json:"campaign,omitempty"`
	Order    *Order    `gorm:"foreignKey:OrderId" json:"order,omitempty"`
}

// TableName 表名
func (PromoterReward) TableName() string {
	return "promoter_rewards"
}
