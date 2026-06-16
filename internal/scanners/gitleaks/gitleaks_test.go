package gitleaks

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/raeseoklee/scanrail/internal/dockerx"
	"github.com/raeseoklee/scanrail/internal/safety"
)

type fakeRunner struct {
	command dockerx.Command
	report  string
	err     error
}

func (f *fakeRunner) Run(_ context.Context, command dockerx.Command) (dockerx.Result, error) {
	f.command = command
	if f.err != nil {
		return dockerx.Result{ExitCode: 1, Stderr: "secret=supersecret"}, f.err
	}
	outDir := mountedSource(command, outputDir)
	if outDir != "" {
		if err := os.WriteFile(filepath.Join(outDir, reportName), []byte(f.report), 0o644); err != nil {
			return dockerx.Result{}, err
		}
	}
	return dockerx.Result{}, nil
}

func TestScanBuildsDockerCommandAndNormalizesFindings(t *testing.T) {
	dir := t.TempDir()
	rawDir := filepath.Join(dir, ".scanrail", "raw", "run")
	runner := &fakeRunner{report: `[
  {
    "Description": "AWS Access Key",
    "File": "app/config.js",
    "StartLine": 3,
    "RuleID": "aws-access-token",
    "Fingerprint": "app/config.js:aws-access-token:3",
    "Secret": "supersecret",
    "Match": "token=supersecret"
  }
]`}

	result, err := Scan(context.Background(), Options{
		WorkspaceDir: dir,
		RawDir:       rawDir,
		Runner:       runner,
		Redactor:     safety.NewRedactor("supersecret"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if runner.command.Image != Image {
		t.Fatalf("image = %q, want %q", runner.command.Image, Image)
	}
	if !slices.Contains(runner.command.Args, "--redact=100") {
		t.Fatalf("args missing --redact=100: %#v", runner.command.Args)
	}
	if !slices.Contains(runner.command.Args, "--report-format") || !slices.Contains(runner.command.Args, "json") {
		t.Fatalf("args missing json report flags: %#v", runner.command.Args)
	}
	if len(runner.command.Mounts) != 2 || !runner.command.Mounts[0].ReadOnly {
		t.Fatalf("mounts = %#v, want read-only workspace and writable output", runner.command.Mounts)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("findings = %d, want 1", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.Tool != ToolName {
		t.Fatalf("finding tool = %q", finding.Tool)
	}
	if strings.Contains(finding.Description, "supersecret") || strings.Contains(finding.Evidence, "supersecret") {
		t.Fatalf("normalized finding leaked secret: %#v", finding)
	}
	if result.Metadata.Image != Image || result.Metadata.Version != ToolVersion {
		t.Fatalf("metadata = %#v", result.Metadata)
	}
	raw, err := os.ReadFile(filepath.Join(rawDir, reportName))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "supersecret") {
		t.Fatalf("raw artifact leaked secret: %s", raw)
	}
}

func TestScanRedactsRunnerErrors(t *testing.T) {
	runner := &fakeRunner{err: os.ErrPermission}
	_, err := Scan(context.Background(), Options{
		WorkspaceDir: t.TempDir(),
		RawDir:       t.TempDir(),
		Runner:       runner,
		Redactor:     safety.NewRedactor("supersecret"),
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "supersecret") {
		t.Fatalf("error leaked secret: %v", err)
	}
}

func mountedSource(command dockerx.Command, target string) string {
	for _, mount := range command.Mounts {
		if mount.Target == target {
			return mount.Source
		}
	}
	return ""
}
