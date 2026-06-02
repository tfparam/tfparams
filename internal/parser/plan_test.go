package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestParsePlan(t *testing.T) {
	in := `{
      "format_version":"1.2",
      "variables":{"instance_type":{"value":"t3.xlarge"},"replica_count":{"value":3}},
      "configuration":{"root_module":{"module_calls":{"app":{"source":"../../modules/app",
        "expressions":{"instance_type":{"references":["var.instance_type"]},"replica_count":{"constant_value":3}}}}}}
    }`
	p, err := ParsePlan(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParsePlan: %v", err)
	}
	if got := FormatValue(p.Variables["instance_type"].Value); got != "t3.xlarge" {
		t.Errorf("instance_type = %q, want t3.xlarge", got)
	}
	if got := FormatValue(p.Variables["replica_count"].Value); got != "3" {
		t.Errorf("replica_count = %q, want 3", got)
	}
	mc, ok := p.Configuration.RootModule.ModuleCalls["app"]
	if !ok {
		t.Fatal("module call app missing")
	}
	if !mc.Expressions["replica_count"].IsConstant() {
		t.Error("replica_count expression should be constant")
	}
	if mc.Expressions["instance_type"].IsConstant() {
		t.Error("instance_type expression should be a reference, not constant")
	}
}

func TestParsePlanRejectsState(t *testing.T) {
	in := `{"format_version":"1.0","values":{"root_module":{"resources":[]}}}`
	_, err := ParsePlan(strings.NewReader(in))
	if err == nil {
		t.Fatal("expected error for state input")
	}
	if !errors.Is(err, ErrLooksLikeState) {
		t.Errorf("want ErrLooksLikeState, got %v", err)
	}
}

func TestFormatValue(t *testing.T) {
	cases := map[string]string{
		`"hello"`:   "hello",
		`3`:         "3",
		`3.5`:       "3.5",
		`true`:      "true",
		`false`:     "false",
		`null`:      "null",
		`["a","b"]`: `["a","b"]`,
		`{"k":"v"}`: `{"k":"v"}`,
	}
	for in, want := range cases {
		if got := FormatValue([]byte(in)); got != want {
			t.Errorf("FormatValue(%s) = %q, want %q", in, got, want)
		}
	}
	if got := FormatValue(nil); got != "" {
		t.Errorf("FormatValue(nil) = %q, want empty", got)
	}
}
