package stemma

import (
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/modes/stemma/styles"
	"github.com/databeast/cyberhud/display/style"
)

// stemmaRegistry is the per-mode StyleRegistry for the stemma display mode.
var stemmaRegistry = style.NewRegistry[source.StemmaSnapshot, source.Policy](
	styles.AllStyles()...,
)

// registryStyleNames returns the ordered list of style names from the stemmaRegistry.
// Used for cmdHandler validation and catalog registration.
func registryStyleNames() []string {
	styles := stemmaRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
