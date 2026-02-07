## Context
- 品牌管理员需要在管理端完成活动海报生成与分销规则配置。
- 现有后端已具备相关数据字段，但活动接口未返回/保存全部字段，海报模板列表接口为空实现。

## Goals / Non-Goals
- Goals:
  - 管理端可配置分销规则并保存到活动
  - 管理端可生成活动海报并查看预览
  - 后端接口返回所需字段供管理端使用
- Non-Goals:
  - 不新增海报模板管理后台
  - 不新增数据库迁移

## Decisions
- 使用既有 `poster_template_configs` 作为模板来源。
- 分销规则以 `distribution_rewards` JSON 字符串保存。

## Risks / Trade-offs
- 管理端未提供模板管理功能，模板来源仍依赖预置数据。

## Migration Plan
- 无迁移，仅更新接口与前端。

## Open Questions
- 是否需要后续补充海报模板管理界面？
