# Task 15: AGENTS.md 测试规范最终落地

## 执行时间
- 开始: 2026-02-26
- 结束: 2026-02-26

## 任务目标

将测试执行规则写入 `AGENTS.md` 与 `backend/AGENTS.md`，包含：
- 测试命令（何时用 -short、何时全量）
- 分层测试职责
- CI 对应关系
- 禁用反模式
- 常见失败排查

## 依赖
- T4: AGENTS 测试规范草案 ✅
- T5: 证据归档命名规范 ✅

---

## Part A: AGENTS.md 测试规范增补内容

以下内容建议替换现有 **第 9 节"常用测试命令"**，并新增相关章节：

### 9. 测试执行规范

#### 9.1 测试命令速查表

**后端测试**：

| 场景 | 命令 | 耗时 | 说明 |
|------|------|------|------|
| 快速验证 | `go test -p 1 ./... -short` | <30s | 跳过集成测试，PR 预检 |
| 提交前验证 | `go test -p 1 ./...` | 2-3min | 全量单元测试 |
| 集成测试 | `DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v -count=1` | 5-10min | 需运行 API |
| 覆盖率 | `go test -p 1 ./... -coverprofile=coverage.out -covermode=atomic` | 3-5min | 生成报告 |

**前端测试**：

| 场景 | 命令 | 耗时 | 说明 |
|------|------|------|------|
| Admin 单元测试 | `cd frontend-admin && npm run test` | <1min | Vitest |
| Admin 覆盖率 | `cd frontend-admin && npm run test:cov` | 1-2min | ≥80% 阈值 |
| Admin E2E | `cd frontend-admin && npm run test:e2e` | 5-10min | Playwright |
| H5 单元测试 | `cd frontend-h5 && npm run test` | <1min | Vitest |
| H5 覆盖率 | `cd frontend-h5 && npm run test:cov` | 1-2min | ≥70% 阈值 |
| H5 E2E | `cd frontend-h5 && npm run test:e2e` | 5-10min | Playwright |

#### 9.2 测试分层职责

**后端分层测试策略**：

| 层级 | 职责 | Mock 策略 | 测试重点 |
|------|------|-----------|----------|
| **Handler** | HTTP 解析与响应 | Mock Logic | 请求参数解析、响应格式、HTTP 状态码 |
| **Logic** | 业务逻辑 | Mock Repository | 业务规则、边界条件、错误处理 |
| **Repository** | 数据访问 | 真实 MySQL8 | SQL 正确性、事务、数据一致性 |

**前端测试分类**：

| 测试类型 | 关注点 | 依赖 | 速度 |
|---------|--------|------|------|
| 单元测试 | 代码正确性 | Mock 隔离 | 毫秒级 |
| E2E 测试 | 用户场景 | 真实/模拟环境 | 秒级 |

#### 9.3 CI 流水线对应关系

| Gate | 执行内容 | 时延预算 | 频率 |
|------|----------|----------|------|
| **PR Gate** | 后端单元 + 前端单元 + Lint | ≤15min | 每次 PR |
| **Merge Gate** | 集成/E2E Smoke | ≤25min | 合并后 |
| **Nightly** | 全量回归 + 性能 + 覆盖率 | ≤45min | 每日 |

**覆盖率阈值**：

| 模块 | 阈值 |
|------|------|
| backend | ≥78% |
| frontend-admin | ≥80% |
| frontend-h5 | ≥70% |

#### 9.4 数据库隔离要求

**必须使用 `-p 1` 标志**：

```bash
go test -p 1 ./...
```

**原因**：
- 集成测试共享 `dmh_test` 测试数据库
- 并行运行会导致主键冲突和竞态条件
- `-p 1` 强制串行执行，确保隔离

#### 9.5 测试数据初始化模式

使用"先删后建"模式确保幂等性：

```go
// ✅ 正确：先删后建（幂等）
suite.db.Exec("DELETE FROM users WHERE id IN (1, 2, 3)")
for _, user := range testUsers {
    suite.db.Create(&user)
}

// ❌ 错误：检查-然后-创建（竞态条件）
if err := suite.db.Where("id = ?", 1).First(&existing).Error; err != nil {
    suite.db.Create(&user)  // 可能因并发导致重复键错误
}
```

#### 9.6 SKIP 原因标签

当测试因环境问题跳过时，使用标准标签：

| 标签 | 触发条件 |
|------|----------|
| `API_UNAVAILABLE` | API 服务不可达或超时 |
| `MYSQL_UNAVAILABLE` | MySQL 连接失败 |
| `REDIS_UNAVAILABLE` | Redis 连接失败（仅依赖测试） |
| `LOGIN_FAILED` | 登录接口返回非 200 |
| `DATA_PREP_FAILED` | 测试数据准备失败 |

### 10. 测试相关反模式（补充）

在现有 **3.1 反模式** 表格中补充：

| 模式 | 原因 |
|------|------|
| Handler 中写业务逻辑 | Handler 只负责解析请求、调用 Logic、返回响应 |
| 集成测试用 mock 数据 | 失去集成测试意义，应使用真实数据库 |
| 测试间共享可变状态 | 导致测试相互干扰，每个测试应独立初始化/清理 |
| 检查-然后-创建模式 | 竞态条件，应使用先删后建（幂等） |
| `t.Skip()` 掩盖失败 | 隐藏真实缺陷，应修复问题或标记为 known issue |
| 无限重试掩盖 flaky | 掩盖不稳定因素，应排查根因、修复或隔离 |

### 11. 测试故障排查（新增章节）

#### 11.1 常见测试失败诊断

| 症状 | 可能原因 | 排查步骤 |
|------|----------|----------|
| `Duplicate entry for key 'PRIMARY'` | 并行测试冲突 | 确认使用 `go test -p 1` |
| `connection refused` | 服务未启动 | 检查 MySQL/API 是否运行 |
| `access denied` | 数据库凭证错误 | 检查 `dmh-api.yaml` 配置 |
| `table doesn't exist` | Schema 未初始化 | 执行 `init.sql` |
| `context deadline exceeded` | 超时 | 检查网络/数据库性能 |
| `login returned 400` | 测试账号问题 | 检查 admin 账号是否存在 |

#### 11.2 诊断命令

```bash
# 检查 MySQL 连接
docker exec -i mysql8 mysql -uroot -p'Admin168' -e "SELECT 1"

# 检查 API 健康
curl http://localhost:8889/health

# 查看测试详细输出
go test -v ./...

# 运行单个测试
go test -v -run TestFunctionName ./path/to/package

# 查看具体文件覆盖率
go tool cover -func=coverage.out | grep -v "100.0%"

# 生成 HTML 覆盖率报告
go tool cover -html=coverage.out -o coverage.html
```

#### 11.3 集成测试修复脚本

```bash
# 一键修复登录问题并重跑回归
backend/scripts/repair_login_and_run_order_regression.sh
```

### 12. 提交前检查清单

提交代码前确认：

- [ ] `go test -p 1 ./... -short` 通过（后端快速验证）
- [ ] 覆盖率达标（后端 ≥78%，Admin ≥80%，H5 ≥70%）
- [ ] 无 `t.Skip()` 掩盖失败
- [ ] Handler 无业务逻辑
- [ ] 集成测试使用真实数据库（非 mock）
- [ ] 测试数据幂等（先删后建）

---

## Part B: backend/AGENTS.md 测试规范增补内容

以下内容建议插入到现有 `backend/AGENTS.md` 的 **## Testing** 章节中，扩展为完整规范：

### Testing (扩展版)

#### Test Commands Quick Reference

| Scenario | Command | Duration | Notes |
|----------|---------|----------|-------|
| Fast verify (PR pre-check) | `go test -p 1 ./... -short` | <30s | Skips integration tests |
| Full unit tests | `go test -p 1 ./...` | 2-3min | All unit tests |
| Integration tests | `DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v -count=1` | 5-10min | Requires running API |
| Coverage report | `go test -p 1 ./... -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out` | 3-5min | Threshold: ≥78% |
| Single test | `go test -v -run TestFunctionName ./path/to/package` | Varies | Debug specific test |

#### Layered Test Strategy

| Layer | Responsibility | Mock Strategy | Test Focus |
|-------|----------------|---------------|------------|
| **Handler** | HTTP parse/response | Mock Logic | Request parsing, response format, HTTP codes |
| **Logic** | Business logic | Mock Repository | Business rules, edge cases, error handling |
| **Repository** | Data access | Real MySQL8 | SQL correctness, transactions, constraints |

#### Test File Locations

```
backend/
├── api/internal/
│   ├── handler/
│   │   └── *_handler_test.go     # Handler tests (mock logic)
│   └── logic/
│       └── *_logic_test.go       # Logic tests (mock repo or real DB)
├── model/
│   └── *_test.go                 # Model tests (real DB)
└── test/
    ├── integration/              # Integration tests (live API)
    └── performance/              # Benchmark tests
```

#### Database Isolation (CRITICAL)

**MUST use `-p 1` flag**:

```bash
go test -p 1 ./...
```

**Why**: Integration tests share `dmh_test` database. Parallel execution causes:
- Primary key conflicts (`Duplicate entry 'X' for key 'users.PRIMARY'`)
- Race conditions in test data setup

#### Test Data Setup Pattern

Use "delete-then-create" for idempotent initialization:

```go
// ✅ CORRECT: Delete first, then create (idempotent)
func (s *UserRepoSuite) SetupTest() {
    s.db.Exec("DELETE FROM users WHERE id IN (1, 2, 3)")
    for _, user := range testUsers {
        s.db.Create(&user)
    }
}

// ❌ WRONG: Check-then-create (race condition)
if err := suite.db.Where("id = ?", user.Id).First(&existing).Error; err != nil {
    suite.db.Create(&user)  // Can fail with duplicate key
}
```

#### SKIP Reason Tags

When tests must skip due to environment issues, use standard tags:

| Tag | Trigger |
|-----|---------|
| `API_UNAVAILABLE` | API service unreachable or timeout |
| `MYSQL_UNAVAILABLE` | MySQL connection failed |
| `REDIS_UNAVAILABLE` | Redis connection failed (dependency tests only) |
| `LOGIN_FAILED` | Login endpoint returned non-200 |
| `DATA_PREP_FAILED` | Test data preparation failed |

Example:
```go
if !isAPIAvailable() {
    t.Skip("SKIP: API_UNAVAILABLE - API service not running")
}
```

#### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DMH_INTEGRATION_BASE_URL` | `http://localhost:8889` | API base URL |
| `DMH_TEST_ADMIN_USERNAME` | `admin` | Test admin username |
| `DMH_TEST_ADMIN_PASSWORD` | `123456` | Test admin password |

#### Coverage Threshold

- **Current**: ~68%
- **Target**: ≥78%

Check coverage:
```bash
go test -p 1 ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | tail -1
```

#### Performance Test Split

| Gate | Command | Duration | Notes |
|------|---------|----------|-------|
| PR Smoke | `go test -short -bench=. -benchtime=100ms ./test/performance/...` | <1s | No backend needed |
| Nightly Full | `go test -v ./test/performance/...` | ~1min | Requires backend |

#### Test Anti-Patterns (Extended)

| Avoid | Reason | Correct Approach |
|-------|--------|------------------|
| Business logic in handlers | Breaks layering | Put in Logic layer |
| Integration tests with mock data | Loses integration value | Use real database |
| Shared mutable state between tests | Test interference | Independent init/cleanup per test |
| Check-then-create pattern | Race condition | Delete-then-create (idempotent) |
| `t.Skip()` to hide failures | Hides real defects | Fix or mark as known issue |
| Infinite retries for flaky tests | Masks instability | Root cause analysis, fix or isolate |

#### Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| `Duplicate entry for key 'PRIMARY'` | Parallel test conflict | Use `go test -p 1` |
| `connection refused` | Service not running | Start MySQL/API |
| `access denied` | DB credentials wrong | Check `dmh-api.yaml` |
| `table doesn't exist` | Schema not initialized | Run `init.sql` |
| `context deadline exceeded` | Timeout | Check network/DB performance |
| `login returned 400` | Test account issue | Run `repair_login_and_run_order_regression.sh` |

#### Diagnostic Commands

```bash
# Check MySQL connection
docker exec -i mysql8 mysql -uroot -p'Admin168' -e "SELECT 1"

# Check API health
curl http://localhost:8889/health

# View test verbose output
go test -v ./...

# Run single test
go test -v -run TestFunctionName ./path/to/package

# Check coverage by file
go tool cover -func=coverage.out | grep -v "100.0%"

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Fix login and rerun regression
./scripts/repair_login_and_run_order_regression.sh
```

#### Pre-Commit Checklist

Before committing:

- [ ] `go test -p 1 ./... -short` passes
- [ ] Coverage ≥78%
- [ ] No `t.Skip()` hiding failures
- [ ] No business logic in handlers
- [ ] Integration tests use real DB (not mocks)
- [ ] Test data is idempotent (delete-then-create)

---

## 合并建议

### AGENTS.md 修改建议

1. **替换第 9 节"常用测试命令"**：用上述 Part A 的完整测试规范替换
2. **扩展第 3.1 节"反模式"**：补充测试相关禁止项
3. **新增第 11 节"测试故障排查"**：添加诊断命令和常见问题
4. **新增第 12 节"提交前检查清单"**：标准化提交流程

### backend/AGENTS.md 修改建议

1. **扩展 `## Testing` 章节**：用上述 Part B 的完整规范替换
2. 保持现有结构，在 Testing 章节内增加子章节

### 内容冲突检查

| 检查项 | 结果 |
|--------|------|
| 与现有 Handler/Logic 分层规则冲突 | ❌ 无冲突，增强说明 |
| 与现有 `-p 1` 要求冲突 | ❌ 无冲突，强调说明 |
| 与现有 go-zero 约定冲突 | ❌ 无冲突 |
| 与现有 GORM 模式冲突 | ❌ 无冲突 |
| 删除原有关键约束 | ❌ 无删除，仅增强 |

---

## 结果

- [x] 创建证据文件: `.sisyphus/evidence/task-15-agents-test-spec-final.md`
- [x] 准备 `AGENTS.md` 测试规范增补内容（文本块）
- [x] 准备 `backend/AGENTS.md` 测试规范增补内容（文本块）
- [x] 包含：何时跑 -short、何时跑全量
- [x] 包含：CI 对应关系（PR Gate / Merge Gate / Nightly）
- [x] 包含：常见失败排查指南
- [x] 无冲突条款

---

*证据版本: 1.0*
*创建日期: 2026-02-26*
*适用范围: DMH 测试优化计划 T15*
