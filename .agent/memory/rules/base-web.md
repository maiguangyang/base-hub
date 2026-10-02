# AIPos Frontend Rules (`/base-web`)

> This document is the frontend appendix of the [AIPos Constitution](../constitution.md). All frontend development **MUST** comply with these rules.
> Default execution entrypoint: [../rules-core/base-web.md](../rules-core/base-web.md). Use the core file first, then return here for full explanations, edge cases, and appendices.

### Web AI Tool Impact Rule (WEB-AI-012)

- Every new or changed Web user-facing function **MUST** record the affected Engine operation IDs and AI tool IDs in its spec, plan, task, and review. State whether the tool was changed, verified against the new behavior, or unaffected with a concrete reason.
- A Web business action backed by a new or changed Engine operation **MUST** ship its permitted AI tool in the same task or prove existing tool coverage. Do not classify the operation `NON_CALLABLE` merely because the AI adapter was deferred; an intrinsic exclusion requires the user's explicit scope decision under Constitution §II-B.
- When the function changes shared operation validation, permissions, workspace or resource scope, results, side effects, or errors, update the affected callable tools and protected Engine behavior tests in the same task. A purely visual or client-local change without Engine business behavior needs only a brief no-impact reason.
- Every source-changing Web PR must include the structured `AI_TOOL_IMPACT` declaration. CI rejects omissions and malformed outcomes; review must verify that claimed no-impact or verified-no-change cases match the actual behavior and Engine contract.

## 1. Technology Stack

*   **Package Manager**: **pnpm**. `npm` (except `npx`) and `yarn` are prohibited.
*   **Core**: Astro (SSG/SSR) + React (Interactive Islands)
*   **Build**: Vite (configured via `astro.config.mjs`)
*   **UI System**: shadcn/ui + Tailwind CSS
*   **Path Aliases**: `@/` maps to `src/`
*   **Proxy**: `server.proxy` in `astro.config.mjs` forwards `/api` or `/graphql` to the local backend
*   **GraphQL Client**: **Apollo Client** (`@apollo/client`). Queries and mutations are defined in `features/{module}/graphql/`. Use `useMutation` and `useQuery` hooks from `@apollo/client/react`. **MUST** import the `gql` wrapper from `@/__generated__` instead of `@apollo/client` to enable Zero-config Type Inference. Wrap interactive islands with `<ApolloProvider>` for context access.
*   **GraphQL Field Selection**: Queries, mutations, subscriptions, and fragments **MUST** select only fields consumed by the current UI or business flow, state/cache updates, cache identity, or pagination. Apply the same rule to nested selections. Do not request whole entities or speculative fields for possible future use; remove a field when its consumer no longer uses it.

**Official Documentation Sources** (always consult before implementation):
*   **Astro**: https://docs.astro.build
*   **shadcn/ui**: https://ui.shadcn.com
*   **React Router**: https://reactrouter.com
*   **Tailwind CSS**: https://tailwindcss.com
*   **Zustand**: https://zustand.docs.pmnd.rs
*   **Apollo Client**: https://www.apollographql.com/docs/react

## 2. Code Quality Standards

*   **Linting**: All code **MUST** pass ESLint checks (`pnpm lint`). **Zero Error Policy**.
*   **Type Safety**: Strict TypeScript. `any` is **PROHIBITED** unless unavoidable (must include `// eslint-disable-next-line` with justification comment).
*   **Browser Overlay**: Development **MUST** have `vite-plugin-checker` enabled. Browser overlay errors must be fixed immediately.
*   **Language**: **TypeScript** is mandatory for all source code. `.js`/`.jsx` files are **PROHIBITED** except root-level config files (e.g., `astro.config.mjs`, `tailwind.config.js`).

### Coding Conventions

| Rule | Description |
|------|-------------|
| **Component Naming** | PascalCase for components (e.g., `UserListPage`), camelCase for hooks (e.g., `useAuth`) |
| **File Naming** | Component files match component name (e.g., `UserListPage.tsx`). Hook files match hook name (e.g., `useAuth.ts`) |
| **CSS Classes** | Use Tailwind utility classes. Custom CSS only when Tailwind cannot express the style. No inline `style={}` unless dynamic |
| **i18n Keys** | Dot-separated namespace format: `{feature}.{page}.{element}` (e.g., `auth.register.submit`) |
| **Export Style** | Named exports only. Default exports are **PROHIBITED** (except Astro pages/layouts which require them) |
| **Props Interface** | Name as `{ComponentName}Props` (e.g., `LoginFormProps`). Define in the same file as the component |

### Code Comment Rules

- 新增或修改的 **类型、接口、枚举、常量、Props、Hook、函数、方法、业务配置对象**，必须补充简洁准确的中文注释。
- 中文注释必须帮助后来者理解职责、业务含义、边界、输入输出、状态变化、副作用或存在原因，不能只机械重复代码字面意思。
- 参数、返回值、状态或 Hook 语义不直观时，声明注释必须说明这些输入输出的业务语义和使用边界。
- `component props should explain behavior`，而不是只描述视觉样式；复杂 Hook、表格列构造器、业务配置对象也必须解释为什么存在以及如何影响调用方。
- 测试里的 helper、mock builder、fixture 工厂、复杂断言函数，只要新增或修改了不直观的业务语义，也必须补充中文注释。
- 触达已有但缺少中文注释的声明时，必须顺手补齐到足以读懂的程度，不得继续扩大“无注释存量”。

| Prohibited | Required |
|------|-------------|
| Leaving touched web declarations undocumented | Add concise Chinese comments before considering the change complete |
| Writing comments that only rename the identifier | Explain responsibility, boundary, business meaning, state semantics, or side effects |
| Omitting non-obvious prop / hook / return semantics | Clarify behavioral meaning in the declaration comment when the code is not self-evident |

## 3. Directory Architecture

```text
base-web/src/
├── __generated__/           # Auto-generated GraphQL types & hooks (DO NOT EDIT MANUALLY)
├── components/
│   ├── ui/                  # shadcn/ui atomic components (business-agnostic)
│   └── shared/              # Cross-feature shared business components
├── features/                # Business logic layer (organized by domain)
│   ├── admin/
│   │   ├── components/      # Shared components (referenced by multiple pages)
│   │   │   ├── AdminLayout/   # Folder pattern
│   │   │   │   ├── index.ts
│   │   │   │   └── AdminLayout.tsx
│   │   │   ├── DataTable/     # Folder pattern
│   │   │   └── ThemeToggle.tsx # Simple components may remain as single files
│   │   ├── pages/           # Pages (Page Colocation)
│   │   │   ├── UserListPage/
│   │   │   │   ├── index.ts
│   │   │   │   ├── UserListPage.tsx
│   │   │   │   └── hooks/   # Page-private hooks
│   │   │   │       └── useUserColumns.ts
│   │   │   └── AdminDashboard/
│   │   ├── hooks/           # Shared hooks (referenced by multiple components)
│   │   ├── graphql/         # GraphQL queries/mutations
│   │   ├── config/          # Module configuration
│   │   ├── routes/          # Route entry points
│   │   └── types.ts         # Module type definitions
│   ├── auth/
│   │   ├── components/      # Shared components
│   │   │   └── LoginForm/
│   │   ├── pages/           # Pages
│   │   │   └── RegisterPage/
│   │   │       ├── index.ts
│   │   │       ├── RegisterPage.tsx
│   │   │       ├── RegisterForm.tsx  # Page-private sub-component
│   │   │       └── hooks/           # Page-private hooks
│   │   │           └── useRegisterForm.ts
│   │   ├── hooks/
│   │   ├── graphql/
│   │   ├── schemas/
│   │   └── types.ts
│   └── storefront/
├── layouts/                 # Astro layouts
├── pages/                   # Astro routing layer (must remain thin)
├── stores/                  # Zustand global state
│   └── slices/
├── lib/                     # General utility functions
└── i18n/                    # Internationalization
```

## 4. Feature-Based Architecture

*   **`src/features/`**: All business logic (components, hooks, API) **MUST** be grouped by domain (e.g., `features/auth`, `features/video`).

### 4.1 Page Colocation (Mandatory)

Every Feature module **MUST** adopt the Page Colocation pattern:
*   **`pages/`**: Each page is an independent folder (`PageName/index.ts` + `PageName.tsx`). Page-private components, hooks, and utilities **MUST** be moved into the page folder.
*   **`components/`**: Only retain components **referenced by multiple pages**. Shared components should use the folder pattern (`ComponentName/index.ts` + `ComponentName.tsx`).
*   **`hooks/`**: Shared hooks referenced by multiple components/pages stay in the module-level `hooks/` directory. Hooks used by only one page should be moved into that page's folder.

#### Ownership Rules

| Reference Count | Semantic Role | Ownership |
|----------------|--------------|-----------|
| Referenced by 1 page only | Page-private | Move into that page folder |
| Referenced by ≥ 2 | Shared | Keep in `components/` or `hooks/` |
| Layout/infrastructure level | Shared | Keep in `components/` |

### 4.2 Folder-as-a-Component (Mandatory)

*   All business components and pages **MUST** use the folder pattern: `ComponentName/index.ts` (re-export) + `ComponentName.tsx` (implementation).
*   Only these exceptions may remain as single files: components under 50 lines with no internal sub-components, `shadcn/ui` atomic elements.
*   `index.ts` is for re-exports only. **NEVER** write business logic in `index.ts`.

### 4.3 Other Architectural Rules

*   **Pages (`src/pages/`)**: **MUST** remain thin containers, only responsible for routing and passing data to Features. Logic complexity in pages is prohibited.
*   **UI Lib (`src/components/ui`)**: **ONLY** generic, business-agnostic shadcn/ui components.
*   **Shared (`src/components/shared`)**: Cross-feature business components.

### 4.3.1 Source Cohesion (Mandatory)

- Run `pnpm lint` after changing handwritten Astro, TypeScript, or TSX source, including tests and closures. New files and functions must meet the shared 250/50 effective-line, complexity-10, and nesting-3 limits. Generated output is excluded.
- Split a page or component by a named UI responsibility, hook, or transformation while keeping one-consumer code in its page/component folder. The folder-as-a-component rule does not authorize one large `PageName.tsx` or `ComponentName.tsx`.
- Promote a helper to a feature or cross-feature shared location when at least two independent consumers need the same semantics, or when a stable infrastructure contract is already shared. Avoid one-line forwarding helpers and speculative `src/lib` accumulation.
- Existing over-limit units use a recorded non-growing baseline. The copied `scripts/translate.ts` also has a per-function frozen baseline checked by `pnpm lint`; preserve that file until translation workflow work explicitly revisits it. A justified exception must name the affected unit and its reason in the plan; do not blanket-exempt handwritten folders.

### 4.4 后台列表状态统一规范（Mandatory）

*   后台列表页的“主状态”必须使用共享 `AdminStatusSwitch`，禁止页面直接使用原始 `Switch` 处理启用/禁用状态。
*   审计状态、删除状态、生命周期状态必须使用共享 `AdminStatusIndicator`，禁止伪装成交互开关。
*   禁止在后台页面中手写状态颜色 class；状态 tone 必须来自共享状态映射。
*   同一资源列表禁止同时出现“状态 Badge”和独立“启停列/按钮”的重复表达。
*   后台用户可见文案不得长期混用英文 fallback；列表页、弹窗、toast 和确认弹窗的 fallback 必须优先使用中文。

### 4.5 Admin Local Time Formatting Rules

*   Admin timestamps must use a shared device-local formatter instead of scattered inline formatting in admin pages.
*   `formatAdminLocalDateTime` and `formatAdminLocalDateTimeInputValue` are the default shared helpers for admin timestamp rendering and `datetime-local` form values.
*   Admin pages must not inline `format(new Date(Number(...)))` for backend millisecond timestamps once the shared formatter exists.
*   Viewer device local timezone is the default display contract for backend timestamps in `base-web`.
*   When a task touches admin timestamp rendering or timestamp form hydration, that touched chain must be brought into compliance with the shared local-time formatter rule before completion.

### 4.6 后台 Admin 列表页结构宪章（Mandatory）

*   后台 Admin 列表页必须统一采用页面层 + 工作台层的双层骨架：页面层负责 `RequirePermission`、`AdminPageHeader`、页面级 `AdminStatsStrip` 与工作台容器；工作台层负责 `AdminTableShell`、`AdminToolbar`、列表内容、空态、分页与行内工作流。
*   `AdminPageHeader` 和页面级 `AdminStatsStrip` 必须留在页面组件中，禁止把标题说明区或 summary 区重新塞回表格组件内部。
*   页面工作台必须使用 `AdminTableShell`，筛选栏必须使用 `AdminToolbar`；新增页面或触达中的后台列表页不得继续手写独立 toolbar 骨架。
*   任何资源型后台列表页只要存在主新增 / 主创建动作，都必须使用 `AdminPrimaryActionButton`，并通过 `AdminToolbar` 的 `rightSlot` 注入筛选栏右侧；禁止继续依赖 `actions` 兼容口子承载主新增入口。
*   只读 / 审计型后台列表页可以没有主新增动作，但仍必须保留统一的 `AdminPageHeader`、页面级 `AdminStatsStrip`、`AdminTableShell` 与 `AdminToolbar` 骨架。
*   `TraceLogListPage` 这类页面级全局控制属于页头级动作，不属于资源主新增入口；此类控制只能放在 `AdminPageHeader.actions`，不能替代 `rightSlot` 主 CTA 规范。
*   后台列表页摘要统一以 `AdminStatsStrip` 为标准，不再允许保留自定义 summary 卡片作为长期结构特例；触达中的历史页面必须同步收口。
*   空态按钮可以作为引导性补充操作保留，但 canonical 主新增入口仍然必须在 `AdminToolbar.rightSlot`。
*   后台列表结构契约必须自动发现 `src/features/**/*ListPage.tsx`，禁止维护容易漏页的手工页面白名单。

### 4.7 后台 Admin 新增 / 编辑弹窗结构宪章（Mandatory）

*   后台新增、创建和编辑表单弹窗必须使用共享 `AdminFormDialogShell`，禁止在页面中手写 `DialogContent`、标题区和底部操作区。
*   自动契约必须递归扫描 `src/features/**/*Dialog.tsx`，并通过 AST 识别含表单控件的弹窗，防止文件命名变化或新页面绕过共享弹窗结构。

### 4.8 后台 Keep-Alive 路由上下文宪章（Mandatory）

*   后台采用多标签页 keep-alive 渲染栈：所有已打开标签的页面组件同时挂载，非活动标签用 `hidden` 隐藏而非卸载。
*   实测依据（`react-router@8.4.0`，声明式 `MemoryRouter` 与 `createMemoryRouter` 行为一致）：
    *   `useLocation`、`useSearchParams`、`useMatches` 返回的是**当前活动标签**的值，与调用方无关。
    *   `useParams` 按标签隔离，返回值正确——`renderMatches` 为每个 match 建立了独立的 `RouteContext`。
*   危险之处在于失败是静默的：隐藏标签调用 `useLocation` 不抛错、不告警、返回值类型正确、页面照常渲染，只是内容属于另一个标签。
*   因此后台页面与页面私有组件**禁止**直接调用 `useParams`、`useLocation`、`useSearchParams`、`useMatch`、`useMatches`，必须通过 `useAdminTab()` 读取自身标签的 `path`、`search`、`params`、`isActive`。
*   **`useParams` 实测安全但仍在禁用之列**，理由不是机制危险，而是两点：其一，页面开发者只需记住「路由上下文一律走 `useAdminTab()`」这一条规矩，不必维护例外；其二，它的正确性依赖 `renderMatches` 建立 `RouteContext` 这一实现细节，渲染栈若改用其他方案，`useParams` 会静默失效且无任何检查能拦住。此处如实记录，避免后人从本条读到错误的机制认知。
*   `useNavigate` 不在禁用范围内。它执行全局跳转动作，不读取调用方上下文。
*   上述 hook 仅允许出现在后台路由外壳内部：`features/admin/routes/**`、`features/admin/components/AdminLayout/**`、`features/admin/components/TabStack/**`、`features/admin/components/AdminTabBar/**`、`features/admin/hooks/useAdminTab.ts`。
*   面包屑、标签栏与侧边栏选中态必须共用标签 store 的 `activePath` 作为唯一真值源。
*   后台页面导航必须使用 react-router 的 `Link` 或 `useNavigate`，禁止裸 `<a href>`：后者触发整页重载，SPA 连同 keep-alive 渲染栈一起重建，滚动位置与表单草稿全部丢失。
*   本规则必须由 ESLint `no-restricted-imports` 强制，禁止仅依赖代码评审与文档重读。
*   `src/components/ui/**` 为 shadcn CLI 拉取的 vendored 源码，已豁免 WEB-COHESION-005 的 `max-lines`、`max-lines-per-function` 与 `complexity` 三条规模类规则；`max-depth` 与其余检查保持生效。豁免严格限定在该目录，不得扩大到 `src/features/**` 或 `src/components/shared/**`。
*   `src/hooks/` 为 shadcn CLI 按 `components.json` 的 `hooks` 别名写入 vendored hook（如 `use-mobile.ts`）的目录，不承载本项目手写业务 hook；业务 hook 仍归 feature 或页面所有。

### 4.9 后台表单控件 Placeholder 规范（Mandatory）

*   后台 Admin 管理端所有单行与多行文本输入控件（包括 `<Input>`、`<textarea>`、搜索栏、弹窗输入框、认证表单输入框等）必须明确配置非空的 `placeholder` 占位文案。
*   占位提示文案必须准确引导用户输入，遵循统一表述规范（例如 `请输入门店名称`、`请输入角色名称`、`搜索待审门店` 等）。
*   在已接入国际化语言包的模块（如 `auth`），`placeholder` 必须提取至对应语言包（如 `zh-CN.json`），禁止页面手写硬编码字符串。
*   触达已有表单或新增弹窗/页面时，必须确保所有输入框补齐 `placeholder`，禁止留空或遗漏。

### 4.10 后台表单与行内标签 flex flex-col gap-5 结构宪章（Mandatory）

*   后台 Admin 管理端所有表单容器（包括 `<form>` 根节点、弹窗内字段列表容器 `<div>` 等）必须统一使用 `flex flex-col gap-5` 规范布局。
*   严禁在表单和字段容器中继续使用遗留的 `space-y-*`（如 `space-y-5`、`space-y-4`、`space-y-3`）选择器类，以避免 margin 塌陷和隐藏元素间距异常。
*   表单行内 `<label>` 元素必须使用 `flex flex-col` 纵向结构，根据场景配置紧凑间距（认证卡片为 `gap-4`，弹窗表单为 `gap-2`），严禁在 `<label>` 内手写 `block space-y-*`。
*   触达已有表单或新增弹窗/页面时，必须确保表单结构遵循 `flex flex-col gap-5` 契约。

### 4.11 后台表单禁止暴露“编号”/“编码”字段规范（Mandatory）

*   后台 Admin 管理端所有表单及弹窗严禁向用户暴露人工录入或维护的“编号”、“编码”（如门店编码、组织编码等）字段。
*   业务实体的编号/编码属于系统底层逻辑标识，人工录入既无业务价值又增加使用与维护负担。
*   当底层 API / GraphQL Mutation 强要求 `code` 字段时，必须由前端适配层（如 action 助手函数）基于 `generateEntityCode` 自动生成合规、唯一的内部编码，对终端用户完全透明无感。
*   表单状态默认仅维护名称等核心业务字段，更新操作时仅提交修改的业务字段，禁止向用户呈现编号输入框。
*   触达已有表单或新增弹窗/页面时，严禁出现任何用户填写的编号/编码类输入控件。

### 4.12 后台列表与表单表面统一视觉契约（Mandatory，WEB-UI-009）

*   新增或修改后台列表、表格、筛选工具栏以及新增/编辑弹窗时，必须使用现有共享组件及其默认视觉规范，不得在页面层复制或另建一套后台设计系统。
*   列表具备搜索能力时必须使用 `AdminListSearch`：输入仅更新本地草稿，用户按 Enter 或点击“搜索”后才提交去除首尾空格的关键词；空关键词用于清除筛选，禁止输入时逐字符触发远端查询；业务 API 不支持搜索时不得为了统一外观添加无效筛选器。
*   `AdminToolbar`、`AdminTableShell`、`AdminPagination` 与 `AdminPrimaryActionButton` 分别拥有工具栏表面、表头和行间距、分页表面以及主动作高度的唯一视觉真值；页面不得通过 `className` 覆盖工具栏外观，也不得重复共享表格间距或 hover 样式。
*   表格操作列必须右对齐，行操作使用紧凑尺寸与语义变体：编辑及可逆工作流使用描边或默认样式，删除等不可逆动作使用危险样式。
*   新增/编辑表单弹窗必须使用 `AdminFormDialogShell`，保持固定标题区、独立滚动主体、固定底部操作区、可见标签、有效 placeholder，以及 `flex flex-col gap-5` 的字段布局；禁止直接使用 `Dialog.Content` 搭建平行外壳。
*   白底黑字契约适用于全部 Web 后台用户输入控件和可见表单项表面，不限于弹窗或新增/编辑场景；页面级新增/编辑、配置/规则页、抽屉、内联动作面板、列表筛选、禁用/只读状态和 Portal 选项层全部必须遵循。输入框、多行文本、下拉触发器、时间控件、可搜索/自定义下拉和复合字段容器均在此范围内。
*   白底黑字必须由 admin 作用域的典型样式或共享表单组件持有唯一视觉真值，让普通消费者自动继承；在业务页重复添加 `bg-white` / `text-black` 不算完成契约，也不得代替共享所有权。确有不同交互语义的例外必须声明明确语义挂钩、说明原因并提供聚焦回归。
*   自定义下拉触发器、Portal 选项层与普通 `div` / `button` 组成的复合表单项必须显式接入共享表单表面语义；Portal 不得依赖弹窗后代样式，复合控件不得依赖评审者从通用标签推测其表单语义。
*   暂停等需要填写原因的确认操作属于动作确认，使用完整的 shadcn `AlertDialog` 结构及已有输入、按钮组件；这类确认弹窗不必伪装成新增/编辑表单，但必须保留校验、取消、提交中防重复和失败反馈。
*   长列表或长选择器中的“全选”等批量控制必须固定在其自身滚动容器顶部，不得随选项一起滚出可视区域。
*   视觉 token、间距和表面样式的调整只能修改共享组件，使所有消费者自动继承；宪章只冻结行为与责任边界，不冻结具体像素或 Tailwind 实现细节。
*   `AdminToolbar` 与 `AdminPrimaryActionButton` 的组件类型不得开放 `className` 等页面级视觉覆盖入口；自动契约必须使用 TypeScript AST 递归发现列表、弹窗与其他表单表面，拦截工具栏内的原始搜索输入、表单弹窗缺少共享壳以及表单弹窗直接渲染 `Dialog.Content`；表单表面检查必须按实际渲染的控件语义发现页面级新增/编辑、配置、规则、抽屉和内联面板，不得仅依赖 `*Dialog.tsx` 等文件名约定。契约同时必须以浏览器计算样式验证浅色、深色与 Portal 渲染的白底黑字；禁止使用手工页面白名单或正则源码扫描替代 AST。

### 4.13 后台列表结构化筛选完整性（Mandatory，WEB-UI-010）

*   每一个后台实体列表，包括页面内嵌的批次、流水、价格历史和发券记录列表，以及文件名不以 `ListPage.tsx` 结尾的页面，都必须依据实体字段、可见列与工作流审查必要且便于操作的筛选；存在有业务价值的状态、类型、关联实体、布尔值或可见时间字段时，禁止只提供通用关键词搜索。
*   列表筛选区必须使用 `AdminListFilters` 与共享筛选控件；关键词继续使用提交式 `AdminListSearch`，结构化选择属于明确用户动作，可以立即生效。
*   关键词与结构化条件必须共同映射到服务端查询并按交集组合，严禁仅过滤当前已加载分页；任何条件变化或统一重置都必须返回第 1 页，重置不得改变页长。
*   关联目录必须确认完整后才能作为“全部选项”使用；超大目录应提供服务端搜索或分页；加载失败必须有可用的重试操作，禁止将固定前 N 条或失败后的残缺目录伪装为完整筛选项。
*   禁止添加 API 不支持、用户不可理解、与页面展示或工作流无关的装饰性筛选器；确实没有有价值结构化字段的页面可以只保留搜索，但必须在设计或测试中明确说明例外。
*   自动列表契约必须递归发现所有后台列表，包括 `*ListPage.tsx` 与其他渲染 `AdminTableShell` 的页面，并要求规范导入、渲染 `AdminListFilters`；AST 只约束共享边界，不得臆测业务筛选是否语义完整。页面测试必须验证条件传参、组合、翻页重置、统一清空及必要的错误恢复；语义完整性由本规则、页面测试和评审共同保证。

### 4.14 shadcn/ui 组件优先（Mandatory，WEB-UI-011）

*   新建或修改交互组件、输入控件、弹窗前，必须先查 [shadcn/ui 官方组件目录](https://ui.shadcn.com/docs/components)，确认是否已有满足需求的官方组件；仓库本地已有时直接复用。
*   官方组件存在但 `base-web/src/components/ui` 尚未安装时，必须从 `base-web` 目录按 `components.json` 配置通过 shadcn CLI 安装，再在业务代码中使用；不得因为本地暂缺就手写同类基础组件。
*   业务专用组件可以组合、封装 shadcn 原语，但交互语义、键盘操作、焦点管理和基础视觉必须由官方组件承担。仅当官方目录没有合适组件，或官方组件无法满足明确的业务约束时，才可自建基础组件，并在设计或变更说明中记录原因。
*   需要用户输入或确认的 Web 操作必须使用 shadcn `Dialog`、`AlertDialog` 等合适组件；存在可用组件时禁止调用浏览器原生 `window.prompt`、`window.confirm`、`window.alert` 作为产品交互。

### 4.15 下拉控件占位文案与业务选项分离（Mandatory，WEB-UI-013）

*   所有 Web 下拉选择器（包括列表筛选、表单单选与可搜索选择器）的 `placeholder` 只负责未选择时的控件提示，严禁自动加入选项列表，也不得作为禁用选项展示。
*   shadcn `Select` 必须通过 `SelectValue.placeholder` 显示未选择提示，选项只渲染真实业务项；禁止用占位文案构造 `SelectItem`、原生 `<option>` 或 `role="option"` 项。
*   “顶级分类”“不设置品牌”“保留现有包装”等确有业务含义的空值选项必须独立显式声明（共享控件使用 `clearLabel`），不得从 `placeholder` 推导。占位提示和业务选项各自描述自己的语义。
*   `AdminFilterSelect` 不再隐式插入“全部”或“清空”选项；列表通过 `AdminListFilters` 的“重置筛选”统一清空条件，恢复占位提示并返回第 1 页。
*   控件测试必须验证占位提示不进入选项、真实项仍可选择、重置恢复提示，以及显式业务空值项的映射。ESLint 必须拦截直接把 `placeholder` 表达式渲染为下拉选项的写法。

### 4.16 复选框未勾选状态统一白底（Mandatory，WEB-UI-014）

*   所有 Web 页面与弹窗的复选框必须使用 `@/components/ui/checkbox`；禁止业务页面手写原生 `input type="checkbox"` 或直接使用其他复选框原语。
*   未勾选状态固定纯白背景，浅色与深色主题、悬停与聚焦时均不得使用页面背景色、灰色或透明背景。统一组件在 `data-state="unchecked"` 时强制白底，业务页面不得通过 className、style 或 CSS 覆盖这一状态。
*   选中、半选、禁用、焦点与键盘交互保持各自既有语义；不要为修复未勾选背景而改变这些行为。
*   ESLint 拦截原生复选框写法；共用组件状态检查与规则契约检查纳入持续验证。此为纯视觉与组件复用约束，不改变 Engine 操作或 AI 工具行为。

## 5. State Management

*   **Solution**: **Zustand** (v5.x) with **Slices Pattern**.
*   **Directory**: Global stores **MUST** be in `src/stores/`. Each slice in `src/stores/slices/`.
*   **Middleware Stack**: `persist` + `devtools` + `immer` in that order.
*   **Selectors**: **MUST** use atomic selectors. Components subscribe to specific properties, **NEVER** the entire store.
*   **Prohibited**: Do NOT introduce Redux, MobX, Jotai, Recoil, or Context-based global state. Zustand is the only allowed solution.

## 6. Development Workflow

### 6.1 Manual Review Commit Protocol

- After completing implementation and verification in `base-web`, the agent must leave all modified files in the local working tree and keep them unstaged for human review.
- The agent must not automatically run `git add`.
- The agent must not automatically run `git commit`.
- The standard handoff summary must include:
  - which files changed
  - which verification commands ran and what passed
  - the current working-tree status
  - a reminder that the developer performs review, staging, and committing manually
- If unrelated local changes already exist in the repository, the agent must not use automatic staging, selective staging, or any other staging workaround.

1.  **Island Strategy**: Default to `.astro`. Use `.tsx` ONLY when client-side interactivity (e.g., `client:load`) is required.
2.  **Component Integration**: shadcn components go in `base-web/src/components/ui`.
3.  **Type Sync**: Use `graphql-codegen` (`pnpm codegen`) to generate TypeScript types from the Backend Schema into `src/__generated__/`. The manual usage of `any` or manual type interfaces for `useMutation`/`useQuery` data is prohibited; rely exclusively on the inference provided by the generated `gql` function.
4.  **Verify**: Run `pnpm exec astro check` (type check) and `pnpm build` (production build) to ensure zero errors before marking any task as complete.

### Error Recovery

If verification fails, follow this troubleshooting sequence:

| Error | Resolution |
|-------|-----------|
| `astro check` type errors | 1) Fix the reported type errors in order 2) Check import paths after file moves 3) Verify `tsconfig.json` path aliases |
| `pnpm build` fails | 1) Fix any SSR/SSG errors in `.astro` files 2) Check for client-only APIs used in server context 3) Verify all `client:*` directives are correct |
| Missing module errors | 1) Run `pnpm install` 2) Verify `@/` alias resolution in `tsconfig.json` 3) Check for circular imports |

## 7. Checklists

### 7.1 New Feature Module Checklist

Example: Adding a new feature module `video`:

- [ ] Create feature directory `src/features/video/`
- [ ] Define module types in `video/types.ts`
- [ ] Create `video/graphql/` with queries and mutations
- [ ] Create `video/hooks/` for shared hooks
- [ ] Create `video/pages/` with page folders (Page Colocation)
- [ ] Create `video/components/` for shared components (Folder-as-a-Component)
- [ ] Register routes in `video/routes/` (if SPA sub-app) or create Astro pages in `src/pages/`
- [ ] Add i18n keys to `src/i18n/locales/{lang}/`
- [ ] Verify with `pnpm exec astro check` and `pnpm build`

### 7.2 New Page in Existing Module Checklist

Example: Adding a new page `VideoDetail` in existing `video` module:

- [ ] Create folder `features/video/pages/VideoDetailPage/`
- [ ] Create `VideoDetailPage.tsx` with component implementation
- [ ] Create `index.ts` with `export { VideoDetailPage } from './VideoDetailPage';`
- [ ] Move any page-private hooks into `VideoDetailPage/hooks/`
- [ ] Move any page-private components into the folder
- [ ] Create Astro page in `src/pages/[lang]/video/[id].astro` (if Astro-routed) OR register route in `features/video/routes/` (if SPA sub-app)
- [ ] Add i18n keys for the new page
- [ ] Verify with `pnpm exec astro check` and `pnpm build`

**Version**: 1.13.0 | **Date**: 2026-09-27
