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

type GetWithdrawalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWithdrawalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWithdrawalLogic {
	return &GetWithdrawalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWithdrawalLogic) GetWithdrawal(withdrawalId int64) (resp *types.WithdrawalResp, err error) {
	withdrawal := &model.Withdrawal{}
	if err := l.svcCtx.DB.First(withdrawal, withdrawalId).Error; err != nil {
		l.Errorf("Withdrawal not found: %v", err)
		return nil, fmt.Errorf("withdrawal not found")
	}

	user := &model.User{}
	l.svcCtx.DB.First(user, withdrawal.UserID)

	resp = &types.WithdrawalResp{
		Id:          withdrawal.ID,
		UserId:      withdrawal.UserID,
		Username:    user.Username,
		Amount:      withdrawal.Amount,
		Status:      withdrawal.Status,
		BankName:    withdrawal.BankName,
		BankAccount: withdrawal.BankAccount,
		AccountName: withdrawal.AccountName,
		Remark:      withdrawal.RejectedReason,
		CreatedAt:   withdrawal.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if withdrawal.ApprovedAt != nil {
		resp.ApprovedAt = withdrawal.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	return resp, nil
}
