package model

import "time"

// Member 会员表（平台唯一，以 unionid 为标识）
type Member struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UnionID   string     `gorm:"column:unionid;type:varchar(100);not null;uniqueIndex" json:"unionid"` // 微信 unionid，平台唯一
	Nickname  string     `gorm:"column:nickname;type:varchar(100)" json:"nickname"`
	Avatar    string     `gorm:"column:avatar;type:varchar(500)" json:"avatar"`
	Phone     string     `gorm:"column:phone;type:varchar(20);index" json:"phone"`
	Gender    int        `gorm:"column:gender;default:0" json:"gender"`                                      // 0:未知 1:男 2:女
	Source    string     `gorm:"column:source;type:varchar(50);index" json:"source"`                         // 首次来源渠道
	Status    string     `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"` // active/disabled
	CreatedAt time.Time  `gorm:"column:created_at;not null;autoCreateTime;index" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

func (Member) TableName() string {
	return "members"
}

// MemberProfile 会员画像扩展表
type MemberProfile struct {
	ID                    int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MemberID              int64      `gorm:"column:member_id;not null;uniqueIndex" json:"memberId"`
	TotalOrders           int        `gorm:"column:total_orders;default:0" json:"totalOrders"`                         // 累计订单数
	TotalPayment          float64    `gorm:"column:total_payment;type:decimal(10,2);default:0.00" json:"totalPayment"` // 累计支付金额
	TotalReward           float64    `gorm:"column:total_reward;type:decimal(10,2);default:0.00" json:"totalReward"`   // 累计奖励金额
	FirstOrderAt          *time.Time `gorm:"column:first_order_at" json:"firstOrderAt"`
	LastOrderAt           *time.Time `gorm:"column:last_order_at" json:"lastOrderAt"`
	FirstPaymentAt        *time.Time `gorm:"column:first_payment_at" json:"firstPaymentAt"`
	LastPaymentAt         *time.Time `gorm:"column:last_payment_at" json:"lastPaymentAt"`
	ParticipatedCampaigns int        `gorm:"column:participated_campaigns;default:0" json:"participatedCampaigns"` // 参与活动数
	CreatedAt             time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt             time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (MemberProfile) TableName() string {
	return "member_profiles"
}

// MemberTag 会员标签表
type MemberTag struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(50);not null;uniqueIndex" json:"name"`
	Category    string    `gorm:"column:category;type:varchar(50);index" json:"category"` // 标签分类
	Color       string    `gorm:"column:color;type:varchar(20)" json:"color"`
	Description string    `gorm:"column:description;type:varchar(200)" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (MemberTag) TableName() string {
	return "member_tags"
}

// MemberTagLink 会员标签关联表
type MemberTagLink struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MemberID  int64     `gorm:"column:member_id;not null;uniqueIndex:uk_member_tag;index" json:"memberId"`
	TagID     int64     `gorm:"column:tag_id;not null;uniqueIndex:uk_member_tag;index" json:"tagId"`
	CreatedBy int64     `gorm:"column:created_by;not null" json:"createdBy"` // 操作人 user_id
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
}

func (MemberTagLink) TableName() string {
	return "member_tag_links"
}

// MemberBrandLink 会员品牌关联表
type MemberBrandLink struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MemberID        int64     `gorm:"column:member_id;not null;uniqueIndex:uk_member_brand;index" json:"memberId"`
	BrandID         int64     `gorm:"column:brand_id;not null;uniqueIndex:uk_member_brand;index" json:"brandId"`
	FirstCampaignID int64     `gorm:"column:first_campaign_id;not null" json:"firstCampaignId"` // 首次参与的活动ID
	CreatedAt       time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
}

func (MemberBrandLink) TableName() string {
	return "member_brand_links"
}

// MemberMergeRequest 会员合并请求表
type MemberMergeRequest struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SourceMemberID int64      `gorm:"column:source_member_id;not null;index" json:"sourceMemberId"`                // 被合并的会员
	TargetMemberID int64      `gorm:"column:target_member_id;not null;index" json:"targetMemberId"`                // 主会员（保留）
	Status         string     `gorm:"column:status;type:varchar(20);not null;default:pending;index" json:"status"` // pending/completed/failed
	Reason         string     `gorm:"column:reason;type:text" json:"reason"`
	ConflictInfo   string     `gorm:"column:conflict_info;type:json" json:"conflictInfo"` // 冲突信息
	CreatedBy      int64      `gorm:"column:created_by;not null" json:"createdBy"`        // 操作人 user_id
	ExecutedAt     *time.Time `gorm:"column:executed_at" json:"executedAt"`
	ErrorMsg       string     `gorm:"column:error_msg;type:text" json:"errorMsg"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (MemberMergeRequest) TableName() string {
	return "member_merge_requests"
}

// ExportRequest 导出申请表
type ExportRequest struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BrandID      int64      `gorm:"column:brand_id;not null;index" json:"brandId"`                               // 申请导出的品牌
	RequestedBy  int64      `gorm:"column:requested_by;not null;index" json:"requestedBy"`                       // 申请人 user_id
	Reason       string     `gorm:"column:reason;type:text;not null" json:"reason"`                              // 导出原因
	Filters      string     `gorm:"column:filters;type:json" json:"filters"`                                     // 筛选条件
	Status       string     `gorm:"column:status;type:varchar(20);not null;default:pending;index" json:"status"` // pending/approved/rejected/completed
	ApprovedBy   *int64     `gorm:"column:approved_by" json:"approvedBy"`                                        // 审批人 user_id
	ApprovedAt   *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	RejectReason string     `gorm:"column:reject_reason;type:text" json:"rejectReason"`
	FileUrl      string     `gorm:"column:file_url;type:varchar(500)" json:"fileUrl"` // 导出文件URL
	RecordCount  int        `gorm:"column:record_count;default:0" json:"recordCount"` // 导出记录数
	CreatedAt    time.Time  `gorm:"column:created_at;not null;autoCreateTime;index" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (ExportRequest) TableName() string {
	return "export_requests"
}
