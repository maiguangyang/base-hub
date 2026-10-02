# Planban Board Worker Workflow

## State Machine

Supported automated state transitions in V1:

- `Triage -> To Do`
- `To Do -> In Progress`
- `Triage -> Need To Discuss`

Not automated in V1:

- `Ready For Test`
- `Done`

## Mode Behavior

### Automatic Mode

- Trigger phrase: `用 $planban-board-worker 处理这个 board: <url>`
- Default when the user does not mention `半自动`
- Scan the whole board, then continue card by card without stopping after each completed classification

### Semi-Automatic Mode

- Trigger phrase: `用 $planban-board-worker 半自动处理这个 board: <url>`
- Use the same scan and classification flow
- Stop after each card and ask whether to continue with the next card

## Remembered Board Context

The worker stores remembered state in:

```text
work/planban-board-worker/state/current-board.json
```

This state remembers:

- `currentBoardUrl`
- `currentMode`
- `currentCardId`
- `processedCardIds`
- `lastRunAt`

Resume behavior:

- If the user says `继续处理当前 board`, load `currentBoardUrl` and `currentMode` from that file.
- If the user gives a new board URL, treat it as the new active board and update the remembered state.
- If there is no remembered board and no explicit board URL, stop and ask for the board URL.
- Validated on `2026-06-25`: a second scan without `--board` successfully reused the remembered board and updated `currentMode` to `semi-auto`.

## Exact Script Sequence

### 1. Scan the Board

With an explicit board URL:

```bash
node .agent/skills/planban-board-worker/scripts/scan-board.mjs \
  --board https://planban.sjfood.us/boards/<board-id> \
  --mode auto \
  --out work/planban-board-worker/runs
```

Resume with remembered board:

```bash
node .agent/skills/planban-board-worker/scripts/scan-board.mjs \
  --mode semi-auto \
  --out work/planban-board-worker/runs
```

Outputs:

- a timestamped run directory under `work/planban-board-worker/runs/`
- `summary.md`
- `manifest.json`
- Validated on `2026-06-25` against board `1793465014108029954`.

### 2. Comment an Executable Summary

Use this after deciding the card can be worked directly:

```bash
node .agent/skills/planban-board-worker/scripts/comment-card.mjs \
  --card <card-id> \
  --text "AI summary: bug is reproducible, scope is clear, ready for implementation review"
```

Add `--apply` only when you are ready to write the comment to Planban.

### 3. Comment a Needs-Discussion Blocker

Use this when evidence is insufficient or conflicting:

```bash
node .agent/skills/planban-board-worker/scripts/comment-card.mjs \
  --card <card-id> \
  --text "AI blocker: screenshot and text conflict; expected behavior needs product clarification"
```

Add `--apply` only when you are ready to write the blocker to Planban.

### 4. Move the Card to `To Do`

Only after the card is judged executable:

```bash
node .agent/skills/planban-board-worker/scripts/move-card.mjs \
  --board https://planban.sjfood.us/boards/<board-id> \
  --card <card-id> \
  --to-list "To Do"
```

Add `--apply` only when you are ready to mutate Planban.

### 5. Invoke `$systematic-debugging`

Before editing code, use `$systematic-debugging` to produce the repair plan and root-cause analysis.

Only continue when the debugging plan is concrete enough to start the fix.

### 6. Move the Card to `In Progress`

Use this only after the debugging plan exists and repair work is actually beginning:

```bash
node .agent/skills/planban-board-worker/scripts/move-card.mjs \
  --board https://planban.sjfood.us/boards/<board-id> \
  --card <card-id> \
  --to-list "In Progress"
```

Add `--apply` only when you are ready to mutate Planban.

## Classification Heuristic

Classify a `Triage` card as executable only when:

- the bug or change request is understandable from title, description, comments, and visuals
- the expected behavior is concrete enough to test
- no missing dependency blocks implementation

Classify a `Triage` card as `Need To Discuss` when:

- screenshots contradict the text
- expected behavior is ambiguous
- the card depends on external product decisions
- the evidence is too weak to produce a safe debugging plan

## Validated Examples

Dry-run comment:

```bash
node .agent/skills/planban-board-worker/scripts/comment-card.mjs \
  --card "$PLANBAN_TEST_CARD_ID" \
  --text "Codex dry-run probe"
```

Dry-run move:

```bash
node .agent/skills/planban-board-worker/scripts/move-card.mjs \
  --board https://planban.sjfood.us/boards/1793465014108029954 \
  --card "$PLANBAN_TEST_CARD_ID" \
  --to-list "$PLANBAN_TEST_FIRST_TARGET_LIST"
```

Live comment and move were verified on `2026-06-25` against the sacrificial test card configured in repo-root `.env`.

Observed command outputs:

- dry-run comment: `DRY RUN: would comment on card <card-id>`
- dry-run move: `DRY RUN: would move card <card-id> to <target-list>`
- live comment: `Commented on card <card-id>`
- first live move: `Moved card <card-id> to To Do`
- second live move: `Moved card <card-id> to In Progress`
