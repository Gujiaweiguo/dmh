# 🛠️ DMH 开发指南

## 目录

* [开发环境搭建](#开发环境搭建)
* [项目结构](#项目结构)
* [开发规范](#开发规范)
* [调试指南](#调试指南)
* [测试指南](#测试指南)
* [常见问题](#常见问题)

***

## 开发环境搭建

### 环境要求

* **Go**: 1.23+
* **Node.js**: 20.19.0+ (建议使用 nvm)
* **MySQL**: 8.0+
* **Git**: 2.0+

> 💡 **提示**: 详细的环境安装步骤（Docker、Go、Node.js）请参考 [SETUP.md](./SETUP.md)

***

### 快速开始

#### 方式一：使用启动脚本（推荐）⭐

```bash
# 克隆项目
git clone https://github.com/Gujiaweiguo/DMH.git
cd DMH

# 一键初始化和启动
./dmh.sh init   # 首次运行（会自动安装 MySQL 容器并初始化数据库）
./dmh.sh start  # 启动所有服务
```

服务启动后：

* 后端 API: http://localhost:8889
* 管理后台: http://localhost:3000
* H5 端: http://localhost:3100

#### 方式二：手动启动

如果需要单独启动某个服务或自定义配置：

**1. 环境准备**

如果还没有安装环境，请参考 [SETUP.md](./SETUP.md) 安装：

* Docker（用于 MySQL）
* Go 1.23+
* Node.js 20+

**2. 初始化数据库**

```bash
# 使用脚本（推荐）
./dmh.sh init

# 或手动启动 MySQL 容器
docker run -d \
  --name mysql8 \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD='#Admin168' \
  -e MYSQL_DATABASE=dmh \
  mysql:8.0

# 导入初始化脚本
docker exec -i mysql8 mysql -uroot -p'#Admin168' < backend/scripts/init.sql
```

**3. 启动后端**

```bash
cd backend
go mod download
go run api/dmh.go -f api/etc/dmh-api.yaml
```

后端服务将在 http://localhost:8889 启动

***

### 容器化开发环境 ⭐

#### 一键启动（推荐）

```bash
cd /opt/code/DMH/deployment/scripts
./quick-start.sh
```

服务启动后访问：

* 📱 H5前端：http://localhost:3100
* 💻 管理后台：http://localhost:3000
* 🔧 后端API：http://localhost:8889

#### 容器内调试

**进入 API 容器**：

```bash
docker exec -it dmh-api sh
```

**查看 API 日志**：

```bash
docker logs -f dmh-api
```

**进入 Nginx 容器**：

```bash
docker exec -it dmh-nginx sh
```

**查看 Nginx 日志**：

```bash
docker logs -f dmh-nginx
```

#### 容器管理命令

**查看容器状态**：

```bash
cd /opt/code/DMH/deployment
docker compose -f docker-compose-simple.yml ps
```

**重启容器**：

```bash
# 重启所有服务
docker compose -f docker-compose-simple.yml restart

# 重启单个容器
docker restart dmh-api
docker restart dmh-nginx
```

**查看日志**：

```bash
# 所有服务
docker compose -f docker-compose-simple.yml logs -f

# 单个服务
docker logs -f dmh-api
docker logs -f dmh-nginx
```

**详细部署文档**：[/deployment/README.md](../deployment/README.md)

***

### 生产环境手动部署

如果需要单独启动某个服务或自定义配置：

**1. 环境准备**

如果还没有安装环境，请参考 [SETUP.md](./SETUP.md) 安装：

* Docker（用于 MySQL）
* Go 1.23+
* Node.js 20+

**2. 初始化数据库**

```bash
# 使用脚本（推荐）
./dmh.sh init

# 或手动启动 MySQL 容器
docker run -d \
  --name mysql8 \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD='#Admin168' \
  -e MYSQL_DATABASE=dmh \
  mysql:8.0

# 导入初始化脚本
docker exec -i mysql8 mysql -uroot -p'#Admin168' < backend/scripts/init.sql
```

**3. 启动后端**

```bash
cd backend
go mod download
go run api/dmh.go -f api/etc/dmh-api.yaml
```

**4. 启动前端**

管理后台：

```bash
cd frontend-admin
npm install
npm run dev
```

H5 端：

```bash
cd frontend-h5
npm install
npm run dev
```

***

### 测试账号

| 用户名 | 密码 | 角色 | 访问地址 |
|--------|------|------|----------|
| admin | 123456 | 平台管理员 | http://localhost:3000 |
| brand\_manager | 123456 | 品牌管理员 | http://localhost:3100/brand/login |
| user001 | 123456 | 普通用户 | http://localhost:3100 |

***

## 项目结构

```
DMH/
├── backend/                    # 后端服务（Go）
│   ├── api/                   # API 定义和入口
│   │   ├── dmh.api           # go-zero API 定义
│   │   ├── dmh.go            # 主入口文件
│   │   ├── etc/              # 配置文件
│   │   └── internal/         # 内部实现
│   │       ├── config/       # 配置结构
│   │       ├── handler/      # HTTP 处理器
│   │       ├── logic/        # 业务逻辑
│   │       ├── middleware/   # 中间件
│   │       ├── svc/          # 服务上下文
│   │       └── types/        # 类型定义
│   ├── common/               # 公共模块
│   │   ├── syncadapter/     # 数据同步适配器
│   │   └── utils/           # 工具函数
│   ├── model/                # 数据模型
│   ├── migrations/           # 数据库迁移
│   ├── scripts/              # 脚本文件
│   ├── test/                 # 测试文件
│   └── storage/              # 文件存储
│
├── frontend-admin/            # 管理后台（Vue 3 + Vite 6）
│   ├── components/           # 组件
│   ├── views/                # 页面视图
│   ├── services/             # API 服务
│   ├── styles/               # 样式文件
│   └── types.ts              # TypeScript 类型
│
├── frontend-h5/               # H5 端（Vue 3 + Vite 5）
│   ├── src/
│   │   ├── components/       # 组件
│   │   ├── views/            # 页面视图
│   │   ├── services/         # API 服务
│   │   ├── router/           # 路由配置
│   │   └── utils/            # 工具函数
│   └── public/               # 静态资源
│
├── docs/                      # 文档
│   ├── api/                  # API 文档
│   ├── deployment/           # 部署文档
│   └── user-manual/          # 用户手册
│
├── openspec/                  # OpenSpec 规范
│   ├── specs/                # 功能规格
│   └── changes/              # 变更提案
│
├── logs/                      # 日志目录
├── .opencode/                # OpenCode 配置
├── README.md                 # 项目说明
├── ARCHITECTURE.md           # 架构文档
├── API.md                    # API 文档
├── DEVELOPMENT.md            # 开发指南（本文件）
└── dmh.sh                    # 启动脚本
```

***

## 开发规范

### Go 代码规范

#### 1. 命名规范

```go
// 包名：小写，简短
package handler

// 接口名：名词，首字母大写
type UserService interface {
    GetUser(id int64) (*User, error)
}

// 结构体：首字母大写
type User struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
}

// 函数名：动词开头，驼峰命名
func CreateUser(req *CreateUserRequest) error {
    // ...
}

// 私有函数：首字母小写
func validateUser(user *User) error {
    // ...
}
```

#### 2. 错误处理

```go
// 推荐：明确的错误处理
user, err := userService.GetUser(id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// 避免：忽略错误
user, _ := userService.GetUser(id)
```

#### 3. 注释规范

```go
// CreateUser 创建新用户
// 参数:
//   - req: 创建用户请求
// 返回:
//   - *User: 创建的用户对象
//   - error: 错误信息
func CreateUser(req *CreateUserRequest) (*User, error) {
    // ...
}
```

#### 4. 代码格式化

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 使用 golangci-lint（推荐）
golangci-lint run
```

### 前端代码规范

#### 1. Vue 组件规范

```vue
<template>
  <!-- 使用语义化的 HTML 标签 -->
  <div class="user-list">
    <h1>用户列表</h1>
    <ul>
      <li v-for="user in users" :key="user.id">
        {{ user.username }}
      </li>
    </ul>
  </div>
</template>

<script>
// 使用 Composition API
import { ref, onMounted } from 'vue';
import { getUserList } from '@/services/userApi';

export default {
  name: 'UserList',
  setup() {
    const users = ref([]);

    const loadUsers = async () => {
      try {
        const data = await getUserList();
        users.value = data;
      } catch (error) {
        console.error('Failed to load users:', error);
      }
    };

    onMounted(() => {
      loadUsers();
    });

    return {
      users,
      loadUsers
    };
  }
};
</script>

<style scoped>
.user-list {
  padding: 20px;
}
</style>
```

#### 2. TypeScript 类型定义

```typescript
// types.ts
export interface User {
  id: number;
  username: string;
  email: string;
  phone: string;
  status: 'active' | 'disabled';
  createdAt: string;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}
```

#### 3. API 服务封装

```typescript
// services/userApi.ts
import axios from 'axios';
import type { User, ApiResponse } from '@/types';

const API_BASE = 'http://localhost:8889/api/v1';

export const getUserList = async (): Promise<User[]> => {
  const response = await axios.get<ApiResponse<User[]>>(
    `${API_BASE}/admin/users`
  );
  return response.data.data;
};

export const createUser = async (user: Partial<User>): Promise<User> => {
  const response = await axios.post<ApiResponse<User>>(
    `${API_BASE}/admin/users`,
    user
  );
  return response.data.data;
};
```

### Git 提交规范

使用 [Conventional Commits](https://conventionalcommits.org/) 规范：

```bash
# 功能开发
git commit -m "feat: 添加用户管理功能"

# Bug 修复
git commit -m "fix: 修复登录失败的问题"

# 文档更新
git commit -m "docs: 更新 API 文档"

# 代码重构
git commit -m "refactor: 重构用户服务代码"

# 性能优化
git commit -m "perf: 优化数据库查询性能"

# 测试相关
git commit -m "test: 添加用户服务单元测试"

# 构建相关
git commit -m "build: 更新依赖版本"

# CI/CD 相关
git commit -m "ci: 添加 GitHub Actions 配置"
```

***

## 调试指南

### 后端调试

#### 1. 使用 Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
cd backend
dlv debug api/dmh.go -- -f api/etc/dmh-api.yaml

# 设置断点
(dlv) break handler.LoginHandler
(dlv) continue
```

#### 2. 日志调试

```go
import "github.com/zeromicro/go-zero/core/logx"

// 在代码中添加日志
logx.Infof("User login: %s", username)
logx.Errorf("Failed to create user: %v", err)
```

#### 3. 查看日志

```bash
# 实时查看日志
tail -f logs/backend.log

# 查看错误日志
grep "ERROR" logs/backend.log
```

### 前端调试

#### 1. 浏览器开发者工具

* **F12** 打开开发者工具
* **Console** 查看日志和错误
* **Network** 查看网络请求
* **Vue DevTools** 查看组件状态

#### 2. 添加调试日志

```javascript
console.log('User data:', user);
console.error('API error:', error);
console.table(users); // 表格形式显示数组
```

#### 3. 使用 debugger

```javascript
const loadUsers = async () => {
  debugger; // 代码会在这里暂停
  const data = await getUserList();
  users.value = data;
};
```

***

## 测试指南

### 后端测试

#### 1. 单元测试

```go
// handler/user_test.go
package handler

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    // 准备测试数据
    req := &CreateUserRequest{
        Username: "testuser",
        Password: "123456",
    }

    // 执行测试
    user, err := CreateUser(req)

    // 断言结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "testuser", user.Username)
}
```

#### 2. 运行测试

```bash
# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./backend/api/internal/handler

# 查看测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### 3. 集成测试

```bash
# 运行集成测试
cd backend/test/integration
go test -v
```

### 前端测试

#### 1. 单元测试（计划中）

```javascript
// 使用 Vitest
import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import UserList from '@/components/UserList.vue';

describe('UserList', () => {
  it('renders user list', () => {
    const wrapper = mount(UserList);
    expect(wrapper.find('.user-list').exists()).toBe(true);
  });
});
```

#### 2. E2E 测试（计划中）

```javascript
// 使用 Playwright
import { test, expect } from '@playwright/test';

test('user can login', async ({ page }) => {
  await page.goto('http://localhost:3000');
  await page.fill('input[name="username"]', 'admin');
  await page.fill('input[name="password"]', '123456');
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL('http://localhost:3000/dashboard');
});
```

***

## 常见问题

### 1. 数据库连接失败

**问题**: `Error 1045: Access denied for user 'root'@'localhost'`

**解决方案**:

```bash
# 检查 MySQL 容器是否启动
docker ps | grep mysql8

# 检查配置文件中的数据库密码
cat backend/api/etc/dmh-api.yaml

# 重新初始化数据库
./dmh.sh init
```

详细的数据库配置请参考 [SETUP.md](./SETUP.md)

### 2. 前端启动失败

**问题**: `Error: Cannot find module 'vue'`

**解决方案**:

```bash
# 删除 node_modules 和 lock 文件
rm -rf node_modules package-lock.json

# 重新安装依赖
npm install

# 或使用 npm ci（推荐）
npm ci
```

### 3. Go 依赖下载慢

**问题**: `go: downloading ... timeout`

**解决方案**:

```bash
# 设置 Go 代理（中国大陆）
go env -w GOPROXY=https://goproxy.cn,direct

# 重新下载依赖
go mod download
```

详细的镜像源配置请参考 [SETUP.md](./SETUP.md)

### 4. 端口被占用

**问题**: `bind: address already in use`

**解决方案**:

```bash
# 查找占用端口的进程
lsof -i :8889  # 后端端口
lsof -i :3000  # 管理后台端口
lsof -i :3100  # H5 端口

# 使用脚本停止服务
./dmh.sh stop

# 或手动杀死进程
kill -9 <PID>
```

***

## 开发工具推荐

### IDE

* **GoLand** - Go 开发（推荐）
* **VS Code** - 通用开发
  * 插件: Go, Vue, ESLint, Prettier

### 数据库工具

* **DBeaver** - 免费开源
* **Navicat** - 商业软件
* **MySQL Workbench** - 官方工具

### API 测试

* **Postman** - API 测试
* **Insomnia** - 轻量级 API 测试
* **curl** - 命令行工具

### 版本控制

* **Git** - 版本控制
* **GitHub Desktop** - Git GUI
* **SourceTree** - Git GUI

***

## 相关文档

* [README.md](./README.md) - 项目介绍
* [SETUP.md](./SETUP.md) - 环境搭建指南
* [ARCHITECTURE.md](./ARCHITECTURE.md) - 系统架构
* [API.md](./API.md) - API 文档
* [CONTRIBUTING.md](./CONTRIBUTING.md) - 贡献指南

***

**文档版本**: v1.0\
**最后更新**: 2025-01-21\
**维护者**: DMH Team
