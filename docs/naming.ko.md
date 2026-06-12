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

결과:

- `scanrail`: `E404 Not Found`
- `@scanrail/cli`: `E404 Not Found`
- `scanrail-darwin-arm64`: `E404 Not Found`
- `@scanrail/cli-darwin-arm64`: `E404 Not Found`
- `npm search scanrail --json`: `[]`

해석:

- unscoped package `scanrail`은 registry상 비어 있는 것으로 보입니다.
- scoped package `@scanrail/cli`도 아직 게시된 package가 없습니다.
- 다만 `@scanrail` scope 자체의 생성 가능 여부는 npm 로그인 후 사용자/조직 생성 단계에서 최종 확인해야 합니다. `npm org ls scanrail --json`은 현재 로컬 npm 인증 문제로 `E401`을 반환했습니다.

## Package Strategy

선호:

```text
@scanrail/cli
@scanrail/cli-darwin-arm64
@scanrail/cli-darwin-x64
@scanrail/cli-win32-x64
@scanrail/cli-win32-arm64
@scanrail/cli-linux-x64
@scanrail/cli-linux-arm64
```

fallback:

```text
scanrail
scanrail-darwin-arm64
scanrail-darwin-x64
scanrail-win32-x64
scanrail-win32-arm64
scanrail-linux-x64
scanrail-linux-arm64
```

scope를 확보할 수 있으면 scoped package를 우선 사용합니다. scope 확보가 어려우면 unscoped `scanrail` package로 전환합니다.
