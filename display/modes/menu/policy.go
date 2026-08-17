package menu

import (
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/menu/source"
)

type Policy = source.Policy

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// normalizePolicy ensures the policy fields contain valid values, falling back to defaults.
func normalizePolicy(p source.Policy) source.Policy {
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && menuRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}

// menuSnapshotter implements catalog.PolicySnapshotter for the menu mode.
type menuSnapshotter struct{}

// SnapshotPolicy returns the current menu policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (menuSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"style": p.Style,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (menuSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}

	SetPolicy(p)
	return nil
}

func init() {
	catalog.RegisterSnapshotter("menu", menuSnapshotter{})
}
