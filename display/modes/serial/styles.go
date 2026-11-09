package serial

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/modes/serial/styles"
	"github.com/databeast/cyberhud/display/style"
)

// serialRegistry is the per-mode StyleRegistry for the serial display mode.
//
// Serial has one visual intent — the console instrument panel — declared per
// panel in the styles subpackage's styles_WxH.go files and listed here
// explicitly. Registration order follows capability ordering (MonoSlow →
// MonoFast → GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast), which also
// defines fitness tie-breaking (first registered wins). The first entry is the
// default returned by Default().
var serialRegistry = style.NewRegistry[source.Snapshot, source.Policy](
	// ─── MonoSlow (e-paper mono) ──────────────────────────────────────────
	styles.MonoSlow104x212Style,
	styles.MonoSlow122x250Style,
	styles.MonoSlow176x264Style,
	styles.MonoSlow200x200Style,
	styles.MonoSlow212x104Style,
	styles.MonoSlow250x122Style,
	styles.MonoSlow264x176Style,
	styles.MonoSlow296x128Style,
	styles.MonoSlow300x400Style,
	styles.MonoSlow400x300Style,
	styles.MonoSlow480x800Style,
	styles.MonoSlow800x480Style,

	// ─── MonoFast (OLED mono) ─────────────────────────────────────────────
	styles.MonoFast128x32Style,
	styles.MonoFast128x64Style,
	styles.MonoFast128x128Style,
	styles.MonoFast160x80Style,
	styles.MonoFast240x135Style,
	styles.MonoFast240x240Style,
	styles.MonoFast320x240Style,
	styles.MonoFast480x320Style,
	styles.MonoFast800x480Style,

	// ─── GrayscaleSlow (grayscale e-paper) ────────────────────────────────
	styles.GrayscaleSlow104x212Style,
	styles.GrayscaleSlow122x250Style,
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

	// ─── GrayscaleFast (grayscale LED matrix) ─────────────────────────────
	styles.GrayscaleFast15x7Style,
	styles.GrayscaleFast16x8Style,
	styles.GrayscaleFast128x128Style,
	styles.GrayscaleFast160x80Style,
	styles.GrayscaleFast240x135Style,
	styles.GrayscaleFast240x240Style,
	styles.GrayscaleFast240x320Style,
	styles.GrayscaleFast320x240Style,
	styles.GrayscaleFast320x480Style,
	styles.GrayscaleFast480x320Style,
	styles.GrayscaleFast480x800Style,
	styles.GrayscaleFast800x480Style,

	// ─── ColorSlow (color e-paper) ────────────────────────────────────────
	styles.ColorSlow104x212Style,
	styles.ColorSlow122x250Style,
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

	// ─── ColorFast (color TFT) ────────────────────────────────────────────
	styles.ColorFast13x9Style,
	styles.ColorFast128x128Style,
	styles.ColorFast128x160Style,
	styles.ColorFast160x80Style,
	styles.ColorFast240x135Style,
	styles.ColorFast240x240Style,
	styles.ColorFast240x320Style,
	styles.ColorFast320x240Style,
	styles.ColorFast320x480Style,
	styles.ColorFast480x320Style,
	styles.ColorFast480x800Style,
	styles.ColorFast800x480Style,
)
