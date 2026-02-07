// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package campaign

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCampaignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCampaignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCampaignLogic {
	return &DeleteCampaignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCampaignLogic) DeleteCampaign(campaignId int64) (resp *types.CommonResp, err error) {
	now := time.Now()
	result := l.svcCtx.DB.Model(&model.Campaign{}).
		Where("id = ?", campaignId).
		Updates(map[string]interface{}{
			"deleted_at": now,
		})
	if result.Error != nil {
		l.Errorf("Failed to soft delete campaign: %v", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		l.Errorf("Campaign not found: id=%d", campaignId)
		return nil, fmt.Errorf("Campaign not found")
	}

	l.Infof("Campaign soft deleted successfully: campaignId=%d", campaignId)

	resp = &types.CommonResp{
		Message: "Campaign deleted successfully",
	}

	return resp, nil
}
