package formatter

import (
	"strings"

	"github.com/tfparam/tfparams/pkg/merger"
)

// CompareOptions controls the environment comparison table.
type CompareOptions struct {
	EnvNames      []string
	GeneratedAt   string
	ShowSensitive bool
	HighlightDiff bool
	Scope         string
	ModuleName    string
}

// CompareMarkdown renders the environment comparison table.
func CompareMarkdown(rows []merger.CompareRow, o CompareOptions) string {
	var b strings.Builder
	b.WriteString("# Environment Comparison\n\n")
	if o.Scope != "" {
		scope := o.Scope
		if o.Scope == "module" && o.ModuleName != "" {
			scope = "module (" + o.ModuleName + ")"
		}
		b.WriteString("**Scope**: " + scope + "\n")
	}
	if o.GeneratedAt != "" {
		b.WriteString("**Generated at**: " + o.GeneratedAt + "\n")
	}
	b.WriteString("\n## Variables\n\n")

	cols := append([]string{"Name", "Description"}, o.EnvNames...)
	if o.HighlightDiff {
		cols = append(cols, "Diff")
	}
	b.WriteString("| " + strings.Join(cols, " | ") + " |\n|")
	for range cols {
		b.WriteString("------|")
	}
	b.WriteString("\n")

	for _, row := range rows {
		cells := []string{row.Name, descOrPlaceholder(row.Description)}
		for _, en := range o.EnvNames {
			cells = append(cells, compareCell(row, en, o.ShowSensitive))
		}
		if o.HighlightDiff {
			cells = append(cells, diffCell(row))
		}
		b.WriteString("| " + strings.Join(cells, " | ") + " |\n")
	}
	return b.String()
}

func descOrPlaceholder(d string) string {
	if d == "" {
		return "(no description)"
	}
	return d
}

func compareCell(row merger.CompareRow, env string, showSensitive bool) string {
	if row.Sensitive && !showSensitive {
		return "(sensitive)"
	}
	switch v := row.Values[env]; v {
	case "(computed)", "(not set)":
		return v
	default:
		return "`" + v + "`"
	}
}

func diffCell(row merger.CompareRow) string {
	if row.Sensitive {
		return "-"
	}
	if row.Differs {
		return "⚠️"
	}
	return "✓"
}
