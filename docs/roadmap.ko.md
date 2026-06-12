[ENGLISH](roadmap.md) | [한국어](roadmap.ko.md)

# 로드맵

## v0.1: 로컬 CLI MVP

목표:

- 개발자가 Docker만으로 기본 보안 점검을 실행할 수 있게 한다.

기능:

- `scanrail doctor`
- `scanrail init`
- `scanrail setup`
- `scanrail run`
- `scanrail update`
- `scanrail.yaml` 생성
- `tools.lock.yaml` 생성
- HTML/JSON/SARIF 리포트

도구:

- Gitleaks
- Trivy
- Semgrep
- basic security headers checker

완료 기준:

- 신규 프로젝트에서 10분 안에 첫 리포트 생성
- high 이상 finding 기준으로 exit code 제어
- secret 원문이 리포트에 노출되지 않음

## v0.2: API 및 웹 스캔 강화

목표:

- staging URL과 OpenAPI 기반 점검 품질을 높인다.

기능:

- ZAP OpenAPI scan
- OWASP ZAP baseline
- Nuclei safe templates
- Schemathesis
- testssl.sh
- bearer/cookie 인증
- allowlist 검증
- blocked paths
- rate limit

완료 기준:

- 인증이 필요한 staging API 점검 가능
- OpenAPI 기반 API 스캔 결과가 HTML 리포트에 통합됨
- full profile 없이 active scan이 실행되지 않음

## v0.3: CI/CD 연동

목표:

- PR과 main branch에서 자동 보안 점검을 실행한다.
- npm 기반 설치와 OSS 배포를 지원한다.

기능:

- `scanrail ci init`
- GitHub Actions 템플릿
- GitLab CI 템플릿
- Jenkins 예시
- SARIF 업로드
- JUnit XML 출력
- cache 최적화
- npm wrapper package
- OS/CPU별 binary package
- release provenance

완료 기준:

- PR마다 quick profile 실행
- main merge 후 standard profile 실행
- policy 위반 시 CI 실패
- `npm install -g`와 `npx` 실행 가능

## v0.4: 결과 품질 개선

목표:

- 오탐과 중복 finding을 줄이고 개발자가 우선순위를 이해하게 한다.

기능:

- finding dedupe
- baseline suppression
- false positive allowlist
- CWE/OWASP/ASVS 매핑
- CVSS/EPSS/KEV 기반 위험도 보정
- remediation 템플릿

완료 기준:

- 같은 이슈가 여러 도구에서 발견돼도 하나로 병합
- 반복적으로 허용된 finding은 명확하게 추적
- 리포트에서 수정 우선순위가 설명됨

## v0.5: 인증과 세션 관리

목표:

- 현실적인 조직 서비스 인증 흐름을 지원한다.

기능:

- form login
- OAuth/OIDC 수동 토큰 입력
- Playwright 기반 recorded session
- role별 scan
- token refresh hook

완료 기준:

- 사용자/관리자 role별 API 접근 결과 비교 가능
- 로그인 세션 만료 시 적절한 오류를 제공

## v1.0: OSS 안정화와 조직 표준화

목표:

- 개인, 오픈소스 프로젝트, 조직 내부 개발팀이 공통 보안 셀프체크 도구로 사용할 수 있게 한다.

기능:

- 중앙 정책 배포
- 조직 scanner image registry mirror
- 프로젝트별 이력
- 조직 표준 profile
- Slack/Jira/GitHub 연동
- ASVS 기준 리포트
- 권한별 API 테스트
- IDOR/BOLA 시나리오 설정
- 공개 plugin/adapter API
- 커뮤니티 룰셋과 조직 전용 룰셋 분리

완료 기준:

- 여러 팀이 동일한 정책으로 실행 가능
- 보안팀이 예외와 정책을 중앙 관리 가능
- 개발자는 로컬/CI에서 같은 결과를 재현 가능
- 공개 OSS 사용자는 내부 인프라 없이 기본 기능을 사용할 수 있음
