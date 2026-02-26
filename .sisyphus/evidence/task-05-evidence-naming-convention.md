# 证据归档命名规范

## 1. 命名模板

```
task-{N}-{slug}.{ext}
```

### 组成部分

| 部分 | 说明 | 示例 |
|------|------|------|
| `N` | 任务编号 | `01`-`17`（主任务）, `F1`-`F4`（最终审核） |
| `slug` | 英文短描述（kebab-case） | `unit-test-report`, `coverage-summary` |
| `ext` | 文件扩展名 | `md`, `json`, `txt` |

### 命名规则

1. **任务编号**：两位数字或字母+数字组合
   - 主任务：`01`, `02`, ..., `17`
   - 最终审核：`F1`, `F2`, `F3`, `F4`

2. **Slug 规则**：
   - 全小写英文字母
   - 单词间用连字符 `-` 连接
   - 长度建议：3-30 字符
   - 避免使用日期（日期在文件元数据中记录）

3. **扩展名规则**：按文件类型选择

---

## 2. 目录结构

```
.sisyphus/evidence/
├── task-01-xxx.md                 # 任务01的证据文档
├── task-02-xxx.json               # 任务02的数据证据
├── ...
├── screenshots/                   # 截图目录
│   ├── task-03-test-output.png
│   └── task-05-coverage-chart.png
└── logs/                          # 日志目录
    ├── task-01-backend-test.log
    └── task-06-frontend-test.log
```

### 目录说明

| 目录 | 用途 | 文件类型 |
|------|------|----------|
| `.sisyphus/evidence/` | 证据根目录 | `.md`, `.json`, `.txt` |
| `.sisyphus/evidence/screenshots/` | 截图文件 | `.png`, `.jpg`, `.webp` |
| `.sisyphus/evidence/logs/` | 日志文件 | `.log`, `.txt` |

---

## 3. 文件类型

| 扩展名 | 类型 | 用途 | 示例 |
|--------|------|------|------|
| `.md` | 文档 | 测试报告、摘要、分析文档 | `task-05-report.md` |
| `.json` | 数据 | 测试结果、覆盖率数据、配置 | `task-03-coverage.json` |
| `.txt` | 文本 | 纯文本日志、命令输出 | `task-01-output.txt` |
| `.log` | 日志 | 运行日志、错误日志 | `logs/task-06-error.log` |
| `.png/.jpg` | 图片 | 截图、图表 | `screenshots/task-05-chart.png` |

---

## 4. 命名示例

### 主任务证据

| 任务 | 文件名 | 说明 |
|------|--------|------|
| T1 | `task-01-current-state.md` | 当前测试状态报告 |
| T2 | `task-02-coverage-baseline.json` | 覆盖率基线数据 |
| T3 | `task-03-backend-unit-tests.md` | 后端单元测试报告 |
| T4 | `task-04-frontend-unit-tests.md` | 前端单元测试报告 |
| T5 | `task-05-evidence-naming-convention.md` | 本规范文档 |
| T6 | `task-06-test-runner-script.md` | 测试运行器脚本 |

### 截图与日志

| 类型 | 文件名 | 说明 |
|------|--------|------|
| 截图 | `screenshots/task-03-test-pass.png` | 后端测试通过截图 |
| 截图 | `screenshots/task-06-coverage-report.png` | 覆盖率报告截图 |
| 日志 | `logs/task-01-backend-test.log` | 后端测试运行日志 |
| 日志 | `logs/task-04-frontend-test.log` | 前端测试运行日志 |

### 最终审核证据

| 任务 | 文件名 | 说明 |
|------|--------|------|
| F1 | `task-F1-backend-verification.md` | 后端验证报告 |
| F2 | `task-F2-frontend-verification.md` | 前端验证报告 |
| F3 | `task-F3-e2e-verification.md` | E2E 验证报告 |
| F4 | `task-F4-final-summary.md` | 最终汇总报告 |

---

## 5. 证据文件内容模板

### Markdown 报告模板

```markdown
# Task {N}: {任务标题}

## 执行时间
- 开始: YYYY-MM-DD HH:MM
- 结束: YYYY-MM-DD HH:MM

## 执行内容
[描述执行的具体操作]

## 结果
[执行结果摘要]

## 证据
- 相关截图: [如有]
- 相关日志: [如有]

## 备注
[其他说明]
```

### JSON 数据模板

```json
{
  "task_id": "T{N}",
  "timestamp": "YYYY-MM-DDTHH:MM:SSZ",
  "type": "coverage|test_result|config",
  "data": {
    // 具体数据
  }
}
```

---

## 6. 验证检查清单

使用本规范时，确保：

- [ ] 文件名符合 `task-{N}-{slug}.{ext}` 格式
- [ ] 任务编号正确（01-17 或 F1-F4）
- [ ] slug 使用 kebab-case 且全小写
- [ ] 扩展名与文件内容类型匹配
- [ ] 截图放入 `screenshots/` 子目录
- [ ] 日志放入 `logs/` 子目录

---

*规范版本: 1.0*  
*创建日期: 2026-02-26*  
*适用范围: DMH 测试优化计划（T1-T17, F1-F4）*
