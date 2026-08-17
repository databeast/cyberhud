package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy holds all runtime-configurable parameters for the GPIO Control mode.
type Policy struct {
	Style string // "list", "compact", or "grid"
	Font  string // "auto" or a registered font ID
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
	}
}

// DefaultPolicy returns the default GPIO Control mode policy.
func DefaultPolicy() Policy {
	return Policy{
		Style: "",
		Font:  "auto",
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s", p.Style, p.Font)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style": p.Style,
		"font":  p.Font,
	}
}
