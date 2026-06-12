[ENGLISH](0001-go-npm-wrapper.md) | [한국어](0001-go-npm-wrapper.ko.md)

# ADR-0001: Go Core with npm Wrapper Distribution

## Status

Accepted.

## Context

`scanrail` is an OSS-first security scan orchestrator. The existing product direction is to provide a simple developer flow, `scanrail init`, `scanrail setup`, and `scanrail run`, while orchestrating Docker-based scanners and producing normalized reports. The current architecture separates configuration loading, setup, scanner adapters, finding normalization, risk scoring, and report generation.

The project must work well on macOS, Windows, Linux, and CI. Many target users already have Node/npm available, but the core program needs stable process execution, Docker orchestration, filesystem handling, predictable exit codes, and long-term release discipline.

## Decision

Build the `scanrail` core CLI in Go and distribute it primarily through an npm wrapper package.

The release model is:

```text
@scanrail/cli
@scanrail/cli-darwin-arm64
@scanrail/cli-darwin-x64
@scanrail/cli-win32-x64
@scanrail/cli-win32-arm64
@scanrail/cli-linux-x64
@scanrail/cli-linux-arm64
```

`@scanrail/cli` exposes the `scanrail` command and depends on platform-specific binary packages through npm `optionalDependencies`. The wrapper locates the installed binary package for the current OS/CPU and forwards all arguments to the Go binary.

## Drivers

- Cross-platform developer install UX must be simple.
- The scanner orchestrator needs reliable subprocess and Docker execution.
- CI usage must have deterministic exit codes and machine-readable reports.
- The OSS package must avoid postinstall binary downloads by default.
- The core should remain usable outside npm through release archives and future Homebrew/Scoop packages.
- The project should keep runtime dependencies small because it is itself a security tool.

## Alternatives Considered

### TypeScript Core Published Directly to npm

Pros:

- Fastest initial CLI development.
- Native npm distribution.
- Strong prompt/config/reporting ecosystem.

Cons:

- Requires Node at runtime.
- Larger dependency surface and supply-chain exposure.
- More cross-platform quoting/path/process edge cases in the core.
- Harder to provide a standalone archive for CI or locked-down environments.

Rejected because installation convenience is not the only requirement. The core orchestrator benefits from Go's single-binary model and predictable process control.

### Go Core with Only Native Package Managers

Pros:

- Clean Go binary distribution.
- Works well with Homebrew, Scoop, winget, and release archives.
- No npm wrapper complexity.

Cons:

- Higher first-run friction for developers who already use npm.
- Requires multiple package-manager-specific channels early.
- Windows adoption is slower without a familiar install path.

Rejected for the initial OSS release because npm gives one common install command across macOS, Windows, and Linux.

### Go Core with postinstall Download npm Package

Pros:

- Only one npm package.
- Smaller npm artifact size.

Cons:

- `postinstall` network downloads are often blocked by security policy.
- Harder to audit and reproduce.
- More fragile in offline and private registry environments.

Rejected because platform-specific npm packages are more transparent and registry-native.

## Consequences

Positive:

- Users can run `npm install -g scanrail`.
- The Go binary remains usable without npm.
- Runtime dependency surface stays small.
- Release artifacts can support npm, archives, Homebrew, and Scoop from the same build matrix.

Negative:

- Release automation must publish multiple npm packages per version.
- Version synchronization across packages is required.
- Windows/macOS binary signing may still be needed for enterprise adoption.
- The wrapper must handle unsupported platforms clearly.

## Implementation Requirements

- The Go module owns all CLI behavior and exit codes.
- The npm wrapper must be thin and must not duplicate business logic.
- Platform packages must declare `os` and `cpu` fields.
- The public npm package must not include organization-specific URLs, secrets, or policies.
- Release CI must build, test, package, checksum, and publish all artifacts from a tagged commit.

## Follow-ups

- Add Go project scaffold.
- Add npm wrapper package scaffold under `packages/npm`.
- Add release automation for cross-compilation and npm publish dry-runs.
- Add Windows/macOS signing after the MVP proves useful.
- Add Homebrew and Scoop channels after npm release stabilizes.
