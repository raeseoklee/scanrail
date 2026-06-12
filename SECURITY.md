[ENGLISH](SECURITY.md) | [한국어](SECURITY.ko.md)

# Security Policy

Scanrail is a security tool, so security reports need a careful disclosure path.

## Supported Versions

Scanrail is pre-1.0. Only the latest published npm version and the current `main` branch receive security fixes.

## Reporting a Vulnerability

Please do not open a public GitHub issue with exploit details, live target URLs, raw secrets, or private scan output.

Preferred reporting path:

1. Use GitHub private vulnerability reporting for this repository if available.
2. If private reporting is unavailable, open a minimal public issue asking for a private disclosure channel. Do not include exploit details.

Include:

- affected Scanrail version
- operating system and install method
- impact summary
- reproduction steps with synthetic targets or sanitized fixtures
- whether the issue can cause unsafe scanning, secret exposure, command execution, or report tampering

## Scope

In scope:

- command injection or arbitrary process execution
- target allowlist bypass
- secret leakage through config, reports, logs, MCP tools, or npm wrappers
- unsafe default behavior that scans targets without explicit user intent
- package integrity or provenance problems

Out of scope:

- vulnerabilities in third-party scanners unless Scanrail makes them worse
- findings produced by Scanrail against unrelated applications
- denial-of-service against a local demo target without a Scanrail-specific impact

## Disclosure

We will try to acknowledge valid reports within 7 days and coordinate a fix or mitigation before public disclosure.
