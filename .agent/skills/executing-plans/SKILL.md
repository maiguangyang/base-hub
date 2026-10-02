---
name: executing-plans
description: Use when you have a written implementation plan and need to execute it in Codex single-flow mode without delegated workers
---

# Executing Plans

## Overview

Load plan, review critically, execute tasks sequentially, and report progress with verification evidence.

**Core principle:** Sequential execution to completion, with checkpoints only when blocked or when the user explicitly asks for an intermediate review.
**Entrypoint principle:** This is the standard execution entrypoint for local single-flow execution when the user has not explicitly asked for delegated workers.

**Announce at start:** "I'm using the executing-plans skill to implement this plan."

## The Process

### Step 1: Load and Review Plan

1. Read plan file
2. Review critically - identify any questions or concerns about the plan
3. If the user explicitly wants delegated workers or subagent-style execution, stop here and use `$subagent-driven-development` instead of this skill
4. If the user has not already chosen or declined `$using-git-worktrees`, and you are not already inside an isolated linked worktree, ask once whether they want isolated workspace execution before starting. If they opt in, use `.agent/skills/using-git-worktrees/SKILL.md` first. If they decline, continue here and do not ask again in downstream execution skills.
5. If concerns: Raise them with your human partner before starting. Batch all pre-flight concerns together instead of interrupting one at a time.
6. If no concerns: follow the single-flow execution model from `.agent/skills/single-flow-task-execution/SKILL.md`
7. Resolve the active-window tracker path and update it before proceeding. Use `<project-root>/docs/plans/task.md` in the default workspace, or `git rev-parse --git-path codex/task.md` inside a linked worktree.

Pre-flight concerns include contradictory task requirements, missing interfaces between tasks, instructions that would violate repo rules, or plan-mandated behavior that conflicts with the review rubric.

### Step 2: Execute Tasks

**Default: Execute all tasks in plan order**

For each task:

1. Mark as in_progress
2. **Check for Figma annotation** — if the task has a `> **Figma:**` block, use `$figma` and follow its required flow (`get_figma_data` → `download_figma_images` if needed) for the linked node(s) **before** writing any code
3. Follow each step exactly (plan has bite-sized steps)
4. Run verifications as specified
5. Mark as completed

### Step 3: Report Progress

At logical checkpoints and after all tasks complete:

- Show what was implemented
- Show verification output
- If work is still in progress because of a blocker or an explicit user-requested pause, say: "Ready for feedback."

### Step 4: Continue Execution

Based on feedback or the current execution state:

- Apply changes if needed
- Continue with the next remaining task
- Repeat until complete

### Step 5: Complete Development

After all tasks complete and verified:

- leave changes local for human review
- Do NOT automatically run git add or git commit
- Stop execution and instruct the user to use `$requesting-code-review` to proceed to the review phase.
- Do NOT invoke `finishing-a-development-branch` or `$archiving-module-memory` until the user actively uses `$requesting-code-review` and approves the review.

## When to Stop and Ask for Help

**STOP executing immediately when:**

- Hit a blocker mid-task (missing dependency, test fails, instruction unclear)
- Plan has critical gaps preventing starting
- You don't understand an instruction
- Verification fails repeatedly

**Ask for clarification rather than guessing.**

## When to Revisit Earlier Steps

**Return to Review (Step 1) when:**

- Partner updates the plan based on your feedback
- Fundamental approach needs rethinking

**Don't force through blockers** - stop and ask.

## Remember

- Review plan critically first
- Follow plan steps exactly
- Don't skip verifications
- Reference skills when plan says to
- Treat the resolved active-window tracker as active-window state only; archive stale execution context under `docs/plans/archive/` in the default workspace, or `git rev-parse --git-path codex/archive` inside a linked worktree
- Do not treat historical plan files as default long-term memory
- Continue through the full plan unless blocked or the user asks to pause for review
- Stop when blocked, don't guess
- Never start implementation on main/master branch without explicit user consent
- Use one clearly scoped task brief for coding tasks; use browser-specific tooling only for browser tasks

## Integration

**Required workflow skills:**

- **`.agent/skills/using-git-worktrees/SKILL.md`** - REQUIRED one-time decision check unless the user already chose or declined isolation, or you are already inside an isolated linked worktree. If the user opts in, load this skill before execution.
- **`.agent/skills/writing-plans/SKILL.md`** - Creates the plan this skill executes
- **`.agent/skills/single-flow-task-execution/SKILL.md`** - REQUIRED: Enforce single-flow execution with two-stage review
- **`.agent/skills/figma/SKILL.md`** - CONDITIONAL: Auto-loaded when a task carries a `> **Figma:**` annotation (see Step 2)
- **`$requesting-code-review`** - Complete development after all tasks

**Alternative execution route:**

- **`$subagent-driven-development`** - Use this instead when the user explicitly wants delegated worker execution in the current session
