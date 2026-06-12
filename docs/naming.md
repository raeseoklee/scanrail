[ENGLISH](naming.md) | [한국어](naming.ko.md)

# Naming

## Decision

The product name is `Scanrail`.

Primary CLI and npm package:

```text
scanrail
@scanrail/cli
```

## Rationale

`Scanrail` suggests a security scanning guardrail built into a developer workflow. It is short, works naturally as a CLI command, and avoids sounding like an offensive security product.

## npmjs Check

Registry checks performed on 2026-06-12:

```bash
npm view scanrail name version description --json
npm view @scanrail/cli name version description --json
npm view scanrail-darwin-arm64 name version description --json
npm view @scanrail/cli-darwin-arm64 name version description --json
npm search scanrail --json
```

Initial observed results before publishing:

- `scanrail`: `E404 Not Found`
- `@scanrail/cli`: `E404 Not Found`
- `scanrail-darwin-arm64`: `E404 Not Found`
- `@scanrail/cli-darwin-arm64`: `E404 Not Found`
- `npm search scanrail --json`: `[]`

Interpretation:

- The unscoped `scanrail` package appeared unclaimed at the time of checking.
- The scoped `@scanrail/cli` package also appeared unpublished.
- The `@scanrail` npm organization was later confirmed under the maintainer account.

## Package Strategy

Use `scanrail` as the recommended user-facing package and keep scoped packages for the wrapper and platform binaries:

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

The unscoped `scanrail` package depends on `@scanrail/cli`. The platform packages remain scoped to avoid occupying six additional unscoped npm names.
