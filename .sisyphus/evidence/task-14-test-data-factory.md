# T14: 测试数据工厂/夹具规范化

## 1. 任务目标
建立统一的 fixture/factory 规则，减少重复造数代码。

## 2. 实施方案

### 2.1 架构设计

采用 **Factory 模式** 构建测试数据工厂系统，提供以下能力：

```
backend/api/internal/testutil/factory/
├── factory.go           # 基础接口和通用工具
├── user_factory.go     # 用户工厂
├── campaign_factory.go  # 品牌+活动工厂
├── order_factory.go    # 订单工厂
├── fixtures.go         # 聚合入口和快捷方法
└── factory_test.go     # 单元测试
```

### 2.2 核心接口

```go
// Factory 定义测试数据工厂的基础接口
type Factory[T any] interface {
    Build() *T                          // 创建实例（不持久化）
    BuildWith(overrides map[string]any) *T  // 创建实例，允许覆盖字段
    BuildList(count int) []*T             // 创建多个实例
    Create(db any) (*T, error)           // 创建并持久化
    CreateWith(db any, overrides map[string]any) (*T, error)
    CreateList(db any, count int) ([]*T, error)
}
```

### 2.3 工厂示例

#### UserFactory
```go
factory := NewUserFactory()

// 基础用法
user := factory.Build()  // 创建默认用户

// 覆盖字段
user := factory.BuildWith(map[string]any{
    "role":     "platform_admin",
    "username": "custom_admin",
})

// 便捷方法
admin := factory.BuildPlatformAdmin()
brandAdmin := factory.BuildBrandAdmin()
locked := factory.BuildLocked()

// 持久化
user, err := factory.Create(db)
admin, err := factory.CreatePlatformAdmin(db)
```

#### CampaignFactory
```go
factory := NewCampaignFactory()

// 必须提供 brand_id
campaign := factory.BuildWith(map[string]any{
    "brand_id": 123,
    "name":     "测试活动",
})

// 便捷方法
active := factory.BuildActive(brandId)
ended := factory.BuildEnded(brandId)
withDist := factory.BuildWithDistribution(brandId, 3)

// 持久化（验证必填字段）
campaign, err := factory.CreateActive(db, brandId)
// 缺少 brand_id 时返回 ErrBrandIdRequired
_, err := factory.CreateWith(db, nil) // err == ErrBrandIdRequired
```

#### OrderFactory
```go
factory := NewOrderFactory()

// 必须提供 campaign_id
order := factory.BuildWith(map[string]any{
    "campaign_id": 456,
    "phone":       "13800138000",
})

// 便捷方法
pending := factory.BuildPending(campaignId)
paid := factory.BuildPaid(campaignId)
verified := factory.BuildVerified(campaignId, verifierId)

// 持久化（验证必填字段）
order, err := factory.CreatePaid(db, campaignId)
// 缺少 campaign_id 时返回 ErrCampaignIdRequired
```

## 3. 使用指南

### 3.1 在测试套件中使用

```go
type RepositoryTestSuite struct {
    suite.Suite
    db       *gorm.DB
    fixtures *factory.Fixtures
}

func (s *RepositoryTestSuite) SetupTest() {
    s.db, _ = testutil.SetupMySQLTestDB(s.T())
    s.fixtures = factory.NewFixtures(s.db)
}

func (s *RepositoryTestSuite) TestUserRepository_FindByRole() {
    // 创建测试用户
    users, _ := s.fixtures.CreateTestUsers(10)
    
    // 创建管理员
    admin, _ := s.fixtures.User.CreatePlatformAdmin(s.db)
    
    // 测试逻辑...
}

func (s *RepositoryTestSuite) TestOrderRepository_Verify() {
    // 快速创建完整数据链
    user, brand, campaign, order, _ := s.fixtures.SetupFullOrderChain()
    
    // 测试逻辑...
}
```

### 3.2 快捷数据链创建

```go
fixtures := factory.NewFixtures(db)

// 完整订单链：用户 -> 品牌 -> 活动 -> 订单
user, brand, campaign, order, err := fixtures.SetupFullOrderChain()

// 已核销订单链
brand, campaign, order, err := fixtures.SetupVerifiedOrder(verifierId)

// 用户+品牌关联
user, brand, userBrand, err := fixtures.SetupUserWithBrand()
```

### 3.3 字段覆盖约定

支持两种命名风格（snake_case 和 camelCase）：
```go
// 以下两种写法等效
factory.BuildWith(map[string]any{"real_name": "张三"})
factory.BuildWith(map[string]any{"realName": "张三"})

// 时间字段支持两种类型
now := time.Now()
factory.BuildWith(map[string]any{"paid_at": &now})   // *time.Time
factory.BuildWith(map[string]any{"paid_at": now})    // time.Time
```

### 3.4 错误处理

必填字段缺失时返回明确错误：
```go
_, err := campaignFactory.CreateWith(db, nil)
// err == ErrBrandIdRequired

_, err := orderFactory.CreateWith(db, nil)
// err == ErrCampaignIdRequired
```

## 4. 测试验证

### 4.1 单元测试覆盖
- ✅ UserFactory: Build, BuildWith, BuildList, BuildPlatformAdmin, BuildBrandAdmin, BuildLocked
- ✅ BrandFactory: Build, BuildWith, Create
- ✅ CampaignFactory: Build, BuildWith, BuildActive, BuildEnded, BuildWithDistribution
- ✅ OrderFactory: Build, BuildWith, BuildPending, BuildPaid, BuildVerified
- ✅ 错误处理: ErrBrandIdRequired, ErrCampaignIdRequired
- ✅ 辅助函数: RandomSuffix, RandomPhone, RandomEmail, TimeHelpers

### 4.2 测试运行结果
```
=== RUN   TestUserFactory_Build
--- PASS: TestUserFactory_Build (0.00s)
=== RUN   TestUserFactory_BuildLocked
--- PASS: TestUserFactory_BuildLocked (0.00s)
... (共 30+ 测试通过)
PASS
ok      dmh/api/internal/testutil/factory    0.003s
```

## 5. 默认值策略

| 工厂 | 字段 | 默认值 |
|------|------|--------|
| User | Username | `testuser_<random>` |
| User | Password | bcrypt hash of "123456" |
| User | Role | `participant` |
| User | Status | `active` |
| Brand | Name | `测试品牌_<random>` |
| Brand | Status | `active` |
| Campaign | Status | `active` |
| Campaign | FormFields | 默认姓名+手机号表单 |
| Campaign | StartTime/EndTime | 当前时间 ~ 7天后 |
| Order | Status | `pending` |
| Order | PayStatus | `unpaid` |
| Order | Phone | `138xxxxxxxxx` |

## 6. 反模式（禁止）

| 模式 | 原因 |
|------|------|
| 在每个测试内手写完整初始化 | 重复代码多，维护困难 |
| 使用硬编码的手机号/用户名 | 唯一约束冲突 |
| 跳过 Build 直接 Create | 不便于仅构建不持久化的场景 |
| 忽略必填字段验证 | 导致数据库插入失败 |

## 7. 扩展指南

添加新实体工厂：

```go
// 1. 创建工厂文件
package factory

type ProductFactory struct {
    BaseFactory
}

func NewProductFactory() *ProductFactory {
    return &ProductFactory{}
}

func (f *ProductFactory) Build() *model.Product {
    return f.BuildWith(nil)
}

func (f *ProductFactory) BuildWith(overrides map[string]any) *model.Product {
    product := &model.Product{
        Name:   "测试产品_" + f.RandomSuffix(),
        Price:  99.99,
        Status: "active",
    }
    // 应用 overrides...
    return product
}

// 2. 在 fixtures.go 中添加
type Fixtures struct {
    // ...
    Product *ProductFactory
}

func NewFixtures(db *gorm.DB) *Fixtures {
    return &Fixtures{
        // ...
        Product: NewProductFactory(),
    }
}
```

## 8. 文件清单

| 文件 | 行数 | 说明 |
|------|------|------|
| factory.go | 117 | 基础接口和常量 |
| user_factory.go | 209 | 用户工厂 |
| campaign_factory.go | 303 | 品牌+活动工厂 |
| order_factory.go | 277 | 订单工厂 |
| fixtures.go | 202 | 聚合入口 |
| factory_test.go | 443 | 单元测试 |

**总计**: ~1550 行代码，覆盖 4 类核心实体

## 9. 变更记录

| 日期 | 版本 | 作者 | 变更内容 |
|------|------|------|----------|
| 2026-02-26 | v1.0 | Sisyphus | 初始版本，建立统一 Factory 模式 |
