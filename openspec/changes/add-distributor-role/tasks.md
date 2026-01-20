# Tasks: 分销商系统实现

## 1. 数据库设计与迁移
- [x] 1.1 创建 `distributors` 表【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:9)
  - id, user_id, brand_id, level, parent_id, status
  - created_at, total_earnings, subordinates_count
  - 索引：user_id, brand_id, parent_id, status
- [x] 1.2 创建 `withdrawals` 表【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:187)
  - id, user_id, brand_id, distributor_id, amount, status
  - pay_type, pay_account, pay_real_name
  - approved_by, approved_at, approved_notes, rejected_reason
  - paid_at, trade_no
  - 索引：user_id, brand_id, distributor_id, status
- [x] 1.3 创建 `poster_templates` 表【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:50)
  - id, type, campaign_id, template_url, created_at
- [x] 1.4 扩展 `campaigns` 表【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:70)
  - enable_distribution (BOOLEAN)
  - distribution_level (INT)
  - distribution_rewards (JSON)
- [x] 1.5 扩展 `orders` 表【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:126)
  - distributor_path (VARCHAR)
- [x] 1.6 修改 `roles` 表，插入 distributor 角色记录【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:162)
- [x] 1.7 创建数据库迁移脚本【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:1)

## 2. 后端API实现
- [x] 2.1 实现自动成为分销商逻辑【已落地】(backend/api/internal/logic/distributor/auto_upgrade_logic.go:30; backend/api/internal/logic/order/payment_callback_logic.go:113)
  - 支付回调中检查用户是否已是分销商
  - 未成为分销商则自动创建分销商记录
- [x] 2.2 实现活动级别分销奖励规则配置【已落地】(backend/model/campaign.go:18; backend/api/internal/logic/campaign/create_campaign_logic.go:72; backend/api/internal/logic/campaign/get_campaign_logic.go:53)
  - 活动创建/编辑时配置分销规则
  - 存储到 campaigns 表的 JSON 字段
- [x] 2.3 实现多级奖励计算逻辑【已落地】(backend/api/internal/logic/distributor/multi_level_reward_logic.go:40; backend/api/internal/logic/distributor/reward_logic.go:29)
  - 查询活动的分销奖励规则
  - 查询分销链（distributor_path）
  - 按级别比例计算奖励
  - 创建多条奖励记录（每级一条）
  - 幂等性检查和事务安全
- [x] 2.4 实现海报生成API【已落地】(backend/api/internal/logic/distributor/poster_logic.go:36; backend/api/internal/handler/distributor/poster_handlers.go:12)
  - POST /api/v1/posters/generate - 生成活动专属海报
  - POST /api/v1/posters/generate - 生成通用分销商海报
  - GET /api/v1/posters/:id - 获取海报
- [x] 2.5 实现分销商数据查询API【已落地】(backend/api/internal/handler/distributor/distributor_handlers.go:63; backend/api/internal/logic/distributor/statistics_logic.go:29)
  - GET /api/v1/distributor/statistics - 推广数据统计
  - GET /api/v1/distributor/rewards - 奖励明细
  - GET /api/v1/distributor/subordinates - 下级列表
- [x] 2.6 实现提现API【已落地】(backend/api/internal/handler/distributor/withdrawal_handlers.go:13; backend/api/internal/logic/distributor/withdrawal_logic.go:31)
  - POST /api/v1/withdrawals/apply - 申请提现
  - GET /api/v1/withdrawals - 提现列表
- [x] 2.7 实现平台管理员提现审批API【已落地】(backend/api/internal/handler/distributor/withdrawal_handlers.go:53)
  - PUT /api/v1/platform/withdrawals/:id/approve - 审批通过
  - PUT /api/v1/platform/withdrawals/:id/reject - 审批拒绝
- [x] 2.8 实现品牌管理员分销商管理API【已落地】(backend/api/internal/handler/distributor/brand_distributor_handlers.go:88; backend/api/internal/logic/distributor/management_logic.go:28)
  - GET /api/v1/brands/:brandId/distributors - 分销商列表
  - GET /api/v1/distributors/:id - 分销商详情
  - PUT /api/v1/brands/:brandId/distributors/:id/level - 调整级别
  - PUT /api/v1/brands/:brandId/distributors/:id/status - 暂停/激活
- [x] 2.9 实现品牌管理员数据查看API【已落地】(backend/api/internal/logic/distributor/management_logic.go:431; backend/api/internal/handler/distributor/brand_distributor_handlers.go:226)
  - GET /api/v1/brands/:brandId/customers - 顾客列表
  - GET /api/v1/brands/:brandId/rewards - 奖励详情
- [x] 2.10 实现平台管理员全局数据查询API【已落地】(backend/api/internal/logic/distributor/global_stats_logic.go:108; backend/api/internal/handler/distributor/platform_handlers.go:33)
  - GET /api/v1/platform/distributors - 全部分销商（按品牌筛选）
  - GET /api/v1/platform/rewards - 全部奖励（按品牌筛选）
  - GET /api/v1/platform/withdrawals - 全部提现（按品牌筛选）

## 3. 前端H5实现（复用现有H5）
- [x] 3.1 创建分销中心页面 (DistributorCenterView.vue)【已落地】(frontend-h5/src/views/distributor/DistributorCenter.vue:1)
  - 入口：个人中心 → 分销中心
  - 展示：当前级别、累计收益、可提现金额、下级数量
- [x] 3.2 实现海报生成页面【已落地】(frontend-h5/src/views/distributor/DistributorPromotion.vue:1)
  - 活动专属海报：选择活动 → 生成海报 → 预览和下载
  - 通用分销商海报：展示所有活动 → 生成海报 → 预览和下载
  - 推广链接和二维码功能
- [x] 3.3 实现推广数据统计页面【已落地】(frontend-h5/src/views/distributor/DistributorCenter.vue:56; backend/api/internal/logic/distributor/statistics_logic.go:29)
  - 数据卡片：订单数、收益、下级数
  - 奖励明细列表（分页）
  - 收益趋势展示
- [x] 3.4 实现下级列表页面【已落地】(frontend-h5/src/views/distributor/DistributorSubordinates.vue:1)
  - 一级下级分销商列表
  - 每个下级的基本信息
- [x] 3.5 实现提现申请页面【已落地】(frontend-h5/src/views/distributor/DistributorWithdrawals.vue:57)
  - 输入提现金额
  - 选择提现方式（微信/支付宝/银行卡）
  - 输入账号和真实姓名
  - 提交申请
- [x] 3.6 实现提现记录页面【已落地】(frontend-h5/src/views/distributor/DistributorWithdrawals.vue:14)
  - 提现记录列表（分页）
  - 显示每笔提现的状态、金额、时间

## 4. 前端管理后台实现（品牌管理员）
- [x] 4.1 创建分销商管理页面 (DistributorManagementView.tsx)【已落地】(frontend-admin/views/DistributorManagementView.tsx:1)
  - 分销商列表（按品牌筛选）
  - 分销商详情查看
  - 级别调整
  - 状态管理（暂停/激活）
- [x] 4.2 创建顾客列表页面 (CustomerListView.tsx)【已落地】(frontend-admin/views/CustomerListView.tsx:1)
  - 顾客列表（按品牌筛选）
  - 顾客基本信息
- [x] 4.3 创建奖励详情页面 (RewardDetailView.tsx)【已落地】(frontend-admin/views/RewardDetailView.tsx:1)
  - 奖励列表（按品牌筛选）
  - 奖励明细
  - 统计数据
- [x] 4.4 在品牌管理菜单中添加"分销管理"入口【已落地】
- [x] 4.5 更新权限控制组件，支持 distributor 角色【已落地】(frontend-admin/components/PermissionGuard.tsx:32)

## 5. 前端管理后台实现（平台管理员）
- [x] 5.1 创建全局分销商查询页面 (PlatformDistributorView.tsx)【已落地】(frontend-admin/views/PlatformDistributorView.tsx:1)
  - 全部分销商列表
  - 按品牌筛选
  - 分销商详情
  - 级别调整和状态管理
- [x] 5.2 创建全局奖励查询页面 (PlatformRewardView.tsx)【已落地】(frontend-admin/views/PlatformRewardView.tsx:1)
  - 全部奖励列表
  - 按品牌筛选
  - 奖励明细
- [x] 5.3 创建提现审批页面 (WithdrawalApprovalView.tsx)【已落地】(frontend-admin/views/WithdrawalApprovalView.tsx:1)
  - 待审批提现列表
  - 审批操作（通过/拒绝）
  - 查看提现详情
  - 输入审批备注或拒绝原因
- [x] 5.4 在平台管理菜单中添加"分销管理"入口【已落地】

## 6. 活动管理系统改造
- [x] 6.1 在活动创建/编辑页面添加分销规则配置【已落地】
  - 启用/禁用分销
  - 设置分销层级（1-3级）
  - 设置各级奖励比例
- [x] 6.2 在活动详情页面展示分销规则【已落地】
- [x] 6.3 保存分销规则到 campaigns 表【已落地】(backend/model/campaign.go:18)

## 7. 奖励系统改造
- [x] 7.1 实现多级奖励计算逻辑【已落地】(backend/api/internal/logic/distributor/multi_level_reward_logic.go:40)
  - 订单支付成功时查询 distributor_path
  - 向上追溯最多3级
  - 按级别比例计算奖励（从活动规则读取）
  - 只有 active 状态的分销商才获得奖励
  - 幂等性检查和事务安全
- [x] 7.2 更新奖励记录，关联分销商级别信息【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:31; backend/api/internal/logic/distributor/multi_level_reward_logic.go:186)
- [x] 7.3 确保奖励计算的幂等性和事务安全【已落地】(backend/api/internal/logic/distributor/multi_level_reward_logic.go:176; backend/api/internal/logic/distributor/reward_logic.go:16)
- [x] 7.4 支持活动级别奖励比例配置【已落地】(backend/api/internal/logic/distributor/multi_level_reward_logic.go:61)

 ## 8. 测试
 - [x] 8.1 单元测试【已落地】(backend/api/internal/logic/distributor/auto_upgrade_logic_test.go:1; backend/api/internal/logic/distributor/multi_level_reward_logic_test.go:1; backend/api/internal/logic/distributor/withdrawal_logic_test.go:1)
   - 自动成为分销商逻辑测试
   - 多级奖励计算测试
   - 提现审批流程测试
   - 权限控制测试
 - [x] 8.2 集成测试【已落地】(backend/test/integration/distributor_integration_test.go:1)
   - 完整的支付→自动升级→推广→奖励→提现流程
   - 多级分销奖励分配验证
   - 品牌管理员管理分销商流程
   - 平台管理员审批提现流程
 - [x] 8.3 边界测试【已落地】(backend/test/integration/distributor_integration_test.go:405; backend/test/integration/distributor_integration_test.go:427)
   - 超过3级不分配奖励
   - 非分销商不获得奖励
   - 暂停状态不获得奖励
   - 提现金额校验
   - 活动未启用分销时不计算奖励

## 9. 文档和部署
- [ ] 9.1 更新API文档【未落地】
- [ ] 9.2 编写分销商使用指南【未落地】
- [ ] 9.3 编写品牌管理员操作指南【未落地】
- [ ] 9.4 编写平台管理员操作指南【未落地】
- [x] 9.5 准备数据库迁移脚本【已落地】(backend/migrations/20250120_create_distributor_tables_final.sql:1)
- [ ] 9.6 部署验证【未落地】

## 验收标准
- [ ] 顾客支付订单后自动成为分销商
- [ ] 品牌管理员可以在活动创建时设置分销奖励规则
- [ ] 分销商可以生成活动专属海报和通用分销商海报
- [ ] 订单支付成功后，最多3级分销商自动获得奖励（按活动规则）
- [ ] 分销商可以查看自己的推广数据和奖励明细
- [ ] 分销商可以申请提现
- [ ] 平台管理员可以审批提现申请
- [ ] 品牌管理员可以管理分销商的级别和状态
- [ ] 品牌管理员可以查看本品牌的分销商、顾客、奖励详情
- [ ] 平台管理员可以查看全局分销商、奖励、提现明细（可按品牌筛选）
- [ ] 数据隔离正确：分销商只能看到自己的数据
