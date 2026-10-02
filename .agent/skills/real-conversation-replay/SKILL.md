---
name: real-conversation-replay
description: Use when the user asks for real app-to-engine simulation, real data replay, iOS simulator/app prefs extraction, runtimeMcpHeaders verification, websocket replay, or end-to-end debugging across app -> realtime -> engine.
---

# Real Conversation Replay

Use this skill when the user asks for any of the following:

- "用真实数据"
- "真实模拟"
- "从 app 到 engine"
- "iOS 模拟器"
- "runtime headers" / "runtimeMcpHeaders"
- "websocket replay"
- approval / interaction / tool loop / runtime-header / latency 的真实链路排查

This skill is repo-specific for `easypos-hub`. It is not a generic replay skill.

Load only the references you need:

- `references/app-prefs-and-simulator.md`
- `references/execution-prereqs.md`
- `references/realtime-websocket-replay.md`
- `references/tmpdebug-index.md`
- `references/evidence-checklist.md`

## Core Rule

Never claim that you performed a real replay unless you can name the real source explicitly, such as:

- `AI_POS_APP_PREFS_PLIST`
- a live iOS simulator app prefs plist
- another clearly identified live source

If the real source is unavailable, say so plainly and treat the task as blocked for true replay.

## Execution Prerequisites

Before running any replay, check these repo-specific prerequisites first:

- Real source is available:
  - prefer `AI_POS_APP_PREFS_PLIST`
  - otherwise resolve the current iOS simulator container for `com.proton.aiPos.dev`
- App-local cached auth exists in that plist:
  - `flutter.dev_token`
  - `flutter.dev_appInstallationId` and `flutter.dev_deviceSkillState`
  - Server-owned MCP grant for the real company member/application; old cached headers are intentionally absent.
- Shell environment is exported when the command needs child processes such as `go test`, `go run`, or `make`:
  - `source .env` alone is not enough
  - use `set -a; source .env; set +a`
- Local services are reachable when the probe touches engine DB-backed code or when starting engine:
  - `DATABASE_URL`
  - MySQL connectivity; Redis only when the selected execution path actually uses it.
- The websocket endpoint is the one you really intend to test:
  - prefer `ws://127.0.0.1:9527/realtime` for local workspace verification
  - if you are validating a fresh code change, make sure the running engine process is from the current workspace, not a stale earlier binary
- Replays with long planning stages may need a longer read deadline:
  - prefer `AI_POS_REALTIME_READ_TIMEOUT_SECONDS=45` unless you know a shorter timeout is safe

If any prerequisite is missing, say so early and fix the environment before debugging business logic.

## Current authorization contract

Runtime authorization is fixed `SERVER_ONLY`. Commands must not send cached `runtimeMcpHeaders` or authorization revisions. Read `flutter.dev_appInstallationId` into `session.init.deviceId`; compute `enabledSkillIds` as manual/group desired Skills intersected with available runtime descriptors, exactly as the App does. The public token authenticates the member; current server grant resolves credentials immediately before MCP I/O. Include the actual `input.agentId`, locale and device timezone.

For model errors, record the failed provider round, HTTP/transport error, serialized request size, message/tool counts and available token usage. Compare an equivalent small request and the failed full request against the same configured model/provider; do not attribute a generic connection error to context size without provider evidence.

## Workflow

1. Confirm the real source first.
   - Prefer `AI_POS_APP_PREFS_PLIST`.
   - Read the real token, installation ID and selected available Skill IDs from App prefs. Never reconstruct or require old `runtimeMcpHeaders`. Set the actual Agent ID and optional conversation ID explicitly.
   - If an old simulator plist path fails with `plutil` or no longer exists, re-resolve the current app data container instead of assuming the old path is still valid.
   - If the real source is missing, do not pretend the replay happened.

2. Normalize the execution environment.
   - For commands that spawn child processes, export repo env first:
     - `cd base-engine && set -a && source .env && set +a`
   - If the probe needs DB-backed engine code, verify `DATABASE_URL` is exported; require Redis only if that path needs it.
   - If you are using websocket replay against localhost, verify what process is listening on `9527` before trusting the result.

3. Choose the replay mode.
   - Use websocket replay for app-visible timing, approval/interaction card arrival, frame order, and wire-level payload verification.
   - Use engine `tmpdebug` in-process replay for `StartRealtimeTurn`, `RunRealtimeTurn`, tool loop, AI usage, and server grant resolution verification.
   - Use both when the user asks for full app -> engine analysis or when websocket behavior and engine behavior diverge.

4. Reuse existing repo probes before inventing new ones.
   - Check the tmpdebug index for the closest existing live debug test.
   - Extend an existing probe only when the repo does not already cover the scenario.

5. Collect evidence before concluding.
   - Capture protocol milestones such as:
     - `session.ready`
     - `response.create`
     - `approval_request`
     - `interaction_request`
     - `tool_execution_result`
     - first assistant text
     - `response.completed`
   - Capture engine milestones when applicable:
     - `StartRealtimeTurn`
     - `RunRealtimeTurn`
     - loop-run state transitions
     - AI usage token / latency records

6. Map the evidence back to code paths.
   - Correlate websocket behavior to app notifier / repository code.
   - Correlate engine behavior to `turn_handler`, tool loop, runtime-header, approval, and interaction code.
   - Do not guess root causes without an evidence-backed path.

7. Report the result with boundaries.
   - State the real data source.
   - State whether the replay was websocket, in-process, or both.
   - State the validated timing / protocol / state findings.
   - State what remains inferred or unverified.

## Replay Mode Guidance

Choose replay mode based on the question:

- Websocket replay:
  - approval card appearance timing
  - interaction card appearance timing
  - app-visible first text timing
  - wire protocol ordering
  - payload transport verification

- Engine in-process replay:
  - `StartRealtimeTurn` latency
  - `RunRealtimeTurn` latency
  - scoped tools construction
  - server grant resolution
  - loop-run state transitions
  - tool result rendering / continuation behavior

- Combined replay:
  - full app -> engine diagnosis
  - websocket timing does not match engine timing
  - UI authorization and server grant readiness disagree

## Evidence Requirements

Every real replay answer should name:

- the real data source
- the replay mode used
- the key ids involved when available:
  - `conversationId`
  - `responseId`
  - `requestId`
  - `interactionId`
- the key timing or protocol milestones collected
- whether the conclusion comes from websocket, engine, or both

## Hard Prohibitions

- Do not claim real replay from code inspection alone.
- Do not use fake tokens or fake runtime headers and call it "real data".
- Do not paste secrets such as raw tokens or sensitive header values into long-lived docs or summary output.
- Do not claim full app -> engine verification if you only ran engine-side replay.
- Do not ignore missing simulator / plist access and continue as if the replay succeeded.
- Do not waste time debugging business logic when the real blocker is missing env export, stale plist path, or localhost engine not actually running the current code.

## Expected Answer Shape

When reporting results from this skill, structure the answer as:

1. `数据来源`
2. `执行方式`
3. `关键证据`
4. `结论与边界`

Keep the conclusions strict:

- say what was directly verified
- say what is inferred
- say what was blocked
