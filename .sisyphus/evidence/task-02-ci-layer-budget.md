# DMH CI 分层测试策略与时延预算

## 1. 概述

本文档定义 DMH 项目的 CI 分层测试策略，明确 PR/Merge/Nightly 三层测试的职责边界、测试范围和时延预算。

## 2. 三层 CI 职责矩阵

| 层级 | 触发条件 | 核心目标 | 时延预算 | 状态 |
|------|---------|---------|---------|------|
| **PR Gate** | PR 创建/更新 | 快速反馈，阻断明显缺陷 | **≤ 15 分钟** | ✅ 已实现 |
| **Merge Gate** | 合并到 main/master | 中等验证，确保主干稳定 | **≤ 25 分钟** | ✅ 已实现 |
| **Nightly** | 每日定时 (UTC 02:00) | 全量回归，深度验证 | **≤ 45 分钟** | ✅ 已实现 |

## 3. PR Gate (快速反馈层)

### 3.1 职责定义
- 提供开发者快速反馈（目标：10 分钟内可见结果）
- 阻断明显缺陷进入代码库
- 强制代码质量门禁（覆盖率、格式化）

### 3.2 测试范围

| 测试类型 | 模块 | 超时配置 | 覆盖率阈值 |
|---------|------|---------|-----------|
| 单元测试 | Backend | 15 min | ≥ 78% |
| 单元测试 | Frontend-Admin | 10 min | ≥ 80% |
| 单元测试 | Frontend-H5 | 10 min | ≥ 70% |
| Lint 检查 | Backend (gofmt) | 5 min | N/A |

### 3.3 Workflow 文件
- `.github/workflows/pr-gate.yml`

### 3.4 设计原则
- **禁止**包含集成测试、E2E 测试、性能测试
- **禁止**启动服务容器（MySQL/Redis）- 单元测试应 mock 外部依赖
- 使用 `concurrency.group` 确保同 PR 只运行最新版本

### 3.5 当前实现验证

```yaml
# pr-gate.yml 中的超时配置
jobs:
  backend-unit:
    timeout-minutes: 15  # ✅ 符合预算
  frontend-unit:
    timeout-minutes: 10  # ✅ 符合预算
  lint:
    timeout-minutes: 5   # ✅ 符合预算
```

**总预估时间**: ~15 分钟（并行执行）

---

## 4. Merge Gate (中等验证层)

### 4.1 职责定义
- 确保 main/master 分支始终处于可发布状态
- 运行需要真实依赖的集成测试 smoke
- 覆盖关键业务路径的 E2E 验证

### 4.2 测试范围

| 测试类型 | 模块 | 预估耗时 | 触发路径 |
|---------|------|---------|---------|
| 集成测试 Smoke | Backend (关键接口) | ~8 min | backend/** |
| E2E Smoke | Frontend-Admin | ~5 min | frontend-admin/** |
| E2E Smoke | Frontend-H5 | ~5 min | frontend-h5/** |
| OpenSpec 验证 | 规格一致性 | ~2 min | openspec/** |

### 4.3 Workflow 文件
- `.github/workflows/system-test-gate.yml` (PR 路径过滤 + Push 触发)
- `.github/workflows/stability-checks.yml` (PR 路径过滤 + Push 触发)

### 4.4 设计原则
- 使用 `paths` 过滤，仅在相关代码变更时触发
- 包含 MySQL/Redis 服务容器
- 启动真实 API 服务器进行集成测试
- 运行 Playwright E2E 测试

### 4.5 当前实现验证

```yaml
# system-test-gate.yml 触发条件
on:
  workflow_call:
  workflow_dispatch:
  pull_request:
    branches: [main]
    paths:
      - backend/**
      - frontend-admin/**
      - frontend-h5/**
      - openspec/**
  push:
    branches: [main]
    paths: [...] # 同上
```

**总预估时间**: ~20-25 分钟

---

## 5. Nightly (全量回归层)

### 5.1 职责定义
- 运行全量测试套件，发现日间开发累积的问题
- 执行长时间运行的测试（性能、压力、安全）
- 生成覆盖率趋势报告

### 5.2 测试范围

| 测试类型 | 模块 | 预估耗时 | 说明 |
|---------|------|---------|------|
| 后端全量测试 | Backend | ~10 min | `go test ./...` |
| 后端集成测试 | Backend | ~8 min | 完整集成套件 |
| 前端 E2E 全量 | Admin + H5 | ~10 min | 所有 E2E 用例 |
| 安全 E2E | Admin | ~3 min | 权限/鉴权场景 |
| 订单回归 | Backend | ~5 min | 订单核销关键路径 |
| 覆盖率门禁 | 全模块 | ~5 min | 独立报告 |
| OpenSpec 验证 | 规格一致性 | ~2 min | `--all --strict` |

### 5.3 Workflow 文件
- `.github/workflows/full-regression.yml` (编排入口)
  - 引用: `stability-checks.yml`
  - 引用: `system-test-gate.yml`
  - 引用: `coverage-gate.yml`
  - 引用: `order-mysql8-regression.yml`

### 5.4 调度配置

```yaml
on:
  workflow_dispatch:  # 手动触发
  schedule:
    - cron: '0 2 * * *'  # UTC 02:00 = 北京时间 10:00
  push:
    branches: [main, master]
    tags:
      - v*.*.*-rc*  # RC 版本发布时触发
```

### 5.5 设计原则
- **并行执行**多个子 workflow 以控制总时长
- 失败时上传 debug artifacts（日志、截图）
- 最终 verdict 汇总所有子任务结果

### 5.6 当前实现验证

```yaml
# full-regression.yml 结构
jobs:
  stability-checks:
    uses: ./.github/workflows/stability-checks.yml
  system-test-gate:
    uses: ./.github/workflows/system-test-gate.yml
  coverage-gate:
    uses: ./.github/workflows/coverage-gate.yml
  order-mysql8-regression:
    uses: ./.github/workflows/order-mysql8-regression.yml
  full-regression-verdict:
    needs: [stability-checks, system-test-gate, coverage-gate, order-mysql8-regression]
```

**总预估时间**: ~35-45 分钟（部分并行）

---

## 6. 时延预算总览

```
┌─────────────────────────────────────────────────────────────┐
│                    CI 时延预算金字塔                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│    ▲ Nightly (≤45m)                                        │
│    │ - 全量 E2E                                             │
│    │ - 性能测试                                             │
│    │ - 安全扫描                                             │
│    │ - 覆盖率趋势                                           │
│    │                                                        │
│   ▲▲ Merge Gate (≤25m)                                     │
│   ││ - 集成测试 Smoke                                       │
│   ││ - E2E 关键路径                                         │
│   ││ - OpenSpec 验证                                        │
│   ││                                                        │
│  ▲▲▲ PR Gate (≤15m)                                        │
│  │││ - 单元测试                                             │
│  │││ - 覆盖率门禁                                           │
│  │││ - Lint 检查                                            │
│  │││                                                        │
└──┴─┴───────────────────────────────────────────────────────┘

   频率: 高 ◄──────────────────────────────────────► 频率: 低
   耗时: 短 ◄──────────────────────────────────────► 耗时: 长
   范围: 窄 ◄──────────────────────────────────────► 范围: 宽
```

## 7. 覆盖率阈值

| 模块 | 阈值 | 位置 | 说明 |
|------|------|------|------|
| Backend | 78% | `pr-gate.yml`, `coverage-gate.yml` | 核心业务逻辑 |
| Frontend-Admin | 80% | `pr-gate.yml`, `coverage-gate.yml` | 管理后台 |
| Frontend-H5 | 70% | `pr-gate.yml`, `coverage-gate.yml` | H5 前端 |

**原则**: 不降低覆盖率阈值，只增不减。

## 8. 禁止事项

### 8.1 PR Gate 禁止
- ❌ 启动 MySQL/Redis 服务容器
- ❌ 运行集成测试 (`./test/integration/...`)
- ❌ 运行 E2E 测试
- ❌ 运行性能测试 (`./test/performance/...`)

### 8.2 覆盖率禁止
- ❌ 降低现有覆盖率阈值
- ❌ 跳过覆盖率检查（除非紧急 hotfix）

## 9. 现有 Workflow 映射

| 层级 | Workflow | 状态 | 备注 |
|------|---------|------|------|
| PR Gate | `pr-gate.yml` | ✅ 完整 | 后端+前端单元测试+Lint |
| Merge Gate | `system-test-gate.yml` | ✅ 完整 | 集成+E2E+OpenSpec |
| Merge Gate | `stability-checks.yml` | ✅ 完整 | 全量后端+安全E2E |
| Merge Gate | `order-mysql8-regression.yml` | ✅ 完整 | 订单关键路径 |
| Nightly | `full-regression.yml` | ✅ 完整 | 编排所有子workflow |
| 覆盖率 | `coverage-gate.yml` | ✅ 完整 | 独立覆盖率检查 |

## 10. 告警与通知

### 10.1 PR Gate 失败
- GitHub PR 状态检查阻塞合并
- 开发者立即收到邮件通知

### 10.2 Nightly 失败
- GitHub Actions 发送失败通知
- 需要次日首个开发者处理

## 11. 文档维护

- 本文档随 CI 配置变更同步更新
- 新增测试类型需评估对时延预算的影响
- 定期回顾（建议每季度）时延预算的合理性

---

**创建时间**: 2026-02-26
**基于版本**: `.github/workflows/` (当前 HEAD)
**下次审查**: 2026-Q2
