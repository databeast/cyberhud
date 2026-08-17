package menu

import (
	"github.com/databeast/cyberhud/display/modes/menu/source"
	"github.com/databeast/cyberhud/display/modes/menu/styles"
	"github.com/databeast/cyberhud/display/style"
)

// menuRegistry is the per-mode StyleRegistry for the menu display mode.
var menuRegistry = style.NewRegistry[source.MenuSnapshot, source.Policy](
	styles.FramedStyle,
	styles.PlainStyle,
	styles.MonoSlow800x480Style,
)

// registeredStyleNames returns the list of style names from the menuRegistry.
func registeredStyleNames() []string {
	styles := menuRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}

// MenuRegistryEnumerate exposes registered menu styles for snapshot tests.
func MenuRegistryEnumerate() []style.Style[source.MenuSnapshot, source.Policy] {
	return menuRegistry.Enumerate()
}
