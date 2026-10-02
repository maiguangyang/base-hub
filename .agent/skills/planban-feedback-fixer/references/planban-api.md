# Planban API Notes

The Planban web UI may not render reliably in automated browsers. Use the API for read-only feedback extraction.

## Endpoints

Base API:

```text
https://planban-api.sjfood.us/api
```

Login:

```http
POST /access-tokens?withHttpOnlyToken=false
Content-Type: application/json

{"emailOrUsername":"...","password":"..."}
```

The response is:

```json
{"item":"JWT_ACCESS_TOKEN"}
```

Use the token as:

```http
Authorization: Bearer JWT_ACCESS_TOKEN
```

Board:

```http
GET /boards/:boardId?subscribe=true
```

Card actions/comments:

```http
GET /cards/:cardId/actions
```

Attachments can be downloaded from the `url` fields in board data. Send the same bearer token.

## Data Shape

`GET /boards/:id?subscribe=true` returns:

- `item`: board metadata.
- `included.labels`: labels on the board.
- `included.lists`: columns/lists.
- `included.cards`: cards.
- `included.cardLabels`: card-to-label joins.
- `included.attachments`: card attachments.
- `included.users`: creators and members.

The exporter resolves names for lists, labels, and users, filters cards by label, and writes one folder per card.

## Safety

Do not log tokens. Prefer environment variables or account files over command-line passwords. Do not update cards through the API unless the user explicitly asks.
