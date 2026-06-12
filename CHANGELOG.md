[ENGLISH](CHANGELOG.md) | [한국어](CHANGELOG.ko.md)

# Changelog

## 0.1.2 - 2026-06-12

- Add `scanrail mcp serve` as a local stdio MCP MVP.
- Expose MCP tools for doctor, config reading, latest report summaries, and a guarded native headers scan.
- Parse generated target allowlists for MCP target validation.
- Add OSS contribution, security, issue template, and demo/tape documentation.
- Add npm smoke workflow and release workflow validation mode.

## 0.1.1 - 2026-06-12

- Publish a no-functional-change patch release through npm Trusted Publishing.
- Keep `scanrail` as the recommended npm entrypoint backed by `@scanrail/cli`.
- Clean up CI and npm publish workflow warnings before the trusted publishing run.

## 0.1.0 - 2026-06-12

- Add initial Go CLI scaffold.
- Add `doctor`, `init`, `setup`, and `run` commands.
- Add native security headers scanner.
- Add JSON and HTML report generation.
- Add npm wrapper package and platform binary package manifests.
- Add `scanrail` unscoped npm alias package.
- Add cross-platform release build dry-run script.
- Add bilingual documentation structure with English as the default language and Korean documents under `.ko.md`.
