# Realtime Websocket Replay

Use websocket replay when the user cares about app-visible behavior, wire ordering, or end-to-end latency from realtime session events.

## Canonical Repo Entry Points

Prefer existing live debug probes first:

- `base-engine/tmpdebug/app_realtime_live_debug_test.go`
- `base-engine/tmpdebug/menu_query_latency_live_debug_test.go`

## Common Environment Inputs

- `AI_POS_APP_PREFS_PLIST`
- `AI_POS_REALTIME_WS_URL`
- `AI_POS_REALTIME_PROMPT`
- `AI_POS_REALTIME_AGENT_ID` (real Agent)
- `AI_POS_REALTIME_CONVERSATION_ID` (omit for fresh replay)
- `AI_POS_REALTIME_DEVICE_TIMEZONE`
- `AI_POS_REALTIME_REPLAY_UNTIL_TERMINAL=1` (multi-tool rounds; stops at approval/interaction, never auto-approves)
- `AI_POS_REALTIME_LOCALE`
- `AI_POS_REALTIME_READ_TIMEOUT_SECONDS`
- `AI_POS_MENU_QUERY_LATENCY_DEBUG`

Preferred local websocket endpoint when verifying the current workspace:

- `ws://127.0.0.1:9527/realtime`

LAN endpoints such as `ws://192.167.167.129:9527/realtime` are valid only when you intentionally want that listener. Do not use them by default for fresh local code verification.

## Before Replaying

Check these first:

- `AI_POS_APP_PREFS_PLIST` points to a current plist, not a stale simulator container path
- the listener on `9527` is the engine process you really want to measure
- if you just changed code, restart the engine from the current workspace before trusting websocket evidence
- if the planning path is slow, set `AI_POS_REALTIME_READ_TIMEOUT_SECONDS=45`

## Example Commands

Run a real menu-query websocket replay using App login/selection and Engine server grants:

```bash
cd base-engine && GOCACHE=/tmp/go-cache AI_POS_REALTIME_READ_TIMEOUT_SECONDS=45 go test ./tmpdebug -run 'TestDebugAppRealtimeQueryWithServerGrants' -count=1
```

Run websocket latency replay with automatic approval:

```bash
cd base-engine && GOCACHE=/tmp/go-cache AI_POS_REALTIME_READ_TIMEOUT_SECONDS=45 go test ./tmpdebug -run 'TestMenuQueryRealtimeWebsocketLatencyWithAutoApproval' -count=1
```

Run the real app replay used for interaction / first-business-node verification:

```bash
cd base-engine && GOCACHE=/tmp/go-cache \
AI_POS_APP_PREFS_PLIST='<plist-path>' \
AI_POS_REALTIME_WS_URL='ws://127.0.0.1:9527/realtime' \
AI_POS_REALTIME_AGENT_ID='<real-agent-id>' \
AI_POS_REALTIME_PROMPT='查询店铺列表' \
AI_POS_REALTIME_LOCALE='zh-CN' \
AI_POS_REALTIME_READ_TIMEOUT_SECONDS=45 \
go test -v ./tmpdebug -run TestDebugAppRealtimeQueryWithServerGrants -count=1
```

## Common Pitfalls

- websocket replay is pointed at an old engine process instead of the code you just changed
- read timeout expires before a slow planning path reaches `interaction_request`
- the plist path is stale, so the replay is not using live app-cached auth at all
- the replay result is treated as "engine truth" even though only the transport path was measured

## Evidence To Capture

For websocket replay, record milestones like:

- websocket dial start
- websocket connected
- `session.init.sent`
- `session.ready`
- `response.create.sent`
- `approval_request.awaiting`
- `approval_request.submit.sent`
- `interaction_request.awaiting`
- `interaction_request.submit.sent`
- first assistant text delta
- `response.completed`

## When Websocket Replay Is Not Enough

Websocket replay alone is not enough when the question is really about:

- `StartRealtimeTurn` latency
- engine state transitions
- loop-run persistence
- server grant resolution
- AI usage token / latency accounting

In those cases, pair websocket replay with engine `tmpdebug` in-process replay so you can separate transport timing from engine timing.
