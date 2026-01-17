# DMH 项目启动指南

## 快速开始

```bash
# 首次运行（初始化环境）
./dmh.sh init

# 日常启动
./dmh.sh start

# 停止服务
./dmh.sh stop

# 查看状态
./dmh.sh status

# 查看日志
./dmh.sh logs
```

---

## 环境要求

- **Go**: 1.19+
- **Node.js**: 16+
- **Docker**: 用于运行 MySQL

### 检查环境
```bash
go version
node --version
docker --version
```

---

## 命令说明

| 命令 | 说明 |
|------|------|
| `./dmh.sh init` | 首次运行，创建 MySQL 容器、初始化数据库、安装依赖 |
| `./dmh.sh start` | 日常启动，自动启动 Docker 和 MySQL，启动所有服务 |
| `./dmh.sh stop` | 停止所有服务 |
| `./dmh.sh restart` | 重启所有服务 |
| `./dmh.sh status` | 查看服务运行状态 |
| `./dmh.sh logs` | 查看服务日志 |

---

## 服务地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 后端 API | http://localhost:8889 | RESTful API |
| H5 前端 | http://localhost:3100 | 用户端 |
| 管理后台 | http://localhost:3000 | 管理员后台 |

---

## 测试账号

| 角色 | 用户名 | 密码 | 访问地址 |
|------|--------|------|----------|
| 平台管理员 | admin | 123456 | http://localhost:3000 |
| 品牌管理员 | brand_manager | 123456 | http://localhost:3100/brand/login |

---

## 测试数据

### 初始化测试数据

```bash
# 导入品牌、活动、会员测试数据
docker exec -i mysql8 mysql -u root -p'#Admin168' dmh < backend/scripts/test_data.sql
docker exec -i mysql8 mysql -u root -p'#Admin168' dmh < backend/scripts/seed_member_campaign_data.sql
```

### 测试品牌

| ID | 品牌名称 | 状态 |
|----|----------|------|
| 1 | 星巴克咖啡 | active |
| 2 | 麦当劳 | active |
| 3 | 耐克运动 | active |
| - | 测试品牌C | active |

### 测试活动

| ID | 活动名称 | 品牌 | 状态 |
|----|----------|------|------|
| 1 | 新年促销活动 | 品牌A | active |
| 2 | 春季招生活动 | - | active |
| 3 | 会员专享活动 | - | active |
| - | 会员招募活动-1 | 测试品牌C | active |
| - | 会员招募活动-2 | 测试品牌C | active |

### 测试会员

| unionid | 昵称 | 手机号 | 状态 |
|---------|------|--------|------|
| test_unionid_0001 | 测试会员01 | 13900000001 | active |
| test_unionid_0002 | 测试会员02 | 13900000002 | active |
| test_unionid_0003 | 测试会员03 | 13900000003 | active |
| test_unionid_0004 | 测试会员04 | 13900000004 | active |
| test_unionid_0005 | 测试会员05 | 13900000005 | disabled |
| test_unionid_0006 | 测试会员06 | 13900000006 | active |

---

## 会员入口链接

### H5 用户端

| 页面 | 链接 | 说明 |
|------|------|------|
| 活动列表 | http://localhost:3100/ | 所有活动列表 |
| 活动详情 | http://localhost:3100/campaign/{id} | 活动详情页 |
| 活动报名 | http://localhost:3100/campaign/{id}/form | 报名表单页 |
| 推广链接 | http://localhost:3100/campaign/{id}?ref={推荐人ID} | 带推荐人的活动链接 |

### 品牌管理端 (H5)

| 页面 | 链接 | 说明 |
|------|------|------|
| 登录 | http://localhost:3100/brand/login | 品牌管理员登录 |
| 仪表板 | http://localhost:3100/brand/dashboard | 数据概览 |
| 活动管理 | http://localhost:3100/brand/campaigns | 活动列表 |
| 活动编辑 | http://localhost:3100/brand/campaigns/edit/{id} | 编辑活动 |
| 页面设计 | http://localhost:3100/brand/campaigns/{id}/designer | 可视化页面设计器 |
| 会员管理 | http://localhost:3100/brand/members | 会员列表 |
| 订单管理 | http://localhost:3100/brand/orders | 订单列表 |
| 推广员 | http://localhost:3100/brand/promoters | 推广员管理 |

### 管理后台

| 页面 | 链接 | 说明 |
|------|------|------|
| 登录 | http://localhost:3000 | 平台管理员登录 |
| 用户管理 | http://localhost:3000 → 用户管理 | 管理所有用户 |
| 品牌管理 | http://localhost:3000 → 品牌管理 | 管理品牌 |
| 会员管理 | http://localhost:3000 → 会员管理 | 查看所有会员 |
| 权限管理 | http://localhost:3000 → 角色权限 | 配置角色权限 |

### 示例链接

```bash
# 查看活动 ID=1 的详情
http://localhost:3100/campaign/1

# 活动 ID=1 的报名表单
http://localhost:3100/campaign/1/form

# 带推荐人 ID=2 的推广链接
http://localhost:3100/campaign/1?ref=2
```

---

## 环境配置

### 自定义 MySQL 配置

通过环境变量覆盖默认配置：

```bash
export MYSQL_HOST=172.17.0.1
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASS='#Admin168'
export MYSQL_DB=dmh
export DOCKER_MYSQL_CONTAINER=mysql8

./dmh.sh start
```

### 后端配置文件

`backend/api/etc/dmh-api.yaml`:

```yaml
Mysql:
  DataSource: root:#Admin168@tcp(172.17.0.1:3306)/dmh?charset=utf8mb4&parseTime=true&loc=Local
```

---

## 常见问题

### Docker 未运行

脚本会自动启动 Docker：
- **WSL2**: 自动运行 `sudo dockerd`
- **Ubuntu**: 自动运行 `sudo systemctl start docker`

### 端口被占用

```bash
# 查看占用进程
lsof -i :8889
lsof -i :3000
lsof -i :3100

# 停止服务后重试
./dmh.sh stop
./dmh.sh start
```

### MySQL 连接失败

```bash
# 检查容器状态
docker ps | grep mysql

# 查看容器日志
docker logs mysql8

# 手动启动容器
docker start mysql8
```

### 依赖安装失败

```bash
# Go 依赖
cd backend && go mod tidy && cd ..

# npm 依赖（使用镜像）
npm config set registry https://registry.npmmirror.com
cd frontend-h5 && npm install && cd ..
cd frontend-admin && npm install && cd ..
```

---

## 日志文件

| 服务 | 日志路径 |
|------|----------|
| 后端 | `logs/backend.log` |
| H5 前端 | `logs/h5.log` |
| 管理后台 | `logs/admin.log` |

```bash
# 实时查看日志
tail -f logs/backend.log
```

---

## 手动操作

### 手动启动 Docker MySQL

```bash
# 创建容器
docker run -d \
  --name mysql8 \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD='#Admin168' \
  -e MYSQL_DATABASE=dmh \
  -v mysql_data:/var/lib/mysql \
  mysql:8.0

# 导入数据
docker exec -i mysql8 mysql -u root -p'#Admin168' dmh < backend/scripts/init.sql
```

### 手动启动服务

```bash
# 后端
cd backend && go run api/dmh.go -f api/etc/dmh-api.yaml

# H5 前端
cd frontend-h5 && npm run dev

# 管理后台
cd frontend-admin && npm run dev
```
