package model

import "time"

type Material struct {
	Id          int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	Type        string     `gorm:"column:type;type:varchar(20);not null;default:image" json:"type"` // image/text
	Url         string     `gorm:"column:url;type:varchar(500)" json:"url"`
	Content     string     `gorm:"column:content;type:text" json:"content"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deletedAt,omitempty"`
}

func (Material) TableName() string {
	return "materials"
}
