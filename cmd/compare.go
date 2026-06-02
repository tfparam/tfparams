package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tfparam/tfparams/pkg/backend"
	"github.com/tfparam/tfparams/pkg/formatter"
	"github.com/tfparam/tfparams/pkg/merger"
	"github.com/tfparam/tfparams/pkg/parser"
)

type compareFlags struct {
	envs          []string
	docsJSON      []string
	out           string
	format        string
	highlightDiff bool
	showSensitive bool
	sortBy        string
	scope         string
	module        string
}

func newCompareCmd() *cobra.Command {
	f := &compareFlags{}
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare applied values across environments from their plan JSONs",
		Long: "compare fetches each environment's plan JSON (local path or s3/gs/azblob URI), " +
			"aligns the variables side by side, and highlights rows that differ.",
		RunE:          func(cmd *cobra.Command, _ []string) error { return runCompare(cmd, f) },
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	fl := cmd.Flags()
	fl.StringArrayVar(&f.envs, "env", nil, "name=<uri-or-path> to a plan JSON (repeat; at least two required)")
	fl.StringArrayVar(&f.docsJSON, "docs-json", nil, "terraform-docs JSON, e.g. <(terraform-docs json .) (required, repeatable)")
	fl.StringVar(&f.out, "out", "", "output file (overwritten); default stdout")
	fl.StringVar(&f.format, "format", "markdown", "output format: markdown (csv/json TBD)")
	fl.BoolVar(&f.highlightDiff, "highlight-diff", true, "highlight rows that differ across environments")
	fl.BoolVar(&f.showSensitive, "show-sensitive", false, "show sensitive values unmasked")
	fl.StringVar(&f.sortBy, "sort-by", "required", "sort key: required (required first, then name) or name")
	fl.StringVar(&f.scope, "scope", "root", "scope: root or module")
	fl.StringVar(&f.module, "module", "", "module call name for --scope module")
	return cmd
}

func runCompare(cmd *cobra.Command, f *compareFlags) error {
	if len(f.envs) < 2 {
		return fmt.Errorf("compare requires at least two --env name=<uri-or-path> entries")
	}
	if len(f.docsJSON) == 0 {
		return fmt.Errorf("--docs-json is required")
	}
	if f.format != "markdown" && f.format != "table" {
		return fmt.Errorf("format %q is not implemented for compare yet (only 'markdown')", f.format)
	}

	var docs []*parser.Docs
	for _, p := range f.docsJSON {
		d, err := readDocs(p)
		if err != nil {
			return err
		}
		docs = append(docs, d)
	}
	inputs := merger.MergeInputs(docs...)

	ctx := context.Background()

	var envResults []merger.EnvResult
	var envNames []string
	moduleName := ""
	for _, entry := range f.envs {
		name, uri, ok := strings.Cut(entry, "=")
		if !ok || name == "" || uri == "" {
			return fmt.Errorf("invalid --env %q (want name=<uri-or-path>)", entry)
		}
		data, err := backend.Fetch(ctx, uri)
		if err != nil {
			return fmt.Errorf("env %s: %w", name, err)
		}
		plan, err := parser.ParsePlan(strings.NewReader(string(data)))
		if err != nil {
			return fmt.Errorf("env %s: %w", name, err)
		}
		params, err := merger.Merge(plan, inputs, merger.Scope(f.scope), f.module)
		if err != nil {
			return fmt.Errorf("env %s: %w", name, err)
		}
		sortParams(params, f.sortBy)
		if f.scope == "module" && moduleName == "" {
			if moduleName, err = merger.ModuleName(plan, f.module); err != nil {
				return err
			}
		}
		envResults = append(envResults, merger.EnvResult{Name: name, Params: params})
		envNames = append(envNames, name)
	}

	rows := merger.Compare(envResults)
	content := formatter.CompareMarkdown(rows, formatter.CompareOptions{
		EnvNames:      envNames,
		GeneratedAt:   time.Now().Format("2006-01-02 15:04:05 MST"),
		ShowSensitive: f.showSensitive,
		HighlightDiff: f.highlightDiff,
		Scope:         f.scope,
		ModuleName:    moduleName,
	})
	return writeOutput(cmd, f.out, content)
}
