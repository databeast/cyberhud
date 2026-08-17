package attract_plasma

import (
	"github.com/databeast/cyberhud/display/modes/attract_plasma/source"
	"github.com/databeast/cyberhud/display/modes/attract_plasma/styles"
	"github.com/databeast/cyberhud/display/style"
)

// Compile-time interface compliance checks.
var (
	_ style.Style[source.Snapshot, source.Policy] = styles.MonoSlowStyle{}
	_ style.Style[source.Snapshot, source.Policy] = styles.MonoFastStyle{}
	_ style.Style[source.Snapshot, source.Policy] = styles.GrayscaleSlowStyle{}
	_ style.Style[source.Snapshot, source.Policy] = styles.GrayscaleFastStyle{}
	_ style.Style[source.Snapshot, source.Policy] = styles.ColorSlowStyle{}
	_ style.Style[source.Snapshot, source.Policy] = styles.ColorFastStyle{}
)

// ─── Registry ─────────────────────────────────────────────────────────────────

// plasmaRegistry is the per-mode StyleRegistry for the plasma display mode.
// Registration order follows capability ordering: MonoSlow → MonoFast →
// GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
var plasmaRegistry = func() *style.StyleRegistry[source.Snapshot, source.Policy] {
	r := style.NewRegistry[source.Snapshot, source.Policy](
		// MonoSlow (e-paper mono)
		styles.MonoSlowStyle{Width: 122, Height: 250},
		styles.MonoSlowStyle{Width: 176, Height: 264},
		styles.MonoSlowStyle{Width: 200, Height: 200},
		styles.MonoSlowStyle{Width: 212, Height: 104},
		styles.MonoSlowStyle{Width: 296, Height: 128},
		styles.MonoSlowStyle{Width: 400, Height: 300},
		styles.MonoSlowStyle{Width: 480, Height: 800},
		styles.MonoSlowStyle{Width: 800, Height: 480},
		styles.MonoSlowStyle{Width: 104, Height: 212},
		styles.MonoSlowStyle{Width: 250, Height: 122},
		styles.MonoSlowStyle{Width: 128, Height: 296},
		styles.MonoSlowStyle{Width: 264, Height: 176},
		styles.MonoSlowStyle{Width: 300, Height: 400},

		// MonoFast (OLED mono)
		styles.MonoFastStyle{Width: 16, Height: 8},
		styles.MonoFastStyle{Width: 8, Height: 16},
		styles.MonoFastStyle{Width: 128, Height: 32},
		styles.MonoFastStyle{Width: 128, Height: 64},
		styles.MonoFastStyle{Width: 128, Height: 128},
		styles.MonoFastStyle{Width: 32, Height: 128},
		styles.MonoFastStyle{Width: 64, Height: 128},

		// GrayscaleSlow (grayscale e-paper)
		styles.GrayscaleSlowStyle{Width: 122, Height: 250},
		styles.GrayscaleSlowStyle{Width: 176, Height: 264},
		styles.GrayscaleSlowStyle{Width: 200, Height: 200},
		styles.GrayscaleSlowStyle{Width: 212, Height: 104},
		styles.GrayscaleSlowStyle{Width: 296, Height: 128},
		styles.GrayscaleSlowStyle{Width: 400, Height: 300},
		styles.GrayscaleSlowStyle{Width: 480, Height: 800},
		styles.GrayscaleSlowStyle{Width: 800, Height: 480},
		styles.GrayscaleSlowStyle{Width: 104, Height: 212},
		styles.GrayscaleSlowStyle{Width: 250, Height: 122},
		styles.GrayscaleSlowStyle{Width: 128, Height: 296},
		styles.GrayscaleSlowStyle{Width: 264, Height: 176},
		styles.GrayscaleSlowStyle{Width: 300, Height: 400},

		// GrayscaleFast (grayscale LED matrix)
		styles.GrayscaleFastStyle{Width: 16, Height: 8},
		styles.GrayscaleFastStyle{Width: 8, Height: 16},
		styles.GrayscaleFastStyle{Width: 160, Height: 80},
		styles.GrayscaleFastStyle{Width: 160, Height: 128},
		styles.GrayscaleFastStyle{Width: 240, Height: 135},
		styles.GrayscaleFastStyle{Width: 240, Height: 240},
		styles.GrayscaleFastStyle{Width: 320, Height: 240},
		styles.GrayscaleFastStyle{Width: 480, Height: 320},
		styles.GrayscaleFastStyle{Width: 800, Height: 480},
		styles.GrayscaleFastStyle{Width: 80, Height: 160},
		styles.GrayscaleFastStyle{Width: 128, Height: 160},
		styles.GrayscaleFastStyle{Width: 135, Height: 240},
		styles.GrayscaleFastStyle{Width: 240, Height: 320},
		styles.GrayscaleFastStyle{Width: 320, Height: 480},
		styles.GrayscaleFastStyle{Width: 480, Height: 800},
		styles.GrayscaleFastStyle{Width: 128, Height: 128},

		// ColorSlow (color e-paper)
		styles.ColorSlowStyle{Width: 122, Height: 250},
		styles.ColorSlowStyle{Width: 176, Height: 264},
		styles.ColorSlowStyle{Width: 200, Height: 200},
		styles.ColorSlowStyle{Width: 212, Height: 104},
		styles.ColorSlowStyle{Width: 296, Height: 128},
		styles.ColorSlowStyle{Width: 400, Height: 300},
		styles.ColorSlowStyle{Width: 480, Height: 800},
		styles.ColorSlowStyle{Width: 800, Height: 480},
		styles.ColorSlowStyle{Width: 104, Height: 212},
		styles.ColorSlowStyle{Width: 250, Height: 122},
		styles.ColorSlowStyle{Width: 128, Height: 296},
		styles.ColorSlowStyle{Width: 264, Height: 176},
		styles.ColorSlowStyle{Width: 300, Height: 400},

		// ColorFast (color TFT)
		styles.ColorFastStyle{Width: 16, Height: 8},
		styles.ColorFastStyle{Width: 8, Height: 16},
		styles.ColorFastStyle{Width: 160, Height: 80},
		styles.ColorFastStyle{Width: 160, Height: 128},
		styles.ColorFastStyle{Width: 240, Height: 135},
		styles.ColorFastStyle{Width: 240, Height: 240},
		styles.ColorFastStyle{Width: 320, Height: 240},
		styles.ColorFastStyle{Width: 480, Height: 320},
		styles.ColorFastStyle{Width: 800, Height: 480},
		styles.ColorFastStyle{Width: 80, Height: 160},
		styles.ColorFastStyle{Width: 128, Height: 160},
		styles.ColorFastStyle{Width: 135, Height: 240},
		styles.ColorFastStyle{Width: 240, Height: 320},
		styles.ColorFastStyle{Width: 320, Height: 480},
		styles.ColorFastStyle{Width: 480, Height: 800},
		styles.ColorFastStyle{Width: 128, Height: 128},
	)
	return r
}()
