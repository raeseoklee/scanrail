[ENGLISH](SECURITY.md) | [한국어](SECURITY.ko.md)

# 보안 정책

Scanrail은 보안 도구이므로 취약점 제보는 신중한 공개 절차가 필요합니다.

## 지원 버전

Scanrail은 아직 pre-1.0입니다. 최신 npm 배포 버전과 현재 `main` branch만 보안 수정을 받습니다.

## 취약점 제보

exploit 세부정보, 실제 target URL, raw secret, private scan output을 공개 GitHub issue에 올리지 마세요.

권장 제보 경로:

1. 이 저장소에서 GitHub private vulnerability reporting을 사용할 수 있으면 그 경로를 사용합니다.
2. 사용할 수 없다면 private disclosure channel을 요청하는 최소 공개 issue만 엽니다. exploit 세부정보는 포함하지 않습니다.

포함하면 좋은 정보:

- 영향받는 Scanrail version
- 운영체제와 설치 방식
- 영향 요약
- synthetic target 또는 sanitized fixture 기반 재현 절차
- unsafe scanning, secret exposure, command execution, report tampering 가능 여부

## 범위

대상:

- command injection 또는 arbitrary process execution
- target allowlist 우회
- config, report, log, MCP tool, npm wrapper를 통한 secret leakage
- 명시적 사용자 의도 없이 target을 scan하는 unsafe default
- package integrity 또는 provenance 문제

비대상:

- Scanrail이 위험을 키우지 않는 third-party scanner 자체 취약점
- Scanrail이 unrelated application에서 발견한 일반 finding
- Scanrail-specific impact가 없는 local demo target denial-of-service

## 공개

유효한 제보는 7일 안에 확인하고, 공개 전 fix 또는 mitigation을 조율하는 것을 목표로 합니다.
