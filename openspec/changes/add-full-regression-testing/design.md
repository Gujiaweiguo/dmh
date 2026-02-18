## Context

DMH 当前已有多层测试资产：
- 后端：`go test ./...`、`go test ./test/integration/...`、订单专项回归脚本。
- 前端：`frontend-admin` 与 `frontend-h5` 的单元测试与 Playwright E2E。
- CI：`stability-checks.yml`、`system-test-gate.yml`、`coverage-gate.yml`、`order-mysql8-regression.yml`。

但“全量回归测试”在规格层仍缺少以下统一约束：
1) 单一编排定义（本地/CI 一致口径）；
2) flaky 重试规则；
3) 证据保留周期与审计字段；
4) 发布阻断与门禁关系。

## Goals / Non-Goals

### Goals
- 定义全量回归的最小必跑矩阵（后端单元+集成+专项，前端单元+E2E，覆盖门禁）。
- 定义失败判定和阻断策略，确保发布决策有统一标准。
- 定义 flaky 控制策略，避免“随机绿灯”。
- 定义可审计证据留存要求，支持回归追溯。

### Non-Goals
- 不引入新测试框架或替换现有工具链。
- 不要求本提案阶段立即实现所有新脚本，仅先固化规格约束。

## Decisions

### Decision 1: 复用并强化 `system-test-execution`，不新建重复 capability
- 原因：已有 capability 已覆盖系统测试矩阵与门禁，新增独立 capability 会产生语义重叠。
- 结果：本次采用 `MODIFIED + ADDED` 方式在现有 spec 上扩展“全量回归”定义。

### Decision 2: 将“全量回归”定义为发布前阻断条件
- 原因：目前部分 workflow 按路径触发，存在“关键检查未被触发即放行”的风险。
- 结果：规格要求发布判定必须基于“全量回归结果集合”，而非单一工作流成功。

### Decision 3: 将 flaky 重试从“默认允许”改为“显式白名单”
- 原因：无约束重试会掩盖真实缺陷。
- 结果：仅允许标注为 flaky 的测试进行有上限重试，超限仍判失败。

### Decision 4: 证据保留采用分级策略
- 原因：兼顾存储成本与审计需要。
- 结果：至少保留最近回归报告与失败证据，且要求包含可定位字段（commit、workflow、用例、日志/截图引用）。

## Risks / Trade-offs

- 风险：将“全量回归”纳入发布阻断后，短期可能增加合并等待时间。
  - 缓解：并行化执行矩阵，优先保留阻断价值高的检查。

- 风险：flaky 白名单治理增加维护成本。
  - 缓解：要求白名单附带过期时间与责任人，定期清理。

## Migration Plan

1. 先更新 OpenSpec requirement（本提案）。
2. 在实现阶段对齐 workflow、脚本与文档。
3. 验证“本地全量回归入口”和“CI 全量回归门禁”口径一致。

## Open Questions

- 证据最短保留天数是否按 30 天统一，还是区分“成功回归”和“失败回归”两档保留。
- 是否在实现阶段引入 nightly 全量回归作为强制门禁前置信号。
