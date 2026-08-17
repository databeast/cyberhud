package ticker

import (
	"github.com/databeast/cyberhud/display/modes/ticker/source"
	"github.com/databeast/cyberhud/display/modes/ticker/styles"
	"github.com/databeast/cyberhud/display/style"
)

// registeredStyleNames returns the list of style names from tickerRegistry
// for use in catalog registration and error messages.
func registeredStyleNames() []string {
	styles := tickerRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}

// ─── Registry ─────────────────────────────────────────────────────────────────

// tickerRegistry is the per-mode StyleRegistry for the ticker display mode.
// Registration order: resolution-specific variants grouped by capability class.
var tickerRegistry = style.NewRegistry[source.TickerSnapshot, source.Policy](
	// MonoSlow (e-paper mono)
	styles.TickerMonoSlow122x250Style,
	styles.TickerMonoSlow176x264Style,
	styles.TickerMonoSlow200x200Style,
	styles.TickerMonoSlow212x104Style,
	styles.TickerMonoSlow296x128Style,
	styles.TickerMonoSlow400x300Style,
	styles.TickerMonoSlow480x800Style,
	styles.TickerMonoSlow104x212Style,
	styles.TickerMonoSlow250x122Style,
	styles.TickerMonoSlow128x296Style,
	styles.TickerMonoSlow264x176Style,
	styles.TickerMonoSlow300x400Style,

	// MonoFast (OLED mono)
	styles.TickerMonoFast128x32Style,
	styles.TickerMonoFast128x64Style,
	styles.TickerMonoFast128x128Style,
	styles.TickerMonoFast32x128Style,
	styles.TickerMonoFast64x128Style,

	// GrayscaleSlow (grayscale e-paper)
	styles.TickerGrayscaleSlow122x250Style,
	styles.TickerGrayscaleSlow176x264Style,
	styles.TickerGrayscaleSlow200x200Style,
	styles.TickerGrayscaleSlow212x104Style,
	styles.TickerGrayscaleSlow296x128Style,
	styles.TickerGrayscaleSlow400x300Style,
	styles.TickerGrayscaleSlow480x800Style,
	styles.TickerGrayscaleSlow800x480Style,
	styles.TickerGrayscaleSlow104x212Style,
	styles.TickerGrayscaleSlow250x122Style,
	styles.TickerGrayscaleSlow128x296Style,
	styles.TickerGrayscaleSlow264x176Style,
	styles.TickerGrayscaleSlow300x400Style,

	// GrayscaleFast (grayscale LED matrix)
	styles.TickerGrayscaleFast160x80Style,
	styles.TickerGrayscaleFast160x128Style,
	styles.TickerGrayscaleFast240x135Style,
	styles.TickerGrayscaleFast240x240Style,
	styles.TickerGrayscaleFast320x240Style,
	styles.TickerGrayscaleFast480x320Style,
	styles.TickerGrayscaleFast800x480Style,
	styles.TickerGrayscaleFast80x160Style,
	styles.TickerGrayscaleFast128x160Style,
	styles.TickerGrayscaleFast135x240Style,
	styles.TickerGrayscaleFast240x320Style,
	styles.TickerGrayscaleFast320x480Style,
	styles.TickerGrayscaleFast480x800Style,
	styles.TickerGrayscaleFast128x128Style,

	// ColorSlow (color e-paper)
	styles.TickerColorSlow122x250Style,
	styles.TickerColorSlow176x264Style,
	styles.TickerColorSlow200x200Style,
	styles.TickerColorSlow212x104Style,
	styles.TickerColorSlow296x128Style,
	styles.TickerColorSlow400x300Style,
	styles.TickerColorSlow480x800Style,
	styles.TickerColorSlow800x480Style,
	styles.TickerColorSlow104x212Style,
	styles.TickerColorSlow250x122Style,
	styles.TickerColorSlow128x296Style,
	styles.TickerColorSlow264x176Style,
	styles.TickerColorSlow300x400Style,

	// ColorFast (color TFT)
	styles.TickerColorFast128x128Style,
	styles.TickerColorFast160x80Style,
	styles.TickerColorFast160x128Style,
	styles.TickerColorFast240x135Style,
	styles.TickerColorFast240x240Style,
	styles.TickerColorFast240x320Style,
	styles.TickerColorFast320x240Style,
	styles.TickerColorFast320x480Style,
	styles.TickerColorFast480x320Style,
	styles.TickerColorFast480x800Style,
	styles.TickerColorFast800x480Style,
	styles.TickerColorFast80x160Style,
	styles.TickerColorFast128x160Style,
	styles.TickerColorFast135x240Style,

	// MonoSlow polished (9-step pipeline)
	styles.MonoSlow800x480Style,
)
