import { readFileSync } from "node:fs";
import { join } from "node:path";
import { spawnSync } from "node:child_process";

const packages = [
  "packages/npm/cli-darwin-arm64",
  "packages/npm/cli-darwin-x64",
  "packages/npm/cli-win32-x64",
  "packages/npm/cli-win32-arm64",
  "packages/npm/cli-linux-x64",
  "packages/npm/cli-linux-arm64",
  "packages/npm/cli"
];

const args = new Set(process.argv.slice(2));
const publish = args.has("--publish");
const dryRun = !publish;
const provenance = args.has("--provenance");
const tag = readOption("--tag") || "latest";

if (publish && process.env.SCANRAIL_ALLOW_NPM_PUBLISH !== "1") {
  fail("Refusing real npm publish without SCANRAIL_ALLOW_NPM_PUBLISH=1.");
}

const manifests = packages.map((dir) => readManifest(dir));
const expectedVersion = manifests.at(-1).version;
for (const manifest of manifests) {
  if (manifest.version !== expectedVersion) {
    fail(`${manifest.name} version ${manifest.version} does not match wrapper version ${expectedVersion}.`);
  }
  if (manifest.publishConfig?.access !== "public") {
    fail(`${manifest.name} must set publishConfig.access to public.`);
  }
}

console.log(`${dryRun ? "Dry-running" : "Publishing"} ${manifests.length} npm packages as ${expectedVersion} with tag ${tag}.`);

for (const manifest of manifests) {
  checkRegistry(manifest);
  publishPackage(manifest);
}

function readOption(name) {
  const index = process.argv.indexOf(name);
  if (index === -1) return null;
  return process.argv[index + 1] || null;
}

function readManifest(dir) {
  const path = join(process.cwd(), dir, "package.json");
  return {
    dir,
    ...JSON.parse(readFileSync(path, "utf8"))
  };
}

function checkRegistry(manifest) {
  const spec = `${manifest.name}@${manifest.version}`;
  const result = spawnSync("npm", ["view", spec, "version", "--json", "--registry", "https://registry.npmjs.org/"], {
    encoding: "utf8"
  });
  const output = `${result.stdout || ""}\n${result.stderr || ""}`;
  if (result.status === 0) {
    fail(`${spec} is already published.`);
  }
  if (!output.includes("E404")) {
    const message = `Could not confirm ${spec} is unpublished.`;
    if (dryRun) {
      console.warn(`WARN: ${message}`);
      return;
    }
    fail(message);
  }
}

function publishPackage(manifest) {
  const publishArgs = ["publish", "--access", "public", "--tag", tag];
  if (dryRun) {
    publishArgs.push("--dry-run");
  }
  if (provenance && !dryRun) {
    publishArgs.push("--provenance");
  }

  console.log(`\n${dryRun ? "Dry-run" : "Publish"} ${manifest.name}`);
  run("npm", publishArgs, manifest.dir);
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

function fail(message) {
  console.error(message);
  process.exit(1);
}
