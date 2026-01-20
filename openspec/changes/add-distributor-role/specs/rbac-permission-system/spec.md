# Spec: RBAC权限管理系统 - 分销商角色扩展

## ADDED Requirements

### Requirement: 分销商角色定义
系统 SHALL 新增 `distributor` 角色作为 `participant` 的高级扩展角色。

#### Scenario: 分销商角色权限定义
- **WHEN** 用户被分配 `distributor` 角色
- **THEN** 用户 SHALL 保留所有 `participant` 角色权限
- **AND** 额外获得分销商专属权限

### Requirement: 分销商权限范围
系统 SHALL 定义分销商的具体权限范围。

#### Scenario: 分销商基础权限
- **WHEN** 用户具有 `distributor` 角色
- **THEN** 系统 SHALL 允许以下操作：
  - 参与/支付活动（继承自 participant）
  - 查看个人订单和奖励（继承自 participant）
  - 申请提现（继承自 participant）
  - 自动成为分销商
  - 生成专属二维码海报
  - 查看推广数据统计
  - 查看奖励明细
  - 查看一级下级分销商列表
  - 申请提现

#### Scenario: 分销商禁止操作
- **WHEN** 用户仅具有 `distributor` 角色
- **THEN** 系统 SHALL 禁止以下操作：
  - 创建/编辑/删除活动
  - 管理品牌信息
  - 管理其他分销商
  - 审批提现
  - 查看平台级数据
  - 访问管理后台

#### Scenario: 分销商数据隔离
- **WHEN** 分销商查询数据
- **THEN** 系统 SHALL 只返回该分销商相关的数据
- **AND** 不能查看其他分销商的数据

### Requirement: 分销商提现权限
系统 SHALL 允许分销商申请提现，但需要平台管理员审批。

#### Scenario: 申请提现权限
- **WHEN** 用户具有 `distributor` 角色
- **THEN** 系统 SHALL 允许用户申请提现
- **AND** 提现申请需要平台管理员审批

#### Scenario: 查看自己提现记录权限
- **WHEN** 用户具有 `distributor` 角色
- **THEN** 系统 SHALL 允许用户查看自己的提现记录
- **AND** 不能查看其他用户的提现记录

## MODIFIED Requirements

### Requirement: 角色权限体系
系统 SHALL 实现5种用户角色（原4种），每种角色具有明确的权限范围和功能边界。

#### Scenario: 平台管理员权限
- **WHEN** 平台管理员(platform_admin)访问任何功能
- **THEN** 系统SHALL允许访问所有功能模块
- **AND** 可以管理所有品牌、活动、用户、分销商和系统配置
- **AND** 可以管理任何品牌的分销商（调整级别、暂停/激活）
- **AND** 可以审批任何品牌的分销商提现申请
- **AND** 可以查看全局分销商、奖励、提现数据（可按品牌筛选）

#### Scenario: 品牌管理员权限
- **WHEN** 品牌管理员(brand_admin)访问品牌相关功能
- **THEN** 系统SHALL只允许访问其管理的品牌数据
- **AND** 可以管理品牌信息(编辑品牌资料、上传品牌logo等)
- **AND** 可以管理品牌素材库(上传、分类、删除素材)
- **AND** 可以创建、编辑、删除、发布本品牌的活动
- **AND** 可以为活动配置分销奖励规则
- **AND** 可以查看本品牌的数据统计和报表
- **AND** 可以管理本品牌的分销商（调整级别、暂停/激活）
- **AND** 可以查看本品牌的分销商、顾客、奖励详情
- **AND** 不能审批提现申请

#### Scenario: 分销商权限（新增）
- **WHEN** 分销商(distributor)访问系统功能
- **THEN** 系统SHALL允许访问所有participant功能
- **AND** 可以自动成为分销商
- **AND** 可以生成专属二维码海报（活动专属 + 通用）
- **AND** 可以查看推广数据统计和奖励明细
- **AND** 可以查看一级下级分销商列表
- **AND** 可以申请提现
- **AND** 不能访问管理后台功能
- **AND** 不能管理其他分销商

#### Scenario: 活动参与者权限
- **WHEN** 活动参与者(participant)访问系统功能
- **THEN** 系统SHALL只允许访问个人相关功能
- **AND** 可以参与活动、查看个人奖励和申请提现
- **AND** 可以自动成为分销商

#### Scenario: 匿名用户权限
- **WHEN** 匿名用户(anonymous)访问系统
- **THEN** 系统SHALL只允许访问公开功能
- **AND** 可以浏览活动列表、查看活动详情和注册账号

### Requirement: 用户注册权限控制
系统 SHALL 根据用户角色类型实现不同的注册方式和权限控制。

#### Scenario: H5注册限制角色
- **WHEN** 用户通过H5页面注册
- **THEN** 系统 SHALL 只允许创建 participant 角色的用户
- **AND** 不允许通过H5注册创建管理员角色或分销商角色

#### Scenario: 分销商角色获得方式
- **WHEN** 用户想获得分销商角色
- **THEN** 系统 SHALL 在用户支付订单成功后自动分配分销商角色
- **AND** 不需要申请审批流程

#### Scenario: 品牌管理员角色分配
- **WHEN** 需要创建品牌管理员用户
- **THEN** 系统 SHALL 要求平台管理员通过后台管理系统操作
- **AND** 同时在 brand_admins 表中建立品牌关联关系

#### Scenario: 平台管理员创建限制
- **WHEN** 需要创建平台管理员用户
- **THEN** 系统 SHALL 只允许现有平台管理员通过后台系统创建
- **AND** 平台管理员角色不能通过任何前端注册方式获得

#### Scenario: 匿名用户转换
- **WHEN** 匿名用户完成H5注册流程
- **THEN** 系统 SHALL 将其转换为 participant 角色
- **AND** 获得相应的登录凭据和基础权限

#### Scenario: 参与者自动转换为分销商
- **WHEN** participant 支付订单成功
- **THEN** 系统 SHALL 自动为其添加 distributor 角色
- **AND** 创建分销商专属记录

### Requirement: API权限控制
系统 SHALL 在 API 层面实现细粒度的权限控制，确保每个接口都有适当的权限检查。

#### Scenario: 分销商API权限检查
- **WHEN** 用户访问分销商专属API
- **THEN** 系统 SHALL 验证用户是否具有 distributor 角色
- **AND** 无该角色时返回 403 错误

#### Scenario: 提现API权限检查
- **WHEN** 用户访问提现API
- **THEN** 系统 SHALL 验证用户是否具有 distributor 角色
- **AND** 申请提现需要 distributor 角色
- **AND** 审批提现需要 platform_admin 角色
- **AND** 无对应角色时返回 403 错误

#### Scenario: 数据级权限隔离
- **WHEN** 品牌管理员查询活动数据
- **THEN** 系统 SHALL 只返回其管理品牌的活动数据
- **AND** 不能访问其他品牌的数据

#### Scenario: 分销商数据隔离
- **WHEN** 分销商查询推广数据
- **THEN** 系统 SHALL 只返回该分销商自己的数据
- **AND** 不能查看其他分销商的数据

#### Scenario: 平台管理员数据访问
- **WHEN** 平台管理员查询全局数据
- **THEN** 系统 SHALL 返回所有数据
- **AND** 支持按品牌筛选

### Requirement: 分销商管理权限
系统 SHALL 实现品牌管理员和平台管理员对分销商的管理权限。

#### Scenario: 品牌管理员管理分销商
- **WHEN** 品牌管理员管理本品牌分销商
- **THEN** 系统 SHALL 允许调整分销商级别
- **AND** 允许暂停/激活分销商
- **AND** 不能管理其他品牌的分销商

#### Scenario: 平台管理员管理分销商
- **WHEN** 平台管理员管理任何品牌分销商
- **THEN** 系统 SHALL 允许调整任何品牌的分销商级别
- **AND** 允许暂停/激活任何品牌的分销商
- **AND** 可以访问所有品牌的数据

#### Scenario: 提现审批权限
- **WHEN** 审批提现申请
- **THEN** 系统 SHALL 只允许平台管理员审批
- **AND** 品牌管理员不能审批提现

### Requirement: 品牌数据查看权限
系统 SHALL 实现品牌管理员和平台管理员对品牌数据的查看权限。

#### Scenario: 品牌管理员查看本品牌数据
- **WHEN** 品牌管理员查看数据
- **THEN** 系统 SHALL 只返回本品牌的数据
- **AND** 包括：分销商列表、顾客列表、奖励详情

#### Scenario: 平台管理员查看全局数据
- **WHEN** 平台管理员查看数据
- **THEN** 系统 SHALL 返回所有品牌的数据
- **AND** 支持按品牌筛选
- **AND** 包括：分销商列表、奖励列表、提现列表
