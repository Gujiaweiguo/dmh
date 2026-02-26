# T17: KPI 看板与治理节奏

## 1. 概述

本文档定义 DMH 项目 CI/CD 质量看板的核心 KPI 指标、采集方法、周报模板和治理节奏。

---

## 2. 核心 KPI 指标定义

### 2.1 KPI 总览

| # | 指标名称 | 英文名 | 目标值 | 采集频率 | 数据源 |
|---|---------|--------|--------|---------|--------|
| 1 | PR 时延 | PR Gate Duration | ≤15 min | 每次 PR | GitHub Actions |
| 2 | Nightly 通过率 | Nightly Pass Rate | ≥95% | 每日 | GitHub Actions |
| 3 | SKIP 率 | Skip Rate | ≤10% | 每次 PR/Merge | 测试日志 |
| 4 | Flaky 率 | Flaky Rate | <2% | 每周 | 多次运行对比 |
| 5 | 覆盖率 | Coverage | Backend≥78%, Admin≥80%, H5≥70% | 每次 PR | Coverage Report |

---

## 3. KPI 详细定义与计算公式

### 3.1 PR 时延 (PR Gate Duration)

**定义**：从 PR 创建/更新到 PR Gate 全部完成的时间。

**计算公式**：
```
PR 时延 = max(backend-unit 耗时, frontend-unit 耗时, lint 耗时)
```

**数据源**：
- GitHub Actions API: `GET /repos/{owner}/{repo}/actions/runs`
- 字段: `created_at` → `updated_at`

**采集脚本**：
```bash
#!/bin/bash
# scripts/collect_pr_duration.sh

REPO="Gujiaweiguo/dmh"
TOKEN="${GITHUB_TOKEN}"

# 获取最近 7 天的 PR Gate 运行记录
curl -s -H "Authorization: token $TOKEN" \
  "https://api.github.com/repos/$REPO/actions/workflows/pr-gate.yml/runs?per_page=50" \
  | jq -r '.workflow_runs[] | select(.status == "completed") | 
      {id, duration: ((.updated_at | fromdateiso8601) - (.created_at | fromdateiso8601))}'
```

**阈值与告警**：

| 级别 | 阈值 | 告警动作 |
|------|------|---------|
| 🟢 正常 | ≤15 min | 无 |
| 🟡 警告 | 15-20 min | 周报记录 |
| 🔴 异常 | >20 min | 立即通知，启动排查 |

**告警条件**：
- 连续 3 次 PR 时延 >15 min
- 单次 PR 时延 >25 min

---

### 3.2 Nightly 通过率 (Nightly Pass Rate)

**定义**：Nightly 全量回归执行成功的比例。

**计算公式**：
```
Nightly 通过率 = (成功次数 / 总执行次数) × 100%
```

**数据源**：
- GitHub Actions: `full-regression.yml` 运行结果
- 子 workflow 结果汇总

**采集脚本**：
```bash
#!/bin/bash
# scripts/collect_nightly_rate.sh

REPO="Gujiaweiguo/dmh"
TOKEN="${GITHUB_TOKEN}"
DAYS=7

# 获取过去 7 天的 Nightly 运行记录
curl -s -H "Authorization: token $TOKEN" \
  "https://api.github.com/repos/$REPO/actions/workflows/full-regression.yml/runs?per_page=50" \
  | jq --arg days "$DAYS" '
    [.workflow_runs[] | select(.status == "completed")] |
    {
      total: length,
      success: [.[] | select(.conclusion == "success")] | length,
      failure: [.[] | select(.conclusion == "failure")] | length
    } |
    .rate = ((.success / .total) * 100 | floor)
  '
```

**阈值与告警**：

| 级别 | 阈值 | 告警动作 |
|------|------|---------|
| 🟢 正常 | ≥95% | 无 |
| 🟡 警告 | 80-95% | 周报记录 |
| 🔴 异常 | <80% | 立即通知，阻断发布 |

**告警条件**：
- 连续 2 次Nightly 失败
- 单周通过率 <80%

---

### 3.3 SKIP 率 (Skip Rate)

**定义**：集成测试中被跳过的测试比例（反映环境稳定性）。

**计算公式**：
```
SKIP 率 = (SKIP 测试数 / 总测试数) × 100%
执行率 = 100% - SKIP 率
```

**数据源**：
- 测试日志中的 `--- SKIP` 统计
- SKIP 原因标签分析

**采集脚本**：
```bash
#!/bin/bash
# scripts/collect_skip_rate.sh

cd backend

# 运行集成测试并分析 SKIP
OUTPUT=$(go test ./test/integration/... -v 2>&1)

TOTAL=$(echo "$OUTPUT" | grep -cE "^--- (PASS|FAIL|SKIP)" || echo 0)
SKIPPED=$(echo "$OUTPUT" | grep -c "^--- SKIP" || echo 0)

if [ "$TOTAL" -gt 0 ]; then
  SKIP_RATE=$(echo "scale=2; $SKIPPED * 100 / $TOTAL" | bc)
  EXEC_RATE=$(echo "scale=2; 100 - $SKIP_RATE" | bc)
  
  echo "Total Tests: $TOTAL"
  echo "Skipped: $SKIPPED"
  echo "Skip Rate: ${SKIP_RATE}%"
  echo "Execution Rate: ${EXEC_RATE}%"
fi

# SKIP 原因分布
echo ""
echo "=== SKIP Reason Distribution ==="
for reason in API_UNAVAILABLE MYSQL_UNAVAILABLE REDIS_UNAVAILABLE LOGIN_FAILED DATA_PREP_FAILED; do
  count=$(echo "$OUTPUT" | grep -c "SKIP_REASON:$reason" || echo 0)
  echo "$reason: $count"
done
```

**阈值与告警**：

| 级别 | 阈值 | 告警动作 |
|------|------|---------|
| 🟢 正常 | ≤10% | 无 |
| 🟡 警告 | 10-20% | 周报记录，检查环境 |
| 🔴 异常 | >20% | 立即修复环境问题 |

**告警条件**：
- SKIP 率连续 3 次超过 10%
- API_UNAVAILABLE 原因占比 >50%

---

### 3.4 Flaky 率 (Flaky Rate)

**定义**：不稳定测试的比例（同一测试多次运行结果不一致）。

**计算公式**：
```
Flaky 率 = (不稳定测试数 / 总测试数) × 100%

不稳定测试 = 连续 3 次运行中，失败次数不一致的测试
```

**数据源**：
- 多次测试运行对比
- 测试历史记录

**采集脚本**：
```bash
#!/bin/bash
# scripts/collect_flaky_rate.sh

cd backend
REPORT_DIR="reports/flaky"
mkdir -p "$REPORT_DIR"

echo "Running tests 3 times to detect flaky behavior..."

# 运行 3 次并记录失败测试
declare -A FAILURES
for i in 1 2 3; do
  echo "=== Run $i/3 ==="
  go test ./test/integration/... -v -count=1 2>&1 | tee "$REPORT_DIR/run_$i.log"
  
  # 提取失败测试名
  grep "^--- FAIL" "$REPORT_DIR/run_$i.log" | awk '{print $3}' > "$REPORT_DIR/failures_$i.txt"
done

# 对比三次运行结果
echo ""
echo "=== Flaky Analysis ==="

# 找出失败次数不一致的测试
comm -3 <(sort "$REPORT_DIR/failures_1.txt") <(sort "$REPORT_DIR/failures_2.txt") > "$REPORT_DIR/flaky_diff.txt"
comm -3 <(sort "$REPORT_DIR/failures_2.txt") <(sort "$REPORT_DIR/failures_3.txt") >> "$REPORT_DIR/flaky_diff.txt"

FLAKY_TESTS=$(sort "$REPORT_DIR/flaky_diff.txt" | uniq | wc -l)
TOTAL_TESTS=$(go test ./test/integration/... -list ".*" 2>/dev/null | wc -l)

if [ "$TOTAL_TESTS" -gt 0 ]; then
  FLAKY_RATE=$(echo "scale=2; $FLAKY_TESTS * 100 / $TOTAL_TESTS" | bc)
  echo "Flaky Tests: $FLAKY_TESTS"
  echo "Total Tests: $TOTAL_TESTS"
  echo "Flaky Rate: ${FLAKY_RATE}%"
fi
```

**阈值与告警**：

| 级别 | 阈值 | 告警动作 |
|------|------|---------|
| 🟢 正常 | <2% | 无 |
| 🟡 警告 | 2-5% | 周报记录，创建修复任务 |
| 🔴 异常 | >5% | 暂停合并，集中修复 |

**告警条件**：
- Flaky 率超过 2%
- 发现新的 Flaky 测试

---

### 3.5 覆盖率 (Coverage)

**定义**：代码被测试覆盖的比例。

**计算公式**：
```
覆盖率 = (被覆盖的代码行数 / 总代码行数) × 100%
```

**数据源**：
- PR Gate 中的覆盖率输出
- Coverage Gate 报告

**阈值（已固化在 CI）**：

| 模块 | 阈值 | 当前值 | 状态 |
|------|------|--------|------|
| Backend | ≥78% | ~78% | ✅ |
| Frontend-Admin | ≥80% | 83.65% | ✅ |
| Frontend-H5 | ≥70% | 87.37% | ✅ |

**告警条件**：
- 覆盖率低于阈值（CI 自动阻断）
- 覆盖率环比下降 >2%

---

## 4. 数据源缺口列表

当前部分指标需要额外配置才能自动采集：

| 指标 | 数据源状态 | 缺口说明 | 解决方案 |
|------|-----------|---------|---------|
| PR 时延 | ⚠️ 部分可用 | 需 GitHub API Token | 配置 `GITHUB_TOKEN` 环境变量 |
| Nightly 通过率 | ⚠️ 部分可用 | 需 GitHub API Token | 配置 `GITHUB_TOKEN` 环境变量 |
| SKIP 率 | ✅ 可用 | 脚本已就绪 | 运行采集脚本 |
| Flaky 率 | ⚠️ 需配置 | 需多次运行 | CI 添加定期 Flaky 检测 |
| 覆盖率 | ✅ 可用 | CI 已集成 | 从 CI 日志提取 |

---

## 5. 周报模板

### 5.1 周报格式

```markdown
# DMH CI/CD 质量周报

**报告周期**: YYYY-MM-DD ~ YYYY-MM-DD
**报告人**: [姓名]

---

## 1. KPI 摘要

| 指标 | 本周值 | 上周值 | 环比 | 目标 | 状态 |
|------|--------|--------|------|------|------|
| PR 时延 | XX min | XX min | ↑/↓ | ≤15 min | 🟢/🟡/🔴 |
| Nightly 通过率 | XX% | XX% | ↑/↓ | ≥95% | 🟢/🟡/🔴 |
| SKIP 率 | XX% | XX% | ↑/↓ | ≤10% | 🟢/🟡/🔴 |
| Flaky 率 | XX% | XX% | ↑/↓ | <2% | 🟢/🟡/🔴 |
| Backend 覆盖率 | XX% | XX% | - | ≥78% | 🟢/🟡/🔴 |
| Admin 覆盖率 | XX% | XX% | - | ≥80% | 🟢/🟡/🔴 |
| H5 覆盖率 | XX% | XX% | - | ≥70% | 🟢/🟡/🔴 |

---

## 2. SKIP 原因分析

| 原因 | 数量 | 占比 | 趋势 |
|------|------|------|------|
| API_UNAVAILABLE | X | XX% | ↑/↓ |
| MYSQL_UNAVAILABLE | X | XX% | ↑/↓ |
| LOGIN_FAILED | X | XX% | ↑/↓ |
| DATA_PREP_FAILED | X | XX% | ↑/↓ |

**分析**:
- [主要问题描述]
- [根因分析]

---

## 3. Flaky 测试清单

| 测试名 | 文件 | 出现次数 | 状态 |
|--------|------|---------|------|
| TestXXX | xxx_test.go | X | 待修复/修复中/已修复 |

---

## 4. 本周亮点

- [积极变化 1]
- [积极变化 2]

---

## 5. 本周问题

- [问题 1]: [描述] - [负责人] - [预计解决时间]
- [问题 2]: [描述] - [负责人] - [预计解决时间]

---

## 6. 下周计划

- [ ] [任务 1]
- [ ] [任务 2]

---

*生成时间: YYYY-MM-DD HH:MM*
```

### 5.2 周报生成脚本

```bash
#!/bin/bash
# scripts/generate_weekly_report.sh

WEEK_START=$(date -d "last monday" +%Y-%m-%d)
WEEK_END=$(date +%Y-%m-%d)
REPORT_FILE="reports/weekly/weekly-report-${WEEK_END}.md"

mkdir -p reports/weekly

cat > "$REPORT_FILE" << EOF
# DMH CI/CD 质量周报

**报告周期**: ${WEEK_START} ~ ${WEEK_END}
**报告人**: [自动生成]

---

## 1. KPI 摘要

EOF

# 采集各项指标并写入报告
echo "Collecting KPI data..."

# PR 时延
echo "| PR 时延 | $(scripts/collect_pr_duration.sh | jq -r '.avg_duration') min | - | - | ≤15 min | - |" >> "$REPORT_FILE"

# Nightly 通过率
echo "| Nightly 通过率 | $(scripts/collect_nightly_rate.sh | jq -r '.rate')% | - | - | ≥95% | - |" >> "$REPORT_FILE"

# SKIP 率
SKIP_DATA=$(scripts/collect_skip_rate.sh)
echo "| SKIP 率 | $(echo "$SKIP_DATA" | grep "Skip Rate" | awk '{print $3}') | - | - | ≤10% | - |" >> "$REPORT_FILE"

echo ""
echo "Report generated: $REPORT_FILE"
```

---

## 6. 治理节奏

### 6.1 治理会议日程

| 会议 | 频率 | 时间 | 参与人 | 议程 |
|------|------|------|--------|------|
| **周度质量评审** | 每周一 | 10:00-10:30 | 全体开发 | KPI 回顾、问题分配 |
| **月度质量复盘** | 每月第一周 | 14:00-15:00 | Tech Lead + QA | 趋势分析、阈值调整 |

### 6.2 周度质量评审议程

```
1. KPI 快速过目 (5 min)
   - PR 时延、Nightly 通过率、SKIP 率、Flaky 率
   - 重点关注异常指标

2. SKIP 原因分析 (5 min)
   - 本周 SKIP 分布
   - 环境问题 vs 代码问题

3. Flaky 测试进展 (5 min)
   - 新发现的 Flaky 测试
   - 修复进展更新

4. 问题分配 (10 min)
   - 本周新增问题认领
   - 遗留问题进度

5. 下周计划 (5 min)
   - 优先级确认
   - 风险预警
```

### 6.3 月度质量复盘议程

```
1. 趋势分析 (15 min)
   - KPI 30 天趋势图
   - 覆盖率变化趋势

2. 阈值评估 (10 min)
   - 当前阈值是否合理
   - 是否需要调整

3. 流程改进 (15 min)
   - CI/CD 流程优化建议
   - 测试策略调整

4. 技术债务 (10 min)
   - Flaky 测试清理进度
   - 测试覆盖率提升计划

5. 下月目标 (10 min)
   - KPI 目标设定
   - 重点工作分配
```

---

## 7. 告警配置

### 7.1 告警渠道

| 渠道 | 用途 | 配置位置 |
|------|------|---------|
| GitHub Actions 通知 | CI 失败通知 | 默认启用 |
| 邮件 | Nightly 失败通知 | GitHub Settings |
| Slack/钉钉 (可选) | 紧急告警 | Webhook 配置 |

### 7.2 告警规则

```yaml
# .github/workflows/kpi-alert.yml (建议新增)
name: KPI Alert

on:
  workflow_run:
    workflows: ["PR Gate", "Full Regression"]
    types: [completed]

jobs:
  alert:
    if: failure()
    runs-on: ubuntu-latest
    steps:
      - name: Send alert
        run: |
          # 可选：发送到 Slack/钉钉
          curl -X POST "${{ secrets.ALERT_WEBHOOK }}" \
            -H 'Content-Type: application/json' \
            -d '{
              "text": "⚠️ CI Alert: ${{ github.workflow }} failed",
              "run_url": "${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}"
            }'
```

### 7.3 告警升级规则

| 场景 | 初次告警 | 15 分钟未响应 | 1 小时未响应 |
|------|---------|--------------|-------------|
| Nightly 失败 | @on-call | @tech-lead | @team |
| PR Gate 超时 | @author | @reviewer | @tech-lead |
| Flaky 检测 | @test-owner | @tech-lead | 创建 Issue |

---

## 8. 实施清单

### 8.1 立即实施

- [x] 定义 KPI 指标和计算公式
- [x] 设计周报模板
- [x] 定义治理会议议程
- [ ] 配置 GitHub API Token
- [ ] 部署采集脚本到 CI

### 8.2 短期实施（1-2 周）

- [ ] 配置 Slack/钉钉告警（可选）
- [ ] 添加 Flaky 定期检测 Job
- [ ] 生成首份周报

### 8.3 长期目标

- [ ] KPI 趋势图表可视化
- [ ] 自动化周报生成
- [ ] 与 Jira/飞书集成

---

## 9. 附录

### 9.1 KPI 目标速查表

| 指标 | 目标 | 警告线 | 异常线 |
|------|------|--------|--------|
| PR 时延 | ≤15 min | 15-20 min | >20 min |
| Nightly 通过率 | ≥95% | 80-95% | <80% |
| SKIP 率 | ≤10% | 10-20% | >20% |
| Flaky 率 | <2% | 2-5% | >5% |
| Backend 覆盖率 | ≥78% | 75-78% | <75% |
| Admin 覆盖率 | ≥80% | 75-80% | <75% |
| H5 覆盖率 | ≥70% | 65-70% | <65% |

### 9.2 采集脚本位置

```
backend/scripts/
├── collect_pr_duration.sh      # PR 时延采集
├── collect_nightly_rate.sh     # Nightly 通过率采集
├── collect_skip_rate.sh        # SKIP 率采集
├── collect_flaky_rate.sh       # Flaky 率采集
├── analyze_skip_reasons.sh     # SKIP 原因分析
└── generate_weekly_report.sh   # 周报生成
```

---

**创建时间**: 2026-02-26
**任务**: T17 - KPI 看板与治理节奏
**依赖**: T2, T10, T12, T13
