package clock

import (
	"github.com/databeast/cyberhud/display/modes/clock/source"
	"github.com/databeast/cyberhud/display/modes/clock/styles"
	"github.com/databeast/cyberhud/display/style"
)

// clockRegistry is the per-mode StyleRegistry for the clock display mode.
// Registration order: MonoSlow → MonoFast → GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
// The first entry (MonoSlow32x128Style) is the default returned by clockRegistry.Default().
var clockRegistry = func() *style.StyleRegistry[source.ClockData, source.Policy] {
	r := style.NewRegistry[source.ClockData, source.Policy](
		// MonoSlow — skeletons (16) + polished (1)
		styles.MonoSlow32x128Style,
		styles.MonoSlow64x128Style,
		styles.MonoSlow80x160Style,
		styles.MonoSlow128x32Style,
		styles.MonoSlow128x64Style,
		styles.MonoSlow128x128Style,
		styles.MonoSlow128x160Style,
		styles.MonoSlow135x240Style,
		styles.MonoSlow160x80Style,
		styles.MonoSlow160x128Style,
		styles.MonoSlow240x135Style,
		styles.MonoSlow240x240Style,
		styles.MonoSlow240x320Style,
		styles.MonoSlow320x240Style,
		styles.MonoSlow320x480Style,
		styles.MonoSlow480x320Style,
		styles.MonoSlow800x480Style,

		// MonoFast — existing mono (5)
		styles.MonoSmall128x32Style,
		styles.MonoSmall128x64Style,
		styles.MonoSmall128x128Style,
		styles.MonoSmall32x128Style,
		styles.MonoSmall64x128Style,

		// GrayscaleSlow — all skeletons (29)
		styles.GrayscaleSlow32x128Style,
		styles.GrayscaleSlow64x128Style,
		styles.GrayscaleSlow80x160Style,
		styles.GrayscaleSlow104x212Style,
		styles.GrayscaleSlow122x250Style,
		styles.GrayscaleSlow128x32Style,
		styles.GrayscaleSlow128x64Style,
		styles.GrayscaleSlow128x128Style,
		styles.GrayscaleSlow128x160Style,
		styles.GrayscaleSlow128x296Style,
		styles.GrayscaleSlow135x240Style,
		styles.GrayscaleSlow160x80Style,
		styles.GrayscaleSlow160x128Style,
		styles.GrayscaleSlow176x264Style,
		styles.GrayscaleSlow200x200Style,
		styles.GrayscaleSlow212x104Style,
		styles.GrayscaleSlow240x135Style,
		styles.GrayscaleSlow240x240Style,
		styles.GrayscaleSlow240x320Style,
		styles.GrayscaleSlow250x122Style,
		styles.GrayscaleSlow264x176Style,
		styles.GrayscaleSlow296x128Style,
		styles.GrayscaleSlow300x400Style,
		styles.GrayscaleSlow320x240Style,
		styles.GrayscaleSlow320x480Style,
		styles.GrayscaleSlow400x300Style,
		styles.GrayscaleSlow480x320Style,
		styles.GrayscaleSlow480x800Style,
		styles.GrayscaleSlow800x480Style,

		// GrayscaleFast (14 existing)
		styles.GrayscaleFast80x160Style,
		styles.GrayscaleFast128x128Style,
		styles.GrayscaleFast128x160Style,
		styles.GrayscaleFast135x240Style,
		styles.GrayscaleFast160x80Style,
		styles.GrayscaleFast160x128Style,
		styles.GrayscaleFast240x135Style,
		styles.GrayscaleFast240x240Style,
		styles.GrayscaleFast240x320Style,
		styles.GrayscaleFast320x240Style,
		styles.GrayscaleFast320x480Style,
		styles.GrayscaleFast480x320Style,
		styles.GrayscaleFast480x800Style,
		styles.GrayscaleFast800x480Style,

		// ColorSlow — all skeletons (29)
		styles.ColorSlow32x128Style,
		styles.ColorSlow64x128Style,
		styles.ColorSlow80x160Style,
		styles.ColorSlow104x212Style,
		styles.ColorSlow122x250Style,
		styles.ColorSlow128x32Style,
		styles.ColorSlow128x64Style,
		styles.ColorSlow128x128Style,
		styles.ColorSlow128x160Style,
		styles.ColorSlow128x296Style,
		styles.ColorSlow135x240Style,
		styles.ColorSlow160x80Style,
		styles.ColorSlow160x128Style,
		styles.ColorSlow176x264Style,
		styles.ColorSlow200x200Style,
		styles.ColorSlow212x104Style,
		styles.ColorSlow240x135Style,
		styles.ColorSlow240x240Style,
		styles.ColorSlow240x320Style,
		styles.ColorSlow250x122Style,
		styles.ColorSlow264x176Style,
		styles.ColorSlow296x128Style,
		styles.ColorSlow300x400Style,
		styles.ColorSlow320x240Style,
		styles.ColorSlow320x480Style,
		styles.ColorSlow400x300Style,
		styles.ColorSlow480x320Style,
		styles.ColorSlow480x800Style,
		styles.ColorSlow800x480Style,

		// ColorFast — existing color (14)
		styles.ColorSmall80x160Style,
		styles.ColorSmall128x128Style,
		styles.ColorSmall128x160Style,
		styles.ColorMedium135x240Style,
		styles.ColorSmall160x80Style,
		styles.ColorSmall160x128Style,
		styles.ColorMedium240x135Style,
		styles.ColorMedium240x240Style,
		styles.ColorMedium240x320Style,
		styles.ColorMedium320x240Style,
		styles.ColorLarge320x480Style,
		styles.ColorLarge480x320Style,
		styles.ColorLarge480x800Style,
		styles.ColorLarge800x480Style,
	)

	return r
}()

// registeredStyleNames returns the names of all registered clock styles
// in registration order from the clockRegistry.
func registeredStyleNames() []string {
	styles := clockRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
