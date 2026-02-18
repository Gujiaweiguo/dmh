# DMH 全量回归基线清单 (Plan v1 - FRP-001)

> 生成日期: 2026-02-18
> 变更: add-full-regression-testing
> 任务: FRP-001 基线盘点冻结

---

## 1. 必跑测试矩阵

### 1.1 后端 (Backend)

| 层级 | 命令入口 | 触发条件 | 覆盖范围 | 状态 |
|------|----------|----------|----------|------|
| **单元测试** | `cd backend && go test ./...` | PR/push to main | 全部 Go 包 | ✅ 已覆盖 |
| **集成测试** | `cd backend && go test ./test/integration/... -v -count=1` | PR/push to main | 需要运行中的 API + MySQL + Redis | ✅ 已覆盖 |
| **订单专项回归** | `backend/scripts/run_order_mysql8_regression.sh` | order 路径变更 PR/push | `TestOrderVerifyRoutesAuthGuard`, `TestOrderCreateDuplicateMessage` | ✅ 已覆盖 |
| **覆盖率门禁** | `go test ./... -coverprofile=coverage.out` | PR to main | 阈值 >= 76% | ✅ 已覆盖 |

### 1.2 前端管理后台 (Frontend-Admin)

| 层级 | 命令入口 | 触发条件 | 覆盖范围 | 状态 |
|------|----------|----------|----------|------|
| **单元测试** | `cd frontend-admin && npm run test` | PR/push to main | `tests/unit/**/*.test.ts` | ✅ 已覆盖 |
| **E2E 测试** | `cd frontend-admin && npm run test:e2e:headless` | PR/push to main | `e2e/**/*.spec.ts` | ✅ 已覆盖 |
| **覆盖率门禁** | `npm run test:cov` | PR to main | 阈值 >= 70% | ✅ 已覆盖 |

### 1.3 前端 H5 (Frontend-H5)

| 层级 | 命令入口 | 触发条件 | 覆盖范围 | 状态 |
|------|----------|----------|----------|------|
| **单元测试** | `cd frontend-h5 && npm run test` | PR/push to main (system-test-gate) | `tests/unit/**/*.test.js` | ⚠️ 仅在 system-test-gate 覆盖 |
| **E2E 测试** | `cd frontend-h5 && npm run test:e2e:headless` | PR/push to main (system-test-gate) | `e2e/**/*.spec.ts` | ⚠️ 仅在 system-test-gate 覆盖 |
| **覆盖率门禁** | `npm run test:cov` | PR to main | 阈值 >= 44% | ✅ 已覆盖 |

### 1.4 OpenSpec 校验

| 层级 | 命令入口 | 触发条件 | 覆盖范围 | 状态 |
|------|----------|----------|----------|------|
| **全量校验** | `openspec validate --all --no-interactive` | PR/push to main (system-test-gate) | 所有 specs 和 changes | ✅ 已覆盖 |

---

## 2. CI Workflow 覆盖矩阵

| Workflow | 后端单元 | 后端集成 | 订单回归 | Admin 单元 | Admin E2E | H5 单元 | H5 E2E | OpenSpec | 触发条件 |
|----------|:--------:|:--------:|:--------:|:----------:|:---------:|:-------:|:------:|:--------:|----------|
| `stability-checks.yml` | ✅ | ✅ | ❌ | ✅ | ✅ (仅 security) | ❌ | ❌ | ❌ | PR/push: backend/**, frontend-admin/** |
| `system-test-gate.yml` | ❌ | ✅ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ | PR/push: backend/**, frontend-admin/**, frontend-h5/**, openspec/** |
| `coverage-gate.yml` | ✅ (覆盖率) | ❌ | ❌ | ✅ (覆盖率) | ❌ | ✅ (覆盖率) | ❌ | ❌ | PR to main |
| `order-mysql8-regression.yml` | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | PR/push: order 相关路径 |
| `feedback-guard.yml` | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | PR/push: feedback 相关路径 |

---

## 3. 缺口分析 (Gaps)

### 3.1 覆盖缺口

| 缺口 | 描述 | 影响 | 建议 |
|------|------|------|------|
| **H5 单元测试不在 stability-checks** | `stability-checks.yml` 未覆盖 `frontend-h5` 单元测试 | H5 变更可能未触发单元测试 | 扩展 stability-checks 路径或合并到 system-test-gate |
| **无全量回归聚合结论** | 各 workflow 独立执行，无统一 PASS/FAIL 结论 | 发布决策需人工汇总 | 引入聚合 job 或独立 workflow |
| **无 flaky 重试策略** | 测试失败即终止，无受控重试 | 环境瞬态问题导致不必要的失败 | 引入 flaky 白名单与重试上限 |
| **无证据保留策略** | artifact 保留依赖 GitHub 默认策略 | 历史回归证据可能丢失 | 定义分级保留策略 |
| **无夜间回归** | 仅 PR/push 触发，无定时回归 | 长时间无变更时回归覆盖不足 | 引入 nightly workflow |
| **Admin E2E 覆盖有限** | stability-checks 仅跑 security-management.spec.ts | 其他 E2E 场景仅在 system-test-gate 覆盖 | 考虑扩展或统一 |
| **性能测试不在 CI** | `test/performance/` 有测试但无 CI job | 性能回归可能被遗漏 | 引入 dedicated performance job |

### 3.2 触发条件缺口

| Workflow | 路径过滤 | 缺失触发 |
|----------|----------|----------|
| `stability-checks.yml` | `backend/**`, `frontend-admin/**` | `frontend-h5/**`, `openspec/**` |
| `coverage-gate.yml` | 无路径过滤（仅 PR to main） | - |
| `order-mysql8-regression.yml` | order 特定路径 | 广泛变更时可能不触发 |

---

## 4. 本地执行入口现状

### 4.1 Makefile 命令

```makefile
make test              # 后端单元 + Admin 单元 + H5 单元 (无 E2E)
make test-backend      # 后端单元
make test-integration  # 后端集成 (需要运行中的服务)
make test-admin        # Admin 单元
make test-h5           # H5 单元
make test-e2e          # Admin E2E + H5 E2E (需 headed 浏览器)
```

### 4.2 独立脚本

- `backend/scripts/run_tests.sh` - RBAC 全套 (单元+集成+性能+基准)
- `backend/scripts/run_order_mysql8_regression.sh` - 订单专项
- `backend/scripts/run_order_logic_tests.sh` - 订单逻辑单元+冒烟
- `backend/scripts/repair_login_and_run_order_regression.sh` - 自动修复+回归

### 4.3 缺口

- **无本地一键全量回归入口** - 需要多个命令串联
- **无与 CI 一致的必跑矩阵** - 本地执行范围与 CI 不完全一致
- **无证据模板** - 结果零散，无统一记录

---

## 5. 现有测试资产统计

| 模块 | 单元测试 | 集成测试 | E2E 测试 | 性能测试 |
|------|:--------:|:--------:|:--------:|:--------:|
| Backend | 72+ 文件 | 22 文件 | - | 3 文件 |
| Frontend-Admin | 30+ 文件 | - | 4 spec | - |
| Frontend-H5 | 40+ 文件 | - | 2 spec | - |
| OpenSpec | - | - | - | - |

---

## 6. 基线冻结确认

- [x] 后端四层（单元/集成/订单回归/覆盖率）已覆盖
- [x] Admin 三层（单元/E2E/覆盖率）已覆盖
- [x] H5 三层（单元/E2E/覆盖率）已覆盖
- [x] OpenSpec 校验已覆盖
- [x] 缺口已识别并记录
- [x] 本地入口现状已盘点

---

## 7. 下一步 (FRP-002+)

1. **FRP-002**: 定义统一口径（全量回归 = 必跑矩阵 + 单一结论）
2. **FRP-003**: 发布阻断规则落地
3. **FRP-004**: flaky 策略落地
4. **FRP-005**: 证据保留窗口定义
5. **FRP-006**: CI 编排对齐实现
6. **FRP-007**: 本地一键入口对齐
7. **FRP-008**: 证据模板与审计落地
8. **FRP-009**: 全量回归演练
9. **FRP-010**: 发布前签收

---

*此文档由 FRP-001 任务生成，作为后续任务的输入基线。*
