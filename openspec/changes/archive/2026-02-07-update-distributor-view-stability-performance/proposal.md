# Change: Stabilize distributor monitoring view and performance

## Why
- 历史归档中的分销监控任务存在未完成项，当前页面在稳定性、交互性能和发布前验证方面仍有缺口。
- 该页面属于平台管理员核心运营入口，需补齐超时处理、回归验证与性能基线，降低上线风险。

## What Changes
- 为分销监控页面补齐请求超时与重试提示，避免长时间无响应。
- 优化筛选与搜索交互（含输入防抖）并建立明确性能目标。
- 增加页面级单元测试、关键路径回归验证与生产构建验证流程。

## Impact
- Affected specs: rbac-permission-system
- Affected code:
  - frontend/views/DistributorManagementView.tsx
  - frontend/index.tsx
  - frontend/tests/** (新增或补充)
- Verification:
  - npm run build
  - 分销监控路由手动验证（筛选、刷新、分页）
  - 关键前端测试用例通过
