package report

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/raeseoklee/scanrail/internal/safety"
)

type Finding struct {
	ID          string `json:"id"`
	Tool        string `json:"tool,omitempty"`
	Title       string `json:"title"`
	Severity    string `json:"severity"`
	Confidence  string `json:"confidence"`
	Target      string `json:"target"`
	Description string `json:"description"`
	Remediation string `json:"remediation"`
	Evidence    string `json:"evidence,omitempty"`
}

type Skipped struct {
	Tool   string `json:"tool"`
	Reason string `json:"reason"`
}

type ToolMetadata struct {
	Tool    string `json:"tool"`
	Version string `json:"version,omitempty"`
	Image   string `json:"image,omitempty"`
	RawPath string `json:"raw_path,omitempty"`
}

type RunReport struct {
	Project   string         `json:"project"`
	Target    string         `json:"target,omitempty"`
	Profile   string         `json:"profile"`
	StartedAt time.Time      `json:"started_at"`
	Findings  []Finding      `json:"findings"`
	Skipped   []Skipped      `json:"skipped,omitempty"`
	Tools     []ToolMetadata `json:"tools,omitempty"`
}

func WriteJSON(path string, report RunReport) error {
	report = report.Redacted(safety.DefaultRedactor())
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func WriteHTML(path string, report RunReport) error {
	report = report.Redacted(safety.DefaultRedactor())
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return htmlTemplate.Execute(f, report)
}

func (r RunReport) Redacted(redactor safety.Redactor) RunReport {
	if r.Findings == nil {
		r.Findings = []Finding{}
	}
	r.Project = redactor.RedactString(r.Project)
	r.Target = redactor.RedactString(r.Target)
	r.Profile = redactor.RedactString(r.Profile)
	for i := range r.Findings {
		r.Findings[i].ID = redactor.RedactString(r.Findings[i].ID)
		r.Findings[i].Title = redactor.RedactString(r.Findings[i].Title)
		r.Findings[i].Severity = redactor.RedactString(r.Findings[i].Severity)
		r.Findings[i].Confidence = redactor.RedactString(r.Findings[i].Confidence)
		r.Findings[i].Target = redactor.RedactString(r.Findings[i].Target)
		r.Findings[i].Description = redactor.RedactString(r.Findings[i].Description)
		r.Findings[i].Remediation = redactor.RedactString(r.Findings[i].Remediation)
		r.Findings[i].Evidence = redactor.RedactString(r.Findings[i].Evidence)
	}
	for i := range r.Skipped {
		r.Skipped[i].Tool = redactor.RedactString(r.Skipped[i].Tool)
		r.Skipped[i].Reason = redactor.RedactString(r.Skipped[i].Reason)
	}
	for i := range r.Tools {
		r.Tools[i].Tool = redactor.RedactString(r.Tools[i].Tool)
		r.Tools[i].Version = redactor.RedactString(r.Tools[i].Version)
		r.Tools[i].Image = redactor.RedactString(r.Tools[i].Image)
		r.Tools[i].RawPath = redactor.RedactString(r.Tools[i].RawPath)
	}
	return r
}

var htmlTemplate = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Scanrail Report - {{.Project}}</title>
  <style>
    body { font-family: system-ui, sans-serif; margin: 2rem; color: #17202a; }
    table { border-collapse: collapse; width: 100%; margin-top: 1rem; }
    th, td { border: 1px solid #d8dee4; padding: .6rem; text-align: left; vertical-align: top; }
    th { background: #f6f8fa; }
    .sev-critical, .sev-high { color: #cf222e; font-weight: 700; }
    .sev-medium { color: #9a6700; font-weight: 700; }
    .sev-low { color: #0969da; font-weight: 700; }
    .empty { color: #57606a; }
  </style>
</head>
<body>
  <h1>Scanrail Report</h1>
  <p><strong>Project:</strong> {{.Project}}</p>
  <p><strong>Profile:</strong> {{.Profile}}</p>
  {{if .Target}}<p><strong>Target:</strong> {{.Target}}</p>{{end}}
  <p><strong>Started:</strong> {{.StartedAt}}</p>

  <h2>Findings</h2>
  {{if .Findings}}
  <table>
    <thead><tr><th>Tool</th><th>Severity</th><th>Title</th><th>Description</th><th>Evidence</th><th>Remediation</th></tr></thead>
    <tbody>
    {{range .Findings}}
      <tr>
        <td>{{.Tool}}</td>
        <td class="sev-{{.Severity}}">{{.Severity}}</td>
        <td>{{.Title}}</td>
        <td>{{.Description}}</td>
        <td>{{.Evidence}}</td>
        <td>{{.Remediation}}</td>
      </tr>
    {{end}}
    </tbody>
  </table>
  {{else}}
  <p class="empty">No findings.</p>
  {{end}}

  {{if .Skipped}}
  <h2>Skipped</h2>
  <table>
    <thead><tr><th>Tool</th><th>Reason</th></tr></thead>
    <tbody>
    {{range .Skipped}}<tr><td>{{.Tool}}</td><td>{{.Reason}}</td></tr>{{end}}
    </tbody>
  </table>
  {{end}}

  {{if .Tools}}
  <h2>Tool Metadata</h2>
  <table>
    <thead><tr><th>Tool</th><th>Version</th><th>Image</th><th>Raw Artifact</th></tr></thead>
    <tbody>
    {{range .Tools}}<tr><td>{{.Tool}}</td><td>{{.Version}}</td><td>{{.Image}}</td><td>{{.RawPath}}</td></tr>{{end}}
    </tbody>
  </table>
  {{end}}
</body>
</html>
`))
