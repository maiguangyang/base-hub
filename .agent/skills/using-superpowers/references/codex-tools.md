| Skill references | Codex equivalent |
| --- | --- |
| `Read` | `exec_command` with `sed`, `cat`, or `rg` |
| `Edit` | `apply_patch` |
| `Bash` | `exec_command` |
| `TodoWrite` | update the resolved active-window tracker |
| `Task` / subagent | stay local by default; use `spawn_agent` only when explicitly allowed |

## Environment Detection

Skills that create worktrees or finish branches should detect their environment with read-only git commands before proceeding:

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
BRANCH=$(git branch --show-current)
```

- `GIT_DIR != GIT_COMMON` → already in a linked worktree unless `git rev-parse --show-superproject-working-tree` shows you are inside a submodule
- `BRANCH` empty → detached HEAD or externally managed workspace; avoid branch cleanup assumptions

Do not stage, unstage, commit, or rewrite git state unless the active repo instructions and user explicitly allow it.
