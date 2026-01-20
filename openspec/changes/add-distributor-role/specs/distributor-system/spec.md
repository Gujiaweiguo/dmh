# Spec: 分销商系统（Distributor System）

## ADDED Requirements

### Requirement: 分销商角色定义
系统 SHALL 定义 `distributor` 角色，作为高级顾客角色，具备分享推广并获得奖励的资格。

#### Scenario: 分销商角色创建
- **WHEN** 用户支付订单成功且首次成为分销商
- **THEN** 系统 SHALL 在 `user_roles` 表中为该用户添加 `distributor` 角色关联
- **AND** 在 `distributors` 表中创建分销商专属记录

#### Scenario: 分销商角色验证
- **WHEN** 用户访问分销商专属功能
- **THEN** 系统 SHALL 验证用户是否具有 `distributor` 角色
- **AND** 无该角色时 SHALL 返回 403 禁止访问错误

### Requirement: 自动成为分销商机制
系统 SHALL 在用户支付订单成功后，自动创建分销商记录，无需申请审批。

#### Scenario: 首次支付自动成为分销商
- **WHEN** 用户首次支付订单成功且尚未成为分销商
- **THEN** 系统 SHALL 自动创建分销商记录
- **AND** 设置分销商级别为 1
- **AND** 设置上级分销商为订单的推荐人（referrer_id）
- **AND** 设置分销商状态为 active
- **AND** 记录成为分销商的时间

#### Scenario: 已是分销商跳过创建
- **WHEN** 用户支付订单成功但已是分销商
- **THEN** 系统 SHALL 跳过分销商创建流程
- **AND** 继续执行奖励计算

#### Scenario: 品牌直接推荐用户
- **WHEN** 用户通过品牌海报扫码支付成功
- **THEN** 系统 SHALL 创建分销商记录
- **AND** 设置上级分销商为 0（品牌直接推荐）

#### Scenario: 分销商推荐用户
- **WHEN** 用户通过分销商海报扫码支付成功
- **THEN** 系统 SHALL 创建分销商记录
- **AND** 设置上级分销商为推荐分销商的ID

### Requirement: 活动级别分销奖励规则
系统 SHALL 支持为每个活动单独配置分销奖励规则。

#### Scenario: 创建活动时设置分销规则
- **WHEN** 品牌管理员创建或编辑活动
- **THEN** 系统 SHALL 允许配置分销奖励规则
- **AND** 可设置启用/禁用分销
- **AND** 可设置分销层级（1-3级）
- **AND** 可设置各级奖励比例
- **AND** 将规则保存到 campaigns 表

#### Scenario: 查询活动分销规则
- **WHEN** 系统需要计算分销奖励
- **THEN** 系统 SHALL 查询活动的分销奖励规则
- **AND** 根据规则计算各级奖励金额

#### Scenario: 活动未启用分销
- **WHEN** 活动的 enable_distribution 为 FALSE
- **THEN** 系统 SHALL 不计算分销奖励
- **AND** 也不创建分销商记录

### Requirement: 多级分销体系
系统 SHALL 支持最多3级分销体系，根据活动的分销奖励规则计算奖励。

#### Scenario: 一级分销商奖励
- **WHEN** 新订单支付成功且有直接推荐人（referrer_id）
- **AND** 推荐人是分销商且状态为 active
- **AND** 活动启用分销
- **THEN** 系统 SHALL 按活动的一级奖励比例计算并分配奖励给直接推荐人

#### Scenario: 二级分销商奖励
- **WHEN** 一级分销商获得奖励且其有上级分销商
- **AND** 上级分销商状态为 active
- **AND** 活动分销层级 >= 2
- **THEN** 系统 SHALL 按活动的二级奖励比例计算并分配奖励给二级分销商

#### Scenario: 三级分销商奖励
- **WHEN** 二级分销商获得奖励且其有上级分销商
- **AND** 上级分销商状态为 active
- **AND** 活动分销层级 >= 3
- **THEN** 系统 SHALL 按活动的三级奖励比例计算并分配奖励给三级分销商

#### Scenario: 超过三级不分配奖励
- **WHEN** 推荐链超过三级
- **THEN** 系统 SHALL 只为前三级分销商分配奖励
- **AND** 第四级及之后不获得奖励

#### Scenario: 非分销商不获得推荐奖励
- **WHEN** 推荐人不是分销商或分销商状态非 active
- **THEN** 系统 SHALL 不为该推荐人分配奖励
- **AND** 继续向上查找是否有分销商推荐人

#### Scenario: 基于分销链计算奖励
- **WHEN** 订单的 distributor_path 包含分销链
- **THEN** 系统 SHALL 解析分销链
- **AND** 为分销链中的各级分销商分配奖励

### Requirement: 分销商级别管理
系统 SHALL 支持品牌管理员或平台管理员调整分销商级别。

#### Scenario: 手动升级分销商级别
- **WHEN** 品牌管理员或平台管理员手动调整分销商级别
- **THEN** 系统 SHALL 更新 `distributors` 表中的 level 字段
- **AND** 记录级别变更日志

#### Scenario: 分销商级别变更影响后续奖励
- **WHEN** 分销商级别发生变更
- **THEN** 新级别 SHALL 适用于级别变更后的新订单
- **AND** 历史奖励保持不变

#### Scenario: 品牌管理员只能管理本品牌分销商
- **WHEN** 品牌管理员调整分销商级别
- **THEN** 系统 SHALL 只允许管理本品牌的分销商
- **AND** 不能操作其他品牌的分销商

### Requirement: 分销商状态管理
系统 SHALL 支持分销商的激活和暂停状态。

#### Scenario: 暂停分销商资格
- **WHEN** 品牌管理员或平台管理员暂停分销商
- **THEN** 系统 SHALL 更新分销商状态为 suspended
- **AND** 暂停后的分销商不再获得新的推荐奖励
- **AND** 历史奖励保持不变

#### Scenario: 重新激活分销商
- **WHEN** 管理员重新激活被暂停的分销商
- **THEN** 系统 SHALL 更新分销商状态为 active
- **AND** 恢复后的分销商可继续获得推荐奖励

#### Scenario: 暂停状态不获得奖励
- **WHEN** 暂停状态的分销商有下级订单支付成功
- **THEN** 系统 SHALL 不为该分销商分配奖励
- **AND** 继续向上查找上级分销商

### Requirement: 二维码海报生成
系统 SHALL 为分销商提供活动专属海报和通用分销商海报的生成功能。

#### Scenario: 生成活动专属海报
- **WHEN** 分销商请求生成活动专属海报
- **THEN** 系统 SHALL 生成包含活动信息和分销商二维码的海报
- **AND** 海报URL格式为 `{poster_url}?campaignId={id}&distributorId={id}`
- **AND** 记录海报生成时间

#### Scenario: 生成通用分销商海报
- **WHEN** 分销商请求生成通用分销商海报
- **THEN** 系统 SHALL 生成包含分销商信息和所有活动二维码的海报
- **AND** 海报URL格式为 `{poster_url}?distributorId={id}`
- **AND** 记录海报生成时间

#### Scenario: 查看海报
- **WHEN** 用户访问海报URL
- **THEN** 系统 SHALL 返回海报图片
- **AND** 自动识别海报类型（活动专属或通用）

#### Scenario: 推广链接访问追踪
- **WHEN** 用户通过分销商海报扫码访问活动
- **THEN** 系统 SHALL 记录访问来源为该分销商
- **AND** 在用户最终下单时将分销商设为推荐人

### Requirement: 分销商数据查看
系统 SHALL 允许分销商查看自己的推广数据和奖励明细。

#### Scenario: 查看推广数据统计
- **WHEN** 分销商访问推广数据页面
- **THEN** 系统 SHALL 返回以下统计数据：
  - 累计订单数
  - 累计奖励金额
  - 可提现金额
  - 下级分销商数量（仅一级）
  - 本月/本周新增订单

#### Scenario: 查看奖励明细
- **WHEN** 分销商访问奖励明细页面
- **THEN** 系统 SHALL 返回分页的奖励记录列表
- **AND** 每条记录包含：订单ID、奖励金额、奖励时间、来源订单

#### Scenario: 查看下级分销商列表
- **WHEN** 分销商访问下级列表页面
- **THEN** 系统 SHALL 返回其直接下级分销商列表（仅一级）
- **AND** 每个下级显示：姓名、级别、加入时间、累计订单数

#### Scenario: 数据隔离
- **WHEN** 分销商查看数据
- **THEN** 系统 SHALL 只返回该分销商自己的数据
- **AND** 不能查看其他分销商的数据

### Requirement: 提现功能
系统 SHALL 允许分销商申请提现，需要平台管理员审批。

#### Scenario: 分销商申请提现
- **WHEN** 分销商提交提现申请
- **THEN** 系统 SHALL 验证提现金额不大于可提现金额
- **AND** 创建提现记录，状态为 pending
- **AND** 记录提现方式、账号、真实姓名
- **AND** 通知平台管理员有待审批申请

#### Scenario: 提现金额校验
- **WHEN** 分销商申请的提现金额大于可提现金额
- **THEN** 系统 SHALL 拒绝提现申请
- **AND** 提示余额不足

#### Scenario: 平台管理员审批通过提现
- **WHEN** 平台管理员批准提现申请
- **THEN** 系统 SHALL 更新提现状态为 approved
- **AND** 记录审批人和审批时间
- **AND** 调用支付接口打款
- **AND** 打款成功后更新状态为 completed
- **AND** 扣除分销商余额
- **AND** 通知分销商提现成功

#### Scenario: 平台管理员拒绝提现
- **WHEN** 平台管理员拒绝提现申请
- **THEN** 系统 SHALL 更新提现状态为 rejected
- **AND** 记录拒绝原因
- **AND** 通知分销商提现被拒绝

#### Scenario: 查看提现记录
- **WHEN** 分销商访问提现记录页面
- **THEN** 系统 SHALL 返回该分销商的提现记录列表
- **AND** 显示每笔提现的状态、金额、时间、原因

#### Scenario: 数据隔离
- **WHEN** 分销商查看提现记录
- **THEN** 系统 SHALL 只返回该分销商自己的提现记录
- **AND** 不能查看其他分销商的提现记录

### Requirement: 分销商与品牌关联
系统 SHALL 支持分销商与特定品牌的关联关系。

#### Scenario: 自动关联品牌
- **WHEN** 用户支付订单成功自动成为分销商
- **THEN** 系统 SHALL 将分销商关联到订单所属的品牌

#### Scenario: 多品牌分销商
- **WHEN** 用户在不同品牌的活动中都支付过
- **THEN** 系统 SHALL 为每个品牌创建独立的分销商记录
- **AND** 奖励按品牌分别计算

#### Scenario: 品牌管理员只能管理本品牌分销商
- **WHEN** 品牌管理员管理分销商
- **THEN** 系统 SHALL 只允许管理本品牌的分销商
- **AND** 不能操作其他品牌的分销商

### Requirement: 品牌管理员数据查看
系统 SHALL 允许品牌管理员查看本品牌的分销商、顾客和奖励数据。

#### Scenario: 查看本品牌分销商列表
- **WHEN** 品牌管理员访问分销商管理页面
- **THEN** 系统 SHALL 返回本品牌的分销商列表
- **AND** 显示每个分销商的基本信息、级别、状态

#### Scenario: 查看本品牌顾客列表
- **WHEN** 品牌管理员访问顾客列表页面
- **THEN** 系统 SHALL 返回本品牌的顾客列表
- **AND** 显示每个顾客的基本信息、订单数量

#### Scenario: 查看本品牌奖励详情
- **WHEN** 品牌管理员访问奖励详情页面
- **THEN** 系统 SHALL 返回本品牌的奖励记录列表
- **AND** 显示奖励的详细信息、统计数据

#### Scenario: 数据隔离
- **WHEN** 品牌管理员查看数据
- **THEN** 系统 SHALL 只返回本品牌的数据
- **AND** 不能查看其他品牌的数据

### Requirement: 平台管理员全局数据查看
系统 SHALL 允许平台管理员查看全局的分销商、奖励和提现数据，并按品牌筛选。

#### Scenario: 查看全部分销商
- **WHEN** 平台管理员访问全局分销商页面
- **THEN** 系统 SHALL 返回所有品牌的分销商列表
- **AND** 支持按品牌ID筛选
- **AND** 显示分销商的基本信息、级别、状态

#### Scenario: 查看全部奖励
- **WHEN** 平台管理员访问全局奖励页面
- **THEN** 系统 SHALL 返回所有品牌的奖励记录列表
- **AND** 支持按品牌ID筛选
- **AND** 显示奖励的详细信息

#### Scenario: 查看全部提现
- **WHEN** 平台管理员访问全局提现页面
- **THEN** 系统 SHALL 返回所有品牌的提现记录列表
- **AND** 支持按品牌ID和状态筛选
- **AND** 显示提现的详细信息

#### Scenario: 筛选数据
- **WHEN** 平台管理员使用筛选条件查询数据
- **THEN** 系统 SHALL 根据筛选条件返回对应的数据
- **AND** 支持多条件组合筛选
