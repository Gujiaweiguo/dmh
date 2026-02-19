# DMH 容器化部署指南

## 🚀 快速开始

### 一键启动（推荐）

```bash
cd /opt/code/dmh/deploy/scripts
./quick-start.sh
```

首次启动需要 2-5 分钟（安装依赖），之后只需 10-30 秒。

---

## 📋 目录结构

```
/opt/code/dmh/deploy/
├── docker-compose.yml           # 完整版Docker编排（包含构建）
├── docker-compose-simple.yml   # 简化版Docker编排（自动安装依赖）⭐
├── nginx/
│   ├── Dockerfile               # Nginx镜像构建文件
│   └── conf.d/
│       └── default.conf         # Nginx配置（3000管理、3100 H5）
└── scripts/
    ├── quick-start.sh          # 一键启动脚本 ⭐
    ├── deploy-containers.sh    # 完整部署脚本（构建镜像）
    ├── quick-restart.sh        # 快速重启脚本
    └── rollback-containers.sh   # 回滚到独立进程脚本
```

---

## 🏗️ 架构说明

```
my-net 网络 (172.19.0.0/16)
├── mysql8 (172.19.0.2)      [已存在] - MySQL数据库
├── redis7 (172.19.0.3)      [已存在] - Redis缓存
├── dataease-app (172.19.0.4) [已存在] - 其他应用
├── dmh-nginx (172.19.0.5)   [新建] - Nginx服务
│   ├── 端口 3000: 管理后台
│   ├── 端口 3100: H5前端
│   └── /api/ 代理 → dmh-api:8889
└── dmh-api (172.19.0.6)     [新建] - 后端API服务
    ├── 端口 8889: API服务
    ├── DB: mysql8:3306
    └── Redis: redis7:6379
```

---

## 📊 服务端口

| 服务 | 容器内端口 | 宿主机端口 | 说明 |
|------|-----------|-----------|------|
| 管理后台 | 3000 | 3000 | Vue 3 管理界面 |
| H5前端 | 3100 | 3100 | Vue 3 移动端界面 |
| 后端API | 8889 | 8889 | Go 后端API服务 |

---

## 🔧 部署方式

### 方式1：简化版（推荐）⭐

**特点**：
- 自动安装依赖
- 无需预构建镜像
- 启动即可使用

**启动命令**：
```bash
cd /opt/code/dmh/deployment
docker compose -f docker-compose-simple.yml up -d
```

**或使用快速启动脚本**：
```bash
cd /opt/code/dmh/deploy/scripts
./quick-start.sh
```

---

### 方式2：构建版

**特点**：
- 预先构建Docker镜像
- 启动速度快
- 适合生产环境

**启动命令**：
```bash
cd /opt/code/dmh/deploy/scripts
./deploy-containers.sh
```

---

## 🚦 服务管理

### 查看容器状态

```bash
cd /opt/code/dmh/deployment
docker compose -f docker-compose-simple.yml ps
```

### 查看日志

```bash
# 所有服务日志
docker compose -f docker-compose-simple.yml logs -f

# 单个服务日志
docker logs -f dmh-nginx
docker logs -f dmh-api
```

### 重启服务

```bash
# 重启所有服务
docker compose -f docker-compose-simple.yml restart

# 重启单个服务
docker restart dmh-nginx
docker restart dmh-api

# 使用快速重启脚本
./scripts/quick-restart.sh
```

### 停止服务

```bash
docker compose -f docker-compose-simple.yml stop
```

### 启动服务

```bash
docker compose -f docker-compose-simple.yml start
```

### 完全清理

```bash
docker compose -f docker-compose-simple.yml down
```

---

## 🔍 故障排查

### 容器启动失败

**检查容器日志**：
```bash
docker logs dmh-nginx
docker logs dmh-api
```

**检查容器状态**：
```bash
docker compose -f docker-compose-simple.yml ps
```

**常见问题**：
1. **端口被占用** - 检查 3000/3100/8889 端口是否被占用
2. **网络问题** - 确认 my-net 网络存在：`docker network inspect my-net`
3. **依赖未安装** - 首次启动需要 2-5 分钟，请耐心等待

---

### API无法访问

**测试数据库连接**：
```bash
docker exec dmh-api wget -q -O - http://mysql8:3306
```

**测试Redis连接**：
```bash
docker exec dmh-api wget -q -O - http://redis7:6379
```

**查看API日志**：
```bash
docker logs dmh-api | grep -E "Error|Starting|api"
```

---

### 前端页面无法加载

**检查前端构建产物**：
```bash
ls -la /opt/code/dmh/frontend-admin/dist
ls -la /opt/code/dmh/frontend-h5/dist
```

**检查容器内的文件**：
```bash
docker exec dmh-nginx ls -la /usr/share/nginx/html/admin
docker exec dmh-nginx ls -la /usr/share/nginx/html/h5
```

**查看Nginx日志**：
```bash
docker logs dmh-nginx | tail -50
```

---

## 🔄 回滚到独立进程

如果需要回滚到原来的独立进程部署方式：

```bash
cd /opt/code/dmh/deploy/scripts
./rollback-containers.sh
```

**回滚后需要**：
1. 单独配置 nginx 托管前端静态文件
2. 后端以独立进程方式运行（使用 `./deploy.sh`）
3. 端口访问地址不变（8889）

---

## 🛠️ 进入容器

### 进入Nginx容器

```bash
docker exec -it dmh-nginx sh
# 查看 nginx 配置
cat /etc/nginx/conf.d/default.conf
# 查看 nginx 日志
tail -f /var/log/nginx/access.log
```

### 进入API容器

```bash
docker exec -it dmh-api sh
# 查看日志
tail -f /var/log/dmh-api/*.log
# 测试数据库连接
wget -q -O - http://mysql8:3306
# 测试Redis连接
wget -q -O - http://redis7:6379
```

---

## 📝 配置说明

### Nginx配置

**文件位置**: `/opt/code/dmh/deploy/nginx/conf.d/default.conf`

**主要配置**：
- 管理后台监听 3000 端口
- H5前端监听 3100 端口
- `/api/` 路径代理到 `dmh-api:8889`
- 静态资源缓存 1 年
- Gzip 压缩已启用

---

### 后端配置

**文件位置**: `/opt/code/dmh/backend/api/etc/dmh-api.docker.yaml`

**主要配置**：
- 数据库: `mysql8:3306`
- Redis: `redis7:6379`
- JWT Secret: `dmh-access-secret-key`
- 日志模式: `file`，路径: `/var/log/dmh-api`
- 频率限制: 使用 Redis 存储

---

## 🌐 访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 管理后台 | http://localhost:3000 | Vue 3 管理界面 |
| H5前端 | http://localhost:3100 | Vue 3 移动端界面 |
| 后端API | http://localhost:8889 | Go 后端API服务 |

---

## 🔒 安全注意事项

1. **生产环境部署前请修改**：
   - JWT Secret（当前：`dmh-access-secret-key`）
   - 数据库密码（当前：`Admin168`）
   - 微信支付配置

2. **网络安全**：
   - 容器在 my-net 网络中，只对内网开放
   - 生产环境建议配置反向代理和HTTPS

3. **日志安全**：
   - 定期清理日志：`docker volume prune`
   - 生产环境建议配置日志收集系统

---

## 📚 相关文档

- `/tmp/dmh-container-deployment-report.md` - 部署完成报告
- `/tmp/docker_migration_guide.md` - 容器化迁移指南
- `/opt/code/dmh/docs/API_Documentation.md` - API文档
- `/opt/code/dmh/docs/Deployment_Checklist.md` - 部署检查清单

---

## 💡 常见问题

### Q: 首次启动为什么这么慢？

A: 首次启动容器时需要安装依赖（ca-certificates, wget, nginx, tzdata 等），需要 2-5 分钟。之后启动会快很多（10-30 秒）。

### Q: 如何查看安装进度？

A: 查看容器日志即可看到安装进度：
```bash
docker logs dmh-nginx | grep apk
docker logs dmh-api | grep apk
```

### Q: 如何更新前端代码？

A: 重新构建前端，然后重启 nginx 容器：
```bash
cd /opt/code/dmh/frontend-admin
npm run build

cd /opt/code/dmh/deployment
docker compose -f docker-compose-simple.yml restart dmh-nginx
```

### Q: 如何更新后端代码？

A: 更新二进制文件和配置，然后重启 api 容器：
```bash
# 更新 /tmp/dmh 二进制文件
# 更新 /tmp/dmh-api.yaml 配置文件

cd /opt/code/dmh/deployment
docker compose -f docker-compose-simple.yml restart dmh-api
```

### Q: 如何扩展服务？

A: 使用 docker compose scale 扩展：
```bash
# 扩展API服务到3个实例
docker compose -f docker-compose-simple.yml up -d --scale dmh-api=3

# 扩展Nginx服务到2个实例
docker compose -f docker-compose-simple.yml up -d --scale dmh-nginx=2
```

---

## 🎯 下一步

1. **启动服务**：`./scripts/quick-start.sh`
2. **验证服务**：访问 http://localhost:3000 和 http://localhost:3100
3. **测试API**：执行登录测试
4. **配置生产环境**：修改密码、JWT密钥等安全配置

---

**部署完成！** 🎉

---

## ⚙️ 配置管理

### 统一配置目录

生产环境配置已迁移到 `/opt/module/dmh/configs/` 目录：

```
/opt/module/dmh/configs/
├── dmh-api.yaml           # 后端 API 配置
├── nginx/conf.d/
│   └── default.conf       # Nginx 反向代理配置
├── frontend/
│   ├── admin.env          # 管理后台环境变量
│   └── h5.env             # H5 前端环境变量
└── backup/                # 配置备份目录
```

### 配置修改流程

```bash
# 1. 修改配置文件
vim /opt/module/dmh/configs/dmh-api.yaml

# 2. 重启服务（自动备份+验证）
cd /opt/code/dmh/deploy/scripts
./restart-services.sh
```

### 可用脚本

| 脚本 | 用途 |
|------|------|
| `sync-configs.sh` | 从项目目录同步配置到统一管理目录 |
| `backup-config.sh` | 备份当前配置 |
| `verify-config.sh` | 验证配置正确性 |
| `restart-services.sh` | 一键重启服务（备份+验证+重启+健康检查） |

### 示例

```bash
cd /opt/code/dmh/deploy/scripts

# 查看备份列表
./backup-config.sh --list

# 恢复最近的备份
./backup-config.sh --restore

# 验证配置
./verify-config.sh

# 完整重启流程
./restart-services.sh
```

### 配置文件详情

详细说明请参阅：`/opt/module/dmh/README.md`

---

**最后更新**: 2026-02-19
