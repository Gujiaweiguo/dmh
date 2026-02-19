# 🚀 DMH 环境搭建指南

> **文档定位**：
>
> * 本文档提供完整的环境安装和部署步骤
> * 适用于首次部署、生产环境部署
> * 如果您是开发人员，建议先阅读 [DEVELOPMENT.md](./DEVELOPMENT.md)

> **快速开始**：
>
> * **开发环境**：推荐使用 `./dmh.sh` 脚本（简单快速）
>
> * **生产环境**：推荐使用 Docker Compose + Nginx（见本文档第十二节）

## 技术栈

### 后端

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.23+ | 编程语言（以 `backend/go.mod` 为准，当前为 `1.23.0`） |
| go-zero | 1.6.0 | 微服务框架 |
| GORM | 1.25.5 | ORM 框架 |
| JWT | - | 身份认证 |
| bcrypt | - | 密码加密 |

### 前端

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.x | 前端框架 |
| Vite | 5.x/6.x | 构建工具 |
| Element Plus | - | UI 组件库 (Admin) |
| Vant | - | UI 组件库 (H5) |
| Pinia | - | 状态管理 |
| Vue Router | 4.x | 路由管理 |
| Axios | - | HTTP 客户端 |

### 基础设施

| 组件 | 版本 | 说明 |
|------|------|------|
| Docker | 20.10+ | 容器运行环境 |
| MySQL | 8.0 | 关系型数据库 |
| Node.js | 20.x | 前端运行环境 |
| npm | 10.x | 包管理器 |

## 一、配置国内镜像源（重要）

### 1.1 Docker 镜像加速

```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<'EOF'
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://docker.1ms.run",
    "https://docker.xuanyuan.me",
    "https://hub.rat.dev"
  ]
}
EOF

# 重启 Docker
sudo systemctl restart docker
# WSL2 中使用
sudo service docker restart

# 验证镜像加速是否生效（能看到 Registry Mirrors）
sudo docker info | grep -i -n mirror || true
```

### 1.2 Go 模块代理

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
```

### 1.3 npm 淘宝镜像

```bash
npm config set registry https://registry.npmmirror.com
```

## 二、安装 Docker

> 说明：若 `curl -fsSL https://get.docker.com | sh` 在国内网络出现 `Connection reset by peer`，请直接使用下面的 APT 安装方式（推荐，稳定）。

### Ubuntu 22.04/20.04（推荐：腾讯云源安装 Docker CE）

```bash
sudo apt update
sudo apt install -y ca-certificates curl gnupg

sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://mirrors.cloud.tencent.com/docker-ce/linux/ubuntu/gpg | sudo gpg --dearmor --yes -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

. /etc/os-release
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://mirrors.cloud.tencent.com/docker-ce/linux/ubuntu ${VERSION_CODENAME} stable" | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

sudo systemctl enable --now docker

# 允许非 root 使用 docker（需要重新登录或 newgrp 生效）
sudo usermod -aG docker $USER
newgrp docker

docker --version
docker compose version
```

### Ubuntu（兜底：系统源 docker.io）

```bash
sudo apt update
sudo apt install -y docker.io docker-compose-plugin
sudo systemctl enable --now docker

sudo usermod -aG docker $USER
newgrp docker

docker --version
docker compose version
```

## 三、安装 Go

> 目录约定（固定）：
>
> * 下载目录：`/opt/software`
> * 安装目录：`/opt/module/go`

> 版本要求：以 `backend/go.mod` 为准（当前 `go 1.23.0`）。

```bash
# 1) 准备目录
sudo mkdir -p /opt/software /opt/module
sudo chown -R $USER:$USER /opt/software /opt/module

# 2) 严格对齐 go.mod：Go 1.23.0
GO_VERSION=1.23.0
arch=$(uname -m)
case "$arch" in
  x86_64) goarch=amd64 ;;
  aarch64|arm64) goarch=arm64 ;;
  *) echo "unsupported arch: $arch"; exit 1 ;;
esac

# 3) 下载（建议国内入口）
cd /opt/software
curl -fLO "https://golang.google.cn/dl/go${GO_VERSION}.linux-${goarch}.tar.gz"

# 4) 校验（避免 .sha256 页面重定向导致校验失败：用 JSON 获取 sha256）
expected=$(curl -fsSL "https://golang.google.cn/dl/?mode=json&include=all" | \
  python3 -c "import sys,json; data=json.load(sys.stdin); \
  f=[x for r in data if r.get('version')=='go${GO_VERSION}' for x in r.get('files',[]) \
  if x.get('os')=='linux' and x.get('arch')=='${goarch}' and x.get('kind')=='archive'][0]; \
  print(f['sha256'])")
echo "${expected}  go${GO_VERSION}.linux-${goarch}.tar.gz" | sha256sum -c -

# 5) 安装到 /opt/module/go
sudo rm -rf /opt/module/go
sudo tar -C /opt/module -xzf "/opt/software/go${GO_VERSION}.linux-${goarch}.tar.gz"

# 验证
/opt/module/go/bin/go version

# 6) 配置环境变量（当前用户）
grep -q '/opt/module/go/bin' ~/.bashrc || cat >> ~/.bashrc <<'EOF'
export PATH=/opt/module/go/bin:$PATH
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH
EOF
source ~/.bashrc
hash -r
go version
```

## 四、安装 Node.js

> 推荐 Node.js 20+。若 `nvm`/GitHub 在国内网络下载慢，可改用腾讯云 Node 镜像二进制安装。

```bash
# 使用 nvm 安装（推荐）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc

# 安装 Node.js 20
nvm install 20
nvm use 20

# 配置淘宝镜像
npm config set registry https://registry.npmmirror.com

# 验证
node -v
npm -v
```

### Node.js 20（腾讯云二进制镜像，兜底）

```bash
sudo apt update
sudo apt install -y xz-utils

NODE_VERSION=20.11.1
arch=$(uname -m)
case "$arch" in
  x86_64) nodearch=x64 ;;
  aarch64|arm64) nodearch=arm64 ;;
  *) echo "unsupported arch: $arch"; exit 1 ;;
esac

cd /opt/software
curl -fLO "https://mirrors.cloud.tencent.com/nodejs-release/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-${nodearch}.tar.xz"
sudo mkdir -p /usr/local/lib/nodejs
sudo tar -xJf "node-v${NODE_VERSION}-linux-${nodearch}.tar.xz" -C /usr/local/lib/nodejs

grep -q '/usr/local/lib/nodejs' ~/.bashrc || cat >> ~/.bashrc <<EOF
export PATH=/usr/local/lib/nodejs/node-v${NODE_VERSION}-linux-${nodearch}/bin:\$PATH
EOF
source ~/.bashrc

npm config set registry https://registry.npmmirror.com
node -v
npm -v
```

## 五、启动 MySQL 8

> 如果 `docker pull mysql:8.0` 连接 `registry-1.docker.io` 超时：
>
> 1. 先确认 Docker 镜像加速已生效（见 1.1）；2) 或直接使用腾讯云镜像：`mirror.ccs.tencentyun.com/library/mysql:8.0`。

```bash
# 创建数据目录
sudo mkdir -p /opt/data/mysql

# 启动 MySQL 容器
sudo docker run -d \
  --name mysql8 \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD='Admin168' \
  -e MYSQL_DATABASE=dmh \
  -v /opt/data/mysql:/var/lib/mysql \
  --restart unless-stopped \
  mysql:8.0 \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci

# 等待启动（约 30 秒）
sleep 30
```

### MySQL 连接信息

| 参数 | 值 |
|------|------|
| Host | `127.0.0.1` 或 `172.17.0.1` (Docker 网关) |
| Port | `3306` |
| User | `root` |
| Password | `Admin168` |
| Database | `dmh` |

## 六、初始化数据库

### SQL 脚本说明

| 文件 | 说明 | 必需 |
|------|------|------|
| `init.sql` | 主初始化（表结构+基础数据+测试用户） | ✅ 是 |
| `create_member_system_tables.sql` | 会员系统表结构 | ✅ 是 |
| `test_data.sql` | 品牌/活动测试数据 | 开发推荐 |
| `seed_member_campaign_data.sql` | 会员活动测试数据 | 开发推荐 |

### 执行初始化

```bash
# 1. 主初始化脚本（必需）
sudo docker exec -i mysql8 mysql -uroot -p'Admin168' \
  --default-character-set=utf8mb4 < backend/scripts/init.sql

# 2. 会员系统表（必需）
sudo docker exec -i mysql8 mysql -uroot -p'Admin168' \
  --default-character-set=utf8mb4 dmh < backend/scripts/create_member_system_tables.sql

# 3. 测试数据（开发环境推荐）
sudo docker exec -i mysql8 mysql -uroot -p'Admin168' \
  --default-character-set=utf8mb4 dmh < backend/scripts/test_data.sql

sudo docker exec -i mysql8 mysql -uroot -p'Admin168' \
  --default-character-set=utf8mb4 dmh < backend/scripts/seed_member_campaign_data.sql
```

## 七、部署后端

```bash
cd backend

# 下载依赖（go-zero, gorm 等自动下载）
go mod download

# 编译
go build -o dmh-api api/dmh.go

# 修改配置（如需要）
# vim api/etc/dmh-api.yaml

# 运行
./dmh-api -f api/etc/dmh-api.yaml
```

后端运行在 `http://localhost:8889`

## 八、部署前端

### 管理后台 (Vue 3 + Element Plus)

```bash
cd frontend-admin
npm install
npm run dev      # 开发模式
# npm run build  # 生产构建
```

运行在 `http://localhost:3000`

### H5 前端 (Vue 3 + Vant)

```bash
cd frontend-h5
npm install
npm run dev      # 开发模式
# npm run build  # 生产构建
```

运行在 `http://localhost:3100`

## 九、测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | 123456 | 平台管理员 |
| brand\_manager | 123456 | 品牌管理员 |
| user001 | 123456 | 普通用户 |

## 十、快速启动脚本

```bash
#!/bin/bash
# start.sh

# 启动 Docker
sudo service docker start

# 启动 MySQL
sudo docker start mysql8

# 等待 MySQL
sleep 5

# 启动后端
cd backend && ./dmh-api -f api/etc/dmh-api.yaml &

# 启动前端
cd frontend-admin && npm run dev &
cd frontend-h5 && npm run dev &

echo "服务已启动"
echo "  后端: http://localhost:8889"
echo "  管理后台: http://localhost:3000"
echo "  H5: http://localhost:3100"
```

## 十一、常见问题

### 1) Docker 安装/更新时报 NO\_PUBKEY 或源未签名

通常是 GPG key 未正确导入或 `docker.list` 写错。建议按本文「二、安装 Docker（腾讯云源）」整段重做，并确保 `/etc/apt/keyrings/docker.gpg` 存在。

### 2) docker pull 连接 registry-1.docker.io 超时

```bash
sudo cat /etc/docker/daemon.json
sudo systemctl restart docker
sudo docker info | grep -i -n mirror || true

# 也可直接拉腾讯云镜像
docker pull mirror.ccs.tencentyun.com/library/mysql:8.0
```

### 3) 登录提示“用户名或密码错误”

优先检查数据库是否初始化成功（必须有 `users` 表和测试账号数据）。

```bash
docker exec -i mysql8 mysql -uroot -p'Admin168' dmh -e "SHOW TABLES LIKE 'users'; SELECT username,role,status FROM users;"
```

### MySQL 连接失败

````bash
# WSL2 中需使用 Docker 网关 IP
# 修改 backend/api/etc/dmh-api.yaml
Mysql:
  DataSource: root:Admin168@tcp(172.17.0.1:3306)/dmh?charset=utf8mb4&parseTime=true&loc=Local

## 生产环境（Docker Compose，推荐）⭐

生产环境推荐使用容器化部署，提供更好的隔离性和可维护性。

### 快速部署

```bash
# 1) 在服务器上准备 docker + docker compose
sudo apt install -y docker.io docker-compose-plugin

# 2) 一键部署（推荐）
cd /opt/code/dmh/deploy/scripts
./quick-start.sh
````

首次启动需要 2-5 分钟（安装依赖），之后只需 10-30 秒。

### 部署后访问

部署成功后：

* 管理后台：`http://<server>/`
* H5前端：`http://<server>/h5/`
* API：`http://<server>/api/v1/...`

### 详细文档

完整的部署说明、故障排查、回滚操作请参考：[deploy/README.md](../deploy/README.md)

### 配置管理

生产环境请不要把密码写死在配置文件里，使用环境变量或 `.env` 管理：

* 修改 docker-compose.yml 中的环境变量
* 或使用 `.env` 文件管理敏感信息

### 回滚方案

如需回滚到独立进程部署方式：

```bash
cd /opt/code/dmh/deploy/scripts
./rollback-containers.sh
```

### 容器管理

常用管理命令：

* 查看状态：`docker compose ps`
* 查看日志：`docker compose logs -f`
* 重启服务：`docker compose restart`
* 停止服务：`docker compose stop`
* 启动服务：`docker compose start`

### Go 编译失败

```bash
go mod tidy
go mod download
```

### npm 安装慢

```bash
npm config set registry https://registry.npmmirror.com
```

### 端口 8889 被占用

```bash
lsof -i :8889
./dmh.sh stop
./dmh.sh start
```
