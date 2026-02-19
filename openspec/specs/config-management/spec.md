# config-management Specification

## Purpose
TBD - created by archiving change add-config-management. Update Purpose after archive.
## Requirements
### Requirement: Unified Config Directory

系统 SHALL 在 `/opt/module/dmh/configs/` 目录下统一管理所有生产环境配置文件。

#### Scenario: Directory structure creation
- **WHEN** 系统初始化配置管理
- **THEN** 创建以下目录结构：
  - `configs/dmh-api.yaml` - 后端配置
  - `configs/nginx/conf.d/default.conf` - Nginx 配置
  - `configs/frontend/admin.env` - 管理后台环境变量
  - `configs/frontend/h5.env` - H5 环境变量
  - `configs/backup/` - 备份目录

#### Scenario: Config file migration
- **WHEN** 运行 `sync-configs.sh` 脚本
- **THEN** 从 `deploy/` 目录复制配置文件到 `/opt/module/dmh/configs/`
- **AND** 保留原文件权限
- **AND** 显示同步结果

---

### Requirement: Config Backup Mechanism

系统 SHALL 在修改配置前自动创建备份，支持快速回滚。

#### Scenario: Automatic backup before modification
- **WHEN** 运行 `restart-services.sh` 或 `backup-config.sh`
- **THEN** 创建带时间戳的备份目录
- **AND** 备份目录格式为 `backup/YYYY-MM-DD_HH-MM-SS/`
- **AND** 复制所有当前配置到备份目录

#### Scenario: Backup retention limit
- **WHEN** 备份数量超过 10 个
- **THEN** 自动删除最旧的备份
- **AND** 保留最近 10 个备份

#### Scenario: Manual backup creation
- **WHEN** 运行 `backup-config.sh` 脚本
- **THEN** 创建当前配置的完整备份
- **AND** 显示备份路径

---

### Requirement: Config Validation

系统 SHALL 提供多层配置验证，确保配置正确性。

#### Scenario: L1 Syntax validation - YAML
- **WHEN** 运行 `verify-config.sh` 
- **THEN** 检查 `dmh-api.yaml` 的 YAML 语法
- **AND** 语法错误时返回非零退出码
- **AND** 显示具体错误信息

#### Scenario: L1 Syntax validation - Nginx
- **WHEN** 运行 `verify-config.sh`
- **THEN** 检查 Nginx 配置语法
- **AND** 使用 `nginx -t` 验证

#### Scenario: L2 Connection validation - Database
- **WHEN** 运行 `verify-config.sh` 且数据库配置有效
- **THEN** 尝试连接 MySQL 数据库
- **AND** 连接失败时显示警告（不阻止操作）

#### Scenario: L2 Connection validation - Redis
- **WHEN** 运行 `verify-config.sh` 且 Redis 配置有效
- **THEN** 尝试连接 Redis
- **AND** 连接失败时显示警告（不阻止操作）

#### Scenario: L3 Function validation - API Health
- **WHEN** 运行 `verify-config.sh` 且服务已启动
- **THEN** 请求 `/api/v1/health` 端点
- **AND** 返回健康状态

---

### Requirement: Service Restart with Validation

系统 SHALL 提供一键重启服务并自动验证的脚本。

#### Scenario: Full restart workflow
- **WHEN** 运行 `restart-services.sh`
- **THEN** 执行以下步骤：
  1. 备份当前配置
  2. 验证配置语法（L1）
  3. 重启 Docker 容器
  4. 等待服务就绪（最长 60 秒）
  5. 执行健康检查（L3）
  6. 显示操作结果

#### Scenario: Restart failure with rollback hint
- **WHEN** 服务重启后健康检查失败
- **THEN** 显示失败原因
- **AND** 提示回滚命令
- **AND** 显示最近备份路径

#### Scenario: Restart success notification
- **WHEN** 服务重启且所有检查通过
- **THEN** 显示成功消息
- **AND** 显示服务访问地址

---

### Requirement: Docker Compose Integration

系统 SHALL 更新 Docker Compose 配置以挂载统一配置目录。

#### Scenario: Config mount path update
- **WHEN** 修改 `docker-compose-simple.yml`
- **THEN** `dmh-api.yaml` 挂载路径从 `./dmh-api.yaml` 改为 `/opt/module/dmh/configs/dmh-api.yaml`
- **AND** 容器重启后使用新配置

#### Scenario: Container restart after config change
- **WHEN** 修改 `/opt/module/dmh/configs/dmh-api.yaml`
- **AND** 运行 `docker restart dmh-api`
- **THEN** 容器使用新配置启动
- **AND** 无需重新构建镜像

---

### Requirement: Documentation

系统 SHALL 提供完整的配置管理文档。

#### Scenario: Main README update
- **WHEN** 查看 `/opt/module/dmh/README.md`
- **THEN** 包含以下内容：
  - 配置目录结构说明
  - 配置修改流程
  - 脚本使用方法
  - 回滚操作指南

#### Scenario: Deploy README update
- **WHEN** 查看 `deploy/README.md`
- **THEN** 包含"配置管理"章节
- **AND** 提供脚本使用示例

---

