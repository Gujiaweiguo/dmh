# Change: 订单相关业务逻辑实现

## 概述

本变更旨在补充 DMH 订单相关的业务逻辑实现，解决端到端测试中发现的问题。

**问题**:
1. `backend/api/internal/logic/order/createOrderLogic.go` 只有 TODO 注释，没有实际业务逻辑
2. 订单核销功能依赖订单创建功能，导致无法完成端到端测试
3. 表单字段验证逻辑需要独立实现

**解决方案**:
1. 实现完整的订单创建业务逻辑
2. 完善订单核销相关业务逻辑
3. 创建独立的表单字段验证服务

---

## Why

在高级功能的端到端测试（`E2E_TEST_REPORT.md`）中发现：

**关键问题 1**: 订单创建逻辑未实现
- 代码位置：`backend/api/internal/logic/order/createOrderLogic.go`
- 当前状态：只有 TODO 注释
- 影响：无法完成订单创建 → 支付 → 核销的完整流程测试

**关键问题 2**: 订单核销依赖缺失
- 虽然 API 接口已定义，但由于缺少订单数据，无法验证核销功能
- 影响：高级功能的订单核销功能无法投入使用

**问题 3**: 表单字段验证需要独立服务
- 前端已支持多种表单字段类型（email, textarea, address）
- 需要后端提供统一的验证服务
- 影响：无法保证表单数据质量

---

## What Changes

### 1. 订单创建逻辑实现

**文件**: `backend/api/internal/logic/order/createOrderLogic.go`

**新增功能**:
- 活动有效性验证
- 手机号格式验证
- 表单数据验证和解析
- 防重复检查（基于 campaign_id + phone）
- 核销码生成（包含签名）
- 订单数据存储
- 推荐人信息处理

### 2. 订单核销逻辑完善

**文件**:
- `backend/api/internal/logic/order/verifyOrderLogic.go` - 核销订单
- `backend/api/internal/logic/order/unverifyOrderLogic.go` - 取消核销
- `backend/api/internal/logic/order/scanOrderLogic.go` - 扫码获取订单

**新增功能**:
- 核销码签名验证
- 订单状态更新（verified/unverified）
- 核销时间和核销人记录
- 核销权限验证（仅品牌管理员）
- 核销操作日志记录

### 3. 表单字段验证服务

**文件**: `backend/api/internal/service/form_field_service.go`

**新增功能**:
- 多种字段类型验证器（text, phone, email, number, textarea, address, select, checkbox）
- 统一的验证错误返回
- 配置化的验证规则

### 4. 单元测试补充

**文件**:
- `backend/api/internal/logic/order/createOrderLogic_test.go`
- `backend/api/internal/logic/order/verifyOrderLogic_test.go`
- `backend/api/internal/service/form_field_service_test.go`

---

## Impact

### 影响的模块

- **订单支付系统** (`specs/002-order-payment-system.md`)
  - 补充订单创建业务逻辑
  - 完善订单核销业务逻辑

- **活动管理模块** (`specs/001-campaign-management.md`)
  - 订单创建依赖活动验证
  - 表单字段依赖活动配置

### API 变更

**已存在的 API（需补充实现）**:
- `POST /api/v1/orders` - 订单创建
- `POST /api/v1/orders/:id/verify` - 核销订单
- `POST /api/v1/orders/:id/unverify` - 取消核销
- `GET /api/v1/orders/scan/:code` - 扫码获取订单

### 数据库变更

**无结构变更**（orders 表结构已存在）**

### 新增文件

1. `backend/api/internal/logic/order/createOrderLogic.go` - 补充实现
2. `backend/api/internal/logic/order/verifyOrderLogic.go` - 完善实现
3. `backend/api/internal/logic/order/unverifyOrderLogic.go` - 完善实现
4. `backend/api/internal/logic/order/scanOrderLogic.go` - 完善实现
5. `backend/api/internal/service/form_field_service.go` - 新增验证服务
6. 4 个测试文件

---

## Breaking Changes

无破坏性变更。所有变更都是补充现有功能的业务逻辑实现，不影响现有接口和数据结构。

---

## Migration Plan

1. **代码实现**（1-2 周）
   - 实现订单创建逻辑
   - 完善订单核销逻辑
   - 实现表单字段验证服务

2. **单元测试**（1 周）
   - 编写订单创建测试
   - 编写订单核销测试
   - 编写表单字段验证测试

3. **集成测试**（1 周）
   - 测试完整的订单流程
   - 测试订单核销流程
   - 验证表单字段验证集成

4. **端到端验证**（0.5 天）
   - 重新运行端到端测试脚本
   - 验证所有功能正常工作

5. **代码审查**（1 天）
   - 代码风格和规范检查
   - 错误处理和日志记录检查
   - 安全性审查

---

## Rollback Plan

如果出现问题，可以：
1. 回滚到修改前的代码版本（git）
2. 保留数据库结构（无结构变更）
3. 新增的测试代码可以保留，不影响现有功能
4. 删除或回滚新创建的 logic 文件

---

## Timeline

- **第 1 周**: 订单创建逻辑实现 + 单元测试
- **第 2 周**: 订单核销逻辑完善 + 单元测试
- **第 3 周**: 表单验证服务实现 + 单元测试
- **第 4 周**: 集成测试 + 端到端验证

---

## Success Metrics

- [ ] 订单创建功能测试通过（正常流程 + 异常流程）
- [ ] 订单核销功能测试通过（核销、取消核销、扫码获取）
- [ ] 表单字段验证测试通过（所有字段类型）
- [ ] 完整的订单流程测试通过
- [ ] 代码覆盖率达到 80% 以上
- [ ] 所有单元测试通过
- [ ] 所有集成测试通过
- [ ] 端到端测试通过（订单创建 → 核销）

---

## Dependencies

- 需要先确保 `orders` 表结构正确
- 需要先确保 `campaigns` 表数据完整
- 需要先确保 `users` 表数据完整
- 需要先确保 `campaign_form_fields` 表数据完整

---

## Notes

1. 本变更专注于补充业务逻辑实现，不涉及架构调整
2. 所有实现需遵循现有的代码规范和模式
3. 测试覆盖要全面，包括正常和异常情况
4. 注意处理并发场景（订单创建和核销）
5. 注意数据一致性（事务处理）
6. 遵循现有的错误处理和日志记录机制

---

**相关文档**:
- [Spec: 订单与支付系统](../specs/002-order-payment-system.md)
- [提案文档](../changes/order-logic-implementation/proposal.md)
- [设计文档](../changes/order-logic-implementation/design.md)
- [任务清单](../changes/order-logic-implementation/tasks.md)
- [端到端测试报告](/tmp/E2E_TEST_REPORT.md)

---

**变更状态**: ✅ 已批准
**审批日期**: 2026-02-01
**审批结果**: 通过（补充性业务逻辑实现，无架构变更）
