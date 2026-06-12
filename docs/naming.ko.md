[ENGLISH](naming.md) | [한국어](naming.ko.md)

# Naming

## Decision

제품명은 `Scanrail`로 사용합니다.

기본 CLI와 npm package:

```text
scanrail
@scanrail/cli
```

## Rationale

`Scanrail`은 보안 스캔을 개발 워크플로의 guardrail처럼 깔아준다는 의미를 담습니다. 이름이 짧고 CLI 명령으로 자연스러우며, 보안 제품처럼 지나치게 무겁거나 공격적으로 들리지 않습니다.

## npmjs Check

2026-06-12 기준 npm registry 조회 결과:

```bash
npm view scanrail name version description --json
npm view @scanrail/cli name version description --json
npm view scanrail-darwin-arm64 name version description --json
npm view @scanrail/cli-darwin-arm64 name version description --json
npm search scanrail --json
```

publish 전 최초 결과:

- `scanrail`: `E404 Not Found`
- `@scanrail/cli`: `E404 Not Found`
- `scanrail-darwin-arm64`: `E404 Not Found`
- `@scanrail/cli-darwin-arm64`: `E404 Not Found`
- `npm search scanrail --json`: `[]`

해석:

- unscoped package `scanrail`은 registry상 비어 있는 것으로 보입니다.
- scoped package `@scanrail/cli`도 아직 게시된 package가 없습니다.
- 이후 maintainer 계정에서 `@scanrail` npm organization을 사용할 수 있음을 확인했습니다.

## Package Strategy

사용자가 설치하는 기본 package는 `scanrail`로 두고, wrapper와 platform binary는 scoped package로 유지합니다.

```text
scanrail
@scanrail/cli
@scanrail/cli-darwin-arm64
@scanrail/cli-darwin-x64
@scanrail/cli-win32-x64
@scanrail/cli-win32-arm64
@scanrail/cli-linux-x64
@scanrail/cli-linux-arm64
```

unscoped `scanrail` package는 `@scanrail/cli`에 의존합니다. platform package는 scoped 이름으로 유지해 불필요하게 unscoped npm 이름 6개를 점유하지 않습니다.
