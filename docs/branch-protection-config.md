# DMH 分支保护规则配置清单

## 当前状态分析

**现有保护规则**：
- Required checks: `Backend full tests + Admin unit + Security E2E`
- PR reviews: 需要至少 1 个 approval
- Admin bypass: ✗ 允许管理员绕过（有问题！）
- Linear history: ✓ 强制

**可用的 workflow checks**：
- `Coverage Gate` - 覆盖率门禁
- `Full Regression` - 完整回归测试
- `Stability Checks` - 稳定性检查
- `System Test Gate` - 系统测试门禁

---

## 推荐配置方案

### 方案：严格保护模式

#### 1. Required Status Checks

启用以下检查，全部通过才能合并：

| Check Name | 说明 | 触发条件 |
|------------|------|----------|
| **PR Gate Verdict** | 快速单元测试+覆盖率（15分钟） | PR 创建/更新 |
| **Coverage Gate** | 后端+Admin+H5 覆盖率阈值 | PR 创建/更新 + Push main |
| **System Test Gate** | 集成测试 + E2E + OpenSpec | PR 创建/更新 |
| **Full Regression** | 完整回归（10分钟定时+PR） | Push main + RC tags + 每日定时 |

**Strict Mode**: ✓ 启用（必须是最新的分支状态）

#### 2. Pull Request Reviews

| 配置项 | 值 | 说明 |
|--------|-----|------|
| Require approvals | 1 人 | 至少需要 1 个 reviewer 批准 |
| Dismiss stale reviews | ✓ 启用 | 有新提交时自动取消过期的批准 |
| Require last push approval | ✓ 启用 | 最后一次推送后需要新的批准 |
| Require code owner reviews | ✗ 禁用 | 不强制代码所有者审核（小型团队） |

#### 3. 其他保护规则

| 配置项 | 值 | 说明 |
|--------|-----|------|
| Enforce administrators | ✓ 启用 | **关键**：管理员也必须遵守规则 |
| Require linear history | ✓ 启用 | 强制线性历史，避免 merge commits |
| Require conversation resolution | ✓ 启用 | 必须解决所有对话才能合并 |
| Allow force pushes | ✗ 禁用 | 禁止强制推送 |
| Allow deletions | ✗ 禁用 | 禁止删除分支 |

---

## 配置方式

### 方式 1：通过 GitHub Web UI（推荐）

1. 进入仓库 Settings → Branches
2. 点击 `Add rule` 或编辑 `main` 分支规则
3. 按照上述方案逐项配置

### 方式 2：通过 GitHub CLI

```bash
# 更新 required status checks
gh api --method PUT \
  repos/Gujiaweiguo/dmh/branches/main/protection \
  --field strict=true \
  --field enforce_admins=true \
  --field required_status_checks='{
    "strict": true,
    "contexts": [
      "PR Gate Verdict",
      "Coverage Gate",
      "System Test Gate",
      "Full Regression"
    ]
  }' \
  --field required_pull_request_reviews='{
    "dismiss_stale_reviews": true,
    "require_last_push_approval": true,
    "required_approving_review_count": 1
  }' \
  --field required_linear_history='{"enabled": true}' \
  --field required_conversation_resolution='{"enabled": true}'
```

### 方式 3：通过 API（高级）

```json
PUT /repos/Gujiaweiguo/dmh/branches/main/protection

{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "PR Gate Verdict",
      "Coverage Gate",
      "System Test Gate",
      "Full Regression"
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "dismiss_stale_reviews": true,
    "require_last_push_approval": true,
    "required_approving_review_count": 1
  },
  "required_linear_history": {
    "enabled": true
  },
  "required_conversation_resolution": {
    "enabled": true
  },
  "allow_force_pushes": false,
  "allow_deletions": false
}
```

---

## 配置验证清单

配置完成后，验证以下项：

- [ ] 管理员无法绕过保护规则（测试：管理员尝试直接 push 到 main）
- [ ] PR 必须通过所有 4 个 checks 才能合并
- [ ] PR 代码更新后，过期的 approval 被自动取消
- [ ] 无法强制推送到 main
- [ ] 无法删除 main 分支
- [ ] 必须解决所有 PR 评论才能合并

---

## 配置后的工作流程

### 正常开发流程

1. 创建 feature 分支
2. 推送代码并创建 PR
3. **PR Gate Verdict** 自动运行（快速检查）
4. 通过后，**Coverage Gate** 运行
5. 通过后，**System Test Gate** 运行（较慢）
6. 所有 checks 通过后，进行 code review
7. Reviewer approve 后合并到 main
8. 合并后触发 **Full Regression**

### 热修复流程（Hotfix）

1. 创建 hotfix 分支
2. 推送代码并创建 PR
3. 同样需要通过所有 checks
4. Review 后合并到 main
5. 手动触发 Full Regression（如需快速发布）

---

## 故障排查

### PR Checks 不运行

- 检查 workflow 文件是否正确触发
- 确认 workflow 状态不是 disabled
- 查看 workflow run logs 确认错误

### 无法合并 PR

- 检查所有 required checks 状态
- 确认是否有未解决的 conversations
- 确认是否有过期的 approvals 需要重新 approve

---

## 附录：Workflow 触发矩阵

| Workflow | PR | Push main | 定时 | RC tags |
|----------|:--:|:---------:|:----:|:--------:|
| **PR Gate** | ✓ | ✗ | ✗ | ✗ |
| **Stability Checks** | ✓ | ✓ | ✗ | ✗ |
| **Coverage Gate** | ✓ | ✓ | ✗ | ✗ |
| **System Test Gate** | ✓ | ✓ | ✗ | ✗ |
| **Full Regression** | ✗ | ✓ | ✓ 10:00 | ✓ |
