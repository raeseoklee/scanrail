[ENGLISH](cli-reference.md) | [한국어](cli-reference.ko.md)

# CLI 명세

이 문서는 `scanrail` CLI의 초기 명령 구조를 정의합니다.

## 명령 목록

```text
scanrail doctor
scanrail init
scanrail setup
scanrail run
scanrail mcp serve
scanrail update
scanrail ci init
scanrail auth setup
scanrail report open
```

현재 구현된 명령은 `doctor`, `init`, `setup`, `run`, `version`, `mcp serve`입니다.

## 공통 exit code

모든 명령은 같은 의미의 exit code를 사용합니다.

```text
0    성공
1    보안 policy 실패 또는 진단 finding 기준 실패
2    사용법/설정 오류
3    실행 환경 오류
4    scanner 실행 또는 파싱 실패
5    안전 정책 위반으로 실행 거부
130  사용자 중단
```

`doctor`는 보안 finding을 만들지 않으므로 환경 불충족은 `3`, 설정 오류는 `2`를 사용합니다.

## scanrail doctor

로컬 실행 환경을 점검합니다.

```bash
scanrail doctor
```

점검 항목:

- Docker daemon 실행 여부
- Docker Compose 사용 가능 여부
- 현재 디렉터리 접근 가능 여부
- Git repository 여부
- 디스크 여유 공간
- scanner image pull 가능 여부
- `scanrail.yaml` 유효성
- `tools.lock.yaml` 유효성

exit code:

```text
0  실행 가능
3  필수 조건 불충족
2  설정 오류
```

## scanrail init

인터랙티브 질문을 통해 `scanrail.yaml`을 생성합니다.

```bash
scanrail init
```

옵션:

```text
--non-interactive       질문 없이 기본값과 옵션으로 설정 생성
--project-name <name>   프로젝트 이름
--target <url>          웹 대상 URL
--openapi <path>        local OpenAPI spec 경로
--profile <name>        기본 profile
--force                 기존 scanrail.yaml 덮어쓰기
```

보호 동작:

- 기존 `scanrail.yaml`이 있으면 기본적으로 덮어쓰지 않습니다.
- secret 원문 입력을 설정 파일에 저장하지 않습니다.
- 운영 URL로 보이는 대상은 경고하고 별도 확인을 요구합니다.

## scanrail setup

`.scanrail` workspace와 Docker 기반 스캐너 실행 환경을 준비합니다.

```bash
scanrail setup
```

옵션:

```text
--pull-policy missing|always|never
--registry <registry-url>
--profile <name>
--offline
```

작업:

- `.scanrail/cache` 생성
- `.scanrail/reports` 생성
- pinned scanner image pull
- `tools.lock.yaml` 생성 또는 갱신

Gitleaks는 `ghcr.io/gitleaks/gitleaks:v8.30.1` 이미지를 사용합니다. Trivy와 Semgrep은 adapter가 구현될 때까지 `tools.lock.yaml`에 placeholder로 남기고 pull은 skip합니다.

## scanrail run

설정된 profile로 보안 점검을 실행합니다.

```bash
scanrail run
scanrail run --profile quick
scanrail run --only headers
scanrail run --only gitleaks
scanrail run --only tls
scanrail run --only openapi --openapi ./openapi.yaml
```

확장 profile을 설정한 경우:

```bash
scanrail run --profile standard
scanrail run --profile full --i-understand-active-scan
```

옵션:

```text
--profile <name>                  실행할 profile
--target <url>                    일회성 target override
--openapi <path>                  일회성 local OpenAPI spec override
--output-dir <path>               report 출력 경로
--format html,json,sarif          report format override
--fail-on critical|high|medium    policy fail 기준
--i-understand-active-scan        active scan 명시 허용
--no-cache                        scanner cache 사용 안 함
--strict-lock                     tools.lock.yaml digest 불일치를 실패 처리
--only <tool>                     특정 도구만 실행
--skip <tool>                     특정 도구 제외
```

exit code:

```text
0  스캔 완료, policy 통과
1  스캔 완료, policy 실패
2  설정 오류
3  실행 환경 오류
4  scanner 실행 실패
5  안전 정책 위반으로 실행 거부
```

실행 동작:

- `gitleaks`는 Docker가 필요하며 local workspace를 read-only bind mount로 스캔합니다.
- `headers`는 native Go scanner라 Docker가 없어도 실행할 수 있습니다.
- `tls`는 native Go scanner이며 HTTPS target에 단일 TLS handshake를 수행합니다.
- `openapi`는 native Go scanner이며 local OpenAPI JSON 또는 일반적인 YAML file을 읽고 API endpoint를 호출하지 않습니다.
- profile에 포함된 scanner가 실행 조건을 만족하지 못하면 skip reason을 report에 남깁니다.
- `--only`로 명시한 scanner가 실행 조건을 만족하지 못하면 실패합니다.

## scanrail mcp serve

local stdio MCP server를 시작합니다.

```bash
scanrail mcp serve
```

구현된 MCP method:

- `initialize`
- `ping`
- `tools/list`
- `tools/call`
- `resources/list`
- `resources/read`

노출 tool:

- `scanrail_doctor`
- `scanrail_config_read`
- `scanrail_report_latest`
- `scanrail_run`

안전 동작:

- stdio only, local HTTP listener 없음
- arbitrary shell execution 없음
- MCP MVP의 `scanrail_run`은 native `headers` scanner만 지원
- active scan 실행은 `confirm_active_scan=true` 필요
- target host는 configured target host 또는 `targets.web.allowlist`와 일치해야 함
- MCP-triggered scan attempt는 `.scanrail/logs/mcp-audit.jsonl`에 기록됨

## scanrail update

scanner 이미지, DB, rule, template을 업데이트합니다.

```bash
scanrail update
```

옵션:

```text
--tools
--rules
--templates
--db
--lock
```

정책:

- private 또는 organization registry mirror를 사용하는 경우 mirror 기준으로만 업데이트합니다.
- 승인된 버전 정책이 있으면 그 범위를 벗어나지 않습니다.

## scanrail ci init

CI 설정 파일을 생성합니다.

```bash
scanrail ci init
```

지원 후보:

- GitHub Actions
- GitLab CI
- Jenkins

생성물:

- CI workflow 파일
- cache 설정
- secret 환경변수 안내
- SARIF 업로드 설정

## scanrail auth setup

인증 설정을 보조합니다.

```bash
scanrail auth setup
```

지원 방식:

- none
- bearer token
- cookie
- form login
- recorded browser session

원칙:

- secret 원문은 저장하지 않습니다.
- 환경변수 이름 또는 CI secret 이름만 저장합니다.
- 인증 실패 시 scanner를 실행하지 않고 명확한 오류를 제공합니다.

## scanrail report open

최근 리포트를 엽니다.

```bash
scanrail report open
scanrail report open --latest
scanrail report open --path .scanrail/reports/report.html
```
