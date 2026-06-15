package safety

import (
	"strings"
	"testing"
)

func TestRedactStringMasksCommonSecretShapes(t *testing.T) {
	redactor := NewRedactor("scanrail-demo-secret-0001")
	input := strings.Join([]string{
		"Authorization: Bearer sk-testsecret0000000000",
		"Cookie: session=abcdef; other=value",
		"https://user:pass@example.com/path?token=abc123&ok=true",
		"password: hunter2",
		"literal scanrail-demo-secret-0001",
	}, "\n")

	output := redactor.RedactString(input)
	for _, leaked := range []string{"sk-testsecret", "session=abcdef", "user:pass", "token=abc123", "hunter2", "scanrail-demo-secret-0001"} {
		if strings.Contains(output, leaked) {
			t.Fatalf("redacted output leaked %q: %s", leaked, output)
		}
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Fatalf("redacted output missing marker: %s", output)
	}
}

func TestRedactValueRecursesStructuredData(t *testing.T) {
	value := map[string]any{
		"target": "https://example.com?api_key=secret-value",
		"nested": []any{
			map[string]any{"token": "token=secret-value"},
		},
	}
	output := NewRedactor().RedactValue(value).(map[string]any)
	if strings.Contains(output["target"].(string), "secret-value") {
		t.Fatalf("target was not redacted: %#v", output)
	}
	nested := output["nested"].([]any)[0].(map[string]any)
	if strings.Contains(nested["token"].(string), "secret-value") {
		t.Fatalf("nested value was not redacted: %#v", output)
	}
}
