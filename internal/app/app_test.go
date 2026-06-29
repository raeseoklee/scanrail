package app

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
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
	code := Run(context.Background(), RunOptions{Only: "semgrep"}, &out)
	if code != exitcode.SafetyViolation {
		t.Fatalf("Run exit code = %d, want SafetyViolation; output: %s", code, out.String())
	}
	if !strings.Contains(out.String(), "not production-ready") {
		t.Fatalf("safety output missing reason: %s", out.String())
	}
}

func TestRunOnlyTLSPersistsFindings(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.Defaults(dir)
	cfg.ProjectName = "demo"
	cfg.TargetURL = server.URL
	cfg.FailOn = "critical"
	if err := config.WriteInitial("scanrail.yaml", cfg, false); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	code := Run(context.Background(), RunOptions{Only: "tls"}, &out)
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
	if !strings.Contains(string(data), `"tool": "tls"`) {
		t.Fatalf("report missing tls finding: %s", string(data))
	}
}

func TestRunOnlyHeadersFailsSafetyOutsideAllowlist(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	cfg := `project:
  name: demo
targets:
  web:
    url: http://127.0.0.1:1
    allowlist:
      - staging.example.com
safety:
  require_allowlist: true
policy:
  fail_on:
    severity: critical
report:
  output_dir: .scanrail/reports
`
	if err := os.WriteFile("scanrail.yaml", []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	code := Run(context.Background(), RunOptions{Only: "headers"}, &out)
	if code != exitcode.SafetyViolation {
		t.Fatalf("Run exit code = %d, want SafetyViolation; output: %s", code, out.String())
	}
	if !strings.Contains(out.String(), "outside the configured Scanrail allowlist") {
		t.Fatalf("safety output missing allowlist reason: %s", out.String())
	}
}

func TestRunOnlyHeadersFailsSafetyForExcludedPath(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	cfg := `project:
  name: demo
targets:
  web:
    url: http://127.0.0.1:1/logout
    allowlist:
      - 127.0.0.1:1
    exclude_paths:
      - /logout
safety:
  require_allowlist: true
policy:
  fail_on:
    severity: critical
report:
  output_dir: .scanrail/reports
`
	if err := os.WriteFile("scanrail.yaml", []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	code := Run(context.Background(), RunOptions{Only: "headers"}, &out)
	if code != exitcode.SafetyViolation {
		t.Fatalf("Run exit code = %d, want SafetyViolation; output: %s", code, out.String())
	}
	if !strings.Contains(out.String(), "blocked by the configured Scanrail path safety policy") {
		t.Fatalf("safety output missing blocked path reason: %s", out.String())
	}
}

func TestRunOnlyOpenAPIPersistsFindings(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	spec := `{
  "openapi": "3.1.0",
  "servers": [{"url": "https://api.example.com"}],
  "paths": {
    "/orders": {
      "post": {
        "responses": {
          "200": {"description": "ok"}
        }
      }
    }
  }
}`
	if err := os.WriteFile("openapi.json", []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := config.Defaults(dir)
	cfg.ProjectName = "demo"
	cfg.OpenAPIPath = "./openapi.json"
	cfg.FailOn = "critical"
	if err := config.WriteInitial("scanrail.yaml", cfg, false); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	code := Run(context.Background(), RunOptions{Only: "openapi"}, &out)
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
	if !strings.Contains(string(data), `"tool": "openapi"`) {
		t.Fatalf("report missing openapi finding: %s", string(data))
	}
	if !strings.Contains(string(data), `"id": "openapi.operation.missing_security"`) {
		t.Fatalf("report missing openapi security finding: %s", string(data))
	}
}
