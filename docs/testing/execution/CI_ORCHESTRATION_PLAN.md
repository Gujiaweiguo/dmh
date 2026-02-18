# DMH CI 编排对齐方案

> 版本: Plan v1 (FRP-006)
> 变更: add-full-regression-testing
> 生成日期: 2026-02-18

---

## 1. 概述

### 1.1 目的

本方案定义 DMH 项目的 CI 编排对齐策略，确保：

- 全量回归覆盖 release candidate 与受保护分支变更
- 产出单一回归结论（PASS/FAIL），可用于阻断发布
- 与现有 PR 验证流程兼容，不破坏日常开发体验

### 1.2 适用范围

| 场景 | 触发条件 | 期望行为 |
|------|----------|----------|
| Release candidate 构建 | tag pattern: `v*.*.*-rc*` | 执行全量回归，产出单一结论 |
| 主干受保护分支变更 | push to `main`/`master` | 执行全量回归，产出单一结论 |
| 手动触发 | workflow_dispatch | 执行全量回归，产出单一结论 |
| 普通 PR 验证 | pull_request | 各 workflow 独立运行，不聚合 |

### 1.3 引用规格

- **全量回归定义**: `docs/testing/execution/FULL_REGRESSION_DEFINITION.md` (FRP-002)
- **发布阻断规则**: `docs/testing/execution/RELEASE_BLOCKING_RULES.md` (FRP-003)
- **Flaky 测试策略**: `docs/testing/execution/FLAKY_TEST_STRATEGY.md` (FRP-004)
- **证据保留策略**: `docs/testing/execution/EVIDENCE_RETENTION_POLICY.md` (FRP-005)
- **Delta 规格**: `openspec/changes/add-full-regression-testing/specs/system-test-execution/spec.md`

---

## 2. 现有 CI 现状

### 2.1 Workflow 清单

| Workflow | 职责 | 触发条件 | 服务依赖 |
|----------|------|----------|----------|
| `stability-checks.yml` | 后端单元 + Admin 单元 + Security E2E | PR/push to main, workflow_dispatch | MySQL, Redis |
| `system-test-gate.yml` | 后端集成 + Admin/H5 E2E + OpenSpec | PR/push to main, workflow_dispatch | MySQL, Redis |
| `coverage-gate.yml` | 三端覆盖率门禁 | PR to main | 无 |
| `order-mysql8-regression.yml` | 订单专项回归 | PR/push (特定路径), workflow_dispatch | MySQL |
| `feedback-guard.yml` | 反馈功能测试 | PR/push (特定路径) | 无 |

### 2.2 覆盖矩阵分析

| 必跑套件 (FRP-002) | 所在 Workflow | 覆盖状态 |
|-------------------|---------------|----------|
| 后端单元测试 | stability-checks.yml | ✅ 覆盖 |
| 后端集成测试 | system-test-gate.yml | ✅ 覆盖 |
| 订单专项回归 | order-mysql8-regression.yml | ⚠️ 仅特定路径触发 |
| Admin 单元测试 | stability-checks.yml | ✅ 覆盖 |
| Admin E2E 测试 | system-test-gate.yml | ✅ 覆盖 |
| H5 单元测试 | ❌ 未覆盖 | ❌ 缺失 |
| H5 E2E 测试 | system-test-gate.yml | ✅ 覆盖 |
| 后端覆盖率 | coverage-gate.yml | ⚠️ 仅 PR 触发 |
| Admin 覆盖率 | coverage-gate.yml | ⚠️ 仅 PR 触发 |
| H5 覆盖率 | coverage-gate.yml | ⚠️ 仅 PR 触发 |
| OpenSpec 校验 | system-test-gate.yml | ✅ 覆盖 |

### 2.3 现有问题

| 问题 | 影响 | 严重程度 |
|------|------|----------|
| 无 RC 触发 | release candidate 无法自动触发全量回归 | 高 |
| 无聚合结论 | 无法产出单一 PASS/FAIL 阻断发布 | 高 |
| H5 单元测试缺失 | 全量回归不完整 | 中 |
| coverage-gate 仅 PR 触发 | push to main 时不检查覆盖率 | 中 |
| order-mysql8-regression 路径限制 | 非订单路径变更时不触发订单回归 | 低 |

### 2.4 触发条件覆盖分析

```
当前触发条件覆盖:

                    ┌─────────────────────────────────────────────────────┐
                    │                  触发场景                            │
                    ├─────────────┬─────────────┬─────────────┬───────────┤
                    │ PR to main  │ Push to main│ RC tag      │ 手动触发  │
Workflow            │             │             │ v*.*.*-rc*  │           │
├───────────────────┼─────────────┼─────────────┼─────────────┼───────────┤
stability-checks    │ ✅          │ ✅          │ ❌          │ ✅        │
system-test-gate    │ ✅          │ ✅          │ ❌          │ ✅        │
coverage-gate       │ ✅          │ ❌          │ ❌          │ ❌        │
order-mysql8-reg    │ ✅ (路径)   │ ✅ (路径)   │ ❌          │ ✅        │
feedback-guard      │ ✅ (路径)   │ ✅ (路径)   │ ❌          │ ❌        │
├───────────────────┼─────────────┼─────────────┼─────────────┼───────────┤
全量回归需求        │ N/A         │ ✅ 必需     │ ✅ 必需     │ ✅ 必需   │
└───────────────────┴─────────────┴─────────────┴─────────────┴───────────┘
```

---

## 3. 全量回归触发条件

### 3.1 触发条件设计

```
全量回归触发条件:

┌─────────────────────────────────────────────────────────────────────┐
│                        触发场景                                      │
├───────────────────────────────┬─────────────────────────────────────┤
│ 1. Release Candidate 构建     │ tag pattern: v*.*.*-rc*            │
│                               │ 示例: v1.0.0-rc1, v2.3.4-rc10      │
├───────────────────────────────┼─────────────────────────────────────┤
│ 2. 主干受保护分支变更          │ push to main 或 master             │
│                               │ (合并 PR 后自动触发)                │
├───────────────────────────────┼─────────────────────────────────────┤
│ 3. 手动触发                   │ workflow_dispatch                  │
│                               │ (发布会议前签收、紧急修复验证)       │
└───────────────────────────────┴─────────────────────────────────────┘
```

### 3.2 Tag Pattern 定义

| 模式 | 正则表达式 | 说明 |
|------|------------|------|
| RC tag | `v[0-9]+\.[0-9]+\.[0-9]+-rc[0-9]+` | Release candidate 版本 |
| 正式版本 | `v[0-9]+\.[0-9]+\.[0-9]+` | 生产发布版本（可选触发） |

### 3.3 触发条件 YAML 配置

```yaml
# 建议的全量回归 workflow 触发配置
on:
  # 1. Release candidate 构建
  push:
    tags:
      - 'v[0-9]+.[0-9]+.[0-9]+-rc[0-9]+'
  
  # 2. 主干受保护分支变更
  push:
    branches:
      - main
      - master
  
  # 3. 手动触发
  workflow_dispatch:
    inputs:
      reason:
        description: '触发原因'
        required: false
        default: 'Manual full regression'
```

### 3.4 不触发全量回归的场景

| 场景 | 原因 | 处理方式 |
|------|------|----------|
| 普通 PR 验证 | 由 stability-checks / system-test-gate 覆盖 | 保持现有行为 |
| Feature 分支 push | 非受保护分支 | 不触发 |
| 文档变更 | 无代码行为变更 | paths-ignore 排除 |
| 配置微调 | 非功能性变更 | paths-ignore 排除 |

---

## 4. 单一结论聚合方案

### 4.1 聚合逻辑定义

```
全量回归结论 =

  IF (stability-checks == success 
      AND system-test-gate == success 
      AND coverage-gate == success 
      AND order-mysql8-regression == success)
  THEN PASS
  ELSE FAIL

特殊情况:
  - 任一必跑 workflow 未执行 → INCONCLUSIVE (需人工排查)
  - 存在隔离中的 flaky 测试 → PASS (含隔离说明)
```

### 4.2 方案对比

| 方案 | 描述 | 优点 | 缺点 | 推荐度 |
|------|------|------|------|--------|
| **方案 A: 新建聚合 Workflow** | 创建 `full-regression.yml`，通过 `workflow_call` 调用其他 workflows | 职责清晰、易维护、不修改现有 workflow | 需要重构现有 workflow 为可调用 | ⭐⭐⭐⭐⭐ |
| **方案 B: 在现有 workflow 中添加聚合 job** | 在各 workflow 末尾添加结论输出，创建汇总 job | 改动较小 | 耦合度高、难以维护 | ⭐⭐ |
| **方案 C: 使用 GitHub 原生 workflow_run 事件** | 监听其他 workflow 完成事件，触发聚合 | 无需修改现有 workflow | 时序复杂、调试困难 | ⭐⭐⭐ |

### 4.3 推荐方案: 方案 A - 新建聚合 Workflow

#### 4.3.1 架构设计

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     full-regression.yml (聚合 Workflow)                  │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                      aggregate-conclusion job                     │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │ 输入: 各子 workflow 的结论                                    │ │   │
│  │  │ 输出: 单一 PASS/FAIL 结论 + 阻断报告                          │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│              │           │           │           │                      │
│              ▼           ▼           ▼           ▼                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │
│  │ stability-  │ │ system-     │ │ coverage-   │ │ order-      │        │
│  │ checks.yml  │ │ test-gate   │ │ gate.yml    │ │ mysql8-reg  │        │
│  │ (callable)  │ │ (callable)  │ │ (callable)  │ │ (callable)  │        │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 4.3.2 子 workflow 改造要点

将现有 workflow 改造为可被调用（`workflow_call`）：

```yaml
# stability-checks.yml 改造示例
on:
  workflow_call:
    inputs:
      trigger-source:
        description: '调用来源'
        required: false
        default: 'direct'
        type: string
    outputs:
      conclusion:
        description: 'Workflow 执行结论'
        value: ${{ jobs.stability-checks.outputs.conclusion }}

jobs:
  stability-checks:
    runs-on: ubuntu-latest
    outputs:
      conclusion: ${{ steps.set-conclusion.outputs.conclusion }}
    # ... 现有步骤 ...
    
    steps:
      # ... 现有步骤 ...
      
      - name: Set workflow conclusion
        id: set-conclusion
        if: always()
        run: |
          echo "conclusion=${{ job.status }}" >> $GITHUB_OUTPUT
```

#### 4.3.3 聚合 Workflow 实现

```yaml
# .github/workflows/full-regression.yml
name: Full Regression

on:
  push:
    tags:
      - 'v[0-9]+.[0-9]+.[0-9]+-rc[0-9]+'
  push:
    branches:
      - main
      - master
  workflow_dispatch:
    inputs:
      reason:
        description: '触发原因'
        required: false
        default: 'Manual full regression'

jobs:
  stability-checks:
    uses: ./.github/workflows/stability-checks.yml
    with:
      trigger-source: 'full-regression'
  
  system-test-gate:
    uses: ./.github/workflows/system-test-gate.yml
    with:
      trigger-source: 'full-regression'
  
  coverage-gate:
    uses: ./.github/workflows/coverage-gate.yml
    with:
      trigger-source: 'full-regression'
  
  order-mysql8-regression:
    uses: ./.github/workflows/order-mysql8-regression.yml
    with:
      trigger-source: 'full-regression'
  
  aggregate-conclusion:
    name: Aggregate Regression Conclusion
    needs: [stability-checks, system-test-gate, coverage-gate, order-mysql8-regression]
    if: always()
    runs-on: ubuntu-latest
    
    outputs:
      verdict: ${{ steps.aggregate.outputs.verdict }}
    
    steps:
      - name: Aggregate workflow results
        id: aggregate
        run: |
          STABILITY="${{ needs.stability-checks.outputs.conclusion || needs.stability-checks.result }}"
          SYSTEM_TEST="${{ needs.system-test-gate.outputs.conclusion || needs.system-test-gate.result }}"
          COVERAGE="${{ needs.coverage-gate.outputs.conclusion || needs.coverage-gate.result }}"
          ORDER_REG="${{ needs.order-mysql8-regression.outputs.conclusion || needs.order-mysql8-regression.result }}"
          
          echo "stability-checks: $STABILITY"
          echo "system-test-gate: $SYSTEM_TEST"
          echo "coverage-gate: $COVERAGE"
          echo "order-mysql8-regression: $ORDER_REG"
          
          # 判定逻辑
          if [[ "$STABILITY" == "success" && \
                "$SYSTEM_TEST" == "success" && \
                "$COVERAGE" == "success" && \
                "$ORDER_REG" == "success" ]]; then
            echo "verdict=PASS" >> $GITHUB_OUTPUT
            echo "## ✅ Full Regression: PASS" >> $GITHUB_STEP_SUMMARY
          else
            echo "verdict=FAIL" >> $GITHUB_OUTPUT
            echo "## ❌ Full Regression: FAIL" >> $GITHUB_STEP_SUMMARY
            echo "" >> $GITHUB_STEP_SUMMARY
            echo "| Workflow | Status |" >> $GITHUB_STEP_SUMMARY
            echo "|----------|--------|" >> $GITHUB_STEP_SUMMARY
            echo "| stability-checks | $STABILITY |" >> $GITHUB_STEP_SUMMARY
            echo "| system-test-gate | $SYSTEM_TEST |" >> $GITHUB_STEP_SUMMARY
            echo "| coverage-gate | $COVERAGE |" >> $GITHUB_STEP_SUMMARY
            echo "| order-mysql8-regression | $ORDER_REG |" >> $GITHUB_STEP_SUMMARY
            exit 1
          fi
      
      - name: Upload regression report
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: full-regression-report-${{ github.run_id }}
          path: |
            regression-report.md
          retention-days: ${{ steps.aggregate.outputs.verdict == 'PASS' && 90 || 180 }}
```

### 4.4 结论到阻断动作的映射

| 聚合结论 | 阻断动作 | 后续流程 |
|----------|----------|----------|
| PASS | 无阻断 | 允许发布，记录证据 |
| FAIL | 立即阻断 | 输出阻断报告，等待修复 |
| INCONCLUSIVE | 立即阻断 | 需人工排查环境问题 |

---

## 5. Workflow 改造建议

### 5.1 改造清单

| Workflow | 改造项 | 优先级 | 影响范围 |
|----------|--------|--------|----------|
| `stability-checks.yml` | 添加 workflow_call + 输出结论 + 添加 H5 单元测试 | P0 | 高 |
| `system-test-gate.yml` | 添加 workflow_call + 输出结论 | P0 | 高 |
| `coverage-gate.yml` | 添加 workflow_call + 输出结论 + 扩展触发条件 | P0 | 高 |
| `order-mysql8-regression.yml` | 添加 workflow_call + 输出结论 | P1 | 中 |
| `full-regression.yml` | 新建聚合 workflow | P0 | 高 |

### 5.2 stability-checks.yml 改造详情

```yaml
# 改造要点:
# 1. 添加 workflow_call 触发器
# 2. 添加 job outputs
# 3. 添加 H5 单元测试 (缺失项)

on:
  workflow_call:
    inputs:
      trigger-source:
        required: false
        default: 'direct'
        type: string
    outputs:
      conclusion:
        value: ${{ jobs.stability-checks.outputs.conclusion }}
  workflow_dispatch:  # 保留原有触发器
  pull_request:       # 保留原有触发器
  push:               # 保留原有触发器

jobs:
  stability-checks:
    outputs:
      conclusion: ${{ steps.final-status.outputs.conclusion }}
    
    steps:
      # ... 现有步骤 ...
      
      # 新增: H5 单元测试
      - name: Setup Node for H5
        uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
          cache-dependency-path: frontend-h5/package-lock.json
      
      - name: Install H5 dependencies
        working-directory: frontend-h5
        run: npm ci
      
      - name: Run frontend-h5 unit tests
        working-directory: frontend-h5
        run: npm run test
      
      # 新增: 输出结论
      - name: Set final status
        id: final-status
        if: always()
        run: |
          echo "conclusion=${{ job.status }}" >> $GITHUB_OUTPUT
```

### 5.3 coverage-gate.yml 改造详情

```yaml
# 改造要点:
# 1. 添加 workflow_call 触发器
# 2. 添加 push to main 触发器
# 3. 添加 job outputs

on:
  workflow_call:
    inputs:
      trigger-source:
        required: false
        default: 'direct'
        type: string
    outputs:
      conclusion:
        value: ${{ jobs.aggregate-coverage.outputs.conclusion }}
  pull_request:  # 保留原有触发器
    branches:
      - main
  push:          # 新增触发器
    branches:
      - main

jobs:
  # ... 现有 jobs ...
  
  # 新增: 聚合结论 job
  aggregate-coverage:
    name: Aggregate Coverage Result
    needs: [backend-coverage, frontend-admin-coverage, frontend-h5-coverage]
    if: always()
    runs-on: ubuntu-latest
    outputs:
      conclusion: ${{ steps.aggregate.outputs.conclusion }}
    
    steps:
      - name: Check all coverage gates
        id: aggregate
        run: |
          if [[ "${{ needs.backend-coverage.result }}" == "success" && \
                "${{ needs.frontend-admin-coverage.result }}" == "success" && \
                "${{ needs.frontend-h5-coverage.result }}" == "success" ]]; then
            echo "conclusion=success" >> $GITHUB_OUTPUT
          else
            echo "conclusion=failure" >> $GITHUB_OUTPUT
            exit 1
          fi
```

### 5.4 order-mysql8-regression.yml 改造详情

```yaml
# 改造要点:
# 1. 添加 workflow_call 触发器
# 2. 全量回归时移除路径限制
# 3. 添加 job outputs

on:
  workflow_call:
    inputs:
      trigger-source:
        required: false
        default: 'direct'
        type: string
    outputs:
      conclusion:
        value: ${{ jobs.order-mysql8-regression.outputs.conclusion }}
  workflow_dispatch:  # 保留
  # PR 触发保留路径限制
  pull_request:
    paths:
      - backend/api/internal/logic/order/**
      # ... 其他路径 ...
  # Push 到 main 保留路径限制 (非全量回归时)
  push:
    branches:
      - main
    paths:
      - backend/api/internal/logic/order/**
      # ... 其他路径 ...
```

---

## 6. 兼容性与回滚

### 6.1 兼容性设计

| 兼容性要求 | 实现方式 |
|------------|----------|
| 保持 PR 验证流程不变 | workflow_call 和原有触发器并存 |
| 不破坏现有 workflow 功能 | 改造采用增量式，保留原有逻辑 |
| 支持独立运行和被调用 | 通过 inputs 区分调用来源 |
| 保持路径过滤功能 | PR 触发时保留 paths 限制 |

### 6.2 触发器共存策略

```yaml
# 触发器优先级和共存策略

on:
  # 1. 被其他 workflow 调用 (最高优先级)
  workflow_call:
    inputs: ...
    outputs: ...
  
  # 2. 手动触发
  workflow_dispatch:
  
  # 3. PR 触发 (带路径过滤)
  pull_request:
    branches: [main]
    paths: [...]
  
  # 4. Push 触发 (带路径过滤)
  push:
    branches: [main]
    paths: [...]

# 在 job 中判断调用来源:
# if: inputs.trigger-source == 'full-regression' || github.event_name == 'push'
```

### 6.3 回滚方案

#### 场景 1: 聚合 workflow 问题

```bash
# 回滚步骤:
# 1. 禁用 full-regression.yml 的自动触发
# 2. 删除或重命名 full-regression.yml
# 3. 恢复到各 workflow 独立运行模式

# 操作:
mv .github/workflows/full-regression.yml .github/workflows/full-regression.yml.disabled
git commit -m "chore: disable full-regression workflow for rollback"
git push
```

#### 场景 2: 子 workflow 改造问题

```bash
# 回滚步骤:
# 1. 移除 workflow_call 触发器
# 2. 移除 outputs 定义
# 3. 恢复原始 workflow 结构

# 每个受影响的 workflow 需单独回滚
git revert <commit-hash>
```

#### 场景 3: 完全回滚

```bash
# 回滚到改造前的 commit
git revert <initial-refactor-commit>..<HEAD>
```

### 6.4 回滚验证清单

- [ ] 各 workflow 可独立触发
- [ ] PR 验证流程正常
- [ ] 手动触发正常
- [ ] 无残留的 workflow_call 依赖

---

## 7. 实施步骤

### 7.1 实施顺序

```
阶段 1: 准备工作 (1-2 天)
├── 1.1 创建 feature 分支
├── 1.2 备份现有 workflow 文件
└── 1.3 设计详细改造清单

阶段 2: 子 workflow 改造 (2-3 天)
├── 2.1 改造 coverage-gate.yml
│   ├── 添加 workflow_call
│   ├── 添加 push to main 触发器
│   └── 添加 outputs
├── 2.2 改造 stability-checks.yml
│   ├── 添加 workflow_call
│   ├── 添加 H5 单元测试
│   └── 添加 outputs
├── 2.3 改造 system-test-gate.yml
│   ├── 添加 workflow_call
│   └── 添加 outputs
└── 2.4 改造 order-mysql8-regression.yml
    ├── 添加 workflow_call
    └── 添加 outputs

阶段 3: 聚合 workflow 创建 (1 天)
├── 3.1 创建 full-regression.yml
├── 3.2 配置触发条件
├── 3.3 实现聚合逻辑
└── 3.4 配置阻断报告输出

阶段 4: 验证与发布 (1-2 天)
├── 4.1 测试 PR 验证流程
├── 4.2 测试 push to main 触发
├── 4.3 测试 RC tag 触发
├── 4.4 测试手动触发
├── 4.5 验证聚合结论正确性
└── 4.6 合并到主分支
```

### 7.2 验证检查点

| 检查点 | 验证方法 | 预期结果 |
|--------|----------|----------|
| PR 验证不变 | 创建测试 PR | 各 workflow 独立运行，路径过滤生效 |
| workflow_call 可用 | 手动触发 full-regression | 所有子 workflow 被正确调用 |
| 聚合结论正确 | 模拟成功/失败场景 | PASS/FAIL 结论与预期一致 |
| 阻断报告生成 | 模拟失败场景 | 报告包含所有失败明细 |
| RC 触发生效 | 推送测试 RC tag | full-regression.yml 自动触发 |

### 7.3 发布检查清单

- [ ] 所有子 workflow 改造完成并测试通过
- [ ] full-regression.yml 创建并测试通过
- [ ] PR 验证流程未受影响
- [ ] 聚合结论逻辑正确
- [ ] 阻断报告格式符合预期
- [ ] 证据保留配置正确 (90/180 天)
- [ ] 文档更新完成

---

## 8. 附录

### 8.1 必跑 Workflow 状态到结论映射表

| stability-checks | system-test-gate | coverage-gate | order-mysql8-regression | 聚合结论 |
|------------------|------------------|---------------|-------------------------|----------|
| ✅ success | ✅ success | ✅ success | ✅ success | ✅ PASS |
| ❌ failure | - | - | - | ❌ FAIL |
| - | ❌ failure | - | - | ❌ FAIL |
| - | - | ❌ failure | - | ❌ FAIL |
| - | - | - | ❌ failure | ❌ FAIL |
| ⚠️ skipped | - | - | - | ⚠️ INCONCLUSIVE |
| - | ⚠️ skipped | - | - | ⚠️ INCONCLUSIVE |

### 8.2 术语表

| 术语 | 定义 |
|------|------|
| **全量回归** | 覆盖必跑矩阵所有层级的发布前验证活动，产出单一结论 |
| **聚合 Workflow** | 调用其他 workflows 并聚合其结果的顶层 workflow |
| **workflow_call** | GitHub Actions 的 workflow 复用机制，允许一个 workflow 调用另一个 |
| **单一结论** | 全量回归的原子输出，值为 PASS/FAIL/INCONCLUSIVE 之一 |
| **RC tag** | Release Candidate 标签，格式为 v*.*.*-rc* |
| **阻断报告** | 发布被阻断时自动生成的失败明细报告 |

### 8.3 参考资料

- **全量回归口径定义**: `docs/testing/execution/FULL_REGRESSION_DEFINITION.md` (FRP-002)
- **发布阻断规则**: `docs/testing/execution/RELEASE_BLOCKING_RULES.md` (FRP-003)
- **Flaky 测试策略**: `docs/testing/execution/FLAKY_TEST_STRATEGY.md` (FRP-004)
- **证据保留策略**: `docs/testing/execution/EVIDENCE_RETENTION_POLICY.md` (FRP-005)
- **Delta 规格**: `openspec/changes/add-full-regression-testing/specs/system-test-execution/spec.md`
- **GitHub Actions Workflow 复用**: https://docs.github.com/en/actions/using-workflows/reusing-workflows

---

*此文档由 FRP-006 任务生成，定义 DMH 项目的 CI 编排对齐方案。*
