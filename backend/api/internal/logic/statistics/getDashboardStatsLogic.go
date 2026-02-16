// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package statistics

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetDashboardStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDashboardStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDashboardStatsLogic {
	return &GetDashboardStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// buildOrderBaseQuery creates a fresh order query with common filters.
// Each call returns a new query to avoid JOIN accumulation.
func (l *GetDashboardStatsLogic) buildOrderBaseQuery(req *types.GetDashboardStatsReq) *gorm.DB {
	query := l.svcCtx.DB.Model(&model.Order{})

	if req.StartDate != "" {
		query = query.Where("orders.created_at >= ?", req.StartDate)
	}

	if req.EndDate != "" {
		query = query.Where("orders.created_at <= ?", req.EndDate)
	}

	if req.BrandId > 0 {
		query = query.Joins("JOIN campaigns ON campaigns.id = orders.campaign_id").
			Where("campaigns.brand_id = ?", req.BrandId)
	}

	return query
}

func (l *GetDashboardStatsLogic) GetDashboardStats(req *types.GetDashboardStatsReq) (resp *types.DashboardStatsResp, err error) {
	if req.BrandId <= 0 {
		return nil, errors.New("brandId is required")
	}

	var totalOrders int64
	if err := l.buildOrderBaseQuery(req).Select("COUNT(*)").Scan(&totalOrders).Error; err != nil {
		return nil, err
	}

	var totalRevenue float64
	if err := l.buildOrderBaseQuery(req).Select("COALESCE(SUM(orders.amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}

	now := l.svcCtx.DB.NowFunc()
	todayStart := now.Format("2006-01-02")
	weekStart := now.Add(-7 * 24 * time.Hour).Format("2006-01-02")
	monthStart := now.Format("2006-01-01")

	var todayOrders int64
	var todayRevenue float64
	if err := l.buildOrderBaseQuery(req).Where("DATE(orders.created_at) = ?", todayStart).
		Select("COUNT(*)").Scan(&todayOrders).Error; err != nil {
		return nil, err
	}
	if err := l.buildOrderBaseQuery(req).Where("DATE(orders.created_at) = ?", todayStart).
		Select("COALESCE(SUM(orders.amount), 0)").Scan(&todayRevenue).Error; err != nil {
		return nil, err
	}

	var weekOrders int64
	var weekRevenue float64
	if err := l.buildOrderBaseQuery(req).Where("orders.created_at >= ?", weekStart).
		Select("COUNT(*)").Scan(&weekOrders).Error; err != nil {
		return nil, err
	}
	if err := l.buildOrderBaseQuery(req).Where("orders.created_at >= ?", weekStart).
		Select("COALESCE(SUM(orders.amount), 0)").Scan(&weekRevenue).Error; err != nil {
		return nil, err
	}

	var monthOrders int64
	var monthRevenue float64
	if err := l.buildOrderBaseQuery(req).Where("orders.created_at >= ?", monthStart).
		Select("COUNT(*)").Scan(&monthOrders).Error; err != nil {
		return nil, err
	}
	if err := l.buildOrderBaseQuery(req).Where("orders.created_at >= ?", monthStart).
		Select("COALESCE(SUM(orders.amount), 0)").Scan(&monthRevenue).Error; err != nil {
		return nil, err
	}

	var orderTrends []types.TrendData
	orderQuery := l.svcCtx.DB.Model(&model.Order{}).
		Select("DATE(orders.created_at) as date, COUNT(*) as value")

	if req.BrandId > 0 {
		orderQuery = orderQuery.Joins("JOIN campaigns ON campaigns.id = orders.campaign_id").
			Where("campaigns.brand_id = ?", req.BrandId)
	}

	var dateCondition string
	if req.Period == "week" {
		dateCondition = "DATE_SUB(NOW(), INTERVAL 7 DAY)"
	} else if req.Period == "month" {
		dateCondition = "DATE_SUB(NOW(), INTERVAL 90 DAY)"
	} else if req.Period == "year" {
		dateCondition = "DATE_SUB(NOW(), INTERVAL 365 DAY)"
	} else {
		dateCondition = "DATE_SUB(NOW(), INTERVAL 90 DAY)"
	}

	if err := orderQuery.Where("orders.created_at >= ?", dateCondition).
		Group("DATE(orders.created_at)").
		Order("DATE(orders.created_at) ASC").
		Scan(&orderTrends).Error; err != nil {
		return nil, err
	}

	var topCampaigns []types.CampaignStats
	campaignQuery := l.svcCtx.DB.Model(&model.Campaign{}).
		Select("campaigns.id, campaigns.name, (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) as order_count, (SELECT COALESCE(SUM(amount), 0) FROM orders WHERE orders.campaign_id = campaigns.id) as revenue")

	if req.BrandId > 0 {
		campaignQuery = campaignQuery.Where("campaigns.brand_id = ?", req.BrandId)
	}

	if err := campaignQuery.Where("campaigns.created_at >= ?", dateFilterQuery(req)).
		Order("(SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) DESC").
		Limit(5).
		Find(&topCampaigns).Error; err != nil {
		return nil, err
	}

	for i := range topCampaigns {
		l.svcCtx.DB.Raw("SELECT COUNT(*) FROM orders WHERE campaign_id = ?", topCampaigns[i].Id).Scan(&topCampaigns[i].OrderCount)
		l.svcCtx.DB.Raw("SELECT COALESCE(SUM(amount), 0) FROM orders WHERE campaign_id = ?", topCampaigns[i].Id).Scan(&topCampaigns[i].Revenue)
	}

	var topDistributors []types.DistributorStats
	distributorQuery := l.svcCtx.DB.Table("distributors").
		Select("distributors.id, users.username, distributors.level, distributors.total_earnings, distributors.subordinates_count")

	if req.BrandId > 0 {
		distributorQuery = distributorQuery.Joins("LEFT JOIN users ON users.id = distributors.user_id").
			Where("distributors.brand_id = ?", req.BrandId)
	}

	if err := distributorQuery.Order("distributors.total_earnings DESC").
		Limit(5).
		Find(&topDistributors).Error; err != nil {
		l.Logger.Errorf("failed to query top distributors: %v", err)
		topDistributors = []types.DistributorStats{}
	}

	for i := range topDistributors {
		topDistributors[i].TotalOrders = 0
		topDistributors[i].TotalRevenue = 0
	}

	resp = &types.DashboardStatsResp{
		TotalOrders:     totalOrders,
		TotalRevenue:    totalRevenue,
		TodayOrders:     todayOrders,
		TodayRevenue:    todayRevenue,
		WeekOrders:      weekOrders,
		WeekRevenue:     weekRevenue,
		MonthOrders:     monthOrders,
		MonthRevenue:    monthRevenue,
		OrderTrend:      orderTrends,
		RevenueTrend:    []types.TrendData{},
		TopCampaigns:    topCampaigns,
		TopDistributors: topDistributors,
	}

	return resp, nil
}

func dateFilterQuery(req *types.GetDashboardStatsReq) string {
	if req.Period == "week" {
		return "DATE_SUB(NOW(), INTERVAL 7 DAY)"
	} else if req.Period == "month" {
		return "DATE_SUB(NOW(), INTERVAL 30 DAY)"
	} else if req.Period == "year" {
		return "DATE_SUB(NOW(), INTERVAL 365 DAY)"
	}
	return "DATE_SUB(NOW(), INTERVAL 90 DAY)"
}
