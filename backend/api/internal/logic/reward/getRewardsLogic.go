// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package reward

import (
	"context"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRewardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRewardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRewardsLogic {
	return &GetRewardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRewardsLogic) GetRewards(req *types.GetRewardsReq) (resp []types.RewardResp, err error) {
	var rewards []model.DistributorReward

	query := l.svcCtx.DB.Model(&model.DistributorReward{})

	if req.UserId > 0 {
		query = query.Where("user_id = ?", req.UserId)
	}
	if req.OrderId > 0 {
		query = query.Where("order_id = ?", req.OrderId)
	}

	if err := query.Order("created_at DESC").
		Limit(100).
		Find(&rewards).Error; err != nil {
		l.Errorf("Failed to get rewards: %v", err)
		return nil, err
	}

	resp = make([]types.RewardResp, 0, len(rewards))
	for _, r := range rewards {
		resp = append(resp, types.RewardResp{
			Id:        r.Id,
			UserId:    r.UserId,
			OrderId:   r.OrderId,
			Amount:    r.Amount,
			Status:    r.Status,
			CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}
