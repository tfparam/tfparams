package merger

import (
	"strings"
	"testing"

	"github.com/tfkit/tfparams/pkg/parser"
)

const planJSON = `{
  "variables": {
    "instance_type": {"value": "t3.xlarge"},
    "replica_count": {"value": 3}
  },
  "configuration": {"root_module": {"module_calls": {"app": {
    "source": "../../modules/app",
    "expressions": {
      "instance_type": {"references": ["var.instance_type", "var.instance_type"]},
      "replica_count": {"constant_value": 3},
      "tags": {"references": ["local.tags"]}
    }
  }}}}
}`

const docsJSON = `{"inputs":[
  {"name":"instance_type","type":"string","description":"EC2","default":"t3.medium","required":false},
  {"name":"replica_count","type":"number","description":"replicas","default":1,"required":false}
]}`

const docsModuleJSON = `{"inputs":[
  {"name":"instance_type","type":"string","description":"type","default":"t3.small","required":false},
  {"name":"replica_count","type":"number","description":"count","default":1,"required":false},
  {"name":"tags","type":"map(string)","description":"tags","default":null,"required":true}
]}`

func mustPlan(t *testing.T, s string) *parser.Plan {
	t.Helper()
	p, err := parser.ParsePlan(strings.NewReader(s))
	if err != nil {
		t.Fatalf("ParsePlan: %v", err)
	}
	return p
}

func mustDocs(t *testing.T, s string) *parser.Docs {
	t.Helper()
	d, err := parser.ParseDocs(strings.NewReader(s))
	if err != nil {
		t.Fatalf("ParseDocs: %v", err)
	}
	return d
}

func find(params []Param, name string) (Param, bool) {
	for _, p := range params {
		if p.Name == name {
			return p, true
		}
	}
	return Param{}, false
}

func TestMergeRoot(t *testing.T) {
	plan := mustPlan(t, planJSON)
	inputs := MergeInputs(mustDocs(t, docsJSON))
	params, err := Merge(plan, inputs, ScopeRoot, "")
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if len(params) != 2 {
		t.Fatalf("want 2 params, got %d", len(params))
	}
	it, _ := find(params, "instance_type")
	if !it.HasApplied || it.Applied != "t3.xlarge" || it.Default != "t3.medium" {
		t.Errorf("instance_type = %+v", it)
	}
	rc, _ := find(params, "replica_count")
	if rc.Applied != "3" {
		t.Errorf("replica_count applied = %q", rc.Applied)
	}
}

func TestMergeModule(t *testing.T) {
	plan := mustPlan(t, planJSON)
	inputs := MergeInputs(mustDocs(t, docsModuleJSON))
	params, err := Merge(plan, inputs, ScopeModule, "app")
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	it, _ := find(params, "instance_type")
	if !it.HasApplied || it.Applied != "t3.xlarge" { // var.instance_type resolved from root
		t.Errorf("instance_type = %+v", it)
	}
	rc, _ := find(params, "replica_count")
	if !rc.HasApplied || rc.Applied != "3" { // constant_value
		t.Errorf("replica_count = %+v", rc)
	}
	tags, _ := find(params, "tags")
	if !tags.Computed || tags.HasApplied { // local.tags is not statically resolvable
		t.Errorf("tags should be computed, got %+v", tags)
	}
}

func TestMergeModuleAutoSelectAndAmbiguity(t *testing.T) {
	plan := mustPlan(t, planJSON)
	if _, err := Merge(plan, nil, ScopeModule, ""); err != nil {
		t.Errorf("single module should auto-select, got %v", err)
	}
}

func TestMergeInputsFirstWins(t *testing.T) {
	a := mustDocs(t, `{"inputs":[{"name":"x","description":"from-a"},{"name":"y","description":"y"}]}`)
	b := mustDocs(t, `{"inputs":[{"name":"x","description":"from-b"},{"name":"z","description":"z"}]}`)
	got := MergeInputs(a, b)
	if len(got) != 3 {
		t.Fatalf("want 3 inputs, got %d", len(got))
	}
	if got[0].Name != "x" || got[0].Description != "from-a" {
		t.Errorf("first-wins failed: %+v", got[0])
	}
	order := []string{got[0].Name, got[1].Name, got[2].Name}
	want := []string{"x", "y", "z"}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("order = %v, want %v", order, want)
		}
	}
}

func TestMergeAppliedOnlyAppended(t *testing.T) {
	plan := mustPlan(t, planJSON)
	// docs only declares instance_type; replica_count is applied-only and appended.
	inputs := MergeInputs(mustDocs(t, `{"inputs":[{"name":"instance_type","type":"string"}]}`))
	params, err := Merge(plan, inputs, ScopeRoot, "")
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if len(params) != 2 || params[0].Name != "instance_type" || params[1].Name != "replica_count" {
		t.Errorf("unexpected params: %+v", params)
	}
}

func TestMergeSensitiveFromPlanConfig(t *testing.T) {
	// terraform plan does NOT redact sensitive values, and terraform-docs often
	// reports sensitive=null. The plan configuration's `sensitive` flag is the
	// reliable signal; verify it marks the param sensitive even when docs do not.
	plan := mustPlan(t, `{
      "variables":{"db_password":{"value":"secret"}},
      "configuration":{"root_module":{
        "variables":{"db_password":{"sensitive":true}},
        "module_calls":{}
      }}
    }`)
	inputs := MergeInputs(mustDocs(t, `{"inputs":[{"name":"db_password","type":"string"}]}`))
	params, err := Merge(plan, inputs, ScopeRoot, "")
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	p, ok := find(params, "db_password")
	if !ok || !p.Sensitive {
		t.Errorf("db_password should be sensitive from plan config: %+v", p)
	}
}
