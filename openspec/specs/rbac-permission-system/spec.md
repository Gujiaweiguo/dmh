# rbac-permission-system Specification

## Purpose
TBD - created by archiving change enhance-rbac-permission-system. Update Purpose after archive.
## Requirements
### Requirement: 用户认证管理
系统SHALL提供完整的用户认证功能，包括注册、登录、密码管理和会话控制。

#### Scenario: H5用户注册成功
- **WHEN** 用户通过H5页面提供有效的用户名、密码和手机号
- **THEN** 系统SHALL创建新用户账号并分配默认角色(participant)
- **AND** 系统SHALL返回JWT token和用户信息

#### Scenario: 平台管理员后台创建用户
- **WHEN** 平台管理员通过后台管理系统创建用户
- **THEN** 系统SHALL创建用户账号并分配指定角色
- **AND** 平台管理员角色只能通过后台系统创建，不允许H5注册

#### Scenario: 用户登录成功
- **WHEN** 用户提供正确的用户名和密码
- **THEN** 系统SHALL验证用户身份并生成JWT token
- **AND** token SHALL包含用户ID、角色信息和有效期

#### Scenario: 密码安全验证
- **WHEN** 用户设置或修改密码
- **THEN** 系统SHALL使用bcrypt加密存储密码
- **AND** 密码SHALL满足最小安全要求(长度≥6位)

#### Scenario: 会话超时控制
- **WHEN** JWT token超过有效期(24小时)
- **THEN** 系统SHALL拒绝请求并返回401未授权错误
- **AND** 前端SHALL清除本地token并跳转到登录页

### Requirement: 角色权限体系
系统SHALL实现4种用户角色，每种角色具有明确的权限范围和功能边界。

#### Scenario: 平台管理员权限
- **WHEN** 平台管理员(platform_admin)访问任何功能
- **THEN** 系统SHALL允许访问所有功能模块
- **AND** 可以管理所有品牌、活动、用户和系统配置

#### Scenario: 品牌管理员权限
- **WHEN** 品牌管理员(brand_admin)访问品牌相关功能
- **THEN** 系统SHALL只允许访问其管理的品牌数据
- **AND** 可以管理品牌信息(编辑品牌资料、上传品牌logo等)
- **AND** 可以管理品牌素材库(上传、分类、删除素材)
- **AND** 可以创建、编辑、删除、发布本品牌的活动
- **AND** 可以查看本品牌的数据统计和报表

#### Scenario: 活动参与者权限
- **WHEN** 活动参与者(participant)访问系统功能
- **THEN** 系统SHALL只允许访问个人相关功能
- **AND** 可以参与活动、查看个人奖励和申请提现

#### Scenario: 匿名用户权限
- **WHEN** 匿名用户(anonymous)访问系统
- **THEN** 系统SHALL只允许访问公开功能
- **AND** 可以浏览活动列表、查看活动详情和注册账号

### Requirement: API权限控制
系统SHALL在API层面实现细粒度的权限控制，确保每个接口都有适当的权限检查。

#### Scenario: JWT token验证
- **WHEN** 客户端请求需要认证的API接口
- **THEN** 系统SHALL验证Authorization header中的JWT token
- **AND** token无效时SHALL返回401未授权错误

#### Scenario: 权限检查机制
- **WHEN** 用户访问受保护的API接口
- **THEN** 系统SHALL根据URL和HTTP方法确定所需权限
- **AND** 检查用户角色是否具有该权限

#### Scenario: 数据级权限隔离
- **WHEN** 品牌管理员查询活动数据
- **THEN** 系统SHALL只返回其管理品牌的活动数据
- **AND** 不能访问其他品牌的数据

### Requirement: 用户注册权限控制
系统SHALL根据用户角色类型实现不同的注册方式和权限控制。

#### Scenario: H5注册限制角色
- **WHEN** 用户通过H5页面注册
- **THEN** 系统SHALL只允许创建participant角色的用户
- **AND** 不允许通过H5注册创建管理员角色

#### Scenario: 品牌管理员角色分配
- **WHEN** 需要创建品牌管理员用户
- **THEN** 系统SHALL要求平台管理员通过后台管理系统操作
- **AND** 同时在brand_admins表中建立品牌关联关系

#### Scenario: 平台管理员创建限制
- **WHEN** 需要创建平台管理员用户
- **THEN** 系统SHALL只允许现有平台管理员通过后台系统创建
- **AND** 平台管理员角色不能通过任何前端注册方式获得

#### Scenario: 匿名用户转换
- **WHEN** 匿名用户完成H5注册流程
- **THEN** 系统SHALL将其转换为participant角色
- **AND** 获得相应的登录凭据和基础权限

### Requirement: 用户管理功能
系统SHALL提供完整的用户管理功能，支持用户创建、查询、更新、状态管理和密码重置。

#### Scenario: 后台创建用户账号
- **WHEN** 平台管理员通过后台管理系统创建新用户
- **THEN** 系统SHALL验证用户名和手机号的唯一性
- **AND** 创建用户记录并分配指定角色
- **AND** 生成初始密码并通知用户

#### Scenario: 用户状态管理
- **WHEN** 平台管理员变更用户状态
- **THEN** 系统SHALL更新用户状态(active/disabled/locked)
- **AND** 禁用用户SHALL立即失去系统访问权限
- **AND** 记录状态变更日志和操作人

#### Scenario: 用户密码重置
- **WHEN** 平台管理员重置用户密码
- **THEN** 系统SHALL生成新的临时密码
- **AND** 强制用户下次登录时修改密码
- **AND** 记录密码重置操作日志

#### Scenario: 用户角色分配
- **WHEN** 平台管理员为用户分配角色
- **THEN** 系统SHALL更新用户角色关联关系
- **AND** 新角色权限SHALL立即生效
- **AND** 记录角色变更日志

#### Scenario: 品牌管理员分配
- **WHEN** 平台管理员指定用户为品牌管理员
- **THEN** 系统SHALL在brand_admins表中创建关联记录
- **AND** 用户SHALL获得该品牌的管理权限
- **AND** 可以同时管理多个品牌

### Requirement: 权限缓存优化
系统SHALL实现权限信息缓存机制，提高权限检查的性能和响应速度。

#### Scenario: 权限信息缓存
- **WHEN** 系统首次查询用户权限信息
- **THEN** 系统SHALL将权限信息缓存到内存中
- **AND** 后续权限检查SHALL优先使用缓存数据

#### Scenario: 缓存失效更新
- **WHEN** 用户角色或权限发生变更
- **THEN** 系统SHALL立即清除相关缓存
- **AND** 下次权限检查SHALL重新查询数据库

### Requirement: 品牌管理员关系管理
系统SHALL为平台管理员提供品牌管理员与品牌关系的完整管理功能，支持绑定、解绑和变更操作。

#### Scenario: 绑定品牌管理员
- **WHEN** 平台管理员为用户绑定品牌管理权限
- **THEN** 系统SHALL在brand_admins表中创建关联记录
- **AND** 用户SHALL立即获得该品牌的管理权限
- **AND** 记录绑定操作日志

#### Scenario: 解绑品牌管理员
- **WHEN** 平台管理员解除用户的品牌管理权限
- **THEN** 系统SHALL删除brand_admins表中的关联记录
- **AND** 用户SHALL立即失去该品牌的管理权限
- **AND** 记录解绑操作日志

#### Scenario: 变更品牌管理员权限
- **WHEN** 平台管理员调整品牌管理员的品牌范围
- **THEN** 系统SHALL更新brand_admins表中的关联记录
- **AND** 新的品牌权限SHALL立即生效
- **AND** 记录权限变更日志

#### Scenario: 多品牌管理支持
- **WHEN** 品牌管理员被分配多个品牌
- **THEN** 系统SHALL支持一个用户管理多个品牌
- **AND** 在数据查询时SHALL正确过滤各品牌数据
- **AND** 权限检查SHALL验证用户对特定品牌的访问权限

#### Scenario: 品牌管理员权限查询
- **WHEN** 查询用户的品牌管理权限
- **THEN** 系统SHALL返回用户管理的所有品牌列表
- **AND** 包含品牌基本信息和权限范围
- **AND** 支持按品牌状态过滤
系统SHALL为品牌管理员提供完整的品牌管理功能，包括品牌信息、素材库、活动和数据管理。

#### Scenario: 品牌信息管理
- **WHEN** 品牌管理员编辑品牌信息
- **THEN** 系统SHALL允许修改品牌名称、描述、logo等基本信息
- **AND** 只能修改其管理的品牌信息

#### Scenario: 品牌素材库管理
- **WHEN** 品牌管理员管理素材库
- **THEN** 系统SHALL允许上传、分类、编辑、删除品牌素材
- **AND** 素材包括图片、视频、文档等多种类型
- **AND** 只能管理本品牌的素材资源

#### Scenario: 品牌活动管理
- **WHEN** 品牌管理员管理活动
- **THEN** 系统SHALL允许创建、编辑、删除、发布本品牌的活动
- **AND** 可以配置活动的动态表单和奖励规则
- **AND** 不能访问其他品牌的活动

#### Scenario: 品牌数据查看
- **WHEN** 品牌管理员查看数据统计
- **THEN** 系统SHALL提供本品牌的完整数据报表
- **AND** 包括活动参与数、订单统计、奖励发放、用户分析等
- **AND** 不能查看其他品牌的数据
系统SHALL实现提现申请和审核的权限控制，确保只有授权用户可以进行相关操作。

#### Scenario: 提现申请权限
- **WHEN** 用户申请提现
- **THEN** 系统SHALL检查用户是否为已认证用户
- **AND** 验证用户余额是否足够

#### Scenario: 提现审核权限
- **WHEN** 用户尝试审核提现申请
- **THEN** 系统SHALL验证用户是否为平台管理员
- **AND** 只有平台管理员可以批准或拒绝提现

#### Scenario: 提现状态更新
- **WHEN** 平台管理员审核提现申请
- **THEN** 系统SHALL使用数据库事务确保数据一致性
- **AND** 记录审核人和审核时间

### Requirement: 品牌管理功能
系统SHALL实现提现申请和审核的权限控制，确保只有授权用户可以进行相关操作。

#### Scenario: 提现申请权限
- **WHEN** 用户申请提现
- **THEN** 系统SHALL检查用户是否为已认证用户
- **AND** 验证用户余额是否足够

#### Scenario: 提现审核权限
- **WHEN** 用户尝试审核提现申请
- **THEN** 系统SHALL验证用户是否为平台管理员
- **AND** 只有平台管理员可以批准或拒绝提现

#### Scenario: 提现状态更新
- **WHEN** 平台管理员审核提现申请
- **THEN** 系统SHALL使用数据库事务确保数据一致性
- **AND** 记录审核人和审核时间

### Requirement: 菜单权限管理
系统SHALL提供完整的菜单权限管理功能，支持菜单结构管理和角色菜单权限配置。

#### Scenario: 菜单结构管理
- **WHEN** 平台管理员管理菜单结构
- **THEN** 系统SHALL支持菜单的增加、删除、修改和排序
- **AND** 支持多级菜单结构和菜单分组
- **AND** 区分后台管理菜单和H5用户菜单

#### Scenario: 页面操作权限配置
- **WHEN** 配置页面操作权限
- **THEN** 系统SHALL支持增删改查、导出、转发等操作权限
- **AND** 每个菜单项可配置多种操作权限
- **AND** 支持按钮级别的权限控制

#### Scenario: 角色菜单权限分配
- **WHEN** 为角色分配菜单权限
- **THEN** 系统SHALL支持选择性分配菜单访问权限
- **AND** 支持为每个菜单配置具体的操作权限
- **AND** 权限变更SHALL立即生效

#### Scenario: 用户菜单权限查询
- **WHEN** 用户登录系统
- **THEN** 系统SHALL根据用户角色返回可访问的菜单列表
- **AND** 包含每个菜单的操作权限信息
- **AND** 前端根据权限动态显示菜单和按钮

#### Scenario: 菜单权限继承
- **WHEN** 配置多级菜单权限
- **THEN** 系统SHALL支持权限继承机制
- **AND** 子菜单可继承父菜单的权限设置
- **AND** 支持覆盖继承的权限配置

### Requirement: 安全审计日志
系统SHALL记录所有重要的安全相关操作，提供完整的审计追踪能力。

#### Scenario: 用户操作日志
- **WHEN** 用户执行重要操作(登录、权限变更、数据修改)
- **THEN** 系统SHALL记录操作日志
- **AND** 日志SHALL包含用户ID、操作类型、时间戳和IP地址

#### Scenario: 权限变更日志
- **WHEN** 管理员修改用户角色或权限
- **THEN** 系统SHALL记录权限变更日志
- **AND** 日志SHALL包含变更前后的权限状态

#### Scenario: 安全事件监控
- **WHEN** 检测到异常登录或权限滥用
- **THEN** 系统SHALL记录安全事件
- **AND** 可选择性地触发安全告警

---

### Requirement: Distributor monitoring view reliability
系统 SHALL 为平台管理员提供稳定可用的分销监控页面，确保加载失败和超时场景可恢复。

#### Scenario: Request timeout handling on distributor monitoring page
- **WHEN** 平台管理员进入分销监控页面且分销数据请求超时
- **THEN** 系统 SHALL 显示明确的超时提示
- **AND** 提供可用的重试入口重新发起请求

### Requirement: Distributor monitoring interaction performance baseline
系统 SHALL 保证分销监控页面在筛选与搜索交互中具备可接受的响应性能。

#### Scenario: Filter and search responsiveness
- **WHEN** 平台管理员在分销监控页面调整状态筛选、等级筛选或输入搜索关键词
- **THEN** 系统 SHALL 在可接受时间内更新结果列表
- **AND** 搜索输入 SHALL 通过防抖避免重复请求造成卡顿

