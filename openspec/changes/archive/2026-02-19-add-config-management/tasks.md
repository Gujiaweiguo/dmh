# Tasks: Add Unified Config Management

> **执行顺序**: 严格按 Phase 顺序执行，Phase 内任务可并行
> **风险等级**: 🟢 Low | 🟡 Medium | 🔴 High
> **状态**: ✅ 已完成 (2026-02-19)

---

## Phase 1: 目录准备与配置迁移 [🟡 Medium] ✅

### Task 1.1: 创建统一配置目录 ✅
- **Task ID**: `CFG-001`
- **输入**: 无
- **输出**: 
  - `/opt/module/dmh/configs/` 目录结构
  - `/opt/module/dmh/configs/backup/` 备份目录
- **验收标准**:
  - [x] 目录结构符合 design.md 定义
  - [x] 权限正确 (755)
- **回滚方案**: `rm -rf /opt/module/dmh/configs/`
- **依赖**: 无
- **风险**: 🟢 Low

### Task 1.2: 迁移后端 API 配置 ✅
- **Task ID**: `CFG-002`
- **输入**: `/opt/code/dmh/deploy/dmh-api.yaml`
- **输出**: `/opt/module/dmh/configs/dmh-api.yaml`
- **验收标准**:
  - [x] 文件内容与源文件一致
  - [x] 文件权限 644
  - [x] YAML 语法正确
- **回滚方案**: 删除目标文件，重新复制
- **依赖**: `CFG-001`
- **风险**: 🟢 Low

### Task 1.3: 迁移 Nginx 配置 ✅
- **Task ID**: `CFG-003`
- **输入**: `/opt/code/dmh/deploy/nginx/conf.d/default.conf`
- **输出**: `/opt/module/dmh/configs/nginx/conf.d/default.conf`
- **验收标准**:
  - [x] 文件内容与源文件一致
  - [x] Nginx 配置语法正确
- **回滚方案**: 删除目标文件，重新复制
- **依赖**: `CFG-001`
- **风险**: 🟢 Low

### Task 1.4: 创建前端环境变量文件 ✅
- **Task ID**: `CFG-004`
- **输入**: 无（使用默认值）
- **输出**: 
  - `/opt/module/dmh/configs/frontend/admin.env`
  - `/opt/module/dmh/configs/frontend/h5.env`
- **验收标准**:
  - [x] 文件存在且格式正确
  - [x] 包含必要的 API 地址配置
- **回滚方案**: 删除文件
- **依赖**: `CFG-001`
- **风险**: 🟢 Low

---

## Phase 2: 脚本开发 [🟡 Medium] ✅

### Task 2.1: 创建配置同步脚本 ✅
- **Task ID**: `CFG-005`
- **输入**: 项目配置文件
- **输出**: `deploy/scripts/sync-configs.sh`
- **验收标准**:
  - [x] 脚本可执行 (chmod +x)
  - [x] 正确复制所有配置文件
  - [x] 显示同步结果
- **回滚方案**: 删除脚本
- **依赖**: Phase 1 完成
- **风险**: 🟢 Low

### Task 2.2: 创建配置备份脚本 ✅
- **Task ID**: `CFG-006`
- **输入**: `/opt/module/dmh/configs/` 当前配置
- **输出**: 
  - `deploy/scripts/backup-config.sh`
  - `/opt/module/dmh/configs/backup/TIMESTAMP/` 备份目录
- **验收标准**:
  - [x] 脚本可执行
  - [x] 创建带时间戳的备份目录
  - [x] 保留最近 10 个备份
  - [x] 自动清理旧备份
- **回滚方案**: 删除脚本，手动删除备份目录
- **依赖**: Phase 1 完成
- **风险**: 🟢 Low

### Task 2.3: 创建配置验证脚本 ✅
- **Task ID**: `CFG-007`
- **输入**: 配置文件
- **输出**: 
  - `deploy/scripts/verify-config.sh`
  - 验证结果（退出码）
- **验收标准**:
  - [x] L1 语法检查：YAML/Nginx 配置
  - [x] L2 连接检查：DB/Redis 可达性
  - [x] L3 功能检查：API 健康端点
  - [x] 返回适当的退出码
- **回滚方案**: 删除脚本
- **依赖**: Phase 1 完成
- **风险**: 🟢 Low

### Task 2.4: 创建服务重启脚本 ✅
- **Task ID**: `CFG-008`
- **输入**: 无
- **输出**: `deploy/scripts/restart-services.sh`
- **验收标准**:
  - [x] 脚本可执行
  - [x] 自动备份当前配置
  - [x] 重启 Docker 容器
  - [x] 等待服务就绪
  - [x] 执行健康检查
  - [x] 失败时提示回滚命令
- **回滚方案**: 删除脚本
- **依赖**: `CFG-006`, `CFG-007`
- **风险**: 🟢 Low

---

## Phase 3: Docker Compose 更新 [🔴 High] ✅

### Task 3.1: 更新 docker-compose 挂载路径 ✅
- **Task ID**: `CFG-009`
- **输入**: `deploy/docker-compose-simple.yml`
- **输出**: 修改后的 `deploy/docker-compose-simple.yml`
- **验收标准**:
  - [x] dmh-api.yaml 挂载路径更新
  - [x] nginx 配置挂载路径更新
  - [x] YAML 语法正确
  - [x] 容器可正常启动
- **回滚方案**: 
  ```bash
  git checkout deploy/docker-compose-simple.yml
  docker compose -f deploy/docker-compose-simple.yml up -d
  ```
- **依赖**: Phase 1, Phase 2 完成
- **风险**: 🔴 High - 影响服务运行

### Task 3.2: 重启服务并验证 ✅
- **Task ID**: `CFG-010`
- **输入**: 更新后的 docker-compose 配置
- **输出**: 运行中的服务
- **验收标准**:
  - [x] `docker compose up -d` 成功
  - [x] API 服务运行中 (HTTP 404 - 端点不存在但服务正常)
  - [x] 管理后台可访问: HTTP 200
  - [x] H5 前端可访问: HTTP 200
  - [x] 数据库连接正常
- **回滚方案**: 
  ```bash
  # 恢复原配置
  git checkout deploy/docker-compose-simple.yml
  docker compose -f deploy/docker-compose-simple.yml down
  docker compose -f deploy/docker-compose-simple.yml up -d
  ```
- **依赖**: `CFG-009`
- **风险**: 🔴 High - 服务中断

---

## Phase 4: 文档更新 [🟢 Low] ✅

### Task 4.1: 更新 /opt/module/dmh/README.md ✅
- **Task ID**: `CFG-011`
- **输入**: 无
- **输出**: `/opt/module/dmh/README.md`
- **验收标准**:
  - [x] 说明配置目录结构
  - [x] 说明配置修改流程
  - [x] 说明脚本使用方法
  - [x] 说明回滚方法
- **回滚方案**: 删除或恢复原文件
- **依赖**: Phase 3 完成
- **风险**: 🟢 Low

### Task 4.2: 更新 deploy/README.md ✅
- **Task ID**: `CFG-012`
- **输入**: 无
- **输出**: `deploy/README.md` (追加章节)
- **验收标准**:
  - [x] 添加"配置管理"章节
  - [x] 说明新的配置位置
  - [x] 提供脚本使用示例
- **回滚方案**: `git checkout deploy/README.md`
- **依赖**: Phase 3 完成
- **风险**: 🟢 Low

---

## 任务依赖图

```
Phase 1 (目录准备) ✅
├── CFG-001 ✅ ─┬─→ CFG-002 ✅ (API配置)
│              ├─→ CFG-003 ✅ (Nginx配置)
│              └─→ CFG-004 ✅ (前端配置)
│
Phase 2 (脚本开发) ✅ - 依赖 Phase 1
├── CFG-005 ✅ (sync)
├── CFG-006 ✅ (backup) ────┐
├── CFG-007 ✅ (verify) ────┼─→ CFG-008 ✅ (restart)
│                          │
Phase 3 (服务切换) ✅       │
└── CFG-009 ✅ (compose) ───┴─→ CFG-010 ✅ (重启验证)
                                  │
Phase 4 (文档) ✅                 │
├── CFG-011 ✅ ←──────────────────┘
└── CFG-012 ✅
```

## 风险汇总

| Phase | 风险等级 | 主要风险 | 状态 |
|-------|----------|----------|------|
| Phase 1 | 🟡 Medium | 配置文件丢失 | ✅ 无问题 |
| Phase 2 | 🟡 Medium | 脚本错误 | ✅ 无问题 |
| Phase 3 | 🔴 High | **服务中断** | ✅ 成功重启 |
| Phase 4 | 🟢 Low | 文档不准确 | ✅ 已更新 |

## 预计时间 vs 实际时间

| Phase | 任务数 | 预计时间 | 实际时间 |
|-------|--------|----------|----------|
| Phase 1 | 4 | 10 分钟 | ~5 分钟 |
| Phase 2 | 4 | 20 分钟 | ~10 分钟 |
| Phase 3 | 2 | 5 分钟 | ~5 分钟 |
| Phase 4 | 2 | 10 分钟 | ~5 分钟 |
| **总计** | **12** | **45 分钟** | **~25 分钟** |

## 执行检查清单

开始前确认：
- [x] 当前服务运行正常
- [x] 有 SSH 访问权限
- [x] 有 Docker 操作权限
- [x] 已备份当前配置

完成后确认：
- [x] 所有脚本可执行
- [x] 服务重启成功
- [x] 所有健康检查通过
- [x] 文档已更新
