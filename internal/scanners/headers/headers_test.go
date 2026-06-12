package headers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScanFindsMissingHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	findings, err := Scan(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) == 0 {
		t.Fatal("expected missing header findings")
	}
}

func TestScanRequiresTarget(t *testing.T) {
	if _, err := Scan(context.Background(), ""); err == nil {
		t.Fatal("expected target error")
	}
}
