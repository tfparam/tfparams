// Package formatter renders merged parameters into the supported output formats.
package formatter

import (
	"strings"

	"github.com/tfparam/tfparams/internal/merger"
)

// DefaultColumns is the standard column order.
var DefaultColumns = []string{"name", "description", "type", "default", "applied_value", "required"}

// Options controls Markdown rendering.
type Options struct {
	Env           string
	Scope         string // "root" or "module"
	ModuleName    string
	GeneratedAt   string // injected for deterministic output
	Source        string
	ShowSensitive bool
	Columns       []string // ordered column keys; nil means DefaultColumns
}

var columnHeaders = map[string]string{
	"name":          "Name",
	"description":   "Description",
	"type":          "Type",
	"default":       "Default",
	"applied_value": "Applied Value",
	"required":      "Required",
}

// Markdown renders the parameter sheet as a Markdown document.
func Markdown(params []merger.Param, o Options) string {
	cols := o.Columns
	if cols == nil {
		cols = DefaultColumns
	}

	var b strings.Builder
	b.WriteString("# Parameter Sheet\n\n")
	if o.Env != "" {
		b.WriteString("**Environment**: " + o.Env + "\n")
	}
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
	if o.Source != "" {
		b.WriteString("**Source**: " + o.Source + "\n")
	}
	b.WriteString("\n")

	heading := "## Variables"
	if o.Scope == "module" && o.ModuleName != "" {
		heading = "## Variables (module: " + o.ModuleName + ")"
	}
	b.WriteString(heading + "\n\n")

	// Header row.
	b.WriteString("| ")
	for i, c := range cols {
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(columnHeaders[c])
	}
	b.WriteString(" |\n|")
	for range cols {
		b.WriteString("------|")
	}
	b.WriteString("\n")

	// Body rows.
	for _, p := range params {
		b.WriteString("| ")
		for i, c := range cols {
			if i > 0 {
				b.WriteString(" | ")
			}
			b.WriteString(cell(c, p, o))
		}
		b.WriteString(" |\n")
	}

	return b.String()
}

func cell(key string, p merger.Param, o Options) string {
	switch key {
	case "name":
		return p.Name
	case "description":
		if p.Description == "" {
			return "(no description)"
		}
		return p.Description
	case "type":
		if p.Type == "" {
			return "-"
		}
		return "`" + p.Type + "`"
	case "default":
		if !p.HasDefault {
			return "-"
		}
		return "`" + p.Default + "`"
	case "applied_value":
		return appliedCell(p, o.ShowSensitive)
	case "required":
		if p.Required {
			return "✓"
		}
		return "-"
	default:
		return ""
	}
}

func appliedCell(p merger.Param, showSensitive bool) string {
	if p.Sensitive && !showSensitive {
		return "(sensitive)"
	}
	if p.Computed {
		return "(computed)"
	}
	if !p.HasApplied {
		return "(not set)"
	}
	return "`" + p.Applied + "`"
}
