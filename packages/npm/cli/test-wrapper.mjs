import { mkdtempSync, writeFileSync, chmodSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { spawnSync } from "node:child_process";

const dir = mkdtempSync(join(tmpdir(), "scanrail-wrapper-"));
try {
  const fake = join(dir, process.platform === "win32" ? "fake.cmd" : "fake");
  if (process.platform === "win32") {
    writeFileSync(fake, "@echo off\r\necho wrapper-ok %*\r\nexit /b 7\r\n");
  } else {
    writeFileSync(fake, "#!/bin/sh\necho wrapper-ok \"$@\"\nexit 7\n");
    chmodSync(fake, 0o755);
  }
  const result = spawnSync(process.execPath, [new URL("./bin/scanrail.js", import.meta.url).pathname, "doctor"], {
    env: { ...process.env, SCANRAIL_BINARY_PATH: fake },
    encoding: "utf8"
  });
  if (result.status !== 7) {
    console.error(result.stdout);
    console.error(result.stderr);
    throw new Error(`expected exit 7, got ${result.status}`);
  }
  if (!result.stdout.includes("wrapper-ok doctor")) {
    console.error(result.stdout);
    throw new Error("wrapper did not forward arguments");
  }
} finally {
  rmSync(dir, { recursive: true, force: true });
}
