package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy holds runtime-configurable parameters for the dashboard mode.
type Policy struct {
	Style       string // Resolution-specific style name or legacy name
	ColorAccent string // "cyan", "green", "amber", "red", "white", "none"
	LoPower     bool   // disable power-consuming fluff
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%v|%t",
		p.Style, p.ColorAccent, p.LoPower)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":        p.Style,
		"color_accent": p.ColorAccent,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{
			Key:     "color_accent",
			Type:    "string",
			Summary: "Uptime heartbeat accent color.",
			Default: "cyan",
			Allowed: AllowedAccents,
		},
	}
}

// allowedAccents is the list of valid ColorAccent values.
var AllowedAccents = []string{"cyan", "green", "amber", "red", "white", "none"}

// DefaultPolicy returns the default dashboard policy.
func DefaultPolicy() Policy {
	return Policy{
		Style:       "",
		ColorAccent: "cyan",
	}
}
