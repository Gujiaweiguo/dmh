# 用户反馈收集系统文档

## 概述

DMH 用户反馈收集系统为活动高级功能提供全方位的用户体验跟踪和问题收集能力，支持以下功能：

* 用户反馈提交与管理
* 功能使用统计与分析
* 满意度调查
* 常见问题（FAQ）管理
* 反馈统计分析

***

## 数据库结构

### 1. 用户反馈表 (user\_feedback)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 反馈ID（主键） |
| user\_id | BIGINT | 用户ID |
| user\_name | VARCHAR(100) | 用户姓名（可选） |
| user\_role | VARCHAR(50) | 用户角色 |
| category | VARCHAR(50) | 反馈类别（poster/payment/verification/other） |
| subcategory | VARCHAR(50) | 子类别 |
| rating | TINYINT | 评分（1-5星） |
| title | VARCHAR(200) | 反馈标题 |
| content | TEXT | 反馈内容 |
| feature\_use\_case | VARCHAR(500) | 使用场景描述 |
| device\_info | VARCHAR(200) | 设备信息 |
| browser\_info | VARCHAR(200) | 浏览器信息 |
| priority | VARCHAR(20) | 优先级（low/medium/high） |
| status | VARCHAR(20) | 处理状态（pending/reviewing/resolved/closed） |
| assignee\_id | BIGINT | 处理人ID |
| response | TEXT | 处理回复 |
| resolved\_at | DATETIME | 解决时间 |
| created\_at | DATETIME | 创建时间 |
| updated\_at | DATETIME | 更新时间 |

### 2. 功能使用统计表 (feature\_usage\_stats)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 统计ID（主键） |
| user\_id | BIGINT | 用户ID |
| user\_role | VARCHAR(50) | 用户角色 |
| feature | VARCHAR(50) | 功能名称（poster/payment/verification） |
| action | VARCHAR(50) | 操作类型 |
| campaign\_id | BIGINT | 关联活动ID |
| success | BOOLEAN | 是否成功 |
| duration\_ms | INT | 耗时（毫秒） |
| error\_message | VARCHAR(500) | 错误信息 |
| created\_at | DATETIME | 创建时间 |

### 3. 功能满意度调查表 (feature\_satisfaction\_survey)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 调查ID（主键） |
| user\_id | BIGINT | 用户ID |
| user\_role | VARCHAR(50) | 用户角色 |
| feature | VARCHAR(50) | 功能名称 |
| ease\_of\_use | TINYINT | 易用性（1-5） |
| performance | TINYINT | 性能满意度（1-5） |
| reliability | TINYINT | 稳定性满意度（1-5） |
| overall\_satisfaction | TINYINT | 整体满意度（1-5） |
| would\_recommend | TINYINT | 推荐意愿（1-5） |
| most\_liked | VARCHAR(500) | 最满意的方面 |
| least\_liked | VARCHAR(500) | 最不满意的方面 |
| improvement\_suggestions | TEXT | 改进建议 |
| would\_like\_more\_features | VARCHAR(500) | 希望增加的功能 |
| created\_at | DATETIME | 创建时间 |

### 4. 常见问题表 (faq\_items)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | FAQ ID（主键） |
| category | VARCHAR(50) | 分类 |
| question | VARCHAR(500) | 问题 |
| answer | TEXT | 答案 |
| sort\_order | INT | 排序序号 |
| view\_count | INT | 浏览次数 |
| helpful\_count | INT | 有帮助次数 |
| not\_helpful\_count | INT | 无帮助次数 |
| is\_published | BOOLEAN | 是否发布 |
| created\_by | BIGINT | 创建人ID |
| created\_at | DATETIME | 创建时间 |
| updated\_at | DATETIME | 更新时间 |

***

## API 接口

### 1. 创建用户反馈

**接口**：`POST /api/v1/feedback`

**请求头**：

```
Authorization: Bearer {token}
Content-Type: application/json
```

**请求体**：

```json
{
  "category": "poster",
  "subcategory": "template",
  "rating": 4,
  "title": "希望增加更多海报模板",
  "content": "当前只有3种模板，希望能够增加更多样化的模板样式，特别是针对不同行业的专业模板。",
  "featureUseCase": "在发布新活动时，发现模板选择有限，无法满足品牌个性化需求",
  "priority": "medium"
}
```

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "userId": 123,
    "userName": "张三",
    "userRole": "brand_admin",
    "category": "poster",
    "subcategory": "template",
    "rating": 4,
    "title": "希望增加更多海报模板",
    "content": "当前只有3种模板，希望...",
    "priority": "medium",
    "status": "pending",
    "createdAt": "2026-02-05T10:00:00Z"
  }
}
```

***

### 2. 查询反馈列表

**接口**：`GET /api/v1/feedback`

**请求头**：

```
Authorization: Bearer {token}
```

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码（默认1） |
| pageSize | int | 否 | 每页数量（默认20） |
| category | string | 否 | 反馈类别 |
| status | string | 否 | 处理状态 |
| priority | string | 否 | 优先级 |

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "feedbacks": [
      {
        "id": 1,
        "userId": 123,
        "userName": "张三",
        "userRole": "brand_admin",
        "category": "poster",
        "title": "希望增加更多海报模板",
        "status": "pending",
        "createdAt": "2026-02-05T10:00:00Z"
      }
    ]
  }
}
```

***

### 3. 获取反馈详情

**接口**：`GET /api/v1/feedback/{id}`

**请求头**：

```
Authorization: Bearer {token}
```

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "userId": 123,
    "userName": "张三",
    "userRole": "brand_admin",
    "category": "poster",
    "subcategory": "template",
    "rating": 4,
    "title": "希望增加更多海报模板",
    "content": "当前只有3种模板，希望...",
    "featureUseCase": "在发布新活动时...",
    "priority": "medium",
    "status": "reviewing",
    "assigneeId": 456,
    "response": "收到您的反馈，我们已经在开发新模板，预计下个版本发布。",
    "resolvedAt": null,
    "createdAt": "2026-02-05T10:00:00Z",
    "tags": [
      {
        "id": 1,
        "name": "功能建议",
        "description": "用户提出的改进建议",
        "color": "#52c41a"
      }
    ]
  }
}
```

***

### 4. 更新反馈状态

**接口**：`PUT /api/v1/feedback/status`

**请求头**：

```
Authorization: Bearer {token}
Content-Type: application/json
```

**权限**：仅管理员（platform\_admin, brand\_admin）

**请求体**：

```json
{
  "id": 1,
  "status": "resolved",
  "assigneeId": 456,
  "response": "问题已解决，新模板已在 v1.1 版本中发布。"
}
```

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "status": "resolved",
    "assigneeId": 456,
    "response": "问题已解决...",
    "resolvedAt": "2026-02-05T15:00:00Z"
  }
}
```

***

### 5. 提交满意度调查

**接口**：`POST /api/v1/feedback/satisfaction`

**请求头**：

```
Authorization: Bearer {token}
Content-Type: application/json
```

**请求体**：

```json
{
  "feature": "poster",
  "easeOfUse": 4,
  "performance": 5,
  "reliability": 4,
  "overallSatisfaction": 4,
  "wouldRecommend": 5,
  "mostLiked": "生成速度快，模板美观",
  "leastLiked": "模板选择较少",
  "improvementSuggestions": "希望能增加自定义模板功能",
  "wouldLikeMoreFeatures": "海报编辑、批量生成"
}
```

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "userId": 123,
    "userRole": "brand_admin",
    "feature": "poster",
    "easeOfUse": 4,
    "performance": 5,
    "overallSatisfaction": 4,
    "mostLiked": "生成速度快...",
    "createdAt": "2026-02-05T10:00:00Z"
  }
}
```

***

### 6. 查询FAQ列表

**接口**：`GET /api/v1/feedback/faq`

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | string | 否 | 分类 |
| keyword | string | 否 | 关键词搜索 |

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 10,
    "faqs": [
      {
        "id": 1,
        "category": "poster",
        "question": "海报生成失败怎么办？",
        "answer": "请检查：1. 网络连接是否正常；2. 后端服务是否运行...",
        "sortOrder": 1,
        "viewCount": 123,
        "helpfulCount": 98,
        "notHelpfulCount": 5
      }
    ]
  }
}
```

***

### 7. 标记FAQ有帮助

**接口**：`POST /api/v1/feedback/faq/helpful`

**请求体**：

```json
{
  "id": 1,
  "type": "helpful"
}
```

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "helpfulCount": 99,
    "notHelpfulCount": 5
  }
}
```

***

### 8. 记录功能使用

**接口**：`POST /api/v1/feedback/usage`

**请求头**：

```
Authorization: Bearer {token}
Content-Type: application/json
```

**请求体**：

```json
{
  "feature": "poster",
  "action": "generate",
  "campaignId": 123,
  "success": true,
  "durationMs": 1800,
  "errorMessage": ""
}
```

**响应**：

```json
{
  "code": 200,
  "message": "使用记录已保存",
  "data": null
}
```

***

### 9. 获取反馈统计

**接口**：`GET /api/v1/feedback/statistics`

**请求头**：

```
Authorization: Bearer {token}
```

**权限**：仅管理员

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| startDate | string | 否 | 开始日期（YYYY-MM-DD） |
| endDate | string | 否 | 结束日期（YYYY-MM-DD） |
| category | string | 否 | 反馈类别 |

**响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "totalFeedbacks": 150,
    "byCategory": {
      "poster": 60,
      "payment": 45,
      "verification": 30,
      "other": 15
    },
    "byStatus": {
      "pending": 30,
      "reviewing": 45,
      "resolved": 60,
      "closed": 15
    },
    "byPriority": {
      "low": 50,
      "medium": 75,
      "high": 25
    },
    "averageRating": 4.2,
    "resolutionRate": 0.85,
    "avgResolutionTime": 48.5,
    "byRating": {
      "1": 5,
      "2": 10,
      "3": 25,
      "4": 60,
      "5": 50
    }
  }
}
```

***

## 前端集成示例

### 用户反馈表单

```vue
<template>
  <div class="feedback-form">
    <h2>提交反馈</h2>

    <van-form @submit="onSubmit">
      <van-field
        v-model="form.category"
        name="category"
        label="反馈类别"
        type="select"
        required
      >
        <template #options>
          <van-option value="poster" text="海报生成" />
          <van-option value="payment" text="支付配置" />
          <van-option value="verification" text="订单核销" />
          <van-option value="other" text="其他" />
        </template>
      </van-field>

      <van-field
        v-model="form.title"
        name="title"
        label="标题"
        placeholder="请输入反馈标题"
        required
      />

      <van-field
        v-model="form.content"
        name="content"
        label="详细描述"
        type="textarea"
        placeholder="请详细描述您的问题或建议"
        required
      />

      <van-field name="rating" label="评分">
        <template #input>
          <van-rate v-model="form.rating" />
        </template>
      </van-field>

      <van-field
        v-model="form.priority"
        name="priority"
        label="优先级"
        type="select"
      >
        <template #options>
          <van-option value="low" text="低" />
          <van-option value="medium" text="中" />
          <van-option value="high" text="高" />
        </template>
      </van-field>

      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit">
          提交反馈
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script>
import { ref } from 'vue';
import { showToast } from 'vant';
import axios from 'axios';

export default {
  setup() {
    const form = ref({
      category: '',
      title: '',
      content: '',
      rating: null,
      priority: 'medium'
    });

    const onSubmit = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.post('/api/v1/feedback', form.value, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (response.data.code === 200) {
          showToast('反馈提交成功！');
          form.value = {
            category: '',
            title: '',
            content: '',
            rating: null,
            priority: 'medium'
          };
        }
      } catch (error) {
        showToast('提交失败，请重试');
      }
    };

    return {
      form,
      onSubmit
    };
  }
};
</script>
```

***

### 功能使用记录

```javascript
// 在功能调用后记录使用情况
import axios from 'axios';

async function recordFeatureUsage(feature, action, campaignId, success, durationMs, errorMessage = '') {
  try {
    const token = localStorage.getItem('token');
    await axios.post('/api/v1/feedback/usage', {
      feature,
      action,
      campaignId,
      success,
      durationMs,
      errorMessage
    }, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
  } catch (error) {
    console.error('记录使用失败:', error);
  }
}

// 使用示例：海报生成
async function generatePoster(campaignId, templateId) {
  const startTime = Date.now();

  try {
    const response = await axios.post(`/api/v1/campaigns/${campaignId}/poster`, {
      templateId
    });

    const durationMs = Date.now() - startTime;

    // 记录成功使用
    await recordFeatureUsage('poster', 'generate', campaignId, true, durationMs);

    return response.data;
  } catch (error) {
    const durationMs = Date.now() - startTime;

    // 记录失败使用
    await recordFeatureUsage('poster', 'generate', campaignId, false, durationMs, error.message);

    throw error;
  }
}
```

***

### 满意度调查表单

```vue
<template>
  <div class="satisfaction-survey">
    <h2>功能满意度调查</h2>
    <p>请对我们提供的功能进行评价，帮助我们改进。</p>

    <van-form @submit="onSubmit">
      <van-field name="feature" label="功能名称" type="select" required>
        <template #options>
          <van-option value="poster" text="海报生成" />
          <van-option value="payment" text="支付配置" />
          <van-option value="verification" text="订单核销" />
        </template>
      </van-field>

      <van-field name="easeOfUse" label="易用性">
        <template #input>
          <van-rate name="easeOfUse" />
        </template>
      </van-field>

      <van-field name="performance" label="性能满意度">
        <template #input>
          <van-rate name="performance" />
        </template>
      </van-field>

      <van-field name="reliability" label="稳定性满意度">
        <template #input>
          <van-rate name="reliability" />
        </template>
      </van-field>

      <van-field name="overallSatisfaction" label="整体满意度">
        <template #input>
          <van-rate name="overallSatisfaction" />
        </template>
      </van-field>

      <van-field name="wouldRecommend" label="推荐意愿">
        <template #input>
          <van-rate name="wouldRecommend" />
        </template>
      </van-field>

      <van-field
        v-model="form.mostLiked"
        name="mostLiked"
        label="最满意的方面"
        type="textarea"
        placeholder="请描述您最满意的方面"
      />

      <van-field
        v-model="form.leastLiked"
        name="leastLiked"
        label="最不满意的方面"
        type="textarea"
        placeholder="请描述您最不满意的方面"
      />

      <van-field
        v-model="form.improvementSuggestions"
        name="improvementSuggestions"
        label="改进建议"
        type="textarea"
        placeholder="请提供改进建议"
      />

      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit">
          提交调查
        </van-button>
      </div>
    </van-form>
  </div>
</template>
```

***

## 使用建议

### 1. 反馈收集时机

* **首次使用后**：用户首次使用某个功能后，提示提交满意度调查
* **关键操作后**：完成海报生成、支付、核销等关键操作后，记录使用情况
* **遇到问题时**：在错误页面或失败操作后，提供反馈入口
* **定期提醒**：每使用功能10次后，提醒用户提交反馈

### 2. 反馈分类标准

| 类别 | 子类别 | 说明 |
|------|--------|------|
| poster | template | 模板相关 |
| poster | generation | 生成相关 |
| poster | performance | 性能相关 |
| payment | wechat | 微信支付相关 |
| payment | amount | 金额相关 |
| payment | refund | 退款相关 |
| verification | qr\_code | 二维码相关 |
| verification | signature | 签名相关 |
| verification | performance | 性能相关 |
| other | general | 其他问题 |

### 3. 优先级判断

* **High**：严重影响使用，阻碍核心功能
* **Medium**：影响使用体验，但可以绕过
* **Low**：改进建议，不影响当前使用

### 4. 反馈处理流程

```
用户提交反馈
    ↓
Pending（待处理）
    ↓
管理员分配（Assignee）
    ↓
Reviewing（处理中）
    ↓
Resolved（已解决）
    ↓
Closed（已关闭）
```

***

## 监控指标

### 关键指标

| 指标 | 说明 | 目标值 |
|------|------|--------|
| 反馈提交率 | 提交反馈的用户比例 | > 30% |
| 平均响应时间 | 反馈的平均处理时间 | < 48小时 |
| 解决率 | 已解决的反馈比例 | > 80% |
| 满意度评分 | 平均满意度评分 | > 4.0 |
| 推荐意愿 | 平均推荐意愿 | > 4.0 |

### 功能使用指标

| 指标 | 说明 |
|------|------|
| 功能使用次数 | 每个功能的使用次数 |
| 成功率 | 操作成功率 |
| 平均耗时 | 操作平均耗时 |
| 错误率 | 操作错误率 |

***

## 最佳实践

1. **及时响应**：尽量在24小时内响应用户反馈
2. **透明沟通**：向用户反馈处理进度
3. **持续改进**：定期分析反馈数据，优化功能
4. **用户参与**：邀请高频用户参与功能测试和设计
5. **数据驱动**：基于数据做出产品决策

***

**文档版本**：1.0
**最后更新**：2026-02-05
**维护人员**：DMH 技术团队
