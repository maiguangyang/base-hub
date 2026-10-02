---
name: requesting-code-review
description: Use when completing tasks, implementing major features, or before merging to run a strictly read-only review against requirements; invoking this skill never authorizes fixes or file changes
---

# Requesting Code Review

Run a structured review pass to catch issues before they cascade.

**Core principle:** Review early, review often, and never confuse review authority with repair authority.

## Authorization Boundary

Invoking this skill authorizes **review only**. It never authorizes modifying code, tests, documentation, configuration, generated files, the working tree, or the staging index.

- Report every finding and stop after the review verdict.
- Do not add a regression test to prove a finding inside the reviewed checkout.
- Do not fix Critical or Important findings automatically.
- A finding is evidence, not repair authorization.
- Proceed with repairs only after the user gives a separate, explicit fix instruction after seeing the review result. That repair is a new task, outside this code-review invocation.

## When to Request Review

**Mandatory:**

- After each task in single-flow task execution
- After completing major feature
- Before merge to main

**Optional but valuable:**

- When stuck (fresh perspective)
- Before refactoring (baseline check)
- After fixing complex bug

## How to Request

**1. Prepare review input:**

Use the smallest complete artifact:

- Commit-based workflow: record `BASE_SHA` and `HEAD_SHA`
- No-commit workflow: provide a review package file, for example `subagent-driven-development/scripts/review-package --worktree final`

**2. Run structured code review checklist:**

Use `requesting-code-review/code-reviewer.md` template and review the diff against requirements. In Codex single-flow mode, do not dispatch generic coding agents.

If the user explicitly asked for delegated review and subagents are available, you may dispatch a reviewer with the same template. Give it precise context only; do not pass the whole conversation history.

**Placeholders:**

- `{WHAT_WAS_IMPLEMENTED}` - What you just built
- `{PLAN_OR_REQUIREMENTS}` - What it should do
- `{REVIEW_PACKAGE}` - Optional path to a generated review package containing status, stat summary, and diff
- `{BASE_SHA}` - Starting commit, or `worktree` when using a no-commit review package
- `{HEAD_SHA}` - Ending commit, or `worktree` when using a no-commit review package
- `{DESCRIPTION}` - Brief summary

**3. Report feedback and stop:**

- Report Critical, Important, and Minor findings with evidence.
- Mark the change not ready when unresolved Critical or Important findings exist.
- Explain questionable findings with technical reasoning, but do not modify files to prove or disprove them.
- Wait for a separate user instruction before starting any repair workflow.

Do not interpret severity, workflow momentum, or “before proceeding” language as permission to fix. Code review has no repair authority.

## Example

```
[Just completed Task 2: Add verification function]

You: Let me request code review before proceeding.

BASE_SHA=$(git log --oneline | grep "Task 1" | head -1 | awk '{print $1}')
HEAD_SHA=$(git rev-parse HEAD)

[Run checklist-based review]
  WHAT_WAS_IMPLEMENTED: Verification and repair functions for conversation index
  PLAN_OR_REQUIREMENTS: Task 2 from docs/plans/deployment-plan.md
  BASE_SHA: a7981ec
  HEAD_SHA: 3df7661
  DESCRIPTION: Added verifyIndex() and repairIndex() with 4 issue types

[Review returns]:
  Strengths: Clean architecture, real tests
  Issues:
    Important: Missing progress indicators
    Minor: Magic number (100) for reporting interval
  Assessment: Ready to proceed

You: Report the findings and stop without changing files.
[Wait for a separate user instruction]

User: Fix the Important issue from the review.
You: Start a new repair task with its own authorization and verification.
```

## Integration with Workflows

**Single-Flow Task Execution:**

- Review after EACH task
- Catch issues before they compound
- Return findings and wait for explicit repair authorization before moving on when fixes are needed

**Executing Plans:**

- Review after each batch (3 tasks)
- Return feedback; apply it only under a separately authorized repair task

**Ad-Hoc Development:**

- Review before merge
- Review when stuck

## Red Flags

**Never:**

- Skip review because "it's simple"
- Ignore Critical issues
- Proceed with unfixed Important issues
- Modify code, tests, docs, configuration, or generated files during review
- Treat the instruction to report or act on feedback as automatic fix authorization
- Argue with valid technical feedback

**If reviewer wrong:**

- Push back with technical reasoning
- Show code/tests that prove it works
- Request clarification

Use existing code and test evidence only. Do not create or edit tests during the review to make the argument.

See template at: requesting-code-review/code-reviewer.md
