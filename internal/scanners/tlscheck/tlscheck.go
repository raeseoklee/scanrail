package tlscheck

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/raeseoklee/scanrail/internal/report"
)

const (
	ToolName    = "tls"
	ToolVersion = "native"
)

const expiryWarningWindow = 30 * 24 * time.Hour

type DialContextFunc func(ctx context.Context, network string, address string, config *tls.Config) (*tls.Conn, error)

type Options struct {
	Now         func() time.Time
	DialContext DialContextFunc
}

func Scan(ctx context.Context, target string) ([]report.Finding, error) {
	return ScanWithOptions(ctx, target, Options{})
}

func ScanWithOptions(ctx context.Context, target string, opts Options) ([]report.Finding, error) {
	parsed, err := parseTarget(target)
	if err != nil {
		return nil, err
	}
	now := time.Now
	if opts.Now != nil {
		now = opts.Now
	}
	dialContext := defaultDialContext
	if opts.DialContext != nil {
		dialContext = opts.DialContext
	}

	conn, err := dialContext(ctx, "tcp", parsed.address, &tls.Config{
		ServerName:         parsed.host,
		InsecureSkipVerify: true, // Collect certificate details, then verify explicitly below.
		MinVersion:         tls.VersionTLS10,
	})
	if err != nil {
		return nil, fmt.Errorf("tls scanner failed to connect to %s: %w", parsed.address, err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	return findingsForState(parsed.displayTarget, parsed.host, state, now()), nil
}

type scanTarget struct {
	host          string
	address       string
	displayTarget string
}

func parseTarget(raw string) (scanTarget, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return scanTarget{}, errors.New("tls scanner requires an HTTPS target")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return scanTarget{}, err
	}
	if !strings.EqualFold(u.Scheme, "https") || u.Hostname() == "" {
		return scanTarget{}, errors.New("tls scanner requires an HTTPS target")
	}
	port := u.Port()
	if port == "" {
		port = "443"
	}
	return scanTarget{
		host:          u.Hostname(),
		address:       net.JoinHostPort(u.Hostname(), port),
		displayTarget: raw,
	}, nil
}

func defaultDialContext(ctx context.Context, network string, address string, config *tls.Config) (*tls.Conn, error) {
	dialer := tls.Dialer{
		NetDialer: &net.Dialer{Timeout: 10 * time.Second},
		Config:    config,
	}
	conn, err := dialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		_ = conn.Close()
		return nil, errors.New("tls dialer returned a non-TLS connection")
	}
	return tlsConn, nil
}

func findingsForState(target string, host string, state tls.ConnectionState, now time.Time) []report.Finding {
	var findings []report.Finding
	add := func(id string, title string, severity string, description string, remediation string, evidence string) {
		findings = append(findings, report.Finding{
			ID:          id,
			Tool:        ToolName,
			Title:       title,
			Severity:    severity,
			Confidence:  "high",
			Target:      target,
			Description: description,
			Remediation: remediation,
			Evidence:    evidence,
		})
	}

	if state.Version != 0 && state.Version < tls.VersionTLS12 {
		add(
			"tls.protocol.legacy",
			"TLS protocol version is below 1.2",
			"high",
			"The target negotiated a TLS protocol version that is no longer appropriate for modern web services.",
			"Disable TLS 1.0 and TLS 1.1, then require TLS 1.2 or newer.",
			"negotiated_version="+tlsVersionName(state.Version),
		)
	}
	if len(state.PeerCertificates) == 0 {
		add(
			"tls.certificate.missing",
			"No peer certificate was presented",
			"critical",
			"The TLS handshake did not provide a peer certificate for validation.",
			"Configure the target to present a valid server certificate.",
			"",
		)
		return findings
	}

	leaf := state.PeerCertificates[0]
	if now.Before(leaf.NotBefore) {
		add(
			"tls.certificate.not_yet_valid",
			"TLS certificate is not valid yet",
			"high",
			"The target certificate validity window starts in the future.",
			"Install a certificate whose validity window includes the current deployment time.",
			"not_before="+leaf.NotBefore.UTC().Format(time.RFC3339),
		)
	}
	if now.After(leaf.NotAfter) {
		add(
			"tls.certificate.expired",
			"TLS certificate is expired",
			"high",
			"The target certificate is outside its validity window.",
			"Renew and deploy a valid certificate before exposing the service.",
			"not_after="+leaf.NotAfter.UTC().Format(time.RFC3339),
		)
	} else if leaf.NotAfter.Sub(now) <= expiryWarningWindow {
		add(
			"tls.certificate.expiring_soon",
			"TLS certificate expires soon",
			"medium",
			"The target certificate expires within 30 days.",
			"Renew the certificate or verify automated renewal before expiry.",
			"not_after="+leaf.NotAfter.UTC().Format(time.RFC3339),
		)
	}
	if err := leaf.VerifyHostname(host); err != nil {
		add(
			"tls.certificate.hostname_mismatch",
			"TLS certificate does not match the target host",
			"high",
			"The target certificate is not valid for the scanned hostname.",
			"Deploy a certificate with a subject alternative name that covers the target host.",
			err.Error(),
		)
	}
	if err := verifyChain(host, state.PeerCertificates, now); err != nil {
		add(
			"tls.certificate.untrusted",
			"TLS certificate chain is not trusted",
			"high",
			"The target certificate chain could not be verified against the system trust store.",
			"Install a certificate chain issued by a trusted CA, including required intermediates.",
			err.Error(),
		)
	}
	return findings
}

func verifyChain(host string, certs []*x509.Certificate, now time.Time) error {
	if len(certs) == 0 {
		return errors.New("no peer certificates")
	}
	intermediates := x509.NewCertPool()
	for _, cert := range certs[1:] {
		intermediates.AddCert(cert)
	}
	roots, err := x509.SystemCertPool()
	if err != nil {
		return err
	}
	_, err = certs[0].Verify(x509.VerifyOptions{
		DNSName:       host,
		Roots:         roots,
		Intermediates: intermediates,
		CurrentTime:   now,
	})
	return err
}

func tlsVersionName(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("0x%04x", version)
	}
}
