[ENGLISH](risk-treatment-plan.md) | [한국어](risk-treatment-plan.ko.md)

# 리스크 처리 계획

이 문서는 release risk register를 실제 운영 gate로 바꿉니다. [릴리스 리스크 레지스터](release-risk-register.ko.md)를 대체하지 않고, 남은 리스크를 release 또는 기능 확장 전에 어떻게 다룰지 설명합니다.

최종 검토일: 2026년 6월 16일.

## 현재 결정

release checklist와 verification gate가 통과한다면 현재 headers + Gitleaks MVP에는 알려진 release blocker가 없습니다.

남은 리스크는 기능 확장 gate로 관리합니다.

- production host compatibility와 audit 동작을 검증하기 전 MCP execution surface를 넓히지 않습니다.
- 각 adapter마다 adapter isolation, redaction, raw artifact handling, version metadata가 구현되기 전 추가 Docker-backed scanner를 활성화하지 않습니다.
- SARIF와 artifact workflow가 ship되기 전 Scanrail을 PR-native CI security gate로 positioning하지 않습니다.
- license metadata 없이 third-party scanner rule pack을 bundle하지 않습니다.

## 릴리스 Gate

다음 리스크는 모든 npm publish 전에 확인해야 합니다.

| 리스크 | Gate | 증거 |
| --- | --- | --- |
| R-001 npm trusted publishing 설정 drift | npm trusted publisher 설정을 `.github/workflows/npm-publish.yml`과 대조합니다. | 작성 완료된 [Release Checklist](release-checklist.ko.md)의 trusted-publisher 섹션 |
| R-002 partial publish로 package set 불일치 발생 | publish 전후 package version을 기록하고 platform package를 wrapper보다 먼저 publish합니다. | 작성 완료된 [Release Checklist](release-checklist.ko.md)의 registry snapshot 섹션 |
| R-004 public npm smoke가 per-commit이 아님 | publish 후 public package version에 대해 scheduled/manual npm smoke를 실행합니다. | `npm Smoke` workflow URL과 결과 |
| R-006 MCP tool이 core policy를 우회할 가능성 | MCP 코드 변경이 있으면 publish 전 MCP Workbench regression을 실행합니다. | `mcp-workbench run examples/mcp-workbench/scanrail-mcp.yaml --verbose` |

## MCP 확장 Gate

`scanrail_setup`, 더 넓은 scanner execution, richer report resource를 MCP에 추가하기 전 다음 gate가 필요합니다.

| 리스크 | Gate | 확장 전 필요 조건 |
| --- | --- | --- |
| R-005 production MCP host compatibility | 최소 1개 production MCP host configuration을 검증합니다. | host name, config snippet, smoke result, 날짜 |
| R-007 sensitive 또는 oversized MCP resource | richer resource마다 size limit과 redaction test를 추가합니다. | maximum payload size와 secret masking test |
| R-008 MCP auditability | `.scanrail/logs/mcp-audit.jsonl` event를 새 setup/scanner tool까지 확장합니다. | denied, started, completed, failed path test |

## Scanner Adapter Gate

Trivy, Semgrep 또는 future Docker-backed adapter를 skipped scaffold에서 실행 상태로 승격하기 전 다음 gate가 필요합니다. Gitleaks는 `0.1.4`에서 첫 gate slice를 통과했지만, OS matrix와 digest-level image policy는 아직 열려 있습니다.

| 리스크 | Gate | 활성화 전 필요 조건 |
| --- | --- | --- |
| R-009 scanner adapter exposure | capability contract를 catalog metadata뿐 아니라 runner에서 강제합니다. | allowlist, redirect, credential, raw output adapter test |
| R-010 Docker host variance | Linux adapter integration을 먼저 실행하고 macOS/Windows Docker behavior를 문서화합니다. | OS compatibility note 또는 CI matrix |
| R-011 scanner image/rule drift | scanner image와 rule version을 pinning합니다. | report metadata에 scanner name, version, 가능한 경우 image digest 기록 |
| R-012 normalization loss | raw artifact를 보존하고 normalized finding을 scanner-specific evidence에 연결합니다. | raw output에서 normalized finding까지 round-trip test |
| R-013 raw output secret leak | central redaction을 raw artifact와 future SARIF에 적용합니다. | raw log, JSON, SARIF, failure output redaction test |
| R-018 license/attribution | bundled rule pack 또는 template의 license metadata를 추적합니다. | third-party asset bundle 전 license manifest entry |

## Roadmap 리스크

현재 MVP blocker는 아니지만 milestone planning에 반영해야 합니다.

| 리스크 | 처리 |
| --- | --- |
| R-003 npm 외 release artifact | npm MVP에서는 수용합니다. stable v1.0 전 checksum과 GitHub release asset을 추가합니다. |
| R-014 현재 coverage 과대평가 | README, docs, report text, release note에 MVP/non-goal 문구를 유지합니다. |
| R-015 accepted-risk workflow | severity override 또는 team policy 전에 만료일 있는 ignore rule을 구현합니다. |
| R-016 SARIF/CI artifact integration | PR-native security gate로 positioning하기 전 필수로 구현합니다. |
| R-017 OSS governance | external issue가 생기면 maintainer response expectation과 label을 추가합니다. |

## 재분류 규칙

다음 조건이 모두 충족될 때만 active register에서 리스크를 제거할 수 있습니다.

- mitigation이 구현됐거나 명시적으로 수용됐습니다.
- git, CI, npm, Codexus artifact 중 하나에 verification evidence가 있습니다.
- 다음 조치가 완료됐거나 milestone gate로 이동했습니다.
- 한국어와 영어 문서가 일치합니다.

하나라도 빠지면 active risk로 유지하고 남은 gate를 정확히 표시합니다.
