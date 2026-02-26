# T5: 证据归档命名规范

**生成时间**: 2026-02-26 20:38 CST

---

## 1. 命名规范

### 1.1 标准格式

```
.sisyphus/evidence/task-{N}-{slug}.{ext}
```

| 部分 | 说明 | 示例 |
|------|------|------|
| `task-{N}` | 任务编号 | `task-1`, `task-12` |
| `{slug}` | 简短描述（kebab-case） | `test-baseline-report` |
| `.{ext}` | 文件扩展名 | `.md`, `.json`, `.log` |

### 1.2 命名示例

```
.sisyphus/evidence/
├── task-1-test-baseline-report.md        # T1 基线报告
├── task-2-ci-layer-budget.md             # T2 CI 分层预算
├── task-3-repository-container-plan.md   # T3 容器方案
├── task-4-agents-test-spec-draft.md      # T4 测试规范草案
├── task-5-evidence-naming-spec.md        # T5 本文档
├── task-7-integration-test-result.log    # T7 集成测试日志
├── task-11-order-module-demo.json        # T11 示范模块数据
└── task-17-kpi-report.md                 # T17 KPI 报告
```

---

## 2. 文件类型约定

| 扩展名 | 用途 | 示例 |
|--------|------|------|
| `.md` | 文档、报告、规范 | 分析报告、方案文档 |
| `.json` | 结构化数据 | 测试结果、KPI 数据 |
| `.log` | 日志输出 | 测试日志、CI 日志 |
| `.txt` | 纯文本 | 命令输出、简单记录 |

---

## 3. 内容模板

### 3.1 文档模板 (`.md`)

```markdown
# T{N}: {任务标题}

**生成时间**: {YYYY-MM-DD HH:MM} CST
**执行者**: Sisyphus AI Agent

---

## 1. {章节}

{内容}

---

## 验收标准

- [x] {已完成项}
- [ ] {待完成项}

---

## 后续任务依赖

此任务解锁：
- T{N}: {依赖任务}
```

### 3.2 数据模板 (`.json`)

```json
{
  "task_id": "T{N}",
  "timestamp": "{ISO8601}",
  "status": "completed|in_progress|failed",
  "metrics": {
    "key": "value"
  },
  "evidence": [
    "path/to/evidence1",
    "path/to/evidence2"
  ]
}
```

---

## 4. 归档流程

### 4.1 创建证据

1. 执行任务
2. 生成证据文件
3. 保存到 `.sisyphus/evidence/task-{N}-{slug}.{ext}`

### 4.2 验证证据

```bash
# 检查命名规范
ls -la .sisyphus/evidence/

# 验证文件完整性
cat .sisyphus/evidence/task-1-*.md | head -20
```

---

## 5. 验收标准

- [x] 命名格式统一
- [x] 文件类型约定明确
- [x] 内容模板可用
- [x] 归档流程清晰

---

## 6. 实际应用

本规范已在以下任务中应用：
- T1: `task-1-test-baseline-report.md`
- T2: `task-2-ci-layer-budget.md`
- T3: `task-3-repository-container-plan.md`
- T4: `task-4-agents-test-spec-draft.md`
- T5: `task-5-evidence-naming-spec.md` (本文档)
