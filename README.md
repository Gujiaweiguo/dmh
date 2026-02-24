# 🎯 DMH Digital Marketing Hub

数字营销中台系统，提供完整的营销活动管理、用户权限管理和数据分析功能。

![DMH Logo](https://img.shields.io/badge/DMH-Digital%20Marketing%20Hub-blue)
![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8)
![Vue Version](https://img.shields.io/badge/Vue.js-3.0+-4FC08D)
![License](https://img.shields.io/badge/License-MIT-green)

## ✨ 主要功能

### 🔐 完整的RBAC权限系统

* **4种用户角色**：平台管理员、品牌管理员、参与者、匿名用户
* **JWT认证**：Token自动刷新、登录状态管理、安全验证
* **权限控制**：API级别、页面级别、数据级别的全方位权限管理
* **菜单权限**：动态菜单生成和按钮权限控制

### 🎨 可视化活动管理

* **页面设计器**：拖拽式组件设计，包含8种常用组件
* **动态表单**：自定义报名字段，支持文本、手机号、邮箱、选择等类型
* **实时预览**：所见即所得的页面设计体验
* **主题配置**：颜色、字体、布局完全自定义

### 👥 用户和品牌管理

* **用户管理**：创建、编辑、禁用用户账号
* **品牌管理**：品牌信息管理和品牌管理员关系绑定
* **数据统计**：用户行为分析、活动参与统计

## 🏗️ 技术架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   H5 Frontend   │    │  Admin Frontend │    │   Backend API   │
│   (Vue.js 3)    │    │   (Vue.js 3)    │    │   (go-zero)     │
│   Port: 3100    │    │   Port: 3000    │    │   Port: 8889    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │     MySQL       │
                    │   Database      │
                    └─────────────────┘
```

### 后端技术栈

* **Go 1.24+** - 高性能后端服务
* **go-zero 1.6+** - 微服务框架，提供API网关、RPC、中间件等完整解决方案
* **GORM** - 强大的ORM框架
* **JWT** - 安全的身份认证
* **MySQL 8.0+** - 可靠的数据存储
* **Redis** - 高性能缓存（可选）

### 前端技术栈

* **Vue.js 3** - 现代化的渐进式框架
* **Vite 5/6** - 极速的前端构建工具
* **Vant UI** - 优秀的移动端组件库（H5端）
* **Lucide Vue** - 精美的图标库（管理后台）
* **TypeScript** - 类型安全的开发体验

## 🚀 快速开始

### 容器化部署（推荐）⭐

#### 一键启动

```bash
cd /opt/code/dmh/deploy/scripts
./quick-start.sh
```

首次启动需要 2-5 分钟（安装依赖），之后只需 10-30 秒。

#### 服务访问

* 📱 H5前端：http://localhost:3100
* 💻 管理后台：http://localhost:3000
* 🔧 后端API：http://localhost:8889

#### 详细文档

详细部署说明：[/deploy/README.md](deploy/README.md)

***

### 开发环境（本地运行）

#### 环境要求

* Go 1.24+
* Node.js 20+
* MySQL 8.0+

> 💡 **提示**: 详细的环境安装步骤（Docker、Go、Node.js）请参考 [SETUP.md](./SETUP.md)

#### 1. 克隆项目

```bash
git clone https://github.com/Gujiaweiguo/dmh.git
cd DMH
```

#### 2. 数据库初始化

```bash
# 初始化数据库（使用 Docker 容器内的 MySQL）
docker exec -i mysql8 mysql -uroot -p'Admin168' < backend/scripts/init.sql
```

#### 3. 启动后端服务

```bash
cd backend
go mod download
go run api/dmh.go -f api/etc/dmh-api.yaml
```

后端服务将在 http://localhost:8889 启动

#### 4. 启动H5前端

```bash
cd frontend-h5
npm install
npm run dev
```

H5前端将在 http://localhost:3100 启动

#### 5. 启动管理后台

```bash
cd frontend-admin
npm install
npm run dev
```

管理后台将在 http://localhost:3000 启动

## 🌐 访问地址和测试账号

| 服务 | 地址 | 用户类型 | 用户名 | 密码 | 功能描述 |
|------|------|----------|--------|------|----------|
| H5前端 | http://localhost:3100 | 普通用户 | - | - | 浏览活动、参与报名 |
| H5前端 | http://localhost:3100/brand/login | 品牌管理员 | brand\_manager | 123456 | 活动管理、页面设计 |
| 管理后台 | http://localhost:3000 | 平台管理员 | admin | 123456 | 系统管理、用户管理 |
| 后端API | http://localhost:8889 | - | - | - | RESTful API服务 |

## 📱 用户角色详解

### 🔧 平台管理员 (platform\_admin)

**访问方式**: http://localhost:3000

* ✅ 用户账号管理（创建、编辑、禁用、重置密码）
* ✅ 品牌信息管理（创建、编辑品牌）
* ✅ 权限配置管理（角色权限、菜单权限）
* ✅ 系统设置和全局数据查看
* ✅ 活动管理（查看、编辑所有品牌的活动）

### 🏢 品牌管理员 (brand\_admin)

**访问方式**: http://localhost:3100/brand/login

* ✅ 营销活动管理（创建、编辑、启用、暂停）
* ✅ 可视化页面设计器（8种组件、主题配置）
* ✅ 动态表单设计（自定义字段、验证规则）
* ✅ 活动数据分析（参与统计、转化率）
* ✅ 报名信息管理（查看、导出报名数据）
* ✅ 素材管理（上传、管理活动素材）

### 👤 普通用户 (participant)

**访问方式**: http://localhost:3100

* ✅ 浏览活动列表（无需登录）
* ✅ 查看活动详情和页面
* ✅ 填写报名表单参与活动
* ✅ 查看个人报名记录
* ✅ 活动筛选和搜索

## 🎨 核心功能展示

### 可视化页面设计器

```
📦 组件库                    ⚙️ 组件配置                   👁️ 实时预览
├── 🖼️ 横幅图片              ├── 图片URL设置               ├── 页面标题
├── 📝 文本内容              ├── 文本内容编辑             ├── 活动描述  
├── 🎬 视频播放              ├── 字体大小/对齐            ├── 组件预览
├── ⏰ 倒计时                ├── 视频URL配置               ├── 表单字段
├── 💬 用户评价              ├── 倒计时设置               ├── 报名按钮
├── ❓ 常见问题              ├── 评价内容管理             └── 实时更新
├── 📞 联系方式              ├── 问答列表编辑
└── 🔗 社交媒体              └── 联系信息配置
```

### RBAC权限体系

```
用户 (User)
├── 拥有角色 (Has Roles)
│   ├── platform_admin (平台管理员)
│   ├── brand_admin (品牌管理员)
│   ├── participant (参与者)
│   └── anonymous (匿名用户)
│
├── 角色权限 (Role Permissions)
│   ├── 资源权限 (Resource Permissions)
│   ├── 操作权限 (Action Permissions)
│   └── 数据权限 (Data Permissions)
│
└── 菜单权限 (Menu Permissions)
    ├── 页面访问权限
    ├── 按钮操作权限
    └── 功能模块权限
```

## 📊 数据库设计

### 核心数据表

* **users** - 用户基础信息
* **roles** - 角色定义
* **permissions** - 权限定义
* **user\_roles** - 用户角色关系
* **role\_permissions** - 角色权限关系
* **menus** - 菜单定义
* **role\_menus** - 角色菜单权限
* **brands** - 品牌信息
* **campaigns** - 营销活动
* **orders** - 报名订单
* **audit\_logs** - 操作审计日志

## 🛠️ 开发指南

### API文档

* **接口定义**: `backend/api/dmh.api`
* **API文档**: `docs/api/RBAC-API-Documentation.md`

### 前端开发

```bash
# H5前端开发
cd frontend-h5
npm run dev

# 管理后台开发
cd frontend-admin  
npm run dev

# 构建生产版本
npm run build
```

### 后端开发

```bash
# 运行开发服务器
cd backend
go run api/dmh.go -f api/etc/dmh-api.yaml

# 运行测试
go test ./...

# 构建生产版本
go build -o dmh-api api/dmh.go
```

### Order MySQL8 回归测试

用于验证订单核销鉴权与重复报名文案回归。

```bash
# 可选：覆盖 API 地址和测试管理员账号
DMH_INTEGRATION_BASE_URL=http://127.0.0.1:8889 \
DMH_TEST_ADMIN_USERNAME=admin \
DMH_TEST_ADMIN_PASSWORD=123456 \
backend/scripts/run_order_mysql8_regression.sh
```

* 本地脚本：`backend/scripts/run_order_mysql8_regression.sh`
* 一键排障+回归：`backend/scripts/repair_login_and_run_order_regression.sh`
* CI 工作流：`.github/workflows/order-mysql8-regression.yml`
* 若用例被跳过，脚本会返回非 0，避免假阳性。

常见失败排查：

* `login returned 400`：先确认管理员账号存在且密码正确（默认见 `backend/scripts/init.sql`）。
  * 可直接执行：`backend/scripts/repair_login_and_run_order_regression.sh`
  * 该脚本会在本地环境自动修复 `admin`/`brand_manager` 的默认测试密码哈希，并重启 `dmh-api` 后重跑回归。
* `Unknown column ...`：检查数据库是否执行了订单/活动相关 migration。
* `field "code" is not set`：确认核销接口请求格式与当前后端解析一致。

## 🤝 贡献指南

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范

* Go代码遵循 `gofmt` 格式
* 前端代码使用 ESLint + Prettier
* 提交信息使用 [Conventional Commits](https://conventionalcommits.org/)

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

## 📞 联系方式

* **GitHub Issues**: [提交问题](https://github.com/Gujiaweiguo/dmh/issues)
* **Email**: weiguogu@163.com

***

⭐ **如果这个项目对你有帮助，请给个星标支持！**

🔗 **项目链接**: https://github.com/Gujiaweiguo/dmh
