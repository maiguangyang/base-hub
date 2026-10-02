<!--
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-03-09 08:55:25
-->
# Figma MCP Config Reference

This document describes how to register the Figma MCP server across popular AI coding assistants and clients (Claude Desktop, Cursor, Gemini/Antigravity, Codex, etc.). Because MCP (Model Context Protocol) is an open standard, only minor format adjustments are needed to work across different clients.

## 1. Common MCP Clients (Claude Desktop / Cursor / etc.)

Most standard MCP-compatible clients (Claude Desktop, Cursor, VSCode-based Roo Code) use **JSON-formatted** config files.
- **Claude Desktop**: `~/Library/Application Support/Claude/claude_desktop_config.json` (Mac)
- **Cursor**: Configure via the UI panel, or edit `.cursor/mcp.json`

### SSE (Streamable HTTP) JSON Config

To connect to the remote Figma MCP endpoint:

```json
{
  "mcpServers": {
    "figma": {
      "type": "sse",
      "url": "https://mcp.figma.com/mcp",
      "env": {
        "FIGMA_OAUTH_TOKEN": "<YOUR_FIGMA_TOKEN>"
      }
    }
  }
}
```
*Note: Support for custom HTTP headers (e.g., `X-Figma-Region`) varies by client — consult each client's docs. In Cursor you can also enter the URL and token directly in the settings UI.*

## 2. Gemini / Antigravity

If Gemini or Antigravity supports MCP server registration, you can register the endpoint the same way. JSON-style declarative configs typically follow the same `mcpServers` schema shown above. Ensure `FIGMA_OAUTH_TOKEN` is available as an environment variable in the host process.

## 3. Codex (TOML Format)

In `~/.codex/config.toml`, register the server as a streamable HTTP service with bearer auth:

```toml
[mcp_servers.figma]
url = "https://mcp.figma.com/mcp"
bearer_token_env_var = "FIGMA_OAUTH_TOKEN"
http_headers = { "X-Figma-Region" = "us-east-1" }
```

- Enable the underlying transport feature at the top level: `[features].rmcp_client = true`.
- Optional timeouts: `startup_timeout_sec` (default 10) and `tool_timeout_sec` (default 60).

## 4. Environment Variable Setup (FIGMA_OAUTH_TOKEN)

Regardless of which AI assistant you use, the host machine must provide the token via an environment variable:

- **One-time (current shell)**: `export FIGMA_OAUTH_TOKEN="<token>"`
- **Persistent**: Add the export line to your shell profile (e.g., `~/.zshrc` or `~/.bashrc`), then **restart your terminal or IDE**.
- **Verify before launch**: Run `echo $FIGMA_OAUTH_TOKEN` — it should print a non-empty value.

## 5. Checklist & Troubleshooting

- **Restart required**: After changing config files or env vars, you **must restart the client** (IDE or standalone app) to reload auth state.
- **Token validity**: If you get OAuth errors, make sure the copied token has no extra quotes or semicolons.
- **Region sync**: If your Figma org uses a specific region, keep the `X-Figma-Region` header consistent across config and requests.
- **Connectivity test**: After restarting, ask the model to "list your available Figma tools." If it can, the server is reachable.

## 6. Interaction Tips

- **How it works**: Most clients cannot actually browse web pages. The correct approach is to **copy the Figma frame/layer link and send it to the model**. The built-in MCP parser extracts the `Node ID` from the URL and fetches data via the server.
- **Deeper analysis**: If the AI output feels too generic, prompt it to follow the required flow: *"Please follow the `get_design_context` → `get_metadata` → `get_screenshot` workflow to re-analyze this Figma node."*
