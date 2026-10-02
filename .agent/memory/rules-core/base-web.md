<!--
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-17 16:00:35
-->
# AIPos Web Core Rules

This is the default execution entrypoint for web admin work. The full source of truth remains [../rules/base-web.md](../rules/base-web.md).

- WEB-AI-012
  Trigger: add or change a user-facing Web function, including an admin action implemented through an existing Engine operation
  Must:
  - deliver or verify the matching AI tool in the same feature task when a Web business action adds or changes Engine behavior; `NON_CALLABLE` cannot stand in for unimplemented AI work
  - record affected Engine operation IDs and AI tool IDs with the tool change, behavior-verification evidence, or a concrete no-impact reason in spec/plan/task/review under Constitution §II-B
  - update affected callable tools and protected behavior tests in the same task when Web changes the operation's validation, permissions, scope, result, effects, or error semantics
  - record only a short no-impact reason for purely visual or client-local changes without Engine business behavior; do not change unrelated tools
  - include the structured `AI_TOOL_IMPACT` declaration on source-changing PRs; CI requires it and reviewers check the claim against the changed UI and API behavior

- WEB-ADMIN-001
  Trigger: touch `base-web/src/features/admin/**`
  Must:
  - preserve the admin list/page charter
  - keep shared list structure, summary, and action patterns aligned
  - discover every `*ListPage` through the list structure contract instead of maintaining a manual page allowlist
  - use `AdminFormDialogShell` for admin create/edit dialogs instead of hand-writing dialog structure

- WEB-GQL-002
  Trigger: add or change GraphQL operations
  Must:
  - colocate operations under feature GraphQL folders
  - use the generated `gql` wrapper instead of raw Apollo imports
  - select only fields the current consumer actually uses, including required cache identity and pagination fields; remove unused fields from operations and fragments

- WEB-TYPE-003
  Trigger: add or refactor frontend source files
  Must:
  - keep strict TypeScript boundaries
  - avoid `any` unless the full rule doc explicitly permits it with justification

- WEB-ADMIN-004
  Trigger: write or modify any page or page-private component under
  `base-web/src/features/admin/pages/**`
  Must:
  - read route context through `useAdminTab()` instead of `useParams`,
    `useLocation`, `useSearchParams`, `useMatch`, or `useMatches`
  - treat those hooks as shell-only APIs, permitted solely inside
    `features/admin/routes/**`, `components/AdminLayout/**`,
    `components/TabStack/**`, `components/AdminTabBar/**`, and `hooks/useAdminTab.ts`
  - back this rule with an ESLint `no-restricted-imports` check rather than
    relying on review

- WEB-COHESION-005
  Trigger: add or modify handwritten Astro, TypeScript, or TSX source
  Must:
  - enforce 250 effective lines per file, 50 per function, complexity 10, and control nesting 3 through ESLint and the source check
  - split page-local UI and hooks by responsibility; promote helpers only after a second semantic consumer or a stable cross-feature contract
  - exclude generated output and keep existing over-limit baselines from growing

- WEB-UI-006
  Trigger: touch or create form controls / dialogs in `base-web`
  Must:
  - ensure all text inputs and textareas define meaningful `placeholder` attributes
  - use concise, descriptive prompt text (e.g., `请输入门店名称`, `请输入角色名称`)
  - source placeholders from localization dictionaries where i18n is used (such as auth pages)

- WEB-UI-007
  Trigger: touch or create forms, dialogs, or form row layouts in `base-web`
  Must:
  - use `flex flex-col gap-5` for form containers and field row vertical layouts
  - avoid legacy `space-y-*` layout classes in form structures

- WEB-UI-008
  Trigger: touch or create forms, dialogs, or form inputs in `base-web`
  Must:
  - prohibit manual "编号" / "编码" (code/number) input fields in all forms
  - generate internal business identifiers automatically via `generateEntityCode` without burdening users

- WEB-UI-009
  Trigger: create or modify any Web admin user-entry control, form-item surface, list/filter control, full-page form, configuration or policy form, dialog, drawer, inline action panel, custom dropdown, or portal-rendered option panel
  Must:
  - use the canonical shared list and dialog components with their default presentation
  - when a list exposes search, use submit-based `AdminListSearch`; use semantic action variants and the fixed-header/scroll-body/fixed-footer `AdminFormDialogShell`
  - keep long-list bulk controls sticky within their own scroll container
  - never override the canonical `AdminToolbar` / `AdminTableShell` / `AdminPrimaryActionButton` appearance or hand-write a parallel form-dialog shell
  - keep every admin input, textarea, select, time control, searchable/custom dropdown trigger, visible composite field surface, and portal option panel white with black text in light and dark themes, including disabled and read-only presentation
  - make canonical admin-scoped styles or shared components own the invariant; do not treat repeated page-local `bg-white` / `text-black` utilities as compliance
  - require custom/composite controls to expose the shared form-surface semantic hook and extend automated coverage in the same change
  - enforce the contract by rendered control semantics rather than `*Dialog.tsx` or other filename conventions, and verify computed styles for ordinary, dark-theme, and portal rendering

- WEB-UI-010
  Trigger: create or modify any admin list or filter toolbar in `base-web`
  Must:
  - review each entity list, including nested lists and pages not named `*ListPage.tsx`, for useful keyword, status, type, relationship, boolean, and date/time conditions
  - use `AdminListFilters` and canonical controls; map every condition to the authorized server query before pagination, compose conditions with search, reset to page 1 on changes, and provide one clear-all action without changing page size
  - load relationship choices completely or use server search/pagination, provide a working retry for failures, and never filter only the loaded page or expose unsupported/decorative controls
  - test filter mapping, composition, reset, and recovery where applicable; document a field-by-field search-only exception

- WEB-UI-011
  Trigger: create or modify any interactive UI component or user-facing dialog in `base-web`
  Must:
  - consult the official shadcn/ui component catalog and use its component when it fits the required behavior
  - if that official component is absent locally, install it with the shadcn CLI into `base-web/src/components/ui` according to `components.json`
  - compose business-specific wrappers from shadcn components; hand-write a primitive only when the official catalog has no suitable component, and record that reason
  - use shadcn Dialog or AlertDialog for user-facing input and confirmation instead of native `window.prompt`, `window.confirm`, or `window.alert`

- WEB-UI-013
  Trigger: create or modify any Web select, dropdown, or filter control
  Must:
  - use `placeholder` only as the unselected control hint; never insert it into dropdown options, including disabled options
  - populate dropdowns from real business options; declare meaningful empty-value choices such as "顶级分类" or "不设置品牌" explicitly and independently with `clearLabel`
  - keep `AdminFilterSelect` options free of automatically generated all/empty items; clear list filters through `AdminListFilters` reset
  - cover shared selection behavior with tests and enforce the prohibition on rendering placeholder expressions as options through ESLint

- WEB-UI-014
  Trigger: create or modify any Web checkbox or checkbox group
  Must:
  - use the shared `@/components/ui/checkbox` component; never hand-write a native `input type="checkbox"` or another checkbox primitive in business pages
  - keep unchecked checkboxes solid white in every theme, including hover/focus; do not inherit page background tokens or override the unchecked background locally
  - preserve checked, indeterminate, disabled, and keyboard behavior; the shared component owns the unchecked white-background invariant
  - enforce canonical checkbox use through ESLint and keep shared state regression checks passing
