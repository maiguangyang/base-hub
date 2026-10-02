# ADK-Go Core Rules

This is the default execution entrypoint for work involving Google ADK-Go, `adk-go/`, or ADK-style agent runtime design. The full source of truth remains [../rules/adk-go.md](../rules/adk-go.md).

- ADK-SOURCE-001
  Trigger: discuss, design, implement, or evaluate anything related to ADK-Go
  Must:
  - inspect the local `adk-go/` source before making implementation claims
  - read the relevant local package, not only README summaries
  - consult official ADK documentation when the task depends on current API usage, setup, deployment, MCP tools, callbacks, workflows, or version-specific behavior
  - treat official docs as version-sensitive and prefer current official docs over model memory

- ADK-HERMES-002
  Trigger: compare ADK-Go with Hermes or attempt Hermes-like runtime behavior
  Must:
  - not assume ADK-Go provides Hermes behavior out of the box
  - explicitly account for missing Hermes-like layers: progressive tool disclosure, session-scoped tool catalog, prompt-cache discipline, result budgeting, safe parallel execution, and production MCP lifecycle
  - keep `tool_search / tool_describe / tool_call` as an added runtime strategy, not an assumed ADK-Go primitive

- ADK-MCP-003
  Trigger: use ADK-Go with MCP tools
  Must:
  - inspect `adk-go/tool/mcptoolset/` before designing MCP behavior
  - verify whether the proposed flow exposes all MCP tools directly or wraps them behind a directory/bridge layer
  - avoid copying a full MCP tool schema surface into every model request when the goal is Hermes-like performance

- ADK-HITL-004
  Trigger: add approval cards, selection cards, confirmation flows, or human input
  Must:
  - inspect ADK-Go confirmation, callback, plugin, workflow, and event primitives before inventing a new control path
  - keep HITL as runtime policy/event handling, not as a prerequisite product catalog configuration

- ADK-VERIFY-005
  Trigger: modify ADK-Go integration code, wrappers, or local experiments
  Must:
  - run the narrowest relevant Go tests or verification command before claiming correctness
  - if tests cannot run, state the blocker and the exact command that should be run
