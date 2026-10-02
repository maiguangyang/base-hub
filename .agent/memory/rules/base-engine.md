# AIPos Backend Rules (`/base-engine`)

**Version**: 2.2.0

> This document is the backend appendix of the [AIPos Constitution](../constitution.md). All backend development **MUST** comply with these rules.
> Default execution entrypoint: [../rules-core/base-engine.md](../rules-core/base-engine.md). Use the core file first, then return here for full explanations, edge cases, and appendices.

## Engine Priority Zero: Entity Relationships Must Use `@relationship`

This rule is the first schema decision for every `@entity`. It takes precedence over convenience, legacy modeling habits, generated field names, and entity-category-specific shortcuts.

- Every relationship between two `@entity` types in handwritten `model/*.graphql` **MUST** be expressed with typed object/list fields and `@relationship(inverse: ...)`.
- A handwritten entity **MUST NOT** declare or use a scalar `xxxId`, `xxxIds`, `<entity>Id`, `<entity>Ids`, or a renamed scalar alias to associate, join, query, or enforce referential linkage to another entity/table.
- This prohibition covers runtime, history, audit, ledger, session, governance, and explicit association entities. If an association entity has independent business facts, its links to both endpoint entities must still be `@relationship`.
- Dolphin-generated FK columns, generated relationship ID projections, and generated mutation input keys are generator-owned implementation details. They may be consumed through the generated contract but **MUST NOT** be copied into handwritten entity fields.
- IDs in `extend.graphql` custom operation inputs or payloads may identify request targets or response resources. They do not define persistence relationships and **MUST NOT** be persisted or reused as substitutes for `@relationship`.
- A field claimed to be an external identifier or immutable snapshot is non-relational only if no code uses it to look up, join, navigate to, or enforce integrity against another project entity/table. Otherwise it is a relationship and must use `@relationship`.
- Schema design, plans, reviews, and implementation **MUST stop** when this rule is violated. Contract tests must reject custom scalar cross-entity links so compliance does not depend on naming conventions or reviewer memory.

```graphql
# Good: the association is explicit and typed.
account: Account! @relationship(inverse: "sessions")

# Bad: a handwritten scalar is acting as a foreign key.
accountId: String! @column(gorm: "type:varchar(36);NOT NULL;index;")
```

## 1. Technology Stack

*   **Framework**: **[dolphin](https://github.com/sj-distributor/dolphin)** (High-Performance Golang base-web base-engine)
    *   *Note*: This framework relies heavily on code generation.
*   **Protocol**: GraphQL
*   **Database**: **MySQL**
*   **Data Access**: Use the framework-generated DAL (Data Access Layer). Do not write raw SQL unless justified for extreme performance.

**Official Documentation Sources** (always consult before implementation):
*   **Go**: https://go.dev/doc/
*   **GORM**: https://gorm.io/docs/
*   **GraphQL**: https://graphql.org/learn/

## 2. Common Commands

> All commands MUST be run from the `base-engine/` directory. AI/agents may run `make generate` only with explicit user authorization for the current task; generated files must never be edited directly.

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `make generate` | Regenerate `gen/` from handwritten `model/` | Only after task-scoped user authorization and schema-first edits; inspect the full diff |
| `make check_structure` | Check handwritten Go source cohesion | After changing Go source |
| `make migrate` | Run database schema migration | After adding/changing entity fields |
| `make start` | Start dev server with CORS enabled (port 8080) | Local development & testing |
| `make redis` | Start Redis Docker container (port 6379) | First-time setup |
| `make mysql` | Start MySQL Docker container (port 3306) | First-time setup |

## 3. Directory Architecture

```text
base-engine/
├── auth/                  # Authentication & Authorization
├── config/                # Global configuration constants
├── docs/                  # API documentation
├── enums/                 # Enums & constants
├── gen/                   # Auto-generated code (DO NOT modify)
├── model/                 # GraphQL Schema definitions
├── src/                   # Core business logic
│   ├── dbup/              #   Database initialization
│   ├── middleware/        #   HTTP middleware
│   ├── plugins/           #   GORM query plugins
│   └── services/          #   Business service layer
│       └── <module>/      #   Organized by business module
├── utils/                 # Utility functions
└── main.go                # Entry point
```

**Directory Rules**:
*   **`gen/`**: Auto-generated code. AI/agents **MUST NOT** edit, format, or delete files here directly. With explicit authorization for the current task, update handwritten `model/`, run `make generate`, then inspect the complete generated diff. Without authorization, stop before changing schema that requires generated output.
*   **`model/`**: GraphQL Schema definitions are the **single source of truth** for the API contract (Schema-First).
*   **`src/services/`**: Business logic **MUST** be organized by module directory (e.g., `services/user/`, `services/order/`). Flat-filing all services in one directory is prohibited.
*   **`auth/`**: Authentication/authorization logic must be centrally managed. Scattering auth code across services is prohibited.

### 3.1 System Auto-Generated Fields

Dolphin automatically generates the following fields for every `@entity` type. **Do NOT define them manually**:

| Field | Type | Description |
|-------|------|-------------|
| `id` | varchar(36) | UUID primary key |
| `createdAt` | bigint(13) | Creation timestamp (milliseconds) |
| `updatedAt` | bigint(13) | Update timestamp (milliseconds) |
| `deletedAt` | bigint(13) | Deletion timestamp (soft delete) |
| `createdBy` | varchar(36) | Creator ID |
| `updatedBy` | varchar(36) | Last updater ID |
| `deletedBy` | varchar(36) | Deleter ID |
| `isDelete` | int(2) | Soft-delete flag: 1=active, 2=deleted |
| `weight` | int(2) | Weight (for sorting) |
| `state` | int(2) | Status: 1=active, 2=disabled |

> [!CAUTION]
> **NEVER** redefine the above fields in `model/*.graphql`. Doing so will cause generation conflicts.

### 3.2 Data Initialization (dbup)

When adding a new `@entity` type, you **MUST** also register its CRUD permissions in `src/dbup/init.go`:

1. Add 4 entries for the new entity module using the standard CRUD action set:
   ```go
   {Name: "查看{Entity}", Action: "{module}:read", Module: "{module}"},
   {Name: "创建{Entity}", Action: "{module}:create", Module: "{module}"},
   {Name: "编辑{Entity}", Action: "{module}:update", Module: "{module}"},
   {Name: "删除{Entity}", Action: "{module}:delete", Module: "{module}"},
   ```
2. Classify the seeds into the correct permission bucket:
   - platform governance permissions **MUST** be registered in `systemPermissionSeeds()`
   - permissions intended to be assignable to `CompanyRole` **MUST** be registered in `tenantBusinessPermissionSeeds()`
   - `tenantBusinessPermissionSeeds()` may be empty until real tenant business modules exist
   - do **NOT** mix tenant-assignable permissions back into the governance seed bucket
   - `Permission.scope` is the backend-authoritative permission range field and **MUST** match the seed bucket classification
   - seeds in `systemPermissionSeeds()` **MUST** use `PermissionScopeSystem`
   - seeds in `tenantBusinessPermissionSeeds()` **MUST** use `PermissionScopeTenant`
2. `InitRolePermissions()` will automatically bind new permissions to the ADMIN role on next server restart.

The migration bootstrap path in `main.go` **MUST** continue to execute the seed initialization chain after `automigrate()`:
- `dbUp.InitRoles()`
- `dbUp.InitPermissions()`
- `dbUp.NormalizePermissionScopes()`
- `dbUp.InitRolePermissions()`

> [!IMPORTANT]
> Failing to register permissions means the new entity's CRUD operations will lack fine-grained permission control, even if `@hasRole` is set.

> [!IMPORTANT]
> `Permission.scope` is **not** a frontend-managed or operator-managed field. Generated `createPermission` / `updatePermission` writes must keep it read-only to clients, and the authoritative value must come from backend seeds plus normalization logic.

> [!IMPORTANT]
> 平台级外部服务注册类实体（例如 `McpHub`）默认属于治理侧实体，CRUD 权限必须进入 `systemPermissionSeeds()`。这类实体只能保存接入元数据（如服务链接、鉴权类型、登录跳转链接），不得在后台持久化用户侧真实凭据，真实授权流程应交给客户端或专用授权模块处理。

> [!IMPORTANT]
> MCP 能力治理必须保持两层分离，并由技能发布链路决定公司可达性：
> - `McpHub` = 平台级注册表，只保存统一元数据
> - `UserMcpPlugin` = 绑定 `CompanyMember` 的成员级授权材料与状态层，只记录启用/停用/待授权/授权失效等状态
>
> 公司是否可见某个 Hub 不再通过独立分配表维护，而是从 `CompanySkill -> SkillGroup -> Skill.primaryMcpHub` 的 `@relationship` 路径派生；不得用 `primaryMcpHubId` 之类的手写标量字段替代该实体关系。
>
> 这两层中：
> - `mcpHub:*` **MUST** 归入 `SYSTEM` 权限范围
> - `userMcpPlugin:*` **MUST** 归入 `TENANT` 权限范围
> - `UserMcpPlugin` **MUST NOT** 绑定裸 `User`
> - 任意一层都不得存储真实凭据，真实授权材料必须留在客户端或专用授权模块处理

## 4. Development Standards

### 4.1 Pattern Selection

| Scenario | Recommended Pattern | Description |
|----------|-------------------|-------------|
| **Custom Query/Mutation** (e.g., login, ossPolicy) | Full Triple Pattern | Requires `model/extend.graphql` + `src/extend.go` + triple |
| **Override Generated CRUD Handler** (e.g., createMerchant) | Handler Only | Just write a `xxx_handler.go` and register to `resolver.Handlers` |

> [!TIP]
> **How to decide**: If the operation is defined in `extend.graphql` (custom API), use **Full Triple**. If you are overriding a CRUD operation that already exists in `gen/`, use **Simplified**.

> [!CAUTION]
> **Transaction Management Rule — Depends on Resolver Entry Source**
>
> | Pattern | Resolver Entry | Who manages transactions? |
> |---------|---------------|---------------------------|
> | **Full Triple (Custom API)** | `src/extend.go` (hand-written) | ✅ **Handler MUST** manage transactions |
> | **Simplified (Override CRUD)** | `gen/resolver-mutations.go` (auto-generated) | ❌ **Handler MUST NOT** call transaction functions |
>
> Key difference: The auto-generated Resolver entry already calls `EnrichContextWithMutations` / `FinishMutationContext` / `RollbackMutationContext`, while the hand-written `extend.go` entry does not.

### 4.2 GraphQL Optional Input Pointer Types

When writing custom handlers or overriding generated handlers that accept a `map[string]interface{}` (usually named `input`), **be extremely careful with type assertions for optional String or ID fields.**

GraphQL base-engines (like `gqlgen` used by Dolphin) often pass **optional** fields as **pointers** (`*string`, `*int`, etc.) rather than underlying values to distinguish between `null` and computationally omitted fields.

**DO NOT** simply assert as `string`:
```go
// ❌ WRONG: Will panic or return empty string "" if the value is a pointer
val, _ := input["optionalField"].(string) 
```

**DO** assert both pointer and value types safely:
```go
// ✅ CORRECT: Handle both generic string and pointer types securely
var val string
if ptr, ok := input["optionalField"].(*string); ok && ptr != nil {
	val = *ptr
} else if str, ok := input["optionalField"].(string); ok {
	val = str
}
```

### 4.3 Template Placeholders

Before using any code template below, resolve these placeholders:

| Placeholder | How to Resolve | Example |
|-------------|---------------|---------|
| `{module_path}` | Read `base-engine/go.mod`, use the `module` line value | `github.com/maiguangyang/triton` |
| `{domain}` | Lowercase service module name | `user`, `order`, `merchant` |
| `{Domain}` | PascalCase service module name | `User`, `Order`, `Merchant` |

### 4.3.1 Conversation Error Code Contract

- Conversation-domain user-visible failures **MUST** originate from one centralized conversation error-code vocabulary under `src/services/conversation/`.
- Conversation-domain GraphQL failures **MUST** surface the stable code through `extensions.code`.
- Conversation-domain realtime failure payloads such as `turn.failed`, `tool_use.failed`, and `session.error` **MUST** include `errorCode`.
- Conversation-domain `errorMessage` text **MUST** be treated as fallback-only text for old clients; it **MUST NOT** be the primary cross-end contract.
- Backend code **MUST NOT** rely on matching localized message text to coordinate app behavior.
- Adding a new conversation user-visible failure mode **MUST** update the centralized error-code vocabulary and add/update regression coverage for both GraphQL and realtime outputs.

### 4.3.2 Language-Neutral Conversation & Prompt Rules

- Engine conversation orchestration, tool routing, prerequisite normalization, selection recovery, confirmation follow-ups, and similar free-form text paths **MUST NOT** depend on hardcoded Chinese/English phrases, locale-specific regex, or handwritten stopword tables as the primary behavior gate.
- When engine behavior must be derived from free-form text, the implementation **MUST** prefer one of these boundaries:
  - schema-driven structured fields
  - published `SkillRuntimeLexicon` data for the active locale
  - validated model output with a strict JSON contract
  - explicit machine-readable markers in schema/config, rather than natural-language prose parsing
- New prompt contracts, built-in prompt examples, and planner templates **MUST NOT** seed one fixed human language for user-facing fields such as `statusNote`, `failureMessage`, `clarification`, follow-up suggestions, or confirmation text.
- User-facing orchestration microcopy generated by prompt planners **MUST** explicitly match the language of the latest user request or active locale.
- New engine code **MUST NOT** silently infer business filters such as IDs, date presets, or entity selections from localized free-form phrases unless that parsing contract is explicit, language-agnostic, and regression-tested as a deliberate product requirement.
- Backward-compat handling for legacy localized phrases may remain temporarily only when migration is in progress, but it **MUST** be documented as compatibility code and **MUST NOT** expand into new features.

### 4.3.3 System Prompt Markdown Location Rule

- Shared System Prompt markdown assets **MUST** live under `base-engine/system_prompt/`.
- Any markdown file that is compiled, embedded, or loaded through the shared `gitlab.sjfood.us/oxygen/base-engine/system_prompt` package **MUST NOT** be stored under `src/services/**`, `tmpdebug/`, or other feature-local directories.
- Service-local prompt assets that are not part of the shared System Prompt package may stay colocated with their owning module, but they **MUST NOT** be treated as shared System Prompt assets by convention alone.
- When promoting an existing module-local prompt into a reusable System Prompt asset, the markdown file **MUST** be moved into `base-engine/system_prompt/` as part of that change.

### 4.3.4 AI Tool Contract Quality Rule (ENG-AI-009)

- Every externally consumable Engine operation **MUST** have a classified contract record under Constitution §II-B, including operations used only by App. A new or changed business capability **MUST** deliver its permitted executable AI tool in the same task, or prove existing tool coverage through protected tests. `NON_CALLABLE` cannot defer tool work; an intrinsic exclusion requires an operation-level reason and the user's explicit scope decision before implementation. An approved `NON_CALLABLE` record does not require an executable tool definition.
- Every `CALLABLE` tool **MUST** provide a stable unique `name` matching its registered invocation name, a concise human-readable `title`, and a separate detailed `description` that explains purpose, workspace or resource scope, and effects. Generic duplicate titles do not satisfy the description requirement.
- The input schema **MUST** specify types, required fields, applicable enum/range bounds, and meaningful descriptions for every top-level and nested parameter. The executable tool **MUST** reject wrong types and unknown fields before reaching the protected Handler. Write-plan drafts **MUST** pass the same reviewed schema before approval-token issuance.
- Adding or editing an Engine operation **MUST** update its callable tool metadata, input schema, protected adapter, and related tests in the same task for every permitted business path. Tests **MUST** enumerate all registered callable tools for metadata/schema completeness and cover representative invalid nested values and unapproved fields, so newly added tools cannot silently skip the gate.
- A business-behavior change **MUST** record affected operation and tool IDs with the same-task tool update, verified-no-change evidence, or concrete no-impact reason even when the GraphQL/HTTP shape is unchanged. The catalog **MUST** reject either side of a callable inventory/tool mismatch; session and identity operations stay non-callable. The task **MUST** test changed authorization, scope, validation, projection, side effects, and errors through the protected Engine path, then run `make check_ai_contracts`.
- Inventory availability accepts only `CALLABLE` and `NON_CALLABLE`. Registered tool permissions must match the independent candidate coverage record and applicable GraphQL `@hasPermission` action after workspace mapping. A source-changing Engine PR must carry a structured `AI_TOOL_IMPACT` declaration; CI validates its shape and reviewers validate its truth.

---

### 4.4 Full Triple Pattern (Custom Query/Mutation)

#### File Structure

```text
src/services/{domain}/
├── {domain}.go              # Interface + Handlers struct + Service struct + Constructor
├── {domain}_handler.go      # Business logic (pure functions)
└── {domain}_service.go      # Service methods (delegates to Handler)
```

#### ① `{domain}.go` — Interface Definition

```go
package {domain}

import (
    "context"
    "{module_path}/gen"
)

// I{Domain}Service defines the service interface
type I{Domain}Service interface {
    MethodA(ctx context.Context, r *gen.GeneratedResolver, input gen.XxxParams) (*gen.XxxData, error)
}

// {Domain}Handlers defines the handler function collection
type {Domain}Handlers struct {
    MethodAHandler func(ctx context.Context, r *gen.GeneratedResolver, input gen.XxxParams) (*gen.XxxData, error)
}

// {Domain}Service is the service implementation
type {Domain}Service struct {
    Handlers {Domain}Handlers
}

// New{Domain}Service creates a new service instance
func New{Domain}Service() I{Domain}Service {
    return &{Domain}Service{
        Handlers: {Domain}Handlers{
            MethodAHandler: MethodAHandler,
        },
    }
}
```

> [!IMPORTANT]
> **Naming Conventions**:
> - Interface: `I{Domain}Service`
> - Handler collection struct: `{Domain}Handlers` (plural, distinct from function name)
> - Service struct: `{Domain}Service`
> - Handler collection field: `Handlers`
> - Function field name = corresponding function name (e.g., `LoginHandler` field maps to `LoginHandler` function)

#### ② `{domain}_handler.go` — Business Logic

For **Query** handlers:
```go
package {domain}

// QueryHandler handles read-only logic (no transaction needed)
func QueryHandler(ctx context.Context, r *gen.GeneratedResolver, input gen.XxxParams) (*gen.XxxData, error) {
    // Read-only business logic
}
```

For **Mutation** handlers (requires transaction management):
```go
package {domain}

// MutationHandler handles write logic with transaction management
func MutationHandler(ctx context.Context, r *gen.GeneratedResolver, input gen.XxxParams) (*gen.XxxData, error) {
    ctx = gen.EnrichContextWithMutations(ctx, r) // Begin transaction

    // Business logic...
    if err != nil {
        gen.RollbackMutationContext(ctx, r) // Rollback
        return nil, err
    }

    err = gen.FinishMutationContext(ctx, r) // Commit
    return result, err
}
```

> [!TIP]
> **Handler functions should be package-level pure functions** for easy testing and flexible replacement.

#### ③ `{domain}_service.go` — Service Proxy

```go
package {domain}

// MethodA delegates to the handler
func (s {Domain}Service) MethodA(ctx context.Context, r *gen.GeneratedResolver, input gen.XxxParams) (*gen.XxxData, error) {
    return s.Handlers.MethodAHandler(ctx, r, input)
}
```

#### ④ Registration (`src/resolver.go`)

```go
// Global Service variable declaration
var {Domain}Service {domain}.I{Domain}Service

// Initialize in New()
func New(...) gen.Config {
    {Domain}Service = {domain}.New{Domain}Service()
    // ...
}
```

#### ⑤ Resolver Entry (`src/extend.go`)

For **Query**:
```go
func (r *QueryResolver) MethodA(ctx context.Context, input gen.XxxParams) (*gen.XxxData, error) {
    return {Domain}Service.MethodA(ctx, r.GeneratedResolver, input)
}
```

For **Mutation**:
```go
func (r *MutationResolver) MethodB(ctx context.Context, input gen.YyyParams) (*gen.YyyData, error) {
    return {Domain}Service.MethodB(ctx, r.GeneratedResolver, input)
}
```

---

### 4.5 Pure Relationship Modeling

纯关系若只承载“实体 A 关联实体 B”，且不承载独立业务事实时，必须使用 engine 自动生成的 many-to-many join table，而不是显式空中间实体。
When a relationship only expresses "entity A is related to entity B" and does **not** carry independent business facts, it **MUST** use the engine-generated many-to-many join table instead of an explicit empty middle entity.

**Pure relationship rule**:

- 只有当关系本身承载独立业务事实时，才允许定义显式中间实体。
- 禁止在业务代码和测试中手写不存在于 schema 契约的列名假设。
- 边界测试必须贴近真实生成结构。
- If a pure relationship only carries "entity A is related to entity B" and has no independent business fields such as `status`, `sort`, `note`, `enabledAt`, `priority`, or `lifecycleState`, it **MUST** use the engine-generated many-to-many join table.
- An explicit middle entity is allowed **only** when the relationship itself carries independent business facts.
- Every endpoint link on an allowed explicit middle entity **MUST** still be modeled with `@relationship`; custom endpoint ID columns are forbidden.
- For engine-generated many-to-many join tables, business code and tests **MUST NOT** hand-write column assumptions that do not exist in the schema contract.
- Any boundary test **MUST** stay close to the real generated structure and **MUST NOT** fake a passing result by adding extra compatibility columns.

**Examples**:

- Good: `CompanyRole.skills <-> Skill.companyRoles`
- Good: `SkillGroup.skills <-> Skill.skillGroups`
- Bad: create `FooBarBinding` only to store two foreign keys and no business fields

> [!IMPORTANT]
> If you discover a middle entity that has become a pure relation over time, prefer simplifying it back to schema-level many-to-many rather than layering more handler logic, CRUD blockers, or fake compatibility fixtures around it.

### 4.6 Simplified Pattern (Override Generated CRUD)

#### File Structure

```text
src/services/{domain}/
└── {domain}_handler.go      # Handler function only
```

#### Handler File

```go
package {domain}

// Create{Domain}Handler creates a {Domain}
func Create{Domain}Handler(ctx context.Context, r *gen.GeneratedResolver, input map[string]interface{}) (*gen.{Domain}, error) {
    // Custom business logic BEFORE default handler...

    item, err := gen.Create{Domain}Handler(ctx, r, input)
    if err != nil {
        return nil, err
    }

    // Custom business logic AFTER default handler...

    return item, nil
}
```

#### Registration (`src/resolver.go`)

```go
resolver.Handlers.Create{Domain} = {domain}.Create{Domain}Handler
```

> [!CAUTION]
> **Transaction Management — Override Handlers MUST NOT call transaction functions!**
> 
> `EnrichContextWithMutations`, `RollbackMutationContext`, `FinishMutationContext` are **already called by the Resolver entry** (`gen/resolver-mutations.go`). Calling them again in the Override Handler is **redundant** and should be avoided.
>
> **Bad example:**
> ```go
> func MyHandler(ctx context.Context, r *gen.GeneratedResolver, input map[string]interface{}) {
>     ctx = gen.EnrichContextWithMutations(ctx, r) // ❌ Redundant! Already called by Resolver entry
>     // ...
>     gen.FinishMutationContext(ctx, r) // ❌ Redundant! Already called by Resolver entry
> }
> ```
>
> **Correct approach:** Handler calls `gen.Create{Domain}Handler(ctx, r, input)` directly and returns `nil, err` on error, letting the Resolver entry handle transaction commit/rollback.

#### Generated CRUD Enum Boundary Rules

*   Custom handlers that forward `map[string]interface{}` inputs that include GraphQL enum fields to generated CRUD handlers **MUST** normalize those enum fields back to GraphQL enum strings first.
*   Any `@entity` with enum fields that exposes generated create/update mutations **MUST** prove that the write path accepts gqlgen-decoded enum values, including optional enum pointers.
*   Use `generatedinput.NormalizeEnumStringFields` for this boundary normalization instead of duplicating one-off enum conversion logic in each handler.
*   This rule applies to generated CRUD handlers such as `gen.Create{Entity}Handler` and `gen.Update{Entity}Handler`, including indirect service flows that call those handlers from custom mutations.
*   Tests for custom write handlers touching enum fields **MUST** verify that generated CRUD receives string enum values, including pointer enum inputs where gqlgen may provide them.
*   Parent write paths that accept nested relationship objects or nested maps containing enum fields **MUST** normalize and verify those nested enum paths too; covering only the top-level entity enum fields is insufficient.
*   Frontend enum string values do not satisfy this backend boundary requirement because gqlgen can decode optional enum inputs before the resolver sees them.
*   Handler tests and contract tests that protect this boundary **MUST** feed gqlgen-decoded runtime shapes (`gen.EnumType`, `*gen.EnumType`, nested relationship structs, nested maps) rather than handwritten string stand-ins when production would not deliver strings at that point.
*   If an enum entity intentionally does not need override normalization, it **MUST** be listed in an explicit contract-test allowlist with a reason.
*   Do not pass already-decoded generated enum values back into generated CRUD. That path can trigger gqlgen / mapstructure failures such as `enums must be strings`.

#### Production-Fidelity Test Boundary Rules

*   Any backend test that claims to protect a GraphQL decode boundary, generated CRUD boundary, GORM persistence path, migration compatibility path, or legacy-schema compatibility path **MUST** preserve the real runtime input shape and the relevant table constraints that can change behavior.
*   SQLite or in-memory fixture schemas used as MySQL surrogates **MUST NOT** silently drop touched production constraints such as `NOT NULL`, default values, uniqueness, required foreign keys, or legacy compatibility columns when those constraints affect the behavior under test.
*   If a production path must remain compatible with a legacy column or legacy constraint, the regression test **MUST** bootstrap that legacy schema shape explicitly instead of assuming the latest clean schema.
*   Handwritten simplified fixtures are allowed only for pure helper/unit tests that do not claim boundary coverage. They **MUST NOT** be the only protection for resolver/handler/service paths that cross gqlgen, generated CRUD, ORM, migration, or persistence boundaries.
*   A test does **not** count as boundary coverage if it replaces gqlgen-decoded typed enums/pointers with plain strings, or if it uses a weakened test table definition that would permit writes rejected by the real database.
*   When a production bug escapes because a fixture was weaker than reality, the fix is incomplete until the suite contains a named regression test that reproduces the real failing boundary first and then proves the repaired path stays green.
*   New or modified test helpers that intentionally diverge from production shape **MUST** document that limitation in comments and **MUST** be paired with at least one higher-fidelity regression or integration test covering the real boundary.

---

### 4.7 General Coding Standards

| Rule | Description |
|------|-------------|
| **File Header Comments** | Every file must include `@Author`, `@Email`, `@Date` annotations |
| **Comment Accuracy** | Function comments must match function name and behavior. Copying stale comments is prohibited |
| **Error Messages** | User-facing errors must use i18n keys or localized strings. System-level errors in English |
| **Password Handling** | Use `utils.EncryptPassword()` uniformly. Storing plaintext passwords is prohibited |
| **Context Usage** | Retrieve values via `config.KeyXxx` from context. Always nil-check before use |
| **Inter-Service Dependencies** | Inject via interface parameters, not direct global variable references |

#### Source Cohesion (Mandatory)

- Run `make check_structure` after changing handwritten Go source, including tests and closures. New files and functions must meet the shared 250/50 effective-line, complexity-10, and nesting-3 limits. Generated code under `gen/` is excluded and immutable.
- A domain handler may be split into multiple cohesive Go files in the same service package while preserving the existing interface, handler registration, transaction owner, and tests. The three-file example above is a starting layout, not a mandate to accumulate every handler in one large file.
- Keep a helper unexported and package-local while it has one consumer. Move it to the nearest shared package only when at least two independent callers need the same semantics. Do not add a generic `utils` wrapper that merely forwards to one call or hides writes, errors, or transaction boundaries.
- Existing measured over-limit units use a recorded non-growing baseline; an exception beyond that baseline requires a concrete reason and review evidence in the plan. Never exempt an entire handwritten directory.

#### Code Comment Rules

*   `base-engine` 中新增或修改的 **类型、常量、结构体字段、接口、函数、方法**，**必须**补充简洁准确的中文注释。
*   注释 **MUST** 说明职责、业务含义、边界、输入输出语义、副作用或存在原因，不能只机械重复代码字面意思。
*   当参数、返回值、状态或布尔语义不直观时，函数或方法的中文注释 **必须**补充说明这些语义，避免调用方靠猜测理解。
*   handler / service / repository 等业务入口函数，注释 **应**优先说明它负责什么业务、会修改什么状态、以及可能产生哪些副作用。
*   测试中的桩对象、辅助函数和关键断言场景，只要引入了新的字段、方法或不直观行为，也 **必须**补充中文注释。
*   当你触达已有未注释声明并对其进行修改时，**必须顺手补齐中文注释**；如果周边声明缺少注释会直接影响理解，也 **必须**一并补齐，不得继续扩大“无注释存量”。

| ❌ Prohibited | ✅ Required |
|--------------|------------|
| Leaving touched declarations without Chinese comments | Add concise Chinese comments for touched declarations before completion |
| Writing comments that only restate the identifier literally | Explain responsibility, boundary, and business meaning in Chinese |
| Omitting non-obvious parameter, return, status, or boolean semantics | Clarify those semantics in the declaration comment |

#### Traceable Logging Rules

*   All newly added or modified backend business flows and infrastructure paths **MUST** integrate Traceable Logging.
*   Engine trace coverage **MUST** stay layered across:
    - request / operation trace
    - business decision logging
    - mutation audit logging
*   Generated entity-change auditing **MUST** continue to use `resolver.Handlers.OnEvent` as the formal mutation audit outlet.
*   Sensitive fields **MUST** use the centralized masking set. Handwritten one-off masking logic scattered across handlers, middleware, or runtime loops is prohibited.
*   Free-form business `log.Printf` style logging is prohibited where structured trace logging is expected.
*   When a task touches middleware, resolver, directive, handler, runtime loop, websocket path, or audit path, that touched chain **MUST** be brought up to the current Traceable Logging standard before the task is complete.

#### SMS E164 Phone Boundary Rules

*   Account identity phone validation and SMS/Twilio phone validation must not share the same validator unless their semantics are intentionally identical.
*   SMS-domain phone fields must use a dedicated `E.164` validator instead of reusing the account-domain `phone` rule.
*   custom SMS mutations must normalize phone input inside handlers before query, persistence, whitelist matching, or delivery.
*   Web or app input UX may remain country-specific for product reasons, but engine SMS storage and comparison must stay canonical at the `E.164` boundary.
*   When compatibility with legacy SMS phone records is required, the compatibility fallback must remain local to the SMS domain and must not weaken the account-domain validator.

#### Timezone Safety Rules

*   Engine business time fields must remain absolute millisecond timestamps unless an explicit contract requires a different absolute epoch format.
*   Engine business identifiers must not encode silent server-local timezone semantics, ambiguous formatted wall-clock date fragments, or locale-dependent calendar text.
*   Readable business identifiers may use explicit UTC formatted business identifiers when their timezone semantics are explicit, reviewable, and stable across deployment regions.
*   When readable UTC prefixes are used, uniqueness must come from node id and sequence or another deterministic anti-collision mechanism rather than probabilistic random suffixes.
*   When human-readable business time strings are unavoidable inside engine code, their timezone semantics must be explicit and reviewable rather than inherited from `time.Local`.
*   New or modified engine helpers, runtime loops, audit outputs, and persistence paths must not rely on hidden server-local timezone assumptions.
*   Prefer timezone-independent identifier strategies or readable UTC prefixes for order numbers or similar business identifiers instead of formatted local wall-clock text.

#### Soft Delete Boundary Rules

*   Engine normal read paths **MUST** exclude soft-deleted rows from entity queries, relationship resolvers, relationship id fields, loaders, and custom aggregate read models unless an explicit recovery or audit contract requires otherwise.
*   If generated relationship resolvers or loaders leak soft-deleted rows, the touched chain **MUST** override that read path instead of relying on write-side soft delete alone.
*   Relationship object fields and their parallel `...Ids` fields **MUST** follow the same soft-delete visibility boundary; they must not disagree about whether a deleted relation still exists.
*   Replace-style relationship writes (for example “replace this entity’s member set”) **MUST** be verified end to end with a reopen/refetch style test so historical soft-deleted relations do not rehydrate back into live UI state.
*   Any handler, resolver, or service that intentionally exposes soft-deleted data **MUST** do so through an explicit, reviewable contract rather than by inheriting generated default behavior accidentally.

---

### 4.8 GraphQL Model Directives

#### `@entity`

Declares a database entity. Dolphin generates full CRUD operations based on this directive.

```graphql
type User @entity(title: "User Management") {
  phone: String! @column(gorm: "type:varchar(32);NOT NULL;") @validator(required: "true", type: "phone")
}
```

#### `@column`

Defines database column properties using GORM tag syntax:

```graphql
phone: String! @column(gorm: "type:varchar(32) comment 'Phone number';NOT NULL;index:phone;")
```

#### `@validator`

Field validation directive, executed in the `Directives.Validator` handler. **Complete parameter list**:

| Parameter | Type | Description | Example |
|-----------|------|-------------|--------|
| `required` | String | Required field | `required: "true"` |
| `immutable` | String | Cannot be modified after creation | `immutable: "true"` |
| `type` | String | Regex validation type (defined in `utils/rule.go`) | `type: "phone"` |
| `minLength` | Int | Minimum string length | `minLength: 6` |
| `maxLength` | Int | Maximum string length | `maxLength: 32` |
| `minValue` | Int | Minimum numeric value | `minValue: 0` |
| `maxValue` | Int | Maximum numeric value | `maxValue: 100` |
| `unique` | String | Uniqueness constraint (no duplicates) | `unique: "true"` |
| `uniqueScope` | String | Conditional uniqueness scope field (optional) | `uniqueScope: "uid"` |

**Built-in validation types** (`type` parameter values):

| Type | Description | Special Behavior |
|------|-------------|------------------|
| `phone` | Phone number format | — |
| `email` | Email format | — |
| `password` | Password format | **Auto-encrypts** via `utils.EncryptPassword` |
| `int` | Integer format | — |
| `justInt` | Numeric-only format | — |

> [!IMPORTANT]
> `type: "password"` not only validates format, but also **auto-encrypts the plaintext password** in the `Directives.Validator` handler in `resolver.go`. This is the single entry point for password encryption.

#### Uniqueness Validation

```graphql
# Table-wide unique
phone: String! @validator(required: "true", unique: "true")

# Scoped unique: name must be unique within the same uid
name: String! @validator(required: "true", unique: "true", uniqueScope: "uid")
```

Uniqueness checks automatically exclude soft-deleted records (`is_delete = 1`) and depend on `db` and `tableName` injected into the context.

**Schema-First Uniqueness Rule**:
- When a field or field-scope pair carries real business uniqueness truth, engineers **SHOULD** express that contract in `model/*.graphql` first with `@validator(unique: "true")`, adding `uniqueScope` when the uniqueness is intentionally scoped.
- Do **not** leave uniqueness enforcement as frontend-only validation, web-form warnings, or scattered service-layer duplicate checks when the schema validator can express the rule directly.
- `@validator(unique: "true")` is appropriate for stable business identifiers and governance fields such as names, actions, phone numbers, singleton keys, or other values whose uniqueness should remain true across normal CRUD paths.
- `@validator(unique: "true")` **MUST NOT** be used for:
  - runtime / history / ledger fields that are expected to repeat
  - computed or mirrored fields that are not the business source of truth
  - uniqueness rules that depend on many-to-many membership, join-table state, or multi-row aggregate conditions
  - delete / recovery / governance semantics that require a dedicated service-layer decision instead of a simple validator lookup
- Adding or changing `@validator(unique: "true")` is **incomplete** unless regression coverage proves both:
  - the schema keeps the expected validator snippet
  - the real write path rejects duplicates using production-fidelity payload shapes and DB fixtures that preserve the relevant nullability/default/soft-delete constraints

#### `@hasRole`

Access control directive, restricts Query/Mutation access by role:

```graphql
type Admin @entity(title: "Admin") @hasRole(role: ADMIN) { ... }
```

Available roles: `ALL` (any authenticated user), `ADMIN` (administrator), `USER` (regular user), `GUEST` (guest).

#### `@hasPermission`

Fine-grained permission directive, restricts specific fields/mutations by permission action string:

```graphql
type Mutation {
  deleteUsers(id: [ID!]!): Boolean! @hasPermission(action: "user:delete")
}
```

The interceptor extracts the `permissions` array from the JWT token payload and checks whether the required `action` string is present. If absent, returns a "Permission denied" error.

#### `@relationship`

Defines entity associations. Dolphin auto-generates a GORM `many2many` join table — **do NOT create manual junction table entities**.

**Syntax**:

```graphql
# Many-to-Many (use [] brackets on both sides)
permissions: [Permission!] @relationship(inverse:"roles")
roles:       [SysRole!]    @relationship(inverse:"permissions")

# Many-to-One / One-to-Many (use [] on the "many" side only)
author:  User!   @relationship(inverse:"posts")
posts:   [Post!] @relationship(inverse:"author")
```

**Rules**:
- The `inverse` parameter **MUST** match the relationship field name on the other entity.
- For many-to-many, the framework auto-creates a join table (e.g., `permission_roles`).
- In Go code, use GORM `Preload("FieldName")` to load associations and `Association("FieldName").Append(...)` to create them.
- **Never** manually create a junction `@entity` for relationships — always use `@relationship`.
- **Never** declare scalar `xxxId` / `xxxIds` fields in a handwritten `@entity` to represent a relationship. If Dolphin emits parallel relationship ID fields or FK columns, those remain generated implementation details.

**Delete-Protection Rule (`master:"yes"`)**:
- When a delete-protection rule is a simple **FK-backed parent -> children** constraint, the parent-side `[]Child` relationship **SHOULD** express it in schema first:
  ```graphql
  systemPromptSkills: [Skill!] @relationship(inverse:"systemPromptTemplate", master:"yes")
  ```
- `master:"yes"` is only appropriate when Dolphin can enforce the guard by checking the child table's `<inverse>_id` column during generated delete flow.
- Before writing a hand-maintained delete guard for a one-to-many / one-to-one inverse relationship, engineers **MUST** first evaluate whether `master:"yes"` can express the rule directly in `model/*.graphql`.
- `master:"yes"` **MUST NOT** be used for:
  - many-to-many or other join-table-backed relationships
  - delete paths that require cascade cleanup instead of simple blocking
  - governance flows with additional business rules beyond "child row exists"
  - runtime / history / ledger entities where historical records would permanently freeze legitimate delete operations
- Adding or changing `master:"yes"` is **incomplete** unless regression coverage proves both:
  - the schema keeps the expected `@relationship(..., master:"yes")` snippet
  - the external generator owner supplies the expected delete guard in `gen/resolver-mutations.go`; AI/agents verify it read-only

### 4.9 One-to-One Bidirectional Binding

Dolphin's one-to-one relationship must first be declared on both entity types with `@relationship` and matching inverses. Dolphin may then generate FK columns and mutation keys for both directions (for example `users.setup_key_id`, `netbird_setup_keys.user_id`, `setupKeyId`, and `userId`). Those names are generator-owned implementation details, not fields that may be handwritten in `model/*.graphql`.

When the generated contract requires both directions to be bound, use its generated mutation keys without duplicating them as custom schema fields:

```go
// 1. Create child record (set forward FK)
setupKeyItem, err := gen.CreateNetbirdSetupKeyHandler(ctx, r, map[string]interface{}{
    "userId": item.ID,  // Forward FK
    // ...
})

// 2. Update parent record (set reverse FK)
_, err = gen.UpdateUserHandler(ctx, r, item.ID, map[string]interface{}{
    "setupKeyId": setupKeyItem.ID,  // Reverse FK
})
```

> [!WARNING]
> Setting only one generated side can cause the reverse relationship query to return `null`. This runtime requirement never permits a handwritten `userId`, `setupKeyId`, or other scalar foreign-key field inside an `@entity`; the schema relationship itself must remain `@relationship`-based.

### 4.10 Permission Activation Workflow

After registering new `@entity` permissions in `InitPermissions()`, the **full activation workflow** is required for them to take effect on the frontend:

1. **`make migrate`** — Writes permission records to the `permissions` table and binds them to the ADMIN role
2. **Restart base-engine** — `make start`
3. **Re-login to the admin panel** — Obtain a new JWT token containing the new permissions

> [!IMPORTANT]
> The frontend sidebar uses `hasPermission()` to check the permissions list in the JWT token. **No re-login = no new permissions in token = menu items invisible**.

---

### 4.11 New Custom API Checklist

Example: Adding a new custom Query `getFoo`:

- [ ] **`model/extend.graphql`** — Define input types, return types, and Query declaration
- [ ] **Generation authorization** — Run `make generate` after schema edits only with task-scoped user authorization; otherwise obtain externally generated output before continuing
- [ ] **`src/services/foo/foo.go`** — Create interface + Handlers + Service + constructor
- [ ] **`src/services/foo/foo_handler.go`** — Implement business logic
- [ ] **`src/services/foo/foo_service.go`** — Implement Service proxy methods
- [ ] **`src/resolver.go`** — Declare global variable and initialize in `New()`
- [ ] **`src/extend.go`** — Add Resolver method, delegate to Service

### 4.12 New Entity Checklist

Every time a new `@entity` is added to `model/model.graphql`, the following steps **MUST** be completed:

- [ ] **`model/model.graphql`** — Define entity with `@entity`, `@column`, `@validator` directives
- [ ] **Relationship audit** — Model every cross-entity link with typed fields plus `@relationship`; verify that the entity contains no custom scalar `xxxId` / `xxxIds` or alias used as a foreign key
- [ ] **Schema contract test** — Add or update a fail-closed check that rejects handwritten scalar cross-entity links, including links on runtime/history/audit/governance/association entities
- [ ] **Generation authorization** — Run `make generate` after schema edits only with task-scoped user authorization; otherwise obtain externally generated output before continuing
- [ ] **`make migrate`** — Create the database table
- [ ] **`src/dbup/init.go`** — Register CRUD permissions for the new entity in the correct permission seed bucket

> [!IMPORTANT]
> **Permission Registration Rule**: For each new entity `{Entity}`, add the following four permission entries to the correct seed bucket in `src/dbup/init.go`:
> ```go
> {Name: "View {Entity}", Action: "{entity}:read",   Module: "{entity}"},
> {Name: "Create {Entity}", Action: "{entity}:create", Module: "{entity}"},
> {Name: "Edit {Entity}", Action: "{entity}:update", Module: "{entity}"},
> {Name: "Delete {Entity}", Action: "{entity}:delete", Module: "{entity}"},
> ```
> Use `systemPermissionSeeds()` for platform governance capabilities, and use `tenantBusinessPermissionSeeds()` for permissions that should be assignable to `CompanyRole`.
> After adding permissions, run `make migrate` again to seed them into the database and auto-bind to the ADMIN role.

## 5. Development Workflow

### 5.1 Manual Review Commit Protocol

- After completing implementation and verification in `base-engine`, the agent must leave all modified files in the local working tree and keep them unstaged for human review.
- The agent must not automatically run `git add`.
- The agent must not automatically run `git commit`.
- The standard handoff summary must include:
  - which files changed
  - which verification commands ran and what passed
  - the current working-tree status
  - a reminder that the developer performs review, staging, and committing manually
- If unrelated local changes already exist in the repository, the agent must not attempt any staging workaround or selective early staging.

0.  **Discover**: Before starting, read `base-engine/go.mod` for the module path, and scan `base-engine/src/services/` for existing service patterns to follow as reference.
1.  **Define**: Before changing GraphQL definitions in `model/`, determine whether `gen/` must change and confirm task-scoped user authorization for generation. Without it, stop and report the dependency; do not leave an incompatible schema edit.
2.  **Generated contract**: With authorization, run `make generate` after handwritten schema changes and inspect the complete diff. Without authorization, wait for externally generated output. Never edit generated files directly.
3.  **Implement**: Implement business logic in `src/services/` using the Triple or Simplified pattern.
4.  **Type Safety**: Leverage Go's strong typing to ensure type-safe conversion between Database Model (MySQL) and GraphQL Model.
5.  **Verify**: Run `make check_structure`, `go build ./...`, and `go vet ./...` before marking a Go source task complete.

### Error Recovery

If verification fails, follow this troubleshooting sequence:

| Error | Resolution |
|-------|-----------|
| Generation is not authorized or unavailable | Stop before schema changes that require `gen/` updates; report the dependency without modifying `gen/` |
| `go build` fails with import errors | 1) Verify `{module_path}` matches `go.mod` 2) Run `go mod tidy` 3) Check for circular imports |
| `go vet` reports issues | Fix all reported issues — these are potential bugs, not style warnings |

**Version**: 2.2.0 | **Date**: 2026-09-25
