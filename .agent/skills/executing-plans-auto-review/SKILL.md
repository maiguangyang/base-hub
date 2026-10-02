---
name: executing-plans-auto-review
description: Use when a written implementation plan must reach code-review approval autonomously without human checkpoints between execution, review, and remediation.
---

# Executing Plans With Automatic Review

## Overview

Execute a written plan once, then repeat read-only review and full repair until the latest review returns no fixes.

**Core principle:** Invocation is standing authorization for execution and an unbounded review-and-repair loop. Review remains read-only; after each review pass ends, repair every actionable fix before reviewing again.

**DIRECTLY ORCHESTRATED SUB-SKILLS:** Use only `$executing-plans` and `$requesting-code-review`. Invoke them in exactly this sequence:

```text
$executing-plans -> (
  $requesting-code-review
  -> if fixes exist: repair ALL fixes -> verify -> repeat review
  -> if no fixes exist: approval gate
)
```

Do not return to plan execution after entering the review loop. Repair is a wrapper-authorized phase, not another sub-skill.

## Authorization Contract

- Treat explicit invocation of this skill as authorization to execute the plan and repair every actionable fix from every review pass without asking again.
- Keep every `$requesting-code-review` pass strictly read-only. Report findings and end that pass before changing files. This wrapper invocation supplies the separate explicit repair authorization required by `$requesting-code-review`; start the repair phase immediately without waiting for another user message.
- Do not stage, commit, merge, push, archive, perform destructive operations, or change external systems unless separately authorized.
- Continue in the current workspace unless already inside a linked worktree. Do not pause to ask about worktree isolation.

## Workflow

1. **Execute** — Run `$executing-plans` through all plan tasks and required verification. When it reaches its normal review handoff, return control here instead of stopping.
2. **Review** — Run `$requesting-code-review` against the complete cumulative change and the original plan or requirements. Use the smallest complete review artifact. Do not modify files during this phase.
3. **Inspect the result** — Treat every actionable requested code, test, documentation, configuration, or behavior change as a `fix`, whether labeled Critical, Important, Minor, suggestion, recommendation, or otherwise. If at least one fix exists, continue to Step 4. If zero fixes exist, evaluate the Approval Gate.
4. **Repair all fixes** — After the review pass ends, repair every returned fix from highest severity to lowest. Do not reject, downgrade, defer, or relabel a returned fix as optional, speculative, invalid, out of scope, or unnecessary. A review verdict such as “approved” or “ready” does not override remaining fixes.
5. **Verify** — Run focused verification for each repair, then run the plan's complete verification set after the entire repair batch.
6. **Re-review** — Create fresh review input for the complete cumulative change and invoke `$requesting-code-review` again. Return to Step 3. Never reuse an earlier verdict, and never stop while the latest review contains a fix.

There is no review-round limit. The required steady-state loop is always:

```text
$requesting-code-review -> repair ALL fixes -> verify -> $requesting-code-review
```

If a returned fix cannot be performed safely because it is contradictory, destructive, external, or requires a materially new product or architecture decision, do not call it resolved or ignore it. Enter the Genuine Blockers path instead; the workflow is not approved.

## Approval Gate

Stop successfully only when all conditions are true:

- the original plan and requirements are fully satisfied;
- the latest review returns zero actionable fixes of every severity and label;
- the latest `$requesting-code-review` verdict is explicitly approved / ready to merge;
- fresh complete verification passes;
- the Git staging index is unchanged from the start of the workflow.

Then leave all changes local, report the review rounds and verification evidence, and stop. Do not enter branch finishing automatically.

## Genuine Blockers

Human intervention is not a normal checkpoint. Stop only when autonomous continuation would be unsafe or impossible, including:

- contradictory or materially ambiguous requirements;
- missing credentials, dependencies, permissions, or unavailable external systems with no safe local substitute;
- a required destructive or external-state action that lacks separate authorization;
- a review finding that requires a new product or architecture decision outside the plan;
- three evidence-backed repair attempts for the same root problem fail, or consecutive review rounds make no measurable progress.

Report the exact blocker, evidence gathered, attempts made, and the smallest decision or authority needed to resume.

## Common Mistakes

| Mistake | Required behavior |
|---|---|
| Stop after the first review because review is read-only | End the review pass, then repair under this wrapper's standing authorization |
| Repair only Critical or Important findings | Repair every returned fix, including Minor findings and suggestions |
| Reclassify a review fix as optional or invalid | A returned actionable change is a fix; repair it or report a genuine blocker |
| Treat an approved verdict as permission to leave fixes open | Any remaining fix forces another repair and review round |
| Treat passing tests as review approval | Require both fresh verification and an explicit approved review verdict |
| Keep looping without progress | Stop on the genuine-blocker rules and show evidence |
