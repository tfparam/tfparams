// Package parser parses the data sources tfparams consumes: the JSON output of
// `terraform show -json <planfile>` (plan representation) and the JSON output of
// `terraform-docs json`.
package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrLooksLikeState is returned when the provided JSON looks like a Terraform
// state representation (it has a top-level "values" key and no "variables"),
// which does not carry input variable values. tfparams requires a plan.
var ErrLooksLikeState = errors.New("input looks like Terraform state (no input variables); pass a plan JSON: terraform show -json <planfile>")

// Plan is the subset of the `terraform show -json <planfile>` output that
// tfparams reads. The full plan representation is large; only these slices are
// decoded and the rest (resource_changes, prior_state, ...) is ignored.
type Plan struct {
	FormatVersion string                  `json:"format_version"`
	Variables     map[string]PlanVariable `json:"variables"`
	Configuration Configuration           `json:"configuration"`

	// values is only used to detect that a state representation was passed by
	// mistake. It is never rendered.
	values json.RawMessage
}

// PlanVariable is a single resolved root input variable from the plan.
type PlanVariable struct {
	Value json.RawMessage `json:"value"`
}

// Configuration mirrors plan.configuration (the static config representation).
type Configuration struct {
	RootModule RootModule `json:"root_module"`
}

// RootModule holds the module calls and variable declarations of the root module.
type RootModule struct {
	ModuleCalls map[string]ModuleCall     `json:"module_calls"`
	Variables   map[string]ConfigVariable `json:"variables"`
}

// ConfigVariable is a variable declaration from
// plan.configuration.root_module.variables. terraform plan output does NOT
// redact sensitive variable values, so this `sensitive` flag is the reliable
// signal for masking (terraform-docs does not always report it).
type ConfigVariable struct {
	Description string          `json:"description"`
	Sensitive   bool            `json:"sensitive"`
	Default     json.RawMessage `json:"default"`
}

// VarSensitive reports whether the named root variable is declared sensitive in
// the plan's configuration block.
func (p *Plan) VarSensitive(name string) bool {
	v, ok := p.Configuration.RootModule.Variables[name]
	return ok && v.Sensitive
}

// ModuleCall is a single `module "<name>" {}` block in the root module.
type ModuleCall struct {
	Source      string                `json:"source"`
	Expressions map[string]Expression `json:"expressions"`
}

// Expression is the value bound to a module input. It is either a constant
// (ConstantValue present) or a set of references (e.g. ["var.instance_type"]).
type Expression struct {
	ConstantValue json.RawMessage `json:"constant_value"`
	References    []string        `json:"references"`
}

// IsConstant reports whether the expression is a literal value.
func (e Expression) IsConstant() bool { return e.ConstantValue != nil }

// ParsePlan decodes a plan JSON from r. It returns ErrLooksLikeState if the
// input is a state representation instead of a plan.
func ParsePlan(r io.Reader) (*Plan, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read plan json: %w", err)
	}
	// Decode the discriminating keys first so we can give a precise error when a
	// state file is passed by mistake.
	var probe struct {
		Variables map[string]json.RawMessage `json:"variables"`
		Values    json.RawMessage            `json:"values"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("parse plan json: %w", err)
	}
	if len(probe.Variables) == 0 && len(probe.Values) > 0 {
		return nil, ErrLooksLikeState
	}

	var p Plan
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse plan json: %w", err)
	}
	p.values = probe.Values
	return &p, nil
}

// FormatValue renders a JSON value as a compact human-readable string. Numbers
// keep their original textual form (3 stays 3, not 3.0).
func FormatValue(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return strings.TrimSpace(string(raw))
	}
	return formatAny(v)
}

func formatAny(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case json.Number:
		return x.String()
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}
