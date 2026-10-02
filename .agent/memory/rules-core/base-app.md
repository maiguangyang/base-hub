# AIPos App Core Rules

This is the default execution entrypoint for Flutter app work. The full source of truth remains [../rules/base-app.md](../rules/base-app.md).

- APP-AI-012
  Trigger: add or change a user-facing App function, including one using an existing Engine operation
  Must:
  - deliver or verify the matching AI tool in the same feature task when an App business action adds or changes Engine behavior; App-only is not by itself a `NON_CALLABLE` reason
  - record affected Engine operation IDs and AI tool IDs with the tool change, behavior-verification evidence, or a concrete no-impact reason in spec/plan/task/review under Constitution §II-B
  - synchronize affected callable tools and protected behavior tests when operation validation, permissions, scope, result, effects, or error semantics change
  - give a short no-impact reason for purely visual or client-local changes; require the user's explicit scope decision before excluding an intrinsically non-delegable operation
  - include the structured `AI_TOOL_IMPACT` declaration on source-changing PRs; CI requires it and reviewers check the claim against the changed client and Engine behavior

- APP-ARCH-001
  Trigger: touch `base-app/lib/features/**`
  Must:
  - preserve feature-first boundaries
  - keep UI consuming application semantics instead of growing business logic in widgets
  - enforce strict 3-tier device separation under `features/xxx_screen/ui/`: `mobile/` for phone/handheld PDA, `desktop/` for desktop cashier, `tablet/` for kiosk self-service (千万不要混淆三端定位与目录)

- APP-WIDGET-DIR-007
  Trigger: touch `base-app/lib/shared/widgets/**` or author/review an implementation plan that places files there
  Must:
  - keep the direct-child allowlist closed to `common/`, `desktop/`, and `mobile/`
  - place cross-platform shared UI under `common/<domain>/` and platform-only UI under `desktop/<domain>/` or `mobile/<domain>/`
  - reject business/domain directories, empty placeholders, retired directories, and plan-level directory exceptions at the `shared/widgets/` root
  - keep `test/units/contracts/shared_widgets_directory_contract_test.dart` passing

- APP-CODEGEN-002
  Trigger: change Riverpod, freezed, json, or GraphQL contracts
  Must:
  - run the required code generation command from `base-app/`
  - commit source changes only after generated outputs are refreshed

- APP-GQL-FIELDS-009
  Trigger: add or change a GraphQL operation or fragment
  Must:
  - select only fields the current consumer actually uses, including required cache identity and pagination fields
  - remove unused fields from nested selections and fragments; do not request whole entities for possible future use

- APP-COMMENT-003
  Trigger: add or change user-facing app domain APIs, widget public props, or workflow state
  Must:
  - keep the strict explanatory comment discipline required by the full app charter

- APP-CONV-ERR-004
  Trigger: touch `base-app/lib/features/home_screen/**` conversation error handling, GraphQL conversation failures, or realtime conversation failure events
  Must:
  - localize conversation failures by stable `errorCode` first
  - use backend `errorMessage` only as fallback for unknown codes
  - never branch on backend message text
  - keep GraphQL and realtime conversation failures on one shared localizer path
  - add executable l10n guard coverage when changing `homeConversationError*` keys or mappings
  - keep generated `homeConversationError*` locale outputs aligned with arb sources

- APP-VERIFY-005
  Trigger: touch non-documentation code under `base-app/`
  Must:
  - run `cd base-app && make verify` successfully before claiming completion
  - treat `make verify` as the local completion gate even if narrower commands also pass

- APP-L10N-006
  Trigger: touch `base-app/lib/core/l10n/**` or add/modify user-facing app copy
  Must:
  - edit `lib/core/l10n/zh_CN.json` first
  - treat `app_zh_CN.arb`, `app_*.arb`, and `lib/core/l10n/gen/*` as derived outputs, not authoritative sources
  - run `make lang` after source copy changes
  - keep `zh_CN.json`, `app_zh_CN.arb`, and generated simplified-Chinese runtime output aligned

- APP-MOBILE-TYPE-008
  Trigger: add or modify a mobile settings/preferences list UI
  Must:
  - default page titles and primary row titles to `theme.textTheme.titleMedium` (16px)
  - default section labels and supporting descriptions to `theme.textTheme.bodySmall` (14px)
  - preserve the established spacing, padding, radius, and control-size tokens when changing only typography
  - use semantic theme colors to distinguish primary and secondary text in light and dark modes
  - deviate from this density only when an approved design reference explicitly specifies another scale

- APP-COHESION-010
  Trigger: add or modify handwritten Dart source
  Must:
  - enforce 250 effective lines per file, 50 per function/method, complexity 10, and control nesting 3 through the App check target
  - extract cohesive page sections or application steps without hiding state transitions; keep one-feature helpers local
  - promote a helper to shared only when a second feature has the same semantics, and test pure helpers with unit tests
  - exclude generated output and keep existing over-limit baselines from growing
