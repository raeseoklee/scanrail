[ENGLISH](open-source-tools.md) | [한국어](open-source-tools.ko.md)

# 오픈소스 도구 검토

`scanrail`은 여러 오픈소스 도구를 Docker 컨테이너로 실행하고 결과를 통합합니다. 각 도구는 특정 영역에 강점이 있으며, 초기 MVP에서는 설치/실행/결과 파싱 부담을 `scanrail`이 흡수합니다.

## 1차 후보

| 영역 | 도구 | 용도 | 도입 단계 |
| --- | --- | --- | --- |
| Headers | native Go checker | 보안 헤더 점검 | v0.1 |
| SAST | Semgrep | 코드 취약점 및 보안 패턴 점검 | v0.1 |
| 의존성/컨테이너/IaC | Trivy | CVE, 컨테이너, IaC, SBOM 점검 | v0.1 |
| Secrets | Gitleaks | 하드코딩된 토큰/키 탐지 | v0.1 |
| DAST | OWASP ZAP | 웹 런타임 취약점 점검 | v0.2 |
| 템플릿 기반 스캔 | Nuclei | CVE/설정/노출 패턴 점검 | v0.2 |
| TLS | testssl.sh | TLS/인증서 설정 점검 | v0.2 |
| API 테스트 | Schemathesis | OpenAPI 기반 API property/fuzz 테스트 | v0.2 |
| 결과 연동 | SARIF | code scanning 결과 연동 | v0.1 partial |

## 추가 후보

| 영역 | 도구 | 용도 | 비고 |
| --- | --- | --- | --- |
| SAST | CodeQL | 정교한 코드 쿼리 기반 분석 | 언어/CI 환경에 따라 도입 |
| 웹 서버 스캔 | Nikto | 서버 설정/취약 파일 점검 | 오탐 관리 필요 |
| 서비스 탐색 | Nmap | 포트/서비스 탐색 | 허용 범위 관리 중요 |
| 공격면 탐색 | OWASP Amass | 서브도메인/외부 노출 탐색 | 조직 정책 필요 |
| 취약점 관리 | DefectDojo | finding 관리/이력/워크플로 | 중앙 서버 단계에서 검토 |
| SBOM 관리 | Dependency-Track | SBOM 기반 공급망 위험 추적 | 조직 단위 운영에 적합 |

## 도구별 역할

### OWASP ZAP

역할:

- baseline scan
- passive scan
- OpenAPI 기반 API scan
- full active scan

초기 정책:

- `quick`, `standard`에서는 baseline/API scan만 사용
- active scan은 `full` profile에서만 허용
- 명시 플래그와 allowlist가 없으면 실행 거부

### Nuclei

역할:

- 안전한 템플릿 기반 점검
- 알려진 CVE, misconfiguration, exposed panel 탐지
- 조직 커스텀 템플릿 확장

초기 정책:

- safe templates 위주 실행
- destructive 또는 intrusive templates는 기본 제외
- template version을 lock하거나 organization mirror 사용

### Semgrep

역할:

- 코드 레벨 취약점 탐지
- 언어별 보안 룰 적용
- 조직 보안 코딩 룰 추가

초기 정책:

- 기본 community rules 사용
- 조직 룰셋은 별도 namespace로 관리
- 결과는 SARIF와 JSON으로 수집

### Trivy

역할:

- 의존성 CVE 점검
- 컨테이너 이미지 점검
- IaC misconfiguration 점검
- SBOM 생성

초기 정책:

- repo filesystem scan
- container image scan은 설정된 경우만 실행
- DB 캐시는 `scanrail setup`과 `scanrail update`에서 관리

### Gitleaks

역할:

- Git history와 working tree의 secret 탐지

초기 정책:

- 기본적으로 working tree 중심
- history scan은 옵션으로 제공
- false positive allowlist를 프로젝트 설정에 저장

### testssl.sh

역할:

- TLS protocol/cipher/certificate 점검

초기 정책:

- HTTPS target이 있을 때만 실행
- 결과는 JSON 또는 파싱 가능한 텍스트로 수집

### Schemathesis

역할:

- OpenAPI spec 기반 API 테스트
- 계약 위반, 서버 오류, edge case 탐지

초기 정책:

- staging 대상에서만 실행
- destructive method는 기본 제한 또는 명시 설정 필요

## Docker 이미지 관리

`scanrail setup`은 도구 이미지를 준비하고 `tools.lock.yaml`을 생성합니다.

예시:

```yaml
tools:
  zap:
    image: ghcr.io/zaproxy/zaproxy:stable
  nuclei:
    image: projectdiscovery/nuclei:<approved-version>
  semgrep:
    image: semgrep/semgrep:<approved-version>
  trivy:
    image: aquasec/trivy:<approved-version>
  gitleaks:
    image: zricethezav/gitleaks:<approved-version>
  testssl:
    image: drwetter/testssl.sh:<approved-version>
  schemathesis:
    image: schemathesis/schemathesis:<approved-version>
```

운영 권장사항:

- 개발 초기에는 tag 기반으로 시작하되 `latest`는 문서 예시와 CI 기본값에서 피한다.
- 조직 배포 전에는 digest 또는 승인된 버전으로 고정한다.
- 조직 규모로 운영할 때는 private registry mirror를 사용할 수 있다.
- 이미지와 룰셋 업데이트는 `scanrail update`로 통제한다.

## 라이선스 검토 필요 항목

조직 배포 전 다음을 확인해야 합니다.

- 각 도구의 라이선스
- Docker image 재배포 가능 여부
- private registry mirror 가능 여부
- template/rule 데이터의 라이선스
- scanner output을 내부 또는 외부 시스템에 저장하는 것에 대한 제약
