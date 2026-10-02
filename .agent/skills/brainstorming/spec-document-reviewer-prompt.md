# Design Document Reviewer Prompt Template

Use this template when reviewing whether a finished design document is ready for planning.

**Purpose:** Verify the design doc is complete, internally consistent, and clear enough to support implementation planning.

## Review Checklist

| Category | What to Look For |
| --- | --- |
| Completeness | TODOs, placeholders, missing sections, incomplete requirements |
| Consistency | Internal contradictions or conflicting constraints |
| Clarity | Requirements ambiguous enough to cause the wrong implementation |
| Scope | Too broad for a single implementation plan, or hiding multiple independent subsystems |
| YAGNI | Unrequested features or clear over-engineering |

## Approval Bar

Only flag issues that would cause real planning problems:

- a missing requirement that blocks task decomposition
- a contradiction that changes architecture
- ambiguity that could produce the wrong plan
- oversized scope that should be decomposed first

Minor wording preferences should not block approval.

## Output Format

### Design Review

**Status:** Approved | Issues Found

**Issues (if any):**
- [Section X]: [specific issue] - [why it matters for planning]

**Recommendations (advisory):**
- [optional improvement suggestions]
