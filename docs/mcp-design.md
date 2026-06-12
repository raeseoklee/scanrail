[ENGLISH](mcp-design.md) | [한국어](mcp-design.ko.md)

# MCP Design

## Recommendation

Build a Scanrail MCP server, but keep it as a thin local adapter over the existing CLI safety model. MCP should help AI assistants inspect configuration, run bounded checks, and summarize reports. It should not become a second execution engine or a way to bypass Scanrail policy.

## Why MCP Fits

The Model Context Protocol is a standard way for AI applications to connect to external tools, resources, and workflows. That maps well to Scanrail:

- Tools can expose bounded actions such as `doctor`, `setup`, and safe scan runs.
- Resources can expose configuration, schemas, and report summaries.
- Prompts can guide finding triage and remediation planning.

References:

- [Model Context Protocol introduction](https://modelcontextprotocol.io/docs/getting-started/intro)
- [MCP security best practices](https://modelcontextprotocol.io/docs/tutorials/security/security_best_practices)
- [Official TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)

## Scope

First MCP release:

- local stdio server only
- no remote HTTP server
- no authentication proxy
- no arbitrary shell command execution
- no active scans unless the project policy explicitly allows them
- no secret values in tool inputs, outputs, resources, logs, or prompts

## 0.1.2 MVP Status

Implemented in `scanrail mcp serve`:

- stdio JSON-RPC server
- `initialize` and `ping`
- `tools/list` and `tools/call`
- `resources/list` and `resources/read`
- `scanrail_doctor`
- `scanrail_config_read`
- `scanrail_report_latest`
- `scanrail_run` for the native headers scanner only
- `scanrail://config`
- `scanrail://reports/latest/summary`
- `scanrail://safety-model`

The implementation intentionally avoids a new dependency and follows the MCP stdio transport requirement that stdout contains only JSON-RPC messages.

## Proposed Tools

| Tool | Purpose | Safety notes |
| --- | --- | --- |
| `scanrail_doctor` | Return environment readiness. | Read-only except normal CLI checks. |
| `scanrail_config_read` | Return normalized project config and validation warnings. | Redact secret references to names only. |
| `scanrail_setup` | Prepare local Scanrail state. | Planned after setup safety logging is added. |
| `scanrail_run` | Run a bounded scan profile. | Implemented for native headers only; enforces configured target allowlists and active-scan opt-in. |
| `scanrail_report_latest` | Return latest report metadata and summary. | Size-limit output and omit raw secrets. |
| `scanrail_findings_explain` | Convert findings into remediation-oriented text. | Uses report data only; no new scan execution. |

## Proposed Resources

| Resource | Description |
| --- | --- |
| `scanrail://config` | Current project configuration with secret values redacted. |
| `scanrail://schema/config` | Config schema and supported fields. Planned. |
| `scanrail://reports/latest/summary` | Latest report summary. |
| `scanrail://reports/latest/json` | Latest report JSON, size-limited. Planned. |
| `scanrail://safety-model` | Effective safety policy and allowlist state. |

## Packaging

Preferred first implementation:

```bash
scanrail mcp serve
```

This keeps the MCP server inside the Go CLI release artifact and avoids introducing a new runtime requirement. If the MCP TypeScript SDK remains the most mature production surface when implementation starts, a later `@scanrail/mcp` npm package can wrap the same CLI capabilities.

Example client command:

```json
{
  "command": "scanrail",
  "args": ["mcp", "serve"]
}
```

## Security Rules

- Use stdio transport first to avoid exposing a listening local HTTP endpoint.
- Validate workspace root before reading files.
- Validate every target URL against Scanrail config before execution.
- Never accept arbitrary scanner commands from MCP input.
- Never pass tokens through the MCP client; credentials remain environment-owned.
- Redact secret values and limit large report payloads.
- Treat MCP tool descriptions as user-facing documentation, not as authorization.
- Record tool calls in Scanrail logs when they execute scans or setup.

## Implementation Plan

1. Add MCP client configuration examples for common hosts.
2. Add `scanrail_setup` only after setup writes execution logs.
3. Add bounded latest report JSON with size limits.
4. Add prompts for finding triage and remediation planning.
5. Add compatibility tests against at least one production MCP client.

## Non-Goals

- Hosted SaaS MCP endpoint
- OAuth proxying to third-party scanners
- Full scanner management through natural language
- Agent-controlled installation of Docker images without explicit user approval
