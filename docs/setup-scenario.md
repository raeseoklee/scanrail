[ENGLISH](setup-scenario.md) | [한국어](setup-scenario.ko.md)

# Setup Scenario

This document shows the expected first-run flow for a developer adopting Scanrail.

## 1. Install

```bash
npm install -g @scanrail/cli
```

The npm wrapper is the primary install path for the first public release. Release archives and installer scripts can be added after release automation is mature.

## 2. Check Environment

```bash
scanrail doctor
```

Example output:

```text
Scanrail Doctor

Docker              OK   Docker Desktop running
Docker Compose      OK
Workspace           OK   /Users/dev/workspace/my-order-api
Git repo            OK
Disk space          OK

Ready.
```

The first release candidate can also run without Docker-backed scanners when using the native headers scanner.

## 3. Initialize Configuration

```bash
scanrail init
```

Example interaction:

```text
? Project name? my-order-api

? Select target types
  ✓ Local code repository
  ✓ Staging web URL
  ✓ OpenAPI/Swagger API
  ○ Container image

? Staging URL?
  https://staging-order.example.com

? Is this production?
  No

? Confirm allowed domains
  ✓ staging-order.example.com

? OpenAPI file or URL?
  ./openapi.yaml

? Authentication required?
  Bearer token

? Where should the token be read from?
  Environment variable

? Environment variable name?
  SCANRAIL_TOKEN

? Default scan profile?
  quick

? CI failure threshold?
  high and above

? Any paths to avoid?
  /logout
```

For CI or scripted setup:

```bash
scanrail init \
  --non-interactive \
  --project-name my-order-api \
  --target https://staging-order.example.com \
  --profile quick
```

## 4. Review `scanrail.yaml`

Example:

```yaml
project:
  name: my-order-api

targets:
  web:
    url: https://staging-order.example.com
    allowlist:
      - staging-order.example.com

profiles:
  default: quick
  quick:
    tools:
      - headers

report:
  output_dir: .scanrail/reports
  formats:
    - html
    - json
```

## 5. Prepare Workspace

```bash
scanrail setup
```

Offline or first-release dry-run setup:

```bash
scanrail setup --pull-policy never
```

## 6. Run a Scan

```bash
scanrail run --profile quick
```

For the current MVP scanner:

```bash
scanrail run --only headers
```

## 7. Review Reports

Default report directory:

```text
.scanrail/reports
```

Expected files:

```text
scanrail-report.json
scanrail-report.html
```

## Docker and Localhost

When Docker-backed scanners are enabled, `localhost` inside a scanner container points to the container itself, not the developer host. Use host mappings such as `host.docker.internal` where supported, or configure a reachable staging URL.
