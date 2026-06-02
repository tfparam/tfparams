// Package formatter renders merged parameters into the supported output formats.
package formatter

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
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

	// Render the table with go-pretty (Markdown mode).
	tw := table.NewWriter()
	header := make(table.Row, len(cols))
	for i, c := range cols {
		header[i] = columnHeaders[c]
	}
	tw.AppendHeader(header)
	for _, p := range params {
		row := make(table.Row, len(cols))
		for i, c := range cols {
			row[i] = cell(c, p, o)
		}
		tw.AppendRow(row)
	}
	b.WriteString(tw.RenderMarkdown())
	b.WriteString("\n")

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
	text, isValue := appliedString(p, showSensitive)
	if isValue {
		return "`" + text + "`"
	}
	return text
}

// appliedString returns the semantic applied value and whether it is a concrete
// value (as opposed to a marker like "(sensitive)"). Shared by all formatters.
func appliedString(p merger.Param, showSensitive bool) (text string, isValue bool) {
	switch {
	case p.Sensitive && !showSensitive:
		return "(sensitive)", false
	case p.Computed:
		return "(computed)", false
	case !p.HasApplied:
		return "(not set)", false
	default:
		return p.Applied, true
	}
}
