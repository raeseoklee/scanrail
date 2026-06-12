[ENGLISH](go-technical-design.md) | [한국어](go-technical-design.ko.md)

# Go 상세 설계

## 목적

이 문서는 Scanrail의 제품 아키텍처를 실제 Go 구현 설계로 옮깁니다. CLI 본체는 Go가 담당하고, npm은 배포와 명령 shim 역할만 수행합니다.

## 설계 원칙

- 제품 동작은 Go 바이너리가 소유합니다.
- npm wrapper는 배포 계층입니다.
- Docker scanner 실행은 runner 인터페이스 뒤에 격리합니다.
- v0.1에서는 동적 plugin이 아니라 코드 구조상 adapter로 scanner를 통합합니다.
- raw scanner output은 디버깅용으로 보존하되, report는 normalized finding을 사용합니다.
- 모든 네트워크성 scanner 실행 전 safety gate를 적용합니다.
- secret은 환경변수 이름으로만 참조하고, 저장 전 redact합니다.
- macOS, Windows, Linux, CI에서 동작을 검증합니다.

## 저장소 구조

```text
.
├─ cmd/scanrail/
├─ internal/
│  ├─ app/
│  ├─ cli/
│  ├─ config/
│  ├─ report/
│  ├─ scanners/
│  └─ workspace/
├─ packages/npm/
├─ docs/
└─ scripts/
```

주요 package 역할:

| Package | 책임 |
| --- | --- |
| `cmd/scanrail` | process entrypoint |
| `internal/cli` | command, flag, argument routing |
| `internal/app` | command use case와 dependency wiring |
| `internal/config` | 설정 구조, 기본값, 검증 |
| `internal/workspace` | `.scanrail` 경로와 OS path 처리 |
| `internal/scanners` | scanner adapter |
| `internal/report` | HTML/JSON/SARIF report |

## CLI 계층

초기 구현은 표준 라이브러리 기반의 얇은 command router를 사용합니다. 장기적으로 command가 복잡해지면 `cobra` 도입을 검토할 수 있습니다.

지원 command:

```text
scanrail doctor
scanrail init
scanrail setup
scanrail run
scanrail version
```

## 설정 계층

`scanrail.yaml`은 프로젝트별 설정입니다. v0.1 후보는 제한된 YAML subset을 다루며, 이후 일반 YAML parser 도입을 검토합니다.

중요 원칙:

- raw secret 저장 금지
- target allowlist 명시
- safe profile 기본값
- profile-selected scanner의 target missing은 skip evidence로 기록
- explicit `--only` scanner의 target missing은 실패 처리

## Scanner Adapter 계약

각 adapter는 다음을 선언해야 합니다.

- 필요한 target type
- intrusiveness level
- safety capability
- 실행 방식
- normalized finding 변환 방식
- raw output 보존 정책

이 계약은 profile이 scanner가 제공할 수 없는 safety guarantee를 약속하지 않도록 막습니다.

## Report 모델

Finding은 공통 필드로 정규화됩니다.

```text
id
scanner
title
severity
category
target
evidence
remediation
references
```

초기 report format:

- JSON
- HTML

계획:

- SARIF
- JUnit XML

## npm 배포

Go binary는 OS/CPU별로 빌드되고 npm platform package에 포함됩니다.

```text
@scanrail/cli
@scanrail/cli-darwin-arm64
@scanrail/cli-darwin-x64
@scanrail/cli-win32-x64
@scanrail/cli-win32-arm64
@scanrail/cli-linux-x64
@scanrail/cli-linux-arm64
```

`@scanrail/cli` wrapper는 현재 platform package를 찾아 argument, signal, exit code를 Go binary에 전달합니다.

## 릴리스 검증

기본 검증:

```bash
go test ./...
npm test
make release-dry-run
```

릴리스 dry-run은 모든 platform binary를 빌드하고 npm workspace package를 dry-run으로 pack합니다.
