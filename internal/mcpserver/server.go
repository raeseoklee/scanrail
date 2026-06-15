package mcpserver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/raeseoklee/scanrail/internal/app"
	"github.com/raeseoklee/scanrail/internal/audit"
	"github.com/raeseoklee/scanrail/internal/config"
	"github.com/raeseoklee/scanrail/internal/exitcode"
	"github.com/raeseoklee/scanrail/internal/report"
	"github.com/raeseoklee/scanrail/internal/safety"
	"github.com/raeseoklee/scanrail/internal/version"
)

const protocolVersion = "2025-06-18"

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type textContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolResult struct {
	Content           []textContent `json:"content"`
	IsError           bool          `json:"isError,omitempty"`
	StructuredContent any           `json:"structuredContent,omitempty"`
}

type server struct {
	in      io.Reader
	out     io.Writer
	stderr  io.Writer
	workdir string
}

func Serve(ctx context.Context, in io.Reader, out io.Writer, stderr io.Writer) int {
	if stderr == nil {
		stderr = io.Discard
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, "scanrail mcp: cannot determine working directory:", err)
		return exitcode.Environment
	}
	s := server{in: in, out: out, stderr: stderr, workdir: wd}
	return s.serve(ctx)
}

func (s server) serve(ctx context.Context) int {
	scanner := bufio.NewScanner(s.in)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	encoder := json.NewEncoder(s.out)
	encoder.SetEscapeHTML(false)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return exitcode.Interrupted
		default:
		}
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			_ = encoder.Encode(errorResponse(json.RawMessage("null"), -32700, "Parse error", err.Error()))
			continue
		}
		if len(req.ID) == 0 {
			s.handleNotification(req)
			continue
		}
		if req.JSONRPC != "2.0" || req.Method == "" {
			_ = encoder.Encode(errorResponse(req.ID, -32600, "Invalid Request", nil))
			continue
		}
		result, rpcErr := s.handleRequest(ctx, req)
		if rpcErr != nil {
			_ = encoder.Encode(rpcResponse{JSONRPC: "2.0", ID: req.ID, Error: rpcErr})
			continue
		}
		_ = encoder.Encode(rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: result})
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(s.stderr, "scanrail mcp: input error:", err)
		return exitcode.Environment
	}
	return exitcode.OK
}

func (s server) handleNotification(req rpcRequest) {
	if req.Method != "notifications/initialized" {
		fmt.Fprintln(s.stderr, "scanrail mcp: ignored notification:", req.Method)
	}
}

func (s server) handleRequest(ctx context.Context, req rpcRequest) (any, *rpcError) {
	switch req.Method {
	case "initialize":
		return map[string]any{
			"protocolVersion": protocolVersion,
			"capabilities": map[string]any{
				"tools":     map[string]any{},
				"resources": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "scanrail",
				"title":   "Scanrail MCP Server",
				"version": version.String(),
			},
		}, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		return map[string]any{"tools": toolDefinitions()}, nil
	case "tools/call":
		return s.callTool(ctx, req.Params)
	case "resources/list":
		return map[string]any{"resources": resourceDefinitions()}, nil
	case "resources/read":
		return s.readResource(req.Params)
	default:
		return nil, &rpcError{Code: -32601, Message: "Method not found"}
	}
}

func (s server) callTool(ctx context.Context, params json.RawMessage) (any, *rpcError) {
	var call struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &call); err != nil {
		return nil, &rpcError{Code: -32602, Message: "Invalid params", Data: err.Error()}
	}
	if len(call.Arguments) == 0 {
		call.Arguments = json.RawMessage("{}")
	}
	switch call.Name {
	case "scanrail_doctor":
		return textToolResult(runDoctor()), nil
	case "scanrail_config_read":
		return s.configRead(call.Arguments)
	case "scanrail_report_latest":
		return s.reportLatest(call.Arguments)
	case "scanrail_run":
		return s.runHeaders(ctx, call.Arguments)
	default:
		return nil, &rpcError{Code: -32602, Message: "Unknown tool", Data: call.Name}
	}
}

func runDoctor() string {
	var out bytes.Buffer
	app.Doctor(&out)
	return out.String()
}

func (s server) configRead(args json.RawMessage) (any, *rpcError) {
	var opts struct {
		ConfigPath string `json:"config_path"`
	}
	if err := json.Unmarshal(args, &opts); err != nil {
		return nil, &rpcError{Code: -32602, Message: "Invalid params", Data: err.Error()}
	}
	cfg, err := config.Load(configPath(opts.ConfigPath, s.workdir), s.workdir)
	if err != nil {
		return nil, &rpcError{Code: -32603, Message: "Could not read config", Data: err.Error()}
	}
	content := redactedConfig(cfg)
	text := mustJSON(content)
	return toolResult{
		Content:           []textContent{{Type: "text", Text: text}},
		StructuredContent: content,
	}, nil
}

func (s server) reportLatest(args json.RawMessage) (any, *rpcError) {
	var opts struct {
		ConfigPath string `json:"config_path"`
		OutputDir  string `json:"output_dir"`
	}
	if err := json.Unmarshal(args, &opts); err != nil {
		return nil, &rpcError{Code: -32602, Message: "Invalid params", Data: err.Error()}
	}
	summary, err := s.latestReport(opts.ConfigPath, opts.OutputDir)
	if err != nil {
		return errorToolResult(err.Error()), nil
	}
	redactor := safety.DefaultRedactor()
	summary = redactor.RedactValue(summary).(map[string]any)
	text := redactor.RedactString(mustJSON(summary))
	return toolResult{
		Content:           []textContent{{Type: "text", Text: text}},
		StructuredContent: summary,
	}, nil
}

func (s server) runHeaders(ctx context.Context, args json.RawMessage) (any, *rpcError) {
	var opts struct {
		ConfigPath        string `json:"config_path"`
		Target            string `json:"target"`
		OutputDir         string `json:"output_dir"`
		Profile           string `json:"profile"`
		Only              string `json:"only"`
		ConfirmActiveScan bool   `json:"confirm_active_scan"`
	}
	if err := json.Unmarshal(args, &opts); err != nil {
		return nil, &rpcError{Code: -32602, Message: "Invalid params", Data: err.Error()}
	}
	if opts.Only != "" && opts.Only != "headers" {
		reason := "MCP MVP only supports the native headers scanner."
		_ = s.writeMCPAudit(audit.Event{
			Action:   "scanrail_run",
			Tool:     opts.Only,
			Decision: "denied",
			Reason:   reason,
			Target:   opts.Target,
			Profile:  opts.Profile,
		}, safety.DefaultRedactor())
		return errorToolResult(reason), nil
	}
	if !opts.ConfirmActiveScan {
		reason := "scanrail_run requires confirm_active_scan=true because it sends an HTTP request to the target."
		_ = s.writeMCPAudit(audit.Event{
			Action:   "scanrail_run",
			Tool:     scanTool(opts.Only),
			Decision: "denied",
			Reason:   reason,
			Target:   opts.Target,
			Profile:  opts.Profile,
		}, safety.DefaultRedactor())
		return errorToolResult(reason), nil
	}
	cfg, err := config.Load(configPath(opts.ConfigPath, s.workdir), s.workdir)
	if err != nil {
		_ = s.writeMCPAudit(audit.Event{
			Action:   "scanrail_run",
			Tool:     scanTool(opts.Only),
			Decision: "denied",
			Reason:   "could not read config: " + err.Error(),
			Target:   opts.Target,
			Profile:  opts.Profile,
		}, safety.DefaultRedactor())
		return nil, &rpcError{Code: -32603, Message: "Could not read config", Data: err.Error()}
	}
	redactor := safety.NewRedactorFromEnv(cfg.TokenEnv)
	target := opts.Target
	if target == "" {
		target = cfg.TargetURL
	}
	if err := allowedTarget(cfg, target); err != nil {
		_ = s.writeMCPAudit(audit.Event{
			Action:     "scanrail_run",
			Tool:       "headers",
			Decision:   "denied",
			Reason:     err.Error(),
			Project:    cfg.ProjectName,
			Target:     target,
			TargetHost: targetHostOrEmpty(target),
			Profile:    opts.Profile,
		}, redactor)
		return errorToolResult(err.Error()), nil
	}
	started := audit.Event{
		Action:     "scanrail_run",
		Tool:       "headers",
		Decision:   "started",
		Project:    cfg.ProjectName,
		Target:     target,
		TargetHost: targetHostOrEmpty(target),
		Profile:    opts.Profile,
	}
	if err := s.writeMCPAudit(started, redactor); err != nil {
		message := "scanrail_run could not write audit log; refusing to execute active scan."
		return errorToolResult(redactor.RedactString(message + " " + err.Error())), nil
	}
	var out bytes.Buffer
	code := app.Run(ctx, app.RunOptions{
		ConfigPath: configPath(opts.ConfigPath, s.workdir),
		Profile:    opts.Profile,
		Target:     target,
		Only:       "headers",
		OutputDir:  opts.OutputDir,
	}, &out)
	completed := started
	completed.Decision = "completed"
	completed.ExitCode = &code
	auditErr := s.writeMCPAudit(completed, redactor)
	redactedOutput := redactor.RedactString(out.String())
	result := map[string]any{
		"exit_code": code,
		"output":    strings.TrimSpace(redactedOutput),
	}
	if auditErr != nil {
		result["audit_logged"] = false
		result["audit_log_error"] = redactor.RedactString(auditErr.Error())
		return toolResult{
			Content:           []textContent{{Type: "text", Text: redactedOutput + "\naudit log failed after scan completion"}},
			IsError:           true,
			StructuredContent: result,
		}, nil
	}
	result["audit_logged"] = true
	if code != exitcode.OK {
		return toolResult{
			Content:           []textContent{{Type: "text", Text: redactedOutput}},
			IsError:           true,
			StructuredContent: result,
		}, nil
	}
	return toolResult{
		Content:           []textContent{{Type: "text", Text: redactedOutput}},
		StructuredContent: result,
	}, nil
}

func (s server) readResource(params json.RawMessage) (any, *rpcError) {
	var req struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, &rpcError{Code: -32602, Message: "Invalid params", Data: err.Error()}
	}
	var (
		text string
		mime = "application/json"
		err  error
	)
	switch req.URI {
	case "scanrail://config":
		var cfg config.Config
		cfg, err = config.Load(configPath("", s.workdir), s.workdir)
		if err == nil {
			text = mustJSON(redactedConfig(cfg))
		}
	case "scanrail://reports/latest/summary":
		var summary any
		summary, err = s.latestReport("", "")
		if err == nil {
			redactor := safety.DefaultRedactor()
			text = redactor.RedactString(mustJSON(redactor.RedactValue(summary)))
		}
	case "scanrail://safety-model":
		mime = "text/markdown"
		text = safetyModelText()
	default:
		return nil, &rpcError{Code: -32002, Message: "Resource not found", Data: map[string]string{"uri": req.URI}}
	}
	if err != nil {
		return nil, &rpcError{Code: -32603, Message: "Could not read resource", Data: err.Error()}
	}
	return map[string]any{
		"contents": []map[string]string{{
			"uri":      req.URI,
			"mimeType": mime,
			"text":     text,
		}},
	}, nil
}

func toolDefinitions() []map[string]any {
	return []map[string]any{
		{
			"name":        "scanrail_doctor",
			"title":       "Run Scanrail Doctor",
			"description": "Check local Scanrail runtime readiness without running a scan.",
			"inputSchema": objectSchema(nil, nil),
		},
		{
			"name":        "scanrail_config_read",
			"title":       "Read Scanrail Config",
			"description": "Return normalized Scanrail configuration with secret values redacted.",
			"inputSchema": objectSchema(map[string]any{
				"config_path": map[string]any{"type": "string", "description": "Optional path to scanrail.yaml."},
			}, nil),
		},
		{
			"name":        "scanrail_report_latest",
			"title":       "Read Latest Scanrail Report",
			"description": "Return a bounded summary of the latest JSON report.",
			"inputSchema": objectSchema(map[string]any{
				"config_path": map[string]any{"type": "string", "description": "Optional path to scanrail.yaml."},
				"output_dir":  map[string]any{"type": "string", "description": "Optional report directory override."},
			}, nil),
		},
		{
			"name":        "scanrail_run",
			"title":       "Run Native Headers Scan",
			"description": "Run the native headers scanner only. Requires explicit active-scan confirmation and target allowlist validation.",
			"inputSchema": objectSchema(map[string]any{
				"config_path":         map[string]any{"type": "string", "description": "Optional path to scanrail.yaml."},
				"target":              map[string]any{"type": "string", "description": "Optional target URL override."},
				"output_dir":          map[string]any{"type": "string", "description": "Optional report output directory."},
				"profile":             map[string]any{"type": "string", "description": "Optional profile name."},
				"only":                map[string]any{"type": "string", "enum": []string{"headers"}, "description": "MCP MVP only supports headers."},
				"confirm_active_scan": map[string]any{"type": "boolean", "description": "Must be true to send an HTTP request."},
			}, []string{"confirm_active_scan"}),
		},
	}
}

func resourceDefinitions() []map[string]any {
	return []map[string]any{
		{
			"uri":         "scanrail://config",
			"name":        "config",
			"title":       "Scanrail Config",
			"description": "Normalized project configuration with secret values redacted.",
			"mimeType":    "application/json",
		},
		{
			"uri":         "scanrail://reports/latest/summary",
			"name":        "latest-report-summary",
			"title":       "Latest Report Summary",
			"description": "Bounded summary of the latest Scanrail JSON report.",
			"mimeType":    "application/json",
		},
		{
			"uri":         "scanrail://safety-model",
			"name":        "safety-model",
			"title":       "Scanrail MCP Safety Model",
			"description": "Effective MCP safety rules for local tool execution.",
			"mimeType":    "text/markdown",
		},
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	if properties == nil {
		properties = map[string]any{}
	}
	schema := map[string]any{
		"type":                 "object",
		"properties":           properties,
		"additionalProperties": false,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func (s server) latestReport(configPathArg string, outputDir string) (map[string]any, error) {
	dir, err := s.reportDir(configPathArg, outputDir)
	if err != nil {
		return nil, err
	}
	path, err := latestJSONFile(dir)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rr report.RunReport
	if err := json.Unmarshal(data, &rr); err != nil {
		return nil, err
	}
	bySeverity := map[string]int{}
	for _, finding := range rr.Findings {
		bySeverity[finding.Severity]++
	}
	skipped := make([]any, 0, len(rr.Skipped))
	for _, item := range rr.Skipped {
		skipped = append(skipped, map[string]any{
			"tool":   item.Tool,
			"reason": item.Reason,
		})
	}
	return map[string]any{
		"path":                 path,
		"project":              rr.Project,
		"target":               rr.Target,
		"profile":              rr.Profile,
		"started_at":           rr.StartedAt,
		"findings_count":       len(rr.Findings),
		"findings_by_severity": bySeverity,
		"skipped":              skipped,
	}, nil
}

func (s server) reportDir(configPathArg string, outputDir string) (string, error) {
	if outputDir == "" {
		cfg, err := config.Load(configPath(configPathArg, s.workdir), s.workdir)
		if err != nil {
			return "", err
		}
		outputDir = cfg.OutputDir
	}
	if outputDir == "" {
		outputDir = ".scanrail/reports"
	}
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(s.workdir, outputDir)
	}
	return outputDir, nil
}

func latestJSONFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}
	if len(files) == 0 {
		return "", errors.New("no JSON reports found")
	}
	sort.Slice(files, func(i, j int) bool {
		left, _ := os.Stat(files[i])
		right, _ := os.Stat(files[j])
		if left == nil || right == nil {
			return files[i] > files[j]
		}
		return left.ModTime().After(right.ModTime())
	})
	return files[0], nil
}

func configPath(path string, workdir string) string {
	if path == "" {
		return filepath.Join(workdir, config.DefaultPath)
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(workdir, path)
}

func redactedConfig(cfg config.Config) map[string]any {
	redactor := safety.NewRedactorFromEnv(cfg.TokenEnv)
	return map[string]any{
		"project_name":        cfg.ProjectName,
		"target_url":          redactor.RedactString(cfg.TargetURL),
		"allowlist":           redactor.RedactValue(cfg.Allowlist),
		"auth":                map[string]string{"token_env": cfg.TokenEnv},
		"output_dir":          redactor.RedactString(cfg.OutputDir),
		"fail_on":             cfg.FailOn,
		"active_scan_default": cfg.ActiveScanDefault,
	}
}

func allowedTarget(cfg config.Config, target string) error {
	targetHost, err := hostname(target)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}
	if targetHost == "" {
		return errors.New("target URL must include a hostname")
	}
	if cfg.TargetURL != "" {
		cfgHost, err := hostname(cfg.TargetURL)
		if err == nil && targetHost == cfgHost {
			return nil
		}
	}
	for _, allowed := range cfg.Allowlist {
		if targetHost == normalizeHost(allowed) {
			return nil
		}
	}
	return fmt.Errorf("target host %q is outside the configured Scanrail allowlist", targetHost)
}

func hostname(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	return strings.ToLower(parsed.Hostname()), nil
}

func normalizeHost(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return ""
	}
	if strings.Contains(raw, "://") {
		host, err := hostname(raw)
		if err == nil {
			return host
		}
	}
	if host, _, err := net.SplitHostPort(raw); err == nil {
		return strings.ToLower(host)
	}
	if host, _, ok := strings.Cut(raw, ":"); ok {
		return host
	}
	return raw
}

func textToolResult(text string) toolResult {
	return toolResult{Content: []textContent{{Type: "text", Text: text}}}
}

func errorToolResult(text string) toolResult {
	return toolResult{Content: []textContent{{Type: "text", Text: text}}, IsError: true}
}

func errorResponse(id json.RawMessage, code int, message string, data any) rpcResponse {
	return rpcResponse{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: message, Data: data}}
}

func mustJSON(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

func safetyModelText() string {
	return strings.TrimSpace(`# Scanrail MCP Safety Model

- The MCP server uses stdio only and does not open a network listener.
- Tools never execute arbitrary shell commands.
- ` + "`scanrail_run`" + ` only runs the native headers scanner in the MVP.
- Active scan execution requires ` + "`confirm_active_scan=true`" + `.
- Scan targets must match the configured target host or ` + "`targets.web.allowlist`" + `.
- Secret values are not accepted in MCP inputs or returned through MCP resources.
- MCP-triggered scan attempts are recorded in ` + "`.scanrail/logs/mcp-audit.jsonl`" + `.
- Reports are summarized by default instead of streaming unbounded raw output.`)
}

func (s server) auditLogPath() string {
	return filepath.Join(s.workdir, ".scanrail", "logs", "mcp-audit.jsonl")
}

func (s server) writeMCPAudit(event audit.Event, redactor safety.Redactor) error {
	event.Source = "mcp"
	return audit.Append(s.auditLogPath(), event, redactor)
}

func scanTool(value string) string {
	if value == "" {
		return "headers"
	}
	return value
}

func targetHostOrEmpty(raw string) string {
	host, err := hostname(raw)
	if err != nil {
		return ""
	}
	return host
}
