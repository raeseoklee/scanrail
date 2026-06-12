[ENGLISH](mcp-design.md) | [한국어](mcp-design.ko.md)

# MCP 설계

## 권장안

Scanrail MCP server를 만드는 것은 좋습니다. 다만 기존 CLI safety model 위의 얇은 local adapter로 유지해야 합니다. MCP는 AI assistant가 configuration을 확인하고, 제한된 check를 실행하고, report를 요약하도록 돕는 역할이어야 합니다. 두 번째 execution engine이 되거나 Scanrail policy를 우회하는 경로가 되면 안 됩니다.

## MCP가 맞는 이유

Model Context Protocol은 AI application이 외부 tool, resource, workflow에 연결되는 표준 방식입니다. Scanrail과는 다음처럼 잘 맞습니다.

- Tool은 `doctor`, `setup`, safe scan run 같은 제한된 action을 노출할 수 있습니다.
- Resource는 configuration, schema, report summary를 노출할 수 있습니다.
- Prompt는 finding triage와 remediation planning을 안내할 수 있습니다.

참고:

- [Model Context Protocol introduction](https://modelcontextprotocol.io/docs/getting-started/intro)
- [MCP security best practices](https://modelcontextprotocol.io/docs/tutorials/security/security_best_practices)
- [Official TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)

## 범위

첫 MCP release:

- local stdio server only
- remote HTTP server 없음
- authentication proxy 없음
- arbitrary shell command execution 없음
- project policy가 명시적으로 허용하지 않는 active scan 없음
- tool input, output, resource, log, prompt에 secret value 없음

## 제안 Tool

| Tool | 목적 | Safety note |
| --- | --- | --- |
| `scanrail_doctor` | environment readiness 반환 | 일반 CLI check 외에는 read-only |
| `scanrail_config_read` | normalized project config와 validation warning 반환 | secret reference는 name만 노출 |
| `scanrail_setup` | local Scanrail state 준비 | 기본은 `--pull-policy never`; image pull은 explicit input 필요 |
| `scanrail_run` | 제한된 scan profile 실행 | configured target allowlist와 active-scan opt-in 강제 |
| `scanrail_report_latest` | latest report metadata와 summary 반환 | output size 제한, raw secret 제외 |
| `scanrail_findings_explain` | finding을 remediation 중심 설명으로 변환 | report data만 사용, 새 scan 실행 없음 |

## 제안 Resource

| Resource | 설명 |
| --- | --- |
| `scanrail://config` | secret value가 redacted된 현재 project configuration |
| `scanrail://schema/config` | config schema와 지원 field |
| `scanrail://reports/latest/summary` | latest report summary |
| `scanrail://reports/latest/json` | size-limited latest report JSON |
| `scanrail://safety-model` | effective safety policy와 allowlist 상태 |

## Packaging

첫 구현 권장 형태:

```bash
scanrail mcp serve
```

이 방식은 MCP server를 Go CLI release artifact 안에 유지하므로 새 runtime requirement를 늘리지 않습니다. 구현 시점에도 MCP TypeScript SDK가 가장 성숙한 production surface라면, 이후 `@scanrail/mcp` npm package가 같은 CLI capability를 감싸는 형태를 추가할 수 있습니다.

## Security Rule

- listening local HTTP endpoint를 피하기 위해 stdio transport부터 사용합니다.
- file read 전 workspace root를 검증합니다.
- 모든 target URL을 실행 전 Scanrail config와 대조합니다.
- MCP input에서 arbitrary scanner command를 받지 않습니다.
- MCP client를 통한 token passthrough를 금지하고, credential은 environment 소유로 둡니다.
- secret value를 redact하고 큰 report payload를 제한합니다.
- MCP tool description은 user-facing documentation일 뿐 authorization으로 취급하지 않습니다.
- scan 또는 setup을 실행하는 tool call은 Scanrail log에 남깁니다.

## 구현 계획

1. `scanrail mcp serve`를 추가하고 tool/resource registration만 먼저 구현합니다.
2. config, safety policy, latest report summary용 read-only resource를 노출합니다.
3. `scanrail_doctor`와 `scanrail_report_latest`를 추가합니다.
4. 이미 구현된 native headers scanner에 한해 `scanrail_run`을 추가합니다.
5. JSON-RPC fixture 기반 MCP integration test를 추가합니다.
6. server가 안정화된 뒤 client configuration example을 문서화합니다.

## 비목표

- Hosted SaaS MCP endpoint
- third-party scanner에 대한 OAuth proxy
- natural language 기반 전체 scanner management
- explicit user approval 없는 Docker image agent 설치
