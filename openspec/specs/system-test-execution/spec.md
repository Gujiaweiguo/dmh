# system-test-execution Specification

## Purpose
TBD - created by archiving change add-system-test-execution-governance. Update Purpose after archive.
## Requirements
### Requirement: Test plan aligned execution scope
系统 SHALL 以 `test-plan.md` 作为系统测试执行基线，并确保测试范围覆盖 RBAC、活动管理、订单支付、反馈系统和规格治理。

#### Scenario: Execute planned module coverage
- **WHEN** 测试负责人发起一次完整系统测试执行
- **THEN** 执行范围 SHALL 覆盖 `test-plan.md` 中定义的全部 In Scope 模块
- **AND** 每个模块 SHALL 绑定至少一个可执行入口（命令、脚本或 CI 工作流）

### Requirement: Unified multi-layer test execution matrix
系统 SHALL 提供统一的分层执行矩阵，覆盖后端单元测试、后端集成测试、订单专项回归、前端单元测试和前端 E2E 测试。

#### Scenario: Run backend and frontend baseline suites
- **WHEN** 变更进入回归验证阶段
- **THEN** 系统 SHALL 执行后端 `go test` 基线与订单回归脚本
- **AND** 系统 SHALL 执行 `frontend-admin` 与 `frontend-h5` 的单元测试和 E2E 测试

### Requirement: Normal and abnormal scenario verification
系统 SHALL 在测试执行中同时覆盖正常与异常场景，异常场景至少包括鉴权失败、参数非法、权限不足、超时重试与重复请求幂等。

#### Scenario: Verify abnormal protections
- **WHEN** 测试执行到异常场景验证阶段
- **THEN** 系统 SHALL 对每类异常场景执行至少一个可重复验证的测试用例
- **AND** 失败结果 SHALL 包含可定位信息（错误码、日志或失败步骤）

### Requirement: Release QA gate for comprehensive testing
系统 SHALL 在发布前应用统一质量门禁：P0 用例通过率必须为 100%，P1 用例通过率必须不低于 95%，且不得存在 Blocker/Critical 未关闭缺陷。

#### Scenario: Enforce release gate
- **WHEN** 测试结果用于发布决策
- **THEN** 若 P0/P1 通过率或缺陷级别未满足门禁条件，系统 SHALL 阻止发布通过
- **AND** 若采用风险豁免，系统 SHALL 记录豁免理由、风险评估与审批信息

#### Scenario: Block release pipeline automatically when full regression is not green
- **WHEN** 全量回归中任一必跑套件失败或未执行
- **THEN** 系统 SHALL 自动阻断发布流水线与主干发布标记
- **AND** 阻断报告 SHALL 包含失败套件、失败用例、对应工作流链接与重跑记录

### Requirement: Traceable evidence and execution records
系统 SHALL 为每次系统测试执行产出可追溯记录，至少包含执行时间、执行范围、通过率、失败用例、关键日志/截图/报文引用。

#### Scenario: Archive execution evidence
- **WHEN** 一轮系统测试执行完成
- **THEN** 系统 SHALL 产出标准化测试结果记录
- **AND** 记录 SHALL 可用于回归审计与发布复盘

#### Scenario: Enforce evidence retention window
- **WHEN** 全量回归证据归档后
- **THEN** 系统 SHALL 保留回归证据至少 90 天
- **AND** 失败回归证据 SHALL 至少保留 180 天并支持按提交号与工作流检索

### Requirement: Full regression orchestration workflow
系统 SHALL 提供全量回归编排能力，统一触发条件、执行顺序、失败重试与结果聚合口径。

#### Scenario: Trigger full regression on release candidate and protected branch changes
- **WHEN** 发生 release candidate 发布候选构建或主干受保护分支变更
- **THEN** 系统 SHALL 自动触发全量回归编排
- **AND** 编排 SHALL 覆盖后端单元、后端集成、订单专项回归、前端双端单元与 E2E、OpenSpec 校验

#### Scenario: Aggregate multi-workflow outcomes into single regression verdict
- **WHEN** 全量回归相关工作流执行完成
- **THEN** 系统 SHALL 生成单一回归结论（PASS/FAIL）
- **AND** 结论 SHALL 可映射到发布门禁判定

### Requirement: Flaky test quarantine and retry policy
系统 SHALL 对 flaky 测试实施受控重试与隔离策略，避免以不稳定结果放行发布。

#### Scenario: Retry transient failures with bounded attempts
- **WHEN** 测试失败被判定为环境瞬态问题
- **THEN** 系统 SHALL 允许自动重试且重试次数不得超过 2 次
- **AND** 每次重试 SHALL 记录独立日志与原因标签

#### Scenario: Quarantine persistent flaky tests with owner and expiry
- **WHEN** 同一测试在连续 3 次回归中出现至少 2 次非确定性失败
- **THEN** 系统 SHALL 将该测试标记为 flaky 并进入隔离清单
- **AND** 隔离项 SHALL 包含责任人、修复截止时间与退出标准

### Requirement: Plan v1 as single execution authority
系统 SHALL 将 `openspec/changes/add-full-regression-testing/tasks.md` 中 Plan v1 作为全量回归改造的唯一执行依据。

#### Scenario: Execute only against Plan v1
- **WHEN** Atlas/Hephaestus 开始实现本变更
- **THEN** 执行任务 SHALL 与 Plan v1 中的任务 ID、依赖顺序、回归批次保持一致
- **AND** 未记录在 Plan v1 的任务 SHALL 不得进入执行

#### Scenario: Revise authority before scope/order change
- **WHEN** 需要变更任务范围、顺序或风险分级
- **THEN** 团队 SHALL 先更新 Plan v1 并完成评审
- **AND** 更新前 SHALL 暂停对应实现任务

