# T4: AGENTS 测试规范草案

**生成时间**: 2026-02-26 20:38 CST
**状态**: 草案 - 待 T15 最终落地

---

## 1. 测试命令速查

### 1.1 后端测试

| 场景 | 命令 | 耗时 | 说明 |
|------|------|------|------|
| **快速验证 (PR)** | `go test -p 1 ./... -short` | <30s | 跳过集成测试 |
| **全量单元测试** | `go test -p 1 ./...` | 2-3min | 所有单元测试 |
| **覆盖率报告** | `go test -p 1 ./... -coverprofile=coverage.out -covermode=atomic` | 3-5min | 阈值 ≥78% |
| **集成测试** | `DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v -count=1` | 5-10min | 需运行 API |
| **单个测试** | `go test -v -run TestName ./path/to/pkg` | 变化 | 调试用 |

### 1.2 前端测试

| 场景 | 命令 | 耗时 | 说明 |
|------|------|------|------|
| **Admin 单测** | `cd frontend-admin && npm run test` | <1min | Vitest |
| **Admin 覆盖率** | `cd frontend-admin && npm run test:cov` | 1-2min | 阈值 ≥80% |
| **H5 单测** | `cd frontend-h5 && npm run test` | <1min | Vitest |
| **H5 覆盖率** | `cd frontend-h5 && npm run test:cov` | 1-2min | 阈值 ≥70% |

---

## 2. 测试分层职责

### 2.1 后端分层

| 层级 | 职责 | Mock 策略 | 测试重点 |
|------|------|-----------|----------|
| **Handler** | HTTP 解析与响应 | Mock Logic | 请求参数解析、响应格式、HTTP 状态码 |
| **Logic** | 业务逻辑 | Mock Repository | 业务规则、边界条件、错误处理 |
| **Repository** | 数据访问 | 真实 MySQL8 | SQL 正确性、事务、数据一致性 |

### 2.2 前端分层

| 类型 | 职责 | 位置 |
|------|------|------|
| **单元测试** | 组件逻辑、工具函数 | `tests/unit/` |
| **E2E 测试** | 跨页面主流程 | `e2e/` |

---

## 3. 何时跑什么测试

### 3.1 本地开发

```bash
# 开发中 - 快速验证
go test -p 1 ./... -short

# 提交前 - 确保通过
go test -p 1 ./... -short
cd frontend-admin && npm run test:cov
cd frontend-h5 && npm run test:cov
```

### 3.2 CI 对应关系

| 触发点 | 执行内容 | 对应命令 |
|--------|----------|----------|
| **PR 提交** | PR Gate | 后端 `-short` + 前端覆盖率 + Lint |
| **合并 main** | Merge Gate | PR Gate + 关键集成测试 |
| **每日 02:00 UTC** | Nightly | 全量测试 + E2E + 性能测试 |

---

## 4. 禁用反模式

### 4.1 后端反模式

| 禁止 | 原因 | 正确做法 |
|------|------|----------|
| Handler 中写业务逻辑 | 破坏分层 | 放入 Logic 层 |
| 集成测试用 mock 数据 | 失去集成价值 | 使用真实数据库 |
| 测试间共享可变状态 | 测试干扰 | 每测试独立初始化/清理 |
| Check-then-create 模式 | 竞态条件 | Delete-then-create |
| `t.Skip()` 掩盖失败 | 隐藏真实缺陷 | 修复或标记已知问题 |
| 无限重试 flaky 测试 | 掩盖不稳定性 | 根因分析，修复或隔离 |

### 4.2 前端反模式

| 禁止 | 原因 |
|------|------|
| `as any` 类型断言 | 使用正确类型 |
| 视觉像素级问题放单测 | 放 E2E |
| Mock 全局状态 | 独立组件测试 |

---

## 5. 数据库隔离 (关键)

### 5.1 必须使用 `-p 1`

```bash
# 正确
go test -p 1 ./...

# 错误 - 会导致主键冲突
go test ./...
```

**原因**：集成测试共享 `dmh_test` 数据库，并行执行会导致：
- 主键冲突 (`Duplicate entry 'X' for key 'users.PRIMARY'`)
- 测试数据设置的竞态条件

### 5.2 测试数据模式

```go
// 正确：先删除再创建（幂等）
func (s *Suite) SetupTest() {
    s.db.Exec("DELETE FROM users WHERE id IN (1, 2, 3)")
    for _, user := range testUsers {
        s.db.Create(&user)
    }
}

// 错误：先检查再创建（竞态）
if err := s.db.Where("id = ?", user.Id).First(&existing).Error; err != nil {
    s.db.Create(&user)  // 可能主键冲突
}
```

---

## 6. SKIP 原因标签

环境问题跳过时，使用标准标签：

| 标签 | 触发条件 |
|------|----------|
| `API_UNAVAILABLE` | API 服务不可达或超时 |
| `MYSQL_UNAVAILABLE` | MySQL 连接失败 |
| `REDIS_UNAVAILABLE` | Redis 连接失败（仅依赖测试） |
| `LOGIN_FAILED` | 登录接口返回非 200 |
| `DATA_PREP_FAILED` | 测试数据准备失败 |

```go
if !isAPIAvailable() {
    t.Skip("SKIP: API_UNAVAILABLE - API service not running")
}
```

---

## 7. 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DMH_INTEGRATION_BASE_URL` | `http://localhost:8889` | API 基础地址 |
| `DMH_TEST_ADMIN_USERNAME` | `admin` | 测试管理员账号 |
| `DMH_TEST_ADMIN_PASSWORD` | `123456` | 测试管理员密码 |
| `MYSQL_TEST_HOST` | `127.0.0.1` | MySQL 测试主机 |
| `MYSQL_TEST_PORT` | `3306` | MySQL 测试端口 |

---

## 8. 覆盖率阈值

| 模块 | 阈值 | 当前 |
|------|------|------|
| Backend | **78%** | ~78% |
| Admin | **80%** | 83.65% |
| H5 | **70%** | ~99% |

检查覆盖率：
```bash
# 后端
go test -p 1 ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -1

# 前端
cd frontend-admin && npm run test:cov
```

---

## 9. 故障排查入口

### 9.1 常见问题

| 问题 | 排查步骤 |
|------|----------|
| 测试超时 | 检查 `-short` 是否生效，排除慢测试 |
| 主键冲突 | 确认使用 `-p 1`，检查测试数据幂等性 |
| MySQL 连接失败 | 检查容器状态：`docker ps`，确认端口 3306 |
| 覆盖率不足 | 运行 `go tool cover -func=coverage.out` 查看详情 |

### 9.2 日志位置

| 类型 | 位置 |
|------|------|
| CI 日志 | GitHub Actions → Workflow → Job |
| 本地日志 | 终端输出 |
| 覆盖率报告 | `coverage.out` / `coverage/` 目录 |

---

## 10. 提交前检查清单

- [ ] `go test -p 1 ./... -short` 通过
- [ ] 覆盖率达标（后端 ≥78%，Admin ≥80%，H5 ≥70%）
- [ ] 无 `t.Skip()` 掩盖失败
- [ ] Handler 无业务逻辑
- [ ] 集成测试使用真实 DB（非 mock）
- [ ] 测试数据幂等（delete-then-create）
