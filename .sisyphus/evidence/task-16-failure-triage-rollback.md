# T16: 失败分诊与回滚机制

## 执行时间
- 开始: 2026-02-26
- 结束: 2026-02-26

---

## 1. 概述

本文档定义 DMH 项目 CI/CD 失败分诊与回滚机制，包括：
- 失败分类标准（代码/环境/数据/外部依赖）
- 分诊决策树与标准处理路径
- 回滚触发条件与步骤
- 重试策略与限制

---

## 2. 失败分类标准

### 2.1 四类失败定义

```
┌─────────────────────────────────────────────────────────────────────┐
│                        失败分类金字塔                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Level 1: 环境问题 (40%) - 最常见                                     │
│  ├── API 服务未启动/未就绪                                            │
│  ├── MySQL 连接失败                                                  │
│  ├── Redis 不可用                                                    │
│  └── CI 容器资源配置问题                                              │
│                                                                      │
│  Level 2: 代码缺陷 (30%) - 需修复                                    │
│  ├── 业务逻辑错误                                                    │
│  ├── 边界条件问题                                                    │
│  ├── 资源泄漏（defer in loop）                                       │
│  └── 并发安全问题                                                    │
│                                                                      │
│  Level 3: 数据问题 (20%) - 需隔离                                    │
│  ├── 测试数据冲突                                                    │
│  ├── 数据准备失败                                                    │
│  ├── 数据库状态不一致                                                 │
│  └── 唯一约束违反                                                    │
│                                                                      │
│  Level 4: 外部依赖 (10%) - 可重试                                    │
│  ├── 网络抖动                                                        │
│  ├── 第三方 API 超时                                                 │
│  ├── DNS 解析失败                                                    │
│  └── 临时服务不可用                                                   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 2.2 分类判断标准

| 分类 | 关键词/错误模式 | SKIP 标签 | 处理优先级 |
|------|----------------|-----------|-----------|
| **环境问题** | `connection refused`, `timeout`, `unavailable`, `not ready` | `*_UNAVAILABLE` | P1: 修复环境 |
| **代码缺陷** | `assertion failed`, `panic`, `nil pointer`, `expected X got Y` | 无 | P0: 立即修复 |
| **数据问题** | `duplicate entry`, `foreign key constraint`, `data not found` | `DATA_PREP_FAILED` | P1: 数据隔离 |
| **外部依赖** | `5xx`, `network error`, `context deadline exceeded` | 无 | P2: 可重试 |

---

## 3. 分诊决策树

### 3.1 决策流程图

```
测试失败
    │
    ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Step 1: 检查 SKIP 标签                                              │
└─────────────────────────────────────────────────────────────────────┘
    │
    ├── 有 SKIP_REASON 标签?
    │   ├── API_UNAVAILABLE ────────────► 环境问题: 检查 API 服务
    │   ├── MYSQL_UNAVAILABLE ──────────► 环境问题: 检查 MySQL 容器
    │   ├── REDIS_UNAVAILABLE ──────────► 环境问题: 检查 Redis（可选）
    │   ├── LOGIN_FAILED ───────────────► 环境问题: 检查认证配置
    │   └── DATA_PREP_FAILED ───────────► 数据问题: 检查数据准备
    │
    └── 无 SKIP 标签
        │
        ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Step 2: 检查错误类型                                                │
└─────────────────────────────────────────────────────────────────────┘
        │
        ├── connection refused / timeout ──► 环境问题
        │
        ├── assertion failed / panic ──────► 代码缺陷
        │
        ├── duplicate entry / constraint ──► 数据问题
        │
        └── 5xx / network error ───────────► 外部依赖（可重试）
```

### 3.2 详细分诊矩阵

| 错误信息模式 | 分类 | 处理路径 | 时限 |
|-------------|------|----------|------|
| `dial tcp: connection refused` | 环境 | 检查服务启动状态 | 15 min |
| `API not ready within timeout` | 环境 | 检查 API 就绪探针 | 10 min |
| `Error 1045: Access denied` | 环境 | 检查数据库凭证 | 5 min |
| `panic: nil pointer dereference` | 代码 | 创建 Bug Issue | 立即 |
| `assertion failed: expected X, got Y` | 代码 | 创建 Bug Issue | 立即 |
| `Error 1062: Duplicate entry` | 数据 | 数据隔离/工厂模式 | 30 min |
| `Error 1452: foreign key constraint` | 数据 | 检查依赖关系 | 30 min |
| `context deadline exceeded` | 外部 | 检查超时配置/重试 | 10 min |
| `http: 502/503/504` | 外部 | 重试或等待恢复 | 15 min |

---

## 4. 标准处理路径

### 4.1 环境问题处理路径

```
环境问题检测
    │
    ├── CI 环境
    │   ├── 检查服务容器状态
    │   ├── 检查健康检查配置
    │   ├── 检查端口映射
    │   └── 重启失败服务
    │
    ├── 本地环境
    │   ├── 运行 docker compose ps
    │   ├── 检查端口占用
    │   ├── 检查环境变量
    │   └── 使用 quick-start.sh 重建
    │
    └── 通用检查
        ├── API: curl ${DMH_INTEGRATION_BASE_URL}/api/v1/auth/login
        ├── MySQL: mysqladmin ping -h ${MYSQL_TEST_HOST}
        └── Redis: redis-cli -h ${REDIS_TEST_HOST} ping
```

**处理时限**: 15 分钟内定位，30 分钟内修复

### 4.2 代码缺陷处理路径

```
代码缺陷检测
    │
    ├── PR Gate 失败
    │   ├── 阻断 PR 合并
    │   ├── 通知 PR 作者
    │   ├── 创建 Bug Issue（可选）
    │   └── 等待修复后重新提交
    │
    ├── Merge Gate 失败
    │   ├── 阻断合并
    │   ├── 通知分支 Owner
    │   ├── 创建 Hotfix 分支
    │   └── 修复后重新触发
    │
    └── Nightly 失败
        ├── 创建 Bug Issue
        ├── 标记优先级（P0/P1）
        ├── 分配责任团队
        └── 下次 Nightly 前修复
```

**处理时限**: 
- PR Gate: 阻断直到修复
- Nightly: 24 小时内修复

### 4.3 数据问题处理路径

```
数据问题检测
    │
    ├── 唯一约束冲突
    │   ├── 使用工厂模式随机化
    │   ├── 添加唯一后缀（时间戳/UUID）
    │   └── 检查测试数据清理
    │
    ├── 外键约束失败
    │   ├── 检查依赖数据创建顺序
    │   ├── 使用事务回滚
    │   └── 使用独立测试数据库
    │
    └── 数据不一致
        ├── TRUNCATE 重置
        ├── 使用 testutil.SetupMySQLTestDB
        └── 确保测试隔离（-p 1）
```

**处理时限**: 30 分钟内定位，1 小时内修复

### 4.4 外部依赖处理路径

```
外部依赖失败
    │
    ├── 可重试错误（5xx, timeout）
    │   ├── 检查幂等性
    │   ├── 执行重试（最多 3 次）
    │   ├── 指数退避等待
    │   └── 记录重试日志
    │
    ├── 不可重试错误（4xx）
    │   ├── 立即报告失败
    │   ├── 不执行重试
    │   └── 检查请求参数
    │
    └── 持续失败
        ├── 检查服务状态
        ├── 联系服务 Owner
        └── 考虑 Mock 替代
```

**处理时限**: 10 分钟内判断，15 分钟内决定下一步

---

## 5. 回滚机制

### 5.1 回滚触发条件

| 触发场景 | 触发条件 | 回滚范围 | 决策者 |
|---------|---------|---------|--------|
| **PR Gate 失败** | 单元测试/覆盖率/Lint 失败 | 不合并 PR | 自动 |
| **Merge Gate 失败** | 集成测试/E2E 失败 | 不合并到 main | 自动 |
| **生产部署后失败** | 健康检查失败/错误率飙升 | 回滚到上一版本 | On-call |
| **Nightly 失败** | 关键路径测试失败 | 不阻断，但创建 Issue | 自动 |
| **手动触发** | 发现严重缺陷 | 根据情况 | Tech Lead |

### 5.2 回滚步骤

#### 5.2.1 PR/Merge Gate 回滚（自动）

```bash
# 自动阻断，无需手动操作
# PR 状态变为 "Changes Requested"
# 合并按钮被禁用
```

#### 5.2.2 生产环境回滚（手动）

```bash
# Step 1: 确认回滚决策
# - 错误率 > 5% 或
# - 关键功能不可用 或
# - 数据完整性风险

# Step 2: 执行回滚
cd deploy
docker compose -f docker-compose-simple.yml down

# Step 3: 恢复上一版本
git checkout HEAD~1 -- backend/
cd backend && go build -o dmh-api api/dmh.go
cp dmh-api ../deploy/dmh-api

# Step 4: 重新部署
cd ../deploy
docker compose -f docker-compose-simple.yml up -d

# Step 5: 验证回滚
curl -s http://localhost:8889/api/v1/health
```

#### 5.2.3 数据库回滚（谨慎）

```bash
# 仅在数据损坏时执行
# 必须有备份

# Step 1: 停止应用
docker stop dmh-api

# Step 2: 恢复数据库备份
mysql -h mysql8 -u root -p dmh < /backup/dmh_YYYYMMDD.sql

# Step 3: 验证数据
mysql -h mysql8 -u root -p dmh -e "SELECT COUNT(*) FROM users;"

# Step 4: 重启应用
docker start dmh-api
```

### 5.3 回滚时限要求

| 回滚类型 | 决策时限 | 执行时限 | 验证时限 |
|---------|---------|---------|---------|
| PR 阻断 | 自动 | N/A | N/A |
| 生产回滚 | 5 min | 10 min | 5 min |
| 数据库回滚 | 15 min | 30 min | 10 min |

---

## 6. 重试策略

### 6.1 幂等性判断矩阵（继承 T12）

| 操作类型 | 幂等性 | 是否可重试 | 最大重试 | 退避策略 |
|---------|--------|-----------|---------|---------|
| HTTP GET | ✅ | ✅ | 3 | 指数退避 |
| HTTP PUT | ✅ | ✅ | 3 | 指数退避 |
| HTTP DELETE | ✅ | ✅ | 3 | 指数退避 |
| HTTP POST (创建) | ❌ | ❌ | 0 | N/A |
| HTTP POST (搜索) | ✅ | ✅ | 3 | 指数退避 |
| 登录请求 | ✅ | ✅ | 3 | 指数退避 |
| 数据库 TRUNCATE | ✅ | ✅ | 3 | 指数退避 |
| 数据库 INSERT | ❌ | ❌ | 0 | N/A |
| 数据库 UPDATE | ⚠️ | ⚠️ | 1 | 固定 1s |

### 6.2 重试配置

```yaml
# 重试配置标准
retry:
  max_attempts: 3
  initial_wait_ms: 100
  max_wait_ms: 2000
  multiplier: 2.0
  
# 退避序列: 100ms → 200ms → 400ms → 800ms → 1600ms (cap at 2000ms)
```

### 6.3 可重试的 HTTP 状态码

```go
// 可重试的状态码
var retryableStatusCodes = []int{
    429, // Too Many Requests
    500, // Internal Server Error
    502, // Bad Gateway
    503, // Service Unavailable
    504, // Gateway Timeout
}

// 不可重试的状态码
var nonRetryableStatusCodes = []int{
    400, // Bad Request
    401, // Unauthorized
    403, // Forbidden
    404, // Not Found
    409, // Conflict
    422, // Unprocessable Entity
}
```

### 6.4 重试日志规范

```
[RETRY] attempt=1/3 operation=GET /api/v1/users error="connection refused" wait=100ms
[RETRY] attempt=2/3 operation=GET /api/v1/users error="timeout" wait=200ms
[RETRY] attempt=3/3 operation=GET /api/v1/users error="timeout" wait=400ms
[RETRY] FAILED after 3 attempts: GET /api/v1/users - timeout
```

---

## 7. 分诊缺项阻断规则

### 7.1 必填分诊信息

每次失败必须记录以下信息，缺项则阻断发布：

| 字段 | 说明 | 示例 |
|------|------|------|
| `failure_id` | 失败唯一标识 | `FAIL-20260226-001` |
| `category` | 失败分类 | `environment/code/data/external` |
| `severity` | 严重程度 | `blocker/critical/major/minor` |
| `workflow` | 触发的 CI 流程 | `pr-gate/merge-gate/nightly` |
| `job_name` | 失败的 Job 名称 | `backend-unit/system-test-gate` |
| `error_pattern` | 错误模式关键词 | `connection refused` |
| `triage_owner` | 分诊责任人 | `@username` |
| `action_taken` | 采取的行动 | `fix/defer/retry/ignore` |

### 7.2 阻断发布条件

以下情况**必须**阻断发布：

1. ❌ 分诊信息不完整
2. ❌ 分类为 "未知"
3. ❌ 严重程度为 `blocker` 且未修复
4. ❌ 代码缺陷未创建 Issue
5. ❌ 环境问题未确认修复方案

### 7.3 分诊完成标准

```
✅ 分类明确（非"未知"）
✅ 责任人已分配
✅ 行动已记录
✅ Issue 已创建（如适用）
✅ 修复方案已确认（如适用）
```

---

## 8. 失败处理时限总结

| 场景 | 分类 | 定位时限 | 修复时限 | 阻断发布 |
|------|------|---------|---------|---------|
| PR Gate | 代码缺陷 | 立即 | 阻断直到修复 | ✅ |
| PR Gate | 环境问题 | 15 min | 30 min | ❌ |
| Merge Gate | 任意 | 15 min | 1 hour | ✅ |
| Nightly | 代码缺陷 | 30 min | 24 hour | ❌ |
| Nightly | 环境问题 | 30 min | 4 hour | ❌ |
| 生产部署 | 任意 | 5 min | 15 min | ✅ |

---

## 9. CI 集成建议

### 9.1 自动分诊脚本

```bash
#!/bin/bash
# backend/scripts/triage_failure.sh
# 自动分析测试失败并生成分诊报告

FAILURE_LOG="$1"
REPORT_FILE="triage_report.md"

# 检测环境问题
if grep -qE "connection refused|timeout|unavailable" "$FAILURE_LOG"; then
    CATEGORY="environment"
elif grep -qE "assertion failed|panic|nil pointer" "$FAILURE_LOG"; then
    CATEGORY="code"
elif grep -qE "duplicate entry|foreign key|constraint" "$FAILURE_LOG"; then
    CATEGORY="data"
elif grep -qE "5[0-9]{2}|network error|deadline exceeded" "$FAILURE_LOG"; then
    CATEGORY="external"
else
    CATEGORY="unknown"
fi

# 生成分诊报告
cat > "$REPORT_FILE" << EOF
# Failure Triage Report

**Failure ID**: FAIL-$(date +%Y%m%d)-$(uuidgen | cut -d'-' -f1)
**Category**: $CATEGORY
**Workflow**: ${GITHUB_WORKFLOW:-local}
**Job**: ${GITHUB_JOB:-unknown}
**Timestamp**: $(date -Iseconds)

## Error Pattern
\`\`\`
$(grep -E "FAIL|ERROR|panic" "$FAILURE_LOG" | head -5)
\`\`\`

## Recommended Action
$(case $CATEGORY in
    environment) echo "检查服务容器状态、网络配置、环境变量" ;;
    code) echo "创建 Bug Issue，分配给相关开发者" ;;
    data) echo "检查数据隔离、工厂模式使用" ;;
    external) echo "检查外部服务状态，考虑重试" ;;
    *) echo "需要人工分诊" ;;
esac)

## Triage Checklist
- [ ] 分类确认
- [ ] 责任人分配
- [ ] Issue 创建（如适用）
- [ ] 修复方案确认
EOF

echo "Triage report generated: $REPORT_FILE"
echo "Category: $CATEGORY"
```

### 9.2 GitHub Actions 集成

```yaml
# 在失败时自动生成分诊报告
- name: Triage Failure
  if: failure()
  run: |
    chmod +x backend/scripts/triage_failure.sh
    ./backend/scripts/triage_failure.sh /tmp/test-output.log
    
- name: Upload Triage Report
  if: failure()
  uses: actions/upload-artifact@v4
  with:
    name: triage-report
    path: triage_report.md
```

---

## 10. 变更记录

| 日期 | 版本 | 作者 | 变更内容 |
|------|------|------|----------|
| 2026-02-26 | v1.0 | Sisyphus | 初始版本，定义失败分诊与回滚机制 |

---

## 附录 A: SKIP 标签速查表

| 标签 | 分类 | 处理路径 |
|------|------|----------|
| `API_UNAVAILABLE` | 环境 | 检查 API 服务 |
| `MYSQL_UNAVAILABLE` | 环境 | 检查 MySQL 容器 |
| `REDIS_UNAVAILABLE` | 环境 | 检查 Redis（可选） |
| `LOGIN_FAILED` | 环境 | 检查认证配置 |
| `DATA_PREP_FAILED` | 数据 | 检查数据准备 |

## 附录 B: 错误模式正则表达式

```regex
# 环境问题
connection refused|timeout|unavailable|not ready|deadline exceeded

# 代码缺陷
assertion failed|panic|nil pointer|expected.*got|index out of range

# 数据问题
duplicate entry|foreign key constraint|data not found|constraint fails

# 外部依赖
5[0-9]{2}|network error|temporary failure|service unavailable
```
