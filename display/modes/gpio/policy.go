package gpio

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/gpio/source"
)

// gpioSnapshotter implements catalog.PolicySnapshotter for the gpio mode.
type gpioSnapshotter struct{}

// SnapshotPolicy returns the current gpio policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (gpioSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"style":   p.Style,
		"color":   p.Color,
		"font":    p.Font,
		"fgcolor": p.FGColor,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (gpioSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}
	if v, ok := data["color"]; ok {
		if b, ok := v.(bool); ok {
			p.Color = b
		}
	}
	if v, ok := data["font"]; ok {
		if s, ok := v.(string); ok {
			p.Font = s
		}
	}
	if v, ok := data["fgcolor"]; ok {
		if s, ok := v.(string); ok {
			p.FGColor = s
		}
	}

	SetPolicy(p)
	return nil
}

// normalizePolicy ensures policy fields contain valid values.
func normalizePolicy(p Policy) Policy {
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && gpioRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	p.FGColor = strings.ToLower(strings.TrimSpace(p.FGColor))
	if !isAllowedAccent(p.FGColor) {
		p.FGColor = "cyan"
	}
	p.Font = strings.TrimSpace(p.Font)
	if p.Font == "" {
		p.Font = "auto"
	}
	return p
}

func init() {
	catalog.RegisterSnapshotter("gpio", gpioSnapshotter{})
}
