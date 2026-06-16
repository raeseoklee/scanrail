[ENGLISH](implementation-plan.md) | [한국어](implementation-plan.ko.md)

# Implementation Plan

## Scope

This plan covers the first OSS implementation of `scanrail` as a Go CLI distributed through npm wrapper packages. It turns the current documentation into a working MVP that supports `doctor`, `init`, `setup`, and `run` with safe defaults.

## Acceptance Criteria

- `scanrail doctor` runs on macOS, Windows, Linux, and CI.
- `scanrail init --non-interactive` creates a valid `scanrail.yaml`.
- `scanrail setup --pull-policy never` can run with a fake Docker runner in tests.
- `scanrail run --only headers` scans a local test HTTP server and writes JSON/HTML reports.
- `scanrail run --profile quick` can orchestrate at least three Docker scanner adapters behind a fake runner.
- Profile-selected scanners with missing targets are skipped with evidence; explicitly selected scanners with missing targets fail.
- Scanner safety capability mismatches are enforced before execution.
- Policy violations produce exit code `1`; safety violations produce exit code `5`.
- Secrets referenced through env vars are redacted from logs, raw command previews, and reports.
- npm wrapper packages can resolve and execute the correct platform binary in tests.
- Release CI can build all target OS/CPU binaries and package npm artifacts in dry-run mode.

## Phase 0: Project Scaffold

Goal:

- Create a minimal Go project with command wiring, tests, lint, and release placeholders.

Files:

```text
go.mod
cmd/scanrail/main.go
internal/cli/root.go
internal/app/app.go
internal/version/version.go
Makefile
.github/workflows/ci.yml
```

Tasks:

1. Initialize Go module.
2. Add root command with `--version`, `--config`, `--workdir`, and `--verbose`.
3. Add version metadata injected by build flags.
4. Add unit test command in CI.
5. Add `make test`, `make lint`, `make build`.

Acceptance:

- `go test ./...` passes.
- `go run ./cmd/scanrail --version` prints version metadata.
- CI runs on Linux at minimum.

## Phase 1: Config and Workspace

Goal:

- Implement typed config loading, defaults, validation, and workspace paths.

Files:

```text
internal/config/config.go
internal/config/defaults.go
internal/config/load.go
internal/config/validate.go
internal/workspace/workspace.go
testdata/config/
```

Tasks:

1. Define structs matching `docs/config-reference.md`.
2. Decode YAML.
3. Apply built-in defaults.
4. Apply CLI overrides.
5. Validate required fields and safe constraints.
6. Resolve `.scanrail/cache`, `.scanrail/raw`, `.scanrail/reports`, `.scanrail/tmp`.

Acceptance:

- Valid sample config loads successfully.
- Unknown scanner names fail validation.
- Secret-looking literal values in `auth` fail validation.
- Windows path tests pass with `filepath`.

## Phase 2: CLI Commands and Init Wizard

Goal:

- Implement `doctor`, `init`, and command option parsing.

Files:

```text
internal/cli/doctor.go
internal/cli/init.go
internal/doctor/doctor.go
internal/cli/prompt.go
internal/config/write.go
```

Tasks:

1. Implement `scanrail doctor`.
2. Implement `scanrail init --non-interactive`.
3. Add prompt interface and fake prompt tests.
4. Generate `scanrail.yaml` and `.env.scanrail.example`.
5. Refuse overwrite unless `--force`.

Acceptance:

- `scanrail doctor` detects missing Docker without panic.
- `scanrail init --non-interactive --project-name demo --target http://localhost:8080 --profile quick` writes config.
- Existing config is not overwritten without `--force`.

## Phase 3: Docker Runner and Setup

Goal:

- Add Docker CLI abstraction and `scanrail setup`.

Files:

```text
internal/dockerx/runner.go
internal/dockerx/cli_runner.go
internal/dockerx/fake_runner.go
internal/setup/setup.go
internal/setup/tools_lock.go
```

Tasks:

1. Define Docker runner interface.
2. Implement Docker CLI runner using `exec.CommandContext`.
3. Implement fake runner for tests.
4. Add image pull policy: `missing`, `always`, `never`.
5. Add Docker network ensure.
6. Generate `tools.lock.yaml`.

Acceptance:

- Docker commands are built as argument arrays, not shell strings.
- `setup --pull-policy never` does not call pull.
- Image digest is recorded when inspect data is available.
- Runner redacts env values in command previews.

## Phase 4: Finding Model, Policy, and Reports

Goal:

- Create normalized finding pipeline before adding many scanners.

Files:

```text
internal/findings/finding.go
internal/findings/fingerprint.go
internal/policy/policy.go
internal/report/json.go
internal/report/html.go
internal/report/sarif.go
internal/report/testdata/
```

Tasks:

1. Define finding model.
2. Add severity and confidence enums.
3. Add stable fingerprinting.
4. Add policy evaluation.
5. Add JSON report writer.
6. Add HTML report writer.
7. Add minimal SARIF writer for file findings.

Acceptance:

- Golden JSON report is stable.
- HTML report escapes all finding content.
- SARIF validates structurally for file findings.
- Policy severity threshold maps to expected exit code.

## Phase 5: Safety Engine

Goal:

- Enforce safe defaults before scanner execution.

Files:

```text
internal/safety/allowlist.go
internal/safety/active_scan.go
internal/safety/blocked_paths.go
internal/safety/redactor.go
internal/safety/preflight.go
internal/safety/capabilities.go
internal/safety/target_resolver.go
```

Tasks:

1. Validate target host allowlist.
2. Detect production-like targets and enforce restrictions.
3. Block active scan without explicit flag.
4. Validate auth env vars exist.
5. Implement redactor for headers, env values, cookies, tokens, and passwords.
6. Define scanner safety capability requirements.
7. Resolve host-local targets for Docker-backed scanners.

Acceptance:

- Target outside allowlist returns exit code `5`.
- Full profile without explicit flag returns exit code `5`.
- Missing auth env var fails before scanner execution.
- Redactor test covers Authorization, Cookie, Set-Cookie, token, password.
- Scanner missing required safety capability is skipped for profile execution and fails for `--only`.
- `http://localhost:<port>` is converted to a Docker-reachable host target on macOS/Windows, and Linux behavior is documented or explicitly rejected.

## Phase 6: Scanner Registry and Native Headers Scanner

Goal:

- Establish adapter contract with one native scanner.

Files:

```text
internal/scanners/scanner.go
internal/scanners/registry.go
internal/scanners/headers/headers.go
internal/app/run.go
```

Tasks:

1. Define scanner interface.
2. Build scanner registry.
3. Add `SafetyCapabilities()` and `Intrusiveness()` to each scanner.
4. Implement `headers` scanner using Go HTTP client.
5. Implement run orchestrator for scanner selection, target skip, capability filtering, and report writing.

Acceptance:

- `scanrail run --only headers` works against local HTTP test server.
- Findings include missing CSP and missing HSTS when applicable.
- `--skip headers` excludes it.
- Missing target behavior matches profile skip vs explicit fail rules.

## Phase 7: Docker Scanner Adapters v0.1

Goal:

- Add the first Docker-backed scanner adapter with golden output normalization.

Files:

```text
internal/scanners/gitleaks/
internal/scanners/trivy/
internal/scanners/semgrep/
testdata/scanners/
```

Tasks:

1. Implement Gitleaks command generation.
2. Keep Trivy and Semgrep behind the production-readiness gate until their command generation and normalization are implemented.
3. Implement Semgrep command generation.
4. Add raw output parsers.
5. Add golden normalization fixtures.
6. Store raw outputs under `.scanrail/raw/<run-id>/<tool>/`.

Acceptance:

- Fake runner verifies mount/env/args for all three scanners.
- Golden fixtures normalize into expected findings.
- Scanner parse failure returns scanner error, not policy failure.

## Phase 8: ZAP Baseline and Network Scanner Controls

Goal:

- Add the first Docker-backed web runtime scanner.

Files:

```text
internal/scanners/zap/baseline.go
internal/scanners/zap/normalize.go
internal/scanners/network/options.go
```

Tasks:

1. Implement ZAP baseline command generation.
2. Inject auth headers where supported.
3. Apply scanner request header configuration where supported.
4. Parse ZAP JSON output into findings.
5. Keep active scan unavailable until full profile is explicitly implemented.

Acceptance:

- ZAP baseline refuses to run without web target and allowlist.
- Fake runner verifies target URL and output mounts.
- Sample ZAP output normalizes into findings.

## Phase 9: npm Wrapper and Platform Packages

Goal:

- Make the Go binary installable through npm without postinstall download.

Files:

```text
packages/npm/cli/package.json
packages/npm/cli/bin/scanrail.js
packages/npm/cli-darwin-arm64/package.json
packages/npm/cli-darwin-x64/package.json
packages/npm/cli-win32-x64/package.json
packages/npm/cli-win32-arm64/package.json
packages/npm/cli-linux-x64/package.json
packages/npm/cli-linux-arm64/package.json
```

Tasks:

1. Add wrapper package.
2. Add platform package manifests with `os` and `cpu`.
3. Add wrapper tests for platform resolution.
4. Add package dry-run script.
5. Copy release binaries into platform packages during release workflow.

Acceptance:

- `npm pack --dry-run` works for wrapper and platform packages.
- Wrapper forwards args and exit code.
- Unsupported platform error is explicit.
- No postinstall script is required.

## Phase 10: Release Automation

Goal:

- Build repeatable OSS release artifacts.

Files:

```text
.github/workflows/release.yml
.goreleaser.yml
scripts/release-check.sh
```

Tasks:

1. Cross-compile Go binaries.
2. Generate checksums.
3. Generate SBOM.
4. Build npm packages.
5. Run npm publish dry-run.
6. Publish archives and npm packages on tags after manual approval.

Acceptance:

- Release dry-run passes in CI.
- All artifacts include version metadata.
- npm wrapper is published after platform packages.
- Checksums are generated for archives.

## Phase 11: CI Templates and Examples

Goal:

- Make adoption easy for OSS users and organizations.

Files:

```text
examples/github-actions.yml
examples/gitlab-ci.yml
examples/scanrail.yaml
docs/setup-scenario.md
docs/cli-reference.md
```

Tasks:

1. Add GitHub Actions example using npm install.
2. Add GitLab CI example.
3. Add sample `scanrail.yaml`.
4. Update docs with real command output once implemented.

Acceptance:

- Example workflows use safe profile by default.
- Examples do not contain real company domains or secrets.
- Docs match implemented flags.

## MVP Cutline

Ship v0.1 when these are complete:

- Phase 0 through the first Phase 7 slice.
- `headers` and `gitleaks` adapters.
- JSON and HTML reports.
- npm wrapper dry-run.
- Linux/macOS/Windows test matrix.

Defer to v0.2:

- ZAP baseline if v0.1 scope slips.
- SARIF for URL-only findings.
- Nuclei, testssl.sh, Schemathesis.
- recorded browser session.
- Homebrew/Scoop.
- binary signing.

## Risks and Mitigations

| Risk | Mitigation |
| --- | --- |
| Docker Desktop path differences on Windows | Use `exec.Command` args, isolate mount construction, add Windows CI tests |
| Scanner output format changes | Golden fixtures and parser version metadata |
| npm platform package version skew | Release workflow publishes platform packages first and checks versions before wrapper |
| Active scan misuse | Safety preflight, explicit flag, production detection, profile gating |
| Secret leakage in reports | Central redactor, hard-fail on redaction errors, tests for sensitive headers |
| Large dependency surface | Keep Go dependencies minimal; npm wrapper has no runtime deps |
| OSS users expect full commercial scanner parity | Documentation states non-goals and safe scanner orchestration scope |

## Verification Gates

Every pull request:

```bash
go test ./...
go vet ./...
npm --prefix packages/npm/cli test
npm --prefix packages/npm/cli pack --dry-run
```

Before release:

```bash
goreleaser release --snapshot --clean
npm publish --dry-run --access public
scanrail doctor
scanrail init --non-interactive --project-name smoke --target http://127.0.0.1:8080 --profile quick
scanrail run --only headers
```

Manual release smoke:

- Install through npm globally on macOS.
- Install through npm globally on Windows.
- Run against a local intentionally insecure test server.
- Confirm report redaction.
- Confirm exit code behavior.

## Open Questions

- Final npm scope: `@scanrail/cli`, another available scope, or unscoped `scanrail`.
- Initial Go dependencies: standard library plus Cobra only, or include a prompt library in v0.1.
- Whether ZAP baseline belongs in v0.1 or v0.2.
- Signing timeline for Windows/macOS artifacts.
- Whether to include a Docker image distribution in v0.1.
