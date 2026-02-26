# DMH 测试资产基线统计

> 执行时间: 2026-02-26 13:39:15 - 13:40:00 (CST)
> 执行命令:
> - Backend: `go test ./... -short -count=1`
> - Admin: `npm run test:cov -- --run`
> - H5: `npm run test:cov -- --run`

---

## 1. Backend (Go)

| 指标 | 数值 |
|------|------|
| 测试文件数 | 106 |
| 通过的包 | 19 |
| 失败的包 | 14 |
| 失败测试用例 | 21 |
| 跳过测试用例 | 0 |
| 覆盖率阈值 | 78% |

### 失败的包列表
- dmh/api/internal/handler/admin
- dmh/api/internal/handler/auth
- dmh/api/internal/handler/brand
- dmh/api/internal/handler/distributor
- dmh/api/internal/handler/feedback
- dmh/api/internal/handler/promoter

### 失败原因摘要
- 外键约束失败 (foreign key constraint)
- 数据库记录未找到 (record not found)
- 重复主键冲突 (Duplicate entry)
- 测试数据隔离问题

---

## 2. Frontend Admin (Vue3 + Vitest)

| 指标 | 数值 |
|------|------|
| 测试文件数 | 41 |
| 测试用例总数 | 394 |
| 通过 | 394 |
| 失败 | 0 |
| 跳过 | 0 |
| 覆盖率 (Statements) | 83.65% |
| 覆盖率 (Branches) | 80.32% |
| 覆盖率 (Functions) | 81.45% |
| 覆盖率 (Lines) | 83.65% |
| 覆盖率阈值 | 80% |
| 执行时间 | 9.75s |

### 覆盖率详情
```
File               | % Stmts | % Branch | % Funcs | % Lines
-------------------|---------|----------|---------|--------
All files          |   83.65 |    80.32 |   81.45 |   83.65
 components        |   74.23 |    98.61 |   69.23 |   74.23
 services          |   93.87 |    90.27 |   95.79 |   93.87
 views             |   78.17 |    63.11 |    68.1 |   78.17
```

---

## 3. Frontend H5 (Vue3 + Vitest)

| 指标 | 数值 |
|------|------|
| 测试文件数 | 56 |
| 测试用例总数 | 1090 |
| 通过 | 1090 |
| 失败 | 0 |
| 跳过 | 0 |
| 覆盖率 (Statements) | 87.37% |
| 覆盖率 (Branches) | 96.5% |
| 覆盖率 (Functions) | 98.79% |
| 覆盖率 (Lines) | 87.37% |
| 覆盖率阈值 | 70% |
| 执行时间 | 10.79s |

### 覆盖率详情
```
File               | % Stmts | % Branch | % Funcs | % Lines
-------------------|---------|----------|---------|--------
All files          |   87.37 |     96.5 |   98.79 |   87.37
 frontend-h5/src   |       0 |        0 |       0 |       0
 ...src/services   |   97.22 |     90.9 |   98.55 |   97.22
 ...src/utils      |   99.54 |    84.93 |     100 |   99.54
 ...src/views      |    98.4 |    94.52 |     100 |    98.4
```

---

## 4. 汇总

| 模块 | 测试文件 | 测试用例 | 通过率 | 覆盖率 | 阈值 | 状态 |
|------|---------|---------|--------|--------|------|------|
| Backend | 106 | - | ~57% (包) | - | 78% | ⚠️ 有失败 |
| Admin | 41 | 394 | 100% | 83.65% | 80% | ✅ 通过 |
| H5 | 56 | 1090 | 100% | 87.37% | 70% | ✅ 通过 |

---

## 5. 问题与风险

### Backend
1. **14个测试包失败** - 主要由于数据库集成测试的数据隔离问题
2. **外键约束问题** - 测试数据清理顺序不当
3. **建议**: 需要修复集成测试的数据准备逻辑

### Admin & H5
1. **全部通过** - 前端单元测试稳定
2. **覆盖率达标** - 均超过阈值要求

---

## 6. 原始输出位置

- Backend: `/tmp/backend-test-output.txt`
- Admin: `/tmp/admin-test-output.txt`
- H5: `/tmp/h5-test-output.txt`

---

*此基线用于后续测试改进的参考起点*
