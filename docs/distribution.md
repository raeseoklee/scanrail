# 배포 전략

`scanrail`은 Go로 CLI 본체를 만들고, npm wrapper를 1차 배포 채널로 사용합니다. 목표는 macOS, Windows, Linux 개발자가 같은 명령으로 설치하고 실행할 수 있게 하는 것입니다.

## 권장 배포 모델

```text
Go binary
  |
  +-- npm wrapper package
  +-- platform-specific npm binary packages
  +-- GitHub/GitLab release archives
  +-- optional Homebrew/Scoop packages
```

## npm 패키지 구조

```text
@scanrail/cli
@scanrail/cli-darwin-arm64
@scanrail/cli-darwin-x64
@scanrail/cli-win32-x64
@scanrail/cli-win32-arm64
@scanrail/cli-linux-x64
@scanrail/cli-linux-arm64
```

최종 package scope와 이름은 npm availability를 확인한 뒤 확정합니다.

## 사용자 설치 UX

전역 설치:

```bash
npm install -g @scanrail/cli
scanrail init
scanrail setup
scanrail run
```

일회성 실행:

```bash
npx @scanrail/cli run --profile quick
```

## wrapper package 예시

```json
{
  "name": "@scanrail/cli",
  "version": "0.1.0",
  "bin": {
    "scanrail": "./bin/scanrail.js"
  },
  "optionalDependencies": {
    "@scanrail/cli-darwin-arm64": "0.1.0",
    "@scanrail/cli-darwin-x64": "0.1.0",
    "@scanrail/cli-win32-x64": "0.1.0",
    "@scanrail/cli-win32-arm64": "0.1.0",
    "@scanrail/cli-linux-x64": "0.1.0",
    "@scanrail/cli-linux-arm64": "0.1.0"
  }
}
```

## platform package 예시

```json
{
  "name": "@scanrail/cli-win32-x64",
  "version": "0.1.0",
  "os": ["win32"],
  "cpu": ["x64"],
  "files": ["scanrail.exe"]
}
```

## wrapper 실행 방식

```js
#!/usr/bin/env node
const { spawnSync } = require("node:child_process");

const platform = process.platform;
const arch = process.arch;
const suffix = `${platform}-${arch}`;
const packageName = `@scanrail/cli-${suffix}`;
const binaryName = platform === "win32" ? "scanrail.exe" : "scanrail";

let binary;
try {
  binary = require.resolve(`${packageName}/${binaryName}`);
} catch {
  console.error(`Unsupported platform: ${platform}/${arch}`);
  process.exit(1);
}

const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

if (result.signal) {
  const signalOffset = 128;
  const signalNumbers = { SIGINT: 2, SIGTERM: 15 };
  process.exit(signalOffset + (signalNumbers[result.signal] ?? 1));
}

process.exit(result.status ?? 1);
```

## npmjs 공개 배포

공개 OSS로 배포할 경우:

```bash
npm publish --access public --provenance
```

원칙:

- scoped package는 최초 public publish 때 `--access public`을 명시합니다.
- CI 기반 trusted publishing 또는 provenance를 사용합니다.
- npm package에는 내부 URL, secret, 조직 전용 정책을 포함하지 않습니다.
- postinstall download 방식은 피합니다.
- platform package에는 `exports` 필드를 두지 않습니다. wrapper가 `require.resolve("<package>/<binary>")`로 바이너리 경로를 찾기 때문입니다.

## fallback 배포

npm을 쓸 수 없는 환경을 위해 release archive도 제공합니다.

```text
scanrail_darwin_arm64.tar.gz
scanrail_darwin_x64.tar.gz
scanrail_win32_x64.zip
scanrail_win32_arm64.zip
scanrail_linux_x64.tar.gz
scanrail_linux_arm64.tar.gz
checksums.txt
checksums.txt.sig
sbom.spdx.json
```

## 추가 채널

MVP 이후 검토:

- Homebrew tap
- Scoop bucket
- winget
- Chocolatey
- Docker image
- GitHub Action

## 서명 정책

초기 OSS 릴리스:

- npm provenance
- checksums
- SBOM

안정화 이후:

- Windows Authenticode signing
- macOS codesign/notarization
- SLSA provenance
- release artifact signature

현실적으로 npm wrapper는 설치 마찰을 낮추지만, 내부에 포함된 Go 바이너리 실행 자체가 모든 기업 EDR/OS 정책을 우회하는 것은 아닙니다. 조직 배포를 목표로 할수록 바이너리 서명은 결국 필요해질 수 있습니다.
