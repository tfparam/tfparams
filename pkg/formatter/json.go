package formatter

import (
	"encoding/json"

	"github.com/tfkit/tfparams/pkg/merger"
)

type jsonDoc struct {
	Environment string      `json:"environment,omitempty"`
	Scope       string      `json:"scope,omitempty"`
	Module      string      `json:"module,omitempty"`
	GeneratedAt string      `json:"generated_at,omitempty"`
	Source      string      `json:"source,omitempty"`
	Variables   []jsonParam `json:"variables"`
}

type jsonParam struct {
	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	Type         string  `json:"type,omitempty"`
	Default      *string `json:"default"`
	AppliedValue *string `json:"applied_value"`
	Computed     bool    `json:"computed,omitempty"`
	Required     bool    `json:"required"`
	Sensitive    bool    `json:"sensitive"`
}

// JSON renders the parameters as an indented JSON document.
func JSON(params []merger.Param, o Options) (string, error) {
	doc := jsonDoc{
		Environment: o.Env,
		Scope:       o.Scope,
		Module:      o.ModuleName,
		GeneratedAt: o.GeneratedAt,
		Source:      o.Source,
		Variables:   make([]jsonParam, 0, len(params)),
	}
	for _, p := range params {
		jp := jsonParam{
			Name:        p.Name,
			Description: p.Description,
			Type:        p.Type,
			Computed:    p.Computed,
			Required:    p.Required,
			Sensitive:   p.Sensitive,
		}
		if p.HasDefault {
			d := p.Default
			jp.Default = &d
		}
		if text, isValue := appliedString(p, o.ShowSensitive); isValue {
			v := text
			jp.AppliedValue = &v
		} else if p.Sensitive && !o.ShowSensitive {
			// masked value is still meaningful in JSON
			v := text
			jp.AppliedValue = &v
		}
		// not-set / computed leave AppliedValue as null
		doc.Variables = append(doc.Variables, jp)
	}

	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}
