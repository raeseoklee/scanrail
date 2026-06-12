[ENGLISH](setup-scenario.md) | [한국어](setup-scenario.ko.md)

# 초기 설치 및 실행 시나리오

이 문서는 개발자가 `scanrail`을 처음 설치하고 프로젝트에 적용하는 흐름을 보여줍니다.

## 1. 설치

```bash
npm install -g scanrail
```

첫 공개 릴리스 전에는 npm 설치를 기본 경로로 사용합니다. release archive와 설치 스크립트는 릴리스 자동화가 준비된 뒤 추가합니다.

## 2. 환경 점검

```bash
scanrail doctor
```

예상 출력:

```text
Scanrail Doctor

Docker              OK   Docker Desktop running
Docker Compose      OK
Network             OK
Workspace           OK   /Users/dev/workspace/my-order-api
Git repo            OK
Disk space          OK
Internet access     OK   scanner image pull available

Ready.
```

## 3. 초기 설정

```bash
scanrail init
```

예상 인터랙션:

```text
? 프로젝트 이름은? my-order-api

? 점검 대상 유형을 선택하세요
  ✓ 로컬 코드 저장소
  ✓ Staging 웹 URL
  ✓ OpenAPI/Swagger API
  ○ 컨테이너 이미지

? Staging URL은?
  https://staging-order.example.com

? 이 URL은 운영 환경인가요?
  No

? 허용할 도메인을 확인하세요
  ✓ staging-order.example.com

? OpenAPI 파일 또는 URL은?
  ./openapi.yaml

? 인증이 필요한가요?
  Bearer token

? 토큰은 어디서 읽을까요?
  환경변수

? 환경변수 이름은?
  SCANRAIL_TOKEN

? 기본 점검 프로파일은?
  quick

? CI에서 실패 처리할 기준은?
  high 이상

? 제외할 경로가 있나요?
  /logout
  /admin/destructive/*
  /payments/real-charge

? Active Scan을 기본 허용할까요?
  No

? 설정 파일을 생성할까요?
  Yes
```

## 4. 생성되는 설정 파일

`scanrail.yaml`

```yaml
project:
  name: my-order-api
  type: web-api

targets:
  web:
    url: https://staging-order.example.com
    allowlist:
      - staging-order.example.com
    exclude_paths:
      - /logout
      - /admin/destructive/*
      - /payments/real-charge

  api:
    openapi: ./openapi.yaml

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

report:
  output_dir: .scanrail/reports
  formats:
    - html
    - json
    - sarif
```

`.env.scanrail.example`

```bash
SCANRAIL_TOKEN=replace-me
```

## 5. 스캐너 준비

```bash
scanrail setup
```

예상 출력:

```text
Preparing Scanrail runtime

Docker network       created scanrail-net
Reports directory    created .scanrail/reports
Cache directory      created .scanrail/cache

Pulling scanners
- gitleaks            OK
- trivy               OK
- semgrep             OK

Updating scanner data
- trivy DB            OK
- semgrep rules       OK

Generated tools.lock.yaml
Setup complete.
```

## 6. 최초 실행

```bash
export SCANRAIL_TOKEN="eyJ..."
scanrail run
```

또는 profile을 명시할 수 있습니다.

```bash
scanrail run --profile quick
```

예상 출력:

```text
Scanrail started

Project       my-order-api
Profile       quick
Target        https://staging-order.example.com
OpenAPI       ./openapi.yaml
Auth          bearer via SCANRAIL_TOKEN

Running checks
[1/4] Gitleaks secrets scan        PASS
[2/4] Trivy dependency scan        WARN   3 findings
[3/4] Semgrep SAST                 WARN   5 findings
[4/4] Security headers             WARN   2 findings

Policy result
FAILED: high severity finding exists

Reports
HTML   .scanrail/reports/my-order-api-20260612-1430.html
JSON   .scanrail/reports/my-order-api-20260612-1430.json
SARIF  .scanrail/reports/my-order-api-20260612-1430.sarif
```

## 7. Full Scan 보호

v0.1 기본 생성 설정에는 `full` profile이 포함되지 않습니다. active scanner adapter가 구현되고 사용자가 profile을 명시적으로 추가한 경우에만 다음 보호 동작을 적용합니다.

```bash
scanrail run --profile full
```

예상 출력:

```text
Refusing to run full profile.

Reason:
- full profile includes active scan
- active scan is disabled by default

Run only against approved staging targets:

scanrail run --profile full --i-understand-active-scan
```

명시적으로 허용한 경우:

```bash
scanrail run --profile full --i-understand-active-scan
```

## 8. CI 설정 생성

```bash
scanrail ci init
```

예상 인터랙션:

```text
? CI 환경은?
  GitHub Actions

? PR마다 실행할 프로파일은?
  quick

? main merge 후 실행할 프로파일은?
  standard

? SARIF 업로드를 활성화할까요?
  Yes
```

생성 예시:

```yaml
name: Security Scan

on:
  pull_request:
  push:
    branches:
      - main

jobs:
  scanrail:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 22
      - run: npm install -g scanrail
      - run: scanrail run --profile quick
        env:
          SCANRAIL_TOKEN: ${{ secrets.SCANRAIL_TOKEN }}
```
