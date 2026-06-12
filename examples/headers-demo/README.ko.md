[ENGLISH](README.md) | [한국어](README.ko.md)

# Headers Demo

이 demo는 Scanrail native headers scanner를 검증하기 위한 작은 local HTTP target입니다.

## 실행

터미널 1:

```bash
node examples/headers-demo/server.mjs
```

터미널 2:

```bash
scanrail init --non-interactive --project-name headers-demo --target http://127.0.0.1:18080 --force
scanrail run --only headers
```

예상 결과:

- `.scanrail/reports/*.json`
- `.scanrail/reports/*.html`
- `/` 경로의 missing security headers finding

`/secure` 경로는 기준 security headers를 설정하므로 이후 scanner 비교 테스트에 사용할 수 있습니다.

## Tape

VHS 시나리오는 `examples/headers-demo/tapes/headers-demo.tape`에 있습니다.

```bash
vhs examples/headers-demo/tapes/headers-demo.tape
```
