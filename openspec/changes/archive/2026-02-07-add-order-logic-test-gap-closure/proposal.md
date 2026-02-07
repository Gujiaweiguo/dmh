# Change: Close remaining order-logic test gaps

## Why
- `order-logic-implementation` 归档任务中仍保留多项未闭环测试条目（单元/集成/E2E），关键订单路径的回归保护不完整。
- 当前订单能力已用于线上迭代，若不补齐自动化测试，后续改动存在较高回归风险。

## What Changes
- 补齐订单创建、核销/取消核销、表单字段校验的缺口单元测试。
- 增加最小可执行的订单关键路径集成测试（创建 -> 核销 -> 状态校验）。
- 增加最小 E2E 冒烟脚本与测试结果记录模板，形成可追溯证据。
- 将测试命令纳入统一执行入口，确保本地与 CI 可重复执行。

## Impact
- Affected specs:
  - order-payment-system
- Affected code:
  - `backend/api/internal/logic/order/*_test.go`
  - `backend/api/internal/service/*_test.go`
  - `backend/api/tests/**`
  - `backend/api/scripts/**`
- Behavior changes: none（仅增强测试覆盖与回归保障）
