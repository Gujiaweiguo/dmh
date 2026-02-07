# Change: Backfill automated tests for order and RBAC critical paths

## Why
- 归档任务中仍有大量订单与权限相关测试项未闭环，导致关键流程缺少持续回归保护。
- 当前核心功能已上线使用，若无自动化测试覆盖，后续迭代容易引入隐性回归。

## What Changes
- 补齐订单创建、订单核销、表单字段校验的单元测试与关键集成测试。
- 补齐 RBAC 权限检查与数据隔离相关单元测试。
- 建立最小端到端回归脚本与测试结果记录模板。
- 统一历史任务状态说明，明确已完成与待补齐边界（仅文档状态，不改历史归档内容）。

## Impact
- Affected specs:
  - order-payment-system
  - rbac-permission-system
- Affected code:
  - backend/api/internal/logic/order/*_test.go
  - backend/api/internal/service/*_test.go
  - backend/api/internal/middleware/*_test.go
  - backend/api/tests/**
- Behavior changes: none（仅增强质量保障与回归能力）
