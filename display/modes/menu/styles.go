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

	// MonoFast
	styles.MonoFast128x32Style,
	styles.MonoFast128x64Style,
	styles.MonoFast168x144Style,
	styles.MonoFast400x240Style,

	// GrayscaleFast
	styles.GrayscaleFast15x7Style,
	styles.GrayscaleFast16x8Style,
	styles.GrayscaleFast128x64Style,

	// ColorSlow
	styles.ColorSlow200x200Style,

	// ColorFast
	styles.ColorFast13x9Style,
	styles.ColorFast128x128Style,
	styles.ColorFast128x160Style,
	styles.ColorFast160x80Style,
	styles.ColorFast240x135Style,
	styles.ColorFast240x240Style,
	styles.ColorFast240x320Style,
	styles.ColorFast320x240Style,
	styles.ColorFast480x320Style,
)

// MenuRegistryEnumerate exposes registered menu styles for snapshot tests.
func MenuRegistryEnumerate() []style.Style[source.MenuSnapshot, source.Policy] {
	return menuRegistry.Enumerate()
}
