[ENGLISH](architecture.md) | [한국어](architecture.ko.md)

# 아키텍처

## 개요

`scanrail`은 오픈소스 보안 스캐너를 직접 내장하지 않고, Docker 컨테이너로 실행한 뒤 결과를 수집/정규화/리포팅하는 오케스트레이터입니다.

```text
Developer / CI
      |
      v
 scanrail CLI
      |
      +-- config loader
      +-- setup manager
      +-- scan orchestrator
      +-- auth/session manager
      +-- scanner adapters
      +-- finding normalizer
      +-- risk engine
      +-- report generator
```

## 주요 구성요소

### CLI

사용자가 직접 호출하는 인터페이스입니다.

```text
scanrail doctor
scanrail init
scanrail setup
scanrail run
scanrail update
scanrail ci init
```

### Config Loader

`scanrail.yaml`, `tools.lock.yaml`, 환경변수, CLI 옵션을 병합합니다.

우선순위:

1. CLI option
2. Environment variable
3. `scanrail.yaml`
4. Organization default (v0.1에서는 예약 레이어)
5. Built-in default

### Setup Manager

Docker 기반 실행 환경을 준비합니다.

담당 작업:

- Docker daemon 확인
- Docker network 생성
- 이미지 pull
- 캐시 볼륨 준비
- scanner DB/rule/template 업데이트
- 도구 버전 lock 파일 생성

### Scan Orchestrator

선택된 profile에 따라 scanner adapter를 순서대로 실행합니다.

실행 원칙:

- 독립적인 스캐너는 병렬 실행 가능
- 같은 target에 과도한 요청을 보내는 스캐너는 rate limit 적용
- active scan은 별도 승인 플래그가 있어야 실행
- allowlist 검증 실패 시 즉시 중단

### Auth/Session Manager

인증이 필요한 대상에 대해 scanner가 사용할 수 있는 인증 정보를 구성합니다.

지원 예정 방식:

- none
- bearer token
- cookie
- form login
- recorded browser session

중요 원칙:

- secret은 `scanrail.yaml`에 저장하지 않는다.
- 환경변수 또는 CI secret 참조만 저장한다.
- 리포트에는 secret 원문을 남기지 않는다.

### Scanner Adapters

각 오픈소스 도구 실행과 결과 파싱을 담당합니다.

v0.1 adapter:

- Gitleaks
- Trivy
- Semgrep
- native security headers checker

후속 adapter:

- OWASP ZAP baseline
- OWASP ZAP API scan
- Nuclei safe templates
- testssl.sh
- Schemathesis
- CodeQL
- Nikto
- Amass
- Nmap

### Finding Normalizer

도구별 결과를 공통 모델로 변환합니다.

```yaml
id: finding-001
title: Missing Content-Security-Policy header
severity: medium
confidence: high
source:
  tool: zap
  rule_id: 10038
target:
  type: url
  value: https://staging.example.com
classification:
  cwe:
    - CWE-693
  owasp:
    - A05:2021 Security Misconfiguration
evidence:
  request: redacted
  response: redacted
remediation: Add a restrictive Content-Security-Policy header.
```

### Risk Engine

finding의 우선순위를 계산합니다.

고려 요소:

- scanner severity
- confidence
- CVSS
- EPSS (v0.4 이후 외부 피드 연동)
- CISA KEV 등재 여부 (v0.4 이후 외부 피드 연동)
- asset criticality
- production exposure
- authentication requirement
- exploitability

v0.1은 scanner severity와 confidence 중심으로 우선순위를 계산합니다. EPSS/KEV 같은 외부 데이터 피드는 갱신 주기, 오프라인 동작, 캐시 무결성 정책을 정의한 뒤 v0.4에서 도입합니다.

### Report Generator

결과를 여러 포맷으로 생성합니다.

- HTML
- JSON
- SARIF
- JUnit XML

## 데이터 흐름

```text
scanrail.yaml
   |
   v
profile selection
   |
   v
scanner execution
   |
   v
raw scanner outputs
   |
   v
normalization
   |
   v
dedupe and risk scoring
   |
   v
reports
```

## 파일 구조 제안

```text
.
├─ scanrail.yaml
├─ tools.lock.yaml
├─ .env.scanrail.example
└─ .scanrail/
   ├─ cache/
   ├─ raw/
   └─ reports/
```

## 안전장치

- allowlist 필수
- active scan 기본 비활성화
- production target 보호
- blocked path 지원
- 요청 rate limit
- scanner 요청 헤더 추가
- secret redaction
- destructive endpoint 기본 차단
- full profile 명시 플래그 요구
