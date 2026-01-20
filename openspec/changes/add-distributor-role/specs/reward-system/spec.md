# Spec: 实时奖励系统 - 多级分销扩展

## MODIFIED Requirements

### Requirement: 奖励计算
系统 SHALL 支持多级分销奖励计算，根据活动的分销奖励规则为最多3级分销商分配奖励。

#### Scenario: 查询活动分销奖励规则
- **WHEN** 系统计算分销奖励
- **THEN** 系统 SHALL 查询活动的分销奖励规则（campaigns.distribution_rewards）
- **AND** 检查活动是否启用分销（campaigns.enable_distribution）
- **AND** 检查活动的分销层级（campaigns.distribution_level）

#### Scenario: 单级分销奖励计算
- **WHEN** 活动启用分销且分销层级为 1
- **AND** 订单支付成功且有直接推荐人
- **THEN** 系统 SHALL 按活动的 level1 奖励比例计算奖励
- **AND** 为直接推荐人创建奖励记录
- **AND** 更新直接推荐人余额

#### Scenario: 二级分销奖励计算
- **WHEN** 活动启用分销且分销层级 >= 2
- **AND** 订单支付成功且推荐链包含二级分销商
- **THEN** 系统 SHALL 按活动的 level1 奖励比例为一级分销商计算奖励
- **AND** 按活动的 level2 奖励比例为二级分销商计算奖励
- **AND** 为两级分销商分别创建奖励记录
- **AND** 更新两级分销商的余额

#### Scenario: 三级分销奖励计算
- **WHEN** 活动启用分销且分销层级 = 3
- **AND** 订单支付成功且推荐链包含三级分销商
- **THEN** 系统 SHALL 按活动的 level1、level2、level3 奖励比例计算奖励
- **AND** 为三级分销商分别创建奖励记录
- **AND** 更新三级分销商的余额

#### Scenario: 超过层级限制不分配奖励
- **WHEN** 推荐链长度超过活动的分销层级
- **THEN** 系统 SHALL 只为前 N 级分销商分配奖励（N = distribution_level）
- **AND** 超过层级的分销商不获得奖励

#### Scenario: 活动未启用分销不计算奖励
- **WHEN** 活动的 enable_distribution 为 FALSE
- **THEN** 系统 SHALL 不计算任何分销奖励
- **AND** 保持原有单级奖励逻辑（如有）

#### Scenario: 基于分销链计算奖励
- **WHEN** 订单的 distributor_path 包含分销链
- **THEN** 系统 SHALL 解析 distributor_path 字段
- **AND** 按顺序获取各级分销商ID
- **AND** 为分销链中的分销商计算奖励

#### Scenario: 非分销商不获得奖励
- **WHEN** 推荐链中的用户不是分销商或分销商状态非 active
- **THEN** 系统 SHALL 跳过该用户
- **AND** 继续向上查找下一级分销商

#### Scenario: 暂停状态分销商不获得奖励
- **WHEN** 分销商状态为 suspended
- **THEN** 系统 SHALL 不为该分销商分配奖励
- **AND** 继续向上查找上级分销商

#### Scenario: 奖励金额计算
- **WHEN** 系统为某级分销商计算奖励
- **THEN** 系统 SHALL 使用公式：奖励金额 = 订单金额 × (奖励比例 / 100)
- **AND** 保留两位小数
- **AND** 记录奖励级别信息

### Requirement: 奖励记录
系统 SHALL 在奖励记录中保存分销商级别信息。

#### Scenario: 记录分销商级别
- **WHEN** 创建分销奖励记录
- **THEN** 系统 SHALL 记录分销商级别（level 1/2/3）
- **AND** 记录订单的 distributor_path

#### Scenario: 多级奖励记录
- **WHEN** 系统为多级分销商创建奖励记录
- **THEN** 系统 SHALL 为每级分销商创建一条独立的奖励记录
- **AND** 每条记录关联同一订单ID
- **AND** 每条记录标记对应的分销商级别

## ADDED Requirements

### Requirement: 订单分销链记录
系统 SHALL 在订单中记录完整的分销链路径。

#### Scenario: 记录分销链路径
- **WHEN** 用户通过分销商海报下单
- **THEN** 系统 SHALL 查询分销链（最多3级）
- **AND** 将分销链ID序列存储到 orders.distributor_path
- **AND** 格式为："一级ID,二级ID,三级ID"

#### Scenario: 品牌直接推荐用户
- **WHEN** 用户通过品牌海报下单
- **THEN** 系统 SHALL 设置 distributor_path 为空
- **AND** 该用户支付后成为分销商时 parent_id = 0

#### Scenario: 查询分销链
- **WHEN** 系统需要计算奖励
- **THEN** 系统 SHALL 解析订单的 distributor_path
- **AND** 获取各级分销商ID
- **AND** 验证各级分销商状态和级别

### Requirement: 余额管理扩展
系统 SHALL 扩展余额管理以支持多级奖励并发更新。

#### Scenario: 并发更新多级分销商余额
- **WHEN** 系统同时为多个分销商更新余额
- **THEN** 系统 SHALL 使用乐观锁确保并发安全
- **AND** 每个分销商的余额更新独立处理
- **AND** 任何一个更新失败不影响其他更新

#### Scenario: 余额更新失败重试
- **WHEN** 某个分销商余额更新因版本冲突失败
- **THEN** 系统 SHALL 重试最多3次
- **AND** 重试后仍失败则记录错误日志
- **AND** 其他分销商的余额更新继续执行
