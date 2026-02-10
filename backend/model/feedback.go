package model

import "time"

// UserFeedback 用户反馈
type UserFeedback struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID         int64      `gorm:"column:user_id;not null;index:idx_user_created" json:"userId"`
	Category       string     `gorm:"column:category;type:varchar(50);not null;index:idx_category_created" json:"category"` // poster, payment, verification, other
	Subcategory    string     `gorm:"column:subcategory;type:varchar(100)" json:"subcategory"`
	Rating         *int       `gorm:"column:rating" json:"rating"` // 1-5星
	Title          string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Content        string     `gorm:"column:content;type:text;not null" json:"content"`
	FeatureUseCase string     `gorm:"column:feature_use_case;type:text" json:"featureUseCase"`
	DeviceInfo     string     `gorm:"column:device_info;type:varchar(500)" json:"deviceInfo"`
	BrowserInfo    string     `gorm:"column:browser_info;type:varchar(500)" json:"browserInfo"`
	Priority       string     `gorm:"column:priority;type:varchar(20);default:medium;index:idx_priority_created" json:"priority"`     // low, medium, high
	Status         string     `gorm:"column:status;type:varchar(20);not null;default:pending;index:idx_status_created" json:"status"` // pending, reviewing, resolved, closed
	AssigneeID     *int64     `gorm:"column:assignee_id" json:"assigneeId"`
	Response       string     `gorm:"column:response;type:text" json:"response"`
	ResolvedAt     *time.Time `gorm:"column:resolved_at" json:"resolvedAt"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index:idx_user_created,idx_status_created,idx_category_created,idx_priority_created" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// 关联
	User     *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Assignee *User         `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Tags     []FeedbackTag `gorm:"many2many:feedback_tag_relations;" json:"tags,omitempty"`
}

func (UserFeedback) TableName() string {
	return "user_feedback"
}

// FeatureSatisfactionSurvey 功能满意度调查
type FeatureSatisfactionSurvey struct {
	ID                     int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID                 int64     `gorm:"column:user_id;not null;index" json:"userId"`
	UserRole               string    `gorm:"column:user_role;type:varchar(50);not null" json:"userRole"`
	Feature                string    `gorm:"column:feature;type:varchar(100);not null;index" json:"feature"` // poster, payment, verification
	EaseOfUse              *int      `gorm:"column:ease_of_use" json:"easeOfUse"`                            // 1-5
	Performance            *int      `gorm:"column:performance" json:"performance"`                          // 1-5
	Reliability            *int      `gorm:"column:reliability" json:"reliability"`                          // 1-5
	OverallSatisfaction    *int      `gorm:"column:overall_satisfaction" json:"overallSatisfaction"`         // 1-5
	WouldRecommend         *int      `gorm:"column:would_recommend" json:"wouldRecommend"`                   // 1-5
	MostLiked              string    `gorm:"column:most_liked;type:text" json:"mostLiked"`
	LeastLiked             string    `gorm:"column:least_liked;type:text" json:"leastLiked"`
	ImprovementSuggestions string    `gorm:"column:improvement_suggestions;type:text" json:"improvementSuggestions"`
	WouldLikeMoreFeatures  string    `gorm:"column:would_like_more_features;type:text" json:"wouldLikeMoreFeatures"`
	CreatedAt              time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (FeatureSatisfactionSurvey) TableName() string {
	return "feature_satisfaction_surveys"
}

// FAQItem 常见问题
type FAQItem struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Category        string    `gorm:"column:category;type:varchar(100);not null;index" json:"category"`
	Question        string    `gorm:"column:question;type:text;not null" json:"question"`
	Answer          string    `gorm:"column:answer;type:text;not null" json:"answer"`
	SortOrder       int       `gorm:"column:sort_order;default:0;index" json:"sortOrder"`
	IsPublished     bool      `gorm:"column:is_published;not null;default:true;index" json:"isPublished"`
	ViewCount       int       `gorm:"column:view_count;default:0" json:"viewCount"`
	HelpfulCount    int       `gorm:"column:helpful_count;default:0" json:"helpfulCount"`
	NotHelpfulCount int       `gorm:"column:not_helpful_count;default:0" json:"notHelpfulCount"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (FAQItem) TableName() string {
	return "faq_items"
}

// FeatureUsageStat 功能使用统计
type FeatureUsageStat struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"column:user_id;not null;index" json:"userId"`
	UserRole     string    `gorm:"column:user_role;type:varchar(50);not null" json:"userRole"`
	Feature      string    `gorm:"column:feature;type:varchar(100);not null;index" json:"feature"` // poster, payment, verification
	Action       string    `gorm:"column:action;type:varchar(100);not null;index" json:"action"`   // generate, refresh, pay, verify, etc.
	CampaignID   *int64    `gorm:"column:campaign_id;index" json:"campaignId"`
	Success      bool      `gorm:"column:success;not null;default:true;index" json:"success"`
	DurationMs   *int      `gorm:"column:duration_ms" json:"durationMs"`
	ErrorMessage string    `gorm:"column:error_message;type:text" json:"errorMessage"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`

	// 关联
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Campaign *Campaign `gorm:"foreignKey:CampaignID" json:"campaign,omitempty"`
}

func (FeatureUsageStat) TableName() string {
	return "feature_usage_stats"
}

// FeedbackTag 反馈标签
type FeedbackTag struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(50);not null;uniqueIndex" json:"name"`
	Color     string    `gorm:"column:color;type:varchar(20);default:#1890ff" json:"color"` // 标签颜色
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

func (FeedbackTag) TableName() string {
	return "feedback_tags"
}

// FeedbackTagRelation 反馈标签关联表
type FeedbackTagRelation struct {
	FeedbackID int64 `gorm:"column:feedback_id;not null;primaryKey" json:"feedbackId"`
	TagID      int64 `gorm:"column:tag_id;not null;primaryKey" json:"tagId"`

	// 关联
	Feedback *UserFeedback `gorm:"foreignKey:FeedbackID" json:"-"`
	Tag      *FeedbackTag  `gorm:"foreignKey:TagID" json:"-"`
}

func (FeedbackTagRelation) TableName() string {
	return "feedback_tag_relations"
}
