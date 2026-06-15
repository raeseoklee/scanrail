package safety

import (
	"os"
	"regexp"
	"strings"
)

var (
	authorizationPattern = regexp.MustCompile(`(?i)\b(authorization\s*[:=]\s*)(bearer|basic)?\s*[A-Za-z0-9._~+/=-]+`)
	cookiePattern        = regexp.MustCompile(`(?i)\b((?:set-cookie|cookie)\s*[:=]\s*)[^\r\n]+`)
	secretKeyPattern     = regexp.MustCompile(`(?i)\b((?:access[_-]?token|refresh[_-]?token|api[_-]?key|client[_-]?secret|secret|password|token)\b\s*[:=]\s*)("[^"]*"|'[^']*'|[^\s,;}]+)`)
	querySecretPattern   = regexp.MustCompile(`(?i)([?&](?:access[_-]?token|refresh[_-]?token|api[_-]?key|key|secret|password|token|auth|authorization)=)[^&#\s]+`)
	urlUserInfoPattern   = regexp.MustCompile(`(https?://)[^/@\s]+@`)
	jwtPattern           = regexp.MustCompile(`\beyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\b`)
	commonTokenPattern   = regexp.MustCompile(`\b(?:sk|ghp|gho|ghu|ghs|glpat)-[A-Za-z0-9_=-]{12,}\b`)
)

type Redactor struct {
	values []string
}

func DefaultRedactor() Redactor {
	return Redactor{}
}

func NewRedactor(values ...string) Redactor {
	seen := map[string]bool{}
	var cleaned []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if len(value) < 4 || seen[value] {
			continue
		}
		seen[value] = true
		cleaned = append(cleaned, value)
	}
	return Redactor{values: cleaned}
}

func NewRedactorFromEnv(names ...string) Redactor {
	var values []string
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if value := os.Getenv(name); value != "" {
			values = append(values, value)
		}
	}
	return NewRedactor(values...)
}

func (r Redactor) RedactString(value string) string {
	if value == "" {
		return value
	}
	for _, secret := range r.values {
		value = strings.ReplaceAll(value, secret, "[REDACTED]")
	}
	value = authorizationPattern.ReplaceAllString(value, "${1}[REDACTED]")
	value = cookiePattern.ReplaceAllString(value, "${1}[REDACTED]")
	value = secretKeyPattern.ReplaceAllString(value, "${1}[REDACTED]")
	value = querySecretPattern.ReplaceAllString(value, "${1}[REDACTED]")
	value = urlUserInfoPattern.ReplaceAllString(value, "${1}[REDACTED]@")
	value = jwtPattern.ReplaceAllString(value, "[REDACTED]")
	value = commonTokenPattern.ReplaceAllString(value, "[REDACTED]")
	return value
}

func (r Redactor) RedactValue(value any) any {
	switch typed := value.(type) {
	case string:
		return r.RedactString(typed)
	case []string:
		out := make([]string, len(typed))
		for i, item := range typed {
			out[i] = r.RedactString(item)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for i, item := range typed {
			out[i] = r.RedactValue(item)
		}
		return out
	case map[string]string:
		out := make(map[string]string, len(typed))
		for key, item := range typed {
			out[key] = r.RedactString(item)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, item := range typed {
			out[key] = r.RedactValue(item)
		}
		return out
	default:
		return value
	}
}
