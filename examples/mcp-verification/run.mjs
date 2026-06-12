import { mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { spawn, spawnSync } from "node:child_process";

const version = process.env.SCANRAIL_VERIFY_VERSION || "0.1.2";
const port = process.env.SCANRAIL_VERIFY_PORT || "18182";
const target = `http://127.0.0.1:${port}`;
const workdir = mkdtempSync(join(tmpdir(), "scanrail-mcp-verify-"));
const server = spawn(process.execPath, ["examples/headers-demo/server.mjs"], {
  cwd: process.cwd(),
  env: { ...process.env, PORT: port },
  stdio: ["ignore", "pipe", "pipe"]
});

process.on("exit", cleanup);
process.on("SIGINT", () => {
  cleanup();
  process.exit(130);
});

try {
  await wait(1000);
  console.log(`Scanrail MCP verification target: ${target}`);
  console.log(`Temporary project: ${workdir}`);
  console.log("");

  run("npm", ["init", "-y"], workdir, { quiet: true });
  console.log(`Installing scanrail@${version} from npm...`);
  run("npm", ["install", `scanrail@${version}`, "--registry", "https://registry.npmjs.org/", "--prefer-online"], workdir);

  const scanrail = join(workdir, "node_modules", ".bin", process.platform === "win32" ? "scanrail.cmd" : "scanrail");
  run(scanrail, ["version"], workdir);
  run(scanrail, ["init", "--non-interactive", "--project-name", "mcp-demo", "--target", target, "--force"], workdir);

  const input = [
    rpc(1, "initialize", {
      protocolVersion: "2025-06-18",
      capabilities: {},
      clientInfo: { name: "scanrail-mcp-tape", version: "1" }
    }),
    JSON.stringify({ jsonrpc: "2.0", method: "notifications/initialized" }),
    rpc(2, "tools/list"),
    rpc(3, "resources/list"),
    rpc(4, "resources/read", { uri: "scanrail://config" }),
    rpc(5, "tools/call", {
      name: "scanrail_run",
      arguments: { only: "headers", target }
    }),
    rpc(6, "tools/call", {
      name: "scanrail_run",
      arguments: { only: "headers", target, confirm_active_scan: true }
    }),
    rpc(7, "tools/call", {
      name: "scanrail_report_latest",
      arguments: {}
    })
  ].join("\n") + "\n";

  writeFileSync(join(workdir, "mcp-input.ndjson"), input);
  const mcp = spawnSync(scanrail, ["mcp", "serve"], {
    cwd: workdir,
    input,
    encoding: "utf8",
    shell: process.platform === "win32"
  });
  if (mcp.status !== 0) {
    process.stdout.write(mcp.stdout || "");
    process.stderr.write(mcp.stderr || "");
    process.exit(mcp.status ?? 1);
  }
  writeFileSync(join(workdir, "mcp-output.ndjson"), mcp.stdout);

  const responses = mcp.stdout.trim().split("\n").map((line) => JSON.parse(line));
  const byId = new Map(responses.map((response) => [response.id, response]));
  const tools = byId.get(2).result.tools.map((tool) => tool.name).sort();
  const resources = byId.get(3).result.resources.map((resource) => resource.uri).sort();
  const denied = byId.get(5).result;
  const runResult = byId.get(6).result;
  const latest = byId.get(7).result.structuredContent;

  assert(byId.get(1).result.protocolVersion === "2025-06-18", "initialize failed");
  assert(tools.includes("scanrail_run"), "scanrail_run tool missing");
  assert(resources.includes("scanrail://reports/latest/summary"), "latest report resource missing");
  assert(byId.get(4).result.contents[0].text.includes(target), "config resource missing target");
  assert(denied.isError === true, "scan without confirmation was not denied");
  assert(runResult.structuredContent.exit_code === 0, "confirmed headers scan failed");
  assert(latest.findings_count >= 1, "latest report did not include findings");

  console.log("");
  console.log("MCP verification results");
  console.log(`PASS initialize -> ${byId.get(1).result.serverInfo.version}`);
  console.log(`PASS tools/list -> ${tools.join(", ")}`);
  console.log(`PASS resources/list -> ${resources.join(", ")}`);
  console.log("PASS scanrail_run without confirmation was denied");
  console.log(`PASS scanrail_run with confirmation completed with exit ${runResult.structuredContent.exit_code}`);
  console.log(`PASS latest report summary -> ${latest.findings_count} findings`);
  console.log("");
  console.log("Scanrail MCP verification passed.");
} finally {
  cleanup();
}

function rpc(id, method, params = undefined) {
  return JSON.stringify({ jsonrpc: "2.0", id, method, ...(params === undefined ? {} : { params }) });
}

function run(command, args, cwd, options = {}) {
  const result = spawnSync(command, args, {
    cwd,
    stdio: options.quiet ? "ignore" : "inherit",
    shell: process.platform === "win32"
  });
  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

function assert(condition, message) {
  if (!condition) {
    throw new Error(message);
  }
}

function wait(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function cleanup() {
  if (!server.killed) {
    server.kill();
  }
  rmSync(workdir, { recursive: true, force: true });
}
