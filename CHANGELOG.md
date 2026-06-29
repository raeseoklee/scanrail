[ENGLISH](CHANGELOG.md) | [한국어](CHANGELOG.ko.md)

# Changelog

## 0.2.2 - 2026-06-29

- Enforce configured web target allowlists before native interactive scanners run.
- Parse and apply `targets.web.exclude_paths` and `safety.blocked_paths` as preflight path guardrails for `headers`, `tls`, and MCP-triggered headers scans.
- Load `safety.require_allowlist` and `safety.max_rps` from `scanrail.yaml`.
- Require `blocked_paths` capability for interactive scanner definitions and declare it for the native headers and TLS scanners.
- Refresh README, CLI/config/safety docs, roadmap, risk docs, npm README files, and release notes for the target guardrail surface.

## 0.2.1 - 2026-06-25

- Add a native local-file-only OpenAPI baseline scanner.
- Add `scanrail run --only openapi` and `--openapi <path>` overrides for init/run.
- Parse `targets.api.openapi` from `scanrail.yaml` and include `openapi` in generated quick profiles.
- Report OpenAPI findings for missing version/server metadata, plain HTTP server URLs, missing effective operation security, and missing client error responses.
- Refresh README, npm README files, CLI/config docs, roadmap, risk docs, and release notes for the OpenAPI baseline surface.

## 0.2.0 - 2026-06-25

- Add a native TLS certificate and protocol baseline scanner.
- Add `scanrail run --only tls` and include `tls` in generated quick profiles.
- Report certificate trust, hostname, expiry, and legacy TLS protocol findings.
- Keep MCP scanner execution limited to headers while broader networked scans remain safety-gated.
- Refresh README, npm README files, CLI/config/safety docs, roadmap, risk docs, and release notes for the TLS baseline surface.

## 0.1.4 - 2026-06-16

- Add the first production-ready Docker-backed scanner adapter for Gitleaks secrets detection.
- Add a shared Docker runner abstraction and normalize Gitleaks findings into the JSON/HTML report model.
- Preserve scanner metadata in reports and redact Gitleaks raw artifacts before writing them to `.scanrail/raw`.
- Update generated quick profiles to run `gitleaks` and `headers`, while keeping Trivy and Semgrep safety-gated.
- Refresh release, safety, setup, roadmap, and risk docs for the expanded scanner surface.

## 0.1.3 - 2026-06-15

- Add MCP Workbench regression evidence and connect the MCP and scanner demo tapes in the public docs.
- Record MCP scan attempts in the audit log so active-scan requests leave reviewable local evidence.
- Centralize target redaction and scanner capability metadata to avoid leaking raw targets or silently weakening safety policy.
- Stabilize the Windows CI runner selection used by the npm smoke path.
- Convert the remaining release risks into explicit release, MCP expansion, scanner adapter, and roadmap gates.
- Add bilingual release checklist and risk treatment documentation for future npm publishes.

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
