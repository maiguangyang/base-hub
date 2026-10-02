# AIPos Engine Core Rules

This is the default execution entrypoint for backend work. The full source of truth remains [../rules/base-engine.md](../rules/base-engine.md).

- ENG-REL-000
  Trigger: add or modify an `@entity`, model an association between entities, or touch a handwritten entity field whose value identifies another project entity/table
  Must:
  - model every entity-to-entity association with typed object/list fields using `@relationship` and a correct inverse
  - never declare or use a custom scalar `xxxId` / `xxxIds`, `<entity>Id` / `<entity>Ids`, or renamed scalar alias to associate, join, query, or enforce referential linkage to another entity/table
  - apply the same rule to runtime, history, audit, ledger, session, governance, and explicit association entities; an association entity's endpoint links must still be `@relationship`
  - treat generator-owned FK columns, generated relationship ID projections, and generated mutation inputs as implementation artifacts, never as precedent for handwritten schema fields
  - allow IDs in custom operation inputs/payloads only as request/response references; never persist or reuse them as substitutes for entity relationships
  - stop planning/review/implementation on a violation and require a schema contract test that fails closed on custom scalar cross-entity links

- ENG-GEN-000
  Trigger: any work in `base-engine`
  Must:
  - never edit, format, or delete `gen/**` directly
  - run `make generate` only after handwritten `model/` changes and explicit user authorization for the current task; inspect the complete generated diff
  - without that authorization, stop before a schema change requiring generated output and report the dependency instead of claiming completion

- ENG-SCHEMA-001
  Trigger: touch `base-engine/model/*.graphql`
  Must:
  - follow schema-first
  - obey ENG-GEN-000; use `make generate` only with task-scoped user authorization, never edit generated output directly
  - when a field carries business uniqueness truth, prefer `@validator(unique: "true")` in schema, and use `uniqueScope` when the uniqueness is intentionally scoped
  - do not rely on UI-only or ad hoc service-layer duplicate checks as the sole source of truth when Dolphin's schema validator can express the uniqueness contract directly
  - do not use `@validator(unique: "true")` for runtime/history/computed fields or conditional uniqueness that depends on join tables or richer delete/governance semantics than the validator can express
  - when adding or changing `@validator(unique: "true")`, add or update regression tests that cover both the schema contract and the production-fidelity write-path uniqueness behavior
  - when a delete-protection rule is a simple FK-backed parent->children constraint, prefer `@relationship(..., master:"yes")` on the parent-side relationship before adding a hand-written delete guard
  - do not use `master:"yes"` on many-to-many relationships, join-table-backed relationships, cascade-cleanup flows, or governance/runtime-history entities that need specialized delete semantics
  - when adding or changing `master:"yes"`, add or update contract tests that cover both the schema snippet and the generated delete guard in `gen/resolver-mutations.go`
  - obey ENG-REL-000 before all other entity modeling choices

- ENG-PERM-002
  Trigger: add a new `@entity`
  Must:
  - add CRUD permission seeds
  - keep bootstrap permission registration complete

- ENG-SERVICE-003
  Trigger: add or refactor business logic under `base-engine/src/services/`
  Must:
  - preserve service-module boundaries
  - follow existing module patterns instead of flat-filing new service code

- ENG-TEST-004
  Trigger: touch GraphQL decode, generated CRUD, GORM persistence, migration compatibility, or legacy schema behavior
  Must:
  - add or update a production-fidelity regression test before the fix
  - use gqlgen/runtime-representative payload shapes instead of simplified stand-ins
  - preserve the touched DB nullability/default/unique/legacy constraints in test fixtures

- ENG-CONV-ERR-005
  Trigger: touch `base-engine/src/services/conversation/**`, GraphQL conversation user-visible errors, or realtime conversation failure payloads
  Must:
  - originate user-visible conversation failures from a centralized conversation error-code vocabulary
  - expose GraphQL conversation codes through `extensions.code`
  - include `errorCode` in realtime conversation failure payloads
  - treat `errorMessage` as fallback-only, never as the cross-end primary contract

- ENG-LANG-006
  Trigger: touch conversation orchestration, tool routing, prompt contracts, confirmation microcopy, or any engine path that interprets free-form user text
  Must:
  - not gate behavior on hardcoded Chinese/English phrases, locale-specific regex, or handwritten natural-language stopword lists
  - prefer schema-driven fields, runtime lexicon data, or validated structured model output over free-form phrase matching
  - require prompt-generated user-facing copy to match the user's language or locale instead of seeding one fixed server-authored language
  - avoid recovering cross-component behavior by matching localized prose when a machine-readable contract can exist

- ENG-PROMPT-007
  Trigger: add, move, rename, or load a System Prompt markdown asset in `base-engine`
  Must:
  - keep markdown assets that belong to the shared `system_prompt` package under `base-engine/system_prompt/`
  - not place shared System Prompt markdown files under `src/services/**`, `tmpdebug/`, or other feature-local directories
  - treat only non-System-Prompt local prompt assets as eligible for service-level colocation

- ENG-COHESION-008
  Trigger: add or modify handwritten Go source
  Must:
  - enforce 250 effective lines per file, 50 per function, complexity 10, and control nesting 3 through the Engine check target
  - split by responsibility and keep helpers package-local until shared semantics justify promotion
  - exclude `gen/**` from measurement and preserve non-growing baselines for existing violations

- ENG-AI-009
  Trigger: add or modify an externally consumable Engine operation, its business behavior, or an AI tool adapter
  Must:
  - deliver an executable AI tool for every permitted new/changed business read or action in the same task, or prove the existing tool covers it through protected tests; do not mark a business operation `NON_CALLABLE` to defer implementation
  - before implementation, record the user's explicit scope decision for each intrinsically prohibited operation; keep sensitive suboperations excluded while implementing tools for safe parts
  - record affected operation/tool IDs and a tool update, verified-no-change, or no-impact reason in spec/plan/task/review; API shape unchanged does not waive the check
  - keep the classified operation inventory synchronized under Constitution §II-B; never turn a prohibited record into a model tool automatically, and reassess an existing non-callable record when its business behavior changes
  - reject both callable inventory records without registered executable tools and tools without callable records; keep session/identity operations non-callable
  - restrict availability to `CALLABLE`/`NON_CALLABLE`; compare every tool permission to reviewed candidate coverage and annotated GraphQL workspace policy
  - supply a truthful `AI_TOOL_IMPACT` declaration for source-changing PRs; Engine CI checks its structure and reviewers verify the stated impact
  - test changed callable behavior through the protected Engine path, including authorization, scope, validation, result, side effects, and errors; run `make check_ai_contracts`
  - give every callable tool a unique stable `name`, human-readable `title`, distinct detailed `description`, and strict typed input schema with meaningful descriptions for nested parameters
  - reject wrong types and unknown fields before the protected call, and validate write-plan drafts against the same reviewed schema before issuing approval
  - enumerate all callable tools in metadata/schema contract tests and cover representative invalid inputs
