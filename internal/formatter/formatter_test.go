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

func TestInjectAppendsWhenNoMarkers(t *testing.T) {
	got, err := Inject("# Title\n\nbody\n", "TABLE")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, BeginMarker) || !strings.Contains(got, "TABLE") || !strings.Contains(got, EndMarker) {
		t.Errorf("append failed:\n%s", got)
	}
	if !strings.HasPrefix(got, "# Title") {
		t.Errorf("original content not preserved:\n%s", got)
	}
}

func TestInjectReplacesBetweenMarkers(t *testing.T) {
	existing := "# Doc\n\n" + BeginMarker + "\nOLD\n" + EndMarker + "\n\n## Keep me\n"
	got, err := Inject(existing, "NEW")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "OLD") {
		t.Errorf("old content should be replaced:\n%s", got)
	}
	if !strings.Contains(got, "NEW") || !strings.Contains(got, "## Keep me") {
		t.Errorf("inject failed to keep surrounding content:\n%s", got)
	}
}

func TestInjectUnbalancedMarkers(t *testing.T) {
	if _, err := Inject("x\n"+BeginMarker+"\nonly begin\n", "NEW"); err != ErrUnbalancedMarkers {
		t.Errorf("want ErrUnbalancedMarkers, got %v", err)
	}
}
