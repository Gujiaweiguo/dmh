# 🤝 贡献指南

感谢您对 DMH 项目的关注！我们欢迎任何形式的贡献。

## 目录

* [行为准则](#行为准则)
* [如何贡献](#如何贡献)
* [开发流程](#开发流程)
* [代码规范](#代码规范)
* [提交规范](#提交规范)
* [问题反馈](#问题反馈)

***

## 行为准则

### 我们的承诺

为了营造一个开放和友好的环境，我们承诺：

* 使用友好和包容的语言
* 尊重不同的观点和经验
* 优雅地接受建设性批评
* 关注对社区最有利的事情
* 对其他社区成员表示同理心

### 不可接受的行为

* 使用性化的语言或图像
* 人身攻击或侮辱性评论
* 公开或私下的骚扰
* 未经许可发布他人的私人信息
* 其他不道德或不专业的行为

***

## 如何贡献

### 贡献方式

您可以通过以下方式为项目做出贡献：

1. **报告 Bug** - 发现问题请提交 Issue
2. **建议功能** - 有好的想法请告诉我们
3. **改进文档** - 帮助完善项目文档
4. **提交代码** - 修复 Bug 或实现新功能
5. **代码审查** - 帮助审查其他人的 PR
6. **分享项目** - 向更多人推荐 DMH

### 第一次贡献？

如果您是第一次为开源项目做贡献，可以从以下开始：

* 查看标记为 `good first issue` 的问题
* 改进文档中的错别字或不清楚的地方
* 添加或改进代码注释
* 编写或改进测试用例

***

## 开发流程

### 1. Fork 项目

点击项目页面右上角的 "Fork" 按钮，将项目 Fork 到您的账号下。

### 2. 克隆仓库

```bash
git clone https://github.com/YOUR_USERNAME/DMH.git
cd DMH
```

### 3. 添加上游仓库

```bash
git remote add upstream https://github.com/Gujiaweiguo/DMH.git
```

### 4. 创建分支

```bash
# 从 main 分支创建新分支
git checkout -b feature/your-feature-name

# 或修复 Bug
git checkout -b fix/your-bug-fix
```

分支命名规范：

* `feature/xxx` - 新功能
* `fix/xxx` - Bug 修复
* `docs/xxx` - 文档更新
* `refactor/xxx` - 代码重构
* `test/xxx` - 测试相关

### 5. 开发和测试

```bash
# 搭建开发环境
./dmh.sh init

# 进行开发
# ...

# 运行测试
cd backend && go test ./...
cd frontend-admin && npm test
```

### 6. 提交更改

```bash
# 添加更改
git add .

# 提交（遵循提交规范）
git commit -m "feat: 添加用户导出功能"
```

### 7. 同步上游更改

```bash
# 获取上游更新
git fetch upstream

# 合并到本地分支
git rebase upstream/main
```

### 8. 推送到 Fork 仓库

```bash
git push origin feature/your-feature-name
```

### 9. 创建 Pull Request

1. 访问您的 Fork 仓库页面
2. 点击 "New Pull Request"
3. 选择您的分支
4. 填写 PR 描述（参考模板）
5. 提交 PR

***

## 代码规范

### Go 代码规范

```go
// 1. 包名使用小写
package handler

// 2. 导出的函数首字母大写
func CreateUser(req *CreateUserRequest) error {
    // 3. 使用有意义的变量名
    user := &User{
        Username: req.Username,
        Email:    req.Email,
    }
    
    // 4. 错误处理
    if err := validateUser(user); err != nil {
        return fmt.Errorf("validate user failed: %w", err)
    }
    
    // 5. 添加必要的注释
    // 保存用户到数据库
    return userRepo.Save(user)
}

// 6. 私有函数首字母小写
func validateUser(user *User) error {
    if user.Username == "" {
        return errors.New("username is required")
    }
    return nil
}
```

**代码检查**：

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 运行 linter
golangci-lint run
```

### 前端代码规范

```vue
<template>
  <!-- 1. 使用语义化标签 -->
  <div class="user-form">
    <!-- 2. 使用 v-bind 简写 -->
    <input 
      :value="username" 
      @input="handleInput"
      placeholder="请输入用户名"
    />
    
    <!-- 3. 使用 v-if/v-show 控制显示 -->
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script>
// 4. 使用 Composition API
import { ref, computed } from 'vue';

export default {
  name: 'UserForm',
  
  setup() {
    // 5. 使用有意义的变量名
    const username = ref('');
    const error = ref('');
    
    // 6. 函数命名使用动词开头
    const handleInput = (event) => {
      username.value = event.target.value;
    };
    
    // 7. 计算属性使用 computed
    const isValid = computed(() => {
      return username.value.length >= 3;
    });
    
    return {
      username,
      error,
      handleInput,
      isValid
    };
  }
};
</script>

<style scoped>
/* 8. 使用 scoped 样式 */
.user-form {
  padding: 20px;
}

.error {
  color: red;
}
</style>
```

**代码检查**：

```bash
# ESLint 检查
npm run lint

# 格式化代码
npm run format
```

***

## 提交规范

### Commit Message 格式

使用 [Conventional Commits](https://conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

* `feat`: 新功能
* `fix`: Bug 修复
* `docs`: 文档更新
* `style`: 代码格式（不影响代码运行）
* `refactor`: 重构（既不是新功能也不是 Bug 修复）
* `perf`: 性能优化
* `test`: 测试相关
* `build`: 构建系统或外部依赖
* `ci`: CI 配置文件和脚本
* `chore`: 其他不修改 src 或 test 的更改
* `revert`: 回退之前的 commit

### Scope 范围（可选）

* `backend`: 后端相关
* `frontend`: 前端相关
* `admin`: 管理后台
* `h5`: H5 端
* `api`: API 相关
* `db`: 数据库相关
* `docs`: 文档相关

### 示例

```bash
# 新功能
git commit -m "feat(backend): 添加用户导出功能"

# Bug 修复
git commit -m "fix(frontend): 修复登录页面样式问题"

# 文档更新
git commit -m "docs: 更新 API 文档"

# 重构
git commit -m "refactor(backend): 重构用户服务代码"

# 性能优化
git commit -m "perf(db): 优化用户查询性能"

# 带详细描述
git commit -m "feat(backend): 添加会员标签功能

- 添加标签创建接口
- 添加标签关联接口
- 添加标签查询接口

Closes #123"
```

***

## Pull Request 规范

### PR 标题

PR 标题应该清晰描述更改内容，格式同 Commit Message：

```
feat(backend): 添加用户导出功能
fix(frontend): 修复登录页面样式问题
```

### PR 描述模板

```markdown
## 变更类型
- [ ] 新功能
- [ ] Bug 修复
- [ ] 文档更新
- [ ] 代码重构
- [ ] 性能优化
- [ ] 测试相关
- [ ] 其他

## 变更说明
<!-- 描述本次 PR 的主要变更内容 -->

## 相关 Issue
<!-- 关联的 Issue 编号，如 #123 -->
Closes #

## 测试说明
<!-- 如何测试本次变更 -->
- [ ] 已添加单元测试
- [ ] 已添加集成测试
- [ ] 已手动测试

## 截图（如果适用）
<!-- 添加截图说明变更效果 -->

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 已添加必要的注释
- [ ] 已更新相关文档
- [ ] 所有测试通过
- [ ] 没有引入新的警告
```

### PR 审查

提交 PR 后，维护者会进行代码审查。请：

1. **及时响应**审查意见
2. **友好讨论**不同观点
3. **积极改进**代码质量
4. **保持耐心**等待审查

***

## 问题反馈

### 报告 Bug

提交 Bug 时，请包含以下信息：

1. **Bug 描述** - 清晰描述问题
2. **复现步骤** - 如何重现问题
3. **预期行为** - 应该发生什么
4. **实际行为** - 实际发生了什么
5. **环境信息** - 操作系统、浏览器、版本等
6. **截图/日志** - 如果适用

### Bug 报告模板

```markdown
## Bug 描述
<!-- 清晰简洁地描述 Bug -->

## 复现步骤
1. 访问 '...'
2. 点击 '...'
3. 滚动到 '...'
4. 看到错误

## 预期行为
<!-- 描述应该发生什么 -->

## 实际行为
<!-- 描述实际发生了什么 -->

## 截图
<!-- 如果适用，添加截图 -->

## 环境信息
- OS: [e.g. macOS 13.0]
- Browser: [e.g. Chrome 120]
- Version: [e.g. 1.0.0]

## 附加信息
<!-- 其他相关信息 -->
```

### 功能建议

提交功能建议时，请说明：

1. **功能描述** - 想要什么功能
2. **使用场景** - 为什么需要这个功能
3. **解决方案** - 如何实现（可选）
4. **替代方案** - 其他可能的方案（可选）

***

## 代码审查指南

### 审查者

作为审查者，请：

1. **及时审查** - 尽快审查 PR
2. **建设性反馈** - 提供具体的改进建议
3. **友好沟通** - 保持尊重和友好
4. **关注重点** - 代码质量、性能、安全性

### 被审查者

作为被审查者，请：

1. **虚心接受** - 认真对待审查意见
2. **积极改进** - 根据反馈改进代码
3. **友好讨论** - 如有不同意见，友好讨论
4. **及时响应** - 尽快回复审查意见

***

## 发布流程

### 版本发布

1. 更新 `CHANGELOG.md`
2. 更新版本号
3. 创建 Git Tag
4. 发布 Release

```bash
# 更新版本号
git tag v1.0.0

# 推送 Tag
git push origin v1.0.0

# 创建 Release（在 GitHub 上）
```

***

## 获取帮助

如果您有任何问题，可以：

1. **查看文档** - [README.md](./README.md), [DEVELOPMENT.md](./DEVELOPMENT.md)
2. **搜索 Issues** - 看看是否有人遇到过类似问题
3. **提交 Issue** - 描述您的问题
4. **联系维护者** - weiguogu@163.com

***

## 致谢

感谢所有为 DMH 项目做出贡献的开发者！

您的贡献让这个项目变得更好！🎉

***

**维护者**: DMH Team\
**最后更新**: 2025-01-21
