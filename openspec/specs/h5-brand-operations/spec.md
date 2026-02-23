# h5-brand-operations Specification

## Purpose
TBD - created by archiving change update-h5-brand-p1-features. Update Purpose after archive.
## Requirements
### Requirement: Campaign editor uses dynamic brand context
H5 品牌端活动编辑页面 SHALL 从当前登录用户的品牌上下文动态获取 `brandId`，不得使用硬编码品牌 ID。

#### Scenario: Create campaign with dynamic brand id
- **WHEN** 品牌管理员在 H5 端创建或编辑活动
- **THEN** 系统 SHALL 使用当前登录用户关联的品牌 ID 作为请求参数
- **AND** 请求体中 SHALL 不出现固定写死的 `brandId` 值

#### Scenario: Missing brand context handling
- **WHEN** 当前登录态无法解析出品牌 ID
- **THEN** 系统 SHALL 阻止提交活动配置
- **AND** 显示明确错误提示引导用户重新登录或联系管理员

### Requirement: Brand settings supports logo upload
H5 品牌设置页面 SHALL 支持品牌 Logo 上传并持久化保存。

#### Scenario: Upload and save logo
- **WHEN** 品牌管理员在设置页上传 Logo 并保存
- **THEN** 系统 SHALL 调用上传接口并获得可访问的 Logo URL
- **AND** 保存品牌设置时 SHALL 持久化该 Logo URL

#### Scenario: Upload failure feedback
- **WHEN** Logo 上传失败（网络异常或格式不合法）
- **THEN** 系统 SHALL 提示上传失败原因
- **AND** 保留当前页面已编辑但未保存的数据

### Requirement: Orders page supports export
H5 品牌端订单页面 SHALL 提供订单导出能力。

#### Scenario: Export filtered order list
- **WHEN** 品牌管理员在订单页设置筛选条件并触发导出
- **THEN** 系统 SHALL 以当前筛选条件请求导出接口
- **AND** 返回可下载文件（如 CSV/Excel）

#### Scenario: Export progress and error handling
- **WHEN** 导出请求处理中或失败
- **THEN** 系统 SHALL 提供加载状态反馈
- **AND** 导出失败时 SHALL 显示错误信息并允许重试

### Requirement: Orders page supports order detail view
H5 品牌端订单页面 SHALL 支持查看订单详情。

#### Scenario: Open order detail from list
- **WHEN** 品牌管理员在订单列表点击某条订单
- **THEN** 系统 SHALL 打开对应订单详情页面或详情面板
- **AND** 展示订单基础信息、支付状态与核销状态

#### Scenario: Handle order detail fetch failure
- **WHEN** 订单详情接口返回失败或订单不存在
- **THEN** 系统 SHALL 显示错误提示
- **AND** 提供返回列表或重试入口

### Requirement: Promoter detail page shows complete profile
H5 品牌端推广员详情页面 SHALL 展示推广员完整信息，包括基础信息、业绩统计、推广链接。

#### Scenario: View promoter detail from list
- **WHEN** 品牌管理员在推广员列表点击某条推广员的"详情"
- **THEN** 系统 SHALL 打开对应推广员详情页面
- **AND** 展示推广员基础信息（姓名/手机/状态/注册时间）
- **AND** 展示业绩统计（推广订单数/奖励金额/推广活动数）
- **AND** 展示推广链接列表

#### Scenario: Handle promoter not found
- **WHEN** 访问不存在的推广员详情
- **THEN** 系统 SHALL 显示错误提示
- **AND** 提供返回列表入口

### Requirement: Promoter detail page links to reward records
H5 品牌端推广员详情页面 SHALL 提供奖励记录查看入口。

#### Scenario: Navigate to reward records from detail
- **WHEN** 品牌管理员在推广员详情页点击"查看奖励记录"
- **THEN** 系统 SHALL 跳转到该推广员的奖励记录页面
- **AND** 筛选条件自动设置为当前推广员

### Requirement: Reward records page shows payment history
H5 品牌端奖励记录页面 SHALL 展示奖励发放历史。

#### Scenario: View reward records list
- **WHEN** 品牌管理员打开奖励记录页面
- **THEN** 系统 SHALL 展示奖励记录列表
- **AND** 每条记录包含订单号、奖励金额、状态、发放时间

#### Scenario: Filter reward records by promoter
- **WHEN** 从推广员详情页进入奖励记录
- **THEN** 系统 SHALL 自动筛选该推广员的奖励记录
- **AND** 可清除筛选查看全部记录

#### Scenario: Handle empty reward records
- **WHEN** 推广员暂无奖励记录
- **THEN** 系统 SHALL 显示空状态提示
- **AND** 提示"暂无奖励记录"

