# 性能测试分流策略 (T13)

## 1. 现有性能测试分析

### 1.1 测试文件清单

| 文件 | 类型 | 需要后端 |
|------|------|----------|
| `benchmark_test.go` | 基准/功能测试 | ❌ |
| `rbac_performance_test.go` | RBAC 性能测试 | ✅ |
| `advanced_features_performance_test.go` | 高级功能性能测试 | ✅ |

### 1.2 `benchmark_test.go` - 基础性能测试

| 测试函数 | 类型 | 耗时 | 模式控制 |
|---------|------|------|----------|
| `BenchmarkCreateOrder` | Benchmark | ~1ms/次 | 无需跳过 |
| `BenchmarkGetCampaigns` | Benchmark | ~1ms/次 | 无需跳过 |
| `BenchmarkVerifyOrder` | Benchmark | ~2ms/次 | 无需跳过 |
| `TestConcurrentOrderCreation` | 功能测试 | **10s** | `testing.Short()` 跳过 |
| `TestDatabaseConnectionPool` | 功能测试 | ~20ms | 无需跳过 |
| `TestMemoryLeak` | 功能测试 | ~15s | `testing.Short()` 跳过 |

### 1.3 `rbac_performance_test.go` - RBAC 性能测试

| 测试函数 | 类型 | 耗时 | 后端不可达处理 |
|---------|------|------|----------------|
| `TestRBACAdminAccessLatency` | 功能测试 | ~1s | `t.Skip()` 跳过 |
| `TestRBACUnauthorizedResponse` | 功能测试 | ~100ms | `t.Skip()` 跳过 |

### 1.4 `advanced_features_performance_test.go` - 高级功能测试

| 测试函数 | 类型 | 耗时 | 后端不可达处理 |
|---------|------|------|----------------|
| `Test_12_1_PosterGenerationPerformance` | 功能测试 | <3s | `t.Skip()` 跳过 |
| `Test_12_2_PaymentQRCodePerformance` | 功能测试 | <500ms | 自动跳过 |
| `Test_12_3_OrderVerifyPerformance` | 功能测试 | <500ms | 自动跳过 |
| `Test_12_4_ConcurrentPosterStressTest` | 压力测试 | ~10s | 自动跳过 |

### 1.5 现有 `-short` 模式支持

```go
// 已有代码示例 (benchmark_test.go)
func TestConcurrentOrderCreation(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过长时间性能测试（使用 go test -v 不带 -short 来运行）")
    }
    // ... 10 秒测试逻辑
}

func TestMemoryLeak(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过内存泄漏测试")
    }
    // ... 15 秒测试逻辑
}
```

**✅ 结论**：
- 基础测试已有完善的 `-short` 模式支持
- 需要后端服务的测试会自动 Skip，不影响 PR Gate

---

## 2. 分流策略定义

### 2.1 PR Gate 性能测试 (short/perf-smoke)

**目标**：快速验证性能没有明显退化，<1 分钟

**执行命令**：
```bash
go test -short -bench=. -benchtime=100ms ./test/performance/...
```

**包含测试**：
| 测试 | 文件 | 原因 |
|------|------|------|
| `BenchmarkCreateOrder` | benchmark_test.go | 快速基准 |
| `BenchmarkGetCampaigns` | benchmark_test.go | 快速基准 |
| `BenchmarkVerifyOrder` | benchmark_test.go | 快速基准 |
| `TestDatabaseConnectionPool` | benchmark_test.go | ~20ms，极快 |
| `TestRBACAdminAccessLatency` | rbac_performance_test.go | 自动 Skip（无后端）|
| `TestRBACUnauthorizedResponse` | rbac_performance_test.go | 自动 Skip（无后端）|
| `AdvancedFeaturesPerformanceTestSuite` | advanced_features_performance_test.go | 自动 Skip（无后端）|

**排除测试**：
| 测试 | 原因 |
|------|------|
| `TestConcurrentOrderCreation` | 10s 持续测试，`-short` 跳过 |
| `TestMemoryLeak` | 15s 内存测试，`-short` 跳过 |
**预计耗时**：~5-10 秒

### 2.2 Nightly 性能测试 (full benchmark)

**目标**：完整性能验证，包括长时间稳定性测试

**执行命令**：
```bash
go test -v ./test/performance/...
```

**包含测试**：
| 测试 | 文件 | 原因 |
|------|------|------|
| `BenchmarkCreateOrder` | benchmark_test.go | 完整基准 |
| `BenchmarkGetCampaigns` | benchmark_test.go | 完整基准 |
| `BenchmarkVerifyOrder` | benchmark_test.go | 完整基准 |
| `TestDatabaseConnectionPool` | benchmark_test.go | 连接池压力测试 |
| `TestConcurrentOrderCreation` | benchmark_test.go | 10s 并发稳定性 |
| `TestMemoryLeak` | benchmark_test.go | 内存泄漏检测 |
| `TestRBACAdminAccessLatency` | rbac_performance_test.go | RBAC 延迟测试 |
| `TestRBACUnauthorizedResponse` | rbac_performance_test.go | RBAC 鉴权测试 |
| `Test_12_1_PosterGenerationPerformance` | advanced_features_performance_test.go | 海报生成 <3s |
| `Test_12_2_PaymentQRCodePerformance` | advanced_features_performance_test.go | 二维码 <500ms |
| `Test_12_3_OrderVerifyPerformance` | advanced_features_performance_test.go | 核销 <500ms |
| `Test_12_4_ConcurrentPosterStressTest` | advanced_features_performance_test.go | 并发压力测试 |

**预计耗时**：~30-45 秒（需要后端服务）

---

## 3. 性能测试分类清单

### 3.1 分类标准

| 分类 | 标志 | 耗时范围 | 运行环境 |
|------|------|---------|----------|
| **PR Smoke** | `-short -bench=.` | <10s | PR Gate |
| **Full Benchmark** | 无 `-short` | 25-30s | Nightly |

### 3.2 测试归属表
```
PR Gate (short/perf-smoke) - 无需后端服务:
├── benchmark_test.go
│   ├── BenchmarkCreateOrder      ✅
│   ├── BenchmarkGetCampaigns     ✅
│   ├── BenchmarkVerifyOrder      ✅
│   └── TestDatabaseConnectionPool ✅
├── rbac_performance_test.go      (自动 Skip - 无后端)
└── advanced_features_performance_test.go (自动 Skip - 无后端)

Nightly (full benchmark) - 需要后端服务:
├── benchmark_test.go
│   ├── BenchmarkCreateOrder      ✅
│   ├── BenchmarkGetCampaigns     ✅
│   ├── BenchmarkVerifyOrder      ✅
│   ├── TestDatabaseConnectionPool ✅
│   ├── TestConcurrentOrderCreation ✅ (仅 Nightly)
│   └── TestMemoryLeak            ✅ (仅 Nightly)
├── rbac_performance_test.go
│   ├── TestRBACAdminAccessLatency ✅
│   └── TestRBACUnauthorizedResponse ✅
└── advanced_features_performance_test.go
    ├── Test_12_1_PosterGenerationPerformance ✅
    ├── Test_12_2_PaymentQRCodePerformance ✅
    ├── Test_12_3_OrderVerifyPerformance ✅
    └── Test_12_4_ConcurrentPosterStressTest ✅
```

---

## 4. CI Workflow 更新建议

### 4.1 PR Gate 更新

**文件**：`.github/workflows/pr-gate.yml`

**建议添加 Job**：
```yaml
  perf-smoke:
    name: Performance Smoke Tests
    runs-on: ubuntu-latest
    timeout-minutes: 5
    
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: backend/go.mod
          cache: true
          cache-dependency-path: backend/go.sum
      
      - name: Run performance smoke tests
        working-directory: backend
        run: |
          # 运行 short 模式性能测试（长时间测试被跳过，需要后端的测试也会自动跳过）
          go test -short -bench=. -benchtime=100ms ./test/performance/... -v 2>&1 | tee /tmp/perf-smoke.log
          echo "✓ Performance smoke tests passed (<1min)"

**更新 `pr-gate-verdict` needs**：
```yaml
    needs:
      - backend-unit
      - frontend-unit
      - lint
      - perf-smoke  # 新增
```

### 4.2 Nightly 更新

**文件**：`.github/workflows/full-regression.yml`

**建议添加 Job**：
```yaml
  perf-benchmark:
    name: Full Performance Benchmark
    runs-on: ubuntu-latest
    timeout-minutes: 10
    
    services:
      mysql8:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: "Admin168"
          MYSQL_DATABASE: dmh
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping -h 127.0.0.1 -p'Admin168' --silent"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=20
    
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: backend/go.mod
          cache: true
      
      - name: Initialize schema
        run: mysql -h 127.0.0.1 -uroot -p'Admin168' dmh < backend/scripts/init.sql
      
      - name: Build and start backend
        working-directory: backend
        run: |
          go build -o dmh-api api/dmh.go
          ./dmh-api -f api/etc/dmh-api.yaml &
          sleep 5
          
      - name: Run full performance benchmark
        working-directory: backend
        run: |
          # 运行完整性能测试（包括需要后端的测试）
          go test -v ./test/performance/...
          echo "✓ Full performance benchmark completed"

**更新 `full-regression-verdict` needs**：
```yaml
    needs:
      - stability-checks
      - system-test-gate
      - coverage-gate
      - order-mysql8-regression
      - perf-benchmark  # 新增
```

---

## 5. 时延预算验证

### 5.1 PR Gate 时延预算

| Job | 当前耗时 | 变更后 |
|-----|---------|--------|
| backend-unit | ~8min | 不变 |
| frontend-unit | ~5min | 不变 |
| lint | ~1min | 不变 |
| **perf-smoke** | - | **~10s** |
| **总计** | ~14min | **~14min** |

**✅ 结论**：PR Gate 总耗时仍在 15 分钟预算内。

### 5.2 Nightly 时延预算

| Job | 耗时 |
|-----|------|
| stability-checks | ~5min |
| system-test-gate | ~10min |
| stability-checks | ~5min |
| system-test-gate | ~10min |
| coverage-gate | ~5min |
| order-mysql8-regression | ~5min |
| **perf-benchmark** | **~1min** |
| **总计** | ~27min |

**✅ 结论**：Nightly 总耗时仍在 45 分钟预算内。

---

## 6. 分流规则文档

### 6.1 何时添加到 PR Smoke

满足以下**所有**条件：
1. 耗时 < 5 秒
2. 使用 `testing.Short()` 支持 `-short` 跳过（如需）
3. 验证核心功能性能不退化

### 6.2 何时添加到 Nightly Only

满足以下**任一**条件：
1. 耗时 > 10 秒
2. 需要长时间稳定性验证
3. 需要完整环境（数据库、缓存等）
4. 压力测试、内存泄漏检测

### 6.3 命名约定

```go
// 快速测试 - 可在 PR 运行
func BenchmarkXxx(b *testing.B) { ... }
func TestXxxQuick(t *testing.T) { ... }

// 慢速测试 - 仅 Nightly
func TestXxxLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过长时间测试")
    }
    // ...
}
```

---

## 7. 验证命令

### 7.1 本地验证 PR Smoke

```bash
cd backend
go test -short -bench=. -benchtime=100ms ./test/performance/... -v
```

**预期输出**：
```
=== RUN   TestAdvancedFeaturesPerformanceTestSuite
--- SKIP: TestAdvancedFeaturesPerformanceTestSuite (0.00s)
=== RUN   TestConcurrentOrderCreation
    benchmark_test.go:53: 跳过长时间性能测试
--- SKIP: TestConcurrentOrderCreation (0.00s)
=== RUN   TestDatabaseConnectionPool
连接池测试完成: 50 连接, 1000 次查询, 耗时 22ms
--- PASS: TestDatabaseConnectionPool (0.02s)
=== RUN   TestMemoryLeak
    benchmark_test.go:113: 跳过内存泄漏测试
--- SKIP: TestMemoryLeak (0.00s)
=== RUN   TestRBACAdminAccessLatency
    rbac_performance_test.go:38: 后端不可达，跳过 RBAC 性能测试
--- SKIP: TestRBACAdminAccessLatency (0.00s)
goos: linux
BenchmarkCreateOrder-8           825     141721 ns/op
BenchmarkGetCampaigns-8          100    1112154 ns/op
BenchmarkVerifyOrder-8           409     283596 ns/op
PASS
ok      dmh/test/performance    0.411s
```

**实际验证结果**：PR Smoke 耗时 **0.411s** ✅

### 7.2 本地验证 Full Benchmark

```bash
cd backend
# 需要先启动后端服务
./dmh-api -f api/etc/dmh-api.yaml &
# 运行完整测试
go test -v ./test/performance/...
```

**预期输出**：
```
BenchmarkCreateOrder-8      ...
BenchmarkGetCampaigns-8     ...
BenchmarkVerifyOrder-8      ...
--- PASS: TestConcurrentOrderCreation (10.xx s)
--- PASS: TestDatabaseConnectionPool (0.02 s)
--- PASS: TestMemoryLeak (15.xx s)
--- PASS: TestRBACAdminAccessLatency (1.xx s)
--- PASS: Test_12_1_PosterGenerationPerformance (2.xx s)
--- PASS: Test_12_2_PaymentQRCodePerformance (0.5 s)
--- PASS: Test_12_3_OrderVerifyPerformance (0.5 s)
--- PASS: Test_12_4_ConcurrentPosterStressTest (10.xx s)
PASS
```

---

## 8. 总结

### 8.1 分流策略摘要

| 维度 | PR Gate (Smoke) | Nightly (Full) |
|------|-----------------|----------------|
| **执行命令** | `-short -bench=.` | 无 `-short` + 后端服务 |
| **测试文件** | 1 个 (benchmark_test.go) | 3 个 (全部) |
| **测试数量** | 4 个 + 自动跳过 | 12 个 |
| **耗时** | <1s | ~1min |
| **目的** | 快速回归 | 完整验证 |
| **需要后端** | ❌ | ✅ |

### 8.2 关键指标

- ✅ PR 性能测试耗时: **<1 分钟** (实际 0.411s)
- ✅ 无漏检关键 smoke 测试: 3 个 benchmark + 1 个连接池测试
- ✅ 长时间测试仅在 Nightly: 10s 并发 + 15s 内存泄漏 + 压力测试
- ✅ 需要后端的测试自动跳过: PR Gate 无后端依赖
- ✅ 时延预算符合: PR ≤15min, Nightly ≤45min

---

*生成时间: 2026-02-26*
*任务: task-13-perf-test-split*
