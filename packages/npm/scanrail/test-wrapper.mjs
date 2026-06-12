import { mkdtempSync, mkdirSync, writeFileSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

const dir = mkdtempSync(join(tmpdir(), "scanrail-alias-wrapper-"));
try {
  const cliBinDir = join(dir, "node_modules", "@scanrail", "cli", "bin");
  mkdirSync(cliBinDir, { recursive: true });
  writeFileSync(
    join(cliBinDir, "scanrail.js"),
    "console.log('alias-ok ' + process.argv.slice(2).join(' ')); process.exit(9);\n"
  );

  const wrapper = fileURLToPath(new URL("./bin/scanrail.js", import.meta.url));
  const result = spawnSync(process.execPath, [wrapper, "doctor"], {
    env: { ...process.env, NODE_PATH: join(dir, "node_modules") },
    encoding: "utf8"
  });
  if (result.status !== 9) {
    console.error(result.stdout);
    console.error(result.stderr);
    throw new Error(`expected exit 9, got ${result.status}`);
  }
  if (!result.stdout.includes("alias-ok doctor")) {
    console.error(result.stdout);
    throw new Error("alias wrapper did not delegate arguments");
  }
} finally {
  rmSync(dir, { recursive: true, force: true });
}
