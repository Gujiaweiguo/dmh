# DMH 活动高级功能 API 文档

## 概述

本文档描述了 DMH 数字营销中台系统新增的活动高级功能相关的 API 接口，包括：

* 海报生成功能
* 支付配置功能
* 订单核销功能

## API 端点

### 1. 海报生成相关

#### 1.1 生成活动海报

**接口地址：** `POST /api/v1/campaigns/{id}/poster`

**功能描述：** 为指定活动生成推广海报

**请求参数：**

```json
{
  "templateId": 1,
  "theme": "default"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| templateId | int64 | 是 | 海报模板ID |
| theme | string | 否 | 主题样式（默认：default） |

**响应示例：**

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "posterUrl": "https://oss.example.com/posters/campaign_123_poster_20250201_1643775200.jpg",
    "thumbnailUrl": "https://oss.example.com/posters/thumb_campaign_123_poster_20250201_1643775200.jpg",
    "fileSize": "256KB",
    "generationTime": 1800
  }
}
```

**频率限制：** 每用户每分钟最多生成 5 张海报

***

#### 1.2 生成分销商海报

**接口地址：** `POST /api/v1/distributors/{id}/poster`

**功能描述：** 为指定分销商生成推广海报

**请求参数：** 同活动海报生成

***

#### 1.3 获取海报模板列表

**接口地址：** `GET /api/v1/poster/templates`

**请求参数：**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int64 | 否 | 页码（默认：1）|
| pageSize | int64 | 否 | 每页数量（默认：20）|
| status | string | 否 | 模板状态（active/inactive）|
| keyword | string | 否 | 搜索关键词 |

**响应示例：**

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 3,
    "templates": [
      {
        "id": 1,
        "name": "经典模板",
        "previewImage": "/templates/classic.jpg",
        "config": {...},
        "status": "active",
        "createdAt": "2025-01-24T10:00:00",
        "updatedAt": "2025-01-24T10:00:00"
      }
    ]
  }
}
```

***

#### 1.4 获取海报生成记录

**接口地址：** `GET /api/v1/poster/records`

**响应示例：**

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 100,
    "records": [...]
  }
}
```

***

### 2. 支付配置相关

#### 2.1 获取支付二维码

**接口地址：** `GET /api/v1/campaigns/{id}/payment-qrcode`

**功能描述：** 获取活动的支付二维码

**请求参数：** 无

**响应示例：**

```json
{
  "qrcodeUrl": "weixin://wxpay/bizpayurl?pr=campaign_123_1643775200",
  "amount": 50.00,
  "campaignName": "测试活动"
}
```

**性能目标：** 二维码生成时间 < 500ms

**缓存机制：** 二维码缓存 2 小时，相同 URL 期间返回缓存

***

### 3. 订单核销相关

#### 3.1 扫码获取订单信息

**接口地址：** `GET /api/v1/orders/scan?code={verificationCode}`

**功能描述：** 通过核销码获取订单信息（扫码后展示订单详情）

**请求参数：**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| code | string | 是 | 核销码 |

**核销码格式：** `{orderId}_{phone}_{timestamp}_{signature}`

**响应示例：**

```json
{
  "orderId": 123,
  "status": "paid",
  "payStatus": "paid",
  "memberId": 456,
  "phone": "13800138000",
  "formData": {
    "name": "张三",
    "age": "25"
  }
}
```

**安全措施：**

* 核销码包含 HMAC-SHA1 签名
* 签名算法：`HMAC-SHA1(orderId_phone_timestamp, "dmh_secret_key")`
* 签名验证失败将返回 400 错误

***

#### 3.2 核销订单

**接口地址：** `POST /api/v1/orders/verify`

**功能描述：** 核销指定订单

**请求参数：**

```json
{
  "code": "{orderId}_{phone}_{timestamp}_{signature}"
}
```

**响应示例：**

```json
{
  "orderId": 123,
  "status": "verified",
  "verifiedAt": "2025-01-25 15:30:00"
}
```

**业务规则：**

* 只有已支付订单可以核销
* 不允许重复核销
* 核销操作需要品牌管理员权限
* 核销后订单状态更新为 `verified`
* 核销操作会记录日志

***

#### 3.3 取消核销

**接口地址：** `POST /api/v1/orders/unverify`

**功能描述：** 取消已核销的订单

**请求参数：** 同核销订单

**响应示例：**

```json
{
  "orderId": 123,
  "status": "unverified"
}
```

***

#### 3.4 查询核销记录

**接口地址：** `GET /api/v1/orders/verification-records`

**请求参数：**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| orderId | int64 | 否 | 订单ID筛选 |
| startTime | string | 否 | 开始时间 |
| endTime | string | 否 | 结束时间 |

**响应示例：**

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 500,
    "records": [
      {
        "id": 1,
        "orderId": 123,
        "verificationStatus": "verified",
        "verifiedAt": "2025-01-25 15:30:00",
        "verifiedBy": 1,
        "verificationCode": "123_13800138000_1643775200_abc123",
        "verificationMethod": "manual",
        "remark": "",
        "createdAt": "2025-01-25 15:30:00"
      }
    ]
  }
}
```

**性能目标：** 核销接口响应时间 < 500ms

***

## 通用说明

### 认证

所有 API 接口（登录接口除外）都需要在请求头中携带 JWT Token：

```
Authorization: Bearer {token}
```

### 错误码

| 错误码 | 说明 |
|---------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证（Token 无效或过期）|
| 403 | 无权限 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁（频率限制）|
| 500 | 服务器内部错误 |

### 响应格式

所有成功响应均采用统一格式：

```json
{
  "code": 200,
  "msg": "success",
  "data": {...}
}
```

错误响应格式：

```json
{
  "code": 400,
  "msg": "错误描述信息",
  "data": null
}
```

***

## 数据模型

### FormField（表单字段）

```json
{
  "type": "text|phone|email|textarea|select|address",
  "name": "字段名（英文）",
  "label": "字段标签（中文）",
  "required": true,
  "placeholder": "占位符",
  "options": ["选项1", "选项2"],
  "validation": {
    "pattern": "正则表达式",
    "minLength": 1,
    "maxLength": 100
  }
}
```

**支持的字段类型：**

* `text` - 文本输入
* `phone` - 手机号
* `email` - 邮箱
* `textarea` - 多行文本
* `select` - 下拉选择
* `address` - 地址

### PaymentConfig（支付配置）

```json
{
  "depositAmount": 50.00,
  "fullAmount": 200.00,
  "paymentType": "deposit|full",
  "wechatMerchant": "1234567890",
  "callbackUrl": "http://example.com/callback"
}
```

***

## 性能指标

| 接口 | 性能目标 | 实测结果 |
|------|---------|---------|
| 生成海报 | < 3 秒 | 待测试 |
| 生成支付二维码 | < 500ms | 待测试 |
| 订单核销 | < 500ms | 待测试 |
| 扫码查询 | < 300ms | 待测试 |

***

## 安全机制

### 1. 频率限制

| 接口 | 限制规则 |
|------|---------|
| 海报生成 | 5 次/分钟/用户 |
| 支付二维码生成 | 10 次/分钟/用户 |
| 登录接口 | 10 次/分钟/IP |

### 2. 签名验证

| 功能 | 签名算法 | 密钥 |
|------|---------|------|
| 核销码 | HMAC-SHA1 | dmh\_secret\_key |
| 支付二维码 | HMAC-SHA1 | dmh\_secret\_key |

### 3. 权限控制

| 接口 | 要求角色 |
|------|---------|
| 核销订单 | 品牌管理员 |
| 取消核销 | 品牌管理员 |
| 查询核销记录 | 品牌管理员、平台管理员 |

***

## 变更记录

| 版本 | 日期 | 变更内容 |
|-------|------|---------|
| 1.0 | 2025-02-01 | 初始版本，包含海报生成、支付二维码、订单核销功能 |
