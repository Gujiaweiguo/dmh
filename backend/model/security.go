package model

import "time"

// PasswordPolicy 密码策略配置
type PasswordPolicy struct {
	ID                    int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MinLength             int       `gorm:"column:min_length;not null;default:8" json:"minLength"`                          // 最小长度
	RequireUppercase      bool      `gorm:"column:require_uppercase;not null;default:true" json:"requireUppercase"`         // 需要大写字母
	RequireLowercase      bool      `gorm:"column:require_lowercase;not null;default:true" json:"requireLowercase"`         // 需要小写字母
	RequireNumbers        bool      `gorm:"column:require_numbers;not null;default:true" json:"requireNumbers"`             // 需要数字
	RequireSpecialChars   bool      `gorm:"column:require_special_chars;not null;default:true" json:"requireSpecialChars"`  // 需要特殊字符
	MaxAge                int       `gorm:"column:max_age;not null;default:90" json:"maxAge"`                               // 密码最大有效期（天）
	HistoryCount          int       `gorm:"column:history_count;not null;default:5" json:"historyCount"`                    // 历史密码记录数量
	MaxLoginAttempts      int       `gorm:"column:max_login_attempts;not null;default:5" json:"maxLoginAttempts"`           // 最大登录尝试次数
	LockoutDuration       int       `gorm:"column:lockout_duration;not null;default:30" json:"lockoutDuration"`             // 锁定时长（分钟）
	SessionTimeout        int       `gorm:"column:session_timeout;not null;default:480" json:"sessionTimeout"`              // 会话超时时间（分钟）
	MaxConcurrentSessions int       `gorm:"column:max_concurrent_sessions;not null;default:3" json:"maxConcurrentSessions"` // 最大并发会话数
	CreatedAt             time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt             time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (PasswordPolicy) TableName() string {
	return "password_policies"
}

// PasswordHistory 密码历史记录
type PasswordHistory struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"column:user_id;not null;index" json:"userId"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (PasswordHistory) TableName() string {
	return "password_histories"
}

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     *int64    `gorm:"column:user_id;index" json:"userId"`                               // 用户ID，可为空（未知用户）
	Username   string    `gorm:"column:username;type:varchar(50);not null;index" json:"username"`  // 尝试登录的用户名
	ClientIP   string    `gorm:"column:client_ip;type:varchar(45);not null;index" json:"clientIP"` // 客户端IP
	UserAgent  string    `gorm:"column:user_agent;type:varchar(500)" json:"userAgent"`             // 用户代理
	Success    bool      `gorm:"column:success;not null;index" json:"success"`                     // 是否成功
	FailReason string    `gorm:"column:fail_reason;type:varchar(200)" json:"failReason"`           // 失败原因
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}

// UserSession 用户会话记录
type UserSession struct {
	ID           string    `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`                                    // 会话ID
	UserID       int64     `gorm:"column:user_id;not null;index" json:"userId"`                                        // 用户ID
	ClientIP     string    `gorm:"column:client_ip;type:varchar(45);not null" json:"clientIP"`                         // 客户端IP
	UserAgent    string    `gorm:"column:user_agent;type:varchar(500)" json:"userAgent"`                               // 用户代理
	LoginAt      time.Time `gorm:"column:login_at;not null;default:CURRENT_TIMESTAMP" json:"loginAt"`                  // 登录时间
	LastActiveAt time.Time `gorm:"column:last_active_at;not null;default:CURRENT_TIMESTAMP;index" json:"lastActiveAt"` // 最后活跃时间
	ExpiresAt    time.Time `gorm:"column:expires_at;not null;index" json:"expiresAt"`                                  // 过期时间
	Status       string    `gorm:"column:status;type:varchar(20);not null;default:active;index" json:"status"`         // 状态: active/expired/revoked
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

// AuditLog 操作审计日志
type AuditLog struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     *int64    `gorm:"column:user_id;index" json:"userId"`                                          // 操作用户ID
	Username   string    `gorm:"column:username;type:varchar(50);index" json:"username"`                      // 操作用户名
	Action     string    `gorm:"column:action;type:varchar(100);not null;index" json:"action"`                // 操作类型
	Resource   string    `gorm:"column:resource;type:varchar(100);index" json:"resource"`                     // 操作资源
	ResourceID string    `gorm:"column:resource_id;type:varchar(100);index" json:"resourceId"`                // 资源ID
	Details    string    `gorm:"column:details;type:text" json:"details"`                                     // 操作详情
	ClientIP   string    `gorm:"column:client_ip;type:varchar(45);index" json:"clientIP"`                     // 客户端IP
	UserAgent  string    `gorm:"column:user_agent;type:varchar(500)" json:"userAgent"`                        // 用户代理
	Status     string    `gorm:"column:status;type:varchar(20);not null;default:success;index" json:"status"` // 操作状态: success/failed
	ErrorMsg   string    `gorm:"column:error_msg;type:varchar(500)" json:"errorMsg"`                          // 错误信息
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// SecurityEvent 安全事件记录
type SecurityEvent struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	EventType   string     `gorm:"column:event_type;type:varchar(50);not null;index" json:"eventType"` // 事件类型
	Severity    string     `gorm:"column:severity;type:varchar(20);not null;index" json:"severity"`    // 严重程度: low/medium/high/critical
	UserID      *int64     `gorm:"column:user_id;index" json:"userId"`                                 // 相关用户ID
	Username    string     `gorm:"column:username;type:varchar(50);index" json:"username"`             // 相关用户名
	ClientIP    string     `gorm:"column:client_ip;type:varchar(45);not null;index" json:"clientIP"`   // 客户端IP
	UserAgent   string     `gorm:"column:user_agent;type:varchar(500)" json:"userAgent"`               // 用户代理
	Description string     `gorm:"column:description;type:text;not null" json:"description"`           // 事件描述
	Details     string     `gorm:"column:details;type:json" json:"details"`                            // 事件详情（JSON格式）
	Handled     bool       `gorm:"column:handled;not null;default:false;index" json:"handled"`         // 是否已处理
	HandledBy   *int64     `gorm:"column:handled_by" json:"handledBy"`                                 // 处理人ID
	HandledAt   *time.Time `gorm:"column:handled_at" json:"handledAt"`                                 // 处理时间
	CreatedAt   time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`
}

func (SecurityEvent) TableName() string {
	return "security_events"
}
