<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# DMH Agent Working Agreement

本文件定义在 DMH 仓库内工作的 AI/自动化代理协作约定。

## 1. 响应语言

* 默认使用中文回复（简体中文）。
* 仅在用户明确要求英文或其他语言时切换。

## 2. 项目上下文（DMH）

* 项目：DMH（Digital Marketing Hub）数字营销中台。
* 主要模块：

  * `backend/`：Go + go-zero + GORM 后端 API（默认端口 `8889`）
  * `frontend-admin/`：Vue 3 管理后台（默认端口 `3000`）
  * `frontend-h5/`：Vue 3 H5 前端（默认端口 `3100`）
  * `openspec/`：需求规格与变更提案体系
  * `deploy/`：容器化和部署脚本

### 2.1 项目结构

```
DMH/
├── backend/              # Go API (入口: api/dmh.go, 端口: 8889)
│   ├── api/internal/     # handler, logic, middleware, svc, types
│   ├── model/            # GORM 数据模型
│   ├── common/           # 共享工具 (poster, wechatpay, utils)
│   ├── migrations/       # SQL 迁移文件
│   └── test/             # 集成测试、性能测试
├── frontend-admin/       # 管理后台 (入口: index.tsx, 端口: 3000)
│   ├── views/            # 页面组件 (混合 .tsx/.vue)
│   ├── services/         # API 服务
│   └── tests/            # 单元测试
├── frontend-h5/          # H5 前端 (入口: src/main.js, 端口: 3100)
│   ├── src/              # 标准结构 (views, components, router)
│   └── tests/unit/       # 逻辑测试 (*.logic.js)
├── openspec/             # 规格驱动开发系统
├── deploy/               # Docker Compose, nginx 配置
└── docs/                 # 文档
```

### 2.2 非标准结构（注意）

| 模块 | 非标准模式 |
|------|-----------|
| `frontend-admin/` | 入口是 `index.tsx`（非 main.js），扁平结构（无 src/），混合 .tsx/.vue 文件 |
| `frontend-h5/` | 每个 view 有 companion `.logic.js` 文件 |
| `backend/` | go-zero 标准结构，handler 薄层 + logic 业务层 |

## 3. 本项目代码与提交流程约定

* 变更应小步、聚焦，避免无关重构。
* 后端代码遵循 `gofmt`。
* 前端代码遵循 ESLint + Prettier 约定。
* 提交信息遵循 Conventional Commits（如 `feat: ...`、`fix: ...`）。
* 优先沿用现有目录结构与命名风格，不随意引入新模式。

### 3.1 反模式（禁止）

| 模式 | 原因 |
|------|------|
| Handler 中写业务逻辑 | Handler 只负责解析请求、调用 Logic、返回响应 |
| 跳过 Logic 层直接访问 DB | 必须通过 Logic 层，保持分层清晰 |
| `as any` 类型断言 | 使用正确的类型定义 |
| 忽略错误返回值 | 必须处理或包装错误 |
| 未经测试的集成代码 | 集成测试必须通过 |
| 在 `frontend-admin` 使用 `.vue` | 优先使用 `.tsx`（现有文件可保留） |
| OpenSpec 场景格式错误 | 必须使用 `#### Scenario:` 格式（4个#） |

### 3.2 代码入口

| 模块 | 入口文件 | 启动命令 |
|------|---------|---------|
| 后端 | `backend/api/dmh.go` | `go run api/dmh.go -f api/etc/dmh-api.yaml` |
| Admin | `frontend-admin/index.tsx` | `npm run dev` |
| H5 | `frontend-h5/src/main.js` | `npm run dev` |

## 4. OpenSpec 约定（必须遵循）

OpenSpec 详细规则见 `openspec/AGENTS.md`，此处给出执行要点。

### 4.1 何时需要先走提案（proposal）

以下场景先创建 OpenSpec change，再进入实现：

* 新功能或新能力
* 架构/模式调整
* 可能影响行为的性能优化
* 接口、数据结构等潜在 breaking change

以下场景可直接改代码（通常不必提案）：

* bug 修复（恢复既有预期行为）
* 文案、注释、格式化、轻量配置调整
* 非破坏性依赖升级

### 4.2 标准流程

1. 先阅读上下文：`openspec/project.md`、`openspec list`、`openspec list --specs`
2. 选择唯一 `change-id`（kebab-case，动词前缀：`add-`/`update-`/`remove-`/`refactor-`）
3. 在 `openspec/changes/<change-id>/` 下编写：

   * `proposal.md`
   * `tasks.md`
   * `design.md`（仅在复杂/跨模块/高风险时需要）
   * 对应 capability 的 spec delta 文件
4. spec delta 必须使用 `ADDED/MODIFIED/REMOVED/RENAMED Requirements`，且每条 Requirement 至少一个 `#### Scenario:`
5. 严格校验：`openspec validate <change-id> --strict --no-interactive`
6. 提案获批后再实现代码

### 4.3 实现与归档

* 实现前先读 `proposal.md`、`tasks.md`（如有 `design.md` 也需阅读）。
* 按 `tasks.md` 顺序执行，并在完成后把任务勾选为 `- [x]`。
* 发布后归档：`openspec archive <change-id> --yes`（仅工具类变更可考虑 `--skip-specs`）。

## 5. 常用命令（按需执行）

* 后端测试：`cd backend && go test ./...`
* 后端启动：`cd backend && go run api/dmh.go -f api/etc/dmh-api.yaml`
* 管理端启动：`cd frontend-admin && npm run dev`
* H5 启动：`cd frontend-h5 && npm run dev`
* OpenSpec 列表：`openspec list`
* OpenSpec 规格：`openspec list --specs`
* OpenSpec 校验：`openspec validate --strict --no-interactive`

## 6. 执行原则

* 先理解需求与现有实现，再动手修改。
* 涉及多模块或高风险改动时，先明确范围与影响。
* 未经用户明确要求，不主动提交 commit 或做破坏性操作。

## 7. 技术栈详细说明

### 7.1 后端技术栈

* **Go 版本**：1.24.0
* **框架**：go-zero 1.6.0（微服务框架，API 网关、RPC、中间件）
* **ORM**：GORM 1.25.5（数据访问层）
* **数据库**：MySQL 8.0+
* **缓存**：Redis 7+
* **JWT**：golang-jwt/jwt/v4 4.5.0（身份认证）
* **数据库驱动**：gorm.io/driver/mysql 1.5.2
* **测试框架**：stretchr/testify 1.10.0
* **代码格式**：gofmt（标准 Go 格式化工具）

### 7.2 前端技术栈（管理后台）

* **Node.js**：>=20.19.0
* **Vue 版本**：^3.4.0
* **构建工具**：Vite ^6.2.0
* **TypeScript**：~5.8.2
* **UI 组件库**：Lucide Vue Next ^0.263.0（图标库）
* **样式**：Tailwind CSS ^3.4.19
* **测试**：
  - Vitest ^2.1.8（单元测试）
  - Playwright ^1.58.2（E2E 测试）
* **代码规范**：ESLint + Prettier

### 7.3 前端技术栈（H5）

* **Node.js**：>=20.19.0
* **Vue 版本**：^3.5.27
* **构建工具**：Vite ^5.4.21
* **UI 组件库**：Vant ^4.9.22（移动端组件）
* **路由**：Vue Router ^4.6.4
* **HTTP 客户端**：Axios ^1.13.4
* **功能库**：
  - qrcode ^1.5.4（二维码生成）
  - html5-qrcode ^2.3.8（二维码扫描）
  - html2canvas ^1.4.1（截图）
  - vuedraggable ^2.24.3（拖拽）
* **测试**：
  - Vitest ^2.1.8（单元测试）
  - Playwright ^1.58.2（E2E 测试）

## 8. Docker Compose 工作约定

### 8.1 容器化部署架构

**网络配置**：`my-net` (172.19.0.0/16)

**服务组成**：
* `dmh-api`：后端 API 容器（镜像：debian:bookworm-slim）
* `dmh-nginx`：前端托管容器（镜像：nginx:1.25-alpine）
* `mysql8`：MySQL 数据库容器（镜像：mysql:8.0，预存在）
* `redis-dmh`：Redis 缓存容器（镜像：redis:7）

**Docker Compose 文件**：
* `deploy/docker-compose-simple.yml`：简化版（推荐开发环境）
* `deploy/docker-compose.yml`：完整版（包含镜像构建，适合生产）

**服务端口**：

| 服务 | 容器内端口 | 宿主机端口 |
|------|-----------|-----------|
| 后端API | 8889 | 8889 |
| 管理后台 | 3000 | 3000 |
| H5前端 | 3100 | 3100 |

### 8.2 容器管理命令

**启动所有服务**：
```bash
cd deploy
docker compose -f docker-compose-simple.yml up -d
# 或使用快速启动脚本
./scripts/quick-start.sh
```

**查看容器状态**：
```bash
docker ps
docker compose -f docker-compose-simple.yml ps
```

**查看日志**：
```bash
# 所有服务
docker compose -f docker-compose-simple.yml logs -f
# 单个服务
docker logs -f dmh-api
docker logs -f dmh-nginx
```

**重启服务**：
```bash
# 单个服务
docker restart dmh-api
docker restart dmh-nginx
# 所有服务
docker compose -f docker-compose-simple.yml restart
```

**停止服务**：
```bash
docker compose -f docker-compose-simple.yml stop
```

### 8.3 代码更新工作流程

**后端代码更新**：
```bash
# 1. 本地编译
cd backend
go build -o dmh-api api/dmh.go

# 2. 复制到 deploy 目录
cp dmh-api ../deploy/dmh-api
chmod +x ../deploy/dmh-api

# 3. 重启容器
cd ../deploy
docker restart dmh-api
```

**前端代码更新**：
```bash
# 1. 构建前端
cd frontend-admin
npm run build

# 2. 重启 nginx 容器
cd ../deploy
docker restart dmh-nginx
```

### 8.4 容器内测试执行

**在容器内运行测试**（推荐用于集成测试）：
```bash
export DMH_INTEGRATION_BASE_URL=http://localhost:8889
go test ./test/integration/... -v -count=1
```

**测试环境变量**：
* `DMH_INTEGRATION_BASE_URL`：API 基础地址（默认：`http://localhost:8889`）
* `DMH_TEST_ADMIN_USERNAME`：测试管理员账号（默认：`admin`）
* `DMH_TEST_ADMIN_PASSWORD`：测试管理员密码（默认：`123456`）

### 8.5 故障排查

**容器启动失败**：
1. 查看容器日志：`docker logs dmh-api`
2. 检查端口占用：`lsof -i :8889` / `lsof -i :3000` / `lsof -i :3100`
3. 检查网络：`docker network inspect my-net`

**API 无法访问**：
1. 测试数据库连接：`docker exec dmh-api wget -q -O - http://mysql8:3306`
2. 查看后端日志：`docker logs dmh-api | grep -E "Error|Starting"`

**前端页面无法加载**：
1. 检查构建产物：`ls -la frontend-admin/dist` / `ls -la frontend-h5/dist`
2. 查看 nginx 日志：`docker logs dmh-nginx`

### 8.6 关键注意事项

1. **二进制文件更新**：容器内 dmh-api 文件被挂载为 volume，必须在容器停止时才能替换
2. **后端编译**：容器内未安装 Go，必须在宿主机编译后复制到容器
3. **网络隔离**：容器在 my-net 网络中，容器间通过服务名通信（mysql8、redis-dmh）
4. **配置文件**：
   - Nginx 配置：`deploy/nginx/conf.d/default.conf`
   - 后端配置：`deploy/dmh-api.yaml`

## 9. 常用测试命令

**后端单元测试**：
```bash
cd backend
go test ./...
```

**后端集成测试**：
```bash
cd backend
export DMH_INTEGRATION_BASE_URL=http://localhost:8889
go test ./test/integration/... -v -count=1
```

## 10. 常见问题 (Troubleshooting)

### 数据库连接失败
**错误**: `Error 1045: Access denied for user 'root'@'localhost'`

**解决**:
```bash
# 检查 MySQL 容器
 docker ps | grep mysql8

# 检查配置密码
cat backend/api/etc/dmh-api.yaml

# 重新初始化
make db-init  # 或 ./dmh.sh init
```

### 前端依赖安装失败
**错误**: `Error: Cannot find module 'xxx'`

**解决**:
```bash
cd frontend-admin  # 或 frontend-h5
rm -rf node_modules package-lock.json
npm install
```

### Go 依赖下载慢
**解决**:
```bash
# 设置代理
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
```

### 端口被占用
**错误**: `bind: address already in use`

**解决**:
```bash
# 查看占用
lsof -i :8889  # 后端
lsof -i :3000  # Admin
lsof -i :3100  # H5

# 停止服务
make down  # 或 ./dmh.sh stop
```

**前端单元测试（管理后台）**：
```bash
cd frontend-admin
npm run test
```

**前端 E2E 测试（管理后台）**：
```bash
cd frontend-admin
npm run test:e2e
```

**前端单元测试（H5）**：
```bash
cd frontend-h5
npm run test
```

**前端 E2E 测试（H5）**：
```bash
cd frontend-h5
npm run test:e2e
```
