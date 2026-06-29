package safety

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"
)

type WebTargetPolicy struct {
	Allowlist        []string
	BlockedPaths     []string
	RequireAllowlist bool
}

func ValidateWebTarget(raw string, policy WebTargetPolicy) error {
	target, err := parseWebTarget(raw)
	if err != nil {
		return err
	}
	if policy.RequireAllowlist && len(nonEmpty(policy.Allowlist)) == 0 {
		return errors.New("target allowlist is required for web target scans")
	}
	if len(nonEmpty(policy.Allowlist)) > 0 && !targetAllowed(target, policy.Allowlist) {
		return fmt.Errorf("target host %q is outside the configured Scanrail allowlist", target.host)
	}
	if blockedByPath(target.path, policy.BlockedPaths) {
		return fmt.Errorf("target path %q is blocked by the configured Scanrail path safety policy", target.path)
	}
	return nil
}

type webTarget struct {
	host string
	port string
	path string
}

func parseWebTarget(raw string) (webTarget, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return webTarget{}, errors.New("web target URL is required")
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return webTarget{}, fmt.Errorf("invalid target URL: %w", err)
	}
	if parsed.Hostname() == "" {
		return webTarget{}, errors.New("target URL must include a hostname")
	}
	port := parsed.Port()
	if port == "" {
		port = defaultPort(parsed.Scheme)
	}
	return webTarget{
		host: normalizeHostname(parsed.Hostname()),
		port: port,
		path: normalizePath(parsed.EscapedPath()),
	}, nil
}

func targetAllowed(target webTarget, allowlist []string) bool {
	for _, raw := range allowlist {
		entry, ok := parseAllowlistEntry(raw)
		if !ok {
			continue
		}
		if entry.port != "" && entry.port != target.port {
			continue
		}
		if entry.wildcard {
			if target.host != entry.host && strings.HasSuffix(target.host, "."+entry.host) {
				return true
			}
			continue
		}
		if target.host == entry.host {
			return true
		}
	}
	return false
}

type allowlistEntry struct {
	host     string
	port     string
	wildcard bool
}

func parseAllowlistEntry(raw string) (allowlistEntry, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return allowlistEntry{}, false
	}
	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Hostname() == "" {
			return allowlistEntry{}, false
		}
		return allowlistEntry{
			host: normalizeHostname(parsed.Hostname()),
			port: parsed.Port(),
		}, true
	}
	host := raw
	port := ""
	if splitHost, splitPort, err := net.SplitHostPort(raw); err == nil {
		host = splitHost
		port = splitPort
	} else if maybeHost, maybePort, ok := strings.Cut(raw, ":"); ok && numeric(maybePort) {
		host = maybeHost
		port = maybePort
	}
	host = normalizeHostname(host)
	wildcard := strings.HasPrefix(host, "*.")
	if wildcard {
		host = strings.TrimPrefix(host, "*.")
	}
	if host == "" {
		return allowlistEntry{}, false
	}
	return allowlistEntry{host: host, port: port, wildcard: wildcard}, true
}

func blockedByPath(targetPath string, blocked []string) bool {
	targetPath = normalizePath(targetPath)
	for _, raw := range blocked {
		pattern := normalizePathPattern(raw)
		if pattern == "" {
			continue
		}
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(targetPath, prefix) {
				return true
			}
			continue
		}
		if targetPath == pattern {
			return true
		}
	}
	return false
}

func normalizePath(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "/"
	}
	unescaped, err := url.PathUnescape(raw)
	if err == nil {
		raw = unescaped
	}
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	cleaned := path.Clean(raw)
	if cleaned == "." {
		return "/"
	}
	return cleaned
}

func normalizePathPattern(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	hasWildcard := strings.HasSuffix(raw, "*")
	if hasWildcard {
		prefix := strings.TrimSuffix(raw, "*")
		trailingSlash := strings.HasSuffix(prefix, "/")
		normalized := normalizePath(prefix)
		if trailingSlash && normalized != "/" {
			normalized += "/"
		}
		return normalized + "*"
	}
	return normalizePath(raw)
}

func normalizeHostname(raw string) string {
	return strings.TrimSuffix(strings.TrimSpace(strings.ToLower(raw)), ".")
}

func nonEmpty(values []string) []string {
	var out []string
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, value)
		}
	}
	return out
}

func numeric(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func defaultPort(scheme string) string {
	switch strings.ToLower(scheme) {
	case "http":
		return "80"
	case "https":
		return "443"
	default:
		return ""
	}
}
