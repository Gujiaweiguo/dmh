package distributor

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

// GlobalStatsLogic 平台全局统计业务逻辑
type GlobalStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGlobalStatsLogic 创建全局统计逻辑实例
func NewGlobalStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GlobalStatsLogic {
	return &GlobalStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetGlobalStats 获取平台全局统计数据
func (l *GlobalStatsLogic) GetGlobalStats(startDate string, endDate string) (*types.GlobalStatsResp, error) {
	// 构建基础查询
	baseQuery := l.svcCtx.DB.Model(&model.Distributor{})
	if startDate != "" {
		baseQuery = baseQuery.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		baseQuery = baseQuery.Where("created_at <= ?", endDate)
	}

	// 1. 统计总分销商数量
	var totalDistributors int64
	if err := baseQuery.Count(&totalDistributors).Error; err != nil {
		return nil, fmt.Errorf("统计总分销商数量失败: %w", err)
	}

	// 2. 统计活跃分销商数量（状态为active）
	var activeDistributors int64
	if err := baseQuery.Where("status = ?", "active").Count(&activeDistributors).Error; err != nil {
		return nil, fmt.Errorf("统计活跃分销商数量失败: %w", err)
	}

	// 3. 统计总奖励金额
	var totalRewards float64
	if err := l.svcCtx.DB.Model(&model.DistributorReward{}).
		Joins("JOIN distributors ON distributors.id = distributor_rewards.distributor_id").
		Where("distributors.status = ?", "active").
		Select("COALESCE(SUM(distributor_rewards.amount), 0)").
		Scan(&totalRewards).Error; err != nil {
		return nil, fmt.Errorf("统计总奖励金额失败: %w", err)
	}

	// 4. 统计总提现金额
	var totalWithdrawals float64
	if err := l.svcCtx.DB.Model(&model.Withdrawal{}).
		Where("status IN (?)", []string{"approved", "processing", "completed"}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalWithdrawals).Error; err != nil {
		return nil, fmt.Errorf("统计总提现金额失败: %w", err)
	}

	// 5. 统计待审批提现数量
	var pendingWithdrawals int64
	if err := l.svcCtx.DB.Model(&model.Withdrawal{}).
		Where("status = ?", "pending").
		Count(&pendingWithdrawals).Error; err != nil {
		return nil, fmt.Errorf("统计待审批提现数量失败: %w", err)
	}

	// 6. 统计总订单数量
	var totalOrders int64
	if err := l.svcCtx.DB.Model(&model.Order{}).
		Where("status = ?", "paid").
		Count(&totalOrders).Error; err != nil {
		return nil, fmt.Errorf("统计总订单数量失败: %w", err)
	}

	// 7. 统计总收入
	var totalRevenue float64
	if err := l.svcCtx.DB.Model(&model.Order{}).
		Where("status = ?", "paid").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, fmt.Errorf("统计总收入失败: %w", err)
	}

	return &types.GlobalStatsResp{
		TotalDistributors:  totalDistributors,
		ActiveDistributors: activeDistributors,
		TotalRewards:       totalRewards,
		TotalWithdrawals:   totalWithdrawals,
		PendingWithdrawals: pendingWithdrawals,
		TotalOrders:        totalOrders,
		TotalRevenue:       totalRevenue,
	}, nil
}

// GetPlatformDistributors 获取全局分销商列表（平台管理员）
func (l *GlobalStatsLogic) GetPlatformDistributors(req *types.GetPlatformDistributorsReq) (*types.DistributorListResp, error) {
	query := l.svcCtx.DB.Model(&model.Distributor{})

	// 按品牌筛选
	if req.BrandId > 0 {
		query = query.Where("brand_id = ?", req.BrandId)
	}

	// 筛选条件
	if req.Level > 0 {
		query = query.Where("level = ?", req.Level)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		query = query.Joins("JOIN users ON distributors.user_id = users.id").
			Where("users.username LIKE ? OR users.phone LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	var total int64
	query.Count(&total)

	// 分页
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var distributors []model.Distributor
	offset := (req.Page - 1) * req.PageSize
	query.Preload("User").Preload("Brand").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&distributors)

	resp := &types.DistributorListResp{
		Total:        total,
		Distributors: make([]types.DistributorResp, 0),
	}

	for _, d := range distributors {
		var parentName string
		if d.ParentId != nil {
			var parent model.Distributor
			if l.svcCtx.DB.Where("id = ?", *d.ParentId).First(&parent).Error == nil {
				if parent.User != nil {
					parentName = parent.User.RealName
					if parentName == "" {
						parentName = parent.User.Username
					}
				}
			}
		}

		username := ""
		if d.User != nil {
			username = d.User.RealName
			if username == "" {
				username = d.User.Username
			}
		}

		brandName := ""
		if d.Brand != nil {
			brandName = d.Brand.Name
		}

		distributorResp := types.DistributorResp{
			Id:                d.Id,
			UserId:            d.UserId,
			Username:          username,
			BrandId:           d.BrandId,
			BrandName:         brandName,
			Level:             d.Level,
			Status:            d.Status,
			TotalEarnings:     d.TotalEarnings,
			SubordinatesCount: d.SubordinatesCount,
			CreatedAt:         d.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if d.ParentId != nil {
			distributorResp.ParentId = *d.ParentId
			distributorResp.ParentName = parentName
		}

		if d.ApprovedBy != nil {
			distributorResp.ApprovedBy = *d.ApprovedBy
		}

		if d.ApprovedAt != nil {
			distributorResp.ApprovedAt = d.ApprovedAt.Format("2006-01-02 15:04:05")
		}

		resp.Distributors = append(resp.Distributors, distributorResp)
	}

	return resp, nil
}

// GetPlatformRewards 获取全局奖励列表（平台管理员）
func (l *GlobalStatsLogic) GetPlatformRewards(req *types.GetPlatformRewardsReq) (*types.DistributorRewardListResp, error) {
	query := l.svcCtx.DB.Table("distributor_rewards dr").
		Select(`
			dr.id,
			dr.order_id,
			dr.amount,
			dr.level,
			dr.reward_rate,
			dr.from_user_id,
			COALESCE(u.username, u.real_name, '未知用户') as from_username,
			dr.status,
			dr.settled_at,
			dr.created_at
		`).
		Joins("LEFT JOIN users u ON dr.from_user_id = u.id")

	// 按品牌筛选
	if req.BrandId > 0 {
		query = query.Where("dr.campaign_id IN (SELECT id FROM campaigns WHERE brand_id = ?)", req.BrandId)
	}

	// 筛选条件
	if req.CampaignId > 0 {
		query = query.Where("dr.campaign_id = ?", req.CampaignId)
	}
	if req.Keyword != "" {
		query = query.Where("(u.username LIKE ? OR u.real_name LIKE ? OR u.phone LIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		query = query.Where("dr.status = ?", req.Status)
	}
	if req.StartDate != "" {
		query = query.Where("dr.created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("dr.created_at <= ?", req.EndDate)
	}

	var total int64
	query.Count(&total)

	// 分页
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	rows, err := query.Order("dr.created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("查询奖励列表失败: %w", err)
	}
	defer rows.Close()

	rewards := []types.DistributorRewardResp{}
	for rows.Next() {
		var r types.DistributorRewardResp
		var fromUserId *int64
		var fromUsername *string
		var settledAt *string
		if err := rows.Scan(&r.Id, &r.OrderId, &r.Amount, &r.Level, &r.RewardRate,
			&fromUserId, &fromUsername, &r.Status, &settledAt, &r.CreatedAt); err != nil {
			logx.Errorf("扫描奖励数据失败: %v", err)
			continue
		}

		if fromUserId != nil {
			r.FromUserId = *fromUserId
		}
		if fromUsername != nil {
			r.FromUsername = *fromUsername
		}
		if settledAt != nil {
			r.SettledAt = *settledAt
		}

		rewards = append(rewards, r)
	}

	return &types.DistributorRewardListResp{
		Total:   total,
		Rewards: rewards,
	}, nil
}

// GetPlatformWithdrawals 获取全局提现列表（平台管理员）
func (l *GlobalStatsLogic) GetPlatformWithdrawals(req *types.GetPlatformWithdrawalsReq) (*types.WithdrawalListResp, error) {
	query := l.svcCtx.DB.Table("withdrawals w").
		Select(`
			w.id,
			w.user_id,
			COALESCE(u.username, u.real_name, '未知用户') as username,
			COALESCE(u.real_name, '') as real_name,
			w.amount,
			w.bank_name,
			w.bank_account,
			w.account_name,
			w.status,
			w.remark,
			w.created_at
		`).
		Joins("LEFT JOIN users u ON w.user_id = u.id").
		Joins("LEFT JOIN distributors d ON w.distributor_id = d.id")

	// 按品牌筛选
	if req.BrandId > 0 {
		query = query.Where("d.brand_id = ?", req.BrandId)
	}

	// 筛选条件
	if req.Keyword != "" {
		query = query.Where("(u.username LIKE ? OR u.real_name LIKE ? OR u.phone LIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		query = query.Where("w.status = ?", req.Status)
	}
	if req.StartDate != "" {
		query = query.Where("w.created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("w.created_at <= ?", req.EndDate)
	}

	var total int64
	query.Count(&total)

	// 分页
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	rows, err := query.Order("w.created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("查询提现列表失败: %w", err)
	}
	defer rows.Close()

	withdrawals := []types.WithdrawalResp{}
	for rows.Next() {
		var w types.WithdrawalResp
		if err := rows.Scan(&w.Id, &w.UserId, &w.Username, &w.RealName,
			&w.Amount, &w.BankName, &w.BankAccount, &w.AccountName,
			&w.Status, &w.Remark, &w.CreatedAt); err != nil {
			logx.Errorf("扫描提现数据失败: %v", err)
			continue
		}

		withdrawals = append(withdrawals, w)
	}

	return &types.WithdrawalListResp{
		Total:       total,
		Withdrawals: withdrawals,
	}, nil
}
