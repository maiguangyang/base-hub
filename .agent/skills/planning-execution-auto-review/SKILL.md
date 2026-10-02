---
name: planning-execution-auto-review
description: Use when implementation requirements must be completed autonomously to final code-review approval, especially when no written plan exists yet.
---

# Planning, Execution, and Automatic Review

## Overview

Turn requirements into the smallest sufficient written plan, then delegate that plan to `$executing-plans-auto-review` for autonomous execution and approval.

**Core principle:** This skill owns planning only. `$executing-plans-auto-review` is the single source of truth for execution, review, repair, verification, blockers, and approval.

**REQUIRED SUB-SKILLS:** Use `$writing-plans`, then `$executing-plans-auto-review`, in exactly this sequence:

```text
$writing-plans -> $executing-plans-auto-review
```

Do not directly orchestrate `$executing-plans` or `$requesting-code-review` here. They are downstream responsibilities of `$executing-plans-auto-review`.

## Authorization Contract

- Treat invocation as authorization to write or revise the plan and complete the entire delegated `$executing-plans-auto-review` workflow without asking again.
- Treat invocation as explicit consent to work on the current branch, including `main` or `master`, while keeping every change local and unstaged. Continue in the current workspace unless already inside a linked worktree; do not pause for branch or worktree confirmation.
- Override only `$writing-plans`' required execution-route header with `EXECUTION ROUTE: Use $planning-execution-auto-review`. Keep every other required plan field. At its normal handoff, continue automatically into `$executing-plans-auto-review`.
- At `$writing-plans`' normal handoff, invoke `$executing-plans-auto-review` with the exact saved plan path. Treat this parent invocation as explicit invocation and standing authorization for that child workflow.
- Do not stage, commit, merge, push, archive, perform destructive operations, or change external systems unless separately authorized.

## No-Overdesign Planning Gate

Always run `$writing-plans` first. If a plan already exists, audit and simplify it before execution.

The plan MUST satisfy all of these constraints before implementation begins:

- Plan only the requested behavior and the verification necessary to prove it.
- Use one cohesive task for a simple change. Split only when separate deliverables can be independently implemented and reviewed.
- Prefer the smallest justified change surface and existing repository patterns.
- State explicit non-goals and identify unrelated behavior that must remain unchanged.
- Do not add abstractions, layers, interfaces, frameworks, dependencies, configuration, migrations, public APIs, or extension points unless a current acceptance criterion or mandatory repository rule directly requires them.
- Do not include speculative future-proofing, unrelated refactors, optional polish, or tasks without a direct requirement mapping.
- Keep the plan executable and precise, but never turn documentation detail into extra implementation scope.

If any planned item cannot be justified by a current requirement, remove it. A plan that fails this gate MUST NOT be delegated to `$executing-plans-auto-review`.

## Workflow

1. **Plan** — Run `$writing-plans` on the requirements or existing plan. Apply the no-overdesign gate, complete its self-review, save the plan, and continue without a human handoff.
2. **Delegate** — Load the current `$executing-plans-auto-review` skill and invoke it with the saved plan. Follow that skill completely without restating, replacing, or weakening any downstream rule.
3. **Finish** — Stop successfully only after `$executing-plans-auto-review` reaches its own approval gate. Leave changes local and report the plan path, planning scope, review rounds, and verification evidence. Do not enter branch finishing automatically.

## Delegation Contract

- `$executing-plans-auto-review` exclusively defines downstream execution order, read-only review boundaries, full-fix behavior, verification, re-review, approval, and blocker handling.
- Do not copy those downstream rules into this skill. Refer to the child skill so future changes apply automatically.
- Do not stop at `$writing-plans`' normal human handoff. Invoke the child workflow immediately with the saved plan.
- If the child reports a genuine blocker, propagate its exact evidence and requested decision; do not reinterpret the blocked workflow as approved.

## Genuine Blockers

Before delegation, human intervention is not a normal checkpoint. Stop only when planning cannot continue safely or correctly:

- requirements are contradictory or materially ambiguous;
- required planning inputs or repository rules are unavailable with no safe substitute;
- creating an executable plan requires a materially new product or architecture decision.

After delegation, use `$executing-plans-auto-review`'s blocker rules without duplicating them here. Report the exact blocker, evidence, and smallest authority or decision needed to resume.

## Common Mistakes

| Mistake | Required behavior |
|---|---|
| Expand a simple request into a reusable platform | Keep one minimal task and use existing patterns |
| Stop after writing the plan | Invoke `$executing-plans-auto-review` automatically with the saved plan |
| Reimplement the child's review loop here | Keep `$executing-plans-auto-review` as the single downstream source of truth |
| Call `$executing-plans` directly | Delegate through `$executing-plans-auto-review` so its complete approval workflow applies |
