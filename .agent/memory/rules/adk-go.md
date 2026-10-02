# ADK-Go Rules

**Version**: 0.1.0

> This document is the full ADK-Go appendix. The default lightweight entrypoint is [../rules-core/adk-go.md](../rules-core/adk-go.md). Use this file when the core rules are insufficient or when designing an ADK-Go-based agent runtime.

## 1. Purpose

These rules exist so future ADK-Go work does not rely on stale model memory.

When a task involves ADK-Go, Google ADK, or an ADK-Go-based agent runtime, the agent must ground itself in:

- local source: `adk-go/`
- official docs: https://adk.dev/get-started/
- package/API reference linked from the official docs
- project analysis notes: `todos/hermes-adk-go-agent-runtime-analysis.md`

## 2. Source-First Rule

ADK-Go is an external library and may evolve quickly. Do not answer implementation questions from memory alone.

Required local source entrypoints:

| Concern | Read first |
| --- | --- |
| Agent definition and subagents | `adk-go/agent/llmagent/` |
| Runner/session lifecycle | `adk-go/runner/`, `adk-go/session/` |
| LLM flow and function calls | `adk-go/internal/llminternal/` |
| Tool abstractions | `adk-go/tool/` |
| Go function tools | `adk-go/tool/functiontool/` |
| MCP integration | `adk-go/tool/mcptoolset/` |
| Human confirmation / callbacks | `adk-go/tool/`, `adk-go/plugin/`, `adk-go/cmd/launcher/console/` |
| Workflow / graph execution | `adk-go/workflow/` |
| REST/SSE/WebSocket/A2A server behavior | `adk-go/server/` |
| Examples | `adk-go/examples/` |

If local source and docs disagree, determine whether the local clone is older/newer than the docs before deciding.

## 3. Official Documentation Rule

The official ADK docs are live and version-sensitive. Consult them when a task depends on:

- installation or quickstart behavior
- current Go API usage
- MCP tools
- custom tools
- action confirmations
- sessions, memory, artifacts, events
- plugins and callbacks
- graph workflows
- runtime/server/deployment
- observability/evaluation/safety

Do not freeze volatile documentation details into this memory file. Link to the docs and re-check them for current behavior.

## 4. ADK-Go Is Not Hermes

ADK-Go is a Go agent SDK. Hermes-agent is a product-level agent runtime with additional strategy and policy layers.

Do not claim ADK-Go can match Hermes out of the box unless the following gaps are explicitly handled:

- progressive tool disclosure
- `tool_search / tool_describe / tool_call`
- session-scoped tool catalog
- raw MCP runtime registry
- prompt-cache-stable request building
- safe parallel tool execution
- result budgeting and large-output spooling
- production MCP lifecycle management
- approval/HITL as runtime policy
- schema sanitization and model/provider compatibility handling

The correct framing is:

```text
ADK-Go = Go SDK foundation
Hermes = runtime strategy and product behavior
```

## 5. Hermes-Like ADK-Go Adoption Pattern

When using ADK-Go to pursue Hermes-like behavior, prefer a wrapper strategy:

```text
ADK Runner / Session / Event
  -> DirectoryToolset
  -> ToolRegistry
  -> ToolDirectory
  -> Bridge tools
       tool_search
       tool_describe
       tool_call
  -> underlying ADK Tool / MCP Tool
```

The model-facing tool surface should be intentionally small:

- visible core tools
- `tool_search`
- `tool_describe`
- `tool_call`

The runtime-owned surface may contain:

- MCP tools
- plugin/internal tools
- full schemas
- toolset scope
- tenant/user authorization scope
- registry generation

`tool_call` must never be a bypass. It must route through the same approval, audit, callback, result-budget, and authorization path as a direct underlying tool call.

## 6. MCP Rule

Before designing ADK-Go MCP behavior, inspect `adk-go/tool/mcptoolset/`.

Known baseline from local source analysis:

- `mcptoolset` can list MCP tools and convert them into ADK tools.
- MCP schemas can be mapped into function declarations.
- tool execution calls MCP `CallTool`.
- confirmation support exists at the tool level.
- connection refresh/retry exists for some transport/session failures.

Known gap for Hermes-like behavior:

- default MCP tool exposure is too direct for large tool surfaces.
- a DirectoryToolset or equivalent wrapper is needed to avoid sending every MCP tool schema to the model each turn.

## 7. HITL and Approval Rule

For approval cards, selection cards, and human input, first inspect ADK-Go's existing primitives:

- tool confirmation
- request confirmation
- plugin callbacks
- workflow pause/resume or request input
- runner/session event stream
- console/server HITL handling

Human approval should be modeled as runtime policy and event handling. It should not require pre-creating product skills, catalogs, or templates before raw MCP tools can be used.

## 8. Prompt Cache and Context Rule

Do not assume ADK-Go provides Hermes-style prompt-cache discipline.

When building on ADK-Go, explicitly design:

- stable system prompt assembly
- dynamic context outside the stable prompt where possible
- small stable model-facing tool schemas
- deterministic tool definitions for cache-friendly providers
- result-size controls before appending tool output back into context

## 9. Parallel Execution Rule

ADK-Go can execute multiple function calls concurrently. Treat that as a performance primitive, not a complete safety policy.

Before enabling or relying on concurrency, classify tools:

- read-only and independent
- path/resource-scoped
- write or mutation tools
- destructive tools
- MCP tools with explicit parallel-safety guarantees

Default to conservative execution for mutation tools.

## 10. Verification Rule

If changing ADK-Go integration code or local wrapper experiments, run the narrowest relevant verification command.

Examples:

```bash
cd adk-go && go test ./tool/mcptoolset/...
cd adk-go && go test ./internal/llminternal/...
cd adk-go && go test ./...
```

If verification cannot run because dependencies, credentials, network, or environment are missing, report the blocker and the exact command that should be run.

## 11. Documentation Boundary

Keep this file as a rule and navigation document.

Do not paste large official documentation excerpts here. Prefer:

- stable principles
- local source entrypoints
- adoption constraints
- links to official docs
- links to project analysis notes
