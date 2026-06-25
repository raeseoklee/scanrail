package openapi

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/raeseoklee/scanrail/internal/report"
)

func TestScanFindsMissingSecurityForDestructiveJSONOperation(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "openapi.json")
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
	if err := os.WriteFile(specPath, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}

	findings, err := Scan(context.Background(), Options{SpecPath: specPath})
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(findings, "openapi.operation.missing_security", "high") {
		t.Fatalf("missing expected high security finding: %#v", findings)
	}
	if !hasFinding(findings, "openapi.operation.missing_client_error_response", "low") {
		t.Fatalf("missing expected response finding: %#v", findings)
	}
}

func TestScanFindsInsecureServer(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "openapi.json")
	spec := `{
  "openapi": "3.1.0",
  "security": [{"bearerAuth": []}],
  "servers": [{"url": "http://api.example.com"}],
  "paths": {
    "/orders": {
      "get": {
        "responses": {
          "200": {"description": "ok"},
          "401": {"description": "unauthorized"}
        }
      }
    }
  }
}`
	if err := os.WriteFile(specPath, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}

	findings, err := Scan(context.Background(), Options{SpecPath: specPath})
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(findings, "openapi.server.insecure_http", "medium") {
		t.Fatalf("missing expected HTTP server finding: %#v", findings)
	}
	if hasFinding(findings, "openapi.operation.missing_security", "medium") {
		t.Fatalf("root security should suppress operation security finding: %#v", findings)
	}
}

func TestScanParsesBasicYAMLSpec(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "openapi.yaml")
	spec := `openapi: 3.1.0
servers:
  - url: https://api.example.com
security:
  - bearerAuth: []
paths:
  /orders:
    get:
      responses:
        '200':
          description: ok
        '404':
          description: missing
`
	if err := os.WriteFile(specPath, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}

	findings, err := Scan(context.Background(), Options{SpecPath: specPath})
	if err != nil {
		t.Fatal(err)
	}
	for _, finding := range findings {
		if strings.HasPrefix(finding.ID, "openapi.operation.") {
			t.Fatalf("did not expect operation finding for secured YAML spec: %#v", findings)
		}
	}
}

func TestScanDoesNotTreatLaterYAMLListsAsEmptyRootSecurity(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "openapi.yaml")
	spec := `openapi: 3.1.0
security: []
servers:
  - url: https://api.example.com
paths:
  /orders:
    post:
      responses:
        '200':
          description: ok
`
	if err := os.WriteFile(specPath, []byte(spec), 0o644); err != nil {
		t.Fatal(err)
	}

	findings, err := Scan(context.Background(), Options{SpecPath: specPath})
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(findings, "openapi.operation.missing_security", "high") {
		t.Fatalf("empty root security should not be satisfied by later YAML lists: %#v", findings)
	}
}

func TestScanRequiresLocalSpecPath(t *testing.T) {
	if _, err := Scan(context.Background(), Options{}); err == nil {
		t.Fatal("expected missing spec path error")
	}
	if _, err := Scan(context.Background(), Options{SpecPath: "https://example.com/openapi.yaml"}); err == nil {
		t.Fatal("expected URL rejection")
	}
}

func TestURLSpecPathDetectionAllowsWindowsAbsolutePaths(t *testing.T) {
	if isURLSpecPath(`C:\work\api\openapi.json`) {
		t.Fatal("Windows absolute paths must be treated as local files")
	}
	if !isURLSpecPath("https://example.com/openapi.yaml") {
		t.Fatal("HTTPS OpenAPI specs must still be rejected")
	}
	if !isURLSpecPath("file:///tmp/openapi.yaml") {
		t.Fatal("file URLs are not accepted by the local path scanner")
	}
}

func hasFinding(findings []report.Finding, id string, severity string) bool {
	for _, finding := range findings {
		if finding.ID == id && finding.Severity == severity {
			return true
		}
	}
	return false
}
