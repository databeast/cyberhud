package zmq

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/modes/zmq/styles"
	"github.com/databeast/cyberhud/display/style"
)

// zmqRegistry is the per-mode StyleRegistry for the ZMQ display mode.
//
// Styles follow the <capability>-<WxH> naming convention and the per-resolution
// pattern shared by the mature modes (serial, thermal, clock): one style per
// (capability floor × physical resolution) cell. Fitness resolution routes each
// panel to the nearest-capability, tightest-resolution style; a configured
// style name always wins, and any lower-capability style is a valid fallback,
// so the mono-slow entries act as the universal backstop for panels with no
// purpose-built cell.
//
// The first registered style is the default (used when there is a fitness tie
// or no panel hints). color-fast-240x240 is a common square TFT profile and
// matches the resolution of the previous default.
var zmqRegistry = style.NewRegistry[content.ZMQData, content.Policy](
	// Default + a common colour profile.
	styles.ColorFast240x240Style,

	// MonoSlow (e-paper mono) — universal fallbacks.
	styles.MonoSlow64x128Style,
	styles.MonoSlow128x64Style,
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

	// MonoFast (OLED mono).
	styles.MonoFast128x32Style,
	styles.MonoFast128x64Style,
	styles.MonoFast128x128Style,
	styles.MonoFast160x80Style,
	styles.MonoFast240x135Style,
	styles.MonoFast240x240Style,
	styles.MonoFast320x240Style,
	styles.MonoFast480x320Style,
	styles.MonoFast800x480Style,

	// GrayscaleSlow (grayscale e-paper).
	styles.GrayscaleSlow160x80Style,
	styles.GrayscaleSlow240x135Style,
	styles.GrayscaleSlow240x240Style,
	styles.GrayscaleSlow320x240Style,
	styles.GrayscaleSlow480x320Style,
	styles.GrayscaleSlow800x480Style,

	// GrayscaleFast.
	styles.GrayscaleFast160x80Style,
	styles.GrayscaleFast240x135Style,
	styles.GrayscaleFast240x240Style,
	styles.GrayscaleFast320x240Style,
	styles.GrayscaleFast480x320Style,
	styles.GrayscaleFast800x480Style,

	// ColorSlow (color e-paper).
	styles.ColorSlow122x250Style,
	styles.ColorSlow176x264Style,
	styles.ColorSlow200x200Style,
	styles.ColorSlow400x300Style,
	styles.ColorSlow800x480Style,

	// ColorFast (color TFT).
	styles.ColorFast128x128Style,
	styles.ColorFast160x80Style,
	styles.ColorFast240x135Style,
	styles.ColorFast240x320Style,
	styles.ColorFast320x240Style,
	styles.ColorFast320x480Style,
	styles.ColorFast480x320Style,
	styles.ColorFast480x800Style,
	styles.ColorFast800x480Style,
)

func registeredStyleNames() []string {
	styles := zmqRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
