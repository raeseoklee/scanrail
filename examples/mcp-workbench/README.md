[ENGLISH](README.md) | [한국어](README.ko.md)

# MCP Workbench Verification

This fixture verifies `scanrail mcp serve` with `mcp-workbench`.

It starts a local HTTP target, writes a temporary `scanrail.yaml`, builds the current Scanrail binary, and exposes the MCP server over stdio. The test spec validates discovery, resources, safety gating, a confirmed native headers scan, and latest report summary retrieval.

Run:

```bash
mcp-workbench run examples/mcp-workbench/scanrail-mcp.yaml --verbose
```

Inspect only:

```bash
mcp-workbench inspect --command node --args "examples/mcp-workbench/serve-fixture.mjs" --json
```
