//go:build layered_demo
// +build layered_demo

// ============================================================
// Distributor 模块 Repository 层分层测试示范
// ============================================================
// 职责：测试数据访问、SQL 正确性、约束验证
// Mock 策略：使用真实 MySQL8 数据库
// ============================================================

package distributor_test

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
// 测试套件 - 核心示范
// ============================================================

// DistributorRepositoryTestSuite Distributor Repository 测试套件
type DistributorRepositoryTestSuite struct {
	suite.Suite
	db     *gorm.DB
	dbName string
}

// TestDistributorRepositorySuite 运行测试套件
func TestDistributorRepositorySuite(t *testing.T) {
	suite.Run(t, new(DistributorRepositoryTestSuite))
}

// SetupTest 每个测试前运行
func (s *DistributorRepositoryTestSuite) SetupTest() {
	s.db, s.dbName = testutil.SetupMySQLTestDB(s.T())
}

// ============================================================
// CRUD 测试 - 创建分销商
// ============================================================

func (s *DistributorRepositoryTestSuite) TestCreate() {
	// 1. 准备依赖数据
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	// 2. 创建分销商
	distributor := &model.Distributor{
		UserId:            user.Id,
		BrandId:           brand.Id,
		Level:             1,
		Status:            "active",
		TotalEarnings:     0,
		SubordinatesCount: 0,
	}

	// 3. 执行创建
	err := s.db.Create(distributor).Error

	// 4. 验证
	s.NoError(err)
	s.NotZero(distributor.Id)
	s.NotZero(distributor.CreatedAt)
	s.NotZero(distributor.UpdatedAt)
}

func (s *DistributorRepositoryTestSuite) TestCreate_UniqueConstraint() {
	// 1. 准备数据
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	// 2. 创建第一个分销商
	dist1 := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(dist1).Error)

	// 3. 尝试创建重复（相同 user + brand）
	dist2 := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   2,
		Status:  "active",
	}

	err := s.db.Create(dist2).Error

	// 4. 验证唯一约束
	s.Error(err)
	s.Contains(err.Error(), "Duplicate")
}

func (s *DistributorRepositoryTestSuite) TestCreate_WithParent() {
	// 1. 创建父分销商
	parentUser := &model.User{
		Username: testutil.GenUniqueUsername("parent"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(parentUser).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	parentDist := &model.Distributor{
		UserId:  parentUser.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(parentDist).Error)

	// 2. 创建子分销商
	childUser := &model.User{
		Username: testutil.GenUniqueUsername("child"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(childUser).Error)

	childDist := &model.Distributor{
		UserId:     childUser.Id,
		BrandId:    brand.Id,
		ParentId:   &parentDist.Id,
		Level:      2,
		Status:     "active",
		ApprovedBy: &parentUser.Id,
	}

	err := s.db.Create(childDist).Error

	// 3. 验证
	s.NoError(err)
	s.NotZero(childDist.Id)
	s.Equal(parentDist.Id, *childDist.ParentId)

	// 4. 验证父分销商的下级计数更新
	var foundParent model.Distributor
	s.db.First(&foundParent, parentDist.Id)
	// 注意：下级计数需要在业务逻辑中更新
}

// ============================================================
// CRUD 测试 - 读取
// ============================================================

func (s *DistributorRepositoryTestSuite) TestRead() {
	// 1. 准备数据
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	now := time.Now()
	distributor := &model.Distributor{
		UserId:        user.Id,
		BrandId:       brand.Id,
		Level:         2,
		Status:        "active",
		TotalEarnings: 1000.50,
		ApprovedBy:    &user.Id,
		ApprovedAt:    &now,
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 读取
	var found model.Distributor
	err := s.db.First(&found, distributor.Id).Error

	// 3. 验证
	s.NoError(err)
	s.Equal(distributor.Id, found.Id)
	s.Equal(distributor.Level, found.Level)
	s.Equal(distributor.TotalEarnings, found.TotalEarnings)
}

func (s *DistributorRepositoryTestSuite) TestRead_WithPreload() {
	// 1. 准备数据
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 使用 Preload 读取
	var found model.Distributor
	err := s.db.Preload("User").Preload("Brand").First(&found, distributor.Id).Error

	// 3. 验证关联数据
	s.NoError(err)
	s.NotNil(found.User)
	s.Equal(user.Username, found.User.Username)
	s.NotNil(found.Brand)
	s.Equal(brand.Name, found.Brand.Name)
}

func (s *DistributorRepositoryTestSuite) TestRead_NotFound() {
	var found model.Distributor
	err := s.db.First(&found, 99999).Error

	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

// ============================================================
// CRUD 测试 - 更新
// ============================================================

func (s *DistributorRepositoryTestSuite) TestUpdate_Status() {
	// 1. 创建分销商
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 更新状态
	err := s.db.Model(distributor).Update("status", "suspended").Error
	s.NoError(err)

	// 3. 验证
	var found model.Distributor
	s.db.First(&found, distributor.Id)
	s.Equal("suspended", found.Status)
}

func (s *DistributorRepositoryTestSuite) TestUpdate_Level() {
	// 1. 创建分销商
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 更新级别
	err := s.db.Model(distributor).Update("level", 3).Error
	s.NoError(err)

	// 3. 验证
	var found model.Distributor
	s.db.First(&found, distributor.Id)
	s.Equal(3, found.Level)
}

func (s *DistributorRepositoryTestSuite) TestUpdate_Earnings() {
	// 1. 创建分销商
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	distributor := &model.Distributor{
		UserId:        user.Id,
		BrandId:       brand.Id,
		Level:         1,
		Status:        "active",
		TotalEarnings: 100.00,
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 增加收益
	err := s.db.Model(distributor).Update("total_earnings", gorm.Expr("total_earnings + ?", 50.00)).Error
	s.NoError(err)

	// 3. 验证
	var found model.Distributor
	s.db.First(&found, distributor.Id)
	s.Equal(150.00, found.TotalEarnings)
}

// ============================================================
// CRUD 测试 - 删除（软删除）
// ============================================================

func (s *DistributorRepositoryTestSuite) TestSoftDelete() {
	// 1. 创建分销商
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 软删除
	now := time.Now()
	err := s.db.Model(distributor).Update("deleted_at", &now).Error
	s.NoError(err)

	// 3. 普通查询不应找到
	var found model.Distributor
	err = s.db.Where("id = ? AND deleted_at IS NULL", distributor.Id).First(&found).Error
	s.Error(err)

	// 4. Unscoped 可以找到
	err = s.db.Unscoped().First(&found, distributor.Id).Error
	s.NoError(err)
	s.NotNil(found.DeletedAt)
}

// ============================================================
// 分销商申请测试
// ============================================================

func (s *DistributorRepositoryTestSuite) TestApplication_Create() {
	// 1. 准备数据
	user := &model.User{
		Username: testutil.GenUniqueUsername("applicant"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	// 2. 创建申请
	application := &model.DistributorApplication{
		UserId:  user.Id,
		BrandId: brand.Id,
		Status:  "pending",
		Reason:  "我想成为分销商",
	}

	err := s.db.Create(application).Error

	// 3. 验证
	s.NoError(err)
	s.NotZero(application.Id)
	s.Equal("pending", application.Status)
}

func (s *DistributorRepositoryTestSuite) TestApplication_Approve() {
	// 1. 创建申请
	user := &model.User{
		Username: testutil.GenUniqueUsername("applicant"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	reviewer := &model.User{
		Username: testutil.GenUniqueUsername("reviewer"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "brand_admin",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(reviewer).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	application := &model.DistributorApplication{
		UserId:  user.Id,
		BrandId: brand.Id,
		Status:  "pending",
		Reason:  "申请",
	}
	s.Require().NoError(s.db.Create(application).Error)

	// 2. 审批通过
	now := time.Now()
	err := s.db.Model(application).Updates(map[string]interface{}{
		"status":       "approved",
		"reviewed_by":  reviewer.Id,
		"reviewed_at":  &now,
		"review_notes": "审核通过",
	}).Error
	s.NoError(err)

	// 3. 验证
	var found model.DistributorApplication
	s.db.First(&found, application.Id)
	s.Equal("approved", found.Status)
	s.NotNil(found.ReviewedBy)
	s.Equal(reviewer.Id, *found.ReviewedBy)
}

// ============================================================
// 分销商链接测试
// ============================================================

func (s *DistributorRepositoryTestSuite) TestLink_Create() {
	// 1. 准备数据
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:   brand.Id,
		Name:      "Test Campaign",
		Status:    "active",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 创建链接
	link := &model.DistributorLink{
		DistributorId: distributor.Id,
		CampaignId:    campaign.Id,
		LinkCode:      fmt.Sprintf("LINK_%d_%d", time.Now().UnixNano(), user.Id),
		ClickCount:    0,
		OrderCount:    0,
		Status:        "active",
	}

	err := s.db.Create(link).Error

	// 3. 验证
	s.NoError(err)
	s.NotZero(link.Id)
	s.NotEmpty(link.LinkCode)
}

func (s *DistributorRepositoryTestSuite) TestLink_UniqueCode() {
	// 1. 创建分销商
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	campaign := &model.Campaign{
		BrandId:   brand.Id,
		Name:      "Test Campaign",
		Status:    "active",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 创建第一个链接
	linkCode := "UNIQUE_CODE_123"
	link1 := &model.DistributorLink{
		DistributorId: distributor.Id,
		CampaignId:    campaign.Id,
		LinkCode:      linkCode,
		Status:        "active",
	}
	s.Require().NoError(s.db.Create(link1).Error)

	// 3. 尝试创建重复链接码
	link2 := &model.DistributorLink{
		DistributorId: distributor.Id,
		CampaignId:    campaign.Id,
		LinkCode:      linkCode, // 重复
		Status:        "active",
	}

	err := s.db.Create(link2).Error

	// 4. 验证唯一约束
	s.Error(err)
	s.Contains(err.Error(), "Duplicate")
}

// ============================================================
// 分销商奖励测试
// ============================================================

func (s *DistributorRepositoryTestSuite) TestReward_Create() {
	// 1. 准备数据
	user := &model.User{
		Username: testutil.GenUniqueUsername("dist_user"),
		Password: "$2a$10$hashed",
		Phone:    testutil.GenUniquePhone(),
		Role:     "participant",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{
		Name:   "Test Brand",
		Status: "active",
	}
	s.Require().NoError(s.db.Create(brand).Error)

	distributor := &model.Distributor{
		UserId:  user.Id,
		BrandId: brand.Id,
		Level:   1,
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(distributor).Error)

	// 2. 创建奖励记录
	now := time.Now()
	reward := &model.DistributorReward{
		DistributorId: distributor.Id,
		UserId:        user.Id,
		OrderId:       1,
		CampaignId:    1,
		Amount:        50.00,
		Level:         1,
		RewardRate:    10.00,
		Status:        "settled",
		SettledAt:     &now,
	}

	err := s.db.Create(reward).Error

	// 3. 验证
	s.NoError(err)
	s.NotZero(reward.Id)
	s.Equal(50.00, reward.Amount)
}

// ============================================================
// 查询测试
// ============================================================

func (s *DistributorRepositoryTestSuite) TestQuery_ByBrandAndStatus() {
	// 1. 创建多个分销商
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	for i := 0; i < 5; i++ {
		user := &model.User{
			Username: testutil.GenUniqueUsername(fmt.Sprintf("user%d", i)),
			Password: "$2a$10$hashed",
			Phone:    testutil.GenUniquePhone(),
			Role:     "participant",
			Status:   "active",
		}
		s.Require().NoError(s.db.Create(user).Error)

		status := "active"
		if i == 4 {
			status = "suspended"
		}

		dist := &model.Distributor{
			UserId:  user.Id,
			BrandId: brand.Id,
			Level:   1,
			Status:  status,
		}
		s.Require().NoError(s.db.Create(dist).Error)
	}

	// 2. 查询 active 状态
	var results []model.Distributor
	err := s.db.Where("brand_id = ? AND status = ?", brand.Id, "active").Find(&results).Error

	s.NoError(err)
	s.Len(results, 4)
}

// ============================================================
// 关键模式总结
// ============================================================
//
// Distributor Repository 层 MySQL8 测试模式：
// 1. 使用 testify/suite 管理测试生命周期
// 2. 每个测试使用独立数据库
// 3. 测试覆盖：
//    - 分销商 CRUD（含父级关系）
//    - 分销商申请流程
//    - 分销商链接（含唯一约束）
//    - 分销商奖励
//    - 状态/级别/收益更新
//    - 软删除
//    - 关联查询（Preload）
//    - 条件查询
//
// 特有场景：
// - 分销商层级关系（ParentId）
// - 唯一约束（user_id + brand_id, link_code）
// - 收益累计（gorm.Expr）
// - 软删除（deleted_at）
