package source

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the system mode.
type Policy struct {
	Style string // "default", "compact", "cores", or "top"
	Font  string // "auto" or a specific font ID
}

// DefaultPolicy returns the default system policy.
func DefaultPolicy() Policy {
	return Policy{
		Style: "",
		Font:  "auto",
	}
}

// Options returns catalog option definitions for the system mode policy.
func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: ""},
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
	}
}

// Fingerprint encodes all policy fields into a deterministic string.
func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s", p.Style, p.Font)
}

// ToMap returns the policy as a JSON-serializable map.
func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style": p.Style,
		"font":  p.Font,
	}
}

// NormalizePolicy performs registry-independent policy normalization.
func NormalizePolicy(p Policy) Policy {
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	p.Font = strings.TrimSpace(p.Font)
	if p.Font == "" {
		p.Font = "auto"
	}
	return p
}
