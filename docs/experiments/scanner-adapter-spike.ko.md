[ENGLISH](scanner-adapter-spike.md) | [한국어](scanner-adapter-spike.ko.md)

# Scanner Adapter 실증

이 실험은 production adapter 코드를 구현하기 전에 다음 Scanrail 구현 방향을 검증하기 위한 spike입니다. Docker로 세 가지 오픈소스 scanner를 실행하고, raw JSON output을 수집한 뒤, 하나의 finding report 형태로 정규화합니다. 또한 VHS tape로 짧은 터미널 데모 영상을 남깁니다.

## 목적

제품 방향은 multi-scanner orchestration이지만 현재 CLI에서 실제 실행되는 scanner는 native headers scanner 하나입니다. Gitleaks, Trivy, Semgrep adapter를 core에 붙이기 전에 실제 command 동작, output shape, exit code, path mount, redaction 필요 지점을 확인합니다.

## 도구

| Tool | Image | 실증 목적 |
| --- | --- | --- |
| Gitleaks | `ghcr.io/gitleaks/gitleaks:v8.30.1` | custom rule 기반 fake secret 탐지 |
| Semgrep | `semgrep/semgrep:1.165.0` | Python insecure-code finding |
| Trivy | `aquasec/trivy:0.71.0` | Terraform misconfiguration finding |

## 실행

```bash
npm run experiment:scanner-spike
```

산출물:

```text
experiments/scanner-adapter-spike/out/raw/
experiments/scanner-adapter-spike/out/scanrail-normalized-report.json
experiments/scanner-adapter-spike/out/scanrail-normalized-report.html
experiments/scanner-adapter-spike/out/summary.md
```

## Summary Demo

```bash
npm run experiment:scanner-spike:summary
```

## 영상

tape file:

```text
experiments/scanner-adapter-spike/tapes/scanner-adapter-spike.tape
```

GIF 렌더링:

```bash
make tape-scanner-spike
```

영상은 짧고 결정적인 데모를 위해 전체 Docker scan이 아니라 summary command를 보여줍니다. 전체 증거를 다시 생성하려면 `npm run experiment:scanner-spike`를 실행합니다.

생성 asset:

![Scanner adapter spike](../assets/scanner-adapter-spike.gif)

## 실증 결과

- 테스트한 Docker 환경에서 세 scanner image 실행이 가능했습니다.
- Gitleaks는 custom config와 scanner-specific expected exit code가 필요합니다.
- Semgrep JSON은 초기 finding model에 필요한 rule, path, severity, message를 제공합니다.
- Trivy misconfiguration output은 ID, title, severity, location, remediation을 제공합니다.
- raw output은 보존하되, report에 노출되는 evidence는 redaction을 거쳐야 합니다.
- adapter 구현은 scanner exit code를 전역 규칙이 아니라 scanner별 계약으로 다뤄야 합니다.

## 구현 반영점

다음 production 구현에는 아래가 필요합니다.

- Docker runner abstraction
- scanner별 expected exit-code handling
- `.scanrail/raw` raw artifact capture
- adapter별 finding normalization
- report 전 central redaction
- safety gate용 adapter capability metadata
