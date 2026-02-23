## 1. Promoter API

- [x] 1.1 创建 promoter handler 目录和基础文件
  - `backend/api/internal/handler/promoter/` 目录
  - 路由注册到 `routes.go`

- [x] 1.2 实现 `GET /promoter/list` 接口
  - 分页查询推广员列表
  - 支持按状态、活动筛选
  - 返回推广员基础信息和统计数据

- [x] 1.3 实现 `GET /promoter/detail/:id` 接口
  - 获取推广员详细信息
  - 包含关联的活动和业绩统计

- [x] 1.4 实现 `POST /promoter/generate-link` 接口
  - 生成推广员专属链接
  - 复用现有链接生成逻辑

- [x] 1.5 实现 `GET /promoter/rewards/:promoterId` 接口
  - 查询推广员奖励记录
  - 支持分页和状态筛选

## 2. Material API

- [x] 2.1 创建 material handler 目录和基础文件
  - `backend/api/internal/handler/material/` 目录
  - 路由注册到 `routes.go`

- [x] 2.2 实现 `POST /material/upload` 接口
  - 支持 multipart/form-data 上传
  - 文件存储和元数据保存
  - 复用 BrandAsset 模型

- [x] 2.3 实现 `GET /material/list` 接口
  - 分页查询素材列表
  - 支持按类型筛选

- [x] 2.4 实现 `DELETE /material/delete/:id` 接口
  - 删除素材文件和元数据
  - 权限校验

## 3. AI API

- [x] 3.1 创建 ai handler 目录和基础文件
  - `backend/api/internal/handler/ai/` 目录
  - 路由注册到 `routes.go`
  - 配置文件添加 AI 服务参数

- [x] 3.2 实现 `POST /ai/generate-copywriting` 接口
  - 接收主题、风格、长度参数
  - 调用 AI 服务生成文案（或 Mock）
  - 返回生成结果

## 4. Verification

- [ ] 4.1 编写 handler 单元测试
  - promoter handler 测试
  - material handler 测试
  - ai handler 测试

- [ ] 4.2 集成测试
  - 前端调用后端 API 验证
  - E2E 测试更新

- [ ] 4.3 文档更新
  - API 文档更新
  - 更新 CHANGELOG.md
