package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/model"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// WithdrawalLogic 提现业务逻辑
type WithdrawalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewWithdrawalLogic 创建提现逻辑实例
func NewWithdrawalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawalLogic {
	return &WithdrawalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ApplyWithdrawal 申请提现
func (l *WithdrawalLogic) ApplyWithdrawal(req *ApplyWithdrawalReq) (*model.Withdrawal, error) {
	// 1. 获取当前用户ID
	userID, ok := l.ctx.Value("userId").(int64)
	if !ok || userID == 0 {
		return nil, errors.New("获取用户信息失败")
	}

	// 2. 查询用户余额
	var userBalance model.UserBalance
	if err := l.svcCtx.DB.Where("user_id = ?", userID).First(&userBalance).Error; err != nil {
		return nil, fmt.Errorf("查询用户余额失败: %w", err)
	}

	// 3. 验证提现金额
	if req.Amount > userBalance.Balance {
		return nil, errors.New("提现金额超过可用余额")
	}
	if req.Amount <= 0 {
		return nil, errors.New("提现金额必须大于0")
	}

	// 4. 查询分销商信息
	var distributor model.Distributor
	if err := l.svcCtx.DB.Where("user_id = ?", userID).First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("查询分销商信息失败: %w", err)
	}

	// 5. 开启事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 6. 创建提现记录
	withdrawal := &model.Withdrawal{
		UserID:        userID,
		BrandId:       distributor.BrandId,
		DistributorId: distributor.Id,
		Amount:        req.Amount,
		Status:        "pending",
		PayType:       req.PayType,
		PayAccount:    req.PayAccount,
		PayRealName:   req.PayRealName,
	}

	if err := tx.Create(withdrawal).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建提现记录失败: %w", err)
	}

	// 7. 扣除用户余额（乐观锁）
	if err := l.deductBalance(tx, userID, req.Amount); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 8. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	return withdrawal, nil
}

// deductBalance 扣除用户余额（乐观锁）
func (l *WithdrawalLogic) deductBalance(tx *gorm.DB, userID int64, amount float64) error {
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		// 查询当前余额和版本号
		var userBalance model.UserBalance
		err := tx.Where("user_id = ?", userID).First(&userBalance).Error
		if err != nil {
			return fmt.Errorf("查询用户余额失败: %w", err)
		}

		// 尝试更新（带版本号检查）
		result := tx.Model(&model.UserBalance{}).
			Where("user_id = ? AND version = ?", userID, userBalance.Version).
			Updates(map[string]interface{}{
				"balance": userBalance.Balance - amount,
				"version": userBalance.Version + 1,
			})

		if result.RowsAffected > 0 {
			return nil // 更新成功
		}

		// 版本冲突，等待后重试
	}

	return errors.New("扣除用户余额失败：达到最大重试次数")
}

// ApproveWithdrawal 审批通过提现申请
func (l *WithdrawalLogic) ApproveWithdrawal(withdrawalID int64, notes string) error {
	// 1. 查询提现记录
	var withdrawal model.Withdrawal
	if err := l.svcCtx.DB.Where("id = ?", withdrawalID).First(&withdrawal).Error; err != nil {
		return fmt.Errorf("查询提现记录失败: %w", err)
	}

	// 2. 检查状态
	if withdrawal.Status != "pending" {
		return errors.New("提现申请已处理")
	}

	// 3. 获取审批人ID
	approverID, ok := l.ctx.Value("userId").(int64)
	if !ok || approverID == 0 {
		return errors.New("获取审批人信息失败")
	}

	// 4. 更新提现状态为approved
	now := time.Now()
	if err := l.svcCtx.DB.Model(&model.Withdrawal{}).Where("id = ?", withdrawalID).Updates(map[string]interface{}{
		"status":         "approved",
		"approved_by":    approverID,
		"approved_at":    &now,
		"approved_notes": notes,
	}).Error; err != nil {
		return fmt.Errorf("更新提现状态失败: %w", err)
	}

	// 5. 调用支付接口（模拟）
	// 实际应用中，这里应该调用微信/支付宝/银行的支付API
	// 这里暂时只更新状态为processing
	if err := l.processWithdrawal(withdrawalID); err != nil {
		return err
	}

	return nil
}

// processWithdrawal 处理提现（调用支付接口）
func (l *WithdrawalLogic) processWithdrawal(withdrawalID int64) error {
	// 模拟支付接口调用
	// 更新状态为processing
	if err := l.svcCtx.DB.Model(&model.Withdrawal{}).Where("id = ?", withdrawalID).Update("status", "processing").Error; err != nil {
		return fmt.Errorf("更新提现状态失败: %w", err)
	}

	// 这里应该调用实际的支付接口
	// 支付成功后更新状态为completed
	// 支付失败后更新状态为failed并记录原因

	// 模拟支付成功
	now := time.Now()
	tradeNo := fmt.Sprintf("WD%d%d", time.Now().Unix(), withdrawalID)
	if err := l.svcCtx.DB.Model(&model.Withdrawal{}).Where("id = ?", withdrawalID).Updates(map[string]interface{}{
		"status":   "completed",
		"paid_at":  &now,
		"trade_no": tradeNo,
	}).Error; err != nil {
		return fmt.Errorf("更新提现完成状态失败: %w", err)
	}

	return nil
}

// RejectWithdrawal 审批拒绝提现申请
func (l *WithdrawalLogic) RejectWithdrawal(withdrawalID int64, reason string) error {
	// 1. 查询提现记录
	var withdrawal model.Withdrawal
	if err := l.svcCtx.DB.Where("id = ?", withdrawalID).First(&withdrawal).Error; err != nil {
		return fmt.Errorf("查询提现记录失败: %w", err)
	}

	// 2. 检查状态
	if withdrawal.Status != "pending" {
		return errors.New("提现申请已处理")
	}

	// 3. 获取审批人ID
	approverID, ok := l.ctx.Value("userId").(int64)
	if !ok || approverID == 0 {
		return errors.New("获取审批人信息失败")
	}

	// 4. 开启事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 5. 更新提现状态为rejected
	now := time.Now()
	if err := tx.Model(&model.Withdrawal{}).Where("id = ?", withdrawalID).Updates(map[string]interface{}{
		"status":          "rejected",
		"approved_by":     approverID,
		"approved_at":     &now,
		"rejected_reason": reason,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新提现状态失败: %w", err)
	}

	// 6. 退还金额到用户余额（乐观锁）
	if err := l.refundBalance(tx, withdrawal.UserID, withdrawal.Amount); err != nil {
		tx.Rollback()
		return err
	}

	// 7. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// refundBalance 退还金额到用户余额（乐观锁）
func (l *WithdrawalLogic) refundBalance(tx *gorm.DB, userID int64, amount float64) error {
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		// 查询当前余额和版本号
		var userBalance model.UserBalance
		err := tx.Where("user_id = ?", userID).First(&userBalance).Error
		if err != nil {
			return fmt.Errorf("查询用户余额失败: %w", err)
		}

		// 尝试更新（带版本号检查）
		result := tx.Model(&model.UserBalance{}).
			Where("user_id = ? AND version = ?", userID, userBalance.Version).
			Updates(map[string]interface{}{
				"balance": userBalance.Balance + amount,
				"version": userBalance.Version + 1,
			})

		if result.RowsAffected > 0 {
			return nil // 更新成功
		}

		// 版本冲突，等待后重试
	}

	return errors.New("退还用户余额失败：达到最大重试次数")
}

// GetWithdrawals 获取提现记录列表
func (l *WithdrawalLogic) GetWithdrawals(brandID int64, status string, page int, pageSize int) (*WithdrawalsResp, error) {
	var withdrawals []model.Withdrawal
	var total int64

	// 构建查询
	query := l.svcCtx.DB.Model(&model.Withdrawal{})
	if brandID > 0 {
		query = query.Where("brand_id = ?", brandID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 查询总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&withdrawals).Error; err != nil {
		return nil, fmt.Errorf("查询提现记录失败: %w", err)
	}

	return &WithdrawalsResp{
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		Withdrawals: withdrawals,
	}, nil
}

// GetMyWithdrawals 获取我的提现记录
func (l *WithdrawalLogic) GetMyWithdrawals(page int, pageSize int) (*WithdrawalsResp, error) {
	// 获取当前用户ID
	userID, ok := l.ctx.Value("userId").(int64)
	if !ok || userID == 0 {
		return nil, errors.New("获取用户信息失败")
	}

	var withdrawals []model.Withdrawal
	var total int64

	// 查询总数
	l.svcCtx.DB.Model(&model.Withdrawal{}).Where("user_id = ?", userID).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := l.svcCtx.DB.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&withdrawals).Error; err != nil {
		return nil, fmt.Errorf("查询提现记录失败: %w", err)
	}

	return &WithdrawalsResp{
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		Withdrawals: withdrawals,
	}, nil
}

// ApplyWithdrawalReq 申请提现请求
type ApplyWithdrawalReq struct {
	Amount      float64 `json:"amount"`
	PayType     string  `json:"payType"`
	PayAccount  string  `json:"payAccount"`
	PayRealName string  `json:"payRealName"`
}

// WithdrawalsResp 提现记录响应
type WithdrawalsResp struct {
	Total       int64              `json:"total"`
	Page        int                `json:"page"`
	PageSize    int                `json:"pageSize"`
	Withdrawals []model.Withdrawal `json:"withdrawals"`
}
