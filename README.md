# Base Hub (企业级中后台与 AI-Native 基础通用底座)

本项目从生产级高可靠系统中提炼而成，提供了一套开箱即用的**现代多租户企业级中后台 + Go 微内核 + AI 智能助理**全栈基础脚手架。

## 🌟 核心特性

- **前端架构 (base-web)**：
  - **Astro 7 + React 19 + Tailwind CSS v4 + Radix UI** 工业级全套原子组件库。
  - **多标签页 Keep-Alive 栈**：多标签浏览、持久化与恢复、表单草稿不丢失、单标签错误隔离。
  - **纯权限驱动动态侧边栏**：后端返回的权限集合驱动菜单显示，无越权风险。
  - **认证闭环**：系统初始化向导、登录、强制首次修改临时密码、多工作区切换、真正的防假退出。
  - **WebSocket 会话重验与强退**：账号禁用、改密或权限变更秒级强退，杜绝幽灵会话。
  - **核心系统页面**：
    - 管理员管理 (AdministratorListPage)：账号开户、角色分配、启停、临时密码复制、防自杀保护。
    - 角色与权限 (RoleListPage + RolePermissionSelector)：树状/模块化细粒度权限勾选矩阵。
    - 操作审计 (AuditLogListPage)：记录操作人、IP、模块、动作与前后快照 Diff。
    - AI 模型配置 (ModelConfigPage)：脱敏存储 Base URL / API Key，内置探针连接测试。
    - 全局 AI 智能助理抽屉 (AiDrawer)：内置 SSE 流式对话、工具调用与授权审批流程。

- **后端架构 (base-engine)**：
  - **Go + GORM + gqlgen + ADK-Go** 微内核。
  - **双数据库驱动支持**：开发环境默认 **SQLite**（零配置一键运行），生产环境一键切 **MySQL**。
  - **端到端类型安全**：基于 GraphQL Schema-First，自动生成前后端强类型。
  - **RBAC 细粒度拦截**：GraphQL Directives (`@hasRole`, `@hasPermission`) 与租户上下文自动注入。

## 🚀 快速启动

### 1. 启动后端 (base-engine)

```bash
cd base-engine

# 1. 运行数据库迁移 (创建 SQLite 数据库)
go run . migrate

# 2. 初始化初始超级管理员账号
go run . bootstrap-admin

# 3. 启动 API 服务 (默认端口 8085)
go run . start --cors
```

GraphQL Playground 入口: `http://localhost:8085/graphql`

### 2. 启动前端 (base-web)

```bash
cd base-web

# 1. 安装依赖
pnpm install

# 2. 启动本地开发服务
pnpm dev
```

打开浏览器访问: `http://localhost:4321/admin`

---
