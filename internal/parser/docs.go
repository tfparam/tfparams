package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// Docs is the subset of `terraform-docs json` output that tfparams reads.
type Docs struct {
	Inputs []Input `json:"inputs"`
}

// Input is the metadata for one variable as reported by terraform-docs.
type Input struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     *DocDefault `json:"default"`
	Required    bool        `json:"required"`
	// Sensitive is set when terraform-docs reports it. Not all versions emit it.
	Sensitive bool `json:"sensitive"`
}

// DocDefault wraps the default value object: { "value": ... }.
type DocDefault struct {
	Value json.RawMessage `json:"value"`
}

// ParseDocs decodes a terraform-docs JSON document from r.
func ParseDocs(r io.Reader) (*Docs, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read docs json: %w", err)
	}
	var d Docs
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("parse docs json: %w", err)
	}
	return &d, nil
}
