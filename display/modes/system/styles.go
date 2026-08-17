package system

import (
	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/modes/system/styles"
	"github.com/databeast/cyberhud/display/style"
)

// systemRegistry is the per-mode StyleRegistry for the system display mode.
var systemRegistry = style.NewRegistry[source.SystemSnapshot, source.Policy](styles.DefaultStyle, styles.CompactStyle, styles.CoresStyle, styles.TopStyle)

// Allowed style values for the system mode.
const (
	StyleDefault = "default"
	StyleCompact = "compact"
	StyleCores   = "cores"
	StyleTop     = "top"
)

// registeredStyleNames returns the list of style names from the systemRegistry.
// Used by catalog registration and HandleCommand for allowed-value validation.
func registeredStyleNames() []string {
	styles := systemRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}

// SystemRegistryEnumerate exposes registered system styles for snapshot tests.
func SystemRegistryEnumerate() []style.Style[source.SystemSnapshot, source.Policy] {
	return systemRegistry.Enumerate()
}
