package promoter

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromoterRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPromoterRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPromoterRewardsLogic {
	return &GetPromoterRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPromoterRewardsLogic) GetPromoterRewards(req *types.GetPromoterRewardsReq) (resp *types.PromoterRewardsListResp, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	query := l.svcCtx.DB.Model(&model.PromoterReward{}).Where("promoter_id = ?", req.PromoterId)

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count rewards: %w", err)
	}

	var rewards []model.PromoterReward
	offset := int((req.Page - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).
		Order("created_at DESC").
		Find(&rewards).Error; err != nil {
		return nil, fmt.Errorf("failed to get rewards: %w", err)
	}

	resp = &types.PromoterRewardsListResp{
		Total:   total,
		Rewards: make([]types.PromoterRewardResp, 0, len(rewards)),
	}

	for _, r := range rewards {
		resp.Rewards = append(resp.Rewards, types.PromoterRewardResp{
			Id:          r.Id,
			Type:        r.Type,
			Status:      r.Status,
			Amount:      r.Amount,
			Description: r.Description,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}
