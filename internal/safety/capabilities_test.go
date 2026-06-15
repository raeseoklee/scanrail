package safety

import "testing"

func TestMissingCapabilities(t *testing.T) {
	missing := MissingCapabilities(InteractiveNetworkRequirements(), Capabilities{AllowlistScope: true})
	if len(missing) != 2 {
		t.Fatalf("missing = %#v, want redirect_scope and rate_limit", missing)
	}
	if JoinMissing(missing) != "redirect_scope, rate_limit" {
		t.Fatalf("JoinMissing = %q", JoinMissing(missing))
	}
}
