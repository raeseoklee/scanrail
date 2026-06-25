[ENGLISH](risk-treatment-plan.md) | [한국어](risk-treatment-plan.ko.md)

# Risk Treatment Plan

This document turns the release risk register into operating gates. It does not replace the [Release Risk Register](release-risk-register.md); it explains how the remaining risks are handled before releases or feature expansion.

Last reviewed: June 16, 2026.

## Current Decision

There is no known remaining release blocker for the current headers, TLS, and Gitleaks MVP when the release checklist and verification gates pass.

The remaining risks are expansion gates:

- Do not broaden MCP execution until production-host compatibility and audit behavior are verified.
- Do not enable additional Docker-backed scanners until adapter isolation, redaction, raw artifact handling, and version metadata are implemented for each adapter.
- Do not position Scanrail as a PR-native CI security gate until SARIF and artifact workflows are shipped.
- Do not bundle third-party scanner rule packs without license metadata.

## Release Gates

These risks must be checked before every npm publish.

| Risk | Gate | Evidence |
| --- | --- | --- |
| R-001 npm trusted publishing settings drift | Validate npm trusted publisher settings against `.github/workflows/npm-publish.yml`. | Completed [Release Checklist](release-checklist.md) trusted-publisher section. |
| R-002 partial publish leaves package set inconsistent | Capture package versions before and after publish; publish platform packages before wrappers. | Completed [Release Checklist](release-checklist.md) registry snapshot section. |
| R-004 public npm smoke is not per-commit | Run scheduled/manual npm smoke after publish for the public package version. | `npm Smoke` workflow URL and result. |
| R-006 MCP tools bypass core policy | Run MCP Workbench regression before publish when MCP code changed. | `mcp-workbench run examples/mcp-workbench/scanrail-mcp.yaml --verbose`. |

## MCP Expansion Gates

These gates must pass before adding `scanrail_setup`, broader scanner execution, or richer report resources to MCP.

| Risk | Gate | Required before expansion |
| --- | --- | --- |
| R-005 production MCP host compatibility | Verify at least one production MCP host configuration. | Host name, config snippet, smoke result, and date. |
| R-007 sensitive or oversized MCP resources | Add size limits and redaction tests for every richer resource. | Tests for maximum payload size and secret masking. |
| R-008 MCP auditability | Extend `.scanrail/logs/mcp-audit.jsonl` events to new setup/scanner tools. | Tests for denied, started, completed, and failed paths. |

## Scanner Adapter Gates

These gates must pass before Trivy, Semgrep, or future Docker-backed adapters move from skipped scaffold to execution. Gitleaks passed the first slice of these gates in `0.1.4`; broader OS matrix coverage and digest-level image policy remain open.

| Risk | Gate | Required before enabling |
| --- | --- | --- |
| R-009 scanner adapter exposure | Enforce capability contracts in the runner, not only catalog metadata. | Adapter tests for allowlist, redirects, credentials, and raw output. |
| R-010 Docker host variance | Run Linux adapter integration first, then document macOS and Windows Docker behavior. | OS compatibility notes or CI matrix. |
| R-011 scanner image/rule drift | Pin scanner image and rule versions. | Report metadata records scanner name, version, image digest where available. |
| R-012 normalization loss | Preserve raw artifacts and link normalized findings to scanner-specific evidence. | Round-trip tests from raw output to normalized finding. |
| R-013 raw output secret leaks | Apply central redaction to raw artifacts and future SARIF. | Redaction tests for raw logs, JSON, SARIF, and failure output. |
| R-018 license/attribution | Track license metadata for bundled rule packs or templates. | License manifest entry before bundling third-party assets. |

## Roadmap Risks

These are not current MVP blockers, but they must shape milestone planning.

| Risk | Treatment |
| --- | --- |
| R-003 non-npm release artifacts | Accepted for npm MVP; add checksums and GitHub release assets before stable v1.0. |
| R-014 user overestimates current coverage | Keep MVP/non-goal language in README, docs, report text, and release notes. |
| R-015 accepted-risk workflow | Implement expiring ignore rules before severity override or team policy features. |
| R-016 SARIF/CI artifact integration | Required before PR-native security gate positioning. |
| R-017 OSS governance | Add maintainer response expectations and labels when external issues appear. |

## Reclassification Rule

A risk can move out of the active register only when all are true:

- the mitigation is implemented or explicitly accepted,
- verification evidence exists in git, CI, npm, or Codexus artifacts,
- the next required action is either complete or moved to a milestone gate,
- the Korean and English docs agree.

If any condition is missing, keep the risk active and mark the exact gate that remains.
