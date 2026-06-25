[ENGLISH](implementation-plan.md) | [한국어](implementation-plan.ko.md)

# 구현 계획

## 범위

이 계획은 Scanrail의 첫 OSS 구현을 다룹니다. 목표는 Go CLI와 npm wrapper를 통해 `doctor`, `init`, `setup`, `run`을 안전한 기본값으로 제공하는 MVP를 만드는 것입니다.

## 인수 기준

- `scanrail doctor`가 macOS, Windows, Linux, CI에서 동작합니다.
- `scanrail init --non-interactive`가 유효한 `scanrail.yaml`을 생성합니다.
- `scanrail setup --pull-policy never`가 offline/dry-run 상황에서 동작합니다.
- `scanrail run --only headers`가 HTTP target을 검사하고 JSON/HTML report를 생성합니다.
- profile-selected scanner의 missing target은 evidence와 함께 skip됩니다.
- explicit `--only` scanner의 missing target은 실패합니다.
- safety capability mismatch는 실행 전에 처리됩니다.
- policy failure는 exit code `1`, safety violation은 `5`를 반환합니다.
- npm wrapper는 현재 platform binary를 resolve하고 실행합니다.
- release CI는 모든 target OS/CPU binary와 npm package dry-run을 검증합니다.

## Phase 0: Project Scaffold

목표:

- 최소 Go project, command wiring, tests, release placeholder를 만듭니다.

주요 파일:

```text
go.mod
cmd/scanrail/main.go
internal/cli/
internal/app/
internal/version/
Makefile
.github/workflows/ci.yml
```

완료 기준:

- `go test ./...` 통과
- `scanrail --version` 동작
- CI 실행

## Phase 1: Config and Workspace

목표:

- typed config loading, defaults, validation, workspace path를 구현합니다.

완료 기준:

- `scanrail init --non-interactive`가 config 생성
- 기존 config overwrite는 `--force` 필요
- secret raw value 저장 금지
- `.scanrail/reports` 등 workspace directory 생성

## Phase 2: Native Headers Scanner

목표:

- Docker 없이 동작하는 첫 scanner를 구현해 end-to-end report pipeline을 검증합니다.

검사 항목:

- Content-Security-Policy
- X-Content-Type-Options
- X-Frame-Options
- Referrer-Policy
- Strict-Transport-Security for HTTPS

완료 기준:

- local HTTP test server 대상으로 finding 생성
- JSON/HTML report 생성

## Phase 3: Docker Adapter Framework

목표:

- Gitleaks, Trivy, Semgrep 같은 Docker-backed scanner를 공통 runner interface 뒤에 둡니다.

완료 기준:

- command preview와 raw output path 기록
- scanner capability 계약 구현
- missing target과 safety mismatch 처리

## Phase 4: npm Wrapper and Release

목표:

- npm wrapper와 platform package를 구성합니다.

완료 기준:

- macOS, Windows, Linux binary package 생성
- wrapper가 argument와 exit code를 보존
- `npm pack --workspaces --dry-run` 통과

## Phase 5: CI and OSS Readiness

목표:

- 공개 저장소에서 반복 가능한 검증과 문서화를 제공합니다.

완료 기준:

- OS matrix test 통과
- release dry-run 통과
- README, docs, license, changelog 준비
- npm publish 직전 상태까지 확인

## 현재 상태

현재 MVP는 Go CLI, native headers scanner, native TLS baseline scanner, Docker 기반 Gitleaks adapter, JSON/HTML report, npm wrapper, platform package, release dry-run을 포함합니다. Trivy와 Semgrep 같은 추가 Docker-backed scanner adapter는 다음 단계 구현 대상입니다.
