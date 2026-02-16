package statistics

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	return gormDB, mock
}

func TestBuildOrderBaseQuery_WithBrandId(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.DB()

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId:   123,
		Period:    "month",
		StartDate: "2024-01-01",
		EndDate:   "2024-12-31",
	}

	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND orders.created_at <= ? AND campaigns.brand_id = ?")).
		WithArgs("2024-01-01", "2024-12-31", int64(123)).
		WillReturnRows(rows)

	var count int64
	err := logic.buildOrderBaseQuery(req).Select("COUNT(*)").Scan(&count).Error

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBuildOrderBaseQuery_WithDateFilters(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.DB()

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId:   0,
		StartDate: "2024-01-01",
		EndDate:   "2024-12-31",
	}

	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(5)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` WHERE orders.created_at >= ? AND orders.created_at <= ?")).
		WithArgs("2024-01-01", "2024-12-31").
		WillReturnRows(rows)

	var count int64
	err := logic.buildOrderBaseQuery(req).Select("COUNT(*)").Scan(&count).Error

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDashboardStats_DBError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.DB()

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId: 123,
		Period:  "month",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ?")).
		WithArgs(int64(123)).
		WillReturnError(sql.ErrConnDone)

	resp, err := logic.GetDashboardStats(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDashboardStats_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.DB()

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId: 123,
		Period:  "month",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ?")).
		WithArgs(int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(100))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ?")).
		WithArgs(int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(5000.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND DATE(orders.created_at) = ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(10))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND DATE(orders.created_at) = ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(500.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(30))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(1500.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(80))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(4000.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE(orders.created_at) as date, COUNT(*) as value FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ? GROUP BY DATE(orders.created_at) ORDER BY DATE(orders.created_at) ASC")).
		WithArgs(int64(123), "DATE_SUB(NOW(), INTERVAL 90 DAY)").
		WillReturnRows(sqlmock.NewRows([]string{"date", "value"}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT campaigns.id, campaigns.name, (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) as order_count, (SELECT COALESCE(SUM(amount), 0) FROM orders WHERE orders.campaign_id = campaigns.id) as revenue FROM `campaigns` WHERE campaigns.brand_id = ? AND campaigns.created_at >= ? ORDER BY (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) DESC LIMIT ?")).
		WithArgs(int64(123), "DATE_SUB(NOW(), INTERVAL 30 DAY)", 5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "order_count", "revenue"}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT distributors.id, users.username, distributors.level, distributors.total_earnings, distributors.subordinates_count FROM `distributors` LEFT JOIN users ON users.id = distributors.user_id WHERE distributors.brand_id = ? ORDER BY distributors.total_earnings DESC LIMIT ?")).
		WithArgs(int64(123), 5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "level", "total_earnings", "subordinates_count"}))

	resp, err := logic.GetDashboardStats(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(100), resp.TotalOrders)
	assert.Equal(t, float64(5000), resp.TotalRevenue)
	assert.Equal(t, int64(10), resp.TodayOrders)
	assert.Equal(t, float64(500), resp.TodayRevenue)
	assert.Equal(t, int64(30), resp.WeekOrders)
	assert.Equal(t, float64(1500), resp.WeekRevenue)
	assert.Equal(t, int64(80), resp.MonthOrders)
	assert.Equal(t, float64(4000), resp.MonthRevenue)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDashboardStats_WithStartDateFilter(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.DB()

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId:   123,
		Period:    "month",
		StartDate: "2024-01-01",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND campaigns.brand_id = ?")).
		WithArgs("2024-01-01", int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(50))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND campaigns.brand_id = ?")).
		WithArgs("2024-01-01", int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(2500.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND campaigns.brand_id = ? AND DATE(orders.created_at) = ?")).
		WithArgs("2024-01-01", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(5))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND campaigns.brand_id = ? AND DATE(orders.created_at) = ?")).
		WithArgs("2024-01-01", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(250.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs("2024-01-01", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(15))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs("2024-01-01", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(750.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs("2024-01-01", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(40))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at >= ? AND campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs("2024-01-01", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(2000.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE(orders.created_at) as date, COUNT(*) as value FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ? GROUP BY DATE(orders.created_at) ORDER BY DATE(orders.created_at) ASC")).
		WithArgs(int64(123), "DATE_SUB(NOW(), INTERVAL 90 DAY)").
		WillReturnRows(sqlmock.NewRows([]string{"date", "value"}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT campaigns.id, campaigns.name, (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) as order_count, (SELECT COALESCE(SUM(amount), 0) FROM orders WHERE orders.campaign_id = campaigns.id) as revenue FROM `campaigns` WHERE campaigns.brand_id = ? AND campaigns.created_at >= ? ORDER BY (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) DESC LIMIT ?")).
		WithArgs(int64(123), "DATE_SUB(NOW(), INTERVAL 30 DAY)", 5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "order_count", "revenue"}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT distributors.id, users.username, distributors.level, distributors.total_earnings, distributors.subordinates_count FROM `distributors` LEFT JOIN users ON users.id = distributors.user_id WHERE distributors.brand_id = ? ORDER BY distributors.total_earnings DESC LIMIT ?")).
		WithArgs(int64(123), 5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "level", "total_earnings", "subordinates_count"}))

	resp, err := logic.GetDashboardStats(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(50), resp.TotalOrders)
	assert.Equal(t, float64(2500), resp.TotalRevenue)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDashboardStats_DistributorError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.DB()

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId: 123,
		Period:  "month",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ?")).
		WithArgs(int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(100))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ?")).
		WithArgs(int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(5000.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND DATE(orders.created_at) = ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(10))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND DATE(orders.created_at) = ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(500.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(30))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(1500.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(80))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs(int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(4000.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE(orders.created_at) as date, COUNT(*) as value FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ? GROUP BY DATE(orders.created_at) ORDER BY DATE(orders.created_at) ASC")).
		WithArgs(int64(123), "DATE_SUB(NOW(), INTERVAL 90 DAY)").
		WillReturnRows(sqlmock.NewRows([]string{"date", "value"}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT campaigns.id, campaigns.name, (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) as order_count, (SELECT COALESCE(SUM(amount), 0) FROM orders WHERE orders.campaign_id = campaigns.id) as revenue FROM `campaigns` WHERE campaigns.brand_id = ? AND campaigns.created_at >= ? ORDER BY (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) DESC LIMIT ?")).
		WithArgs(int64(123), "DATE_SUB(NOW(), INTERVAL 30 DAY)", 5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "order_count", "revenue"}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT distributors.id, users.username, distributors.level, distributors.total_earnings, distributors.subordinates_count FROM `distributors` LEFT JOIN users ON users.id = distributors.user_id WHERE distributors.brand_id = ? ORDER BY distributors.total_earnings DESC LIMIT ?")).
		WithArgs(int64(123), 5).
		WillReturnError(sql.ErrConnDone)

	resp, err := logic.GetDashboardStats(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.TopDistributors)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDashboardStats_WithEndDateFilter(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.DB()

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId: 123,
		Period:  "month",
		EndDate: "2024-12-31",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at <= ? AND campaigns.brand_id = ?")).
		WithArgs("2024-12-31", int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(75))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at <= ? AND campaigns.brand_id = ?")).
		WithArgs("2024-12-31", int64(123)).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(3750.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at <= ? AND campaigns.brand_id = ? AND DATE(orders.created_at) = ?")).
		WithArgs("2024-12-31", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(8))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at <= ? AND campaigns.brand_id = ? AND DATE(orders.created_at) = ?")).
		WithArgs("2024-12-31", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(400.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at <= ? AND campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs("2024-12-31", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(25))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at <= ? AND campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs("2024-12-31", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(1250.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at <= ? AND campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs("2024-12-31", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(60))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(orders.amount), 0) FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE orders.created_at <= ? AND campaigns.brand_id = ? AND orders.created_at >= ?")).
		WithArgs("2024-12-31", int64(123), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(3000.00))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE(orders.created_at) as date, COUNT(*) as value FROM `orders` JOIN campaigns ON campaigns.id = orders.campaign_id WHERE campaigns.brand_id = ? AND orders.created_at >= ? GROUP BY DATE(orders.created_at) ORDER BY DATE(orders.created_at) ASC")).
		WithArgs(int64(123), "DATE_SUB(NOW(), INTERVAL 90 DAY)").
		WillReturnRows(sqlmock.NewRows([]string{"date", "value"}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT campaigns.id, campaigns.name, (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) as order_count, (SELECT COALESCE(SUM(amount), 0) FROM orders WHERE orders.campaign_id = campaigns.id) as revenue FROM `campaigns` WHERE campaigns.brand_id = ? AND campaigns.created_at >= ? ORDER BY (SELECT COUNT(*) FROM orders WHERE orders.campaign_id = campaigns.id) DESC LIMIT ?")).
		WithArgs(int64(123), "DATE_SUB(NOW(), INTERVAL 30 DAY)", 5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "order_count", "revenue"}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT distributors.id, users.username, distributors.level, distributors.total_earnings, distributors.subordinates_count FROM `distributors` LEFT JOIN users ON users.id = distributors.user_id WHERE distributors.brand_id = ? ORDER BY distributors.total_earnings DESC LIMIT ?")).
		WithArgs(int64(123), 5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "level", "total_earnings", "subordinates_count"}))

	resp, err := logic.GetDashboardStats(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(75), resp.TotalOrders)
	assert.Equal(t, float64(3750), resp.TotalRevenue)
	assert.NoError(t, mock.ExpectationsWereMet())
}
