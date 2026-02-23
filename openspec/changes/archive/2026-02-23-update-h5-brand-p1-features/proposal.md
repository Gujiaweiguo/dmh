# Change: H5 品牌端 P1 功能补齐

## Why

当前 `docs/H5_PENDING_FEATURES.md` 中的 4 个 P1 项（活动编辑 brandId 动态获取、品牌 Logo 上传、订单导出、订单详情）仍为待办，导致 H5 品牌端在核心运营链路上存在功能断点与临时实现。

## What Changes

- 明确并规划 H5 品牌端 4 个 P1 能力的交付范围与验收标准。
- 为活动编辑页补齐 brandId 从登录态/用户上下文动态获取能力，移除硬编码依赖。
- 为品牌设置页补齐 Logo 上传能力（上传、回填、保存）。
- 为订单页补齐订单导出能力与订单详情查看能力。
- 为上述能力定义前端逻辑测试与最小 E2E 覆盖要求，确保回归可验证。

## Impact

- Affected specs:
  - `h5-brand-operations` (new)
- Affected code:
  - `frontend-h5/src/views/brand/CampaignEditor.vue`
  - `frontend-h5/src/views/brand/Settings.vue`
  - `frontend-h5/src/views/brand/Orders.vue`
  - `frontend-h5/src/services/*`
  - `frontend-h5/tests/unit/*`
  - `frontend-h5/e2e/*`
