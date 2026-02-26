# T3: Repository 容器测试方案定稿

**生成时间**: 2026-02-26 20:38 CST

---

## 1. 选型结论

### 1.1 方案选择

| 方案 | 优点 | 缺点 | 决策 |
|------|------|------|------|
| **MySQL8 现有容器** | 零额外依赖、启动快、CI 已集成 | 需手动管理清理 | ✅ **采用** |
| testcontainers-go | 自动生命周期 | 冷启动慢(~10s)、CI 资源消耗大 | ❌ 不采用 |
| SQLite 内存库 | 最快 | 与生产环境不一致 | ❌ 禁止 |

**选型理由**：
1. CI 已配置 MySQL8 service container，无需额外基础设施
2. 现有 `mysql_test_helper.go` 已实现完善的隔离机制
3. 生产环境使用 MySQL8，测试环境必须保持一致

---

## 2. 现有方案分析

### 2.1 核心工具：`mysql_test_helper.go`

```go
// 位置：backend/api/internal/testutil/mysql_test_helper.go

// 主要功能：
// 1. SetupMySQLTestDB(t) - 创建隔离测试数据库
// 2. CleanupMySQLTestDB(t, dbName) - 清理测试数据库
// 3. SkipIfNoMySQL(t) - 环境不可用时跳过
// 4. migrateTestSchema(db) - 自动迁移表结构
```

### 2.2 隔离机制

| 特性 | 实现 |
|------|------|
| 数据库名 | `t_<testName>_<timestamp>_<random>` |
| 自动清理 | `t.Cleanup()` 注册清理函数 |
| Schema 迁移 | 自动创建所有必要表 |
| 连接池 | 10 连接，5 空闲，5 分钟生命周期 |

### 2.3 环境变量配置

```bash
MYSQL_TEST_HOST=127.0.0.1    # 默认
MYSQL_TEST_PORT=3306          # 默认
MYSQL_TEST_USER=root          # 默认
MYSQL_TEST_PASSWORD=Admin168  # 默认
MYSQL_TEST_DB=dmh             # 默认
```

---

## 3. 分层测试示范

### 3.1 Repository 层测试模板

基于 `order_repository_mysql8_test.go` 的示范：

```go
//go:build layered_demo

type OrderRepositoryTestSuite struct {
    suite.Suite
    db     *gorm.DB
    dbName string
}

func (s *OrderRepositoryTestSuite) SetupTest() {
    // 每个测试独立数据库
    s.db, s.dbName = testutil.SetupMySQLTestDB(s.T())
}

func (s *OrderRepositoryTestSuite) TestCreate() {
    // 1. 准备依赖数据
    brand := &model.Brand{Name: "Test Brand"}
    s.db.Create(brand)
    
    // 2. 执行测试
    order := &model.Order{BrandId: brand.Id, ...}
    err := s.db.Create(order).Error
    
    // 3. 验证结果
    s.NoError(err)
    s.NotZero(order.Id)
}
```

### 3.2 Handler/Logic 层 Mock 策略

| 层级 | Mock 对象 | 工具 |
|------|-----------|------|
| Handler | Logic | gomock / 手动 mock |
| Logic | Repository | gomock / 接口 mock |
| Repository | 无 (真实 DB) | - |

---

## 4. 事务回滚策略

### 4.1 当前方案：独立数据库

- 每个测试创建独立数据库 `t_<name>_<time>_<rand>`
- 测试结束后自动 `DROP DATABASE`
- **优点**：完全隔离，无竞态条件
- **缺点**：创建/删除数据库有开销

### 4.2 可选优化：事务回滚

```go
// 适用于简单测试的快速方案
func (s *Suite) SetupTest() {
    s.db = testutil.GetSharedDB()
    s.tx = s.db.Begin()
}

func (s *Suite) TearDownTest() {
    s.tx.Rollback()
}
```

**建议**：当前独立数据库方案已足够，无需优化。

---

## 5. CI 集成状态

### 5.1 现有配置 (pr-gate.yml)

```yaml
services:
  mysql8:
    image: mysql:8.0
    env:
      MYSQL_ROOT_PASSWORD: "Admin168"
      MYSQL_DATABASE: dmh
    ports:
      - 3306:3306
    options: >-
      --health-cmd="mysqladmin ping -h 127.0.0.1"
      --health-interval=10s
      --health-retries=20
```

### 5.2 测试命令

```bash
# 排除集成测试
go test -p 1 $(go list ./... | grep -v -E 'dmh/test/integration|dmh/test/performance') -coverprofile=coverage.out
```

---

## 6. 性能权衡

| 操作 | 耗时 | 说明 |
|------|------|------|
| 创建测试数据库 | ~100ms | 包含 Schema 迁移 |
| 单个 CRUD 测试 | ~10ms | GORM 操作 |
| 清理数据库 | ~50ms | DROP DATABASE |
| **单测试套件** | **~2s** | 10-20 个测试 |

---

## 7. 回退方案

| 场景 | 回退策略 |
|------|----------|
| MySQL 容器不可用 | `SkipIfNoMySQL(t)` 跳过并标记 |
| CI 资源不足 | 减少并行度 `-p 1` |
| Schema 迁移失败 | 记录错误，跳过测试 |

---

## 8. 验收标准

- [x] MySQL8 作为唯一数据引擎
- [x] 每测试独立数据库隔离
- [x] 自动清理机制
- [x] 环境不可用时优雅跳过
- [x] CI 集成完成

---

## 9. 后续任务依赖

此方案解锁：
- T7: 集成测试环境标准化
- T11: 后端模块分层示范
- T14: 测试数据工厂/夹具规范
