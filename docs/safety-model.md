[ENGLISH](safety-model.md) | [한국어](safety-model.ko.md)

# Safety Model

Scanrail is a developer self-check tool, but some scanners can generate intrusive or attack-like traffic. The safety model is therefore part of the product, not an add-on.

## Principles

- Scan only authorized targets.
- Block active production scanning by default.
- Do not claim allowlist enforcement for scanners that cannot provide it.
- Exclude or skip scanners that cannot enforce required destructive-path controls.
- Do not store secrets in config or reports.
- Validate rate-limit support through scanner adapter capabilities.
- Preserve enough execution metadata for reproducibility.

## Capability-Based Enforcement

The v0.1 safety model is based on scanner capability contracts. Scanrail does not claim it can network-block every redirect, method, path, or request rate inside third-party containers.

Each scanner adapter must declare support for:

```text
allowlist_scope
redirect_scope
blocked_paths
blocked_methods
rate_limit
header_injection
auth_injection
```

If a profile requires a capability the scanner cannot provide, Scanrail must:

- skip the scanner and record the reason, or
- fail when the user explicitly selected that scanner, or
- keep the scanner out of default profiles until support is clear.

Current MVP enforcement:

- Scanner capability metadata is declared in `internal/scanners`.
- Gitleaks is the first production-ready Docker-backed passive scanner and mounts the workspace read-only.
- Trivy and Semgrep remain non-production-ready and are skipped in profile execution.
- Explicit execution of a non-production-ready scanner fails with safety exit code `5`.
- The native headers scanner does not follow redirects and declares `allowlist_scope`, `redirect_scope`, `rate_limit`, and `header_injection`.
- The native TLS scanner performs one TLS handshake for HTTPS targets, sends no HTTP payload, and declares `allowlist_scope`, `redirect_scope`, and `rate_limit`.
- The native OpenAPI scanner reads a local spec file only. It does not fetch remote specs, use credentials, or call API endpoints.

Network-level enforcement through a dedicated Docker network and egress proxy is a later v0.x design.

## Target Restrictions

All web and API scans must pass allowlist validation.

```yaml
targets:
  web:
    url: https://staging.example.com
    allowlist:
      - staging.example.com
```

Block conditions:

- target host is not in the allowlist
- OpenAPI server URL is outside the allowlist
- scanner redirects outside the allowlist when redirect scoping is required
- scanner lacks allowlist scoping while the profile requires it

Allowlist validation has two layers:

1. Validate configured target URLs before execution.
2. Check whether the scanner adapter can maintain scope during runtime.

## Intrusiveness Levels

Scanners are classified by behavior, not only by profile name.

```text
passive      observes responses or scans local code/dependencies/secrets
interactive  creates requests with low expected state change, such as crawling or API property tests
active       sends attack payloads, active DAST checks, or intrusive templates
```

Policy:

- `passive` can be allowed in default profiles.
- `interactive` requires staging targets, allowlists, blocked paths, and rate-limit capabilities.
- `active` requires explicit opt-in, such as an `--i-understand-active-scan` style flag.

The native TLS scanner is classified as `interactive` because it opens a network connection, but its current implementation only performs a single TLS handshake. The native OpenAPI scanner is `passive` because it reads only local files. Schemathesis, ZAP API scans, and intrusive Nuclei templates are at least `interactive`, regardless of profile name.

## Production Targets

Production active scanning is blocked by default. Enabling it should require:

- explicit CLI flag
- explicit config setting
- allowlist
- low rate limits
- clear report evidence
- organization approval where applicable

## Secrets

Secrets must be referenced, not stored.

Allowed:

```yaml
auth:
  type: bearer
  token_env: SCANRAIL_TOKEN
```

Disallowed:

```yaml
auth:
  token: raw-secret-value
```

All logs, raw command previews, and reports must redact secret values. The current implementation centralizes masking in `internal/safety` and applies it before report JSON/HTML persistence, MCP report/run outputs, and MCP audit events. Gitleaks raw artifacts are also rewritten with secret and match fields redacted before being exposed. Future SARIF output must pass through the same boundary before being exposed.

MCP-triggered scan attempts are written to `.scanrail/logs/mcp-audit.jsonl` as JSON Lines. Denied, started, and completed events include the tool, decision, redacted target, target host, profile, and exit code when available.

## Exit Codes

Safety violations return exit code `5`. Policy failures based on findings return exit code `1`. Configuration mistakes return `2`, and runtime environment errors return `3`.
