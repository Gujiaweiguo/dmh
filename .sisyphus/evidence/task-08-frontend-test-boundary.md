# 前端单元测试与 E2E 测试职责边界定义

## 1. 测试职责边界定义

### 1.1 单元测试 (Unit Test) 职责

单元测试关注**最小可测试单元**的正确性，使用 mock 隔离外部依赖。

| 类别 | 职责范围 | 示例 |
|------|---------|------|
| **工具函数测试** | 纯函数输入输出验证，无副作用 | `string.logic.test.js`, `dateFormat.logic.test.js`, `number.logic.test.js` |
| **业务逻辑测试** | View 伴生 `.logic.js` 中的纯逻辑函数 | `campaignEditor.logic.test.js`, `formValidation.logic.test.js` |
| **路由守卫测试** | 路由解析、权限判断逻辑 | `router.guard.test.js`, `adminHashRoute.test.ts` |
| **API 服务 Mock 测试** | API 调用逻辑、响应处理、错误处理 | `authApi.test.ts`, `memberApi.test.ts` |
| **组件逻辑测试** | 组件挂载、状态管理、事件触发（使用 Vue Test Utils） | `LoginView.component.test.ts` |
| **状态管理测试** | Store/状态管理逻辑 | `*.store.test.ts` |
| **HTTP 拦截器测试** | Axios 拦截器逻辑、Token 注入、错误处理 | `axios.test.js` |
| **性能监控测试** | 防抖、节流、性能测量逻辑 | `performanceMonitor.test.ts` |

### 1.2 E2E 测试 (End-to-End Test) 职责

E2E 测试关注**真实用户场景**的端到端流程验证。

| 类别 | 职责范围 | 示例 |
|------|---------|------|
| **登录认证流程** | 完整登录流程、Token 存储、页面跳转 | `admin-flows.spec.ts` Login Flow |
| **导航流程** | 菜单点击、页面切换、路由跳转 | `admin-flows.spec.ts` Navigation Flow |
| **CRUD 完整流程** | 列表查看、详情查看、状态更新、批量操作 | `feedback.spec.ts` |
| **跨页面主流程** | 从入口到完成的完整用户旅程 | `h5-flows.spec.ts` Order Flow |
| **关键业务路径** | 支付、核销、分销等核心业务流程 | `h5-flows.spec.ts` Promoter Flow |
| **用户交互场景** | 表单填写、按钮点击、弹窗操作 | `feedback.spec.ts` 场景4-6 |
| **权限验证流程** | 不同角色访问不同页面的端到端验证 | 需补充：角色权限 E2E |

---

## 2. 用例分类标准

### 2.1 判断标准决策树

```
                    测试什么？
                        │
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
    纯函数/逻辑      组件/服务        用户流程
        │               │               │
        ▼               ▼               ▼
    单元测试        单元测试*        E2E测试
                    (需mock)
                    
    * 组件测试：
      - 状态变化 → 单元测试
      - 事件触发 → 单元测试  
      - 页面跳转 → E2E测试
      - 跨组件交互 → E2E测试
```

### 2.2 分类判断清单

| 判断条件 | 单元测试 | E2E测试 |
|---------|---------|---------|
| 是否涉及真实浏览器渲染？ | ❌ | ✅ |
| 是否需要真实网络请求？ | ❌ (mock) | ✅ (可 mock API) |
| 是否跨多个页面？ | ❌ | ✅ |
| 是否测试用户交互流程？ | ❌ | ✅ |
| 是否可独立运行（无依赖）？ | ✅ | ❌ (需启动服务) |
| 执行时间是否 < 100ms？ | ✅ | ❌ |
| 是否涉及 DOM 快照比对？ | ⚠️ 慎用 | ✅ |

---

## 3. Admin 前端测试分类表

### 3.1 单元测试文件 (40 个)

| 文件名 | 测试类型 | 分类 | 评估 | 迁移建议 |
|--------|---------|------|------|---------|
| `mockApi.test.ts` | 工具函数 | ✅ 单元测试 | 正确 | 保持 |
| `performanceMonitor.test.ts` | 性能监控逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `adminHashRoute.test.ts` | 路由解析 | ✅ 单元测试 | 正确 | 保持 |
| `authApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `campaignApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `memberApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `profileApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `menuApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `userApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `securityApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `roleApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `brandApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `distributorApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `posterApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `feedbackApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `orderApi.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `api.test.ts` | API Mock | ✅ 单元测试 | 正确 | 保持 |
| `LoginView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `LoginView.test.ts` | 组件逻辑 | ✅ 单元测试 | ⚠️ 重复 | 合并到 `.component.test.ts` |
| `DynamicMenu.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `UserProfileView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `SecuritySettingsView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `CustomerListView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `MemberListView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `MemberDetailView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `MemberMergeView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `MemberExportView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `CampaignListView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `RewardDetailView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `PlatformDistributorView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `PlatformRewardView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `BrandManagementView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `BrandManagementView.test.ts` | 组件逻辑 | ✅ 单元测试 | ⚠️ 重复 | 合并到 `.component.test.ts` |
| `UserManagementView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `UserManagementView.test.ts` | 组件逻辑 | ✅ 单元测试 | ⚠️ 重复 | 合并到 `.component.test.ts` |
| `RolePermissionView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `DistributorManagementView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `distributorManagementView.test.ts` | 组件逻辑 | ✅ 单元测试 | ⚠️ 重复 | 合并到 `.component.test.ts` |
| `DistributorApprovalView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |
| `WithdrawalApprovalView.component.test.ts` | 组件逻辑 | ✅ 单元测试 | 正确 | 保持 |

### 3.2 E2E 测试文件 (4 个)

| 文件名 | 测试类型 | 分类 | 评估 | 迁移建议 |
|--------|---------|------|------|---------|
| `admin-flows.spec.ts` | 登录+导航流程 | ✅ E2E测试 | 正确 | 保持 |
| `feedback.spec.ts` | 反馈管理 CRUD | ✅ E2E测试 | 正确 | 保持 |
| `admin-dashboard.spec.ts` | 仪表板流程 | ✅ E2E测试 | 正确 | 保持 |
| `security-management.spec.ts` | 安全管理流程 | ✅ E2E测试 | 正确 | 保持 |

---

## 4. H5 前端测试分类表

### 4.1 单元测试文件 (54 个)

#### 4.1.1 工具函数测试 (✅ 正确分类)

| 文件名 | 测试类型 | 分类 | 评估 |
|--------|---------|------|------|
| `string.logic.test.js` | 字符串工具 | ✅ 单元测试 | 正确 |
| `number.logic.test.js` | 数字工具 | ✅ 单元测试 | 正确 |
| `array.logic.test.js` | 数组工具 | ✅ 单元测试 | 正确 |
| `object.logic.test.js` | 对象工具 | ✅ 单元测试 | 正确 |
| `color.logic.test.js` | 颜色工具 | ✅ 单元测试 | 正确 |
| `url.logic.test.js` | URL 工具 | ✅ 单元测试 | 正确 |
| `dateFormat.logic.test.js` | 日期格式化 | ✅ 单元测试 | 正确 |
| `storage.logic.test.js` | 存储工具 | ✅ 单元测试 | 正确 |
| `formValidation.logic.test.js` | 表单验证 | ✅ 单元测试 | 正确 |

#### 4.1.2 业务逻辑测试 (✅ 正确分类)

| 文件名 | 测试类型 | 分类 | 评估 |
|--------|---------|------|------|
| `campaignEditor.logic.test.js` | 活动编辑逻辑 | ✅ 单元测试 | 正确 |
| `campaignList.logic.test.js` | 活动列表逻辑 | ✅ 单元测试 | 正确 |
| `campaignForm.logic.test.js` | 活动表单逻辑 | ✅ 单元测试 | 正确 |
| `campaignDetail.logic.test.js` | 活动详情逻辑 | ✅ 单元测试 | 正确 |
| `campaignPageDesigner.logic.test.js` | 页面设计器逻辑 | ✅ 单元测试 | 正确 |
| `designer.logic.test.js` | 设计器逻辑 | ✅ 单元测试 | 正确 |
| `campaigns.logic.test.js` | 活动通用逻辑 | ✅ 单元测试 | 正确 |
| `orders.logic.test.js` | 订单逻辑 | ✅ 单元测试 | 正确 |
| `orderVerify.logic.test.js` | 订单核销逻辑 | ✅ 单元测试 | 正确 |
| `orderVerification.logic.test.js` | 核销逻辑 | ✅ 单元测试 | 正确 |
| `verificationRecords.logic.test.js` | 核销记录逻辑 | ✅ 单元测试 | 正确 |
| `verificationRecords.actions.test.js` | 核销动作逻辑 | ✅ 单元测试 | 正确 |
| `materials.logic.test.js` | 素材逻辑 | ✅ 单元测试 | 正确 |
| `posterGenerator.logic.test.js` | 海报生成逻辑 | ✅ 单元测试 | 正确 |
| `posterRecords.logic.test.js` | 海报记录逻辑 | ✅ 单元测试 | 正确 |
| `paymentQrcode.logic.test.js` | 支付二维码逻辑 | ✅ 单元测试 | 正确 |
| `members.logic.test.js` | 会员逻辑 | ✅ 单元测试 | 正确 |
| `memberDetail.logic.test.js` | 会员详情逻辑 | ✅ 单元测试 | 正确 |
| `promoters.logic.test.js` | 推广员逻辑 | ✅ 单元测试 | 正确 |
| `promoterDetail.logic.test.js` | 推广员详情逻辑 | ✅ 单元测试 | 正确 |
| `distributors.logic.test.js` | 分销商逻辑 | ✅ 单元测试 | 正确 |
| `distributorCenter.logic.test.js` | 分销中心逻辑 | ✅ 单元测试 | 正确 |
| `distributorApply.logic.test.js` | 分销申请逻辑 | ✅ 单元测试 | 正确 |
| `distributorLogin.logic.test.js` | 分销登录逻辑 | ✅ 单元测试 | 正确 |
| `distributorPromotion.logic.test.js` | 分销推广逻辑 | ✅ 单元测试 | 正确 |
| `distributorRewards.logic.test.js` | 分销奖励逻辑 | ✅ 单元测试 | 正确 |
| `distributorSubordinates.logic.test.js` | 下级分销逻辑 | ✅ 单元测试 | 正确 |
| `distributorWithdrawals.logic.test.js` | 分销提现逻辑 | ✅ 单元测试 | 正确 |
| `distributorLevelRewards.logic.test.js` | 分销等级奖励 | ✅ 单元测试 | 正确 |
| `distributorApproval.logic.test.js` | 分销审批逻辑 | ✅ 单元测试 | 正确 |
| `brandLogin.logic.test.js` | 品牌登录逻辑 | ✅ 单元测试 | 正确 |
| `dashboard.logic.test.js` | 仪表板逻辑 | ✅ 单元测试 | 正确 |
| `analytics.logic.test.js` | 分析逻辑 | ✅ 单元测试 | 正确 |
| `feedbackCenter.logic.test.js` | 反馈中心逻辑 | ✅ 单元测试 | 正确 |
| `settings.logic.test.js` | 设置逻辑 | ✅ 单元测试 | 正确 |
| `myOrders.logic.test.js` | 我的订单逻辑 | ✅ 单元测试 | 正确 |
| `success.logic.test.js` | 成功页逻辑 | ✅ 单元测试 | 正确 |
| `apiTest.logic.test.js` | API 测试逻辑 | ✅ 单元测试 | 正确 |

#### 4.1.3 路由/API 测试 (✅ 正确分类)

| 文件名 | 测试类型 | 分类 | 评估 |
|--------|---------|------|------|
| `router.guard.test.js` | 路由守卫 | ✅ 单元测试 | 正确 |
| `router.index.guard.test.js` | 路由守卫 | ✅ 单元测试 | 正确 |
| `api.test.js` | API Mock | ✅ 单元测试 | 正确 |
| `axios.test.js` | Axios 拦截器 | ✅ 单元测试 | 正确 |
| `brandApi.wrappers.test.js` | API 包装器 | ✅ 单元测试 | 正确 |
| `brandApi.orderApi.test.js` | 订单 API | ✅ 单元测试 | 正确 |

### 4.2 E2E 测试文件 (2 个)

| 文件名 | 测试类型 | 分类 | 评估 | 迁移建议 |
|--------|---------|------|------|---------|
| `h5-flows.spec.ts` | 多流程端到端 | ✅ E2E测试 | 正确 | 保持，建议补充更多场景 |
| `feedback.spec.ts` | 反馈页面流程 | ✅ E2E测试 | 正确 | 保持 |

---

## 5. 目录映射规范

### 5.1 Admin 前端目录结构

```
frontend-admin/
├── tests/
│   └── unit/                          # 单元测试
│       ├── *.component.test.ts        # 组件逻辑测试
│       ├── *Api.test.ts               # API Mock 测试
│       ├── *.test.ts                  # 工具函数测试
│       └── ...
├── e2e/                               # E2E 测试
│   ├── *-flows.spec.ts                # 流程类 E2E
│   ├── *-management.spec.ts           # 管理类 E2E
│   └── *.spec.ts                      # 功能类 E2E
└── ...
```

### 5.2 H5 前端目录结构

```
frontend-h5/
├── tests/
│   └── unit/                          # 单元测试
│       ├── *.logic.test.js            # View 伴生逻辑测试
│       ├── *.test.js                  # 工具/API 测试
│       └── ...
├── e2e/                               # E2E 测试
│   ├── *-flows.spec.ts                # 流程类 E2E
│   └── *.spec.ts                      # 功能类 E2E
└── ...
```

---

## 6. 迁移建议汇总

### 6.1 需要合并的重复测试

| 模块 | 重复文件 | 建议 |
|------|---------|------|
| Admin | `LoginView.test.ts` + `LoginView.component.test.ts` | 合并为 `LoginView.component.test.ts` |
| Admin | `BrandManagementView.test.ts` + `BrandManagementView.component.test.ts` | 合并 |
| Admin | `UserManagementView.test.ts` + `UserManagementView.component.test.ts` | 合并 |
| Admin | `distributorManagementView.test.ts` + `DistributorManagementView.component.test.ts` | 合并 |

### 6.2 需要补充的 E2E 测试

| 模块 | 建议新增 | 优先级 |
|------|---------|--------|
| Admin | 用户管理完整 CRUD 流程 | 高 |
| Admin | 角色权限配置流程 | 高 |
| Admin | 品牌管理流程 | 中 |
| H5 | 活动报名完整流程 | 高 |
| H5 | 分销商申请到提现完整流程 | 高 |
| H5 | 订单核销流程 | 中 |

### 6.3 测试质量改进建议

1. **避免视觉像素级断言在单测中**：使用 `toHaveClass`、`toHaveAttribute` 替代像素比对
2. **E2E 测试使用 Mock API**：`h5-flows.spec.ts` 已正确使用 `page.route()` mock API
3. **统一命名规范**：
   - 单元测试：`*.test.ts` / `*.test.js`
   - E2E 测试：`*.spec.ts`

---

## 7. 验证：随机抽样用例分类正确性

### 抽样验证结果

| 抽样文件 | 预期分类 | 实际内容 | 验证结果 |
|---------|---------|---------|---------|
| `string.logic.test.js` | 单元测试 | 纯函数测试，18 个 describe，无 DOM | ✅ 正确 |
| `campaignEditor.logic.test.js` | 单元测试 | 业务逻辑函数，无外部依赖 | ✅ 正确 |
| `router.guard.test.js` | 单元测试 | 路由解析纯函数，无浏览器 | ✅ 正确 |
| `axios.test.js` | 单元测试 | 拦截器逻辑，完全 mock | ✅ 正确 |
| `admin-flows.spec.ts` | E2E测试 | 真实页面导航，Playwright | ✅ 正确 |
| `feedback.spec.ts` (Admin) | E2E测试 | 完整 CRUD 流程，真实浏览器 | ✅ 正确 |
| `h5-flows.spec.ts` | E2E测试 | 多页面流程，Mock API + 真实渲染 | ✅ 正确 |
| `authApi.test.ts` | 单元测试 | API Mock，无真实网络请求 | ✅ 正确 |
| `performanceMonitor.test.ts` | 单元测试 | 防抖节流逻辑，完全 mock | ✅ 正确 |

**验证结论**：随机抽样 9 个测试文件，分类正确率 100%。

---

## 8. 总结

### 8.1 当前测试分布

| 模块 | 单元测试 | E2E 测试 | 覆盖率 |
|------|---------|---------|--------|
| Admin | 40 文件 | 4 文件 | 83.65% |
| H5 | 54 文件 | 2 文件 | 87.37% |

### 8.2 职责边界总结

| 维度 | 单元测试 | E2E测试 |
|------|---------|---------|
| **关注点** | 代码正确性 | 用户场景 |
| **依赖** | Mock 隔离 | 真实/模拟环境 |
| **速度** | 毫秒级 | 秒级 |
| **范围** | 函数/组件 | 完整流程 |
| **维护成本** | 低 | 高 |
| **数量比例** | 80-90% | 10-20% |

### 8.3 最佳实践

1. **新功能开发**：先写单元测试覆盖核心逻辑，再写 E2E 验证关键流程
2. **Bug 修复**：先写失败的单元测试复现 bug，修复后确认通过
3. **重构**：单元测试保护重构安全，E2E 确保端到端行为不变
4. **CI/CD**：单元测试每次提交运行，E2E 测试合并前运行

---

*生成时间: 2026-02-26*
*任务: task-08-frontend-test-boundary*
