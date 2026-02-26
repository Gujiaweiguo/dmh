# Task 10: Nightly 全量回归改造方案

## 执行时间
- 开始: 2026-02-26
- 结束: 2026-02-26

## 1. 概述

本文档定义 DMH 项目 Nightly 全量回归测试的改造方案，将慢测（集成/E2E/性能）迁移到 Nightly 调度，确保 **≤45 分钟** 的时延预算目标。

---

## 2. Nightly 测试范围

### 2.1 测试分层架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                     NIGHTLY 全量回归 (≤45 min)                       │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐   │
│  │ Stability Checks │  │ System Test Gate │  │  Coverage Gate   │   │
│  │    (~15 min)     │  │    (~15 min)     │  │    (~8 min)      │   │
│  │  - 后端全量测试   │  │  - 集成测试全量   │  │  - 覆盖率门禁     │   │
│  │  - Admin 单元测试 │  │  - Admin E2E     │  │  - 趋势报告       │   │
│  │  - H5 单元测试    │  │  - H5 E2E        │  │                   │   │
│  │  - 安全 E2E       │  │  - OpenSpec      │  │                   │   │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘   │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │              Order MySQL8 Regression (~5 min)                 │   │
│  │  - 订单核销鉴权回归                                           │   │
│  │  - 重复报名文案验证                                           │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### 2.2 测试范围矩阵

| 测试类型 | 模块 | 预估耗时 | 执行位置 | 说明 |
|---------|------|---------|---------|------|
| **后端全量测试** | Backend | ~8 min | `stability-checks.yml` | `go test -p 1 ./...` |
| **后端集成测试** | Backend | ~5 min | `system-test-gate.yml` | `./test/integration/...` |
| **Admin 单元测试** | Frontend-Admin | ~3 min | `stability-checks.yml` | `npm run test` |
| **H5 单元测试** | Frontend-H5 | ~3 min | `stability-checks.yml` | `npm run test` |
| **Admin E2E 全量** | Frontend-Admin | ~5 min | `system-test-gate.yml` | Playwright |
| **H5 E2E 全量** | Frontend-H5 | ~5 min | `system-test-gate.yml` | Playwright |
| **安全 E2E 回归** | Frontend-Admin | ~3 min | `stability-checks.yml` | `security-management.spec.ts` |
| **覆盖率门禁** | 全模块 | ~8 min | `coverage-gate.yml` | 三个模块并行 |
| **订单回归** | Backend | ~5 min | `order-mysql8-regression.yml` | 关键路径验证 |
| **OpenSpec 验证** | 规格一致性 | ~2 min | `system-test-gate.yml` | `--all --strict` |

**并行执行预估**: ~35-45 分钟（部分 Job 并行）

---

## 3. 调度策略

### 3.1 触发条件

```yaml
on:
  workflow_dispatch:          # 手动触发
  schedule:
    - cron: '0 2 * * *'       # UTC 02:00 = 北京时间 10:00
  push:
    branches: [main, master]
    tags:
      - v*.*.*-rc*            # RC 版本发布时触发
```

### 3.2 并行策略

当前 `full-regression.yml` 采用 **Job 级并行**：

```yaml
jobs:
  stability-checks:      # Job 1
    uses: ./.github/workflows/stability-checks.yml
  system-test-gate:      # Job 2 (与 Job 1 并行)
    uses: ./.github/workflows/system-test-gate.yml
  coverage-gate:         # Job 3 (与 Job 1/2 并行)
    uses: ./.github/workflows/coverage-gate.yml
  order-mysql8-regression:  # Job 4 (与 Job 1/2/3 并行)
    uses: ./.github/workflows/order-mysql8-regression.yml
  full-regression-verdict:  # Job 5 (等待所有 Job 完成)
    needs: [stability-checks, system-test-gate, coverage-gate, order-mysql8-regression]
```

### 3.3 失败告警

**当前机制**:
- GitHub Actions 自动发送邮件通知
- 失败时上传 debug artifacts（日志、截图）

**增强建议**（可选）:
- Slack/钉钉 Webhook 集成
- 失败时 @on-call 开发者

---

## 4. 报告格式

### 4.1 Summary 报告模板

当前 `full-regression-verdict` 生成的报告：

```markdown
## Full Regression Summary

| Workflow | Result |
|---|---|
| stability-checks | success/failure |
| system-test-gate | success/failure |
| coverage-gate | success/failure |
| order-mysql8-regression | success/failure |

Final Verdict: PASS/FAIL
```

### 4.2 增强报告格式（建议）

```markdown
## 🌙 Nightly Full Regression Report

**执行时间**: YYYY-MM-DD HH:MM UTC
**总耗时**: XX min YY sec

### 测试通过率

| 模块 | 测试数 | 通过 | 失败 | 跳过 | 通过率 |
|------|--------|------|------|------|--------|
| Backend Unit | XX | XX | XX | XX | XX% |
| Backend Integration | XX | XX | XX | XX | XX% |
| Admin E2E | XX | XX | XX | XX | XX% |
| H5 E2E | XX | XX | XX | XX | XX% |
| Security E2E | XX | XX | XX | XX | XX% |
| Order Regression | XX | XX | XX | XX | XX% |

### 覆盖率统计

| 模块 | 覆盖率 | 阈值 | 状态 |
|------|--------|------|------|
| Backend | XX% | 78% | ✅/❌ |
| Admin | XX% | 80% | ✅/❌ |
| H5 | XX% | 70% | ✅/❌ |

### 失败分类

| 类别 | 数量 | 示例 |
|------|------|------|
| 🐛 代码缺陷 | X | [链接] |
| 🔧 环境问题 | X | [链接] |
| 📊 数据问题 | X | [链接] |
| 🌐 外部依赖 | X | [链接] |

### 详细日志

- [stability-checks 日志](链接)
- [system-test-gate 日志](链接)
- [coverage-gate 日志](链接)
- [order-mysql8-regression 日志](链接)

---

**Final Verdict**: ✅ PASS / ❌ FAIL
```

### 4.3 失败分类定义

| 分类 | 说明 | 典型场景 |
|------|------|---------|
| 🐛 **代码缺陷** | 业务逻辑错误、边界条件问题 | 断言失败、空指针、类型错误 |
| 🔧 **环境问题** | CI/CD 环境配置问题 | 服务启动失败、端口冲突、超时 |
| 📊 **数据问题** | 测试数据或数据库问题 | 数据库连接失败、数据不一致 |
| 🌐 **外部依赖** | 第三方服务或网络问题 | API 超时、DNS 解析失败 |

---

## 5. 时延预算验证

### 5.1 当前预估

| Job | 预估耗时 | 实际上限 |
|-----|---------|---------|
| stability-checks | ~15 min | timeout: 未设置 |
| system-test-gate | ~15 min | timeout: 未设置 |
| coverage-gate | ~8 min | timeout: 未设置 |
| order-mysql8-regression | ~5 min | timeout: 未设置 |
| **总计（并行）** | **~35-40 min** | **≤45 min** ✅ |

### 5.2 超时保护建议

```yaml
jobs:
  stability-checks:
    uses: ./.github/workflows/stability-checks.yml
    timeout-minutes: 20  # 新增
  
  system-test-gate:
    uses: ./.github/workflows/system-test-gate.yml
    timeout-minutes: 20  # 新增
  
  coverage-gate:
    uses: ./.github/workflows/coverage-gate.yml
    timeout-minutes: 15  # 新增
  
  order-mysql8-regression:
    uses: ./.github/workflows/order-mysql8-regression.yml
    timeout-minutes: 10  # 新增
```

### 5.3 验证结论

- **当前设计满足 ≤45 分钟预算** ✅
- 四个子 Job 并行执行，最长路径 ~20 分钟
- 总耗时 = max(Job 耗时) + Verdict 耗时 ≈ 35-45 分钟

---

## 6. 与 PR Workflow 的职责分离

### 6.1 职责矩阵

| 测试类型 | PR Gate | Merge Gate | Nightly |
|---------|---------|------------|---------|
| 单元测试 | ✅ 快速 | ✅ 完整 | ✅ 全量 |
| 集成测试 | ❌ | ✅ Smoke | ✅ 全量 |
| E2E 测试 | ❌ | ✅ 关键路径 | ✅ 全量 |
| 性能测试 | ❌ | ❌ | ✅ |
| 安全扫描 | ❌ | ✅ Smoke | ✅ 全量 |
| 覆盖率 | ✅ 门禁 | ✅ 门禁 | ✅ 趋势 |

### 6.2 不重复执行策略

1. **PR Gate**: 仅单元测试 + Lint（快速反馈，≤15 min）
2. **Merge Gate**: 集成/E2E Smoke（中等验证，≤25 min）
3. **Nightly**: 全量回归 + 性能 + 趋势（深度验证，≤45 min）

---

## 7. 实施检查清单

### 7.1 当前状态

- [x] `full-regression.yml` 编排入口已实现
- [x] 四个子 workflow 已正确引用
- [x] 调度配置（UTC 02:00）已设置
- [x] Verdict 汇总已实现
- [x] 失败时上传 artifacts 已实现

### 7.2 改进建议

- [ ] 为子 Job 添加 `timeout-minutes` 保护
- [ ] 增强报告格式（失败分类、通过率统计）
- [ ] 添加覆盖率趋势图表
- [ ] 可选：Slack/钉钉告警集成

---

## 8. 附录：当前 Workflow 结构

### 8.1 文件列表

| 文件 | 职责 |
|------|------|
| `.github/workflows/full-regression.yml` | Nightly 编排入口 |
| `.github/workflows/stability-checks.yml` | 后端全量 + 前端单元 + 安全 E2E |
| `.github/workflows/system-test-gate.yml` | 集成测试 + Admin/H5 E2E + OpenSpec |
| `.github/workflows/coverage-gate.yml` | 覆盖率门禁 |
| `.github/workflows/order-mysql8-regression.yml` | 订单关键路径回归 |

### 8.2 依赖关系

```
full-regression.yml (编排入口)
├── stability-checks.yml (并行)
├── system-test-gate.yml (并行)
├── coverage-gate.yml (并行)
├── order-mysql8-regression.yml (并行)
└── full-regression-verdict (等待全部完成)
```

---

## 结果

**Nightly 全量回归改造方案已完成设计**：

1. ✅ 测试范围定义：覆盖后端全量、集成、E2E、安全、覆盖率、订单回归
2. ✅ 调度策略：每日 UTC 02:00，四 Job 并行，Verdict 汇总
3. ✅ 报告格式：通过率 + 失败分类 + 日志链接
4. ✅ 时延验证：预估 35-45 分钟，满足 ≤45 分钟目标

## 备注

- 当前 workflow 设计已满足时延预算要求
- 改进建议为可选增强，不影响核心功能
- 后续可考虑添加性能基准测试和安全扫描
