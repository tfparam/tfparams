package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func runCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	c := NewRootCmd()
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	c.SetArgs(args)
	err := c.Execute()
	return out.String(), err
}

func TestRootCommandTableRoot(t *testing.T) {
	out, err := runCmd(t, "--plan-json", "../testdata/plan.json", "--docs-json", "../testdata/docs.json", "--env", "production")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	for _, want := range []string{
		"# Parameter Sheet",
		"**Environment**: production",
		"**Scope**: root",
		"| instance_type |",
		"`t3.xlarge`",
		"(sensitive)",
		"`true`",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
}

func TestRootCommandModuleScope(t *testing.T) {
	out, err := runCmd(
		t,
		"--plan-json", "../testdata/plan.json",
		"--docs-json", "../testdata/docs_module.json",
		"--scope", "module", "--module", "app",
	)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(out, "## Variables (module: app)") {
		t.Errorf("missing module heading:\n%s", out)
	}
	if !strings.Contains(out, "(computed)") { // tags -> local.tags
		t.Errorf("expected computed for tags:\n%s", out)
	}
}

func TestRootCommandRejectsState(t *testing.T) {
	_, err := runCmd(t, "--plan-json", "../testdata/state.json", "--docs-json", "../testdata/docs.json")
	if err == nil {
		t.Fatal("expected error when state JSON is passed")
	}
	if !strings.Contains(err.Error(), "state") {
		t.Errorf("error should mention state, got: %v", err)
	}
}

func TestRootCommandRequiresDocs(t *testing.T) {
	_, err := runCmd(t, "--plan-json", "../testdata/plan.json")
	if err == nil {
		t.Fatal("expected error when --docs-json is missing")
	}
}

func TestRootCommandNoDefaultColumn(t *testing.T) {
	out, err := runCmd(t, "--plan-json", "../testdata/plan.json", "--docs-json", "../testdata/docs.json", "--no-default-col")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if strings.Contains(out, " Default ") {
		t.Errorf("Default column should be hidden:\n%s", out)
	}
}
