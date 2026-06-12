import { existsSync, mkdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { dirname, join, relative, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "../..");
const experimentDir = join(root, "experiments", "scanner-adapter-spike");
const fixtureDir = join(experimentDir, "fixture");
const outDir = join(experimentDir, "out");
const rawDir = join(outDir, "raw");
const cacheDir = join(outDir, "cache");
const reportPath = join(outDir, "scanrail-normalized-report.json");
const htmlPath = join(outDir, "scanrail-normalized-report.html");
const summaryPath = join(outDir, "summary.md");
const testSecret = "scanrail-demo-" + "secret-0001";

const images = {
  gitleaks: process.env.SCANRAIL_SPIKE_GITLEAKS_IMAGE || "ghcr.io/gitleaks/gitleaks:v8.30.1",
  semgrep: process.env.SCANRAIL_SPIKE_SEMGREP_IMAGE || "semgrep/semgrep:1.165.0",
  trivy: process.env.SCANRAIL_SPIKE_TRIVY_IMAGE || "aquasec/trivy:0.71.0"
};

if (process.argv.includes("--summary")) {
  printSummary();
  process.exit(0);
}

rmSync(outDir, { recursive: true, force: true });
mkdirSync(rawDir, { recursive: true });
mkdirSync(cacheDir, { recursive: true });

const commands = [];
commands.push(runScanner("gitleaks", [
  "docker", "run", "--rm",
  "-v", `${fixtureDir}:/src:ro`,
  "-v", `${rawDir}:/out`,
  images.gitleaks,
  "detect",
  "--source=/src",
  "--no-git",
  "--config=/src/gitleaks.toml",
  "--report-format=json",
  "--report-path=/out/gitleaks.json",
  "--redact=100",
  "--exit-code=7"
], [0, 7]));

commands.push(runScanner("semgrep", [
  "docker", "run", "--rm",
  "-v", `${fixtureDir}:/src:ro`,
  "-v", `${rawDir}:/out`,
  images.semgrep,
  "semgrep", "scan",
  "--config", "/src/semgrep.yml",
  "--json",
  "--output", "/out/semgrep.json",
  "/src"
], [0]));

commands.push(runScanner("trivy", [
  "docker", "run", "--rm",
  "-v", `${fixtureDir}:/src:ro`,
  "-v", `${rawDir}:/out`,
  "-v", `${cacheDir}:/root/.cache`,
  images.trivy,
  "config",
  "--format", "json",
  "--output", "/out/trivy.json",
  "/src"
], [0]));

const findings = [
  ...normalizeGitleaks(join(rawDir, "gitleaks.json")),
  ...normalizeSemgrep(join(rawDir, "semgrep.json")),
  ...normalizeTrivy(join(rawDir, "trivy.json"))
];

const report = {
  schema_version: 1,
  experiment: "scanner-adapter-spike",
  generated_at: new Date().toISOString(),
  fixture: relative(root, fixtureDir),
  commands,
  findings,
  summary: summarize(findings)
};

writeJSON(reportPath, report);
writeFileSync(htmlPath, renderHTML(report));
writeFileSync(summaryPath, renderMarkdownSummary(report));
printSummary(report);

function runScanner(tool, command, allowedExitCodes) {
  const startedAt = new Date().toISOString();
  const result = spawnSync(command[0], command.slice(1), {
    cwd: root,
    encoding: "utf8",
    maxBuffer: 20 * 1024 * 1024
  });
  const completedAt = new Date().toISOString();
  writeFileSync(join(rawDir, `${tool}.stdout.log`), cleanLog(result.stdout || ""));
  writeFileSync(join(rawDir, `${tool}.stderr.log`), cleanLog(result.stderr || ""));

  const exitCode = result.status ?? 1;
  const ok = allowedExitCodes.includes(exitCode);
  if (!ok) {
    throw new Error(`${tool} exited ${exitCode}; see ${relative(root, rawDir)}/${tool}.stderr.log`);
  }
  return {
    tool,
    image: images[tool],
    command: command.map(redactArg),
    exit_code: exitCode,
    expected_exit_codes: allowedExitCodes,
    started_at: startedAt,
    completed_at: completedAt,
    raw_artifacts: rawArtifacts(tool)
  };
}

function rawArtifacts(tool) {
  return [
    relative(root, join(rawDir, `${tool}.json`)),
    relative(root, join(rawDir, `${tool}.stdout.log`)),
    relative(root, join(rawDir, `${tool}.stderr.log`))
  ].filter((p) => existsSync(join(root, p)));
}

function normalizeGitleaks(path) {
  if (!existsSync(path)) return [];
  const data = parseJSON(path);
  const rows = Array.isArray(data) ? data : [];
  return rows.map((item) => ({
    tool: "gitleaks",
    rule_id: item.RuleID || "unknown",
    title: item.Description || "Potential secret detected",
    severity: "high",
    confidence: "high",
    file: normalizePath(item.File),
    line: item.StartLine || 0,
    description: "Gitleaks detected a secret-like value. The spike uses a custom fake-secret rule.",
    remediation: "Remove the secret value, rotate real credentials, and use an environment-backed secret store.",
    raw_artifact: relative(root, path)
  }));
}

function normalizeSemgrep(path) {
  if (!existsSync(path)) return [];
  const data = parseJSON(path);
  return (data.results || []).map((item) => ({
    tool: "semgrep",
    rule_id: item.check_id || "unknown",
    title: item.extra?.message || "Semgrep finding",
    severity: mapSemgrepSeverity(item.extra?.severity),
    confidence: "medium",
    file: normalizePath(item.path),
    line: item.start?.line || 0,
    description: item.extra?.message || "",
    remediation: "Review the matched code path and replace unsafe dynamic execution with explicit, validated behavior.",
    raw_artifact: relative(root, path)
  }));
}

function normalizeTrivy(path) {
  if (!existsSync(path)) return [];
  const data = parseJSON(path);
  const findings = [];
  for (const result of data.Results || []) {
    for (const item of result.Misconfigurations || []) {
      findings.push({
        tool: "trivy",
        rule_id: item.ID || "unknown",
        title: item.Title || "Trivy misconfiguration",
        severity: normalizeSeverity(item.Severity),
        confidence: "medium",
        file: normalizePath(result.Target),
        line: item.CauseMetadata?.StartLine || 0,
        description: item.Description || item.Message || "",
        remediation: item.Resolution || "Review the infrastructure configuration and apply the recommended control.",
        raw_artifact: relative(root, path)
      });
    }
  }
  return findings;
}

function parseJSON(path) {
  const text = readFileSync(path, "utf8").trim();
  if (!text) return {};
  return JSON.parse(text);
}

function summarize(findings) {
  const byTool = {};
  const bySeverity = {};
  for (const finding of findings) {
    byTool[finding.tool] = (byTool[finding.tool] || 0) + 1;
    bySeverity[finding.severity] = (bySeverity[finding.severity] || 0) + 1;
  }
  return { total_findings: findings.length, by_tool: byTool, by_severity: bySeverity };
}

function printSummary(report = null) {
  if (!report) {
    report = parseJSON(reportPath);
  }
  console.log("Scanrail scanner adapter spike");
  console.log(`fixture: ${report.fixture}`);
  console.log(`findings: ${report.summary.total_findings}`);
  console.log("");
  for (const [tool, count] of Object.entries(report.summary.by_tool || {})) {
    console.log(`${tool.padEnd(9)} ${String(count).padStart(2)} finding(s)`);
  }
  console.log("");
  console.log(`json: ${relative(root, reportPath)}`);
  console.log(`html: ${relative(root, htmlPath)}`);
  console.log(`summary: ${relative(root, summaryPath)}`);
}

function renderMarkdownSummary(report) {
  const lines = [
    "# Scanner Adapter Spike Summary",
    "",
    `Generated: ${report.generated_at}`,
    "",
    "## Commands",
    ""
  ];
  for (const command of report.commands) {
    lines.push(`- ${command.tool}: exit ${command.exit_code}, image \`${command.image}\``);
  }
  lines.push("", "## Findings", "");
  for (const finding of report.findings) {
    lines.push(`- [${finding.severity}] ${finding.tool}:${finding.rule_id} at ${finding.file}:${finding.line}`);
  }
  lines.push("");
  return lines.join("\n");
}

function renderHTML(report) {
  const rows = report.findings.map((f) => `<tr><td>${escapeHTML(f.tool)}</td><td>${escapeHTML(f.severity)}</td><td>${escapeHTML(f.rule_id)}</td><td>${escapeHTML(f.file)}:${f.line}</td><td>${escapeHTML(f.title)}</td></tr>`).join("\n");
  return `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Scanrail Scanner Adapter Spike</title>
  <style>
    body { font-family: system-ui, sans-serif; margin: 2rem; color: #17202a; }
    table { border-collapse: collapse; width: 100%; }
    th, td { border: 1px solid #d8dee4; padding: .55rem; text-align: left; }
    th { background: #f6f8fa; }
  </style>
</head>
<body>
  <h1>Scanrail Scanner Adapter Spike</h1>
  <p>Fixture: ${escapeHTML(report.fixture)}</p>
  <p>Total findings: ${report.summary.total_findings}</p>
  <table>
    <thead><tr><th>Tool</th><th>Severity</th><th>Rule</th><th>Location</th><th>Title</th></tr></thead>
    <tbody>${rows}</tbody>
  </table>
</body>
</html>
`;
}

function writeJSON(path, data) {
  writeFileSync(path, `${JSON.stringify(data, null, 2)}\n`);
}

function mapSemgrepSeverity(value) {
  switch ((value || "").toUpperCase()) {
    case "ERROR":
      return "high";
    case "WARNING":
      return "medium";
    case "INFO":
      return "low";
    default:
      return "unknown";
  }
}

function normalizeSeverity(value) {
  return String(value || "unknown").toLowerCase();
}

function normalizePath(value) {
  return String(value || "").replace(/^\/src\/?/, "").replace(/^\.\//, "");
}

function redact(value) {
  return String(value).replaceAll(testSecret, "[REDACTED_TEST_SECRET]");
}

function cleanLog(value) {
  return redact(value).replace(/[ \t]+$/gm, "");
}

function redactArg(value) {
  return redact(String(value)).replace(root, "$REPO");
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}
