package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGeneratedConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scanrail.yaml")
	cfg := Defaults(dir)
	cfg.ProjectName = "demo"
	cfg.TargetURL = "http://localhost:8080"
	if err := WriteInitial(path, cfg, false); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path, dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ProjectName != "demo" {
		t.Fatalf("ProjectName = %q", loaded.ProjectName)
	}
	if loaded.TargetURL != "http://localhost:8080" {
		t.Fatalf("TargetURL = %q", loaded.TargetURL)
	}
	if loaded.TokenEnv != "SCANRAIL_TOKEN" {
		t.Fatalf("TokenEnv = %q", loaded.TokenEnv)
	}
}

func TestWriteInitialRefusesOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scanrail.yaml")
	if err := os.WriteFile(path, []byte("project:\n  name: existing\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := WriteInitial(path, Defaults(dir), false)
	if err == nil {
		t.Fatal("expected overwrite error")
	}
}

func TestValidateRejectsLiteralToken(t *testing.T) {
	cfg := Defaults(".")
	cfg.TokenEnv = "eyJ.this.looks.like.a.jwt.token.literal.and.should.not.be.stored"
	if err := Validate(cfg); err == nil {
		t.Fatal("expected validation error")
	}
}
