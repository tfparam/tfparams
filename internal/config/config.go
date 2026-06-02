// Package config loads .tfparams.yml and exposes its settings. CLI flags
// override file values; see the Resolve helper used by cmd.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config mirrors the .tfparams.yml schema.
type Config struct {
	Format    string    `yaml:"format"`
	Env       string    `yaml:"env"`
	Scope     string    `yaml:"scope"`
	Module    string    `yaml:"module"`
	Output    Output    `yaml:"output"`
	Columns   Columns   `yaml:"columns"`
	Sort      Sort      `yaml:"sort"`
	Sensitive Sensitive `yaml:"sensitive"`
	Recursive Recursive `yaml:"recursive"`
}

// Output configures where and how the sheet is written.
type Output struct {
	File     string `yaml:"file"`
	Mode     string `yaml:"mode"`
	Template string `yaml:"template"`
}

// Columns configures which columns appear.
type Columns struct {
	Show []string `yaml:"show"`
}

// Sort configures row ordering.
type Sort struct {
	Enabled bool   `yaml:"enabled"`
	By      string `yaml:"by"`
}

// Sensitive configures sensitive value handling.
type Sensitive struct {
	Show bool   `yaml:"show"`
	Mask string `yaml:"mask"`
}

// Recursive configures recursive mode.
type Recursive struct {
	Enabled  bool   `yaml:"enabled"`
	Path     string `yaml:"path"`
	PlanFile string `yaml:"plan_file"`
}

// Default returns a Config populated with tfparams' built-in defaults.
func Default() Config {
	return Config{
		Format: "table",
		Scope:  "root",
		Output: Output{Mode: "standalone"},
		Columns: Columns{Show: []string{
			"name", "description", "type", "default", "applied_value", "required",
		}},
		Sort:      Sort{Enabled: false, By: "name"},
		Sensitive: Sensitive{Show: false, Mask: "(sensitive)"},
		Recursive: Recursive{Enabled: false, Path: ".", PlanFile: "tfplan.json"},
	}
}

// SearchPaths is the ordered list of locations probed when no explicit
// --config path is given.
func SearchPaths() []string {
	paths := []string{
		".tfparams.yml",
		filepath.Join(".config", ".tfparams.yml"),
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".tfparams.d", ".tfparams.yml"))
	}
	return paths
}

// Load reads the config file. When explicit is non-empty it must exist;
// otherwise the SearchPaths are tried in order and the first existing file is
// used. If no file is found, defaults are returned with found=false.
func Load(explicit string) (cfg Config, found bool, err error) {
	cfg = Default()

	if explicit != "" {
		if err := loadInto(explicit, &cfg); err != nil {
			return cfg, false, err
		}
		return cfg, true, nil
	}

	for _, p := range SearchPaths() {
		if _, statErr := os.Stat(p); statErr == nil {
			if err := loadInto(p, &cfg); err != nil {
				return cfg, false, err
			}
			return cfg, true, nil
		}
	}
	return cfg, false, nil
}

func loadInto(path string, cfg *Config) error {
	data, err := os.ReadFile(path) //nolint:gosec // path is user-provided config
	if err != nil {
		return fmt.Errorf("read config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}
	return nil
}
