# T2: CI 分层目标与时延预算定义

**生成时间**: 2026-02-26 20:35 CST

---

## 1. 分层策略概览

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CI 测试分层架构                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌──────────────────┐                                              │
│   │   PR Gate        │  ← 开发者提交 PR 时触发                       │
│   │   目标: ≤15min   │  ← 快速反馈，阻断低质量代码                   │
│   └────────┬─────────┘                                              │
│            │                                                         │
│   ┌────────▼─────────┐                                              │
│   │   Merge Gate     │  ← 合并到 main/master 时触发                 │
│   │   目标: ≤20min   │  ← 确保主干稳定                              │
│   └────────┬─────────┘                                              │
│            │                                                         │
│   ┌────────▼─────────┐                                              │
│   │   Nightly        │  ← 每日定时 + 发布候选标签触发               │
│   │   目标: ≤45min   │  ← 全量回归，深度验证                        │
│   └──────────────────┘                                              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. 三层测试责任矩阵

| 测试类型 | PR Gate | Merge Gate | Nightly |
|----------|:-------:|:----------:|:-------:|
| **单元测试** | ✅ 必须 | ✅ 必须 | ✅ 必须 |
| **覆盖率门禁** | ✅ 必须 | ✅ 必须 | ✅ 必须 |
| **Lint 检查** | ✅ 必须 | ✅ 必须 | ✅ 必须 |
| **集成测试** | ❌ 跳过 | ⚠️ 关键用例 | ✅ 全量 |
| **E2E 测试** | ❌ 跳过 | ⚠️ Smoke | ✅ 全量 |
| **性能测试** | ⚠️ Smoke | ❌ 跳过 | ✅ 全量 |
| **安全扫描** | ❌ 跳过 | ❌ 跳过 | ✅ 必须 |

---

## 3. 时延预算分配

### 3.1 PR Gate (目标: ≤15min)

| Job | 预算 | 说明 |
|-----|------|------|
| backend-unit | 8min | 单元测试 + 覆盖率 |
| frontend-unit | 5min | Admin + H5 单测 |
| lint | 2min | Go + 前端 lint |
| **总计** | **15min** | 并行执行 |

### 3.2 Merge Gate (目标: ≤20min)

| Job | 预算 | 说明 |
|-----|------|------|
| PR Gate 全部 | 15min | 继承 PR Gate |
| integration-smoke | 3min | 关键集成用例 |
| e2e-smoke | 2min | 登录/支付关键路径 |
| **总计** | **20min** | 串行叠加 |

### 3.3 Nightly (目标: ≤45min)

| Job | 预算 | 说明 |
|-----|------|------|
| coverage-gate | 10min | 全量覆盖率 |
| system-test-gate | 15min | 系统测试 |
| stability-checks | 10min | 稳定性检查 |
| order-mysql8-regression | 5min | 订单回归 |
| e2e-full | 5min | 全量 E2E |
| **总计** | **45min** | 并行 + 串行混合 |

---

## 4. 测试命令分层映射

### 4.1 PR Gate 命令

```bash
# 后端 (排除集成/性能)
go test -p 1 $(go list ./... | grep -v -E 'dmh/test/integration|dmh/test/performance') -coverprofile=coverage.out

# Admin
cd frontend-admin && npm run test:cov

# H5
cd frontend-h5 && npm run test:cov
```

### 4.2 Merge Gate 命令

```bash
# 继承 PR Gate + 集成 Smoke
go test -p 1 ./test/integration/... -run "Smoke|Critical" -short

# E2E Smoke
cd frontend-admin && npm run test:e2e -- --grep "Smoke"
```

### 4.3 Nightly 命令

```bash
# 全量后端
go test -p 1 ./... -coverprofile=coverage.out

# 集成测试
DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v

# 性能测试
go test ./test/performance/... -bench=. -benchtime=5s

# E2E 全量
cd frontend-admin && npm run test:e2e
cd frontend-h5 && npm run test:e2e
```

---

## 5. 覆盖率阈值

| 模块 | 阈值 | 当前 | 状态 |
|------|------|------|------|
| Backend | 78% | ~78% | ✅ |
| Admin | 80% | 83.65% | ✅ |
| H5 | 70% | ~99% | ✅ |

---

## 6. 失败处理策略

| 场景 | PR Gate | Merge Gate | Nightly |
|------|---------|------------|---------|
| 单元测试失败 | ❌ 阻断 | ❌ 阻断 | 📧 告警 |
| 覆盖率不足 | ❌ 阻断 | ❌ 阻断 | 📧 告警 |
| 集成测试失败 | - | ⚠️ 警告 | 📧 告警 |
| E2E 失败 | - | ⚠️ 警告 | 📧 告警 |
| 环境问题 | ⏭️ Skip+Tag | ⏭️ Skip+Tag | 📧 告警+重试 |

---

## 7. 当前状态评估

### 7.1 现有配置分析

| Workflow | 超时设置 | 实际情况 | 差距 |
|----------|----------|----------|------|
| pr-gate.yml | 15min | ✅ 合理 | 无 |
| coverage-gate.yml | 未设置 | ⚠️ 需添加 | 添加 timeout |
| full-regression.yml | 未设置 | ⚠️ 需添加 | 添加 timeout |

### 7.2 改进建议

1. **后端测试优化** (优先级: 🔴 高)
   - 当前 `-short` 超时 >3min，需优化到 <2min
   - 排除慢测试，移至 Nightly

2. **添加 Merge Gate** (优先级: 🟡 中)
   - 当前无独立 Merge Gate
   - 建议在 main/master 合并时触发关键集成测试

3. **统一超时配置** (优先级: 🟡 中)
   - 为所有 workflow 添加显式 timeout-minutes

---

## 8. 验收标准

- [x] PR/Merge/Nightly 三层职责定义完成
- [x] 时延预算表落地
- [x] 测试命令分层映射完成
- [x] 失败处理策略定义完成
- [ ] CI 配置实际改造 (T9/T10)
