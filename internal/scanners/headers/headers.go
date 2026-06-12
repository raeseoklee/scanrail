package headers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/raeseoklee/scanrail/internal/report"
)

func Scan(ctx context.Context, target string) ([]report.Finding, error) {
	if strings.TrimSpace(target) == "" {
		return nil, errors.New("headers scanner requires a web target")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "scanrail/0")
	req.Header.Set("X-Scanrail-Scanner", "headers")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var findings []report.Finding
	addMissing := func(header, severity, remediation string) {
		if resp.Header.Get(header) == "" {
			id := "headers.missing." + strings.ToLower(header)
			id = strings.ReplaceAll(id, "-", "_")
			findings = append(findings, report.Finding{
				ID:          id,
				Title:       "Missing " + header + " header",
				Severity:    severity,
				Confidence:  "high",
				Target:      target,
				Description: "The response does not include the " + header + " security header.",
				Remediation: remediation,
			})
		}
	}
	addMissing("Content-Security-Policy", "medium", "Add a restrictive Content-Security-Policy header.")
	addMissing("X-Content-Type-Options", "low", "Set X-Content-Type-Options to nosniff.")
	addMissing("X-Frame-Options", "low", "Set X-Frame-Options or frame-ancestors in CSP.")
	addMissing("Referrer-Policy", "low", "Set a Referrer-Policy appropriate for the application.")
	if strings.HasPrefix(strings.ToLower(target), "https://") {
		addMissing("Strict-Transport-Security", "medium", "Set Strict-Transport-Security on HTTPS responses.")
	}
	return findings, nil
}
