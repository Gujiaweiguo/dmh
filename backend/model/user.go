package model

import "time"

// User 用户表
type User struct {
	Id            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username      string     `gorm:"column:username;type:varchar(50);not null;uniqueIndex" json:"username"`
	Password      string     `gorm:"column:password;type:varchar(255);not null" json:"-"` // 密码不返回给前端
	Phone         string     `gorm:"column:phone;type:varchar(20);not null;uniqueIndex" json:"phone"`
	Email         string     `gorm:"column:email;type:varchar(100)" json:"email"`
	Avatar        string     `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	RealName      string     `gorm:"column:real_name;type:varchar(50)" json:"realName"`
	Role          string     `gorm:"column:role;type:varchar(50);not null;default:participant;index" json:"role"` // platform_admin/participant
	Status        string     `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"`  // active/disabled/locked
	LoginAttempts int        `gorm:"column:login_attempts;default:0" json:"loginAttempts"`
	LockedUntil   *time.Time `gorm:"column:locked_until" json:"lockedUntil"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

// Role 角色表
type Role struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Code        string    `gorm:"column:code;type:varchar(50);not null;uniqueIndex" json:"code"` // platform_admin/participant/anonymous
	Description string    `gorm:"column:description;type:varchar(200)" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (Role) TableName() string {
	return "roles"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null;uniqueIndex:uk_user_role" json:"userId"`
	RoleID    int64     `gorm:"column:role_id;not null;uniqueIndex:uk_user_role;index" json:"roleId"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// Permission 权限表
type Permission struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Code        string    `gorm:"column:code;type:varchar(100);not null;uniqueIndex" json:"code"`
	Resource    string    `gorm:"column:resource;type:varchar(100);not null" json:"resource"` // campaign/order/user/brand
	Action      string    `gorm:"column:action;type:varchar(50);not null" json:"action"`      // create/read/update/delete
	Description string    `gorm:"column:description;type:varchar(200)" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID       int64     `gorm:"column:role_id;not null;uniqueIndex:uk_role_permission;index" json:"roleId"`
	PermissionID int64     `gorm:"column:permission_id;not null;uniqueIndex:uk_role_permission;index" json:"permissionId"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// Brand 品牌表
type Brand struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Logo        string    `gorm:"column:logo;type:varchar(255)" json:"logo"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Status      string    `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"` // active/disabled
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (Brand) TableName() string {
	return "brands"
}

// UserBrand 用户品牌关联表 (替代原来的BrandAdmin)
type UserBrand struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"column:user_id;not null;uniqueIndex:uk_user_brand;index" json:"userId"`
	BrandId   int64     `gorm:"column:brand_id;not null;uniqueIndex:uk_user_brand;index" json:"brandId"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (UserBrand) TableName() string {
	return "user_brands"
}

// Withdrawal 提现申请表
type Withdrawal struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID         int64      `gorm:"column:user_id;not null;index" json:"userId"`
	BrandId        int64      `gorm:"column:brand_id;not null;index" json:"brandId"`
	DistributorId  int64      `gorm:"column:distributor_id;not null;index" json:"distributorId"`
	Amount         float64    `gorm:"column:amount;type:decimal(10,2);not null" json:"amount"`
	Status         string     `gorm:"column:status;type:varchar(20);not null;default:pending;index" json:"status"` // pending/approved/rejected/processing/completed/failed
	PayType        string     `gorm:"column:pay_type;not null" json:"payType"`                                     // wechat/alipay/bank
	PayAccount     string     `gorm:"column:pay_account;not null" json:"payAccount"`
	PayRealName    string     `gorm:"column:pay_real_name" json:"payRealName"`
	ApprovedBy     *int64     `gorm:"column:approved_by" json:"approvedBy"`
	ApprovedAt     *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	ApprovedNotes  string     `gorm:"column:approved_notes;type:text" json:"approvedNotes"`
	RejectedReason string     `gorm:"column:rejected_reason;type:text" json:"rejectedReason"`
	PaidAt         *time.Time `gorm:"column:paid_at" json:"paidAt"`
	TradeNo        string     `gorm:"column:trade_no;type:varchar(100)" json:"tradeNo"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// 兼容旧字段（保留用于向后兼容）
	BankName    string `gorm:"column:bank_name;type:varchar(100)" json:"bankName"`
	BankAccount string `gorm:"column:bank_account;type:varchar(50)" json:"bankAccount"`
	AccountName string `gorm:"column:account_name;type:varchar(50)" json:"accountName"`
	Remark      string `gorm:"column:remark;type:text" json:"remark"`
}

func (Withdrawal) TableName() string {
	return "withdrawals"
}

// Menu 菜单表
type Menu struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Code      string    `gorm:"column:code;type:varchar(100);not null;uniqueIndex" json:"code"`
	Path      string    `gorm:"column:path;type:varchar(200)" json:"path"`
	Icon      string    `gorm:"column:icon;type:varchar(100)" json:"icon"`
	ParentID  *int64    `gorm:"column:parent_id;index" json:"parentId"`
	Sort      int       `gorm:"column:sort;default:0" json:"sort"`
	Type      string    `gorm:"column:type;type:varchar(20);not null;default:menu" json:"type"`             // menu/button
	Platform  string    `gorm:"column:platform;type:varchar(20);not null;default:admin" json:"platform"`    // admin/h5
	Status    string    `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"` // active/disabled
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (Menu) TableName() string {
	return "menus"
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID    int64     `gorm:"column:role_id;not null;uniqueIndex:uk_role_menu;index" json:"roleId"`
	MenuID    int64     `gorm:"column:menu_id;not null;uniqueIndex:uk_role_menu;index" json:"menuId"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (RoleMenu) TableName() string {
	return "role_menus"
}

// BrandAsset 品牌素材表
type BrandAsset struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BrandID     int64     `gorm:"column:brand_id;not null;index" json:"brandId"`
	Name        string    `gorm:"column:name;type:varchar(200);not null" json:"name"`
	Type        string    `gorm:"column:type;type:varchar(50);not null;index" json:"type"` // image/video/document
	Category    string    `gorm:"column:category;type:varchar(100);index" json:"category"`
	Tags        string    `gorm:"column:tags;type:varchar(500)" json:"tags"`
	FileUrl     string    `gorm:"column:file_url;type:varchar(500);not null" json:"fileUrl"`
	FileSize    int64     `gorm:"column:file_size;default:0" json:"fileSize"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Status      string    `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"` // active/disabled
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (BrandAsset) TableName() string {
	return "brand_assets"
}
