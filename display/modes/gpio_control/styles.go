package gpio_control

import (
	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/modes/gpio_control/styles"
	"github.com/databeast/cyberhud/display/style"
)

// gpioControlRegistry is the per-mode StyleRegistry for the gpio_control display mode.
// Registration order: MonoSlow → MonoFast → GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
// The first entry (ListStyle) is the default returned by gpioControlRegistry.Default().
var gpioControlRegistry = style.NewRegistry[source.Data, source.Policy](
	// MonoSlow — existing generic styles + per-resolution
	styles.ListStyle,            // "list" (unconstrained, MonoSlow)
	styles.CompactStyle,         // "compact" (unconstrained, MonoSlow)
	styles.MonoSlow800x480Style, // "mono-slow-800x480" (bespoke BuildFn)
	styles.MonoSlow128x32Style,
	styles.MonoSlow128x64Style,
	styles.MonoSlow160x80Style,
	styles.MonoSlow128x128Style,
	styles.MonoSlow240x135Style,
	styles.MonoSlow240x240Style,
	styles.MonoSlow240x320Style,
	styles.MonoSlow320x240Style,
	styles.MonoSlow320x480Style,
	styles.MonoSlow480x320Style,

	// MonoFast
	styles.MonoFast128x32Style,
	styles.MonoFast128x64Style,
	styles.MonoFast160x80Style,
	styles.MonoFast128x128Style,
	styles.MonoFast240x135Style,
	styles.MonoFast240x240Style,
	styles.MonoFast240x320Style,
	styles.MonoFast320x240Style,
	styles.MonoFast320x480Style,
	styles.MonoFast480x320Style,
	styles.MonoFast800x480Style,

	// GrayscaleSlow
	styles.GrayscaleSlow128x32Style,
	styles.GrayscaleSlow128x64Style,
	styles.GrayscaleSlow160x80Style,
	styles.GrayscaleSlow128x128Style,
	styles.GrayscaleSlow240x135Style,
	styles.GrayscaleSlow240x240Style,
	styles.GrayscaleSlow240x320Style,
	styles.GrayscaleSlow320x240Style,
	styles.GrayscaleSlow320x480Style,
	styles.GrayscaleSlow480x320Style,
	styles.GrayscaleSlow800x480Style,

	// GrayscaleFast
	styles.GrayscaleFast128x32Style,
	styles.GrayscaleFast128x64Style,
	styles.GrayscaleFast160x80Style,
	styles.GrayscaleFast128x128Style,
	styles.GrayscaleFast240x135Style,
	styles.GrayscaleFast240x240Style,
	styles.GrayscaleFast240x320Style,
	styles.GrayscaleFast320x240Style,
	styles.GrayscaleFast320x480Style,
	styles.GrayscaleFast480x320Style,
	styles.GrayscaleFast800x480Style,

	// ColorSlow — existing grid + per-resolution
	styles.GridStyle, // "grid" (64x64 min, ColorSlow, bespoke BuildFn)
	styles.ColorSlow128x32Style,
	styles.ColorSlow128x64Style,
	styles.ColorSlow160x80Style,
	styles.ColorSlow128x128Style,
	styles.ColorSlow240x135Style,
	styles.ColorSlow240x240Style,
	styles.ColorSlow240x320Style,
	styles.ColorSlow320x240Style,
	styles.ColorSlow320x480Style,
	styles.ColorSlow480x320Style,
	styles.ColorSlow800x480Style,

	// ColorFast
	styles.ColorFast128x32Style,
	styles.ColorFast128x64Style,
	styles.ColorFast160x80Style,
	styles.ColorFast128x128Style,
	styles.ColorFast240x135Style,
	styles.ColorFast240x240Style,
	styles.ColorFast240x320Style,
	styles.ColorFast320x240Style,
	styles.ColorFast320x480Style,
	styles.ColorFast480x320Style,
	styles.ColorFast800x480Style,
)

// registeredStyleNames returns the list of style names from the registry.
// Used by catalog registration and cmdHandler for allowed-value validation.
func registeredStyleNames() []string {
	styles := gpioControlRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
