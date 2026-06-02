package merger

// EnvResult is one environment's merged parameters, labeled by environment name.
type EnvResult struct {
	Name   string
	Params []Param
}

// CompareRow is one variable's values across all environments.
type CompareRow struct {
	Name        string
	Description string
	Type        string
	Sensitive   bool
	Values      map[string]string // env name -> semantic value (value / "(not set)" / "(computed)")
	Differs     bool              // true if non-sensitive values differ across environments
}

// Compare builds comparison rows across environments. Row order follows the
// first appearance of each variable; names only seen in later environments are
// appended. Missing values are rendered as "(not set)".
func Compare(envs []EnvResult) []CompareRow {
	order := []string{}
	byName := map[string]*CompareRow{}

	for _, env := range envs {
		for _, p := range env.Params {
			row, ok := byName[p.Name]
			if !ok {
				row = &CompareRow{Name: p.Name, Description: p.Description, Type: p.Type, Values: map[string]string{}}
				byName[p.Name] = row
				order = append(order, p.Name)
			}
			if row.Description == "" && p.Description != "" {
				row.Description = p.Description
			}
			if p.Sensitive {
				row.Sensitive = true
			}
			row.Values[env.Name] = semanticValue(p)
		}
	}

	envNames := make([]string, len(envs))
	for i, e := range envs {
		envNames[i] = e.Name
	}

	rows := make([]CompareRow, 0, len(order))
	for _, name := range order {
		row := byName[name]
		for _, en := range envNames {
			if _, ok := row.Values[en]; !ok {
				row.Values[en] = "(not set)"
			}
		}
		row.Differs = differs(row, envNames)
		rows = append(rows, *row)
	}
	return rows
}

func semanticValue(p Param) string {
	switch {
	case p.Computed:
		return "(computed)"
	case !p.HasApplied:
		return "(not set)"
	default:
		return p.Applied
	}
}

func differs(row *CompareRow, envNames []string) bool {
	if row.Sensitive {
		return false // sensitive rows are not diff-highlighted
	}
	var first string
	for i, en := range envNames {
		if i == 0 {
			first = row.Values[en]
			continue
		}
		if row.Values[en] != first {
			return true
		}
	}
	return false
}
