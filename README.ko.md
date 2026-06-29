[ENGLISH](README.md) | [한국어](README.ko.md)

# Scanrail

[![CI](https://github.com/raeseoklee/scanrail/actions/workflows/ci.yml/badge.svg)](https://github.com/raeseoklee/scanrail/actions/workflows/ci.yml)
[![npm](https://img.shields.io/npm/v/scanrail.svg)](https://www.npmjs.com/package/scanrail)
[![License](https://img.shields.io/github/license/raeseoklee/scanrail.svg)](LICENSE)

개발자가 직접 웹서비스 보안진단을 실행할 수 있도록, 검증된 오픈소스 보안 도구들을 하나의 CLI로 묶는 보안진단 오케스트레이터입니다.

목표는 상용 보안 제품을 그대로 재구현하는 것이 아니라, 개인 개발자와 조직이 각자의 개발 흐름에 맞게 다음을 표준화하는 것입니다.

- 보안 도구 설치와 실행 방식
- 프로젝트별 스캔 설정
- 인증과 허용 범위 관리
- 결과 정규화와 위험도 산정
- 개발자가 이해할 수 있는 리포트
- CI/CD 연동

## 핵심 방향

`scanrail`은 자체 취약점 탐지 엔진을 처음부터 만들지 않습니다. 대신 OWASP ZAP, Nuclei, Semgrep, Trivy, Gitleaks 같은 오픈소스 도구를 격리된 adapter로 실행하고, 결과를 하나의 형식으로 모아 개발자에게 제공합니다. 현재 MVP는 native security headers scanner, native TLS certificate baseline scanner, native OpenAPI baseline scanner, configured target guardrail, Docker 기반 Gitleaks secrets adapter를 제공합니다.

```text
scanrail init
scanrail setup
scanrail run --profile quick
scanrail mcp serve
```

개발자 관점의 기본 흐름은 다음과 같습니다.

1. Docker 기반 스캐너 사용 시 Docker 설치
2. `scanrail init`으로 프로젝트 설정 생성
3. `scanrail setup`으로 스캐너 이미지와 캐시 준비
4. `scanrail run`으로 점검 실행
5. HTML, JSON 리포트 확인

## 설치

```bash
npm install -g scanrail
scanrail doctor
```

기본 quick profile은 Gitleaks로 로컬 secret을 점검하고, web target이 있으면 security headers와 TLS 인증서 상태를, local OpenAPI spec이 있으면 API contract baseline도 함께 점검합니다.

```bash
scanrail init --non-interactive --project-name demo --target https://example.com --openapi ./openapi.yaml
scanrail run --profile quick
```

Docker를 사용할 수 없으면 `scanrail run --only headers`, `scanrail run --only tls`, `scanrail run --only openapi`로 native 점검만 실행할 수 있고, `scanrail run --only gitleaks`로 Docker 기반 secret scan만 실행할 수 있습니다.

npm package는 얇은 JavaScript command wrapper와 현재 OS/CPU에 맞는 Go binary package를 설치합니다. scoped `@scanrail/cli` package는 내부 wrapper package로 계속 사용할 수 있습니다. Docker 기반 Trivy, Semgrep adapter는 다음 구현 대상이며, 실제 command/output 계약은 [Scanner Adapter 실증](docs/experiments/scanner-adapter-spike.ko.md)에 남겨두었습니다.

stdio MCP를 지원하는 AI client에서는 local MCP server를 실행할 수 있습니다.

```bash
scanrail mcp serve
```

MCP MVP는 `doctor`, config read, latest report summary, 그리고 명시적 active-scan 확인이 필요한 native headers scan tool을 제공합니다.
stdio MCP 경로는 [MCP Workbench 검증](examples/mcp-workbench/README.ko.md)으로 회귀 테스트할 수 있습니다.

## 문서

- [제품 요구사항](docs/product-requirements.ko.md)
- [아키텍처](docs/architecture.ko.md)
- [Go 상세 설계](docs/go-technical-design.ko.md)
- [구현 계획](docs/implementation-plan.ko.md)
- [ADR-0001: Go Core with npm Wrapper](docs/adr/0001-go-npm-wrapper.ko.md)
- [제품명과 npm 확인](docs/naming.ko.md)
- [초기 설치 및 실행 시나리오](docs/setup-scenario.ko.md)
- [CLI 명세](docs/cli-reference.ko.md)
- [설정 파일 명세](docs/config-reference.ko.md)
- [안전 모델](docs/safety-model.ko.md)
- [오픈소스 도구 검토](docs/open-source-tools.ko.md)
- [OSS 전략](docs/oss-strategy.ko.md)
- [배포 전략](docs/distribution.ko.md)
- [npm Publish Runbook](docs/npm-publish.ko.md)
- [릴리스 노트 0.2.2](docs/releases/0.2.2.ko.md)
- [릴리스 노트 0.2.1](docs/releases/0.2.1.ko.md)
- [릴리스 노트 0.2.0](docs/releases/0.2.0.ko.md)
- [Release Checklist](docs/release-checklist.ko.md)
- [릴리스 리스크 레지스터](docs/release-risk-register.ko.md)
- [리스크 처리 계획](docs/risk-treatment-plan.ko.md)
- [MCP 설계](docs/mcp-design.ko.md)
- [MCP Workbench 검증](examples/mcp-workbench/README.ko.md)
- [Demo Tape Scenario](docs/demo-tape-scenario.ko.md)
- [기여 가이드](CONTRIBUTING.ko.md)
- [보안 정책](SECURITY.ko.md)
- [로드맵](docs/roadmap.ko.md)
- [Scanner Adapter 실증](docs/experiments/scanner-adapter-spike.ko.md)

## 데모 녹화

MCP 검증:

![MCP verification tape](docs/assets/mcp-verification.gif)

Scanner adapter 실증:

![Scanner adapter spike tape](docs/assets/scanner-adapter-spike.gif)

## MVP 범위

현재 MVP는 다음 기능을 지원합니다.

- `scanrail doctor`
- `scanrail init --non-interactive`
- `scanrail setup`
- `scanrail run --only headers`
- `scanrail run --only gitleaks`
- `scanrail run --only tls`
- `scanrail run --only openapi`
- `scanrail run --profile quick`
- `scanrail mcp serve`
- JSON, HTML 리포트 생성
- macOS, Windows, Linux용 npm wrapper와 platform binary package
- release dry-run 자동화

Gitleaks는 `ghcr.io/gitleaks/gitleaks:v8.30.1`로 고정된 첫 Docker 기반 adapter입니다. native headers/TLS scanner는 network contact 전에 configured web target allowlist와 blocked/excluded path를 강제합니다. native TLS scanner는 HTTPS target에 단일 TLS handshake를 수행하고 certificate trust, hostname, expiry, legacy protocol finding을 리포트합니다. native OpenAPI scanner는 local JSON 또는 일반적인 YAML OpenAPI file을 읽고 version/server metadata, plain HTTP server URL, effective security requirement가 없는 operation, client error response 문서 누락을 리포트합니다. Trivy, Semgrep adapter와 SARIF 리포트는 다음 단계에서 구현합니다.

## 개발

```bash
go test ./...
npm test
node scripts/build-release.mjs
npm pack --workspaces --dry-run
```

현재 MVP는 `doctor`, `init --non-interactive`, `setup`, `run --only headers`, `run --only gitleaks`, `run --only tls`, `run --only openapi`, `mcp serve`를 지원합니다. Trivy, Semgrep adapter는 패키징 골격과 실행 정책을 유지하고, scanner별 command generation과 normalization을 다음 단계에서 채웁니다.

## 라이선스

Apache-2.0. [LICENSE](LICENSE)를 참고하세요.

## 기본 안전 원칙

- Active Scan은 기본 비활성화합니다.
- native interactive scanner는 allowlist 밖 도메인을 요청하기 전에 거부합니다.
- configured `targets.web.exclude_paths`와 `safety.blocked_paths`는 native interactive scanner 실행 전에 차단합니다.
- 토큰과 비밀번호는 설정 파일에 저장하지 않습니다.
- 모든 스캔 요청에는 식별 가능한 스캔 헤더를 추가할 수 있습니다.
- 운영 환경 대상 스캔은 별도 명시 플래그를 요구합니다.

## 참조 표준

- OWASP Top 10
- OWASP ASVS
- OWASP WSTG
- OWASP API Security Top 10
- CWE
- CVSS
- EPSS
- CISA KEV
- SARIF
- CycloneDX / SPDX
