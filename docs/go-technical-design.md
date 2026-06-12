[ENGLISH](go-technical-design.md) | [한국어](go-technical-design.ko.md)

# Go Technical Design

## Purpose

This document translates the product architecture into a concrete Go implementation design. It assumes the OSS-first product shape in `README.md`, the component split in `docs/architecture.md`, the config schema in `docs/config-reference.md`, and the npm wrapper distribution model in `docs/distribution.md`.

## Design Principles

- The Go binary owns all product behavior.
- npm is a distribution and command shim layer only.
- Docker scanner execution is isolated behind a runner interface.
- Scanner integrations are plugins in code structure, not dynamic plugins in v0.1.
- Raw scanner outputs are preserved for debugging, but reports use normalized findings.
- Safety gates run before any network-active scanner.
- Secrets are referenced by environment variable names and redacted before persistence.
- Cross-platform behavior is tested on macOS, Windows, Linux, and CI.

## Repository Layout

```text
.
├─ cmd/
│  └─ scanrail/
│     └─ main.go
├─ internal/
│  ├─ app/
│  ├─ cli/
│  ├─ config/
│  ├─ doctor/
│  ├─ dockerx/
│  ├─ findings/
│  ├─ policy/
│  ├─ report/
│  ├─ safety/
│  ├─ scanners/
│  │  ├─ gitleaks/
│  │  ├─ headers/
│  │  ├─ semgrep/
│  │  ├─ trivy/
│  │  └─ zap/
│  ├─ setup/
│  └─ workspace/
├─ packages/
│  └─ npm/
│     ├─ cli/
│     ├─ cli-darwin-arm64/
│     ├─ cli-darwin-x64/
│     ├─ cli-win32-x64/
│     ├─ cli-win32-arm64/
│     ├─ cli-linux-x64/
│     └─ cli-linux-arm64/
├─ docs/
├─ examples/
└─ testdata/
```

Package intent:

| Package | Responsibility |
| --- | --- |
| `cmd/scanrail` | Process entrypoint only |
| `internal/cli` | Command definitions, flags, prompt entrypoints |
| `internal/app` | Command use cases and dependency wiring |
| `internal/config` | YAML schema, defaults, validation, merge |
| `internal/workspace` | Path resolution, `.scanrail` directories, OS path handling |
| `internal/dockerx` | Docker CLI runner abstraction |
| `internal/setup` | Image pull, cache, network, lock file |
| `internal/scanners` | Scanner registry and adapters |
| `internal/findings` | Common finding model and normalization helpers |
| `internal/safety` | allowlist, active scan, path/method, secret policy |
| `internal/policy` | fail-on severity and ignore rules |
| `internal/report` | HTML, JSON, SARIF, JUnit output |
| `internal/doctor` | Environment diagnostics |

## CLI Layer

Use `cobra` for command routing.

Commands:

```text
scanrail doctor
scanrail init
scanrail setup
scanrail run
scanrail update
scanrail ci init
scanrail auth setup
scanrail report open
```

Command handlers should call `internal/app` services rather than embedding business logic. This keeps command parsing thin and makes behavior easier to test.

Suggested shape:

```go
type RootOptions struct {
    ConfigPath string
    Workdir    string
    Verbose    bool
}

type RunOptions struct {
    Profile                  string
    Target                   string
    OpenAPI                  string
    OutputDir                string
    Formats                  []string
    FailOn                   string
    UnderstandsActiveScan    bool
    OnlyTools                []string
    SkipTools                []string
}
```

Interactive prompts:

- Use a small prompt interface owned by `internal/cli`.
- Implement it with a terminal prompt library in the real CLI.
- Use a fake implementation in tests.

```go
type Prompter interface {
    Text(label string, def string) (string, error)
    Select(label string, options []Choice) (string, error)
    MultiSelect(label string, options []Choice) ([]string, error)
    Confirm(label string, def bool) (bool, error)
}
```

## Config Model

`scanrail.yaml` is decoded into typed Go structs and validated before use.

```go
type Config struct {
    Project ProjectConfig `yaml:"project"`
    Targets TargetConfig  `yaml:"targets"`
    Auth    AuthConfig    `yaml:"auth"`
    Profiles Profiles     `yaml:"profiles"`
    Safety  SafetyConfig  `yaml:"safety"`
    Policy  PolicyConfig  `yaml:"policy"`
    Report  ReportConfig  `yaml:"report"`
}
```

Merge order:

```text
CLI options
Environment variables
scanrail.yaml
Organization defaults (reserved after v0.1)
Built-in defaults
```

Validation rules:

- `project.name` is required.
- At least one target is required for network scanners.
- `auth` may store env var names, not secret values.
- `full` profile with active scanners must require explicit opt-in.
- Web/API scanners require allowlist unless a safe local-only mode is selected.
- Unknown tools fail validation unless `--allow-unknown-tools` is introduced later.
- Unsupported reserved auth types fail validation until implemented.
- v0.1 does not load organization defaults; the merge layer is reserved for a later policy feature.

## Workspace Model

The workspace package owns filesystem paths.

Default paths:

```text
.scanrail/
├─ cache/
├─ raw/
├─ reports/
└─ tmp/
```

Rules:

- Resolve paths with `filepath.Abs`.
- Use `filepath` for host paths and POSIX paths only inside containers.
- Keep all generated artifacts under `.scanrail` by default.
- Never write scanner raw output outside the workspace output directory.

Windows handling:

- MVP supports native Windows with Docker Desktop.
- Use `exec.Command` argument arrays, not shell strings.
- Convert host paths only at the Docker boundary.
- WSL path conversion is a follow-up unless explicitly detected and tested.

Localhost target handling:

- Host-side `http://localhost:<port>` and `http://127.0.0.1:<port>` work for native Go scanners.
- Docker-backed scanners must not receive container-local `localhost` for host services.
- On macOS and Windows Docker Desktop, convert host-local targets to `host.docker.internal`.
- On Linux, prefer a documented `--docker-host-gateway` mode that maps `host.docker.internal` to the host gateway, or require an explicit reachable target URL.
- The target resolver records both user-facing URL and scanner-facing URL in run metadata.

## Docker Runner

Use Docker CLI wrapping for v0.1. The Docker Go SDK can be introduced later if a real need appears.

Reasons:

- Smaller initial dependency surface.
- Behavior matches user-visible Docker commands.
- Easier to debug by printing equivalent commands.
- Easier to fake in unit tests.

Interface:

```go
type Runner interface {
    Version(ctx context.Context) (*VersionInfo, error)
    Pull(ctx context.Context, image string) error
    ImageInspect(ctx context.Context, image string) (*ImageInfo, error)
    NetworkEnsure(ctx context.Context, name string) error
    Run(ctx context.Context, req RunRequest) (*RunResult, error)
}
```

`RunRequest`:

```go
type RunRequest struct {
    Image       string
    Name        string
    Workdir     string
    Env         map[string]string
    Volumes     []Volume
    Network     string
    Args        []string
    Timeout     time.Duration
    StdoutPath  string
    StderrPath  string
}
```

Execution rules:

- Do not invoke through a shell.
- Capture stdout/stderr to raw files.
- Redact command previews before logging.
- Return scanner exit code separately from runner infrastructure errors.
- Apply timeouts through `context.Context`.

## Scanner Adapter Contract

Each scanner adapter has three responsibilities: prepare, run, normalize.

```go
type Scanner interface {
    ID() ID
    Metadata() Metadata
    RequiredTargets() []TargetKind
    Intrusiveness() Intrusiveness
    SafetyCapabilities() SafetyCapabilities
    Validate(ctx context.Context, cfg config.Config) error
    Prepare(ctx context.Context, env Environment) error
    Run(ctx context.Context, env Environment) (*RawResult, error)
    Normalize(ctx context.Context, raw *RawResult) ([]findings.Finding, error)
}
```

Safety capability contract:

```go
type Intrusiveness string

const (
    Passive     Intrusiveness = "passive"
    Interactive Intrusiveness = "interactive"
    Active      Intrusiveness = "active"
)

type SafetyCapabilities struct {
    AllowlistScope  bool
    RedirectScope   bool
    BlockedPaths    bool
    BlockedMethods  bool
    RateLimit       bool
    HeaderInjection bool
    AuthInjection   bool
}
```

The orchestrator must compare profile requirements with scanner capabilities before execution. If a scanner cannot provide a required guarantee, profile execution skips it with a report entry; explicit `--only <tool>` execution fails with a safety error.

`Environment` contains:

- resolved config
- workspace paths
- Docker runner
- logger
- redactor
- auth material references
- scanner image lock data

Initial adapters:

| Adapter | v0.1 role |
| --- | --- |
| `headers` | Native Go HTTP header checks, no Docker needed |
| `gitleaks` | Docker scanner for secrets |
| `trivy` | Docker scanner for filesystem dependency/IaC scan |
| `semgrep` | Docker scanner for SAST |

v0.2 adapters:

- `zap-baseline`
- `nuclei-safe`
- `zap-api`
- `testssl`
- `schemathesis`

## Scan Orchestrator

The orchestrator converts a profile into executable scanner tasks.

Flow:

```text
load config
apply CLI overrides
validate config
select profile
resolve scanners
run safety preflight
filter scanners by target availability and safety capabilities
prepare workspace
prepare scanners
execute scanner graph
normalize raw outputs
dedupe findings
apply policy
write reports
return exit code
```

Concurrency:

- Code-only scanners can run in parallel.
- Network-active scanners default to sequential execution in v0.1.
- Add bounded worker pools only after rate limits are proven.

Failure policy:

- Config errors stop before any scanner.
- Safety violations return exit code `5`.
- Scanner infrastructure errors return exit code `4`.
- Scanner findings that violate policy return exit code `1`.
- Successful scans return exit code `0`.
- Environment errors return exit code `3`.
- Usage and config errors return exit code `2`.

Target handling:

- A profile-selected scanner with missing required target is skipped and recorded as `SKIP`.
- An explicitly selected scanner through `--only` with missing target fails with exit code `2`.
- A scanner missing required safety capability fails with exit code `5` when explicit, and is skipped when inherited from a profile.

## Finding Model

```go
type Finding struct {
    ID             string
    Fingerprint    string
    Title          string
    Description    string
    Severity       Severity
    Confidence     Confidence
    Source         Source
    Target         Target
    Location       *Location
    Evidence       []Evidence
    Classification Classification
    Remediation    string
    References     []Reference
}
```

Severity:

```text
critical
high
medium
low
info
```

Fingerprint inputs:

- source tool
- rule id
- target URL/file/image
- normalized location
- normalized evidence key

The fingerprint must be stable across runs so suppressions and baselines can work.

Cross-tool dedupe in later versions must use a separate merge key. The per-tool fingerprint intentionally includes source tool so scanner-specific suppressions remain stable; cross-tool merging should use normalized CWE/category, target, location, and evidence shape.

## Safety Gates

Safety preflight runs before scanner execution.

Checks:

- target host is in allowlist
- redirect policy is safe
- active scan requires explicit flag
- production target blocks intrusive tools
- blocked paths are present for network profiles
- auth env vars exist when required
- scanner request headers are configured when available

Redaction:

- Central `Redactor` masks headers, cookies, tokens, passwords, and configured env values.
- Redaction runs before logs, raw command preview files, and reports.
- Redaction failure is a hard failure.

## Reports

Report writers consume normalized findings and run metadata.

Writers:

- JSON: full internal model
- HTML: developer readable report
- SARIF: code scanning integration
- JUnit: CI summary follow-up

HTML v0.1 should be static and dependency-light:

- `html/template`
- embedded CSS
- no external assets
- grouped by severity and tool
- finding details include evidence, affected target, and remediation

SARIF v0.1 scope:

- file-based findings from Semgrep, Trivy, and Gitleaks
- URL-only findings mapped to logical locations later

## Logging and Telemetry

Use Go `log/slog`.

Default:

- no network telemetry
- local logs only
- secret redaction on all structured fields

Future telemetry must be opt-in and documented.

## Release Design

Build matrix:

```text
darwin/arm64
darwin/amd64
windows/amd64
windows/arm64
linux/amd64
linux/arm64
```

Artifacts:

```text
scanrail_darwin_arm64.tar.gz
scanrail_darwin_x64.tar.gz
scanrail_win32_x64.zip
scanrail_win32_arm64.zip
scanrail_linux_x64.tar.gz
scanrail_linux_arm64.tar.gz
checksums.txt
sbom.spdx.json
```

npm publish:

- publish platform packages first
- publish wrapper package last
- use provenance/trusted publishing where available
- keep versions synchronized

## Testing Strategy

Unit tests:

- config decode/default/validation
- safety allowlist and active scan policy
- Docker command construction
- scanner normalization
- finding fingerprinting
- policy exit code calculation
- redaction

Integration tests:

- fake Docker runner for scanner command requests
- golden raw scanner output normalization
- HTML/JSON/SARIF golden reports
- npm wrapper resolves mock platform package

End-to-end tests:

- `scanrail doctor`
- `scanrail init --non-interactive`
- `scanrail setup --pull-policy never` with fake runner
- `scanrail run --only headers` against a local test server
- CI matrix on macOS, Windows, Linux

Manual smoke tests before release:

- npm global install
- npx execution
- Docker Desktop on macOS
- Docker Desktop on Windows
- Linux Docker engine in CI

## Deferred Decisions

- Docker SDK adoption.
- Dynamic plugin API.
- Windows and macOS binary signing timing.
- Central dashboard storage model.
- WSL-native path support.
- Advanced authentication session recording.
