package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/raeseoklee/scanrail/internal/safety"
)

type Event struct {
	Timestamp  time.Time `json:"timestamp"`
	Source     string    `json:"source"`
	Action     string    `json:"action"`
	Tool       string    `json:"tool,omitempty"`
	Decision   string    `json:"decision"`
	Reason     string    `json:"reason,omitempty"`
	Project    string    `json:"project,omitempty"`
	Target     string    `json:"target,omitempty"`
	TargetHost string    `json:"target_host,omitempty"`
	Profile    string    `json:"profile,omitempty"`
	ExitCode   *int      `json:"exit_code,omitempty"`
}

func Append(path string, event Event, redactor safety.Redactor) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	event = event.Redacted(redactor)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(event)
}

func (e Event) Redacted(redactor safety.Redactor) Event {
	e.Action = redactor.RedactString(e.Action)
	e.Tool = redactor.RedactString(e.Tool)
	e.Decision = redactor.RedactString(e.Decision)
	e.Reason = redactor.RedactString(e.Reason)
	e.Project = redactor.RedactString(e.Project)
	e.Target = redactor.RedactString(e.Target)
	e.TargetHost = redactor.RedactString(e.TargetHost)
	e.Profile = redactor.RedactString(e.Profile)
	return e
}
