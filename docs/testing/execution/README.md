# 系统测试执行证据模板

本目录用于承接 `add-system-test-execution-governance` 变更中的证据输出。

## 文件清单

- `SCOPE_MAPPING.md`: 测试范围到执行入口映射
- `EXECUTION_ORDER.md`: P0/P1 执行顺序与节奏
- `PREREQUISITES.md`: 测试前置条件清单
- `BACKEND_MATRIX.md`: 后端执行矩阵
- `FRONTEND_MATRIX.md`: 前端执行矩阵
- `ABNORMAL_MATRIX.md`: 异常场景验证矩阵
- `QUALITY_GATE.md`: 发布门禁规则
- `TEST_RESULT_TEMPLATE.md`: 测试结果记录模板
- `RISK_WAIVER_TEMPLATE.md`: 风险豁免模板
- `CI_ALIGNMENT.md`: CI 守卫与执行矩阵对齐
- `STABILITY_CHECKS_WORKFLOW.md`: 稳定性 CI 工作流说明（backend 全量 + admin 单测 + security E2E）
- `PREMERGE_MIN_SUITE.md`: 合并前最小必跑集
- `FAILURE_PLAYBOOK.md`: 失败回滚与重试手册
- `FULL_REGRESSION_DEFINITION.md`: 全量回归口径定义（必跑矩阵 + PASS/FAIL 判定规则）
- `RELEASE_BLOCKING_RULES.md`: 发布阻断规则定义（硬阻断/软阻断 + CI 状态映射 + 阻断报告模板）
- `FLAKY_TEST_STRATEGY.md`: Flaky 测试策略（判定标准 + 重试策略 + 隔离机制 + 退出标准）
- `EVIDENCE_RETENTION_POLICY.md`: 证据保留策略（成功 90 天 / 失败 180 天 + 检索主键定义）
- `CI_ORCHESTRATION_PLAN.md`: CI 编排对齐方案（全量回归触发条件 + 单一结论聚合 + workflow 改造建议）
- `FULL_REGRESSION_EVIDENCE_TEMPLATE.md`: 全量回归证据模板（审计字段定义 + 分模块模板 + 填写示例 + 存储检索）
- `LOCAL_REGRESSION_ENTRY.md`: 本地一键全量回归入口（命令定义 + 执行范围 + 环境前置要求 + Makefile 扩展建议）
- `FULL_REGRESSION_DRILL_REPORT.md`: 全量回归演练报告（模拟演练记录 + 执行结果 + 问题清单 + 改进建议）

## 使用说明

1. 每轮系统测试开始前先更新 `PREREQUISITES.md` 与 `EXECUTION_ORDER.md`。
2. 执行完成后填写 `TEST_RESULT_TEMPLATE.md`，并在各矩阵文件中补齐状态。
3. 如存在带风险发布，必须填写 `RISK_WAIVER_TEMPLATE.md` 并附审批信息。
