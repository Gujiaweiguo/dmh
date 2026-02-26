# 后端分层测试模板 (T6)

## 概述

本文档提供 DMH 后端三层测试模板，遵循以下分层职责：

| 层级 | Mock 策略 | 测试重点 |
|------|----------|----------|
| Handler | Mock Logic | 请求解析、响应组装、HTTP 状态码 |
| Logic | Mock DB (sqlmock) | 业务逻辑、数据验证、错误处理 |
| Repository | Real MySQL8 | 数据访问、SQL 正确性、约束验证 |

**关键约束**：
- 使用 `go test -p 1` 避免数据库冲突
- Handler 测试不包含业务逻辑
- Logic 测试使用 sqlmock 隔离数据库
- Repository 测试使用真实 MySQL8（testcontainers 或现有容器）

---

## 1. Handler 层测试模板

### 1.1 职责边界

Handler 层**只负责**：
- HTTP 请求解析 (`httpx.Parse`)
- 调用 Logic 层方法
- 返回 HTTP 响应 (`httpx.OkJsonCtx` / `httpx.ErrorCtx`)

Handler 层**不包含**：
- 业务逻辑
- 数据验证
- 数据库操作

### 1.2 完整模板

```go
// backend/api/internal/handler/<module>/<handler>_test.go
package <module>

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// 测试工具函数
// ============================================================

// setupHandlerTestEnv 创建 Handler 测试环境
// Handler 测试应该 mock Logic 层，不依赖真实数据库
func setupHandlerTestEnv(t *testing.T) (*svc.ServiceContext, *Mock<Module>Logic) {
	t.Helper()

	// 创建 mock logic
	mockLogic := NewMock<Module>Logic()

	// 创建最小的 ServiceContext（Handler 只需要能创建 Logic）
	svcCtx := &svc.ServiceContext{
		// Handler 测试通常不需要完整配置
		Config: config.Config{
			Auth: config.AuthConfig{
				AccessSecret: "test-secret",
				AccessExpire: 3600,
			},
		},
	}

	return svcCtx, mockLogic
}

// Mock<Module>Logic 是 Logic 层的 mock 实现
type Mock<Module>Logic struct {
	// 控制测试行为
	ShouldError bool
	ErrorMsg    string
	// 返回数据
	ReturnData *types.<Response>Type
}

func NewMock<Module>Logic() *Mock<Module>Logic {
	return &Mock<Module>Logic{
		ShouldError: false,
		ReturnData:  &types.<Response>Type{},
	}
}

// 实现 Logic 接口方法
func (m *Mock<Module>Logic) <Action>(req *types.<Request>Type) (*types.<Response>Type, error) {
	if m.ShouldError {
		return nil, fmt.Errorf(m.ErrorMsg)
	}
	return m.ReturnData, nil
}

// ============================================================
// 测试用例模板
// ============================================================

func Test<Handler>_Success(t *testing.T) {
	svcCtx, mockLogic := setupHandlerTestEnv(t)

	// 配置 mock 返回
	mockLogic.ReturnData = &types.<Response>Type{
		// 填充期望的返回数据
		Id:   1,
		Name: "test",
	}

	// 构造请求
	reqBody := types.<Request>Type{
		// 填充请求参数
		Name: "test",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/<module>/<action>", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// 如果需要认证上下文
	ctx := context.WithValue(req.Context(), "userId", int64(1))
	req = req.WithContext(ctx)

	resp := httptest.NewRecorder()

	// 执行 Handler
	handler := <Handler>Handler(svcCtx)
	handler(resp, req)

	// 验证结果
	assert.Equal(t, http.StatusOK, resp.Code)

	var result types.<Response>Type
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, mockLogic.ReturnData.Id, result.Id)
}

func Test<Handler>_InvalidJSON(t *testing.T) {
	svcCtx, _ := setupHandlerTestEnv(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/<module>/<action>", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler := <Handler>Handler(svcCtx)
	handler(resp, req)

	// JSON 解析失败应返回错误
	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func Test<Handler>_LogicError(t *testing.T) {
	svcCtx, mockLogic := setupHandlerTestEnv(t)

	// 配置 mock 返回错误
	mockLogic.ShouldError = true
	mockLogic.ErrorMsg = "业务错误：记录不存在"

	reqBody := types.<Request>Type{
		Id: 999,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/<module>/<action>", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler := <Handler>Handler(svcCtx)
	handler(resp, req)

	// Logic 层错误应正确传递
	assert.NotEqual(t, http.StatusOK, resp.Code)
}

func Test<Handler>_MissingAuth(t *testing.T) {
	svcCtx, _ := setupHandlerTestEnv(t)

	reqBody := types.<Request>Type{}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/<module>/<action>", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// 不设置 userId context
	resp := httptest.NewRecorder()

	handler := <Handler>Handler(svcCtx)
	handler(resp, req)

	// 需要认证的接口应拒绝
	assert.NotEqual(t, http.StatusOK, resp.Code)
}
```

### 1.3 Handler 测试最佳实践

1. **不依赖真实数据库**：Handler 测试应该完全隔离
2. **Mock Logic 接口**：通过接口 mock，不关心 Logic 内部实现
3. **覆盖边界情况**：
   - 无效 JSON
   - 缺少必填字段
   - 认证/授权失败
   - Logic 层错误

---

## 2. Logic 层测试模板

### 2.1 职责边界

Logic 层**负责**：
- 核心业务逻辑
- 数据验证
- 调用 Repository/Model 层
- 错误处理和包装

Logic 层**不包含**：
- HTTP 相关逻辑
- 请求/响应格式化

### 2.2 完整模板（使用 sqlmock）

```go
// backend/api/internal/logic/<module>/<logic>_test.go
package <module>

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ============================================================
// 测试工具函数
// ============================================================

// setupLogicTestWithMock 创建带 sqlmock 的测试环境
func setupLogicTestWithMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *svc.ServiceContext) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true, // 简化 mock
	})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	svcCtx := &svc.ServiceContext{
		DB: gormDB,
		Config: config.Config{
			Auth: config.AuthConfig{
				AccessSecret: "test-secret",
				AccessExpire: 3600,
			},
		},
	}

	// 注册清理
	t.Cleanup(func() {
		sqlDB.Close()
	})

	return gormDB, mock, svcCtx
}

// ============================================================
// 测试用例模板
// ============================================================

func Test<Logic>_<Action>_Success(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithMock(t)

	// 期望查询用户
	rows := sqlmock.NewRows([]string{"id", "username", "password", "status"}).
		AddRow(1, "testuser", "$2a$10$hashedpassword", "active")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	// 期望查询角色
	roleRows := sqlmock.NewRows([]string{"id", "code", "name"}).
		AddRow(1, "admin", "管理员")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT roles.* FROM `roles` INNER JOIN user_roles ur ON roles.id = ur.role_id WHERE ur.user_id = ?")).
		WithArgs(1).
		WillReturnRows(roleRows)

	// 创建 Logic 并执行
	ctx := context.Background()
	logic := New<Logic>Logic(ctx, svcCtx)

	req := &types.<Request>Type{
		Username: "testuser",
		// ... 其他字段
	}

	resp, err := logic.<Action>(req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.UserId)

	// 验证所有 SQL 期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test<Logic>_<Action>_RecordNotFound(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithMock(t)

	// 期望查询返回空结果
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("nonexistent", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "status"}))

	ctx := context.Background()
	logic := New<Logic>Logic(ctx, svcCtx)

	req := &types.<Request>Type{
		Username: "nonexistent",
	}

	resp, err := logic.<Action>(req)

	// 验证错误
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "不存在")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test<Logic>_<Action>_DBError(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithMock(t)

	// 期望数据库错误
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("testuser", 1).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	logic := New<Logic>Logic(ctx, svcCtx)

	req := &types.<Request>Type{
		Username: "testuser",
	}

	resp, err := logic.<Action>(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test<Logic>_<Action>_ValidationFailed(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithMock(t)

	ctx := context.Background()
	logic := New<Logic>Logic(ctx, svcCtx)

	// 传入无效请求（如空字符串）
	req := &types.<Request>Type{
		Username: "", // 无效
	}

	resp, err := logic.<Action>(req)

	// 验证错误（不需要 DB 调用）
	assert.Error(t, err)
	assert.Nil(t, resp)

	// 不应该有任何 SQL 调用
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test<Logic>_<Action>_BusinessRuleViolation(t *testing.T) {
	db, mock, svcCtx := setupLogicTestWithMock(t)

	// 测试业务规则违规（如状态不正确）
	rows := sqlmock.NewRows([]string{"id", "username", "password", "status"}).
		AddRow(1, "disableduser", "$2a$10$hashedpassword", "disabled") // 禁用状态
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("disableduser", 1).
		WillReturnRows(rows)

	ctx := context.Background()
	logic := New<Logic>Logic(ctx, svcCtx)

	req := &types.<Request>Type{
		Username: "disableduser",
	}

	resp, err := logic.<Action>(req)

	// 业务规则应该阻止操作
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已禁用")

	assert.NoError(t, mock.ExpectationsWereMet())
}
```

### 2.3 使用真实数据库的 Logic 测试（可选）

对于复杂的 Logic 逻辑，也可以使用真实数据库：

```go
// backend/api/internal/logic/<module>/<logic>_integration_test.go
package <module>

import (
	"context"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/testutil"
	"dmh/api/internal/types"

	"github.com/stretchr/testify/assert"
)

func Test<Logic>_<Action>_WithRealDB(t *testing.T) {
	// 使用现有的 MySQL 测试助手
	db, _ := testutil.SetupMySQLTestDB(t)

	svcCtx := &svc.ServiceContext{
		DB: db,
		Config: config.Config{
			Auth: config.AuthConfig{
				AccessSecret: "test-secret",
				AccessExpire: 3600,
			},
		},
	}

	// 准备测试数据
	// ... 使用 db 创建测试数据

	ctx := context.Background()
	logic := New<Logic>Logic(ctx, svcCtx)

	req := &types.<Request>Type{
		// ...
	}

	resp, err := logic.<Action>(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
```

---

## 3. Repository 层测试模板

### 3.1 职责边界

**注意**：DMH 项目当前没有独立的 Repository 层，Logic 直接使用 GORM。以下模板适用于：

1. **Model 层测试**（`backend/model/model_test.go` 风格）
2. **未来引入 Repository 层时的测试**

Repository 层**负责**：
- 数据访问（CRUD 操作）
- SQL 查询正确性
- 事务管理
- 数据库约束验证

### 3.2 完整模板（使用真实 MySQL8）

```go
// backend/model/<model>_test.go
// 或 backend/api/internal/repository/<repo>_test.go

package <repository>

import (
	"context"
	"testing"
	"time"

	"dmh/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// ============================================================
// 测试套件模式（推荐）
// ============================================================

type <Model>RepositoryTestSuite struct {
	suite.Suite
	db      *gorm.DB
	cleanup func()
	ctx     context.Context
}

func Test<Model>RepositorySuite(t *testing.T) {
	suite.Run(t, new(<Model>RepositoryTestSuite))
}

func (s *<Model>RepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	// 使用共享容器（详见 task-03 方案）
}

func (s *<Model>RepositoryTestSuite) SetupTest() {
	// 每个测试获取独立数据库
	s.db, s.cleanup = testutil.SetupMySQLTestDB(s.T())
}

func (s *<Model>RepositoryTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

// ============================================================
// CRUD 测试
// ============================================================

func (s *<Model>RepositoryTestSuite) TestCreate() {
	// 准备依赖数据
	brand := &model.Brand{Name: "Test Brand", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	// 创建测试对象
	entity := &model.<Model>{
		BrandId: brand.Id,
		Name:    "Test Name",
		Status:  "active",
		// ... 其他字段
	}

	// 执行创建
	err := s.db.Create(entity).Error

	// 验证
	s.NoError(err)
	s.NotZero(entity.Id)
	s.NotZero(entity.CreatedAt)
}

func (s *<Model>RepositoryTestSuite) TestCreate_UniqueConstraint() {
	// 创建第一个
	entity1 := &model.<Model>{
		Username: "uniqueuser",
		Password: "hash",
		Phone:    "13800138000",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(entity1).Error)

	// 尝试创建重复
	entity2 := &model.<Model>{
		Username: "uniqueuser", // 重复
		Password: "hash",
		Phone:    "13800138001",
		Status:   "active",
	}

	err := s.db.Create(entity2).Error
	s.Error(err) // 应该失败
	s.Contains(err.Error(), "Duplicate")
}

func (s *<Model>RepositoryTestSuite) TestRead() {
	// 创建测试数据
	brand := &model.Brand{Name: "Test", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	entity := &model.<Model>{
		BrandId: brand.Id,
		Name:    "Test Name",
		Status:  "active",
	}
	s.Require().NoError(s.db.Create(entity).Error)

	// 读取
	var found model.<Model>
	err := s.db.First(&found, entity.Id).Error

	s.NoError(err)
	s.Equal(entity.Name, found.Name)
}

func (s *<Model>RepositoryTestSuite) TestRead_NotFound() {
	var found model.<Model>
	err := s.db.First(&found, 99999).Error

	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *<Model>RepositoryTestSuite) TestUpdate() {
	// 创建
	entity := &model.<Model>{
		Username: "testuser",
		Password: "oldhash",
		Phone:    "13800138000",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(entity).Error)

	// 更新
	newStatus := "inactive"
	err := s.db.Model(entity).Update("status", newStatus).Error
	s.NoError(err)

	// 验证
	var found model.<Model>
	s.db.First(&found, entity.Id)
	s.Equal(newStatus, found.Status)
}

func (s *<Model>RepositoryTestSuite) TestDelete() {
	// 创建
	entity := &model.<Model>{
		Username: "testuser",
		Password: "hash",
		Phone:    "13800138000",
		Status:   "active",
	}
	s.Require().NoError(s.db.Create(entity).Error)

	// 删除
	err := s.db.Delete(entity).Error
	s.NoError(err)

	// 验证
	var found model.<Model>
	err = s.db.First(&found, entity.Id).Error
	s.Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func (s *<Model>RepositoryTestSuite) TestSoftDelete() {
	// 创建（假设 model 有 gorm.Model 或 DeletedAt 字段）
	campaign := &model.Campaign{
		BrandId:   1,
		Name:      "Test",
		Status:    "active",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}
	s.Require().NoError(s.db.Create(campaign).Error)

	// 软删除
	deletedAt := time.Now()
	err := s.db.Model(&model.Campaign{}).
		Where("id = ?", campaign.Id).
		Update("deleted_at", deletedAt).Error
	s.NoError(err)

	// 普通查询不应找到
	var found model.Campaign
	err = s.db.Where("id = ? AND deleted_at IS NULL", campaign.Id).First(&found).Error
	s.Error(err)

	// Unscoped 应该找到
	err = s.db.Unscoped().First(&found, campaign.Id).Error
	s.NoError(err)
	s.NotNil(found.DeletedAt)
}

// ============================================================
// 查询测试
// ============================================================

func (s *<Model>RepositoryTestSuite) TestQuery_WithConditions() {
	// 创建多条测试数据
	for i := 0; i < 5; i++ {
		entity := &model.<Model>{
			Username: fmt.Sprintf("user%d", i),
			Password: "hash",
			Phone:    fmt.Sprintf("1380013800%d", i),
			Status:   "active",
		}
		s.Require().NoError(s.db.Create(entity).Error)
	}

	// 查询
	var results []model.<Model>
	err := s.db.Where("status = ?", "active").Find(&results).Error

	s.NoError(err)
	s.Len(results, 5)
}

func (s *<Model>RepositoryTestSuite) TestQuery_WithJoin() {
	// 创建关联数据
	user := &model.User{Username: "test", Password: "hash", Phone: "13800138000", Status: "active"}
	s.Require().NoError(s.db.Create(user).Error)

	brand := &model.Brand{Name: "Test", Status: "active"}
	s.Require().NoError(s.db.Create(brand).Error)

	userBrand := &model.UserBrand{UserId: user.Id, BrandId: brand.Id}
	s.Require().NoError(s.db.Create(userBrand).Error)

	// 关联查询
	var result struct {
		UserId   int64
		BrandId  int64
		Username string
	}
	err := s.db.Table("user_brands").
		Select("user_brands.user_id, user_brands.brand_id, users.username").
		Joins("LEFT JOIN users ON users.id = user_brands.user_id").
		Where("user_brands.user_id = ?", user.Id).
		Scan(&result).Error

	s.NoError(err)
	s.Equal(user.Username, result.Username)
}

// ============================================================
// 事务测试
// ============================================================

func (s *<Model>RepositoryTestSuite) TestTransaction_Commit() {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 在事务中创建
		entity := &model.<Model>{
			Username: "txuser",
			Password: "hash",
			Phone:    "13800138000",
			Status:   "active",
		}
		if err := tx.Create(entity).Error; err != nil {
			return err
		}

		// 更新
		if err := tx.Model(entity).Update("status", "inactive").Error; err != nil {
			return err
		}

		return nil
	})

	s.NoError(err)

	// 验证数据已提交
	var found model.<Model>
	err = s.db.Where("username = ?", "txuser").First(&found).Error
	s.NoError(err)
	s.Equal("inactive", found.Status)
}

func (s *<Model>RepositoryTestSuite) TestTransaction_Rollback() {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		entity := &model.<Model>{
			Username: "rollbackuser",
			Password: "hash",
			Phone:    "13800138000",
			Status:   "active",
		}
		if err := tx.Create(entity).Error; err != nil {
			return err
		}

		// 故意返回错误触发回滚
		return fmt.Errorf("simulated error")
	})

	s.Error(err)

	// 验证数据已回滚
	var found model.<Model>
	err = s.db.Where("username = ?", "rollbackuser").First(&found).Error
	s.Error(err) // 应该找不到
}
```

### 3.3 使用 testcontainers 的 Repository 测试（推荐）

```go
// backend/api/internal/repository/<repo>_container_test.go
package repository_test

import (
	"context"
	"testing"

	"dmh/api/internal/testutil/container"
	"dmh/model"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type <Model>RepositoryContainerTestSuite struct {
	suite.Suite
	db      *gorm.DB
	cleanup func()
	ctx     context.Context
}

func Test<Model>RepositoryContainerSuite(t *testing.T) {
	suite.Run(t, new(<Model>RepositoryContainerTestSuite))
}

func (s *<Model>RepositoryContainerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	// 启动共享 MySQL 容器
	container.SetupSharedMySQLContainer(s.ctx, s.T())
}

func (s *<Model>RepositoryContainerTestSuite) SetupTest() {
	// 每个测试独立数据库
	s.db, s.cleanup = container.SetupIsolatedTestDB(s.ctx, s.T())
}

func (s *<Model>RepositoryContainerTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func (s *<Model>RepositoryContainerTestSuite) TestCRUD() {
	// 与上面的测试套件模式相同
}
```

---

## 4. 测试运行命令

### 4.1 单元测试（Handler + Logic with mock）

```bash
# 运行所有单元测试（快速，不需要真实数据库）
cd backend
go test -v ./api/internal/handler/... ./api/internal/logic/... -short

# 运行特定模块
go test -v ./api/internal/handler/auth/...
go test -v ./api/internal/logic/auth/...
```

### 4.2 集成测试（Logic with real DB + Repository）

```bash
# 必须使用 -p 1 避免数据库冲突
cd backend
go test -p 1 -v ./model/...
go test -p 1 -v ./api/internal/logic/... -run Integration

# 使用 testcontainers
go test -p 1 -v ./api/internal/repository/...
```

### 4.3 完整测试套件

```bash
# 单元测试 + 集成测试
cd backend
go test -p 1 -v -count=1 ./...
```

---

## 5. Mock 库推荐

| 场景 | 推荐库 | 用途 |
|------|--------|------|
| Handler 测试 | 手动 Mock 结构体 | 简单直接 |
| Logic 测试 | `github.com/DATA-DOG/go-sqlmock` | Mock GORM 底层 SQL |
| Repository 测试 | testcontainers-go + MySQL8 | 真实数据库隔离 |
| 接口 Mock | `github.com/stretchr/testify/mock` | 复杂接口 mock |

---

## 6. 常见陷阱

### 6.1 Handler 测试依赖真实数据库

```go
// ❌ 错误：Handler 测试不应该依赖真实数据库
func TestHandler(t *testing.T) {
    db := setupTestDB(t) // 不应该这样做
    // ...
}

// ✅ 正确：Handler 测试应该 mock Logic
func TestHandler(t *testing.T) {
    mockLogic := NewMockLogic()
    // ...
}
```

### 6.2 Logic 测试并行运行导致数据库冲突

```bash
# ❌ 错误：并行运行会导致冲突
go test ./api/internal/logic/...

# ✅ 正确：使用 -p 1 串行执行
go test -p 1 ./api/internal/logic/...
```

### 6.3 sqlmock 期望不匹配

```go
// ❌ 错误：SQL 不完全匹配
mock.ExpectQuery("SELECT * FROM users").WillReturnRows(rows)

// ✅ 正确：使用 QuoteMeta 精确匹配
mock.ExpectQuery(regexp.QuoteMeta(
    "SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?",
)).WithArgs("testuser", 1).WillReturnRows(rows)
```

### 6.4 忘记验证 mock 期望

```go
// ❌ 错误：忘记验证
func TestLogic(t *testing.T) {
    db, mock, _ := setupMock(t)
    // ... 测试代码
    // 没有验证 mock 期望
}

// ✅ 正确：总是验证
func TestLogic(t *testing.T) {
    db, mock, _ := setupMock(t)
    // ... 测试代码
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

## 7. 变更记录

| 日期 | 版本 | 作者 | 变更内容 |
|------|------|------|----------|
| 2026-02-26 | v1.0 | Sisyphus | 初始版本，三层测试模板 |
