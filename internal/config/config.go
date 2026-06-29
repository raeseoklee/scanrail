package config

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const DefaultPath = "scanrail.yaml"

type Config struct {
	ProjectName       string
	TargetURL         string
	OpenAPIPath       string
	Allowlist         []string
	ExcludePaths      []string
	BlockedPaths      []string
	TokenEnv          string
	OutputDir         string
	FailOn            string
	ActiveScanDefault bool
	RequireAllowlist  bool
	MaxRPS            int
}

func Defaults(workdir string) Config {
	name := filepath.Base(workdir)
	if name == "." || name == string(filepath.Separator) || name == "" {
		name = "scanrail-project"
	}
	return Config{
		ProjectName: name,
		TokenEnv:    "SCANRAIL_TOKEN",
		OutputDir:   ".scanrail/reports",
		FailOn:      "high",
		MaxRPS:      5,
	}
}

func Load(path string, workdir string) (Config, error) {
	cfg := Defaults(workdir)
	if path == "" {
		path = filepath.Join(workdir, DefaultPath)
	}
	data, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	defer data.Close()

	var stack []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		raw := scanner.Text()
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		indent := leadingSpaces(raw) / 2
		if indent < 0 {
			indent = 0
		}
		if strings.HasPrefix(line, "- ") {
			value := cleanValue(strings.TrimPrefix(line, "- "))
			switch strings.Join(stack, ".") {
			case "targets.web.allowlist":
				cfg.Allowlist = append(cfg.Allowlist, value)
			case "targets.web.exclude_paths":
				cfg.ExcludePaths = append(cfg.ExcludePaths, value)
			case "safety.blocked_paths":
				cfg.BlockedPaths = append(cfg.BlockedPaths, value)
			}
			continue
		}
		if strings.HasSuffix(line, ":") {
			key := strings.TrimSuffix(line, ":")
			stack = setStack(stack, indent, key)
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = cleanValue(value)
		path := appendPath(stack, indent, key)
		switch strings.Join(path, ".") {
		case "project.name":
			cfg.ProjectName = value
		case "targets.web.url":
			cfg.TargetURL = value
		case "targets.api.openapi":
			cfg.OpenAPIPath = value
		case "auth.token_env":
			cfg.TokenEnv = value
		case "report.output_dir":
			cfg.OutputDir = value
		case "policy.fail_on.severity":
			cfg.FailOn = value
		case "safety.active_scan_default":
			cfg.ActiveScanDefault = value == "true"
		case "safety.require_allowlist":
			cfg.RequireAllowlist = value == "true"
		case "safety.max_rps":
			if parsed, err := strconv.Atoi(value); err == nil {
				cfg.MaxRPS = parsed
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func Validate(cfg Config) error {
	if strings.TrimSpace(cfg.ProjectName) == "" {
		return errors.New("project.name is required")
	}
	if looksLikeSecret(cfg.TokenEnv) {
		return errors.New("auth token must be referenced by environment variable name, not stored as a literal")
	}
	return nil
}

func WriteInitial(path string, cfg Config, force bool) error {
	if path == "" {
		path = DefaultPath
	}
	if !force {
		if _, err := os.Stat(path); err == nil {
			return errors.New(path + " already exists; use --force to overwrite")
		}
	}
	apiBlock := ""
	if strings.TrimSpace(cfg.OpenAPIPath) != "" {
		apiBlock = `
  api:
    openapi: ` + cfg.OpenAPIPath + `
`
	}
	content := `project:
  name: ` + cfg.ProjectName + `
  type: web-api

targets:
  web:
    url: ` + cfg.TargetURL + `
    allowlist:
      - ` + hostOnly(cfg.TargetURL) + `
    exclude_paths:
      - /logout
` + apiBlock + `

auth:
  type: bearer
  token_env: ` + cfg.TokenEnv + `

profiles:
  default: quick

  quick:
    tools:
      - gitleaks
      - headers
      - tls
      - openapi

safety:
  active_scan_default: false
  require_allowlist: true
  max_rps: 5
  add_header:
    X-Scanrail-Project: ` + cfg.ProjectName + `
  blocked_paths:
    - /logout

policy:
  fail_on:
    severity: ` + cfg.FailOn + `

report:
  output_dir: ` + cfg.OutputDir + `
  formats:
    - html
    - json
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	envPath := ".env.scanrail.example"
	if force {
		_ = os.Remove(envPath)
	}
	if _, err := os.Stat(envPath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(envPath, []byte(cfg.TokenEnv+"=replace-me\n"), 0o644); err != nil {
			return err
		}
	}
	return nil
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

func setStack(stack []string, indent int, key string) []string {
	if indent < len(stack) {
		stack = stack[:indent]
	}
	for len(stack) < indent {
		stack = append(stack, "")
	}
	return append(stack, key)
}

func appendPath(stack []string, indent int, key string) []string {
	if indent > len(stack) {
		indent = len(stack)
	}
	out := append([]string{}, stack[:indent]...)
	return append(out, key)
}

func cleanValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	return value
}

func looksLikeSecret(value string) bool {
	if value == "" {
		return false
	}
	if strings.Contains(value, " ") {
		return true
	}
	if strings.HasPrefix(value, "eyJ") || strings.HasPrefix(value, "sk-") {
		return true
	}
	return strings.Contains(value, ".") && len(value) > 40
}

func hostOnly(raw string) string {
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	raw, _, _ = strings.Cut(raw, "/")
	return raw
}
