# Stability Checks Workflow 使用说明

本文档说明新增 CI 流水线 `.github/workflows/stability-checks.yml` 的触发条件、执行内容和排障方式。

## 1. 目标

- 在单条流水线中固化三类稳定性检查：
  - backend 全量测试：`go test ./...`
  - frontend-admin 单测：`npm run test`
  - frontend-admin 安全回归 E2E：`npm run test:e2e:headless -- e2e/security-management.spec.ts --reporter=line`
- 让安全管理页回归与后端全量回归成为默认门禁。

## 2. 触发条件

- 手动触发：`workflow_dispatch`
- PR 到 `main` 且改动命中以下路径之一：
  - `backend/**`
  - `frontend-admin/**`
  - `.github/workflows/stability-checks.yml`
- Push 到 `main` 且改动命中同样路径。

## 3. 流水线执行步骤

单个 job：`stability-checks`（`ubuntu-latest`）

1. Checkout 代码，安装 Go/Node 环境。
2. 启动 MySQL 8 与 Redis 7 服务容器。
3. 等待数据库与缓存健康检查通过。
4. 初始化数据库：
   - `backend/scripts/init.sql`
   - `backend/scripts/create_member_system_tables.sql`
   - `backend/migrations/20250120_create_distributor_tables_final.sql`
   - `backend/migrations/20250124_add_advanced_features.sql`
   - `backend/migrations/2026_01_29_add_record_tables.sql`
5. 修正订单表兼容字段（`orders.paid_at`、`verification_code`）。
6. 启动 backend API 并等待 `http://127.0.0.1:8889/api/v1/auth/login` 可达。
7. 执行 backend 全量测试：`go test ./...`。
8. 安装 frontend-admin 依赖并安装 Playwright Chrome 浏览器。
9. 执行 frontend-admin 单测与安全 E2E 回归。
10. 失败时上传调试产物（`/tmp/dmh-api.log`、Playwright 报告）。
11. 无论成败，最后停止后台 API 进程。

## 4. 环境变量与凭据

工作流内置：

- `DMH_INTEGRATION_BASE_URL=http://127.0.0.1:8889`
- `DMH_TEST_ADMIN_USERNAME=${{ vars.DMH_TEST_ADMIN_USERNAME || 'admin' }}`
- `DMH_TEST_ADMIN_PASSWORD=${{ secrets.DMH_TEST_ADMIN_PASSWORD || '123456' }}`

建议仓库配置：

- `secrets.DMH_TEST_ADMIN_PASSWORD`：避免使用默认密码。
- `vars.DMH_TEST_ADMIN_USERNAME`：如管理员账号非 `admin`，应显式配置。

## 5. 失败排障

- **API 不可达**
  - 看失败产物中的 `/tmp/dmh-api.log`。
  - 优先检查 DB 初始化失败、配置加载失败、端口占用。
- **安全 E2E 失败**
  - 下载 `frontend-admin/playwright-report` 与 `frontend-admin/test-results`。
  - 检查登录态、页面渲染与接口响应（尤其 `/api/v1/security/*`）。
- **数据库迁移失败**
  - 重点查看 `Ensure orders.paid_at and verification_code schema` 步骤日志。
  - 确认 SQL 在 MySQL 8 下兼容。

## 6. 本地复现建议

若要本地复现 CI 关键路径，建议先启动依赖服务（MySQL/Redis + backend），再执行：

```bash
cd backend
go test ./...

cd ../frontend-admin
npm ci
npm run test
npm run test:e2e:headless -- e2e/security-management.spec.ts --reporter=line
```

如需严格复刻 CI，可按 `.github/workflows/stability-checks.yml` 中的 DB 初始化 SQL 顺序执行。
