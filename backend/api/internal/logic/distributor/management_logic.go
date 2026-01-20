package distributor

import (
	"context"
	"fmt"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManagementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManagementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManagementLogic {
	return &ManagementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetDistributors 获取分销商列表
func (l *ManagementLogic) GetDistributors(req *types.GetDistributorsReq) (*types.DistributorListResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	query := l.svcCtx.DB.Model(&model.Distributor{})

	// 品牌管理员只能查看自己品牌的分销商
	if user.Role != "platform_admin" {
		if req.BrandId == 0 {
			// 获取用户管理的品牌列表
			var userBrands []model.UserBrand
			l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands)

			if len(userBrands) == 0 {
				return &types.DistributorListResp{
					Total:        0,
					Distributors: []types.DistributorResp{},
				}, nil
			}

			brandIds := make([]int64, len(userBrands))
			for i, ub := range userBrands {
				brandIds[i] = ub.BrandId
			}
			query = query.Where("brand_id IN ?", brandIds)
		} else {
			// 检查是否有权限查看该品牌
			var count int64
			l.svcCtx.DB.Model(&model.UserBrand{}).
				Where("user_id = ? AND brand_id = ?", userId, req.BrandId).
				Count(&count)
			if count == 0 {
				return nil, fmt.Errorf("无权限查看该品牌的分销商")
			}
			query = query.Where("brand_id = ?", req.BrandId)
		}
	} else if req.BrandId > 0 {
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

// GetDistributor 获取分销商详情
func (l *ManagementLogic) GetDistributor(id int64) (*types.DistributorResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	var distributor model.Distributor
	query := l.svcCtx.DB.Preload("User").Preload("Brand").Preload("Parent")

	// 品牌管理员只能查看自己品牌的分销商
	if user.Role != "platform_admin" {
		var userBrands []model.UserBrand
		l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands)

		if len(userBrands) == 0 {
			return nil, fmt.Errorf("无权限查看该分销商")
		}

		brandIds := make([]int64, len(userBrands))
		for i, ub := range userBrands {
			brandIds[i] = ub.BrandId
		}
		query = query.Where("distributors.id = ? AND brand_id IN ?", id, brandIds)
	} else {
		query = query.Where("distributors.id = ?", id)
	}

	if err := query.First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("分销商不存在")
	}

	return l.buildDistributorResp(distributor), nil
}

// UpdateDistributorLevel 更新分销商级别
func (l *ManagementLogic) UpdateDistributorLevel(id int64, req *types.UpdateDistributorLevelReq) error {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证级别
	if req.Level < 1 || req.Level > 3 {
		return fmt.Errorf("级别必须在1-3之间")
	}

	query := l.svcCtx.DB.Model(&model.Distributor{}).Where("id = ?", id)

	// 品牌管理员只能修改自己品牌的分销商
	if user.Role != "platform_admin" {
		var userBrands []model.UserBrand
		l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands)

		if len(userBrands) == 0 {
			return fmt.Errorf("无权限修改该分销商")
		}

		brandIds := make([]int64, len(userBrands))
		for i, ub := range userBrands {
			brandIds[i] = ub.BrandId
		}
		query = query.Where("brand_id IN ?", brandIds)
	}

	result := query.Update("level", req.Level)
	if result.Error != nil {
		return fmt.Errorf("更新失败")
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("分销商不存在或无权限修改")
	}

	return nil
}

// UpdateDistributorStatus 更新分销商状态
func (l *ManagementLogic) UpdateDistributorStatus(id int64, req *types.UpdateDistributorStatusReq) error {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证状态
	if req.Status != "active" && req.Status != "suspended" {
		return fmt.Errorf("状态必须是 active 或 suspended")
	}

	query := l.svcCtx.DB.Model(&model.Distributor{}).Where("id = ?", id)

	// 品牌管理员只能修改自己品牌的分销商
	if user.Role != "platform_admin" {
		var userBrands []model.UserBrand
		l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands)

		if len(userBrands) == 0 {
			return fmt.Errorf("无权限修改该分销商")
		}

		brandIds := make([]int64, len(userBrands))
		for i, ub := range userBrands {
			brandIds[i] = ub.BrandId
		}
		query = query.Where("brand_id IN ?", brandIds)
	}

	result := query.Update("status", req.Status)
	if result.Error != nil {
		return fmt.Errorf("更新失败")
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("分销商不存在或无权限修改")
	}

	return nil
}

// buildDistributorResp 构建分销商响应
func (l *ManagementLogic) buildDistributorResp(distributor model.Distributor) *types.DistributorResp {
	username := ""
	if distributor.User != nil {
		username = distributor.User.RealName
		if username == "" {
			username = distributor.User.Username
		}
	}

	brandName := ""
	if distributor.Brand != nil {
		brandName = distributor.Brand.Name
	}

	var parentName string
	var parentId int64
	if distributor.ParentId != nil {
		parentId = *distributor.ParentId
		if distributor.Parent != nil && distributor.Parent.User != nil {
			parentName = distributor.Parent.User.RealName
			if parentName == "" {
				parentName = distributor.Parent.User.Username
			}
		}
	}

	resp := &types.DistributorResp{
		Id:                distributor.Id,
		UserId:            distributor.UserId,
		Username:          username,
		BrandId:           distributor.BrandId,
		BrandName:         brandName,
		Level:             distributor.Level,
		ParentId:          parentId,
		ParentName:        parentName,
		Status:            distributor.Status,
		TotalEarnings:     distributor.TotalEarnings,
		SubordinatesCount: distributor.SubordinatesCount,
		CreatedAt:         distributor.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if distributor.ApprovedBy != nil {
		resp.ApprovedBy = *distributor.ApprovedBy
	}

	if distributor.ApprovedAt != nil {
		resp.ApprovedAt = distributor.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	return resp
}

// GetDistributorsByBrand 根据品牌获取分销商列表（简化版，用于下拉选择）
func (l *ManagementLogic) GetDistributorsByBrand(brandId int64) ([]types.DistributorResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 品牌管理员只能查看自己品牌的分销商
	if user.Role != "platform_admin" {
		var count int64
		l.svcCtx.DB.Model(&model.UserBrand{}).
			Where("user_id = ? AND brand_id = ?", userId, brandId).
			Count(&count)
		if count == 0 {
			return nil, fmt.Errorf("无权限查看该品牌的分销商")
		}
	}

	var distributors []model.Distributor
	l.svcCtx.DB.Preload("User").
		Where("brand_id = ? AND status = ?", brandId, "active").
		Order("level ASC, created_at ASC").
		Find(&distributors)

	resp := make([]types.DistributorResp, 0, len(distributors))
	for _, d := range distributors {
		username := ""
		if d.User != nil {
			username = d.User.RealName
			if username == "" {
				username = d.User.Username
			}
		}

		resp = append(resp, types.DistributorResp{
			Id:       d.Id,
			UserId:   d.UserId,
			Username: username,
			BrandId:  d.BrandId,
			Level:    d.Level,
			Status:   d.Status,
		})
	}

	return resp, nil
}

// GetMyDistributorStatus 获取当前用户的分销商状态
func (l *ManagementLogic) GetMyDistributorStatus(brandId int64) (*types.DistributorResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	var distributor model.Distributor
	if err := l.svcCtx.DB.Preload("User").Preload("Brand").
		Where("user_id = ? AND brand_id = ?", userId, brandId).
		First(&distributor).Error; err != nil {
		return nil, fmt.Errorf("未找到分销商记录")
	}

	return l.buildDistributorResp(distributor), nil
}

// GetCustomers 获取顾客列表（参与活动但不是分销商的用户）
func (l *ManagementLogic) GetCustomers(req *types.GetCustomersReq) (*types.CustomerListResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 构建基础查询
	query := l.svcCtx.DB.Table("orders o").
		Select(`
			o.id as id,
			MAX(o.member_id) as user_id,
			COALESCE(MAX(u.username), MAX(u.real_name), o.phone) as username,
			o.phone,
			MAX(o.campaign_id) as campaign_id,
			COALESCE(MAX(c.name), '') as campaign_name,
			COUNT(DISTINCT o.id) as order_count,
			SUM(CASE WHEN o.pay_status = 'paid' THEN o.amount ELSE 0 END) as total_amount,
			MIN(o.created_at) as first_order_at,
			MAX(o.created_at) as last_order_at,
			MAX(o.created_at) as created_at
		`).
		Joins("LEFT JOIN campaigns c ON o.campaign_id = c.id").
		Joins("LEFT JOIN members m ON o.member_id = m.id").
		Joins("LEFT JOIN users u ON m.user_id = u.id").
		Where("o.deleted_at IS NULL")

	// 品牌管理员只能查看自己品牌的顾客
	if user.Role != "platform_admin" {
		var userBrands []model.UserBrand
		l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands)

		if len(userBrands) == 0 {
			return &types.CustomerListResp{
				Total:     0,
				Customers: []types.CustomerResp{},
			}, nil
		}

		brandIds := make([]int64, len(userBrands))
		for i, ub := range userBrands {
			brandIds[i] = ub.BrandId
		}
		query = query.Where("o.campaign_id IN (SELECT id FROM campaigns WHERE brand_id IN ?)", brandIds)
	} else if req.BrandId > 0 {
		query = query.Where("o.campaign_id IN (SELECT id FROM campaigns WHERE brand_id = ?)", req.BrandId)
	}

	// 筛选条件
	if req.CampaignId > 0 {
		query = query.Where("o.campaign_id = ?", req.CampaignId)
	}
	if req.Keyword != "" {
		query = query.Where("(o.phone LIKE ? OR u.username LIKE ? OR u.real_name LIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		query = query.Where("o.pay_status = ?", req.Status)
	}

	// 排除非分销商（排除已成为分销商的用户）
	query = query.Where("o.member_id NOT IN (SELECT user_id FROM distributors)")

	// 统计总数
	var total int64
	countQuery := l.svcCtx.DB.Table("(?) as sub", query)
	countQuery.Count(&total)

	// 分页
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	query.Group("o.phone, o.member_id").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	// 执行查询
	rows, err := query.Rows()
	if err != nil {
		return nil, fmt.Errorf("查询顾客列表失败: %w", err)
	}
	defer rows.Close()

	customers := []types.CustomerResp{}
	for rows.Next() {
		var c types.CustomerResp
		if err := rows.Scan(&c.Id, &c.UserId, &c.Username, &c.Phone,
			&c.CampaignId, &c.CampaignName, &c.OrderCount, &c.TotalAmount,
			&c.FirstOrderAt, &c.LastOrderAt, &c.CreatedAt); err != nil {
			logx.Errorf("扫描顾客数据失败: %v", err)
			continue
		}

		// 查询品牌名称
		if c.CampaignId > 0 {
			var campaign model.Campaign
			if err := l.svcCtx.DB.Where("id = ?", c.CampaignId).First(&campaign).Error; err == nil {
				c.BrandId = campaign.BrandId
				var brand model.Brand
				if err := l.svcCtx.DB.Where("id = ?", campaign.BrandId).First(&brand).Error; err == nil {
					c.BrandName = brand.Name
				}
			}
		}

		customers = append(customers, c)
	}

	return &types.CustomerListResp{
		Total:     total,
		Customers: customers,
	}, nil
}

// GetBrandRewards 获取品牌奖励详情
func (l *ManagementLogic) GetBrandRewards(req *types.GetBrandRewardsReq) (*types.BrandRewardListResp, error) {
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok || userId == 0 {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 检查权限
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	query := l.svcCtx.DB.Table("rewards r").
		Select(`
			r.id,
			r.user_id,
			COALESCE(u.username, u.real_name, '未知用户') as username,
			r.order_id,
			r.campaign_id,
			COALESCE(c.name, '') as campaign_name,
			r.amount,
			r.status,
			r.settled_at,
			r.created_at
		`).
		Joins("LEFT JOIN users u ON r.user_id = u.id").
		Joins("LEFT JOIN campaigns c ON r.campaign_id = c.id").
		Where("r.deleted_at IS NULL")

	// 品牌管理员只能查看自己品牌的奖励
	if user.Role != "platform_admin" {
		var userBrands []model.UserBrand
		l.svcCtx.DB.Where("user_id = ?", userId).Find(&userBrands)

		if len(userBrands) == 0 {
			return &types.BrandRewardListResp{
				Total:   0,
				Rewards: []types.BrandRewardResp{},
			}, nil
		}

		brandIds := make([]int64, len(userBrands))
		for i, ub := range userBrands {
			brandIds[i] = ub.BrandId
		}
		query = query.Where("r.campaign_id IN (SELECT id FROM campaigns WHERE brand_id IN ?)", brandIds)
	} else if req.BrandId > 0 {
		query = query.Where("r.campaign_id IN (SELECT id FROM campaigns WHERE brand_id = ?)", req.BrandId)
	}

	// 筛选条件
	if req.CampaignId > 0 {
		query = query.Where("r.campaign_id = ?", req.CampaignId)
	}
	if req.Keyword != "" {
		query = query.Where("(u.username LIKE ? OR u.real_name LIKE ? OR u.phone LIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		query = query.Where("r.status = ?", req.Status)
	}
	if req.StartDate != "" {
		query = query.Where("r.created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("r.created_at <= ?", req.EndDate)
	}

	// 统计总数
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

	// 执行查询
	rows, err := query.Order("r.created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("查询奖励列表失败: %w", err)
	}
	defer rows.Close()

	rewards := []types.BrandRewardResp{}
	for rows.Next() {
		var r types.BrandRewardResp
		var settledAt *string
		if err := rows.Scan(&r.Id, &r.UserId, &r.Username, &r.OrderId,
			&r.CampaignId, &r.CampaignName, &r.Amount, &r.Status,
			&settledAt, &r.CreatedAt); err != nil {
			logx.Errorf("扫描奖励数据失败: %v", err)
			continue
		}

		if settledAt != nil {
			r.SettledAt = *settledAt
		}

		rewards = append(rewards, r)
	}

	return &types.BrandRewardListResp{
		Total:   total,
		Rewards: rewards,
	}, nil
}
