package model

import "time"

// Distributor 分销商表
type Distributor struct {
	Id                int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId            int64      `gorm:"column:user_id;not null;index:idx_distributor_user;uniqueIndex:idx_user_brand" json:"userId"`
	BrandId           int64      `gorm:"column:brand_id;not null;index:idx_distributor_brand;uniqueIndex:idx_user_brand" json:"brandId"`
	Level             int        `gorm:"column:level;not null;default:1;index" json:"level"`                          // 1/2/3 级别
	ParentId          *int64     `gorm:"column:parent_id;index" json:"parentId"`                                      // 上级分销商ID
	Status            string     `gorm:"column:status;type:varchar(20);not null;default:pending;index" json:"status"` // pending/active/suspended
	ApprovedBy        *int64     `gorm:"column:approved_by" json:"approvedBy"`
	ApprovedAt        *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	TotalEarnings     float64    `gorm:"column:total_earnings;type:decimal(10,2);not null;default:0.00" json:"totalEarnings"`
	SubordinatesCount int        `gorm:"column:subordinates_count;not null;default:0" json:"subordinatesCount"`
	CreatedAt         time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;index" json:"updatedAt"`
	DeletedAt         *time.Time `gorm:"column:deleted_at" json:"deletedAt,omitempty"`

	// 关联
	User         *User         `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Brand        *Brand        `gorm:"foreignKey:BrandId" json:"brand,omitempty"`
	Parent       *Distributor  `gorm:"foreignKey:ParentId" json:"parent,omitempty"`
	Subordinates []Distributor `gorm:"foreignKey:ParentId" json:"subordinates,omitempty"`
}

// TableName 表名
func (Distributor) TableName() string {
	return "distributors"
}

// DistributorApplication 分销商申请表
type DistributorApplication struct {
	Id          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId      int64      `gorm:"column:user_id;not null;index" json:"userId"`
	BrandId     int64      `gorm:"column:brand_id;not null;index" json:"brandId"`
	Status      string     `gorm:"column:status;type:varchar(20);not null;default:pending;index" json:"status"` // pending/approved/rejected
	Reason      string     `gorm:"column:reason;type:text" json:"reason"`
	ReviewedBy  *int64     `gorm:"column:reviewed_by" json:"reviewedBy"`
	ReviewedAt  *time.Time `gorm:"column:reviewed_at" json:"reviewedAt"`
	ReviewNotes string     `gorm:"column:review_notes;type:text" json:"reviewNotes"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// 关联
	User     *User  `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Brand    *Brand `gorm:"foreignKey:BrandId" json:"brand,omitempty"`
	Reviewer *User  `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

// TableName 表名
func (DistributorApplication) TableName() string {
	return "distributor_applications"
}

// DistributorLevelReward 分销商级别奖励配置表
type DistributorLevelReward struct {
	Id               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BrandId          int64     `gorm:"column:brand_id;not null;index:idx_level_reward_brand;uniqueIndex:idx_brand_level" json:"brandId"`
	Level            int       `gorm:"column:level;not null;uniqueIndex:idx_brand_level" json:"level"`                           // 1/2/3
	RewardPercentage float64   `gorm:"column:reward_percentage;type:decimal(5,2);not null;default:0.00" json:"rewardPercentage"` // 奖励百分比
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// 关联
	Brand *Brand `gorm:"foreignKey:BrandId" json:"brand,omitempty"`
}

// TableName 表名
func (DistributorLevelReward) TableName() string {
	return "distributor_level_rewards"
}

// DistributorReward 分销商奖励记录表（扩展原有奖励表）
type DistributorReward struct {
	Id            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DistributorId int64      `gorm:"column:distributor_id;not null;index" json:"distributorId"`
	UserId        int64      `gorm:"column:user_id;not null;index" json:"userId"`
	OrderId       int64      `gorm:"column:order_id;not null;index" json:"orderId"`
	CampaignId    int64      `gorm:"column:campaign_id;not null;index" json:"campaignId"`
	Amount        float64    `gorm:"column:amount;type:decimal(10,2);not null;default:0.00" json:"amount"`
	Level         int        `gorm:"column:level;not null" json:"level"`                              // 奖励级别 1/2/3
	RewardRate    float64    `gorm:"column:reward_rate;type:decimal(5,2);not null" json:"rewardRate"` // 奖励比例
	FromUserId    *int64     `gorm:"column:from_user_id" json:"fromUserId"`                           // 购买用户ID
	Status        string     `gorm:"column:status;type:varchar(20);not null;default:settled;index" json:"status"`
	SettledAt     *time.Time `gorm:"column:settled_at" json:"settledAt"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// 关联
	Distributor *Distributor `gorm:"foreignKey:DistributorId" json:"distributor,omitempty"`
	Order       *Order       `gorm:"foreignKey:OrderId" json:"order,omitempty"`
}

// TableName 表名
func (DistributorReward) TableName() string {
	return "distributor_rewards"
}

// DistributorLink 分销商推广链接表
type DistributorLink struct {
	Id            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DistributorId int64      `gorm:"column:distributor_id;not null;index" json:"distributorId"`
	CampaignId    int64      `gorm:"column:campaign_id;not null;index" json:"campaignId"`
	LinkCode      string     `gorm:"column:link_code;type:varchar(50);not null;uniqueIndex" json:"linkCode"` // 推广码
	ClickCount    int        `gorm:"column:click_count;not null;default:0" json:"clickCount"`
	OrderCount    int        `gorm:"column:order_count;not null;default:0" json:"orderCount"`
	Status        string     `gorm:"column:status;type:varchar(20);not null;default:active" json:"status"`
	ExpiresAt     *time.Time `gorm:"column:expires_at" json:"expiresAt"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// 关联
	Distributor *Distributor `gorm:"foreignKey:DistributorId" json:"distributor,omitempty"`
	Campaign    *Campaign    `gorm:"foreignKey:CampaignId" json:"campaign,omitempty"`
}

// TableName 表名
func (DistributorLink) TableName() string {
	return "distributor_links"
}

// DistributorStatistics 分销商统计数据结构
type DistributorStatistics struct {
	DistributorId     int64   `json:"distributorId"`
	TotalOrders       int64   `json:"totalOrders"`
	TotalEarnings     float64 `json:"totalEarnings"`
	TodayEarnings     float64 `json:"todayEarnings"`
	MonthEarnings     float64 `json:"monthEarnings"`
	SubordinatesCount int     `json:"subordinatesCount"`
	ClickCount        int     `json:"clickCount"`
	ConversionRate    float64 `json:"conversionRate"`
}
