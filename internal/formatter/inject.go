package formatter

import (
	"errors"
	"strings"
)

// BeginMarker and EndMarker delimit the tfparams-managed region in inject mode.
const (
	BeginMarker = "<!-- BEGIN_TFPARAMS -->"
	EndMarker   = "<!-- END_TFPARAMS -->"
)

// ErrUnbalancedMarkers is returned when exactly one of the markers is present.
var ErrUnbalancedMarkers = errors.New("inject: found only one of BEGIN_TFPARAMS/END_TFPARAMS markers")

// Inject replaces the content between the markers in existing with content.
// If neither marker is present, a marker block is appended to the end. If only
// one marker is present, ErrUnbalancedMarkers is returned.
func Inject(existing, content string) (string, error) {
	begin := strings.Index(existing, BeginMarker)
	end := strings.Index(existing, EndMarker)

	switch {
	case begin == -1 && end == -1:
		return appendBlock(existing, content), nil
	case begin == -1 || end == -1:
		return "", ErrUnbalancedMarkers
	case end < begin:
		return "", ErrUnbalancedMarkers
	}

	body := strings.TrimRight(content, "\n")
	var b strings.Builder
	b.WriteString(existing[:begin])
	b.WriteString(BeginMarker)
	b.WriteString("\n")
	b.WriteString(body)
	b.WriteString("\n")
	b.WriteString(EndMarker)
	b.WriteString(existing[end+len(EndMarker):])
	return b.String(), nil
}

func appendBlock(existing, content string) string {
	body := strings.TrimRight(content, "\n")
	var b strings.Builder
	b.WriteString(existing)
	if existing != "" && !strings.HasSuffix(existing, "\n") {
		b.WriteString("\n")
	}
	if existing != "" {
		b.WriteString("\n")
	}
	b.WriteString(BeginMarker)
	b.WriteString("\n")
	b.WriteString(body)
	b.WriteString("\n")
	b.WriteString(EndMarker)
	b.WriteString("\n")
	return b.String()
}
