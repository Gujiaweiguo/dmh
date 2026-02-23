# Proposal: Add Backend Promoter, Material, and AI APIs

## Summary

实现 H5 品牌管理端所需的后端 API，包括推广员管理、素材管理和 AI 文案生成接口。

## Motivation

H5 品牌管理端已实现以下功能的前端代码，但后端 API 尚未实现：

1. **推广员管理** (`/promoter/*`)
   - 推广员列表
   - 推广员详情
   - 推广链接生成
   - 奖励记录

2. **素材管理** (`/material/*`)
   - 素材上传
   - 素材列表
   - 素材删除

3. **AI 文案生成** (`/ai/*`)
   - 营销文案生成

当前前端使用 Mock 数据或降级处理，需要后端支持完整功能。

## Proposed Changes

### 1. Promoter API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/promoter/list` | 获取推广员列表（分页） |
| GET | `/promoter/detail/:id` | 获取推广员详情 |
| POST | `/promoter/generate-link` | 生成推广链接 |
| GET | `/promoter/rewards/:promoterId` | 获取推广员奖励记录 |

### 2. Material API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/material/upload` | 上传素材 |
| GET | `/material/list` | 获取素材列表（分页） |
| DELETE | `/material/delete/:id` | 删除素材 |

### 3. AI API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/ai/generate-copywriting` | 生成营销文案 |

## Implementation Notes

### Promoter API

- 推广员数据可复用现有的 Distributor 模型（筛选 role=promoter）
- 推广链接生成逻辑与分销商链接生成类似
- 奖励记录复用现有的 Reward 模型

### Material API

- 素材存储可使用现有的 BrandAsset 模型
- 文件上传使用 multipart/form-data
- 支持 image 和 text 类型素材

### AI API

- 集成第三方 AI 服务（如 OpenAI、通义千问等）
- 配置化 API Key 和模型参数
- 降级处理：API 不可用时返回模板文案

## Scope

### In Scope

- 上述 8 个 API 的实现
- 对应的 handler 和 logic 层代码
- 单元测试

### Out of Scope

- AI 服务选型和采购
- 素材 CDN 加速
- 推广员高级统计报表

## Success Criteria

1. 所有 8 个 API 通过单元测试
2. H5 前端可以正常调用后端 API（不再使用 Mock）
3. API 响应时间 < 500ms（AI API 除外，允许 < 5s）

## Dependencies

- 无新增外部依赖
- AI API 需配置第三方服务（可使用 Mock 实现）

## Timeline

- Phase 1: Promoter API (2 tasks)
- Phase 2: Material API (2 tasks)
- Phase 3: AI API (2 tasks)
- Phase 4: Testing & Documentation (2 tasks)

## Related Documents

- `frontend-h5/src/services/brandApi.js` - 前端 API 调用定义
- `backend/api/internal/handler/routes.go` - 路由配置
- `docs/H5_PENDING_FEATURES.md` - 功能清单
