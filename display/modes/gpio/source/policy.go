package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// AllowedFGColors lists all valid fgcolor policy values.
var AllowedFGColors = []string{"cyan", "green", "amber", "red", "white", "none"}

// Policy captures all runtime-configurable parameters for the GPIO mode.
type Policy struct {
	Style     string         // registry style name (legacy or resolution-specific)
	Color     bool           // whether to emit per-row Colors
	Font      string         // "auto" or a registered font ID
	PinLabels map[int]string // BCM pin number -> user label (for "detail" style)
	FGColor   string         // "cyan", "green", "amber", "red", "white", "none"
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "style", Type: "string", Summary: "Visual presentation style."},
		{Key: "color", Type: "bool", Summary: "Whether to color-code pin rows by level.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
		{Key: "fgcolor", Type: "string", Summary: "Foreground color for GPIO display elements.", Default: "cyan", Allowed: AllowedFGColors},
	}
}

// DefaultPolicy returns the default GPIO policy.
func DefaultPolicy() Policy {
	return Policy{Style: "", Color: true, Font: "auto", FGColor: "cyan"}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%v|%s|%s|%d", p.Style, p.Color, p.Font, p.FGColor, len(p.PinLabels))
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":   p.Style,
		"color":   p.Color,
		"font":    p.Font,
		"fgcolor": p.FGColor,
	}
}
