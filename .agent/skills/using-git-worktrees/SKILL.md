---
name: using-git-worktrees
description: Use when the user wants isolation from the current workspace, or opts in after a one-time reminder before implementation - creates isolated git worktrees on consistently namespaced worktrees/task-slug branches, migrates current-task artifacts from the coordinating parent repository, and verifies safe isolation
---

# Using Git Worktrees

## Overview

Git worktrees create isolated workspaces sharing the same repository, allowing work on multiple branches simultaneously without switching.

**Core principle:** Standardized `worktrees/<task-slug>` branch naming + systematic directory selection + task-artifact migration + safety verification = reliable isolation.

**Announce at start:** "I'm using the using-git-worktrees skill to set up an isolated workspace."

## Trigger Policy

Use this skill only after one of these is true:

- The user explicitly asks for an isolated workspace or mentions `git worktree`
- Another execution workflow asks once whether the user wants isolation and the user says yes
- The user explicitly chooses `$using-git-worktrees` as the planning-entry option after brainstorming or proposal review

Do not trigger this skill just because implementation is starting. A brief one-time reminder is allowed; repeated nudging is not.

## Step 0: Detect Existing Isolation

Before creating anything, check whether you are already inside an isolated linked worktree or an externally managed workspace.

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
BRANCH=$(git branch --show-current)
git rev-parse --show-superproject-working-tree 2>/dev/null
```

Submodule guard: `GIT_DIR != GIT_COMMON` can also be true inside git submodules. If `git rev-parse --show-superproject-working-tree` returns a path, treat this as a normal repo checkout, not as an already isolated worktree.

If `GIT_DIR != GIT_COMMON` and you are not inside a submodule, report the existing isolated workspace and skip creation:

- If `BRANCH` is set: "Already in isolated workspace at `<path>` on branch `<branch>`."
- If `BRANCH` is empty: "Already in isolated workspace at `<path>` (detached HEAD, externally managed). Branch creation/cleanup may need platform-native controls."

If this skill was chosen as the planning-entry option after brainstorming or proposal review, treat the existing isolated workspace as a successful setup and continue with the same automatic handoff rules described below.

## Step 0.5: Identify the Coordinating Repository and Task Files

Before creating a worktree, identify where task-level artifacts currently live. In a repository with submodules, the superproject is the coordinating repository and commonly owns design documents, plans, review notes, module memory, and final submodule gitlink changes.

```bash
repo_root=$(git rev-parse --show-toplevel)
super_root=$(git rev-parse --show-superproject-working-tree 2>/dev/null)

if [ -n "$super_root" ]; then
  coordination_root=$super_root
else
  coordination_root=$repo_root
fi

git -C "$coordination_root" status --short --untracked-files=all
```

Build an explicit task-file manifest before creating anything. Include only files already created or changed for the current task, such as:

- design, design-review, specification, or implementation-plan documents
- current-task notes or memory updates intended to be committed
- implementation or test files changed before isolation began
- participating submodule changes and the parent repository's eventual gitlink updates

Exclude unrelated pre-existing changes and ambiguous files. If ownership is unclear, report the candidate path and ask before moving it. Record each selected path's repository, tracked/untracked state, and source location. Do not rely on filename globbing or move the entire dirty worktree.

If the task spans a parent repository and submodules, create a coordinated isolation layout for all participating repositories from the user-selected base branches. The parent worktree is the task's primary workspace; task documents and final gitlink changes belong there. Do not leave task-owned parent files in the original checkout merely because implementation happens in a submodule.

## Directory Selection Process

Follow this priority order:

### 1. Check Existing Directories

```bash
# Check in priority order
ls -d .worktrees 2>/dev/null     # Preferred (hidden)
ls -d worktrees 2>/dev/null      # Alternative
```

**If found:** Use that directory. If both exist, `.worktrees` wins.

### 2. Check AGENTS.md

```bash
grep -i "worktree.*director" .agent/AGENTS.md 2>/dev/null
```

**If preference specified:** Use it without asking.

### 3. Ask User

If no directory exists and no AGENTS.md preference:

```
No worktree directory found. Where should I create worktrees?

1. .worktrees/ (project-local, hidden)
2. ~/.config/superpowers/worktrees/<project-name>/ (global location)

Which would you prefer?
```

If the user gave a clear preference earlier in the conversation or project instructions, honor it without asking again. If there is no preference and interaction is not possible, default to `.worktrees/` at the project root.

## Safety Verification

### For Project-Local Directories (.worktrees or worktrees)

**MUST verify directory is ignored before creating worktree:**

```bash
# Check if directory is ignored (respects local, global, and system gitignore)
git check-ignore -q .worktrees 2>/dev/null || git check-ignore -q worktrees 2>/dev/null
```

**If NOT ignored:**

Per Jesse's rule "Fix broken things immediately":

1. Add appropriate line to .gitignore
2. leave the .gitignore change local for human review
3. Proceed with worktree creation only after the ignore-file update is preserved locally

**Why critical:** Prevents accidentally committing worktree contents to repository.

### For Global Directory (~/.config/superpowers/worktrees)

No .gitignore verification needed - outside project entirely.

## Creation Steps

### 1. Derive the Required Branch Name

Every branch created by this skill **MUST** use exactly this format:

```text
worktrees/<task-slug>
```

Derive one non-empty, lowercase kebab-case `task_slug` from the task name. If the requested name already starts with `worktrees/`, strip that prefix before normalizing so it is never duplicated. Convert spaces, underscores, slashes, and other separators to single hyphens, then remove leading or trailing hyphens.

```bash
raw_name=${BRANCH_NAME:-$TASK_NAME}
task_slug=${raw_name#worktrees/}
task_slug=$(printf '%s' "$task_slug" | tr '[:upper:]_ ' '[:lower:]--' | sed -E 's#[^a-z0-9-]+#-#g; s#-+#-#g; s#^-|-$##g')

if [ -z "$task_slug" ]; then
  echo "Cannot create worktree: task slug is empty" >&2
  exit 1
fi

BRANCH_NAME="worktrees/$task_slug"
git check-ref-format --branch "$BRANCH_NAME"
```

Examples:

| Requested task or branch | Created branch |
| ------------------------ | -------------- |
| `auth`                   | `worktrees/auth` |
| `Feature/Auth Flow`      | `worktrees/feature-auth-flow` |
| `worktrees/billing-fix`  | `worktrees/billing-fix` |

Do not create `feature/*`, `fix/*`, bare task-name branches, or any other namespace. The `worktrees/` prefix is a branch namespace; it does not determine the filesystem directory.

### 2. Detect Project Name

```bash
project=$(basename "$(git rev-parse --show-toplevel)")
```

### 3. Create Worktree

```bash
# Determine full path
case $LOCATION in
  .worktrees|worktrees)
    path="$LOCATION/$task_slug"
    ;;
  ~/.config/superpowers/worktrees/*)
    path="~/.config/superpowers/worktrees/$project/$task_slug"
    ;;
esac

# Create worktree with new branch
git worktree add "$path" -b "$BRANCH_NAME"
cd "$path"
```

If `git worktree add` fails with a sandbox or permission error, report that isolation was blocked and work in the current directory only after the user confirms that fallback.

### 3.1 Migrate Existing Task Files

Immediately after creating all required worktrees, migrate every path in the task-file manifest into the matching repository worktree while preserving its relative path.

For untracked task files:

1. Refuse to overwrite an existing destination file silently.
2. Create the destination parent directory.
3. Move the file from the original checkout into the new worktree.
4. Verify the destination exists and the original path no longer exists.

For tracked task modifications:

1. Transfer both staged and unstaged changes without losing either state; use Git patches or an equivalent lossless method.
2. Verify the destination diff matches the source diff.
3. Restore only the transferred task paths in the original checkout after verification.
4. Never restore, reset, stash, or clean unrelated paths.

When a destination already contains a tracked version of the same path, compare it with the source and preserve both versions until their differences are understood. Do not delete the source copy merely because a same-named file exists on the base branch.

After migration, run status in the original checkout and every new worktree:

```bash
git -C "$coordination_root" status --short --untracked-files=all
git -C "$path" status --short --untracked-files=all
```

The original checkout must no longer list any manifest path. Every manifest path must appear in, or be represented by a diff in, the correct isolated worktree. If either check fails, stop and reconcile before implementation begins.

### 3.5 Establish Worktree-Local Tracker

Inside a linked worktree, do **not** keep using `<project-root>/docs/plans/task.md` as the active tracker.

Use a worktree-local tracker instead:

```bash
tracker=$(git rev-parse --git-path codex/task.md)
archive_dir=$(git rev-parse --git-path codex/archive)
mkdir -p "$(dirname "$tracker")" "$archive_dir"
```

**Why:** `git rev-parse --git-path ...` stores the tracker inside the linked worktree's git dir, so it stays isolated from other worktrees and cannot create branch merge conflicts.

The local tracker supplements task documents; it does not replace design or plan files that are intended to be committed. Commit-bound documents must live inside the parent task worktree.

### 4. Run Project Setup

Auto-detect and run appropriate setup:

```bash
# Node.js
if [ -f package.json ]; then npm install; fi

# Rust
if [ -f Cargo.toml ]; then cargo build; fi

# Python
if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
if [ -f pyproject.toml ]; then poetry install; fi

# Go
if [ -f go.mod ]; then go mod download; fi
```

### 5. Verify Clean Baseline

Run tests to ensure worktree starts clean:

```bash
# Examples - use project-appropriate command
npm test
cargo test
pytest
go test ./...
```

**If tests fail:** Report failures, ask whether to proceed or investigate.

**If tests pass:** Report ready.

## Automatic Continuation Into Writing Plans

When the user explicitly chose `$using-git-worktrees` from a design-to-planning handoff:

1. Finish the isolation workflow first.
2. If setup succeeds and no new user decision is pending, automatically continue into `$writing-plans` in the resulting workspace.
3. Do not ask the user to choose again between `$using-git-worktrees` and `$writing-plans`; that choice was already made.
4. If setup needs a new decision point, such as choosing a directory location or deciding what to do about failing baseline tests, stop and wait for the user before continuing.
5. If setup fails or remains unresolved, do not continue into `$writing-plans`.

### 6. Report Location

```
Worktree ready at <full-path>
Local tracker at <tracker-path>
Tests passing (<N> tests, 0 failures)
Next approved step:
- continue with the workflow that requested isolation
- if this was a planning-entry handoff, continue into `$writing-plans`
```

## Quick Reference

| Situation                  | Action                              |
| -------------------------- | ----------------------------------- |
| Creating any new branch    | Require `worktrees/<task-slug>`     |
| Requested name has prefix  | Strip once, then rebuild the canonical branch name |
| `.worktrees/` exists       | Use it (verify ignored)             |
| `worktrees/` exists        | Use it (verify ignored)             |
| Both exist                 | Use `.worktrees/`                   |
| Neither exists             | Check `.agent/AGENTS.md` → Ask user |
| Directory not ignored      | Add to .gitignore + leave local     |
| Tests fail during baseline | Report failures + ask               |
| No package.json/Cargo.toml | Skip dependency install             |
| Need active tracker        | Use `git rev-parse --git-path codex/task.md` |
| Task files predate worktree | Manifest, migrate, and verify them          |
| Parent + submodule task     | Isolate both; keep task docs in parent worktree |
| Already in linked worktree | Skip creation and continue setup    |
| Detached HEAD workspace    | Treat as externally managed; do not create nested worktree |

## Common Mistakes

### Skipping ignore verification

- **Problem:** Worktree contents get tracked, pollute git status
- **Fix:** Always use `git check-ignore` before creating project-local worktree

### Skipping existing-isolation detection

- **Problem:** Creates nested worktrees or fights a harness-managed workspace
- **Fix:** Always run Step 0 before creating anything

### Assuming directory location

- **Problem:** Creates inconsistency, violates project conventions
- **Fix:** Follow priority: existing > `.agent/AGENTS.md` > ask

### Creating an inconsistently named branch

- **Problem:** Bare names or mixed namespaces such as `feature/*` and `fix/*` make worktree branches difficult to identify and manage
- **Fix:** Normalize the task slug first and create only `worktrees/<task-slug>`; keep the filesystem path based on `<task-slug>` alone

### Proceeding with failing tests

- **Problem:** Can't distinguish new bugs from pre-existing issues
- **Fix:** Report failures, get explicit permission to proceed

### Hardcoding setup commands

- **Problem:** Breaks on projects using different tools
- **Fix:** Auto-detect from project files (package.json, etc.)

### Continuing to use `docs/plans/task.md` inside a linked worktree

- **Problem:** Multiple worktrees update different branch copies of the same tracked file and create avoidable merge conflicts
- **Fix:** Switch to `git rev-parse --git-path codex/task.md` immediately after entering the linked worktree

### Leaving pre-worktree task files in the original checkout

- **Problem:** Design documents, review notes, or early task changes remain as untracked or dirty files after the isolated work is merged
- **Fix:** Inventory task files before creation, migrate them into the matching parent or submodule worktree, and verify both source and destination status before implementation

### Moving every dirty file into the worktree

- **Problem:** Unrelated user work becomes mixed into the task branch or is accidentally removed from the original checkout
- **Fix:** Use an explicit task-file manifest; move only paths whose ownership is clear

## Example Workflow

```
You: I'm using the using-git-worktrees skill to set up an isolated workspace.

[Check .worktrees/ - exists]
[Verify ignored - git check-ignore confirms .worktrees/ is ignored]
[Inventory current-task files in the parent repository and participating submodules]
[Normalize requested task `Auth Flow` to slug `auth-flow`]
[Create worktree: git worktree add .worktrees/auth-flow -b worktrees/auth-flow]
[Move the manifest paths into the corresponding worktrees and verify source/destination status]
[Run npm install]
[Run npm test - 47 passing]

Worktree ready at /Users/jesse/myproject/.worktrees/auth-flow on branch worktrees/auth-flow
Local tracker at /Users/jesse/myproject/.git/worktrees/auth-flow/codex/task.md
Tests passing (47 tests, 0 failures)
Next approved step:
- continue with the workflow that requested isolation
- if this was a planning-entry handoff, continue into `$writing-plans`
```

## Red Flags

**Never:**

- Create a worktree when Step 0 detects an existing linked worktree
- Treat a detached HEAD workspace as a normal named branch
- Create a new branch outside the `worktrees/<task-slug>` namespace
- Use the full `worktrees/<task-slug>` branch name as the directory suffix
- Create worktree without verifying it's ignored (project-local)
- Leave current-task files or parent-repository task documents in the original checkout
- Bulk-move, stash, restore, reset, or clean unrelated dirty files
- Delete a source file before verifying its destination copy or diff
- Skip baseline test verification
- Proceed with failing tests without asking
- Assume directory location when ambiguous
- Skip `.agent/AGENTS.md` check

**Always:**

- Follow directory priority: existing > `.agent/AGENTS.md` > ask
- Normalize every new branch to `worktrees/<task-slug>` and validate it with `git check-ref-format`
- Verify directory is ignored for project-local
- Auto-detect and run project setup
- Verify clean test baseline

## Integration

**Called by:**

- **single-flow-task-execution** - OPTIONAL after a one-time reminder if the user opts into isolation
- **executing-plans** - OPTIONAL after a one-time reminder if the user opts into isolation
- **brainstorming / evaluating-technical-proposals handoff** - OPTIONAL when the user explicitly chooses isolated planning setup before `$writing-plans`
- Any skill needing isolated workspace after explicit user approval

**Pairs with:**

- **finishing-a-development-branch** - REQUIRED for cleanup after work complete
- **merging-reviewed-worktree** - OPTIONAL stricter closeout flow when isolated work must wait for `$requesting-code-review` and `$archiving-module-memory` before commit, merge, and cleanup
