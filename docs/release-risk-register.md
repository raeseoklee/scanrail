[ENGLISH](release-risk-register.md) | [한국어](release-risk-register.ko.md)

# Release Risk Register

This register tracks release risks that remain after the first npm publication and the `0.1.1` trusted publishing verification.

## Current Status

| Risk | Status | Mitigation |
| --- | --- | --- |
| GitHub Actions Node 20 runtime deprecation warnings | Mitigated | Workflows use `actions/checkout@v5`, `actions/setup-node@v5`, and `actions/setup-go@v6`, which declare `node24` action runtimes. |
| npm `always-auth` warning during publish workflow | Mitigated | The publish workflow no longer asks `setup-node` to create npm registry auth configuration. Trusted publishing uses OIDC instead of `.npmrc` token auth. |
| Published package does not run on every target OS | Mitigated by workflow | `.github/workflows/npm-smoke.yml` installs the public `scanrail` package on Ubuntu, macOS, and Windows, then runs `version`, `doctor`, and `npm audit signatures`. |
| Trusted publishing settings drift in npm | Residual external risk | Settings live in npm, not git. Each release must keep the workflow filename, repository, owner, and allowed action aligned with the npm package trusted publisher configuration. |
| Partial publish leaves package set inconsistent | Residual operational risk | Platform packages publish before wrappers, package versions are checked before publish, and failed versions must be fixed with a later patch version. |
| Release artifacts outside npm are unsigned or absent | Accepted for npm MVP | npm provenance is in place. GitHub release archives, checksums, and optional package-manager channels remain roadmap work. |
| Scanner adapters can expand network or credential exposure | Product safety risk | Active scans stay opt-in, target allowlists are enforced, and secrets are referenced by environment variable names rather than values. |

## Verification Gates

Before a publish:

```bash
npm run docs:check-links
npm test
npm run publish:dry-run
make release-dry-run
```

For the current already-published version, use `mode=validate` in the publish workflow instead of `mode=dry-run`.

After a publish:

```bash
npm view scanrail version
npm run smoke:npm -- <version>
```

Run the `npm Smoke` workflow for OS matrix validation.

## External Checks

- Confirm npm trusted publisher settings for every package before changing the workflow filename or repository owner.
- Confirm package pages show provenance for the published version.
- Confirm GitHub Actions has no deprecation or auth warnings after workflow changes.
