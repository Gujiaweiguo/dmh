# 测试覆盖度提升至 80% 工作计划

## TL;DR

> **Quick Summary**: 系统性提升 DMH 项目测试覆盖率至 80%，通过补充后端 Handler 测试、前端 Admin 视图测试，并配置 CI 覆盖率门禁，完成测试闭环。
> 
> **Deliverables**:
> - 后端测试覆盖率: 45.6% → 80%+
> - 前端 Admin 测试覆盖率: 0.94% → 70%+
> - CI 覆盖率门禁配置
> - OpenSpec change proposal 创建与归档
> 
> **Estimated Effort**: Large (2-3 周)
> **Parallel Execution**: YES - 3 waves (后端 Handlers / 前端 Views / CI 基础设施)
> **Critical Path**: 基础设施 → 后端 Handlers → 前端 Views → 集成验证

---

## Context

### Original Request
> 生成一个openspec的changes,目标是提高测试覆盖度到80%以上，生成详细执行计划，目标是完成测试闭环。
> 计划必须包含：任务拆解、每步输入/输出文件、验收标准与失败回滚策略、风险点与优先级

### Interview Summary

**当前覆盖率状态**:
| 模块 | 当前 | 目标 | 差距 |
|------|------|------|------|
| 后端整体 | 45.6% | 80% | +34.4% |
| 后端 Handler 层 | 0-67% | 80% | +13-80% |
| 后端 Logic 层 | 75-100% | 80% | ✅ 基本达标 |
| frontend-admin views | 0.94% | 70% | +69% |
| frontend-h5 logic | ~100% | 80% | ✅ 达标 |

**关键低覆盖模块**:
- 后端: `api/internal/handler/*` (111 文件, 仅 17 测试文件), `cmd/` (0%), `svc/` (0%)
- 前端: `frontend-admin/views/*` (几乎所有组件 0%)

**技术选择**:
- 后端: Go testing + testify + httptest + sqlite in-memory
- 前端: Vitest + Vue Test Utils + Playwright
- CI: GitHub Actions

### Metis Review

**识别的关键差距** (已解决):
- 后端 Handler: 94 个文件无独立测试
- frontend-admin: 12 个 views 无单元测试，vitest 配置未收集 views 目录
- CI 缺少覆盖率门禁
- 缺少覆盖率报告追踪

**建议**:
- 设置阶段性目标: 60% → 70% → 80%
- 优先核心业务 API (auth, order, campaign, reward)
- 使用现有测试模式 (handler/auth/handler_test.go) 作为模板

---

## Work Objectives

### Core Objective
系统性提升 DMH 项目测试覆盖率至 80%，建立可持续的测试基础设施，完成从需求到归档的完整测试闭环。

### Concrete Deliverables
1. **后端测试文件**: 40+ 新增 handler 测试文件
2. **前端测试文件**: 12+ 新增 view 测试文件
3. **CI 配置**: 覆盖率门禁 workflow
4. **OpenSpec**: 完整的 change proposal 创建与归档

### Definition of Done
- [ ] `cd backend && go test ./... -cover` 显示 ≥80%
- [ ] `cd frontend-admin && npm run test:cov` 显示 ≥70%
- [ ] CI PR 检查中包含覆盖率门禁
- [ ] OpenSpec change 已创建并验证通过
- [ ] 所有测试通过，无 flaky tests

### Must Have
- 后端 Handler 层核心 API 测试 (auth, order, campaign, reward, distributor)
- 前端核心视图测试 (Login, CampaignList, UserManagement)
- CI 覆盖率门禁阻断低覆盖率 PR
- OpenSpec change proposal

### Must NOT Have (Guardrails)
- ❌ 重构 handler 代码以"使其更可测试"
- ❌ 为自动生成代码 (types, routes) 写测试
- ❌ 添加性能/压力测试
- ❌ 添加安全渗透测试
- ❌ 创建新的测试框架或抽象
- ❌ 范围膨胀到其他模块

---

## Verification Strategy (MANDATORY)

### Test Decision
- **Infrastructure exists**: YES (testify, Vitest, Playwright)
- **Automated tests**: YES (Tests-after 模式)
- **Framework**: Go testing + testify / Vitest

### Agent-Executed QA Scenarios (MANDATORY — ALL tasks)

**覆盖率验证场景**:

```
Scenario: Backend coverage meets 80% threshold
  Tool: Bash
  Preconditions: Backend tests implemented
  Steps:
    1. cd /opt/code/DMH/backend
    2. go test ./... -coverprofile=coverage.out -covermode=atomic
    3. go tool cover -func=coverage.out | grep total
    4. Assert: percentage >= 80.0
  Expected Result: "total: (statements) 80.x%"
  Evidence: coverage.out file

Scenario: Frontend-admin coverage meets 70% threshold
  Tool: Bash
  Preconditions: Frontend tests implemented
  Steps:
    1. cd /opt/code/DMH/frontend-admin
    2. npm run test:cov 2>&1 | tail -20
    3. Extract "All files" line percentage
    4. Assert: Stmts >= 70%
  Expected Result: Coverage report shows >= 70%
  Evidence: coverage/coverage-final.json

Scenario: CI coverage gate blocks low-coverage PR
  Tool: Bash
  Preconditions: CI workflow configured
  Steps:
    1. Check .github/workflows/coverage-gate.yml exists
    2. Verify threshold configuration (80% for backend, 70% for frontend)
    3. Assert: workflow has "Check coverage" step
  Expected Result: Workflow file contains coverage gate
  Evidence: .github/workflows/coverage-gate.yml content
```

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately - 基础设施):
├── Task 1.1: 修复 vitest 覆盖率收集配置
├── Task 1.2: 创建 handler 测试模板
└── Task 1.3: 配置 CI 覆盖率门禁

Wave 2 (After Wave 1 - 并行执行):
├── Task 2.1: 后端 auth handler 测试 [depends: 1.2]
├── Task 2.2: 后端 order handler 测试 [depends: 1.2]
├── Task 2.3: 后端 campaign handler 测试 [depends: 1.2]
├── Task 2.4: 后端 reward/distributor handler 测试 [depends: 1.2]
└── Task 2.5: 前端核心 views 测试 [depends: 1.1]

Wave 3 (After Wave 2 - 集成验证):
├── Task 3.1: 验证后端覆盖率 ≥80%
├── Task 3.2: 验证前端覆盖率 ≥70%
├── Task 3.3: 创建 OpenSpec change proposal
└── Task 3.4: 归档与文档

Critical Path: 1.2 → 2.1-2.4 → 3.1 → 3.3
Parallel Speedup: ~40% faster than sequential
```

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 1.1 | None | 2.5 | 1.2, 1.3 |
| 1.2 | None | 2.1-2.4 | 1.1, 1.3 |
| 1.3 | None | 3.3 | 1.1, 1.2 |
| 2.1 | 1.2 | 3.1 | 2.2, 2.3, 2.4, 2.5 |
| 2.2 | 1.2 | 3.1 | 2.1, 2.3, 2.4, 2.5 |
| 2.3 | 1.2 | 3.1 | 2.1, 2.2, 2.4, 2.5 |
| 2.4 | 1.2 | 3.1 | 2.1, 2.2, 2.3, 2.5 |
| 2.5 | 1.1 | 3.2 | 2.1, 2.2, 2.3, 2.4 |
| 3.1 | 2.1-2.4 | 3.3 | 3.2 |
| 3.2 | 2.5 | 3.3 | 3.1 |
| 3.3 | 1.3, 3.1, 3.2 | 3.4 | None |
| 3.4 | 3.3 | None | None |

### Agent Dispatch Summary

| Wave | Tasks | Recommended Agents |
|------|-------|-------------------|
| 1 | 1.1-1.3 | `task(category='quick')` - 配置修改 |
| 2 | 2.1-2.5 | `task(category='unspecified-low')` - 测试编写，可并行 |
| 3 | 3.1-3.4 | `task(category='quick')` - 验证与归档 |

---

## TODOs

- [ ] 1. 基础设施准备

  **What to do**:
  - 修复 vitest 配置，确保收集 views 目录覆盖率
  - 创建可复用的 handler 测试模板
  - 配置 CI 覆盖率门禁 workflow

  **Must NOT do**:
  - 修改生产代码
  - 引入新的测试框架

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 配置修改和小规模代码生成
  - **Skills**: None required

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with 1.1, 1.2, 1.3)
  - **Blocks**: Wave 2 tasks
  - **Blocked By**: None

  **References**:
  - `backend/api/internal/handler/auth/handler_test.go:23-76` - 现有 handler 测试模式
  - `frontend-admin/vitest.config.ts:9-12` - 现有 coverage 配置
  - `.github/workflows/system-test-gate.yml` - 现有 CI workflow

  **输入/输出文件**:
  | 任务 | 输入 | 输出 |
  |------|------|------|
  | 1.1 | `frontend-admin/vitest.config.ts` | 修改后的 vitest.config.ts |
  | 1.2 | `handler/auth/handler_test.go` | 新建 `backend/api/internal/handler/testutil/testutil.go` |
  | 1.3 | `.github/workflows/*.yml` | 新建 `.github/workflows/coverage-gate.yml` |

  **Acceptance Criteria**:
  - [ ] `frontend-admin/vitest.config.ts` 包含 views 目录
  - [ ] `backend/api/internal/handler/testutil/testutil.go` 存在且可导入
  - [ ] `.github/workflows/coverage-gate.yml` 存在且语法正确

  **失败回滚策略**:
  - 1.1 失败: 恢复原始 vitest.config.ts
  - 1.2 失败: 删除 testutil 目录
  - 1.3 失败: 删除 coverage-gate.yml

  **风险点与优先级**:
  | 风险 | 级别 | 缓解措施 |
  |------|------|----------|
  | vitest 配置不兼容 | 低 | 参考 frontend-h5 配置 |
  | CI workflow 语法错误 | 低 | 本地验证后提交 |

  **Commit**: YES (groups with 1)
  - Message: `chore(test): setup coverage infrastructure`
  - Files: vitest.config.ts, testutil.go, coverage-gate.yml

---

- [ ] 2. 后端 Handler 测试 (核心模块)

  **What to do**:
  - 2.1: auth handler 测试补充
  - 2.2: order handler 测试补充
  - 2.3: campaign handler 测试补充
  - 2.4: reward/distributor handler 测试补充

  **测试文件清单**:
  | 模块 | 现有测试 | 需补充 |
  |------|---------|--------|
  | auth | ✅ handler_test.go | 扩展边界场景 |
  | order | ✅ handler_test.go | 补充 verify/complete 场景 |
  | campaign | ✅ handler_test.go | 补充 create/update 场景 |
  | reward | ⚠️ 部分覆盖 | 新增完整测试 |
  | distributor | ⚠️ 部分覆盖 | 新增完整测试 |

  **Must NOT do**:
  - 修改 handler 业务逻辑
  - 为已弃用 API 写测试
  - 重复集成测试覆盖的场景

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: 大量重复性测试编写工作
  - **Skills**: None required

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with 2.1-2.5)
  - **Blocks**: 3.1
  - **Blocked By**: 1.2

  **References**:
  - `backend/api/internal/handler/auth/handler_test.go` - 测试模式模板
  - `backend/api/internal/logic/auth/auth_logic_test.go` - Logic 层测试参考
  - `backend/test/integration/*_test.go` - 集成测试参考

  **输入/输出文件**:
  | 任务 | 输入 | 输出 |
  |------|------|------|
  | 2.1 | `handler/auth/*.go` | `handler/auth/handler_test.go` (扩展) |
  | 2.2 | `handler/order/*.go` | `handler/order/handler_test.go` (扩展) |
  | 2.3 | `handler/campaign/*.go` | `handler/campaign/handler_test.go` (扩展) |
  | 2.4 | `handler/reward/*.go`, `handler/distributor/*.go` | 新建测试文件 |

  **Acceptance Criteria**:
  - [ ] `go test ./api/internal/handler/... -v` 通过
  - [ ] 每个模块至少 5 个测试场景 (正常 + 异常)
  - [ ] Handler 测试文件数量 ≥ 20

  **失败回滚策略**:
  - 单文件失败: 跳过该文件，继续其他
  - 模块失败: 标记 TODO，创建 follow-up issue
  - 整体失败: 回退到现有测试，降低目标

  **风险点与优先级**:
  | 风险 | 级别 | 优先级 | 缓解措施 |
  |------|------|--------|----------|
  | Mock 复杂度高 | 中 | P1 | 使用 testutil 共享 fixture |
  | 外部依赖 (WeChat Pay) | 高 | P2 | 使用 interface mock |
  | 测试运行时间过长 | 低 | P3 | 并行执行 |

  **Commit**: YES (groups with 2)
  - Message: `test(backend): add handler tests for core modules`
  - Files: `backend/api/internal/handler/*_test.go`

---

- [ ] 3. 前端 Admin Views 测试

  **What to do**:
  - 为核心视图组件添加单元测试
  - 优先级: Login > CampaignList > UserManagement > 其他

  **测试文件清单**:
  | 视图 | 当前状态 | 优先级 |
  |------|---------|--------|
  | LoginView.tsx | ❌ 0% | P0 |
  | CampaignListView.tsx | ❌ 0% | P0 |
  | UserManagementView.tsx | ⚠️ 部分 | P0 |
  | BrandManagementView.tsx | ⚠️ 部分 | P1 |
  | OrderListView.tsx | ❌ 0% | P1 |

  **Must NOT do**:
  - 测试 DOM 结构细节
  - 重复 API service 测试
  - 为纯展示组件写复杂测试

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: 前端组件测试编写
  - **Skills**: None required

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with 2.1-2.4)
  - **Blocks**: 3.2
  - **Blocked By**: 1.1

  **References**:
  - `frontend-admin/tests/unit/LoginView.test.ts` - 现有测试模式
  - `frontend-admin/tests/unit/UserManagementView.test.ts` - 现有测试模式
  - `frontend-h5/tests/unit/*.test.js` - H5 测试模式参考

  **输入/输出文件**:
  | 任务 | 输入 | 输出 |
  |------|------|------|
  | 2.5 | `views/*.tsx` | `tests/unit/*View.test.ts` |

  **Acceptance Criteria**:
  - [ ] `npm run test` 通过
  - [ ] 至少 5 个核心视图有测试
  - [ ] 覆盖率 ≥ 70%

  **失败回滚策略**:
  - 单组件失败: 跳过，创建 follow-up issue
  - 整体未达标: 降低目标到 60%

  **风险点与优先级**:
  | 风险 | 级别 | 优先级 | 缓解措施 |
  |------|------|--------|----------|
  | Vue Test Utils 兼容性 | 中 | P1 | 参考 H5 项目模式 |
  | 组件依赖复杂 | 高 | P2 | 使用 shallowMount + mock |

  **Commit**: YES (groups with 2)
  - Message: `test(frontend): add view component tests`
  - Files: `frontend-admin/tests/unit/*.test.ts`

---

- [ ] 4. 集成验证与 OpenSpec 归档

  **What to do**:
  - 验证后端覆盖率 ≥80%
  - 验证前端覆盖率 ≥70%
  - 创建 OpenSpec change proposal
  - 归档变更

  **Must NOT do**:
  - 跳过覆盖率验证
  - 提交未通过验证的 proposal

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 验证和文档工作
  - **Skills**: None required

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential (final)
  - **Blocks**: None
  - **Blocked By**: 3.1, 3.2

  **References**:
  - `openspec/AGENTS.md` - OpenSpec 流程
  - `openspec/specs/system-test-execution/spec.md` - 系统测试规范
  - `.sisyphus/drafts/test-coverage-improvement.md` - 本草案

  **输入/输出文件**:
  | 任务 | 输入 | 输出 |
  |------|------|------|
  | 3.1 | backend test results | coverage.out |
  | 3.2 | frontend test results | coverage/ |
  | 3.3 | 本草案 | `openspec/changes/improve-test-coverage-80-percent/` |
  | 3.4 | change proposal | 归档到 `openspec/changes/archive/` |

  **Acceptance Criteria**:
  - [ ] `go tool cover -func=coverage.out | grep total` 显示 ≥80%
  - [ ] `npm run test:cov` 显示覆盖率 ≥70%
  - [ ] `openspec validate improve-test-coverage-80-percent --strict` 通过
  - [ ] 所有测试通过，无 flaky tests

  **失败回滚策略**:
  - 覆盖率未达标: 创建 follow-up issue，降低目标
  - OpenSpec 验证失败: 修复 proposal 内容
  - 测试失败: 回退相关测试文件

  **风险点与优先级**:
  | 风险 | 级别 | 优先级 | 缓解措施 |
  |------|------|--------|----------|
  | 覆盖率计算差异 | 低 | P3 | 使用统一工具链 |
  | OpenSpec 格式错误 | 低 | P3 | 参考 AGENTS.md |

  **Commit**: YES (groups with 4)
  - Message: `docs(openspec): archive test coverage improvement change`
  - Files: `openspec/changes/improve-test-coverage-80-percent/`, `openspec/specs/`

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 1 | `chore(test): setup coverage infrastructure` | vitest.config.ts, testutil.go, coverage-gate.yml | CI workflow syntax check |
| 2 | `test(backend): add handler tests for core modules` | `backend/api/internal/handler/*_test.go` | `go test ./...` |
| 3 | `test(frontend): add view component tests` | `frontend-admin/tests/unit/*.test.ts` | `npm run test` |
| 4 | `docs(openspec): archive test coverage improvement change` | `openspec/` | `openspec validate --strict` |

---

## Success Criteria

### Verification Commands

```bash
# 后端覆盖率检查 (MUST show ≥80%)
cd /opt/code/DMH/backend
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | grep total
# Expected: total: (statements) 8X.X%

# 前端 Admin 覆盖率检查 (MUST show ≥70%)
cd /opt/code/DMH/frontend-admin
npm run test:cov 2>&1 | tail -5
# Expected: All files | 7X.X | ...

# 前端 H5 覆盖率检查 (MUST show ≥80%)
cd /opt/code/DMH/frontend-h5
npm run test:cov 2>&1 | tail -5
# Expected: All files | 8X.X | ...

# OpenSpec 验证
openspec validate improve-test-coverage-80-percent --strict
# Expected: PASS

# CI 门禁检查
grep -A5 "coverage" .github/workflows/coverage-gate.yml
# Expected: threshold configuration exists
```

### Final Checklist
- [ ] 后端覆盖率 ≥80%
- [ ] 前端 Admin 覆盖率 ≥70%
- [ ] 前端 H5 覆盖率 ≥80%
- [ ] CI 覆盖率门禁配置完成
- [ ] OpenSpec change 创建并验证通过
- [ ] 所有 "Must Have" 完成
- [ ] 所有 "Must NOT Have" 未发生
- [ ] 无 flaky tests
- [ ] 测试运行时间 < 5 分钟

---

## 附录

### A. Handler 测试模板

```go
// backend/api/internal/handler/testutil/testutil.go
package testutil

import (
    "database/sql"
    "net/http"
    "net/http/httptest"
    
    _ "github.com/mattn/go-sqlite3"
)

// SetupTestDB 创建内存数据库用于测试
func SetupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open test db: %v", err)
    }
    // 执行 schema 初始化...
    return db
}

// MakeRequest 创建测试 HTTP 请求
func MakeRequest(method, path string, body io.Reader) *http.Request {
    req := httptest.NewRequest(method, path, body)
    req.Header.Set("Content-Type", "application/json")
    return req
}

// ExecuteRequest 执行请求并返回响应
func ExecuteRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    handler(rr, req)
    return rr
}
```

### B. Vitest 配置修复

```typescript
// frontend-admin/vitest.config.ts
export default defineConfig({
  test: {
    include: [
      'tests/unit/**/*.{test,spec}.{ts,tsx,js,jsx}',
      'views/**/*.{test,spec}.{ts,tsx}'  // 新增: 收集 views 目录
    ],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      include: ['views/', 'services/', 'components/'],  // 明确包含
      exclude: [
        'node_modules/',
        'dist/',
        'e2e/',
        '**/*.config.ts',
        '**/types/**'
      ],
      thresholds: {
        lines: 70,
        functions: 70,
        branches: 70,
        statements: 70
      }
    }
  }
});
```

### C. CI 覆盖率门禁 Workflow

```yaml
# .github/workflows/coverage-gate.yml
name: Coverage Gate

on:
  pull_request:
    branches: [main]

jobs:
  backend-coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Run tests with coverage
        working-directory: backend
        run: |
          go test ./... -coverprofile=coverage.out -covermode=atomic
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Backend coverage: ${COVERAGE}%"
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "::error::Backend coverage ${COVERAGE}% is below 80% threshold"
            exit 1
          fi

  frontend-coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install dependencies
        working-directory: frontend-admin
        run: npm ci
      - name: Run tests with coverage
        working-directory: frontend-admin
        run: |
          npm run test:cov
          # 解析覆盖率并检查阈值...
```

### D. 风险点与优先级汇总

| # | 风险 | 级别 | 优先级 | 缓解措施 | 负责 |
|---|------|------|--------|----------|------|
| 1 | 时间超预期 | 高 | P0 | 分波执行，MVP 优先 | 全体 |
| 2 | Mock 复杂度高 | 中 | P1 | 使用 testutil 共享 fixture | 后端 |
| 3 | 组件测试困难 | 高 | P2 | E2E 替代单元测试 | 前端 |
| 4 | 外部依赖 Mock | 高 | P2 | 使用 interface mock | 后端 |
| 5 | CI 配置问题 | 中 | P3 | 本地验证后合并 | DevOps |
| 6 | 覆盖率计算差异 | 低 | P3 | 使用统一工具链 | 全体 |
