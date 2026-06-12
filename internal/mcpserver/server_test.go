package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/raeseoklee/scanrail/internal/config"
	"github.com/raeseoklee/scanrail/internal/report"
)

func TestServeListsToolsAndReadsResources(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	cfg := config.Defaults(dir)
	cfg.ProjectName = "demo"
	cfg.TargetURL = "http://localhost:8080"
	if err := config.WriteInitial("scanrail.yaml", cfg, false); err != nil {
		t.Fatal(err)
	}
	writeReport(t, filepath.Join(dir, ".scanrail", "reports", "demo-20260101.json"))

	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":3,"method":"resources/list"}`,
		`{"jsonrpc":"2.0","id":4,"method":"resources/read","params":{"uri":"scanrail://config"}}`,
		`{"jsonrpc":"2.0","id":5,"method":"resources/read","params":{"uri":"scanrail://reports/latest/summary"}}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if code := Serve(context.Background(), strings.NewReader(input), &out, ioDiscard{}); code != 0 {
		t.Fatalf("Serve exit code = %d", code)
	}
	responses := decodeResponses(t, out.String())
	if len(responses) != 5 {
		t.Fatalf("responses = %d, want 5\n%s", len(responses), out.String())
	}
	if !strings.Contains(out.String(), `"scanrail_run"`) {
		t.Fatalf("tools/list missing scanrail_run: %s", out.String())
	}
	if !strings.Contains(out.String(), `"scanrail://safety-model"`) {
		t.Fatalf("resources/list missing safety model: %s", out.String())
	}
	configText := contentText(t, responses[3])
	if !strings.Contains(configText, `"token_env": "SCANRAIL_TOKEN"`) {
		t.Fatalf("config resource missing redacted token env: %s", configText)
	}
	reportText := contentText(t, responses[4])
	if !strings.Contains(reportText, `"findings_count": 1`) {
		t.Fatalf("latest report summary missing finding count: %s", reportText)
	}
}

func TestRunToolRequiresActiveScanConfirmation(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	cfg := config.Defaults(dir)
	cfg.ProjectName = "demo"
	cfg.TargetURL = "http://localhost:8080"
	if err := config.WriteInitial("scanrail.yaml", cfg, false); err != nil {
		t.Fatal(err)
	}

	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"scanrail_run","arguments":{"only":"headers","target":"http://localhost:8080"}}}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if code := Serve(context.Background(), strings.NewReader(input), &out, ioDiscard{}); code != 0 {
		t.Fatalf("Serve exit code = %d", code)
	}
	if !strings.Contains(out.String(), `"isError":true`) {
		t.Fatalf("scanrail_run without confirmation should return tool error: %s", out.String())
	}
	if !strings.Contains(out.String(), "confirm_active_scan=true") {
		t.Fatalf("scanrail_run error should explain confirmation: %s", out.String())
	}
}

func TestAllowedTargetUsesConfigAllowlist(t *testing.T) {
	cfg := config.Config{
		TargetURL: "https://staging.example.com",
		Allowlist: []string{"api.example.com:443"},
	}
	if err := allowedTarget(cfg, "https://api.example.com/v1"); err != nil {
		t.Fatalf("expected allowlisted target: %v", err)
	}
	if err := allowedTarget(cfg, "https://prod.example.com"); err == nil {
		t.Fatal("expected target outside allowlist to fail")
	}
}

func writeReport(t *testing.T, path string) {
	t.Helper()
	rr := report.RunReport{
		Project: "demo",
		Target:  "http://localhost:8080",
		Profile: "quick",
		Findings: []report.Finding{{
			ID:       "headers.missing.content_security_policy",
			Title:    "Missing Content-Security-Policy header",
			Severity: "medium",
		}},
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(rr)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func decodeResponses(t *testing.T, output string) []map[string]any {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var responses []map[string]any
	for _, line := range lines {
		var response map[string]any
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			t.Fatalf("invalid JSON response %q: %v", line, err)
		}
		responses = append(responses, response)
	}
	return responses
}

func contentText(t *testing.T, response map[string]any) string {
	t.Helper()
	result, ok := response["result"].(map[string]any)
	if !ok {
		t.Fatalf("response result missing: %#v", response)
	}
	contents, ok := result["contents"].([]any)
	if !ok || len(contents) == 0 {
		t.Fatalf("response contents missing: %#v", response)
	}
	content, ok := contents[0].(map[string]any)
	if !ok {
		t.Fatalf("response content invalid: %#v", contents[0])
	}
	text, ok := content["text"].(string)
	if !ok {
		t.Fatalf("response text missing: %#v", content)
	}
	return text
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
