package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/raeseoklee/scanrail/internal/config"
	"github.com/raeseoklee/scanrail/internal/exitcode"
)

func TestRunRedactsSecretsBeforeWritingReports(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	target := "http://127.0.0.1:1/?token=supersecret"
	cfg := config.Defaults(dir)
	cfg.ProjectName = "demo"
	cfg.TargetURL = target
	if err := config.WriteInitial("scanrail.yaml", cfg, false); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	code := Run(context.Background(), RunOptions{Only: "headers"}, &out)
	if code != exitcode.ConfigError {
		t.Fatalf("Run exit code = %d, want ConfigError from unreachable test target", code)
	}
	if strings.Contains(out.String(), "supersecret") {
		t.Fatalf("stdout leaked target secret: %s", out.String())
	}
}

func TestRunRedactsSuccessfulReports(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	target := startHeadersTestServer(t) + "/?token=supersecret"
	cfg := config.Defaults(dir)
	cfg.ProjectName = "demo"
	cfg.TargetURL = target
	if err := config.WriteInitial("scanrail.yaml", cfg, false); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	code := Run(context.Background(), RunOptions{Only: "headers"}, &out)
	if code != exitcode.OK {
		t.Fatalf("Run exit code = %d, output: %s", code, out.String())
	}
	reports, err := filepath.Glob(filepath.Join(dir, ".scanrail", "reports", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 {
		t.Fatalf("reports = %#v", reports)
	}
	data, err := os.ReadFile(reports[0])
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "supersecret") {
		t.Fatalf("report leaked target secret: %s", string(data))
	}
	if !strings.Contains(string(data), "[REDACTED]") {
		t.Fatalf("report missing redaction marker: %s", string(data))
	}
}

func TestExplicitUnreadyAdapterFailsSafety(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	cfg := config.Defaults(dir)
	cfg.ProjectName = "demo"
	if err := config.WriteInitial("scanrail.yaml", cfg, false); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	code := Run(context.Background(), RunOptions{Only: "gitleaks"}, &out)
	if code != exitcode.SafetyViolation {
		t.Fatalf("Run exit code = %d, want SafetyViolation; output: %s", code, out.String())
	}
	if !strings.Contains(out.String(), "not production-ready") {
		t.Fatalf("safety output missing reason: %s", out.String())
	}
}
