---
name: subagent-driven-development
description: Use when executing an implementation plan with independent tasks and the user explicitly wants delegated worker execution in the current session
---

# Subagent-Driven Development

Execute a plan by dispatching a fresh delegated worker for each task, handing task context through files, running one task reviewer that checks spec compliance and code quality together, then running one broad final review.

**Activation rule:** Use this skill only when the user explicitly asks for subagents, delegation, or parallel agent-style execution. If the user did not explicitly ask for delegated workers, prefer `$single-flow-task-execution` or `$executing-plans`.

**Core principle:** Fresh delegated worker per task + file handoffs + one task review (spec + quality) + final broad review = faster iteration with controlled context.

**Narration:** Between tool calls, narrate at most one short line. The ledger, tracker, and tool results carry the record.

**Continuous execution:** Once the user has explicitly chosen this skill, do not pause for "should I continue?" style check-ins between tasks. Continue until:
- all tasks complete
- a real blocker requires human input
- the plan is ambiguous enough that safe execution cannot continue

## When to Use

Use this skill when all of the following are true:

- You have a written implementation plan
- The plan tasks are mostly independent
- The user explicitly wants delegated worker execution
- You want to stay in the current session instead of switching to a separate local execution workflow

Prefer `$executing-plans` when:

- The user wants local sequential execution without delegated workers
- The work is tightly coupled and better handled sequentially in one thread
- Delegated workers would create more coordination overhead than value

## The Process

0. If the user has not already chosen or declined `$using-git-worktrees`, and you are not already inside an isolated linked worktree, ask once whether they want an isolated workspace before dispatching workers. If they opt in, use `$using-git-worktrees` first. If they decline, continue here and do not ask again in downstream execution skills.
1. Resolve the SDD workspace with `scripts/sdd-workspace` from this skill directory. It creates `.superpowers/sdd`, an ignored scratch area for briefs, reports, review packages, and `progress.md`.
2. Check `.superpowers/sdd/progress.md` if it exists. Tasks already marked complete are done unless `git status` proves otherwise.
3. Read the plan once, extract task numbers, global constraints, file ownership, interfaces, and verification commands.
4. Run the pre-flight plan review below before dispatching Task 1.
5. Resolve the active-window tracker path, then update it with the active task list before dispatching workers.
6. For each task, create a task brief with `scripts/task-brief PLAN_FILE N`. Name the report file beside it, usually `task-N-report.md`.
7. Dispatch one implementer worker using `./implementer-prompt.md`. The prompt points at the brief and report file instead of pasting the full task text.
8. When the implementer returns `DONE` or `DONE_WITH_CONCERNS`, create a review package:
   - Default for this repo: `scripts/review-package --worktree task-N -- PATH...` using the task's changed files or owned paths
   - If path scope is unknown and the worktree is clean enough: `scripts/review-package --worktree task-N`
   - Only if the user/project explicitly permits task commits: `scripts/review-package BASE HEAD`
9. Dispatch one task reviewer using `./task-reviewer-prompt.md` with the brief file, report file, review package path, and binding global constraints.
10. If the reviewer finds Critical or Important issues, send all actionable findings back to a fixer worker, require focused verification, append the fix report to the same report file, regenerate the review package, and re-run the task reviewer.
11. Mark the task complete in both the active-window tracker and `.superpowers/sdd/progress.md` only after the task reviewer approves it and any cannot-verify items have been resolved by the controller.
12. After all tasks complete, create a final review package and dispatch one final broad reviewer using `../requesting-code-review/code-reviewer.md`, then hand off to `$finishing-a-development-branch`.

## Pre-Flight Plan Review

Before dispatching Task 1, scan the plan once for conflicts:

- tasks that contradict each other or the plan's global constraints
- task boundaries that require overlapping file ownership
- plan instructions that would violate repo rules, such as staging, committing, or runtime hardcoding
- missing interface definitions between tasks
- anything the plan explicitly mandates that the review rubric treats as a defect, such as a test that asserts nothing or verbatim duplication of a logic block

If the scan finds conflicts, present all of them as one batched question before execution begins. Put each concern beside the plan text that mandates it and ask which governs. If the scan is clean, proceed without a status checkpoint.

## Model Choice

Use the least expensive worker that can reliably handle the role, but obey the current platform's model-selection rules.

- Mechanical implementation tasks with complete specs and 1-2 touched files can use a faster, cheaper worker when model override is allowed.
- Integration, debugging, and multi-file coordination need a standard or stronger worker.
- Reviewers need enough judgment for the diff size, complexity, and risk.
- In Codex, set a model override only when the user explicitly asked for it or the task has a clear task-specific reason. Otherwise inherit the current model.

Turn count beats token price. If a cheap worker is likely to need repeated clarification or rework, choose a stronger worker up front when policy permits.

## Worker Statuses

Implementer workers should report one of these statuses:

- `DONE` — task implemented and ready for review
- `DONE_WITH_CONCERNS` — implemented, but the worker has correctness or maintainability concerns
- `NEEDS_CONTEXT` — missing information blocked safe execution
- `BLOCKED` — task cannot be completed without a different approach or human guidance

Handle them like this:

- `DONE` — create the review package and dispatch the task reviewer
- `DONE_WITH_CONCERNS` — read the concerns first; address correctness or scope concerns before review, otherwise note them and continue to review
- `NEEDS_CONTEXT` — provide the missing context and re-dispatch
- `BLOCKED` — either provide more context, break down the task further, switch approach, or escalate to the user

Never ignore a worker escalation and never blindly retry with the same instructions.

## File Handoffs

Everything pasted into a dispatch prompt and everything printed back stays in the controller context. Hand large artifacts over as files.

- **Task brief:** Run `scripts/task-brief PLAN_FILE N` and pass the printed path. The dispatch prompt should contain only where the task fits, the brief path, relevant interfaces, ambiguity resolutions, and the report-file path.
- **Report file:** Put the report beside the brief. The worker writes full details there and returns only status, one-line test summary, concerns, and report path.
- **Review package:** Use `scripts/review-package --worktree task-N -- PATH...` for this repo's no-commit workflow so unrelated local files stay out of the package. Use `scripts/review-package BASE HEAD` only when commits are explicitly allowed. Pass the printed path to the reviewer.
- **Reviewer prompt:** Give the reviewer exactly three files: brief, report, and review package. Do not ask the reviewer to infer requirements from conversation history.
- **Fix dispatches:** Send the complete actionable finding list in one dispatch, name the focused tests to run, and append fix results to the same report file before re-review.

## Reviewer Prompt Discipline

Per-task reviews are task-scoped gates. The broad review happens once at the end.

- Use `./task-reviewer-prompt.md`; do not run separate spec and quality reviewers for routine tasks.
- Do not pre-judge findings for reviewers. If your prompt contains "do not flag", "ignore this", or "the plan chose this so it is fine", stop and rewrite it.
- Do not ask a reviewer to re-run tests the implementer already ran unless there is a concrete reason.
- If a reviewer says "cannot verify from diff", resolve that item yourself before marking the task complete. If it is a real gap, send it back for a fix and re-review.
- Critical and Important findings must be fixed and re-reviewed before continuing. Minor findings may be noted in the tracker for final triage.
- If a finding conflicts with plan text, present the finding and plan text to the user. Do not dismiss the finding because the plan mandated it, and do not fix against the plan without asking.

## No-Commit Workflow

This repository keeps changes local for human review. Workers must not run `git add` or `git commit`.

Because no per-task commits exist, `scripts/review-package --worktree task-N -- PATH...` packages the current working-tree status, staged diff, unstaged diff, and untracked file contents for the provided path scope. Prefer the task's changed files or owned directories as `PATH...`; fall back to the whole worktree only when the tree is clean enough. If earlier completed task changes are still present in the package, the controller narrows the review by pointing at the current task's changed files and interfaces.

Use commit-range review packages only when the user and project rules explicitly allow task commits.

## Durable Progress

Conversation memory can be compacted. The ledger and tracker are the recovery map.

- At skill start, read `.superpowers/sdd/progress.md` if it exists and trust completed entries unless `git status` contradicts them.
- After a task review passes, append one line: `Task N: complete (review package PATH, verification SUMMARY)`.
- Update the resolved active-window tracker with task status, worker result, review result, changed files, and verification summary as work progresses.
- After a context transition, trust the ledger, active-window tracker, and `git status` over memory.
- Do not re-dispatch a completed task unless the ledger/tracker or git state proves it is incomplete.

## Prompt Templates

- `./implementer-prompt.md` — implement one task with a brief file, report file, self-review, and explicit status
- `./task-reviewer-prompt.md` — review one task for spec compliance and code quality together
- `../requesting-code-review/code-reviewer.md` — final broad review after all tasks complete

## Red Flags

Never:

- Use this skill unless the user explicitly asked for delegated workers
- Dispatch multiple implementer workers that edit overlapping files
- Dispatch separate spec and quality reviewers for a routine task
- Let "self-review" replace formal review
- Paste the whole plan into every worker prompt
- Allow workers to auto-stage or auto-commit changes
- Pre-judge reviewer findings in the prompt
- Ignore reviewer "cannot verify" items
- Start implementation on `main` or `master` without explicit user consent

Always:

- Keep one task active at a time unless the user explicitly asks for broader parallelism
- Provide full task text through a brief file and exact constraints up front
- Keep changes local for human review
- Re-run the task reviewer after fixes
- Update durable progress in the ledger and active-window tracker

## Integration

**Required workflow skills:**

- `$using-git-worktrees` — REQUIRED one-time decision check unless the user already chose or declined isolation, or you are already inside an isolated linked worktree. If the user opts in, use it before dispatching workers.
- `$writing-plans` — creates the plan this skill executes
- `$requesting-code-review` — provides the code review template used by reviewer workers
- `$finishing-a-development-branch` — handles final completion after all tasks are done

**Delegated workers should follow:**

- `$test-driven-development` — when the task requires new behavior or a bug fix

**Alternative workflow:**

- `$executing-plans` — use for local sequential execution instead of same-session delegated workers
