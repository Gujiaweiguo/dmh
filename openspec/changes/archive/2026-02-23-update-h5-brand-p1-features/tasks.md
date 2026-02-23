## 1. Implementation

- [x] 1.1 活动编辑页改造：移除 `brandId` 硬编码，改为从当前用户品牌上下文动态注入。
- [x] 1.2 品牌设置页补齐 Logo 上传：接入上传接口、回填 Logo URL、保存品牌配置。
- [x] 1.3 订单页补齐订单导出：提供筛选条件下的导出入口与文件下载反馈。
- [x] 1.4 订单页补齐订单详情：支持从列表进入详情并展示关键字段（支付、状态、核销等）。

## 2. Verification

- [x] 2.1 新增/更新 H5 单元测试，覆盖上述 4 个 P1 功能的核心逻辑分支。
- [x] 2.2 新增/更新 H5 E2E 用例，覆盖订单详情查看与导出主路径。
- [x] 2.3 执行 `cd frontend-h5 && npm run test` 并通过。
- [x] 2.4 执行 `cd frontend-h5 && npm run test:e2e:headless` 并通过。

## 3. Documentation

- [x] 3.1 更新 `docs/H5_PENDING_FEATURES.md`，将本变更覆盖的 4 个 P1 项标记为由本 change 跟踪。
- [x] 3.2 在实现完成后补充回归证据到 `docs/testing/execution/runs/`。
