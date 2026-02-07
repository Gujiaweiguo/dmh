# Project Context

## Purpose

DMH (Digital Marketing Hub) 是一个轻量级、高性能的数字营销中台系统，旨在提供快速的营销活动管理和实时奖励结算功能。

### 核心目标

* **极速部署**：1分钟内完成营销活动上线
* **动态表单**：支持灵活的数据采集，无需后端改动
* **实时激励**：支付即结算，显著提升推广员积极性
* **无缝集成**：内置外网数据库适配器

## Tech Stack

### Frontend

* **框架**: Vue 3.3.4 (Composition API)
* **构建工具**: Vite 6.2.0
* **语言**: TypeScript 5.8.2
* **样式**: Tailwind CSS
* **图标**: Lucide Vue Next 0.263.0
* **AI集成**: Google Generative AI 1.34.0

### Backend (模拟)

* **架构参考**: Go-Zero 微服务架构
* **核心逻辑**:
  * `api.createOrder`: 封装事务处理、奖励计算、同步任务入队
  * `SyncAdapter`: 模拟异步外网直连驱动

## Project Conventions

### Code Style

* 使用 TypeScript 严格模式
* Vue 3 Composition API 编写组件
* 文件命名：PascalCase for components (e.g., `CampaignEditorView.tsx`)
* 变量命名：camelCase
* 常量命名：UPPER\_SNAKE\_CASE
* 代码缩进：2 spaces

### Architecture Patterns

* **前端架构**:
  * Views: 页面级组件（如 DashboardView, CampaignListView）
  * Components: 可复用组件（如 Sidebar, MobilePreview）
  * Services: API 调用和业务逻辑（如 mockApi.ts）
  * Types: TypeScript 类型定义集中管理
* **状态管理**: 计划接入 Vuex/Pinia (Beta 1.1)
* **移动端**: Uni-app 风格 H5 容器，未来支持多端（Beta 1.2）

### Testing Strategy

* MVP 阶段以功能验证为主
* 重点测试：
  * 营销活动创建和管理流程
  * 实时奖励结算逻辑
  * 数据同步功能
  * 移动端预览功能

### Git Workflow

* 主分支: `main`
* 功能开发: feature branches
* 提交信息格式: `[类型] 简短描述`
  * 类型: feat, fix, docs, style, refactor, test, chore

## Domain Context

### 数字营销业务知识

* **营销活动**: 包含活动配置、动态表单、奖励规则
* **推广员**: 通过推荐链接推广活动，获得实时奖励
* **奖励结算**: 支付成功后2秒内自动结算到推荐人账户
* **数据同步**: 支持将数据同步到外部数据库（客户自有系统）

### 核心业务规则

1. **实时奖励结算**:
   * 触发时机：支付回调成功
   * 结算速度：推荐人余额在2秒内自动更新
   * 防重复：同一活动/同一手机号限报一次

2. **动态表单**:
   * 管理员可自定义表单字段
   * 无需后端代码修改即可上线新表单

## Important Constraints

### 技术约束

* MVP 版本，功能以快速验证为主
* 前端项目位于 `/opt/code/DMH/frontend/` 目录
* 使用模拟 API (mockApi.ts) 而非真实后端

### 业务约束

* 极速部署要求：活动上线时间 < 1分钟
* 实时性要求：奖励结算延迟 < 2秒
* 安全性：同一用户同一活动只能参与一次

## External Dependencies

### AI 服务

* **Google Generative AI (Gemini)**: 用于智能功能增强
* API Key 配置：需在 `.env.local` 中设置 `GEMINI_API_KEY`

### 外部系统集成（计划中）

* 外网数据库适配器：支持客户自有数据库直连
* 支付回调系统：接收支付成功通知
* 推送服务：实时通知推广员奖励到账

## MVP 演进路线

* **Beta 1.1**: 接入 Vuex/Pinia 状态管理
* **Beta 1.2**: Uni-app 分配至多端（微信小程序/字节小程序）
