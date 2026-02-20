# Test Coverage Specification

## Overview

测试覆盖率规格定义了 DMH 项目的测试覆盖标准和 CI 门禁要求。

---

## ADDED Requirements

### Requirement: Backend Test Coverage Threshold

后端代码测试覆盖率 SHALL 达到 70% 最低阈值。

#### Scenario: Backend coverage meets threshold
- **WHEN** 运行 `go test ./... -cover`
- **THEN** 总覆盖率 ≥ 70%
- **AND** 输出显示覆盖率百分比

#### Scenario: CI blocks low coverage PR
- **WHEN** PR 提交到 main 分支
- **AND** 后端覆盖率 < 70%
- **THEN** CI 检查显示警告

---

### Requirement: Frontend Admin Test Coverage Threshold

前端 Admin 测试覆盖率 SHALL 达到 70% 最低阈值。

#### Scenario: Frontend Admin coverage meets threshold
- **WHEN** 运行 `npm run test:cov`
- **THEN** 总覆盖率 ≥ 70%
- **AND** 输出显示 All files 覆盖率

---

### Requirement: CI Coverage Gate

项目 SHALL 配置 CI 覆盖率门禁 workflow。

#### Scenario: Coverage gate workflow exists
- **WHEN** 查看 `.github/workflows/coverage-gate.yml`
- **THEN** 文件存在
- **AND** 包含后端覆盖率检查 job
- **AND** 包含前端覆盖率检查 job

---

### Requirement: Test Database Isolation

测试 SHALL 使用独立的数据库实例，不影响生产数据。

#### Scenario: Test database created per test
- **WHEN** 运行需要数据库的测试
- **THEN** 创建独立的测试数据库
- **AND** 数据库名唯一且短暂

#### Scenario: Test database cleaned up after test
- **WHEN** 测试完成
- **THEN** 测试数据库被删除
- **AND** 不影响其他测试
