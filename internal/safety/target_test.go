package safety

import (
	"strings"
	"testing"
)

func TestValidateWebTargetAllowsHostAndPort(t *testing.T) {
	err := ValidateWebTarget("https://api.example.com:8443/v1/orders", WebTargetPolicy{
		Allowlist:        []string{"api.example.com:8443"},
		RequireAllowlist: true,
	})
	if err != nil {
		t.Fatalf("ValidateWebTarget returned error: %v", err)
	}
}

func TestValidateWebTargetMatchesDefaultHTTPSPort(t *testing.T) {
	err := ValidateWebTarget("https://api.example.com/v1/orders", WebTargetPolicy{
		Allowlist:        []string{"api.example.com:443"},
		RequireAllowlist: true,
	})
	if err != nil {
		t.Fatalf("ValidateWebTarget returned error: %v", err)
	}
}

func TestValidateWebTargetAllowsWildcardHost(t *testing.T) {
	err := ValidateWebTarget("https://staging.api.example.com", WebTargetPolicy{
		Allowlist:        []string{"*.api.example.com"},
		RequireAllowlist: true,
	})
	if err != nil {
		t.Fatalf("ValidateWebTarget returned error: %v", err)
	}
}

func TestValidateWebTargetRejectsMissingAllowlistWhenRequired(t *testing.T) {
	err := ValidateWebTarget("https://api.example.com", WebTargetPolicy{RequireAllowlist: true})
	if err == nil {
		t.Fatal("expected missing allowlist error")
	}
	if !strings.Contains(err.Error(), "allowlist is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateWebTargetRejectsHostOutsideAllowlist(t *testing.T) {
	err := ValidateWebTarget("https://prod.example.com", WebTargetPolicy{
		Allowlist:        []string{"staging.example.com"},
		RequireAllowlist: true,
	})
	if err == nil {
		t.Fatal("expected outside allowlist error")
	}
	if !strings.Contains(err.Error(), "outside the configured Scanrail allowlist") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateWebTargetRejectsBlockedPath(t *testing.T) {
	err := ValidateWebTarget("https://api.example.com/admin/destructive/delete", WebTargetPolicy{
		Allowlist:        []string{"api.example.com"},
		BlockedPaths:     []string{"/admin/destructive/*"},
		RequireAllowlist: true,
	})
	if err == nil {
		t.Fatal("expected blocked path error")
	}
	if !strings.Contains(err.Error(), "blocked by the configured Scanrail path safety policy") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateWebTargetDoesNotOvermatchSlashWildcardPath(t *testing.T) {
	err := ValidateWebTarget("https://api.example.com/admin/destructive-action", WebTargetPolicy{
		Allowlist:        []string{"api.example.com"},
		BlockedPaths:     []string{"/admin/destructive/*"},
		RequireAllowlist: true,
	})
	if err != nil {
		t.Fatalf("ValidateWebTarget returned error: %v", err)
	}
}
