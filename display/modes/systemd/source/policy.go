package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// AllowedAccents lists the valid color accent values for boot status theming.
var AllowedAccents = []string{"cyan", "green", "amber", "red", "white", "none"}

// Policy captures all runtime-configurable parameters for the systemd mode.
type Policy struct {
	Style       string // Resolution-specific style name (e.g., "color-240x240", "mono-128x64")
	ColorAccent string // Boot-in-progress color accent: "cyan", "green", "amber", "red", "white", "none"
}

// DefaultPolicy returns the default systemd policy.
func DefaultPolicy() Policy {
	return Policy{
		Style:       "",
		ColorAccent: "amber",
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: nil},
		{Key: "color_accent", Type: "string", Summary: "Boot progress accent color.", Default: "amber", Allowed: AllowedAccents},
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s", p.Style, p.ColorAccent)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":        p.Style,
		"color_accent": p.ColorAccent,
	}
}
