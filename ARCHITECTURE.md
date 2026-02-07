# 🏗️ DMH 系统架构文档

## 📋 目录

* [系统概述](#系统概述)
* [技术架构](#技术架构)
* [模块设计](#模块设计)
* [数据库设计](#数据库设计)
* [API设计](#api设计)
* [安全架构](#安全架构)
* [部署架构](#部署架构)

***

## 系统概述

DMH (Digital Marketing Hub) 是一个企业级数字营销中台系统，提供完整的营销活动管理、多级分销、会员管理和数据分析功能。

### 核心特性

* **极速部署**：1分钟内创建并上线营销活动
* **动态表单**：无需后端改动即可配置表单字段
* **实时奖励**：支付成功后2秒内完成奖励结算
* **多级分销**：支持最多3级分销体系
* **会员系统**：基于UnionID的跨活动会员管理
* **数据同步**：自动同步到客户既有系统（Oracle/SQL Server）

***

## 技术架构

### 容器化部署架构图 ⭐

```
my-net 网络 (172.19.0.0/16)
├── mysql8 (172.19.0.2)              [已存在]
│   └── MySQL 8.0 - 业务数据库
│
├── redis7 (172.19.0.3)              [已存在]
│   └── Redis 7 - 缓存和会话
│
├── dataease-app (172.19.0.4)         [已存在]
│   └── DataEase 应用
│
├── dmh-nginx (172.19.0.5)           [新建] ⭐
│   ├── 端口 3000: 殡理后台
│   ├── 端口 3100: H5前端
│   └── /api/ 代理 → dmh-api:8889
│       └── 托管前端静态文件
│
└── dmh-api (172.19.0.6)             [新建] ⭐
    ├── 端口 8889: 后端API
    ├── DB: mysql8:3306
    └── Redis: redis7:6379
        └── 提供业务API服务
```

***

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         前端层                                │
├─────────────────────────────────────────────────────────────┤
│  管理后台 (Vue3)    │    H5端 (Vue3)    │   小程序 (规划中)  │
│  - 平台管理         │    - 活动浏览      │                    │
│  - 品牌管理         │    - 报名支付      │                    │
│  - 分销商管理       │    - 分销推广      │                    │
│  - 会员管理         │    - 会员中心      │                    │
└─────────────────────────────────────────────────────────────┘
                              ↓ HTTP/HTTPS
┌─────────────────────────────────────────────────────────────┐
│                      API网关层 (go-zero)                      │
├─────────────────────────────────────────────────────────────┤
│  - JWT认证          │  - 权限校验        │  - 限流熔断        │
│  - 请求路由         │  - 日志记录        │  - 监控告警        │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                        业务服务层                             │
├─────────────────────────────────────────────────────────────┤
│  活动管理  │  订单支付  │  奖励系统  │  分销系统  │  会员系统  │
│  品牌管理  │  权限管理  │  安全管理  │  数据同步  │  提现管理  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                        数据存储层                             │
├─────────────────────────────────────────────────────────────┤
│  MySQL 8.0         │  Redis (可选)      │  文件存储          │
│  - 业务数据        │  - 缓存            │  - 海报图片        │
│  - 用户数据        │  - 会话            │  - 品牌素材        │
│  - 订单数据        │  - 队列            │                    │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                      外部系统集成层                           │
├─────────────────────────────────────────────────────────────┤
│  微信支付          │  外部数据库        │  短信/邮件服务      │
│  - 支付回调        │  - Oracle          │  - 验证码          │
│  - 退款            │  - SQL Server      │  - 通知            │
└─────────────────────────────────────────────────────────────┘
```

### 技术栈

#### 后端技术栈

* **语言**: Go 1.23+
* **框架**: go-zero 1.6+ (微服务框架)
* **ORM**: GORM
* **数据库**: MySQL 8.0+
* **缓存**: Redis (可选)
* **认证**: JWT
* **图像处理**: fogleman/gg (海报生成)

#### 前端技术栈

* **框架**: Vue.js 3 (Composition API)
* **构建工具**:
  * Vite 6 (管理后台)
  * Vite 5 (H5端)
* **语言**: TypeScript 5.8+
* **UI组件**:
  * Vant UI (H5端)
  * Lucide Vue (管理后台图标)
* **状态管理**: 计划接入 Pinia
* **路由**: Vue Router

***

## 模块设计

### 1. 营销活动管理模块

**职责**：

* 活动的完整生命周期管理
* 动态表单配置
* 页面设计器
* 活动数据统计

**核心组件**：

```
CampaignService
├── CreateCampaign()      // 创建活动
├── UpdateCampaign()      // 更新活动
├── GetCampaigns()        // 查询活动列表
├── GetCampaignDetail()   // 查询活动详情
└── GetCampaignStats()    // 活动统计
```

**数据表**：

* `campaigns` - 活动主表
* `page_configs` - 页面配置

### 2. 订单与支付模块

**职责**：

* 订单创建与管理
* 支付集成（微信支付）
* 支付回调处理
* 订单状态管理

**核心组件**：

```
OrderService
├── CreateOrder()         // 创建订单
├── GetOrder()            // 查询订单
├── PaymentCallback()     // 支付回调
└── GetOrderList()        // 订单列表
```

**关键流程**：

```
用户提交表单
    ↓
防重复检查 (campaign_id + phone)
    ↓
创建订单 (status: pending)
    ↓
发起支付
    ↓
支付回调
    ↓
[事务开始]
├── 更新订单状态 (paid)
├── 计算奖励
├── 更新余额 (乐观锁)
└── 创建奖励记录
[事务提交]
    ↓
异步同步到外部系统
```

**数据表**：

* `orders` - 订单表
* `external_orders` - 外部订单表

### 3. 实时奖励系统

**职责**：

* 奖励自动计算
* 多级分销奖励
* 余额管理（乐观锁）
* 奖励记录查询

**核心组件**：

```
RewardService
├── CalculateReward()     // 计算奖励
├── SettleReward()        // 结算奖励
├── UpdateBalance()       // 更新余额 (乐观锁)
├── GetBalance()          // 查询余额
└── GetRewardList()       // 奖励列表
```

**乐观锁实现**：

```sql
UPDATE user_balances 
SET 
    balance = balance + ?,
    total_reward = total_reward + ?,
    version = version + 1
WHERE 
    user_id = ? 
    AND version = ?
```

**数据表**：

* `rewards` - 奖励记录
* `user_balances` - 用户余额
* `distributor_rewards` - 分销商奖励

### 4. 分销商系统

**职责**：

* 分销商申请与审批
* 多级分销体系（最多3级）
* 推广链接生成
* 分销数据统计

**核心组件**：

```
DistributorService
├── ApplyDistributor()    // 申请成为分销商
├── ApproveDistributor()  // 审批分销商
├── GenerateLink()        // 生成推广链接
├── GetStatistics()       // 分销统计
└── GetSubordinates()     // 下级分销商
```

**分销层级**：

```
一级分销商 (Level 1)
    ↓ 推荐
二级分销商 (Level 2)
    ↓ 推荐
三级分销商 (Level 3)
```

**数据表**：

* `distributors` - 分销商表
* `distributor_links` - 推广链接
* `distributor_rewards` - 分销奖励

### 5. 会员系统

**职责**：

* 会员档案管理（基于UnionID）
* 会员标签系统
* 会员数据分析
* 会员合并功能

**核心组件**：

```
MemberService
├── CreateMember()        // 创建会员
├── GetMember()           // 查询会员
├── UpdateMember()        // 更新会员
├── MergeMembers()        // 合并会员
└── ExportMembers()       // 导出会员
```

**数据表**：

* `members` - 会员主表
* `member_profiles` - 会员档案
* `member_tags` - 会员标签
* `member_tag_links` - 标签关联
* `member_brand_links` - 品牌关联

### 6. RBAC权限系统

**职责**：

* 用户认证与授权
* 角色权限管理
* 菜单权限管理
* 数据权限控制

**角色定义**：

* `platform_admin` - 平台管理员（全局权限）
* `brand_admin` - 品牌管理员（品牌范围权限）
* `distributor` - 分销商（推广权限）
* `participant` - 普通用户（基础权限）

**核心组件**：

```
AuthService
├── Login()               // 用户登录
├── Register()            // 用户注册
├── CheckPermission()     // 权限检查
└── GetUserMenus()        // 获取用户菜单
```

**数据表**：

* `users` - 用户表
* `roles` - 角色表
* `permissions` - 权限表
* `user_roles` - 用户角色关联
* `role_permissions` - 角色权限关联
* `menus` - 菜单表
* `role_menus` - 角色菜单关联

### 7. 安全管理系统

**职责**：

* 密码策略管理
* 登录尝试监控
* 会话管理
* 审计日志
* 安全事件监控

**核心组件**：

```
SecurityService
├── CheckPasswordPolicy() // 密码策略检查
├── RecordLoginAttempt()  // 记录登录尝试
├── CreateSession()       // 创建会话
├── LogAudit()            // 记录审计日志
└── HandleSecurityEvent() // 处理安全事件
```

**数据表**：

* `password_policies` - 密码策略
* `login_attempts` - 登录尝试
* `user_sessions` - 用户会话
* `audit_logs` - 审计日志
* `security_events` - 安全事件

### 8. 外部数据同步适配器

**职责**：

* 连接外部数据库（Oracle/SQL Server）
* 异步数据同步
* 同步状态监控
* 失败重试机制

**核心组件**：

```
SyncAdapter
├── Connect()             // 连接外部数据库
├── SyncOrder()           // 同步订单
├── SyncReward()          // 同步奖励
├── AsyncSync()           // 异步同步
└── HealthCheck()         // 健康检查
```

**同步流程**：

```
支付成功
    ↓
入队同步任务 (Redis/RabbitMQ)
    ↓
Worker消费任务
    ↓
查询DMH数据
    ↓
字段映射转换
    ↓
执行INSERT到外部数据库
    ↓
更新同步状态
    ↓
失败重试 (最多3次)
```

**数据表**：

* `sync_logs` - 同步日志
* `external_orders` - 外部订单表
* `external_rewards` - 外部奖励表

***

## 数据库设计

### ER图概览

```
users (用户)
  ├─→ user_roles (用户角色)
  ├─→ user_brands (用户品牌)
  ├─→ user_balances (用户余额)
  ├─→ orders (订单)
  ├─→ rewards (奖励)
  └─→ distributors (分销商)

brands (品牌)
  ├─→ campaigns (活动)
  ├─→ brand_assets (品牌素材)
  └─→ distributors (分销商)

campaigns (活动)
  ├─→ orders (订单)
  ├─→ page_configs (页面配置)
  └─→ distributor_links (推广链接)

members (会员)
  ├─→ member_profiles (会员档案)
  ├─→ member_tag_links (标签关联)
  └─→ member_brand_links (品牌关联)
```

### 核心表结构

#### 1. campaigns (营销活动表)

```sql
CREATE TABLE campaigns (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    brand_id BIGINT NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    form_fields JSON,              -- 动态表单配置
    reward_rule DECIMAL(10,2),     -- 奖励规则
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_brand_id (brand_id),
    INDEX idx_status (status)
);
```

#### 2. orders (订单表)

```sql
CREATE TABLE orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    campaign_id BIGINT NOT NULL,
    member_id BIGINT,              -- 会员ID
    phone VARCHAR(20) NOT NULL,
    form_data JSON,                -- 表单数据
    referrer_id BIGINT DEFAULT 0,  -- 推荐人ID
    status VARCHAR(20) DEFAULT 'pending',
    amount DECIMAL(10,2),
    pay_status VARCHAR(20) DEFAULT 'unpaid',
    trade_no VARCHAR(100),
    sync_status VARCHAR(20) DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    UNIQUE KEY uk_campaign_phone (campaign_id, phone, deleted_at),
    INDEX idx_member_id (member_id),
    INDEX idx_referrer_id (referrer_id)
);
```

#### 3. user\_balances (用户余额表)

```sql
CREATE TABLE user_balances (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL UNIQUE,
    balance DECIMAL(10,2) DEFAULT 0.00,
    total_reward DECIMAL(10,2) DEFAULT 0.00,
    version BIGINT DEFAULT 0,      -- 乐观锁版本号
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 4. distributors (分销商表)

```sql
CREATE TABLE distributors (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    brand_id BIGINT NOT NULL,
    level INT DEFAULT 1,           -- 分销商级别 1-3
    parent_id BIGINT,              -- 上级分销商ID
    status VARCHAR(20) DEFAULT 'active',
    total_earnings DECIMAL(10,2) DEFAULT 0.00,
    approved_by BIGINT,
    approved_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_brand (user_id, brand_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_level (level)
);
```

#### 5. members (会员表)

```sql
CREATE TABLE members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    unionid VARCHAR(100) NOT NULL UNIQUE,  -- 微信UnionID
    phone VARCHAR(20),
    nickname VARCHAR(100),
    avatar VARCHAR(255),
    gender VARCHAR(10),
    status VARCHAR(20) DEFAULT 'active',
    first_source VARCHAR(50),      -- 首次来源
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_phone (phone)
);
```

***

## API设计

### RESTful API规范

**基础URL**: `https://api.dmh.com/api/v1`

**认证方式**: JWT Bearer Token

**请求头**:

```
Authorization: Bearer <token>
Content-Type: application/json
```

### API分组

#### 1. 认证相关 `/auth`

```
POST   /auth/register          # 用户注册
POST   /auth/login             # 用户登录
POST   /auth/refresh-token     # 刷新Token
GET    /auth/userinfo          # 获取用户信息
```

#### 2. 活动管理 `/campaigns`

```
POST   /campaigns              # 创建活动
GET    /campaigns              # 活动列表
GET    /campaigns/:id          # 活动详情
PUT    /campaigns/:id          # 更新活动
DELETE /campaigns/:id          # 删除活动
```

#### 3. 订单管理 `/orders`

```
POST   /orders                 # 创建订单
GET    /orders/:id             # 订单详情
POST   /orders/payment/callback # 支付回调
```

#### 4. 奖励管理 `/rewards`

```
GET    /rewards/balance/:userId    # 查询余额
GET    /rewards/:userId            # 奖励列表
```

#### 5. 分销商管理 `/distributor`

```
POST   /distributor/apply          # 申请成为分销商
GET    /distributor/statistics/:brandId  # 分销统计
POST   /distributor/link/generate  # 生成推广链接
GET    /distributor/subordinates/:brandId # 下级分销商
```

#### 6. 会员管理 `/members`

```
GET    /members                # 会员列表
GET    /members/:id            # 会员详情
POST   /members/merge          # 合并会员
POST   /members/export         # 导出会员
```

### 响应格式

**成功响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 业务数据
  }
}
```

**错误响应**:

```json
{
  "code": 40001,
  "message": "您已参与过该活动",
  "data": null
}
```

***

## 安全架构

### 1. 认证与授权

**JWT认证流程**:

```
用户登录
    ↓
验证用户名密码
    ↓
生成JWT Token (有效期24小时)
    ↓
返回Token给客户端
    ↓
客户端请求携带Token
    ↓
服务端验证Token
    ↓
解析用户信息和权限
    ↓
执行业务逻辑
```

**权限检查**:

```go
// 中间件检查权限
func CheckPermission(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userId := c.GetInt64("userId")
        hasPermission := authService.CheckPermission(userId, resource, action)
        if !hasPermission {
            c.JSON(403, gin.H{"error": "无权限"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 2. 数据安全

**敏感数据加密**:

* 密码使用 bcrypt 加密存储
* 外部数据库密码使用 AES 加密
* 支付回调签名验证

**SQL注入防护**:

* 使用预编译语句
* 参数化查询
* ORM框架自动转义

### 3. 业务安全

**防重复提交**:

* 数据库唯一索引约束
* 业务逻辑双重校验
* 幂等性设计

**并发控制**:

* 乐观锁（余额更新）
* 分布式锁（关键业务）
* 事务隔离级别控制

### 4. 监控与审计

**审计日志**:

* 记录所有关键操作
* 包含用户、时间、IP、操作内容
* 日志不可篡改

**安全事件监控**:

* 异常登录检测
* 暴力破解防护
* 异常操作告警

***

## 部署架构

### 开发环境

```
┌─────────────────────────────────────┐
│  开发机器                            │
│  ├── MySQL 8.0 (Docker)             │
│  ├── Backend (go run)               │
│  ├── Frontend-Admin (npm run dev)  │
│  └── Frontend-H5 (npm run dev)     │
└─────────────────────────────────────┘
```

**启动命令**:

```bash
# 初始化数据库
./dmh.sh init

# 启动后端
cd backend && go run api/dmh.go -f api/etc/dmh-api.yaml

# 启动管理后台
cd frontend-admin && npm run dev

# 启动H5端
cd frontend-h5 && npm run dev
```

### 生产环境

```
┌─────────────────────────────────────────────────────────┐
│  Nginx (反向代理 + 静态文件)                             │
│  ├── https://admin.dmh.com → Frontend-Admin (静态)     │
│  ├── https://h5.dmh.com → Frontend-H5 (静态)           │
│  └── https://api.dmh.com → Backend API                 │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│  Backend API (多实例 + 负载均衡)                         │
│  ├── dmh-api-1 (Docker)                                │
│  ├── dmh-api-2 (Docker)                                │
│  └── dmh-api-3 (Docker)                                │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│  MySQL 8.0 (主从复制)                                    │
│  ├── Master (读写)                                      │
│  └── Slave (只读)                                       │
└─────────────────────────────────────────────────────────┘
```

**Docker Compose 部署**:

```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    volumes:
      - mysql-data:/var/lib/mysql
  
  backend:
    build: ./backend
    depends_on:
      - mysql
    environment:
      - DB_HOST=mysql
  
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./frontend-admin/dist:/usr/share/nginx/html/admin
      - ./frontend-h5/dist:/usr/share/nginx/html/h5
```

***

## 性能优化

### 1. 数据库优化

* 合理使用索引
* 查询优化（避免全表扫描）
* 连接池配置
* 读写分离

### 2. 缓存策略

* Redis缓存热点数据
* 活动信息缓存（5分钟）
* 用户会话缓存
* 查询结果缓存

### 3. 异步处理

* 外部数据同步异步化
* 消息队列解耦
* 后台任务处理

### 4. 前端优化

* 代码分割
* 懒加载
* CDN加速
* 图片压缩

***

## 监控与运维

### 监控指标

**系统指标**:

* CPU使用率
* 内存使用率
* 磁盘IO
* 网络流量

**业务指标**:

* API响应时间
* 订单创建量
* 支付成功率
* 奖励结算延迟
* 数据同步成功率

**告警规则**:

* API响应时间 > 1秒
* 支付成功率 < 95%
* 数据同步失败 > 10条
* 系统错误率 > 1%

***

## 扩展性设计

### 水平扩展

* 无状态API设计
* 支持多实例部署
* 负载均衡

### 垂直扩展

* 模块化设计
* 插件化架构
* 配置化驱动

### 未来规划

* 微服务拆分
* 服务网格
* 容器编排（Kubernetes）
* 消息队列（RabbitMQ/Kafka）

***

## 相关文档

* [README.md](./README.md) - 项目介绍
* [SETUP.md](./SETUP.md) - 环境搭建指南
* [API.md](./API.md) - API文档
* [DEVELOPMENT.md](./DEVELOPMENT.md) - 开发指南
* [OpenSpec规格文档](./openspec/specs/) - 详细规格

***

**文档版本**: v1.0\
**最后更新**: 2025-01-21\
**维护者**: DMH Team
