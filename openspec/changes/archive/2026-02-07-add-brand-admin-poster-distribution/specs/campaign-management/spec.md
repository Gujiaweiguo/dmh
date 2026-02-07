## ADDED Requirements

### Requirement: Brand admin generates campaign posters
系统 SHALL 允许品牌管理员在管理端为活动生成海报。

#### Scenario: Generate poster from campaign list
- **WHEN** 品牌管理员在活动列表点击“生成海报”
- **THEN** 系统 SHALL 使用该活动的海报模板生成海报
- **AND** 返回可预览/下载的海报 URL

### Requirement: Brand admin configures distribution rules in campaign editor
系统 SHALL 允许品牌管理员在活动创建/编辑时配置分销规则。

#### Scenario: Configure distribution rules
- **WHEN** 品牌管理员启用分销
- **THEN** 系统 SHALL 允许选择分销层级（1-3级）
- **AND** 允许填写各级奖励比例
- **AND** 配置结果 SHALL 保存到活动数据中
