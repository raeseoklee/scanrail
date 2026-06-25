package scanners

import (
	"slices"
	"testing"
)

func TestDefaultToolsOnlyIncludesProductionReadyAdapters(t *testing.T) {
	tools := DefaultTools()
	if !slices.Equal(tools, []string{"gitleaks", "headers", "tls"}) {
		t.Fatalf("DefaultTools() = %#v", tools)
	}
	for _, tool := range tools {
		def, ok := DefinitionFor(tool)
		if !ok {
			t.Fatalf("missing definition for %s", tool)
		}
		if !def.ProductionReady {
			t.Fatalf("%s is in default tools but is not production-ready", tool)
		}
	}
}
