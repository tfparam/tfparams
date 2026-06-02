package formatter

import (
	"strings"
	"testing"

	"github.com/tfparam/tfparams/internal/merger"
)

func sampleRows() []merger.CompareRow {
	return []merger.CompareRow{
		{Name: "instance_type", Description: "type", Values: map[string]string{"dev": "t3.small", "prd": "t3.xlarge"}, Differs: true},
		{Name: "region", Description: "region", Values: map[string]string{"dev": "ap-northeast-1", "prd": "ap-northeast-1"}, Differs: false},
		{Name: "db_password", Description: "pw", Sensitive: true, Values: map[string]string{"dev": "dev-pw", "prd": "prd-pw"}},
	}
}

func TestCompareMarkdown(t *testing.T) {
	out := CompareMarkdown(sampleRows(), CompareOptions{EnvNames: []string{"dev", "prd"}, HighlightDiff: true})
	for _, want := range []string{
		"# Environment Comparison",
		"| Name | Description | dev | prd | Diff |",
		"| instance_type | type | `t3.small` | `t3.xlarge` | ⚠️ |",
		"| region | region | `ap-northeast-1` | `ap-northeast-1` | ✓ |",
		"| db_password | pw | (sensitive) | (sensitive) | - |",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("compare markdown missing %q\n%s", want, out)
		}
	}
}

func TestCompareMarkdownNoHighlight(t *testing.T) {
	out := CompareMarkdown(sampleRows(), CompareOptions{EnvNames: []string{"dev", "prd"}, HighlightDiff: false})
	if strings.Contains(out, "Diff") || strings.Contains(out, "⚠️") {
		t.Errorf("Diff column should be hidden when HighlightDiff is false:\n%s", out)
	}
}
