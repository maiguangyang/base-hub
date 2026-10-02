# Local Workflow Overview

This document summarizes the current local `.agent` workflow for design, planning, implementation, review, archival, and branch completion.

## Scope Routing

Classify work before entering the full workflow.

### Lightweight Change

Use for an explicit, low-risk change that follows an existing pattern and requires no architectural, migration, compatibility, permission, security, billing, or irreversible-operation decision. Touching multiple files or generated artifacts does not by itself make a change substantial.

```text
inspect the existing pattern
  ↓
implement the minimal change
  ↓
run focused verification
  ↓
report concisely
```

Do not require a design document, implementation plan, formal proposal review, module archive, or branch-completion ceremony for this path unless the user explicitly invokes the corresponding skill or inspection reveals a concrete reason to escalate.

### Substantive Change

Use the primary flow below when the work contains real product or architecture decisions, ambiguous requirements, migration or compatibility choices, high-risk boundaries, or irreversible effects.

## Primary Flow For Substantive Changes

```text
brainstorming
  ↓
using-git-worktrees (optional; ask once, only if the user opts in)
  ↓
writing-plans
  ↓
subagent-driven-development (when the user explicitly wants delegated workers)
or
executing-plans (default local single-flow execution)
  ↓
test-driven-development (implementation discipline during execution, not a separate gate)
  ↓
requesting-code-review
  ↓
archiving-module-memory
  ↓
finishing-a-development-branch
```

## Stage Meanings

### 1. `$brainstorming`

Use for design and requirement shaping before substantive implementation. Its lightweight branch handles clear mechanical changes without entering the rest of the primary flow.

- Loads constitution and relevant rule files first
- Reads module memory before proposing designs
- Produces or validates the design direction
- Hands off to `$writing-plans` as the next manual step

### 2. `$using-git-worktrees`

Optional workspace isolation step.

- Not mandatory in this repo
- Ask once before substantial implementation or plan execution
- Only use if the user explicitly opts in
- If the user declines, continue in the current workspace and do not ask again unless they later reverse that choice

### 3. `$writing-plans`

Turns the approved design or requirements into an implementation plan.

- Writes plan files under `docs/plans/`
- Defines exact files, steps, tests, and verification commands
- Chooses the execution route:
  - `$subagent-driven-development` if the user explicitly wants delegated workers
  - `$executing-plans` otherwise

### 4A. `$subagent-driven-development`

Use only when the user explicitly asks for delegated worker execution.

- Dispatches one focused worker per task
- Runs spec review before code quality review
- Reuses the same task loop until the task is approved
- Keeps changes local for human review

### 4B. `$executing-plans`

Default local execution path when the user has not explicitly asked for delegated workers.

- Executes plan tasks in local single-flow mode
- Uses batch checkpoints instead of delegated worker orchestration
- Stops for blockers or failed verification
- Hands off to `$requesting-code-review` after implementation is complete

### 5. `$test-driven-development`

Implementation discipline that applies during execution.

- Not usually a standalone handoff point in the main flow
- Guides how new behavior and bug fixes should be implemented
- Should be followed by implementers during `$subagent-driven-development`, `$executing-plans`, or `$single-flow-task-execution`

### 6. `$requesting-code-review`

Structured review pass before completion.

- Review after major implementation work
- Fix critical and important findings before proceeding
- If the work is ready, the next step is `$archiving-module-memory`

### 7. `$archiving-module-memory`

Updates long-term project memory after the implementation is accepted.

- Updates `docs/modules/<module>.md`
- Writes dated archive entries under `docs/modules-archive/<module>/`
- Keeps the changes local for human review

### 8. `$finishing-a-development-branch`

Final branch completion and cleanup workflow.

- Verifies tests before offering completion options
- Presents merge / PR / keep / discard choices
- Cleans up worktrees only when appropriate

## Decision Rules

### Execution Route Selection

- If the user explicitly wants delegated workers, subagents, or parallel-agent-style execution:
  - use `$subagent-driven-development`
- Otherwise:
  - use `$executing-plans`

### Worktree Selection

- If the user opts into workspace isolation:
  - use `$using-git-worktrees`
- Otherwise:
  - stay in the current workspace

### Review and Completion

- Require `$requesting-code-review` for substantive or high-risk implementation, or when the user explicitly requests it; focused self-review and verification are sufficient for lightweight changes.
- Use `$archiving-module-memory` when accepted work changes durable module constraints or architecture, or when the user explicitly requests it; do not archive every mechanical addition.
- Use `$finishing-a-development-branch` only when branch integration or cleanup is actually needed.

## Practical Notes

- `$single-flow-task-execution` remains a local execution helper pattern and is still used internally by some execution-oriented skills.
- `$subagent-driven-development` is now part of the local flow, but it is not the default route.
- The repository's current default implementation path is still conservative:
  - current workspace by default
  - local single-flow execution by default
  - delegated workers only when explicitly requested
