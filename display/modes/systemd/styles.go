package systemd

import (
	"github.com/databeast/cyberhud/display/modes/systemd/source"
	"github.com/databeast/cyberhud/display/modes/systemd/styles"
	"github.com/databeast/cyberhud/display/style"
)

// systemdRegistry is the per-mode StyleRegistry for the systemd display mode.
// Registration order: MonoSlow → MonoFast → GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
// The first entry (EinkStyle104x212) is the default style.
var systemdRegistry = func() *style.StyleRegistry[source.Snapshot, source.Policy] {
	r := style.NewRegistry[source.Snapshot, source.Policy](
		// MonoSlow
		styles.EinkStyle104x212,
		styles.EinkStyle122x250,
		styles.EinkStyle128x296,
		styles.EinkStyle176x264,
		styles.EinkStyle200x200,
		styles.EinkStyle212x104,
		styles.EinkStyle250x122,
		styles.EinkStyle264x176,
		styles.EinkStyle296x128,
		styles.EinkStyle300x400,
		styles.EinkStyle400x300,
		styles.EinkStyle480x800,
		styles.EinkStyle800x480,
		styles.MonoSlow128x32Style,
		styles.MonoSlow128x64Style,
		styles.MonoSlow128x128Style,
		styles.MonoSlow32x128Style,
		styles.MonoSlow64x128Style,
		styles.MonoSlow160x80Style,
		styles.MonoSlow128x160Style,
		styles.MonoSlow160x128Style,
		styles.MonoSlow240x135Style,
		styles.MonoSlow135x240Style,
		styles.MonoSlow240x240Style,
		styles.MonoSlow240x320Style,
		styles.MonoSlow320x240Style,
		styles.MonoSlow320x480Style,
		styles.MonoSlow480x320Style,
		styles.MonoSlow480x800Style,
		styles.MonoSlow800x480Style,
		// MonoFast
		styles.MonoStyle128x32,
		styles.MonoStyle128x64,
		styles.MonoStyle128x128,
		styles.MonoStyle32x128,
		styles.MonoStyle64x128,
		styles.MonoFast104x212Style,
		styles.MonoFast122x250Style,
		styles.MonoFast128x296Style,
		styles.MonoFast176x264Style,
		styles.MonoFast200x200Style,
		styles.MonoFast212x104Style,
		styles.MonoFast250x122Style,
		styles.MonoFast264x176Style,
		styles.MonoFast296x128Style,
		styles.MonoFast300x400Style,
		styles.MonoFast400x300Style,
		styles.MonoFast160x80Style,
		styles.MonoFast128x160Style,
		styles.MonoFast160x128Style,
		styles.MonoFast240x135Style,
		styles.MonoFast135x240Style,
		styles.MonoFast240x240Style,
		styles.MonoFast240x320Style,
		styles.MonoFast320x240Style,
		styles.MonoFast320x480Style,
		styles.MonoFast480x320Style,
		styles.MonoFast480x800Style,
		styles.MonoFast800x480Style,
		// GrayscaleSlow
		styles.GrayscaleSlow128x32Style,
		styles.GrayscaleSlow128x64Style,
		styles.GrayscaleSlow128x128Style,
		styles.GrayscaleSlow32x128Style,
		styles.GrayscaleSlow64x128Style,
		styles.GrayscaleSlow104x212Style,
		styles.GrayscaleSlow122x250Style,
		styles.GrayscaleSlow128x296Style,
		styles.GrayscaleSlow176x264Style,
		styles.GrayscaleSlow200x200Style,
		styles.GrayscaleSlow212x104Style,
		styles.GrayscaleSlow250x122Style,
		styles.GrayscaleSlow264x176Style,
		styles.GrayscaleSlow296x128Style,
		styles.GrayscaleSlow300x400Style,
		styles.GrayscaleSlow400x300Style,
		styles.GrayscaleSlow480x800Style,
		styles.GrayscaleSlow800x480Style,
		styles.GrayscaleSlow160x80Style,
		styles.GrayscaleSlow128x160Style,
		styles.GrayscaleSlow160x128Style,
		styles.GrayscaleSlow240x135Style,
		styles.GrayscaleSlow135x240Style,
		styles.GrayscaleSlow240x240Style,
		styles.GrayscaleSlow240x320Style,
		styles.GrayscaleSlow320x240Style,
		styles.GrayscaleSlow320x480Style,
		styles.GrayscaleSlow480x320Style,
		// GrayscaleFast
		styles.GrayscaleFast160x80Style,
		styles.GrayscaleFast128x128Style,
		styles.GrayscaleFast128x160Style,
		styles.GrayscaleFast160x128Style,
		styles.GrayscaleFast240x135Style,
		styles.GrayscaleFast135x240Style,
		styles.GrayscaleFast240x240Style,
		styles.GrayscaleFast240x320Style,
		styles.GrayscaleFast320x240Style,
		styles.GrayscaleFast320x480Style,
		styles.GrayscaleFast480x320Style,
		styles.GrayscaleFast480x800Style,
		styles.GrayscaleFast800x480Style,
		styles.GrayscaleFast128x32Style,
		styles.GrayscaleFast128x64Style,
		styles.GrayscaleFast32x128Style,
		styles.GrayscaleFast64x128Style,
		styles.GrayscaleFast104x212Style,
		styles.GrayscaleFast122x250Style,
		styles.GrayscaleFast128x296Style,
		styles.GrayscaleFast176x264Style,
		styles.GrayscaleFast200x200Style,
		styles.GrayscaleFast212x104Style,
		styles.GrayscaleFast250x122Style,
		styles.GrayscaleFast264x176Style,
		styles.GrayscaleFast296x128Style,
		styles.GrayscaleFast300x400Style,
		styles.GrayscaleFast400x300Style,
		// ColorSlow
		styles.ColorSlow128x32Style,
		styles.ColorSlow128x64Style,
		styles.ColorSlow128x128Style,
		styles.ColorSlow32x128Style,
		styles.ColorSlow64x128Style,
		styles.ColorSlow104x212Style,
		styles.ColorSlow122x250Style,
		styles.ColorSlow128x296Style,
		styles.ColorSlow176x264Style,
		styles.ColorSlow200x200Style,
		styles.ColorSlow212x104Style,
		styles.ColorSlow250x122Style,
		styles.ColorSlow264x176Style,
		styles.ColorSlow296x128Style,
		styles.ColorSlow300x400Style,
		styles.ColorSlow400x300Style,
		styles.ColorSlow480x800Style,
		styles.ColorSlow800x480Style,
		styles.ColorSlow160x80Style,
		styles.ColorSlow128x160Style,
		styles.ColorSlow160x128Style,
		styles.ColorSlow240x135Style,
		styles.ColorSlow135x240Style,
		styles.ColorSlow240x240Style,
		styles.ColorSlow240x320Style,
		styles.ColorSlow320x240Style,
		styles.ColorSlow320x480Style,
		styles.ColorSlow480x320Style,
		// ColorFast
		styles.ColorStyle240x240, // default color style
		styles.ColorStyle160x80,
		styles.ColorStyle128x128,
		styles.ColorStyle128x160,
		styles.ColorStyle160x128,
		styles.ColorStyle240x135,
		styles.ColorStyle135x240,
		styles.ColorStyle240x320,
		styles.ColorStyle320x240,
		styles.ColorStyle320x480,
		styles.ColorStyle480x320,
		styles.ColorStyle480x800,
		styles.ColorStyle800x480,
		styles.ColorFast128x32Style,
		styles.ColorFast128x64Style,
		styles.ColorFast32x128Style,
		styles.ColorFast64x128Style,
		styles.ColorFast104x212Style,
		styles.ColorFast122x250Style,
		styles.ColorFast128x296Style,
		styles.ColorFast176x264Style,
		styles.ColorFast200x200Style,
		styles.ColorFast212x104Style,
		styles.ColorFast250x122Style,
		styles.ColorFast264x176Style,
		styles.ColorFast296x128Style,
		styles.ColorFast300x400Style,
		styles.ColorFast400x300Style,
	)

	return r
}()

// allEmptyItems returns true if items is nil, empty, or contains only empty strings.
func allEmptyItems(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item != "" {
			return false
		}
	}
	return true
}

// registeredStyleNames returns the list of style names from the registry.
// Used by catalog registration and cmdHandler for allowed-value validation.
func registeredStyleNames() []string {
	styles := systemdRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
