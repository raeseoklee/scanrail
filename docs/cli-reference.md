[ENGLISH](cli-reference.md) | [한국어](cli-reference.ko.md)

# CLI Reference

This document defines the initial `scanrail` CLI shape.

## Commands

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

The first release candidate implements `doctor`, `init`, `setup`, `run`, and `version`.

## Common Exit Codes

```text
0    success
1    policy failure or finding threshold failure
2    usage or configuration error
3    runtime environment error
4    scanner execution or parsing failure
5    safety policy violation
130  interrupted by user
```

`doctor` does not produce findings. Missing runtime prerequisites return `3`; invalid configuration returns `2`.

## `scanrail doctor`

Checks whether the local environment can run Scanrail.

```bash
scanrail doctor
```

Checks include:

- Docker command availability
- Docker daemon availability when needed
- current directory access
- Git repository status
- disk space
- scanner image pull availability in future adapter phases
- `scanrail.yaml` validity
- `tools.lock.yaml` validity

Exit codes:

```text
0  ready
3  prerequisite missing
2  configuration error
```

## `scanrail init`

Creates `scanrail.yaml`.

```bash
scanrail init
```

Options:

```text
--non-interactive       create config from defaults and supplied flags
--project-name <name>   project name
--target <url>          web target URL
--openapi <path-or-url> OpenAPI spec path or URL
--profile <name>        default profile
--force                 overwrite existing scanrail.yaml
```

Safety behavior:

- never overwrite existing config unless `--force` is set
- never write raw secrets
- infer allowlist from the target hostname
- default to a safe `quick` profile

## `scanrail setup`

Prepares local workspace state and scanner assets.

```bash
scanrail setup
```

Options:

```text
--pull-policy <policy>  missing, always, or never
```

The first release candidate supports `--pull-policy never` for offline/dry-run setup. Docker image pulls are skipped for placeholder scanner images until real adapters are implemented.

## `scanrail run`

Runs a scan profile or a selected scanner.

```bash
scanrail run --profile quick
scanrail run --only headers
scanrail run --target https://staging.example.com
```

Options:

```text
--profile <name>        profile name, default quick
--only <scanner>        run one scanner explicitly
--target <url>          override web target URL
--output-dir <path>     report output directory
```

Target behavior:

- If a profile-selected scanner has no required target, it is skipped with evidence.
- If a scanner was selected with `--only` and its target is missing, the command fails.
- Safety capability mismatches fail explicit execution and are recorded for profile execution.

## `scanrail version`

Prints version metadata.

```bash
scanrail version
scanrail --version
```

Build metadata is injected through Go linker flags during release builds.

## Planned Commands

### `scanrail update`

Updates scanner metadata, templates, or lock files.

### `scanrail ci init`

Writes CI templates for GitHub Actions, GitLab CI, or other supported systems.

### `scanrail auth setup`

Guides users through safe authentication references without storing raw tokens.

### `scanrail report open`

Opens the latest HTML report from `.scanrail/reports`.
