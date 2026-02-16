//go:build mysql
// +build mysql

package statistics

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMySQLTestDB(t *testing.T) *gorm.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "#Admin168"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "dmh_test"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("无法连接到 MySQL: %v", err)
	}

	// Auto migrate
	db.AutoMigrate(&model.Order{}, &model.Campaign{}, &model.Distributor{}, &model.Brand{}, &model.User{})

	return db
}

func TestGetDashboardStatsLogic_WithMySQL(t *testing.T) {
	db := setupMySQLTestDB(t)

	// Clean up test data
	defer func() {
		db.Exec("DELETE FROM orders WHERE campaign_id = 99999")
		db.Exec("DELETE FROM campaigns WHERE id = 99999")
		db.Exec("DELETE FROM brands WHERE id = 99999")
	}()

	// Create test data
	brand := &model.Brand{
		Id:   99999,
		Name: "Test Brand",
	}
	db.Create(brand)

	now := time.Now()
	campaign := &model.Campaign{
		Id:         99999,
		BrandId:    99999,
		Name:       "Test Campaign",
		Status:     "active",
		FormFields: "[]",
		StartTime:  now.Add(-24 * time.Hour),
		EndTime:    now.Add(24 * time.Hour),
	}
	db.Create(campaign)

	order := &model.Order{
		CampaignId: 99999,
		Amount:     100.00,
		Status:     "completed",
		PayStatus:  "paid",
		FormData:   "{}",
	}
	db.Create(order)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetDashboardStatsLogic(ctx, svcCtx)

	req := &types.GetDashboardStatsReq{
		BrandId: 99999,
		Period:  "month",
	}

	resp, err := logic.GetDashboardStats(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.GreaterOrEqual(t, resp.TotalOrders, int64(1))
	assert.GreaterOrEqual(t, resp.TotalRevenue, float64(100))
}
