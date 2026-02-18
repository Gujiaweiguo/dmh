# Plan v1 - add-full-regression-testing

## Authority (唯一执行依据)
- 本文件即 **Plan v1**，是 Atlas/Hephaestus 执行 `add-full-regression-testing` 的唯一依据。
- 禁止生成并使用独立计划文档；如需变更顺序/范围/风险，必须先修改本文件并评审通过。

## Dependency Order / Regression Batches / Risk

| Task ID | Depends On | 回归批次 | 风险等级 |
|---|---|---|---|
| FRP-001 | None | Batch-R0 基线 | 中 |
| FRP-002 | FRP-001 | Batch-R1 规格口径 | 中 |
| FRP-003 | FRP-002 | Batch-R1 规格口径 | 高 |
| FRP-004 | FRP-002 | Batch-R1 规格口径 | 中 |
| FRP-005 | FRP-002 | Batch-R1 规格口径 | 低 |
| FRP-006 | FRP-003, FRP-004, FRP-005 | Batch-R2 CI 对齐 | 高 |
| FRP-007 | FRP-006 | Batch-R2 本地入口 | 中 |
| FRP-008 | FRP-005, FRP-006 | Batch-R2 证据治理 | 低 |
| FRP-009 | FRP-007, FRP-008 | Batch-R3 全量回归演练 | 高 |
| FRP-010 | FRP-009 | Batch-R4 发布前签收 | 中 |

## Tasks

- [x] **FRP-001 基线盘点冻结**
  - **输入**: `test-plan.md`、现有 workflow（`system-test-gate.yml`/`stability-checks.yml`/`coverage-gate.yml`/`order-mysql8-regression.yml`）、`Makefile`、`backend/scripts/*.sh`
  - **输出**: 全量回归基线清单（必跑套件、触发条件、当前缺口）
  - **验收标准**:
    - 明确 backend/admin/h5/OpenSpec 四层必跑范围
    - 明确当前“已覆盖 vs 未覆盖”项及文件级引用
  - **回滚方案**: 放弃本任务产出并恢复到既有基线定义，不改动任何执行脚本

- [x] **FRP-002 统一全量回归口径**
  - **输入**: FRP-001 基线清单、`openspec/specs/system-test-execution/spec.md`
  - **输出**: 统一口径映射（全量回归 = 必跑矩阵 + 单一 PASS/FAIL 结论）
  - **验收标准**:
    - 口径覆盖：后端单元、后端集成、订单专项、前端双端单元+E2E、OpenSpec 校验
    - 与现有 `system-test-execution` requirement 无重复冲突
  - **回滚方案**: 恢复到“仅分散 workflow 语义”，不启用统一结论口径

- [x] **FRP-003 发布阻断规则落地**
  - **输入**: FRP-002 统一口径、现有发布门禁规则（P0=100%，P1>=95%）
  - **输出**: 自动阻断规则定义（必跑套件失败/未执行即阻断）
  - **验收标准**:
    - 发布阻断条件可映射到 CI 状态
    - 阻断报告字段包含：失败套件、失败用例、workflow 链接、重跑记录
  - **回滚方案**: 回退为原“人工判断发布”模式，保留旧门禁阈值不变

- [x] **FRP-004 flaky 策略落地**
  - **输入**: FRP-002 统一口径、现有测试重跑能力
  - **输出**: flaky 白名单/隔离与重试上限策略（受控重试 <=2）
  - **验收标准**:
    - 明确可重试范围（仅环境瞬态失败）
    - 隔离项包含责任人、截止时间、退出标准
  - **回滚方案**: 关闭 flaky 自动重试与隔离策略，恢复一次失败即失败

- [x] **FRP-005 证据保留窗口定义**
  - **输入**: FRP-002 统一口径、现有 artifact 上传策略
  - **输出**: 证据保留与检索规则（成功>=90天、失败>=180天）
  - **验收标准**:
    - 明确保留时长、证据类型、检索主键（提交号/workflow）
    - 与审计/复盘流程对齐
  - **回滚方案**: 回退到当前默认 artifact 保留策略，不启用分层保留

- [x] **FRP-006 CI 编排对齐实现**
  - **输入**: FRP-003/004/005 的规则定义
  - **输出**: CI 工作流对齐方案（`.github/workflows/*.yml`）
  - **验收标准**:
    - 全量回归触发条件覆盖 release candidate 与受保护分支变更
    - 产出单一回归结论 PASS/FAIL，且可用于阻断发布
  - **回滚方案**: 回退 workflow 变更到上一稳定版本并恢复原触发逻辑

- [x] **FRP-007 本地一键入口对齐**
  - **输入**: FRP-006 CI 编排方案、`Makefile`、`backend/scripts/*.sh`
  - **输出**: 本地一键全量回归入口定义（与 CI 口径一致）
  - **验收标准**:
    - 本地入口执行范围与 CI 必跑矩阵一致
    - 文档给出标准调用方式与环境前置要求
  - **回滚方案**: 保留原分散命令入口，移除新增聚合入口定义

- [x] **FRP-008 证据模板与审计落地**
  - **输入**: FRP-005 保留策略、FRP-006 输出字段
  - **输出**: 证据模板与审计清单（执行时间/范围/通过率/失败明细/日志引用）
  - **验收标准**:
    - 模板可直接覆盖 backend/admin/h5/OpenSpec 四类结果
    - 审计字段可追溯到具体 workflow run
  - **回滚方案**: 继续沿用现有零散证据记录方式

- [x] **FRP-009 Batch-R3 全量回归演练**
  - **输入**: FRP-007 本地入口、FRP-008 证据模板
  - **输出**: 一轮完整演练结果与问题清单
  - **验收标准**:
    - 必跑矩阵全部执行完成并生成统一结论
    - 演练问题具备“定位信息+修复建议+优先级”
  - **回滚方案**: 暂停新编排，回退到旧回归链路进行发布前验证

- [x] **FRP-010 发布前签收与严格校验**
  - **输入**: FRP-009 演练报告、OpenSpec 变更文件
  - **输出**: Plan v1 签收结论 + OpenSpec 严格校验结果
  - **验收标准**:
    - `openspec validate add-full-regression-testing --strict` 通过
    - Atlas/Hephaestus 签收记录引用 Plan v1 任务 ID 完整闭环
  - **回滚方案**: 若校验或签收失败，冻结发布并回退至 FRP-009 前的稳定方案
