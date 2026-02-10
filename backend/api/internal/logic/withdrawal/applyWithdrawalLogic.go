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

type ApplyWithdrawalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyWithdrawalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyWithdrawalLogic {
	return &ApplyWithdrawalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyWithdrawalLogic) ApplyWithdrawal(req *types.WithdrawalApplyReq, userId int64) (resp *types.WithdrawalResp, err error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	if req.Amount < 10 {
		return nil, fmt.Errorf("amount must be at least 10")
	}

	userBalance := &model.UserBalance{}
	if err := l.svcCtx.DB.Where("user_id = ?", userId).First(userBalance).Error; err != nil {
		l.Errorf("Failed to get user balance: %v", err)
		return nil, fmt.Errorf("failed to get user balance")
	}

	if userBalance.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	var pendingCount int64
	if err := l.svcCtx.DB.Model(&model.Withdrawal{}).
		Where("user_id = ? AND status = ?", userId, "pending").
		Count(&pendingCount).Error; err != nil {
		l.Errorf("Failed to check pending withdrawal: %v", err)
		return nil, fmt.Errorf("failed to check pending withdrawal")
	}

	if pendingCount > 0 {
		return nil, fmt.Errorf("pending withdrawal exists")
	}

	withdrawal := &model.Withdrawal{
		UserID:      userId,
		Amount:      req.Amount,
		Status:      "pending",
		BankName:    req.BankName,
		BankAccount: req.BankAccount,
		AccountName: req.AccountName,
	}

	if err := l.svcCtx.DB.Create(withdrawal).Error; err != nil {
		l.Errorf("Failed to create withdrawal: %v", err)
		return nil, fmt.Errorf("failed to create withdrawal")
	}

	userBalance.Balance -= req.Amount
	if err := l.svcCtx.DB.Save(userBalance).Error; err != nil {
		l.Errorf("Failed to deduct balance: %v", err)
		return nil, fmt.Errorf("failed to deduct balance")
	}

	l.Infof("Withdrawal applied successfully: id=%d, userId=%d, amount=%.2f", withdrawal.ID, userId, req.Amount)

	resp = &types.WithdrawalResp{
		Id:          withdrawal.ID,
		UserId:      withdrawal.UserID,
		Amount:      withdrawal.Amount,
		Status:      withdrawal.Status,
		BankName:    withdrawal.BankName,
		BankAccount: withdrawal.BankAccount,
		AccountName: withdrawal.AccountName,
		CreatedAt:   withdrawal.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}
