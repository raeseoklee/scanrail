[ENGLISH](architecture.md) | [한국어](architecture.ko.md)

# Architecture

## Overview

Scanrail is an orchestrator, not a scanner engine. It runs scanner adapters, gathers outputs, normalizes findings, applies policy, and writes reports.

```text
Developer / CI
      |
      v
 scanrail CLI
      |
      +-- config loader
      +-- setup manager
      +-- scan orchestrator
      +-- auth/session manager
      +-- scanner adapters
      +-- finding normalizer
      +-- policy engine
      +-- report generator
```

## Components

### CLI

The CLI is the user-facing entrypoint.

```text
scanrail doctor
scanrail init
scanrail setup
scanrail run
scanrail update
scanrail ci init
```

The current MVP implements `doctor`, `init`, `setup`, `run`, and `version`.

### Config Loader

The config loader merges values from CLI flags, environment variables, `scanrail.yaml`, organization defaults, and built-in defaults.

Precedence:

1. CLI option
2. Environment variable
3. `scanrail.yaml`
4. Organization default
5. Built-in default

The public v0.1 candidate implements a small typed config surface and keeps organization defaults as a future layer.

### Setup Manager

The setup manager prepares local Scanrail workspace state and, in later adapter phases, Docker-backed scanner assets.

Responsibilities:

- check Docker when needed
- create `.scanrail` directories
- prepare cache and report paths
- manage scanner version lock data
- avoid pulling placeholder scanner images

### Scan Orchestrator

The orchestrator selects scanners from a profile, validates target availability, checks scanner safety capabilities, executes adapters, and records evidence.

Rules:

- independent scanners may run in parallel in future phases
- active or interactive scanners require stronger safety capabilities
- missing profile-selected targets are skipped with evidence
- explicitly selected scanners fail when their required target is missing
- safety violations return exit code `5`

### Auth and Session Manager

Authentication is represented by references, not raw secrets. Tokens should be read from environment variables and redacted before any persistence.

Supported future modes:

- none
- bearer token
- cookie
- custom header

### Scanner Adapters

Each scanner integration is isolated behind an adapter contract.

Adapter responsibilities:

- declare required targets
- declare safety capabilities
- build commands or native checks
- execute scanner logic
- preserve raw output where useful
- normalize findings into the common model

### Finding Normalizer

The normalizer maps scanner-specific output into a shared finding shape:

- id
- scanner
- title
- severity
- category
- target
- evidence
- remediation
- references

### Policy Engine

The policy engine decides whether findings should fail the command. It supports severity thresholds and planned ignore rules with reasons and expiration dates.

### Report Generator

Reports are written to `.scanrail/reports` by default.

Initial formats:

- JSON
- HTML

Planned formats:

- SARIF
- JUnit XML

## Data Flow

```text
scanrail.yaml
    |
    v
config loader -> safety validation -> scanner adapters
    |                                      |
    |                                      v
    +------------------------------ normalized findings
                                           |
                                           v
                              policy + reports + exit code
```

## Safety Boundary

Scanrail does not assume every scanner can enforce every safety policy. Each adapter declares what it can enforce. If a profile requires stronger guarantees than the adapter provides, Scanrail skips the scanner or fails explicit execution.

Network-level enforcement through a dedicated proxy and Docker network is a later v0.x design.
