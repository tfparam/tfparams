package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func TestRecursiveMode(t *testing.T) {
	root := t.TempDir()
	planDev := `{"variables":{"instance_type":{"value":"t3.small"}},"configuration":{"root_module":{"module_calls":{}}}}`
	planProd := `{"variables":{"instance_type":{"value":"t3.xlarge"}},"configuration":{"root_module":{"module_calls":{}}}}`
	docs := `{"inputs":[{"name":"instance_type","type":"string","description":"type","default":{"value":"t3.medium"}}]}`

	writeFile(t, filepath.Join(root, "dev", "tfplan.json"), planDev)
	writeFile(t, filepath.Join(root, "dev", ".tfparams.yml"), "env: development\n")
	writeFile(t, filepath.Join(root, "prod", "tfplan.json"), planProd)
	if err := os.MkdirAll(filepath.Join(root, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	docsPath := filepath.Join(root, "docs.json")
	writeFile(t, docsPath, docs)

	out, err := runCmd(t, "--recursive", "--recursive-path", root, "--docs-json", docsPath, "--out", "PARAMETERS.md")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	devContent := readFile(t, filepath.Join(root, "dev", "PARAMETERS.md"))
	if !strings.Contains(devContent, "`t3.small`") {
		t.Errorf("dev applied value missing:\n%s", devContent)
	}
	if !strings.Contains(devContent, "**Environment**: development") {
		t.Errorf("dev subdir config override not applied:\n%s", devContent)
	}

	prodContent := readFile(t, filepath.Join(root, "prod", "PARAMETERS.md"))
	if !strings.Contains(prodContent, "`t3.xlarge`") {
		t.Errorf("prod applied value missing:\n%s", prodContent)
	}
	if strings.Contains(prodContent, "Environment") {
		t.Errorf("prod should have no environment header:\n%s", prodContent)
	}

	if fileExists(filepath.Join(root, "empty", "PARAMETERS.md")) {
		t.Error("empty dir without a plan file should be skipped")
	}
	if !strings.Contains(out, "skipping") {
		t.Errorf("expected a skip warning for the empty dir:\n%s", out)
	}
}

func TestRecursiveStdout(t *testing.T) {
	root := t.TempDir()
	plan := `{"variables":{"x":{"value":"1"}},"configuration":{"root_module":{"module_calls":{}}}}`
	docs := `{"inputs":[{"name":"x","type":"string"}]}`
	writeFile(t, filepath.Join(root, "a", "tfplan.json"), plan)
	docsPath := filepath.Join(root, "docs.json")
	writeFile(t, docsPath, docs)

	out, err := runCmd(t, "--recursive", "--recursive-path", root, "--docs-json", docsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "===== ") || !strings.Contains(out, "# Parameter Sheet") {
		t.Errorf("expected stdout sheet with dir header:\n%s", out)
	}
}
