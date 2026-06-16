package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteJSONKeepsEmptyFindingsArray(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.json")
	err := WriteJSON(path, RunReport{
		Project:   "demo",
		Profile:   "quick",
		StartedAt: time.Unix(0, 0).UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	findings, ok := decoded["findings"].([]any)
	if !ok {
		t.Fatalf("findings = %#v, want JSON array", decoded["findings"])
	}
	if len(findings) != 0 {
		t.Fatalf("findings length = %d, want 0", len(findings))
	}
}
