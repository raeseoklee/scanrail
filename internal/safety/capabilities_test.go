package safety

import "testing"

func TestMissingCapabilities(t *testing.T) {
	missing := MissingCapabilities(InteractiveNetworkRequirements(), Capabilities{AllowlistScope: true})
	if len(missing) != 3 {
		t.Fatalf("missing = %#v, want redirect_scope, blocked_paths, and rate_limit", missing)
	}
	if JoinMissing(missing) != "redirect_scope, blocked_paths, rate_limit" {
		t.Fatalf("JoinMissing = %q", JoinMissing(missing))
	}
}
