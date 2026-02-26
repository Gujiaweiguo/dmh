# Task 09: PR Gate 轻量路径改造方案

## 执行时间
- 开始: 2026-02-26
- 结束: 2026-02-26

## 1. 概述

本文档定义 PR Gate 轻量路径改造方案，目标是在保持覆盖率门禁的前提下，将 PR 必跑集合收敛至快反馈路径，确保 PR Gate 目标时间 ≤15 分钟。

## 2. 设计原则

1. **保留覆盖率门禁**：不降低覆盖率阈值（Backend 78%, Admin 80%, H5 70%）
2. **快速反馈优先**：开发者应在 15 分钟内获得 PR 结果
3. **分层策略**：重量级测试移至 Merge Gate 或 Nightly

## 3. PR 必跑集合

### 3.1 后端单元测试（-short 模式）

| 配置项 | 值 | 说明 |
|--------|-----|------|
| 命令 | `go test -short ./...` | 跳过标记为 `// +build !short` 的测试 |
| 时延预算 | ≤ 5 分钟 | 当前约 2-3 分钟 |
| 覆盖率阈值 | ≥ 78% | 保持不变 |
| 超时配置 | 10 分钟 | CI job timeout |

**改造要点**：
- 使用 `-short` 标志跳过长时间运行的测试
- 保持覆盖率收集：`-coverprofile=coverage.out -covermode=atomic`
- 排除集成测试和性能测试目录

**示例命令**：
```bash
go test -short $(go list ./... | grep -v -E 'dmh/test/integration|dmh/test/performance') \
  -coverprofile=coverage.out -covermode=atomic
```

### 3.2 前端单元测试

| 模块 | 时延预算 | 覆盖率阈值 | 超时配置 |
|------|---------|-----------|---------|
| Frontend-Admin | ≤ 5 分钟 | ≥ 80% | 8 分钟 |
| Frontend-H5 | ≤ 5 分钟 | ≥ 70% | 8 分钟 |

**改造要点**：
- 使用 Vitest 的 `--reporter=verbose` 提供清晰输出
- 并行执行两个前端测试（当前实现已并行）
- 保持覆盖率门禁检查

**示例命令**：
```bash
# Admin
cd frontend-admin && npm run test:cov

# H5
cd frontend-h5 && npm run test:cov
```

### 3.3 Lint 检查

| 配置项 | 值 |
|--------|-----|
| 工具 | gofmt |
| 时延预算 | ≤ 2 分钟 |
| 超时配置 | 5 分钟 |

**检查内容**：
- Go 代码格式化：`gofmt -d .`
- 未来可扩展：ESLint（前端）、golangci-lint（后端）

### 3.4 覆盖率门禁

| 模块 | 阈值 | 位置 |
|------|------|------|
| Backend | 78% | pr-gate.yml |
| Frontend-Admin | 80% | pr-gate.yml |
| Frontend-H5 | 70% | pr-gate.yml |

**门禁逻辑**：
```bash
# Backend
if (( $(echo "$COVERAGE < 78" | bc -l) )); then
  echo "::error::Backend coverage ${COVERAGE}% is below 78% threshold"
  exit 1
fi

# Admin
if [ "${PCT%.*}" -ge 80 ]; then
  echo "✓ Admin coverage meets 80% threshold"
else
  exit 1
fi

# H5
if [ "${PCT%.*}" -ge 70 ]; then
  echo "✓ H5 coverage meets 70% threshold"
else
  exit 1
fi
```

## 4. 跳过策略

以下测试类型不在 PR Gate 中运行，移至相应层级：

### 4.1 集成测试 → Merge Gate / Nightly

| 测试类型 | 移至位置 | 触发条件 |
|---------|---------|---------|
| Backend 集成测试 | `system-test-gate.yml` | backend/** 变更 |
| 完整集成套件 | `full-regression.yml` | Nightly |

**原因**：
- 需要 MySQL/Redis 服务容器
- 执行时间较长（~8 分钟）
- 适合在 Merge Gate 或 Nightly 中验证

### 4.2 E2E 测试 → Merge Gate / Nightly

| 测试类型 | 移至位置 | 触发条件 |
|---------|---------|---------|
| Admin E2E Smoke | `system-test-gate.yml` | frontend-admin/** 变更 |
| H5 E2E Smoke | `system-test-gate.yml` | frontend-h5/** 变更 |
| 完整 E2E 套件 | `full-regression.yml` | Nightly |

**原因**：
- 需要 Playwright 浏览器环境
- 执行时间不稳定
- 适合在较低频次运行

### 4.3 性能测试 → Nightly

| 测试类型 | 移至位置 | 触发条件 |
|---------|---------|---------|
| Backend 性能测试 | `full-regression.yml` | Nightly |
| 负载测试 | `full-regression.yml` | Nightly |

**原因**：
- 执行时间长
- 资源消耗大
- 不适合每次 PR 都运行

## 5. 时延预算明细

### 5.1 目标时延

```
┌─────────────────────────────────────────────────────────────┐
│                  PR Gate 时延预算（≤15 分钟）                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────────────────┐   ┌─────────────────┐                │
│   │ Backend Unit    │   │ Frontend Unit   │                │
│   │ (≤5 min)        │   │ Admin (≤5 min)  │                │
│   │ + Coverage 78%  │   │ + Coverage 80%  │                │
│   └─────────────────┘   └─────────────────┘                │
│                                                             │
│   ┌─────────────────┐   ┌─────────────────┐                │
│   │ Frontend Unit   │   │ Lint Check      │                │
│   │ H5 (≤5 min)     │   │ (≤2 min)        │                │
│   │ + Coverage 70%  │   │                 │                │
│   └─────────────────┘   └─────────────────┘                │
│                                                             │
│   并行执行预估: ~5-8 分钟（取决于 CI 并行度）                  │
│   最坏情况: 5 (backend) + 5 (admin) + 5 (h5) + 2 (lint)     │
│           = 17 分钟 → 通过并行控制在 15 分钟内                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 5.2 时延预算分配

| Job | 预估时间 | 超时配置 | 并行度 |
|-----|---------|---------|--------|
| backend-unit | 2-5 分钟 | 10 分钟 | 1（需 MySQL） |
| frontend-unit (Admin) | 2-5 分钟 | 8 分钟 | 与 H5 并行 |
| frontend-unit (H5) | 2-5 分钟 | 8 分钟 | 与 Admin 并行 |
| lint | 1-2 分钟 | 5 分钟 | 独立 |

### 5.3 优化建议

1. **后端优化**：
   - 使用 `-short` 模式跳过长时间测试
   - 考虑移除 MySQL 服务容器（单元测试应 mock 外部依赖）

2. **前端优化**：
   - 保持当前并行执行策略
   - 考虑使用 Vitest 的 `--pool=threads` 提升并行度

3. **缓存优化**：
   - Go 模块缓存：已启用
   - npm 依赖缓存：已启用

## 6. 当前实现状态

### 6.1 pr-gate.yml 现状

| Job | 当前状态 | 改造需求 |
|-----|---------|---------|
| backend-unit | ✅ 包含覆盖率门禁 | 添加 `-short` 标志 |
| frontend-unit | ✅ 包含覆盖率门禁 | 无需修改 |
| lint | ✅ gofmt 检查 | 无需修改 |
| pr-gate-verdict | ✅ 汇总结果 | 无需修改 |

### 6.2 改造清单

- [ ] 后端单元测试添加 `-short` 标志
- [ ] 验证覆盖率门禁仍然生效
- [ ] 确认时延在预算内

## 7. 禁止事项

### 7.1 禁止降低覆盖率阈值

```
❌ Backend < 78%
❌ Admin < 80%
❌ H5 < 70%
```

### 7.2 禁止跳过核心单元测试

```
❌ 跳过所有单元测试
❌ 使用 -coverpkg 虚假覆盖率
❌ 禁用覆盖率门禁检查
```

### 7.3 禁止在 PR Gate 中包含

```
❌ 集成测试 (./test/integration/...)
❌ E2E 测试 (Playwright)
❌ 性能测试 (./test/performance/...)
❌ 需要外部服务容器的测试（除必要的 MySQL 用于 schema 验证）
```

## 8. 验证标准

### 8.1 功能验证

- [ ] 后端单元测试全部通过
- [ ] 前端 Admin 单元测试全部通过
- [ ] 前端 H5 单元测试全部通过
- [ ] Lint 检查通过
- [ ] 覆盖率门禁生效

### 8.2 时延验证

```bash
# 预期结果
Backend:  ≤ 5 分钟 ✓
Admin:   ≤ 5 分钟 ✓
H5:      ≤ 5 分钟 ✓
Lint:    ≤ 2 分钟 ✓
Total:   ≤ 15 分钟 ✓
```

### 8.3 覆盖率验证

```bash
# 预期结果
Backend: ≥ 78% ✓
Admin:   ≥ 80% ✓
H5:      ≥ 70% ✓
```

## 9. 相关文档

- `.sisyphus/evidence/task-02-ci-layer-budget.md` - CI 分层测试策略
- `.github/workflows/pr-gate.yml` - 当前 PR Gate 实现
- `.github/workflows/full-regression.yml` - Nightly 全量回归

---

**创建时间**: 2026-02-26
**目标时延**: PR Gate ≤ 15 分钟
**覆盖率阈值**: 保持不变（78%/80%/70%）
