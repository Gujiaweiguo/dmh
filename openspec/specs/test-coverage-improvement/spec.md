# test-coverage-improvement Specification

## Purpose
TBD - created by archiving change improve-test-coverage-80-percent. Update Purpose after archive.
## Requirements
### Requirement: Coverage Infrastructure Setup
系统 SHALL 提供完整的测试覆盖率基础设施，包括配置文件、测试工具包和 CI 门禁。

#### Scenario: Vitest coverage configuration includes views
- **WHEN** 运行 frontend-admin 测试覆盖率检查
- **THEN** 系统 SHALL 收集 views 目录下的代码覆盖率
- **AND** 覆盖率阈值 SHALL 设置为 70%

#### Scenario: Handler test utility package available
- **WHEN** 开发者编写 handler 测试
- **THEN** 系统 SHALL 提供 testutil 包包含 SetupTestDB、MakeRequest、ExecuteRequest 函数
- **AND** 函数 SHALL 使用 sqlite in-memory 数据库模式

#### Scenario: CI coverage gate configured
- **WHEN** 创建 PR 到 main 分支
- **THEN** CI SHALL 运行覆盖率检查
- **AND** 后端覆盖率 < 80% SHALL 阻断合并
- **AND** 前端覆盖率 < 70% SHALL 阻断合并

### Requirement: Test Coverage Targets
系统 SHALL 定义并追踪各模块的测试覆盖率目标。

#### Scenario: Backend coverage target met
- **WHEN** 后端测试执行完成
- **THEN** 总体覆盖率 SHALL ≥ 80%
- **AND** 覆盖率报告 SHALL 保存为 coverage.out

#### Scenario: Frontend admin coverage target met
- **WHEN** 前端 admin 测试执行完成
- **THEN** 总体覆盖率 SHALL ≥ 70%
- **AND** 覆盖率报告 SHALL 保存为 coverage/coverage-final.json

### Requirement: Coverage Tracking and Reporting
系统 SHALL 提供覆盖率追踪和报告功能。

#### Scenario: Coverage report generated
- **WHEN** 测试执行完成
- **THEN** 系统 SHALL 生成覆盖率报告 (text, json, html, lcov)
- **AND** 报告 SHALL 包含各文件和总体覆盖率

#### Scenario: CI coverage threshold enforcement
- **WHEN** PR 提交
- **THEN** CI SHALL 检查覆盖率是否达到阈值
- **AND** 未达标 SHALL 返回非零退出码

