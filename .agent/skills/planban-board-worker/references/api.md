# Planban Board Worker API Contract

Captured on `2026-06-25` from:

- live Planban frontend bundle: `https://planban.sjfood.us/assets/index-CSzEiYMP.js`
- successful live write probe against a sacrificial card on board `1793465014108029954`

## Auth

- Live UI bundle auth route:
  - `POST /access-tokens?withHttpOnlyToken=true`
  - transport: `credentials: include`
- Repo worker bearer-token auth route used by the shared client:
  - `POST /access-tokens?withHttpOnlyToken=false`
  - transport: `authorization: Bearer <token>`

## Read

- `GET /boards/:boardId?subscribe=true`
- `GET /cards/:cardId`
- `GET /cards/:cardId/actions`

## Write

### Comment Card

- Method: `POST`
- Path pattern: `/cards/:cardId/comment-actions`
- Required headers:
  - `authorization`
  - `content-type: application/json`
- Body example:

```json
{
  "text": "Codex live write probe <timestamp>",
  "plainText": "Codex live write probe <timestamp>",
  "atUserIds": []
}
```

### Move Card

- Method: `PATCH`
- Path pattern: `/cards/:cardId`
- Required headers:
  - `authorization`
  - `content-type: application/json`
- Body example:

```json
{
  "listId": "<target-list-id>",
  "position": 131070
}
```

Confirmed on `2026-06-25` with a successful live probe:

- comment write succeeded on the sacrificial test card
- move write succeeded for `Triage -> To Do`
- move write succeeded for `To Do -> In Progress`
- a follow-up script-level live probe also succeeded with:
  - comment text: `Codex write probe comment 2026-06-25`
  - final card state: `In Progress`

## Notes

- The live bundle uses `ht.post('/cards/${cardId}/comment-actions', body)` for comment creation.
- The live bundle uses `ht.patch('/cards/${cardId}', body)` for card updates, including cross-list moves.
- The UI computes move positions client-side and sends `listId` plus numeric `position`.
