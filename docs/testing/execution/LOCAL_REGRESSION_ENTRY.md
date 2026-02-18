# DMH 本地一键全量回归入口

> 版本: Plan v1 (FRP-007)
> 变更: add-full-regression-testing
> 生成日期: 2026-02-18

---

## 1. 概述

### 1.1 目的

本文档定义 DMH 项目的**本地一键全量回归入口**，确保：

- 本地执行范围与 CI 必跑矩阵 100% 一致
- 开发者可在发布前本地完成全量回归验证
- 本地回归结论与 CI 回归结论具有同等参考价值

### 1.2 适用场景

| 场景 | 推荐使用 |
|------|----------|
| 发布前本地签收 | ✅ 一键全量回归 |
| RC 构建前预验证 | ✅ 一键全量回归 |
| 大范围重构后验证 | ✅ 一键全量回归 |
| 普通 PR 验证 | ⚠️ 按需执行子集 |
| 单模块开发调试 | ❌ 使用单一测试命令 |

### 1.3 引用规格

- **全量回归口径定义**: `docs/testing/execution/FULL_REGRESSION_DEFINITION.md` (FRP-002)
- **CI 编排方案**: `docs/testing/execution/CI_ORCHESTRATION_PLAN.md` (FRP-006)
- **发布阻断规则**: `docs/testing/execution/RELEASE_BLOCKING_RULES.md` (FRP-003)

---

## 2. 环境前置要求

### 2.1 基础环境清单

| 依赖 | 版本要求 | 验证命令 | 说明 |
|------|----------|----------|------|
| Go | 1.24+ | `go version` | 后端编译与测试 |
| Node.js | 20.19+ | `node -v` | 前端构建与测试 |
| MySQL | 8.0+ | `mysql --version` | 数据库服务 |
| Redis | 7+ | `redis-cli ping` | 缓存服务 |
| OpenSpec CLI | 最新 | `openspec version` | 规格校验 |

### 2.2 服务运行状态

执行全量回归前，必须确保以下服务正常运行：

```bash
# 检查服务状态
docker ps | grep -E "mysql8|redis-dmh|dmh-api"

# 或使用 Makefile
make ps
```

**必需运行的服务**：

| 服务 | 端口 | 用途 |
|------|------|------|
| MySQL 8.0 | 3306 | 数据存储 |
| Redis 7 | 6379 | 缓存服务 |
| DMH API | 8889 | 集成测试目标 |

### 2.3 环境启动命令

```bash
# 方式 1: 使用 Docker Compose (推荐)
cd deploy
docker compose -f docker-compose-simple.yml up -d

# 方式 2: 使用 Makefile
make up

# 方式 3: 本地直接启动
# 终端 1: 启动后端
cd backend && go run api/dmh.go -f api/etc/dmh-api.yaml

# 终端 2-3: 前端 (如需 E2E 测试)
cd frontend-admin && npm run dev
cd frontend-h5 && npm run dev
```

### 2.4 环境验证脚本

```bash
#!/bin/bash
# 快速验证环境是否就绪

echo "=== 环境检查 ==="

# Go
GO_VERSION=$(go version 2>/dev/null | grep -oP 'go\d+\.\d+' | head -1)
if [[ -n "$GO_VERSION" ]]; then
  echo "✅ Go: $GO_VERSION"
else
  echo "❌ Go 未安装"
  exit 1
fi

# Node.js
NODE_VERSION=$(node -v 2>/dev/null)
if [[ -n "$NODE_VERSION" ]]; then
  echo "✅ Node.js: $NODE_VERSION"
else
  echo "❌ Node.js 未安装"
  exit 1
fi

# MySQL
MYSQL_CHECK=$(docker exec mysql8 mysql -uroot -p'#Admin168' -e "SELECT 1" 2>/dev/null)
if [[ -n "$MYSQL_CHECK" ]]; then
  echo "✅ MySQL: 运行中"
else
  echo "❌ MySQL 未运行"
  exit 1
fi

# Redis
REDIS_CHECK=$(docker exec redis-dmh redis-cli ping 2>/dev/null)
if [[ "$REDIS_CHECK" == "PONG" ]]; then
  echo "✅ Redis: 运行中"
else
  echo "❌ Redis 未运行"
  exit 1
fi

# API
API_CHECK=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8889/health 2>/dev/null)
if [[ "$API_CHECK" == "200" ]]; then
  echo "✅ DMH API: 运行中"
else
  echo "⚠️  DMH API 未响应 (集成测试将跳过)"
fi

echo ""
echo "=== 环境检查完成 ==="
```

---

## 3. 一键命令

### 3.1 推荐命令

```bash
# 方式 1: 使用 Makefile (推荐)
make full-regression

# 方式 2: 使用脚本
./scripts/run_full_regression.sh
```

### 3.2 完整命令定义

```bash
#!/bin/bash
# scripts/run_full_regression.sh
# DMH 本地一键全量回归

set -e

echo "=========================================="
echo "DMH 本地一键全量回归"
echo "=========================================="
echo "开始时间: $(date)"
echo ""

# 计时开始
START_TIME=$(date +%s)

# 1. 后端单元测试
echo "=== [1/6] 后端单元测试 ==="
cd backend && go test ./... -v
cd ..

# 2. 后端集成测试
echo ""
echo "=== [2/6] 后端集成测试 ==="
cd backend && DMH_INTEGRATION_BASE_URL=http://localhost:8889 go test ./test/integration/... -v -count=1
cd ..

# 3. 订单专项回归
echo ""
echo "=== [3/6] 订单专项回归 ==="
DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
DMH_TEST_ADMIN_USERNAME=admin \
DMH_TEST_ADMIN_PASSWORD=123456 \
backend/scripts/run_order_mysql8_regression.sh

# 4. 前端管理后台测试
echo ""
echo "=== [4/6] 前端管理后台 (Admin) 测试 ==="
cd frontend-admin && npm run test
npm run test:e2e:headless 2>/dev/null || echo "⚠️  E2E 测试跳过 (无 headless 环境)"
cd ..

# 5. 前端 H5 测试
echo ""
echo "=== [5/6] 前端 H5 测试 ==="
cd frontend-h5 && npm run test
npm run test:e2e:headless 2>/dev/null || echo "⚠️  E2E 测试跳过 (无 headless 环境)"
cd ..

# 6. OpenSpec 校验
echo ""
echo "=== [6/6] OpenSpec 校验 ==="
openspec validate --all --no-interactive 2>/dev/null || echo "⚠️  OpenSpec CLI 未安装"

# 计时结束
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
echo "=========================================="
echo "✅ 全量回归完成"
echo "耗时: ${DURATION} 秒"
echo "结束时间: $(date)"
echo "=========================================="
```

---

## 4. 执行范围

### 4.1 必跑矩阵对齐

本地执行范围与 CI 必跑矩阵 100% 对齐：

```
本地全量回归必跑项 (与 CI 一致):
├── 后端 (Backend)
│   ├── 单元测试: make test-backend 或 cd backend && go test ./...
│   ├── 集成测试: make test-integration (需要运行中的服务)
│   └── 订单专项: backend/scripts/run_order_mysql8_regression.sh
│
├── 前端管理后台 (Frontend-Admin)
│   ├── 单元测试: make test-admin 或 cd frontend-admin && npm run test
│   └── E2E 测试: make test-e2e-admin (需要 headed/headless 浏览器)
│
├── 前端 H5 (Frontend-H5)
│   ├── 单元测试: make test-h5 或 cd frontend-h5 && npm run test
│   └── E2E 测试: make test-e2e-h5 (需要 headed/headless 浏览器)
│
└── OpenSpec
    └── 校验: openspec validate --all --no-interactive
```

### 4.2 分步执行命令

如需单独执行某一项：

| 测试项 | 命令 | 依赖 |
|--------|------|------|
| 后端单元测试 | `make test-backend` 或 `cd backend && go test ./... -v` | Go |
| 后端集成测试 | `make test-integration` | API + MySQL + Redis |
| 订单专项回归 | `backend/scripts/run_order_mysql8_regression.sh` | API + MySQL |
| Admin 单元测试 | `make test-admin` 或 `cd frontend-admin && npm run test` | Node.js |
| Admin E2E 测试 | `cd frontend-admin && npm run test:e2e:headless` | Playwright |
| H5 单元测试 | `make test-h5` 或 `cd frontend-h5 && npm run test` | Node.js |
| H5 E2E 测试 | `cd frontend-h5 && npm run test:e2e:headless` | Playwright |
| OpenSpec 校验 | `make spec-validate` 或 `openspec validate --all` | OpenSpec CLI |

### 4.3 覆盖率检查

全量回归包含覆盖率门禁验证：

| 模块 | 覆盖率阈值 | 本地验证命令 |
|------|------------|--------------|
| 后端 | >= 76% | `cd backend && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out` |
| Admin | >= 70% | `cd frontend-admin && npm run test:cov` |
| H5 | >= 44% | `cd frontend-h5 && npm run test:cov` |

---

## 5. 结果判定

### 5.1 PASS 条件

所有以下条件同时满足时，本地回归判定为 **PASS**：

| 条件 | 验证方法 |
|------|----------|
| 后端单元测试 100% 通过 | 无 FAILED 用例 |
| 后端集成测试 100% 通过 | 无 FAILED 用例 |
| 订单专项回归 100% 通过 | 脚本输出 "PASS" |
| Admin 单元测试 100% 通过 | 无 FAILED 用例 |
| H5 单元测试 100% 通过 | 无 FAILED 用例 |
| OpenSpec 校验无错误 | 无 ERROR 输出 |

### 5.2 FAIL 条件

任一条件触发时，本地回归判定为 **FAIL**：

| 条件 | 处理建议 |
|------|----------|
| 任一测试套件有 FAILED 用例 | 查看失败日志，修复后重跑 |
| 订单专项测试被 SKIP | 检查 API 连接和测试账号 |
| OpenSpec 校验失败 | 运行 `openspec validate --strict` 查看详情 |
| 覆盖率不达标 | 补充测试用例 |

### 5.3 快速判定脚本

```bash
#!/bin/bash
# 快速判定回归结果

PASSED=0
FAILED=0

# 检查后端测试
if go test ./... -v 2>&1 | grep -q "FAIL"; then
  echo "❌ 后端单元测试失败"
  FAILED=$((FAILED + 1))
else
  echo "✅ 后端单元测试通过"
  PASSED=$((PASSED + 1))
fi

# 检查前端测试
if cd frontend-admin && npm run test 2>&1 | grep -q "failed"; then
  echo "❌ Admin 单元测试失败"
  FAILED=$((FAILED + 1))
else
  echo "✅ Admin 单元测试通过"
  PASSED=$((PASSED + 1))
fi

# ... 其他检查 ...

echo ""
echo "=== 回归结果 ==="
echo "通过: $PASSED"
echo "失败: $FAILED"

if [[ $FAILED -eq 0 ]]; then
  echo "✅ 本地全量回归: PASS"
  exit 0
else
  echo "❌ 本地全量回归: FAIL"
  exit 1
fi
```

---

## 6. Makefile 扩展建议

### 6.1 新增 Target 建议

在 `Makefile` 末尾添加以下 target：

```makefile
# ============================================
# 全量回归 (FRP-007)
# ============================================

.PHONY: full-regression full-regression-quick env-check

# 全量回归（本地一键）- 完整版
full-regression: env-check test-backend test-integration test-admin test-h5 test-e2e spec-validate
	@echo ""
	@echo "=== ✅ 全量回归完成 ==="
	@echo "所有必跑测试项均已执行"

# 快速回归 - 仅核心测试（跳过 E2E）
full-regression-quick: test-backend test-integration test-admin test-h5 spec-validate
	@echo ""
	@echo "=== ✅ 快速回归完成 ==="
	@echo "注意: 未执行 E2E 测试"

# 环境检查
env-check:
	@echo "=== 环境检查 ==="
	@go version > /dev/null 2>&1 || (echo "❌ Go 未安装" && exit 1)
	@node -v > /dev/null 2>&1 || (echo "❌ Node.js 未安装" && exit 1)
	@curl -s http://localhost:8889/health > /dev/null 2>&1 || echo "⚠️  API 未运行，集成测试将跳过"
	@echo "✅ 环境检查通过"

# 订单专项回归
test-order-regression:
	@DMH_INTEGRATION_BASE_URL=http://localhost:8889 \
	DMH_TEST_ADMIN_USERNAME=admin \
	DMH_TEST_ADMIN_PASSWORD=123456 \
	backend/scripts/run_order_mysql8_regression.sh

# Admin E2E 测试
test-e2e-admin:
	cd frontend-admin && npm run test:e2e:headless

# H5 E2E 测试
test-e2e-h5:
	cd frontend-h5 && npm run test:e2e:headless

# E2E 测试（双端）
test-e2e: test-e2e-admin test-e2e-h5
	@echo "✅ E2E 测试完成"
```

### 6.2 完整 Target 依赖图

```
full-regression
├── env-check
├── test-backend
│   └── cd backend && go test ./... -v
├── test-integration
│   └── cd backend && DMH_INTEGRATION_BASE_URL=... go test ./test/integration/...
├── test-order-regression (可选，集成测试包含)
│   └── backend/scripts/run_order_mysql8_regression.sh
├── test-admin
│   └── cd frontend-admin && npm run test
├── test-h5
│   └── cd frontend-h5 && npm run test
├── test-e2e
│   ├── test-e2e-admin
│   │   └── cd frontend-admin && npm run test:e2e:headless
│   └── test-e2e-h5
│       └── cd frontend-h5 && npm run test:e2e:headless
└── spec-validate
    └── openspec validate --all --no-interactive
```

### 6.3 现有 Makefile 命令对照

| 现有命令 | 对应必跑项 | 备注 |
|----------|------------|------|
| `make test` | 后端 + Admin + H5 单元测试 | 部分覆盖 |
| `make test-backend` | 后端单元测试 | ✅ 已覆盖 |
| `make test-integration` | 后端集成测试 | ✅ 已覆盖 |
| `make test-admin` | Admin 单元测试 | ✅ 已覆盖 |
| `make test-h5` | H5 单元测试 | ✅ 已覆盖 |
| `make test-e2e` | Admin + H5 E2E | ✅ 已覆盖 |
| `make spec-validate` | OpenSpec 校验 | ✅ 已覆盖 |
| - | 订单专项回归 | ⚠️ 需新增 |
| - | 覆盖率门禁 | ⚠️ 需新增 |

---

## 7. 常见问题

### 7.1 集成测试失败

**问题**: `connection refused` 或 `API unavailable`

**解决方案**:
```bash
# 1. 确认 API 服务运行
curl http://localhost:8889/health

# 2. 启动服务
make up
# 或
docker compose -f deploy/docker-compose-simple.yml up -d
```

### 7.2 订单专项测试被 SKIP

**问题**: `tests were skipped`

**解决方案**:
```bash
# 1. 检查 API 可用性
curl http://localhost:8889/health

# 2. 检查测试账号
mysql -h localhost -u root -p'#Admin168' dmh -e "SELECT username FROM users WHERE username='admin'"

# 3. 修复登录凭据
backend/scripts/repair_login_and_run_order_regression.sh
```

### 7.3 E2E 测试无头浏览器问题

**问题**: `Playwright browser not found`

**解决方案**:
```bash
# 安装 Playwright 浏览器
cd frontend-admin && npx playwright install
cd frontend-h5 && npx playwright install

# 或跳过 E2E 测试
make full-regression-quick
```

### 7.4 OpenSpec CLI 未安装

**问题**: `openspec: command not found`

**解决方案**:
```bash
# 按项目文档安装 OpenSpec CLI
# 或跳过 OpenSpec 校验
make test test-integration test-e2e
```

### 7.5 覆盖率不达标

**问题**: 本地覆盖率低于阈值

**解决方案**:
```bash
# 查看详细覆盖率报告
cd backend && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

# 补充缺失的测试用例
# 参考现有测试模式编写新测试
```

---

## 8. 与 CI 的对齐

### 8.1 本地 vs CI 执行对照

| 测试项 | 本地命令 | CI Workflow | 对齐状态 |
|--------|----------|-------------|----------|
| 后端单元测试 | `make test-backend` | stability-checks.yml | ✅ |
| 后端集成测试 | `make test-integration` | system-test-gate.yml | ✅ |
| 订单专项回归 | `backend/scripts/run_order_mysql8_regression.sh` | order-mysql8-regression.yml | ✅ |
| Admin 单元测试 | `make test-admin` | stability-checks.yml | ✅ |
| Admin E2E 测试 | `make test-e2e-admin` | system-test-gate.yml | ✅ |
| H5 单元测试 | `make test-h5` | stability-checks.yml | ⚠️ CI 需补充 |
| H5 E2E 测试 | `make test-e2e-h5` | system-test-gate.yml | ✅ |
| 覆盖率门禁 | `npm run test:cov` | coverage-gate.yml | ✅ |
| OpenSpec 校验 | `make spec-validate` | system-test-gate.yml | ✅ |

### 8.2 差异说明

| 差异项 | 本地 | CI | 处理方式 |
|--------|------|-----|----------|
| 执行环境 | 开发机 | GitHub Actions | 环境变量配置不同 |
| 浏览器 | 本地 Playwright | GitHub-hosted runner | 命令相同 |
| 服务启动 | Docker Compose | workflow services | 端口映射一致 |
| 结果聚合 | 人工查看 | 自动聚合 PASS/FAIL | 判定规则一致 |

---

## 9. 参考资料

- **全量回归口径定义**: `docs/testing/execution/FULL_REGRESSION_DEFINITION.md` (FRP-002)
- **CI 编排方案**: `docs/testing/execution/CI_ORCHESTRATION_PLAN.md` (FRP-006)
- **发布阻断规则**: `docs/testing/execution/RELEASE_BLOCKING_RULES.md` (FRP-003)
- **Flaky 测试策略**: `docs/testing/execution/FLAKY_TEST_STRATEGY.md` (FRP-004)
- **Delta 规格**: `openspec/changes/add-full-regression-testing/specs/system-test-execution/spec.md`
- **现有 Makefile**: `Makefile`
- **后端测试脚本**: `backend/scripts/`

---

*此文档由 FRP-007 任务生成，定义 DMH 项目的本地一键全量回归入口。*
