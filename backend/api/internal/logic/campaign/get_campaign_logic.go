// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"encoding/json"
	"strconv"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCampaignLogic {
	return &GetCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCampaignLogic) GetCampaign(id string) (resp *types.CampaignResp, err error) {
	campaignId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	var campaign model.Campaign
	if err := l.svcCtx.DB.First(&campaign, campaignId).Error; err != nil {
		return nil, err
	}

	// 解析表单字段
	var formFields []types.FormField
	if campaign.FormFields != "" {
		json.Unmarshal([]byte(campaign.FormFields), &formFields)
	}

	// 解析分销奖励配置
	var distributionRewards []types.DistributorLevelRewardResp
	if campaign.DistributionRewards != nil {
		json.Unmarshal([]byte(*campaign.DistributionRewards), &distributionRewards)
	}

	// 获取品牌名称
	var brand model.Brand
	l.svcCtx.DB.First(&brand, campaign.BrandId)

	return &types.CampaignResp{
		Id:                  campaign.Id,
		BrandId:             campaign.BrandId,
		BrandName:           brand.Name,
		Name:                campaign.Name,
		Description:         campaign.Description,
		FormFields:          formFields,
		RewardRule:          campaign.RewardRule,
		StartTime:           campaign.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:             campaign.EndTime.Format("2006-01-02 15:04:05"),
		Status:              campaign.Status,
		EnableDistribution:  campaign.EnableDistribution,
		DistributionLevel:   campaign.DistributionLevel,
		DistributionRewards: distributionRewards,
		CreatedAt:           campaign.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
