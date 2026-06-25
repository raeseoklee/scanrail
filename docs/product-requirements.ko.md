[ENGLISH](product-requirements.md) | [한국어](product-requirements.ko.md)

# 제품 요구사항

## 배경

상용 보안진단 제품을 도입하기 어렵거나 승인 과정이 오래 걸리는 팀은 자체적으로 최소한의 보안 점검을 반복 수행할 수 있는 표준 도구가 필요합니다. 이 문제는 특정 회사 내부에만 국한되지 않고, 스타트업, 소규모 개발팀, 오픈소스 프로젝트, 엔터프라이즈 내부 플랫폼 팀에서 반복적으로 발생합니다.

개별 오픈소스 보안 도구는 이미 강력하지만, 개발자가 각각의 설치법, 옵션, 결과 포맷, 오탐 처리 방식을 모두 이해하기는 어렵습니다. `scanrail`은 이 도구들을 하나의 표준 인터페이스로 묶어 개발자가 쉽게 실행하고 결과를 해석할 수 있게 합니다.

## 목표

- 개발자가 PR 전 또는 배포 전 보안 셀프체크를 수행할 수 있게 한다.
- 여러 오픈소스 스캐너를 Docker 기반으로 자동 구성한다.
- 프로젝트별 설정을 인터랙티브하게 생성한다.
- 스캐너별 결과를 하나의 finding 모델로 정규화한다.
- OWASP, CWE, CVSS, EPSS, KEV 기준으로 위험도를 설명한다.
- HTML, JSON, SARIF 리포트를 생성한다.
- CI/CD에서 자동 실행할 수 있게 한다.
- 조직별 정책, registry mirror, 조직 전용 룰셋은 공개 코어 위에 설정으로 얹을 수 있게 한다.

## 비목표

- 상용 DAST/SAST 제품의 모든 기능을 즉시 대체하지 않는다.
- SQL injection, XSS, 인증 우회 탐지 엔진을 처음부터 구현하지 않는다.
- 허가되지 않은 외부 대상에 대한 공격성 스캔을 지원하지 않는다.
- 운영 환경에 대한 active scan을 기본 허용하지 않는다.
- 비즈니스 로직 취약점을 완전 자동으로 탐지한다고 주장하지 않는다.

## 주요 사용자

### 개발자

- 로컬 또는 CI에서 보안 점검을 실행한다.
- 리포트를 보고 취약점을 수정한다.
- 보안 도구별 세부 사용법을 몰라도 된다.

### 보안 담당자

- 기본 프로파일과 정책을 관리한다.
- 조직 보안 기준에 맞는 룰과 예외를 정의한다.
- 결과를 중앙에서 추적하고 개선 방향을 제시한다.

### 플랫폼/DevOps 담당자

- CI 템플릿과 내부 registry mirror를 관리한다.
- Docker 이미지 버전과 캐시 정책을 관리한다.

### 오픈소스 메인테이너

- 프로젝트에 기본 보안 점검을 추가한다.
- 공개 CI에서 실행 가능한 safe profile을 사용한다.
- SARIF 리포트를 code scanning에 연동한다.

## 핵심 기능

### 1. 인터랙티브 초기 설정

`scanrail init`은 프로젝트 특성에 맞는 `scanrail.yaml`을 생성합니다.

수집 정보:

- 프로젝트 이름
- 점검 대상 유형
- staging URL
- OpenAPI/Swagger 파일 또는 URL
- 컨테이너 이미지 이름
- 인증 방식
- allowlist 도메인
- 제외 경로
- 기본 프로파일
- CI 실패 기준

### 2. 도구 구성

`scanrail setup`은 workspace와 필요한 도구 이미지를 준비합니다.

준비 항목:

- Docker 실행 상태 확인
- 결과/캐시 디렉터리 생성
- pinned scanner image pull 또는 검증
- `tools.lock.yaml` 생성

현재 MVP는 native security headers scanner, native TLS certificate baseline scanner, Docker 기반 Gitleaks secrets adapter를 구현합니다. Docker 기반 Trivy와 Semgrep adapter는 planned v0.x surface로 남아 있습니다.

### 3. 스캔 프로파일

```text
quick    PR 전 빠른 점검
standard staging 기준 권장 점검
full     active scan 포함, 명시 승인 필요
```

### 4. 결과 통합

각 도구의 결과를 공통 finding 모델로 변환합니다.

공통 필드:

- id
- title
- description
- severity
- confidence
- source tool
- affected target
- evidence
- remediation
- CWE
- OWASP category
- CVSS
- EPSS
- KEV 여부

### 5. 리포트

지원 포맷:

- HTML: 개발자 읽기용
- JSON: 외부/내부 시스템 연동용
- SARIF: GitHub/GitLab code scanning 연동용, 이후 릴리스에서 안정화

## 성공 기준

- 신규 개발자가 10분 안에 첫 리포트를 생성할 수 있다.
- 보안 도구별 설치 없이 Docker만으로 실행 가능하다.
- PR에서 high 이상 finding을 자동 차단할 수 있다.
- 스캔 설정이 코드 저장소에 남아 재현 가능하다.
- 토큰과 비밀번호는 설정 파일에 저장되지 않는다.

## 주요 리스크

- 오픈소스 도구별 라이선스와 재배포 조건 검토 필요
- 스캐너 이미지 업데이트로 인한 결과 변동
- DAST active scan의 서비스 영향 가능성
- 인증 세션 처리 복잡도
- 오탐/중복 finding 관리 비용
- 비즈니스 로직 취약점 자동화 한계
- 공개 OSS에 내부 정책과 조직 전용 URL이 섞이지 않도록 구성 경계 필요
