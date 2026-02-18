# DMH 证据保留策略

> 版本: Plan v1 (FRP-005)
> 变更: add-full-regression-testing
> 生成日期: 2026-02-18

---

## 1. 证据保留概述

### 1.1 目的

本策略定义 DMH 项目测试执行证据的保留规则，确保：

- 回归测试结果可追溯、可审计
- 失败证据支持发布后复盘与根因分析
- 证据检索支持多维度查询
- 与审计周期和发布流程对齐

### 1.2 适用范围

本策略适用于以下测试活动的证据保留：

| 测试活动 | 触发场景 | 归档位置 |
|----------|----------|----------|
| 全量回归 | Release candidate / 主干受保护分支变更 | GitHub Actions Artifacts |
| PR 验证 | Pull Request / Push 事件 | GitHub Actions Artifacts |
| 订单专项回归 | CI 工作流触发 | GitHub Actions Artifacts |
| OpenSpec 校验 | CI 工作流触发 | GitHub Actions Artifacts |

### 1.3 引用规格

- `openspec/changes/add-full-regression-testing/specs/system-test-execution/spec.md` - Enforce evidence retention window
- `docs/testing/execution/FULL_REGRESSION_DEFINITION.md` (FRP-002) - 全量回归口径定义

---

## 2. 证据类型

### 2.1 证据分类

| 证据类型 | 说明 | 来源 | 格式 |
|----------|------|------|------|
| **测试日志** | 测试执行的标准输出/错误输出 | CI workflow | `.log`, `.txt` |
| **覆盖率报告** | 代码覆盖率详情与汇总 | coverage-gate.yml | `.xml` (Cobertura), `.html` |
| **截图/录屏** | E2E 测试失败时的视觉证据 | Playwright artifacts | `.png`, `.webm` |
| **请求/响应报文** | API 测试的网络记录 | 集成测试日志 | `.json`, `.har` |
| **OpenSpec 校验结果** | 规格校验输出 | system-test-gate.yml | `.txt`, `.json` |
| **阻断报告** | 发布被阻断时的失败明细 | 自动生成 | `.md`, `.json` |
| **测试结果汇总** | JUnit 格式测试报告 | 各测试框架 | `.xml` (JUnit) |

### 2.2 证据与测试层级映射

```
证据包结构:
├── backend/
│   ├── unit-tests/
│   │   ├── test-output.log
│   │   ├── coverage.out
│   │   └── junit-report.xml
│   ├── integration-tests/
│   │   ├── test-output.log
│   │   └── network-traces.json
│   └── order-regression/
│       ├── test-output.log
│       └── api-responses.json
│
├── frontend-admin/
│   ├── unit-tests/
│   │   ├── test-output.log
│   │   └── coverage/
│   │       └── lcov.info
│   └── e2e-tests/
│       ├── test-output.log
│       ├── screenshots/
│       │   └── failure-*.png
│       └── videos/
│           └── test-*.webm
│
├── frontend-h5/
│   ├── unit-tests/
│   │   ├── test-output.log
│   │   └── coverage/
│   │       └── lcov.info
│   └── e2e-tests/
│       ├── test-output.log
│       ├── screenshots/
│       │   └── failure-*.png
│       └── videos/
│           └── test-*.webm
│
└── openspec/
    └── validation-output.log
```

---

## 3. 保留时长规则

### 3.1 分级保留策略

基于回归结论的保留时长：

| 回归结论 | 保留时长 | 最小保留 | 说明 |
|----------|----------|----------|------|
| **PASS** | 90 天 | >= 90 天 | 成功回归证据，覆盖季度审计周期 |
| **FAIL** | 180 天 | >= 180 天 | 失败回归证据，支持发布后复盘 |
| **INCONCLUSIVE** | 180 天 | >= 180 天 | 需人工裁定的证据，保留更长时间 |

### 3.2 保留时长依据

| 来源 | 要求 | 本策略实现 |
|------|------|------------|
| Delta Spec: Enforce evidence retention window | 成功 >= 90 天 | PASS: 90 天 |
| Delta Spec: Enforce evidence retention window | 失败 >= 180 天 | FAIL/INCONCLUSIVE: 180 天 |
| 审计周期 | 季度审计 | 90 天覆盖一个完整季度 |
| 发布复盘 | 发布后 6 个月内 | 180 天覆盖发布后复盘窗口 |

### 3.3 特殊场景

| 场景 | 保留时长 | 说明 |
|------|----------|------|
| 生产事故关联证据 | 永久 | 涉及生产问题排查的证据单独归档 |
| 安全审计证据 | 365 天 | 安全相关测试证据保留一年 |
| 合规审计证据 | 按合规要求 | 遵循外部合规要求 |

---

## 4. 检索方式

### 4.1 检索主键

支持以下检索维度：

| 检索主键 | 格式 | 示例 | 说明 |
|----------|------|------|------|
| **提交号** | `commit:[SHA]` | `commit:abc123def456` | Git commit SHA |
| **Workflow Run ID** | `workflow:[run-id]` | `workflow:1234567890` | GitHub Actions Run ID |
| **日期范围** | `date:[from]-[to]` | `date:2026-01-01-2026-02-01` | ISO 日期格式 |
| **回归结论** | `verdict:[PASS/FAIL]` | `verdict:FAIL` | 回归结论过滤 |
| **分支名** | `branch:[name]` | `branch:main` | Git 分支名 |
| **标签** | `tag:[name]` | `tag:v1.2.0` | Git 标签 |

### 4.2 检索示例

```bash
# GitHub CLI 检索示例

# 按提交号检索
gh run list --commit abc123def456

# 按 Workflow 名称和日期检索
gh run list --workflow system-test-gate.yml --created "2026-01-01..2026-02-01"

# 按 Run ID 获取 Artifacts
gh run view 1234567890 --log

# 下载特定 Run 的 Artifacts
gh run download 1234567890 --dir ./evidence
```

### 4.3 GitHub Actions Artifacts 检索

| 操作 | 命令/路径 |
|------|----------|
| 列出 Artifacts | `gh api /repos/{owner}/{repo}/actions/artifacts` |
| 下载 Artifact | `gh run download <run-id>` |
| 查看 Run 日志 | `gh run view <run-id> --log` |
| Web UI | `https://github.com/{owner}/{repo}/actions/runs/<run-id>` |

---

## 5. 存储位置

### 5.1 主要存储

| 存储位置 | 用途 | 访问方式 |
|----------|------|----------|
| **GitHub Actions Artifacts** | CI/CD 执行证据 | GitHub CLI / Web UI |
| **GitHub Actions Logs** | 实时执行日志 | GitHub CLI / Web UI |
| **本地归档** | 特殊证据长期保存 | 文件系统 |

### 5.2 GitHub Actions Artifacts 配置

```yaml
# 示例：配置 retention-days
- name: Upload test artifacts
  uses: actions/upload-artifact@v4
  with:
    name: test-evidence-${{ github.sha }}
    path: |
      test-results/
      coverage/
      screenshots/
    retention-days: 90  # 默认值，失败证据可设为 180
```

### 5.3 Artifacts 命名约定

| 模式 | 示例 | 说明 |
|------|------|------|
| `{module}-test-{sha}` | `backend-test-abc123` | 按模块和提交号 |
| `{module}-coverage-{sha}` | `admin-coverage-abc123` | 覆盖率报告 |
| `{module}-e2e-{sha}` | `h5-e2e-abc123` | E2E 测试证据 |
| `regression-{run-id}` | `regression-1234567890` | 全量回归汇总 |

---

## 6. 与审计/复盘的对齐

### 6.1 审计周期覆盖

| 审计类型 | 周期 | 证据窗口 | 覆盖情况 |
|----------|------|----------|----------|
| 季度质量审计 | 每季度 | 90 天 | ✅ 完全覆盖 |
| 发布复盘 | 发布后 30 天内 | 180 天 | ✅ 完全覆盖 |
| 年度合规审计 | 每年 | 365 天（合规证据） | ✅ 合规证据单独处理 |

### 6.2 证据追溯矩阵

```
发布决策 ──→ 回归结论 ──→ Workflow Run ──→ Artifacts
    │              │              │              │
    │              │              │              ├── 测试日志
    │              │              │              ├── 覆盖率报告
    │              │              │              └── 失败截图
    │              │              │
    │              │              └── Commit SHA
    │              │
    │              └── PASS/FAIL/INCONCLUSIVE
    │
    └── 发布时间戳
```

### 6.3 复盘流程支持

| 复盘场景 | 证据需求 | 保留保障 |
|----------|----------|----------|
| 发布失败复盘 | 失败用例、日志、截图 | 180 天保留 |
| 性能回归分析 | 性能指标、基准对比 | 包含在测试日志中 |
| Flaky 测试分析 | 历史执行记录 | 通过 Run ID 追溯 |
| 审计合规检查 | 完整执行记录 | 90/180 天保留 |

---

## 7. CI 配置建议

### 7.1 推荐的 retention-days 配置

```yaml
# 全量回归工作流 - 使用条件判断设置保留时长
- name: Determine retention days
  id: retention
  run: |
    if [ "${{ steps.test.outcome }}" == "failure" ]; then
      echo "days=180" >> $GITHUB_OUTPUT
    else
      echo "days=90" >> $GITHUB_OUTPUT
    fi

- name: Upload test evidence
  uses: actions/upload-artifact@v4
  with:
    name: regression-evidence-${{ github.run_id }}
    path: evidence/
    retention-days: ${{ steps.retention.outputs.days }}
```

### 7.2 各 Workflow 建议配置

| Workflow | 建议 retention-days | 理由 |
|----------|---------------------|------|
| system-test-gate.yml | 90 (默认) / 180 (失败) | 全量回归门禁 |
| coverage-gate.yml | 90 | 覆盖率趋势分析 |
| stability-checks.yml | 90 | PR 验证 |
| order-mysql8-regression.yml | 90 / 180 (失败) | 订单专项 |

### 7.3 Artifacts 清理策略

```yaml
# 定期清理过期 Artifacts（GitHub 自动处理）
# 可通过以下方式验证保留状态：
# gh api /repos/{owner}/{repo}/actions/artifacts --jq '.artifacts[] | {name, expires_at}'
```

---

## 8. 术语表

| 术语 | 定义 |
|------|------|
| **证据** | 测试执行过程中产生的可追溯记录，包括日志、报告、截图等 |
| **保留窗口** | 证据从生成到删除的时间跨度 |
| **检索主键** | 用于查找特定证据的标识符或查询条件 |
| **Artifact** | GitHub Actions 中存储的工作流产出物 |
| **回归结论** | 全量回归的原子输出，值为 PASS/FAIL/INCONCLUSIVE 之一 |

---

## 9. 参考资料

- **Delta 规格**: `openspec/changes/add-full-regression-testing/specs/system-test-execution/spec.md`
- **全量回归定义**: `docs/testing/execution/FULL_REGRESSION_DEFINITION.md` (FRP-002)
- **GitHub Actions Artifacts 文档**: https://docs.github.com/en/actions/using-workflows/storing-workflow-data-as-artifacts
- **GitHub CLI 文档**: https://cli.github.com/manual/

---

*此文档由 FRP-005 任务生成，定义 DMH 项目的证据保留策略。*
