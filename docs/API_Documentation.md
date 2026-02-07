# DMH API 文档

## 基础信息

* **API 版本**: v1.0.0
* **基础URL**: `http://localhost:8889/api/v1`
* **协议**: HTTP/HTTPS
* **数据格式**: JSON
* **认证方式**: JWT Bearer Token

## 认证

### 登录获取 Token

**端点**: `POST /auth/login`

**请求体**：

```json
{
  "username": "admin",
  "密码": "your_password"
}
```

**响应**：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJuYW1lIjoiYWRtaW4iLCJyb2xlcyI6WyJwbGF0Zm9ybV9hZG1pbiJdLCJleHAiOjE3Njk4ODU44iLCJpc3MiI6ImRtaC1zeS1zdGVtIiwibmJmYiOjEzNjk3MTU44iLCJpYXQiIjE3Njk3MTY1Nz0fQ.EYyCCt8FRivPC8i59PBojpGa2v6_yrYazFaENOTfC9o"
  },
  "userId": 1,
  "username": "admin",
  "phone": "13800000001",
  "roles": ["platform_admin"],
  "brandIds": null
}
```

**状态码说明**：

* `200`: 成功
* `400`: 请求参数错误
* `401`: 未授权
* `403`: 权限不足
* `500`: 服务器内部错误

## 用户角色

| 角色 | 说明 | 权限 |
|------|------|------|
| platform\_admin | 平台管理员 | 所有权限 |
| brand\_admin | 品牌管理员 | 品牌数据管理 |
| participant | 参与者 | 活动报名、订单查询 |

***

## 活动管理

### 获取活动列表

**端点**: `GET /campaigns`

**查询参数**：
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| page | int64 | 否 | 页码，从 1 开始，默认 1 |
| pageSize | int64 | 否 | 每页数量，默认 20 |
| status | string | 否 | 活动状态筛选（active/paused/ended） |
| keyword | string | 否 | 关键词搜索 |

**请求示例**：

```bash
curl -X GET "http://localhost:8889/api/v1/campaigns?page=1&pageSize=20&status=active" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应**：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 10,
    "campaigns": [
      {
        "id": 1,
        "brandId": 1,
        "name": "测试活动",
        "description": "这是一个测试活动",
        "formFields": "[{\"type\":\"text\",\"name\":\"name\",\"label\":\"姓名\",\"required\":true}]",
        "rewardRule": 10.0,
        "startTime": "2026-02-01T10:00:00",
        "endTime": "2026-03-01T10:00:00",
        "status": "active",
        "createdAt": "2026-01-30T12:00:00",
        "updatedAt": "2026-01-30T12:00:00"
      }
    ]
  }
}
```

### 获取单个活动

**端点**: `GET /campaigns/:id`

**请求示例**：

```bash
curl -X GET "http://localhost:8889/api/v1/campaigns/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应**：同上单个活动对象

### 创建活动

**端点**: `POST /campaigns`

**请求体**：

```json
{
  "brandId": 1,
  "name": "新活动",
  "description": "活动描述",
  "formFields": [
    {
      "type": "text",
      "name": "name",
      "label": "姓名",
      "required": true
    },
    {
      "type": "phone",
      "name": "phone",
      "label": "手机号",
      "required": true
    },
    {
      "type": "select",
      "name": "gender",
      "label": "性别",
      "required": false,
      "options": ["男", "女", "其他"]
    }
  ],
  "rewardRule": 10.0,
  "startTime": "2026-02-01T10:00:00",
  "endTime": "2026-03-01T10:00:00"
}
```

**formFields 字段类型**：
| 类型 | 说明 | 必填 | 特殊属性 |
|------|------|------|------|----------|
| text | 文本框 | 可选 | - |
| phone | 手机号 | 可选 | 格式验证 |
| email | 邮箱 | 可选 | 格式验证 |
| select | 下拉选择 | 可选 | 需要 `options` 数组 |
| textarea | 多行文本 | 可选 | - |
| address | 地址 | 可选 | - |

### 更新活动

**端点**: `PUT /campaigns/:id`

**请求体**：

```json
{
  "brandId": 1,
  "name": "更新后的名称",
  "description": "更新后的描述",
  "formFields": [...],
  "rewardRule": 20.0,
  "startTime": "2026-02-01T10:00:00",
  "endTime": "2026-03-01T10:00:00",
  "status": "active"
}
```

### 删除活动

**端点**: `DELETE /campaigns/:id`

**请求示例**：

```bash
curl -X DELETE "http://localhost:8889/api/v1/campaigns/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

***

## 支付二维码

### 生成支付二维码

**端点**: `GET /campaigns/:id/payment-qrcode`

**请求示例**：

```bash
curl -X GET "http://localhost:8889/api/v1/campaigns/1/payment-qrcode" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应**：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "qrcodeUrl": "weixin://wxpay/bizpayurl?pr=campaign_1_1738267320",
    "amount": 100.00,
    "campaignName": "测试活动"
  }
}
```

**说明**：

* `amount`: 支付金额 = 活动奖励规则 \* 10
* `qrcodeUrl`: 微信支付二维码链接
* 二维码包含时间戳，每次调用都会返回不同的 URL（用于刷新）

***

## 订单管理

### 创建订单

**端点**: `POST /orders`

**请求体**：

```json
{
  "campaignId": 1,
  "phone": "13800138001",
  "formData": {
    "name": "张三",
    "remark": "备注信息"
  },
  "referrerId": 1
}
```

### 获取订单列表

**端点**: `GET /orders`

**查询参数**：
| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| campaignId | int64 | 否 | 按活动ID筛选 |
| phone | string | 否 | 按手机号搜索 |
| status | string | 否 | 订单状态（pending/paid/cancelled） |

**请求示例**：

```bash
curl -X GET "http://localhost:8889/api/v1/orders?status=paid&page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 查询订单详情

**端点**: `GET /orders/:id`

### 生成核销码

**端点**: `POST /orders/:id/generate-verification-code`

**请求示例**：

```bash
curl -X POST "http://localhost:8889/api/v1/orders/1/generate-verification-code" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应**：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "orderId": 1,
    "verificationCode": "1_13800138001_1738267320_abcd1234",
    "expiresAt": 1738267620
  }
}
```

**核销码格式**: `{orderId}_{phone}_{timestamp}_{signature}`

### 订单核销

**端点**: `POST /orders/verify`

**请求体**：

```json
{
  "code": "1_13800138001_1738267320_abcd1234"
}
```

**响应**：

```json
{
  "code": 200,
  "msg": "核销成功",
  "data": {
    "orderId": 1,
    "verificationStatus": "verified",
    "verifiedAt": "2026-01-30T12:00:00",
    "verifiedBy": 1
  }
}
```

**核销码验证机制**：

* 签名算法：HMAC-SHA1
* 密钥：`dmh_secret_key`
* 拒绝重复核销
* 订单已核销后无法再次核销

***

## 错误码说明

| 错误码 | HTTP | 说明 |
|--------|------|------|
| 200 | OK | 请求成功 |
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | Token 无效或过期 |
| 403 | Forbidden | 权限不足 |
| 404 | Not Found | 资源不存在 |
| 500 | Internal Server Error | 服务器内部错误 |

常见错误响应示例：

```json
{
  "code": 400,
  "msg": "参数验证失败：startTime 不能早于 endTime",
  "data": null
}
```

***

## 分页参数

所有列表接口都支持分页：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int64 | 1 | 页码 |
| pageSize | int64 | 20 | 每页数量 |
| total | int64 | - | 返回总记录数 |

**分页响应格式**：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 100,
    "list": [...]
  }
}
```

***

## 速率限制

### 海报生成

* **频率**: 5 次/60 秒
* **窗口**: 滑动窗口

### 通用API

* **频率**: 100 次/分钟
* **窗口**: 60 秒

超限响应：

```json
{
  "code": 429,
  "msg": "请求过于频繁，请稍后再试"
}
```

***

## 安全机制

### JWT Token 认证

* **算法**: HS256
* **过期时间**: 86400 秒（24 小时）
* **刷新机制**: 过期后需要重新登录

### 核销码签名

* **格式**: `{orderId}_{phone}_{timestamp}_{signature}`
* **算法**: HMAC-SHA1
* **密钥**: `dmh_secret_key`
* **防护措施**：
  * 伪造签名检测
  * 重复核销拒绝
  * 核销码过期检查

### 权限控制

* **基于角色的 RBAC**
* **资源级别权限控制**
* **未授权访问返回 401/403**

***

## 时间格式

所有时间字段使用 ISO 8601 格式：

* 示例：`2026-01-30T12:00:00`
* 时区：本地时间

***

## 测试工具

### cURL 示例

```bash
# 登录并获取 token
TOKEN=$(curl -s -X POST 'http://localhost:8889/api/v1/auth/login' \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"your_password"}' | jq -r .data.token)

# 创建活动
curl -X POST 'http://localhost:8889/api/v1/campaigns' \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN' \
  -d '{
    "brandId": 1,
    "name": "新活动",
    "rewardRule": 10.0,
    "startTime": "2026-02-01T10:00:00",
    "endTime": "2026-03-01T10:00:00"
  }'

# 获取活动列表
curl -X GET 'http://localhost:8889/api/v1/campaigns?page=1&pageSize=10' \
  -H "Authorization: Bearer $TOKEN"

# 生成支付二维码
curl -X GET 'http://localhost:8889/api/v1/campaigns/1/payment-qrcode' \
  -H "Authorization: Bearer $TOKEN'
```

### Postman 集合测试

建议使用 Postman 进行 API 测试，导入以下环境变量：

* `base_url`: `http://localhost:8889/api/v1`
* `token`: 从登录响应中获取

***

## 版本历史

| 版本 | 日期 | 变更内容 |
|------|------|----------|
| v1.0.0 | 2025-01-30 | 初始版本，包含：

* 活动管理（增删改查）
* 支付二维码生成
* 订单管理（创建、查询、核销）
* JWT 认证
* RBAC 权限控制
* 速率限制
* 核销码签名验证
* FormFields 动态表单字段

***

## 联系与支持

如有问题，请联系：

* 技术支持：\[待补充]
* 技术文档：\[待补充]

***

*最后更新：2025年1月30日*
