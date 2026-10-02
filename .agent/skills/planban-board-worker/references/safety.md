# Planban Board Worker Safety

## Credentials

- Read Planban credentials from the repo-root `.env` file or the supported account-file fallbacks.
- Never print usernames, passwords, bearer tokens, cookies, or private attachment URLs into chat output.
- Do not duplicate credentials into additional files unless the workflow explicitly requires it.

## Attachments and Privacy

- Treat screenshots, videos, and attachments as potentially private.
- Inspect only the artifacts needed to classify the card or debug the repair.
- Do not paste private image URLs or raw attachment payloads into the final answer.
- If an attachment cannot be downloaded safely, say that evidence is missing and bias to `Need To Discuss`.

## Dry-Run vs Apply

- `comment-card.mjs` defaults to dry-run.
- `move-card.mjs` defaults to dry-run.
- Only add `--apply` when the intended mutation is correct and the user flow requires a real Planban update.
- Prefer dry-run first when validating a new board, new card, or new workflow branch.

## Screenshot Conflict Rule

- If a screenshot contradicts the title, description, or comments, trust neither side blindly.
- Write a short blocker comment describing the conflict.
- Move the card to `Need To Discuss` instead of forcing it through `To Do`.

## Stop Conditions

Stop and ask for intervention when any of these happen:

- missing board URL and no remembered board exists
- required board lists are missing
- a live write fails
- the card is not present on the fetched board snapshot
- the debugging plan from `$systematic-debugging` is not concrete enough to start implementation
- the card depends on missing screenshots or conflicting evidence

## Mutation Boundaries

Allowed automation in V1:

- comment on a card
- move a card to `To Do`
- move a card to `In Progress`
- move a card to `Need To Discuss`

Not allowed in V1:

- auto-move to `Ready For Test`
- auto-move to `Done`
- edit titles
- edit labels
- edit assignees
- perform arbitrary status rewrites

## Human-Friendly Defaults

- Default to automatic mode unless the user explicitly says `半自动`.
- If the user says `继续处理当前 board`, reuse remembered board state rather than asking again.
- If a remembered board does not exist, ask for the board URL instead of guessing.
