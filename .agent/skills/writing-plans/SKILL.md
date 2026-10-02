---
name: writing-plans
description: Use when you have a spec or requirements for a multi-step task, before touching code
---

# Writing Plans

## Overview

Write comprehensive implementation plans assuming the engineer has zero context for our codebase and questionable taste. Document everything they need to know: which files to touch for each task, code, testing, docs they might need to check, how to test it. Give them the whole plan as bite-sized tasks. DRY. YAGNI. TDD. Leave changes local for human review.

Comprehensive means precise enough to execute, not broader than the approved scope.

Assume they are a skilled developer, but know almost nothing about our toolset or problem domain. Assume they don't know good test design very well.

**Announce at start:** "I'm using the writing-plans skill to create the implementation plan."

**Context:** Run in the current workspace by default. If the user already chose `$using-git-worktrees`, or you are already inside an isolated linked worktree, continue there and do not ask about isolation again. Otherwise, you MUST ask once whether the user wants an isolated worktree before continuing. If they opt in, use `$using-git-worktrees` first. If they decline or have already told you to skip worktrees, continue in the current workspace and do not ask again in downstream planning or execution skills.
**Prerequisite:** Read `.agent/memory/constitution-core.md`, then the relevant `rules-core/*`, then `docs/memory/index.yaml`, then the relevant `docs/modules/<module>.md` brief files. Open `docs/modules-details/...` only when the brief is insufficient.

**Save plans to:** `docs/plans/YYYY-MM-DD-<feature-name>.md`

## Delegation Policy

Default to single-agent execution for all planning, review, and execution routes defined by this skill.

- If the user has not explicitly requested Sub Agents or delegated-worker execution for the current task, NEVER create, spawn, invoke, delegate to, or recommend a Sub Agent.
- Use `$subagent-driven-development` only when the user explicitly requests delegated workers for the current task.
- Do not infer delegation permission from task size, independent subtasks, expected speed, review value, prior requests, or available tools.
- Delegation permission is task-scoped and does not carry forward to later requests.
- Every generated plan MUST repeat the applicable execution route in its **Global Constraints** section.

## Scope Check

If the spec covers multiple independent subsystems, suggest splitting it into separate plans so each plan produces working, testable software on its own. Each plan should be independently buildable, reviewable, and verifiable.

## No Overdesign or Redundant Design

Plan only the smallest implementation that satisfies the approved requirements and verified project constraints. Every task, file change, abstraction, interface, dependency, configuration option, compatibility path, and refactor must map to a current requirement, an observed codebase constraint, or a concrete risk that must be handled now.

- Reuse existing project patterns and boundaries before creating new ones.
- Do not plan speculative extension points, generic frameworks, future-proofing layers, duplicate fallbacks, parallel implementations, or abstractions for hypothetical reuse.
- Do not include unrelated cleanup, architecture modernization, optional enhancements, or refactors that are not required for the deliverable.
- Do not create extra files or layers solely to make the solution appear cleaner; split only when the current change gains clear ownership, contract clarity, or independent testability.
- Do not add tasks for optional hardening, observability, documentation, or configurability unless the spec, project rules, or a demonstrated risk requires them.
- When multiple implementations meet the requirements, plan the one with fewer new concepts, files, dependencies, and runtime paths.

Before finalizing, remove any task or design element whose absence would not block a stated requirement, violate a project rule, or leave a demonstrated risk unaddressed.

## File Structure

Before defining tasks, map out which files will be created or modified and what each one is responsible for. This is where decomposition decisions get locked in.

- Design units with clear boundaries and well-defined interfaces. Each file should have one clear responsibility.
- Prefer smaller, focused files over large files that do too much; if an existing file you must modify is already unwieldy, include only the targeted split needed for this work.
- Do not introduce a new file or interface when an existing boundary can own the behavior cleanly.
- Files that change together should live together. Split by responsibility, not by technical layer alone.
- In existing codebases, follow established patterns instead of forcing unrelated restructuring.

This structure informs task decomposition. Each task should produce self-contained changes that make sense independently.

## Task Right-Sizing

A task is the smallest unit that carries its own test cycle and is worth a fresh review gate. Fold setup, configuration, scaffolding, and documentation into the task whose deliverable needs them; split only where a reviewer could meaningfully reject one task while approving its neighbor. Each task ends with an independently testable deliverable.

## Engine Hardcoding Guardrail

For any plan that touches `base-engine` or engine-side conversation/runtime code, include an explicit **No Engine Hardcoding** requirement before the task breakdown.

Plans MUST prohibit runtime hardcoding in engine code:

- no hardcoded user-facing or business prose in runtime implementation files
- no product, store, tenant, or skill-specific literals to drive behavior
- no phrase deny/allow lists or substring matching against natural-language model output
- no fixed natural-language failure reasons, summaries, or status notes
- no model-output post-processing based on human-language wording

Plans MAY allow protocol/config constants when they are not user-facing prose:

- error codes, enum values, feature flag names, JSON field names, event/item IDs, database column names, and test fixtures
- prompt assets under `system_prompt/` when the goal is to constrain model behavior

For engine behavior that needs determinism, plan structured metadata instead of text matching: error codes, reason codes, typed enums, schema fields, or persisted state. User-facing wording should come from System Prompt constrained LLM output or existing localization/UI layers, not ad hoc runtime strings.

Every engine-touching plan MUST include a hardcoding audit step with exact commands, for example:

```bash
rg '"[^"]*[\p{Han}][^"]*"' base-engine/src --glob '!**/*_test.go'
rg '"[^"]{12,}"' base-engine/src --glob '!**/*_test.go'
rg 'strings\.(Contains|ContainsAny|HasPrefix|HasSuffix|EqualFold)\(.*(summary|message|text|delta|content|reason|output)' base-engine/src --glob '!**/*_test.go'
rg '\[\]string\{|map\[string\](bool|struct\{\}|string)\{' base-engine/src --glob '!**/*_test.go'
```

Expected: no runtime business/user-facing hardcoding in engine implementation files; hits are either removed or explicitly classified as protocol/config constants. Tests are reviewed separately.

## Bite-Sized Task Granularity

**Each step is one action (2-5 minutes):**

- "Write the failing test" - step
- "Run it to make sure it fails" - step
- "Implement the minimal code to make the test pass" - step
- "Run the tests and make sure they pass" - step
- "Leave changes local for human review" - step

## Plan Document Header

**Every plan MUST start with this header:**

```markdown
# [Feature Name] Implementation Plan

> **For this repo:** EXECUTION ROUTE: Default to the current agent with `$executing-plans`. Use `$subagent-driven-development` only when the user explicitly requests Sub Agents or delegated-worker execution for the current task; never infer delegation permission.

**Goal:** [One sentence describing what this builds]

**Architecture:** [2-3 sentences about approach]

**Tech Stack:** [Key technologies/libraries]

## Global Constraints

[Execution route for this task: state whether the user explicitly requested Sub Agents or delegated-worker execution. If not, Sub Agents are prohibited and every task must run sequentially in the current agent with `$executing-plans`.]

[Project-wide requirements from the spec: version floors, dependency limits, naming/copy rules, platform requirements, security constraints, and repo-specific rules. Copy exact values verbatim. Every task implicitly inherits these constraints.]

---
```

## Task Structure

````markdown
### Task N: [Component Name]

**Files:**

- Create: `exact/path/to/file.py`
- Modify: `exact/path/to/existing.py:123-145`
- Test: `tests/exact/path/to/test.py`

**Interfaces:**

- Consumes: [what this task uses from earlier tasks — exact signatures, schemas, or contract names]
- Produces: [what later tasks rely on — exact function names, parameter/return types, events, fields, or artifacts]

**Step 1: Write the failing test**

```python
def test_specific_behavior():
    result = function(input)
    assert result == expected
```

**Step 2: Run test to verify it fails**

Run: `pytest tests/path/test.py::test_name -v`
Expected: FAIL with "function not defined"

**Step 3: Write minimal implementation**

```python
def function(input):
    return expected
```

**Step 4: Run test to verify it passes**

Run: `pytest tests/path/test.py::test_name -v`
Expected: PASS

**Step 5: Leave changes local for human review**

Run:

```bash
git diff --cached --name-only
git status --short
```

Expected:

- no staged files
- only local working tree changes for the task remain
- human review and manual staging/commit happen later
````

## Figma-Aware Task Generation

When the input design doc contains a **Figma References** section (produced by the brainstorming skill):

1. For every task that implements UI covered by a Figma link, add a **Figma** annotation at the top of the task:
   ```markdown
   ### Task N: [Component Name]

   > **Figma:** Use `$figma` and follow its required flow for [link].
   ```
2. The annotation tells the executing engineer (or AI) to invoke the Figma MCP skill (`get_design_context` → `get_metadata` → `get_screenshot`) before writing code for that task.
3. If a single Figma link covers multiple tasks, reference it in each relevant task — do not assume the executor will remember context from a previous task.
4. Tasks without UI work should **not** receive a Figma annotation.

## No Placeholders

These are plan failures and must not appear in the final plan:

- `TBD`, `TODO`, `implement later`, `fill in details`
- `add appropriate error handling`, `add validation`, `handle edge cases`
- `write tests for the above` without actual test code
- `similar to Task N` instead of repeating the required code or command
- steps that describe what to do without showing how
- references to types, functions, or methods not defined in the plan
- later tasks consuming names, shapes, or contracts that no earlier task produces

## Remember

- Exact file paths always
- Complete code in plan (not "add validation")
- Exact commands with expected output
- Explicit interfaces for cross-task dependencies
- Reference relevant skills with `$skill-name` syntax
- DRY, YAGNI, TDD, explicit human-review handoff

## Self-Review

After writing the plan, review it against the source spec or approved requirements:

1. **Coverage:** every requirement maps to at least one task
2. **Placeholder scan:** remove vagueness and incomplete steps
3. **Type and naming consistency:** later tasks must match earlier definitions
4. **Buildability:** an implementer should be able to follow the plan without guessing
5. **Engine hardcoding:** if the plan touches `base-engine`, confirm the **No Engine Hardcoding** requirement and audit commands are present
6. **Overdesign audit:** every task and new design element maps to a current requirement, project constraint, or demonstrated risk; remove speculative, redundant, or unrelated work

If useful, use `writing-plans/plan-document-reviewer-prompt.md` as a structured review checklist.

## Execution Handoff

After saving the plan, hand off to exactly one execution path.

When the user explicitly requested delegated workers for the current task:

**"Plan complete and saved to `docs/plans/<filename>.md`.**
**Next step: use `$subagent-driven-development` for the explicitly requested delegated-worker execution."**

Otherwise:

**"Plan complete and saved to `docs/plans/<filename>.md`.**
**Next step: use `$executing-plans` for local single-flow execution. Sub Agents are prohibited because delegation was not explicitly requested for this task."**

Execution requirements:

- Treat the absence of an explicit current-task delegation request as a prohibition on Sub Agents.
- Use or recommend `$subagent-driven-development` only when the user explicitly requests delegated workers for the current task.
- Otherwise instruct the user to use `$executing-plans` and do not offer delegated execution as an option.
- Do NOT proceed to execute the plan yourself until the user explicitly chooses the applicable execution skill.
