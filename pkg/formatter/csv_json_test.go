package formatter

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCSV(t *testing.T) {
	out, err := CSV(sampleParams(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 6 { // header + 5 params
		t.Fatalf("lines = %d, want 6:\n%s", len(lines), out)
	}
	if lines[0] != "Name,Description,Type,Default,Applied Value,Required" {
		t.Errorf("header = %q", lines[0])
	}
	if !strings.Contains(out, "(sensitive)") {
		t.Errorf("expected masked sensitive value:\n%s", out)
	}
	if !strings.Contains(out, "(computed)") {
		t.Errorf("expected computed marker:\n%s", out)
	}
	if !strings.Contains(out, "instance_type,EC2インスタンスタイプ,string,t3.medium,t3.xlarge,false") {
		t.Errorf("unexpected instance_type row:\n%s", out)
	}
}

func TestJSON(t *testing.T) {
	out, err := JSON(sampleParams(), Options{Env: "prod", Scope: "root"})
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	for _, want := range []string{
		`"environment": "prod"`,
		`"applied_value": "t3.xlarge"`,
		`"applied_value": "(sensitive)"`,
		`"computed": true`,
		`"applied_value": null`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("json missing %q\n%s", want, out)
		}
	}
}
