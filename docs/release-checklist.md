[ENGLISH](release-checklist.md) | [한국어](release-checklist.ko.md)

# Release Checklist

Use this checklist for every public npm publish. It is the operating control for `R-001` and `R-002` in the [Release Risk Register](release-risk-register.md).

Release version:

Release owner:

Date:

## 1. Git State

- [ ] `main` is clean and up to date with `origin/main`.
- [ ] Version changes are committed before publish.
- [ ] No generated `dist/` or temporary report artifacts are unstaged.

Evidence:

```bash
git status --short --branch
git rev-parse HEAD
```

## 2. Package Set

Publish order:

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

- [ ] Every package has the same version.
- [ ] Every package has `publishConfig.access=public`.
- [ ] The target version does not already exist for any package.

Evidence:

```bash
npm run publish:dry-run
```

## 3. Trusted Publisher Settings

For every package above, confirm the npm trusted publisher configuration:

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

Evidence:

Record npm UI screenshots or notes in the release issue or internal release log. This state lives outside git.

## 4. Pre-Publish Registry Snapshot

Record the currently published versions before publishing:

```bash
for package in scanrail @scanrail/cli @scanrail/cli-darwin-arm64 @scanrail/cli-darwin-x64 @scanrail/cli-win32-x64 @scanrail/cli-win32-arm64 @scanrail/cli-linux-x64 @scanrail/cli-linux-arm64; do
  npm view "$package" version dist-tags.latest --registry https://registry.npmjs.org/ --prefer-online
done
```

- [ ] Snapshot captured.
- [ ] Target version is absent for all packages.

## 5. Verification Before Publish

- [ ] Markdown links pass.
- [ ] Unit and wrapper tests pass.
- [ ] Release dry-run passes.
- [ ] MCP Workbench passes if MCP code changed.

Evidence:

```bash
npm run docs:check-links
npm test
make release-dry-run
mcp-workbench inspect --command node --args "examples/mcp-workbench/serve-fixture.mjs" --json
mcp-workbench run examples/mcp-workbench/scanrail-mcp.yaml --verbose
```

## 6. Publish

Preferred path:

1. Run `.github/workflows/npm-publish.yml` with `mode=validate`.
2. Run `.github/workflows/npm-publish.yml` with `mode=dry-run` after version bump.
3. Run `.github/workflows/npm-publish.yml` with `mode=publish`.

Fallback local path:

```bash
SCANRAIL_ALLOW_NPM_PUBLISH=1 npm run publish:npm
```

- [ ] Publish path recorded.
- [ ] Workflow run URL or local command log recorded.

## 7. Post-Publish Verification

```bash
npm view scanrail version dist-tags.latest --registry https://registry.npmjs.org/ --prefer-online
npm run smoke:npm -- <version>
```

For each package, record provenance state:

```bash
for package in scanrail @scanrail/cli @scanrail/cli-darwin-arm64 @scanrail/cli-darwin-x64 @scanrail/cli-win32-x64 @scanrail/cli-win32-arm64 @scanrail/cli-linux-x64 @scanrail/cli-linux-arm64; do
  npm view "$package@<version>" dist.attestations.provenance.predicateType --registry https://registry.npmjs.org/ --prefer-online
done
```

- [ ] `scanrail` installs from the public registry.
- [ ] `scanrail version` works from the installed package.
- [ ] `scanrail doctor` works from the installed package.
- [ ] npm signatures audit passes where supported.
- [ ] npm Smoke workflow passes for Ubuntu, macOS, and Windows.

## 8. Failure Handling

- [ ] Do not overwrite or unpublish a normal release version.
- [ ] If a partial publish happens, publish a fixed patch version.
- [ ] Record which packages reached the registry before failure.
- [ ] Update the release issue and risk register if the failure exposes a new class of risk.
