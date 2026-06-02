package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	c := Default()
	if c.Format != "table" {
		t.Errorf("Format = %q, want table", c.Format)
	}
	if c.Scope != "root" {
		t.Errorf("Scope = %q, want root", c.Scope)
	}
	if c.Recursive.Path != "." {
		t.Errorf("Recursive.Path = %q, want .", c.Recursive.Path)
	}
	if c.Sensitive.Mask != "(sensitive)" {
		t.Errorf("Sensitive.Mask = %q", c.Sensitive.Mask)
	}
}

func TestLoadExplicit(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "x.yml")
	body := "format: csv\nenv: prod\nscope: module\nmodule: app\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	c, found, err := Load(p)
	if err != nil || !found {
		t.Fatalf("Load: found=%v err=%v", found, err)
	}
	if c.Format != "csv" || c.Env != "prod" || c.Scope != "module" || c.Module != "app" {
		t.Errorf("loaded config = %+v", c)
	}
	// Defaults must be preserved for unspecified keys.
	if c.Sensitive.Mask != "(sensitive)" {
		t.Errorf("default mask lost: %q", c.Sensitive.Mask)
	}
}

func TestLoadExplicitMissing(t *testing.T) {
	_, _, err := Load(filepath.Join(t.TempDir(), "nope.yml"))
	if err == nil {
		t.Error("expected error for missing explicit config")
	}
}

func TestLoadSearchNotFound(t *testing.T) {
	// Run from an empty dir (and isolated HOME) so no .tfparams.yml is discovered.
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	wd, _ := os.Getwd()
	defer func() { _ = os.Chdir(wd) }()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	c, found, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if found {
		t.Errorf("did not expect to find a config in empty dir")
	}
	if c.Format != "table" {
		t.Errorf("expected defaults, got %+v", c)
	}
}
