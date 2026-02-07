# Design: 活动高级功能技术设计

## Context

当前系统已经实现了基础的活动管理和订单支付功能，但缺少关键的营销工具和运营功能。本次设计旨在补齐这些功能，提升系统的商业化能力和用户体验。

## Goals

1. 提供一键生成海报功能，降低分销员的推广门槛
2. 支持灵活的支付配置，满足不同的商业模式
3. 增强表单字段配置能力，适应多样化的业务场景
4. 实现订单核销功能，打通线上线下闭环

## Non-Goals

1. 不支持自定义海报模板（MVP 阶段使用预设模板）
2. 不支持多种支付方式（MVP 阶段仅支持微信支付）
3. 不支持批量核销（MVP 阶段仅支持单个核销）

## Technical Decisions

### 1. 海报生成方案

**选择：服务端生成 + Canvas 渲染**

**方案对比：**

| 方案 | 优点 | 缺点 | 选择 |
|------|------|------|------|
| 纯前端生成 | 无服务器压力，响应快 | 模板管理困难，样式不统一 | ❌ |
| 服务端生成 | 模板统一管理，质量可控 | 服务器压力大，响应慢 | ✅ |
| 第三方服务 | 功能强大，维护成本低 | 依赖外部服务，成本高 | ❌ |

**实现细节：**
- 使用 Go 的 `gg` 库进行图片合成
- 使用 `qrcode` 库生成二维码
- 海报模板配置存储在数据库中（JSON 格式）
- 生成的海报缓存到 OSS/本地存储

**海报模板数据结构：**
```json
{
  "id": 1,
  "name": "经典模板",
  "width": 750,
  "height": 1334,
  "background": "https://example.com/bg.jpg",
  "elements": [
    {
      "type": "text",
      "content": "{{campaignName}}",
      "x": 50,
      "y": 100,
      "fontSize": 32,
      "color": "#333333",
      "fontWeight": "bold"
    },
    {
      "type": "qrcode",
      "content": "{{distributorLink}}",
      "x": 300,
      "y": 1000,
      "size": 200
    },
    {
      "type": "image",
      "content": "{{campaignImage}}",
      "x": 0,
      "y": 200,
      "width": 750,
      "height": 400
    }
  ]
}
```

### 2. 支付配置方案

**选择：活动级别配置 + 订单级别记录**

**数据模型：**
```go
type PaymentConfig struct {
    DepositAmount  float64 `json:"depositAmount"`  // 订金金额
    FullAmount     float64 `json:"fullAmount"`     // 全款金额
    PaymentType    string  `json:"paymentType"`    // 支付类型：deposit/full
    WechatMerchant string  `json:"wechatMerchant"` // 微信商户号
    CallbackURL    string  `json:"callbackUrl"`    // 支付回调地址
}
```

**支付二维码生成：**
- 使用微信支付 Native 支付方式
- 二维码内容：`weixin://wxpay/bizpayurl?pr=xxxxx`
- 二维码有效期：2 小时
- 支持动态刷新二维码

### 3. 表单字段增强方案

**新增字段类型：**

```typescript
interface FormField {
  type: 'text' | 'phone' | 'email' | 'address' | 'textarea' | 'select';
  name: string;
  label: string;
  required: boolean;
  placeholder?: string;
  validation?: {
    pattern?: string;      // 正则表达式
    minLength?: number;    // 最小长度
    maxLength?: number;    // 最大长度
    message?: string;      // 错误提示
  };
  options?: string[];      // 下拉选项（仅 select 类型）
  order: number;           // 排序
}
```

**字段验证规则：**
- **email**: `/^[^\s@]+@[^\s@]+\.[^\s@]+$/`
- **phone**: `/^1[3-9]\d{9}$/`
- **address**: 最小长度 5，最大长度 200

**前端实现：**
- 使用拖拽排序（sortable.js）
- 实时预览表单效果
- 支持字段复制和删除

### 4. 订单核销方案

**选择：二维码扫描 + 状态更新**

**核销流程：**
```
1. 品牌管理员打开扫码页面
2. 扫描订单二维码（包含订单 ID）
3. 后端验证订单状态（已支付 + 未核销）
4. 显示订单详情，确认核销
5. 更新订单状态，记录核销信息
6. 返回核销成功提示
```

**订单二维码内容：**
```
https://dmh.example.com/verify?orderId=123456&code=abc123
```

**安全措施：**
- 二维码包含签名，防止伪造
- 核销操作需要品牌管理员权限
- 记录核销操作日志
- 支持取消核销（仅限当天）

**数据模型：**
```go
type Order struct {
    // ... 现有字段
    VerificationStatus string     `json:"verificationStatus"` // unverified/verified/cancelled
    VerifiedAt         *time.Time `json:"verifiedAt"`
    VerifiedBy         int64      `json:"verifiedBy"`         // 核销人 user_id
    VerificationCode   string     `json:"verificationCode"`   // 核销码（用于验证）
}
```

## Data Model Changes

### 1. campaigns 表变更

```sql
ALTER TABLE campaigns 
ADD COLUMN payment_config JSON COMMENT '支付配置' AFTER reward_rule,
ADD COLUMN poster_template_id INT DEFAULT 1 COMMENT '海报模板ID' AFTER payment_config;
```

### 2. orders 表变更

```sql
ALTER TABLE orders
ADD COLUMN verification_status VARCHAR(20) DEFAULT 'unverified' COMMENT '核销状态' AFTER pay_status,
ADD COLUMN verified_at DATETIME NULL COMMENT '核销时间' AFTER verification_status,
ADD COLUMN verified_by BIGINT NULL COMMENT '核销人ID' AFTER verified_at,
ADD COLUMN verification_code VARCHAR(50) NULL COMMENT '核销码' AFTER verified_by,
ADD INDEX idx_verification_status (verification_status),
ADD INDEX idx_verified_at (verified_at);
```

### 3. 新增 poster_templates 表

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

## API Design

### 1. 生成海报

```
POST /api/v1/campaigns/:id/poster
Authorization: Bearer {token}

Request:
{
  "templateId": 1,
  "distributorId": 123  // 可选，用于生成分销二维码
}

Response:
{
  "posterUrl": "https://cdn.example.com/posters/abc123.jpg",
  "qrcodeUrl": "https://dmh.example.com/c/123?d=456",
  "expiresAt": "2025-01-31T23:59:59Z"
}
```

### 2. 获取支付二维码

```
GET /api/v1/campaigns/:id/payment-qrcode?amount=99.00&orderId=123
Authorization: Bearer {token}

Response:
{
  "qrcodeUrl": "weixin://wxpay/bizpayurl?pr=xxxxx",
  "qrcodeImage": "data:image/png;base64,iVBORw0KGgo...",
  "expiresAt": "2025-01-24T14:00:00Z"
}
```

### 3. 扫码获取订单信息

```
GET /api/v1/orders/scan/:code
Authorization: Bearer {token}

Response:
{
  "orderId": 123,
  "campaignName": "新年促销",
  "customerName": "张三",
  "customerPhone": "138****8888",
  "amount": 99.00,
  "payStatus": "paid",
  "verificationStatus": "unverified",
  "formData": {...}
}
```

### 4. 核销订单

```
POST /api/v1/orders/:id/verify
Authorization: Bearer {token}

Request:
{
  "verificationCode": "abc123"
}

Response:
{
  "message": "核销成功",
  "verifiedAt": "2025-01-24T12:00:00Z",
  "verifiedBy": "brand_manager"
}
```

### 5. 取消核销

```
POST /api/v1/orders/:id/unverify
Authorization: Bearer {token}

Request:
{
  "reason": "用户要求退款"
}

Response:
{
  "message": "取消核销成功"
}
```

## Frontend Architecture

### H5 端新增页面

```
frontend-h5/src/views/brand/
├── PosterGenerator.vue      # 海报生成页面
├── OrderVerification.vue    # 订单核销页面
└── CampaignEditorEnhanced.vue  # 增强的活动编辑页面
```

### 核心组件

```vue
<!-- 海报生成组件 -->
<template>
  <div class="poster-generator">
    <div class="template-selector">
      <!-- 模板选择 -->
    </div>
    <div class="poster-preview">
      <!-- 海报预览 -->
    </div>
    <div class="actions">
      <button @click="generatePoster">生成海报</button>
      <button @click="downloadPoster">下载海报</button>
    </div>
  </div>
</template>

<!-- 扫码核销组件 -->
<template>
  <div class="order-verification">
    <div class="scanner">
      <!-- 二维码扫描器 -->
    </div>
    <div class="order-info" v-if="order">
      <!-- 订单信息展示 -->
    </div>
    <div class="actions">
      <button @click="verifyOrder">确认核销</button>
    </div>
  </div>
</template>
```

## Performance Considerations

### 1. 海报生成性能

- **目标**: 生成时间 < 3 秒
- **优化措施**:
  - 使用 goroutine 并发处理
  - 缓存生成的海报（24 小时）
  - 使用 CDN 加速图片访问
  - 限制并发生成数量（最多 10 个）

### 2. 二维码生成性能

- **目标**: 生成时间 < 500ms
- **优化措施**:
  - 使用内存缓存（Redis）
  - 预生成常用二维码
  - 使用 CDN 缓存二维码图片

### 3. 扫码核销性能

- **目标**: 响应时间 < 500ms
- **优化措施**:
  - 数据库索引优化
  - 使用缓存减少数据库查询
  - 异步记录操作日志

## Security Considerations

### 1. 海报生成安全

- 限制生成频率（每用户每分钟最多 5 次）
- 验证模板 ID 有效性
- 防止恶意上传图片

### 2. 支付二维码安全

- 二维码包含签名，防止伪造
- 二维码有效期限制（2 小时）
- 支付回调验证签名

### 3. 订单核销安全

- 核销码包含时间戳和签名
- 验证核销人权限（仅品牌管理员）
- 记录所有核销操作日志
- 防止重复核销

## Testing Strategy

### 1. 单元测试

- 海报生成逻辑测试
- 二维码生成测试
- 表单字段验证测试
- 订单核销逻辑测试

### 2. 集成测试

- 完整的海报生成流程测试
- 支付二维码生成和验证测试
- 扫码核销完整流程测试

### 3. 性能测试

- 并发生成海报压力测试
- 二维码生成性能测试
- 核销接口响应时间测试

### 4. 安全测试

- 核销码伪造测试
- 权限验证测试
- 频率限制测试

## Risks and Mitigations

### Risk 1: 海报生成服务器压力过大

**Mitigation**:
- 使用队列异步处理
- 限制并发数量
- 使用缓存减少重复生成

### Risk 2: 支付二维码过期导致用户无法支付

**Mitigation**:
- 提供刷新二维码功能
- 延长二维码有效期到 2 小时
- 显示倒计时提醒用户

### Risk 3: 核销码被盗用

**Mitigation**:
- 核销码包含时间戳和签名
- 限制核销码有效期（24 小时）
- 记录核销 IP 和设备信息

## Migration Plan

### Phase 1: 数据库迁移（第 1 天）

```bash
# 1. 备份数据库
mysqldump -u root -p dmh > backup_$(date +%Y%m%d).sql

# 2. 执行迁移脚本
mysql -u root -p dmh < migrations/20250124_add_advanced_features.sql

# 3. 验证迁移结果
mysql -u root -p dmh -e "DESCRIBE campaigns;"
mysql -u root -p dmh -e "DESCRIBE orders;"
mysql -u root -p dmh -e "DESCRIBE poster_templates;"
```

### Phase 2: 后端部署（第 2-3 天）

```bash
# 1. 编译新版本
cd backend/api
go build -o ../../logs/dmh-api-new dmh.go

# 2. 停止旧服务
./stop.sh

# 3. 替换可执行文件
mv logs/dmh-api logs/dmh-api-old
mv logs/dmh-api-new logs/dmh-api

# 4. 启动新服务
./start.sh

# 5. 验证服务
curl http://localhost:8889/health
```

### Phase 3: 前端部署（第 4-5 天）

```bash
# 1. 构建前端
cd frontend-h5
npm run build

# 2. 部署到服务器
rsync -avz dist/ user@server:/var/www/dmh-h5/

# 3. 重启 Nginx
ssh user@server "sudo systemctl reload nginx"
```

### Phase 4: 功能验证（第 6-7 天）

- 测试海报生成功能
- 测试支付二维码生成
- 测试表单字段配置
- 测试订单核销功能

## Rollback Plan

如果出现严重问题：

```bash
# 1. 回滚后端
./stop.sh
mv logs/dmh-api-old logs/dmh-api
./start.sh

# 2. 回滚前端
ssh user@server "cd /var/www/dmh-h5 && git checkout HEAD~1"
ssh user@server "sudo systemctl reload nginx"

# 3. 数据库回滚（如果需要）
mysql -u root -p dmh < backup_20250124.sql
```

## Open Questions

1. 海报模板是否需要支持用户自定义？
   - **决策**: MVP 阶段不支持，使用预设模板
   
2. 是否需要支持多种支付方式（支付宝、银联等）？
   - **决策**: MVP 阶段仅支持微信支付
   
3. 订单核销是否需要支持批量操作？
   - **决策**: MVP 阶段不支持，仅支持单个核销
   
4. 海报生成是否需要支持视频海报？
   - **决策**: MVP 阶段不支持，仅支持静态图片

## References

- [微信支付 Native 支付文档](https://pay.weixin.qq.com/wiki/doc/api/native.php)
- [Go gg 图像处理库](https://github.com/fogleman/gg)
- [QR Code 生成库](https://github.com/skip2/go-qrcode)
