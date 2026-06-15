package audit

import (
	"os"
	"strings"
	"testing"

	"github.com/raeseoklee/scanrail/internal/safety"
)

func TestAppendRedactsEvent(t *testing.T) {
	path := t.TempDir() + "/logs/mcp-audit.jsonl"
	code := 0
	err := Append(path, Event{
		Source:   "mcp",
		Action:   "scanrail_run",
		Tool:     "headers",
		Decision: "completed",
		Reason:   "Authorization: Bearer sk-testsecret0000000000",
		Target:   "https://example.com?token=supersecret",
		ExitCode: &code,
	}, safety.NewRedactor("literal-secret"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)
	for _, leaked := range []string{"sk-testsecret", "token=supersecret"} {
		if strings.Contains(output, leaked) {
			t.Fatalf("audit log leaked %q: %s", leaked, output)
		}
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Fatalf("audit log missing redaction marker: %s", output)
	}
	if !strings.Contains(output, `"exit_code":0`) {
		t.Fatalf("audit log missing exit code: %s", output)
	}
}
