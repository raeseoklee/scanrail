package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/raeseoklee/scanrail/internal/config"
	"github.com/raeseoklee/scanrail/internal/exitcode"
	"github.com/raeseoklee/scanrail/internal/report"
	"github.com/raeseoklee/scanrail/internal/scanners/headers"
	"github.com/raeseoklee/scanrail/internal/version"
	"github.com/raeseoklee/scanrail/internal/workspace"
)

type InitOptions struct {
	ConfigPath  string
	ProjectName string
	Target      string
	Profile     string
	Force       bool
}

type SetupOptions struct {
	PullPolicy string
}

type RunOptions struct {
	ConfigPath string
	Profile    string
	Target     string
	Only       string
	OutputDir  string
}

func Doctor(stdout io.Writer) int {
	fmt.Fprintln(stdout, "Scanrail Doctor")
	fmt.Fprintln(stdout)
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Fprintln(stdout, "Docker              WARN docker command not found")
	} else if err := exec.Command("docker", "version", "--format", "{{.Server.Version}}").Run(); err != nil {
		fmt.Fprintln(stdout, "Docker              WARN docker command found but daemon unavailable")
	} else {
		fmt.Fprintln(stdout, "Docker              OK")
	}
	if _, err := os.Getwd(); err != nil {
		fmt.Fprintln(stdout, "Workspace           FAIL", err)
		return exitcode.Environment
	}
	fmt.Fprintln(stdout, "Workspace           OK")
	if _, err := os.Stat(config.DefaultPath); errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(stdout, "Config              WARN scanrail.yaml not found")
	} else if err != nil {
		fmt.Fprintln(stdout, "Config              FAIL", err)
		return exitcode.ConfigError
	} else {
		fmt.Fprintln(stdout, "Config              OK")
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Ready.")
	return exitcode.OK
}

func Init(opts InitOptions, stdout io.Writer) int {
	wd, _ := os.Getwd()
	cfg := config.Defaults(wd)
	if opts.ProjectName != "" {
		cfg.ProjectName = opts.ProjectName
	}
	if opts.Target != "" {
		cfg.TargetURL = opts.Target
	}
	if opts.Profile == "" {
		opts.Profile = "quick"
	}
	if err := config.Validate(cfg); err != nil {
		fmt.Fprintln(stdout, "Config error:", err)
		return exitcode.ConfigError
	}
	if err := config.WriteInitial(opts.ConfigPath, cfg, opts.Force); err != nil {
		fmt.Fprintln(stdout, "Init failed:", err)
		return exitcode.ConfigError
	}
	fmt.Fprintln(stdout, "Generated scanrail.yaml")
	fmt.Fprintln(stdout, "Generated .env.scanrail.example")
	return exitcode.OK
}

func Setup(opts SetupOptions, stdout io.Writer) int {
	if opts.PullPolicy == "" {
		opts.PullPolicy = "missing"
	}
	ws, err := workspace.New("", "")
	if err != nil {
		fmt.Fprintln(stdout, "Workspace error:", err)
		return exitcode.Environment
	}
	if err := ws.Ensure(); err != nil {
		fmt.Fprintln(stdout, "Setup failed:", err)
		return exitcode.Environment
	}
	fmt.Fprintln(stdout, "Preparing Scanrail runtime")
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Reports directory    created", rel(ws.ReportsDir, ws.Root))
	fmt.Fprintln(stdout, "Cache directory      created", rel(ws.CacheDir, ws.Root))

	images := map[string]string{
		"gitleaks": "zricethezav/gitleaks:<approved-version>",
		"trivy":    "aquasec/trivy:<approved-version>",
		"semgrep":  "semgrep/semgrep:<approved-version>",
	}
	if opts.PullPolicy != "never" {
		if _, err := exec.LookPath("docker"); err != nil {
			fmt.Fprintln(stdout, "Docker unavailable:", err)
			return exitcode.Environment
		}
		for name, image := range images {
			if strings.Contains(image, "<approved-version>") {
				fmt.Fprintf(stdout, "- %-10s SKIP image version not pinned yet\n", name)
				continue
			}
			if err := exec.Command("docker", "pull", image).Run(); err != nil {
				fmt.Fprintln(stdout, "Pull failed:", image, err)
				return exitcode.Environment
			}
		}
	}
	lockPath := filepath.Join(ws.Root, "tools.lock.yaml")
	if err := writeToolsLock(lockPath, images); err != nil {
		fmt.Fprintln(stdout, "Failed to write tools.lock.yaml:", err)
		return exitcode.Environment
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Generated tools.lock.yaml")
	return exitcode.OK
}

func Run(ctx context.Context, opts RunOptions, stdout io.Writer) int {
	wd, _ := os.Getwd()
	cfg, err := config.Load(opts.ConfigPath, wd)
	if err != nil {
		fmt.Fprintln(stdout, "Config error:", err)
		return exitcode.ConfigError
	}
	if opts.Target != "" {
		cfg.TargetURL = opts.Target
	}
	if opts.OutputDir != "" {
		cfg.OutputDir = opts.OutputDir
	}
	if opts.Profile == "" {
		opts.Profile = "quick"
	}
	if err := config.Validate(cfg); err != nil {
		fmt.Fprintln(stdout, "Config error:", err)
		return exitcode.ConfigError
	}
	ws, err := workspace.New(wd, cfg.OutputDir)
	if err != nil {
		fmt.Fprintln(stdout, "Workspace error:", err)
		return exitcode.Environment
	}
	if err := ws.Ensure(); err != nil {
		fmt.Fprintln(stdout, "Workspace error:", err)
		return exitcode.Environment
	}

	runReport := report.RunReport{
		Project:   cfg.ProjectName,
		Target:    cfg.TargetURL,
		Profile:   opts.Profile,
		StartedAt: time.Now().UTC(),
	}
	tools := []string{"gitleaks", "trivy", "semgrep", "headers"}
	if opts.Only != "" {
		tools = []string{opts.Only}
	}
	for _, tool := range tools {
		switch tool {
		case "headers":
			findings, err := headers.Scan(ctx, cfg.TargetURL)
			if err != nil {
				if opts.Only == "headers" {
					fmt.Fprintln(stdout, "headers failed:", err)
					return exitcode.ConfigError
				}
				runReport.Skipped = append(runReport.Skipped, report.Skipped{Tool: "headers", Reason: err.Error()})
				continue
			}
			runReport.Findings = append(runReport.Findings, findings...)
		case "gitleaks", "trivy", "semgrep":
			runReport.Skipped = append(runReport.Skipped, report.Skipped{Tool: tool, Reason: "docker adapter command generation is scaffolded for the first release candidate"})
		default:
			fmt.Fprintln(stdout, "Unknown tool:", tool)
			return exitcode.ConfigError
		}
	}
	base := filepath.Join(ws.ReportsDir, cfg.ProjectName+"-"+ws.RunID)
	if err := report.WriteJSON(base+".json", runReport); err != nil {
		fmt.Fprintln(stdout, "Report failed:", err)
		return exitcode.Environment
	}
	if err := report.WriteHTML(base+".html", runReport); err != nil {
		fmt.Fprintln(stdout, "Report failed:", err)
		return exitcode.Environment
	}
	fmt.Fprintln(stdout, "Scanrail completed")
	fmt.Fprintln(stdout, "JSON ", base+".json")
	fmt.Fprintln(stdout, "HTML ", base+".html")
	if hasPolicyFailure(runReport.Findings, cfg.FailOn) {
		return exitcode.PolicyFailed
	}
	return exitcode.OK
}

func Version(stdout io.Writer) int {
	fmt.Fprintln(stdout, "scanrail", version.String())
	return exitcode.OK
}

func writeToolsLock(path string, images map[string]string) error {
	var b strings.Builder
	b.WriteString("tools:\n")
	for name, image := range images {
		b.WriteString("  " + name + ":\n")
		b.WriteString("    image: " + image + "\n")
	}
	b.WriteString("generated_at: " + time.Now().UTC().Format(time.RFC3339) + "\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func hasPolicyFailure(findings []report.Finding, failOn string) bool {
	threshold := severityRank(failOn)
	if threshold == 0 {
		return false
	}
	for _, finding := range findings {
		if severityRank(finding.Severity) >= threshold {
			return true
		}
	}
	return false
}

func severityRank(sev string) int {
	switch strings.ToLower(sev) {
	case "critical":
		return 5
	case "high":
		return 4
	case "medium":
		return 3
	case "low":
		return 2
	case "info":
		return 1
	default:
		return 0
	}
}

func rel(path, root string) string {
	if r, err := filepath.Rel(root, path); err == nil {
		return r
	}
	return path
}
