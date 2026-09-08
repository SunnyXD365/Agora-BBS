# Agora-BBS 项目 Git 开发规范与协作指南

## 1. 分支管理策略

代码库遵循简化版的 Feature Branch 工作流，分支划分如下：

| **分支类型** | **命名规则**         | **说明**                            | **允许直接 Push？**           |
| -------- | ---------------- | --------------------------------- | ------------------------ |
| **主分支**  | `main`           | 生产级稳定代码，始终保持随时可部署/可运行状态           | ❌ **严格禁止**（仅限 PR 合并）     |
| **功能分支** | `feat/<模块名>`     | 用于开发新功能（如 `feat/user-auth`）       | ✅ 允许（推送至远程供 Code Review） |
| **修复分支** | `fix/<缺陷名>`      | 用于修复特定 Bug（如 `fix/cooling-timer`） | ✅ 允许                     |
| **重构分支** | `refactor/<模块名>` | 用于代码重构、性能优化或架构调整                  | ✅ 允许                     |
| **文档分支** | `docs/<内容>`      | 用于更新项目设计或开发文档                     | ✅ 允许                     |
| 等等等等     | 这里不赘述            | 一般的功能性的分支                         | 都允许直接提交push              |

## 2. Commit Message 提交规范

本项目遵循 **Conventional Commits** 规范。统一的提交信息有助于快速 Code Review 以及自动生成 Changelog。
### 2.1 格式模板

```Plaintext
<type>(<scope>): <summary>

[可选的详细描述 Body]
```

#### Type（必填，改动性质）

- `feat`: 新增功能
- `fix`: 修复 Bug
- `docs`: 文档变更
- `style`: 代码格式调整（不影响逻辑的空格、分号等）
- `refactor`: 重构（既不修复 Bug 也不增加功能的代码改动）
- `perf`: 性能优化
- `test`: 增加或修改测试用例
- `chore`: 构建过程、依赖更新或辅助工具变动（如 Docker、.gitignore 修改）
#### Scope（选填，改动模块）

- `backend`: 后端 Go 相关改动
- `frontend`: 前端 Next.js 相关改动
- `workflow`: Temporal 工作流相关改动
- `docker`: Docker / Nginx / 网关配置改动

#### Summary（必填，简述）

- 用简短的语言描述改动，建议控制在 **50 字以内**，句末不加句号。

### 2.2 提交示例

```Bash
# 示例 1: 后端实现结构化发帖 API
git commit -m "feat(backend): 实现结构化主楼发帖 API 及 JSONB 校验"

# 示例 2: 前端修复冷静期倒计时 Bug
git commit -m "fix(frontend): 修复评论冷静期倒计时在页面刷新后重置的问题"

# 示例 3: 基础设施配置更新
git commit -m "chore(docker): 在 docker-compose 中添加 Temporal UI 环境变量"
```

## 3. 标准开发流程（8 步工作流）

```Plaintext
[拉取最新 main] ──► [新建功能分支] ──► [本地编码与提交]
                                             │
[合并到 main] ◄── [Code Review/合并] ◄── [推送远程并建 PR]
```

### 步骤 1：同步远程 `main` 最新代码

开始任何新任务前，先确保本地 `main` 分支是最新的：

```Bash
git checkout main
git pull origin main
```

### 步骤 2：创建并切换到功能分支

```Bash
git checkout -b feat/cooling-period
```

### 步骤 3：本地开发与频繁小步提交

在本地编写代码，按功能逻辑小步提交：

```Bash
git add .
git commit -m "feat(frontend): 完成评论冷静期倒计时 UI 组件"
```

### 步骤 4：推送分支到远程仓库

```Bash
	git push -u origin feat/cooling-period
```

### 步骤 5：在 GitHub 发起 Pull Request (PR)

1. 打开 GitHub 仓库页面：[https://github.com/SunnyXD365/Agora-BBS](https://github.com/SunnyXD365/Agora-BBS)
2. 点击 **Compare & pull request** 按钮。
3. 标题填写清晰明了（如 `feat(frontend): 新增评论冷静期撤回功能`）。
4. 在右侧 **Reviewers** 处指派组员进行审核。

### 步骤 6：Code Review

组员收到 PR 通知后，审查代码逻辑与规范：
- 若有修改建议：在 GitHub PR 页面提出 Inline Comment，开发者修改后重新提交推送。
- 若无问题：点击 **Approve** 通过审查。

### 步骤 7：合并 PR 到 `main`

- 审核通过后，推荐使用 **Squash and merge**（将该分支上的多次提交合并为一个干净的 Commit 写入 `main`），保持主分支历史整洁。

### 步骤 8：清理已完成的本地与远程分支

```Bash
git checkout main
git pull origin main
git branch -d feat/cooling-period               # 删除本地分支
git push origin --delete feat/cooling-period    # 删除远程分支
```

## 4. 冲突处理指南

当多人修改了同一个文件的同一位置时，合并时会产生冲突。请在**本地分支**上解决冲突后再提交 PR。
  

### 处理步骤：

1. **拉取远程 `main` 到本地分支**：
```bash
git checkout feat/your-feature
git fetch origin
git merge origin/main
```
2. **解决冲突**：   
    - 打开编辑器（VS Code / GoLand / WebStorm）。
    - 查找包含 `<<<<<<<`、`=======`、`>>>>>>>` 的冲突标记。
    - 与相关组员沟通，确认最终保留的代码 logic。

3. **提交解决后的代码**：
```Bash
git add .
git commit -m "chore: 解决与 main 分支的合并冲突"
git push origin feat/your-feature
```
## 5. 项目通用 `.gitignore` 模板

在项目根目录下创建 `.gitignore` 文件，确保敏感数据、编译产物与日志不被提交：

代码段

```
# ==========================================
# 操作系统与 IDE 缓存
# ==========================================
.DS_Store
Thumbs.db
.idea/
.vscode/
*.swp

# ==========================================
# 环境变量与敏感凭据 (严禁提交)
# ==========================================
.env
.env.local
.env.development.local
.env.test.local
.env.production.local
*.pem
*.key

# ==========================================
# Go 后端配置与编译产物
# ==========================================
backend/tmp/
backend/bin/
*.exe
*.exe-
*.dll
*.so
*.dylib

# Go 测试产物
*.test
*.out
coverage.txt

# ==========================================
# Next.js 前端配置与构建产物
# ==========================================
frontend/node_modules/
frontend/.next/
frontend/out/
frontend/build/
frontend/dist/

# 日志文件
*.log
npm-debug.log*
yarn-debug.log*
pnpm-debug.log*

# ==========================================
# Docker 与临时数据
# ==========================================
.dockerdata/
pgdata/
```

## 6. Git 常用命令速查表

|**使用场景**|**命令**|
|---|---|
|**查看修改状态**|`git status`|
|**查看修改内容差异**|`git diff`|
|**查看提交历史记录**|`git log --oneline --graph`|
|**撤销本地未暂存的修改**|`git restore <file>`|
|**修改上一次 Commit 的 Message**|`git commit --amend`|
|**暂存当前修改（紧急切换分支）**|`git stash` （恢复时使用 `git stash pop`）|
|**放弃本地所有未提交修改（慎用）**|`git reset --hard HEAD`|