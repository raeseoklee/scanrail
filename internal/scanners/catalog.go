package scanners

import "github.com/raeseoklee/scanrail/internal/safety"

type Definition struct {
	Name            string
	Intrusiveness   safety.Intrusiveness
	Capabilities    safety.Capabilities
	ProductionReady bool
	SkipReason      string
}

func DefaultTools() []string {
	return []string{"gitleaks", "trivy", "semgrep", "headers"}
}

func DefinitionFor(name string) (Definition, bool) {
	def, ok := definitions[name]
	return def, ok
}

var definitions = map[string]Definition{
	"headers": {
		Name:          "headers",
		Intrusiveness: safety.IntrusivenessInteractive,
		Capabilities: safety.Capabilities{
			AllowlistScope:  true,
			RedirectScope:   true,
			RateLimit:       true,
			HeaderInjection: true,
		},
		ProductionReady: true,
	},
	"gitleaks": {
		Name:            "gitleaks",
		Intrusiveness:   safety.IntrusivenessPassive,
		ProductionReady: false,
		SkipReason:      "docker adapter is not production-ready; command generation and central redaction are required before execution",
	},
	"trivy": {
		Name:            "trivy",
		Intrusiveness:   safety.IntrusivenessPassive,
		ProductionReady: false,
		SkipReason:      "docker adapter is not production-ready; command generation and central redaction are required before execution",
	},
	"semgrep": {
		Name:            "semgrep",
		Intrusiveness:   safety.IntrusivenessPassive,
		ProductionReady: false,
		SkipReason:      "docker adapter is not production-ready; command generation and central redaction are required before execution",
	},
}
