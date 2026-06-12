[ENGLISH](oss-strategy.md) | [한국어](oss-strategy.ko.md)

# OSS Strategy

Scanrail is intended as a public OSS project, not a company-internal utility. It targets developers and teams that cannot easily adopt commercial security testing products but still need a practical baseline.

## Positioning

One-line description:

```text
A developer-friendly security scan orchestrator that runs Docker-backed open-source scanners and normalizes their results into useful reports.
```

Key messages:

- It is a developer self-check tool, not only a security-expert tool.
- It does not claim to replace every commercial DAST/SAST capability.
- It combines proven open-source scanners instead of rebuilding every detection engine.
- It defaults to safe behavior and requires explicit opt-in for risky scans.
- Organization policy can be layered through configuration and future plugins.

## Target Users

- development teams blocked from commercial security product approval
- OSS projects that want security self-checks before PRs
- platform teams that want to standardize multiple scanners in CI
- security engineers who want HTML/JSON/SARIF output from one flow
- backend developers who want safer repeated checks against staging APIs

## Public Core and Organization Extensions

Included in the public core:

- CLI
- interactive and non-interactive setup
- scanner adapter framework
- common finding model
- HTML/JSON/SARIF report plan
- safe default profiles
- public documentation and examples
- CI templates

Provided externally by organizations:

- internal registry mirrors
- standard organization profiles
- private security policies
- private Semgrep/Nuclei rules
- Slack/Jira endpoints
- asset criticality data
- internal allowlist templates

## License

Recommended license:

```text
Apache-2.0
```

Rationale:

- friendly to company adoption
- includes patent language
- fits CLI, SDK, and plugin ecosystems

Alternatives:

- MIT: simple, but weaker patent protection
- MPL-2.0: useful when file-level copyleft is desired
- AGPL-3.0: strong SaaS restriction, but high adoption friction

The initial repository uses Apache-2.0.

## Repository Principles

- Keep the public core useful without private infrastructure.
- Keep defaults safe enough for public CI.
- Do not store secrets in configuration.
- Make scanner versions and reports reproducible.
- Document unsupported capabilities honestly.
- Prefer small, explicit adapters over hidden behavior.

## Community Roadmap

Early community work should focus on:

- adapter quality
- report usability
- safe defaults
- CI examples
- Windows/macOS/Linux compatibility
- npm package reliability

Later work can add:

- central dashboard integrations
- plugin contracts
- organization policy bundles
- richer SARIF and issue tracker integration
