# Change: 建立全量回归测试执行标准

## Why
- 现有仓库已具备多条测试链路（`stability-checks.yml`、`system-test-gate.yml`、`coverage-gate.yml`、`order-mysql8-regression.yml`、本地脚本与 Makefile），但“全量回归”尚未被统一定义为单一可执行标准。
- 当前执行口径在本地与 CI 间存在分散：有门禁但缺少统一编排、失败重试策略、证据保留周期和发布阻断判定的一致约束。

## What Changes
- 在 `system-test-execution` 规格中补充“全量回归执行编排”要求，定义必须覆盖的测试层级与统一入口。
- 在发布门禁要求中补充“全量回归必过”的阻断规则，明确与现有 CI 门禁的一致性要求。
- 新增 flakiness 控制要求，约束可重试范围、重试上限和失败判定口径。
- 新增证据保留与审计要求，明确 artifact 最低保留时长、内容和可追溯字段。

## Plan v1 Authority
- 本变更的唯一执行计划为 `openspec/changes/add-full-regression-testing/tasks.md` 中的 **Plan v1**。
- Atlas/Hephaestus 在实现阶段 SHALL 严格按 Plan v1 执行，不得使用独立计划文档或临时口头计划替代。
- 若需变更执行顺序/范围，必须先更新 Plan v1，再进入实现。

## Impact
- Affected specs:
  - `system-test-execution` (modified)
- Affected code/docs (implementation stage):
  - `.github/workflows/system-test-gate.yml`
  - `.github/workflows/stability-checks.yml`
  - `.github/workflows/coverage-gate.yml`
  - `.github/workflows/order-mysql8-regression.yml`
  - `Makefile`
  - `backend/scripts/*.sh`
  - `docs/testing/execution/*`
- Behavior changes:
  - 测试与发布流程治理增强（非业务功能变更）

## Out of Scope
- 不在本提案中新增业务接口或修改业务数据模型。
- 不在本提案中替换现有测试框架（Go test / Vitest / Playwright）。
