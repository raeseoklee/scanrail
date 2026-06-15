[ENGLISH](README.md) | [한국어](README.ko.md)

# MCP Workbench Verification

이 fixture는 `mcp-workbench`로 `scanrail mcp serve`를 검증합니다.

local HTTP target을 시작하고, temporary `scanrail.yaml`을 생성하고, 현재 Scanrail binary를 build한 뒤 stdio MCP server로 노출합니다. test spec은 discovery, resource, safety gating, confirmed native headers scan, audit logging status, latest report summary retrieval을 검증합니다.

실행:

```bash
mcp-workbench run examples/mcp-workbench/scanrail-mcp.yaml --verbose
```

inspect만 실행:

```bash
mcp-workbench inspect --command node --args "examples/mcp-workbench/serve-fixture.mjs" --json
```
