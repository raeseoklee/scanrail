package gitleaks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/raeseoklee/scanrail/internal/dockerx"
	"github.com/raeseoklee/scanrail/internal/report"
	"github.com/raeseoklee/scanrail/internal/safety"
)

const (
	ToolName              = "gitleaks"
	Image                 = "ghcr.io/gitleaks/gitleaks:v8.30.1"
	ToolVersion           = "v8.30.1"
	containerWorkspaceDir = "/scan/workspace"
	outputDir             = "/scan/out"
	reportName            = "gitleaks.json"
)

type Options struct {
	WorkspaceDir string
	RawDir       string
	Runner       dockerx.Runner
	Redactor     safety.Redactor
}

type Result struct {
	Findings []report.Finding
	Metadata report.ToolMetadata
}

type Finding struct {
	Description string  `json:"Description"`
	File        string  `json:"File"`
	StartLine   int     `json:"StartLine"`
	EndLine     int     `json:"EndLine"`
	StartColumn int     `json:"StartColumn"`
	EndColumn   int     `json:"EndColumn"`
	RuleID      string  `json:"RuleID"`
	Fingerprint string  `json:"Fingerprint"`
	Secret      string  `json:"Secret,omitempty"`
	Match       string  `json:"Match,omitempty"`
	Entropy     float64 `json:"Entropy,omitempty"`
}

func Scan(ctx context.Context, opts Options) (Result, error) {
	if opts.Runner == nil {
		opts.Runner = dockerx.CLIRunner{}
	}
	if opts.WorkspaceDir == "" {
		return Result{}, errors.New("gitleaks scanner requires a workspace directory")
	}
	if opts.RawDir == "" {
		return Result{}, errors.New("gitleaks scanner requires a raw output directory")
	}
	workspaceDir, err := filepath.Abs(opts.WorkspaceDir)
	if err != nil {
		return Result{}, err
	}
	rawDir, err := filepath.Abs(opts.RawDir)
	if err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		return Result{}, err
	}

	rawPath := filepath.Join(rawDir, reportName)
	command := dockerx.Command{
		Image:   Image,
		Workdir: containerWorkspaceDir,
		Mounts: []dockerx.Mount{
			{Source: workspaceDir, Target: containerWorkspaceDir, ReadOnly: true},
			{Source: rawDir, Target: outputDir},
		},
		Args: []string{
			"dir",
			"--no-banner",
			"--redact=100",
			"--exit-code", "0",
			"--report-format", "json",
			"--report-path", filepath.ToSlash(filepath.Join(outputDir, reportName)),
			containerWorkspaceDir,
		},
	}
	runResult, err := opts.Runner.Run(ctx, command)
	if err != nil {
		message := strings.TrimSpace(runResult.Stderr)
		if message == "" {
			message = err.Error()
		}
		return Result{}, fmt.Errorf("gitleaks docker scan failed: %s: %w", opts.Redactor.RedactString(message), err)
	}

	findings, err := readFindings(rawPath)
	if err != nil {
		return Result{}, err
	}
	if err := sanitizeRawReport(rawPath, findings, opts.Redactor); err != nil {
		return Result{}, err
	}
	return Result{
		Findings: normalizeFindings(findings, opts.Redactor),
		Metadata: report.ToolMetadata{
			Tool:    ToolName,
			Version: ToolVersion,
			Image:   Image,
			RawPath: opts.Redactor.RedactString(rawPath),
		},
	}, nil
}

func readFindings(path string) ([]Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var findings []Finding
	if err := json.Unmarshal(data, &findings); err != nil {
		return nil, fmt.Errorf("parse gitleaks report: %w", err)
	}
	return findings, nil
}

func sanitizeRawReport(path string, findings []Finding, redactor safety.Redactor) error {
	if findings == nil {
		findings = []Finding{}
	}
	for i := range findings {
		findings[i].Secret = "[REDACTED]"
		findings[i].Match = "[REDACTED]"
		findings[i].Description = redactor.RedactString(findings[i].Description)
		findings[i].File = redactor.RedactString(findings[i].File)
		findings[i].RuleID = redactor.RedactString(findings[i].RuleID)
		findings[i].Fingerprint = redactor.RedactString(findings[i].Fingerprint)
	}
	data, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func normalizeFindings(findings []Finding, redactor safety.Redactor) []report.Finding {
	out := make([]report.Finding, 0, len(findings))
	for _, finding := range findings {
		ruleID := redactor.RedactString(finding.RuleID)
		if ruleID == "" {
			ruleID = "unknown"
		}
		location := location(finding)
		id := "gitleaks." + ruleID
		if finding.Fingerprint != "" {
			id = "gitleaks." + redactor.RedactString(finding.Fingerprint)
		}
		out = append(out, report.Finding{
			ID:          id,
			Tool:        ToolName,
			Title:       "Potential secret detected by Gitleaks",
			Severity:    "high",
			Confidence:  "high",
			Target:      redactor.RedactString(location),
			Description: "Gitleaks rule " + ruleID + " matched " + redactor.RedactString(location) + ". Secret value redacted.",
			Remediation: "Remove the secret from source, rotate the exposed credential, and move the value into an approved secret manager.",
			Evidence:    redactor.RedactString(finding.Description),
		})
	}
	return out
}

func location(finding Finding) string {
	file := finding.File
	if file == "" {
		file = "workspace"
	}
	if finding.StartLine > 0 {
		return file + ":" + strconv.Itoa(finding.StartLine)
	}
	return file
}
