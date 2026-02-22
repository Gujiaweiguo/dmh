# DMH 发布前全量回归测试计划

> **版本**: v1.1
> **创建日期**: 2026-02-15
> **更新日期**: 2026-02-15
> **目标**: 完成所有历史变更的测试闭环，确保发布质量
> **环境**: Docker Compose 容器化环境

---

## 0. 环境配置 (Docker Compose)

### 0.1 容器状态验证

**执行**:
```bash
cd /opt/code/dmh/deploy
docker compose -f docker-compose-simple.yml ps
```

**预期输出**:
```
NAMES            STATUS                          PORTS
dmh-nginx        Up                              0.0.0.0:3000->3000/tcp, 0.0.0.0:3100->3100/tcp
dmh-api          Up                              0.0.0.0:8889->8889/tcp
mysql8           Up (healthy)                    0.0.0.0:3306->3306/tcp
redis-dmh        Up (healthy)                    0.0.0.0:6379->6379/tcp
```

### 0.2 Playwright Docker 环境配置

**需要修改的文件**:
- `frontend-admin/playwright.config.ts`
- `frontend-h5/playwright.config.ts`

**修改内容** (添加 `DOCKER_ENV` 支持):

```typescript
// frontend-admin/playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  timeout: 30 * 1000,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    // 支持自定义 baseURL，Docker 环境可通过环境变量覆盖
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    // Docker 环境或 CI 环境使用 headless 模式
    headless: process.env.CI === 'true' || process.env.DOCKER_ENV === 'true',
  },

  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        channel: 'chrome', // Use system chrome browser
      },
    },
  ],

  // Docker 环境下禁用 webServer（使用 nginx 托管的静态文件）
  // 使用命令: DOCKER_ENV=true npm run test:e2e
  webServer: process.env.DOCKER_ENV === 'true' ? undefined : {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
  },
});
```

```typescript
// frontend-h5/playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  timeout: 30 * 1000,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    // 支持自定义 baseURL，Docker 环境可通过环境变量覆盖
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3100',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    // Docker 环境或 CI 环境使用 headless 模式
    headless: process.env.CI === 'true' || process.env.DOCKER_ENV === 'true',
  },

  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        channel: 'chrome',
      },
    },
  ],

  // Docker 环境下禁用 webServer（使用 nginx 托管的静态文件）
  // 使用命令: DOCKER_ENV=true npm run test:e2e
  webServer: process.env.DOCKER_ENV === 'true' ? undefined : {
    command: 'npm run dev',
    url: 'http://localhost:3100',
    reuseExistingServer: !process.env.CI,
  },
});
```

**Docker 环境 E2E 测试命令**:
```bash
# Admin E2E (Docker 环境)
cd frontend-admin && DOCKER_ENV=true npm run test:e2e

# H5 E2E (Docker 环境)
cd frontend-h5 && DOCKER_ENV=true npm run test:e2e
```

---

## 0.3 执行任务清单

> **注意**: 以下任务需要通过 `/start-work` 执行

- [ ] **Task 0.1**: 修改 `frontend-admin/playwright.config.ts`
  - 添加 `DOCKER_ENV` 环境变量支持
  - 条件禁用 `webServer` 配置
  - 见上文完整代码

- [ ] **Task 0.2**: 修改 `frontend-h5/playwright.config.ts`
  - 添加 `DOCKER_ENV` 环境变量支持
  - 条件禁用 `webServer` 配置
  - 见上文完整代码

- [ ] **Task 0.3**: 验证 E2E 测试在 Docker 环境运行
  - 执行: `cd frontend-admin && DOCKER_ENV=true npm run test:e2e`
  - 执行: `cd frontend-h5 && DOCKER_ENV=true npm run test:e2e`

---

## 1. 需求基线汇总策略

### 1.1 基线来源

| 来源 | 文件路径 | 内容 |
|------|----------|------|
| 需求总览 | `all-requirements.md` | 5 个活跃规格 + 16 个归档变更 |
| 测试计划 | `test-plan.md` | 50+ 测试用例（正常 + 异常场景） |
| 规格目录 | `openspec/specs/*/spec.md` | 30 需求 / 96 场景 |
| 执行矩阵 | `docs/testing/execution/*.md` | 执行入口与验收标准 |

### 1.2 规格覆盖矩阵

| 规格 | 需求数 | 场景数 | P0 需求 | P1 需求 | 测试资产状态 |
|------|--------|--------|---------|---------|--------------|
| `rbac-permission-system` | 12 | 38 | 8 | 4 | ✅ 单元/集成/E2E 完备 |
| `campaign-management` | 6 | 24 | 4 | 2 | ✅ E2E 完备 |
| `order-payment-system` | 1 | 3 | 1 | 0 | ✅ 回归脚本完备 |
| `feedback-system` | 9 | 29 | 0 | 9 | ⚠️ 集成测试部分 |
| `spec-governance` | 2 | 2 | 0 | 2 | ✅ OpenSpec 校验 |
| **总计** | **30** | **96** | **13** | **17** | - |

### 1.3 历史变更追溯

从 16 个归档变更中提取关键回归点：

| 归档日期 | 变更ID | 回归重点 | 风险等级 |
|----------|--------|----------|----------|
| 2026-02-04 | dmh-mvp-core-features | MVP 核心功能全量回归 | 🔴 高 |
| 2026-02-10 | add-feedback-system | 反馈系统端到端验证 | 🟡 中 |
| 2026-02-07 | add-order-logic-test-gap-closure | 订单创建/核销闭环 | 🔴 高 |
| 2026-01-28 | add-distributor-role | 分销商权限边界 | 🟡 中 |
| 2026-01-28 | add-member-system | 会员合并/导出审批 | 🟢 低 |
| 2026-02-07 | add-brand-admin-poster-distribution | 海报生成/分销配置 | 🟢 低 |

### 1.4 需求基线同步策略

```
┌─────────────────────────────────────────────────────────────┐
│                    需求基线同步流程                           │
├─────────────────────────────────────────────────────────────┤
│  1. all-requirements.md (权威来源)                           │
│     ↓                                                        │
│  2. test-plan.md (测试范围映射)                              │
│     ↓                                                        │
│  3. docs/testing/execution/SCOPE_MAPPING.md (执行入口)       │
│     ↓                                                        │
│  4. 验证报告 (verification-report.md)                        │
│     ↓                                                        │
│  5. 回归结论 → 发布决策                                       │
└─────────────────────────────────────────────────────────────┘
```

**同步触发条件**：
- OpenSpec 新增/修改规格
- 归档变更后
- 发布前回归前

---

## 2. 系统级测试计划与分批执行策略

### 2.1 分批原则

| 批次 | 优先级 | 目标 | 时间窗口 | 通过标准 |
|------|--------|------|----------|----------|
| **Batch 1** | P0 核心链路 | 认证/订单/活动 | 30min | 100% 通过 |
| **Batch 2** | P0 扩展 | RBAC 完整覆盖 | 20min | 100% 通过 |
| **Batch 3** | P1 业务 | 反馈/分销/会员 | 25min | ≥95% 通过 |
| **Batch 4** | 规格治理 | OpenSpec 校验 | 5min | 100% 通过 |

### 2.2 Batch 1: P0 核心链路 (必须全部通过)

**执行顺序**: RBAC 认证 → 活动管理 → 订单支付

```
┌─────────────────────────────────────────────────────────────┐
│                    Batch 1 执行流程                          │
├─────────────────────────────────────────────────────────────┤
│  Phase 1.1: 后端服务启动与健康检查 (5min)                     │
│    ├── docker compose up -d                                  │
│    ├── 健康检查: curl http://localhost:8889/health           │
│    └── 数据库连接验证                                         │
│                                                              │
│  Phase 1.2: RBAC 认证授权 (10min)                            │
│    ├── TC-RBAC-01: H5 用户注册                               │
│    ├── TC-RBAC-05: 用户登录 JWT 校验                         │
│    ├── TC-RBAC-06: Token 过期处理                            │
│    └── TC-RBAC-07: 品牌管理员数据隔离                        │
│                                                              │
│  Phase 1.3: 活动管理 (10min)                                 │
│    ├── TC-CAMP-01: 页面设计器访问                            │
│    ├── TC-CAMP-02~04: 组件操作                               │
│    ├── TC-CAMP-05: 海报生成                                  │
│    └── TC-CAMP-09: 分销规则配置                              │
│                                                              │
│  Phase 1.4: 订单支付与核销 (10min)                           │
│    ├── TC-ORDER-01: 订单创建                                 │
│    ├── TC-ORDER-02: 重复报名拦截                             │
│    ├── TC-ORDER-05: 订单核销                                 │
│    ├── TC-ORDER-07: 重复核销幂等                             │
│    └── TC-ORDER-08: 取消核销回滚                             │
└─────────────────────────────────────────────────────────────┘
```

**执行命令**:
```bash
# Phase 1.1: 环境启动
cd deploy && docker compose -f docker-compose-simple.yml up -d
docker ps | grep -E "dmh-api|dmh-nginx|mysql8|redis-dmh"
curl -s http://localhost:8889/api/health || echo "API not ready"

# Phase 1.2: RBAC 认证
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
  go test ./test/integration/... -v -run "TestAuth|TestLogin" -count=1

# Phase 1.3: 活动管理
cd frontend-h5 && npm run test:e2e -- --grep "campaign|designer"

# Phase 1.4: 订单支付
cd backend && ./scripts/run_order_mysql8_regression.sh
```

### 2.3 Batch 2: P0 扩展 - RBAC 完整覆盖

```
┌─────────────────────────────────────────────────────────────┐
│                    Batch 2 执行流程                          │
├─────────────────────────────────────────────────────────────┤
│  Phase 2.1: 用户管理 (8min)                                  │
│    ├── TC-RBAC-03: 平台管理员创建用户                        │
│    ├── TC-RBAC-04: 越权创建拦截                              │
│    ├── TC-RBAC-09: 提现审批                                  │
│    └── TC-RBAC-11: 安全审计日志                              │
│                                                              │
│  Phase 2.2: 菜单权限 (6min)                                  │
│    ├── TC-RBAC-11: 登录后返回菜单权限树                      │
│    └── 菜单权限继承验证                                      │
│                                                              │
│  Phase 2.3: 缓存与性能 (6min)                                │
│    ├── TC-RBAC-08: 权限变更后缓存失效                        │
│    └── TC-RBAC-13: 分销监控超时处理                          │
└─────────────────────────────────────────────────────────────┘
```

**执行命令**:
```bash
# Phase 2.1~2.3: 后端全量集成测试
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
  go test ./test/integration/... -v -count=1

# Admin E2E 覆盖
cd frontend-admin && npm run test:e2e
```

### 2.4 Batch 3: P1 业务功能

```
┌─────────────────────────────────────────────────────────────┐
│                    Batch 3 执行流程                          │
├─────────────────────────────────────────────────────────────┤
│  Phase 3.1: 反馈系统 (10min)                                 │
│    ├── TC-FEED-01: 用户提交反馈                              │
│    ├── TC-FEED-04: 反馈列表查询                              │
│    ├── TC-FEED-06: 管理员更新反馈状态                        │
│    ├── TC-FEED-08: 满意度调查                                │
│    └── TC-FEED-11: 反馈统计分析                              │
│                                                              │
│  Phase 3.2: 分销商系统 (8min)                                │
│    ├── 分销商角色验证                                        │
│    ├── 多级分销奖励计算                                      │
│    └── 提现申请与审批                                        │
│                                                              │
│  Phase 3.3: 会员系统 (7min)                                  │
│    ├── UnionID 唯一性验证                                    │
│    ├── 会员合并流程                                          │
│    └── 导出审批流程                                          │
└─────────────────────────────────────────────────────────────┘
```

**执行命令**:
```bash
# 反馈系统集成测试
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
  go test ./test/integration/... -v -run "Feedback" -count=1

# H5 E2E 完整覆盖
cd frontend-h5 && npm run test:e2e

# 性能压测（可选）
cd backend && go test ./test/performance/... -v -count=1
```

### 2.5 Batch 4: 规格治理

```
┌─────────────────────────────────────────────────────────────┐
│                    Batch 4 执行流程                          │
├─────────────────────────────────────────────────────────────┤
│  Phase 4.1: OpenSpec 校验 (3min)                            │
│    ├── openspec validate --all --strict --no-interactive    │
│    └── 确认 5 个规格 + 1 个活跃变更                          │
│                                                              │
│  Phase 4.2: 归档一致性 (2min)                                │
│    ├── 检查归档状态索引                                      │
│    └── tasks.md 状态一致性                                   │
└─────────────────────────────────────────────────────────────┘
```

**执行命令**:
```bash
# OpenSpec 全量校验
openspec validate --all --strict --no-interactive

# 归档状态索引检查
cat openspec/changes/archive/ARCHIVE_STATUS_INDEX.md
```

---

## 3. 每阶段输入/输出与验收标准

### 3.1 阶段定义

| 阶段 | 输入 | 输出 | 验收标准 |
|------|------|------|----------|
| **S0: 准备** | 环境配置 | 服务就绪 | 所有服务健康 |
| **S1: Batch 1** | P0 用例 | 执行日志 | 100% 通过 |
| **S2: Batch 2** | P0 扩展 | 执行日志 | 100% 通过 |
| **S3: Batch 3** | P1 用例 | 执行日志 | ≥95% 通过 |
| **S4: Batch 4** | 规格文件 | 校验报告 | 100% 通过 |
| **S5: 汇总** | 所有日志 | verification-report.md | 发布决策 |

### 3.2 详细输入/输出规范

#### S0: 环境准备阶段

**输入**:
- `deploy/docker-compose-simple.yml`
- `backend/api/etc/dmh-api.yaml`
- `backend/scripts/init.sql`

**执行**:
```bash
cd deploy && docker compose -f docker-compose-simple.yml up -d
sleep 10
curl -sf http://localhost:8889/api/health || exit 1
```

**输出**:
- 服务状态: `docker ps` 输出
- 健康检查: API 响应 200

**验收标准**:
- [ ] dmh-api 运行中
- [ ] dmh-nginx 运行中
- [ ] mysql8 运行中
- [ ] redis-dmh 运行中
- [ ] API 健康检查返回 200

---

#### S1: Batch 1 - P0 核心链路

**输入**:
- `test-plan.md` Section 5.1~5.3
- `backend/test/integration/*.go`
- `frontend-h5/tests/e2e/*.spec.ts`

**执行**:
```bash
# 记录开始时间
START_TIME=$(date +%s)

# 后端集成测试 - 认证
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
  go test ./test/integration/... -v -run "TestAuth" -count=1 \
  | tee ../../docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch1-auth.log

# 后端集成测试 - 订单
./scripts/run_order_mysql8_regression.sh \
  | tee ../../docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch1-order.log

# H5 E2E - 活动相关
cd frontend-h5 && npm run test:e2e -- --grep "campaign|order" \
  | tee ../docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch1-h5-e2e.log
```

**输出**:
- `batch1-auth.log`: 认证测试日志
- `batch1-order.log`: 订单回归日志
- `batch1-h5-e2e.log`: H5 E2E 日志

**验收标准**:
- [ ] TC-RBAC-01~06 全部 PASS
- [ ] TC-CAMP-01~05 全部 PASS
- [ ] TC-ORDER-01~08 全部 PASS
- [ ] 无 Blocker/Critical 缺陷

---

#### S2: Batch 2 - P0 扩展

**输入**:
- `test-plan.md` Section 5.1 (扩展用例)
- `frontend-admin/tests/e2e/*.spec.ts`

**执行**:
```bash
# 后端全量集成
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
  go test ./test/integration/... -v -count=1 \
  | tee ../../docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch2-backend.log

# Admin E2E
cd frontend-admin && npm run test:e2e \
  | tee ../docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch2-admin-e2e.log
```

**输出**:
- `batch2-backend.log`: 后端全量日志
- `batch2-admin-e2e.log`: Admin E2E 日志

**验收标准**:
- [ ] 后端集成测试 27 suites 全部 PASS
- [ ] Admin E2E 21/21 PASS
- [ ] TC-RBAC-07~13 全部 PASS
- [ ] 无 Blocker/Critical 缺陷

---

#### S3: Batch 3 - P1 业务功能

**输入**:
- `test-plan.md` Section 5.4
- `backend/test/integration/feedback_test.go`

**执行**:
```bash
# 反馈系统测试
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
  go test ./test/integration/... -v -run "Feedback" -count=1 \
  | tee ../../docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch3-feedback.log

# H5 E2E 全量
cd frontend-h5 && npm run test:e2e \
  | tee ../docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch3-h5-full.log
```

**输出**:
- `batch3-feedback.log`: 反馈测试日志
- `batch3-h5-full.log`: H5 E2E 全量日志

**验收标准**:
- [ ] TC-FEED-01~13 通过率 ≥95%
- [ ] H5 E2E 7/7 PASS
- [ ] 无 Blocker 缺陷
- [ ] Critical 缺陷有风险豁免

---

#### S4: Batch 4 - 规格治理

**输入**:
- `openspec/specs/*/spec.md`
- `openspec/changes/*/proposal.md`

**执行**:
```bash
# OpenSpec 全量校验
openspec validate --all --strict --no-interactive \
  | tee docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch4-openspec.log

# 归档索引检查
cat openspec/changes/archive/ARCHIVE_STATUS_INDEX.md \
  | tee docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/batch4-archive-index.log
```

**输出**:
- `batch4-openspec.log`: OpenSpec 校验日志
- `batch4-archive-index.log`: 归档索引

**验收标准**:
- [ ] openspec validate 返回 "valid"
- [ ] 5 个规格全部校验通过
- [ ] 归档索引完整

---

#### S5: 汇总与报告

**输入**:
- 所有 batch 日志
- `docs/testing/execution/TEST_RESULT_TEMPLATE.md`

**执行**:
```bash
# 生成汇总报告
cat > docs/testing/execution/runs/$(date +%Y-%m-%d)-run001/verification-report.md << 'EOF'
# DMH 发布前全量回归验证报告

> 执行日期: $(date +%Y-%m-%d)
> 执行人: AI Agent
> 版本: $(git rev-parse --short HEAD)

## 1. 执行摘要

| 批次 | 状态 | 通过率 | 耗时 |
|------|------|--------|------|
| Batch 1 (P0 核心) | ✅/❌ | X% | Xmin |
| Batch 2 (P0 扩展) | ✅/❌ | X% | Xmin |
| Batch 3 (P1 业务) | ✅/❌ | X% | Xmin |
| Batch 4 (规格治理) | ✅/❌ | 100% | Xmin |

## 2. 质量门禁检查

- [ ] P0 用例 100% 通过
- [ ] P1 用例 ≥95% 通过
- [ ] 无 Blocker 缺陷
- [ ] 无未豁免 Critical 缺陷

## 3. 缺陷统计

| 严重级 | 数量 | 状态 |
|--------|------|------|
| Blocker | X | 全部关闭 |
| Critical | X | X 关闭 / Y 豁免 |
| Major | X | X 关闭 / Y 待修 |
| Minor | X | 记录 |

## 4. 发布决策

- **建议**: 可以发布 / 阻断发布
- **条件**: [如有阻断项，列出解除条件]

EOF
```

**输出**:
- `verification-report.md`: 完整验证报告

**验收标准**:
- [ ] 报告包含所有 batch 结果
- [ ] 质量门禁检查项已填写
- [ ] 发布决策明确

---

## 4. 失败回滚策略、风险分级与发布阻断条件

### 4.1 失败分级与响应

| 失败级别 | 定义 | 响应动作 | SLA |
|----------|------|----------|-----|
| **Level 1 (阻断)** | Batch 1/2 任何失败 | 停止执行，修复后重跑 | 立即 |
| **Level 2 (严重)** | Batch 3 通过率 <95% | 评估影响，决定继续或阻断 | 30min 内决策 |
| **Level 3 (一般)** | Minor 缺陷或 P1 失败 | 记录缺陷，继续执行 | 记录即可 |
| **Level 4 (警告)** | 非关键功能异常 | 记录，不阻断发布 | 记录即可 |

### 4.2 回滚策略

```
┌─────────────────────────────────────────────────────────────┐
│                    失败回滚决策树                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  失败发生                                                    │
│      │                                                       │
│      ▼                                                       │
│  ┌─────────────────┐                                         │
│  │ 是否 Level 1?   │──YES──► 立即停止                        │
│  └────────┬────────┘         │                               │
│           NO                 ▼                               │
│           │          ┌─────────────────┐                     │
│           │          │ 记录失败详情    │                     │
│           │          │ 通知相关人员    │                     │
│           ▼          │ 准备修复方案    │                     │
│  ┌─────────────────┐ └────────┬────────┘                     │
│  │ 是否 Level 2?   │          │                               │
│  └────────┬────────┘          ▼                               │
│           YES         ┌─────────────────┐                     │
│           │           │ 评估影响范围    │                     │
│           │           │ 决定继续/阻断   │                     │
│           │           │ 如阻断: 申请豁免│                     │
│           ▼           └─────────────────┘                     │
│  ┌─────────────────┐                                          │
│  │ 记录并继续      │◄── Level 3/4                            │
│  └─────────────────┘                                          │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 4.3 回滚操作清单

**代码回滚**:
```bash
# 回滚到上一版本
git revert HEAD
git push origin main

# 或回滚到指定版本
git reset --hard <commit-hash>
git push origin main --force  # 谨慎使用
```

**数据库回滚**:
```bash
# 恢复数据库快照（需要预先备份）
mysql -u root -p dmh < backup/pre-release-backup.sql
```

**服务回滚**:
```bash
# 重新部署上一版本
cd deploy
docker compose -f docker-compose-simple.yml down
# 替换 dmh-api 二进制为上一版本
docker compose -f docker-compose-simple.yml up -d
```

### 4.4 风险分级

| 风险等级 | 影响范围 | 示例 | 发布条件 |
|----------|----------|------|----------|
| 🔴 **Critical** | 核心业务中断 | 认证失败、订单无法创建、支付失败 | **必须修复** |
| 🟠 **High** | 主要功能受损 | 活动编辑失败、核销异常、数据泄露 | **必须修复或豁免** |
| 🟡 **Medium** | 部分功能异常 | 反馈提交失败、统计不准确 | 修复或有计划 |
| 🟢 **Low** | 体验问题 | UI 错位、提示文案错误 | 可发布后修复 |

### 4.5 发布阻断条件

**自动阻断** (任何一条触发):
1. ❌ Batch 1 (P0 核心) 任何测试失败
2. ❌ Batch 2 (P0 扩展) 任何测试失败
3. ❌ 存在未关闭的 Blocker 缺陷
4. ❌ OpenSpec 校验失败

**人工阻断** (需评估决策):
1. ⚠️ Batch 3 (P1) 通过率 <95%
2. ⚠️ 存在未关闭的 Critical 缺陷
3. ⚠️ 性能指标不达标 (订单创建 <100 QPS)

**可发布条件**:
- ✅ Batch 1~2 全部通过
- ✅ Batch 3 通过率 ≥95%
- ✅ Batch 4 全部通过
- ✅ 无 Blocker 缺陷
- ✅ Critical 缺陷全部关闭或有有效豁免

### 4.6 风险豁免流程

```
┌─────────────────────────────────────────────────────────────┐
│                    风险豁免流程                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. 发现不可修复/暂缓修复的缺陷                              │
│      │                                                       │
│      ▼                                                       │
│  2. 填写风险豁免申请                                         │
│     - 缺陷描述与影响评估                                     │
│     - 补救计划与时间表                                       │
│     - 风险接受理由                                           │
│      │                                                       │
│      ▼                                                       │
│  3. 审批链                                                   │
│     - 技术负责人审批 (Critical)                              │
│     - 产品负责人审批 (业务影响)                              │
│      │                                                       │
│      ▼                                                       │
│  4. 豁免生效                                                 │
│     - 记录到 verification-report.md                          │
│     - 设置有效期 (如: 下个版本必须修复)                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**豁免模板**:
```markdown
## 风险豁免申请

- **缺陷ID**: BUG-XXX
- **严重级**: Critical
- **影响范围**: [描述]
- **风险评估**: [低/中/高]
- **补救计划**: [具体措施和时间]
- **申请人**: [姓名]
- **审批人**: [姓名]
- **有效期**: [日期]
- **状态**: [待审批/已批准/已拒绝]
```

---

## 5. 执行时间表

| 时间 | 阶段 | 内容 | 负责人 |
|------|------|------|--------|
| T+0 | S0 | 环境准备与启动 | AI Agent |
| T+5min | S1 | Batch 1: P0 核心 | AI Agent |
| T+35min | S2 | Batch 2: P0 扩展 | AI Agent |
| T+55min | S3 | Batch 3: P1 业务 | AI Agent |
| T+80min | S4 | Batch 4: 规格治理 | AI Agent |
| T+85min | S5 | 汇总与报告 | AI Agent |
| T+90min | - | 发布决策 | 人工审核 |

**预计总耗时**: 90 分钟

---

## 6. 附录

### A. 测试命令速查

```bash
# 后端全量
cd backend && go test ./... -v -cover

# 后端集成
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
  go test ./test/integration/... -v -count=1

# 订单回归
cd backend && ./scripts/run_order_mysql8_regression.sh

# Admin E2E
cd frontend-admin && npm run test:e2e

# H5 E2E
cd frontend-h5 && npm run test:e2e

# OpenSpec 校验
openspec validate --all --strict --no-interactive
```

### B. 联系人

| 角色 | 负责范围 | 联系方式 |
|------|----------|----------|
| 技术负责人 | Critical 缺陷决策 | - |
| 产品负责人 | 业务影响评估 | - |
| QA 负责人 | 测试执行与报告 | - |

### C. 相关文档

- `test-plan.md`: 测试计划详细用例
- `all-requirements.md`: 需求基线
- `docs/testing/execution/`: 执行矩阵与模板
- `openspec/AGENTS.md`: OpenSpec 操作指南

---

*计划版本: v1.0 | 创建: 2026-02-15*
