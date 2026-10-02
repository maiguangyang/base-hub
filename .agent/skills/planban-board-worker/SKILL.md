---
name: planban-board-worker
description: Use when the user wants Codex to auto-login to Planban, process a board URL, remember the current board, scan Triage cards, read descriptions/comments/screenshots, move executable cards from Triage to To Do to In Progress, move non-executable cards to Need To Discuss, and invoke $systematic-debugging before any code edits.
---

# Planban Board Worker

Use this skill when the user wants Codex to treat a Planban board as a working queue instead of a passive reference.

## Entry Phrases

Accept these user-facing entry phrases:

- `用 $planban-board-worker 处理这个 board: <url>`
- `用 $planban-board-worker 半自动处理这个 board: <url>`
- `继续处理当前 board`

Default to automatic mode unless the user explicitly says `半自动`.

## Ground Rules

- Auto-login through the shared Planban client and local repo-root `.env` credentials. Never expose credentials, access tokens, cookies, or private attachment URLs in chat.
- Remember the current board context in `work/planban-board-worker/state/current-board.json`.
- If the user does not provide `board <url>` and no remembered board exists, stop and ask for the board URL before doing anything else.
- Always start by running `scripts/scan-board.mjs`.
- When screenshots or image attachments exist for a card, inspect those exported artifacts before classifying the card.
- If text evidence and screenshot evidence conflict, bias toward `Need To Discuss`.
- Move `Triage -> To Do` only after the card is judged executable.
- Invoke `$systematic-debugging` before any code edits.
- Move `To Do -> In Progress` only after the debugging plan is written and you are actually beginning the repair.
- Never auto-move cards to `Ready For Test` or `Done` in V1.

## Workflow

1. Resolve the board context.
   - If the user provided a board URL, use it and remember it.
   - If the user said `继续处理当前 board`, load `work/planban-board-worker/state/current-board.json`.
   - If neither exists, ask for the board URL.
2. Run the read-only scan:

```bash
node .agent/skills/planban-board-worker/scripts/scan-board.mjs \
  --board <board-url-if-provided> \
  --mode <auto|semi-auto> \
  --out work/planban-board-worker/runs
```

3. Read the generated `summary.md` and `manifest.json`, then inspect screenshots if the run exported any visual evidence.
4. For each `Triage` card:
   - summarize the problem, expected behavior, affected area, and missing information
   - if evidence is insufficient or conflicting, comment the blocker and move the card to `Need To Discuss`
   - if evidence is sufficient, comment the executable summary and move the card to `To Do`
5. Before changing code for a `To Do` card, invoke `$systematic-debugging` and write the repair plan.
6. Only after the debugging plan is clear, move the card from `To Do` to `In Progress`.
7. In automatic mode, continue with the next eligible `Triage` card.
8. In semi-automatic mode, stop after each card and ask whether to continue.

## Operational Commands

Comment an executable summary:

```bash
node .agent/skills/planban-board-worker/scripts/comment-card.mjs \
  --card <card-id> \
  --text "<short executable summary>"
```

Comment a blocker before `Need To Discuss`:

```bash
node .agent/skills/planban-board-worker/scripts/comment-card.mjs \
  --card <card-id> \
  --text "<what information is missing or conflicting>"
```

Move an executable card to `To Do`:

```bash
node .agent/skills/planban-board-worker/scripts/move-card.mjs \
  --board <board-url> \
  --card <card-id> \
  --to-list "To Do"
```

Move a started repair to `In Progress`:

```bash
node .agent/skills/planban-board-worker/scripts/move-card.mjs \
  --board <board-url> \
  --card <card-id> \
  --to-list "In Progress"
```

When actually mutating Planban, re-run the same command with `--apply`.

## Classification Rules

- `Executable` means the card has enough concrete evidence to start a repair without additional product clarification.
- `Need To Discuss` means any of the following are true:
  - the expected behavior is unclear
  - screenshots contradict the text
  - the scope depends on another team decision
  - the screenshots or attachments needed for safe judgment are missing
- Do not silently infer business rules from weak evidence. When in doubt, write a precise blocker comment and move to `Need To Discuss`.

## V1 Boundaries

- Supported automated transitions:
  - `Triage -> To Do`
  - `To Do -> In Progress`
  - `Triage -> Need To Discuss`
- Not automated in V1:
  - `Ready For Test`
  - `Done`
  - arbitrary column hopping
  - label edits, assignee edits, or title edits

## References

- Read `references/workflow.md` for the exact state machine and command sequence.
- Read `references/safety.md` for stop conditions, privacy rules, and dry-run behavior.
- Read `references/api.md` for the live-probed Planban API contract.
