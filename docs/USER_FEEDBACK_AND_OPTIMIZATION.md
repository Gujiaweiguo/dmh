# DMH 高级功能用户反馈收集与优化指南

## 📋 概述

本文档提供了 DMH 高级功能（海报生成、支付配置、表单增强、订单核销）的用户反馈收集方法和基于反馈的优化建议。

**相关文档**:

* [性能测试报告](../backend/test/performance/PERFORMANCE_TEST_REPORT.md)
* [生产部署指南](../deployment/PRODUCTION_DEPLOYMENT_GUIDE.md)
* [监控配置指南](../monitoring/MONITORING_SETUP_GUIDE.md)

***

## 🎯 反馈收集目标

### 关键指标

| 指标 | 目标值 | 测量方法 |
|------|--------|---------|
| **用户满意度** | > 4.5/5 | 问卷调查 |
| **功能使用率** | > 80% | 后台统计 |
| **功能缺陷率** | < 5% | Bug 跟踪 |
| **功能请求响应时间** | < 24 小时 | 工单系统 |
| **功能采用率** | > 70% | 使用统计 |
| **用户留存率** | > 85% | 活跃用户统计 |

***

## 📊 反馈收集渠道

### 1. 应用内反馈

#### 1.1 添加反馈按钮

在前端应用中添加便捷的反馈入口：

**H5 前端** - 在用户菜单和活动详情页添加：

```vue
<!-- UserMenu.vue -->
<template>
  <div class="user-menu">
    <van-button @click="showFeedbackModal">
      💬 意见反馈
    </van-button>
  </div>
</template>

<script>
export default {
  methods: {
    showFeedbackModal() {
      this.$router.push('/feedback');
    }
  }
}
</script>
```

**管理后台** - 在设置菜单添加：

```vue
<!-- SettingsMenu.vue -->
<template>
  <el-menu-item index="/admin/feedback">
    <i class="el-icon-s-comment"></i>
    <span>用户反馈</span>
  </el-menu-item>
</template>
```

#### 1.2 反馈表单设计

创建反馈收集表单组件：

**文件**: `frontend-h5/src/components/FeedbackForm.vue`

```vue
<template>
  <van-dialog
    v-model="showFeedback"
    title="意见反馈"
    show-cancel-button
    confirm-button-text="提交"
    @confirm="submitFeedback"
  >
    <van-form @submit="submitFeedback">
      <van-field
        v-model="feedback.category"
        name="category"
        label="反馈类型"
        is-link
        readonly
        @click="showCategoryPicker = true"
      />

      <van-field
        v-model="feedback.feature"
        name="feature"
        label="功能模块"
        is-link
        readonly
        @click="showFeaturePicker = true"
      />

      <van-field
        v-model="feedback.content"
        name="content"
        label="反馈内容"
        type="textarea"
        rows="4"
        placeholder="请详细描述您的反馈意见..."
        required
      />

      <van-field name="rating" label="满意度">
        <template #input>
          <van-rate v-model="feedback.rating" />
        </template>
      </van-field>

      <van-field
        v-model="feedback.contact"
        name="contact"
        label="联系方式（可选）"
        placeholder="邮箱或手机号，方便我们回复您"
      />
    </van-form>
  </van-dialog>
</template>

<script>
export default {
  data() {
    return {
      showFeedback: false,
      showCategoryPicker: false,
      showFeaturePicker: false,
      categories: ['功能建议', 'Bug 报告', '性能问题', 'UI/UX 改进', '其他'],
      features: ['海报生成', '支付配置', '表单管理', '订单核销', '其他'],
      feedback: {
        category: '',
        feature: '',
        content: '',
        rating: 5,
        contact: ''
      }
    };
  },
  methods: {
    async submitFeedback() {
      try {
        await this.$http.post('/api/v1/feedback', {
          ...this.feedback,
          userAgent: navigator.userAgent,
          pageUrl: window.location.href,
          timestamp: new Date().toISOString()
        });
        this.$toast.success('感谢您的反馈！');
        this.showFeedback = false;
        this.resetFeedback();
      } catch (error) {
        this.$toast.fail('提交失败，请稍后重试');
      }
    },
    resetFeedback() {
      this.feedback = {
        category: '',
        feature: '',
        content: '',
        rating: 5,
        contact: ''
      };
    }
  }
};
</script>
```

#### 1.3 后端反馈 API

创建反馈数据表和 API：

**数据库迁移**:

```sql
CREATE TABLE IF NOT EXISTS user_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    category VARCHAR(50) NOT NULL COMMENT '反馈类型：feature_suggestion, bug_report, performance, ui_ux, other',
    feature VARCHAR(50) COMMENT '功能模块：poster, payment, form, order_verify, other',
    content TEXT NOT NULL COMMENT '反馈内容',
    rating TINYINT DEFAULT 5 COMMENT '满意度评分 1-5',
    contact VARCHAR(100) COMMENT '联系方式',
    user_agent VARCHAR(500) COMMENT '用户代理',
    page_url VARCHAR(500) COMMENT '反馈页面URL',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态：pending, reviewing, resolved, closed',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_category (category),
    INDEX idx_feature (feature),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户反馈表';
```

**API 定义** (`backend/api/dmh.api`):

```api
type (
    // 提交反馈请求
    CreateFeedbackReq {
        Category string `json:"category"`
        Feature string `json:"feature"`
        Content string `json:"content"`
        Rating int8 `json:"rating"`
        Contact string `json:"contact,optional"`
        UserAgent string `json:"userAgent,optional"`
        PageUrl string `json:"pageUrl,optional"`
    }

    // 提交反馈响应
    CreateFeedbackResp {
        Id int64 `json:"id"`
    }

    // 反馈列表项
    FeedbackItem {
        Id int64 `json:"id"`
        UserId int64 `json:"userId"`
        Category string `json:"category"`
        Feature string `json:"feature"`
        Content string `json:"content"`
        Rating int8 `json:"rating"`
        Contact string `json:"contact"`
        Status string `json:"status"`
        CreatedAt string `json:"createdAt"`
    }

    // 反馈列表响应
    FeedbackListResp {
        Total int64 `json:"total"`
        List []FeedbackItem `json:"list"`
    }
)

// 提交用户反馈
@server(
    prefix: /api/v1
    group: feedback
    middleware: AuthMiddleware
)
service dmh-api {
    @doc "提交用户反馈"
    @handler createFeedback
    post /feedback (CreateFeedbackReq) returns (CreateFeedbackResp)

    @doc "获取反馈列表（管理员）"
    @handler getFeedbackList
    get /feedback (GetFeedbackListReq) returns (FeedbackListResp)
}
```

### 2. 问卷调查

#### 2.1 上线后问卷调查

在功能上线 1 周后发送问卷调查：

```markdown
# DMH 高级功能满意度调查

## 基本信息
- 角色：[ ] 品牌管理员 [ ] 平台管理员 [ ] 普通用户
- 使用时长：[ ] 1周以内 [ ] 1-2周 [ ] 2-4周 [ ] 1个月以上

## 功能使用情况

### 1. 海报生成功能
**使用频率**:
[ ] 每天使用 [ ] 每周 2-3 次 [ ] 偶尔使用 [ ] 从未使用

**满意度评分** (1-5 分):
- 生成速度: [1] [2] [3] [4] [5]
- 海报质量: [1] [2] [3] [4] [5]
- 操作便利性: [1] [2] [3] [4] [5]

**需要改进的地方**:
___________________________________________________________

**建议新增功能**:
___________________________________________________________

### 2. 支付配置功能
**使用频率**:
[ ] 每次活动都配置 [ ] 部分活动配置 [ ] 从未使用

**满意度评分** (1-5 分):
- 配置便捷性: [1] [2] [3] [4] [5]
- 支付二维码质量: [1] [2] [3] [4] [5]
- 支付成功率: [1] [2] [3] [4] [5]

**遇到的问题**:
___________________________________________________________

**建议改进**:
___________________________________________________________

### 3. 表单字段增强功能
**使用频率**:
[ ] 每次活动都使用 [ ] 部分活动使用 [ ] 从未使用

**满意度评分** (1-5 分):
- 字段类型丰富度: [1] [2] [3] [4] [5]
- 验证规则灵活性: [1] [2] [3] [4] [5]
- 字段排序功能: [1] [2] [3] [4] [5]

**需要新增的字段类型**:
___________________________________________________________

**使用场景描述**:
___________________________________________________________

### 4. 订单核销功能
**使用频率**:
[ ] 每天使用 [ ] 每周 2-3 次 [ ] 偶尔使用 [ ] 从未使用

**满意度评分** (1-5 分):
- 扫码识别速度: [1] [2] [3] [4] [5]
- 核销操作便捷性: [1] [2] [3] [4] [5]
- 核销记录完整性: [1] [2] [3] [4] [5]

**遇到的问题**:
___________________________________________________________

**建议改进**:
___________________________________________________________

## 整体评价

**整体满意度**: [1] [2] [3] [4] [5]

**最满意的功能**: ___________________________

**最不满意的功能**: ___________________________

**是否愿意推荐给他人**: [ ] 肯定推荐 [ ] 可能推荐 [ ] 不会推荐

**其他建议**:
___________________________________________________________

感谢您的反馈！
```

#### 2.2 问卷调查收集方式

* **邮件推送**: 使用系统邮件发送问卷调查链接
* **应用内弹窗**: 上线后首次使用功能时弹出问卷
* **公众号推送**: 在微信公众号推送调查通知

### 3. 数据分析反馈

#### 3.1 功能使用统计

创建数据查询，分析功能使用情况：

```sql
-- 海报生成统计
SELECT
    DATE(created_at) as date,
    COUNT(*) as poster_count,
    COUNT(DISTINCT user_id) as unique_users
FROM poster_records
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- 支付二维码生成统计
SELECT
    DATE(created_at) as date,
    COUNT(*) as qrcode_count,
    AVG(generation_time) as avg_generation_time
FROM payment_qrcode_records
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- 订单核销统计
SELECT
    DATE(verified_at) as date,
    COUNT(*) as verified_count,
    AVG(TIMESTAMPDIFF(SECOND, created_at, verified_at)) as avg_verify_time_seconds
FROM orders
WHERE verification_status = 'verified'
  AND verified_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY DATE(verified_at)
ORDER BY date DESC;

-- 表单字段类型使用统计
SELECT
    type,
    COUNT(DISTINCT campaign_id) as campaign_count,
    COUNT(*) as field_count
FROM campaign_form_fields
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY type
ORDER BY field_count DESC;
```

#### 3.2 性能数据统计

```sql
-- 海报生成性能统计
SELECT
    MIN(generation_time) as min_time,
    MAX(generation_time) as max_time,
    AVG(generation_time) as avg_time,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY generation_time) as p95_time
FROM poster_records
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY);

-- 海报生成成功率
SELECT
    COUNT(CASE WHEN status = 'success' THEN 1 END) * 100.0 / COUNT(*) as success_rate
FROM poster_records
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY);
```

### 4. 用户访谈

#### 4.1 访谈计划

**访谈对象选择**:

* 高频用户（每周使用 3 次以上）
* 0 评分用户（给予特别关注）
* 活跃用户（登录频次高但功能使用率低）
* 新用户（上线后注册的用户）

**访谈提纲**:

```markdown
# DMH 高级功能用户访谈提纲

## 开场
- 介绍访谈目的：了解功能使用情况和改进方向
- 说明访谈时长：30-45 分钟
- 确认是否可以录音

## 基本信息
- 您的角色和主要工作内容
- 使用 DMH 系统的时长
- 主要负责哪些活动

## 功能体验

### 海报生成功能
- 您使用海报生成功能的频率？
- 您通常在什么场景下使用？
- 海报生成速度和质量是否满足需求？
- 海报模板是否够用？希望有什么新模板？
- 您是否遇到过生成失败的情况？

### 支付配置功能
- 您配置过支付功能吗？为什么配置/不配置？
- 支付二维码生成是否顺畅？
- 您遇到过支付相关的问题吗？
- 支付配置流程是否清晰？

### 表单字段增强功能
- 您使用过哪些新增的字段类型？
- 字段验证规则是否灵活？
- 字段排序功能是否实用？
- 还需要什么类型的字段？

### 订单核销功能
- 您的团队如何使用核销功能？
- 扫码识别速度和准确度如何？
- 核销操作流程是否顺畅？
- 核销记录查询是否方便？

## 整体评价
- 这些高级功能解决了您的什么问题？
- 还有哪些痛点没有解决？
- 功能之间的协作是否顺畅？
- 系统整体易用性如何？

## 未来期望
- 希望新增哪些功能？
- 希望改进哪些功能？
- 对系统整体有什么建议？

## 结束
- 感谢用户参与
- 告知后续行动和联系方式
- 询问是否还有其他补充
```

***

## 📈 反馈分析方法

### 1. 分类汇总

创建反馈分类统计表：

| 分类 | 数量 | 占比 | 优先级 | 处理状态 |
|------|------|------|--------|---------|
| 功能建议 | 45 | 35% | 高 | 处理中 |
| Bug 报告 | 28 | 22% | 高 | 已解决 20 |
| 性能问题 | 18 | 14% | 高 | 处理中 |
| UI/UX 改进 | 25 | 19% | 中 | 计划中 |
| 其他 | 12 | 10% | 低 | 已归档 |

### 2. 优先级评估

使用 RICE 模型评估反馈优先级：

**RICE 模型**:

* **Reach (影响范围)**: 有多少用户会受益
* **Impact (影响程度)**: 对用户有多大价值
* **Confidence (信心程度)**: 对评估结果的信心
* **Effort (投入成本)**: 实现所需工作量

**评分表**:

| 反馈项 | Reach | Impact | Confidence | Effort | RICE 得分 | 优先级 |
|--------|--------|---------|-----------|----------|-----------|--------|
| 海报生成太慢 | 85% | 高 | 90% | 2周 | 34 | P0 |
| 海报模板太少 | 70% | 中 | 80% | 1周 | 28 | P1 |
| 核销记录查询不便 | 60% | 中 | 70% | 1周 | 14 | P2 |
| 希望支持视频海报 | 40% | 低 | 50% | 4周 | 2 | P3 |

### 3. 趋势分析

按时间维度分析反馈趋势：

* **功能采用率趋势**: 每周功能使用用户数变化
* **满意度趋势**: 满意度评分随时间变化
* **Bug 数量趋势**: Bug 报告数量变化
* **功能请求趋势**: 功能建议数量变化

***

## 🔄 基于反馈的优化建议

### 优化 1: 海报生成性能优化

**反馈内容**: "海报生成速度慢，特别是并发生成时"

**优化方案**:

1. **实现海报缓存机制**:

```go
// backend/api/internal/logic/poster/generateCampaignPosterLogic.go

// 检查缓存
cacheKey := fmt.Sprintf("poster:campaign:%d:template:%d", campaignID, templateID)
if cachedPoster, err := redis.Get(cacheKey); err == nil {
    return &types.GeneratePosterResp{
        PosterUrl: cachedPoster,
        GenerationTime: 0,
    }, nil
}

// 生成海报...
// posterURL := ...

// 设置缓存（1小时TTL）
redis.Set(cacheKey, posterURL, time.Hour)
```

2. **使用 CDN 加速海报访问**:

```yaml
# backend/api/internal/svc/service_context.go
PosterService:
  StorageType: oss  # oss, local
  OSS:
    Endpoint: "oss-cn-hangzhou.aliyuncs.com"
    Bucket: "dmh-posters"
    AccessKeyID: "your_access_key"
    AccessKeySecret: "your_secret_key"
    CDN: "https://cdn.dmh.com"
```

3. **实现异步海报生成**:

```go
// 使用消息队列异步生成海报
func (l *GenerateCampaignPosterLogic) GenerateCampaignPosterAsync(req *types.GeneratePosterReq) (resp *types.GeneratePosterResp, err error) {
    // 生成任务ID
    taskID := generateTaskID()

    // 发送到消息队列
    queue.Publish("poster:generate", map[string]interface{}{
        "task_id": taskID,
        "campaign_id": req.CampaignID,
        "template_id": req.TemplateID,
        "user_id": l.ctx.Value("user_id"),
    })

    // 立即返回任务ID
    return &types.GeneratePosterResp{
        TaskId: taskID,
        Status: "processing",
    }, nil
}
```

### 优化 2: 海报模板扩展

**反馈内容**: "海报模板太少，希望有更多样式"

**优化方案**:

1. **新增海报模板类型**:

```sql
-- 插入新模板
INSERT INTO poster_templates (name, config, status) VALUES
('简约风格', '{"width":750,"height":1334,"background":"#FFFFFF","elements":[...]}', 'active'),
('商务风格', '{"width":750,"height":1334,"background":"#2C3E50","elements":[...]}', 'active'),
('节日风格', '{"width":750,"height":1334,"background":"#E74C3C","elements":[...]}', 'active'),
('简约风格-黑色', '{"width":750,"height":1334,"background":"#000000","elements":[...]}', 'active'),
('科技风格', '{"width":750,"height":1334","background":"linear-gradient(135deg, #667eea 0%, #764ba2 100%)","elements":[...]}', 'active');
```

2. **实现自定义模板功能**:

```go
// backend/api/internal/logic/poster/templateLogic.go

// 创建自定义模板
func (l *TemplateLogic) CreateCustomTemplate(req *types.CreateTemplateReq) (resp *types.CreateTemplateResp, err error) {
    // 验证用户权限
    if !l.isAdmin() && !l.isBrandAdmin() {
        return nil, errors.New("权限不足")
    }

    // 保存模板
    template := model.PosterTemplate{
        Name: req.Name,
        PreviewImage: req.PreviewImage,
        Config: req.Config,
        Status: "pending",  // 需要管理员审核
        CreatedBy: l.ctx.Value("user_id"),
    }

    if err := l.svcCtx.DB.Create(&template).Error; err != nil {
        return nil, err
    }

    return &types.CreateTemplateResp{
        Id: template.Id,
        Status: template.Status,
    }, nil
}
```

3. **提供模板预览功能**:

```vue
<!-- TemplateSelector.vue -->
<template>
  <div class="template-selector">
    <div class="template-grid">
      <div
        v-for="template in templates"
        :key="template.id"
        class="template-item"
        :class="{ active: selectedTemplate === template.id }"
        @click="selectTemplate(template)"
      >
        <img :src="template.previewImage" :alt="template.name" />
        <div class="template-info">
          <h3>{{ template.name }}</h3>
          <p v-if="template.isCustom">自定义</p>
        </div>
      </div>
    </div>
  </div>
</template>
```

### 优化 3: 核销记录查询优化

**反馈内容**: "核销记录查询不方便，希望支持多条件筛选"

**优化方案**:

1. **增强核销记录查询 API**:

```api
// backend/api/dmh.api

type (
    // 核销记录查询请求
    QueryVerificationRecordsReq {
        CampaignId int64 `json:"campaignId,optional"`
        OrderId int64 `json:"orderId,optional"`
        UserId int64 `json:"userId,optional"`
        VerifiedBy int64 `json:"verifiedBy,optional"`
        Status string `json:"status,optional"`
        VerificationMethod string `json:"verificationMethod,optional"`
        StartTime string `json:"startTime,optional"`
        EndTime string `json:"endTime,optional"`
        Page int `json:"page,default=1"`
        PageSize int `json:"pageSize,default=20"`
    }

    // 核销记录项
    VerificationRecord {
        Id int64 `json:"id"`
        CampaignId int64 `json:"campaignId"`
        CampaignName string `json:"campaignName"`
        OrderId int64 `json:"orderId"`
        UserId int64 `json:"userId"`
        UserName string `json:"userName"`
        VerifiedBy int64 `json:"verifiedBy"`
        VerifiedByName string `json:"verifiedByName"`
        VerificationCode string `json:"verificationCode"`
        Status string `json:"status"`
        VerificationMethod string `json:"verificationMethod"`
        VerifiedAt string `json:"verifiedAt"`
        CreatedAt string `json:"createdAt"`
    }

    // 核销记录列表响应
    QueryVerificationRecordsResp {
        Total int64 `json:"total"`
        List []VerificationRecord `json:"list"`
    }
)

@server(
    prefix: /api/v1
    group: verification
    middleware: AuthMiddleware, BrandAdminMiddleware
)
service dmh-api {
    @doc "查询核销记录"
    @handler queryVerificationRecords
    get /verification/records (QueryVerificationRecordsReq) returns (QueryVerificationRecordsResp)

    @doc "导出核销记录"
    @handler exportVerificationRecords
    get /verification/records/export (QueryVerificationRecordsReq)
}
```

2. **实现数据导出功能**:

```go
// backend/api/internal/logic/verification/exportRecordsLogic.go

func (l *ExportRecordsLogic) ExportRecords(req *types.QueryVerificationRecordsReq) (string, error) {
    // 查询数据
    var records []VerificationRecord
    query := l.svcCtx.DB.Model(&VerificationRecord{})

    // 应用筛选条件
    if req.CampaignId > 0 {
        query = query.Where("campaign_id = ?", req.CampaignId)
    }
    if req.StartTime != "" {
        query = query.Where("verified_at >= ?", req.StartTime)
    }
    // ... 其他条件

    query.Find(&records)

    // 生成 Excel 文件
    file := excelize.NewFile()
    sheetName := "核销记录"
    file.NewSheet(sheetName)

    // 设置表头
    headers := []string{"订单ID", "用户", "核销码", "核销人", "核销方式", "核销时间", "状态"}
    for i, header := range headers {
        file.SetCellValue(sheetName, fmt.Sprintf("%c1", 'A'+i), header)
    }

    // 填充数据
    for i, record := range records {
        row := i + 2
        file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), record.OrderId)
        file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), record.UserName)
        file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), record.VerificationCode)
        file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), record.VerifiedByName)
        file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), record.VerificationMethod)
        file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), record.VerifiedAt)
        file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), record.Status)
    }

    // 保存文件
    fileName := fmt.Sprintf("核销记录_%s.xlsx", time.Now().Format("20060102_150405"))
    filePath := fmt.Sprintf("/tmp/%s", fileName)
    if err := file.SaveAs(filePath); err != nil {
        return "", err
    }

    return filePath, nil
}
```

### 优化 4: 表单字段类型扩展

**反馈内容**: "希望支持更多字段类型，如日期选择、文件上传等"

**优化方案**:

1. **新增字段类型**:

```sql
-- 更新表单字段类型
ALTER TABLE campaign_form_fields MODIFY COLUMN type ENUM(
    'text', 'phone', 'email', 'number', 'select',
    'textarea', 'address', 'date', 'time', 'datetime',
    'file', 'image', 'qrcode', 'checkbox', 'radio',
    'divider', 'header'
) NOT NULL;
```

2. **实现新字段类型的前端组件**:

```vue
<!-- DateField.vue -->
<template>
  <van-field
    :label="field.label"
    :required="field.required"
    readonly
    :placeholder="field.placeholder"
    @click="showDatePicker = true"
  >
    <template #input>
      <van-date-picker
        v-model="dateValue"
        :min-date="minDate"
        :max-date="maxDate"
        @confirm="onConfirm"
        v-if="showDatePicker"
      />
    </template>
  </van-field>
</template>

<!-- FileField.vue -->
<template>
  <van-field name="file" :label="field.label" :required="field.required">
    <template #input>
      <van-uploader
        v-model="files"
        :max-count="field.maxCount"
        :max-size="field.maxSize"
        accept="image/*"
      />
    </template>
  </van-field>
</template>
```

3. **支持拖拽式表单设计**:

```vue
<!-- FormDesigner.vue -->
<template>
  <div class="form-designer">
    <div class="field-library">
      <h3>字段库</h3>
      <draggable
        :list="fieldTypes"
        :group="{ name: 'fields', pull: 'clone', put: false }"
        item-key="type"
      >
        <template #item="{ element: field }">
          <div class="field-type-item">
            <i :class="field.icon"></i>
            <span>{{ field.name }}</span>
          </div>
        </template>
      </draggable>
    </div>

    <div class="form-canvas">
      <h3>表单设计区</h3>
      <draggable
        v-model="formFields"
        group="fields"
        item-key="id"
      >
        <template #item="{ element: field }">
          <FormField
            :field="field"
            @edit="editField"
            @delete="deleteField"
          />
        </template>
      </draggable>
    </div>

    <div class="field-properties">
      <h3>字段属性</h3>
      <FieldProperties
        v-if="selectedField"
        :field="selectedField"
        @update="updateField"
      />
    </div>
  </div>
</template>
```

***

## 📋 反馈管理流程

### 1. 反馈接收流程

```
用户提交反馈
    ↓
自动分类和优先级评估
    ↓
通知相关负责人
    ↓
进入处理队列
```

### 2. 反馈处理流程

```
新建反馈
    ↓
评估和确认 (1-2 工作日)
    ↓
分配给开发/产品
    ↓
实现和测试
    ↓
发布上线
    ↓
通知用户反馈已解决
```

### 3. 反馈跟踪

**创建反馈跟踪表**:

```sql
CREATE TABLE IF NOT EXISTS feedback_tracking (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    feedback_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL COMMENT '状态：new, assigned, in_progress, resolved, closed',
    assigned_to BIGINT COMMENT '负责人ID',
    assigned_at DATETIME,
    resolved_at DATETIME,
    comment TEXT COMMENT '处理备注',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_feedback_id (feedback_id),
    INDEX idx_status (status),
    FOREIGN KEY (feedback_id) REFERENCES user_feedback(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='反馈跟踪表';
```

***

## ✅ 检查清单

### 反馈收集

* \[ ] 应用内反馈入口已添加
* \[ ] 反馈表单已实现
* \[ ] 反馈 API 已开发
* \[ ] 反馈数据表已创建
* \[ ] 问卷调查已准备
* \[ ] 用户访谈计划已制定

### 反馈分析

* \[ ] 反馈分类汇总已完成
* \[ ] 优先级评估已完成
* \[ ] 趋势分析已完成
* \[ ] 功能使用统计已完成

### 优化实施

* \[ ] 优化方案已设计
* \[ ] 开发计划已制定
* \[ ] 优化进度已跟踪
* \[ ] 优化效果已评估

***

**文档状态**: 待执行
**最后更新**: 2026-02-01
**负责人**: 产品团队
