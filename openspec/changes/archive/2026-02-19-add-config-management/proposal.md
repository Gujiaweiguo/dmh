# Change: Add Unified Config Management

## Why

当前 DMH 项目的配置文件分散在多处（`backend/api/etc/`、`deploy/`），生产环境修改配置需要：
1. 定位配置文件位置（开发配置 vs 生产配置）
2. 修改后重启容器
3. 手动验证服务是否正常

缺乏统一的配置管理入口和变更验证机制，增加了运维复杂度和出错风险。

## What Changes

### 新增配置管理能力
- **统一配置目录**: `/opt/module/dmh/configs/` 作为所有生产环境配置的唯一来源
- **配置文件范围**: 
  - `dmh-api.yaml` - 后端 API 配置（数据库、Redis、JWT、微信支付等）
  - `nginx/conf.d/default.conf` - Nginx 反向代理配置
  - `frontend/admin.env` - 管理后台环境变量
  - `frontend/h5.env` - H5 前端环境变量
  - `docker-compose.yml` - Docker 编排配置（可选）
- **备份机制**: 修改配置前自动备份到 `/opt/module/dmh/configs/backup/`
- **验证脚本**: 配置变更后自动验证服务健康状态

### Docker Compose 挂载调整
- 将配置文件挂载路径从 `./deploy/` 改为 `/opt/module/dmh/configs/`
- 修改后重启容器即可生效，无需重新构建镜像

## Impact

### Affected Specs
- **NEW**: `config-management` - 新增配置管理能力规格

### Affected Code
- `deploy/docker-compose-simple.yml` - 修改配置文件挂载路径（第 96 行）
- `deploy/scripts/` - 新增配置管理脚本

### 新建目录结构
```
/opt/module/dmh/configs/
├── dmh-api.yaml
├── nginx/conf.d/default.conf
├── frontend/admin.env
├── frontend/h5.env
└── backup/
```

### 风险等级
- 🟡 中等风险：涉及生产环境配置路径变更
- 🟢 可回滚：备份机制确保可快速回退

## Acceptance Criteria
- [x] 配置文件已迁移到 `/opt/module/dmh/configs/`
- [x] Docker Compose 挂载路径已更新
- [x] 备份/验证/重启脚本可用
- [x] 文档更新
