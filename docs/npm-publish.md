[ENGLISH](npm-publish.md) | [한국어](npm-publish.ko.md)

# npm Publish Runbook

This runbook covers public npm publishes for Scanrail.

## Package Set

Publish platform packages first, then the wrapper package:

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

`@scanrail/cli` depends on the platform packages through `optionalDependencies`, so publishing it after the platform packages avoids a broken install window. The unscoped `scanrail` package depends on `@scanrail/cli` and is published last as the recommended user-facing entrypoint.

## Current Registry State

As of June 12, 2026, `@scanrail/cli`, the six platform packages, and `scanrail` were first published at `0.1.0`. Version `0.1.1` later verified GitHub Actions trusted publishing with npm provenance. Registry checks before a later publish should fail the release if any target package/version already exists.

## Prerequisites

- Own the `@scanrail` npm organization or be a member with publish rights.
- Use an npm account with 2FA enabled, or a granular access token that is allowed to publish.
- Keep every package `publishConfig.access` set to `public`.
- Run the dry-run from a clean `main` branch.
- Do not publish if any package version already exists in the registry.

## Local Dry-Run

```bash
npm run publish:dry-run
```

This command:

1. runs `go test ./...`
2. builds all macOS, Windows, and Linux binaries
3. confirms package versions are aligned
4. confirms each package is not already published
5. runs `npm publish --dry-run --access public` in release order

## Local First Publish

Use this path for the first publish if trusted publishing cannot be configured before the packages exist.

```bash
npm login --scope=@scanrail --registry=https://registry.npmjs.org
npm whoami --registry=https://registry.npmjs.org
npm run publish:dry-run
SCANRAIL_ALLOW_NPM_PUBLISH=1 npm run publish:npm
```

The guard variable is intentional. Without `SCANRAIL_ALLOW_NPM_PUBLISH=1`, the publish script refuses to run a real publish.

To publish only one package, for example an alias package added after the scoped packages already exist:

```bash
node scripts/publish-npm.mjs --dry-run --only scanrail
SCANRAIL_ALLOW_NPM_PUBLISH=1 node scripts/publish-npm.mjs --publish --only scanrail
```

## GitHub Actions Publish

The workflow is `.github/workflows/npm-publish.yml`.

Before using `mode=publish`, configure a trusted publisher for each npm package:

```text
Publisher: GitHub Actions
Organization or user: raeseoklee
Repository: scanrail
Workflow filename: npm-publish.yml
Allowed action: npm publish
```

Then run the workflow manually:

1. Run `mode=validate` when checking workflow infrastructure without touching npm.
2. Run `mode=dry-run` after bumping to a version that is not in the registry.
3. Confirm every package dry-run succeeds.
4. Run `mode=publish`.

The workflow uses GitHub OIDC (`id-token: write`) and passes `--provenance` for the real publish path. It intentionally does not configure `NODE_AUTH_TOKEN` or an action-level npm registry auth file; trusted publishing supplies the publish identity through OIDC.

## Post-Publish Smoke

```bash
npm view scanrail version
npm install -g scanrail
scanrail version
scanrail doctor
npm audit signatures
```

On a clean project:

```bash
scanrail init --non-interactive --project-name demo --target https://example.com
scanrail run --only headers
```

For OS matrix validation against the public npm registry, run `.github/workflows/npm-smoke.yml` with the version or dist-tag to install.

## Rollback Notes

npm package versions are immutable in normal release practice. If a bad package is published, publish a fixed patch version instead of overwriting `0.1.0`.

## Known First-Publish Blockers

- A normal npm login may require OTP when account-level 2FA is set to `auth-and-writes`.
- Automation tokens with `bypass_2fa` or trusted publishing avoid interactive OTP in release automation.
- Trusted publishing settings live in npm outside this repository. If the publish workflow fails with authentication or provenance errors, verify the npm trusted publisher settings for every package.
