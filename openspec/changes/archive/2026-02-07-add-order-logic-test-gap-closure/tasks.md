## 1. Baseline and Mapping
- [x] 1.1 对照 `archive/2026-02-02-order-logic-implementation/.../tasks.md` 建立“未完成测试项 -> 代码位置”映射表
- [x] 1.2 标记每项测试的优先级（P0: 业务核心，P1: 边界扩展，P2: 可选）

## 2. Unit Test Backfill
- [x] 2.1 补齐 `createOrderLogic` 的异常分支测试（活动无效、重复订单、字段缺失/格式错误）
- [x] 2.2 补齐 `verifyOrderLogic`/`unverifyOrderLogic` 的权限与状态分支测试
- [x] 2.3 补齐表单字段验证的关键类型测试（phone/email/number/select）与不支持类型分支

## 3. Integration and Smoke
- [x] 3.1 新增订单关键路径集成测试：创建订单 -> 核销订单 -> 查询状态一致
- [x] 3.2 增加最小并发场景校验（重复创建或重复核销至少覆盖 1 条）
- [x] 3.3 新增最小 E2E 冒烟脚本并记录执行结果

## 4. Verification and Documentation
- [x] 4.1 执行统一测试命令并保存结果摘要（通过率、失败用例、耗时）
- [x] 4.2 更新变更内测试说明文档，记录覆盖范围与已知未覆盖项
- [x] 4.3 运行 `openspec validate add-order-logic-test-gap-closure --strict --no-interactive`
