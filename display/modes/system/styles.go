package system

import (
	"fmt"
	"regexp"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
)

// styleNamePattern is the regex that all registered system style names must
// match. Format: <category_prefix><WxH>[-<variant>], where variant is one of
// default|compact|cores|top. Mirrors the thermal/clock naming convention.
var styleNamePattern = regexp.MustCompile(`^(color-|color-slow-|mono-|mono-slow-|grayscale-fast-|grayscale-slow-)\d+x\d+(-(default|compact|cores|top))?$`)

func init() {
	for _, s := range systemRegistry.Enumerate() {
		if !styleNamePattern.MatchString(s.Name()) {
			panic(fmt.Sprintf("system: invalid style name %q does not match expected pattern %s", s.Name(), styleNamePattern.String()))
		}
	}
}

// Allowed style variant tokens for the system mode. The full style name is
// <capability-prefix>-<WxH>[-<variant>]; these are the variant suffixes.
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
