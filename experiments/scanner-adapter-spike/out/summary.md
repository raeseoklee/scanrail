# Scanner Adapter Spike Summary

Generated: 2026-06-12T10:13:52.598Z

## Commands

- gitleaks: exit 7, image `ghcr.io/gitleaks/gitleaks:v8.30.1`
- semgrep: exit 0, image `semgrep/semgrep:1.165.0`
- trivy: exit 0, image `aquasec/trivy:0.71.0`

## Findings

- [high] gitleaks:scanrail-fake-secret at app.py:3
- [high] semgrep:scanrail.python.eval at app.py:7
- [medium] semgrep:scanrail.python.shell-true at app.py:11
- [high] trivy:AWS-0086 at main.tf:8
- [high] trivy:AWS-0087 at main.tf:9
- [low] trivy:AWS-0089 at main.tf:1
- [medium] trivy:AWS-0090 at main.tf:1
- [high] trivy:AWS-0091 at main.tf:10
- [high] trivy:AWS-0093 at main.tf:11
- [high] trivy:AWS-0132 at main.tf:1
