package source

import "github.com/databeast/cyberhud/display/catalog"

// Policy captures all runtime-configurable parameters for the Stemma mode.
type Policy struct {
	Style string // registered style name; empty uses automatic resolution
}

func (p Policy) Options() []catalog.OptionDefinition {
	return nil
}

// DefaultPolicy returns the default Stemma policy.
func DefaultPolicy() Policy {
	return Policy{
		Style: "",
	}
}

func (p Policy) Fingerprint() string {
	return p.Style
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style": p.Style,
	}
}
