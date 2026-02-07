## 1. Baseline Audit
- [x] 1.1 复核 archive 中 `tasks.md` 的未勾选项与冲突项清单
- [x] 1.2 将问题分类为：重复条目 / 状态冲突 / 可选未执行

## 2. Archive Normalization
- [x] 2.1 清理 `enhance-rbac-permission-system/tasks.md` 的重复章节与冲突 7.1 状态
- [x] 2.2 修正 `order-logic-implementation/tasks.md` 中说明文字与任务勾选的语义冲突
- [x] 2.3 为 `fix-distributor-view-architecture/tasks.md` 的可选未执行任务补充归档说明

## 3. Traceability
- [x] 3.1 新增 `openspec/changes/archive/ARCHIVE_STATUS_INDEX.md` 作为归档状态索引
- [x] 3.2 在上述 3 个归档任务文件中添加索引引用

## 4. Validation
- [x] 4.1 运行严格校验：`openspec validate refactor-openspec-archive-task-normalization --strict --no-interactive`
- [x] 4.2 复核 active changes 不受影响
