# Test Coverage 80% - Learnings

## 项目上下文

**目标**: 后端 80%+ / 前端 Admin 70%+

**当前状态**:
- 后端整体: 45.6%
- 前端 Admin views: 0.94%
- 前端 H5 logic: ~100%

## 技术选择

- 后端测试: Go testing + testify + httptest
- 前端测试: Vitest + Vue Test Utils
- CI: GitHub Actions

## 关键参考

- Handler 测试模板: `backend/api/internal/handler/auth/handler_test.go`
- Vitest 配置: `frontend-admin/vitest.config.ts`
- CI workflow: `.github/workflows/system-test-gate.yml`

## 执行记录

### 2026-02-15 Wave 1 开始

Session: ses_39f7d8991ffeMi2Gg4p7bAvyLW
任务: 1.1-1.3 并行执行
- 修复 frontend-admin vitest.coverage 收集：覆盖率从仅 unit tests 的覆盖范围扩展至 views、services、components 目录，确保 views 子目录的测试能够被收集。具体变更：
- 在 frontend-admin/vitest.config.ts 中：新增 test.include 条目 views/**/*.{test,spec}.{ts,tsx}；在 coverage.include 新增 ['views/', 'services/', 'components/']；在 coverage 增加 thresholds：lines/functions/branches/statements 均为 70。
- 变更保持向后兼容，不移除原有配置，仅添加必要的收集路径与阈值。
- Created a CI coverage gate workflow: .github/workflows/coverage-gate.yml
- Triggers on pull_request to main with two jobs:
  - backend-coverage: runs backend tests with coverage, enforces >= 80%
  - frontend-coverage: runs frontend tests with coverage, enforces >= 70%
- Uses standard GitHub Actions: checkout@v4, setup-go@v5 (Go 1.24), setup-node@v4 (Node 20)
- Backend command: go test ./... -coverprofile=coverage.out -covermode=atomic
- Frontend command: npm run test:cov
- Coverage checks implemented by parsing coverage outputs and exiting non-zero on failure
- YAML syntax is validated (via a lightweight runtime check during CI) to ensure proper workflow configuration
