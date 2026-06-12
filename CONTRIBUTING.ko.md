[ENGLISH](CONTRIBUTING.md) | [한국어](CONTRIBUTING.ko.md)

# 기여 가이드

Scanrail 개선에 기여해 주셔서 감사합니다. 아직 초기 프로젝트라서 실제 개발자 워크플로우에 연결된 작고 재현 가능한 기여가 가장 가치 있습니다.

## 좋은 첫 기여

- 최소 fixture로 scanner adapter 문제 재현
- 설치, 설정, 리포트 문서 개선
- config parsing, report generation, MCP JSON-RPC 동작 테스트 추가
- 새 adapter 구현 전에 scanner adapter contract 제안

## 개발 환경

```bash
git clone https://github.com/raeseoklee/scanrail.git
cd scanrail
go test ./...
npm test
make release-dry-run
```

## Pull Request 체크리스트

- 변경 범위를 작고 되돌리기 쉽게 유지합니다.
- 동작이 바뀌면 테스트를 추가하거나 수정합니다.
- 사용자 대상 문서는 영문을 먼저 수정하고, 한국어 문서는 `.ko.md`로 추가합니다.
- 리뷰 요청 전 `npm run docs:check-links`, `npm test`, `make release-dry-run`을 실행합니다.
- 표준 라이브러리나 기존 도구로 충분하지 않은 이유가 명확하지 않다면 새 dependency를 추가하지 않습니다.

## Scanner Adapter 변경

scanner adapter를 추가하거나 변경하기 전에 다음을 문서화합니다.

- scanner version 및 container image
- 정확한 command line
- 필요한 input과 secret
- target safety control
- raw output 위치
- normalized finding mapping
- 알려진 gap과 skipped behavior

## Commit Message

commit message는 intent-first로 작성합니다. 변경 이유, 제약, 검증 내용을 남깁니다.

## 보안

공개 issue나 pull request에 exploit 세부정보, 실제 secret, 고객 target URL, private scan result를 포함하지 마세요. [SECURITY.ko.md](SECURITY.ko.md)를 참고하세요.
