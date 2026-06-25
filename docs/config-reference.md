[ENGLISH](config-reference.md) | [한국어](config-reference.ko.md)

# Configuration Reference

`scanrail.yaml` stores project-level security scan configuration. It should be safe to commit and must not contain raw secrets.

## Full Example

```yaml
project:
  name: my-order-api
  type: web-api
  owner: payments-platform
  criticality: high

targets:
  web:
    url: https://staging-order.example.com
    environment: staging
    allowlist:
      - staging-order.example.com
    exclude_paths:
      - /logout
      - /admin/destructive/*
      - /payments/real-charge

  api:
    openapi: ./openapi.yaml

  container:
    image: registry.example.com/order-api:staging

auth:
  type: bearer
  token_env: SCANRAIL_TOKEN

profiles:
  default: quick

  quick:
    tools:
      - gitleaks
      - headers
      - tls
      - openapi

safety:
  active_scan_default: false
  require_allowlist: true
  max_rps: 5
  add_header:
    X-Scanrail-Project: my-order-api
  blocked_paths:
    - /logout
    - /admin/destructive/*
    - /payments/real-charge

policy:
  fail_on:
    severity: high
  ignore:
    - id: finding-id
      reason: accepted risk
      expires: 2026-12-31

report:
  output_dir: .scanrail/reports
  formats:
    - html
    - json
    - sarif
```

## `project`

Project metadata.

```yaml
project:
  name: my-order-api
  type: web-api
  owner: payments-platform
  criticality: high
```

Required:

- `name`

Recommended:

- `type`
- `owner`
- `criticality`

## `targets`

Targets define what scanners can inspect.

### Web Target

```yaml
targets:
  web:
    url: https://staging.example.com
    environment: staging
    allowlist:
      - staging.example.com
    exclude_paths:
      - /logout
```

`exclude_paths` describes paths that scanners should avoid when they can. It is advisory unless the scanner adapter declares enforcement support.

### API Target

```yaml
targets:
  api:
    openapi: ./openapi.yaml
```

The native OpenAPI scanner currently reads local JSON or common YAML OpenAPI files only. It does not fetch remote specs and does not probe API endpoints. OpenAPI servers must still pass allowlist validation before future interactive API scanning.

### Container Target

```yaml
targets:
  container:
    image: registry.example.com/order-api:staging
```

Container scanning is planned for Trivy-backed adapter phases.

## `auth`

Authentication references must avoid raw secrets.

```yaml
auth:
  type: bearer
  token_env: SCANRAIL_TOKEN
```

Supported shape:

- `type: none`
- `type: bearer` with `token_env`

Planned shapes:

- cookie
- custom header
- browser/session export

## `profiles`

Profiles select scanner groups.

```yaml
profiles:
  default: quick
  quick:
    tools:
      - gitleaks
      - headers
      - tls
      - openapi
```

The current generated `quick` profile only includes implemented adapters: Docker-backed `gitleaks`, native `headers`, native `tls`, and native local-file-only `openapi`. Extended profiles belong in examples until their adapters are executable.

## `safety`

Safety policy controls target boundaries and scanner behavior.

```yaml
safety:
  active_scan_default: false
  require_allowlist: true
  max_rps: 5
  blocked_paths:
    - /logout
```

`blocked_paths` is a required safety boundary for scanners that can enforce it. If a scanner cannot enforce a required capability, Scanrail skips it or fails explicit execution.

## `policy`

Policy decides whether findings fail the run.

```yaml
policy:
  fail_on:
    severity: high
```

Ignore rules should include reasons and expiration dates.

## `report`

Report output configuration.

```yaml
report:
  output_dir: .scanrail/reports
  formats:
    - html
    - json
```

Default output directory:

```text
.scanrail/reports
```

## Secret Rules

- Do not store tokens, passwords, cookies, or API keys in `scanrail.yaml`.
- Reference secret environment variable names instead.
- Redact secret values from logs, reports, and raw command previews.
