package safety

import "strings"

type Intrusiveness string

const (
	IntrusivenessPassive     Intrusiveness = "passive"
	IntrusivenessInteractive Intrusiveness = "interactive"
	IntrusivenessActive      Intrusiveness = "active"
)

type Capabilities struct {
	AllowlistScope  bool `json:"allowlist_scope"`
	RedirectScope   bool `json:"redirect_scope"`
	BlockedPaths    bool `json:"blocked_paths"`
	BlockedMethods  bool `json:"blocked_methods"`
	RateLimit       bool `json:"rate_limit"`
	HeaderInjection bool `json:"header_injection"`
	AuthInjection   bool `json:"auth_injection"`
}

func InteractiveNetworkRequirements() Capabilities {
	return Capabilities{
		AllowlistScope: true,
		RedirectScope:  true,
		RateLimit:      true,
	}
}

func MissingCapabilities(required Capabilities, actual Capabilities) []string {
	var missing []string
	if required.AllowlistScope && !actual.AllowlistScope {
		missing = append(missing, "allowlist_scope")
	}
	if required.RedirectScope && !actual.RedirectScope {
		missing = append(missing, "redirect_scope")
	}
	if required.BlockedPaths && !actual.BlockedPaths {
		missing = append(missing, "blocked_paths")
	}
	if required.BlockedMethods && !actual.BlockedMethods {
		missing = append(missing, "blocked_methods")
	}
	if required.RateLimit && !actual.RateLimit {
		missing = append(missing, "rate_limit")
	}
	if required.HeaderInjection && !actual.HeaderInjection {
		missing = append(missing, "header_injection")
	}
	if required.AuthInjection && !actual.AuthInjection {
		missing = append(missing, "auth_injection")
	}
	return missing
}

func JoinMissing(missing []string) string {
	return strings.Join(missing, ", ")
}
