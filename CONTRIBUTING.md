[ENGLISH](CONTRIBUTING.md) | [한국어](CONTRIBUTING.ko.md)

# Contributing

Thanks for helping improve Scanrail. This project is early, so the highest-value contributions are small, reproducible, and tied to a real developer workflow.

## Good First Contributions

- Reproduce a scanner adapter issue with a minimal fixture.
- Improve documentation for installation, setup, or reports.
- Add tests around config parsing, report generation, or MCP JSON-RPC behavior.
- Propose a scanner adapter contract before implementing a new adapter.

## Development Setup

```bash
git clone https://github.com/raeseoklee/scanrail.git
cd scanrail
go test ./...
npm test
make release-dry-run
```

## Pull Request Checklist

- Keep the change focused and reversible.
- Add or update tests when behavior changes.
- Update English docs first; add Korean `.ko.md` docs when user-facing docs change.
- Run `npm run docs:check-links`, `npm test`, and `make release-dry-run` before asking for review.
- Do not add a new dependency unless the PR explains why the standard library or existing tooling is insufficient.

## Scanner Adapter Changes

Before adding or changing a scanner adapter, document:

- scanner version and container image
- exact command line
- required inputs and secrets
- target safety controls
- raw output locations
- normalized finding mapping
- known gaps and skipped behavior

## Commit Messages

Use intent-first commit messages. Explain why the change exists, what constraints shaped it, and what was tested.

## Security

Do not include exploit details, real secrets, customer target URLs, or private scan results in public issues or pull requests. See [SECURITY.md](SECURITY.md).
