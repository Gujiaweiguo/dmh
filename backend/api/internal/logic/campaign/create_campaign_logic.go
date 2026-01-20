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
	// 验证FormFields结构
	for _, field := range req.FormFields {
		if field.Type == "" || field.Name == "" || field.Label == "" {
			return nil, fmt.Errorf("FormField必须包含type、name和label字段")
		}

		// 验证字段类型
		if field.Type != "text" && field.Type != "phone" && field.Type != "select" {
			return nil, fmt.Errorf("不支持的字段类型: %s", field.Type)
		}

		// select类型必须有选项
		if field.Type == "select" && len(field.Options) == 0 {
			return nil, fmt.Errorf("select类型字段必须提供选项列表")
		}
	}

	// 验证分销规则配置
	if req.EnableDistribution {
		if req.DistributionLevel < 1 || req.DistributionLevel > 3 {
			return nil, fmt.Errorf("分销层级必须在1-3之间")
		}
		if len(req.DistributionRewards) != req.DistributionLevel {
			return nil, fmt.Errorf("分销奖励配置数量必须与分销层级一致")
		}
	}

	// 序列化FormFields为JSON
	formFieldsJSON, err := json.Marshal(req.FormFields)
	if err != nil {
		return nil, fmt.Errorf("序列化FormFields失败: %w", err)
	}

	// 序列化分销奖励配置为JSON
	var distributionRewardsJSON *string
	if req.EnableDistribution && len(req.DistributionRewards) > 0 {
		rewardsJSON, err := json.Marshal(req.DistributionRewards)
		if err != nil {
			return nil, fmt.Errorf("序列化分销奖励配置失败: %w", err)
		}
		jsonStr := string(rewardsJSON)
		distributionRewardsJSON = &jsonStr
	}

	// 创建活动记录
	campaign := &model.Campaign{
		BrandId:             req.BrandId,
		Name:                req.Name,
		Description:         req.Description,
		FormFields:          string(formFieldsJSON),
		RewardRule:          req.RewardRule,
		StartTime:           time.Now(),
		EndTime:             time.Now(),
		Status:              "active",
		EnableDistribution:  req.EnableDistribution,
		DistributionLevel:   req.DistributionLevel,
		DistributionRewards: distributionRewardsJSON,
	}

	if err := l.svcCtx.DB.Create(campaign).Error; err != nil {
		return nil, fmt.Errorf("创建活动失败: %w", err)
	}

	// 查询品牌名称
	var brand model.Brand
	if err := l.svcCtx.DB.Where("id = ?", req.BrandId).First(&brand).Error; err == nil {
		return &types.CampaignResp{
			Id:                 campaign.Id,
			BrandId:            campaign.BrandId,
			BrandName:          brand.Name,
			Name:               campaign.Name,
			Description:        campaign.Description,
			FormFields:         req.FormFields,
			RewardRule:         campaign.RewardRule,
			StartTime:          campaign.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:            campaign.EndTime.Format("2006-01-02 15:04:05"),
			Status:             campaign.Status,
			EnableDistribution: campaign.EnableDistribution,
			DistributionLevel:  campaign.DistributionLevel,
			CreatedAt:          campaign.CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil
	}

	return &types.CampaignResp{
		Id:                 campaign.Id,
		BrandId:            campaign.BrandId,
		Name:               campaign.Name,
		Description:        campaign.Description,
		FormFields:         req.FormFields,
		RewardRule:         campaign.RewardRule,
		StartTime:          campaign.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:            campaign.EndTime.Format("2006-01-02 15:04:05"),
		Status:             campaign.Status,
		EnableDistribution: campaign.EnableDistribution,
		DistributionLevel:  campaign.DistributionLevel,
		CreatedAt:          campaign.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
