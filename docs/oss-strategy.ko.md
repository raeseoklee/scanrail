[ENGLISH](oss-strategy.md) | [한국어](oss-strategy.ko.md)

# OSS 전략

`scanrail`은 특정 회사 내부 도구가 아니라, 상용 보안진단 제품 도입이 어렵거나 여러 오픈소스 보안 도구를 직접 조합해야 하는 개발자를 위한 공개 OSS로 설계합니다.

## 포지셔닝

한 줄 설명:

```text
Docker 기반 오픈소스 보안 스캐너들을 하나의 개발자 친화 CLI로 실행하고, 결과를 정규화해 리포트하는 보안진단 오케스트레이터.
```

중요한 메시지:

- 보안 전문가 전용 툴이 아니라 개발자 셀프체크 툴입니다.
- 상용 DAST/SAST 제품 전체를 대체한다고 주장하지 않습니다.
- 취약점 탐지 엔진을 직접 재구현하지 않고 검증된 오픈소스를 조합합니다.
- 안전한 기본값을 제공하고, 위험한 scan은 명시적으로 opt-in합니다.
- 조직별 정책은 공개 코어 위에 설정과 plugin으로 얹습니다.

## 대상 사용자

- 상용 보안 제품 도입이 어려운 개발팀
- PR 전에 보안 셀프체크를 추가하려는 오픈소스 프로젝트
- 여러 보안 도구를 CI에 붙이고 싶은 플랫폼 팀
- 스캐너 결과를 SARIF/HTML/JSON으로 통합하고 싶은 보안 담당자
- staging API를 안전하게 반복 점검하고 싶은 백엔드 개발자

## 공개 코어와 조직 확장

공개 OSS에 포함할 것:

- CLI
- 인터랙티브 설정
- Docker scanner adapter
- 공통 finding model
- HTML/JSON/SARIF 리포트
- safe 기본 profile
- 공개 문서와 예제
- GitHub/GitLab CI 템플릿

조직별로 외부에서 주입할 것:

- 조직 registry mirror
- 조직 표준 profile
- 조직 보안 정책
- 비공개 Semgrep/Nuclei 룰셋
- Slack/Jira endpoint
- asset criticality 정보
- internal allowlist template

## 라이선스 후보

권장 후보:

```text
Apache-2.0
```

이유:

- 기업 사용자가 채택하기 쉽습니다.
- 특허 조항이 있어 조직 도입에 유리합니다.
- CLI, SDK, plugin 생태계와 잘 맞습니다.

대안:

- MIT: 단순하지만 특허 보호가 약합니다.
- MPL-2.0: 파일 단위 copyleft가 필요할 때 적합합니다.
- AGPL-3.0: SaaS 상용화를 강하게 제한하려는 경우 적합하지만 기업 채택 장벽이 높습니다.

초기에는 Apache-2.0이 가장 현실적입니다.

## 저장소 구성

```text
.
├─ cmd/scanrail/
├─ internal/
│  ├─ config/
│  ├─ dockerx/
│  ├─ scanners/
│  ├─ findings/
│  ├─ policy/
│  └─ report/
├─ packages/
│  └─ npm/
│     ├─ cli/
│     ├─ cli-darwin-arm64/
│     ├─ cli-darwin-x64/
│     ├─ cli-win32-x64/
│     ├─ cli-win32-arm64/
│     ├─ cli-linux-x64/
│     └─ cli-linux-arm64/
├─ docs/
├─ examples/
└─ .github/
```

## 커뮤니티 운영 원칙

- 기본 profile은 안전해야 합니다.
- destructive scan은 기본 profile에 넣지 않습니다.
- exploit payload 추가 PR은 별도 검토 기준을 둡니다.
- scanner 결과를 과장하지 않습니다.
- false positive와 noisy rule은 명확하게 관리합니다.
- 보안 취약점 제보는 public issue가 아니라 security policy로 받습니다.

## MVP 공개 기준

첫 공개 릴리스는 다음을 만족해야 합니다.

- `npm install -g`로 설치 가능
- macOS, Windows, Linux에서 `scanrail doctor` 동작
- `scanrail init`, `setup`, `run` 기본 동작
- 최소 3개 scanner adapter 동작
- HTML/JSON/SARIF 리포트 생성
- active scan opt-in 보호
- README quickstart 완성
- LICENSE, SECURITY.md, CONTRIBUTING.md 존재

## 피해야 할 것

- 내부 회사명, 도메인, registry URL을 기본값에 넣지 않습니다.
- scanner output을 과장해 “완전한 보안진단”이라고 표현하지 않습니다.
- 무분별한 공격성 template을 기본 활성화하지 않습니다.
- postinstall에서 외부 바이너리를 내려받는 설치 방식을 기본으로 삼지 않습니다.
- private API key나 telemetry를 요구하지 않습니다.
