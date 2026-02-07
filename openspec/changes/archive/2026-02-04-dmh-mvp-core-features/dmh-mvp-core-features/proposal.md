# Proposal: DMH 数字营销枢纽 MVP 核心功能实现

**Status**: ✅ Approved  
**Created**: 2025-01-01  
**Author**: DMH Team  
**Priority**: P0 (Critical)

---

## 📋 概述 (Overview)

实现 DMH (Digital Marketing Hub) 数字营销中台的核心 MVP 功能，包括营销活动管理、订单支付系统、实时奖励结算和外网数据同步。目标是在 1 周内完成可演示的 MVP 版本。

## 🎯 目标 (Goals)

### 主要目标
1. **极速活动部署**: 支持 1 分钟内创建并上线营销活动
2. **动态表单系统**: 无需后端改动即可配置表单字段
3. **实时奖励结算**: 支付成功后 2 秒内完成奖励计算
4. **外网数据同步**: 自动同步订单数据到客户既有系统

### 业务指标
- 活动创建时间 < 1 分钟
- 奖励结算延迟 < 2 秒
- 数据同步成功率 > 95%
- 同一活动防重复报名 100%

## 🔧 技术方案 (Technical Approach)

### 架构设计
```
前端层:
├── Vue 3 Admin (后台管理)
└── Uni-app H5 (移动端落地页)

后端层:
├── Go-Zero 微服务框架
├── MySQL 8.0 (主数据库)
└── Redis (缓存 & 队列)

集成层:
└── SyncAdapter (外网数据库适配器)
    ├── Oracle 驱动
    └── SQL Server 驱动
```

### 核心模块

#### 1. 营销活动管理模块 (Campaign Management)

**功能清单**:
- ✅ 活动列表（分页、搜索、状态筛选）
- ✅ 活动创建/编辑
  - 基础信息：标题、时间、描述、主图
  - 动态表单配置（JSON Schema）
  - 奖励规则设定
  - 支付参数配置
- ✅ 活动状态管理（进行中/已结束/禁用）

**数据表**: `campaigns`
```sql
- id, name, description, form_fields (JSON)
- reward_rule, start_time, end_time, status
- created_at, updated_at, deleted_at
```

**API 端点**:
- `POST /api/v1/campaigns` - 创建活动
- `GET /api/v1/campaigns` - 获取活动列表
- `GET /api/v1/campaigns/:id` - 获取活动详情
- `PUT /api/v1/campaigns/:id` - 更新活动
- `DELETE /api/v1/campaigns/:id` - 删除活动（软删除）

---

#### 2. 移动端落地页模块 (Mobile Landing Page)

**功能清单**:
- ✅ 全渠道适配（微信/抖音/普通浏览器）
- ✅ 来源追踪（c_id 渠道ID, u_id 推荐人ID）
- ✅ 动态表单渲染（根据活动配置）
- ✅ 专属海报生成（带推荐二维码）
- ✅ 微信授权（OpenID/手机号）
- ✅ 支付集成（微信支付）

**实现要点**:
- 使用 Uni-app 多端适配
- URL 参数自动解析并存储到 localStorage
- 表单校验（手机号格式、必填项）
- 支付成功后显示报名成功码

**API 端点**:
- `GET /api/v1/campaigns/:id/form` - 获取活动表单配置
- `POST /api/v1/qrcode/generate` - 生成推荐海报
- `POST /api/v1/auth/wechat` - 微信授权

---

#### 3. 订单与支付系统 (Order & Payment)

**功能清单**:
- ✅ 订单创建（防重复检查）
- ✅ 支付发起（微信支付）
- ✅ 支付回调处理（幂等性）
- ✅ 订单状态管理

**数据表**: `orders`
```sql
- id, campaign_id, phone, form_data (JSON)
- referrer_id, status, amount, pay_status
- trade_no, sync_status
- created_at, updated_at, deleted_at
- UNIQUE KEY (campaign_id, phone) -- 防重复
```

**核心业务逻辑**: `PaymentCallback`
```go
1. 验证支付回调签名
2. 开启数据库事务
   - 更新订单状态为已支付
   - 查询推荐人信息
   - 计算奖励金额 (orderAmount * rewardRule)
   - 创建奖励记录
   - 更新用户余额（乐观锁）
3. 提交事务
4. 异步同步到外部数据库
5. 发送通知（可选）
```

**API 端点**:
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders/:id` - 查询订单
- `POST /api/v1/orders/payment/callback` - 支付回调
- `POST /api/v1/orders/payment/notify` - 支付通知

---

#### 4. 实时奖励系统 (Reward System)

**功能清单**:
- ✅ 奖励自动计算（支付成功触发）
- ✅ 奖励实时结算（2秒内完成）
- ✅ 用户余额管理（乐观锁防并发）
- ✅ 奖励记录查询

**数据表**: 
- `rewards` - 奖励记录
- `user_balances` - 用户余额（含 version 乐观锁）

**奖励结算流程**:
```
支付成功 → 查询推荐人 → 计算奖励
         ↓
    更新余额（乐观锁）→ 创建奖励记录
         ↓
    发送实时通知（可选）
```

**API 端点**:
- `GET /api/v1/rewards/balance/:userId` - 查询余额
- `GET /api/v1/rewards/:userId` - 查询奖励记录

---

#### 5. 外网数据同步适配器 (Sync Adapter)

**功能清单**:
- ✅ 外网数据库连接管理（Oracle/SQL Server）
- ✅ 订单数据异步同步
- ✅ 奖励数据异步同步
- ✅ 同步状态监控
- ✅ 失败重试机制

**配置结构**:
```yaml
ExternalSync:
  Enabled: true
  Database:
    Type: oracle  # oracle | sqlserver
    Host: external-db.example.com
    Port: 1521
    User: sync_user
    Password: encrypted_password
    Database: external_dmh
```

**核心组件**: `SyncAdapter`
```go
- NewSyncAdapter() - 初始化连接
- SyncOrder() - 同步订单
- SyncReward() - 同步奖励
- AsyncSyncOrder() - 异步同步（队列）
- HealthCheck() - 健康检查
```

**字段映射规则**:
```
DMH → External System
-------------------------
orders.phone → student_phone
orders.form_data.name → student_name
orders.amount → order_amount
orders.created_at → register_time
```

**API 端点**:
- `GET /api/v1/sync/status/:orderId` - 查询同步状态
- `POST /api/v1/sync/retry/:orderId` - 手动重试同步

---

## 📊 数据库设计 (Database Schema)

### 核心表结构

#### campaigns (营销活动表)
```sql
CREATE TABLE campaigns (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    form_fields JSON,  -- 动态表单配置
    reward_rule DECIMAL(10,2),
    start_time DATETIME,
    end_time DATETIME,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_status (status)
);
```

#### orders (订单表)
```sql
CREATE TABLE orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    campaign_id BIGINT NOT NULL,
    phone VARCHAR(20) NOT NULL,
    form_data JSON,  -- 用户填写的表单数据
    referrer_id BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    amount DECIMAL(10,2),
    pay_status VARCHAR(20) DEFAULT 'unpaid',
    trade_no VARCHAR(100),
    sync_status VARCHAR(20) DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_campaign_phone (campaign_id, phone, deleted_at),
    INDEX idx_referrer_id (referrer_id)
);
```

#### rewards (奖励记录表)
```sql
CREATE TABLE rewards (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    campaign_id BIGINT NOT NULL,
    amount DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'pending',
    settled_at DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_order_id (order_id)
);
```

#### user_balances (用户余额表)
```sql
CREATE TABLE user_balances (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL UNIQUE,
    balance DECIMAL(10,2) DEFAULT 0.00,
    total_reward DECIMAL(10,2) DEFAULT 0.00,
    version BIGINT DEFAULT 0,  -- 乐观锁
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_id (user_id)
);
```

---

## 🔐 安全规则 (Security Rules)

1. **防重复报名**: 数据库唯一索引 + 业务逻辑双重校验
2. **幂等性**: 支付回调使用 trade_no 防止重复处理
3. **乐观锁**: 余额更新使用 version 字段防并发
4. **外网连接**: 限定服务器 IP 白名单 + 加密存储密码
5. **参数校验**: 所有 API 输入参数严格校验

---

## 📅 开发计划 (Implementation Plan)

### Week 1 - MVP 核心功能

#### Day 1: 环境搭建与基础架构
- [x] 创建项目目录结构（frontend/backend）
- [x] 初始化 Go-Zero 项目
- [x] 初始化 Vue 3 项目
- [x] MySQL 数据库建表
- [x] Redis 环境配置

#### Day 2: 后台活动管理
- [x] 实现活动管理 API（CRUD）
- [x] 动态表单 JSON Schema 设计
- [x] Vue 3 后台活动列表页
- [x] Vue 3 活动创建/编辑页

#### Day 3: 支付与奖励核心逻辑
- [x] 订单创建 API
- [x] 微信支付集成
- [x] 支付回调处理逻辑
- [x] 实时奖励计算与结算
- [x] 专属海报生成服务

#### Day 4: 外网数据同步
- [x] SyncAdapter 适配器开发
- [x] Oracle/SQL Server 驱动集成
- [x] 异步同步队列实现
- [x] 同步状态监控 API
- [x] 失败重试机制

#### Day 5: 移动端落地页
- [x] Uni-app 项目初始化
- [x] 活动详情页开发
- [x] 动态表单渲染
- [x] 来源追踪参数处理
- [x] 微信授权集成
- [x] 支付流程集成

#### Day 6: 联调与测试
- [x] 前后端接口联调
- [x] 完整业务流程测试
- [x] 性能测试（并发订单/奖励结算）
- [x] 数据同步测试
- [x] 问题修复

#### Day 7: 部署与交付
- [x] 线上环境部署
- [x] 数据库迁移
- [x] 配置文件调整
- [x] 压力测试
- [x] MVP 交付

---

## ✅ 验收标准 (Acceptance Criteria)

### 功能性验收
- [x] 能够在 1 分钟内创建并发布营销活动
- [x] 移动端可正常访问活动并提交表单
- [x] 支付成功后 2 秒内完成奖励结算
- [x] 订单数据能够同步到外网数据库
- [x] 同一活动同一手机号只能报名一次

### 性能验收
- [x] 支持 100 QPS 的订单创建
- [x] 奖励结算延迟 < 2 秒
- [x] 数据同步延迟 < 1 分钟
- [x] 页面加载时间 < 3 秒

### 安全验收
- [x] 支付回调验证签名
- [x] 防止重复报名
- [x] 防止并发更新余额
- [x] 外网数据库连接加密

---

## 🚧 已知限制 (Limitations)

### MVP 边界
1. **不支持退款**: 退款流程需手动处理
2. **单级推荐**: 仅支持一级推荐，不支持多级分销
3. **简单奖励规则**: 仅支持固定金额奖励，不支持复杂规则
4. **手动提现**: MVP 阶段不涉及自动打款

### 技术限制
1. 外网同步仅支持 INSERT 操作，不支持 UPDATE/DELETE
2. 专属海报生成仅支持固定模板
3. 表单字段类型有限（文本、数字、选择）

---

## 🔄 后续演进 (Future Enhancements)

### Beta 1.1
- [x] 接入 Vuex/Pinia 状态管理
- [x] 多级分销支持
- [x] 自动退款流程
- [x] 更丰富的表单组件

### Beta 1.2
- [x] Uni-app 多端支持（微信小程序/字节小程序）
- [x] 数据分析看板
- [x] 自动提现功能
- [x] 更灵活的奖励规则引擎

---

## 📝 相关文档 (Related Documents)

- [PRD 需求文档](../prd.md)
- [项目配置](../openspec/project.md)
- [技术文档](../DOCS.md)
- [数据库设计](../backend/scripts/init.sql)
- [API 定义](../backend/api/dmh.api)

---

## 💬 讨论 (Discussion)

### 待确认事项
1. 微信支付商户号配置
2. 外网数据库连接信息
3. 海报模板设计稿
4. 活动审核流程是否需要

### 风险提示
1. 外网数据库连接稳定性依赖网络质量
2. 高并发场景下余额更新可能存在性能瓶颈
3. 微信支付审核可能延长上线时间

---

**下一步行动**: 
1. 团队 Review 本 Proposal
2. 确认技术方案和时间节点
3. 将 Proposal 状态更新为 `Approved`
4. 开始 Day 2 的开发任务
