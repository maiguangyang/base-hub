# Superpowers for Codex

You have superpowers.

This profile adapts Superpowers conventions for Codex with strict single-flow execution.

## Core Rules

1. Prefer local skills in `.agent/skills/<skill-name>/SKILL.md`.
2. Execute one core task at a time with one focused step or bounded edit set.
3. Use browser-specific tooling only for genuine browser automation tasks.
4. Track checklist progress in the resolved active-window tracker: default to `<project-root>/docs/plans/task.md`, but inside a linked worktree use `git rev-parse --git-path codex/task.md` instead.
5. Keep changes scoped to the requested task and verify before completion claims.
6. Keep all work in the main session unless the user explicitly requests Sub Agents or delegated-worker execution for the current task. Do not infer delegation permission from task size, parallelizability, expected speed, review value, prior requests, or available tools.

## Tool Translation Contract

When source skills reference legacy tool names, use these Codex equivalents:

- Legacy assistant/platform names -> `Codex`
- `Task` tool -> stay in the main session by default; `spawn_agent` is prohibited unless the user explicitly requests Sub Agents or delegated-worker execution for the current task
- `Skill` tool -> invoke the skill by name (`$skill-name`) and load `.agent/skills/<skill-name>/SKILL.md` when needed; second preference is `~/.codex/skills/<skill-name>/SKILL.md`
- `TodoWrite` -> update the resolved active-window tracker
- File reads -> `exec_command` (`sed`, `cat`, `rg`)
- File edits -> `apply_patch`
- Directory listing -> `exec_command` (`find`, `rg --files`, `ls`)
- Search -> `exec_command` (`rg`)
- Shell -> `exec_command`
- Web fetch/search -> `web`
- Image generation -> `generate_image`
- User communication during tasks -> commentary updates
- MCP tools -> `mcp_*` tool family

## Skill Loading

- First preference: project skills at `.agent/skills`.
- Second preference: user skills at `~/.codex/skills`.
- If both exist, project-local skills win for this profile.
- User-facing entry convention is `$skill-name`.

## Single-Flow Execution Model

- Default mode: do not create, spawn, invoke, or delegate to Sub Agents; keep all work in the current main session.
- Explicit delegation mode: enter it only when the user explicitly requests Sub Agents or delegated-worker execution for the current task.
- Delegation permission is task-scoped and does not carry forward to later requests.
- Do not infer delegation permission from task complexity, independent subtasks, time pressure, review needs, or the presence of a subagent-oriented skill.
- Decompose large work into ordered, explicit steps.
- Keep exactly one active task at a time in the resolved active-window tracker.
- Move stale execution windows and no-longer-active trackers under `<project-root>/docs/plans/archive/`, or under `git rev-parse --git-path codex/archive` when inside a linked worktree.
- If browser work is required, isolate it in a dedicated browser step.

## Verification Discipline

Before saying a task is done:

1. Run the relevant verification command(s).
2. Confirm exit status and key output.
3. Update the resolved active-window tracker.
4. Report evidence, then claim completion.
