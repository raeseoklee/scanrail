[ENGLISH](CHANGELOG.md) | [한국어](CHANGELOG.ko.md)

# 변경 기록

## 0.1.1 - 2026-06-12

- npm Trusted Publishing으로 기능 변경 없는 patch release를 publish합니다.
- `@scanrail/cli`를 내부 wrapper로 유지하고, 사용자가 설치하는 기본 npm entrypoint는 `scanrail`로 둡니다.

## 0.1.0 - 2026-06-12

- 초기 Go CLI scaffold 추가.
- `doctor`, `init`, `setup`, `run` 명령 추가.
- native security headers scanner 추가.
- JSON 및 HTML report 생성 추가.
- npm wrapper package와 platform binary package manifest 추가.
- unscoped npm alias package `scanrail` 추가.
- cross-platform release build dry-run script 추가.
- 기본 문서는 영문, 한국어 문서는 `.ko.md`로 제공하는 bilingual 문서 구조 추가.
