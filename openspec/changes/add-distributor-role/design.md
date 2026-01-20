# Design: 分销商系统

## Context

DMH平台需要支持分销商（大C端）功能，这是一种高级顾客角色，具备分享推广并获得奖励的资格。分销商在支付订单后自动生效。

系统当前已有：
- 普通顾客（participant）可参与活动、支付、获得优惠券
- 推荐人机制（referrer_id）已在订单中实现，用于单级推荐奖励

需要在此基础上扩展：
- 多级分销体系（1-3级）
- 自动成为分销商机制
- 活动级别分销奖励规则配置
- 分销商专属功能（二维码海报、提现、数据查看）

## Goals / Non-Goals

**Goals**:
- 实现分销商角色定义与权限管理
- 实现支付后自动成为分销商机制
- 实现多级分销奖励计算（最多3级，活动级别配置）
- 提供分销商推广工具（二维码海报：活动专属 + 通用）
- 提供分销数据查看功能
- 实现提现功能和平台审批流程

**Non-Goals**:
- 分销商独立H5/小程序（复用现有H5）
- 复杂的团队管理系统
- 自动升级机制（级别需品牌管理员手动调整）
- 分销商PC管理后台
- 品牌管理员审批提现（由平台管理员审批）

## Decisions

### 1. 分销商角色设计

**Decision**: 分销商作为 participant 的扩展角色，而非完全独立的新角色

**Rationale**:
- 分销商本质上仍是顾客，只是具有额外的推广权限
- 复用现有的用户表和会员系统
- 降低系统复杂度

**Implementation**:
```
users 表保持不变
user_roles 表新增 distributor 角色关联
新增 distributors 表存储分销商专属信息（级别、上级等）
```

### 2. 自动成为分销商机制

**Decision**: 用户支付订单成功后自动成为分销商，无需审批

**Rationale**:
- 简化流程，降低用户门槛
- 提高分销商转化率
- 已支付用户具有一定质量保障

**Implementation**:
```
支付回调成功 →
  检查用户是否已是分销商 →
    否：创建分销商记录（level=1, parent_id=referrer_id）
    是：跳过
  → 原有奖励计算
```

### 3. 活动级别分销奖励规则

**Decision**: 每个活动可单独设置分销奖励规则

**Rationale**:
- 不同活动可能需要不同的激励策略
- 品牌可灵活调整分销政策
- 提高运营灵活性

**Implementation**:
```
campaigns 表新增字段：
- enable_distribution: 是否启用分销
- distribution_level: 分销层级（1-3级）
- distribution_rewards: 各级奖励比例 JSON
  示例：{"level1": 10, "level2": 8, "level3": 5}
```

### 4. 多级分销层级限制

**Decision**: 最多支持3级分销

**Rationale**:
- 合规考虑：中国法规禁止超过3级的传销模式
- 业务平衡：足够激励推广，同时控制复杂度

**Implementation**:
```
一级分销：直接推荐人（referrer）
二级分销：推荐人的推荐人
三级分销：推荐人的推荐人的推荐人
超过3级不分配奖励
```

### 5. 二维码海报双模式

**Decision**: 支持活动专属海报和通用分销商海报

**Rationale**:
- 活动专属海报：针对性强，转化率高
- 通用分销商海报：方便分销商推广多个活动

**Implementation**:
```
活动专属海报：
- URL: /posters/campaign/{campaignId}?distributorId={id}
- 包含：活动信息 + 分销商二维码

通用分销商海报：
- URL: /posters/distributor/{distributorId}
- 包含：分销商信息 + 所有活动二维码
```

### 6. 提现审批流程

**Decision**: 提现需平台管理员审批

**Rationale**:
- 确保资金安全
- 防止恶意提现
- 平台可整体控制资金流向

**Implementation**:
```
分销商申请提现 →
  创建提现记录 →
  平台管理员审批 →
    通过：调用支付接口打款 → 更新状态
    拒绝：更新状态并记录原因
```

### 7. 数据查看双视角

**Decision**: 品牌管理员查看本品牌数据，平台管理员可按品牌筛选

**Rationale**:
- 数据隔离，保护隐私
- 平台可全局监控
- 满足不同角色需求

**Implementation**:
```
品牌管理员：
- 查看本品牌分销商列表
- 查看本品牌顾客列表
- 查看本品牌奖励详情
- 管理本品牌分销商（级别、状态）

平台管理员：
- 查看全部分销商（可按品牌筛选）
- 查看全部奖励（可按品牌筛选）
- 查看全部提现（可按品牌筛选）
- 审批提现申请
```

### 8. 分销商级别调整

**Decision**: 由品牌管理员手动调整

**Rationale**:
- 保持人工控制，避免自动升级的风险
- 品牌可根据分销商表现灵活调整

**Implementation**:
```
品牌管理员在后台可：
- 调整分销商级别（1-3级）
- 暂停/激活分销商
- 记录变更日志
```

### 9. 前端实现

**Decision**: 在现有H5基础上增加"分销中心"模块

**Rationale**:
- 用户无需安装额外应用
- 开发成本最低
- 与现有顾客功能无缝集成

## Architecture

### 数据模型

```
distributors 表：
- id: 分销商ID
- user_id: 关联用户ID
- brand_id: 关联品牌ID
- level: 分销级别(1/2/3)
- parent_id: 上级分销商ID
- status: 状态(active/suspended)
- created_at: 成为分销商时间
- total_earnings: 累计收益
- subordinates_count: 下级人数

campaigns 表扩展：
- enable_distribution: 是否启用分销
- distribution_level: 分销层级(1/2/3)
- distribution_rewards: 各级奖励比例 JSON
  示例：{"level1": 10, "level2": 8, "level3": 5}

orders 表扩展：
- distributor_path: 分销链路径
  示例："100,99,98"（一级,二级,三级）

withdrawals 表：
- id: 提现ID
- user_id: 用户ID
- brand_id: 品牌ID
- distributor_id: 分销商ID
- amount: 提现金额
- status: 状态(pending/approved/rejected/processing/completed/failed)
- pay_type: 提现方式(wechat/alipay/bank)
- pay_account: 提现账号
- pay_real_name: 真实姓名
- approved_by: 审批人ID
- approved_at: 审批时间
- approved_notes: 审批备注
- rejected_reason: 拒绝原因
- paid_at: 打款时间
- trade_no: 交易流水号

poster_templates 表：
- id: 模板ID
- type: 类型(campaign/distributor)
- campaign_id: 活动ID（活动海报才有）
- template_url: 模板URL
- created_at: 创建时间
```

### 自动成为分销商流程

```
订单支付成功 ->
  检查用户是否已是分销商 ->
    否：创建分销商记录（level=1, parent_id=referrer_id）
    是：跳过
  -> 原有奖励计算
```

### 多级奖励分配流程

```
订单支付成功 ->
  查询活动的分销奖励规则（campaigns.distribution_rewards）->
  查询分销链（distributor_path 或推荐关系）->
  最多3级，按级别比例计算奖励 ->
  创建多条奖励记录（每级一条）->
  更新各级分销商余额 ->
  发送通知
```

### 提现审批流程

```
分销商申请提现 ->
  创建提现记录（状态：pending）->
  平台管理员审批 ->
    通过 ->
      更新状态为 approved/processing ->
      调用支付接口打款 ->
      更新状态为 completed -> 通知分销商
    拒绝 ->
      更新状态为 rejected ->
      记录拒绝原因 -> 通知分销商
```

### 权限矩阵

| 功能 | participant | distributor | brand_admin | platform_admin |
|------|-------------|-------------|-------------|----------------|
| 参与活动 | ✓ | ✓ | ✓ | ✓ |
| 支付订单 | ✓ | ✓ | ✓ | ✓ |
| 自动成为分销商 | ✓ | - | - | - |
| 生成二维码海报 | - | ✓ | - | - |
| 查看推广数据 | - | ✓(自己) | ✓(本品牌) | ✓(全部) |
| 查看下级列表 | - | ✓(一级) | ✓(本品牌) | ✓(全部) |
| 申请提现 | - | ✓ | - | - |
| 审批提现 | - | - | - | ✓ |
| 管理分销商 | - | - | ✓(本品牌) | ✓(全部) |
| 调整分销商级别 | - | - | ✓(本品牌) | ✓(全部) |
| 暂停/激活分销商 | - | - | ✓(本品牌) | ✓(全部) |
| 查看顾客列表 | - | - | ✓(本品牌) | ✓(全部) |
| 查看奖励详情 | - | ✓(自己) | ✓(本品牌) | ✓(全部) |
| 查看提现明细 | - | ✓(自己) | - | ✓(全部) |

## API Design

### 生成二维码海报
```
POST /api/v1/posters/generate
Request: {
  type: "campaign" | "distributor",
  campaignId?: number,  // 活动海报必填
}
Response: {
  posterUrl: string,
  qrcodeUrl: string
}
```

### 查看推广数据
```
GET /api/v1/distributor/statistics
Response: {
  total_orders: number,
  total_earnings: number,
  subordinates_count: number,
  recent_orders: [],
  recent_earnings: []
}
```

### 申请提现
```
POST /api/v1/withdrawals/apply
Request: {
  amount: number,
  payType: "wechat" | "alipay" | "bank",
  payAccount: string,
  payRealName: string
}
Response: {
  withdrawalId: number,
  status: "pending"
}
```

### 查询提现记录
```
GET /api/v1/withdrawals?page=1&pageSize=20
Response: {
  total: number,
  withdrawals: [...]
}
```

### 平台审批提现
```
PUT /api/v1/platform/withdrawals/:id/approve
Request: { notes?: string }
Response: { status: "approved" }

PUT /api/v1/platform/withdrawals/:id/reject
Request: { reason: string }
Response: { status: "rejected" }
```

### 品牌管理员管理分销商
```
PUT /api/v1/brands/:brandId/distributors/:id/level
Request: { level: 1 | 2 | 3 }
Response: { level: number }

PUT /api/v1/brands/:brandId/distributors/:id/status
Request: { status: "active" | "suspended" }
Response: { status: string }
```

### 查看品牌数据
```
GET /api/v1/brands/:brandId/distributors
Response: {
  total: number,
  distributors: [...]
}

GET /api/v1/brands/:brandId/customers
Response: {
  total: number,
  customers: [...]
}

GET /api/v1/brands/:brandId/rewards
Response: {
  total: number,
  rewards: [...]
}
```

### 平台查看全局数据
```
GET /api/v1/platform/distributors?brandId={brandId}
Response: {
  total: number,
  distributors: [...]
}

GET /api/v1/platform/rewards?brandId={brandId}
Response: {
  total: number,
  rewards: [...]
}

GET /api/v1/platform/withdrawals?brandId={brandId}&status={status}
Response: {
  total: number,
  withdrawals: [...]
}
```

## Migration Plan

1. 创建数据库表：distributors, withdrawals, poster_templates
2. 扩展现有表：campaigns, orders
3. 修改 user_roles 表支持 distributor 角色
4. 部署后端API
5. 更新H5前端添加分销中心、海报生成、提现功能
6. 更新品牌管理后台添加分销商管理、数据查看功能
7. 更新平台管理后台添加提现审批、全局数据查看功能
8. 数据迁移：现有推荐人可选择升级为分销商（可选）

## Open Questions

1. ~~分销商级别是自动升级还是手动调整？~~ -> 品牌管理员手动调整
2. ~~分销商是否有有效期？~~ -> 无有效期，除非手动暂停
3. ~~是否支持分销商转让/更换上级？~~ -> 暂不支持
4. ~~提现是否需要审核？谁审核？~~ -> 平台管理员审批
5. ~~分销奖励规则是全局统一还是活动级别？~~ -> 活动级别自定义
