[ENGLISH](scanner-adapter-spike.md) | [한국어](scanner-adapter-spike.ko.md)

# Scanner Adapter Spike

This experiment proves the next Scanrail implementation step before adding production adapter code. It runs three open-source scanners through Docker, captures their raw JSON outputs, normalizes findings into one report shape, and records a short terminal demo with VHS.

## Why This Exists

The product direction is multi-scanner orchestration. This spike was recorded before the first Gitleaks production adapter and validates the real command behavior, output shape, exit codes, path mounts, and redaction needs for Gitleaks, Trivy, and Semgrep.

## Tools

| Tool | Image | Purpose in spike |
| --- | --- | --- |
| Gitleaks | `ghcr.io/gitleaks/gitleaks:v8.30.1` | fake secret detection through a custom rule |
| Semgrep | `semgrep/semgrep:1.165.0` | Python insecure-code findings |
| Trivy | `aquasec/trivy:0.71.0` | Terraform misconfiguration findings |

## Run

```bash
npm run experiment:scanner-spike
```

Outputs:

```text
experiments/scanner-adapter-spike/out/raw/
experiments/scanner-adapter-spike/out/scanrail-normalized-report.json
experiments/scanner-adapter-spike/out/scanrail-normalized-report.html
experiments/scanner-adapter-spike/out/summary.md
```

## Summary Demo

```bash
npm run experiment:scanner-spike:summary
```

## Recording

The tape file is:

```text
experiments/scanner-adapter-spike/tapes/scanner-adapter-spike.tape
```

Render the GIF:

```bash
make tape-scanner-spike
```

The recording intentionally shows the summary command instead of the full Docker scan so the demo stays short and deterministic. Run `npm run experiment:scanner-spike` to regenerate the full evidence set.

Generated asset:

![Scanner adapter spike](../assets/scanner-adapter-spike.gif)

## Findings From The Spike

- Docker image execution works on the tested Docker environment.
- Gitleaks can be integrated with a custom config and a non-default expected finding exit code.
- Semgrep JSON contains enough rule, path, severity, and message data for the initial finding model.
- Trivy misconfiguration output contains enough ID, title, severity, location, and remediation data for normalization.
- Raw outputs should be preserved, but report evidence must be redacted before being surfaced.
- Adapter implementations should treat scanner exit codes as scanner-specific contracts, not as a single global rule.

## Implementation Implications

The next production implementation should add:

- a Docker runner abstraction
- scanner-specific expected exit-code handling
- raw artifact capture under `.scanrail/raw`
- finding normalization per adapter
- central redaction before reports
- adapter capability metadata for safety gates
