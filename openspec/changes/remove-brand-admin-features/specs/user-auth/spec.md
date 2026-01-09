# 用户认证权限规范变更

## REMOVED Requirements

### Requirement: 品牌管理员认证支持
**Reason**: 移除品牌管理员在管理后台的登录支持
**Migration**: 品牌管理员需要通过H5端进行登录和操作

#### Scenario: 品牌管理员管理后台登录
- **WHEN** 用户使用brand_admin角色凭据登录管理后台
- **THEN** 系统应拒绝登录并返回权限不足错误

## MODIFIED Requirements

### Requirement: 管理后台登录权限验证
系统SHALL只允许platform_admin角色用户登录管理后台。

#### Scenario: 平台管理员登录
- **WHEN** platform_admin角色用户提供正确凭据
- **THEN** 系统允许登录并提供完整的管理后台访问权限

#### Scenario: 非平台管理员登录尝试
- **WHEN** 非platform_admin角色用户尝试登录管理后台
- **THEN** 系统拒绝登录并提示"仅限平台管理员访问"

### Requirement: 用户角色验证中间件
系统SHALL在所有管理后台API请求中验证用户角色为platform_admin。

#### Scenario: API权限检查
- **WHEN** 用户访问管理后台API端点
- **THEN** 系统验证JWT token中的角色必须为platform_admin

#### Scenario: 无效角色访问
- **WHEN** 非platform_admin角色用户访问管理后台API
- **THEN** 系统返回403 Forbidden错误