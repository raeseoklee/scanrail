package openapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/raeseoklee/scanrail/internal/report"
)

const (
	ToolName    = "openapi"
	ToolVersion = "native"
)

type Options struct {
	SpecPath string
	WorkDir  string
	ReadFile func(string) ([]byte, error)
}

func Scan(ctx context.Context, opts Options) ([]report.Finding, error) {
	path := strings.TrimSpace(opts.SpecPath)
	if path == "" {
		return nil, errors.New("openapi scanner requires targets.api.openapi")
	}
	if isURLSpecPath(path) {
		return nil, errors.New("openapi scanner supports local OpenAPI files only in this release")
	}
	if !filepath.IsAbs(path) {
		workdir := opts.WorkDir
		if workdir == "" {
			workdir = "."
		}
		path = filepath.Join(workdir, path)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	readFile := os.ReadFile
	if opts.ReadFile != nil {
		readFile = opts.ReadFile
	}
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	spec, err := parseSpec(data)
	if err != nil {
		return nil, err
	}
	return findingsForSpec(path, spec), nil
}

type specDocument struct {
	Version      string
	Servers      []string
	RootSecurity bool
	Operations   []operation
}

type operation struct {
	Method    string
	Path      string
	Security  bool
	Responses []string
}

func parseSpec(data []byte) (specDocument, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return specDocument{}, errors.New("OpenAPI spec is empty")
	}
	if strings.HasPrefix(trimmed, "{") {
		return parseJSONSpec(data)
	}
	return parseYAMLishSpec(trimmed), nil
}

func parseJSONSpec(data []byte) (specDocument, error) {
	var raw struct {
		OpenAPI  string                     `json:"openapi"`
		Swagger  string                     `json:"swagger"`
		Servers  []struct{ URL string }     `json:"servers"`
		Security json.RawMessage            `json:"security"`
		Paths    map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return specDocument{}, fmt.Errorf("parse OpenAPI JSON: %w", err)
	}
	spec := specDocument{
		Version:      firstNonEmpty(raw.OpenAPI, raw.Swagger),
		RootSecurity: hasEffectiveSecurity(raw.Security),
	}
	for _, server := range raw.Servers {
		if url := strings.TrimSpace(server.URL); url != "" {
			spec.Servers = append(spec.Servers, url)
		}
	}
	for path, itemRaw := range raw.Paths {
		var pathItem map[string]json.RawMessage
		if err := json.Unmarshal(itemRaw, &pathItem); err != nil {
			continue
		}
		for method, operationRaw := range pathItem {
			method = strings.ToLower(method)
			if !isHTTPMethod(method) {
				continue
			}
			var rawOperation map[string]json.RawMessage
			if err := json.Unmarshal(operationRaw, &rawOperation); err != nil {
				continue
			}
			op := operation{
				Method:    strings.ToUpper(method),
				Path:      path,
				Security:  hasEffectiveSecurity(rawOperation["security"]),
				Responses: responseKeys(rawOperation["responses"]),
			}
			spec.Operations = append(spec.Operations, op)
		}
	}
	sortSpec(&spec)
	return spec, nil
}

func hasEffectiveSecurity(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var requirements []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &requirements); err != nil {
		return false
	}
	return len(requirements) > 0
}

func responseKeys(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	var responses map[string]json.RawMessage
	if err := json.Unmarshal(raw, &responses); err != nil {
		return nil
	}
	keys := make([]string, 0, len(responses))
	for key := range responses {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func parseYAMLishSpec(data string) specDocument {
	var spec specDocument
	top := ""
	pathIndent := -1
	methodIndent := -1
	currentPath := ""
	currentOp := -1
	rootSecurityIndent := -1
	operationSecurityIndent := -1
	responsesIndent := -1

	for _, raw := range strings.Split(data, "\n") {
		raw = trimYAMLComment(strings.ReplaceAll(raw, "\t", "  "))
		if strings.TrimSpace(raw) == "" {
			continue
		}
		indent := leadingSpaces(raw)
		line := strings.TrimSpace(raw)

		if rootSecurityIndent >= 0 && indent > rootSecurityIndent && strings.HasPrefix(line, "- ") {
			spec.RootSecurity = true
			continue
		}
		if operationSecurityIndent >= 0 && indent > operationSecurityIndent && currentOp >= 0 && strings.HasPrefix(line, "- ") {
			spec.Operations[currentOp].Security = true
			continue
		}
		if responsesIndent >= 0 && indent > responsesIndent && currentOp >= 0 {
			if key, _, ok := cutYAMLKeyValue(line); ok {
				code := cleanYAMLKey(key)
				if isResponseKey(code) {
					spec.Operations[currentOp].Responses = append(spec.Operations[currentOp].Responses, code)
				}
				continue
			}
		}
		if indent == 0 {
			pathIndent = -1
			methodIndent = -1
			currentPath = ""
			currentOp = -1
			operationSecurityIndent = -1
			responsesIndent = -1
		}

		key, value, ok := cutYAMLKeyValue(line)
		if !ok {
			continue
		}
		key = cleanYAMLKey(key)
		value = cleanYAMLValue(value)

		if indent == 0 {
			switch key {
			case "openapi", "swagger":
				spec.Version = value
				top = ""
				rootSecurityIndent = -1
			case "servers", "paths":
				top = key
				rootSecurityIndent = -1
			case "security":
				top = "security"
				rootSecurityIndent = indent
				spec.RootSecurity = value != "" && value != "[]"
			default:
				top = key
				rootSecurityIndent = -1
			}
			continue
		}

		switch top {
		case "servers":
			if strings.HasPrefix(line, "- ") {
				if key, value, ok := cutYAMLKeyValue(strings.TrimSpace(strings.TrimPrefix(line, "- "))); ok && cleanYAMLKey(key) == "url" {
					spec.Servers = append(spec.Servers, cleanYAMLValue(value))
				}
				continue
			}
			if key == "url" && value != "" {
				spec.Servers = append(spec.Servers, value)
			}
		case "paths":
			if strings.HasPrefix(key, "/") {
				currentPath = key
				pathIndent = indent
				methodIndent = -1
				currentOp = -1
				operationSecurityIndent = -1
				responsesIndent = -1
				continue
			}
			if currentPath != "" && indent > pathIndent && isHTTPMethod(strings.ToLower(key)) {
				spec.Operations = append(spec.Operations, operation{
					Method: strings.ToUpper(key),
					Path:   currentPath,
				})
				currentOp = len(spec.Operations) - 1
				methodIndent = indent
				operationSecurityIndent = -1
				responsesIndent = -1
				continue
			}
			if currentOp >= 0 && indent > methodIndent {
				switch key {
				case "security":
					operationSecurityIndent = indent
					spec.Operations[currentOp].Security = value != "" && value != "[]"
				case "responses":
					responsesIndent = indent
				}
			}
		}
	}
	sortSpec(&spec)
	return spec
}

func findingsForSpec(target string, spec specDocument) []report.Finding {
	var findings []report.Finding
	add := func(id string, title string, severity string, findingTarget string, description string, remediation string, evidence string) {
		findings = append(findings, report.Finding{
			ID:          id,
			Tool:        ToolName,
			Title:       title,
			Severity:    severity,
			Confidence:  "medium",
			Target:      findingTarget,
			Description: description,
			Remediation: remediation,
			Evidence:    evidence,
		})
	}

	if strings.TrimSpace(spec.Version) == "" {
		add(
			"openapi.version.missing",
			"OpenAPI version is missing",
			"medium",
			target,
			"The API contract does not declare an OpenAPI or Swagger version.",
			"Add an openapi or swagger version field to the API contract.",
			"",
		)
	}
	if len(spec.Servers) == 0 {
		add(
			"openapi.servers.missing",
			"OpenAPI servers are missing",
			"low",
			target,
			"The API contract does not declare a server URL, which makes environment and transport review harder.",
			"Declare at least one servers entry for the documented deployment environment.",
			"",
		)
	}
	for _, server := range spec.Servers {
		if isPlainHTTPServer(server) {
			add(
				"openapi.server.insecure_http",
				"OpenAPI server uses plain HTTP",
				"medium",
				server,
				"The API contract advertises a plain HTTP server URL.",
				"Use HTTPS for non-local API environments or document why plaintext is limited to local development.",
				"server="+server,
			)
		}
	}
	if len(spec.Operations) == 0 {
		add(
			"openapi.paths.missing_operations",
			"OpenAPI paths do not define operations",
			"medium",
			target,
			"The API contract does not contain HTTP operations under paths.",
			"Define documented operations under paths before relying on API scan coverage.",
			"",
		)
	}
	for _, operation := range spec.Operations {
		operationTarget := operation.Method + " " + operation.Path
		if !spec.RootSecurity && !operation.Security {
			severity := "medium"
			if isDestructiveMethod(operation.Method) {
				severity = "high"
			}
			add(
				"openapi.operation.missing_security",
				"OpenAPI operation has no effective security requirement",
				severity,
				operationTarget,
				"The operation does not inherit root security and does not define operation-level security.",
				"Define a root security requirement or add operation-level security for this route. If the operation is intentionally public, document that exception in policy.",
				"root_security=false operation_security=false",
			)
		}
		if len(operation.Responses) > 0 && !hasClientErrorResponse(operation.Responses) {
			add(
				"openapi.operation.missing_client_error_response",
				"OpenAPI operation lacks documented client error responses",
				"low",
				operationTarget,
				"The operation documents responses but does not include a 4xx or default response.",
				"Document expected authentication, authorization, validation, or not-found error responses.",
				"responses="+strings.Join(operation.Responses, ","),
			)
		}
	}
	return findings
}

func sortSpec(spec *specDocument) {
	sort.Strings(spec.Servers)
	for i := range spec.Operations {
		sort.Strings(spec.Operations[i].Responses)
	}
	sort.Slice(spec.Operations, func(i, j int) bool {
		left := spec.Operations[i].Path + " " + spec.Operations[i].Method
		right := spec.Operations[j].Path + " " + spec.Operations[j].Method
		return left < right
	})
}

func cutYAMLKeyValue(line string) (string, string, bool) {
	key, value, ok := strings.Cut(line, ":")
	if !ok {
		return "", "", false
	}
	return key, value, true
}

func cleanYAMLKey(key string) string {
	return strings.Trim(strings.TrimSpace(key), `"'`)
}

func cleanYAMLValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	return value
}

func trimYAMLComment(line string) string {
	for i := 0; i < len(line); i++ {
		if line[i] == '#' && (i == 0 || line[i-1] == ' ' || line[i-1] == '\t') {
			return line[:i]
		}
	}
	return line
}

func leadingSpaces(s string) int {
	count := 0
	for _, r := range s {
		if r != ' ' {
			return count
		}
		count++
	}
	return count
}

func isPlainHTTPServer(raw string) bool {
	parsed, err := url.Parse(raw)
	return err == nil && strings.EqualFold(parsed.Scheme, "http")
}

func isURLSpecPath(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" {
		return false
	}
	if len(parsed.Scheme) == 1 && len(raw) > 1 && raw[1] == ':' {
		return false
	}
	return true
}

func hasClientErrorResponse(responses []string) bool {
	for _, response := range responses {
		if response == "default" || strings.HasPrefix(response, "4") {
			return true
		}
	}
	return false
}

func isResponseKey(key string) bool {
	if key == "default" {
		return true
	}
	if len(key) != 3 {
		return false
	}
	for _, r := range key {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isDestructiveMethod(method string) bool {
	switch strings.ToUpper(method) {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}

func isHTTPMethod(method string) bool {
	switch strings.ToLower(method) {
	case "get", "put", "post", "delete", "options", "head", "patch", "trace":
		return true
	default:
		return false
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
