// Package cmd wires the tfparams command-line interface.
package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/tfparam/tfparams/internal/config"
	"github.com/tfparam/tfparams/internal/formatter"
	"github.com/tfparam/tfparams/internal/merger"
	"github.com/tfparam/tfparams/internal/parser"
)

var version = "dev"

// SetVersion injects the build version (set from main).
func SetVersion(v string) { version = v }

type rootFlags struct {
	planJSON      string
	docsJSON      []string
	scope         string
	module        string
	out           string
	format        string
	env           string
	showSensitive bool
	noDefaultCol  bool
	sortBy        string
	recursive     bool
	recursivePath string
	configPath    string
}

// NewRootCmd builds the root cobra command.
func NewRootCmd() *cobra.Command {
	f := &rootFlags{}
	cmd := &cobra.Command{
		Use:   "tfparams",
		Short: "Generate Terraform parameter sheets from a plan JSON and terraform-docs metadata",
		Long: "tfparams merges the applied values from a Terraform plan " +
			"(terraform show -json <planfile>) with variable metadata from terraform-docs " +
			"and renders a Markdown parameter sheet.",
		RunE:          func(cmd *cobra.Command, _ []string) error { return run(cmd, f) },
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
	}
	fl := cmd.Flags()
	fl.StringVar(&f.planJSON, "plan-json", "", "plan JSON file (terraform show -json <planfile>); default stdin")
	fl.StringArrayVar(&f.docsJSON, "docs-json", nil, "terraform-docs JSON, e.g. <(terraform-docs json .) (required, repeatable)")
	fl.StringVar(&f.scope, "scope", "root", "scope: root or module")
	fl.StringVar(&f.module, "module", "", "module call name for --scope module")
	fl.StringVar(&f.out, "out", "", "output file (overwritten); default stdout")
	fl.StringVar(&f.format, "format", "table", "output format: table, csv, or json")
	fl.StringVar(&f.env, "env", "", "environment name shown in the header")
	fl.BoolVar(&f.showSensitive, "show-sensitive", false, "show sensitive values unmasked")
	fl.BoolVar(&f.noDefaultCol, "no-default-col", false, "hide the Default column")
	fl.StringVar(&f.sortBy, "sort-by", "name", "sort key: name, required, or type")
	fl.BoolVar(&f.recursive, "recursive", false, "process subdirectories recursively")
	fl.StringVar(&f.recursivePath, "recursive-path", ".", "root to scan in recursive mode")
	fl.StringVar(&f.configPath, "config", "", "config file path (default: search .tfparams.yml)")
	return cmd
}

// Execute runs the root command.
func Execute() error { return NewRootCmd().Execute() }

func run(cmd *cobra.Command, f *rootFlags) error {
	changed := func(name string) bool { return cmd.Flags().Changed(name) }

	explicitConfig := ""
	if changed("config") {
		explicitConfig = f.configPath
	}
	cfg, _, err := config.Load(explicitConfig)
	if err != nil {
		return err
	}

	env := pick(changed("env"), f.env, cfg.Env)
	scope := pick(changed("scope"), f.scope, cfg.Scope)
	if scope == "" {
		scope = "root"
	}
	module := pick(changed("module"), f.module, cfg.Module)
	format := pick(changed("format"), f.format, cfg.Format)
	if format == "" {
		format = "table"
	}
	out := pick(changed("out"), f.out, cfg.Output.File)
	sortBy := pick(changed("sort-by"), f.sortBy, cfg.Sort.By)
	showSensitive := cfg.Sensitive.Show
	if changed("show-sensitive") {
		showSensitive = f.showSensitive
	}
	sortEnabled := cfg.Sort.Enabled || changed("sort-by")

	if f.recursive || changed("recursive") {
		return fmt.Errorf("recursive mode is not implemented yet")
	}
	if len(f.docsJSON) == 0 {
		return fmt.Errorf("--docs-json is required\n\n%s", cmd.UsageString())
	}

	cols := cfg.Columns.Show
	if len(cols) == 0 {
		cols = formatter.DefaultColumns
	}
	if f.noDefaultCol {
		cols = removeString(cols, "default")
	}

	plan, err := loadPlan(f)
	if err != nil {
		return err
	}

	var docsList []*parser.Docs
	for _, path := range f.docsJSON {
		d, derr := readDocs(path)
		if derr != nil {
			return derr
		}
		docsList = append(docsList, d)
	}
	inputs := merger.MergeInputs(docsList...)

	params, err := merger.Merge(plan, inputs, merger.Scope(scope), module)
	if err != nil {
		return err
	}
	if sortEnabled {
		sortParams(params, sortBy)
	}

	moduleName := ""
	if scope == "module" {
		if moduleName, err = merger.ModuleName(plan, module); err != nil {
			return err
		}
	}

	opts := formatter.Options{
		Env:           env,
		Scope:         scope,
		ModuleName:    moduleName,
		GeneratedAt:   time.Now().Format("2006-01-02 15:04:05 MST"),
		Source:        "terraform show -json tfplan (plan)",
		ShowSensitive: showSensitive,
		Columns:       cols,
	}

	var content string
	switch format {
	case "table":
		content = formatter.Markdown(params, opts)
	case "csv":
		c, cerr := formatter.CSV(params, opts)
		if cerr != nil {
			return cerr
		}
		content = c
	case "json":
		c, cerr := formatter.JSON(params, opts)
		if cerr != nil {
			return cerr
		}
		content = c
	default:
		return fmt.Errorf("unknown format %q (want table, csv, or json)", format)
	}

	return writeOutput(cmd, out, content)
}

// pick returns flagVal when the flag was changed, otherwise cfgVal.
func pick(changed bool, flagVal, cfgVal string) string {
	if changed {
		return flagVal
	}
	return cfgVal
}

func loadPlan(f *rootFlags) (*parser.Plan, error) {
	if f.planJSON == "" {
		return parser.ParsePlan(os.Stdin)
	}
	file, err := os.Open(f.planJSON) //nolint:gosec // path is user-provided input
	if err != nil {
		return nil, fmt.Errorf("open plan json: %w", err)
	}
	defer func() { _ = file.Close() }()
	return parser.ParsePlan(file)
}

func readDocs(path string) (*parser.Docs, error) {
	file, err := os.Open(path) //nolint:gosec // path is user-provided input
	if err != nil {
		return nil, fmt.Errorf("open docs json %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	return parser.ParseDocs(file)
}

// writeOutput writes to stdout when out is empty, otherwise overwrites the file.
func writeOutput(cmd *cobra.Command, out, content string) error {
	if out == "" {
		_, err := io.WriteString(cmd.OutOrStdout(), content)
		return err
	}
	return os.WriteFile(out, []byte(content), 0o644) //nolint:gosec // sheet is meant to be world-readable
}

func sortParams(params []merger.Param, by string) {
	switch by {
	case "required":
		sort.SliceStable(params, func(i, j int) bool {
			if params[i].Required != params[j].Required {
				return params[i].Required
			}
			return params[i].Name < params[j].Name
		})
	case "type":
		sort.SliceStable(params, func(i, j int) bool {
			if params[i].Type != params[j].Type {
				return params[i].Type < params[j].Type
			}
			return params[i].Name < params[j].Name
		})
	default:
		sort.SliceStable(params, func(i, j int) bool { return params[i].Name < params[j].Name })
	}
}

func removeString(s []string, target string) []string {
	out := make([]string, 0, len(s))
	for _, v := range s {
		if v != target {
			out = append(out, v)
		}
	}
	return out
}
