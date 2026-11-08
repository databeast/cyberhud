package schedule

import (
	"github.com/databeast/cyberhud/display/modes/schedule/source"
	"github.com/databeast/cyberhud/display/modes/schedule/styles"
	"github.com/databeast/cyberhud/display/style"
)

// scheduleRegistry is the per-mode StyleRegistry for the schedule display mode.
//
// Schedule has one visual intent — the "what's on now" card — declared per panel
// in the styles subpackage and listed here explicitly. Registration order
// follows capability ordering (MonoSlow → MonoFast → GrayscaleSlow →
// GrayscaleFast → ColorSlow → ColorFast), which also defines fitness
// tie-breaking (first registered wins). The first entry is the default.
var scheduleRegistry = style.NewRegistry[source.Snapshot, source.Policy](
	// MonoSlow
	styles.MonoSlow128x64Style,
	styles.MonoSlow128x128Style,
	styles.MonoSlow200x200Style,
	styles.MonoSlow250x122Style,
	styles.MonoSlow296x128Style,
	styles.MonoSlow400x300Style,
	styles.MonoSlow800x480Style,

	// MonoFast
	styles.MonoFast128x32Style,
	styles.MonoFast128x64Style,
	styles.MonoFast128x128Style,
	styles.MonoFast240x240Style,
	styles.MonoFast320x240Style,

	// GrayscaleSlow
	styles.GrayscaleSlow200x200Style,
	styles.GrayscaleSlow296x128Style,
	styles.GrayscaleSlow400x300Style,
	styles.GrayscaleSlow480x800Style,
	styles.GrayscaleSlow800x480Style,

	// GrayscaleFast
	styles.GrayscaleFast128x128Style,
	styles.GrayscaleFast240x240Style,
	styles.GrayscaleFast320x240Style,
	styles.GrayscaleFast480x320Style,

	// ColorSlow
	styles.ColorSlow200x200Style,
	styles.ColorSlow296x128Style,
	styles.ColorSlow400x300Style,
	styles.ColorSlow480x800Style,
	styles.ColorSlow800x480Style,

	// ColorFast
	styles.ColorFast128x128Style,
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
