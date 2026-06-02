package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tfparam/tfparams/internal/config"
	"github.com/tfparam/tfparams/internal/parser"
)

// runRecursive processes every immediate subdirectory of recursivePath that
// contains the configured plan-JSON file. Each subdirectory may override the
// root config via its own .tfparams.yml.
func runRecursive(cmd *cobra.Command, f *rootFlags, rootCfg config.Config, root settings) error {
	log := slog.New(slog.NewTextHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{Level: slog.LevelInfo}))

	entries, err := os.ReadDir(root.recursivePath)
	if err != nil {
		return fmt.Errorf("scan %s: %w", root.recursivePath, err)
	}

	processed := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(root.recursivePath, e.Name())
		planPath := filepath.Join(dir, root.planFile)
		if !fileExists(planPath) {
			log.Warn("skipping subdirectory: no plan file", "dir", dir, "plan_file", root.planFile)
			continue
		}

		sub := root
		if subCfgPath := filepath.Join(dir, ".tfparams.yml"); fileExists(subCfgPath) {
			merged, oerr := config.Overlay(rootCfg, subCfgPath)
			if oerr != nil {
				return oerr
			}
			sub = resolveSettings(cmd, f, merged)
			sub.docs = root.docs
		}

		plan, perr := readPlanFile(planPath)
		if perr != nil {
			return fmt.Errorf("%s: %w", planPath, perr)
		}
		content, berr := buildContent(plan, sub)
		if berr != nil {
			return fmt.Errorf("%s: %w", dir, berr)
		}

		if sub.out == "" {
			// Program output (the sheets) goes to stdout, not the diagnostic logger.
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "===== %s =====\n%s\n", dir, content)
		} else {
			target := filepath.Join(dir, sub.out)
			if werr := writeToFile(target, content); werr != nil {
				return werr
			}
			log.Info("wrote parameter sheet", "path", target)
		}
		processed++
	}

	if processed == 0 {
		log.Warn("no subdirectories with a plan file", "plan_file", root.planFile, "root", root.recursivePath)
	}
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func readPlanFile(path string) (*parser.Plan, error) {
	file, err := os.Open(path) //nolint:gosec // path is derived from the scanned directory
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	return parser.ParsePlan(file)
}
