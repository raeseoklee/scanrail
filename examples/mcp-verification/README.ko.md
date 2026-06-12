[ENGLISH](README.md) | [한국어](README.ko.md)

# MCP Verification

이 시나리오는 npm에 배포된 package를 stdio MCP server 경로로 검증합니다.

검증 항목:

- `initialize`
- `tools/list`
- `resources/list`
- `scanrail://config`에 대한 `resources/read`
- `confirm_active_scan` 없는 `scanrail_run` 거부
- `confirm_active_scan=true`를 넣은 `scanrail_run` 성공
- `scanrail_report_latest`

실행:

```bash
node examples/mcp-verification/run.mjs
```

녹화:

```bash
vhs examples/mcp-verification/tapes/mcp-verification.tape
```
