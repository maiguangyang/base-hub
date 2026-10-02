# Tmpdebug Index

Use this index to find the closest existing live debug probe before creating a new one.

## Preconditions

Many `tmpdebug` tests are not pure unit tests. Before running them, check:

- real source plist is available when the probe needs app-cached auth
- `.env` has been exported with `set -a; source .env; set +a` when the probe touches DB-backed engine code
- `DATABASE_URL` is exported for DB probes; Redis is required only for paths using it
- the engine listener on `9527` is the intended process when the probe talks to a live websocket endpoint

## Core Realtime / Conversation Probes

- `base-engine/tmpdebug/app_realtime_live_debug_test.go`
  - real websocket replay using App login, device Skill selection, and server grants
- `base-engine/tmpdebug/menu_query_live_debug_test.go`
  - menu query live conversation probing
- `base-engine/tmpdebug/menu_query_latency_live_debug_test.go`
  - websocket + engine latency decomposition for menu query
- `base-engine/tmpdebug/menu_query_latency_assertions_test.go`
  - guardrails for measured latency expectations
- `base-engine/tmpdebug/menu_query_latency_thresholds_test.go`
  - threshold-oriented verification helpers
- `base-engine/tmpdebug/menu_query_prompt_template_live_debug_test.go`
  - prompt-template behavior under live menu-query flow
- `base-engine/tmpdebug/menu_query_tool_result_rendering_live_debug_test.go`
  - live rendering / compaction behavior for tool result JSON
- `base-engine/tmpdebug/greeting_latency_live_debug_test.go`
  - simple greeting / no-tool path latency decomposition

## Runtime Header Probes

- `base-engine/tmpdebug/runtime_header_live_debug_test.go`
  - runtime header binding and lookup diagnostics
- `base-engine/tmpdebug/runtime_header_state_debug_test.go`
  - runtime header state debugging and persistence checks

## Suggested Pairings

Use these pairings as a default:

- Approval / first-card / first-text timing:
  - websocket replay + `menu_query_latency_live_debug_test.go`
- Greeting or no-tool slow path:
  - `greeting_latency_live_debug_test.go`
- Runtime header missing / incorrect binding:
  - `runtime_header_live_debug_test.go`
  - `runtime_header_state_debug_test.go`
- Tool result formatting / markdown table issues:
  - `menu_query_tool_result_rendering_live_debug_test.go`

## Rule Of Thumb

Reuse the closest existing tmpdebug test before inventing a new probe.

Only create a new probe when:

- the repo has no close match, or
- the existing probe cannot emit the evidence the user actually needs
