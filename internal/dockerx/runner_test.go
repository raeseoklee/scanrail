package dockerx

import (
	"context"
	"errors"
	"testing"
)

func TestMountSpecUsesDockerMountSyntax(t *testing.T) {
	got := mountSpec(Mount{Source: `/tmp/work:with-colon`, Target: "/scan/workspace", ReadOnly: true})
	want := `type=bind,source=/tmp/work:with-colon,target=/scan/workspace,readonly`
	if got != want {
		t.Fatalf("mountSpec() = %q, want %q", got, want)
	}
}

func TestIsDockerUnavailable(t *testing.T) {
	err := errors.Join(ErrDockerUnavailable, context.Canceled)
	if !IsDockerUnavailable(err) {
		t.Fatal("expected docker unavailable error to be detected")
	}
}
