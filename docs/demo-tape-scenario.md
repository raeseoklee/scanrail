[ENGLISH](demo-tape-scenario.md) | [한국어](demo-tape-scenario.ko.md)

# Demo Tape Scenario

Scanrail uses small demo/tape scenarios to make OSS behavior visible and reproducible.

## Headers Demo

The first reusable demo target lives in `examples/headers-demo`.

```bash
node examples/headers-demo/server.mjs
scanrail init --non-interactive --project-name headers-demo --target http://127.0.0.1:18080 --force
scanrail run --only headers
```

The demo proves:

- npm-installed `scanrail` starts on the host platform
- config generation works
- the native headers scanner can scan a local allowlisted target
- JSON and HTML reports are produced

## VHS Recording

The recording script is:

```bash
vhs examples/headers-demo/tapes/headers-demo.tape
```

The generated GIF should be committed only when it is short, readable, and does not contain local secrets or private paths.

## Existing Adapter Spike Tape

The scanner adapter spike has a separate tape:

```bash
vhs experiments/scanner-adapter-spike/tapes/scanner-adapter-spike.tape
```

Use this tape to explain future Docker-backed adapter normalization work.
