[ENGLISH](README.md) | [한국어](README.ko.md)

# Scanrail

개발자가 직접 웹서비스 보안진단을 실행할 수 있도록, 검증된 오픈소스 보안 도구들을 하나의 CLI로 묶는 보안진단 오케스트레이터입니다.

목표는 상용 보안 제품을 그대로 재구현하는 것이 아니라, 개인 개발자와 조직이 각자의 개발 흐름에 맞게 다음을 표준화하는 것입니다.

- 보안 도구 설치와 실행 방식
- 프로젝트별 스캔 설정
- 인증과 허용 범위 관리
- 결과 정규화와 위험도 산정
- 개발자가 이해할 수 있는 리포트
- CI/CD 연동

## 핵심 방향

`scanrail`은 자체 취약점 탐지 엔진을 처음부터 만들지 않습니다. 대신 OWASP ZAP, Nuclei, Semgrep, Trivy, Gitleaks 같은 오픈소스 도구를 Docker 기반으로 실행하고, 결과를 하나의 형식으로 모아 개발자에게 제공합니다.

```text
scanrail init
scanrail setup
scanrail run --profile quick
```

개발자 관점의 기본 흐름은 다음과 같습니다.

1. Docker 설치
2. `scanrail init`으로 프로젝트 설정 생성
3. `scanrail setup`으로 스캐너 이미지와 캐시 준비
4. `scanrail run`으로 점검 실행
5. HTML, JSON, SARIF 리포트 확인

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
- [로드맵](docs/roadmap.ko.md)
- [Scanner Adapter 실증](docs/experiments/scanner-adapter-spike.ko.md)

## MVP 범위

초기 버전은 다음 기능을 우선 지원합니다.

- 인터랙티브 초기 설정: `scanrail init`
- Docker 기반 스캐너 준비: `scanrail setup`
- 기본 진단 실행: `scanrail run`
- 코드/의존성/시크릿/보안 헤더 점검
- staging 웹 런타임/API/TLS 점검은 v0.2 이후 단계적으로 확장
- HTML, JSON, SARIF 리포트 생성
- CI 실패 기준 설정
- Active Scan 안전장치

## 개발

```bash
go test ./...
npm test
node scripts/build-release.mjs
npm pack --workspaces --dry-run
```

현재 첫 배포 후보는 `doctor`, `init --non-interactive`, `setup --pull-policy never`, `run --only headers`를 지원합니다. Docker 기반 Gitleaks, Trivy, Semgrep adapter는 패키징 골격과 실행 정책을 먼저 잡고, scanner별 command generation과 normalization을 다음 단계에서 채웁니다.

## 기본 안전 원칙

- Active Scan은 기본 비활성화합니다.
- allowlist 밖 도메인은 요청하지 않습니다.
- 토큰과 비밀번호는 설정 파일에 저장하지 않습니다.
- destructive path는 기본 차단하거나 명시 경고합니다.
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
