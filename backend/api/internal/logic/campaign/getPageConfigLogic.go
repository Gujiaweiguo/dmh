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

type GetPageConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPageConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPageConfigLogic {
	return &GetPageConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPageConfigLogic) GetPageConfig(req *types.GetPageConfigReq) (resp *types.PageConfigResp, err error) {
	var config model.PageConfig
	campaignId := req.Id

	if err := l.svcCtx.DB.Where("campaign_id = ? AND deleted_at IS NULL", campaignId).First(&config).Error; err != nil {
		if err.Error() == "record not found" {
			// 如果没有配置，返回空配置
			resp = &types.PageConfigResp{
				CampaignId: campaignId,
				Components: []map[string]interface{}{},
				Theme:      map[string]interface{}{},
			}
			return resp, nil
		}
		l.Errorf("Failed to get page config: %v", err)
		return nil, fmt.Errorf("Failed to get page config: %w", err)
	}

	// 解析JSON
	var components []map[string]interface{}
	if err := json.Unmarshal([]byte(config.Components), &components); err != nil {
		l.Errorf("Failed to parse components JSON: %v", err)
		return nil, fmt.Errorf("Failed to parse components JSON: %w", err)
	}

	var theme map[string]interface{}
	if err := json.Unmarshal([]byte(config.Theme), &theme); err != nil {
		l.Errorf("Failed to parse theme JSON: %v", err)
		return nil, fmt.Errorf("Failed to parse theme JSON: %w", err)
	}

	resp = &types.PageConfigResp{
		Id:         config.Id,
		CampaignId: config.CampaignId,
		Components: components,
		Theme:      theme,
		CreatedAt:  config.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt:  config.UpdatedAt.Format("2006-01-02T15:04:05"),
	}

	l.Infof("Successfully got page config: campaignId=%d", campaignId)
	return resp, nil
}
