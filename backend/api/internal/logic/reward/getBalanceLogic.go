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

type GetBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBalanceLogic {
	return &GetBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBalanceLogic) GetBalance(userId int64) (resp *types.BalanceResp, err error) {
	balance := &model.UserBalance{}
	if err := l.svcCtx.DB.Where("user_id = ?", userId).First(balance).Error; err != nil {
		l.Errorf("Failed to get user balance: %v", err)
		return nil, err
	}

	resp = &types.BalanceResp{
		UserId:      balance.UserId,
		Balance:     balance.Balance,
		TotalReward: balance.TotalReward,
	}

	return resp, nil
}
