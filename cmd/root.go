// Package cmd wires the tfparams command-line interface.
package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/tfkit/tfparams/pkg/config"
	"github.com/tfkit/tfparams/pkg/formatter"
	"github.com/tfkit/tfparams/pkg/merger"
	"github.com/tfkit/tfparams/pkg/parser"
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
	fl.StringVar(&f.format, "format", "markdown", "output format: markdown, csv, or json")
	fl.StringVar(&f.env, "env", "", "environment name shown in the header")
	fl.BoolVar(&f.showSensitive, "show-sensitive", false, "show sensitive values unmasked")
	fl.BoolVar(&f.noDefaultCol, "no-default-col", false, "hide the Default column")
	fl.StringVar(&f.sortBy, "sort-by", "required", "sort key: required (required first, then name) or name")
	fl.BoolVar(&f.recursive, "recursive", false, "process subdirectories recursively")
	fl.StringVar(&f.recursivePath, "recursive-path", ".", "root to scan in recursive mode")
	fl.StringVar(&f.configPath, "config", "", "config file path (default: search .tfparams.yml)")

	cmd.AddCommand(newCompareCmd())
	return cmd
}

// Execute runs the root command.
func Execute() error { return NewRootCmd().Execute() }

// settings holds the effective configuration after merging file config and flags.
type settings struct {
	env, scope, module, format, out, sortBy, source string
	showSensitive, recursive                        bool
	recursivePath, planFile                         string
	cols                                            []string
	docs                                            []*parser.Docs
}

func run(cmd *cobra.Command, f *rootFlags) error {
	cfg, err := loadRootConfig(cmd, f)
	if err != nil {
		return err
	}
	if len(f.docsJSON) == 0 {
		return fmt.Errorf("--docs-json is required\n\n%s", cmd.UsageString())
	}
	docs, err := loadDocs(f)
	if err != nil {
		return err
	}

	s := resolveSettings(cmd, f, cfg)
	s.docs = docs

	if s.recursive {
		return runRecursive(cmd, f, cfg, s)
	}
	return runSingle(cmd, f, s)
}

func loadRootConfig(cmd *cobra.Command, f *rootFlags) (config.Config, error) {
	explicit := ""
	if cmd.Flags().Changed("config") {
		explicit = f.configPath
	}
	cfg, _, err := config.Load(explicit)
	return cfg, err
}

func loadDocs(f *rootFlags) ([]*parser.Docs, error) {
	var docs []*parser.Docs
	for _, path := range f.docsJSON {
		d, err := readDocs(path)
		if err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

// resolveSettings computes effective settings from cfg and CLI flags (flags win).
func resolveSettings(cmd *cobra.Command, f *rootFlags, cfg config.Config) settings {
	changed := func(name string) bool { return cmd.Flags().Changed(name) }
	s := settings{
		env:           pick(changed("env"), f.env, cfg.Env),
		scope:         pick(changed("scope"), f.scope, cfg.Scope),
		module:        pick(changed("module"), f.module, cfg.Module),
		format:        pick(changed("format"), f.format, cfg.Format),
		out:           pick(changed("out"), f.out, cfg.Output.File),
		sortBy:        pick(changed("sort-by"), f.sortBy, cfg.Sort.By),
		source:        "terraform show -json tfplan (plan)",
		recursive:     f.recursive || changed("recursive") || cfg.Recursive.Enabled,
		recursivePath: pick(changed("recursive-path"), f.recursivePath, cfg.Recursive.Path),
		planFile:      cfg.Recursive.PlanFile,
	}
	if s.scope == "" {
		s.scope = "root"
	}
	if s.format == "" {
		s.format = "markdown"
	}
	if s.sortBy == "" {
		s.sortBy = "required"
	}
	if s.recursivePath == "" {
		s.recursivePath = "."
	}
	if s.planFile == "" {
		s.planFile = "tfplan.json"
	}
	s.showSensitive = cfg.Sensitive.Show
	if changed("show-sensitive") {
		s.showSensitive = f.showSensitive
	}

	cols := cfg.Columns.Show
	if len(cols) == 0 {
		cols = formatter.DefaultColumns
	}
	if f.noDefaultCol {
		cols = removeString(cols, "default")
	}
	s.cols = cols
	return s
}

func buildContent(plan *parser.Plan, s settings) (string, error) {
	inputs := merger.MergeInputs(s.docs...)
	params, err := merger.Merge(plan, inputs, merger.Scope(s.scope), s.module)
	if err != nil {
		return "", err
	}
	sortParams(params, s.sortBy)
	moduleName := ""
	if s.scope == "module" {
		if moduleName, err = merger.ModuleName(plan, s.module); err != nil {
			return "", err
		}
	}
	opts := formatter.Options{
		Env:           s.env,
		Scope:         s.scope,
		ModuleName:    moduleName,
		GeneratedAt:   time.Now().Format("2006-01-02 15:04:05 MST"),
		Source:        s.source,
		ShowSensitive: s.showSensitive,
		Columns:       s.cols,
	}
	switch s.format {
	case "markdown", "table": // "table" kept as a backward-compatible alias
		return formatter.Markdown(params, opts), nil
	case "csv":
		return formatter.CSV(params, opts)
	case "json":
		return formatter.JSON(params, opts)
	default:
		return "", fmt.Errorf("unknown format %q (want markdown, csv, or json)", s.format)
	}
}

func runSingle(cmd *cobra.Command, f *rootFlags, s settings) error {
	plan, err := loadPlan(f)
	if err != nil {
		return err
	}
	content, err := buildContent(plan, s)
	if err != nil {
		return err
	}
	return writeOutput(cmd, s.out, content)
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
	return writeToFile(out, content)
}

func writeToFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644) //nolint:gosec // sheet is meant to be world-readable
}

// sortParams orders rows by the chosen key. "required" (default) lists required
// variables first, then alphabetically; "name" is plain alphabetical.
func sortParams(params []merger.Param, by string) {
	if by == "name" {
		sort.SliceStable(params, func(i, j int) bool { return params[i].Name < params[j].Name })
		return
	}
	sort.SliceStable(params, func(i, j int) bool {
		if params[i].Required != params[j].Required {
			return params[i].Required
		}
		return params[i].Name < params[j].Name
	})
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
