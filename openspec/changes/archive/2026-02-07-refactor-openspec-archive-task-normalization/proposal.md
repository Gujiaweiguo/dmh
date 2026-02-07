# Change: Normalize OpenSpec archive task status records

## Why
- 当前 archive 中部分 `tasks.md` 存在状态表达不一致（如重复章节、同一任务 `[ ]` 与 `[x]` 并存、说明文字与勾选状态冲突），影响历史可读性与追溯准确性。
- 这些问题会增加后续排查成本，也会干扰对 active change 的真实进度判断。

## What Changes
- 规范化 archive 中已识别问题文件的任务状态表达，清理重复标题与冲突条目。
- 对“可选未执行”与“后续补齐”任务增加明确注释，避免误读为遗漏执行。
- 新增 archive 状态索引文档，集中记录归档变更的最终口径与关联说明。

## Impact
- Affected specs: spec-governance
- Affected docs:
  - openspec/changes/archive/2026-02-02-order-logic-implementation/order-logic-implementation/tasks.md
  - openspec/changes/archive/2026-01-28-fix-distributor-view-architecture/tasks.md
  - openspec/changes/archive/2026-01-24-enhance-rbac-permission-system/tasks.md
  - openspec/changes/archive/ARCHIVE_STATUS_INDEX.md (new)
- Verification:
  - openspec validate refactor-openspec-archive-task-normalization --strict --no-interactive
