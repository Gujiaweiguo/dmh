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

func (l *GetDashboardStatsLogic) GetDashboardStats(req *types.GetDashboardStatsReq) (resp *types.DashboardStatsResp, err error) {
	if req.BrandId <= 0 {
		return nil, errors.New("brandId is required")
	}

	baseQuery := l.svcCtx.DB.Model(&model.Order{})

	if req.StartDate != "" {
		baseQuery = baseQuery.Where("created_at >= ?", req.StartDate)
	}

	if req.EndDate != "" {
		baseQuery = baseQuery.Where("created_at <= ?", req.EndDate)
	}

	if req.BrandId > 0 {
		baseQuery = baseQuery.Where("brand_id = ?", req.BrandId)
	}

	var totalOrders int64
	var totalRevenue float64
	if err := baseQuery.Select("COUNT(*)").Scan(&totalOrders).Error; err != nil {
		return nil, err
	}

	if err := baseQuery.Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}

	now := l.svcCtx.DB.NowFunc()
	todayStart := now.Format("2006-01-02")
	weekStart := now.AddDate(-7 * 24 * time.Hour).Format("2006-01-02")
	monthStart := now.Format("2006-01-01")

	var todayOrders int64
	var todayRevenue float64
	if err := baseQuery.Where("DATE(created_at) = ?", todayStart).
		Select("COUNT(*)").Scan(&todayOrders).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Where("DATE(created_at) = ?", todayStart).
		Select("COALESCE(SUM(amount), 0)").Scan(&todayRevenue).Error; err != nil {
		return nil, err
	}

	var weekOrders int64
	var weekRevenue float64
	if err := baseQuery.Where("created_at >= ?", weekStart).
		Select("COUNT(*)").Scan(&weekOrders).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Where("created_at >= ?", weekStart).
		Select("COALESCE(SUM(amount), 0)").Scan(&weekRevenue).Error; err != nil {
		return nil, err
	}

	var monthOrders int64
	var monthRevenue float64
	if err := baseQuery.Where("created_at >= ?", monthStart).
		Select("COUNT(*)").Scan(&monthOrders).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Where("created_at >= ?", monthStart).
		Select("COALESCE(SUM(amount), 0)").Scan(&monthRevenue).Error; err != nil {
		return nil, err
	}

	var orderTrends []types.TrendData
	orderQuery := l.svcCtx.DB.Model(&model.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as value")

	if req.BrandId > 0 {
		orderQuery = orderQuery.Where("brand_id = ?", req.BrandId)
	}

	dateFilter := "7 DAY"
	if req.Period == "week" {
		dateFilter = "30 DAY"
	} else if req.Period == "month" {
		dateFilter = "90 DAY"
	}

	if err := orderQuery.Where("created_at >= ?", dateFilterQuery(req)).
		Group("DATE(created_at)").
		Order("DATE(created_at) ASC").
		Scan(&orderTrends).Error; err != nil {
		return nil, err
	}

	var topCampaigns []types.CampaignStats
	campaignQuery := l.svcCtx.DB.Model(&model.Campaign{}).
		Select("id, name, (SELECT COUNT(*) FROM orders WHERE campaign_id = campaigns.id) as order_count, (SELECT COALESCE(SUM(amount), 0) FROM orders WHERE campaign_id = campaigns.id) as revenue")

	if req.BrandId > 0 {
		campaignQuery = campaignQuery.Where("brand_id = ?", req.BrandId)
	}

	if err := campaignQuery.Where("created_at >= ?", dateFilterQuery(req)).
		Order("(SELECT COUNT(*) FROM orders WHERE campaign_id = campaigns.id) DESC").
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
		return nil, err
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
