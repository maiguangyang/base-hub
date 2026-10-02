---
name: planban-feedback-fixer
description: Use when fixing Planban board feedback from labels such as backend, 后台, or 後台, especially when cards include bug descriptions, screenshots, videos, comments, or a board URL.
---

# Planban Feedback Fixer

Use this skill to turn Planban feedback cards into local repair tasks and then fix them in a codebase.

## Ground Rules

- Treat Planban as read-only unless the user explicitly asks to comment, move cards, or update status.
- Do not expose credentials, access tokens, cookies, or private attachment URLs in chat.
- Always keep the Planban card id in notes, commits, PR summaries, and final status.
- Prefer the API exporter when the Planban UI spins, redirects, or hides DOM content.

## Quick Start

From the repo or current workspace:

Export feedback:

```bash
node .agent/skills/planban-feedback-fixer/scripts/fetch-feedback.mjs \
  --board https://planban.sjfood.us/boards/1793465014108029954 \
  --label 後台 \
  --out work/planban-feedback
```

On first run, if credentials are missing, the exporter creates `planban-account` in the current directory and asks the user to fill it in. The file is ignored by git and auto-loaded on the next run.

The generated account file is plain text with two non-comment lines:

```text
username_or_email
password
```

Credentials can also be provided with `PLANBAN_USERNAME` and `PLANBAN_PASSWORD`, or with an explicit file:

```bash
--account-file /absolute/path/to/account-file
```

Never put credentials in a final answer.

## Workflow

1. Export the issue pack with `scripts/fetch-feedback.mjs`.
2. Read `summary.md` and `manifest.json` to understand scope and grouping.
3. Group related cards by feature area before editing: user management, role management, company management, skill management, or UI polish.
4. For each selected card, read its `issue.md` and inspect downloaded screenshots in `attachments/`.
5. If a screenshot and description conflict, trust the screenshot and call out the ambiguity.
6. Fix in small batches. Use existing project patterns and tests.
7. For functional bugs, write or update tests before implementation when feasible.
8. For visual/UI feedback, verify with browser screenshots at relevant desktop/mobile sizes.
9. Report completed card ids and any cards that need product clarification.

## Priority Heuristic

Handle data/permission/state bugs before cosmetic issues:

1. Blocking workflow failures, failed generation, failed save, login/session issues.
2. Deleted/disabled records still editable or visible.
3. Required field, validation, and missing data-column mismatches.
4. Incorrect aggregate counts or filtered totals.
5. Empty-state copy, alignment, spacing, and loading polish.

## References

- Read `references/planban-api.md` when changing the exporter or manually querying Planban.
- Read `references/repair-workflow.md` when deciding how to batch and verify feedback fixes.
