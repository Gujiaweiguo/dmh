# Project Context

## Current Status

**版本**: 1.3.0-dev
**最后更新**: 2026-02-19

### 最近完成（2026-02-19）

* H5 品牌设置页真实接口化（`/brands/:id` + 本地持久化策略）
* Admin 个人中心真实接口化（新增 `profileApi` 服务）
* 微信支付服务可配置 Mock/实调（`MockEnabled` 开关）
* 前端 lint 门禁打通（ESLint 9 + Makefile check 强制）

## Purpose

DMH (Digital Marketing Hub) 数字营销中台系统，提供营销活动管理、会员管理、分销商管理、奖励结算等功能。

### 核心目标

* **极速部署**：1分钟内完成营销活动上线
* **动态表单**：支持灵活的数据采集，无需后端改动
* **实时激励**：支付即结算，显著提升推广员积极性
* **无缝集成**：内置外网数据库适配器

## Tech Stack

### Backend

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.24.0 | 主要语言 |
| go-zero | 1.6.0 | REST 框架 |
| GORM | 1.25.5 | ORM |
| MySQL | 8.0+ | 主数据库 |
| Redis | 7+ | 缓存 |
| golang-jwt | 4.5.0 | 认证 |
| testify | 1.10.0 | 测试框架 |

### Frontend Admin

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | ^3.4.0 | 框架 |
| TypeScript | ~5.8.2 | 语言 |
| Vite | ^6.2.0 | 构建 |
| Tailwind CSS | ^3.4.19 | 样式 |
| Lucide Vue Next | ^0.263.0 | 图标 |
| Vitest | ^2.1.8 | 单元测试 |
| Playwright | ^1.58.2 | E2E 测试 |

### Frontend H5

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | ^3.5.27 | 框架 |
| Vite | ^5.4.21 | 构建 |
| Vant | ^4.9.22 | 移动端 UI |
| Vue Router | ^4.6.4 | 路由 |
| Axios | ^1.13.4 | HTTP |
| Vitest | ^2.1.8 | 单元测试 |
| Playwright | ^1.58.2 | E2E 测试 |

## Project Structure

```
DMH/
├── backend/              # Go API (端口: 8889)
│   ├── api/              # 入口、API定义、配置
│   │   ├── dmh.go        # 主入口
│   │   ├── dmh.api       # go-zero API 规范
│   │   └── internal/     # handler, logic, middleware, svc, types
│   ├── model/            # GORM 模型
│   ├── common/           # 共享工具
│   ├── migrations/       # SQL 迁移
│   └── test/             # 集成/性能测试
├── frontend-admin/       # 管理后台 (端口: 3000)
│   ├── index.tsx         # 入口 (非标准)
│   ├── views/            # 页面组件
│   └── services/         # API 服务
├── frontend-h5/          # H5 前端 (端口: 3100)
│   └── src/              # 标准 Vue 结构
├── openspec/             # 规格驱动开发
│   ├── specs/            # 已实现的能力
│   └── changes/          # 变更提案
└── deploy/              # Docker Compose 部署
```

## Project Conventions

### Code Style

* Go: `gofmt` 标准格式化
* TypeScript: ESLint + Prettier
* 提交信息: Conventional Commits (`feat:`, `fix:`, `refactor:`)
* 代码缩进: 2 spaces (前端), tabs (Go)

### Architecture Patterns

#### Backend (go-zero)
* **Handler**: 薄层，仅解析请求、调用 Logic、返回响应
* **Logic**: 业务逻辑层，所有业务代码在此
* **Model**: GORM 数据模型
* **Middleware**: 认证、授权、限流、CORS

#### Frontend Admin (非标准结构)
* 入口是 `index.tsx` (非 `main.js`)
* 扁平结构，无 `src/` 目录
* 混合 `.tsx` 和 `.vue` 文件

#### Frontend H5 (标准结构)
* 标准 Vue 3 结构 (`src/views/`, `src/components/`)
* 每个 view 有 companion `.logic.js` 文件

### Testing Strategy

| 类型 | 后端 | 前端 |
|------|------|------|
| 单元测试 | `go test ./...` | Vitest |
| 集成测试 | `test/integration/` | - |
| E2E 测试 | - | Playwright |
| 覆盖率目标 | 70%+ | 70%+ |

### Git Workflow

* 主分支: `main`
* 提交格式: Conventional Commits
  * `feat:` 新功能
  * `fix:` Bug 修复
  * `refactor:` 重构
  * `test:` 测试
  * `docs:` 文档
  * `chore:` 杂项

## Domain Context

### 核心模块

| 模块 | 说明 |
|------|------|
| Campaign | 营销活动管理 |
| Order | 订单管理 |
| Reward | 奖励结算 |
| Distributor | 分销商管理 |
| Member | 会员管理 |
| Brand | 品牌管理 |
| Feedback | 反馈系统 |
| RBAC | 权限管理 |

### 核心业务规则

1. **实时奖励结算**
   * 触发时机：支付回调成功
   * 结算速度：推荐人余额在2秒内自动更新
   * 防重复：同一活动/同一手机号限报一次

2. **动态表单**
   * 管理员可自定义表单字段
   * 无需后端代码修改即可上线新表单

3. **分销商体系**
   * 多级分销支持
   * 实时佣金结算
   * 推广链接/二维码生成

## Important Constraints

### 技术约束

* 后端 API 端口: 8889
* 管理后台端口: 3000
* H5 前端端口: 3100
* 数据库: MySQL 8.0+ (容器名: mysql8)
* 缓存: Redis 7+ (容器名: redis-dmh)

### 部署约束

* Docker Compose 部署 (网络: my-net)
* 后端二进制需在宿主机编译后复制到容器
* 容器间通过服务名通信 (mysql8, redis-dmh)

## Current Specs

| Capability | Status |
|------------|--------|
| campaign-management | ✅ Active |
| order-payment-system | ✅ Active |
| rbac-permission-system | ✅ Active |
| feedback-system | ✅ Active |
| spec-governance | ✅ Active |
| system-test-execution | ✅ Active |
| test-coverage-improvement | ✅ Active |
| brand-settings-real-api | ✅ Active |
| admin-profile-real-api | ✅ Active |
| wechat-pay-mock-switch | ✅ Active |
| frontend-lint-gate | ✅ Active |

## External Dependencies

* **MySQL 8.0**: 主数据存储
* **Redis 7**: 缓存、会话、队列
* **WeChat Pay**: 支付集成 (`backend/common/wechatpay/`)
* **Poster Generation**: 海报生成 (`backend/common/poster/`)
