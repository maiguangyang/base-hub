# Plan Document Reviewer Prompt Template

Use this template when reviewing whether a finished implementation plan is ready for execution.

**Purpose:** Verify the plan is complete, matches the approved design/spec, and decomposes work into actionable tasks.

## Review Checklist

| Category | What to Look For |
| --- | --- |
| Completeness | TODOs, placeholders, missing steps, incomplete tasks |
| Alignment | Plan covers approved requirements without major scope creep |
| Decomposition | Tasks have clear boundaries and actionable steps |
| Buildability | An implementer can follow the plan without guessing |

## Approval Bar

Only flag issues that would cause real implementation problems:

- missing requirements
- contradictory steps
- vague instructions that block execution
- placeholders that hide real work

Minor wording preferences should not block approval.

## Output Format

### Plan Review

**Status:** Approved | Issues Found

**Issues (if any):**
- [Task X, Step Y]: [specific issue] - [why it matters for implementation]

**Recommendations (advisory):**
- [optional improvement suggestions]
