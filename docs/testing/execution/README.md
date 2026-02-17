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

## 使用说明

1. 每轮系统测试开始前先更新 `PREREQUISITES.md` 与 `EXECUTION_ORDER.md`。
2. 执行完成后填写 `TEST_RESULT_TEMPLATE.md`，并在各矩阵文件中补齐状态。
3. 如存在带风险发布，必须填写 `RISK_WAIVER_TEMPLATE.md` 并附审批信息。
