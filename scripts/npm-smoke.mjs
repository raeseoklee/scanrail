import { mkdtempSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { spawnSync } from "node:child_process";

const requestedVersion = process.argv[2] || process.env.SCANRAIL_NPM_SMOKE_VERSION || "latest";
const spec = requestedVersion.startsWith("scanrail@") ? requestedVersion : `scanrail@${requestedVersion}`;
const tmp = mkdtempSync(join(tmpdir(), "scanrail-npm-smoke-"));

try {
  console.log(`Smoke-testing ${spec} in ${tmp}`);
  run("npm", ["init", "-y"], tmp);
  run("npm", ["install", spec, "--registry", "https://registry.npmjs.org/", "--prefer-online"], tmp);
  run("npm", ["exec", "--", "scanrail", "version"], tmp);
  run("npm", ["exec", "--", "scanrail", "doctor"], tmp);
  run("npm", ["audit", "signatures", "--registry", "https://registry.npmjs.org/"], tmp);
} finally {
  rmSync(tmp, { recursive: true, force: true });
}

function run(command, args, cwd) {
  const result = spawnSync(command, args, {
    cwd,
    stdio: "inherit",
    env: process.env
  });
  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}
