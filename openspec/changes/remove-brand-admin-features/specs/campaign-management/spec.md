# 活动管理权限规范变更

## REMOVED Requirements

### Requirement: 品牌管理员活动管理权限
**Reason**: 移除品牌管理员在管理后台的活动管理权限
**Migration**: 品牌管理员通过H5端管理活动，admin通过管理后台查看和控制活动

#### Scenario: 品牌管理员管理后台活动操作
- **WHEN** 品牌管理员尝试在管理后台操作活动
- **THEN** 系统不再提供此访问权限

## MODIFIED Requirements

### Requirement: 活动查看和控制权限
系统SHALL只允许platform_admin在管理后台查看所有活动详情并进行控制操作。

#### Scenario: 平台管理员查看活动
- **WHEN** platform_admin访问活动管理页面
- **THEN** 可以查看所有品牌的所有活动详细信息

#### Scenario: 平台管理员控制活动
- **WHEN** platform_admin对活动进行操作
- **THEN** 可以停止、启用或删除任何活动

#### Scenario: 活动数据统计查看
- **WHEN** platform_admin查看活动详情
- **THEN** 可以看到完整的参与人数、转化率等统计数据

### Requirement: 活动创建权限简化
系统SHALL允许platform_admin创建活动，不再区分品牌归属权限。

#### Scenario: 创建活动
- **WHEN** platform_admin创建新活动
- **THEN** 可以为任何品牌创建活动，不受品牌管理员关系限制