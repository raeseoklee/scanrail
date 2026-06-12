import { existsSync, readdirSync, readFileSync } from "node:fs";
import { dirname, extname, resolve } from "node:path";

const root = process.cwd();
const ignoredDirs = new Set([".git", ".omc", ".omx", ".scanrail", ".codexus", "dist", "node_modules"]);
const markdownFiles = walk(root).filter((file) => extname(file) === ".md");
let failures = 0;

for (const file of markdownFiles) {
  const text = readFileSync(file, "utf8");
  const links = text.matchAll(/!?\[[^\]]*\]\(([^)]+)\)/g);
  for (const match of links) {
    const target = parseTarget(match[1]);
    if (!target || shouldSkip(target)) continue;

    const absoluteTarget = resolve(dirname(file), target);
    if (!existsSync(absoluteTarget)) {
      console.error(`${relativePath(file)}: missing link target ${target}`);
      failures += 1;
    }
  }
}

if (failures > 0) {
  process.exit(1);
}

console.log(`markdown links ok (${markdownFiles.length} files)`);

function walk(dir) {
  const files = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    if (entry.isDirectory()) {
      if (!ignoredDirs.has(entry.name)) {
        files.push(...walk(resolve(dir, entry.name)));
      }
      continue;
    }
    if (entry.isFile()) {
      files.push(resolve(dir, entry.name));
    }
  }
  return files;
}

function parseTarget(value) {
  const trimmed = value.trim();
  if (trimmed.startsWith("<") && trimmed.includes(">")) {
    return trimmed.slice(1, trimmed.indexOf(">"));
  }
  return trimmed.split(/\s+/)[0].split("#")[0];
}

function shouldSkip(target) {
  return (
    target === "" ||
    target.startsWith("#") ||
    target.startsWith("mailto:") ||
    /^[a-z][a-z0-9+.-]*:/i.test(target)
  );
}

function relativePath(file) {
  return file.replace(`${root}/`, "");
}
