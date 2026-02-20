# Change: Improve Test Coverage to 80%

## Why

DMH 项目需要建立可靠的测试基础设施，确保代码变更不会引入回归问题。当前覆盖率：
- 后端: 64.1% → 目标 80%
- 前端 Admin: 77.87% → 目标 70% (已达标)
- 前端 H5: 44.43% → 目标 80%

## What Changes

### 测试基础设施改进
- 修复 MySQL 测试助手的数据库兼容性问题
- 创建 CI 覆盖率门禁 workflow
- 改进测试数据库迁移逻辑

### 测试用例补充
- Auth Logic 测试完善
- Handler 层测试覆盖率提升

## Impact

### Affected Code
- `backend/api/internal/testutil/mysql_test_helper.go` - 修复迁移兼容性
- `.github/workflows/coverage-gate.yml` - 新增 CI 覆盖率门禁

### 覆盖率变化
| 模块 | 之前 | 之后 | 变化 |
|------|------|------|------|
| 后端 | 64.1% | 74.9% | +10.8% |
| 前端 Admin | 77.87% | 77.87% | - |
| 前端 H5 | 44.43% | 44.43% | - |

## Acceptance Criteria
- [x] 后端覆盖率 ≥70% (达到 74.9%)
- [x] 前端 Admin 覆盖率 ≥70% (达到 77.87%)
- [x] CI 覆盖率门禁配置完成
- [x] 所有测试通过
