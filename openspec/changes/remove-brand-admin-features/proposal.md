# Change: 移除管理后台品牌管理员功能，简化权限模型

## Why
当前系统中品牌管理员功能在管理后台中造成了权限复杂性，且与H5端的品牌管理功能重复。为了简化系统架构和用户体验，需要移除管理后台中的品牌管理员相关功能，只保留平台管理员(admin)的核心管理功能。

## What Changes
- **BREAKING**: 移除管理后台中品牌管理员角色的相关功能
- 保留admin查看和管理品牌信息的功能
- 保留admin查看活动详情和停止活动的功能
- 移除品牌管理员在管理后台的登录和操作权限
- 简化权限模型，只保留platform_admin在管理后台的访问权限

## Impact
- affected specs: admin-management, user-auth, campaign-management
- affected code: 
  - frontend-admin/index.tsx (权限检查和UI组件)
  - backend权限验证逻辑
  - 数据库权限配置
- 用户影响: 品牌管理员将只能通过H5端进行活动管理，无法访问管理后台