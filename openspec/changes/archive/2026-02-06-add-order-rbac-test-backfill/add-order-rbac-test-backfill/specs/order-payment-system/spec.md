## ADDED Requirements

### Requirement: Automated regression coverage for order core flows
系统 SHALL 为订单核心流程提供自动化回归测试覆盖，以确保后续迭代不破坏既有行为。

#### Scenario: Order creation regression suite
- **WHEN** 代码变更涉及订单创建逻辑
- **THEN** 系统 SHALL 运行订单创建单元测试覆盖正常与异常分支
- **AND** 测试 SHALL 覆盖活动有效性、重复订单防护与关键字段校验

#### Scenario: Order verification regression suite
- **WHEN** 代码变更涉及订单核销逻辑
- **THEN** 系统 SHALL 运行订单核销相关单元测试
- **AND** 测试 SHALL 覆盖核销、重复核销、权限不足与取消核销路径
