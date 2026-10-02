# `.agent/memory` 速查

## 这是什么

`.agent/memory` 是规则入口层，不是业务记忆层。

它负责回答三件事：

- 先读什么
- 按什么范围读
- 什么时候再展开读完整规则

## 默认读法

1. 先读 [constitution-core.md](./constitution-core.md)
2. 再按任务范围读 `rules-core/` 里的对应文件
3. 再读 `docs/memory/index.yaml`
4. 再进入对应的 `docs/modules/<module>.md`
5. 只有核心规则不够时，才回到 [constitution.md](./constitution.md) 或 `rules/` 里的完整附录

规则句：

- `constitution-core.md` 顶部提到 `constitution.md`，不等于默认一起加载
- 默认策略是 `先读 core，必要时才回 full`

## 目录作用

| 路径 | 作用 | 什么时候读 |
| --- | --- | --- |
| `.agent/memory/constitution-core.md` | 全仓库最小执行入口，定义默认读序和底线规则 | 每次新任务先读 |
| `.agent/memory/constitution.md` | 全仓库完整宪法，提供完整原则、治理规则和背景 | `constitution-core.md` 不够时再读 |
| `.agent/memory/rules-core/` | 各技术栈核心入口，只保留高频、强约束规则 | 已判断任务属于 backend / web / app / ADK-Go 后立即读 |
| `.agent/memory/rules-core/base-engine.md` | backend 核心规则入口 | 任务涉及 `base-engine`、GraphQL、数据库、服务层时 |
| `.agent/memory/rules-core/base-web.md` | web 核心规则入口 | 任务涉及 `base-web`、Astro/React、前端 GraphQL 接入时 |
| `.agent/memory/rules-core/base-app.md` | Flutter app 核心规则入口 | 任务涉及 `base-app`、Riverpod、l10n、App 工作流时 |
| `.agent/memory/rules-core/adk-go.md` | ADK-Go 核心规则入口 | 任务涉及 `adk-go`、Google ADK、agent runtime、MCP Toolset、Hermes-like ADK 方案时 |
| `.agent/memory/rules/` | 各技术栈完整附录，放细则、命令、边界说明和示例 | `rules-core/` 不够时再进入 |
| `.agent/memory/rules/base-engine.md` | backend 完整附录 | 需要细查 backend 命令、架构、边界规则时 |
| `.agent/memory/rules/base-web.md` | web 完整附录 | 需要细查前端架构、页面组织、注释标准、GraphQL 约束时 |
| `.agent/memory/rules/base-app.md` | app 完整附录 | 需要细查 Flutter 分层、命令、测试、缓存、平台细节时 |
| `.agent/memory/rules/adk-go.md` | ADK-Go 完整附录 | 需要细查 ADK-Go 源码入口、官方文档入口、Hermes-like runtime 采用边界时 |

## 和其他目录的边界

- `.agent/memory`：规则入口、治理原则、读取顺序
- `docs/memory/index.yaml`：把任务路由到真正相关的业务模块
- `docs/modules/`：模块当前事实，是业务记忆稳定入口
- `docs/modules-details/`：专题深挖
- `docs/modules-archive/`：历史演进记录
- `docs/plans/task.md` 或 worktree-local tracker：活跃任务追踪，不属于长期规则记忆

## 什么时候回到完整附录

出现下面这些情况时，再从 core 进入 full：

- 当前问题涉及复杂命令或代码生成链路
- 需要查看详细目录架构或实现边界
- 需要确认某条规则的背景、例外或细化解释
- 需要修改规则本身，而不是实现业务代码

## 使用原则

- 先读最小必要规则，不要默认通读整套长文档
- 先按技术栈进入 `rules-core/`，再决定要不要进 `rules/`
- 先走模块路由，再读业务模块 memory
- 把 `.agent/memory` 当规则入口，不要把它当业务事实仓库
