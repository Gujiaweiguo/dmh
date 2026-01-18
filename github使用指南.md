# GitHub 使用指南（单人多机同步，SSH 免密）

适用场景：你一个人开发，但会在多台机器之间切换；使用 GitHub 作为唯一真源；仓库 remote 为 SSH（如：`git@github.com:Gujiaweiguo/DMH.git`）。

## 1. 核心原则（最重要）

- GitHub 远程仓库（`origin`）是唯一真源（source of truth）。
- 每台机器各自 `git clone` 一份工作目录；**不需要**为了“拿最新代码”去删除本地目录再重新下载。
- 切换机器前：必须把改动变成可同步状态，通常是 `commit + push`。
- 到另一台机器开工前：先 `pull --rebase` 拿到最新，再开始改。

## 2. 你刚才这次操作的含义（为什么显示 nothing to commit）

你执行的是：

- `git status`
- `git add -A`
- `git commit -m "xxx"`
- `git push`

当 `git status` 显示：

```
nothing to commit, working tree clean
```

含义是：工作区没有任何文件变更，所以：

- `git commit -m "xxx"` 会提示 `nothing to commit`（没有可提交的内容）
- `git push` 会提示 `Everything up-to-date`（远程也没有需要更新的提交）

这是正常现象；只有你实际修改了文件（或新增/删除文件）之后，才会产生可提交的变更。

## 3. 日常固定流程（强烈建议照做）

### 3.1 开工（进入项目目录后）

```
git switch main
git pull --rebase
git status
```

### 3.2 收工 / 切换机器前（必须）

```
git status
git add -A
git commit -m "描述本次改动"   # 不想写也可临时用 "wip"
git push
```

### 3.3 想看看“到底改了什么”（可选）

- 看未暂存变更：`git diff`
- 看已暂存（已 `add`）内容：`git diff --staged`

## 4. 第二台机器：首次配置 SSH key（一次性）

> 每台机器建议各自生成自己的 SSH key，并分别添加到 GitHub（更安全；不建议拷贝私钥到另一台机器）。

### 4.1 配置 Git 身份（每台机器都要）

```
git config --global user.name "你的名字"
git config --global user.email "你的邮箱"
```

### 4.2 建议开启的默认行为（每台机器都建议设置）

```
git config --global pull.rebase true
git config --global rebase.autoStash true
```

### 4.3 生成 SSH key（没有就生成）

```
ssh-keygen -t ed25519 -C "你的邮箱"
```

按提示选择保存路径（默认 `~/.ssh/id_ed25519`），是否设置 passphrase 取决于你的安全习惯。

### 4.4 启动 ssh-agent 并加载 key（推荐）

```
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519
```

### 4.5 将公钥添加到 GitHub

```
cat ~/.ssh/id_ed25519.pub
```

复制输出内容到 GitHub：

- GitHub 头像 → `Settings`
- `SSH and GPG keys`
- `New SSH key`

### 4.6 验证 SSH 登录

```
ssh -T git@github.com
```

出现 “successfully authenticated” 类似提示即可（提示不提供 shell 是正常的）。

### 4.7 克隆仓库（第二台机器）

```
git clone git@github.com:Gujiaweiguo/DMH.git
cd DMH
git status
```

## 5. “拿最新代码”时的正确做法（不删除目录）

进入项目目录后先看状态：

```
git status
```

### 5.1 工作区干净（推荐情况）

```
git pull --rebase
```

### 5.2 工作区不干净（有未提交改动）

推荐：先提交再拉（最稳，适合跨机器同步）

```
git add -A
git commit -m "wip"
git pull --rebase
git push
```

备选：临时 stash（不建议长期依赖来跨机器同步）

```
git stash -u
git pull --rebase
git stash pop
```

## 6. 冲突处理（偶尔会遇到）

如果 `git pull --rebase` 提示冲突：

1. 打开冲突文件手工解决（保留你要的内容）
2. 标记并继续：

```
git add <冲突文件>
git rebase --continue
```

放弃这次 rebase：

```
git rebase --abort
```

完成后推送：

```
git push
```

## 7. 强制让本地完全等于远程（危险操作）

当你明确：本机所有本地改动都不要，只要远程 `main` 最新：

```
git fetch origin
git switch main
git reset --hard origin/main
```

如果还要删除未跟踪文件/目录（更危险）：

```
git clean -nd   # 先预览会删什么
git clean -fd   # 确认后执行
```

