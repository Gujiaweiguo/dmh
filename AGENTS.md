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
  * `deployment/`：容器化和部署脚本

## 3. 本项目代码与提交流程约定

* 变更应小步、聚焦，避免无关重构。
* 后端代码遵循 `gofmt`。
* 前端代码遵循 ESLint + Prettier 约定。
* 提交信息遵循 Conventional Commits（如 `feat: ...`、`fix: ...`）。
* 优先沿用现有目录结构与命名风格，不随意引入新模式。

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
