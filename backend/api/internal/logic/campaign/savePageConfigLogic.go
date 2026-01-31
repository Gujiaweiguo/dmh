// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"encoding/json"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePageConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSavePageConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePageConfigLogic {
	return &SavePageConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SavePageConfigLogic) SavePageConfig(req *types.PageConfigReq) (resp *types.PageConfigResp, err error) {
	campaignId := req.Id

	// 将components和theme序列化为JSON
	componentsJSON, err := json.Marshal(req.Components)
	if err != nil {
		l.Errorf("Failed to marshal components: %v", err)
		return nil, fmt.Errorf("Failed to marshal components: %w", err)
	}

	themeJSON, err := json.Marshal(req.Theme)
	if err != nil {
		l.Errorf("Failed to marshal theme: %v", err)
		return nil, fmt.Errorf("Failed to marshal theme: %w", err)
	}

	// 检查是否已存在配置
	var existingConfig model.PageConfig
	err = l.svcCtx.DB.Where("campaign_id = ? AND deleted_at IS NULL", campaignId).First(&existingConfig).Error

	if err != nil && err.Error() != "record not found" {
		l.Errorf("Failed to query page config: %v", err)
		return nil, fmt.Errorf("Failed to query page config: %w", err)
	}

	// 转换为string类型
	componentsStr := string(componentsJSON)
	themeStr := string(themeJSON)

	if err.Error() == "record not found" {
		// 创建新配置
		newConfig := &model.PageConfig{
			CampaignId: campaignId,
			Components: componentsStr,
			Theme:      themeStr,
		}
		if err := l.svcCtx.DB.Create(newConfig).Error; err != nil {
			l.Errorf("Failed to create page config: %v", err)
			return nil, fmt.Errorf("Failed to create page config: %w", err)
		}

		resp = &types.PageConfigResp{
			Id:         newConfig.Id,
			CampaignId: newConfig.CampaignId,
			Components: req.Components,
			Theme:      req.Theme,
			CreatedAt:  newConfig.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:  newConfig.UpdatedAt.Format("2006-01-02T15:04:05"),
		}
		l.Infof("Successfully created page config: campaignId=%d", campaignId)
	} else {
		// 更新现有配置
		existingConfig.Components = componentsStr
		existingConfig.Theme = themeStr
		if err := l.svcCtx.DB.Save(&existingConfig).Error; err != nil {
			l.Errorf("Failed to update page config: %v", err)
			return nil, fmt.Errorf("Failed to update page config: %w", err)
		}

		resp = &types.PageConfigResp{
			Id:         existingConfig.Id,
			CampaignId: existingConfig.CampaignId,
			Components: req.Components,
			Theme:      req.Theme,
			CreatedAt:  existingConfig.CreatedAt.Format("2006-01-02T15:04:05"),
			UpdatedAt:  existingConfig.UpdatedAt.Format("2006-01-02T15:04:05"),
		}
		l.Infof("Successfully updated page config: campaignId=%d", campaignId)
	}

	return resp, nil
}
