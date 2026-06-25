[ENGLISH](roadmap.md) | [한국어](roadmap.ko.md)

# Roadmap

## v0.1: Local CLI MVP

Goal:

- Let developers run a basic security check with a local CLI and safe defaults.

Features:

- `scanrail doctor`
- `scanrail init`
- `scanrail setup`
- `scanrail run`
- `scanrail.yaml`
- local report directory
- HTML and JSON reports
- npm wrapper package topology
- platform binary packages

Tools:

- native security headers checker
- Docker-backed Gitleaks secrets adapter
- planned adapter surfaces for Trivy and Semgrep

Completion criteria:

- a new project can produce a first report quickly
- exit codes distinguish policy, config, runtime, scanner, and safety failures
- raw secrets are not exposed in reports
- npm wrapper works on macOS, Windows, and Linux

## v0.2: API and Web Scan Expansion

Goal:

- Improve staging URL and OpenAPI-based scan quality.

Features:

- native TLS certificate baseline (`0.2.0` shipped)
- OWASP ZAP baseline
- ZAP OpenAPI scan
- Nuclei safe templates
- Schemathesis
- testssl.sh
- bearer/cookie authentication references
- allowlist validation
- blocked paths
- rate limits

Completion criteria:

- authenticated staging API checks are possible
- OpenAPI scan results appear in HTML reports
- active scans do not run without explicit opt-in

Status note: `0.2.0` ships the native TLS baseline only. ZAP, Nuclei, Schemathesis, testssl.sh, and OpenAPI scanning remain planned under this milestone.

## v0.3: CI/CD Integration

Goal:

- Run Scanrail automatically on PRs and main branches.

Features:

- `scanrail ci init`
- GitHub Actions template
- GitLab CI template
- Jenkins example
- SARIF upload
- JUnit XML output
- cache optimization
- release provenance

Completion criteria:

- quick profile runs on every PR
- stronger profile can run after merge or on schedule
- reports are easy to archive as CI artifacts

## v0.4: Policy and Workflow

Goal:

- Make findings manageable for teams.

Features:

- ignore rules with reasons and expiration
- policy packs
- organization defaults
- severity overrides
- Jira/GitHub issue export
- Slack/webhook notifications

Completion criteria:

- accepted risk is auditable
- stale ignores are visible
- platform teams can provide shared policy without forking the tool

## v0.5: Agent Integration

Goal:

- Let AI coding assistants use Scanrail safely without bypassing the CLI safety model.
- The first stdio MCP slice shipped early in `0.1.2`; this milestone tracks hardening and broader client compatibility.

Features:

- local MCP server over stdio
- read-only resources for configuration, report summaries, and schemas
- bounded tools for `doctor`, safe setup, native headers scans, and report summarization
- explicit target allowlist enforcement
- no secret values in MCP payloads

Completion criteria:

- MCP clients can inspect the current Scanrail state
- tool calls cannot execute arbitrary shell commands
- active or network-heavy scans require explicit user opt-in through Scanrail policy

## v1.0: Stable OSS Release

Goal:

- Provide a stable developer security scan orchestrator with clear guarantees.

Requirements:

- stable config schema
- stable finding model
- stable exit codes
- documented adapter contracts
- reproducible release process
- signed/provenance-aware release artifacts
- mature safety documentation

Non-requirements:

- commercial scanner parity
- central SaaS dashboard
- automatic detection of every business-logic vulnerability
