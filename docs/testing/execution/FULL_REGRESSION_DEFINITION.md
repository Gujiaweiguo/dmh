# DMH 全量回归口径定义

> 版本: Plan v1 (FRP-002)
> 变更: add-full-regression-testing
> 生成日期: 2026-02-18

---

## 1. 全量回归定义

**全量回归** 是 DMH 项目发布前的最高级别验证活动，定义为：

> **全量回归 = 必跑测试矩阵执行完成 + 单一 PASS/FAIL 结论**

核心特征：
- **完整性**：覆盖后端、前端双端、OpenSpec 的所有必跑层级
- **原子性**：产生且仅产生一个回归结论，无中间态
- **阻断性**：FAIL 结论自动阻断发布流水线
- **可追溯性**：证据保留至少 90 天（失败回归 180 天）

全量回归与日常 PR 验证的区别：

| 维度 | PR 验证 | 全量回归 |
|------|---------|----------|
| 触发 | PR/push 事件 | Release candidate / 主干受保护分支变更 |
| 范围 | 受影响模块相关测试 | 必跑矩阵 100% 覆盖 |
| 结论 | 各 workflow 独立状态 | 单一聚合结论 |
| 阻断 | 阻断 PR 合并 | 阻断发布 |

---

## 2. 必跑测试矩阵

以下测试套件为全量回归的**必跑项**，缺一不可：

```
全量回归必跑矩阵:
├── 后端 (Backend)
│   ├── 单元测试: go test ./...
│   ├── 集成测试: go test ./test/integration/... -v -count=1
│   ├── 订单专项回归: backend/scripts/run_order_mysql8_regression.sh
│   └── 覆盖率门禁: >= 76%
│
├── 前端管理后台 (Frontend-Admin)
│   ├── 单元测试: npm run test
│   ├── E2E 测试: npm run test:e2e:headless
│   └── 覆盖率门禁: >= 70%
│
├── 前端 H5 (Frontend-H5)
│   ├── 单元测试: npm run test
│   ├── E2E 测试: npm run test:e2e:headless
│   └── 覆盖率门禁: >= 44%
│
└── OpenSpec
    └── 校验: openspec validate --all --no-interactive
```

### 2.1 后端必跑项

| 层级 | 命令入口 | 执行环境要求 |
|------|----------|--------------|
| 单元测试 | `cd backend && go test ./...` | Go 1.24+ |
| 集成测试 | `cd backend && go test ./test/integration/... -v -count=1` | 运行中的 API + MySQL + Redis |
| 订单专项回归 | `backend/scripts/run_order_mysql8_regression.sh` | 运行中的 API + MySQL |
| 覆盖率 | `go test ./... -coverprofile=coverage.out` | 阈值 >= 76% |

### 2.2 前端管理后台必跑项

| 层级 | 命令入口 | 执行环境要求 |
|------|----------|--------------|
| 单元测试 | `cd frontend-admin && npm run test` | Node.js 20+ |
| E2E 测试 | `cd frontend-admin && npm run test:e2e:headless` | Headless 浏览器 |
| 覆盖率 | `npm run test:cov` | 阈值 >= 70% |

### 2.3 前端 H5 必跑项

| 层级 | 命令入口 | 执行环境要求 |
|------|----------|--------------|
| 单元测试 | `cd frontend-h5 && npm run test` | Node.js 20+ |
| E2E 测试 | `cd frontend-h5 && npm run test:e2e:headless` | Headless 浏览器 |
| 覆盖率 | `npm run test:cov` | 阈值 >= 44% |

### 2.4 OpenSpec 必跑项

| 层级 | 命令入口 | 执行环境要求 |
|------|----------|--------------|
| 全量校验 | `openspec validate --all --no-interactive` | OpenSpec CLI |

---

## 3. 单一结论判定规则

全量回归产出**唯一结论**，遵循以下判定规则：

### 3.1 PASS 条件

所有以下条件同时满足：

| 条件 | 说明 |
|------|------|
| 后端单元测试 | 100% 通过，无跳过用例 |
| 后端集成测试 | 100% 通过 |
| 订单专项回归 | 100% 通过 |
| 后端覆盖率 | >= 76% |
| Admin 单元测试 | 100% 通过 |
| Admin E2E 测试 | 100% 通过 |
| Admin 覆盖率 | >= 70% |
| H5 单元测试 | 100% 通过 |
| H5 E2E 测试 | 100% 通过 |
| H5 覆盖率 | >= 44% |
| OpenSpec 校验 | 无错误 |

**判定公式**：
```
PASS = (∀ 套件 ∈ 必跑矩阵: 套件.通过) ∧ (∀ 覆盖率检查: 覆盖率 >= 阈值)
```

### 3.2 FAIL 条件

任一以下条件触发：

| 条件 | 说明 |
|------|------|
| 任一必跑套件失败 | 包含断言失败、运行时错误、超时 |
| 任一必跑套件未执行 | 因配置错误、环境问题导致未运行 |
| 任一覆盖率不达标 | 低于定义的阈值 |
| OpenSpec 校验失败 | 存在规格冲突或格式错误 |

**判定公式**：
```
FAIL = (∃ 套件 ∈ 必跑矩阵: ¬套件.通过 ∨ ¬套件.已执行) ∨ (∃ 覆盖率检查: 覆盖率 < 阈值)
```

### 3.3 INCONCLUSIVE 条件

以下情况需人工裁定：

| 条件 | 处理方式 |
|------|----------|
| 部分套件因基础设施问题无法完成 | 标记 INCONCLUSIVE，记录失败原因，触发人工审核 |
| Flaky 测试隔离中 | 隔离清单内的测试不计入 FAIL，但需记录隔离状态 |
| 环境瞬态问题重试中 | 等待重试完成，最多 2 次 |

**判定公式**：
```
INCONCLUSIVE = (部分套件.环境故障) ∧ (¬FAIL) ∧ (需要人工裁定)
```

### 3.4 结论映射到发布门禁

| 回归结论 | 发布门禁动作 |
|----------|--------------|
| PASS | 允许发布，记录证据 |
| FAIL | 阻断发布，输出失败报告 |
| INCONCLUSIVE | 暂停发布，等待人工裁定 |

---

## 4. 与现有规格的关系

本定义是**聚合口径层**，不重复定义已有规格内容：

### 4.1 引用的现有 Requirement

| 来源 | Requirement | 本定义如何使用 |
|------|-------------|----------------|
| `openspec/specs/system-test-execution/spec.md` | Unified multi-layer test execution matrix | 引用其分层结构定义必跑矩阵 |
| `openspec/specs/system-test-execution/spec.md` | Release QA gate for comprehensive testing | 引用 P0=100%、P1>=95% 门禁作为质量基线 |
| `openspec/specs/system-test-execution/spec.md` | Traceable evidence and execution records | 引用证据追溯要求作为证据保留依据 |

### 4.2 新增的 Delta Requirement

| 来源 | Requirement | 本定义如何使用 |
|------|-------------|----------------|
| delta spec | Full regression orchestration workflow | 本定义是其口径落地 |
| delta spec | Flaky test quarantine and retry policy | 本定义引用其重试上限（2次） |
| delta spec | Plan v1 as single execution authority | 本定义遵循 Plan v1 的执行依据 |

### 4.3 关系图

```
┌─────────────────────────────────────────────────────────────┐
│                  全量回归口径定义 (本文档)                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐│
│  │   必跑矩阵定义    │  │  单一结论规则   │  │  触发条件    ││
│  └────────┬────────┘  └────────┬────────┘  └──────┬───────┘│
└───────────┼────────────────────┼──────────────────┼────────┘
            │                    │                  │
            ▼                    ▼                  ▼
┌───────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ 现有 Spec          │  │ Delta Spec       │  │ 基线文档        │
│ (分层执行矩阵)      │  │ (编排/Flaky/证据)│  │ (FRP-001)       │
└───────────────────┘  └─────────────────┘  └─────────────────┘
```

---

## 5. 触发条件

全量回归在以下场景触发：

### 5.1 自动触发

| 触发事件 | 说明 |
|----------|------|
| Release candidate 构建 | 版本发布候选创建时 |
| 主干受保护分支变更 | main/master 分支有新合并 |

### 5.2 手动触发

| 触发场景 | 说明 |
|----------|------|
| 发布前签收 | 发布会议前执行全量回归 |
| 紧急修复验证 | Hotfix 合并后验证 |
| 定期回归 | 夜间/周末定时回归（建议） |

### 5.3 不触发全量回归的场景

| 场景 | 原因 |
|------|------|
| 普通 PR 验证 | 由 stability-checks / system-test-gate 覆盖 |
| 文档变更 | 无代码行为变更 |
| 配置微调 | 非功能性变更 |

---

## 6. 术语表

| 术语 | 定义 |
|------|------|
| **全量回归** | 覆盖必跑矩阵所有层级的发布前验证活动，产出单一结论 |
| **必跑矩阵** | 全量回归必须执行的测试套件集合，缺一不可 |
| **单一结论** | 全量回归的原子输出，值为 PASS/FAIL/INCONCLUSIVE 之一 |
| **阻断性** | FAIL 结论自动阻止发布流水线继续 |
| **Flaky 测试** | 同一条件下结果不确定的测试，需隔离处理 |
| **覆盖率门禁** | 代码覆盖率低于阈值时阻断发布的规则 |
| **Plan v1** | `openspec/changes/add-full-regression-testing/tasks.md` 中定义的执行计划 |
| **风险豁免** | 在未满足门禁条件下经审批允许发布的例外机制 |

---

## 7. 参考资料

- **基线文档**: `docs/testing/execution/FULL_REGRESSION_BASELINE.md` (FRP-001)
- **现有规格**: `openspec/specs/system-test-execution/spec.md`
- **Delta 规格**: `openspec/changes/add-full-regression-testing/specs/system-test-execution/spec.md`
- **执行计划**: `openspec/changes/add-full-regression-testing/tasks.md` (Plan v1)

---

*此文档由 FRP-002 任务生成，定义全量回归的统一口径。*
