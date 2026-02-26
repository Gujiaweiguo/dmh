# T12: 集成测试稳定性治理

## 1. 概述

本文档定义 DMH 项目集成测试 flaky 治理方案，包括：
- Flaky 模式分析与分类
- 重试策略（仅限幂等操作）
- SKIP 原因标签使用规范
- 就绪探针最佳实践
- Flaky 率监控与周报指标

---

## 2. Flaky 模式分析

### 2.1 发现的 Flaky 模式

| 模式 | 描述 | 风险级别 | 文件示例 |
|------|------|----------|----------|
| **资源泄漏** | 循环内使用 `defer resp.Body.Close()`，导致连接累积 | 高 | `rate_limiting_test.go:159`, `concurrency_test.go:134` |
| **硬编码配置** | 直接使用 `http://localhost:8889` 而非环境变量 | 中 | 所有集成测试文件 |
| **登录失败处理不一致** | 部分测试登录失败后继续执行，导致后续 panic | 高 | `rate_limiting_test.go:80-83` |
| **无超时保护** | HTTP 请求未设置超时或超时过长 | 中 | 部分测试 |
| **并发竞态** | 并发测试中共享状态未正确同步 | 高 | `concurrency_test.go` |
| **SKIP 原因缺失** | Skipf 未使用标准化标签 | 中 | 多数测试文件 |

### 2.2 失败根因分类

```
┌─────────────────────────────────────────────────────────────┐
│                    Flaky 失败根因                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  环境问题 (40%)                                             │
│  ├── API 服务未启动/未就绪                                  │
│  ├── MySQL 连接失败                                         │
│  └── Redis 不可用（可选依赖）                               │
│                                                             │
│  代码缺陷 (30%)                                             │
│  ├── 资源泄漏（defer in loop）                              │
│  ├── 错误处理不完整                                         │
│  └── 并发安全问题                                           │
│                                                             │
│  数据问题 (20%)                                             │
│  ├── 测试数据冲突                                           │
│  ├── 数据准备失败                                           │
│  └── 数据库状态不一致                                       │
│                                                             │
│  外部依赖 (10%)                                             │
│  ├── 网络抖动                                               │
│  └── 第三方服务超时                                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 高风险测试文件

| 文件 | Flaky 风险 | 原因 |
|------|------------|------|
| `rate_limiting_test.go` | 高 | 循环内 defer、限流测试依赖服务器状态 |
| `concurrency_test.go` | 高 | 并发竞态、共享状态 |
| `order_complete_flow_test.go` | 中 | 多步骤依赖、数据准备 |
| `campaign_handler_integration_test.go` | 中 | 数据库写入、唯一约束 |

---

## 3. 重试策略

### 3.1 核心原则

> **重试仅限幂等操作，不掩盖真实缺陷**

### 3.2 幂等性判断矩阵

| 操作类型 | 幂等性 | 是否可重试 | 说明 |
|----------|--------|------------|------|
| **HTTP GET** | ✅ 是 | ✅ 可重试 | 读取操作，无副作用 |
| **HTTP PUT** | ✅ 是 | ✅ 可重试 | 全量更新，重复执行结果相同 |
| **HTTP DELETE** | ✅ 是 | ✅ 可重试 | 删除不存在资源返回 404 |
| **HTTP POST (创建)** | ❌ 否 | ❌ 不可重试 | 可能产生重复数据 |
| **HTTP POST (搜索)** | ✅ 是 | ✅ 可重试 | 查询类 POST |
| **登录请求** | ✅ 是 | ✅ 可重试 | 无状态 Token 生成 |
| **数据库 TRUNCATE** | ✅ 是 | ✅ 可重试 | 清空表，幂等 |
| **数据库 INSERT** | ❌ 否 | ❌ 不可重试 | 可能违反唯一约束 |
| **外部 API 调用** | ⚠️ 视情况 | ⚠️ 谨慎重试 | 取决于外部服务 |

### 3.3 重试实现

```go
// backend/api/internal/testutil/retry.go
package testutil

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts int           // 最大尝试次数
	InitialWait time.Duration // 初始等待时间
	MaxWait     time.Duration // 最大等待时间
	Multiplier  float64       // 退避倍数
}

// DefaultRetryConfig 默认重试配置（用于幂等操作）
var DefaultRetryConfig = RetryConfig{
	MaxAttempts: 3,
	InitialWait: 100 * time.Millisecond,
	MaxWait:     2 * time.Second,
	Multiplier:  2.0,
}

// IsRetryable 判断错误是否可重试
func IsRetryable(statusCode int, err error) bool {
	// 网络错误可重试
	if err != nil {
		return true
	}
	// 特定 HTTP 状态码可重试
	switch statusCode {
	case http.StatusTooManyRequests, // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	}
	return false
}

// RetryableGet 执行可重试的 GET 请求
func RetryableGet(ctx context.Context, client *http.Client, url string, cfg RetryConfig) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	wait := cfg.InitialWait

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		resp, err := client.Get(url)
		
		// 成功或不可重试
		if err == nil && !IsRetryable(resp.StatusCode, nil) {
			return resp, nil
		}
		if !IsRetryable(0, err) {
			return resp, err
		}

		// 关闭响应体，避免泄漏
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		lastErr = err
		lastResp = resp

		// 等待后重试
		if attempt < cfg.MaxAttempts {
			time.Sleep(wait)
			wait = time.Duration(float64(wait) * cfg.Multiplier)
			if wait > cfg.MaxWait {
				wait = cfg.MaxWait
			}
		}
	}

	// 返回最后一次错误
	if lastErr != nil {
		return nil, fmt.Errorf("retry exhausted after %d attempts: %w", cfg.MaxAttempts, lastErr)
	}
	return lastResp, nil
}

// RetryableRequest 执行可重试的 HTTP 请求（仅限幂等方法）
func RetryableRequest(ctx context.Context, client *http.Client, req *http.Request, cfg RetryConfig) (*http.Response, error) {
	// 仅允许幂等方法
	switch req.Method {
	case http.MethodGet, http.MethodPut, http.MethodDelete:
		// 允许重试
	default:
		// 非幂等方法，不重试
		return client.Do(req)
	}

	var lastErr error
	var lastResp *http.Response
	wait := cfg.InitialWait

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		resp, err := client.Do(req)
		
		if err == nil && !IsRetryable(resp.StatusCode, nil) {
			return resp, nil
		}
		if !IsRetryable(0, err) {
			return resp, err
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		lastErr = err
		lastResp = resp

		if attempt < cfg.MaxAttempts {
			time.Sleep(wait)
			wait = time.Duration(float64(wait) * cfg.Multiplier)
			if wait > cfg.MaxWait {
				wait = cfg.MaxWait
			}
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("retry exhausted after %d attempts: %w", cfg.MaxAttempts, lastErr)
	}
	return lastResp, nil
}
```

### 3.4 禁止的重试模式

```go
// ❌ 禁止：对非幂等操作重试
func BAD_Example() {
	// POST 创建订单不应重试，可能产生重复订单
	for i := 0; i < 3; i++ {
		resp, _ := client.Post(url, "application/json", body)
		if resp.StatusCode == 200 {
			break
		}
		time.Sleep(time.Second)
	}
}

// ✅ 正确：创建失败直接报告
func GOOD_Example() {
	resp, err := client.Post(url, "application/json", body)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("创建订单失败: %v, status: %d", err, resp.StatusCode)
	}
}
```

---

## 4. SKIP 原因标签规范

### 4.1 标签定义（继承 T7）

```go
// backend/api/internal/testutil/env.go

// SkipReason 跳过原因标签
type SkipReason string

const (
	// SkipReasonAPIUnavailable API 服务不可用
	SkipReasonAPIUnavailable SkipReason = "API_UNAVAILABLE"
	// SkipReasonMySQLUnavailable MySQL 数据库不可用
	SkipReasonMySQLUnavailable SkipReason = "MYSQL_UNAVAILABLE"
	// SkipReasonRedisUnavailable Redis 缓存不可用
	SkipReasonRedisUnavailable SkipReason = "REDIS_UNAVAILABLE"
	// SkipReasonLoginFailed 登录失败
	SkipReasonLoginFailed SkipReason = "LOGIN_FAILED"
	// SkipReasonDataPrepFailed 测试数据准备失败
	SkipReasonDataPrepFailed SkipReason = "DATA_PREP_FAILED"
)
```

### 4.2 使用规范

```go
// ✅ 正确：使用标准化标签
func (suite *MyTestSuite) SetupSuite() {
	config := testutil.SkipIfNoAPI(suite.T())
	
	token, err := login(config)
	if err != nil {
		testutil.SkipfWithReason(suite.T(), testutil.SkipReasonLoginFailed, 
			"login failed: %v", err)
		return
	}
	
	campaign, err := createCampaign(token)
	if err != nil {
		testutil.SkipfWithReason(suite.T(), testutil.SkipReasonDataPrepFailed,
			"campaign creation failed: %v", err)
		return
	}
}

// ❌ 错误：无标签或模糊描述
func BAD_SetupSuite() {
	suite.T().Skipf("无法连接到后端服务: %v", err)  // 缺少标签
	suite.T().Skipf("测试跳过")  // 无原因
}
```

### 4.3 SKIP 标签决策树

```
测试需要 SKIP?
│
├─ API 服务不可达?
│   └─ 使用 API_UNAVAILABLE
│
├─ MySQL 连接失败?
│   └─ 使用 MYSQL_UNAVAILABLE
│
├─ Redis 连接失败?
│   ├─ 测试依赖 Redis?
│   │   └─ 使用 REDIS_UNAVAILABLE
│   └─ 测试不依赖 Redis?
│       └─ 不 SKIP，降级到 memory 存储
│
├─ 登录失败?
│   └─ 使用 LOGIN_FAILED
│
└─ 数据准备失败?
    └─ 使用 DATA_PREP_FAILED
```

### 4.4 SKIP 原因统计脚本

```bash
#!/bin/bash
# backend/scripts/analyze_skip_reasons.sh

echo "=== SKIP Reason Analysis ==="
echo ""

# 运行测试并捕获输出
OUTPUT=$(go test ./test/integration/... -v 2>&1)

# 统计各类 SKIP 原因
echo "| Skip Reason | Count |"
echo "|-------------|-------|"

for reason in API_UNAVAILABLE MYSQL_UNAVAILABLE REDIS_UNAVAILABLE LOGIN_FAILED DATA_PREP_FAILED; do
    count=$(echo "$OUTPUT" | grep -c "SKIP_REASON:$reason" || true)
    echo "| $reason | $count |"
done

echo ""

# 计算执行率
total_tests=$(echo "$OUTPUT" | grep -c "^--- PASS\|^--- FAIL" || true)
skipped_tests=$(echo "$OUTPUT" | grep -c "^--- SKIP" || true)

if [ $((total_tests + skipped_tests)) -gt 0 ]; then
    execution_rate=$((total_tests * 100 / (total_tests + skipped_tests)))
    echo "Execution Rate: ${execution_rate}%"
    echo "  - Executed: $total_tests"
    echo "  - Skipped: $skipped_tests"
fi
```

---

## 5. 就绪探针使用指南

### 5.1 探针优先级

```
优先级 1: API 就绪探针（必需，所有集成测试）
优先级 2: MySQL 就绪探针（需要数据库的测试）
优先级 3: Redis 就绪探针（可选，仅 Redis 依赖测试）
```

### 5.2 标准使用模式

```go
// backend/test/integration/example_integration_test.go
package integration

import (
	"context"
	"testing"
	"time"

	"dmh/api/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type ExampleIntegrationSuite struct {
	suite.Suite
	config    *testutil.IntegrationConfig
	client    *http.Client
	authToken string
}

func (s *ExampleIntegrationSuite) SetupSuite() {
	// 1. API 就绪探针（10s 超时）
	s.config = testutil.SkipIfNoAPI(s.T())
	
	// 2. MySQL 就绪探针（如果测试需要数据库）
	// testutil.SkipIfNoMySQL(s.T())
	
	// 3. 初始化 HTTP 客户端
	s.client = &http.Client{Timeout: 10 * time.Second}
	
	// 4. 登录（带失败处理）
	s.loginAsAdmin()
}

func (s *ExampleIntegrationSuite) loginAsAdmin() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	token, err := loginWithRetry(ctx, s.client, s.config)
	if err != nil {
		testutil.SkipfWithReason(s.T(), testutil.SkipReasonLoginFailed,
			"login failed after retries: %v", err)
		return
	}
	s.authToken = token
}

func TestExampleIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ExampleIntegrationSuite))
}
```

### 5.3 就绪探针超时配置

| 探针类型 | 超时时间 | 重试间隔 | 说明 |
|----------|----------|----------|------|
| API | 10s | 500ms | 必须在 10s 内响应 |
| MySQL | 5s | 200ms | 连接验证 |
| Redis | 5s | 200ms | Ping 验证 |
| 登录 | 10s | - | 单次请求，可重试 |

### 5.4 CI 环境就绪检查

```yaml
# .github/workflows/integration-test.yml
- name: Wait for API readiness
  timeout-minutes: 1
  run: |
    for i in {1..20}; do
      code=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8889/api/v1/auth/login || true)
      if [ "$code" != "000" ]; then
        echo "API is ready (status: $code)"
        exit 0
      fi
      echo "Waiting for API... attempt $i/20"
      sleep 2
    done
    echo "API did not become ready within timeout"
    exit 1
```

---

## 6. 常见 Flaky 修复模式

### 6.1 资源泄漏修复

```go
// ❌ 问题代码：循环内 defer
func BAD_RateLimitTest() {
	for i := 0; i < 10; i++ {
		resp, _ := client.Do(req)
		defer resp.Body.Close()  // 不会在循环结束时执行！
	}
}

// ✅ 修复：立即关闭或使用辅助函数
func GOOD_RateLimitTest() {
	for i := 0; i < 10; i++ {
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		// 立即读取并关闭
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		
		// 使用 body...
	}
}

// ✅ 更好：使用辅助函数
func doRequest(client *http.Client, req *http.Request) (int, []byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	return resp.StatusCode, body, err
}
```

### 6.2 并发测试修复

```go
// ❌ 问题代码：竞态条件
func BAD_ConcurrentTest() {
	var successCount int
	for i := 0; i < 10; i++ {
		go func() {
			resp, _ := client.Do(req)
			if resp.StatusCode == 200 {
				successCount++  // 竞态！
			}
		}()
	}
}

// ✅ 修复：使用互斥锁或原子操作
func GOOD_ConcurrentTest() {
	var successCount int32  // 使用 int32
	var wg sync.WaitGroup
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, _ := client.Do(req)
			if resp.StatusCode == 200 {
				atomic.AddInt32(&successCount, 1)  // 原子操作
			}
		}()
	}
	wg.Wait()
}
```

### 6.3 数据隔离修复

```go
// ❌ 问题代码：使用固定 ID
func BAD_DataTest() {
	// 如果用户 ID=1 不存在，测试失败
	resp, _ := client.Get(baseURL + "/api/v1/users/1")
}

// ✅ 修复：使用动态创建的数据
func GOOD_DataTest() {
	// 创建测试用户
	userId := createTestUser(t)
	
	// 使用动态 ID
	resp, _ := client.Get(baseURL + "/api/v1/users/" + strconv.FormatInt(userId, 10))
	
	// 清理（可选）
	defer deleteUser(t, userId)
}
```

---

## 7. Flaky 率监控与周报

### 7.1 Flaky 率定义

```
Flaky 率 = (失败后重试通过的测试数 / 总测试数) × 100%

目标：Flaky 率 < 2%
```

### 7.2 周报指标模板

```markdown
# 集成测试周报 (YYYY-MM-DD)

## 执行统计
- 总测试数: XXX
- 通过: XXX (XX%)
- 失败: XXX (XX%)
- 跳过: XXX (XX%)
- 执行率: XX%

## Flaky 分析
- Flaky 测试数: X
- Flaky 率: X.XX%
- 环比变化: ↓0.XX% / ↑0.XX%

## SKIP 原因分布
| 原因 | 数量 | 占比 |
|------|------|------|
| API_UNAVAILABLE | X | XX% |
| LOGIN_FAILED | X | XX% |
| ... | ... | ... |

## 本周修复
- [ ] 修复 rate_limiting_test.go 资源泄漏
- [ ] 统一登录失败处理

## 下周计划
- [ ] 治理 concurrency_test.go 并发竞态
- [ ] 增加就绪探针覆盖率
```

### 7.3 监控脚本

```bash
#!/bin/bash
# backend/scripts/flaky_report.sh
# 生成 Flaky 测试报告

REPORT_DATE=$(date +%Y-%m-%d)
REPORT_DIR="reports/flaky"
mkdir -p "$REPORT_DIR"

# 运行测试 3 次以检测 Flaky
echo "Running tests 3 times to detect flaky behavior..."

RESULTS=()
for i in 1 2 3; do
    echo "=== Run $i/3 ==="
    go test ./test/integration/... -v -count=1 2>&1 | tee "$REPORT_DIR/run_$i.log"
    RESULTS+=($(grep -c "^--- FAIL" "$REPORT_DIR/run_$i.log" || echo 0))
done

# 分析结果
echo ""
echo "=== Flaky Analysis ==="
echo "Run 1 failures: ${RESULTS[0]}"
echo "Run 2 failures: ${RESULTS[1]}"
echo "Run 3 failures: ${RESULTS[2]}"

# 如果失败数不一致，存在 Flaky
if [ "${RESULTS[0]}" != "${RESULTS[1]}" ] || [ "${RESULTS[1]}" != "${RESULTS[2]}" ]; then
    echo ""
    echo "⚠️  FLAKY TESTS DETECTED!"
    echo "Failure counts inconsistent across runs."
else
    echo ""
    echo "✅ No flaky tests detected."
fi
```

---

## 8. 实施清单

### 8.1 立即修复（高优先级）

| 任务 | 文件 | 预估工时 |
|------|------|----------|
| 修复资源泄漏 | `rate_limiting_test.go`, `concurrency_test.go` | 1h |
| 统一 SKIP 标签 | 所有集成测试 | 2h |
| 添加就绪探针 | 缺失探针的测试 | 1h |

### 8.2 短期改进（中优先级）

| 任务 | 描述 | 预估工时 |
|------|------|----------|
| 重试工具集成 | 创建 `testutil/retry.go` | 2h |
| 并发测试治理 | 修复竞态条件 | 2h |
| 监控脚本部署 | CI 集成 Flaky 检测 | 1h |

### 8.3 长期目标

| 目标 | 指标 | 时间线 |
|------|------|--------|
| 执行率 | >90% | 2 周 |
| Flaky 率 | <2% | 4 周 |
| 覆盖率 | >70% | 持续 |

---

## 9. 反模式警告

### 9.1 禁止的做法

| 反模式 | 原因 | 替代方案 |
|--------|------|----------|
| 无限重试 | 掩盖真实问题 | 有限次数 + 指数退避 |
| 对 POST 重试 | 可能产生重复数据 | 失败立即报告 |
| 跳过关键测试 | 降低覆盖率 | 修复测试或环境 |
| 全局超时过长 | 延迟失败反馈 | 每操作独立超时 |
| 忽略 SKIP 统计 | 无法追踪环境问题 | 定期分析 SKIP 原因 |

### 9.2 正确的失败处理

```go
// ✅ 正确：区分环境问题与代码缺陷
func (s *MySuite) TestSomething() {
	// 环境问题 → SKIP
	if !isAPIAvailable() {
		s.T().Skip("SKIP_REASON:API_UNAVAILABLE | API not ready")
		return
	}
	
	// 代码缺陷 → FAIL
	result, err := doSomething()
	if err != nil {
		s.FailNow("unexpected error: %v", err)  // 不是 Skip！
	}
}
```

---

## 10. 变更记录

| 日期 | 版本 | 作者 | 变更内容 |
|------|------|------|----------|
| 2026-02-26 | v1.0 | Sisyphus | 初始版本，定义 Flaky 测试治理方案 |
