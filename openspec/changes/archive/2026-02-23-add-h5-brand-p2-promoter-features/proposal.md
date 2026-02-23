# Change: H5 品牌端 P2 推广员功能增强

## Why

当前 `docs/H5_PENDING_FEATURES.md` 中推广员相关的 P2 功能（推广员详情、奖励记录）仍为待办，导致品牌管理员无法深入了解推广员画像和奖励发放历史，影响运营效率。

## What Changes

- 为推广员管理页增加详情查看能力，展示推广员基础信息、业绩统计、推广链接等。
- 为推广员管理页增加奖励记录查看能力，展示奖励发放历史、状态、金额等。
- 定义前端逻辑测试与最小 E2E 覆盖要求，确保可验证。

## Impact

- Affected specs:
  - `h5-brand-operations` (extend)
- Affected code:
  - `frontend-h5/src/views/brand/Promoters.vue`
  - `frontend-h5/src/views/brand/PromoterDetail.vue` (new)
  - `frontend-h5/src/views/brand/RewardRecords.vue` (new)
  - `frontend-h5/src/views/brand/promoters.logic.js`
  - `frontend-h5/src/router/index.js`
  - `frontend-h5/tests/unit/*`
  - `frontend-h5/e2e/*`
