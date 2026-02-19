# DMH 数字营销中台 - 系统需求文档

> 生成日期: 2026-02-14
> 文档版本: 1.0
> 数据来源: OpenSpec 规格与归档变更汇总

---

## 目录

1. [文档概述](#1-文档概述)
2. [系统概述](#2-系统概述)
3. [当前有效规格](#3-当前有效规格)
   - [3.1 活动管理系统 (campaign-management)](#31-活动管理系统-campaign-management)
   - [3.2 订单支付系统 (order-payment-system)](#32-订单支付系统-order-payment-system)
   - [3.3 RBAC 权限系统 (rbac-permission-system)](#33-rbac-权限系统-rbac-permission-system)
   - [3.4 反馈系统 (feedback-system)](#34-反馈系统-feedback-system)
   - [3.5 规格治理 (spec-governance)](#35-规格治理-spec-governance)
4. [历史变更记录](#4-历史变更记录)
5. [测试指导建议](#5-测试指导建议)

---

## 1. 文档概述

### 1.1 目的

本文档汇总 DMH (Digital Marketing Hub) 数字营销中台系统的所有需求规格，包括：
- 当前有效的 5 个能力规格
- 16 个已归档的历史变更
- 完整的需求场景与验收标准

旨在为全面测试提供权威的需求参考。

### 1.2 规格统计

| 规格 | 需求数 | 场景数 | 状态 |
|------|--------|--------|------|
| campaign-management | 6 | 24 | ✅ Active |
| order-payment-system | 1 | 3 | ✅ Active |
| rbac-permission-system | 12 | 38 | ✅ Active |
| feedback-system | 9 | 29 | ✅ Active |
| spec-governance | 2 | 2 | ✅ Active |
| **总计** | **30** | **96** | - |

---

## 2. 系统概述

### 2.1 项目背景

DMH (Digital Marketing Hub) 是一个数字营销中台系统，提供营销活动管理、会员管理、分销商管理、奖励结算等功能。

### 2.2 核心目标

- **极速部署**: 1分钟内完成营销活动上线
- **动态表单**: 支持灵活的数据采集，无需后端改动
- **实时激励**: 支付即结算，显著提升推广员积极性
- **无缝集成**: 内置外网数据库适配器

### 2.3 技术架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   H5 Frontend   │    │  Admin Frontend │    │   Backend API   │
│   (Vue.js 3)    │    │   (Vue.js 3)    │    │   (go-zero)     │
│   Port: 3100    │    │   Port: 3000    │    │   Port: 8889    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐    ┌─────────────────┐
                    │     MySQL       │    │     Redis       │
                    │   Database      │    │     Cache       │
                    └─────────────────┘    └─────────────────┘
```

### 2.4 用户角色体系

| 角色 | 标识 | 描述 |
|------|------|------|
| 平台管理员 | platform_admin | 全系统管理权限 |
| 品牌管理员 | brand_admin | 管理本品牌活动和数据 |
| 分销商 | distributor | 高级顾客，可推广获奖励 |
| 活动参与者 | participant | 普通用户，参与活动 |
| 匿名用户 | anonymous | 未登录访客 |

---

## 3. 当前有效规格

### 3.1 活动管理系统 (campaign-management)

**能力概述**: H5 活动页面设计器，允许品牌管理员设计活动落地页

#### 需求清单

| ID | 需求名称 | 场景数 |
|----|----------|--------|
| CM-01 | H5 Campaign Page Designer | 7 |
| CM-02 | Supported Page Components | 11 |
| CM-03 | Mobile-Optimized Interaction | 4 |
| CM-04 | Data Persistence | 3 |
| CM-05 | Brand admin generates campaign posters | 1 |
| CM-06 | Brand admin configures distribution rules | 1 |

#### 详细需求

##### CM-01: H5 Campaign Page Designer

H5 应用应提供活动页面设计器，允许品牌管理员设计活动落地页。

**场景**:

| 场景 | WHEN | THEN |
|------|------|------|
| Access page designer | 品牌管理员编辑活动 | 可访问页面设计器，看到移动端优化设计界面 |
| Add page components | 添加组件 | 可从组件库选择，组件添加到画布，收到视觉反馈 |
| Edit component content | 点击组件 | 底部弹出编辑器，可编辑属性，更改即时反映 |
| Reorder components | 更改组件顺序 | 可使用上下箭头，组件移动到新位置 |
| Delete components | 删除组件 | 可删除，确认后删除，从画布移除 |
| Preview page design | 点击预览 | 全屏预览，显示实际移动端外观 |
| Save page configuration | 保存设计 | 配置持久化到后端，收到确认，关联到活动 |
| Load existing configuration | 打开现有活动 | 加载保存的配置，组件正确恢复 |

##### CM-02: Supported Page Components

页面设计器应支持以下组件类型：

| 组件 | 功能 |
|------|------|
| Activity poster | 上传/选择图片，正确宽高比 |
| Activity title | 输入标题/副标题，AI调色 |
| Activity time | 选择日期时间，时间范围显示 |
| Activity location | 输入位置，可选地图集成 |
| Activity highlights | 多个亮点，图标选择，重排序 |
| Activity details | 富文本，基本格式，图片插入 |
| Ticket information | 票种定义，价格，可用性 |
| Refund policy | 政策文本，多规则，格式显示 |
| Payment and invoice | 支付方式配置，发票选项 |
| Organizer information | 主办方详情，logo，联系方式 |
| Divider | 视觉分隔符，自定义样式 |
| Registration button | 自定义文本，颜色样式，动作配置 |

##### CM-03: Mobile-Optimized Interaction

移动端优化交互模式：

- 组件库通过底部动作表显示
- 编辑器通过底部弹出显示，占屏幕高度 70%
- 触摸目标至少 44x44 像素
- 响应式布局适配不同屏幕

##### CM-04: Data Persistence

页面配置数据持久化：

- JSON 格式存储所有组件数据
- 保存前验证必填字段
- 定期自动保存
- 退出前警告未保存更改

##### CM-05: 品牌管理员生成活动海报

系统允许品牌管理员为活动生成海报。

- 点击"生成海报"后使用模板生成
- 返回可预览/下载的海报 URL

##### CM-06: 品牌管理员配置分销规则

系统允许品牌管理员在活动编辑时配置分销规则。

- 启用分销后可选择 1-3 级
- 可填写各级奖励比例
- 配置保存到活动数据

---

### 3.2 订单支付系统 (order-payment-system)

**能力概述**: 订单核心业务逻辑的自动化回归测试闭环

#### 需求清单

| ID | 需求名称 | 场景数 |
|----|----------|--------|
| OP-01 | Automated regression closure for order logic | 3 |

#### 详细需求

##### OP-01: 订单逻辑自动化回归闭环

系统应为订单核心业务逻辑提供可重复执行的自动化回归测试闭环，覆盖创建、核销和关键字段校验路径。

**场景**:

| 场景 | WHEN | THEN |
|------|------|------|
| Order creation regression | 代码变更涉及订单创建逻辑 | 执行订单创建回归单元测试，覆盖活动有效性、重复订单防护与关键字段格式校验 |
| Order verification regression | 代码变更涉及订单核销或取消核销逻辑 | 执行对应回归单元测试，覆盖权限不足、重复核销和状态回滚分支 |
| Minimal integration guardrail | 变更准备合入主干 | 执行最小订单关键路径集成/冒烟测试，产出可追溯的测试结果记录 |

---

### 3.3 RBAC 权限系统 (rbac-permission-system)

**能力概述**: 完整的基于角色的访问控制系统

#### 需求清单

| ID | 需求名称 | 场景数 |
|----|----------|--------|
| RBAC-01 | 用户认证管理 | 5 |
| RBAC-02 | 角色权限体系 | 4 |
| RBAC-03 | API权限控制 | 3 |
| RBAC-04 | 用户注册权限控制 | 4 |
| RBAC-05 | 用户管理功能 | 5 |
| RBAC-06 | 权限缓存优化 | 2 |
| RBAC-07 | 品牌管理员关系管理 | 5 |
| RBAC-08 | 品牌管理功能 | 4 |
| RBAC-09 | 提现权限控制 | 3 |
| RBAC-10 | 菜单权限管理 | 5 |
| RBAC-11 | 安全审计日志 | 3 |
| RBAC-12 | 分销监控视图可靠性 | 2 |

#### 详细需求

##### RBAC-01: 用户认证管理

系统提供完整的用户认证功能，包括注册、登录、密码管理和会话控制。

**场景**:

| 场景 | WHEN | THEN |
|------|------|------|
| H5用户注册成功 | 通过H5提供有效用户名、密码、手机号 | 创建账号并分配 participant 角色，返回 JWT token |
| 平台管理员后台创建用户 | 平台管理员通过后台创建用户 | 创建账号并分配指定角色，平台管理员角色只能通过后台创建 |
| 用户登录成功 | 提供正确用户名和密码 | 验证身份并生成 JWT token，token 包含用户ID、角色、有效期 |
| 密码安全验证 | 设置或修改密码 | bcrypt 加密存储，长度≥6位 |
| 会话超时控制 | JWT token 超过24小时 | 拒绝请求返回 401，前端清除 token 跳转登录 |

##### RBAC-02: 角色权限体系

系统实现 4 种用户角色，每种角色具有明确的权限范围。

**角色权限矩阵**:

| 角色 | 权限范围 |
|------|----------|
| platform_admin | 所有功能模块，管理所有品牌、活动、用户、系统配置 |
| brand_admin | 只能访问管理的品牌数据，可管理品牌信息、素材库、活动、数据统计 |
| participant | 只能访问个人相关功能，参与活动、查看奖励、申请提现 |
| anonymous | 只能访问公开功能，浏览活动、查看详情、注册账号 |

##### RBAC-03: API权限控制

系统在 API 层面实现细粒度权限控制。

- JWT token 验证（无效返回 401）
- 根据 URL 和 HTTP 方法确定所需权限
- 数据级权限隔离（品牌管理员只能查询自己品牌的数据）

##### RBAC-04: 用户注册权限控制

根据用户角色类型实现不同的注册方式和权限控制。

- H5 注册只能创建 participant 角色
- 品牌管理员必须由平台管理员通过后台创建
- 平台管理员只能由现有平台管理员创建
- 匿名用户完成 H5 注册后转换为 participant

##### RBAC-05: 用户管理功能

完整的用户管理功能，支持创建、查询、更新、状态管理和密码重置。

- 后台创建用户时验证用户名和手机号唯一性
- 用户状态管理（active/disabled/locked）
- 密码重置生成临时密码，强制下次登录修改
- 角色分配立即生效
- 品牌管理员可同时管理多个品牌

##### RBAC-06: 权限缓存优化

权限信息缓存机制提高性能。

- 首次查询权限信息缓存到内存
- 权限变更时立即清除相关缓存

##### RBAC-07: 品牌管理员关系管理

品牌管理员与品牌关系的完整管理。

- 绑定/解绑品牌管理权限
- 变更立即生效
- 支持一个用户管理多个品牌
- 记录所有操作日志

##### RBAC-08: 品牌管理功能

品牌管理员完整的品牌管理功能。

- 品牌信息管理（名称、描述、logo）
- 品牌素材库管理（上传、分类、删除）
- 品牌活动管理（创建、编辑、删除、发布）
- 品牌数据查看（完整数据报表）

##### RBAC-09: 提现权限控制

提现申请和审核的权限控制。

- 提现申请：检查用户是否已认证、余额是否足够
- 提现审核：只有平台管理员可审批
- 提现状态更新：使用数据库事务确保一致性

##### RBAC-10: 菜单权限管理

完整的菜单权限管理功能。

- 菜单结构管理（增删改排序，多级菜单）
- 页面操作权限配置（增删改查、导出、转发）
- 角色菜单权限分配
- 用户登录后返回可访问菜单列表
- 支持权限继承机制

##### RBAC-11: 安全审计日志

所有重要安全操作记录审计日志。

- 用户操作日志（用户ID、操作类型、时间戳、IP）
- 权限变更日志（变更前后状态）
- 安全事件监控（异常登录、权限滥用）

##### RBAC-12: 分销监控视图可靠性

为平台管理员提供稳定可用的分销监控页面。

- 请求超时时显示明确提示并提供重试入口
- 筛选和搜索通过防抖避免重复请求

---

### 3.4 反馈系统 (feedback-system)

**能力概述**: 用户反馈从提交、流转处理到统计分析的完整闭环

#### 需求清单

| ID | 需求名称 | 场景数 |
|----|----------|--------|
| FB-01 | 用户反馈提交 | 4 |
| FB-02 | 反馈列表查询 | 3 |
| FB-03 | 反馈详情查看 | 3 |
| FB-04 | 反馈状态更新 | 4 |
| FB-05 | 满意度调查 | 3 |
| FB-06 | FAQ管理 | 3 |
| FB-07 | 功能使用统计 | 3 |
| FB-08 | 反馈统计分析 | 5 |
| FB-09 | 反馈数据模型 | 6 |

#### 详细需求

##### FB-01: 用户反馈提交

系统允许用户提交功能反馈。

**字段要求**:
- 必填：标题、内容、类别
- 可选：评分（1-5）、子分类
- 自动记录：用户ID、设备信息、浏览器信息

**支持分类**: poster(海报), payment(支付), verification(核销), other(其他)

**优先级**: low, medium, high（默认 medium）

##### FB-02: 反馈列表查询

根据用户角色和筛选条件返回反馈列表。

- 普通用户只能查看自己的反馈
- 管理员可查看所有反馈
- 支持按类别、状态、优先级筛选
- 支持分页查询
- 按创建时间倒序

##### FB-03: 反馈详情查看

允许用户查看单个反馈详情。

- 用户只能查看自己的反馈详情
- 管理员可查看任意反馈
- 普通用户访问他人反馈返回 403

##### FB-04: 反馈状态更新

允许管理员更新反馈状态并添加处理回复。

**支持状态**: pending(待处理), reviewing(审核中), resolved(已解决), closed(已关闭)

- 更新状态时记录处理人和时间
- 解决时记录解决时间用于统计
- 非管理员更新返回 403

##### FB-05: 满意度调查

允许用户对系统功能进行满意度评分。

**评分维度**: 易用性、性能、稳定性、整体满意度、推荐意愿

- 所有评分在 1-5 范围
- 可填写最喜欢/最不喜欢的功能
- 可填写改进建议和新功能需求

##### FB-06: FAQ管理

常见问题的查询和反馈功能。

- 查询所有已发布的 FAQ
- 支持按分类筛选和关键词搜索
- 支持点赞/踩反馈（防重复）
- 按 sort_order 升序排列

##### FB-07: 功能使用统计

记录用户功能使用情况。

- 记录功能名称、操作类型、执行结果、耗时
- 失败时记录错误信息
- 关联活动ID
- 异步记录避免影响性能

##### FB-08: 反馈统计分析

为管理员提供反馈数据统计分析。

- 反馈总数和按类别/状态/优先级统计
- 平均评分和解决率
- 平均解决时间
- 评分分布（1-5星）
- 支持时间范围筛选
- 非管理员访问返回 403

##### FB-09: 反馈数据模型

| 模型 | 用途 | 关键字段 |
|------|------|----------|
| UserFeedback | 用户反馈 | user_id, category, rating, title, content, status, resolved_at |
| FeatureSatisfactionSurvey | 满意度调查 | user_id, feature, ease_of_use, performance, overall_satisfaction |
| FAQItem | FAQ | category, question, answer, sort_order, is_published |
| FeatureUsageStat | 功能使用统计 | user_id, feature, action, success, duration_ms |
| FeedbackTag | 反馈标签 | name, color |
| FeedbackTagRelation | 标签关联 | feedback_id, tag_id |

---

### 3.5 规格治理 (spec-governance)

**能力概述**: OpenSpec 归档变更的状态一致性管理

#### 需求清单

| ID | 需求名称 | 场景数 |
|----|----------|--------|
| SG-01 | Archive task status consistency | 1 |
| SG-02 | Archive status traceability index | 1 |

#### 详细需求

##### SG-01: 归档任务状态一致性

OpenSpec 归档变更中保持 tasks.md 状态表达一致。

- 同一任务出现冲突状态或重复定义时统一为单一权威状态
- 保留解释性说明体现历史真实性

##### SG-02: 归档状态可追溯索引

为归档变更提供统一索引。

- 提供可导航的索引入口
- 标注状态口径来源与关联变更

---

## 4. 历史变更记录

### 4.1 归档变更列表

| 归档日期 | 变更ID | 主要内容 | 影响规格 |
|----------|--------|----------|----------|
| 2026-01-24 | add-h5-campaign-page-designer | H5 活动页面设计器 | campaign-management |
| 2026-01-24 | enhance-rbac-permission-system | RBAC 权限系统完善 | rbac-permission-system |
| 2026-01-28 | add-distributor-role | 分销商角色与多级分销系统 | rbac-permission-system, campaign-management |
| 2026-01-28 | add-member-system | 会员系统（UnionID 平台唯一） | member-system, campaign-management, order-payment-system |
| 2026-01-28 | fix-campaign-formfields-api | 活动表单字段 API 修复 | campaign-management |
| 2026-01-28 | fix-distributor-view-architecture | 分销商视图架构修复 | rbac-permission-system |
| 2026-01-28 | remove-brand-admin-features | 移除品牌管理员部分功能 | campaign-management, admin-management |
| 2026-02-02 | order-logic-implementation | 订单业务逻辑实现 | order-payment-system |
| 2026-02-03 | add-campaign-advanced-features | 活动高级功能 | campaign-management, order-payment-system |
| 2026-02-04 | dmh-mvp-core-features | MVP 核心功能实现 | 全系统 |
| 2026-02-06 | add-order-rbac-test-backfill | 订单 RBAC 测试补填 | rbac-permission-system, order-payment-system |
| 2026-02-07 | add-brand-admin-poster-distribution | 品牌管理员海报生成与分销配置 | campaign-management |
| 2026-02-07 | add-order-logic-test-gap-closure | 订单逻辑测试闭环 | order-payment-system |
| 2026-02-07 | refactor-openspec-archive-task-normalization | OpenSpec 归档任务规范化 | spec-governance |
| 2026-02-07 | update-distributor-view-stability-performance | 分销商视图稳定性与性能 | rbac-permission-system |
| 2026-02-10 | add-feedback-system | 反馈系统 | feedback-system |

### 4.2 重大变更详情

#### 4.2.1 MVP 核心功能 (2026-02-04)

**变更ID**: dmh-mvp-core-features

**核心模块**:
1. 营销活动管理模块
2. 移动端落地页模块
3. 订单与支付系统
4. 实时奖励系统
5. 外网数据同步适配器

**验收标准**:
- 1分钟内创建并发布营销活动
- 支付成功后 2 秒内完成奖励结算
- 支持 100 QPS 订单创建
- 同一活动防重复报名 100%

#### 4.2.2 分销商系统 (2026-01-28)

**变更ID**: add-distributor-role

**新增功能**:
- `distributor` 角色定义
- 支付后自动成为分销商机制
- 活动级别多级分销体系（最多3级）
- 分销商专属功能（海报生成、推广数据、提现）

**合规要求**: 分销层级不超过3级

#### 4.2.3 会员系统 (2026-01-28)

**变更ID**: add-member-system

**核心能力**:
- 以 `unionid` 作为平台唯一标识
- 会员档案沉淀（基础资料、标签、来源渠道）
- 人工合并流程
- 导出审批流程

#### 4.2.4 订单逻辑实现 (2026-02-02)

**变更ID**: order-logic-implementation

**实现内容**:
- 订单创建逻辑（活动验证、防重复、核销码生成）
- 订单核销逻辑（签名验证、状态更新、权限检查）
- 表单字段验证服务（7种字段类型）
- 完整单元测试

#### 4.2.5 反馈系统 (2026-02-10)

**变更ID**: add-feedback-system

**实现内容**:
- 6个数据模型
- 9个 API 函数实现
- 完整权限控制
- 数据统计分析

---

## 5. 测试指导建议

### 5.1 测试覆盖矩阵

| 模块 | 单元测试 | 集成测试 | E2E测试 | 优先级 |
|------|----------|----------|---------|--------|
| RBAC 权限系统 | ✅ | ✅ | ✅ | P0 |
| 订单支付系统 | ✅ | ✅ | ✅ | P0 |
| 活动管理系统 | ✅ | ✅ | ✅ | P0 |
| 反馈系统 | ⚠️ | ⚠️ | ⚠️ | P1 |
| 分销商系统 | ⚠️ | ⚠️ | ⚠️ | P1 |
| 会员系统 | ⚠️ | ⚠️ | ⚠️ | P1 |

### 5.2 关键测试场景

#### 5.2.1 认证与授权 (P0)

```
TC-AUTH-01: H5 用户注册并获取 JWT token
TC-AUTH-02: 平台管理员后台创建用户
TC-AUTH-03: JWT token 过期处理
TC-AUTH-04: 品牌管理员数据隔离验证
TC-AUTH-05: 分销商权限边界验证
```

#### 5.2.2 订单流程 (P0)

```
TC-ORDER-01: 订单创建（正常流程）
TC-ORDER-02: 订单创建（重复手机号拦截）
TC-ORDER-03: 订单核销（品牌管理员）
TC-ORDER-04: 订单核销（权限不足拦截）
TC-ORDER-05: 取消核销
TC-ORDER-06: 扫码获取订单
```

#### 5.2.3 活动管理 (P0)

```
TC-CAMP-01: H5 页面设计器组件添加
TC-CAMP-02: 页面配置保存与加载
TC-CAMP-03: 海报生成
TC-CAMP-04: 分销规则配置
TC-CAMP-05: 动态表单验证
```

#### 5.2.4 奖励系统 (P1)

```
TC-REWARD-01: 支付回调后奖励计算
TC-REWARD-02: 多级分销奖励分配
TC-REWARD-03: 余额更新并发控制
TC-REWARD-04: 提现申请与审批
```

#### 5.2.5 反馈系统 (P1)

```
TC-FEED-01: 用户提交反馈
TC-FEED-02: 管理员更新反馈状态
TC-FEED-03: 满意度调查提交
TC-FEED-04: FAQ 查询与点赞
TC-FEED-05: 反馈统计分析
```

### 5.3 测试命令参考

```bash
# 后端单元测试
cd backend && go test ./... -v

# 后端集成测试
cd backend
export DMH_INTEGRATION_BASE_URL=http://localhost:8889
go test ./test/integration/... -v -count=1

# 订单回归测试
backend/scripts/run_order_mysql8_regression.sh

# 前端单元测试
cd frontend-admin && npm run test
cd frontend-h5 && npm run test

# 前端 E2E 测试
cd frontend-admin && npm run test:e2e
cd frontend-h5 && npm run test:e2e

# OpenSpec 验证
openspec validate --strict --no-interactive
```

### 5.4 测试数据准备

```bash
# 初始化测试数据库
cd backend/scripts && mysql -u root -p < init.sql

# 创建测试用户
# 平台管理员: admin / 123456
# 品牌管理员: brand_manager / 123456
```

### 5.5 测试环境

| 服务 | 地址 | 用途 |
|------|------|------|
| 后端 API | http://localhost:8889 | API 测试 |
| 管理后台 | http://localhost:3000 | E2E 测试 |
| H5 前端 | http://localhost:3100 | E2E 测试 |
| MySQL | localhost:3306 | 数据验证 |
| Redis | localhost:6379 | 缓存验证 |

---

## 附录

### A. 参考资料

- [OpenSpec 项目规范](/opt/code/dmh/openspec/project.md)
- [OpenSpec AGENTS.md](/opt/code/dmh/openspec/AGENTS.md)
- [归档状态索引](/opt/code/dmh/openspec/changes/archive/ARCHIVE_STATUS_INDEX.md)

### B. 规格文件路径

| 规格 | 路径 |
|------|------|
| campaign-management | `/opt/code/dmh/openspec/specs/campaign-management/spec.md` |
| order-payment-system | `/opt/code/dmh/openspec/specs/order-payment-system/spec.md` |
| rbac-permission-system | `/opt/code/dmh/openspec/specs/rbac-permission-system/spec.md` |
| feedback-system | `/opt/code/dmh/openspec/specs/feedback-system/spec.md` |
| spec-governance | `/opt/code/dmh/openspec/specs/spec-governance/spec.md` |

### C. 修订历史

| 版本 | 日期 | 描述 |
|------|------|------|
| 1.0 | 2026-02-14 | 初始版本，汇总 5 个有效规格和 16 个归档变更 |

---

*本文档由 OpenSpec 系统自动生成，如需更新请重新运行生成脚本。*
