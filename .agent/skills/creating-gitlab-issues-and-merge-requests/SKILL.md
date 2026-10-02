---
name: creating-gitlab-issues-and-merge-requests
description: Use when users ask in natural language to prepare, create, publish, recover, or update GitLab Issues and Merge Requests from local uncommitted work, pushed feature branches, or completed release changes in the EasyPos code repositories.
---

# Creating GitLab Issues and Merge Requests

## Overview

根据用户的自然语言判断代码所处状态，为 EasyPos 的代码仓库分别准备或发布 GitLab 变更请求。根仓库只负责协调，禁止为根仓库创建 Issue、MR、分支或提交。

这是严格的单代理工作流：禁止使用 sub agent、委派执行或并行 agent 工作流。

**REQUIRED SUB-SKILL:** Use $generating-git-changelogs for every repository selected for publication or an Issue-only release record.

执行 agent 必须在创建 manifest 前自动串联并遵循 `$generating-git-changelogs`。Node 脚本只做确定性的 Git/GitLab 检查与发布，不会自行调用 Skill，也不会编写 ChangeLog。

完成该必需子 Skill 后，才可创建 UTF-8 临时 manifest。

## 自然语言输入与固定范围

用户无需填写 `BASE=` 或 `TARGET=`。可接受例如：

- “把当前未提交的修改基于当前分支创建 Issue 和 MR。”
- “把当前未提交的修改基于 1.1.0-rc 创建 Issue 和 MR。”
- “为已经推送的 feature/order 向 release/1.2 创建 Issue 和 MR。”
- “检查 feature/order 是否已经合并；如果没有最终差异，只整理 ChangeLog。”
- “代码已经合并了，另外创建一条发布记录 Issue。”

例如“当前分支向 1.1.0-rc 发布”表示分别解析三个子仓库的当前分支。

仅允许以下精确映射：

- `base-app` → `oxygen/base-app`
- `base-engine` → `oxygen/base-engine`
- `base-web` → `oxygen/base-web`

三个仓库各自创建 Issue + MR，不创建跨项目总 Issue。干净且无相关差异的仓库标记为 `skipped`。

## 三状态路由

先以只读方式检查三个代码仓库，再选择且只选择一种流程：

1. **本地未提交代码（local-worktree）**：当前基线分支上存在 staged、unstaged、untracked、删除或重命名内容。执行本地预检、ChangeLog、preview、单独确认，再创建 Issue、Issue 派生分支、提交、普通推送和 MR。
2. **已推送但未合并的功能分支（pushed-branch）**：工作区干净，源分支同一 SHA 已在 GitLab，相对目标分支仍有最终 `diffs`，且当前 `sourceSha` 未被相同源/目标的 merged MR 记录为已合并。沿用只读本地 Git 的 `preflight` / `publish` 流程。
3. **已合并或无最终差异（issue-only）**：禁止创建空 MR。默认只输出 ChangeLog；只有用户明确要求发布记录 Issue 时，才走显式 Issue-only 流程。

普通 commit 区间可以生成 ChangeLog，但它本身不是可发布的 MR 源分支。若无法归入以上状态，说明阻断原因，不猜测或重写历史。

## 固定 ChangeLog 合同

对每个待发布仓库单独调用 `$generating-git-changelogs`，输出必须完整包含：

- 新增功能
- 功能优化
- 问题修复
- 升级注意事项
- 测试验证清单

“测试验证清单”是给外部测试人员的验收建议，不得写成实际已经执行或通过的验证证据。Issue 和 MR 使用同一份完整 ChangeLog；MR 描述首行固定为 `Closes #<issue_iid>`。

## 本地未提交流程

### 1. 只读预检

从根目录运行：

```bash
node .agent/skills/creating-gitlab-issues-and-merge-requests/scripts/publish-change-requests.mjs \
  local-preflight --target CURRENT
```

若用户点名基线分支，用该分支替换 `CURRENT`。每个 dirty 仓库必须位于该基线分支；基线 `baseSha` 必须以相同 SHA 存在于 GitLab。预检还必须拒绝 detached HEAD、冲突、merge/rebase/cherry-pick/revert 中间态、缺失 Git `user.name`/`user.email`、错误 origin、Developer 以下权限及未覆盖的本地变更。

预检用隔离临时 index/object 计算 `snapshotTreeSha` 和精确路径，不得改变真实 index、HEAD、分支或工作区。

### 2. 生成 manifest 与 preview

对每个 `ready` 仓库，让 `$generating-git-changelogs` 比较 `baseSha` 与完整 WORKTREE 快照，并生成版本 2 的 `local-worktree` manifest。每项保存 `title`、完整 `changelog`、`baseBranch`、`baseSha`、`snapshotTreeSha`、精确 `paths` 和 `branchSlug`。

`branchSlug` 只能是 1–63 位小写 ASCII 字母、数字或连字符；无法得到有意义英文名时使用 `change`。临时 manifest 必须 UTF-8、仅当前用户可读写，且不能包含 Token。

运行不带 `--apply` 的 preview：

```bash
node .agent/skills/creating-gitlab-issues-and-merge-requests/scripts/publish-change-requests.mjs \
  local-publish --manifest <临时文件>
```

向用户完整展示每个仓库的路径、Issue 标题、提交消息、目标分支、`feature/<issue_iid>-<branchSlug>` 分支模式和返回的 `confirmationDigest`。即使用户最初说“创建”或“发布”，此时也只能完成准备和 preview。

### 3. 单独确认后 apply

preview 之后必须取得用户单独、明确的确认；旧消息或最初的“创建”意图不能代替这次确认。确认后使用原样 digest：

```bash
node .agent/skills/creating-gitlab-issues-and-merge-requests/scripts/publish-change-requests.mjs \
  local-publish --manifest <临时文件> --apply --confirm <confirmationDigest>
```

apply 会重新检查 branch、SHA、树和路径是否漂移，然后按固定顺序逐仓执行：

1. 创建或复用带 ownership marker 的 Issue。
2. 创建/恢复 `feature/<issue_iid>-<branchSlug>` 并切换到它。
3. 只暂存确认的路径，创建一个 `Refs #<issue_iid>: <title>` 提交，并写入 `EasyPos-Change-Request: <snapshotTreeSha>` Trailer。
4. 校验父提交、树、Trailer、干净状态后执行普通非 force push。
5. 校验远端 SHA，创建或复用首行为 `Closes #<issue_iid>` 的 MR。

成功后保持在对应 feature 分支。宪章例外仅适用于这条 `local-worktree` apply 路径，且仅在 preview 后的单独确认生效；其他流程无权修改本地 Git 状态。

发生部分失败时立即停止，不自动 reset、clean、stash、rebase、amend、删除分支或回滚。报告已经存在的 Issue、分支、commit、push 和 MR；失败后保留 manifest 供安全续跑，仅在全部成功或用户明确取消后删除 manifest。

## 已推送功能分支流程

运行 `preflight`，其中当前分支在三个仓库分别解析：

```bash
node .agent/skills/creating-gitlab-issues-and-merge-requests/scripts/publish-change-requests.mjs \
  preflight --source CURRENT --target 1.1.0-rc
```

预检要求工作区干净、精确 origin、Developer+ 权限、源 `localSha` 与远端源 SHA 相等、目标 `localTargetSha` 与远端目标 SHA 相等，并拒绝 Compare 超时、缺失 `diffs` 或外来同源/目标 MR。还必须查询相同源/目标的 merged MR：其 head SHA 与当前 `sourceSha` 相同时直接标记为已合并；仅存在不同 head SHA 的旧 MR 时，允许把源分支上的新提交作为新变更继续发布。

对每个 `ready` 仓库用 `$generating-git-changelogs` 比较精确 `localTargetSha..localSha`，创建版本 1 manifest。每项除 `title`、`changelog` 外，必须保存预检解析后的 `sourceBranch`、`sourceSha`、`targetBranch`、`targetSha`；preview 和 apply 任一阶段发现这些值漂移都必须在写入前停止。先运行 `publish --manifest <临时文件>` 预览；用户已明确要求“创建/发布”时，再运行 `publish --manifest <临时文件> --apply`。该流程不得 checkout、commit 或 push。

## 已合并或无差异流程

GitLab Compare 的最终 `diffs` 为空时，不创建 MR。先按用户描述的历史范围生成 ChangeLog。只有用户明确要求 Issue-only 发布记录时，创建版本 2、`mode: issue-only` 的 manifest；每个仓库必须同时保存生成该 ChangeLog 时解析出的 `sourceSha` 和 `targetSha`，后续 preview/apply 必须拒绝任何 SHA 漂移。然后依次运行：

```bash
node .agent/skills/creating-gitlab-issues-and-merge-requests/scripts/publish-change-requests.mjs \
  issue-record --manifest <临时文件>
node .agent/skills/creating-gitlab-issues-and-merge-requests/scripts/publish-change-requests.mjs \
  issue-record --manifest <临时文件> --apply
```

Issue-only 描述必须明确它是已合并/无差异的发布或变更记录，不得包含 `Closes #...`，也不得调用 MR API。有当前最终差异时必须拒绝 Issue-only，改走 MR 流程。

## 安全与恢复合同

- 所有模式默认 preview；只更新带匹配 ownership marker 的对象，绝不覆盖外来分支或手工 MR。
- 三项目不具备跨项目事务性；某仓失败即停止，并列出此前已完成对象。
- 本地续跑只接受与 Issue IID、父 SHA、树、Trailer 和远端 SHA 全部吻合的 branch/commit/push/MR 状态。
- 禁止 force-push、amend、rebase、reset、restore、stash、clean、merge、branch deletion 和 tag mutation。
- GitLab API Token 只从进程环境或根 `.env` 的 `GITLAB_API_TOKEN` 读取；不得打印、写入 manifest、URL、Issue、MR 或 Git 命令。
- Git push 只使用仓库既有 `origin` 凭据，不把 API Token 拼入远端 URL。
- 所有 Git 子进程最长运行 120 秒，并设置 `GIT_TERMINAL_PROMPT=0`；凭据或网络不可用时必须失败并进入可恢复报告，不得无限等待。

## 配置

根 `.env` 必须包含：

```dotenv
GITLAB_BASE_URL=https://gitlab.sjfood.us
GITLAB_API_TOKEN=<具有三个项目 Developer+ 权限和 api scope 的令牌>
```
