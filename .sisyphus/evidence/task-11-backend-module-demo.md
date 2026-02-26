# T11: 两个后端模块分层示范改造

## 任务完成状态

✅ 已完成两个高频模块（Order、Distributor）的分层测试示范改造

## 选择的模块

| 模块 | 原因 | 复杂度 | 业务重要性 |
|------|------|--------|------------|
| Order (订单) | 高频核心业务，包含核销、支付回调等复杂逻辑 | 高 | 高 |
| Distributor (分销商) | 高频核心业务，包含申请、审批、层级关系等复杂逻辑 | 高 | 高 |

## 创建的测试文件

### Order 模块

| 层级 | 文件 | 行数 | 测试用例数 |
|------|------|------|------------|
| Handler | `backend/api/internal/handler/order/verify_order_handler_mock_test.go` | ~300 | 7 |
| Logic | `backend/api/internal/logic/order/verify_order_logic_sqlmock_test.go` | ~425 | 8 |
| Repository | `backend/api/internal/logic/order/order_repository_mysql8_test.go` | ~603 | 16 |

### Distributor 模块

| 层级 | 文件 | 行数 | 测试用例数 |
|------|------|------|------------|
| Handler | `backend/api/internal/handler/distributor/distributor_apply_handler_mock_test.go` | ~396 | 9 |
| Logic | `backend/api/internal/logic/distributor/distributor_apply_logic_sqlmock_test.go` | ~427 | 10 |
| Repository | `backend/api/internal/logic/distributor/distributor_repository_mysql8_test.go` | ~745 | 18 |

## 分层测试模式总结

### Handler 层（Mock Logic）

```
┌─────────────────────────────────────────────────────────────┐
│                     Handler 测试模式                         │
├─────────────────────────────────────────────────────────────┤
│  1. 创建 Mock Logic 结构体，实现 Logic 接口                  │
│  2. Mock Logic 控制 ShouldError 和 ReturnData               │
│  3. Handler 测试不依赖真实数据库                             │
│  4. 测试覆盖：成功路径、无效 JSON、Logic 层错误、边界情况     │
└─────────────────────────────────────────────────────────────┘
```

**核心代码模式**：
```go
type MockVerifyOrderLogic struct {
    ShouldError bool
    ErrorMsg    string
    ReturnData  *types.VerifyOrderResp
    CalledWith  *types.VerifyOrderReq
}

func (m *MockVerifyOrderLogic) VerifyOrder(req *types.VerifyOrderReq) (*types.VerifyOrderResp, error) {
    m.CalledWith = req
    if m.ShouldError {
        return nil, &MockError{Message: m.ErrorMsg}
    }
    return m.ReturnData, nil
}
```

### Logic 层（SQLMock）

```
┌─────────────────────────────────────────────────────────────┐
│                     Logic 测试模式                           │
├─────────────────────────────────────────────────────────────┤
│  1. 使用 sqlmock.New() 创建 mock 数据库连接                 │
│  2. 配置 GORM 使用 mock 连接                                 │
│  3. 使用 mock.ExpectQuery/ExpectExec 设置 SQL 期望          │
│  4. 执行 Logic 方法                                          │
│  5. 验证返回结果和错误                                        │
│  6. 使用 mock.ExpectationsWereMet() 确保所有期望被满足       │
└─────────────────────────────────────────────────────────────┘
```

**核心代码模式**：
```go
func setupLogicTestWithSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *svc.ServiceContext) {
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err)
    
    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn: sqlDB,
        SkipInitializeWithVersion: true,
    }), &gorm.Config{
        SkipDefaultTransaction: true,
    })
    
    return gormDB, mock, &svc.ServiceContext{DB: gormDB}
}

// 使用示例
mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ?")).
    WithArgs(orderId, 1).
    WillReturnRows(rows)
```

### Repository 层（真实 MySQL8）

```
┌─────────────────────────────────────────────────────────────┐
│                   Repository 测试模式                        │
├─────────────────────────────────────────────────────────────┤
│  1. 使用 testify/suite 管理测试生命周期                      │
│  2. 每个测试使用独立数据库（testutil.SetupMySQLTestDB）      │
│  3. 测试覆盖：CRUD、唯一约束、关联查询、分页、事务           │
│  4. 测试隔离：SetupTest 创建，TearDownTest 自动清理         │
└─────────────────────────────────────────────────────────────┘
```

**核心代码模式**：
```go
type OrderRepositoryTestSuite struct {
    suite.Suite
    db      *gorm.DB
    dbName  string
}

func (s *OrderRepositoryTestSuite) SetupTest() {
    s.db, s.dbName = testutil.SetupMySQLTestDB(s.T())
}

func TestOrderRepositorySuite(t *testing.T) {
    suite.Run(t, new(OrderRepositoryTestSuite))
}
```

## 测试覆盖的场景

### Order 模块

| 场景类型 | Handler | Logic | Repository |
|----------|---------|-------|------------|
| 成功路径 | ✅ | ✅ | ✅ |
| 无效 JSON | ✅ | - | - |
| 订单不存在 | ✅ | ✅ | ✅ |
| 订单已核销 | ✅ | ✅ | - |
| 无效核销码 | ✅ | ✅ | - |
| 权限不足 | ✅ | ✅ | - |
| 数据库错误 | - | ✅ | ✅ |
| 事务回滚 | - | ✅ | ✅ |
| 唯一约束 | - | - | ✅ |
| 关联查询 | - | - | ✅ |

### Distributor 模块

| 场景类型 | Handler | Logic | Repository |
|----------|---------|-------|------------|
| 成功路径 | ✅ | ✅ | ✅ |
| 无效 JSON | ✅ | - | - |
| 重复申请 | ✅ | ✅ | - |
| 已是分销商 | ✅ | ✅ | ✅ |
| 未登录 | ✅ | ✅ | - |
| 无效品牌ID | ✅ | ✅ | - |
| 审批通过/拒绝 | - | ✅ | ✅ |
| 分销商层级 | - | - | ✅ |
| 链接唯一约束 | - | - | ✅ |
| 收益累计 | - | - | ✅ |

## 运行命令

```bash
# 运行 Order 模块分层测试
cd backend
go test -v ./api/internal/handler/order/verify_order_handler_mock_test.go
go test -v ./api/internal/logic/order/verify_order_logic_sqlmock_test.go
go test -p 1 -v ./api/internal/logic/order/order_repository_mysql8_test.go -run RepositorySuite

# 运行 Distributor 模块分层测试
go test -v ./api/internal/handler/distributor/distributor_apply_handler_mock_test.go
go test -v ./api/internal/logic/distributor/distributor_apply_logic_sqlmock_test.go
go test -p 1 -v ./api/internal/logic/distributor/distributor_repository_mysql8_test.go -run RepositorySuite
```

## 可复用模板

### Mock Logic 模板

```go
type Mock<Name>Logic struct {
    ShouldError bool
    ErrorMsg    string
    ReturnData  *types.<Response>
    CalledWith  *types.<Request>
}

func (m *Mock<Name>Logic) <Method>(req *types.<Request>) (*types.<Response>, error) {
    m.CalledWith = req
    if m.ShouldError {
        return nil, fmt.Errorf(m.ErrorMsg)
    }
    return m.ReturnData, nil
}
```

### SQLMock Setup 模板

```go
func setup<Module>LogicTestWithSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *svc.ServiceContext) {
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err)
    
    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn: sqlDB,
        SkipInitializeWithVersion: true,
    }), &gorm.Config{SkipDefaultTransaction: true})
    
    t.Cleanup(func() { sqlDB.Close() })
    return gormDB, mock, &svc.ServiceContext{DB: gormDB}
}
```

### Repository Suite 模板

```go
type <Module>RepositoryTestSuite struct {
    suite.Suite
    db     *gorm.DB
    dbName string
}

func (s *<Module>RepositoryTestSuite) SetupTest() {
    s.db, s.dbName = testutil.SetupMySQLTestDB(s.T())
}

func Test<Module>RepositorySuite(t *testing.T) {
    suite.Run(t, new(<Module>RepositoryTestSuite))
}
```

## 改造收益

| 维度 | 改造前 | 改造后 |
|------|--------|--------|
| 测试隔离 | Handler 依赖真实数据库 | Handler 完全隔离 |
| 测试速度 | 所有测试依赖数据库 | Handler/Logic 无数据库依赖 |
| 测试覆盖 | 主要覆盖成功路径 | 覆盖成功、错误、边界 |
| 可维护性 | 测试代码耦合度高 | 清晰分层，易于维护 |
| CI 友好 | 需要数据库环境 | Handler/Logic 无需数据库 |

## 验证结果

| 检查项 | 状态 |
|--------|------|
| Order Handler 测试文件创建 | ✅ |
| Order Logic sqlmock 测试创建 | ✅ |
| Order Repository MySQL8 测试创建 | ✅ |
| Distributor Handler 测试文件创建 | ✅ |
| Distributor Logic sqlmock 测试创建 | ✅ |
| Distributor Repository MySQL8 测试创建 | ✅ |
| Mock 隔离数据库 | ✅ |
| 可复用脚手架 | ✅ |

## 变更记录

| 日期 | 版本 | 作者 | 变更内容 |
|------|------|------|----------|
| 2026-02-26 | v1.0 | Sisyphus | 初始版本，完成 Order 模块 Handler 层测试示范 |
| 2026-02-26 | v2.0 | Sisyphus | 完成 Order 和 Distributor 模块三层完整测试示范 |
