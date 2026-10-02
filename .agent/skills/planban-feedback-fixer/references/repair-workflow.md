# Repair Workflow

## Batch Selection

Start with one cohesive module:

- Skill management: Skeleton generation, Capability requirements, prompt-only skill creation.
- User management: create/edit/delete users, roles on users, disabled/deleted state, search empty states, overview counts.
- Role management: role code, required fields, permission selection, delete semantics, status field.
- Company management: create/edit company fields, list formatting.
- UI polish: alignment, spacing, loading, copy.

Avoid mixing unrelated modules in one patch unless the fix is a shared component.

## Card Handling

For each card:

1. Read `issue.md`.
2. Inspect every image attachment.
3. Open video URL files only when the static description and images are insufficient.
4. Extract actual behavior, expected behavior, and acceptance criteria.
5. If expected behavior is empty, infer narrowly from title/problem and mark the inference.

## Fix Discipline

- Locate the existing component, route, API call, validation schema, or state selector before editing.
- Reuse existing UI components and form validation patterns.
- Add tests for data logic, validation, permissions, and state transitions when the repo has a test harness.
- For visual changes, use browser verification and screenshots.
- Keep commits or final summaries mapped to Planban card ids.

## Completion Criteria

A card is done only when:

- The described bug no longer reproduces.
- Existing tests still pass, or the unrun tests are clearly reported.
- Visual/UI fixes are checked in a browser if the app can run locally.
- Any unclear requirements are listed as open questions instead of silently guessed.
