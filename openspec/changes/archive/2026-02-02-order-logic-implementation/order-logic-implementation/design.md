# Design: 订单相关业务逻辑实现

## Context

本设计文档详细说明了如何补充订单相关的业务逻辑实现。根据端到端测试发现的问题，主要需要补充以下内容：

1. 订单创建逻辑（`createOrderLogic.go`）
2. 订单核销逻辑（`verifyOrderLogic.go`, `unverifyOrderLogic.go`, `scanOrderLogic.go`）
3. 表单字段验证服务

**相关文档**:
- Spec: `openspec/specs/002-order-payment-system.md`
- Proposal: `openspec/changes/order-logic-implementation/proposal.md`
- 端到端测试报告: `/tmp/E2E_TEST_REPORT.md`

---

## Goals / Non-Goals

### Goals
- 实现完整的订单创建业务逻辑
- 完善订单核销相关业务逻辑
- 实现表单字段验证服务
- 确保事务处理正确
- 确保数据一致性
- 提供完整的单元测试覆盖

### Non-Goals
- 不修改现有的 API 接口定义
- 不修改现有的数据库表结构
- 不实现支付功能（已在高级功能中实现）
- 不实现外部同步适配器（已存在）

---

## Decisions

### Decision 1: 订单创建流程

**选择**: 实现完整的订单创建流程，包括验证、存储、推荐人处理

**流程**:
```
1. 接收订单创建请求
2. 验证用户身份（JWT）
3. 验证活动是否存在且有效
4. 验证手机号格式
5. 防重复检查（campaign_id + phone）
6. 验证表单数据
7. 生成核销码（带签名）
8. 保存订单到数据库
9. 返回订单信息
```

**验证规则**:
- 活动状态：必须为 `active`
- 活动时间：必须在当前时间范围内
- 手机号：必须是 11 位数字
- 表单必填字段：必须全部填写
- 防重复：同一活动、同一手机号只能创建一次订单

**核销码生成算法**:
```
verification_code = {order_id}_{phone}_{timestamp}_{signature}
signature = md5(order_id + phone + timestamp + secret_key)
```

**替代方案**: 不做防重复检查
**拒绝理由**: 会导致重复订单，影响活动统计和奖励发放

---

### Decision 2: 订单核销流程

**选择**: 实现基于核销码的核销流程，包含权限验证和操作日志

**核销流程**:
```
1. 接收核销请求
2. 验证用户身份（JWT + 角色检查）
3. 验证订单是否存在
4. 验证核销码签名
5. 验证订单状态（必须是 pending 或 unverified）
6. 更新订单状态为 verified
7. 记录核销时间、核销人
8. 记录核销操作日志
9. 返回核销成功
```

**取消核销流程**:
```
1. 接收取消核销请求
2. 验证用户身份（JWT + 角色检查）
3. 验证订单是否存在
4. 验证订单状态（必须是 verified）
5. 更新订单状态为 unverified
6. 清除核销时间、核销人
7. 记录取消核销操作日志
8. 返回成功
```

**权限规则**:
- 只有品牌管理员可以核销订单
- 只能核销自己品牌下的活动订单

**替代方案**: 核销不需要权限验证
**拒绝理由**: 会存在安全风险，非品牌管理员可能核销他人订单

---

### Decision 3: 表单字段验证

**选择**: 创建独立的表单字段验证服务，支持多种字段类型和验证规则

**验证规则**:
```
text: 非空字符串
phone: 11 位数字，符合手机号格式
email: 符合邮箱格式
number: 必须是数字
textarea: 非空字符串
address: 长度 10-200 字符
select: 值必须在 options 中
checkbox: 至少选择一个选项
```

**实现方式**:
- 创建 `FormFieldValidator` 结构体
- 为每种字段类型实现验证方法
- 返回详细的验证错误信息

**替代方案**: 在每个 Logic 中直接实现验证
**拒绝理由**: 代码重复，难以维护，验证规则不统一

---

### Decision 4: 事务处理

**选择**: 使用数据库事务确保数据一致性

**事务边界**:
```
开始事务
  ├─ 验证订单状态
  ├─ 更新订单状态
  ├─ 记录核销操作日志（如需要）
  └─ 其他相关操作
提交事务
```

**错误处理**:
```
出现错误
  ├─ 回滚事务
  ├─ 记录错误日志
  └─ 返回错误信息
```

**替代方案**: 不使用事务，分别执行数据库操作
**拒绝理由**: 无法保证数据一致性，可能导致数据不一致

---

## Implementation Details

### 1. 订单创建逻辑实现

**文件**: `backend/api/internal/logic/order/createOrderLogic.go`

**关键方法**:
```go
func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (*types.OrderResp, error) {
    // 1. 获取用户信息
    userId := l.ctx.Value("userId").(int64)

    // 2. 验证活动
    campaign, err := l.validateCampaign(req.CampaignId)
    if err != nil {
        return nil, err
    }

    // 3. 防重复检查
    if err := l.checkDuplicate(req.CampaignId, req.Phone); err != nil {
        return nil, err
    }

    // 4. 验证表单数据
    formFields, err := l.getFormFields(campaign.Id)
    if err != nil {
        return nil, err
    }
    if err := l.validateFormData(req.FormData, formFields); err != nil {
        return nil, err
    }

    // 5. 生成核销码
    verificationCode := l.generateVerificationCode(campaign.Id, req.Phone)

    // 6. 创建订单
    order := &model.Order{
        CampaignId:   campaign.Id,
        Phone:        req.Phone,
        FormData:     req.FormData,
        UserId:       userId,
        ReferrerId:   req.ReferrerId,
        Status:       "pending",
        VerificationCode: verificationCode,
    }

    if err := l.svcCtx.DB.Create(order).Error; err != nil {
        l.Errorf("创建订单失败: %v", err)
        return nil, err
    }

    // 7. 返回订单信息
    return &types.OrderResp{
        Id:             order.Id,
        CampaignId:     order.CampaignId,
        Phone:          order.Phone,
        Status:          order.Status,
        VerificationCode: order.VerificationCode,
        CreatedAt:      order.CreatedAt,
    }, nil
}
```

**辅助方法**:
- `validateCampaign(campaignId int64)` - 验证活动有效性
- `checkDuplicate(campaignId int64, phone string)` - 防重复检查
- `getFormFields(campaignId int64)` - 获取表单字段配置
- `validateFormData(formData string, formFields []model.CampaignFormField)` - 验证表单数据
- `generateVerificationCode(campaignId int64, phone string)` - 生成核销码

---

### 2. 订单核销逻辑实现

**文件**: `backend/api/internal/logic/order/verifyOrderLogic.go`

**关键方法**:
```go
func (l *VerifyOrderLogic) VerifyOrder(req *types.VerifyOrderReq) (*types.VerifyOrderResp, error) {
    // 1. 获取用户信息和角色
    userId := l.ctx.Value("userId").(int64)
    userRole := l.ctx.Value("userRole").(string)

    // 2. 权限检查
    if userRole != "brand_admin" {
        return nil, errors.New("权限不足")
    }

    // 3. 查询订单
    var order model.Order
    if err := l.svcCtx.DB.Where("id = ?", req.OrderId).First(&order).Error; err != nil {
        return nil, err
    }

    // 4. 验证核销码签名
    if !l.verifyVerificationCode(order.VerificationCode, req.Code) {
        return nil, errors.New("核销码无效")
    }

    // 5. 验证订单状态
    if order.VerificationStatus == "verified" {
        return nil, errors.New("订单已核销")
    }

    // 6. 使用事务更新订单状态
    tx := l.svcCtx.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 更新订单状态
    now := time.Now()
    if err := tx.Model(&order).Updates(map[string]interface{}{
        "verification_status": "verified",
        "verified_at":        now,
        "verified_by":        userId,
        "verification_method": "manual",
    }).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // 记录核销日志
    log := &model.VerificationRecord{
        OrderId:       order.Id,
        CampaignId:    order.CampaignId,
        UserId:        order.UserId,
        Phone:         order.Phone,
        VerificationCode: order.VerificationCode,
        Status:         "verified",
        VerifiedBy:      userId,
        VerifiedAt:      now,
        VerificationMethod: "manual",
    }
    if err := tx.Create(log).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return &types.VerifyOrderResp{
        Id:             order.Id,
        VerificationStatus: "verified",
        VerifiedAt:      now.Format("2006-01-02T15:04:05"),
    }, nil
}
```

**辅助方法**:
- `verifyVerificationCode(storedCode, inputCode)` - 验证核销码签名

---

### 3. 表单字段验证服务实现

**文件**: `backend/api/internal/service/form_field_service.go`

**验证器结构**:
```go
type FormFieldValidator struct {
    validators map[string]func(value string, field *model.CampaignFormField) error
}

func NewFormFieldValidator() *FormFieldValidator {
    v := &FormFieldValidator{
        validators: make(map[string]func(string, *model.CampaignFormField) error),
    }

    // 注册验证器
    v.validators["text"] = v.validateText
    v.validators["phone"] = v.validatePhone
    v.validators["email"] = v.validateEmail
    v.validators["number"] = v.validateNumber
    v.validators["textarea"] = v.validateTextarea
    v.validators["address"] = v.validateAddress

    return v
}

func (v *FormFieldValidator) Validate(value string, field *model.CampaignFormField) error {
    if validator, ok := v.validators[field.Type]; ok {
        return validator(value, field)
    }
    return fmt.Errorf("不支持的字段类型: %s", field.Type)
}
```

---

## Risks / Trade-offs

### Risk 1: 核销码安全性

**风险**: 核销码可能被伪造

**缓解措施**:
- 使用签名机制（MD5(order_id + phone + timestamp + secret_key)）
- 核销码包含时间戳，设置有效期（如 24 小时）
- 验证时检查时间戳是否在有效期内

---

### Risk 2: 并发订单创建

**风险**: 同一用户可能同时发起多个订单创建请求

**缓解措施**:
- 数据库唯一索引约束：`UNIQUE KEY uk_campaign_phone (campaign_id, phone, deleted_at)`
- 创建前先查询是否存在
- 使用事务确保原子性

---

### Risk 3: 核销操作并发

**风险**: 多个管理员同时核销同一订单

**缓解措施**:
- 使用数据库事务
- 使用乐观锁（在 order 表添加 version 字段）
- 在更新前验证订单状态

---

## Migration Plan

### 阶段 1: 数据库准备（无需变更）
- 验证 orders 表结构
- 验证索引是否完整
- 验证核销记录表是否存在

### 阶段 2: 代码实现
- 实现订单创建逻辑
- 实现订单核销逻辑
- 实现表单字段验证服务

### 阶段 3: 测试
- 单元测试
- 集成测试
- 端到端测试

### 阶段 4: 部署
- 代码审查
- 合并代码
- 部署到生产环境

### Rollback Plan
如果出现问题：
1. 回滚到修改前的代码版本
2. 数据库无需回滚（无结构变更）
3. 保留已创建的订单数据（不影响现有功能）

---

## Open Questions

1. **Q**: 核销码有效期设置为多长？
   **A**: 建议 24 小时，可在配置文件中配置

2. **Q**: 是否需要实现批量核销功能？
   **A**: 当前不实现，后续可根据需求添加

3. **Q**: 核销记录是否需要删除？
   **A**: 不删除，永久保留用于审计

4. **Q**: 表单字段验证规则是否需要可配置？
   **A**: 当前硬编码，后续可根据需求改为可配置

---

## Testing Strategy

### 单元测试覆盖
- 订单创建：正常流程、重复订单、活动不存在、活动已结束
- 订单核销：正常核销、重复核销、权限不足、核销码无效
- 表单验证：所有字段类型的正常和异常情况

### 集成测试覆盖
- 完整的订单创建 → 核销流程
- 多个并发订单创建
- 多个并发核销操作

### 端到端测试覆盖
- 使用现有的测试脚本验证所有功能

---

## Dependencies

### 外部依赖
- 无新增外部依赖

### 内部依赖
- `model.Order` - 订单模型
- `model.Campaign` - 活动模型
- `model.CampaignFormField` - 表单字段模型
- `model.VerificationRecord` - 核销记录模型
- `model.User` - 用户模型
- `model.UserRole` - 用户角色模型

---

## Success Criteria

- [ ] 订单创建功能正常工作
- [ ] 订单核销功能正常工作
- [ ] 防重复检查生效
- [ ] 权限控制生效
- [ ] 核销码验证生效
- [ ] 表单字段验证生效
- [ ] 所有单元测试通过
- [ ] 所有集成测试通过
- [ ] 端到端测试通过
- [ ] 代码审查通过
- [ ] 部署到生产环境
