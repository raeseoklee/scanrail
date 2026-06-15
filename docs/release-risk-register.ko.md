[ENGLISH](release-risk-register.md) | [한국어](release-risk-register.ko.md)

# 릴리스 및 제품 리스크 레지스터

이 문서는 첫 public npm release, trusted publishing 설정, `0.1.2` MCP MVP, MCP Workbench 검증 이후 남은 주요 리스크를 추적합니다.

최종 검토일: 2026년 6월 15일.

## 리스크 수준

| 수준 | 의미 |
| --- | --- |
| High | release adoption을 막거나, unsafe scan behavior를 만들거나, public package trust를 훼손할 수 있음 |
| Medium | user confusion, reliability 저하, 다음 milestone 지연으로 이어질 수 있음 |
| Low | blast radius가 작거나 manual workaround가 명확한 known gap |

## 현재 요약

| 영역 | 상태 | 주요 잔여 리스크 |
| --- | --- | --- |
| npm distribution | MVP 기준 안정 | trusted publisher 설정과 npm package 상태가 git 밖에 있음 |
| cross-platform install | smoke workflow로 커버 | public npm smoke는 scheduled/manual이며 모든 commit에서 실행되지는 않음 |
| MCP MVP | 구현 및 Workbench 검증 완료 | 실제 host compatibility와 MCP security hardening은 미완료 |
| scanner adapter | spike로만 실증 | Docker runner, adapter contract, redaction, safety gate가 production code가 아님 |
| reporting | native headers용 JSON/HTML 제공 | multi-scanner normalization, SARIF, false-positive workflow는 계획 단계 |
| OSS operations | public repo 사용 가능 | community triage, support boundary, governance는 아직 가벼움 |

## 즉시 우선순위

1. Docker-backed scanner를 default profile에 넣기 전에 adapter safety와 redaction gap을 먼저 닫습니다 (`R-009`, `R-013`).
2. broader MCP execution tool을 추가하기 전에 MCP resource, audit log, host compatibility를 harden합니다 (`R-005`, `R-007`, `R-008`).
3. 각 publish 전후 npm package state를 기록하는 release checklist를 추가합니다 (`R-001`, `R-002`).
4. normalized report를 확장하기 전에 raw scanner evidence와 version metadata를 보존합니다 (`R-011`, `R-012`).
5. Scanrail을 PR-native security gate로 positioning하기 전에 SARIF를 required milestone로 취급합니다 (`R-016`).

## 활성 리스크

| ID | 리스크 | 수준 | 상태 | 완화 | 다음 조치 |
| --- | --- | --- | --- | --- | --- |
| R-001 | npm trusted publishing 설정 drift | High | 외부 잔여 리스크 | 설정은 git이 아니라 npm에 있습니다. release 문서에 필요한 workflow filename, repository, owner, allowed action을 기록했습니다. | 각 publish 전 모든 package의 trusted publisher 설정을 `.github/workflows/npm-publish.yml`과 대조합니다. |
| R-002 | partial publish로 package set 불일치 발생 | High | 운영 잔여 리스크 | platform package를 wrapper보다 먼저 publish하고, publish 전에 version을 확인합니다. 실패한 version은 이후 patch version으로 수정합니다. | publish 전후 package version을 기록하는 release checklist 항목을 추가합니다. |
| R-003 | npm 외 release artifact가 아직 unsigned 또는 부재 | Medium | npm MVP에서 수용 | npm provenance는 적용됐습니다. GitHub release archive, checksum, 추가 package manager channel은 roadmap 작업입니다. | stable `v1.0` release path 선언 전 checksum과 GitHub release asset 생성을 추가합니다. |
| R-004 | public npm smoke가 모든 commit에서 실행되지 않음 | Medium | 부분 완화 | `.github/workflows/npm-smoke.yml`이 scheduled/manual로 public registry 대상 Ubuntu, macOS, Windows 검증을 수행합니다. | release branch에 lightweight post-publish required smoke gate를 둘지 결정합니다. |
| R-005 | MCP server compatibility가 client ecosystem 전체 대비 좁음 | Medium | 부분 완화 | `mcp-workbench inspect`와 Workbench regression spec이 stdio protocol behavior를 검증합니다. | MCP tool 확장 전 production MCP host 최소 1개에 대한 smoke note 또는 fixture를 추가합니다. |
| R-006 | MCP tool이 core policy와 diverge하면 CLI safety를 우회할 수 있음 | High | 설계로 통제 | MCP는 CLI safety model 위의 thin adapter로 유지하고, active scan은 `confirm_active_scan=true`와 allowlist를 요구합니다. | MCP tool implementation이 CLI execution과 같은 config, exit-code, safety validation path를 사용하게 유지합니다. |
| R-007 | MCP resource가 sensitive config 또는 oversized report data를 노출할 수 있음 | High | 부분 완화 | 현재 config resource는 secret value가 아니라 environment variable name만 노출합니다. full report resource는 아직 defer 상태입니다. | `scanrail://reports/latest/json` 또는 richer resource 추가 전 size limit과 redaction test를 추가합니다. |
| R-008 | MCP tool call audit이 team environment에 충분하지 않음 | Medium | known gap | 설계 문서는 scan/setup tool call 기록을 요구하지만 audit logging은 아직 완성된 product feature가 아닙니다. | `scanrail_setup` 또는 broader scanner execution 추가 전 MCP-triggered scan에 대한 local execution log를 추가합니다. |
| R-009 | scanner adapter 확장으로 network 또는 credential exposure 증가 | High | 제품 safety 리스크 | active scan은 opt-in, target allowlist는 강제, secret은 environment variable name으로만 참조, adapter capability는 문서화했습니다. | Docker-backed scanner를 default profile에 넣기 전 adapter capability metadata와 fail/skip behavior를 구현합니다. |
| R-010 | Docker-backed scanner가 host OS별로 다르게 동작할 수 있음 | Medium | spike로만 검증 | scanner adapter spike는 한 Docker 환경에서 command shape과 output을 검증했으며 full OS matrix는 아닙니다. | Linux integration test를 먼저 만들고 macOS/Windows Docker compatibility note 또는 matrix coverage를 추가합니다. |
| R-011 | scanner image 또는 rule update가 finding을 예고 없이 바꿀 수 있음 | Medium | known gap | spike image는 pinning되어 있습니다. production adapter versioning과 update policy는 아직 없습니다. | default scanner image/template을 pinning하고 모든 report에 version을 기록합니다. |
| R-012 | normalized report가 scanner-specific context를 숨길 수 있음 | Medium | known gap | spike는 raw output을 보존하고 core field를 normalize합니다. production report schema는 아직 evolving 상태입니다. | raw artifact를 별도 저장하고 normalized finding이 scanner, rule, evidence, remediation field로 되돌아갈 수 있게 합니다. |
| R-013 | secret redaction이 raw scanner output path를 놓칠 수 있음 | High | known gap | product doc은 redaction을 요구하고, spike는 raw artifact capture를 implementation implication으로 식별했습니다. | log, raw artifact, MCP resource, HTML, JSON, future SARIF를 노출하기 전 central redaction을 적용합니다. |
| R-014 | 사용자가 현재 native headers scan의 coverage를 과대평가할 수 있음 | Medium | 문서 리스크 | README와 docs는 MVP가 commercial scanner parity가 아니며 Docker-backed adapter가 planned 상태라고 설명합니다. | README, report, release note에 capability와 non-goal 문구를 계속 노출합니다. |
| R-015 | team-scale accepted risk management workflow가 아직 없음 | Medium | roadmap item | ignore rule, policy pack, issue export, workflow integration은 이후 milestone에 있습니다. | severity override 또는 team policy feature 전에 expiring ignore rule을 구현합니다. |
| R-016 | SARIF와 CI artifact integration이 아직 ship되지 않음 | Low | 계획됨 | CI/CD milestone에 SARIF, JUnit XML, CI template이 포함되어 있습니다. | Scanrail을 PR-native security gate로 positioning하기 전 SARIF를 required로 취급합니다. |
| R-017 | OSS support와 governance가 adoption을 따라가지 못할 수 있음 | Medium | lightweight process | contributing, security policy, issue template, roadmap이 있습니다. | external issue가 들어오기 시작하면 maintainer response expectation과 label workflow를 추가합니다. |
| R-018 | bundled scanner rule/template의 license와 attribution obligation을 놓칠 수 있음 | Medium | future integration risk | 현재 package는 third-party scanner rule pack을 bundle하지 않습니다. | Semgrep/Nuclei/config asset을 bundle할 때 license metadata를 추적합니다. |

## 완화된 리스크

| 리스크 | 완화 |
| --- | --- |
| GitHub Actions Node 20 runtime deprecation warning | workflow가 `node24` action runtime을 선언하는 `actions/checkout@v5`, `actions/setup-node@v5`, `actions/setup-go@v6`를 사용합니다. |
| publish workflow의 npm `always-auth` warning | publish workflow가 더 이상 `setup-node`에 npm registry auth 설정 생성을 요청하지 않습니다. Trusted publishing은 `.npmrc` token auth 대신 OIDC를 사용합니다. |
| published package가 모든 target OS에서 실행되는지 불확실 | npm smoke workflow가 Ubuntu, macOS, Windows에서 public `scanrail` package를 설치하고 `version`, `doctor`, signature audit check를 실행합니다. |
| Windows hosted runner label migration | 2026년 6월 GitHub runner image notice 이후 Windows job은 `windows-2025-vs2026`을 명시합니다. |
| MCP MVP에 real client-style regression coverage 부족 | `examples/mcp-workbench/scanrail-mcp.yaml`이 MCP discovery, resource, safety gating, confirmed native headers scan execution, latest report retrieval을 검증합니다. |

## 검증 게이트

publish 전:

```bash
npm run docs:check-links
npm test
npm run publish:dry-run
make release-dry-run
```

현재 이미 publish된 version으로는 publish workflow의 `mode=validate`를 사용하고 `mode=dry-run`은 새 version으로 bump한 뒤 사용합니다.

publish 후:

```bash
npm view scanrail version
npm run smoke:npm -- <version>
```

public npm registry 대상 OS matrix 검증은 `npm Smoke` workflow를 실행합니다.

MCP regression check:

```bash
mcp-workbench inspect --command node --args "examples/mcp-workbench/serve-fixture.mjs" --json
mcp-workbench run examples/mcp-workbench/scanrail-mcp.yaml --verbose
```

## 외부 확인

- workflow filename, repository owner, npm organization setting을 변경하기 전 모든 package의 npm trusted publisher 설정을 확인합니다.
- published version의 npm package page에 provenance가 표시되는지 확인합니다.
- workflow 변경 후 GitHub Actions에 deprecation warning 또는 auth warning이 없는지 확인합니다.
- third-party scanner asset을 bundle하기 전 scanner image license, rule license, update policy를 확인합니다.
- MCP execution tool을 확장하기 전 production AI client 최소 1개에서 MCP host behavior를 확인합니다.
