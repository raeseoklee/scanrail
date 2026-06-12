[ENGLISH](distribution.md) | [한국어](distribution.ko.md)

# Distribution Strategy

Scanrail uses a Go binary for the core CLI and an npm wrapper as the primary developer distribution channel. The goal is one familiar install path across macOS, Windows, and Linux.

## Recommended Model

```text
Go binary
  |
  +-- npm wrapper package
  +-- platform-specific npm binary packages
  +-- GitHub release archives
  +-- optional Homebrew/Scoop packages
```

## npm Package Layout

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

`scanrail` is the recommended user-facing npm package. It depends on `@scanrail/cli`, which exposes the `scanrail` command and depends on platform packages through `optionalDependencies`.

## User Installation

Global install:

```bash
npm install -g scanrail
scanrail init
scanrail setup
scanrail run
```

One-time execution:

```bash
npx scanrail run --profile quick
```

The scoped `@scanrail/cli` package remains available for users who want to install the underlying wrapper directly.

## Wrapper Package

```json
{
  "name": "scanrail",
  "version": "0.1.0",
  "bin": {
    "scanrail": "./bin/scanrail.js"
  },
  "dependencies": {
    "@scanrail/cli": "0.1.0"
  }
}
```

The scoped wrapper should:

- detect `process.platform` and `process.arch`
- resolve the matching platform package
- forward all arguments and signals
- preserve the Go binary exit code
- avoid `postinstall` binary downloads

No `exports` field is needed for the wrapper package because users invoke the CLI through `bin`.

## Platform Package

```json
{
  "name": "@scanrail/cli-win32-x64",
  "version": "0.1.0",
  "os": ["win32"],
  "cpu": ["x64"],
  "files": ["scanrail.exe"]
}
```

Platform packages contain only the binary and minimal package metadata.

## Release Dry-Run

The release dry-run should:

1. run Go tests
2. build all target binaries
3. copy binaries into npm platform packages
4. test the wrapper
5. run `npm pack --workspaces --dry-run`

## Future Channels

- GitHub release archives
- Homebrew tap
- Scoop bucket
- winget package
- container image for CI-only execution
