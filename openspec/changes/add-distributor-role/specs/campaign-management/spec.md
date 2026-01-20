# Spec: 营销活动管理模块 - 分销奖励规则配置

## ADDED Requirements

### Requirement: 活动分销奖励规则配置
系统 SHALL 支持在活动创建和编辑时配置分销奖励规则。

#### Scenario: 活动创建时配置分销规则
- **WHEN** 品牌管理员创建活动
- **THEN** 系统 SHALL 允许配置以下分销规则：
  - enable_distribution: 是否启用分销（布尔值）
  - distribution_level: 分销层级（1/2/3级）
  - distribution_rewards: 各级奖励比例（JSON格式）
- **AND** 将规则保存到 campaigns 表

#### Scenario: 活动编辑时修改分销规则
- **WHEN** 品牌管理员编辑活动
- **THEN** 系统 SHALL 允许修改分销规则
- **AND** 更新 campaigns 表对应字段
- **AND** 新规则适用于新订单
- **AND** 历史订单不受影响

#### Scenario: 分销规则格式验证
- **WHEN** 品牌管理员提交分销规则
- **THEN** 系统 SHALL 验证规则格式：
  - distribution_level 必须是 1、2 或 3
  - distribution_rewards 必须包含对应数量的奖励比例
  - 每个奖励比例必须是 0-100 之间的数字
- **AND** 验证失败时提示用户

#### Scenario: 分销规则默认值
- **WHEN** 品牌管理员未配置分销规则
- **THEN** 系统 SHALL 使用默认值：
  - enable_distribution: FALSE
  - distribution_level: 1
  - distribution_rewards: {"level1": 10}

#### Scenario: 禁用分销的活动
- **WHEN** 活动的 enable_distribution 为 FALSE
- **THEN** 系统 SHALL 不计算分销奖励
- **AND** 用户支付后不会成为分销商
- **AND** 保持原有单级奖励逻辑

#### Scenario: 启用分销的活动
- **WHEN** 活动的 enable_distribution 为 TRUE
- **THEN** 系统 SHALL 根据分销规则计算奖励
- **AND** 用户支付后自动成为分销商
- **AND** 按分销层级和奖励比例分配奖励

### Requirement: 分销奖励规则展示
系统 SHALL 在活动详情和管理页面展示分销奖励规则。

#### Scenario: 活动列表显示分销状态
- **WHEN** 品牌管理员查看活动列表
- **THEN** 系统 SHALL 在列表中显示活动是否启用分销
- **AND** 显示分销层级

#### Scenario: 活动详情显示分销规则
- **WHEN** 品牌管理员查看活动详情
- **THEN** 系统 SHALL 展示完整的分销奖励规则：
  - 是否启用分销
  - 分销层级
  - 各级奖励比例
- **AND** 格式化显示，便于阅读

#### Scenario: 活动统计包含分销数据
- **WHEN** 品牌管理员查看活动统计
- **THEN** 系统 SHALL 在统计数据中包含：
  - 分销商数量
  - 分销奖励总额
  - 各级分销商数量
  - 各级奖励金额

### Requirement: 分销规则模板
系统 SHALL 提供常用的分销规则模板，方便品牌管理员快速配置。

#### Scenario: 一级分销模板
- **WHEN** 品牌管理员选择一级分销模板
- **THEN** 系统 SHALL 自动填充：
  - distribution_level: 1
  - distribution_rewards: {"level1": 10}

#### Scenario: 二级分销模板
- **WHEN** 品牌管理员选择二级分销模板
- **THEN** 系统 SHALL 自动填充：
  - distribution_level: 2
  - distribution_rewards: {"level1": 10, "level2": 5}

#### Scenario: 三级分销模板
- **WHEN** 品牌管理员选择三级分销模板
- **THEN** 系统 SHALL 自动填充：
  - distribution_level: 3
  - distribution_rewards: {"level1": 10, "level2": 5, "level3": 3}

#### Scenario: 自定义模板
- **WHEN** 品牌管理员创建自定义模板
- **THEN** 系统 SHALL 允许保存模板
- **AND** 模板包含名称和分销规则
- **AND** 品牌管理员可以重用自定义模板

## MODIFIED Requirements

### Requirement: 活动创建/编辑
系统 SHALL 扩展活动创建/编辑功能，增加分销奖励规则配置。

#### Scenario: 创建活动完整流程（扩展）
- **WHEN** 品牌管理员创建活动
- **THEN** 系统 SHALL 支持以下配置步骤：
  1. 基础信息（已有）
  2. 动态表单（已有）
  3. 奖励规则（已有）
  4. 支付参数（已有）
  5. **分销奖励规则（新增）**
     - 启用/禁用分销
     - 选择分销层级（1-3级）
     - 设置各级奖励比例
     - 选择模板或自定义

#### Scenario: 编辑活动修改分销规则
- **WHEN** 品牌管理员编辑活动
- **THEN** 系统 SHALL 允许修改分销规则
- **AND** 显示当前规则配置
- **AND** 提示新规则仅影响新订单

### Requirement: 活动数据存储
系统 SHALL 扩展 campaigns 表，增加分销相关字段。

#### Scenario: 分销规则字段
- **WHEN** 活动保存分销规则
- **THEN** 系统 SHALL 存储到 campaigns 表：
  - enable_distribution: BOOLEAN
  - distribution_level: INT
  - distribution_rewards: JSON

#### Scenario: 分销规则 JSON Schema
- **WHEN** 存储 distribution_rewards
- **THEN** 系统 SHALL 使用以下 JSON Schema：
  ```json
  {
    "type": "object",
    "properties": {
      "level1": {"type": "number", "minimum": 0, "maximum": 100},
      "level2": {"type": "number", "minimum": 0, "maximum": 100},
      "level3": {"type": "number", "minimum": 0, "maximum": 100}
    },
    "required": ["level1"]
  }
  ```

### Requirement: 奖励规则配置（扩展）
系统 SHALL 扩展奖励规则配置，与分销奖励规则共存。

#### Scenario: 单级奖励与分销奖励共存
- **WHEN** 活动配置了奖励规则（reward_rule）和分销规则
- **THEN** 系统 SHALL 支持：
  - reward_rule: 单级固定金额奖励（原有功能）
  - distribution_rewards: 多级比例奖励（新增功能）
- **AND** 两种奖励规则互斥或共存根据业务需求决定

#### Scenario: 仅配置单级奖励
- **WHEN** 活动仅配置 reward_rule 且未启用分销
- **THEN** 系统 SHALL 使用原有的单级奖励逻辑
- **AND** 不计算分销奖励

#### Scenario: 仅配置分销奖励
- **WHEN** 活动启用分销且配置了 distribution_rewards
- **THEN** 系统 SHALL 使用多级分销奖励逻辑
- **AND** 忽略 reward_rule（或根据业务需求）

#### Scenario: 同时配置两种奖励
- **WHEN** 活动同时配置了 reward_rule 和 distribution_rewards
- **THEN** 系统 SHALL 根据业务需求决定优先级
- **AND** 推荐优先使用分销奖励逻辑
