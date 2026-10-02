# Evidence Checklist

Every real replay result should include four sections:

1. `数据来源`
2. `执行方式`
3. `关键证据`
4. `结论与边界`

## Minimum Evidence By Scenario

### App -> Engine Replay

Always state:

- where the token came from
- where the runtime headers came from
- whether the source was `AI_POS_APP_PREFS_PLIST`, simulator prefs, or another live source
- whether replay was websocket, engine in-process, or both

### Latency Analysis

Capture as many of these milestones as the scenario supports:

- `session.ready`
- `response.create.sent`
- `approval_request.awaiting`
- `approval_request.submit.sent`
- `interaction_request.awaiting`
- `interaction_request.submit.sent`
- tool execution start / completion
- first assistant text
- `response.completed`

When available, also capture engine-side timing such as:

- `StartRealtimeTurn`
- `RunRealtimeTurn`
- AI usage token counts
- AI usage latency

### Runtime Header Diagnostics

State clearly:

- which hub / capability / tool path is affected
- whether runtime headers came from app-local cached state
- whether the failure is caused by:
  - missing source data
  - transport omission
  - preflight / orchestration behavior
  - downstream MCP execution

### Approval / Interaction Diagnostics

Capture identifiers when available:

- request id
- interaction id
- conversation id
- response id

Also state which layer failed:

- app render / notifier layer
- websocket protocol layer
- engine orchestration layer
- downstream tool execution layer

## Boundary Rules

Always say which of these categories the conclusion belongs to:

- websocket only
- engine only
- websocket + engine
- code inspection only

If the simulator source or plist is unavailable, say explicitly that true replay was blocked.

## Preferred Closing Sentence

Use a boundary-conscious close such as:

`已验证的部分是 X；当前仍未直接验证的边界是 Y。`
