package source

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// LineDirective describes one ticker line and its rendering behavior.
type LineDirective struct {
	Text     string `json:"text"`
	Font     string `json:"font,omitempty"`
	LineMode string `json:"line_mode,omitempty"`
	Scaling  string `json:"scaling,omitempty"`
	Scroll   string `json:"scroll,omitempty"`
}

// FormattedLine is the per-line render output for ViewData.
type FormattedLine struct {
	Text        string
	Tier        tiercatalog.Tier
	FontWarning string
}

// ParseJSONFeed parses a JSON payload into validated LineDirectives.
// Returns an error with line index context on validation failure.
func ParseJSONFeed(payload string) ([]LineDirective, error) {
	// Step 1: Check top-level is array.
	payload = strings.TrimSpace(payload)
	if len(payload) == 0 || payload[0] != '[' {
		ch := ""
		if len(payload) > 0 {
			ch = payload[:1]
		}
		return nil, fmt.Errorf("expected JSON array, got %q", ch)
	}

	// Step 2: Unmarshal into raw messages for per-element validation.
	var raw []json.RawMessage
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	// Step 3: Parse and validate each element.
	directives := make([]LineDirective, 0, len(raw))
	for i, elem := range raw {
		var d LineDirective
		if err := json.Unmarshal(elem, &d); err != nil {
			return nil, fmt.Errorf("line %d: %v", i, err)
		}
		// Validate required text field.
		if strings.TrimSpace(d.Text) == "" {
			var m map[string]interface{}
			json.Unmarshal(elem, &m)
			if _, hasText := m["text"]; !hasText {
				return nil, fmt.Errorf("line %d: missing required field \"text\"", i)
			}
			return nil, fmt.Errorf("line %d: \"text\" must be non-empty", i)
		}
		d.Text = strings.TrimSpace(d.Text)
		directives = append(directives, d)
	}
	return directives, nil
}
