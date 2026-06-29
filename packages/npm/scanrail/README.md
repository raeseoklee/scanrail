# Scanrail

[![npm](https://img.shields.io/npm/v/scanrail.svg)](https://www.npmjs.com/package/scanrail)
[![CI](https://github.com/raeseoklee/scanrail/actions/workflows/ci.yml/badge.svg)](https://github.com/raeseoklee/scanrail/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/raeseoklee/scanrail.svg)](https://github.com/raeseoklee/scanrail/blob/main/LICENSE)

Developer-first security scan orchestration from one CLI.

This package installs the `scanrail` command. It delegates to `@scanrail/cli`, which installs the matching platform-specific Go binary package for macOS, Windows, or Linux.

## Install

```bash
npm install -g scanrail
scanrail doctor
```

You can also run it without a global install:

```bash
npx scanrail doctor
```

## First Scan

```bash
scanrail init --non-interactive --project-name demo --target https://example.com --openapi ./openapi.yaml
scanrail run --profile quick
```

The current MVP includes the CLI scaffold, config generation, workspace setup, JSON/HTML reporting, native security headers, TLS certificate, and local OpenAPI baseline scanners, configured web target guardrails, and a Docker-backed Gitleaks secrets adapter. Use `scanrail run --only headers`, `scanrail run --only tls`, or `scanrail run --only openapi` without Docker, or `scanrail run --only gitleaks` for the secrets scan only. Trivy and Semgrep adapters are planned.

Native interactive scanners reject targets outside `targets.web.allowlist` and block configured `targets.web.exclude_paths` / `safety.blocked_paths` before making network contact.

## MCP

Scanrail includes a local stdio MCP server for AI clients:

```bash
scanrail mcp serve
```

The MCP MVP exposes bounded tools for `doctor`, config reading, latest report summaries, and the native headers scan with explicit active-scan confirmation.

## Package Layout

`scanrail` is the recommended npm entrypoint. It depends on `@scanrail/cli`, which installs one optional platform package:

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
- Security: https://github.com/raeseoklee/scanrail/blob/main/SECURITY.md

## License

Apache-2.0
