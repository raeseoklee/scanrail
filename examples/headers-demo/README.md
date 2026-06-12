[ENGLISH](README.md) | [한국어](README.ko.md)

# Headers Demo

This demo starts a tiny local HTTP target for the native Scanrail headers scanner.

## Run

Terminal 1:

```bash
node examples/headers-demo/server.mjs
```

Terminal 2:

```bash
scanrail init --non-interactive --project-name headers-demo --target http://127.0.0.1:18080 --force
scanrail run --only headers
```

Expected result:

- `.scanrail/reports/*.json`
- `.scanrail/reports/*.html`
- findings for missing security headers on `/`

The `/secure` path sets the expected baseline headers and can be used later for scanner comparison tests.

## Tape

The VHS scenario is stored at `examples/headers-demo/tapes/headers-demo.tape`.

```bash
vhs examples/headers-demo/tapes/headers-demo.tape
```
