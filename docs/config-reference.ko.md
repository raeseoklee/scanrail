[ENGLISH](config-reference.md) | [한국어](config-reference.ko.md)

# 설정 파일 명세

`scanrail.yaml`은 프로젝트별 보안 점검 설정을 담습니다. 이 파일은 코드 저장소에 커밋할 수 있어야 하며, secret 원문을 포함하면 안 됩니다.

## 전체 예시

```yaml
project:
  name: my-order-api
  type: web-api
  owner: payments-platform
  criticality: high

targets:
  web:
    url: https://staging-order.example.com
    environment: staging
    allowlist:
      - staging-order.example.com
    exclude_paths:
      - /logout
      - /admin/destructive/*
      - /payments/real-charge

  api:
    openapi: ./openapi.yaml

  container:
    image: registry.example.com/order-api:staging

auth:
  type: bearer
  token_env: SCANRAIL_TOKEN

profiles:
  default: quick

  quick:
    tools:
      - gitleaks
      - trivy
      - semgrep
      - headers

safety:
  active_scan_default: false
  require_allowlist: true
  max_rps: 5
  add_header:
    X-Scanrail-Project: my-order-api
  blocked_paths:
    - /logout
    - /admin/destructive/*
    - /payments/real-charge

policy:
  fail_on:
    severity: high
  ignore:
    - id: finding-id
      reason: accepted risk
      expires: 2026-12-31

report:
  output_dir: .scanrail/reports
  formats:
    - html
    - json
    - sarif
```

## project

프로젝트 메타데이터입니다.

```yaml
project:
  name: my-order-api
  type: web-api
  owner: payments-platform
  criticality: high
```

필드:

| 필드 | 필수 | 설명 |
| --- | --- | --- |
| name | 예 | 프로젝트 이름 |
| type | 아니오 | `web`, `api`, `web-api`, `service` 등 |
| owner | 아니오 | 담당 팀 |
| criticality | 아니오 | `low`, `medium`, `high`, `critical` |

## targets

스캔 대상을 정의합니다.

```yaml
targets:
  web:
    url: https://staging.example.com
    environment: staging
    allowlist:
      - staging.example.com
    exclude_paths:
      - /logout
```

지원 대상:

- `web.url`
- `api.openapi`
- `container.image`
- `repo.path`

## auth

인증 정보를 정의합니다. secret 원문은 저장하지 않습니다.

```yaml
auth:
  type: bearer
  token_env: SCANRAIL_TOKEN
```

지원 타입:

```text
none
bearer
cookie
form-login
recorded-session
```

v0.1 구현 상태:

- `none`, `bearer`, `cookie`를 우선 지원합니다.
- `form-login`, `recorded-session`은 예약된 타입입니다. 구현 전에는 설정 검증에서 명확한 오류를 반환해야 합니다.

예시:

```yaml
auth:
  type: cookie
  cookie_env: SCANRAIL_COOKIE
```

```yaml
auth:
  type: form-login
  login_url: https://staging.example.com/login
  username_env: SCANRAIL_USERNAME
  password_env: SCANRAIL_PASSWORD
```

## profiles

실행할 도구 묶음을 정의합니다.

```yaml
profiles:
  default: quick
  quick:
    tools:
      - gitleaks
      - trivy
      - semgrep
      - headers
```

기본 profile:

- `quick`
- `standard`
- `full`

프로파일 의미:

- `quick`: 코드 저장소만 있어도 동작해야 하는 기본 셀프체크입니다.
- `standard`: `quick`을 포함하고, 웹 target이 있으면 passive DAST를 추가합니다.
- `full`: `standard`를 포함하고, 명시 플래그가 있을 때만 active scanner를 추가합니다.

`extends`는 부모 profile의 `tools`를 먼저 포함한 뒤 현재 profile의 `tools`를 뒤에 추가합니다. 같은 도구가 중복되면 한 번만 실행합니다.

대상 없는 스캐너 처리:

- profile에 포함된 도구가 필요한 target을 찾지 못하면 기본적으로 `SKIP` 처리하고 리포트에 남깁니다.
- 사용자가 `--only <tool>`로 명시 실행한 도구가 target을 찾지 못하면 설정 오류로 실패합니다.
- `strict_targets: true`를 도입하면 profile 실행에서도 skip 대신 실패하도록 확장할 수 있습니다.

v0.1에서 생성되는 기본 프로파일은 실제 구현된 adapter만 포함해야 합니다. Nuclei, ZAP API, Schemathesis, testssl.sh 같은 도구는 v0.2 이후 프로파일 예시로 문서화하되, v0.1 `scanrail init`의 기본 출력에는 포함하지 않습니다.

확장 profile 예시:

```yaml
profiles:
  default: quick

  quick:
    tools:
      - gitleaks
      - trivy
      - semgrep
      - headers

  standard:
    extends: quick
    tools:
      - zap-baseline

  full:
    extends: standard
    requires_explicit_flag: true
    tools:
      - zap-active
```

이 예시는 해당 adapter가 구현된 버전에서만 생성해야 합니다.

## safety

스캔 안전장치를 정의합니다.

```yaml
safety:
  active_scan_default: false
  require_allowlist: true
  max_rps: 5
```

권장 기본값:

```yaml
active_scan_default: false
require_allowlist: true
max_rps: 5
```

`targets.web.exclude_paths`와 `safety.blocked_paths`는 의미가 다릅니다.

- `targets.web.exclude_paths`: crawler나 scanner가 탐색하지 않아야 하는 경로입니다.
- `safety.blocked_paths`: 요청 자체가 금지되는 경로입니다. scanner가 이 차단을 강제할 수 없으면 해당 scanner는 제외되거나 경고 후 skip되어야 합니다.

## policy

CI 실패 기준과 예외를 정의합니다.

```yaml
policy:
  fail_on:
    severity: high
```

지원 severity:

```text
critical
high
medium
low
info
never
```

예외 예시:

```yaml
policy:
  ignore:
    - id: semgrep.javascript.express.security.audit.xss.mustache
      reason: false positive, output is escaped by template engine
      expires: 2026-12-31
```

예외 정책:

- `reason` 필수
- `expires` 권장
- 만료된 예외는 실패 처리 가능

## report

리포트 출력 설정입니다.

```yaml
report:
  output_dir: .scanrail/reports
  formats:
    - html
    - json
    - sarif
```

지원 포맷:

- `html`
- `json`
- `sarif`
- `junit`

## tools.lock.yaml

`tools.lock.yaml`은 실제 사용한 scanner 이미지와 버전을 기록합니다.

```yaml
tools:
  zap:
    image: ghcr.io/zaproxy/zaproxy:stable
    digest: sha256:...
  trivy:
    image: aquasec/trivy:0.71.1
    digest: sha256:...
generated_at: 2026-06-12T00:00:00Z
```

운영 정책:

- CI에서는 lock 파일 기준으로 실행합니다.
- 보안팀이 승인한 버전만 허용할 수 있습니다.
- `scanrail update --lock`으로 갱신합니다.

lock 검증 정책:

- 로컬 기본 실행은 digest 불일치를 경고로 처리합니다.
- CI 또는 `--strict-lock` 실행은 digest 불일치를 실패로 처리합니다.
- `latest` 태그는 개발 초기 예시에만 허용하고, 재현 가능한 실행이 필요한 환경에서는 고정 tag 또는 digest를 사용해야 합니다.

## 조직 기본값

config 병합 순서에는 organization defaults 레이어가 있지만, v0.1에서는 예약된 확장 지점입니다. v0.1 구현은 built-in defaults, `scanrail.yaml`, 환경변수, CLI 옵션만 병합합니다.
