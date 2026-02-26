# Repository 层测试容器方案：MySQL8 定稿

## 1. 选型结论

**推荐方案：testcontainers-go + MySQL 8.0 镜像**

### 决策依据

| 维度 | MySQL8 容器 | SQLite | 结论 |
|------|-------------|--------|------|
| 语法兼容性 | 100% | ~60% | MySQL8 完胜 |
| 行为一致性 | 完全一致 | 有差异 | MySQL8 完胜 |
| 外键约束 | 完整支持 | 有限支持 | MySQL8 完胜 |
| 事务隔离 | 4 种级别 | 1 种 | MySQL8 完胜 |
| 启动速度 | ~2-3s | <100ms | SQLite 优 |
| 资源占用 | 中等 | 极低 | SQLite 优 |
| 维护成本 | 中等 | 低 | SQLite 优 |

**结论：DMH 项目大量使用 MySQL 特有语法，SQLite 无法作为可靠的测试替代品。生产环境一致性优先于性能。**

---

## 2. 项目 MySQL 特有功能分析

### 2.1 表级特性（SQLite 不支持）

```sql
-- 所有表均使用以下 MySQL 特有语法
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='表注释'
```

**影响**：SQLite 不支持这些子句，需要修改所有 migration 文件或创建单独的测试 schema。

### 2.2 列级特性

| 特性 | MySQL | SQLite | 影响 |
|------|-------|--------|------|
| `AUTO_INCREMENT` | ✅ | ❌ (用 AUTOINCREMENT) | 需替换 |
| `ON UPDATE CURRENT_TIMESTAMP` | ✅ | ❌ | 功能缺失 |
| `COMMENT '注释'` | ✅ | ❌ | 需移除 |

### 2.3 索引特性

```sql
-- MySQL 特有写法
UNIQUE KEY uk_name (columns)
INDEX idx_name (columns)

-- SQLite 需要单独 CREATE INDEX
```

### 2.4 DML 特性

```sql
-- MySQL 特有，SQLite 使用 INSERT OR REPLACE
INSERT INTO ... ON DUPLICATE KEY UPDATE ...
```

### 2.5 动态 SQL（高级特性）

```sql
-- MySQL 特有预处理语句
SET @exist_col := (SELECT ...);
SET @sql := IF(@exist_col = 0, 'ALTER TABLE ...', 'SELECT ...');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
```

**来源**：`20250219_add_distributor_tables.sql` 第 112-119 行

### 2.6 外键约束

```sql
-- 项目大量使用级联删除
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
```

**SQLite 默认不启用外键约束**，需要额外配置且行为可能有差异。

---

## 3. 实施方案

### 3.1 依赖引入

```go
// go.mod
require (
    github.com/testcontainers/testcontainers-go v0.31.0
    github.com/testcontainers/testcontainers-go/modules/mysql v0.31.0
)
```

### 3.2 测试容器管理器

```go
// backend/api/internal/testutil/container/mysql_container.go
package container

import (
    "context"
    "fmt"
    "sync"
    "testing"
    "time"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/mysql"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var (
    sharedContainer *mysql.MySQLContainer
    containerOnce   sync.Once
    containerMu     sync.Mutex
)

// MySQLContainerConfig 容器配置
type MySQLContainerConfig struct {
    Image       string
    Database    string
    Username    string
    Password    string
    Charset     string
    Collation   string
}

// DefaultConfig 默认配置
func DefaultConfig() MySQLContainerConfig {
    return MySQLContainerConfig{
        Image:     "mysql:8.0",
        Database:  "dmh_test",
        Username:  "test",
        Password:  "test123",
        Charset:   "utf8mb4",
        Collation: "utf8mb4_unicode_ci",
    }
}

// SetupSharedMySQLContainer 创建共享的 MySQL 容器（测试套件级别）
// 容器在所有测试结束后自动清理
func SetupSharedMySQLContainer(ctx context.Context, t testing.TB) *mysql.MySQLContainer {
    t.Helper()
    
    containerOnce.Do(func() {
        config := DefaultConfig()
        
        var err error
        sharedContainer, err = mysql.Run(ctx,
            config.Image,
            mysql.WithDatabase(config.Database),
            mysql.WithUsername(config.Username),
            mysql.WithPassword(config.Password),
            mysql.WithScripts(getInitScripts()...),
            testcontainers.WithLogger(testcontainers.TestLogger(t)),
        )
        if err != nil {
            t.Fatalf("Failed to start MySQL container: %v", err)
        }
        
        // 注册清理
        t.Cleanup(func() {
            if err := sharedContainer.Terminate(ctx); err != nil {
                t.Logf("Warning: Failed to terminate MySQL container: %v", err)
            }
            sharedContainer = nil
            containerOnce = sync.Once{}
        })
    })
    
    return sharedContainer
}

// SetupIsolatedTestDB 为单个测试创建隔离数据库
// 使用共享容器，但每个测试有独立的数据库
func SetupIsolatedTestDB(ctx context.Context, t testing.TB) (*gorm.DB, func()) {
    t.Helper()
    
    containerMu.Lock()
    defer containerMu.Unlock()
    
    // 确保共享容器已启动
    if sharedContainer == nil {
        SetupSharedMySQLContainer(ctx, t)
    }
    
    config := DefaultConfig()
    dbName := generateTestDBName(t)
    
    // 获取连接字符串
    connStr, err := sharedContainer.ConnectionString(ctx, "parseTime=true&loc=Local")
    if err != nil {
        t.Fatalf("Failed to get connection string: %v", err)
    }
    
    // 连接到默认数据库
    adminDB, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        t.Fatalf("Failed to connect to admin database: %v", err)
    }
    
    // 创建隔离数据库
    createSQL := fmt.Sprintf(
        "CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s",
        dbName, config.Charset, config.Collation,
    )
    if err := adminDB.Exec(createSQL).Error; err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }
    
    sqlDB, _ := adminDB.DB()
    sqlDB.Close()
    
    // 连接到测试数据库
    testConnStr, _ := sharedContainer.ConnectionString(ctx, 
        fmt.Sprintf("parseTime=true&loc=Local&dbname=%s", dbName))
    
    testDB, err := gorm.Open(mysql.Open(testConnStr), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // 运行 migration
    if err := runMigrations(testDB); err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }
    
    // 清理函数
    cleanup := func() {
        dropTestDB(t, dbName)
        sqlDB, _ := testDB.DB()
        if sqlDB != nil {
            sqlDB.Close()
        }
    }
    
    return testDB, cleanup
}

// generateTestDBName 生成唯一的数据库名
func generateTestDBName(t testing.TB) string {
    return fmt.Sprintf("t_%d_%s", time.Now().UnixNano()%1000000, 
        sanitizeName(t.Name()))
}

// dropTestDB 删除测试数据库
func dropTestDB(t testing.TB, dbName string) {
    config := DefaultConfig()
    connStr, _ := sharedContainer.ConnectionString(ctx, "parseTime=true&loc=Local")
    
    db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        t.Logf("Warning: Failed to connect for cleanup: %v", err)
        return
    }
    
    db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
    sqlDB, _ := db.DB()
    if sqlDB != nil {
        sqlDB.Close()
    }
}

// getInitScripts 获取初始化脚本路径
func getInitScripts() []string {
    return []string{
        "backend/migrations/init_schema.sql", // 预编译的 schema
    }
}
```

### 3.3 测试套件集成

```go
// backend/api/internal/logic/repository/suite_test.go
package repository_test

import (
    "context"
    "testing"
    
    "dmh/api/internal/testutil/container"
    "github.com/stretchr/testify/suite"
    "gorm.io/gorm"
)

type RepositoryTestSuite struct {
    suite.Suite
    db      *gorm.DB
    cleanup func()
    ctx     context.Context
}

func TestRepositorySuite(t *testing.T) {
    suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupSuite() {
    s.ctx = context.Background()
    // 共享容器在套件级别创建
}

func (s *RepositoryTestSuite) SetupTest() {
    // 每个测试获取独立数据库
    s.db, s.cleanup = container.SetupIsolatedTestDB(s.ctx, s.T())
}

func (s *RepositoryTestSuite) TearDownTest() {
    if s.cleanup != nil {
        s.cleanup()
    }
}

func (s *RepositoryTestSuite) TestUserRepository_Create() {
    repo := NewUserRepository(s.db)
    // 测试逻辑...
}
```

### 3.4 CI/CD 集成

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      
      - name: Run tests with testcontainers
        env:
          TESTCONTAINERS_RYUK_DISABLED: true  # CI 优化
        run: |
          cd backend
          go test -p 1 -v -count=1 ./...
```

---

## 4. 事务回滚策略

### 4.1 策略选择

| 策略 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| 数据库级隔离 | 完全隔离、无污染 | 资源占用略高 | 推荐方案 |
| 事务回滚 | 速度快 | 外键约束问题 | 简单表 |
| TRUNCATE | 简单 | 不重置 AUTO_INCREMENT | 特定场景 |

### 4.2 推荐方案：数据库级隔离

```go
// 每个测试创建独立数据库，测试完成后删除
// 优点：
// 1. 完全隔离，无交叉污染
// 2. AUTO_INCREMENT 重置
// 3. 外键约束正常工作
// 4. 测试可并行执行

func (s *RepositoryTestSuite) SetupTest() {
    s.db, s.cleanup = container.SetupIsolatedTestDB(s.ctx, s.T())
}
```

### 4.3 备选方案：事务回滚（仅限简单场景）

```go
// 仅适用于无复杂外键关系的表
func (s *SimpleTestSuite) SetupTest() {
    s.tx = s.db.Begin()
}

func (s *SimpleTestSuite) TearDownTest() {
    s.tx.Rollback()
}
```

---

## 5. 性能优化方案

### 5.1 共享容器策略

```
┌─────────────────────────────────────────────────────────┐
│                    测试套件生命周期                        │
├─────────────────────────────────────────────────────────┤
│  1. SetupSuite: 启动 1 个 MySQL 容器 (约 2-3s)            │
│  2. SetupTest: 创建隔离数据库 (< 100ms)                   │
│  3. Test 执行                                             │
│  4. TearDownTest: 删除测试数据库 (< 50ms)                 │
│  ... 重复 2-4 ...                                        │
│  N. TearDownSuite: 停止容器                               │
└─────────────────────────────────────────────────────────┘
```

**性能收益**：
- 单容器启动：~2.5s（仅一次）
- 每个测试数据库创建：~50-100ms
- 100 个测试总耗时：~2.5s + 100 × 0.1s = ~12.5s

### 5.2 预编译 Schema

```go
// 生成预编译 schema 脚本
//go:generate go run scripts/compile_schema.go

func getInitScripts() []string {
    return []string{
        "backend/migrations/init_schema.sql", // 预编译
    }
}
```

### 5.3 并行执行限制

```bash
# 必须使用 -p 1 避免数据库冲突
go test -p 1 ./...

# 或使用 t.Parallel() 配合隔离数据库
func (s *RepositoryTestSuite) TestParallel() {
    s.T().Parallel() // 安全：每个测试有独立数据库
}
```

### 5.4 CI 优化

```yaml
env:
  # 禁用 Ryuk（CI 中不需要）
  TESTCONTAINERS_RYUK_DISABLED: true
  # 使用 Docker 层缓存
  DOCKER_BUILDKIT: 1
```

---

## 6. 回退方案

### 6.1 场景一：CI 环境不支持 Docker

**方案**：使用预置 MySQL 服务

```yaml
# GitHub Actions
services:
  mysql:
    image: mysql:8.0
    env:
      MYSQL_ROOT_PASSWORD: test
      MYSQL_DATABASE: dmh_test
    ports:
      - 3306:3306
```

```go
// 检测环境自动切换
func SetupTestDB(t testing.TB) *gorm.DB {
    if os.Getenv("CI") != "" {
        return setupCIMySQL(t)
    }
    return setupTestcontainers(t)
}
```

### 6.2 场景二：testcontainers 不可用

**方案**：使用现有的 MySQL 容器（已在 Docker Compose 中）

```go
// 复用 deploy/docker-compose-simple.yml 中的 mysql8 容器
func SetupLocalTestDB(t testing.TB) *gorm.DB {
    config := testutil.GetMySQLTestConfig()
    // 使用现有的 mysql8 容器...
}
```

### 6.3 性能回退阈值

| 场景 | 阈值 | 回退方案 |
|------|------|----------|
| 容器启动 > 10s | 开发体验差 | 使用预置容器 |
| 单测试 > 5s | CI 超时 | 拆分测试 |
| 内存 > 2GB | 资源限制 | 串行执行 |

---

## 7. 迁移计划

### 7.1 阶段一：基础设施（1-2 天）

- [ ] 添加 testcontainers 依赖
- [ ] 创建 `container/mysql_container.go`
- [ ] 创建预编译 schema 脚本
- [ ] 更新 CI 配置

### 7.2 阶段二：Repository 测试迁移（2-3 天）

- [ ] 迁移 `UserRepository` 测试
- [ ] 迁移 `CampaignRepository` 测试
- [ ] 迁移 `OrderRepository` 测试
- [ ] 迁移 `DistributorRepository` 测试

### 7.3 阶段三：验收（1 天）

- [ ] 确保所有测试通过
- [ ] 验证 CI 执行时间可接受（< 10min）
- [ ] 文档更新

---

## 8. 参考资料

- [testcontainers-go 官方文档](https://golang.testcontainers.org/)
- [MySQL Module 文档](https://golang.testcontainers.org/modules/mysql/)
- [DMH 现有测试工具](../backend/api/internal/testutil/mysql_test_helper.go)
- [项目 Migration 文件](../backend/migrations/)

---

## 9. 变更记录

| 日期 | 版本 | 作者 | 变更内容 |
|------|------|------|----------|
| 2026-02-26 | v1.0 | Sisyphus | 初始版本，定稿 MySQL8 容器方案 |
