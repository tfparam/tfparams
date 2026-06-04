// Package merger joins applied values (from a plan) with variable metadata
// (from terraform-docs) into a flat list of parameters ready for formatting.
package merger

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tfkit/tfparams/pkg/parser"
)

// Scope selects which applied values are used.
type Scope string

const (
	// ScopeRoot uses the plan's root input variables.
	ScopeRoot Scope = "root"
	// ScopeModule uses the values a module call was given.
	ScopeModule Scope = "module"
)

// Param is one row of the parameter sheet.
type Param struct {
	Name        string
	Description string
	Type        string
	HasDefault  bool
	Default     string
	HasApplied  bool // a concrete, resolved applied value is present
	Applied     string
	Computed    bool // module expression could not be resolved statically
	Required    bool
	Sensitive   bool
}

// applied is an intermediate resolved value keyed by variable name.
type applied struct {
	value     string
	computed  bool
	sensitive bool
}

// MergeInputs combines several terraform-docs documents, first-wins on name,
// preserving the order of first appearance.
func MergeInputs(docs ...*parser.Docs) []parser.Input {
	var out []parser.Input
	seen := map[string]bool{}
	for _, d := range docs {
		for _, in := range d.Inputs {
			if seen[in.Name] {
				continue
			}
			seen[in.Name] = true
			out = append(out, in)
		}
	}
	return out
}

// Merge produces the parameter rows. inputs is the (already merged) terraform-docs
// metadata; scope/module select the applied-value source from the plan.
func Merge(plan *parser.Plan, inputs []parser.Input, scope Scope, module string) ([]Param, error) {
	values, err := appliedValues(plan, scope, module)
	if err != nil {
		return nil, err
	}

	var params []Param
	used := map[string]bool{}
	for _, in := range inputs {
		p := Param{
			Name:        in.Name,
			Description: in.Description,
			Type:        in.Type,
			Required:    in.Required,
			Sensitive:   in.Sensitive,
		}
		if in.HasDefault() {
			p.HasDefault = true
			p.Default = parser.FormatValue(in.Default)
		}
		if a, ok := values[in.Name]; ok {
			applyValue(&p, a)
			if a.sensitive {
				p.Sensitive = true
			}
			used[in.Name] = true
		}
		params = append(params, p)
	}

	// Applied values that have no metadata are appended at the end, sorted by name.
	var extra []string
	for name := range values {
		if !used[name] {
			extra = append(extra, name)
		}
	}
	sort.Strings(extra)
	for _, name := range extra {
		p := Param{Name: name, Sensitive: values[name].sensitive}
		applyValue(&p, values[name])
		params = append(params, p)
	}

	return params, nil
}

func applyValue(p *Param, a applied) {
	if a.computed {
		p.Computed = true
		return
	}
	p.HasApplied = true
	p.Applied = a.value
}

// appliedValues resolves the applied value for each variable name in the chosen scope.
func appliedValues(plan *parser.Plan, scope Scope, module string) (map[string]applied, error) {
	switch scope {
	case ScopeRoot, "":
		m := make(map[string]applied, len(plan.Variables))
		for name, v := range plan.Variables {
			m[name] = applied{value: parser.FormatValue(v.Value), sensitive: plan.VarSensitive(name)}
		}
		return m, nil
	case ScopeModule:
		_, mc, err := selectModule(plan, module)
		if err != nil {
			return nil, err
		}
		m := make(map[string]applied, len(mc.Expressions))
		for name, expr := range mc.Expressions {
			m[name] = resolveExpression(plan, expr)
		}
		return m, nil
	default:
		return nil, fmt.Errorf("unknown scope %q (want root or module)", scope)
	}
}

// resolveExpression turns a module-call argument into a concrete value when possible.
func resolveExpression(plan *parser.Plan, expr parser.Expression) applied {
	if expr.IsConstant() {
		return applied{value: parser.FormatValue(expr.ConstantValue)}
	}
	if name, ok := singleVarRef(expr.References); ok {
		if pv, exists := plan.Variables[name]; exists {
			return applied{value: parser.FormatValue(pv.Value), sensitive: plan.VarSensitive(name)}
		}
	}
	return applied{computed: true}
}

// singleVarRef returns the variable name when the references all point to a
// single `var.<name>`. Plan output sometimes repeats the reference.
func singleVarRef(refs []string) (string, bool) {
	uniq := map[string]bool{}
	for _, r := range refs {
		uniq[r] = true
	}
	if len(uniq) != 1 {
		return "", false
	}
	var only string
	for r := range uniq {
		only = r
	}
	if strings.HasPrefix(only, "var.") {
		return strings.TrimPrefix(only, "var."), true
	}
	return "", false
}

// ModuleName resolves the effective module call name for the given selector,
// applying the same auto-selection rules as Merge.
func ModuleName(plan *parser.Plan, module string) (string, error) {
	name, _, err := selectModule(plan, module)
	return name, err
}

// selectModule picks the target module call. When module is empty it auto-selects
// the only call, or errors if there are zero or many.
func selectModule(plan *parser.Plan, module string) (string, parser.ModuleCall, error) {
	calls := plan.Configuration.RootModule.ModuleCalls
	if len(calls) == 0 {
		return "", parser.ModuleCall{}, fmt.Errorf("no module calls found in plan configuration")
	}
	if module != "" {
		mc, ok := calls[module]
		if !ok {
			return "", parser.ModuleCall{}, fmt.Errorf("module call %q not found in plan", module)
		}
		return module, mc, nil
	}
	if len(calls) == 1 {
		for name, mc := range calls {
			return name, mc, nil
		}
	}
	names := make([]string, 0, len(calls))
	for n := range calls {
		names = append(names, n)
	}
	sort.Strings(names)
	return "", parser.ModuleCall{}, fmt.Errorf("multiple module calls (%s); specify --module", strings.Join(names, ", "))
}
