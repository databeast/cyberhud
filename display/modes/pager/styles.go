package pager

import (
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/modes/pager/styles"
	"github.com/databeast/cyberhud/display/style"
)

// ─── Registry ─────────────────────────────────────────────────────────────────

// pagerRegistry is the per-mode StyleRegistry for the pager display mode.
// Registration order follows capability ordering: MonoSlow → MonoFast →
// GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
var pagerRegistry = func() *style.StyleRegistry[source.PagerSnapshot, source.Policy] {
	r := style.NewRegistry[source.PagerSnapshot, source.Policy](
		// MonoSlow (e-paper mono)
		styles.PagerMonoSlow122x250Style,
		styles.PagerMonoSlow176x264Style,
		styles.PagerMonoSlow200x200Style,
		styles.PagerMonoSlow212x104Style,
		styles.PagerMonoSlow296x128Style,
		styles.PagerMonoSlow400x300Style,
		styles.PagerMonoSlow480x800Style,
		styles.PagerMonoSlow800x480Style,
		styles.PagerMonoSlow104x212Style,
		styles.PagerMonoSlow250x122Style,
		styles.PagerMonoSlow128x296Style,
		styles.PagerMonoSlow264x176Style,
		styles.PagerMonoSlow300x400Style,

		// MonoFast (OLED mono)
		styles.PagerMonoFast128x32Style,
		styles.PagerMonoFast128x64Style,
		styles.PagerMonoFast128x128Style,
		styles.PagerMonoFast32x128Style,
		styles.PagerMonoFast64x128Style,

		// GrayscaleSlow (grayscale e-paper)
		styles.PagerGrayscaleSlow122x250Style,
		styles.PagerGrayscaleSlow176x264Style,
		styles.PagerGrayscaleSlow200x200Style,
		styles.PagerGrayscaleSlow212x104Style,
		styles.PagerGrayscaleSlow296x128Style,
		styles.PagerGrayscaleSlow400x300Style,
		styles.PagerGrayscaleSlow480x800Style,
		styles.PagerGrayscaleSlow800x480Style,
		styles.PagerGrayscaleSlow104x212Style,
		styles.PagerGrayscaleSlow250x122Style,
		styles.PagerGrayscaleSlow128x296Style,
		styles.PagerGrayscaleSlow264x176Style,
		styles.PagerGrayscaleSlow300x400Style,

		// GrayscaleFast (grayscale LED matrix)
		styles.PagerGrayscaleFast160x80Style,
		styles.PagerGrayscaleFast160x128Style,
		styles.PagerGrayscaleFast240x135Style,
		styles.PagerGrayscaleFast240x240Style,
		styles.PagerGrayscaleFast320x240Style,
		styles.PagerGrayscaleFast480x320Style,
		styles.PagerGrayscaleFast800x480Style,
		styles.PagerGrayscaleFast80x160Style,
		styles.PagerGrayscaleFast128x160Style,
		styles.PagerGrayscaleFast135x240Style,
		styles.PagerGrayscaleFast240x320Style,
		styles.PagerGrayscaleFast320x480Style,
		styles.PagerGrayscaleFast480x800Style,
		styles.PagerGrayscaleFast128x128Style,

		// ColorSlow (color e-paper)
		styles.PagerColorSlow122x250Style,
		styles.PagerColorSlow176x264Style,
		styles.PagerColorSlow200x200Style,
		styles.PagerColorSlow212x104Style,
		styles.PagerColorSlow296x128Style,
		styles.PagerColorSlow400x300Style,
		styles.PagerColorSlow480x800Style,
		styles.PagerColorSlow800x480Style,
		styles.PagerColorSlow104x212Style,
		styles.PagerColorSlow250x122Style,
		styles.PagerColorSlow128x296Style,
		styles.PagerColorSlow264x176Style,
		styles.PagerColorSlow300x400Style,

		// ColorFast (color TFT)
		styles.PagerColorFast160x80Style,
		styles.PagerColorFast160x128Style,
		styles.PagerColorFast240x135Style,
		styles.PagerColorFast240x240Style,
		styles.PagerColorFast320x240Style,
		styles.PagerColorFast480x320Style,
		styles.PagerColorFast800x480Style,
		styles.PagerColorFast80x160Style,
		styles.PagerColorFast128x160Style,
		styles.PagerColorFast135x240Style,
		styles.PagerColorFast240x320Style,
		styles.PagerColorFast320x480Style,
		styles.PagerColorFast480x800Style,
		styles.PagerColorFast128x128Style,

		// MonoSlow polished (9-step pipeline)
		styles.MonoSlow800x480Style,
	)
	return r
}()

func registeredStyleNames() []string {
	styles := pagerRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
