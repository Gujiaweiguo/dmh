# Design: Unified Config Management

## Context

DMH 项目当前配置分散在多个位置，需要统一管理以简化运维。

**Stakeholders**:
- 运维人员：需要快速定位和修改配置
- 开发人员：需要了解配置与代码的对应关系

**Constraints**:
- 必须兼容现有 Docker Compose 部署方式
- 不能影响开发环境的配置结构
- 配置变更需要可追溯、可回滚

## Goals / Non-Goals

### Goals
- ✅ 统一生产环境配置到 `/opt/module/dmh/configs/`
- ✅ 提供配置备份和回滚能力
- ✅ 提供配置变更验证脚本
- ✅ 简化运维操作流程

### Non-Goals
- ❌ 不改变开发环境配置结构（`backend/api/etc/`）
- ❌ 不实现配置热更新（需要重启容器）
- ❌ 不实现配置版本控制系统（仅简单备份）

## Decisions

### D1: 配置目录位置
**Decision**: 使用 `/opt/module/dmh/configs/` 作为统一配置目录

**Rationale**:
- `/opt/module/` 是系统约定用于第三方模块的目录
- 已存在 `/opt/module/dmh/` 目录（用于数据和日志）
- 配置与数据分离，便于备份和迁移

**Alternatives Considered**:
| 方案 | 优点 | 缺点 | 决策 |
|------|------|------|------|
| `/opt/module/dmh/configs/` | 符合现有结构 | - | ✅ 采用 |
| `/etc/dmh/` | Linux 标准配置路径 | 需要权限 | ❌ |
| `deploy/configs/` | 靠近代码 | 与代码耦合 | ❌ |

### D2: 挂载方式
**Decision**: Docker Compose 直接挂载 `/opt/module/dmh/configs/` 目录

**Rationale**:
- 修改配置后重启容器即可生效
- 无需复制文件到 deploy 目录
- 配置与代码完全分离

**Implementation**:
```yaml
# docker-compose-simple.yml
volumes:
  # 之前
  - ./dmh-api.yaml:/app/etc/dmh-api.yaml:ro
  # 之后
  - /opt/module/dmh/configs/dmh-api.yaml:/app/etc/dmh-api.yaml:ro
```

### D3: 备份策略
**Decision**: 简单时间戳备份，保留最近 10 个版本

**Rationale**:
- 实现简单，满足基本回滚需求
- 不引入额外依赖（如 git）
- 自动清理旧备份，避免磁盘占用过大

**Implementation**:
```
/opt/module/dmh/configs/backup/
├── 2026-02-19_14-30-00/    # 备份时间戳
│   ├── dmh-api.yaml
│   └── nginx/conf.d/default.conf
├── 2026-02-19_15-00-00/
└── ...（最多保留 10 个）
```

### D4: 验证策略
**Decision**: 分层验证（语法 → 连接 → 功能）

**Rationale**:
- 逐步验证，快速定位问题
- 每层独立，可选择性执行

**Validation Levels**:
| Level | 检查项 | 失败处理 |
|-------|--------|----------|
| L1 语法 | YAML/Nginx 配置语法 | 阻止重启 |
| L2 连接 | DB/Redis 连接性 | 警告 |
| L3 功能 | API 健康检查 | 警告 |

## Risks / Trade-offs

### Risk 1: 配置路径变更导致服务无法启动
- **Likelihood**: 低
- **Impact**: 高
- **Mitigation**: 
  - 提供 `sync-configs.sh` 脚本自动迁移配置
  - 启动前验证配置文件存在

### Risk 2: 备份磁盘占用
- **Likelihood**: 中
- **Impact**: 低
- **Mitigation**: 
  - 限制备份数量为 10 个
  - 提供清理脚本

### Risk 3: 配置同步问题
- **Likelihood**: 中
- **Impact**: 中
- **Mitigation**: 
  - 文档明确说明配置位置
  - 提供 `sync-configs.sh` 脚本

## Migration Plan

### Phase 1: 准备（不影响服务）
1. 创建 `/opt/module/dmh/configs/` 目录结构
2. 复制现有配置到新目录
3. 创建管理脚本

### Phase 2: 切换（短暂中断）
1. 修改 `docker-compose-simple.yml` 挂载路径
2. 重启服务
3. 验证服务正常

### Phase 3: 清理
1. 保留原配置文件作为备份
2. 更新文档

### Rollback Plan
```bash
# 回滚步骤
1. 恢复 docker-compose-simple.yml 挂载路径
2. docker compose down && docker compose up -d
3. 验证服务正常
```

## Open Questions
- 无（用户已确认所有关键决策）
