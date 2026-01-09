# 管理后台权限规范变更

## REMOVED Requirements

### Requirement: 品牌管理员后台访问权限
**Reason**: 简化权限模型，移除品牌管理员在管理后台的访问权限
**Migration**: 品牌管理员将只能通过H5端进行活动管理

#### Scenario: 品牌管理员登录管理后台
- **WHEN** 品牌管理员尝试登录管理后台
- **THEN** 系统应拒绝登录并提示权限不足

### Requirement: 品牌管理员关系管理
**Reason**: 移除品牌管理员功能，不再需要管理品牌管理员与品牌的关系
**Migration**: 相关数据将被清理，功能由admin直接管理品牌替代

#### Scenario: 管理品牌管理员关系
- **WHEN** admin尝试配置品牌管理员关系
- **THEN** 系统不再提供此功能

## MODIFIED Requirements

### Requirement: 管理后台用户角色管理
系统SHALL只支持platform_admin角色在管理后台的访问，移除brand_admin角色支持。

#### Scenario: 用户角色创建
- **WHEN** admin创建新用户
- **THEN** 只能选择platform_admin或participant角色，不能选择brand_admin

#### Scenario: 现有品牌管理员用户处理
- **WHEN** 系统中存在brand_admin角色的用户
- **THEN** 这些用户将无法访问管理后台，需要通过H5端进行操作

### Requirement: 品牌管理权限
系统SHALL只允许platform_admin角色管理品牌信息和活动。

#### Scenario: 品牌信息管理
- **WHEN** platform_admin访问品牌管理功能
- **THEN** 可以查看、创建、编辑、删除所有品牌信息

#### Scenario: 活动管理权限
- **WHEN** platform_admin访问活动管理功能  
- **THEN** 可以查看所有活动详情并进行停止/启用操作