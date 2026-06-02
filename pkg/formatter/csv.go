package formatter

import (
	"encoding/csv"
	"strings"

	"github.com/tfparam/tfparams/pkg/merger"
)

// CSV renders the parameters as CSV using the configured column order.
func CSV(params []merger.Param, o Options) (string, error) {
	cols := o.Columns
	if cols == nil {
		cols = DefaultColumns
	}

	var b strings.Builder
	w := csv.NewWriter(&b)

	header := make([]string, len(cols))
	for i, c := range cols {
		header[i] = columnHeaders[c]
	}
	if err := w.Write(header); err != nil {
		return "", err
	}

	for _, p := range params {
		rec := make([]string, len(cols))
		for i, c := range cols {
			rec[i] = csvCell(c, p, o)
		}
		if err := w.Write(rec); err != nil {
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return b.String(), nil
}

func csvCell(key string, p merger.Param, o Options) string {
	switch key {
	case "name":
		return p.Name
	case "description":
		return p.Description
	case "type":
		return p.Type
	case "default":
		if !p.HasDefault {
			return ""
		}
		return p.Default
	case "applied_value":
		text, _ := appliedString(p, o.ShowSensitive)
		return text
	case "required":
		if p.Required {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
