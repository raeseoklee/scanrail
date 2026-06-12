#!/usr/bin/env node
const { spawnSync } = require("node:child_process");

const platform = process.platform;
const arch = process.arch;
const suffix = `${platform}-${arch}`;
const packageName = `@scanrail/cli-${suffix}`;
const binaryName = platform === "win32" ? "scanrail.exe" : "scanrail";

let binary = process.env.SCANRAIL_BINARY_PATH;
if (!binary) {
  try {
    binary = require.resolve(`${packageName}/${binaryName}`);
  } catch {
    console.error(`Unsupported platform or missing package: ${platform}/${arch}`);
    console.error(`Expected package: ${packageName}`);
    process.exit(1);
  }
}

const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

if (result.signal) {
  const signalOffset = 128;
  const signalNumbers = { SIGINT: 2, SIGTERM: 15 };
  process.exit(signalOffset + (signalNumbers[result.signal] ?? 1));
}

process.exit(result.status ?? 1);
