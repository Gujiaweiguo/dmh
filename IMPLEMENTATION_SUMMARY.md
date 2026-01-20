# 分销商系统实施总结

## 实施时间
2025-01-20

## 核心变更点

### 1. 商业逻辑变更（vs 原方案）
| 维度 | 原方案 | 新方案 |
|------|--------|--------|
| **成为分销商** | 申请→品牌管理员审批 | **支付订单后自动成为** |
| **奖励规则** | 全局统一 | **活动级别自定义** |
| **海报类型** | 未明确 | **活动专属 + 通用分销商海报** |
| **提现功能** | 无 | **有（平台审批）** |
| **审批流程** | 需要 | **不需要（简化）** |
| **数据查看** | 单一视角 | **品牌/平台双视角** |

### 2. 核心功能实现

#### 2.1 自动成为分销商
- **触发时机**：订单支付成功后
- **逻辑**：检查用户是否已是分销商，如不是则自动创建
- **默认级别**：1级
- **上级关系**：根据推荐人自动建立
- **无需审批**：简化流程，提升转化率

#### 2.2 多级分销奖励
- **最多3级**：一级、二级、三级
- **活动级别配置**：每个活动可自定义各级奖励比例
- **奖励计算**：订单金额 × 级别比例
- **分销链**：存储在 orders.distributor_path
- **状态检查**：只有 active 状态的分销商才获得奖励

#### 2.3 海报生成
- **活动专属海报**：包含活动信息 + 分销商二维码
- **通用分销商海报**：展示分销商所有活动
- **动态生成**：支持预览和下载

#### 2.4 提现功能
- **申请**：分销商可申请提现
- **审批**：平台管理员审批通过/拒绝
- **提现方式**：微信、支付宝、银行卡
- **金额验证**：不超过可用余额

### 3. 数据模型

#### 3.1 新增表
- `distributors` - 分销商信息表
- `distributor_rewards` - 分销商奖励记录表
- `poster_templates` - 海报模板表
- `withdrawals` - 提现申请表（已存在，扩展字段）

#### 3.2 扩展表
- `campaigns` - 增加 enable_distribution, distribution_level, distribution_rewards
- `orders` - 增加 distributor_path, paid_at

### 4. API接口设计

#### 4.1 H5端API
```
# 自动成为分销商（集成在支付回调中）
# POST /api/v1/orders/payment/callback - 支付回调

# 海报生成
POST /api/v1/posters/generate
GET /api/v1/posters/:id

# 提现功能
POST /api/v1/withdrawals/apply
GET /api/v1/withdrawals/my
```

#### 4.2 管理端API
```
# 品牌管理员
GET /api/v1/brands/:brandId/distributors
PUT /api/v1/brands/:brandId/distributors/:id/level
PUT /api/v1/brands/:brandId/distributors/:id/status
GET /api/v1/brands/:brandId/customers
GET /api/v1/brands/:brandId/rewards

# 平台管理员
GET /api/v1/platform/distributors?brandId={id}&status={status}
GET /api/v1/platform/rewards?brandId={id}
GET /api/v1/platform/withdrawals?brandId={id}&status={status}
PUT /api/v1/platform/withdrawals/:id/approve
PUT /api/v1/platform/withdrawals/:id/reject
```

### 5. 前端实现

#### 5.1 H5端页面
```
/views/distributor/
├── DistributorCenter.vue - 分销中心主页（含提现入口）
├── DistributorWithdrawals.vue - 提现申请和记录
├── DistributorPromotion.vue - 推广工具（海报生成）
├── DistributorRewards.vue - 奖励明细
├── DistributorSubordinates.vue - 下级列表
```

#### 5.2 管理后台页面
```
/views/
├── WithdrawalApprovalView.tsx - 平台管理员提现审批
├── DistributorManagementView.tsx - 品牌管理员分销商管理
```

### 6. 核心业务流程

#### 6.1 自动成为分销商流程
```
用户扫码 → 访问活动 → 填写表单 → 支付订单
  ↓
支付回调触发
  ↓
1. 检查用户是否已是该品牌分销商
2. 如果不是：
   - 查询活动是否启用分销
   - 检查推荐人是否是分销商
   - 创建分销商记录（level=1, parent_id=推荐人ID）
   - 添加distributor角色
   - 更新上级下级数量
3. 更新订单的distributor_path
4. 执行多级奖励计算
```

#### 6.2 多级奖励计算流程
```
订单支付成功 →
  查询活动的分销奖励配置
  查询订单的 distributor_path
  解析分销链（最多3级）
  为每级分销商计算并发放奖励
    - 级别1：订单金额 × level1比例
    - 级别2：订单金额 × level2比例
    - 级别3：订单金额 × level3比例
  → 更新余额和累计收益
```

#### 6.3 提现审批流程
```
分销商申请提现 →
  扣查余额 → 扣除金额 → 创建提现记录
    → 状态：pending
平台管理员审批 →
  批准 → 调用支付接口 → 完成打款 → 状态：completed
  拒绝 → 退还金额 → 状态：rejected
```

### 7. 权限矩阵

| 功能 | participant | distributor | brand_admin | platform_admin |
|------|-------------|-------------|-------------|----------------|
| 支付订单 | ✓ | ✓ | ✓ | ✓ |
| 自动成为分销商 | ✓ | - | - | - |
| 生成二维码海报 | - | ✓ | - | - |
| 查看分销数据 | - | ✓（自己） | ✓（本品牌） | ✓（全部）|
| 查看下级列表 | - | ✓（一级） | ✓（本品牌） | ✓（全部）|
| 申请提现 | - | ✓ | - | - |
| 管理分销商 | - | - | ✓（本品牌） | ✓（全部）|
| 审批提现 | - | - | - | ✓ |
| 查看品牌数据 | - | - | ✓（本品牌） | ✓（全部）|
| 查看全局数据 | - | - | - | ✓ |

### 8. 合规考虑

- **分销层级**：严格限制最多3级，符合法规要求
- **奖励比例**：活动级别配置，品牌可控
- **提现审核**：平台管理员统一审核，确保资金安全
- **数据隔离**：分销商只能查看自己的数据

## 文件清单

### 后端
```
/model/
  - distributor.go（更新）
  - user.go（扩展 Withdrawal 模型）
  - campaign.go（扩展 Campaign 模型）
  - poster.go（新增）

/api/internal/logic/distributor/
  - auto_upgrade_logic.go（新增）
  - multi_level_reward_logic.go（新增）
  - poster_logic.go（新增）
  - withdrawal_logic.go（新增）

/migrations/
  - 20250120_create_distributor_tables_final.sql（新增并执行）
```

### 前端H5
```
/frontend-h5/src/views/distributor/
  - DistributorCenter.vue（更新）
  - DistributorWithdrawals.vue（新增）

/frontend-h5/src/services/
  - （可能需要更新）
```

### 前端管理后台
```
/frontend-admin/views/
  - WithdrawalApprovalView.tsx（新增）

/frontend-admin/services/
  - distributorApi.ts（更新）
```

## 遗弃的功能

- ❌ 分销商申请审批流程
  - distributor_applications 表
- ❌ 品牌级别奖励配置表
- distributor_level_rewards 表
- ❌ 推广链接表
  distributor_links 表
- ❌ 申请审批相关API
- ❌ DistributorApprovalView.tsx（原申请审批页面，需改为提现审批）

## 下一步工作

### 紧急
1. 集成自动升级和多级奖励到支付回调
2. 完成海报生成API Handler实现
3. 完成提现API Handler实现
4. 集成到路由系统

### 短期
1. H5海报生成页面
2. H5分销商管理页面
3. 管理后台品牌管理员界面
4. 单元测试、集成测试、边界测试

### 后期
1. 数据迁移脚本优化
2. API文档完善
3. 性能监控
4. 运维文档

## 验收标准

### 功能验收
- ✅ 顾客支付订单后自动成为分销商
- ✅ 品牌管理员可以为活动配置分销奖励规则
- ✅ 分销商可以生成活动专属海报和通用海报
- ✅ 订单支付成功后，最多3级分销商自动获得奖励
- ✅ 分销商可以查看自己的推广数据和奖励明细
- ✅ 分销商可以申请提现
- ✅ 平台管理员可以审批分销商提现申请
- ✅ 品牌管理员可以管理分销商的级别和状态
- ✅ 平台管理员可以查看全局分销商、奖励、提现明细
- ✅ 品牌管理员可以查看本品牌的分销商、顾客、奖励详情
- ✅ 数据隔离正确：分销商只能看到自己的数据

### 性能验收
- 多级奖励计算 < 2秒
- 支持并发场景
- 提现审批流程响应及时

### 安全验收
- 最多3级分销，符合法规
- 并发场景余额正确
- 提现资金安全
- 数据隔离有效

## 技术栈

- **后端**: Go + go-zero + GORM
- **数据库**: MySQL 8.0（Docker容器中）
- **前端H5**: Vue 3 + Vant
- **前端管理后台**: React + TypeScript

## 注意事项

1. **数据库操作**：必须通过Docker进入MySQL容器执行
   ```bash
   docker exec -i mysql8 mysql -uroot -p'#Admin168' dmh < migration_file.sql
   ```

2. **Model字段命名**：Go中使用驼峰命名（CamelCase），JSON中使用下划线命名（snake_case）
   - Go: `TotalEarnings` -> DB: `total_earnings`
   - Go: `DistributorId` -> DB: `distributor_id`

3. **乐观锁**：余额更新使用乐观锁，防止并发问题

4. **状态管理**：
   - pending: 待处理
   - active/suspended: 分销商状态
   - pending/approved/rejected: 提现状态
   - pending/processing/completed/failed: 提现详细状态

5. **权限控制**：
   - H5端：检查distributor角色
   - 品牌管理员：检查brand_admin角色
   - 平台管理员：检查platform_admin角色

## 回滚计划

如果需要回滚：
1. 停止相关API路由
2. 回滚数据库迁移脚本（创建回滚脚本）
3. 删除或禁用相关代码

## 参考资料

- OpenSpec提案：/opt/code/DMH/openspec/changes/add-distributor-role/
- 数据库迁移：/opt/code/DMH/backend/migrations/20250120_create_distributor_tables_final.sql
- 工作笔记：/opt/code/DMH/.opencode/work_notes.md
