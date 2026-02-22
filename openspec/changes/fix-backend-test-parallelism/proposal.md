# Change: Fix Backend Test Parallelism for Database Isolation

## Why

后端集成测试在并行运行时出现数据库竞争问题，导致主键冲突失败：
- `distributor_integration_test.go` 使用共享测试数据库 `dmh_test`
- 多个测试包并行执行时，`createTestData()` 的 "检查再创建" 模式存在竞态条件
- 导致 `Error 1062: Duplicate entry 'X' for key 'users.PRIMARY'` 错误

当前临时解决方案是使用 `-p 1` 顺序运行测试，但这增加了 CI 时间。

## What Changes

1. **规范化测试运行模式**：在 `system-test-execution` 规格中添加后端测试并行性要求
2. **CI/CD 配置更新**：确保 CI 工作流使用正确的测试运行模式
3. **文档更新**：在 AGENTS.md 和 README 中记录测试运行规范

## Impact

- **Affected specs**: `system-test-execution`
- **Affected code**: 
  - `backend/test/integration/*.go` (已修复 distributor_integration_test.go)
  - `.github/workflows/*.yml` (可能需要更新)
  - `AGENTS.md` (添加测试运行说明)

## Risks

- 低风险：仅涉及测试配置和文档，不影响生产代码
- CI 时间可能略有增加（使用 `-p 1` 顺序运行）

## Success Criteria

1. `go test -p 1 ./...` 在后端目录稳定通过
2. CI 工作流使用规范化测试命令
3. 文档清晰记录测试运行规范
