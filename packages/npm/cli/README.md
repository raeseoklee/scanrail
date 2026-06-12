# Scanrail CLI

[![npm](https://img.shields.io/npm/v/@scanrail/cli.svg)](https://www.npmjs.com/package/@scanrail/cli)
[![CI](https://github.com/raeseoklee/scanrail/actions/workflows/ci.yml/badge.svg)](https://github.com/raeseoklee/scanrail/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/raeseoklee/scanrail.svg)](https://github.com/raeseoklee/scanrail/blob/main/LICENSE)

Developer-first security scan orchestration from one CLI.

This package installs the `scanrail` command and delegates to the platform-specific Go binary package for macOS, Windows, or Linux.

## Install

```bash
npm install -g @scanrail/cli
scanrail doctor
```

You can also run it without a global install:

```bash
npx @scanrail/cli doctor
```

## First Scan

```bash
scanrail init --non-interactive --project-name demo --target https://example.com
scanrail run --only headers
```

The first release candidate includes the CLI scaffold, config generation, workspace setup, JSON/HTML reporting, and a native security headers scanner. Docker-backed adapters for Gitleaks, Trivy, and Semgrep are planned next.

## Package Layout

`@scanrail/cli` is the wrapper package. It installs one optional platform package:

- `@scanrail/cli-darwin-arm64`
- `@scanrail/cli-darwin-x64`
- `@scanrail/cli-win32-x64`
- `@scanrail/cli-win32-arm64`
- `@scanrail/cli-linux-x64`
- `@scanrail/cli-linux-arm64`

## Links

- Repository: https://github.com/raeseoklee/scanrail
- Documentation: https://github.com/raeseoklee/scanrail#readme
- Issues: https://github.com/raeseoklee/scanrail/issues

## License

Apache-2.0
