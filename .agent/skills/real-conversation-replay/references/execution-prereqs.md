# Execution Prerequisites

Use this checklist before running any real replay in `easypos-hub`.

## 1. Export Repo Env Correctly

Many live debug commands spawn child processes such as `go test` or `go run`.

This repo's `.env` file mainly contains shell assignments. That means:

- `source .env` loads variables into the current shell
- but they are not automatically exported to child processes
- so direct `go test` / `go run` can still fail with missing env such as `DATABASE_URL`

Use this pattern instead:

```bash
cd base-engine
set -a
source .env
set +a
```

After that, confirm the critical env exists:

```bash
test -n "$DATABASE_URL" # Check presence without printing credentials.
```

## 2. Know When `DATABASE_URL` Is Required

Assume `DATABASE_URL` is required whenever the replay path touches any of:

- `StartRealtimeTurn`
- `RunRealtimeTurn`
- `gen.NewDBFromEnvVars`
- `tmpdebug` tests that construct a real resolver / DB
- starting local `base-engine`

If `DATABASE_URL` is missing, fix env export first before investigating code behavior.

## 3. Local Services Must Be Reachable

For DB-backed engine replay, make sure these dependencies are reachable from the shell you are using:

- MySQL from `DATABASE_URL`
- Redis only if the chosen path uses it; it is not a universal replay prerequisite.

Typical failure signs:

- `dial tcp 127.0.0.1:3306`
- `connect: operation not permitted`
- startup panic in `gen.NewDBWithString`

If the failure is sandbox or permission related, rerun with the required escalation instead of treating it as an application bug.

## 4. Simulator Plist Paths Go Stale

Do not assume an old app container UUID is still valid.

The simulator app data container can change after:

- reinstalling the app
- deleting / recreating simulator state
- certain app relaunch / deployment flows

If an old plist path fails, re-resolve it:

```bash
xcrun simctl get_app_container booted com.proton.aiPos.dev data
```

Then point `AI_POS_APP_PREFS_PLIST` at:

```text
<container>/Library/Preferences/com.proton.aiPos.dev.plist
```

If `plutil` fails with `exit status 1`, suspect a stale path first.

## 5. Validate the Real Source Before Replay

The plist must contain the current login and selection state:

- `flutter.dev_token`
- `flutter.dev_appInstallationId`
- `flutter.dev_deviceSkillState`

Old cached headers are intentionally absent. Use real member/application server grants.

Useful check:

```bash
test -f "$AI_POS_APP_PREFS_PLIST" # Parse privately in the replay helper.
```

If these values are missing, the app may not have logged in or cached the runtime state yet.

## 6. Prefer Loopback For Local Workspace Verification

When the goal is to verify your local code changes, prefer:

```bash
AI_POS_REALTIME_WS_URL='ws://127.0.0.1:9527/realtime'
```

Do not rely on a LAN IP unless you intentionally want a remote / shared engine process.

## 7. Verify What Is Listening On `9527`

Before trusting websocket replay evidence for a fresh code change, check the current listener:

```bash
lsof -n -P -iTCP:9527 -sTCP:LISTEN
ps -p <PID> -o pid,lstart,command
```

If you just changed code, do not assume the long-running listener is current.

Use a separately built current-workspace engine on an unused port if the existing listener is shared; stop only the process you started.

## 8. Use A Longer Read Timeout For Slow Planning Paths

Some real flows need more than the default short timeout before the first `interaction_request` or assistant text.

Preferred default for live replay:

```bash
AI_POS_REALTIME_READ_TIMEOUT_SECONDS=45
```

## 9. Known Good Command Skeletons

Export env for DB-backed engine commands:

```bash
cd base-engine
set -a
source .env
set +a
GOCACHE=/tmp/go-cache go test ./tmpdebug -run '<TestName>' -count=1
```

Replay websocket with real app prefs:

```bash
cd base-engine
GOCACHE=/tmp/go-cache \
AI_POS_APP_PREFS_PLIST='<plist-path>' \
AI_POS_REALTIME_WS_URL='ws://127.0.0.1:9527/realtime' \
AI_POS_REALTIME_AGENT_ID='<real-agent-id>' \
AI_POS_REALTIME_PROMPT='查询店铺列表' \
AI_POS_REALTIME_LOCALE='zh-CN' \
AI_POS_REALTIME_READ_TIMEOUT_SECONDS=45 \
go test -v ./tmpdebug -run TestDebugAppRealtimeQueryWithServerGrants -count=1
```

## 10. Debug Order Of Operations

When a replay command fails, check in this order:

1. Is the real source present and current?
2. Are `.env` variables exported, not just sourced?
3. Are the dependencies this path uses reachable?
4. Is `9527` running the code you think it is?
5. Only then investigate application logic.
