## 1. Implementation

- [x] 1.1 新建推广员详情页组件 `PromoterDetail.vue`，展示基础信息、业绩统计、推广链接。
- [x] 1.2 在 `router/index.js` 添加 `/brand/promoter-detail/:id` 路由。
- [x] 1.3 在 `Promoters.vue` 中为每行推广员添加"详情"入口，跳转详情页。
- [x] 1.4 新建奖励记录页组件 `RewardRecords.vue`，展示奖励列表、状态、金额、时间。
- [x] 1.5 在 `router/index.js` 添加 `/brand/reward-records/:promoterId?` 路由。
- [x] 1.6 在推广员详情页增加"查看奖励记录"入口。

## 2. Verification

- [x] 2.1 新增/更新 H5 单元测试，覆盖详情与奖励记录的核心逻辑分支。
- [x] 2.2 新增/更新 H5 E2E 用例，覆盖推广员详情查看与奖励记录跳转。
- [x] 2.3 执行 `cd frontend-h5 && npm run test` 并通过。
- [x] 2.4 执行 `cd frontend-h5 && npm run test:e2e:headless` 并通过。

## 3. Documentation

- [x] 3.1 更新 `docs/H5_PENDING_FEATURES.md`，将本变更覆盖的 P2 项标记为由本 change 跟踪。
- [ ] 3.2 在实现完成后补充回归证据到 `docs/testing/execution/runs/`。
