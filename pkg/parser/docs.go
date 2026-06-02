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
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	// Default is the raw JSON default value (e.g. "t3.medium", 1, false, null).
	// terraform-docs emits the value directly, not wrapped in { "value": ... }.
	Default  json.RawMessage `json:"default"`
	Required bool            `json:"required"`
	// Sensitive is set when terraform-docs reports it. Not all versions emit it
	// (it is often null); sensitivity is also derived from the plan configuration.
	Sensitive bool `json:"sensitive"`
}

// HasDefault reports whether terraform-docs provided a non-null default.
func (in Input) HasDefault() bool {
	return len(in.Default) > 0 && string(in.Default) != "null"
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
