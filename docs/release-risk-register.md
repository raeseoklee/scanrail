[ENGLISH](release-risk-register.md) | [한국어](release-risk-register.ko.md)

# Release and Product Risk Register

This register tracks the material risks that remain after the first public npm release, trusted publishing setup, the `0.1.2` MCP MVP, MCP audit logging, the MCP Workbench verification pass, and the `0.1.4` Gitleaks adapter slice.

Last reviewed: June 16, 2026.

## Risk Levels

| Level | Meaning |
| --- | --- |
| High | Can block release adoption, create unsafe scan behavior, or break public package trust. |
| Medium | Can confuse users, reduce reliability, or delay the next milestone if not managed. |
| Low | Known gap with limited blast radius or clear manual workaround. |

## Current Snapshot

| Area | Status | Main remaining risk |
| --- | --- | --- |
| npm distribution | Stable for MVP with checklist gate | Trusted publisher settings and npm package state live outside git. |
| Cross-platform install | Covered by smoke workflow | Public npm smoke is scheduled/manual, not run on every commit. |
| MCP MVP | Implemented, audit-logged, and Workbench-verified | Real production host compatibility is not complete. |
| Scanner adapters | First Docker-backed adapter implemented for Gitleaks | Trivy/Semgrep execution, OS matrix coverage, and digest-level image policy are still planned. |
| Reporting | JSON/HTML for headers and Gitleaks | SARIF and false-positive workflow are still planned. |
| OSS operations | Public repo is usable | Community triage, support boundaries, and governance are still lightweight. |

## Immediate Priorities

1. Run the [Release Checklist](release-checklist.md) for every npm publish (`R-001`, `R-002`, `R-004`).
2. Confirm at least one production MCP host before adding broader MCP execution tools (`R-005`).
3. Extend the redaction boundary to future SARIF before exposing richer report data (`R-007`, `R-013`, `R-016`).
4. Add Docker adapter integration coverage before promoting Trivy or Semgrep from skipped scaffolds (`R-009`, `R-010`).
5. Preserve scanner version metadata and raw artifact links as normalized reports expand (`R-011`, `R-012`).

The operating gates for all remaining risks are tracked in the [Risk Treatment Plan](risk-treatment-plan.md).

## Active Risks

| ID | Risk | Level | Status | Mitigation | Next action |
| --- | --- | --- | --- | --- | --- |
| R-001 | npm trusted publishing settings drift | High | Checklist-gated external risk | Settings live in npm, not git. Release docs and the release checklist record the required workflow filename, repository, owner, allowed action, and per-package confirmation step. | Run the trusted-publisher section of the [Release Checklist](release-checklist.md) before each publish. |
| R-002 | Partial publish leaves package set inconsistent | High | Checklist-gated operational risk | Platform packages publish before wrappers, versions are checked before publish, and the release checklist records pre/post package state. Failed versions must be fixed with a later patch version. | Record the registry snapshot and post-publish package state in the [Release Checklist](release-checklist.md). |
| R-003 | Release artifacts outside npm are unsigned or absent | Medium | Accepted for npm MVP | npm provenance is in place. GitHub release archives, checksums, and additional package-manager channels remain roadmap work. | Add checksum and GitHub release asset generation before declaring a stable `v1.0` release path. |
| R-004 | Public npm smoke does not run on every commit | Medium | Partially mitigated | `.github/workflows/npm-smoke.yml` covers Ubuntu, macOS, and Windows against the public registry on schedule or manual dispatch. | Decide whether to add a lightweight post-publish required smoke gate for release branches. |
| R-005 | MCP server compatibility is narrower than the client ecosystem | Medium | Partially mitigated | `mcp-workbench inspect` and the Workbench regression spec validate stdio protocol behavior. | Add smoke notes or fixtures for at least one production MCP host before expanding MCP tools. |
| R-006 | MCP tools could bypass CLI safety if they diverge from core policy | High | Controlled by design | MCP stays a thin adapter over the CLI safety model, active scans require `confirm_active_scan=true`, and allowlists are enforced. | Keep MCP tool implementations bound to the same config, exit-code, and safety validation paths as CLI execution. |
| R-007 | MCP resources can leak sensitive configuration or oversized report data | High | Partially mitigated | Current config resources expose secret environment variable names, not values. Planned full report resources are still deferred. | Add size limits and redaction tests before adding `scanrail://reports/latest/json` or richer resources. |
| R-008 | MCP tool calls are not yet auditable enough for team environments | Medium | Partially mitigated | MCP-triggered scan attempts now write denied, started, and completed events to `.scanrail/logs/mcp-audit.jsonl` with redacted targets and exit codes. Confirmed active scans refuse to run when the start audit event cannot be written. | Treat audit coverage as an MCP expansion gate for `scanrail_setup` and future scanner execution. |
| R-009 | Scanner adapters can expand network or credential exposure | High | Partially mitigated | Gitleaks is enabled as a passive local scanner with a read-only workspace mount and writable raw output mount. Unready Docker adapters are skipped in profile runs and fail with exit code `5` when explicitly selected. The native headers scanner declares interactive network capabilities and does not follow redirects. | Add adapter-specific integration coverage before enabling Trivy, Semgrep, or networked Docker scanners. |
| R-010 | Docker-backed scanners can behave differently across host OSes | Medium | Partially mitigated | The scanner adapter spike and Gitleaks unit coverage validate command shape, but full Docker behavior is not yet covered across Linux, macOS, and Windows. | Add macOS/Windows Docker compatibility notes or CI matrix coverage for Gitleaks and future adapters. |
| R-011 | Scanner image or rule updates can change findings unexpectedly | Medium | Partially mitigated | Gitleaks uses a pinned `ghcr.io/gitleaks/gitleaks:v8.30.1` image and records tool version/image metadata in reports. Digest-level pinning and broader update policy are not implemented. | Add digest policy or lockfile validation before treating scanner versions as stable release inputs. |
| R-012 | Normalized reports can hide scanner-specific context | Medium | Partially mitigated | Gitleaks preserves a redacted raw artifact and normalized findings include tool and evidence fields. Broader schema mapping is still evolving. | Add SARIF and false-positive workflow without dropping raw scanner context. |
| R-013 | Secret redaction may miss raw scanner output paths | High | Partially mitigated | Central redaction now masks configured env values, auth headers, cookies, token/password fields, URL userinfo, and secret-like query parameters before report JSON/HTML persistence and MCP report/run outputs. Gitleaks raw artifacts are rewritten with secret and match fields redacted. | Extend the same redaction boundary to future SARIF and any new raw scanner formats before those outputs ship. |
| R-014 | Users may overestimate coverage from the current headers plus Gitleaks scan | Medium | Documentation risk | README and docs state that the MVP is not commercial scanner parity and that Trivy, Semgrep, SARIF, and broader DAST remain planned. | Keep capability and non-goal language visible in README, reports, and release notes. |
| R-015 | Findings workflow is not ready for team-scale accepted risk management | Medium | Roadmap item | Ignore rules, policy packs, issue export, and workflow integrations are listed for later milestones. | Implement expiring ignore rules before adding severity override or team policy features. |
| R-016 | SARIF and CI artifact integrations are not shipped yet | Low | Planned | CI/CD milestone includes SARIF, JUnit XML, and CI templates. | Treat SARIF as required before positioning Scanrail as a PR-native security gate. |
| R-017 | OSS support and governance can lag adoption | Medium | Lightweight process | Contributing, security policy, issue templates, and roadmap exist. | Add maintainer response expectations and label workflow once external issues start arriving. |
| R-018 | License and attribution obligations for bundled scanner rules/templates can be missed | Medium | Future integration risk | The current package does not bundle third-party scanner rule packs. | Track license metadata when adding bundled Semgrep/Nuclei/config assets. |

## Mitigated Risks

| Risk | Mitigation |
| --- | --- |
| GitHub Actions Node 20 runtime deprecation warnings | Workflows use `actions/checkout@v5`, `actions/setup-node@v5`, and `actions/setup-go@v6`, which declare `node24` action runtimes. |
| npm `always-auth` warning during publish workflow | The publish workflow no longer asks `setup-node` to create npm registry auth configuration. Trusted publishing uses OIDC instead of `.npmrc` token auth. |
| Published package does not run on every target OS | The npm smoke workflow installs the public `scanrail` package on Ubuntu, macOS, and Windows, then runs `version`, `doctor`, and signature audit checks. |
| Windows hosted runner label migration | Windows jobs pin `windows-2025-vs2026` explicitly after the June 2026 GitHub runner image notice. |
| MCP MVP lacks real client-style regression coverage | `examples/mcp-workbench/scanrail-mcp.yaml` now verifies MCP discovery, resources, safety gating, confirmed native headers scan execution, and latest report retrieval. |
| npm release package state is not captured in git | [Release Checklist](release-checklist.md) records trusted publisher checks, pre-publish package state, publish workflow evidence, post-publish package state, and failure handling. |

## Verification Gates

Before a publish:

```bash
npm run docs:check-links
npm test
npm run publish:dry-run
make release-dry-run
```

Also complete the [Release Checklist](release-checklist.md).

For the current already-published version, use `mode=validate` in the publish workflow instead of `mode=dry-run`.

After a publish:

```bash
npm view scanrail version
npm run smoke:npm -- <version>
```

Run the `npm Smoke` workflow for OS matrix validation against the public npm registry.

MCP regression check:

```bash
mcp-workbench inspect --command node --args "examples/mcp-workbench/serve-fixture.mjs" --json
mcp-workbench run examples/mcp-workbench/scanrail-mcp.yaml --verbose
```

## External Checks

- Confirm npm trusted publisher settings for every package before changing the workflow filename, repository owner, or npm organization settings.
- Confirm package pages show provenance for the published version.
- Confirm GitHub Actions has no deprecation or auth warnings after workflow changes.
- Confirm scanner image licenses, rule licenses, and update policies before bundling third-party scanner assets.
- Confirm MCP host behavior with at least one production AI client before expanding MCP execution tools.
