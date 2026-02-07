# 添加活动高级功能 - OpenSpec Change

## 📋 变更概述

本变更为 DMH 系统添加四个关键的营销和运营功能：
1. **一键生成海报** - 分销员快速生成推广海报
2. **支付配置** - 活动级别的支付参数配置和二维码生成
3. **表单字段增强** - 支持更多字段类型和验证规则
4. **订单核销** - 品牌管理员扫码核销订单

## 📁 文件结构

```
openspec/changes/add-campaign-advanced-features/
├── README.md                           # 本文件
├── proposal.md                         # 变更提案（Why, What, Impact）
├── design.md                           # 技术设计文档
├── tasks.md                            # 实施任务清单（91 个任务）
└── specs/                              # 增量规范
    ├── campaign-management/
    │   └── spec.md                     # 活动管理模块增量变更
    └── order-payment-system/
        └── spec.md                     # 订单支付系统增量变更
```

## ✅ 验证状态

- [x] OpenSpec 格式验证通过
- [x] 所有需求包含场景（Scenario）
- [x] 增量变更正确标记（ADDED/MODIFIED）
- [x] 等待用户审批

## 🎯 影响的模块

### 1. 活动管理模块 (campaign-management)
**新增需求：**
- 海报生成功能（4 个场景）
- 支付配置功能（4 个场景）
- 表单字段类型扩展（4 个场景）
- 表单字段排序功能（2 个场景）
- 表单实时预览（2 个场景）

**修改需求：**
- 动态表单配置（扩展支持新字段类型）
- 活动创建/编辑（增加支付配置）
- 活动详情展示（增加海报和支付二维码入口）

### 2. 订单支付系统 (order-payment-system)
**新增需求：**
- 订单核销功能（6 个场景）
- 订单核销码生成（3 个场景）
- 订单核销权限控制（3 个场景）
- 订单核销操作日志（3 个场景）

**修改需求：**
- 订单创建（增加核销码生成）
- 订单查询（返回核销状态信息）

## 🗄️ 数据库变更

### campaigns 表
```sql
ALTER TABLE campaigns 
ADD COLUMN payment_config JSON COMMENT '支付配置',
ADD COLUMN poster_template_id INT DEFAULT 1 COMMENT '海报模板ID';
```

### orders 表
```sql
ALTER TABLE orders
ADD COLUMN verification_status VARCHAR(20) DEFAULT 'unverified' COMMENT '核销状态',
ADD COLUMN verified_at DATETIME NULL COMMENT '核销时间',
ADD COLUMN verified_by BIGINT NULL COMMENT '核销人ID',
ADD COLUMN verification_code VARCHAR(50) NULL COMMENT '核销码',
ADD INDEX idx_verification_status (verification_status),
ADD INDEX idx_verified_at (verified_at);
```

### 新增 poster_templates 表
```sql
CREATE TABLE poster_templates (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '模板名称',
    preview_image VARCHAR(255) COMMENT '预览图',
    config JSON NOT NULL COMMENT '模板配置',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='海报模板表';
```

## 🔌 新增 API 接口

1. `POST /api/v1/campaigns/:id/poster` - 生成海报
2. `GET /api/v1/campaigns/:id/payment-qrcode` - 获取支付二维码
3. `GET /api/v1/orders/scan/:code` - 扫码获取订单信息
4. `POST /api/v1/orders/:id/verify` - 核销订单
5. `POST /api/v1/orders/:id/unverify` - 取消核销

## 📦 依赖变更

### 后端（Go）
- `github.com/fogleman/gg` - 图像处理和合成
- `github.com/skip2/go-qrcode` - 二维码生成

### 前端（H5）
- `html5-qrcode` 或 `jsQR` - 二维码扫描
- `sortablejs` - 拖拽排序

## 📝 实施任务

共 **91 个任务**，分为 15 个阶段：

1. 数据库迁移（7 个任务）
2. 后端 - 海报生成（8 个任务）
3. 后端 - 支付配置（7 个任务）
4. 后端 - 表单字段增强（5 个任务）
5. 后端 - 订单核销（8 个任务）
6. 前端 - H5 海报生成页面（8 个任务）
7. 前端 - H5 订单核销页面（8 个任务）
8. 前端 - 活动编辑页面增强（8 个任务）
9. 前端 - 活动详情页面增强（4 个任务）
10. 管理后台开发（3 个任务）
11. 集成测试（6 个任务）
12. 性能测试（4 个任务）
13. 安全测试（4 个任务）
14. 文档和部署（7 个任务）
15. 监控和优化（4 个任务）

详见 [tasks.md](./tasks.md)

## 🎯 性能目标

- 海报生成时间 < 3 秒
- 支付二维码生成时间 < 1 秒
- 订单核销响应时间 < 500ms
- 海报生成成功率 > 95%

## 🔒 安全措施

1. **海报生成**：频率限制（每用户每分钟 5 次）
2. **支付二维码**：包含签名防伪造，有效期 2 小时
3. **订单核销**：核销码包含签名，仅品牌管理员可操作
4. **操作日志**：记录所有核销操作

## 📅 时间线

- **第 1 周**：数据库设计和后端 API 开发
- **第 2 周**：前端页面开发和联调
- **第 3 周**：测试和优化
- **第 4 周**：上线和用户培训

## 🚀 下一步

1. **等待审批**：用户审查 proposal.md 和 design.md
2. **开始实施**：审批通过后，按 tasks.md 顺序执行
3. **测试验证**：完成开发后进行完整测试
4. **部署上线**：测试通过后部署到生产环境
5. **归档变更**：上线后使用 `openspec archive add-campaign-advanced-features`

## 📚 相关文档

- [proposal.md](./proposal.md) - 变更提案
- [design.md](./design.md) - 技术设计
- [tasks.md](./tasks.md) - 任务清单
- [specs/campaign-management/spec.md](./specs/campaign-management/spec.md) - 活动管理增量规范
- [specs/order-payment-system/spec.md](./specs/order-payment-system/spec.md) - 订单支付增量规范

## ✨ 验证命令

```bash
# 验证变更格式
openspec validate add-campaign-advanced-features --strict --no-interactive

# 查看变更详情
openspec show add-campaign-advanced-features

# 查看任务列表
openspec list
```
