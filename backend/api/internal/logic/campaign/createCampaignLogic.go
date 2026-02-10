// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCampaignLogic {
	return &CreateCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCampaignLogic) CreateCampaign(req *types.CreateCampaignReq) (resp *types.CampaignResp, err error) {
	var startTime, endTime time.Time
	var err1, err2 error

	startTime, err1 = time.Parse(time.RFC3339, req.StartTime)
	endTime, err2 = time.Parse(time.RFC3339, req.EndTime)

	if err1 != nil || err2 != nil {
		startTime, err1 = time.Parse("2006-01-02T15:04:05", req.StartTime)
		endTime, err2 = time.Parse("2006-01-02T15:04:05", req.EndTime)

		if err1 != nil || err2 != nil {
			startTime, err1 = time.Parse("2006-01-02", req.StartTime)
			endTime, err2 = time.Parse("2006-01-02", req.EndTime)
		}
	}
	if err1 != nil || err2 != nil {
		l.Errorf("Time format error: startTime=%v, endTime=%v", err1, err2)
		return nil, fmt.Errorf("Time format error")
	}

	level := req.DistributionLevel
	if level == 0 {
		level = 1
	}
	if level < 1 || level > 3 {
		return nil, fmt.Errorf("Distribution level must be between 1 and 3")
	}

	posterTemplateId := req.PosterTemplateId
	if posterTemplateId == 0 {
		posterTemplateId = 1
	}

	var distributionRewards *string
	if req.DistributionRewards != "" {
		distributionRewards = &req.DistributionRewards
	}

	var paymentConfig *string
	if req.PaymentConfig != "" {
		paymentConfig = &req.PaymentConfig
	}

	newCampaign := model.Campaign{
		BrandId:             req.BrandId,
		Name:                req.Name,
		Description:         req.Description,
		RewardRule:          req.RewardRule,
		StartTime:           startTime,
		EndTime:             endTime,
		Status:              "active",
		EnableDistribution:  req.EnableDistribution,
		DistributionLevel:   level,
		DistributionRewards: distributionRewards,
		PaymentConfig:       paymentConfig,
		PosterTemplateId:    posterTemplateId,
	}

	// 序列化formFields数组为JSON字符串存储到数据库
	if len(req.FormFields) > 0 {
		formFieldsJSON, err := json.Marshal(req.FormFields)
		if err == nil {
			newCampaign.FormFields = string(formFieldsJSON)
			l.Infof("FormFields JSON: %s", newCampaign.FormFields)
		}
	}

	if err := l.svcCtx.DB.Create(&newCampaign).Error; err != nil {
		l.Errorf("Failed to create campaign: %v", err)
		return nil, fmt.Errorf("Failed to create campaign: %w", err)
	}

	l.Infof("Campaign created successfully: campaignId=%d, name=%s", newCampaign.Id, newCampaign.Name)

	var formFieldsStr string
	if newCampaign.FormFields != "" {
		formFieldsStr = newCampaign.FormFields
	} else {
		formFieldsStr = "[]"
	}

	distributionRewardsResp := ""
	if newCampaign.DistributionRewards != nil {
		distributionRewardsResp = *newCampaign.DistributionRewards
	}
	paymentConfigResp := ""
	if newCampaign.PaymentConfig != nil {
		paymentConfigResp = *newCampaign.PaymentConfig
	}

	resp = &types.CampaignResp{
		Id:                  newCampaign.Id,
		BrandId:             newCampaign.BrandId,
		Name:                newCampaign.Name,
		Description:         newCampaign.Description,
		FormFields:          formFieldsStr,
		RewardRule:          newCampaign.RewardRule,
		StartTime:           newCampaign.StartTime.Format("2006-01-02T15:04:05"),
		EndTime:             newCampaign.EndTime.Format("2006-01-02T15:04:05"),
		Status:              newCampaign.Status,
		EnableDistribution:  newCampaign.EnableDistribution,
		DistributionLevel:   newCampaign.DistributionLevel,
		DistributionRewards: distributionRewardsResp,
		PaymentConfig:       paymentConfigResp,
		PosterTemplateId:    newCampaign.PosterTemplateId,
		CreatedAt:           newCampaign.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt:           newCampaign.UpdatedAt.Format("2006-01-02T15:04:05"),
	}

	return resp, nil
}
