package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the menu mode.
type Policy struct {
	Style string // registry style name
}

// DefaultPolicy returns the default menu policy.
func DefaultPolicy() Policy {
	return Policy{Style: ""}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return nil
}

func (p Policy) Fingerprint() string {
	return p.Style
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style": p.Style,
	}
}

// String returns a human-readable representation of the policy.
func (p Policy) String() string {
	return fmt.Sprintf("style=%s", p.Style)
}
