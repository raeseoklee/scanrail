package workspace

import (
	"os"
	"path/filepath"
	"time"
)

type Workspace struct {
	Root       string
	StateDir   string
	CacheDir   string
	RawDir     string
	ReportsDir string
	RunID      string
}

func New(root string, reportsDir string) (Workspace, error) {
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return Workspace{}, err
		}
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return Workspace{}, err
	}
	if reportsDir == "" {
		reportsDir = ".scanrail/reports"
	}
	if !filepath.IsAbs(reportsDir) {
		reportsDir = filepath.Join(abs, reportsDir)
	}
	runID := time.Now().UTC().Format("20060102-150405")
	return Workspace{
		Root:       abs,
		StateDir:   filepath.Join(abs, ".scanrail"),
		CacheDir:   filepath.Join(abs, ".scanrail", "cache"),
		RawDir:     filepath.Join(abs, ".scanrail", "raw", runID),
		ReportsDir: reportsDir,
		RunID:      runID,
	}, nil
}

func (w Workspace) Ensure() error {
	for _, path := range []string{w.StateDir, w.CacheDir, w.RawDir, w.ReportsDir} {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return err
		}
	}
	return nil
}
