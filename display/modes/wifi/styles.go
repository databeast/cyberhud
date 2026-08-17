package wifi

import (
	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/modes/wifi/styles"
	"github.com/databeast/cyberhud/display/style"
)

// ─── Registry ─────────────────────────────────────────────────────────────────

// wifiRegistry is the per-mode StyleRegistry for the WiFi display mode.
// Registration order follows capability ordering: MonoSlow → MonoFast →
// GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
// Within each tier, resolutions are ordered small→large, landscape before portrait.
// The first style (mono-slow-128x128) is the default returned by wifiRegistry.Default().
var wifiRegistry = func() *style.StyleRegistry[source.WifiData, source.Policy] {
	r := style.NewRegistry[source.WifiData, source.Policy](
		// ─── MonoSlow (e-paper mono) ──────────────────────────────────────
		styles.MonoSlow128x128Style, // default
		styles.MonoSlow128x32Style,
		styles.MonoSlow32x128Style,
		styles.MonoSlow128x64Style,
		styles.MonoSlow160x80Style,
		styles.MonoSlow80x160Style,
		styles.MonoSlow160x128Style,
		styles.MonoSlow128x160Style,
		styles.MonoSlow200x200Style,
		styles.MonoSlow212x104Style,
		styles.MonoSlow240x135Style,
		styles.MonoSlow135x240Style,
		styles.MonoSlow240x240Style,
		styles.MonoSlow250x122Style,
		styles.MonoSlow264x176Style,
		styles.MonoSlow296x128Style,
		styles.MonoSlow128x296Style,
		styles.MonoSlow320x240Style,
		styles.MonoSlow240x320Style,
		styles.MonoSlow400x300Style,
		styles.MonoSlow480x320Style,
		styles.MonoSlow320x480Style,
		styles.MonoSlow800x480Style,
		styles.MonoSlow480x800Style,
		styles.MonoSlow800x600Style,

		// ─── MonoFast (OLED mono) ─────────────────────────────────────────
		styles.MonoFast128x32Style,
		styles.MonoFast32x128Style,
		styles.MonoFast128x64Style,
		styles.MonoFast128x128Style,
		styles.MonoFast160x80Style,
		styles.MonoFast80x160Style,
		styles.MonoFast160x128Style,
		styles.MonoFast128x160Style,
		styles.MonoFast200x200Style,
		styles.MonoFast212x104Style,
		styles.MonoFast240x135Style,
		styles.MonoFast135x240Style,
		styles.MonoFast240x240Style,
		styles.MonoFast250x122Style,
		styles.MonoFast264x176Style,
		styles.MonoFast296x128Style,
		styles.MonoFast128x296Style,
		styles.MonoFast320x240Style,
		styles.MonoFast240x320Style,
		styles.MonoFast400x300Style,
		styles.MonoFast480x320Style,
		styles.MonoFast320x480Style,
		styles.MonoFast800x480Style,
		styles.MonoFast480x800Style,
		styles.MonoFast800x600Style,

		// ─── GrayscaleSlow (grayscale e-paper) ────────────────────────────
		styles.GrayscaleSlow128x32Style,
		styles.GrayscaleSlow32x128Style,
		styles.GrayscaleSlow128x64Style,
		styles.GrayscaleSlow128x128Style,
		styles.GrayscaleSlow160x80Style,
		styles.GrayscaleSlow80x160Style,
		styles.GrayscaleSlow160x128Style,
		styles.GrayscaleSlow128x160Style,
		styles.GrayscaleSlow200x200Style,
		styles.GrayscaleSlow212x104Style,
		styles.GrayscaleSlow240x135Style,
		styles.GrayscaleSlow135x240Style,
		styles.GrayscaleSlow240x240Style,
		styles.GrayscaleSlow250x122Style,
		styles.GrayscaleSlow264x176Style,
		styles.GrayscaleSlow296x128Style,
		styles.GrayscaleSlow128x296Style,
		styles.GrayscaleSlow320x240Style,
		styles.GrayscaleSlow240x320Style,
		styles.GrayscaleSlow400x300Style,
		styles.GrayscaleSlow480x320Style,
		styles.GrayscaleSlow320x480Style,
		styles.GrayscaleSlow800x480Style,
		styles.GrayscaleSlow480x800Style,
		styles.GrayscaleSlow800x600Style,

		// ─── GrayscaleFast (grayscale LED matrix) ─────────────────────────
		styles.GrayscaleFast128x32Style,
		styles.GrayscaleFast32x128Style,
		styles.GrayscaleFast128x64Style,
		styles.GrayscaleFast128x128Style,
		styles.GrayscaleFast160x80Style,
		styles.GrayscaleFast80x160Style,
		styles.GrayscaleFast160x128Style,
		styles.GrayscaleFast128x160Style,
		styles.GrayscaleFast200x200Style,
		styles.GrayscaleFast212x104Style,
		styles.GrayscaleFast240x135Style,
		styles.GrayscaleFast135x240Style,
		styles.GrayscaleFast240x240Style,
		styles.GrayscaleFast250x122Style,
		styles.GrayscaleFast264x176Style,
		styles.GrayscaleFast296x128Style,
		styles.GrayscaleFast128x296Style,
		styles.GrayscaleFast320x240Style,
		styles.GrayscaleFast240x320Style,
		styles.GrayscaleFast400x300Style,
		styles.GrayscaleFast480x320Style,
		styles.GrayscaleFast320x480Style,
		styles.GrayscaleFast800x480Style,
		styles.GrayscaleFast480x800Style,
		styles.GrayscaleFast800x600Style,

		// ─── ColorSlow (color e-paper) ────────────────────────────────────
		styles.ColorSlow128x32Style,
		styles.ColorSlow32x128Style,
		styles.ColorSlow128x64Style,
		styles.ColorSlow128x128Style,
		styles.ColorSlow160x80Style,
		styles.ColorSlow80x160Style,
		styles.ColorSlow160x128Style,
		styles.ColorSlow128x160Style,
		styles.ColorSlow200x200Style,
		styles.ColorSlow212x104Style,
		styles.ColorSlow240x135Style,
		styles.ColorSlow135x240Style,
		styles.ColorSlow240x240Style,
		styles.ColorSlow250x122Style,
		styles.ColorSlow264x176Style,
		styles.ColorSlow296x128Style,
		styles.ColorSlow128x296Style,
		styles.ColorSlow320x240Style,
		styles.ColorSlow240x320Style,
		styles.ColorSlow400x300Style,
		styles.ColorSlow480x320Style,
		styles.ColorSlow320x480Style,
		styles.ColorSlow800x480Style,
		styles.ColorSlow480x800Style,
		styles.ColorSlow800x600Style,

		// ─── ColorFast (color TFT) ────────────────────────────────────────
		styles.ColorFast128x32Style,
		styles.ColorFast32x128Style,
		styles.ColorFast128x64Style,
		styles.ColorSmall128x128Style, // bespoke
		styles.ColorFast160x80Style,
		styles.ColorFast80x160Style,
		styles.ColorFast160x128Style,
		styles.ColorFast128x160Style,
		styles.ColorFast200x200Style,
		styles.ColorFast212x104Style,
		styles.ColorFast240x135Style,
		styles.ColorFast135x240Style,
		styles.Color240x240Style, // bespoke
		styles.ColorFast250x122Style,
		styles.ColorFast264x176Style,
		styles.ColorFast296x128Style,
		styles.ColorFast128x296Style,
		styles.Color320x240Style, // bespoke
		styles.ColorFast240x320Style,
		styles.ColorFast400x300Style,
		styles.ColorFast480x320Style,
		styles.ColorFast320x480Style,
		styles.ColorFast800x480Style,
		styles.ColorFast480x800Style,
		styles.ColorFast800x600Style,
	)

	return r
}()

// registeredStyleNames returns the list of style names from the registry.
// Used by catalog registration and cmdHandler for allowed-value validation.
func registeredStyleNames() []string {
	styles := wifiRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
