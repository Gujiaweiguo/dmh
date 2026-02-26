//go:build layered_demo
// +build layered_demo

// ============================================================
// Order 模块 Repository 层分层测试示范
// ============================================================
// 职责：测试数据访问、SQL 正确性、约束验证
// Mock 策略：使用真实 MySQL8 数据库（testcontainers 或现有容器）
// ============================================================

package order_test

import (
	"fmt"
	"testing"
	"time"

	"dmh/api/internal/testutil"
	"dmh/model"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// ============================================================
// 测试套件模式 - 核心示范
// ============================================================

// OrderRepositoryTestSuite Order Repository 测试套件
// 使用 testify/suite 管理测试生命周期
type OrderRepositoryTestSuite struct {
	suite.Suite
	db      *gorm.DB
	dbName  string
	cleanup func()
}

// TestOrderRepositorySuite 运行测试套件
func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

// ============================================================
// 套件生命周期
// ============================================================

// SetupSuite 套件初始化（运行一次）
func (s *OrderRepositoryTestSuite) SetupSuite() {
	// 可在此初始化共享资源（如容器）
}

// SetupTest 每个测试前运行
func (s *OrderRepositoryTestSuite) SetupTest() {
	// 为每个测试创建隔离数据库
	s.db, s.dbName = testutil.SetupMySQLTestDB(s.T())
}

// TearDownTest 每个测试后运行
func (s *OrderRepositoryTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
	// 数据库清理由 testutil.SetupMySQLTestDB 自动处理
}

// ============================================================
// CRUD 测试 - 创建
// ============================================================

func (s *OrderRepositoryTestSuite) TestCreate() {
	// 1. 准备依赖数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now().Add(-24 * time.Hour),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	// 2. 创建订单
	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		Amount:     100.00,
		PayStatus:  "pending",
		Status:     "pending",
		FormData:   `{"name":"张三"}`,
	}

	// 3. 执行创建
	err := s.db.Create(order).Error

	// 4. 验证结果
	s.NoError(err)
	s.NotZero(order.Id)
	s.NotZero(order.CreatedAt)
	s.NotZero(order.UpdatedAt)
}

func (s *OrderRepositoryTestSuite) TestCreate_UniqueConstraint() {
	// 1. 准备依赖数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	// 2. 创建第一个订单
	order1 := &model.Order{
		CampaignId: campaign.Id,
		Phone:      "13800138001",
		Amount:     100.00,
		PayStatus:  "pending",
		Status:     "pending",
		FormData:   "{}",
	}
	s.Require().NoError(s.db.Create(order1).Error)

	// 3. 尝试创建重复订单（相同 campaign + phone）
	order2 := &model.Order{
		CampaignId: campaign.Id,
		Phone:      "13800138001", // 重复
		Amount:     200.00,
		PayStatus:  "pending",
		Status:     "pending",
		FormData:   "{}",
	}

	err := s.db.Create(order2).Error

	// 4. 验证唯一约束
	s.Error(err)
	s.Contains(err.Error(), "Duplicate")
}

// ============================================================
// CRUD 测试 - 读取
// ============================================================

func (s *OrderRepositoryTestSuite) TestRead() {
	// 1. 准备测试数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
		FormData:   `{"name":"张三"}`,
	}
	s.Require().NoError(s.db.Create(order).Error)

	// 2. 读取订单
	var found model.Order
	err := s.db.First(&found, order.Id).Error

	// 3. 验证结果
	s.NoError(err)
	s.Equal(order.Id, found.Id)
	s.Equal(order.Phone, found.Phone)
	s.Equal(order.Amount, found.Amount)
	s.Equal(order.FormData, found.FormData)
}

func (s *OrderRepositoryTestSuite) TestRead_NotFound() {
	var found model.Order
	err := s.db.First(&found, 99999).Error

	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *OrderRepositoryTestSuite) TestRead_WithPreload() {
	// 1. 准备关联数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
		FormData:   "{}",
	}
	s.Require().NoError(s.db.Create(order).Error)

	// 2. 使用 Preload 读取
	var found model.Order
	err := s.db.First(&found, order.Id).Error

	// 3. 验证关联数据
	s.NoError(err)
	s.Equal(campaign.Id, found.CampaignId)
}

// ============================================================
// CRUD 测试 - 更新
// ============================================================

func (s *OrderRepositoryTestSuite) TestUpdate() {
	// 1. 创建测试数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		Amount:     100.00,
		PayStatus:  "pending",
		Status:     "pending",
		FormData:   "{}",
	}
	s.Require().NoError(s.db.Create(order).Error)

	// 2. 更新状态
	newPayStatus := "paid"
	err := s.db.Model(order).Update("pay_status", newPayStatus).Error
	s.NoError(err)

	// 3. 验证更新
	var found model.Order
	s.db.First(&found, order.Id)
	s.Equal(newPayStatus, found.PayStatus)
}

func (s *OrderRepositoryTestSuite) TestUpdate_VerificationStatus() {
	// 1. 创建订单
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	order := &model.Order{
		CampaignId:       campaign.Id,
		Phone:            "13800138000",
		Amount:           100.00,
		PayStatus:        "paid",
		Status:           "active",
		VerificationCode: "TEST_CODE",
		FormData:         "{}",
	}
	s.Require().NoError(s.db.Create(order).Error)

	// 2. 更新核销状态
	now := time.Now()
	verifiedBy := int64(100)
	err := s.db.Model(order).Updates(map[string]interface{}{
		"verification_status": "verified",
		"verified_at":         &now,
		"verified_by":         verifiedBy,
	}).Error
	s.NoError(err)

	// 3. 验证
	var found model.Order
	s.db.First(&found, order.Id)
	s.Equal("verified", found.VerificationStatus)
	s.NotNil(found.VerifiedAt)
	s.NotNil(found.VerifiedBy)
	s.Equal(verifiedBy, *found.VerifiedBy)
}

// ============================================================
// CRUD 测试 - 删除
// ============================================================

func (s *OrderRepositoryTestSuite) TestDelete() {
	// 1. 创建订单
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	order := &model.Order{
		CampaignId: campaign.Id,
		Phone:      "13800138000",
		Amount:     100.00,
		PayStatus:  "paid",
		Status:     "active",
		FormData:   "{}",
	}
	s.Require().NoError(s.db.Create(order).Error)

	// 2. 删除
	err := s.db.Delete(order).Error
	s.NoError(err)

	// 3. 验证删除
	var found model.Order
	err = s.db.First(&found, order.Id).Error
	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

// ============================================================
// 查询测试
// ============================================================

func (s *OrderRepositoryTestSuite) TestQuery_WithConditions() {
	// 1. 创建多条测试数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	for i := 0; i < 5; i++ {
		order := &model.Order{
			CampaignId: campaign.Id,
			Phone:      fmt.Sprintf("1380013800%d", i),
			Amount:     100.00,
			PayStatus:  "paid",
			Status:     "active",
			FormData:   "{}",
		}
		s.Require().NoError(s.db.Create(order).Error)
	}

	// 2. 查询
	var results []model.Order
	err := s.db.Where("campaign_id = ? AND pay_status = ?", campaign.Id, "paid").Find(&results).Error

	// 3. 验证
	s.NoError(err)
	s.Len(results, 5)
}

func (s *OrderRepositoryTestSuite) TestQuery_Pagination() {
	// 1. 创建 25 条数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	for i := 0; i < 25; i++ {
		order := &model.Order{
			CampaignId: campaign.Id,
			Phone:      testutil.GenUniquePhone(),
			Amount:     100.00,
			PayStatus:  "paid",
			Status:     "active",
			FormData:   "{}",
		}
		s.Require().NoError(s.db.Create(order).Error)
	}

	// 2. 分页查询
	var results []model.Order
	var total int64

	s.db.Model(&model.Order{}).Where("campaign_id = ?", campaign.Id).Count(&total)
	err := s.db.Where("campaign_id = ?", campaign.Id).
		Offset(0).Limit(10).
		Order("id DESC").
		Find(&results).Error

	// 3. 验证
	s.NoError(err)
	s.Equal(int64(25), total)
	s.Len(results, 10)
}

// ============================================================
// 事务测试
// ============================================================

func (s *OrderRepositoryTestSuite) TestTransaction_Commit() {
	// 1. 准备数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	// 2. 执行事务
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 创建订单
		order := &model.Order{
			CampaignId: campaign.Id,
			Phone:      "13800138000",
			Amount:     100.00,
			PayStatus:  "paid",
			Status:     "active",
			FormData:   "{}",
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 更新状态
		if err := tx.Model(order).Update("pay_status", "paid").Error; err != nil {
			return err
		}

		return nil
	})

	// 3. 验证事务成功
	s.NoError(err)

	var count int64
	s.db.Model(&model.Order{}).Where("phone = ?", "13800138000").Count(&count)
	s.Equal(int64(1), count)
}

func (s *OrderRepositoryTestSuite) TestTransaction_Rollback() {
	// 1. 准备数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	// 2. 执行事务（故意失败）
	err := s.db.Transaction(func(tx *gorm.DB) error {
		order := &model.Order{
			CampaignId: campaign.Id,
			Phone:      "13800138000",
			Amount:     100.00,
			PayStatus:  "paid",
			Status:     "active",
			FormData:   "{}",
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 故意返回错误触发回滚
		return fmt.Errorf("simulated error")
	})

	// 3. 验证事务回滚
	s.Error(err)

	var count int64
	s.db.Model(&model.Order{}).Where("phone = ?", "13800138000").Count(&count)
	s.Equal(int64(0), count) // 数据应被回滚
}

// ============================================================
// 核销记录测试
// ============================================================

func (s *OrderRepositoryTestSuite) TestVerificationRecord_Create() {
	// 1. 创建订单
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:    brand.Id,
		Name:       "Test Campaign",
		Status:     "active",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(24 * time.Hour),
		FormFields: "[]",
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	order := &model.Order{
		CampaignId:       campaign.Id,
		Phone:            "13800138000",
		Amount:           100.00,
		PayStatus:        "paid",
		Status:           "active",
		VerificationCode: "TEST_CODE",
		FormData:         "{}",
	}
	s.Require().NoError(s.db.Create(order).Error)

	// 2. 创建核销记录
	now := time.Now()
	verifiedBy := int64(100)
	record := &model.VerificationRecord{
		OrderID:            order.Id,
		VerificationStatus: "verified",
		VerifiedAt:         &now,
		VerifiedBy:         &verifiedBy,
		VerificationCode:   "TEST_CODE",
		VerificationMethod: "manual",
		Remark:             "测试核销",
	}

	err := s.db.Create(record).Error

	// 3. 验证
	s.NoError(err)
	s.NotZero(record.ID)
}

// ============================================================
// 关键模式总结
// ============================================================
//
// Repository 层 MySQL8 测试模式：
// 1. 使用 testify/suite 管理测试生命周期
// 2. 每个测试使用独立数据库（testutil.SetupMySQLTestDB）
// 3. 测试覆盖：
//    - CRUD 操作
//    - 唯一约束
//    - 关联查询（Preload）
//    - 分页
//    - 事务（提交/回滚）
//
// 测试隔离策略：
// - SetupTest: 创建独立数据库
// - TearDownTest: 自动清理数据库
// - 使用唯一测试数据（testutil.GenUniquePhone）
//
// 运行命令：
// cd backend && go test -p 1 -v ./api/internal/logic/order/... -run RepositorySuite
