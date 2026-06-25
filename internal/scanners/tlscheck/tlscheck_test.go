package tlscheck

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/raeseoklee/scanrail/internal/report"
)

func TestScanRequiresHTTPSTarget(t *testing.T) {
	if _, err := Scan(context.Background(), "http://example.com"); err == nil || !strings.Contains(err.Error(), "HTTPS") {
		t.Fatalf("expected HTTPS target error, got %v", err)
	}
}

func TestScanFindsUntrustedCertificate(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	findings, err := Scan(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(findings, "tls.certificate.untrusted") {
		t.Fatalf("findings missing untrusted certificate: %#v", findings)
	}
}

func TestScanFindsCertificateExpiringSoon(t *testing.T) {
	now := time.Date(2026, 6, 25, 0, 0, 0, 0, time.UTC)
	server := newTLSServer(t, now.Add(-time.Hour), now.Add(10*24*time.Hour))
	defer server.Close()

	findings, err := ScanWithOptions(context.Background(), server.URL, Options{Now: func() time.Time {
		return now
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(findings, "tls.certificate.expiring_soon") {
		t.Fatalf("findings missing expiring soon certificate: %#v", findings)
	}
}

func newTLSServer(t *testing.T, notBefore time.Time, notAfter time.Time) *httptest.Server {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "127.0.0.1"},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	cert := tls.Certificate{Certificate: [][]byte{certDER}, PrivateKey: privateKey}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	server.StartTLS()
	return server
}

func hasFinding(findings []report.Finding, id string) bool {
	for _, finding := range findings {
		if finding.ID == id {
			return true
		}
	}
	return false
}
