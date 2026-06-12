[ENGLISH](README.md) | [한국어](README.ko.md)

# MCP Verification

This scenario verifies the published npm package through the stdio MCP server.

It checks:

- `initialize`
- `tools/list`
- `resources/list`
- `resources/read` for `scanrail://config`
- `scanrail_run` denial without `confirm_active_scan`
- `scanrail_run` success with `confirm_active_scan=true`
- `scanrail_report_latest`

Run:

```bash
node examples/mcp-verification/run.mjs
```

Record:

```bash
vhs examples/mcp-verification/tapes/mcp-verification.tape
```
