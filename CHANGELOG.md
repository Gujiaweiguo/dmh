# 📝 变更日志

本文档记录 DMH 项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

***

## \[1.2.1] - 2026-02-11

### 修复

* 🐛 前端构建稳定性
  * frontend-admin 构建压缩从 `terser` 切换为 `esbuild`
  * 解决 Vite 可选依赖缺失导致的构建失败（`terser not found`）

### 测试

* ✅ 测试覆盖补齐（后端 + 前端）
  * 新增 backend handler/logic/middleware/common/model/config/types/svc/api/cmd 测试
  * 新增 frontend-admin `performanceMonitor` 单测
  * 新增 frontend-h5 `api` 单测
  * 后端 `go test ./... -run TestDoesNotExist` 下 `no test files` 归零

### 验证

* ✅ 本地回归全绿
  * backend: `go test ./...`
  * frontend-admin: `npm run test && npm run build`
  * frontend-h5: `npm run test && npm run build`

* ✅ 远端回归全绿
  * `Order MySQL8 Regression` 手动触发成功

***

## \[1.2.0] - 2026-02-10

### 新增

* ✨ 反馈管理系统
  * 用户反馈提交（支持 bug、feature、other 分类）
  * 反馈列表查看（分页、筛选、搜索）
  * 反馈详情查看
  * 反馈状态管理（pending、reviewing、resolved、closed）
  * 反馈统计面板（总数、解决率、平均评分、平均解决时长）
  * 满意度调查功能
  * FAQ 管理功能
  * 功能使用记录
  * 品牌管理员和平台管理员权限控制

* ✨ 完整的测试覆盖
  * 后端单元测试（feedback handler 和 logic）
  * 后端集成测试（feedback CRUD 操作）
  * 前端单元测试（admin 和 h5）
  * 前端 E2E 测试（反馈管理完整流程）
  * 分销商逻辑测试套件
  * CI guard workflow（feedback 核心测试）

### 修复

* 🐛 反馈系统数据库 schema
  * 修复 user_feedback 表结构与 GORM 模型对齐
  * 移除无效的索引定义
  * 修复 TableName() 返回值

* 🐛 JWT 认证上下文
  * 硬化 feedback handler 的 JWT 解析
  * 添加 getUserIDFromRequest/getUserRoleFromRequest 辅助函数
  * 改进错误处理和回退机制

* 🐛 E2E 测试配置
  * 修复 Playwright baseURL 配置（3000 → 正确端口）
  * 修正登录表单选择器（placeholder vs name）
  * 修复反馈管理元素选择器（.table tbody tr）
  * 使用系统 Chrome 浏览器绕过下载限制

### 变更

* 🔄 数据库初始化
  * 反馈系统 6 张表合并到 init.sql
  * 提供默认反馈数据

* 🔄 部署脚本
  * 更新本地重启脚本（支持反馈系统）
  * 更新迁移脚本说明

### 技术栈

* **后端**: Go 1.23 + go-zero 1.6 + GORM + MySQL 8.0
* **前端**: Vue 3 + Vite 5/6 + TypeScript
* **测试**: Playwright 1.58.2 + Vitest

### 文档

* 📚 更新 OpenSpec 归档
* 📚 添加测试执行报告和计划文档
* 📚 更新 CHANGELOG.md（本条目）

***

## \[1.1.0] - 2026-01-30

### 新增

* ✨ 容器化部署支持（dmh-nginx + dmh-api）✨
  * 独立 Nginx 容器托管前端静态文件
  * 独立 API 容器托管后端服务
  * 加入 my-net 网络（172.19.0.0/16）
  * 固定 IP：dmh-nginx (172.19.0.5), dmh-api (172.19.0.6)
* ✨ 部署脚本和配置文件
  * 快速启动脚本（quick-start.sh）
  * 完整部署脚本（deploy-containers.sh）
  * 回滚脚本（rollback-containers.sh）
  * Docker Compose 编排文件
* ✨ 部署文档
  * deploy/README.md - 容器化部署指南
  * 统一的部署管理文档
* ✨ Nginx 配置优化
  * API 代理到后端容器
  * 静态资源缓存（1年）
  * Gzip 压缩支持

### 变更

* 🔄 部署方式从独立进程改为容器化
* 🔄 文档结构优化，移除临时文档
* 🔄 前端构建产物路径适配容器化

### 移除

* ❌ AGENTS.md - 与 openspec/AGENTS.md 重复
* ❌ API.md - 与 docs/API\_Documentation.md 重复
* ❌ SCRIPTS.md - 已废弃（改用容器化）
* ❌ STARTUP.md - 已废弃（改用容器化）
* ❌ API\_TEST\_AND\_ISSUES\_ANALYSIS.md - 临时测试报告
* ❌ Deployment\_Checklist.md - 被 deploy/README.md 替代
* ❌ JWT\_FIX\_AND\_API\_TEST\_REPORT.md - 临时测试报告
* ❌ P0\_COMPLETION\_REPORT.md - 临时任务报告
* ❌ P0\_TEST\_REPORT.md - 临时测试报告
* ❌ deployment\_summary\_report.md - 临时报告
* ❌ dmh-container-deployment-report.md - 临时报告
* ❌ docker\_migration\_guide.md - 已整合到 deploy/README.md
* ❌ docker\_quick\_reference.md - 已整合到 deploy/README.md

### 优化

* ⚡ 文档结构优化（文档数量减少53%）
* ⚡ 部署方式更现代化、更易维护
* ⚡ 服务隔离，故障影响更小
* ⚡ 易于横向扩展

### 技术栈

* **后端**: Go 1.23 + go-zero 1.6 + GORM + MySQL 8.0
* **前端**: Vue 3 + Vite 5/6 + TypeScript
* **容器化**: Docker + Docker Compose
* **反向代理**: Nginx 1.25
* **认证**: JWT
* **部署**: 容器化部署（Docker Compose）

### 文档

* 📚 完整的项目文档
* 📚 部署文档：deploy/README.md（容器化）
* 📚 架构文档：ARCHITECTURE.md
* 📚 API 文档：docs/API\_Documentation.md
* 📚 用户手册：docs/User\_Manual.md
* 📚 开发指南：DEVELOPMENT.md
* 📚 更新日志：CHANGELOG.md（本文件）

***

## \[Unreleased]

### 计划中

* \[ ] 微信小程序端
* \[ ] 数据分析看板
* \[ ] 消息推送系统
* \[ ] 优惠券系统
* \[ ] 积分系统

***

## \[1.0.0] - 2025-01-21

### 新增

* ✨ 完整的 RBAC 权限系统
  * 4 种用户角色（平台管理员、品牌管理员、分销商、普通用户）
  * JWT 认证和 Token 刷新
  * 动态菜单权限
  * 数据权限控制

* ✨ 营销活动管理
  * 活动创建和编辑
  * 动态表单配置
  * 页面设计器（8 种组件）
  * 活动数据统计

* ✨ 多级分销系统
  * 分销商申请和审批
  * 3 级分销体系
  * 推广链接生成
  * 分销数据统计
  * 实时奖励结算

* ✨ 会员管理系统
  * 基于 UnionID 的会员档案
  * 会员标签系统
  * 会员数据分析
  * 会员合并功能
  * 会员数据导出

* ✨ 订单与支付
  * 订单创建和管理
  * 微信支付集成
  * 支付回调处理
  * 防重复提交

* ✨ 安全管理
  * 密码策略管理
  * 登录尝试监控
  * 会话管理
  * 审计日志
  * 安全事件监控

* ✨ 数据同步适配器
  * 连接外部数据库（Oracle/SQL Server）
  * 异步数据同步
  * 同步状态监控
  * 失败重试机制

* ✨ 提现管理
  * 提现申请
  * 提现审批
  * 提现记录查询

### 技术栈

* **后端**: Go 1.23 + go-zero 1.6 + GORM + MySQL 8.0
* **前端**: Vue 3 + Vite 5/6 + TypeScript
* **认证**: JWT
* **部署**: Docker + Nginx

### 文档

* 📚 完整的项目文档
  * README.md - 项目介绍
  * ARCHITECTURE.md - 系统架构
  * API.md - API 文档
  * DEVELOPMENT.md - 开发指南
  * CHANGELOG.md - 变更日志（本文件）

***

## \[0.9.0] - 2025-01-20

### 新增

* 会员系统基础功能
* 会员标签管理
* 会员数据导出

### 优化

* 优化数据库查询性能
* 改进前端用户体验
* 完善错误处理

### 修复

* 修复分销商审批流程问题
* 修复订单重复提交问题
* 修复权限检查逻辑

***

## \[0.8.0] - 2025-01-19

### 新增

* 分销商系统
* 推广链接生成
* 分销数据统计

### 优化

* 优化活动页面设计器
* 改进表单验证逻辑

***

## \[0.7.0] - 2025-01-18

### 新增

* H5 页面设计器
* 8 种页面组件
* 主题配置功能

### 优化

* 优化移动端适配
* 改进组件拖拽体验

***

## \[0.6.0] - 2025-01-15

### 新增

* 安全管理系统
* 密码策略配置
* 审计日志功能

### 修复

* 修复登录会话过期问题
* 修复权限检查漏洞

***

## \[0.5.0] - 2025-01-10

### 新增

* 数据同步适配器
* 外部数据库连接
* 异步同步机制

### 优化

* 优化同步性能
* 改进错误重试逻辑

***

## \[0.4.0] - 2025-01-05

### 新增

* 订单管理功能
* 微信支付集成
* 支付回调处理

### 修复

* 修复订单状态更新问题
* 修复支付金额计算错误

***

## \[0.3.0] - 2025-01-01

### 新增

* 活动管理功能
* 动态表单配置
* 活动数据统计

### 优化

* 优化活动列表查询
* 改进表单验证

***

## \[0.2.0] - 2024-12-25

### 新增

* RBAC 权限系统
* 用户角色管理
* 菜单权限配置

### 优化

* 优化权限检查性能
* 改进菜单加载逻辑

***

## \[0.1.0] - 2024-12-20

### 新增

* 项目初始化
* 基础架构搭建
* 用户认证功能
* 品牌管理功能

***

## 版本说明

### 版本号格式

版本号格式：`主版本号.次版本号.修订号`

* **主版本号**：重大架构变更或不兼容的 API 修改
* **次版本号**：新增功能，向下兼容
* **修订号**：Bug 修复，向下兼容

### 变更类型

* **新增** (Added): 新功能
* **修改** (Changed): 现有功能的变更
* **弃用** (Deprecated): 即将移除的功能
* **移除** (Removed): 已移除的功能
* **修复** (Fixed): Bug 修复
* **安全** (Security): 安全相关的修复
* **优化** (Optimized): 性能优化

***

## 相关链接

* [项目主页](https://github.com/Gujiaweiguo/DMH)
* [问题反馈](https://github.com/Gujiaweiguo/DMH/issues)
* [贡献指南](./CONTRIBUTING.md)

***

**维护者**: DMH Team\
**最后更新**: 2025-01-21
