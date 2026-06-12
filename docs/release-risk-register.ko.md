[ENGLISH](release-risk-register.md) | [한국어](release-risk-register.ko.md)

# 릴리스 리스크 레지스터

이 문서는 첫 npm publish와 `0.1.1` trusted publishing 검증 이후 남은 release risk를 추적합니다.

## 현재 상태

| 리스크 | 상태 | 완화 |
| --- | --- | --- |
| GitHub Actions Node 20 runtime deprecation warning | 완화 | workflow가 `node24` action runtime을 선언하는 `actions/checkout@v5`, `actions/setup-node@v5`, `actions/setup-go@v6`를 사용합니다. |
| publish workflow의 npm `always-auth` warning | 완화 | publish workflow가 더 이상 `setup-node`에 npm registry auth 설정 생성을 요청하지 않습니다. Trusted publishing은 `.npmrc` token auth 대신 OIDC를 사용합니다. |
| published package가 모든 target OS에서 실행되는지 불확실 | workflow로 완화 | `.github/workflows/npm-smoke.yml`이 Ubuntu, macOS, Windows에서 public `scanrail` package를 설치하고 `version`, `doctor`, `npm audit signatures`를 실행합니다. |
| npm trusted publishing 설정 drift | 외부 잔여 리스크 | 설정은 git이 아니라 npm에 있습니다. release마다 workflow filename, repository, owner, allowed action이 npm package trusted publisher 설정과 맞아야 합니다. |
| partial publish로 package set 불일치 발생 | 운영 잔여 리스크 | platform package를 wrapper보다 먼저 publish하고, publish 전에 version을 확인합니다. 실패한 version은 이후 patch version으로 수정합니다. |
| npm 외 release artifact가 아직 unsigned 또는 부재 | npm MVP에서 수용 | npm provenance는 적용됐습니다. GitHub release archive, checksum, 추가 package manager channel은 roadmap 작업입니다. |
| scanner adapter 확장으로 network 또는 credential exposure 증가 | 제품 safety 리스크 | active scan은 opt-in으로 유지하고, target allowlist를 강제하며, secret value 대신 environment variable name만 참조합니다. |

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

OS matrix 검증은 `npm Smoke` workflow를 실행합니다.

## 외부 확인

- workflow filename 또는 repository owner를 변경하기 전 모든 package의 npm trusted publisher 설정을 확인합니다.
- published version의 npm package page에 provenance가 표시되는지 확인합니다.
- workflow 변경 후 GitHub Actions에 deprecation warning 또는 auth warning이 없는지 확인합니다.
