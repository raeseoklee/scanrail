import { copyFileSync, mkdirSync, rmSync } from "node:fs";
import { join } from "node:path";
import { spawnSync } from "node:child_process";

const targets = [
  ["darwin", "arm64", "cli-darwin-arm64", "scanrail"],
  ["darwin", "amd64", "cli-darwin-x64", "scanrail"],
  ["windows", "amd64", "cli-win32-x64", "scanrail.exe"],
  ["windows", "arm64", "cli-win32-arm64", "scanrail.exe"],
  ["linux", "amd64", "cli-linux-x64", "scanrail"],
  ["linux", "arm64", "cli-linux-arm64", "scanrail"]
];

run("go", ["test", "./..."]);
rmSync("dist", { recursive: true, force: true });
mkdirSync("dist", { recursive: true });

for (const [goos, goarch, pkg, binaryName] of targets) {
  const outDir = join("dist", `${goos}-${goarch}`);
  mkdirSync(outDir, { recursive: true });
  const out = join(outDir, binaryName);
  run("go", ["build", "-trimpath", "-ldflags", ldflags(), "-o", out, "./cmd/scanrail"], {
    GOOS: goos,
    GOARCH: goarch
  });
  copyFileSync(out, join("packages", "npm", pkg, binaryName));
}

function run(command, args, env = {}) {
  const result = spawnSync(command, args, {
    stdio: "inherit",
    env: { ...process.env, ...env }
  });
  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

function ldflags() {
  const version = process.env.SCANRAIL_VERSION || "0.1.0";
  const commit = process.env.SCANRAIL_COMMIT || "snapshot";
  const date = new Date().toISOString();
  return [
    `-X github.com/raeseoklee/scanrail/internal/version.Version=${version}`,
    `-X github.com/raeseoklee/scanrail/internal/version.Commit=${commit}`,
    `-X github.com/raeseoklee/scanrail/internal/version.Date=${date}`
  ].join(" ");
}
