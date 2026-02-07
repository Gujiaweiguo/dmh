# Change: Add brand-admin poster generation and distribution rules

## Why
- 品牌管理员缺少生成活动海报与分销规则配置入口，影响活动上线与分销推广。

## What Changes
- 后端补齐海报模板列表查询，活动接口返回/保存分销规则与海报模板字段。
- 管理端活动编辑页新增分销规则配置。
- 管理端活动列表新增生成海报按钮。

## Impact
- Affected specs: campaign-management
- APIs:
  - GET /api/v1/poster/templates
  - POST /api/v1/campaigns/:id/poster
  - GET/POST/PUT /api/v1/campaigns (返回/保存分销字段)
- UI:
  - frontend-admin/views/CampaignEditorView.tsx
  - frontend-admin/views/CampaignListView.tsx
- Data:
  - campaigns.enable_distribution
  - campaigns.distribution_level
  - campaigns.distribution_rewards
  - campaigns.poster_template_id
  - poster_template_configs (读取)
- Migrations: none (复用现有字段)
