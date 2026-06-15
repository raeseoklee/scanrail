[ENGLISH](release-checklist.md) | [한국어](release-checklist.ko.md)

# Release Checklist

모든 public npm publish에서 이 checklist를 사용합니다. 이 문서는 [릴리스 리스크 레지스터](release-risk-register.ko.md)의 `R-001`, `R-002` 운영 통제입니다.

Release version:

Release owner:

Date:

## 1. Git 상태

- [ ] `main`이 clean 상태이고 `origin/main`과 동기화되어 있습니다.
- [ ] version 변경이 publish 전에 commit되어 있습니다.
- [ ] 생성된 `dist/` 또는 임시 report artifact가 unstaged 상태로 남아 있지 않습니다.

증거:

```bash
git status --short --branch
git rev-parse HEAD
```

## 2. Package Set

Publish 순서:

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

- [ ] 모든 package version이 같습니다.
- [ ] 모든 package가 `publishConfig.access=public`입니다.
- [ ] 대상 version이 어떤 package에도 이미 publish되어 있지 않습니다.

증거:

```bash
npm run publish:dry-run
```

## 3. Trusted Publisher 설정

위 모든 package에 대해 npm trusted publisher 설정을 확인합니다.

```text
Publisher: GitHub Actions
Organization or user: raeseoklee
Repository: scanrail
Workflow filename: npm-publish.yml
Allowed action: npm publish
```

- [ ] `scanrail`
- [ ] `@scanrail/cli`
- [ ] `@scanrail/cli-darwin-arm64`
- [ ] `@scanrail/cli-darwin-x64`
- [ ] `@scanrail/cli-win32-x64`
- [ ] `@scanrail/cli-win32-arm64`
- [ ] `@scanrail/cli-linux-x64`
- [ ] `@scanrail/cli-linux-arm64`

증거:

release issue 또는 내부 release log에 npm UI screenshot이나 확인 메모를 남깁니다. 이 상태는 git 밖에 있습니다.

## 4. Publish 전 Registry Snapshot

publish 전 현재 registry version을 기록합니다.

```bash
for package in scanrail @scanrail/cli @scanrail/cli-darwin-arm64 @scanrail/cli-darwin-x64 @scanrail/cli-win32-x64 @scanrail/cli-win32-arm64 @scanrail/cli-linux-x64 @scanrail/cli-linux-arm64; do
  npm view "$package" version dist-tags.latest --registry https://registry.npmjs.org/ --prefer-online
done
```

- [ ] snapshot을 기록했습니다.
- [ ] target version이 모든 package에서 아직 없습니다.

## 5. Publish 전 검증

- [ ] Markdown link 검사가 통과했습니다.
- [ ] unit test와 wrapper test가 통과했습니다.
- [ ] release dry-run이 통과했습니다.
- [ ] MCP 코드가 변경됐다면 MCP Workbench가 통과했습니다.

증거:

```bash
npm run docs:check-links
npm test
make release-dry-run
mcp-workbench inspect --command node --args "examples/mcp-workbench/serve-fixture.mjs" --json
mcp-workbench run examples/mcp-workbench/scanrail-mcp.yaml --verbose
```

## 6. Publish

권장 경로:

1. `.github/workflows/npm-publish.yml`을 `mode=validate`로 실행합니다.
2. version bump 후 `.github/workflows/npm-publish.yml`을 `mode=dry-run`으로 실행합니다.
3. `.github/workflows/npm-publish.yml`을 `mode=publish`로 실행합니다.

Fallback local path:

```bash
SCANRAIL_ALLOW_NPM_PUBLISH=1 npm run publish:npm
```

- [ ] publish 경로를 기록했습니다.
- [ ] workflow run URL 또는 local command log를 기록했습니다.

## 7. Publish 후 검증

```bash
npm view scanrail version dist-tags.latest --registry https://registry.npmjs.org/ --prefer-online
npm run smoke:npm -- <version>
```

각 package의 provenance 상태를 기록합니다.

```bash
for package in scanrail @scanrail/cli @scanrail/cli-darwin-arm64 @scanrail/cli-darwin-x64 @scanrail/cli-win32-x64 @scanrail/cli-win32-arm64 @scanrail/cli-linux-x64 @scanrail/cli-linux-arm64; do
  npm view "$package@<version>" dist.attestations.provenance.predicateType --registry https://registry.npmjs.org/ --prefer-online
done
```

- [ ] `scanrail`이 public registry에서 설치됩니다.
- [ ] 설치된 package에서 `scanrail version`이 동작합니다.
- [ ] 설치된 package에서 `scanrail doctor`가 동작합니다.
- [ ] 지원되는 환경에서 npm signatures audit이 통과합니다.
- [ ] npm Smoke workflow가 Ubuntu, macOS, Windows에서 통과합니다.

## 8. 실패 처리

- [ ] 일반 release version을 overwrite하거나 unpublish하지 않습니다.
- [ ] partial publish가 발생하면 수정된 patch version을 publish합니다.
- [ ] 실패 전에 registry에 올라간 package를 기록합니다.
- [ ] 실패가 새로운 risk class를 드러내면 release issue와 risk register를 업데이트합니다.
