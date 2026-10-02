---
name: merging-reviewed-worktree
description: Use when work completed inside a linked git worktree must land locally only after `$requesting-code-review` and `$archiving-module-memory`; commits the reviewed worktree changes, merges them into the chosen target branch, verifies the merged result, and cleans up the worktree.
---

# Merging Reviewed Worktree

## Overview

Land completed linked-worktree changes with a strict finishing sequence:

`requesting-code-review` -> `archiving-module-memory` -> commit in worktree -> merge into target branch -> verify merged result -> clean up worktree

**Core principle:** no commit, merge, or worktree cleanup before review and memory archival are complete.

**Announce at start:** "I'm using the merging-reviewed-worktree skill to land this reviewed worktree safely."

## When to Use

Use this skill when all of the following are true:

- The implementation lives in a linked git worktree
- The user wants a local merge back into a specific target branch
- The repo workflow expects code review before landing changes
- The repo workflow expects module-memory archival before the merge is finalized
- The worktree should be removed after a successful local merge

Do not use this skill when:

- The user wants a PR or push-first workflow instead of a local merge
- The work is not in a linked worktree
- The user is still iterating on implementation and has not entered the review/archive phase

For broader "what should we do next?" branch-completion choices, use `finishing-a-development-branch` instead.

## Required Preconditions

Before any `git add`, `git commit`, `git merge`, or worktree cleanup:

1. Confirm you are inside a linked worktree:

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
```

If `GIT_DIR == GIT_COMMON`, stop and use a non-worktree finishing flow instead.

2. Run fresh verification following `verification-before-completion`.
3. Confirm `$requesting-code-review` has already been used for the current implementation, and that any Critical or Important issues are fixed.
4. Confirm `$archiving-module-memory` has already been used when the change affects module memory, workflow memory, or durable process documentation in this repo.
5. Identify the target branch that should receive the merge.
6. Inspect both the worktree branch and the target checkout for unrelated local changes before touching Git state.

If any precondition is not satisfied, stop and tell the user exactly which required skill or decision comes next.

## The Process

### Step 1: Re-verify the Worktree State

Run the same test or verification command that proves the work is complete.

- Do not rely on old output
- Do not claim readiness without fresh evidence
- Record the exact command that passed

If verification fails, stop here and fix the issue before continuing.

### Step 2: Gate on Review Completion

Check whether the current thread already completed `$requesting-code-review`.

If not:

```text
Review has not been completed for this worktree yet.
Use `$requesting-code-review` first, resolve any Critical/Important findings, then return here.
```

Do not commit before this gate is satisfied.

### Step 3: Gate on Memory Archival

Check whether the current thread already completed `$archiving-module-memory` for any affected repo memory.

If not:

```text
Module memory has not been archived yet.
Use `$archiving-module-memory` next, then return here before commit/merge.
```

Do not commit before this gate is satisfied.

### Step 4: Commit the Worktree Changes

Inside the linked worktree:

1. Review `git status --short`
2. Stage only the intended files
3. Create a focused commit

Preferred pattern:

```bash
git add <explicit-file-list>
git commit -m "<focused message>"
```

Rules:

- Never use `git add .` in a dirty repo
- Never stage unrelated user changes
- Never amend an existing commit unless the user explicitly asks

### Step 5: Prepare the Target Checkout

Move to the checkout that owns the target branch and inspect its state:

```bash
git status --short --branch
```

If unrelated local changes would block the merge:

- Prefer a narrow stash of only the conflicting file(s)
- Restore the stash after the merge
- If the local changes are ambiguous or risky, stop and ask the user

Never overwrite unrelated changes just to force the merge through.

### Step 6: Merge Back to the Target Branch

Use a normal local merge that preserves branch history:

```bash
git merge --no-ff <worktree-branch>
```

If the merge conflicts:

- Resolve only the relevant conflicts
- Re-run verification after resolution
- Do not clean up the worktree until the merged result is proven good

### Step 7: Verify the Merged Result

Run fresh verification again on the target branch after the merge.

If the repo keeps linked worktrees inside a project-local `.worktrees/` directory and your test runner globs recursively, explicitly exclude worktree paths to avoid duplicate test discovery.

Example:

```bash
pnpm exec vitest run --exclude .worktrees/** <test-files>
```

If post-merge verification fails:

- Report the failure honestly with command output
- Do not delete the worktree
- Do not delete the worktree branch
- Fix forward or ask the user how they want to proceed

### Step 8: Clean Up the Worktree

Only after the merged result verifies cleanly:

```bash
git worktree remove <worktree-path>
git branch -d <worktree-branch>
git worktree prune
```

If you created a temporary stash while preparing the target checkout, restore it after cleanup.

### Step 9: Report Final State

Report:

- The commit created in the worktree
- The merge commit or fast-forward result on the target branch
- The verification command and whether it passed
- Whether the worktree and branch were deleted
- Any remaining top-level or parent-repo changes that were intentionally left outside this merge

## Nested Repo Reminder

In this monorepo, linked worktree code may live inside a nested repo such as `base-web`, while module-memory archives may live in the parent repo.

Do not assume a successful nested-repo merge also commits parent-repo documentation.

After the worktree merge, explicitly check whether any parent-repo files still need separate staging or commit, for example:

- `docs/modules/...`
- `docs/modules-archive/...`
- submodule pointer updates in the parent repo

## Common Mistakes

### Committing before review

- **Problem:** The work lands before issues are surfaced
- **Fix:** Require `$requesting-code-review` before any commit

### Skipping module-memory archival

- **Problem:** Durable repo knowledge diverges from current implementation
- **Fix:** Require `$archiving-module-memory` before merge finalization when repo memory changed

### Cleaning up too early

- **Problem:** Post-merge failures become harder to debug or recover
- **Fix:** Keep the worktree until merged verification passes

### Overly broad staging or stashing

- **Problem:** Unrelated local changes get swept into the merge flow
- **Fix:** Stage and stash only explicit files

### Forgetting parent-repo fallout

- **Problem:** Nested repo merge succeeds but parent repo still has unlanded docs or submodule pointers
- **Fix:** Always inspect both the nested repo and parent repo before declaring the flow complete

## Red Flags

**Never:**

- Commit before `$requesting-code-review`
- Merge before `$archiving-module-memory` when memory changed
- Claim success without fresh verification
- Delete the worktree before merged verification passes
- Use destructive Git commands to bulldoze through conflicts

**Always:**

- Verify fresh before commit and after merge
- Stage only intended files
- Preserve unrelated local changes
- Report any parent-repo leftovers explicitly

## Integration

**Called after:**

- `$requesting-code-review`
- `$archiving-module-memory`

**Pairs with:**

- `using-git-worktrees` - creates the isolated workspace this skill closes out
- `verification-before-completion` - required before readiness claims, commits, and merge success claims
- `finishing-a-development-branch` - use that skill instead when the user wants broader completion choices rather than this strict reviewed-merge path
