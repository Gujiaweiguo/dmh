# DMH 测试执行规范草案

> 本文档为 AGENTS.md 测试章节的补充规范草案，待审核后合并至 AGENTS.md。

---

## 1. 测试分层与执行命令

### 1.1 后端测试

#### 单元测试（快速验证）

```bash
cd backend
go test -p 1 ./... -short
```

- **适用场景**：日常开发、快速反馈循环、PR 预检
- **特点**：跳过集成测试，仅运行单元测试
- **耗时**：通常 < 30 秒

#### 全量测试（含集成测试）

```bash
cd backend
go test -p 1 ./...
```

- **适用场景**：提交前最终验证、CI 流水线
- **注意**：需要 MySQL8 容器运行

#### 集成测试（需运行 API）

```bash
cd backend
export DMH_INTEGRATION_BASE_URL=http://localhost:8889
go test ./test/integration/... -v -count=1
```

- **适用场景**：验证 API 端到端行为、回归测试
- **前置条件**：
  - MySQL8 容器运行
  - dmh-api 服务启动
- **环境变量**：
  - `DMH_INTEGRATION_BASE_URL`：API 基础地址（默认 `http://localhost:8889`）
  - `DMH_TEST_ADMIN_USERNAME`：测试管理员账号（默认 `admin`）
  - `DMH_TEST_ADMIN_PASSWORD`：测试管理员密码（默认 `123456`）

#### 覆盖率报告

```bash
cd backend
go test -p 1 ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out
```

**阈值要求**：≥ 78%

### 1.2 前端测试

#### 管理后台（frontend-admin）

```bash
cd frontend-admin
npm run test          # 单元测试
npm run test:cov      # 带覆盖率
npm run test:e2e      # E2E 测试
```

**覆盖率阈值**：≥ 80%

#### H5 前端（frontend-h5）

```bash
cd frontend-h5
npm run test          # 单元测试
npm run test:cov      # 带覆盖率
npm run test:e2e      # E2E 测试
```

**覆盖率阈值**：≥ 70%

---

## 2. 测试分层职责

### 2.1 后端分层测试策略

| 层级 | 职责 | Mock 策略 | 测试重点 |
|------|------|-----------|----------|
| **Handler** | HTTP 解析与响应 | Mock Logic | 请求参数解析、响应格式、HTTP 状态码 |
| **Logic** | 业务逻辑 | Mock Repository | 业务规则、边界条件、错误处理 |
| **Repository** | 数据访问 | 真实 MySQL8 | SQL 正确性、事务、数据一致性 |

### 2.2 Handler 测试模式

Handler 层应保持"薄"，仅测试：
- 请求参数正确解析
- 正确调用 Logic 方法
- 返回正确的 HTTP 状态码和响应格式

```go
// Handler 测试示例：Mock Logic，验证调用
func TestLoginHandler(t *testing.T) {
    mockLogic := new(MockAuthLogic)
    mockLogic.On("Login", mock.Anything, &types.LoginReq{
        Username: "test",
        Password: "123456",
    }).Return(&types.LoginResp{Token: "xxx"}, nil)
    
    // 验证 Handler 正确调用 Logic
}
```

### 2.3 Logic 测试模式

Logic 层是业务逻辑核心，测试应覆盖：
- 正常业务流程
- 边界条件
- 错误处理路径

```go
// Logic 测试示例：Mock Repository
func TestCreateUserLogic(t *testing.T) {
    mockRepo := new(MockUserRepo)
    mockRepo.On("Create", mock.Anything, &model.User{...}).Return(nil)
    
    logic := NewCreateUserLogic(ctx, svc, mockRepo)
    // 验证业务逻辑
}
```

### 2.4 Repository 测试模式

Repository 层使用真实 MySQL8 容器：
- 使用 `dmh_test` 测试数据库
- 遵循"先删后建"模式确保幂等性

```go
// Repository 测试示例：真实数据库
func (s *UserRepoSuite) SetupTest() {
    // 先清理，再创建（幂等）
    s.db.Exec("DELETE FROM users WHERE id IN (1, 2, 3)")
    s.db.Create(&testUsers)
}
```

---

## 3. 测试执行时机

### 3.1 本地开发

| 场景 | 命令 | 说明 |
|------|------|------|
| 快速验证 | `go test -p 1 ./... -short` | 跳过集成测试，秒级反馈 |
| 提交前 | `go test -p 1 ./...` | 全量测试 |
| 修改 API | 集成测试 | 验证端到端行为 |

### 3.2 CI 流水线（pr-gate.yml）

| Job | 命令 | 阈值 |
|-----|------|------|
| backend-unit | `go test -p 1 ./... -coverprofile=coverage.out` | ≥ 78% |
| frontend-unit (admin) | `npm run test:cov` | ≥ 80% |
| frontend-unit (h5) | `npm run test:cov` | ≥ 70% |
| lint | `gofmt -d .` | 必须通过 |

---

## 4. 禁止项（Anti-Patterns）

| 禁止行为 | 原因 | 正确做法 |
|----------|------|----------|
| Handler 中写业务逻辑 | 破坏分层架构 | 业务逻辑放 Logic 层 |
| 跳过测试（`t.Skip()`）掩盖问题 | 隐藏真实缺陷 | 修复问题或标记为 known issue |
| 无限重试掩盖 flaky 测试 | 掩盖不稳定因素 | 排查根因，修复或隔离 |
| 集成测试用 mock 数据 | 失去集成测试意义 | 使用真实数据库 |
| 测试间共享可变状态 | 导致测试相互干扰 | 每个测试独立初始化/清理 |
| 不处理错误返回值 | 测试不完整 | 验证错误路径 |
| 检查-然后-创建模式 | 竞态条件 | 先删后建（幂等） |

### 4.1 禁止：检查-然后-创建

```go
// ❌ 错误：竞态条件
if err := db.Where("id = ?", 1).First(&user).Error; err != nil {
    db.Create(&user)  // 可能因并发导致重复键错误
}

// ✅ 正确：先删后建
db.Exec("DELETE FROM users WHERE id = 1")
db.Create(&user)
```

### 4.2 禁止：跳过失败测试

```go
// ❌ 错误：掩盖问题
func TestSomething(t *testing.T) {
    t.Skip("flaky, will fix later")  // 永远不会被修复
    // ...
}

// ✅ 正确：记录并修复
func TestSomething(t *testing.T) {
    // 如果暂时无法修复，在 issue tracker 中记录
    // 并设置明确的修复期限
}
```

---

## 5. 数据库隔离要求

**必须使用 `-p 1` 标志**：

```bash
go test -p 1 ./...
```

**原因**：
- 集成测试共享 `dmh_test` 测试数据库
- 并行运行会导致主键冲突和竞态条件
- `-p 1` 强制串行执行，确保隔离

---

## 6. 故障排查

### 6.1 测试失败诊断

| 症状 | 可能原因 | 排查步骤 |
|------|----------|----------|
| `Duplicate entry for key 'PRIMARY'` | 并行测试冲突 | 确认使用 `-p 1` |
| `connection refused` | 服务未启动 | 检查 MySQL/API 是否运行 |
| `access denied` | 数据库凭证错误 | 检查 `dmh-api.yaml` 配置 |
| `table doesn't exist` | Schema 未初始化 | 执行 `init.sql` |
| `context deadline exceeded` | 超时 | 检查网络/数据库性能 |

### 6.2 常用诊断命令

```bash
# 检查 MySQL 连接
docker exec -i mysql8 mysql -uroot -p'Admin168' -e "SELECT 1"

# 检查 API 健康
curl http://localhost:8889/health

# 查看测试详细输出
go test -v ./...

# 运行单个测试
go test -v -run TestFunctionName ./path/to/package
```

### 6.3 覆盖率不达标

```bash
# 查看具体文件覆盖率
go tool cover -func=coverage.out | grep -v "100.0%"

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html
```

---

## 7. 测试文件组织

### 7.1 后端

```
backend/
├── api/internal/
│   ├── handler/
│   │   └── login_handler_test.go    # Handler 测试（mock logic）
│   └── logic/
│       └── login_logic_test.go      # Logic 测试（mock repo）
├── model/
│   └── user_test.go                 # Model 测试
└── test/
    ├── integration/                  # 集成测试（真实 DB）
    └── performance/                  # 性能测试
```

### 7.2 前端

```
frontend-admin/
├── tests/
│   └── unit/                        # 单元测试
└── e2e/                             # E2E 测试

frontend-h5/
├── tests/
│   └── unit/                        # 单元测试
└── e2e/                             # E2E 测试
```

---

## 8. 检查清单

提交代码前确认：

- [ ] `go test -p 1 ./... -short` 通过
- [ ] 覆盖率达标（后端 ≥78%，Admin ≥80%，H5 ≥70%）
- [ ] 无 `t.Skip()` 掩盖失败
- [ ] Handler 无业务逻辑
- [ ] 集成测试使用真实数据库
- [ ] 测试数据幂等（先删后建）

---

## 9. 合并建议

建议将本规范内容整合到 AGENTS.md 的以下位置：

1. **第 9 节"常用测试命令"**：扩展为完整测试规范
2. **第 3.1 节"反模式"**：补充测试相关禁止项
3. **新增"测试分层职责"章节**：定义 Handler/Logic/Repository 测试策略
4. **第 10 节"常见问题"**：补充测试故障排查

---

*草案版本：2024-02-26*
*状态：待审核*
