# App Prefs And Simulator

This repo's iOS app bundle id is:

- `com.proton.aiPos.dev`

The preferred real source for live replay is:

- `AI_POS_APP_PREFS_PLIST`

Remember that simulator app container UUIDs are not stable. A previously valid plist path can become stale after reinstall / redeploy / simulator churn.

Important cached plist keys:

- `flutter.dev_token`
- `flutter.dev_appInstallationId`
- `flutter.dev_deviceSkillState`
- `flutter.dev_currentCompanyId` / `flutter.dev_companyMemberId`
- `flutter.dev_agent_floating_sessions_v1_item_current` for actual Agent/conversation IDs

## Preferred Access Pattern

1. Check whether `AI_POS_APP_PREFS_PLIST` is already provided.
2. If it is provided, treat it as the primary real source.
3. If the provided path fails to open or `plutil` fails, re-resolve the current simulator app container for `com.proton.aiPos.dev`.
4. If it is not provided, try to locate the simulator app container for `com.proton.aiPos.dev`.
5. If simulator access is unavailable, explicitly report that real replay is blocked.

## Useful Commands

Check the environment variable:

```bash
printenv AI_POS_APP_PREFS_PLIST
```

Parse the plist inside the probe; do not dump tokens into terminal logs. Resolve the source path:

```bash
test -f "$AI_POS_APP_PREFS_PLIST"
```

Try to locate the app container on the booted simulator:

```bash
xcrun simctl get_app_container booted com.proton.aiPos.dev data
```

Try to locate the app container on a specific simulator device:

```bash
xcrun simctl get_app_container <device-udid> com.proton.aiPos.dev data
```

Typical prefs location inside the app data container:

```text
.../Library/Preferences/com.proton.aiPos.dev.plist
```

Fallback search when the container path is unknown:

```bash
rg --files --hidden --no-ignore ~/Library/Developer/CoreSimulator/Devices -g com.proton.aiPos.dev.plist
```

Export the resolved plist path for reuse:

```bash
export AI_POS_APP_PREFS_PLIST='<container>/Library/Preferences/com.proton.aiPos.dev.plist'
```

## Existing Repo Readers

Existing tmpdebug probes already know how to read these values:

- `base-engine/tmpdebug/app_realtime_live_debug_test.go`
- `base-engine/tmpdebug/menu_query_latency_live_debug_test.go`

Prefer reusing those readers instead of rewriting plist parsing logic from scratch.

## Common Failure Modes

- `plutil` exits with status `1`
  - usually means the plist path is stale
  - re-run `xcrun simctl get_app_container ... data`
- plist exists but `flutter.dev_token` is empty
  - app has not logged in or cached auth yet
- plist exists but `flutter.dev_mcpRuntimeHeaders` is empty
  - expected in the current App; Engine uses server grants and the replay must not require this cache
- old hard-coded app container path from a previous session
  - do not trust it without re-checking

## Reporting Rules

When reporting findings, describe the source as one of:

- `AI_POS_APP_PREFS_PLIST`
- `iOS simulator app prefs plist`
- `real source unavailable`

Do not print raw token values or sensitive runtime-header values in durable outputs.
