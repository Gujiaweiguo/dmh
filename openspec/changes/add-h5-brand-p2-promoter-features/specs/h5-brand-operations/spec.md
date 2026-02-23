## ADDED Requirements

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
