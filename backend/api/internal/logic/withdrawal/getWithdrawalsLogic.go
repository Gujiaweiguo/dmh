// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package withdrawal

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWithdrawalsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWithdrawalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWithdrawalsLogic {
	return &GetWithdrawalsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWithdrawalsLogic) GetWithdrawals(req *types.WithdrawalListReq) (resp *types.WithdrawalListResp, err error) {
	var withdrawals []model.Withdrawal
	var total int64

	query := l.svcCtx.DB.Model(&model.Withdrawal{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.UserId > 0 {
		query = query.Where("user_id = ?", req.UserId)
	}
	if req.BrandId > 0 {
		query = query.Where("brand_id = ?", req.BrandId)
	}

	if err := query.Count(&total).Error; err != nil {
		l.Errorf("Failed to count withdrawals: %v", err)
		return nil, fmt.Errorf("failed to count withdrawals")
	}

	offset := int((req.Page - 1) * req.PageSize)
	limit := int(req.PageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&withdrawals).Error; err != nil {
		l.Errorf("Failed to get withdrawals: %v", err)
		return nil, fmt.Errorf("failed to get withdrawals")
	}

	respList := make([]types.WithdrawalResp, 0, len(withdrawals))
	for _, w := range withdrawals {
		user := &model.User{}
		l.svcCtx.DB.First(user, w.UserID)

		withdrawalResp := types.WithdrawalResp{
			Id:          w.ID,
			UserId:      w.UserID,
			Username:    user.Username,
			RealName:    user.RealName,
			Amount:      w.Amount,
			Status:      w.Status,
			BankName:    w.BankName,
			BankAccount: w.BankAccount,
			AccountName: w.AccountName,
			Remark:      w.RejectedReason,
			CreatedAt:   w.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if w.ApprovedBy != nil {
			withdrawalResp.ApprovedBy = *w.ApprovedBy
		}
		if w.ApprovedAt != nil {
			withdrawalResp.ApprovedAt = w.ApprovedAt.Format("2006-01-02 15:04:05")
		}

		respList = append(respList, withdrawalResp)
	}

	resp = &types.WithdrawalListResp{
		Total:       total,
		Withdrawals: respList,
	}

	return resp, nil
}
