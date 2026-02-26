# T7: 集成测试环境标准化

## 1. 概述

本文档定义 DMH 项目集成测试环境标准化方案，包括：
- API 启动检查
- MySQL 初始化策略
- Redis 可用性探测（可选）
- SKIP 条件与原因标签
- 快速失败策略

---

## 2. 环境变量规范

### 2.1 标准环境变量

| 变量名 | 默认值 | 用途 | 必需 |
|--------|--------|------|------|
| `DMH_INTEGRATION_BASE_URL` | `http://localhost:8889` | API 服务地址 | 是 |
| `DMH_TEST_ADMIN_USERNAME` | `admin` | 测试管理员账号 | 是 |
| `DMH_TEST_ADMIN_PASSWORD` | `123456` | 测试管理员密码 | 是 |
| `MYSQL_TEST_HOST` | `127.0.0.1` | MySQL 主机 | MySQL 测试必需 |
| `MYSQL_TEST_PORT` | `3306` | MySQL 端口 | MySQL 测试必需 |
| `MYSQL_TEST_USER` | `root` | MySQL 用户 | MySQL 测试必需 |
| `MYSQL_TEST_PASSWORD` | `Admin168` | MySQL 密码 | MySQL 测试必需 |
| `MYSQL_TEST_DB` | `dmh` | MySQL 基础数据库 | MySQL 测试必需 |
| `REDIS_TEST_HOST` | `localhost:6379` | Redis 地址 | 可选 |

### 2.2 环境变量读取代码

```go
// backend/api/internal/testutil/env.go
package testutil

import (
	"os"
	"strings"
)

// IntegrationConfig 集成测试配置
type IntegrationConfig struct {
	BaseURL       string
	AdminUsername string
	AdminPassword string
}

// GetIntegrationConfig 获取集成测试配置
func GetIntegrationConfig() *IntegrationConfig {
	return &IntegrationConfig{
		BaseURL:       getEnv("DMH_INTEGRATION_BASE_URL", "http://localhost:8889"),
		AdminUsername: getEnv("DMH_TEST_ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("DMH_TEST_ADMIN_PASSWORD", "123456"),
	}
}

// getEnv 获取环境变量或默认值
func getEnv(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}
```

---

## 3. API 启动检查

### 3.1 就绪探针设计

**超时设置**：≤10 秒（符合 MUST DO 要求）

```go
// backend/api/internal/testutil/api_probe.go
package testutil

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	// APIReadinessTimeout API 就绪检查超时
	APIReadinessTimeout = 10 * time.Second
	// APIProbeInterval 探测间隔
	APIProbeInterval = 500 * time.Millisecond
)

// APIProbeResult API 探测结果
type APIProbeResult struct {
	Available bool
	Status    int
	Error     error
	Latency   time.Duration
}

// ProbeAPI 检查 API 是否就绪
func ProbeAPI(baseURL string) APIProbeResult {
	start := time.Now()
	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(baseURL + "/api/v1/auth/login")
	latency := time.Since(start)

	if err != nil {
		return APIProbeResult{
			Available: false,
			Error:     err,
			Latency:   latency,
		}
	}
	defer resp.Body.Close()

	return APIProbeResult{
		Available: true,
		Status:    resp.StatusCode,
		Latency:   latency,
	}
}

// WaitForAPI 等待 API 就绪（带超时）
func WaitForAPI(ctx context.Context, baseURL string) error {
	ticker := time.NewTicker(APIProbeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("API readiness check timed out: %v", ctx.Err())
		case <-ticker.C:
			result := ProbeAPI(baseURL)
			if result.Available {
				return nil
			}
		}
	}
}

// SkipIfNoAPI 跳过测试如果 API 不可用
func SkipIfNoAPI(t *testing.T) *IntegrationConfig {
	t.Helper()
	config := GetIntegrationConfig()

	ctx, cancel := context.WithTimeout(context.Background(), APIReadinessTimeout)
	defer cancel()

	if err := WaitForAPI(ctx, config.BaseURL); err != nil {
		t.Skipf("SKIP_REASON:API_UNAVAILABLE | API at %s not ready within %v: %v",
			config.BaseURL, APIReadinessTimeout, err)
		return nil
	}

	return config
}
```

### 3.2 使用示例

```go
func TestMyIntegration(t *testing.T) {
	config := testutil.SkipIfNoAPI(t)
	// API 已就绪，继续测试...
}
```

---

## 4. MySQL 初始化策略

### 4.1 现有工具复用

项目已有 `backend/api/internal/testutil/mysql_test_helper.go`，提供：
- `SkipIfNoMySQL(t)` - 检查 MySQL 可用性
- `SetupMySQLTestDB(t)` - 创建隔离数据库
- `CleanupMySQLTestDB(t, dbName)` - 清理测试数据库

### 4.2 SKIP 条件标签

```go
// SkipReason 跳过原因标签
type SkipReason string

const (
	SkipReasonAPIUnavailable  SkipReason = "API_UNAVAILABLE"
	SkipReasonMySQLUnavailable SkipReason = "MYSQL_UNAVAILABLE"
	SkipReasonRedisUnavailable SkipReason = "REDIS_UNAVAILABLE"
	SkipReasonLoginFailed      SkipReason = "LOGIN_FAILED"
	SkipReasonDataPrepFailed   SkipReason = "DATA_PREP_FAILED"
)

// SkipfWithReason 带原因标签的跳过
func SkipfWithReason(t *testing.T, reason SkipReason, format string, args ...interface{}) {
	t.Helper()
	t.Skipf("SKIP_REASON:%s | %s", reason, fmt.Sprintf(format, args...))
}
```

### 4.3 快速失败策略

**原则**：
- 连接失败立即跳过（不重试）
- 登录失败立即跳过（不重试）
- 数据准备失败立即跳过（不重试）
- 超时 ≤10 秒

```go
// EnsureTestEnvironment 确保测试环境就绪
// 快速失败：任一依赖不可用立即跳过
func EnsureTestEnvironment(t *testing.T) (*IntegrationConfig, *gorm.DB) {
	t.Helper()

	// 1. 检查 API（10s 超时）
	config := SkipIfNoAPI(t)

	// 2. 检查 MySQL（仅对需要 DB 的测试）
	var db *gorm.DB
	if needsDatabase(t) {
		SkipIfNoMySQL(t)
		db, _ = SetupMySQLTestDB(t)
	}

	return config, db
}

// needsDatabase 判断测试是否需要数据库
func needsDatabase(t *testing.T) bool {
	// 通过 build tag 或测试名称判断
	name := t.Name()
	return strings.Contains(name, "Repository") ||
		strings.Contains(name, "Database") ||
		strings.Contains(name, "Integration")
}
```

---

## 5. Redis 可用性探测（可选）

### 5.1 设计原则

Redis 在 DMH 中是**可选依赖**：
- 默认配置使用 `memory` 存储
- 仅生产/高负载环境需要 Redis
- 集成测试不应因 Redis 不可用而失败

### 5.2 探测实现

```go
// backend/api/internal/testutil/redis_probe.go
package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// RedisProbeTimeout Redis 探测超时
	RedisProbeTimeout = 5 * time.Second
)

// ProbeRedis 检查 Redis 是否可用
func ProbeRedis(addr string) error {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), RedisProbeTimeout)
	defer cancel()

	return client.Ping(ctx).Err()
}

// SkipIfNoRedis 跳过测试如果 Redis 不可用（仅用于 Redis 依赖的测试）
func SkipIfNoRedis(t *testing.T) {
	t.Helper()
	addr := getEnv("REDIS_TEST_HOST", "localhost:6379")

	if err := ProbeRedis(addr); err != nil {
		t.Skipf("SKIP_REASON:REDIS_UNAVAILABLE | Redis at %s not available: %v", addr, err)
	}
}

// IsRedisAvailable 检查 Redis 是否可用（不跳过测试）
func IsRedisAvailable() bool {
	addr := getEnv("REDIS_TEST_HOST", "localhost:6379")
	return ProbeRedis(addr) == nil
}
```

### 5.3 使用场景

```go
// 仅当 Redis 可用时运行的测试
func TestRateLimitingWithRedis(t *testing.T) {
	testutil.SkipIfNoRedis(t)
	// Redis 依赖的测试...
}

// 降级到内存存储的测试
func TestRateLimiting(t *testing.T) {
	storage := "memory"
	if testutil.IsRedisAvailable() {
		storage = "redis"
	}
	// 使用 storage 进行测试...
}
```

---

## 6. 标准化集成测试模板

### 6.1 完整模板

```go
// backend/test/integration/example_integration_test.go
package integration

import (
	"net/http"
	"testing"
	"time"

	"dmh/api/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestExampleIntegration(t *testing.T) {
	// 1. 环境检查（快速失败 ≤10s）
	config := testutil.SkipIfNoAPI(t)

	// 2. 登录获取 Token
	client := &http.Client{Timeout: 10 * time.Second}
	token, err := loginAndGetToken(client, config.BaseURL, config.AdminUsername, config.AdminPassword)
	if err != nil {
		testutil.SkipfWithReason(t, testutil.SkipReasonLoginFailed, "login failed: %v", err)
		return
	}

	// 3. 测试逻辑
	t.Run("some test case", func(t *testing.T) {
		// ...
	})
}
```

### 6.2 Suite 模式模板

```go
type ExampleIntegrationSuite struct {
	suite.Suite
	config    *testutil.IntegrationConfig
	client    *http.Client
	authToken string
}

func (s *ExampleIntegrationSuite) SetupSuite() {
	// 快速失败：API 不可用立即跳过
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.config = testutil.GetIntegrationConfig()
	if err := testutil.WaitForAPI(ctx, s.config.BaseURL); err != nil {
		s.T().Skipf("SKIP_REASON:API_UNAVAILABLE | %v", err)
		return
	}

	s.client = &http.Client{Timeout: 10 * time.Second}
	s.loginAsAdmin()
}

func (s *ExampleIntegrationSuite) loginAsAdmin() {
	token, err := loginAndGetToken(s.client, s.config.BaseURL,
		s.config.AdminUsername, s.config.AdminPassword)
	if err != nil {
		s.T().Skipf("SKIP_REASON:LOGIN_FAILED | %v", err)
		return
	}
	s.authToken = token
}

func TestExampleIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ExampleIntegrationSuite))
}
```

---

## 7. CI 环境集成

### 7.1 GitHub Actions 服务配置

```yaml
# .github/workflows/integration-test.yml
services:
  mysql8:
    image: mysql:8.0
    env:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: dmh
    ports:
      - 3306:3306
    options: >-
      --health-cmd="mysqladmin ping -h 127.0.0.1 -proot --silent"
      --health-interval=10s
      --health-timeout=5s
      --health-retries=20

  redis:
    image: redis:7
    ports:
      - 6379:6379
    options: >-
      --health-cmd="redis-cli ping"
      --health-interval=10s
      --health-timeout=5s
      --health-retries=10
```

### 7.2 环境变量配置

```yaml
env:
  DMH_INTEGRATION_BASE_URL: http://127.0.0.1:8889
  DMH_TEST_ADMIN_USERNAME: admin
  DMH_TEST_ADMIN_PASSWORD: "123456"
  MYSQL_TEST_HOST: 127.0.0.1
  MYSQL_TEST_PORT: "3306"
  MYSQL_TEST_USER: root
  MYSQL_TEST_PASSWORD: root
  MYSQL_TEST_DB: dmh
  REDIS_TEST_HOST: 127.0.0.1:6379
```

### 7.3 API 启动与就绪检查

```yaml
- name: Start API server
  run: |
    cd backend
    nohup go run ./api/dmh.go -f ./api/etc/dmh-api.yaml > /tmp/dmh-api.log 2>&1 &
    echo $! > /tmp/dmh-api.pid

- name: Wait for API readiness
  timeout-minutes: 1
  run: |
    for i in {1..20}; do
      code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8889/api/v1/auth/login || true)
      if [ "$code" != "000" ]; then
        echo "API is ready (status: $code)"
        exit 0
      fi
      sleep 2
    done
    echo "API did not become ready"
    cat /tmp/dmh-api.log || true
    exit 1
```

---

## 8. SKIP 条件总结

| SKIP_REASON | 触发条件 | 处理方式 |
|-------------|----------|----------|
| `API_UNAVAILABLE` | API 服务不可达或超时 | 立即跳过，不重试 |
| `MYSQL_UNAVAILABLE` | MySQL 连接失败 | 立即跳过，不重试 |
| `REDIS_UNAVAILABLE` | Redis 连接失败（仅 Redis 依赖测试） | 立即跳过，不重试 |
| `LOGIN_FAILED` | 登录接口返回非 200 | 立即跳过，不重试 |
| `DATA_PREP_FAILED` | 测试数据准备失败 | 立即跳过，不重试 |

---

## 9. 现有脚本改造建议

### 9.1 `run_order_mysql8_regression.sh` 改造

现有脚本已有 SKIP 检测（第 26 行），建议增强：

```bash
# 增强版 SKIP 检测
if grep -E -q -- '--- SKIP:' "$OUT_FILE"; then
  echo ""
  echo "[order-mysql8-regression] ERROR: tests were skipped."
  echo ""
  echo "Skip reasons found:"
  grep -oE 'SKIP_REASON:[^|]+' "$OUT_FILE" | sort | uniq -c
  echo ""
  echo "Check environment:"
  echo "  DMH_INTEGRATION_BASE_URL=${DMH_INTEGRATION_BASE_URL}"
  echo "  DMH_TEST_ADMIN_USERNAME=${DMH_TEST_ADMIN_USERNAME}"
  echo "  DMH_TEST_ADMIN_PASSWORD=<your-password>"
  exit 2
fi
```

### 9.2 统一启动脚本

```bash
#!/usr/bin/env bash
# backend/scripts/run_integration_tests.sh
# 集成测试统一启动脚本

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

# 环境变量默认值
export DMH_INTEGRATION_BASE_URL="${DMH_INTEGRATION_BASE_URL:-http://localhost:8889}"
export DMH_TEST_ADMIN_USERNAME="${DMH_TEST_ADMIN_USERNAME:-admin}"
export DMH_TEST_ADMIN_PASSWORD="${DMH_TEST_ADMIN_PASSWORD:-123456}"

echo "=== Integration Test Environment ==="
echo "API URL: ${DMH_INTEGRATION_BASE_URL}"
echo "Admin User: ${DMH_TEST_ADMIN_USERNAME}"
echo ""

# API 就绪检查（10s 超时）
echo "Checking API readiness..."
for i in {1..20}; do
  if curl -s -o /dev/null -w "%{http_code}" "${DMH_INTEGRATION_BASE_URL}/api/v1/auth/login" | grep -qv "000"; then
    echo "API is ready"
    break
  fi
  if [ $i -eq 20 ]; then
    echo "ERROR: API not ready within 10 seconds"
    exit 1
  fi
  sleep 0.5
done

# 运行测试
echo ""
echo "Running integration tests..."
go test ./test/integration/... -v -count=1 -timeout 5m
```

---

## 10. 集成测试执行率目标

### 10.1 当前状态

| 测试类型 | 文件数 | SKIP 风险 |
|----------|--------|-----------|
| 集成测试 | 25 | 中（依赖 API/MySQL） |

### 10.2 执行率目标

**目标：>90%**

实现策略：
1. CI 环境确保 API/MySQL/Redis 服务就绪
2. 本地开发提供 Docker Compose 快速启动
3. SKIP 必须带原因标签，便于分析
4. 定期统计 SKIP 比例，识别环境问题

### 10.3 监控指标

```bash
# 统计 SKIP 原因分布
go test ./test/integration/... -v 2>&1 | grep -oE 'SKIP_REASON:[^|]+' | sort | uniq -c

# 计算执行率
total=$(go test ./test/integration/... -v 2>&1 | grep -c "^--- PASS\|^--- FAIL")
skipped=$(go test ./test/integration/... -v 2>&1 | grep -c "^--- SKIP")
execution_rate=$((total * 100 / (total + skipped)))
echo "Execution rate: ${execution_rate}%"
```

---

## 11. 快速失败策略总结

| 场景 | 超时 | 行为 |
|------|------|------|
| API 探测 | ≤10s | 超时立即跳过 |
| MySQL 连接 | ≤5s | 失败立即跳过 |
| Redis 连接 | ≤5s | 失败立即跳过（仅 Redis 测试） |
| 登录请求 | ≤10s | 失败立即跳过 |
| 数据准备 | ≤30s | 失败立即跳过 |

---

## 12. 文件清单

| 文件 | 状态 | 说明 |
|------|------|------|
| `backend/api/internal/testutil/mysql_test_helper.go` | ✅ 已存在 | MySQL 测试工具 |
| `backend/api/internal/testutil/env.go` | 🆕 新增 | 环境变量配置 |
| `backend/api/internal/testutil/api_probe.go` | 🆕 新增 | API 就绪探测 |
| `backend/api/internal/testutil/redis_probe.go` | 🆕 新增 | Redis 可用性探测 |
| `backend/scripts/run_integration_tests.sh` | 🆕 新增 | 统一启动脚本 |

---

## 13. 变更记录

| 日期 | 版本 | 作者 | 变更内容 |
|------|------|------|----------|
| 2026-02-26 | v1.0 | Sisyphus | 初始版本，定义集成测试环境标准化方案 |
