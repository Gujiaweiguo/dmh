package model

import "time"

// PosterTemplate 海报模板模型
type PosterTemplate struct {
	Id            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Type          string    `gorm:"column:type;type:varchar(20);not null;index" json:"type"` // campaign/distributor
	CampaignId    *int64    `gorm:"column:campaign_id;index" json:"campaignId,omitempty"`
	DistributorId *int64    `gorm:"column:distributor_id;index" json:"distributorId,omitempty"`
	TemplateUrl   string    `gorm:"column:template_url;type:varchar(500);not null" json:"templateUrl"`
	PosterData    string    `gorm:"column:poster_data;type:json" json:"posterData"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`

	// 关联
	Campaign    *Campaign    `gorm:"foreignKey:CampaignId" json:"campaign,omitempty"`
	Distributor *Distributor `gorm:"foreignKey:DistributorId" json:"distributor,omitempty"`
}

// TableName 表名
func (PosterTemplate) TableName() string {
	return "poster_templates"
}
