package cmd

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRootCommandCSV(t *testing.T) {
	out, err := runCmd(t, "--plan-json", "../testdata/plan.json", "--docs-json", "../testdata/docs.json", "--format", "csv")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.HasPrefix(out, "Name,Description,Type,Default,Applied Value,Required") {
		t.Errorf("csv header missing:\n%s", out)
	}
	if !strings.Contains(out, "instance_type,EC2インスタンスタイプ,string,t3.medium,t3.xlarge,false") {
		t.Errorf("csv instance_type row missing:\n%s", out)
	}
}

func TestRootCommandJSON(t *testing.T) {
	out, err := runCmd(t, "--plan-json", "../testdata/plan.json", "--docs-json", "../testdata/docs.json", "--format", "json")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("invalid json output: %v\n%s", err, out)
	}
	if _, ok := doc["variables"]; !ok {
		t.Errorf("json missing variables key:\n%s", out)
	}
}
