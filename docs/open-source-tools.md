[ENGLISH](open-source-tools.md) | [한국어](open-source-tools.ko.md)

# Open Source Tool Review

Scanrail combines specialized open-source scanners and presents their results through one developer workflow. Each tool keeps its own strength; Scanrail absorbs setup, execution, normalization, and reporting friction.

## Primary Candidates

| Area | Tool | Purpose | Phase |
| --- | --- | --- | --- |
| Headers | native Go checker | security header checks | v0.1 |
| SAST | Semgrep | code vulnerability and pattern checks | v0.x |
| Dependencies, containers, IaC | Trivy | CVE, container, IaC, SBOM checks | v0.x |
| Secrets | Gitleaks | hardcoded token/key detection | v0.1 |
| DAST | OWASP ZAP | runtime web vulnerability checks | v0.2 |
| Template scanning | Nuclei | CVE, misconfiguration, exposure checks | v0.2 |
| TLS | testssl.sh | TLS and certificate checks | v0.2 |
| API testing | Schemathesis | OpenAPI-based property and fuzz testing | v0.2 |
| Result integration | SARIF | code scanning integration | v0.1 partial |

The current MVP executes the native headers scanner and the Docker-backed Gitleaks adapter. Other scanner names remain part of the planned adapter surface.

## Additional Candidates

| Area | Tool | Purpose | Notes |
| --- | --- | --- | --- |
| SAST | CodeQL | deep code query analysis | depends on language and CI environment |
| Web server scan | Nikto | server configuration and known files | false-positive control required |
| Service discovery | Nmap | port and service discovery | strict scope control required |
| Attack surface | OWASP Amass | subdomain and exposure discovery | organization policy required |
| Vulnerability management | DefectDojo | finding workflow and history | useful at central-server stage |
| SBOM management | Dependency-Track | SBOM supply-chain risk | better for organization-level operation |

## Tool Roles

### OWASP ZAP

Role:

- baseline scan
- passive scan
- OpenAPI-based API scan
- full active scan

Initial policy:

- use baseline/API scans before active scans
- require explicit flags and allowlists for active behavior
- treat intrusive behavior as at least `interactive` or `active`

### Nuclei

Role:

- template-based checks
- known CVE and misconfiguration detection
- organization-specific template extension

Initial policy:

- default to safe templates
- exclude destructive or intrusive templates by default
- lock template versions or use organization mirrors

### Semgrep

Role:

- code-level vulnerability detection
- language-specific security rules
- organization coding policy rules

Initial policy:

- start with community rules
- support organization rules as a separate namespace later
- collect JSON and SARIF output

### Trivy

Role:

- dependency CVE checks
- container image checks
- IaC checks
- SBOM generation

Initial policy:

- pin image versions
- cache vulnerability databases
- clearly separate dependency and container findings

### Gitleaks

Role:

- detect hardcoded secrets
- scan Git history or working tree depending on profile

Initial policy:

- never print raw secret values in reports
- redact evidence
- mount the workspace read-only
- pin the Docker image to `ghcr.io/gitleaks/gitleaks:v8.30.1`
- support allowlisted false positives with expiration dates

### Native Headers Checker

Role:

- check common HTTP response security headers
- provide a zero-Docker first scanner
- validate the reporting pipeline

Initial checks:

- Content-Security-Policy
- X-Content-Type-Options
- X-Frame-Options
- Referrer-Policy
- Strict-Transport-Security for HTTPS targets

## Adapter Requirements

Every scanner adapter must declare:

- required target type
- intrusiveness level
- safety capabilities
- supported output formats
- normalization strategy
- raw evidence handling
- skip/failure behavior

This contract prevents the profile from promising safety guarantees that a scanner cannot enforce.
