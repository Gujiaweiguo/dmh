# DMH 测试体系完整优化计划

## TL;DR

> **Quick Summary**: 通过“分层测试重构 + CI 分层执行 + 测试环境标准化 + AGENTS 规范固化”四条主线，降低测试耗时与不稳定性，同时保持覆盖率门禁与回归质量。
>
> **Deliverables**:
> - 后端/前端测试分层执行蓝图（Unit / Integration / E2E / Performance）
> - Repository 测试容器方案（MySQL8）与落地步骤
> - CI 分层工作流改造方案（PR 快速反馈 + Nightly 全量）
> - `AGENTS.md` 测试执行规范增补内容
>
> **Estimated Effort**: Large
> **Parallel Execution**: YES - 4 waves + Final Verification
> **Critical Path**: T1 → T3 → T7 → T12 → F1-F4

---

## Context

### Original Request
用户要求：给出完整测试优化计划，并明确纳入“把测试相关信息写入 `AGENTS.md`”的执行项。

### Interview Summary
**Key Discussions**:
- 需要完整计划，不是零散建议。
- 要覆盖后端/前端/CI/集成/E2E/性能测试。
- 需要把测试约束固化到 `AGENTS.md` 作为长期协作规则。

**Research Findings**:
- 已有 CI 覆盖率门禁（backend 78%、admin 80%、h5 70%）。
- 后端存在较多 DB/Redis 依赖导致测试慢与易波动。
- 集成测试受运行环境影响，易 SKIP 或超时。

### Metis Review
**Identified Gaps** (addressed):
- Metis 会话超时，采用保守缺口策略：将“范围锁定、验收标准、边界风险、回滚策略”显式写入任务与最终验证波次。

---

## Work Objectives

### Core Objective
构建一套可持续、可扩展、可审计的 DMH 测试体系，确保“开发期快反馈、合并前稳门禁、夜间全量回归”三层质量闭环。

### Concrete Deliverables
- 单元测试职责边界与重构路径（handler/logic/repository）。
- 集成测试环境标准化方案（MySQL8 + 可选 Redis + API 启动检查）。
- E2E 与性能测试调度策略（PR/Nightly 分流）。
- `AGENTS.md` 测试约定更新文本（命令、分层、禁用反模式）。

### Definition of Done
- [ ] `AGENTS.md` 中新增测试执行规范段落并通过评审。
- [ ] PR 流程完成“快速门禁”并能在目标时长内稳定通过。
- [ ] Nightly 流程完成“全量验证”并有失败告警路径。
- [ ] 至少 2 个代表模块完成分层示范（含可复用模板）。

### Must Have
- 继续维持现有覆盖率阈值，不降门槛。
- 新增策略不破坏现有主干发布节奏。
- 所有优化任务必须可量化验收。

### Must NOT Have (Guardrails)
- 不以“跳过测试”替代问题修复。
- 不在 handler 层新增业务逻辑。
- 不将本应 nightly 的慢测强塞到 PR gate。
- 不把测试计划拆成多个独立计划文件。

---

## Verification Strategy (MANDATORY)

> **ZERO HUMAN INTERVENTION** — 所有验收均可由执行代理通过命令/工具验证。

### Test Decision
- **Infrastructure exists**: YES
- **Automated tests**: YES (Tests-after，按任务补齐)
- **Framework**: Go test / Vitest / Playwright / GitHub Actions
- **If TDD**: 非强制全局 TDD；对新增或高风险模块建议局部 TDD。

### QA Policy
每个任务必须产出可执行验证步骤，证据落盘：`.sisyphus/evidence/task-{N}-{slug}.{ext}`。

- **Backend**: `go test -p 1 ...` + 目标包验证
- **Frontend**: `npm run test:cov`（admin/h5 分开）
- **E2E**: Playwright 场景执行 + 报告上传
- **CI**: workflow dry-run/触发结果与 summary 校验

---

## Execution Strategy

### Parallel Execution Waves

Wave 1（基础策略与基线）:
- T1 测试资产盘点与基线指标冻结
- T2 CI 分层目标与时延预算定义
- T3 Repository 测试容器方案定稿（MySQL8）
- T4 AGENTS.md 测试规范草案
- T5 测试证据归档规范（evidence naming）

Wave 2（核心改造）:
- T6 后端分层测试模板（handler/logic/repository）
- T7 集成测试环境标准化（API/MySQL/可选Redis）
- T8 前端单测与E2E职责边界重划
- T9 PR Gate 轻量路径改造
- T10 Nightly 全量回归工作流改造

Wave 3（样板模块与扩展）:
- T11 选择2个后端模块做分层示范改造
- T12 集成测试稳定性治理（跳过原因治理/重试策略）
- T13 性能测试分流（PR short + Nightly full）
- T14 测试数据工厂/夹具规范化

Wave 4（落地与制度化）:
- T15 AGENTS.md 最终落地与团队宣贯检查项
- T16 失败分诊与回滚机制固化
- T17 KPI 看板（时长、稳定性、跳过率、失败率）

Wave FINAL（并行评审）:
- F1 计划符合性审计
- F2 质量与反模式审计
- F3 场景执行与证据审计
- F4 范围与变更污染审计

Critical Path: T1 → T3 → T7 → T12 → F1-F4

### Dependency Matrix (FULL)
- T1: — → T6,T7,T8,T9,T10
- T2: — → T9,T10,T17
- T3: — → T7,T11,T14
- T4: — → T15
- T5: — → T9,T10,T15,F3
- T6: T1 → T11
- T7: T1,T3 → T12,T16
- T8: T1 → T13
- T9: T1,T2,T5 → T16
- T10: T1,T2,T5 → T16,T17
- T11: T3,T6 → T17
- T12: T7 → T16,T17
- T13: T8 → T17
- T14: T3 → T11,T12
- T15: T4,T5 → F1,F4
- T16: T7,T9,T10,T12 → F1,F2,F4
- T17: T2,T10,T11,T12,T13 → F1,F3
- F1: T15,T16,T17 → done
- F2: T16 → done
- F3: T17 → done
- F4: T15,T16 → done

### Agent Dispatch Summary
- Wave1: `quick/unspecified-high/writing`
- Wave2: `deep/quick/writing/unspecified-high`
- Wave3: `deep/quick/unspecified-high`
- Wave4: `writing/deep/unspecified-high`
- Final: `oracle/unspecified-high/deep`

---

## TODOs

- [ ] 1. 测试资产盘点与基线冻结（T1）
  **What to do**: 统计 backend/admin/h5 的测试数量、耗时、跳过率、失败率，固化为基线表。
  **Must NOT do**: 不调整阈值、不删除现有测试。
  **Recommended Agent Profile**: Category=`quick`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave1；Blocks: T6/T7/T8/T9/T10。
  **References**: `.github/workflows/pr-gate.yml`, `docs/testing/TEST_IMPROVEMENT_PLAN.md`。
  **Acceptance Criteria**: 产出基线表（含时间戳、命令、结果）。
  **QA Scenarios**: Happy=`go test -p 1 ./... -short` 汇总成功；Error=缺失模块时返回明确缺项清单。

- [ ] 2. CI 分层目标与时延预算定义（T2）
  **What to do**: 定义 PR/merge/nightly 三层测试责任和时延预算。
  **Must NOT do**: 不把 nightly 慢测塞进 PR gate。
  **Recommended Agent Profile**: Category=`writing`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave1；Blocks: T9/T10/T17。
  **References**: `.github/workflows/pr-gate.yml`, `.github/workflows/full-regression.yml`。
  **Acceptance Criteria**: 预算表（PR≤15m、Nightly≤45m）落地到计划文档。
  **QA Scenarios**: Happy=workflow 对应分层正确；Error=发现职责冲突时列冲突矩阵。

- [ ] 3. Repository 容器测试方案定稿（T3）
  **What to do**: 明确使用 MySQL8 test container（非 sqlite），并给出事务回滚策略。
  **Must NOT do**: 不引入与生产不一致的数据引擎替代。
  **Recommended Agent Profile**: Category=`deep`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave1；Blocks: T7/T11/T14。
  **References**: `backend/migrations/*.sql`, `backend/api/internal/testutil/mysql_test_helper.go`。
  **Acceptance Criteria**: 文档写明选型理由、性能权衡、回退方案。
  **QA Scenarios**: Happy=容器启动+迁移+事务回滚通过；Error=容器不可用时降级为跳过并报警。

- [ ] 4. AGENTS 测试规范草案（T4）
  **What to do**: 起草测试执行约定（命令、分层职责、禁止项、故障排查入口）。
  **Must NOT do**: 不覆盖项目既有开发规范。
  **Recommended Agent Profile**: Category=`writing`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave1；Blocks: T15。
  **References**: `AGENTS.md`, `backend/AGENTS.md`。
  **Acceptance Criteria**: 草案包含“何时跑 -short / 何时跑全量 / CI与本地差异”。
  **QA Scenarios**: Happy=团队按草案可一键复现；Error=命令缺失时有补充占位清单。

- [ ] 5. 证据归档命名规范（T5）
  **What to do**: 统一 `.sisyphus/evidence/task-{N}-{slug}` 证据命名。
  **Must NOT do**: 不使用临时无规则文件名。
  **Recommended Agent Profile**: Category=`quick`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave1；Blocks: T9/T10/T15/F3。
  **References**: `.sisyphus/plans/test-optimization-master-plan.md`。
  **Acceptance Criteria**: 证据模板被各任务复用。
  **QA Scenarios**: Happy=抽样3项证据命名正确；Error=命名不合规时校验失败。

- [ ] 6. 后端分层测试模板（T6）
  **What to do**: 给出 handler(mock logic)/logic(mock repo)/repo(real mysql8) 三层模板。
  **Must NOT do**: 不新增 handler 业务逻辑。
  **Recommended Agent Profile**: Category=`deep`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave2；Blocked By: T1；Blocks: T11。
  **References**: `backend/api/internal/handler`, `backend/api/internal/logic`。
  **Acceptance Criteria**: 模板示例可直接复制到新模块。
  **QA Scenarios**: Happy=示例模块通过三层测试；Error=mock 断言失败能定位到层级。

- [ ] 7. 集成测试环境标准化（T7）
  **What to do**: 固化 API 启动检查、MySQL 初始化、可选 Redis 可用性探测。
  **Must NOT do**: 不默默 SKIP 关键场景。
  **Recommended Agent Profile**: Category=`unspecified-high`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave2；Blocked By: T1,T3；Blocks: T12,T16。
  **References**: `backend/test/integration`, `backend/scripts/run_order_mysql8_regression.sh`。
  **Acceptance Criteria**: 集成测试执行率目标 >90%（可解释性失败除外）。
  **QA Scenarios**: Happy=服务就绪后集成套件可运行；Error=服务不可达时快速失败并提示修复命令。

- [ ] 8. 前端单测/E2E 边界重划（T8）
  **What to do**: 明确组件逻辑放单测、跨页面主流程放 E2E。
  **Must NOT do**: 不把视觉像素级问题塞到单测。
  **Recommended Agent Profile**: Category=`visual-engineering`, Skills=`["playwright"]`。
  **Parallelization**: 并行=YES，Wave2；Blocked By: T1；Blocks: T13。
  **References**: `frontend-admin/tests/unit`, `frontend-admin/e2e`, `frontend-h5/tests/unit`, `frontend-h5/e2e`。
  **Acceptance Criteria**: 用例分类表完成并映射到目录。
  **QA Scenarios**: Happy=随机抽样用例分类正确；Error=分类冲突时输出迁移建议。

- [ ] 9. PR Gate 轻量路径改造（T9）
  **What to do**: 保留覆盖率门禁，收敛 PR 必跑集合至快反馈路径。
  **Must NOT do**: 不降低覆盖率阈值。
  **Recommended Agent Profile**: Category=`quick`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave2；Blocked By: T1,T2,T5；Blocks: T16。
  **References**: `.github/workflows/pr-gate.yml`, `.github/workflows/coverage-gate.yml`。
  **Acceptance Criteria**: PR gate 在目标时间内稳定运行。
  **QA Scenarios**: Happy=PR workflow 全绿；Error=阈值低于门槛时明确 fail。

- [ ] 10. Nightly 全量回归改造（T10）
  **What to do**: 将慢测（集成/E2E/性能）迁移到 nightly，产出汇总报告。
  **Must NOT do**: 不与 PR workflow 重复执行同一重任务。
  **Recommended Agent Profile**: Category=`unspecified-high`, Skills=`["playwright"]`。
  **Parallelization**: 并行=YES，Wave2；Blocked By: T1,T2,T5；Blocks: T16,T17。
  **References**: `.github/workflows/full-regression.yml`, `.github/workflows/system-test-gate.yml`。
  **Acceptance Criteria**: Nightly 报告包含失败分类与链接。
  **QA Scenarios**: Happy=nightly 产出完整 artifacts；Error=任一子任务失败时 summary 明确归因。

- [ ] 11. 两个后端模块分层示范（T11）
  **What to do**: 选择2个高频模块（建议 order/distributor）完成分层改造样板。
  **Must NOT do**: 不跨模块做无关重构。
  **Recommended Agent Profile**: Category=`deep`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave3；Blocked By: T3,T6；Blocks: T17。
  **References**: `backend/api/internal/handler/order`, `backend/api/internal/handler/distributor`。
  **Acceptance Criteria**: 样板模块具备可复制脚手架。
  **QA Scenarios**: Happy=样板模块测试全部通过；Error=分层边界被破坏时静态检查报错。

- [ ] 12. 集成测试稳定性治理（T12）
  **What to do**: 治理 flaky：重试仅限幂等步骤、增加就绪探针、明确 SKIP 原因标签。
  **Must NOT do**: 不使用无限重试掩盖缺陷。
  **Recommended Agent Profile**: Category=`unspecified-high`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave3；Blocked By: T7；Blocks: T16,T17。
  **References**: `backend/test/integration/*.go`。
  **Acceptance Criteria**: flaky 率下降并有周报指标。
  **QA Scenarios**: Happy=重跑结果稳定；Error=外部依赖故障时在10分钟内失败并给行动建议。

- [ ] 13. 性能测试分流策略（T13）
  **What to do**: PR 仅跑 short/perf-smoke，full benchmark 放 nightly。
  **Must NOT do**: 不在 PR 执行长达10s+基准用例。
  **Recommended Agent Profile**: Category=`quick`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave3；Blocked By: T8；Blocks: T17。
  **References**: `backend/test/performance/benchmark_test.go`。
  **Acceptance Criteria**: PR 性能测试耗时显著下降且无漏检关键 smoke。
  **QA Scenarios**: Happy=`-short` 路径<1m；Error=full 用例误入 PR 时 gate 拒绝。

- [ ] 14. 测试数据工厂/夹具规范（T14）
  **What to do**: 建立统一 fixture/factory 规则，减少重复造数代码。
  **Must NOT do**: 不在每个测试内手写重复初始化。
  **Recommended Agent Profile**: Category=`quick`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave3；Blocked By: T3；Blocks: T11,T12。
  **References**: `backend/api/internal/testutil`, `backend/model/*_test.go`。
  **Acceptance Criteria**: 至少两类实体有 factory 模板。
  **QA Scenarios**: Happy=factory 可复用；Error=字段缺失时构造器返回明确错误。

- [ ] 15. AGENTS.md 测试规范最终落地（T15）
  **What to do**: 将测试执行规则写入 `AGENTS.md` 与 `backend/AGENTS.md`（命令、分层、禁用反模式、故障排查）。
  **Must NOT do**: 不与现有规则冲突，不删除原有关键约束。
  **Recommended Agent Profile**: Category=`writing`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave4；Blocked By: T4,T5；Blocks: F1,F4。
  **References**: `AGENTS.md`, `backend/AGENTS.md`。
  **Acceptance Criteria**: 两份 AGENTS 均包含“何时跑 -short、何时跑全量、CI 对应关系、常见失败排查”。
  **QA Scenarios**: Happy=按 AGENTS 可无歧义执行；Error=发现冲突条款时列冲突并回滚到安全版本。

- [ ] 16. 失败分诊与回滚机制（T16）
  **What to do**: 建立失败分类（代码/环境/数据/外部依赖）与回滚/重试策略。
  **Must NOT do**: 不把所有失败统一标记为 flaky。
  **Recommended Agent Profile**: Category=`deep`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave4；Blocked By: T7,T9,T10,T12；Blocks: F1,F2,F4。
  **References**: `.github/workflows/*.yml`, `docs/testing/*`。
  **Acceptance Criteria**: 每类失败都有标准处理路径与时限。
  **QA Scenarios**: Happy=模拟失败可被正确分诊；Error=分诊缺项时阻断发布流程。

- [ ] 17. KPI 看板与治理节奏（T17）
  **What to do**: 输出周度 KPI（PR 时长、nightly 通过率、skip率、flaky率）。
  **Must NOT do**: 不只给覆盖率单指标。
  **Recommended Agent Profile**: Category=`writing`, Skills=`[]`。
  **Parallelization**: 并行=YES，Wave4；Blocked By: T2,T10,T11,T12,T13；Blocks: F1,F3。
  **References**: CI summary artifacts, `docs/testing/TEST_IMPROVEMENT_PLAN.md`。
  **Acceptance Criteria**: KPI 面板模板固定并可追溯。
  **QA Scenarios**: Happy=可生成周报；Error=指标数据缺失时给数据源缺口列表。

---

## Final Verification Wave (MANDATORY)

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT`

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Output: `Build/Lint/Test summary | Anti-pattern scan | VERDICT`

- [ ] F3. **Real QA Replay** — `unspecified-high` (+`playwright` if UI)
  Output: `Scenarios [N/N] | Evidence completeness | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`
  Output: `Tasks compliant [N/N] | Scope creep [0/N] | VERDICT`

---

## Commit Strategy
- 以“测试体系优化”为主题分组提交：配置类、测试类、文档类分开。

## Success Criteria
### Verification Commands
```bash
cd backend && go test -p 1 ./... -short -count=1
cd frontend-admin && npm run test:cov
cd frontend-h5 && npm run test:cov
```

### Final Checklist
- [ ] 分层策略文档可执行且可复用
- [ ] AGENTS 测试规范更新完成
- [ ] PR 与 Nightly 分层执行稳定
- [ ] 关键模块示范完成并可复制
