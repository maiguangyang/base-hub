---
name: using-superpowers
description: Use when starting any conversation - establishes how to find and use skills, requiring skill invocation before ANY response including clarifying questions
---

<SUBAGENT-STOP>
If you were dispatched as a subagent to execute a specific task, skip this skill.
</SUBAGENT-STOP>

<EXTREMELY-IMPORTANT>
If you think there is even a 1% chance a skill might apply to what you are doing, you ABSOLUTELY MUST invoke the skill.

This threshold requires checking and loading the skill. It does not require choosing the skill's heaviest workflow. When a skill provides scope-based branches, use the smallest branch that fits the request and evidence.

IF A SKILL APPLIES TO YOUR TASK, YOU DO NOT HAVE A CHOICE. YOU MUST USE IT.

This is not negotiable. This is not optional. You cannot rationalize your way out of this.
</EXTREMELY-IMPORTANT>

## Instruction Priority

Skills define process, but explicit project and user instructions still win:

1. **User and project instructions** (`AGENTS.md`, direct user requests, repo-specific rules)
2. **Skills** - these override default behavior when applicable
3. **Default system behavior**

If repository instructions conflict with a generic skill expectation, follow the repository/user instruction and adapt the skill accordingly.

## How to Access Skills

In this repo, use `$skill-name` as the user-facing shorthand. Internally, load `.agent/skills/<skill-name>/SKILL.md` or `~/.codex/skills/<skill-name>/SKILL.md` when needed, then follow it directly.

## The Rule

Invoke relevant or requested skills before any response or action, including clarifying questions, code exploration, or file checks. A 1% chance is enough to load the skill and check it. If the skill turns out not to apply, say so briefly and continue.

Announce the active skill in one short line: `Using $skill-name to ...`. Create or update an active-window tracker only when the selected branch is genuinely multi-step and the skill requires tracking; do not create one for a lightweight path.

Before updating any checklist, resolve the active-window tracker path:

- Default workspace: `<project-root>/docs/plans/task.md`
- Linked worktree: `git rev-parse --git-path codex/task.md`
- Default archive path: `<project-root>/docs/plans/archive/`
- Linked worktree archive path: `git rev-parse --git-path codex/archive`

If the tracker file is missing, create it at the resolved location as a table-only active-window task list. Archive stale execution windows under the resolved archive path instead of letting old task trackers accumulate in the default read path.

## Platform Adaptation

Skill bodies may mention tool names from other hosts. When needed, map them using:

- `references/codex-tools.md`
- `references/copilot-tools.md`
- `references/gemini-tools.md`

## Red Flags

These thoughts mean STOP—you're rationalizing:

| Thought | Reality |
|---------|---------|
| "This is just a simple question" | Questions are tasks. Check for skills. |
| "I need more context first" | Skill check comes BEFORE clarifying questions. |
| "Let me explore the codebase first" | Skills tell you HOW to explore. Check first. |
| "I can check git/files quickly" | Files lack conversation context. Check for skills. |
| "Let me gather information first" | Skills tell you HOW to gather information. |
| "This doesn't need a formal skill" | If a skill exists, use it. |
| "I remember this skill" | Skills evolve. Read current version. |
| "This doesn't count as a task" | Action = task. Check for skills. |
| "The skill is overkill" | Do not skip checking the skill; select its smallest valid branch and escalate only for concrete risk or ambiguity. |
| "I'll just do this one thing first" | Check BEFORE doing anything. |
| "This feels productive" | Undisciplined action wastes time. Skills prevent this. |
| "I know what that means" | Knowing the concept ≠ using the skill. Invoke it. |

## Skill Priority

When multiple skills could apply, use this order:

1. **Process skills first** (brainstorming, debugging) - these determine HOW to approach the task
2. **Implementation skills second** (`$ui-ux-pro-max`, `$figma`, `$executing-plans`, `$subagent-driven-development`) - these guide execution

"Let's build X" → brainstorming first; scope classification decides between its lightweight path and full design path.
"Fix this bug" → debugging first, then domain-specific skills.

For plan execution:
- If the user explicitly wants delegated workers or subagent-style execution, prefer `$subagent-driven-development`
- Otherwise prefer `$executing-plans` or `$single-flow-task-execution`, depending on whether batch checkpoints are needed

## Skill Types

**Rigid** (TDD, debugging): Follow exactly. Don't adapt away discipline.

**Flexible** (patterns): Adapt principles to context.

The skill itself tells you which.

## User Instructions

Instructions usually say WHAT, not HOW, but explicit constraints such as “minimal change”, “no overdesign”, or a direct skill invocation also constrain the process. Follow them unless they conflict with a higher-priority safety or repository rule.
