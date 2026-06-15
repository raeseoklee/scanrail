import { createServer } from "node:http";
import { mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, join, resolve } from "node:path";
import { spawn, spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(here, "../..");
const workdir = mkdtempSync(join(tmpdir(), "scanrail-mcp-workbench-"));
const binary = join(workdir, process.platform === "win32" ? "scanrail.exe" : "scanrail");
const host = "127.0.0.1";
const requestedPort = Number.parseInt(process.env.SCANRAIL_WORKBENCH_PORT || "0", 10);
const signalExitCodes = new Map([
  ["SIGINT", 130],
  ["SIGTERM", 143]
]);
let cleanedUp = false;
let shuttingDown = false;
let shutdownTimer;

const httpServer = createServer((req, res) => {
  if (req.url === "/secure") {
    res.setHeader("Content-Security-Policy", "default-src 'self'");
    res.setHeader("X-Content-Type-Options", "nosniff");
    res.setHeader("X-Frame-Options", "DENY");
    res.setHeader("Referrer-Policy", "no-referrer");
  }
  res.setHeader("Content-Type", "text/plain; charset=utf-8");
  res.end("Scanrail MCP Workbench target\n");
});

await listen(httpServer, requestedPort, host);
const port = httpServer.address().port;
const target = `http://${host}:${port}`;
writeConfig(join(workdir, "scanrail.yaml"), target);
buildScanrail();

const child = spawn(binary, ["mcp", "serve"], {
  cwd: workdir,
  stdio: ["pipe", "pipe", "pipe"]
});

process.stdin.pipe(child.stdin);
child.stdout.pipe(process.stdout);
child.stderr.pipe(process.stderr);

child.on("exit", (code, signal) => {
  if (shutdownTimer) {
    clearTimeout(shutdownTimer);
  }
  cleanup();
  if (shuttingDown) {
    process.exit(0);
  }
  if (signal) {
    process.exit(signalExitCodes.get(signal) ?? 1);
  }
  process.exit(code ?? 0);
});

for (const signal of ["SIGINT", "SIGTERM"]) {
  process.on(signal, () => {
    shuttingDown = true;
    child.kill(signal);
    cleanup();
    shutdownTimer = setTimeout(() => process.exit(0), 1000);
    shutdownTimer.unref();
  });
}

function buildScanrail() {
  const result = spawnSync("go", ["build", "-o", binary, "./cmd/scanrail"], {
    cwd: repoRoot,
    stdio: ["ignore", "ignore", "inherit"]
  });
  if (result.status !== 0) {
    cleanup();
    process.exit(result.status ?? 1);
  }
}

function writeConfig(path, url) {
  writeFileSync(path, `project:
  name: workbench-mcp

targets:
  web:
    url: ${url}
    allowlist:
      - ${host}:${port}

auth:
  token_env: SCANRAIL_TOKEN

safety:
  active_scan_default: false
  require_allowlist: true

policy:
  fail_on:
    severity: high

report:
  output_dir: .scanrail/reports
`, "utf8");
}

function listen(server, listenPort, listenHost) {
  return new Promise((resolveListen, reject) => {
    server.once("error", reject);
    server.listen(listenPort, listenHost, resolveListen);
  });
}

function cleanup() {
  if (cleanedUp) {
    return;
  }
  cleanedUp = true;
  httpServer.close();
  rmSync(workdir, { recursive: true, force: true });
}
