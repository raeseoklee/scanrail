[ENGLISH](safety-model.md) | [한국어](safety-model.ko.md)

# 안전 모델

`scanrail`은 개발자 셀프체크 도구이지만, 일부 스캐너는 실제 공격성 트래픽을 발생시킬 수 있습니다. 따라서 안전 모델은 제품의 핵심 기능입니다.

## 기본 원칙

- 허가된 대상만 스캔합니다.
- 운영 환경 active scan은 기본 차단합니다.
- allowlist 밖 요청을 막을 수 없는 scanner는 해당 보장을 제공한다고 표시하지 않습니다.
- destructive endpoint 차단을 강제할 수 없는 scanner는 interactive/active profile에서 제외하거나 명시 skip합니다.
- secret을 설정 파일이나 리포트에 저장하지 않습니다.
- scanner 요청량 제한은 scanner adapter capability로 검증합니다.
- 실행 결과와 사용한 도구 버전을 재현 가능하게 남깁니다.

## 강제 모델

v0.1 안전 모델은 scanner capability 계약 기반입니다. `scanrail`은 컨테이너 안의 ZAP, Nuclei, Schemathesis가 런타임에 수행하는 모든 redirect, method, path, RPS를 네트워크 레벨에서 강제 차단하지 않습니다.

따라서 각 scanner adapter는 다음 capability를 선언해야 합니다.

```text
allowlist_scope
redirect_scope
blocked_paths
blocked_methods
rate_limit
header_injection
auth_injection
```

프로파일이 요구하는 안전 capability를 scanner가 충족하지 못하면 오케스트레이터는 다음 중 하나를 선택해야 합니다.

- scanner를 `SKIP` 처리하고 리포트에 이유를 남깁니다.
- 사용자가 명시 실행한 scanner라면 안전 정책 위반으로 실패합니다.
- 구현 상태가 불명확한 scanner는 기본 profile에 포함하지 않습니다.

현재 MVP 강제 상태:

- scanner capability metadata는 `internal/scanners`에 선언합니다.
- Gitleaks는 첫 production-ready Docker 기반 passive scanner이며 workspace를 read-only로 mount합니다.
- Trivy와 Semgrep은 아직 production-ready가 아니므로 profile 실행에서 skip합니다.
- production-ready가 아닌 scanner를 명시 실행하면 safety exit code `5`로 실패합니다.
- native headers scanner는 redirect를 따르지 않고 `allowlist_scope`, `redirect_scope`, `rate_limit`, `header_injection`을 선언합니다.

네트워크 레벨 강제 모델은 v0.x 후속 과제입니다. 이 방식은 전용 Docker network와 egress proxy를 두고 allowlist, path, method, RPS를 프록시에서 강제하는 구조입니다.

## 대상 제한

모든 웹/API 스캔은 allowlist를 통과해야 합니다.

```yaml
targets:
  web:
    url: https://staging.example.com
    allowlist:
      - staging.example.com
```

차단 조건:

- target host가 allowlist에 없음
- scanner가 redirect 제한 capability를 지원하고 redirect 대상이 allowlist 밖으로 벗어남
- OpenAPI server URL이 allowlist 밖임
- scanner가 allowlist scoping capability를 지원하지 않는데 profile이 해당 보장을 요구함

allowlist 검증은 두 단계입니다.

1. 실행 전 target URL과 OpenAPI server URL을 검증합니다.
2. scanner adapter capability를 확인해 런타임 이탈을 제한할 수 있는지 판단합니다.

## Intrusiveness 등급

scanner는 profile 이름이 아니라 동작 성격으로 분류합니다.

```text
passive      응답 관찰, 헤더 검사, 코드/의존성/시크릿 스캔
interactive  API 호출, crawling, property/fuzz 테스트처럼 상태 변화 가능성이 낮지만 요청을 생성하는 스캔
active       공격성 payload, active DAST, intrusive template 실행
```

보호 정책:

- `passive`는 기본 profile에서 허용할 수 있습니다.
- `interactive`는 staging target, allowlist, blocked paths, rate limit capability를 요구합니다.
- `active`는 `--i-understand-active-scan` 같은 명시 플래그가 없으면 실행하지 않습니다.

Schemathesis, ZAP API scan, intrusive Nuclei templates는 profile 이름과 무관하게 최소 `interactive`로 취급합니다.

## Active Scan 보호

`full` profile은 기본적으로 실행되지 않습니다.

```bash
scanrail run --profile full
```

이 명령은 다음 조건을 만족하지 않으면 실패해야 합니다.

- `--i-understand-active-scan` 플래그 존재
- target environment가 `production`이 아님
- allowlist 설정 존재
- blocked paths 설정 검토 완료

## 운영 환경 보호

다음 값은 운영 환경으로 간주합니다.

```yaml
environment: production
```

또는 URL/host에 다음 패턴이 있으면 경고합니다.

```text
prod
production
www
api.example.com
```

운영 환경에서는 기본적으로 다음만 허용합니다.

- secrets scan
- dependency scan
- SAST
- passive headers check
- TLS check

운영 환경에서 금지되는 기본 작업:

- ZAP active scan
- intrusive nuclei templates
- destructive method fuzzing
- high-rate crawling

## 경로 차단

다음 경로는 기본 차단 또는 경고 대상입니다.

```yaml
safety:
  blocked_paths:
    - /logout
    - /delete
    - /remove
    - /destroy
    - /payment
    - /admin/destructive/*
```

차단 대상 HTTP method:

- `DELETE`
- `PATCH`
- destructive semantic을 가진 `POST`

OpenAPI 기반 스캔에서는 operationId, summary, path를 함께 보고 destructive 가능성을 판단합니다.

## 요청 제한

기본 rate limit:

```yaml
safety:
  max_rps: 5
```

scanner별 요청량이 직접 제어되지 않는 경우, 해당 scanner adapter는 `rate_limit` capability를 false로 선언해야 합니다. profile이 rate limit 보장을 요구하면 해당 scanner는 제외하거나 명시 경고 후 skip합니다.

## 식별 헤더

모든 웹/API 요청에는 식별 가능한 스캔 헤더를 추가할 수 있습니다.

```yaml
safety:
  add_header:
    X-Scanrail-Project: my-order-api
```

목적:

- 서버 로그에서 스캔 트래픽 식별
- 장애 분석
- WAF/SIEM 예외 처리

## Secret 처리

금지:

- `scanrail.yaml`에 토큰 원문 저장
- 리포트에 Authorization header 원문 저장
- raw request/response에 cookie 원문 저장

허용:

- 환경변수 이름 저장
- CI secret 이름 저장
- redacted evidence 저장

예시:

```yaml
auth:
  type: bearer
  token_env: SCANRAIL_TOKEN
```

## 리포트 redaction

다음 값은 리포트 생성 전에 마스킹합니다.

- `Authorization`
- `Cookie`
- `Set-Cookie`
- `X-Api-Key`
- access token
- refresh token
- password
- session id

마스킹 예시:

```text
Authorization: Bearer [REDACTED]
Cookie: SESSION=[REDACTED]
```

현재 구현은 `internal/safety`의 central redactor를 사용하며 report JSON/HTML persistence, MCP report/run output, MCP audit event 전에 masking합니다. Gitleaks raw artifact는 저장 전 secret과 match field를 redaction한 형태로 다시 씁니다. 이후 SARIF output을 노출할 때도 같은 boundary를 반드시 통과해야 합니다.

MCP-triggered scan attempt는 JSON Lines 형식으로 `.scanrail/logs/mcp-audit.jsonl`에 기록됩니다. denied, started, completed event에는 tool, decision, redacted target, target host, profile, 가능한 경우 exit code가 포함됩니다.

## 실패 우선 원칙

다음 상황에서는 경고 후 계속하지 않고 실패해야 합니다.

- allowlist 없음
- secret 환경변수 없음
- active scan 플래그 없음
- production target에 intrusive scan 요청
- CI 또는 `--strict-lock` 실행에서 Docker image digest 불일치
- scanner 결과 파싱 실패
- 리포트 redaction 실패

## 감사 로그

각 실행은 최소한 다음 정보를 남깁니다.

- 실행 시간
- 실행자 또는 CI job id
- profile
- target host
- scanner 목록
- scanner image digest
- policy 결과
- report path

secret, token, cookie 원문은 로그에 남기지 않습니다.
