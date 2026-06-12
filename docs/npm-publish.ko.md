[ENGLISH](npm-publish.md) | [한국어](npm-publish.ko.md)

# npm Publish Runbook

이 문서는 Scanrail의 첫 public npm publish 절차를 정리합니다.

## Package Set

platform package를 먼저 publish하고, wrapper package를 마지막에 publish합니다.

```text
@scanrail/cli-darwin-arm64
@scanrail/cli-darwin-x64
@scanrail/cli-win32-x64
@scanrail/cli-win32-arm64
@scanrail/cli-linux-x64
@scanrail/cli-linux-arm64
@scanrail/cli
scanrail
```

`@scanrail/cli`는 platform package를 `optionalDependencies`로 참조하므로 platform package 뒤에 publish해야 설치 중 깨진 의존성 창을 줄일 수 있습니다. unscoped `scanrail` package는 `@scanrail/cli`에 의존하며, 사용자가 설치하는 기본 entrypoint로 마지막에 publish합니다.

## 현재 Registry 상태

2026년 6월 12일 기준 `@scanrail/cli`, 6개 platform package, `scanrail`은 `0.1.0` public publish 대상입니다. 이후 publish 전 registry check에서 같은 package/version이 이미 존재하면 release를 중단해야 합니다.

## 사전 조건

- `@scanrail` npm organization을 소유하거나 publish 권한이 있어야 합니다.
- npm 계정에 2FA가 켜져 있거나, publish 가능한 granular access token을 사용해야 합니다.
- 모든 package의 `publishConfig.access`는 `public`이어야 합니다.
- clean `main` branch에서 dry-run을 실행합니다.
- 같은 package version이 이미 registry에 있으면 publish하지 않습니다.

## Local Dry-Run

```bash
npm run publish:dry-run
```

이 명령은 다음을 수행합니다.

1. `go test ./...`
2. macOS, Windows, Linux binary build
3. package version 정합성 확인
4. 각 package가 아직 publish되지 않았는지 확인
5. release 순서대로 `npm publish --dry-run --access public` 실행

## Local First Publish

package가 존재하기 전 trusted publishing 설정이 어렵다면 첫 publish는 이 경로를 사용합니다.

```bash
npm login --scope=@scanrail --registry=https://registry.npmjs.org
npm whoami --registry=https://registry.npmjs.org
npm run publish:dry-run
SCANRAIL_ALLOW_NPM_PUBLISH=1 npm run publish:npm
```

guard variable은 의도된 안전장치입니다. `SCANRAIL_ALLOW_NPM_PUBLISH=1`이 없으면 실제 publish는 실행되지 않습니다.

이미 publish된 package를 건너뛰고 특정 package만 publish해야 하는 경우:

```bash
node scripts/publish-npm.mjs --dry-run --only scanrail
SCANRAIL_ALLOW_NPM_PUBLISH=1 node scripts/publish-npm.mjs --publish --only scanrail
```

## GitHub Actions Publish

workflow 파일은 `.github/workflows/npm-publish.yml`입니다.

`mode=publish`를 사용하기 전에 각 npm package에 trusted publisher를 설정합니다.

```text
Publisher: GitHub Actions
Organization or user: raeseoklee
Repository: scanrail
Workflow filename: npm-publish.yml
Allowed action: npm publish
```

실행 순서:

1. `mode=dry-run` 실행
2. 모든 package dry-run 성공 확인
3. `mode=publish` 실행

workflow는 GitHub OIDC(`id-token: write`)를 사용하고, 실제 publish 경로에서는 `--provenance`를 전달합니다.

## Publish 후 Smoke Test

```bash
npm view scanrail version
npm install -g scanrail
scanrail version
scanrail doctor
```

빈 project에서:

```bash
scanrail init --non-interactive --project-name demo --target https://example.com
scanrail run --only headers
```

## Rollback Notes

npm package version은 일반적인 release 운용에서 immutable로 다룹니다. 잘못 publish했다면 `0.1.0`을 덮어쓰지 말고 수정된 patch version을 publish합니다.

## 알려진 First-Publish Blocker

- npm 계정의 2FA가 `auth-and-writes`이면 일반 npm login으로는 OTP가 필요할 수 있습니다.
- release automation에서는 `bypass_2fa` automation token 또는 trusted publishing을 사용해야 interactive OTP를 피할 수 있습니다.
- trusted publishing은 첫 package page가 생성된 뒤 package별 설정이 필요할 수 있습니다.
