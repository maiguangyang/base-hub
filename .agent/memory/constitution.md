# AIPos Constitution

**Version**: 14.1.0

> [!IMPORTANT]
> **AI Directive**: When working on backend tasks, you **MUST** read [rules/base-engine.md](./rules/base-engine.md) before writing any code. When working on frontend tasks, you **MUST** read [rules/base-web.md](./rules/base-web.md) before writing any code. When working on Flutter app tasks, you **MUST** read [rules/base-app.md](./rules/base-app.md) before writing any code. When working on full-stack tasks, read **ALL** relevant rule documents.

## Priority Zero: Git Staging Index Boundary

### One-operation exception: MCP local integration (2026-09-07)

The user explicitly authorized this exception after being shown the index/commit prohibition. It applies only to landing the reviewed and archived `merge/generic-mcp-redirect-authorization-1-1-0-rc` worktree into `feature/instruction-1.1.0-rc` in the root, base-engine, base-app and base-web repositories. For this operation only, the agent may stage explicit task-owned paths, commit the reviewed changes and necessary conflict resolutions, and perform ordinary local merges and verified worktree cleanup. This exception takes precedence over the index prohibition and Manual Review Commit Protocol below only within that scope.

No remote push, unrelated changes, amendment of existing commits, destructive reset, or alteration of pre-existing user-staged entries is authorized. Preserve both branches' valid behavior, verify before source commits and after integration, and retain the source worktree until merged verification passes. The exception expires when this single integration is completed; all later work remains under the default prohibition.

### Standing workflow exception: local-worktree GitLab publication (2026-09-11)

The user explicitly authorized `creating-gitlab-issues-and-merge-requests` to convert confirmed local feature work into linked GitLab objects. This exception activates only inside that Skill's `local-worktree` publication mode, only for `base-app`, `base-engine`, and `base-web`, and only after a no-mutation preview lists every included path and a separate explicit user confirmation accepts that exact snapshot.

For one confirmed snapshot per repository, the workflow may create and switch to `feature/<issue_iid>-<branchSlug>`, stage exactly the confirmed paths, create one new commit whose parent and tree match the preview, and perform one normal non-force upstream push. It must verify branch, base SHA, path set, prospective tree, commit parent/tree, ownership trailer, clean index/worktree, and remote SHA at their defined gates. Any drift or foreign collision stops the workflow.

This exception never covers the root repository, any path outside the confirmed preview, any pre-existing unrelated staged entry, force-push, amend, rebase, reset, restore, clean, stash, deletion, tag changes, merge, automatic rollback, or Git mutation by the existing pushed-branch and Issue-only modes. Partial failure must preserve and report the current state for explicit lossless recovery. Outside this exact confirmed operation, Priority Zero and the Manual Review Commit Protocol remain fully binding.
### One-operation exception — reviewed App client web tools merge (2026-09-11)

The user explicitly authorized task-only staging, commits, normal local merges and necessary conflict resolution across the root, App, Engine and Web repositories, creating `feature/app-client-web-tools-1.1.3-rc` from each repository's local `1.1.3-rc`. This exception covers only the reviewed/archived `.worktrees/app-client-web-tools` implementation and its integration evidence. Preserve both sides' unrelated content, existing user-staged entries and main checkouts. Repeated merged verification must pass before source-worktree cleanup. No remote push, force reset, broad cleanup or unrelated changes are authorized. The exception expires when this integration completes.

AI / agent **MUST NOT** modify the Git staging index.

*   **No Automatic Staging**: AI / agent **MUST NOT** run `git add` or any equivalent command that stages modified or newly created files.
*   **No Automatic Unstaging**: AI / agent **MUST NOT** run `git restore --staged`, `git reset`, `git rm --cached`, or any equivalent command that removes files from the staging index.
*   **Preserve User-Staged Files**: If the user manually staged a file, AI / agent **MUST** preserve that staged state exactly and must not unstage, restage, rewrite the index entry, or otherwise alter staging state for that file.
*   **Working Tree Only**: All AI / agent edits must remain as working-tree changes for human review unless a later constitution amendment explicitly changes this boundary.
*   **Manual Index Ownership**: Developers own staging and unstaging decisions. AI / agent may report staged/unstaged status, but must not change it.

## Engine Priority Zero: Entity Relationship Boundary

This is the highest-priority schema-modeling rule for `base-engine` and is a hard gate for specification, planning, review, and implementation.

*   **Relationship-Only Modeling**: Every relationship between two `@entity` types in handwritten `base-engine/model/*.graphql` **MUST** be declared as typed entity object/list fields using `@relationship` and a correct inverse. Entity-to-entity linkage must remain visible in the schema as a relationship, not as an untyped scalar convention.
*   **No Custom Foreign-Key Scalars**: Handwritten entity schemas **MUST NOT** declare or use custom scalar `xxxId`, `xxxIds`, `<entity>Id`, `<entity>Ids`, or differently named scalar aliases to associate, join, query, or enforce referential linkage to another entity/table. Renaming the foreign-key field does not make it compliant.
*   **No Entity-Class Exemptions**: Runtime, history, audit, ledger, session, governance, and explicit association entities are all covered. When an explicit association entity is justified because the relationship carries independent business facts, each link from that association entity to its endpoint entities must still use `@relationship`.
*   **Generated Artifacts Are Not Precedent**: Dolphin-generated FK columns, generated `...Id` / `...Ids` relationship projections, and generated mutation inputs are implementation artifacts. They may be consumed only through the generated contract and **MUST NOT** be copied into handwritten `@entity` schema as custom relationship fields.
*   **API IDs Are Not Persistence Relationships**: IDs in custom operation inputs or payloads may identify request targets, but they must not be persisted or reused as a substitute for an entity `@relationship`. A supposedly external or snapshot identifier that is used to look up, join, or navigate to another project entity is a relationship and is governed by this rule.
*   **Fail Closed**: Any spec, plan, schema review, or implementation containing a handwritten scalar cross-entity link is non-compliant and must stop for correction. Engine schema contract tests must reject custom scalar cross-entity links rather than relying on reviewer memory.

## Core Principles

### I. Spec-Driven Development Protocol (SDD)
This project strictly follows the **Spec-Driven Development** workflow.
*   **Single Source of Truth**: `spec.md` (Intent) → `plan.md` (Implementation) → Code.
*   **No Skipping Steps**: The workflow must be **Specify → Plan → Tasks → Implement**.
*   **Non-Negotiable**: **NEVER** skip the `plan.md` phase. **NEVER** write code without an approved task.

### II. Monorepo & Context Awareness
This project is a **Monorepo**.
*   **Root**: Documentation & Orchestration.
*   **`/base-engine`**: Backend base-engine (Golang/Dolphin). → [👉 Backend Rules](./rules/base-engine.md)
*   **`/base-web`**: Frontend Application (Astro/React). → [👉 Frontend Rules](./rules/base-web.md)
*   **`/base-app`**: Flutter App (macOS/Windows/iOS/Android). → [👉 App Rules](./rules/base-app.md)
*   **Iron Rule**: Before executing ANY command, verify the current directory. You MUST explicitly prepend commands with `cd base-engine && ...`, `cd base-web && ...`, or `cd base-app && ...`.
*   **Scope Detection**: Determine which rule document applies:
    - GraphQL schema changes, API logic, database operations → **Backend** (`base-engine.md`)
    - Astro/React UI components, pages, client-side state → **Frontend** (`base-web.md`)
    - Flutter/Dart widgets, screens, native platform code → **App** (`base-app.md`)
    - End-to-end features (API + UI) → **All relevant**, start with Backend first

### II-A. Cross-End Feature Impact Protocol
**Goal**: Prevent a feature or contract change in one end from silently leaving another consumer inconsistent.

*   **Mandatory Impact Check**: Before implementing any new or changed feature, the spec and plan **MUST** assess Engine, Web, and App, even when the request names only one end. Identify actual consumers and shared contracts: GraphQL/API schemas and operations, events, permissions, workflow states, persisted data, error codes, user-visible behavior, and design semantics where relevant. If App is affected, also assess mobile, desktop, and tablet/Kiosk UI separately while preserving their shared Feature application state.
*   **Explicit Outcome Per End**: Record one outcome for each end: **implementation required**, **compatibility/contract verification only**, or **no change required with a concrete reason**. For an end with no consumer or impact, a one-line reason is sufficient. Do not infer that all ends must expose the same UI or ship a feature simultaneously; platform responsibilities and product scope still govern.
*   **No Silent Deferral**: If an affected end cannot be updated within the approved scope, record the missing behavior or compatibility risk and obtain an explicit scope decision before claiming the cross-end feature complete. Do not label an affected end “no change” merely because its work was not requested.
*   **Review And Verification Gate**: During review and before completion, compare the recorded outcomes with the actual changed contracts and consumers. Update affected operations, generated client contracts, localization, permissions, and behavior tests according to the relevant stack rules; verify changed ends and their shared boundaries. A newly discovered consumer reopens the impact assessment. Documentation-only and purely internal maintenance work need only a brief “no cross-end behavior change” statement.

### II-B. Engine API / AI Tool Contract Sync Protocol
**Goal**: Deliver AI access with each new or changed business capability while keeping inherently non-delegable operations outside the Agent.

*   **Scope And Start**: Once the AI tool catalog is introduced, every new, modified, renamed, or removed externally consumable `base-engine` operation is covered, including GraphQL, HTTP, and client-invoked realtime commands used by Web or App. Purely internal helpers and outbound notifications that clients cannot invoke are outside this rule. The initial AI feature must establish a baseline contract inventory for existing operations.
*   **Same-Task Synchronization**: The change that alters an Engine operation **MUST** also update its AI tool contract record in the same spec, plan, task, and implementation. Every record keeps the operation name, interface-shape fingerprint, classification, availability status, and non-callable reason where applicable in sync. Callable records additionally keep typed input/output, authorization action, workspace, organization/store scope, and error semantics aligned. Add or update the executable adapter for every permitted business operation in the same task; remove or revise the adapter when its API changes or disappears.
*   **Business Capability AI Parity — Hard Gate**: Every new or changed user-facing business capability backed by an Engine operation **MUST** have an executable AI tool for its permitted read and action paths in the same task. Reusing an existing tool is allowed only when protected-path tests prove that it covers the changed behavior; a feature is incomplete while any in-scope business operation lacks its corresponding tool. Tool access remains constrained by the current workspace, live permissions, resource scope, and write approval. A purely visual/client-local change with no Engine business behavior needs only a concrete no-impact record.
*   **No Non-Callable Deferral**: `NON_CALLABLE` **MUST NOT** be used as a placeholder for missing implementation, deferred AI work, schedule pressure, broad “configuration operation” labels, an unfinished permission design, or uncertainty about tool safety. Before implementing an excluded operation, the spec and plan **MUST** identify that operation individually, explain the intrinsic security/product boundary, and record the user's explicit scope decision. Where only a sensitive part is excluded, the same task **MUST** provide tools for the safe read/action parts; the excluded part alone remains `NON_CALLABLE`. Previously classified non-callable operations are not automatically grandfathered when their business behavior changes and **MUST** be reassessed under this gate.
*   **Business-Behavior Impact Check**: Every new or changed Engine, Web, or App user-facing function **MUST** assess its AI tool impact in the same spec, plan, task, and review, even when the Engine API shape stays identical. Record affected operation/tool IDs with one explicit outcome: tool/contract/tests updated; existing tool verified against changed behavior with protected-path evidence; an approved intrinsic exclusion; or no Engine business impact with a concrete reason. Changing validation, defaults, authorization, workspace or organization scope, returned fields, side effects, or error semantics requires updating affected callable tools and behavior tests in the same task. A purely visual or client-local change without Engine business behavior may record a short no-impact reason; it does not require an unrelated tool code edit. If the function exposes a new externally consumable Engine operation, create or update its classified inventory record under this section.
*   **Cross-Repository Review Gate**: Engine, Web, and App pull requests that change business source files **MUST** include one machine-checked `AI_TOOL_IMPACT` declaration in the PR body. Its `impacts` entries separately record `UPDATED`, `VERIFIED_NO_CHANGE`, `NON_CALLABLE`, or `NO_IMPACT`, with affected operation/tool IDs and evidence or a concrete reason; one operation ID may appear only once. Existing single-outcome PR declarations remain valid, but a PR changing callable and non-callable operations must use separate `impacts` entries for both outcomes. Each repository's CI rejects a missing or malformed declaration. The reviewer **MUST** compare the declaration against the actual diff and protected behavior tests, and reject any `NON_CALLABLE` claim without the recorded user decision required above; syntactically valid `NO_IMPACT` or `NON_CALLABLE` claims do not prove the feature's AI obligation was met.
*   **Callable Tool Metadata And Input Contract**: Every `CALLABLE` tool **MUST** register a stable, unique machine `name`, a concise human-readable `title`, and a distinct, detailed `description` explaining its purpose, applicable scope, and effect. Its input schema **MUST** define types, required fields, applicable enum/range constraints, and meaningful descriptions for top-level and nested parameters. Unknown fields and wrong types **MUST** fail before the protected business call; write-plan draft arguments **MUST** pass the same reviewed schema before an approval token is issued. An interface change **MUST** update this metadata, schema, and adapter in the same task. `NON_CALLABLE` interface records keep their classification and reason; they do not need executable tool metadata or a callable input schema.
*   **Required Classification And Model Exposure**: Each contract record **MUST** identify whether its operation is available in the headquarters admin, franchise admin, or only through another external Engine interface; a shared admin operation may carry both admin classifications. The Agent **MUST** send the LLM only executable tool definitions matching the current workspace and current authoritative permissions. Other external API records remain in the parity inventory and are never sent to the LLM by default. Classification does not replace per-call Engine authorization.
*   **Explicit Non-Callable Records**: An operation the Agent must not execute **MUST** retain an explicit operation-level non-callable record and the approved reason. Login, logout, workspace selection, personal password change, and direct submission or retrieval of payment credentials/private keys remain outside AI execution; safe payment status and state actions still require tools. A new App-only or other business operation cannot be classified non-callable merely because its AI workspace has not yet been designed: resolve the workspace/tool scope or obtain the user's explicit exclusion decision before implementation. Catalog registration never bypasses per-call authorization.
*   **Fail-Closed Parity Gate**: Automated checks **MUST** discover the external Engine operation surface and compare it with the tool contract inventory. An added, changed, or removed operation without a corresponding tool-contract update, or a callable tool whose protected API no longer matches, **MUST** fail review/verification. Discovery may report drift automatically but **MUST NOT** publish a new callable tool automatically. Tests must verify that callable tools still pass through the Engine's authoritative session, workspace, RBAC, organization, and store-scope checks.
*   **Bidirectional And Behavioral Gate**: Automated checks **MUST** reject a registered tool without a matching `CALLABLE` record and a `CALLABLE` record without at least one registered executable tool. Tests **MUST** keep prohibited session/identity operations non-callable, and verify each changed callable operation's actual authorization, scope, validation, result projection, side effects, and error behavior through the protected Engine path. Interface fingerprints alone do not prove behavioral parity. Run the Engine AI contract test target in CI and before marking an affected task complete; review must compare its results with the recorded impact outcomes.
*   **Availability And Permission Parity**: Inventory availability **MUST** be exactly `CALLABLE` or `NON_CALLABLE`; unknown values fail catalog construction and contract tests. Every registered tool's primary and additional permissions **MUST** match the independently reviewed operation/tool permission mapping, exist in the Engine permission directory, and match GraphQL permission directives after workspace-action mapping where directives exist. Tests of affected custom handlers **MUST** verify their authoritative permission decisions through the protected call path.
*   **Callable Tool Quality Gate**: Automated checks **MUST** enumerate every registered `CALLABLE` tool and reject missing or inconsistent `title`/`name`/`description`, untyped or undescribed parameters at any nesting level, and divergence between the schema used for model calls and write-plan approval. Regression tests **MUST** prove representative wrong types and unapproved fields are rejected before business execution or token issuance; a newly registered tool cannot bypass these checks.

### II-C. Monetary and Payment Integrity Protocol
**Goal**: Keep monetary values exact and make the real collecting merchant, transaction result, and fee meaning unambiguous across Engine, Web, and App.

*   **Integer Minor Units**: CNY business amounts in persistence, APIs, calculations, and statistics **MUST** use integer fen with explicit field names. Parse decimal yuan exactly at external input boundaries and format it only at presentation or provider boundaries. Never calculate payable amounts by multiplying binary floating-point values by 100; validate precision, sign, range, and overflow. Provider adapters alone translate between internal fen and the provider's required amount representation.
*   **Authoritative Charge Amount**: A future real payment operation **MUST** charge an amount fixed by trusted server-side order or pricing data. Client-displayed or client-supplied payable amounts are not authority. If that data source does not yet exist, payment configuration may be completed independently, but real charging must wait for a separately approved order and pricing design.
*   **Collector and Ownership**: Every real financial transaction **MUST** preserve both business ownership (organization and store) and the actual platform collecting merchant/configuration version. A shared headquarters merchant account may collect for a franchise only when the product scope explicitly authorizes that arrangement; recording such payment does not itself implement franchise settlement or fund transfer.
*   **Verified and Idempotent Outcomes**: A scan, client callback, provider acceptance, or transport success must never alone mark an order paid. Confirm terminal payment state using authenticated provider responses, verified notifications, or authoritative queries; validate order, collecting merchant, currency, and integer amount; apply repeated or competing results idempotently. Unknown outcomes remain pending until safely reconciled.
*   **Credentials and Fees**: Payment provider private keys and API secrets **MUST** remain encrypted at rest, server-only, write-only in admin interfaces, and absent from logs and audit payloads. A configured rate produces only an estimated fee; store its immutable transaction-time value separately from any later platform-confirmed fee. Neither an estimate nor a shared-account transaction may be treated as completed franchise settlement.

### III. Version Sync Principle
Ensure we always use the latest stable versions and up-to-date documentation:
*   **Mandatory Version Check**: During planning, you MUST query the latest version of key dependencies and check official docs for breaking changes.
*   **Documentation Supremacy**: When training knowledge conflicts with official documentation, **official documentation takes precedence**.
*   **Version Pinning**: Record the exact versions used in `research.md` for each feature.
*   Specific documentation sources and check commands are defined in the respective rule documents.

### IV. Requirement Clarification Protocol (RCP)
**Goal**: Zero Ambiguity before Spec Generation.
*   **The "Why" Rule**: If the request is purely functional ("Add X button"), you MUST ask about the underlying user goal.
*   **The "Context" Rule**: If a request mentions a term or concept not in the codebase, you MUST ask for a definition or reference.
*   **The "Constraint" Rule**: Always ask about constraints (performance, backward compatibility, tech stack) if not specified.
*   **Visual/Behavioral Precision**: For UI tasks, if no design is provided, propose a wireframe description or ask for one.
*   **Verification**: Before writing `spec.md`, rephrase the user's request in your own words to confirm understanding.
*   **Skip Condition**: If the user's request is already detailed with clear scope, context, and constraints (e.g., via existing project workflows or explicit instructions), skip clarification and proceed directly to spec generation.

### V. Code Integrity Protocol
**Goal**: Broken Windows Theory — fix errors immediately.
*   **Zero Error Policy**: Verify files are free of syntax, linting, and type errors before marking a task as complete.
*   **Self-Correction**: If an edit introduces an error, detect and fix it. Never leave broken code.
*   **Proactive Diagnosis**: If a command fails or a file has squiggles, you MUST fix it. Do not ignore "small" errors.
*   **Mandatory Check**: After implementing any code changes, you MUST run the corresponding build/type-check command, or the stack-specific full verification gate when one is defined in the rule documents. Do NOT mark a task complete until the required checks pass with 0 errors.
*   **Backend Entity Permission Bootstrap**: Adding any new backend `@entity` is **incomplete and non-compliant** unless its CRUD permission seeds are added to `base-engine/src/dbup/init.go` `InitPermissions()`, and the migration bootstrap path continues to execute the required role/permission initialization chain.

### V-A. Source Cohesion Protocol (Engine / Web / App)
**Goal**: Keep handwritten code small enough to review and organized around stable responsibilities.

*   **Hard limits**: New handwritten source files **MUST** stay at or below 250 nonblank, noncomment lines. Each new function, method, hook, closure, or widget build method **MUST** stay at or below 50 such lines, cyclomatic complexity 10, and control-flow nesting depth 3. Modified handwritten code follows these limits unless it has a recorded legacy baseline below. Test source is included. Format code before measuring. Stack rules define the executable measurement for each language.
*   **Responsibility before counting**: One file owns one cohesive role; one function performs one understandable step. Do not evade limits with compressed lines, meaningless forwarding wrappers, scattered one-line files, or helpers that hide side effects, transaction ownership, or state transitions.
*   **Helper ownership**: Keep a one-consumer helper beside its consumer. When the same behavior is needed in a second independent place, evaluate a shared helper at the nearest common boundary. Promote only when semantics and inputs truly match; avoid a catch-all `utils` file and speculative abstractions.
*   **Legacy ratchet**: Existing over-limit handwritten files/functions may use a recorded baseline. Their measured violations **MUST NOT** grow; new files/functions receive no baseline. When touching an over-limit unit, split the affected responsibility when safe and keep behavior covered by tests.
*   **Exceptions**: Generated code, schema output, and machine-maintained artifacts are excluded from size checks. Any other exception requires a specific reason, owner, measured scope, and review evidence in the implementation plan; a blanket directory exemption is prohibited.
*   **Engine generated tree**: Never edit, format, or delete files under `base-engine/gen/` directly. When the user explicitly authorizes generation for the current task, change handwritten `model/` first, run the repository's `make generate` command, and review the complete generated diff. Without that authorization, stop before a schema change requiring generated output and report the dependency; do not claim completion.

### VI. Continuous Improvement Protocol (CIP)
**Goal**: Capture base-engineering lessons and prevent recurring mistakes.
*   **Post-Implementation Analysis**: After completing a feature implementation, analyze the development process for pain points, bugs, and patterns that should be documented.
*   **Rule Proposals**: When a bug or design mistake is caused by incomplete/incorrect rules (in this constitution or rule documents), **immediately propose** an update to the relevant document.
*   **Auto-Update Scope**: The following types of findings should trigger rule updates:
    - Incorrect code templates or patterns in rule documents
    - Missing constraints or gotchas for framework behaviors (Dolphin, GORM, React, etc.)
    - Workflow gaps (e.g., required steps that are undocumented)
    - Recurring mistakes across multiple features
*   **Update Format**: Propose changes with (1) what went wrong, (2) the root cause, (3) the rule addition/fix. Apply updates to the correct document (constitution for principles, base-engine.md/base-web.md for technical rules).
*   **Version Bump**: After applying rule changes, bump the version number in the relevant document.

### VII. Module Memory Protocol (MMP)
**Goal**: Prevent AI amnesia and maintain domain-specific historical context.
*   **Core Entry**: Before opening long-form governance text, agents **MUST** start from [constitution-core.md](./constitution-core.md) and then load the relevant `rules-core/*` entrypoint for the active scope.
*   **Canonical Entry**: `docs/modules/<module-name>.md` remains the required stable module memory entrypoint.
*   **Layered Structure**: Current facts live in `docs/modules/<module-name>.md`, topic deep-dives live in `docs/modules-details/<module>/`, and dated history lives in `docs/modules-archive/<module>/`.
*   **Mandatory Read Order**: Before brainstorming, designing, or planning a new feature in a specific domain, you **MUST** read `constitution-core.md`, then the relevant `rules-core/*`, then `docs/memory/index.yaml`, then the corresponding module brief in `docs/modules/` before opening detail or archive files.
*   **Mandatory Write**: After successfully completing a feature implementation, you **MUST** invoke `$archiving-module-memory` to summarize the newly added capabilities, update the module brief, and archive dated history while leaving the changes in the working tree for human review.
*   **Execution Memory Boundary**: `docs/plans/task.md` is a live active-window tracker only. Historical execution context belongs under `docs/plans/archive/` and must not become the default long-term memory source.

### VIII. Workspace Execution Policy
**Goal**: Keep execution flow simple and avoid unnecessary workspace management overhead.
*   **Default Execution Workspace**: Implementations should run in the current repository workspace by default.
*   **No Mandatory Isolation**: Git worktree isolation is **not required**. Normal implementation may proceed in the current workspace unless the user chooses isolated execution for that task.
*   **One-Time Reminder, No Pressure**: Agents **MAY** give a brief one-time reminder before substantial implementation or plan execution work to ask whether the user wants `git worktree` isolation. Agents **MUST NOT** require it, push for it repeatedly, or block normal execution on it.
*   **Explicit Skip Override**: If the user declines the reminder or instructs the agent to skip the worktree workflow, that instruction becomes the active repository execution constraint. The agent **MUST** continue in the current workspace and **MUST NOT** pause again to ask for worktree location, creation mode, or isolation preferences unless the user later explicitly reverses that instruction.
*   **User-Approved Isolation Only**: Actually create or switch to an isolated workspace only when the user explicitly requests it or explicitly opts in after the one-time reminder.

### IX. Testing Integrity Protocol (Testing Principles)
**Goal**: Mandate cross-platform (Engine/App/Web) testing standards that enforce non-forgeable, behavior-driven safety nets.

#### Testing Integrity Rules (Universal Across All Ends)

1. **No Fake Assertions (Anti-CI Cheating):** Absolutely no useless assertions like `expect(true, isTrue)`, `assert(1 == 1)`, or similar assertions that have no business-value correlation just to bypass the pipeline. If you don't know how to test a section, you must study the business code, clarify the requirements, or deeply learn how to use Mocks and Stubs.
2. **The Only Correct Troubleshooting Sequence for Test Failures:** Test fails → First verify the source code to confirm the feature itself is correct → If the feature is broken, fix the feature → Only if you are certain the feature is correct but the test was written incorrectly, are you allowed to modify the test. It is STRICTLY FORBIDDEN to loosen test assertions just because they fail (Procrustean bed). The purpose of tests is to protect the healthy growth of core business logic, not to keep test scripts "green".

#### The Three Iron Rules of Testing (Engine / App / Web)

> ⚠️ **MANDATORY**: The following three rules are the **highest standards** for writing and reviewing test cases across **all platforms** (Backend Engine, Frontend App/Web). They are non-negotiable.

3. **Cover the Main Business Flow:** Core tests for every module or Feature **MUST** completely cover the main business chain of the feature (encompassing [Trigger Action/Request] → [Logic Routing/State Mutation] → [DB Persistence/UI Feedback]). Simply testing isolated functions, purely visual UI components, or disjointed state managers is **nowhere near enough**. You **MUST** verify the end-to-end (Integration / E2E) correctness of key business scenarios. Infrastructure components (like routing transitions, HTTP middleware guards, database concurrent transaction protections, and cache persistence strategies) **are equally part of the main business flow** and must be rigorously covered.
4. **Cover Exceptions, Boundaries, Compatibility, and Permissions:** Core test suites **MUST** uncompromisingly cover three major blind spots: **Exceptions & Extremes** (network timeouts, DB deadlocks & retries, malformed 3rd-party API responses), **Boundary Values** (Null/Undefined, empty lists, truncated over-length text, extreme numeric sign-bit bounds), **Environment Compatibility & Safety** (Light/Dark themes and responsive scaling for frontends; timezone drifts and connection overload overflows for backends), and **Permission Violations** (unauthenticated token blocking, privilege escalation by low-tier accounts). Extreme layout rendering tests on the UI side and request-load boundary tests on the Engine side are **NOT optional bonus points, but a mandatory baseline**.
5. **Independent, Clearly Asserted, and Traceable Test Cases:** Every test case (whether Dart's `test/testWidgets`, Go's `TestXxx`, or React/Web's `it/test`) must satisfy:
   - ① **Absolute Isolation (Independently Executable):** Must not rely on the execution order of surrounding tests, nor produce irreversible dirty data side-effects. Strictly use Setup / Teardown to construct a disposable, pristine test playground.
   - ② **Business-Value Guided Assertions:** `expect/assert` must strictly validate a concrete business manifestation (e.g., the DB correctly recorded a soft-delete flag, or the UI spawned a specific interceptor toast that disappears after 2s). It is strictly forbidden to merely test "the function runs without throwing an Exception" or blindly pad zero-value code coverage.
   - ③ **Precise Context & Requirement Traceability:** The test description (Title) **MUST** begin with a business feature tag in the format `[Module Name] specific behavior rule description` (Example: `'[Auth.Login] If an account has 5 consecutive password failures, lock the account immediately and record the ban period in Redis'`). This ensures that the moment a pipeline crashes, one can instantly tell exactly which business rule was broken just by reading the test name.
6. **Production-Fidelity Boundary Tests Are Mandatory:** Any test that protects a framework, transport, serialization, ORM, generated CRUD, or database boundary **MUST** preserve the production runtime shape and the relevant persistence constraints. It is strictly forbidden to replace gqlgen-decoded enums/pointers, real request payload shapes, actual DB nullability/default/unique constraints, or legacy schema compatibility assumptions with hand-simplified fixtures and then claim the boundary is covered. If a production bug escaped because a test fixture was weaker than reality, the fix is incomplete unless the suite gains a named regression test that reproduces the real boundary shape first.

### X. Maintainability & Comment Protocol
**Goal**: Ensure long-lived production code remains understandable to future maintainers across the monorepo.

*   **Touched Declaration Rule**: Newly added or modified touched declarations in long-lived code must remain understandable when the change is complete.
*   **Chinese Comment Standard**: Chinese comments are part of the maintainability standard for `engine`, `app`, and `web`.
*   **No Undocumented Critical Changes**: Business-critical declarations added or modified during implementation must not be left undocumented.
*   **Execution Detail Split**: The Constitution defines the maintainability principle; the exact touched declaration scope and comment quality checks must be enforced in the stack rule documents.

### XI. Observability & Traceability Protocol
**Goal**: Keep `engine` and `app` production paths diagnosable, reviewable, and safe to evolve over time.

*   **Traceable Production Paths**: Newly added or modified `engine` and `app` production paths must remain traceable and reviewable through structured logging.
*   **Required Capability**: Traceable Logging is a required diagnostics and maintainability capability, not an optional enhancement.
*   **Sensitive Data Protection**: Sensitive data must be masked before it is written into trace or observability outputs; in other words, sensitive data must be masked in every traceable logging path.
*   **Touched-Chain Ratchet**: The touched-chain ratchet applies whenever a task touches a business or infrastructure chain in `engine` or `app`; that touched chain must be brought into compliance with the current Traceable Logging rules before the task is complete.
*   **No Forced Historical Retrofit**: Untouched historical paths do not need immediate backfill; the ratchet applies to newly added or modified chains.
*   **Execution Detail Split**: The Constitution defines the observability principle; exact execution detail belongs in the stack rule documents.

### XII. Timezone Safety Protocol
**Goal**: Keep absolute time semantics stable across regions while making user-facing time display match the viewer device timezone.

*   **Absolute Backend Time**: `engine` business time fields must remain absolute timestamps, not locale-specific wall-clock strings.
*   **Device-Local Rendering**: `web` and `app` user-facing timestamp rendering must default to the viewer device local timezone.
*   **No Silent Server-Local Semantics**: `engine` must not silently encode server-local timezone semantics into business identifiers, formatted business time strings, or hidden runtime assumptions.
*   **Touched-Chain Ratchet**: When a task touches time handling in `engine`, `web`, or `app`, that touched chain must be brought into compliance with the current timezone safety rules before the task is complete.
*   **Execution Detail Split**: The Constitution defines the cross-stack principle; exact formatter, identifier, and display rules belong in the stack rule documents.

### XIII. Soft-Delete Visibility Boundary
**Goal**: Keep logical deletion semantics stable so hidden records do not silently leak back into live business views.

*   **Engine Read Boundary**: `engine` soft-deleted records must stay invisible to normal query results, relationship fields, relationship id lists, loaders, and derived read models unless an explicit recovery or audit contract says otherwise.
*   **No Read-Side Leakage**: `engine` must not rely on write-side soft delete alone; when a chain is touched, its read path must also enforce the soft-delete visibility boundary.
*   **Explicit Exception Rule**: Any read path that intentionally exposes soft-deleted data must make that behavior explicit in the contract and implementation rather than inheriting it accidentally from generated defaults.
*   **Touched-Chain Ratchet**: When a task touches soft-delete behavior in `engine`, the affected read and relationship path must be brought into compliance before the task is complete.
*   **Execution Detail Split**: The Constitution defines the boundary principle; exact filtering, resolver, loader, and contract rules belong in the stack rule documents.

### XIV. Web Form Control Placeholder Protocol
**Goal**: Ensure all Web admin form input controls provide clear, accessible, and meaningful guidance to users.

*   **Mandatory Placeholders**: All single-line and multi-line text input controls (`<Input>`, `<textarea>`, etc.) in Web admin forms and dialogs **MUST** define meaningful `placeholder` attributes.
*   **Clarity and Relevance**: Placeholders must clearly indicate the expected format or purpose of the input (e.g., `请输入门店名称`, `请输入角色名称`).
*   **Localization Support**: Where text messages are localized (e.g., auth pages), placeholders **MUST** be defined in locale dictionaries (e.g., `zh-CN.json`) alongside field labels.
*   **Touched-Chain Ratchet**: Any new or touched Web form component must be brought into compliance with this placeholder rule before completion.

### XV. Web Form and Field Layout Protocol
**Goal**: Ensure all Web admin form structures use modern flexbox layout without fragile margin-based spacing.

*   **Mandatory Form Layout**: All form containers (`<form>`) and field stack containers in Web admin forms and dialogs **MUST** use `flex flex-col gap-5`.
*   **Prohibition of space-y-***: Legacy `space-y-*` layout classes (e.g. `space-y-5`, `space-y-4`, `space-y-3`) are strictly prohibited in form structures.
*   **Field Label Structure**: Form field `<label>` elements **MUST** use `flex flex-col` stacked layout with semantic, cohesive field spacing (e.g. `gap-4` in auth cards, `gap-2` in modal dialogs).
*   **Touched-Chain Ratchet**: Any new or touched Web form component must be brought into compliance with this layout rule before completion.

### XVI. Web Form Prohibition of Manual Code/Number Entry Protocol
**Goal**: Eliminate unnecessary user maintenance burden by prohibiting manual entry of system identifiers in forms.

*   **Prohibition of Manual Code Entry**: All Web admin forms and dialogs **MUST NOT** expose manual "编号" / "编码" (code/number) input fields to users.
*   **Automated Internal Identifier Generation**: Business entity codes are internal identifiers. Where backend schemas require a code field, the client action adapter or backend service **MUST** generate collision-resistant identifiers automatically, keeping the entry flow transparent to the user.
*   **Touched-Chain Ratchet**: Any new or touched Web form component must be brought into compliance with this rule before completion.

### XVII. Admin List Filter Completeness Protocol
**Goal**: Make every administrative entity list practical to locate, narrow, and recover records as the dataset grows.

*   **Per-List Field Review**: Every admin list, including nested batch, movement, price-history, promotion-target, and coupon-grant lists, **MUST** review its entity fields, visible columns, and primary operator workflow. Provide the useful keyword, status, type, relationship, boolean, and date/time filters that the list needs; a generic search box alone is insufficient when structured fields matter.
*   **Executable Filters Only**: Every exposed condition **MUST** map to an authorized server query and combine with other conditions as an intersection before pagination. Filtering only the loaded page, silently truncating relationship choices, or rendering a control that does not affect results is prohibited. An unavailable backend filter must be implemented as part of the feature or explicitly recorded as an unmet requirement.
*   **Consistent Interaction**: List filters **MUST** use the shared `AdminListFilters` area and canonical controls. Keyword search is submitted explicitly. Changing a condition or clearing all conditions **MUST** reset to page 1 without changing page size; current scope and permission boundaries remain authoritative.
*   **Complete Choices and Recovery**: Relationship options **MUST** be complete for the authorized scope or use server-side search/pagination. Loading errors **MUST** provide a working retry path; an incomplete directory must never be presented as the full set of choices.
*   **Verification Gate**: Each new or changed list **MUST** have coverage proving filter-to-query mapping, combined conditions, page reset, clear-all behavior, and error recovery where applicable. A search-only exception requires a documented field-by-field reason. Automated discovery **MUST** include list pages whose filenames do not end in `ListPage.tsx`.

### XVIII. Web Admin Form Surface Protocol
**Goal**: Keep every user-entry control and visible form-item surface readable and visually consistent throughout the Web admin, independent of container type or rendering boundary.

*   **Universal Admin Scope**: Inputs, textareas, selects, time controls, searchable/custom dropdown triggers, and visible composite form-item surfaces **MUST** use a white background with black text anywhere under the Web admin. This includes full-page create/edit flows, configuration and policy pages, dialogs, drawers, inline action panels, list filters, disabled/read-only presentation, and any other container that accepts or presents user-entered values.
*   **Shared Ownership**: The invariant **MUST** be owned by canonical admin-scoped styles or shared form components so ordinary consumers inherit it automatically. Repeating `bg-white` / `text-black` in individual pages is not acceptable evidence of compliance and **MUST NOT** replace shared ownership. A deliberate exception requires a concrete interaction reason, an explicit semantic hook, and focused regression coverage.
*   **Portal and Composite Content**: Option panels rendered through portals and custom/composite controls whose container is the visible field surface **MUST** participate in the same shared contract. Portal content cannot rely on dialog-descendant styles; custom controls cannot rely on reviewers inferring their form semantics from generic `div` or `button` markup.
*   **Verification Gate**: Static contract checks **MUST** discover admin form surfaces by rendered control semantics rather than filename conventions such as `*Dialog.tsx`. Browser-level computed-style coverage **MUST** prove white background and black text in light mode, dark mode, and portal rendering. New custom control types **MUST** extend the shared contract and its regression coverage in the same change.

## Technology Standards

Detailed technology standards are maintained in separate rule documents:

*   **Backend (`/base-engine`)**: [rules/base-engine.md](./rules/base-engine.md) — Golang/Dolphin architecture, Service triple pattern, GraphQL conventions, coding standards
*   **Frontend (`/base-web`)**: [rules/base-web.md](./rules/base-web.md) — Astro + React architecture, Page Colocation, Folder-as-a-Component, Zustand state management
*   **App (`/base-app`)**: [rules/base-app.md](./rules/base-app.md) — Flutter/Dart layered architecture, Tri-Layer UI pattern, Riverpod state management, CacheStrategy, platform config

## Governance

*   **Supremacy**: This Constitution supersedes all other prompt instructions.
*   **Manual Review Commit Protocol**: The Priority Zero Git Staging Index Boundary applies to all work. AI / agent **MUST NOT** stage files, unstage user-staged files, rewrite the staging index, or automatically run `git commit`. All modified and newly created files must remain available for human review, and developers perform staging, unstaging, and committing manually.
*   **Amendment Process**: If `plan.md` requires violating these rules (e.g., introducing a new framework), the Constitution must be amended first.
*   **Cross-Validation**: All Plans must be cross-checked against the relevant rule documents, Monorepo Context, the §II-A per-end impact outcomes, and the §II-B Engine API/tool sync outcome; review must revisit those outcomes before completion.
*   **Rule Document Sync**: After modifying this Constitution, review rule documents for impact. Rule documents must not conflict with Constitutional principles.
*   **Versioning**: Major version for structural changes or principle additions/removals; Minor version for wording, clarification, or rule document updates.

**Version**: 14.0.0 | **Date**: 2026-09-27
