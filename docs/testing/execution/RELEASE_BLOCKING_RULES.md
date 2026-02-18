# DMH 发布阻断规则定义

> 版本: Plan v1 (FRP-003)
> 变更: add-full-regression-testing
> 生成日期: 2026-02-18

---

## 1. 阻断规则概述

本文档定义 DMH 项目的**发布阻断机制**，将全量回归判定规则映射到 CI workflow 状态，确保：

- 自动阻断条件可被 CI 系统识别并执行
- 阻断报告包含可追溯的失败明细
- 硬阻断与软阻断有明确的处理流程

**适用范围**：
- Release candidate 构建
- 主干受保护分支（main/master）变更
- 发布流水线签收点

**引用基线**：
- 全量回归口径定义：`docs/testing/execution/FULL_REGRESSION_DEFINITION.md`
- 现有门禁规格：`openspec/specs/system-test-execution/spec.md`
- Delta 规格：`openspec/changes/add-full-regression-testing/specs/system-test-execution/spec.md`

---

## 2. 阻断条件分类

### 2.1 硬阻断（Hard Block）

**定义**：立即阻断发布，不允许人工豁免，必须修复后重新触发回归。

| 阻断级别 | 触发条件 | 阻断动作 | 修复要求 |
|----------|----------|----------|----------|
| **硬阻断** | 任一必跑套件失败 | 立即阻断 | 修复代码，重新触发全量回归 |
| **硬阻断** | 任一必跑套件未执行 | 立即阻断 | 修复配置/环境问题，重新触发 |
| **硬阻断** | 后端覆盖率 < 76% | 立即阻断 | 补充测试或调整代码 |
| **硬阻断** | Admin 覆盖率 < 70% | 立即阻断 | 补充测试或调整代码 |
| **硬阻断** | H5 覆盖率 < 44% | 立即阻断 | 补充测试或调整代码 |
| **硬阻断** | OpenSpec 校验存在错误 | 立即阻断 | 修复规格冲突或格式错误 |

**必跑套件清单**（来自 FRP-002）：

```
必跑矩阵:
├── 后端单元测试: go test ./...
├── 后端集成测试: go test ./test/integration/... -v -count=1
├── 订单专项回归: backend/scripts/run_order_mysql8_regression.sh
├── Admin 单元测试: npm run test
├── Admin E2E 测试: npm run test:e2e:headless
├── H5 单元测试: npm run test
├── H5 E2E 测试: npm run test:e2e:headless
└── OpenSpec 校验: openspec validate --all --no-interactive
```

### 2.2 软阻断（Soft Block）

**定义**：阻断发布，但可申请风险豁免，需审批通过后方可放行。

| 阻断级别 | 触发条件 | 阻断动作 | 豁免要求 |
|----------|----------|----------|----------|
| **软阻断** | P1 用例通过率 < 95% | 阻断，可申请豁免 | 需填写风险豁免单，获得 Tech Lead + QA 审批 |
| **软阻断** | 存在 Critical 缺陷未关闭 | 阻断，可申请豁免 | 需评估缺陷影响范围，获得 Tech Lead + Product 审批 |
| **软阻断** | 隔离清单内有 flaky 测试未修复 | 阻断，可申请豁免 | 需确认 flaky 测试不影响核心功能，记录技术债 |

**软阻断豁免流程**：

1. 提交风险豁免申请（使用 `RISK_WAIVER_TEMPLATE.md`）
2. Tech Lead 评估技术风险
3. QA/Product 评估业务风险
4. 双方签字批准后放行
5. 发布后 7 天内完成遗留问题修复

---

## 3. CI 状态映射

### 3.1 Workflow 到阻断条件映射

| Workflow | 成功状态 | 失败状态 | 映射到阻断条件 |
|----------|----------|----------|----------------|
| `stability-checks.yml` | ✅ 通过 | ❌ 失败 | 后端单元/Admin 单元失败 → **硬阻断** |
| `system-test-gate.yml` | ✅ 通过 | ❌ 失败 | 集成测试/E2E/OpenSpec 失败 → **硬阻断** |
| `coverage-gate.yml` | ✅ 通过 | ❌ 失败 | 覆盖率不达标 → **硬阻断** |
| `order-mysql8-regression.yml` | ✅ 通过 | ❌ 失败 | 订单回归失败 → **硬阻断** |

### 3.2 Workflow 状态到结论映射

| Workflow 组合 | 聚合结论 | 阻断级别 |
|---------------|----------|----------|
| 所有 workflow ✅ | PASS | 无阻断 |
| 任一必跑 workflow ❌ | FAIL | 硬阻断 |
| workflow 因环境问题未执行 | INCONCLUSIVE | 硬阻断（需排查） |

### 3.3 状态判定逻辑

```
if (all_workflows == success):
    verdict = PASS
    block_level = NONE
    
elif (any_mandatory_workflow == failed OR any_mandatory_workflow == not_run):
    verdict = FAIL
    block_level = HARD_BLOCK
    
elif (coverage_threshold_not_met):
    verdict = FAIL
    block_level = HARD_BLOCK
    
elif (p1_pass_rate < 95% OR has_critical_defects):
    verdict = FAIL_WITH_WAIVER_OPTION
    block_level = SOFT_BLOCK
    
else:
    verdict = INCONCLUSIVE
    block_level = HARD_BLOCK  # 需人工排查
```

---

## 4. 阻断报告模板

当发布被阻断时，系统自动生成以下格式的报告：

```markdown
# 发布阻断报告

## 基本信息
- 报告ID: [UUID 或 CI Run ID]
- 触发时间: [ISO 8601 格式，如 2026-02-18T14:30:00Z]
- 触发提交: [Git SHA，如 abc123def456]
- 触发分支: [分支名，如 main]
- 触发事件: [push / release / manual]

## 阻断原因
- 阻断级别: [硬阻断 / 软阻断]
- 阻断条件: [具体条件描述，如 "后端单元测试失败"]
- 失败套件列表:
  - [套件名称]: [状态] - [workflow 链接]
  - 示例: 后端单元测试: FAILED - https://github.com/xxx/actions/runs/12345

## 失败明细
| 用例ID | 用例名称 | 失败原因 | 所属套件 | 日志链接 |
|--------|----------|----------|----------|----------|
| test-001 | TestOrderCreate | 断言失败: expected 200, got 500 | 后端单元 | [link] |
| test-002 | TestUserLogin | 超时: connection refused | 后端集成 | [link] |

## 覆盖率状态（如适用）
| 模块 | 当前覆盖率 | 阈值 | 状态 |
|------|-----------|------|------|
| 后端 | 72.5% | 76% | ❌ 不达标 |
| Admin | 68.0% | 70% | ❌ 不达标 |
| H5 | 45.0% | 44% | ✅ 达标 |

## 重跑记录
| 重跑时间 | 重跑原因 | 结果 | 操作人 |
|----------|----------|------|--------|
| 2026-02-18T14:35:00Z | 环境瞬态问题 | FAILED | [automated] |
| 2026-02-18T14:40:00Z | 第二次重试 | FAILED | [automated] |

## 处理建议
[针对当前阻断的具体修复建议]

### 常见修复路径
- 后端单元测试失败 → 检查代码变更，修复断言或实现
- 集成测试失败 → 检查 API 服务、数据库连接状态
- E2E 测试失败 → 检查前端构建、路由配置
- 覆盖率不达标 → 补充单元测试覆盖新代码
- OpenSpec 校验失败 → 检查规格文件格式与冲突

## 证据保留
- 证据保留至: [日期，失败回归保留 180 天]
- 证据检索: commit:[SHA] workflow:[workflow-name]
```

---

## 5. 放行与豁免流程

### 5.1 硬阻断处理

**不允许豁免**，必须完成以下步骤：

1. **定位问题**：根据阻断报告定位失败原因
2. **修复代码**：提交修复 PR
3. **重新触发回归**：PR 合并后自动触发，或手动触发
4. **确认通过**：所有必跑套件 100% 通过

### 5.2 软阻断处理

**可申请风险豁免**，流程如下：

```
┌─────────────────┐
│ 1. 填写豁免申请  │
│ RISK_WAIVER_    │
│ TEMPLATE.md     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 2. Tech Lead    │
│ 评估技术风险    │
│ 签字确认        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 3. QA/Product   │
│ 评估业务风险    │
│ 签字确认        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 4. 双方批准后   │
│ 临时放行发布    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 5. 发布后 7 天  │
│ 内完成遗留修复  │
└─────────────────┘
```

### 5.3 豁免申请必填字段

```markdown
# 风险豁免申请

## 基本信息
- 申请日期: [日期]
- 申请人: [姓名]
- 发布版本: [版本号]
- 阻断原因: [引用阻断报告 ID]

## 风险评估
- 遗留问题: [具体描述未满足的门禁条件]
- 影响范围: [评估对用户/系统的影响]
- 缓解措施: [发布后如何降低风险]

## 审批记录
- Tech Lead: [姓名] - [批准/拒绝] - [日期]
- QA/ Product: [姓名] - [批准/拒绝] - [日期]

## 承诺修复
- 修复截止: [日期]
- 责任人: [姓名]
```

---

## 6. 与现有门禁的关系

### 6.1 引用关系

本文档是**阻断机制落地层**，不重复定义已有门禁阈值：

| 来源 | 门禁/规则 | 本文档如何使用 |
|------|-----------|----------------|
| `openspec/specs/system-test-execution/spec.md` | Release QA gate: P0=100%, P1>=95%, 无 Blocker/Critical | 作为软阻断的触发条件 |
| `docs/testing/execution/FULL_REGRESSION_DEFINITION.md` | 必跑矩阵 + PASS/FAIL 判定规则 | 作为硬阻断的触发条件 |
| Delta spec | 自动阻断 + 阻断报告字段 | 定义报告模板与 CI 映射 |

### 6.2 层级关系

```
┌─────────────────────────────────────────────────────────────┐
│                  发布阻断规则定义 (本文档)                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐│
│  │  硬阻断/软阻断   │  │  CI 状态映射    │  │  阻断报告    ││
│  └────────┬────────┘  └────────┬────────┘  └──────┬───────┘│
└───────────┼────────────────────┼──────────────────┼────────┘
            │                    │                  │
            ▼                    ▼                  ▼
┌───────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ 全量回归口径定义   │  │ 现有门禁规格     │  │ Delta 规格      │
│ (FRP-002)         │  │ (P0/P1 阈值)    │  │ (自动阻断场景)  │
└───────────────────┘  └─────────────────┘  └─────────────────┘
```

### 6.3 不重复定义的内容

以下内容由引用文档定义，本文档不重复：

- P0/P1 用例分级标准 → 引用 `test-plan.md`
- 覆盖率阈值 → 引用 FRP-002
- 必跑套件清单 → 引用 FRP-002
- Flaky 测试隔离策略 → 引用 Delta spec

---

## 7. 术语表

| 术语 | 定义 |
|------|------|
| **硬阻断** | 立即阻断发布，不允许人工豁免，必须修复后重新触发回归 |
| **软阻断** | 阻断发布，但可申请风险豁免，需审批通过后放行 |
| **必跑套件** | 全量回归必须执行的测试套件，缺一不可 |
| **阻断报告** | 发布被阻断时自动生成的失败明细报告 |
| **风险豁免** | 在未满足门禁条件下经审批允许发布的例外机制 |
| **CI 状态映射** | GitHub Actions workflow 状态到阻断条件的对应关系 |

---

## 8. 参考资料

- **全量回归口径定义**: `docs/testing/execution/FULL_REGRESSION_DEFINITION.md` (FRP-002)
- **现有规格**: `openspec/specs/system-test-execution/spec.md`
- **Delta 规格**: `openspec/changes/add-full-regression-testing/specs/system-test-execution/spec.md`
- **风险豁免模板**: `docs/testing/execution/RISK_WAIVER_TEMPLATE.md`
- **测试结果模板**: `docs/testing/execution/TEST_RESULT_TEMPLATE.md`

---

*此文档由 FRP-003 任务生成，定义发布阻断规则的落地机制。*
