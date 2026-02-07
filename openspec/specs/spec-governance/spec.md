# spec-governance Specification

## Purpose
TBD - created by archiving change refactor-openspec-archive-task-normalization. Update Purpose after archive.
## Requirements
### Requirement: Archive task status consistency
系统 SHALL 在 OpenSpec 归档变更中保持 `tasks.md` 的状态表达一致，避免同一任务出现冲突状态或重复定义。

#### Scenario: Resolve contradictory checklist states in archived tasks
- **WHEN** 归档任务文档中同一任务出现 `[ ]` 与 `[x]` 并存或重复章节
- **THEN** 维护者 SHALL 统一为单一权威状态表达
- **AND** 必须保留解释性说明以体现历史真实性（如“可选未执行”或“后续补齐”）

### Requirement: Archive status traceability index
系统 SHALL 为归档变更提供统一索引，记录关键状态口径和跨文档关联，便于后续审计与追溯。

#### Scenario: Locate final status rationale for archived changes
- **WHEN** 维护者查看任一归档变更的任务状态
- **THEN** 系统 SHALL 提供可导航的索引入口
- **AND** 索引 SHALL 标注状态口径来源与关联变更

