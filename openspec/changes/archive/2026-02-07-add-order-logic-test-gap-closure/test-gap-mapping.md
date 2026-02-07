# Order Logic Test Gap Mapping

## Source
- Archive task file: `openspec/changes/archive/2026-02-02-order-logic-implementation/order-logic-implementation/tasks.md`
- Focus sections: 4 (create), 5 (verify/unverify), 6 (field validation), 7 (integration), 9 (E2E)

## Mapping and Priority

| Archive Item | Current Coverage | Target File(s) | Priority |
|---|---|---|---|
| 4.1-4.9 订单创建测试 | 已有核心覆盖（success/活动不存在/活动结束/重复/必填/字段错误/手机号）但缺少“字段类型异常分支”显式用例 | `backend/api/internal/logic/order/create_order_logic_test.go` | P0 |
| 5.1-5.10 订单核销/取消核销测试 | 已有核心覆盖（invalid code/order not found/success/already verified/unverify success）但缺少“权限上下文分支”和部分异常补齐 | `backend/api/internal/logic/order/verify_order_logic_test.go` | P0 |
| 6.1-6.8 字段验证测试 | 已覆盖 phone/email/number/select；需补齐不支持字段类型与边界说明 | `backend/api/internal/logic/order/create_order_logic_test.go` | P0 |
| 7.1-7.8 集成测试 | 现状缺失独立集成测试入口 | `backend/api/internal/logic/order/order_integration_test.go` | P1 |
| 9.1-9.5 端到端验证 | 现状缺失最小可执行脚本与结果记录模板 | `backend/scripts/run_order_logic_tests.sh`, `backend/scripts/order_smoke.sh`, `openspec/changes/add-order-logic-test-gap-closure/test-results.md` | P1 |

## Execution Scope in This Change
- P0: 先补齐关键单元测试缺口，确保创建/核销关键分支具备稳定回归保护。
- P1: 补齐最小集成与冒烟执行入口，形成“可执行 + 可记录”的闭环。
- P2: 暂不引入更重的全链路外部依赖测试（如真实支付回调/外部数据库同步）。
