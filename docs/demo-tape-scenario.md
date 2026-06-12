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

## MCP Verification Tape

The MCP verification scenario installs the published npm package, starts a local target, talks to `scanrail mcp serve` over stdio JSON-RPC, checks safety confirmation behavior, runs the native headers scanner, and reads the latest report summary.

```bash
node examples/mcp-verification/run.mjs
vhs examples/mcp-verification/tapes/mcp-verification.tape
```

The generated GIF is stored at `docs/assets/mcp-verification.gif`.

![MCP verification tape](assets/mcp-verification.gif)

## Existing Adapter Spike Tape

The scanner adapter spike has a separate tape:

```bash
vhs experiments/scanner-adapter-spike/tapes/scanner-adapter-spike.tape
```

Use this tape to explain future Docker-backed adapter normalization work.

![Scanner adapter spike tape](assets/scanner-adapter-spike.gif)
