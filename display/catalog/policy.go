package catalog

type ConfigPolicy interface {
	Fingerprint() string
	ToMap() map[string]interface{}
	Options() []OptionDefinition
}

// DefaultModePolicy resolves the default mode for a panel given
// the available modes and whether input hardware is present.
type DefaultModePolicy interface {
	ResolveDefault(available []string, inputEnabled bool) string
}

// StandardPolicy implements the current UX rule:
// input-enabled → "menu", else → "dashboard", else → first available.
type StandardPolicy struct{}

func (StandardPolicy) ResolveDefault(available []string, inputEnabled bool) string {
	preferred := "dashboard"
	if inputEnabled {
		preferred = "menu"
	}

	for _, m := range available {
		if m == preferred {
			return m
		}
	}
	if len(available) > 0 {
		return available[0]
	}
	return preferred
}
