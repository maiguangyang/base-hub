---
name: evaluating-technical-proposals
description: Use when reviewing a technical proposal, architecture design, migration plan, or implementation approach and a rigorous feasibility, risk, scalability, security, team-fit, or ROI assessment is needed before approval.
---

# Evaluating Technical Proposals

## Overview

When the user provides a technical proposal and wants a Principal / Senior Staff level review, treat the task as an approval gate, not a polite summary. Your job is to identify the failure modes that could cause delivery delays, runaway cost, operational fragility, security exposure, or long-term maintenance pain.

Do not invent facts. If traffic levels, data volume, dependency versions, team capacity, budget, or compliance scope are missing, state the gaps and continue under explicit assumptions.

When this skill is invoked on a design produced by `.agent/skills/brainstorming/SKILL.md`, it acts as the planning gate before the user's planning-entry choice: direct `$writing-plans`, or `$using-git-worktrees` followed by automatic continuation into `$writing-plans`. Report whether the design is blocked, needs revision, or is ready for manual planning handoff, but never trigger either option automatically before the user chooses.

## The Process

### Step 1: Establish Missing Context and Review Assumptions

1. Extract the stated business goal, MVP boundary, timeline, and success criteria.
2. List the missing information required to judge feasibility, such as scale, SLOs, data sensitivity, team size, dependency versions, budget, and compliance constraints.
3. If important facts are missing, create a `Missing Information and Review Assumptions` section before the main analysis.
4. Never fabricate specific CVEs, precise performance numbers, or detailed cost totals without supporting inputs.

### Step 2: Judge MVP and Business Fit First

1. Assess whether the proposal actually satisfies the MVP instead of solving an imaginary future platform problem.
2. Identify over-engineering, hidden scope expansion, or roadmap assumptions disguised as hard requirements.
3. Only after judging MVP fit, evaluate whether the design can evolve over the next one to three years without major rewrites.

### Step 3: Evaluate the Seven Core Dimensions

For each dimension, provide both `Feasibility Highlights` and `Potential Dangers`.

1. `Business and Product Fit`
   - Does the proposal solve the real user and business problem?
   - Is the MVP boundary clear and appropriately sized?
   - Is the design overbuilt for the current phase?
   - How well can it evolve over one to three years?

2. `Technical Architecture and Ecosystem`
   - Is the stack mature enough for the team and use case?
   - How well does it integrate with upstream and downstream systems?
   - Are there distributed consistency risks, single points of failure, or vendor lock-in concerns?

3. `Non-Functional Requirements (SLA / SLO)`
   - Where are the likely performance bottlenecks?
   - Does the design scale horizontally or vertically in a realistic way?
   - Are rate limiting, degradation paths, circuit breaking, failover, and disaster recovery addressed?

4. `Security and Compliance`
   - Is data classification clear?
   - Are encryption, masking, access control, and privilege boundaries handled correctly?
   - What security validation points must be checked if dependency versions are missing?
   - Are industry or regional compliance requirements likely to block launch?

5. `Operations and Observability (Day 2)`
   - How hard will it be to add logs, metrics, and tracing?
   - Does the deployment model support safe rollback and low-risk releases?
   - Is the system testable at the unit, integration, and production-readiness levels?

6. `Team and Execution`
   - Does the design match the actual skill profile of the team?
   - Is the learning curve acceptable for the delivery timeline?
   - Is there a high bus-factor risk around any critical subsystem?

7. `Cost and ROI (TCO)`
   - What is the likely implementation effort?
   - What infrastructure or managed-service costs are implied?
   - What long-term maintenance burden does the design create?
   - Is the expected value worth the complexity and spend?

### Step 4: Rank Risks by Severity

1. Classify each meaningful risk as `P0`, `P1`, or `P2`.
2. Every risk should include:
   - Trigger condition
   - Consequence
   - Severity
   - Mitigation
3. Use these severity definitions:
   - `P0`: Likely to cause project failure, system collapse, a major data incident, or a compliance blocker
   - `P1`: The solution may be viable, but this must be resolved before the milestone
   - `P2`: Important, but can be deferred as hardening or optimization work

### Step 5: Extract the Critical Red Flags

1. Pull out only the one to three risks most likely to kill the project or destabilize the system.
2. Explain why each red flag matters and what failure path it creates.
3. Prefer concrete failure scenarios over generic warnings.

### Step 6: Produce an Actionable Verdict

1. Choose exactly one conclusion:
   - `Strongly Recommended`
   - `Conditionally Viable`
   - `Not Recommended`
   - `Contains Fatal Flaws`
2. Match the conclusion to the evidence:
   - `Strongly Recommended`: No unresolved `P0` risk and only manageable `P1` risks
   - `Conditionally Viable`: Workable, but dependent on explicit preconditions or risk retirement
   - `Not Recommended`: Another simpler or lower-risk path is clearly better
   - `Contains Fatal Flaws`: An unresolved `P0` risk makes execution or operation unacceptable

### Step 7: Emit the Planning Gate Decision

1. End the review with exactly one planning gate status:
   - `BLOCKED`
   - `REVISE`
   - `READY_FOR_MANUAL_WRITE_PLAN`
2. Use these gate rules:
   - `BLOCKED`: Any unresolved `P0` exists, or the proposal is not credible enough to move into planning
   - `REVISE`: No `P0`, but one or more unresolved critical `P1` issues must be corrected before planning
   - `READY_FOR_MANUAL_WRITE_PLAN`: No unresolved `P0`, and no unresolved critical `P1` blocks remain
3. If the gate is `BLOCKED` or `REVISE`, explicitly state that the design should not proceed into either planning-entry option yet.
4. If the gate is `READY_FOR_MANUAL_WRITE_PLAN`, explicitly state that the user may manually choose either direct `$writing-plans`, or `$using-git-worktrees` first and then automatic continuation into `$writing-plans`.

## Review Rules

- Judge MVP fit before long-term extensibility.
- Every risk must include trigger, impact, severity, and mitigation.
- If dependency versions are missing, list the security checks instead of inventing CVE IDs.
- State the assumptions behind performance, capacity, and cost estimates.
- Prefer concrete, executable judgments over abstract commentary.
- If the proposal lacks key information, say so explicitly instead of filling the gaps with guesswork.
- When this review is acting as a pre-planning gate, never execute or imply automatic progression into either planning-entry option before the user chooses.

## Output Format

Return a structured Markdown report with these sections:

- `One-Sentence Executive Conclusion`
  - Choose one: `Strongly Recommended`, `Conditionally Viable`, `Not Recommended`, or `Contains Fatal Flaws`
- `Missing Information and Review Assumptions`
- `Detailed Analysis Across Seven Dimensions`
  - For each dimension, include `Feasibility Highlights` and `Potential Dangers`
- `Critical Red Flags`
  - List only the top one to three risks and label each as `P0` or `P1`
- `Prioritized Action Items`
  - Group them into `Immediate`, `This Phase`, and `Can Be Deferred`
- `PoC Recommendations`
  - Include validation goal, validation method, and pass criteria
- `Planning Gate Decision`
  - Choose one: `BLOCKED`, `REVISE`, or `READY_FOR_MANUAL_WRITE_PLAN`
- `Write-Plan Handoff`
  - State whether both planning-entry options should wait, whether the design must be revised first, or whether the user may manually choose direct `$writing-plans` versus `$using-git-worktrees` first
- `Final Recommendation`
  - State whether the proposal should proceed and what preconditions must be satisfied first
