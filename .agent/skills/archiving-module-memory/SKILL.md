---
name: archiving-module-memory
description: Use at the end of a feature implementation (via `$archiving-module-memory`) to summarize and archive the executed plan into the long-term module memory docs.
---

# Archiving Module Memory

## Overview

When the user is satisfied with a feature implementation and invokes `$archiving-module-memory`, you should document the completed work into the layered module memory system. `docs/modules/<module-name>.md` remains the canonical module entry, while dated history now lives under `docs/modules-archive/<module>/`. This preserves current facts in the brief, keeps deep history searchable, and prevents AI amnesia in future tasks.

## The Process

### Step 1: Identify the Affected Modules

1. Review the just-completed `plan` (e.g. from `docs/plans/`).
2. Identify which core business modules were touched or expanded (e.g., `user`, `auth`, `device`, `payment`).
3. If the plan spans multiple distinct modules, you will update multiple module files.

### Step 2: Locate or Create the Module Brief, and Update Current Facts

For each identified module:
1. Check if `docs/modules/<module-name>.md` exists.
2. If it does not exist, create it with a standard header (Title and Description), plus empty `## 领域核心约束 (Domain Constraints)`, `## 依赖关系 (Dependencies)`, `## 当前核心能力图谱 (Core Capabilities Index)`, and `## 历史实现功能日志 (Feature Logs)` blocks.
3. Review the newly completed feature. Did it establish any **Absolute Rules, Anti-Patterns, or Red Lines**? If so, append them to the `## 领域核心约束 (Domain Constraints)` list to warn future AIs.
4. Does this feature introduce a hard coupling with another module? If so, append it to `## 依赖关系 (Dependencies)`.
5. Update the brief-only sections with the **current, still-valid facts** introduced by the feature:
   - `## 当前核心能力图谱 (Core Capabilities Index)`
   - `## 已知坑点 (Known Pitfalls)` if present
   - `## 详细主题 (Detail Links)` if the feature adds or changes a durable topic deep-dive
6. Do **not** paste raw dated log bodies back into the brief file.

### Step 3: Archive the Dated Feature Log and Link It From the Brief

Create a new dated archive file under `docs/modules-archive/<module>/`.

The archive entry format MUST use Status Tags `[Active]`, `[Refactored]`, or `[Deprecated]`:
```markdown
### [Active] YYYY-MM-DD: [Feature Title]
- **关联计划**: [Plan File Name](../../plans/plan-file-name.md)
- **能力总结**:
  - [Point 1 summarizing the specific technical capabilities or business logic introduced]
  - [Point 2 ...]
```

Then add or update a link under the brief file's `## 历史实现功能日志 (Feature Logs)` section:

```markdown
- [YYYY-MM-DD: Feature Title](../modules-archive/<module>/YYYY-MM-DD-feature-title.md)
```

**Handling Old Logs (Knowledge Refactoring)**:
- If your current feature rewrites, deprecates, or overrides a historical feature, you **MUST** find the old archive file and change its tag from `[Active]` to `[Refactored]` or `[Deprecated]`.
- Add a bullet point under the old archive entry: `- **更新指引**: 详见 [YYYY-MM-DD: New Feature Title] 计划。`
- Update the brief's `## 当前核心能力图谱 (Core Capabilities Index)` so it reflects the new current truth, not the outdated one.

**Guidelines for the Summary**:
- Be concise but highly technical.
- Focus on *what* the system can now do and *what* rules were implemented (e.g., "Interceptor added for X", "Cron job pulls Y").
- Do not just list the files changed; focus on the **capabilities**.

### Step 4: Leave The Memory Update For Human Review

1. Save the updated module files.
2. Save the new or updated archive files in `docs/modules-archive/<module>/`.
3. Keep `## 历史实现功能日志 (Feature Logs)` in the brief as a link index, not as raw log bodies.
4. Leave the memory update in the local working tree for human review.
5. Do not run `git add` or `git commit`; the developer decides manual staging and commit timing.
6. Stop execution and notify the user that the task has been successfully archived.
