package storage

import (
	"fmt"
	"regexp"

	"github.com/databeast/cyberhud/display/modes/storage/source"
	"github.com/databeast/cyberhud/display/modes/storage/styles"
	"github.com/databeast/cyberhud/display/style"
)

// styleNamePattern is the regex that all registered storage style names must
// match. Format: <category_prefix><WxH>[-<variant>] where variant is bars|focus.
// Mirrors the thermal/system/clock naming convention.
var styleNamePattern = regexp.MustCompile(`^(mono-fast-|mono-slow-|grayscale-fast-|grayscale-slow-|color-fast-|color-slow-)\d+x\d+(-(bars|focus))?$`)

// storageRegistry is the StyleRegistry for the storage display mode. Styles are
// registered grouped by capability class; registration order determines fitness
// tie-breaking (first registered wins), and the skeleton (adaptive) style is
// registered first within each resolution group.
var storageRegistry = style.NewRegistry[source.StorageSnapshot, source.Policy](styles.All()...)

func init() {
	for _, s := range storageRegistry.Enumerate() {
		if !styleNamePattern.MatchString(s.Name()) {
			panic(fmt.Sprintf("storage: invalid style name %q does not match expected pattern %s", s.Name(), styleNamePattern.String()))
		}
	}
}

// allowedStyleNames returns the list of registered style names, used by catalog
// registration and the command handler for allowed-value validation.
func allowedStyleNames() []string {
	enum := storageRegistry.Enumerate()
	names := make([]string, len(enum))
	for i, s := range enum {
		names[i] = s.Name()
	}
	return names
}

// StorageRegistry exposes the registry for tests and tooling.
var StorageRegistry = storageRegistry
