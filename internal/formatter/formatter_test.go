package formatter

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tfparam/tfparams/internal/merger"
)

var update = flag.Bool("update", false, "update golden files")

func sampleParams() []merger.Param {
	return []merger.Param{
		{Name: "instance_type", Description: "EC2インスタンスタイプ", Type: "string", HasDefault: true, Default: "t3.medium", HasApplied: true, Applied: "t3.xlarge"},
		{Name: "replica_count", Description: "RDSレプリカ数", Type: "number", HasDefault: true, Default: "1", HasApplied: true, Applied: "3"},
		{Name: "db_password", Description: "DBパスワード", Type: "string", Required: true, Sensitive: true, HasApplied: true, Applied: "s3cr3t"},
		{Name: "multi_az", Description: "Multi-AZ有効化", Type: "bool", HasDefault: true, Default: "false", HasApplied: true, Applied: "true"},
		{Name: "extra", Computed: true}, // no docs, unresolved module expr
	}
}

func TestMarkdownGolden(t *testing.T) {
	got := Markdown(sampleParams(), Options{
		Env:         "production",
		Scope:       "root",
		GeneratedAt: "2025-01-15 10:30:00 JST",
		Source:      "terraform show -json tfplan (plan)",
	})
	golden := filepath.Join("testdata", "root_sheet.golden.md")
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden (run with -update first): %v", err)
	}
	if got != string(want) {
		t.Errorf("markdown mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestMarkdownSensitiveShown(t *testing.T) {
	out := Markdown(sampleParams(), Options{ShowSensitive: true})
	if !strings.Contains(out, "`s3cr3t`") {
		t.Errorf("expected unmasked password, got:\n%s", out)
	}
	if strings.Contains(out, "(sensitive)") {
		t.Errorf("should not mask when ShowSensitive is true")
	}
}

func TestMarkdownNoDefaultColumn(t *testing.T) {
	cols := []string{"name", "description", "type", "applied_value", "required"}
	out := Markdown(sampleParams(), Options{Columns: cols})
	if strings.Contains(out, "Default") {
		t.Errorf("Default column should be hidden:\n%s", out)
	}
}
