---
name: generating-git-changelogs
description: Use when users ask in natural language for release notes, a ChangeLog, a version upgrade summary, or a user-visible comparison between branches, tags, commits, commit ranges, or current uncommitted code.
---

# Generating Git Changelogs

## Overview

Generate an evidence-backed, user-facing ChangeLog from two Git states described conversationally. Reconstruct complete feature lineages instead of mechanically rewriting commit subjects.

This is a strictly single-agent, read-only workflow. Do not use sub-agents.

## Natural-Language Input

Accept an ordinary sentence naming the two states to compare. Never require parameter names, key-value syntax, or special invocation grammar. The states may be local branches, remote-tracking branches, tags, full or abbreviated commit SHAs, an `A..B` commit range, or phrases such as “当前分支”“当前代码”“当前未提交修改”.

Interpret direction in this order:

1. “从 A 到 B” and “A 升级到 B” mean A is the earlier state and B is the destination; “A 相比 B” means B is the earlier state and A is the destination.
2. “A 到当前未提交代码” compares the committed state at A with the current checkout, including intervening commits plus staged, unstaged, and untracked files.
3. For “A 和 B 之间” or “比较 A、B”, inspect Git ancestry first: the ancestor is the earlier state. If neither is an ancestor, use an unambiguous release/version order when available.
4. If direction still changes the meaning and cannot be inferred reliably, ask one concise question: “哪个是升级前版本，哪个是升级后版本？” Do not ask the user to rewrite the request as parameters.

If one of the two states is missing and cannot be inferred from an explicit phrase such as “当前未提交代码”, ask only for the missing state.

## Safety Contract

- Stay on the current checkout. Do not checkout, switch, reset, restore, stash, clean, fetch, pull, merge, or edit files.
- Preserve all dirty and untracked work. Read committed content with object-qualified commands such as `git show <ref>:<path>`.
- Do not use sub-agents, delegated workers, or parallel agent workflows.
- If `.codegraph/` exists, use CodeGraph before grep/find or broad source reads when understanding code behavior.
- State when refs are local snapshots; do not silently refresh remote refs.

## Evidence Workflow

1. Resolve each committed state with `git rev-parse --verify <ref>^{commit}` and record its OID. For two committed states, inspect ancestry, merge-base, `git diff <from-ref>..<to-ref>`, and left/right history. Use two-dot tree differences for the release result; use three-dot history only to understand divergence.
2. When the destination is the current uncommitted code, compare tracked content with `git diff <from-ref>` and enumerate untracked content with `git ls-files --others --exclude-standard`. Include relevant untracked files in the feature analysis.
3. Detect gitlinks with `git ls-tree`. For every changed or dirty submodule, compare the exact endpoint pointers inside that repository. When comparing with current code, also inspect the submodule's current `HEAD`, staged/unstaged changes, and untracked files. Never substitute the current submodule checkout for a committed endpoint pointer.
4. Build a feature inventory from commit history, changed paths, schemas, migrations, configuration, tests, runbooks, and product/module documentation. In easypos-hub, prefer `docs/memory/index.yaml`, `docs/modules/`, and `docs/modules-archive/` as routing evidence. Treat generated files, tests, plans, and docs as evidence; do not present them as product features.
5. Reconcile feature lineages across the comparison boundary. For each substantial destination-side feature family, inspect nearby earlier-state history and archived feature records for its root. When the destination completes or materially refactors an unreleased/pre-release feature whose initial merge landed in the earlier state, describe the complete final capability under **新增功能** and note the cross-boundary evidence in the comparison basis. This prevents an accidentally early merge from hiding a new release feature. If the evidence does not establish an unreleased lineage, classify only the destination-side delta as an improvement.
6. Describe the final destination behavior. Exclude intermediate designs, superseded contracts, reverted UI, merge noise, backup commits, and implementation-only churn. Prefer user and operator outcomes over class names or file lists.
7. Cross-check every claim against at least one code/config/schema artifact and one independent history, test, runbook, or module-memory artifact when available. Mark uncertain conclusions as inferences instead of presenting them as facts.
8. Derive an external QA checklist from the final feature inventory. Cover the main success path for each added or improved capability, regression checks for fixes, and compatibility or migration checks for upgrade risks. Add permission, failure/recovery, state-transition, platform, or historical-data cases when the change affects them.

## Output Contract

Write in the user's language; default to Chinese when unspecified. Lead with a short comparison basis containing the two interpreted states, resolved OIDs or current-worktree status, direction, topology, submodule coverage, and any cross-boundary feature lineage included.

Then emit this release-ready structure:

```markdown
# <新版版本号或比较目标>

## 新增功能
- <complete user-visible capability>

## 功能优化
- <meaningful improvement to an already released capability>

## 问题修复
- <observable corrected behavior>

## 升级注意事项
- <migration, configuration, compatibility, rollout, or rollback concern>

## 测试验证清单
- [ ] [P0] <功能或平台>：<外部测试人员执行的操作>；预期：<可观察结果>
```

Always emit all five sections; write `- 无` when one of the first four has no supported item. Group related App, service, admin, web, and integration changes into one product capability. Do not dump commits, line counts, generated-code volume, internal test names, or speculative marketing claims.

The **测试验证清单** is a forward-looking handoff for external testers, not evidence of checks already performed. Write actionable Markdown checkboxes, assign `P0`/`P1`/`P2` by release risk, and state both the tester action and observable expected result. Keep product wording and required setup visible; do not expose code symbols, test filenames, internal logs, credentials, or implementation details. Never phrase an item as “已验证”, “测试通过”, or CI/automated-test status.

## Natural-Language Examples

- `$generating-git-changelogs 帮我整理 1.0.0-rc 和 1.1.0-rc 之间的版本升级记录。`
- `$generating-git-changelogs 看一下 abc123 到 def456 这段提交做了哪些功能。`
- `$generating-git-changelogs 整理 release/1.0 到当前未提交代码的 ChangeLog。`

## Quick Reference

| Question | Evidence |
|---|---|
| What changed in the final trees? | `git diff --stat/--name-status <from-ref>..<to-ref>` |
| How are committed states related? | `git merge-base`, `git merge-base --is-ancestor`, left/right log |
| What changed in submodules? | `git ls-tree`, then nested `git log` and `git diff` using exact gitlink OIDs |
| What is currently uncommitted? | `git status --short`, `git diff`, `git diff --cached`, untracked-file listing |
| What is the final feature meaning? | Current source/schema/config plus tests, runbooks, and module memory |
| Was a feature root merged too early? | Source-side feature merge/docs plus target-side continuation/refactor evidence |

## Common Mistakes

- Comparing only commit subjects and missing behavior hidden in submodules.
- Reading the dirty checkout as though it were a committed target.
- Ignoring untracked files when the request includes current uncommitted code.
- Requiring users to learn parameter names instead of interpreting their natural-language request.
- Reporting migrations, generated code, tests, or documentation as features.
- Calling a cross-boundary unreleased feature an “optimization” because its first merge predates the earlier state.
- Reporting temporary behavior that was later reverted before the destination state.
- Assuming an accidental early merge without corroborating lineage evidence.
