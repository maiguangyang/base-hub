---
name: recall-project-memory
description: Use when long conversations, context compaction, resumed tasks, hallucination risk, unarchived current changes, or explicit requests such as “唤起记忆”, “当前方案”, “项目记忆”, “宪章约束”, “上下文过大” require Codex to rebuild the relevant easypos-hub project context, archive unarchived current work first when needed, then recover the active plan, live workspace state, module memory, and constitution constraints before continuing.
---

# Recall Project Memory

## Overview

Use this skill to rebuild a compact, evidence-backed working memory for `easypos-hub` before continuing a task. It prevents long-context drift by loading the repository's canonical rule and memory entrypoints in the intended order, then summarizing only the facts needed for the current work.

This skill is a context-recovery gate. After the memory snapshot is rebuilt, continue with any other task-specific skill that applies.

## Core Rule

Prefer small, authoritative reads over broad history scans. Mark facts as verified only when they come from files or command output in the current workspace. Treat older conversation context as a hint, not truth.

Current workspace state is fresher than long-term memory. If relevant implementation changes exist and `docs/modules/*` / `docs/modules-archive/*` have not yet been updated by `$archiving-module-memory`, invoke `$archiving-module-memory` before continuing the memory recall. Until that archival step completes, treat module memory as stale baseline only.

## Workflow

### 1. Establish The Active Workspace

Run light inspection commands first:

```bash
pwd
git status --short
git rev-parse --show-toplevel
git branch --show-current
```

If the workspace is a linked worktree, resolve the active tracker with:

```bash
git rev-parse --git-path codex/task.md
```

Use that worktree-local tracker if it exists. Otherwise use `docs/plans/task.md`.

Do not revert or overwrite dirty files. Treat existing changes as user or prior-agent work unless the user explicitly says otherwise.

### 2. Detect Unarchived Current Work

Before reading long-term module memory, inspect the live workspace for work that may not have been archived yet:

```bash
git status --short
git diff --name-only
git diff --stat
```

Use freshness precedence in this order:

1. The user's latest request in the current turn
2. Dirty and untracked workspace files
3. The active tracker and current plan files
4. Newly edited module/archive docs that are still uncommitted
5. `docs/modules/<module>.md` brief files
6. `docs/modules-archive/...`, older plans, and older conversation context

If dirty or untracked files match the active topic, read the relevant source files and `git diff -- <path>` before trusting module memory. If new archive entries or module brief edits are present but uncommitted, treat them as current in-progress memory, not as stale noise.

#### Archiving Gate

Before loading `docs/modules/<module>.md` as current memory, decide whether the active work has already been archived:

- If relevant source/runtime/test files changed and there are matching current module brief or `docs/modules-archive/<module>/YYYY-MM-DD-*.md` changes, continue; use those uncommitted memory docs as current in-progress memory.
- If relevant source/runtime/test files changed and there are no matching module/archive memory changes, stop this skill's normal workflow and invoke `$archiving-module-memory` first.
- After `$archiving-module-memory` finishes, resume this skill from Step 1 so the snapshot includes the newly updated memory.
- If `$archiving-module-memory` cannot run because the current方案, completed plan, or affected modules cannot be identified, ask one concise clarification and do not treat old module memory as current truth.

This gate applies before any decision that depends on long-term module memory. Do not skip it just because the user asked for fast recall.

When module memory contradicts live workspace evidence, write the snapshot as:

```markdown
- Module memory baseline:
- Live workspace override:
- Archiving status: archived / pending / blocked
```

Do not conclude that old module memory is the current truth until you have checked for unarchived work.

### 3. Load Constitution And Rule Entrypoints

Read these in order:

1. `.agent/memory/constitution-core.md`
2. The relevant `.agent/memory/rules-core/*.md` files
3. `.agent/memory/README.md` only if the memory layout is unclear

Choose rules-core files by scope:

- `base-engine`, GraphQL, database, service, billing, realtime backend, conversation runtime -> `.agent/memory/rules-core/base-engine.md`
- `base-app`, Flutter, Riverpod, l10n, app UI, app auth/home/settings -> `.agent/memory/rules-core/base-app.md`
- `base-web`, Astro, React admin, web GraphQL -> `.agent/memory/rules-core/base-web.md`
- `adk-go`, Google ADK, Hermes-like runtime, MCP Toolset, agent runtime internals -> `.agent/memory/rules-core/adk-go.md`

Open `.agent/memory/constitution.md` or `.agent/memory/rules/*.md` only when the core entry is insufficient for the active decision.

### 4. Recover The Current Plan

Read the active tracker:

- Worktree-local: output of `git rev-parse --git-path codex/task.md`
- Default workspace: `docs/plans/task.md`

Then identify the current方案 or implementation thread:

- Prefer explicit filenames or topics mentioned by the user.
- If absent, list recent non-archive plan files under `docs/plans/`.
- Read only the most relevant design, design-review, or implementation-plan file.
- If multiple plausible plans remain, ask one concise clarifying question before making edits.

Do not treat `docs/plans/task.md` as long-term memory. It is an active-window tracker.

### 5. Route Module Memory

Read `docs/memory/index.yaml` before opening module files. Match by:

- User request terms
- Active plan title and filenames
- Dirty paths from `git status --short`
- Code paths already identified during the task

Read the matching `docs/modules/<module>.md` brief files first. Keep the initial set to one to three modules unless the task is truly full-stack.

Read module memory as baseline after live workspace inspection. If live changes are more recent than the module brief, annotate the difference instead of silently overwriting the active方案 with older memory.

If `docs/memory/index.yaml` is incomplete for the visible module set, inspect `docs/modules/` filenames and choose the closest brief by module name. Open `docs/modules-details/...` only when the brief does not answer the active question.

Leave `docs/modules-archive/...` closed unless the task depends on historical evolution, regressions, or recently archived work.

### 6. Optional Deepening

Use deeper sources only after the entrypoints above:

- `docs/modules-details/...` for current detailed domain facts
- `docs/modules-archive/...` for dated history or regressions
- `git log --oneline -- <path>` for recent file evolution
- `rg` for exact symbols, errors, tests, or plan topic names

Do not broaden the search just to feel safer. Broaden only to answer a concrete uncertainty.

### 7. Produce A Context Snapshot

Before continuing substantial work, write a compact snapshot in the conversation or your working notes:

```markdown
**Context Snapshot**

- Current task /方案:
- Active plan:
- Live workspace overrides:
- Archiving status:
- Constitution constraints:
- Relevant module memory:
- Dirty workspace facts:
- Verified from files:
- Inferences / unknowns:
- Next skill or next action:
```

Keep the snapshot short. Its job is to anchor decisions, not preserve every detail.

## Hallucination Guardrails

- Never rely on memory when a local file can answer the question.
- If file evidence conflicts with prior conversation, trust the file and call out the conflict.
- Separate verified facts from inferences.
- Use exact filenames, module IDs, dates, and command outputs when summarizing.
- Ask a concise question when the current方案 cannot be identified after reading the tracker and recent plans.
- Do not paste secrets, tokens, runtime headers, or private credentials into long-lived docs or summaries.
- After rebuilding context, continue with the required task-specific skill: debugging, planning, implementation, review, replay, Figma, or UI/UX as appropriate.
