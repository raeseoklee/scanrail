[ENGLISH](demo-tape-scenario.md) | [한국어](demo-tape-scenario.ko.md)

# Demo Tape Scenario

Scanrail은 작은 demo/tape scenario로 OSS 사용자가 동작을 눈으로 확인하고 재현할 수 있게 합니다.

## Headers Demo

첫 재사용 demo target은 `examples/headers-demo`에 있습니다.

```bash
node examples/headers-demo/server.mjs
scanrail init --non-interactive --project-name headers-demo --target http://127.0.0.1:18080 --force
scanrail run --only headers
```

이 demo가 증명하는 것:

- npm으로 설치한 `scanrail`이 host platform에서 실행됨
- config generation이 동작함
- native headers scanner가 allowlist된 local target을 scan할 수 있음
- JSON 및 HTML report가 생성됨

## VHS Recording

recording script:

```bash
vhs examples/headers-demo/tapes/headers-demo.tape
```

생성된 GIF는 짧고 읽기 쉬우며 local secret이나 private path가 포함되지 않을 때만 commit합니다.

## 기존 Adapter Spike Tape

scanner adapter spike에는 별도 tape가 있습니다.

```bash
vhs experiments/scanner-adapter-spike/tapes/scanner-adapter-spike.tape
```

이 tape는 이후 Docker-backed adapter normalization 작업을 설명할 때 사용합니다.
