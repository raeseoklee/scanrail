[ENGLISH](0001-go-npm-wrapper.md) | [한국어](0001-go-npm-wrapper.ko.md)

# ADR-0001: Go Core와 npm Wrapper 배포

## 상태

Accepted.

## Context

Scanrail은 OSS-first 보안 스캔 오케스트레이터입니다. 목표는 `scanrail init`, `scanrail setup`, `scanrail run`으로 이어지는 단순한 개발자 흐름을 제공하면서 Docker 기반 scanner를 실행하고 결과를 정규화된 report로 만드는 것입니다.

프로젝트는 macOS, Windows, Linux, CI에서 잘 동작해야 합니다. 대상 사용자는 Node/npm을 이미 갖고 있을 가능성이 높지만, core program은 안정적인 process 실행, Docker orchestration, filesystem handling, 명확한 exit code, release discipline이 필요합니다.

## Decision

Scanrail core CLI는 Go로 만들고, 1차 배포 채널은 npm wrapper package로 사용합니다.

배포 모델:

```text
@scanrail/cli
@scanrail/cli-darwin-arm64
@scanrail/cli-darwin-x64
@scanrail/cli-win32-x64
@scanrail/cli-win32-arm64
@scanrail/cli-linux-x64
@scanrail/cli-linux-arm64
```

`@scanrail/cli`는 `scanrail` command를 노출하고, OS/CPU별 binary package를 `optionalDependencies`로 참조합니다. wrapper는 현재 platform에 맞는 binary를 찾아 모든 argument를 Go binary로 전달합니다.

## Drivers

- cross-platform 설치 UX가 단순해야 합니다.
- scanner orchestrator는 안정적인 subprocess/Docker 실행이 필요합니다.
- CI는 deterministic exit code와 machine-readable report가 필요합니다.
- npm package는 기본적으로 postinstall binary download를 피해야 합니다.
- core는 npm 밖에서도 release archive나 future Homebrew/Scoop package로 사용할 수 있어야 합니다.
- 보안 도구 자체이므로 runtime dependency surface를 작게 유지해야 합니다.

## Alternatives Considered

### TypeScript Core Published Directly to npm

장점:

- 초기 CLI 개발 속도가 빠릅니다.
- npm 배포가 자연스럽습니다.

단점:

- runtime Node가 필요합니다.
- dependency surface와 supply-chain exposure가 커집니다.
- core에서 cross-platform quoting/path/process edge case가 늘어납니다.

결론: 설치 편의성만으로는 충분하지 않습니다. orchestrator core는 Go의 single-binary model과 process control이 더 적합합니다.

### Go Core with Only Native Package Managers

장점:

- 깨끗한 Go binary 배포가 가능합니다.
- Homebrew, Scoop, winget, release archive와 잘 맞습니다.

단점:

- npm에 익숙한 개발자에게 초기 설치 friction이 큽니다.
- 초기부터 여러 package manager channel을 관리해야 합니다.

결론: 첫 OSS release에서는 npm이 macOS, Windows, Linux에서 하나의 공통 설치 명령을 제공합니다.

### Go Core with postinstall Download npm Package

장점:

- npm package를 하나로 줄일 수 있습니다.

단점:

- postinstall network download는 보안 도구 배포에 부적합합니다.
- locked-down CI와 offline 환경에서 실패하기 쉽습니다.
- provenance와 checksum 검증 부담이 커집니다.

결론: platform-specific optional package가 더 투명하고 npm 생태계에 맞습니다.

## Consequences

좋은 점:

- core runtime dependency가 작습니다.
- npm 설치 UX를 유지합니다.
- platform binary를 명시적으로 관리할 수 있습니다.
- CI와 release archive 배포가 쉬워집니다.

비용:

- platform package가 여러 개 필요합니다.
- release automation이 중요해집니다.
- npm scope 확보와 publish 순서가 필요합니다.

## Follow-up

- npm scope 확보
- release provenance와 checksum 추가
- Homebrew/Scoop 배포 검토
- adapter별 Docker image lock 정책 구체화
