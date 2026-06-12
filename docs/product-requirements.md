[ENGLISH](product-requirements.md) | [한국어](product-requirements.ko.md)

# Product Requirements

## Background

Many teams cannot adopt commercial security testing products quickly. Approval, budget, procurement, or internal security review can take longer than the delivery cycle. Those teams still need a repeatable way to run baseline security checks before pull requests, releases, and staging deployments.

Individual open-source scanners are powerful, but each one has its own installation model, flags, output format, false-positive workflow, and safety model. Scanrail provides one developer-facing interface around those tools.

## Goals

- Let developers run security self-checks before PRs and deployments.
- Automatically prepare open-source scanners through a reproducible setup flow.
- Generate project-specific configuration through an interactive or non-interactive flow.
- Normalize scanner outputs into a common finding model.
- Explain risk with OWASP, CWE, CVSS, EPSS, and CISA KEV context where available.
- Produce HTML, JSON, and later SARIF reports.
- Run cleanly in CI/CD.
- Allow organizations to layer internal policies, mirrors, and private rules on top of the public core.

## Non-Goals

- Do not claim full parity with commercial SAST or DAST suites.
- Do not build SQL injection, XSS, or auth-bypass engines from scratch.
- Do not support unauthorized scanning of external targets.
- Do not enable active production scanning by default.
- Do not claim complete automated detection of business-logic vulnerabilities.

## Primary Users

### Developers

Developers run checks locally or in CI, inspect reports, and fix findings without needing deep scanner-specific knowledge.

### Security Engineers

Security engineers define default profiles, policies, rules, exceptions, and reporting expectations.

### Platform and DevOps Teams

Platform teams maintain CI templates, Docker image mirrors, scanner versions, and cache policies.

### Open Source Maintainers

Maintainers add a safe security baseline to public projects and publish machine-readable results for code scanning integrations.

## Core Capabilities

### 1. Initial Configuration

`scanrail init` creates `scanrail.yaml` from project metadata and scan targets.

Collected inputs include:

- project name
- scan target types
- staging URL
- OpenAPI or Swagger path/URL
- container image name
- authentication type
- allowlisted domains
- excluded or blocked paths
- default profile
- CI failure threshold

### 2. Scanner Setup

`scanrail setup` prepares the local execution environment.

Responsibilities:

- check Docker availability when Docker-backed tools are selected
- create workspace directories
- prepare scanner cache directories
- pull or validate scanner images in future adapter phases
- write or validate a scanner lock file

### 3. Scan Execution

`scanrail run` executes the selected profile and writes normalized reports.

The first release candidate implements the native security headers scanner. Docker-backed adapters for Gitleaks, Trivy, and Semgrep are part of the planned v0.x surface.

### 4. Reporting

Scanrail writes developer-readable HTML and machine-readable JSON. SARIF is part of the broader release plan.

Reports should include:

- findings grouped by severity
- scanner evidence
- remediation guidance
- skipped scanners and skip reasons
- policy result and exit code

### 5. CI Integration

CI usage must be deterministic:

- stable exit codes
- no interactive prompts
- reproducible scanner versions
- reports written to known paths
- policy failures separated from setup/runtime failures

## Success Metrics

- A new project can generate its first report in under ten minutes.
- Developers can understand high-risk findings without reading scanner manuals.
- CI failures clearly distinguish findings, configuration errors, runtime errors, and safety violations.
- Secrets are not written to config, logs, raw outputs, or reports.
- Organizations can adopt the OSS core without forking it for basic policy needs.
