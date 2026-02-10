// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package withdrawal

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveWithdrawalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveWithdrawalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveWithdrawalLogic {
	return &ApproveWithdrawalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveWithdrawalLogic) ApproveWithdrawal(withdrawalId int64, req *types.WithdrawalApproveReq, adminId int64) (resp *types.WithdrawalResp, err error) {
	withdrawal := &model.Withdrawal{}
	if err := l.svcCtx.DB.First(withdrawal, withdrawalId).Error; err != nil {
		l.Errorf("Withdrawal not found: %v", err)
		return nil, fmt.Errorf("withdrawal not found")
	}

	if withdrawal.Status != "pending" {
		return nil, fmt.Errorf("withdrawal can only be approved when pending")
	}

	if req.Status != "approved" && req.Status != "rejected" {
		return nil, fmt.Errorf("invalid status")
	}

	now := time.Now()
	withdrawal.Status = req.Status
	withdrawal.ApprovedBy = &adminId
	withdrawal.ApprovedAt = &now

	if req.Status == "rejected" {
		withdrawal.RejectedReason = req.Remark
		userBalance := &model.UserBalance{}
		if err := l.svcCtx.DB.Where("user_id = ?", withdrawal.UserID).First(userBalance).Error; err != nil {
			l.Errorf("Failed to get user balance: %v", err)
			return nil, fmt.Errorf("failed to get user balance")
		}
		userBalance.Balance += withdrawal.Amount
		if err := l.svcCtx.DB.Save(userBalance).Error; err != nil {
			l.Errorf("Failed to refund balance: %v", err)
			return nil, fmt.Errorf("failed to refund balance")
		}
	}

	if err := l.svcCtx.DB.Save(withdrawal).Error; err != nil {
		l.Errorf("Failed to update withdrawal: %v", err)
		return nil, fmt.Errorf("failed to update withdrawal")
	}

	l.Infof("Withdrawal approved: id=%d, status=%s, adminId=%d", withdrawalId, req.Status, adminId)

	resp = &types.WithdrawalResp{
		Id:          withdrawal.ID,
		UserId:      withdrawal.UserID,
		Amount:      withdrawal.Amount,
		Status:      withdrawal.Status,
		BankName:    withdrawal.BankName,
		BankAccount: withdrawal.BankAccount,
		AccountName: withdrawal.AccountName,
		Remark:      req.Remark,
		CreatedAt:   withdrawal.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if withdrawal.ApprovedAt != nil {
		resp.ApprovedAt = withdrawal.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	return resp, nil
}
