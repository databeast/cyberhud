package attract_bokeh

import (
	"github.com/databeast/cyberhud/display/modes/attract_bokeh/source"
	"github.com/databeast/cyberhud/display/modes/attract_bokeh/styles"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// ─── Registry ─────────────────────────────────────────────────────────────────

// bokehRegistry is the per-mode StyleRegistry for the bokeh attract display mode.
// Registration order follows capability ordering: MonoSlow → MonoFast →
// GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast.
var bokehRegistry = func() *style.StyleRegistry[source.BokehFrame, source.Policy] {
	r := style.NewRegistry[source.BokehFrame, source.Policy](
		// MonoSlow (e-paper mono)
		styles.BokehMonoSlowStyle{Width: 122, Height: 250},
		styles.BokehMonoSlowStyle{Width: 176, Height: 264},
		styles.BokehMonoSlowStyle{Width: 200, Height: 200},
		styles.BokehMonoSlowStyle{Width: 212, Height: 104},
		styles.BokehMonoSlowStyle{Width: 296, Height: 128},
		styles.BokehMonoSlowStyle{Width: 400, Height: 300},
		styles.BokehMonoSlowStyle{Width: 480, Height: 800},
		styles.BokehMonoSlowStyle{Width: 800, Height: 480},
		styles.BokehMonoSlowStyle{Width: 104, Height: 212},
		styles.BokehMonoSlowStyle{Width: 250, Height: 122},
		styles.BokehMonoSlowStyle{Width: 128, Height: 296},
		styles.BokehMonoSlowStyle{Width: 264, Height: 176},
		styles.BokehMonoSlowStyle{Width: 300, Height: 400},

		// MonoFast (OLED mono)
		styles.BokehMonoFastStyle{Width: 16, Height: 8},
		styles.BokehMonoFastStyle{Width: 8, Height: 16},
		styles.BokehMonoFastStyle{Width: 128, Height: 32},
		styles.BokehMonoFastStyle{Width: 128, Height: 64},
		styles.BokehMonoFastStyle{Width: 128, Height: 128},
		styles.BokehMonoFastStyle{Width: 32, Height: 128},
		styles.BokehMonoFastStyle{Width: 64, Height: 128},

		// GrayscaleSlow (grayscale e-paper)
		styles.BokehGrayscaleSlowStyle{Width: 122, Height: 250},
		styles.BokehGrayscaleSlowStyle{Width: 176, Height: 264},
		styles.BokehGrayscaleSlowStyle{Width: 200, Height: 200},
		styles.BokehGrayscaleSlowStyle{Width: 212, Height: 104},
		styles.BokehGrayscaleSlowStyle{Width: 296, Height: 128},
		styles.BokehGrayscaleSlowStyle{Width: 400, Height: 300},
		styles.BokehGrayscaleSlowStyle{Width: 480, Height: 800},
		styles.BokehGrayscaleSlowStyle{Width: 800, Height: 480},
		styles.BokehGrayscaleSlowStyle{Width: 104, Height: 212},
		styles.BokehGrayscaleSlowStyle{Width: 250, Height: 122},
		styles.BokehGrayscaleSlowStyle{Width: 128, Height: 296},
		styles.BokehGrayscaleSlowStyle{Width: 264, Height: 176},
		styles.BokehGrayscaleSlowStyle{Width: 300, Height: 400},

		// GrayscaleFast (grayscale LED matrix)
		styles.BokehGrayscaleFastStyle{Width: 16, Height: 8},
		styles.BokehGrayscaleFastStyle{Width: 8, Height: 16},
		styles.BokehGrayscaleFastStyle{Width: 160, Height: 80},
		styles.BokehGrayscaleFastStyle{Width: 160, Height: 128},
		styles.BokehGrayscaleFastStyle{Width: 240, Height: 135},
		styles.BokehGrayscaleFastStyle{Width: 240, Height: 240},
		styles.BokehGrayscaleFastStyle{Width: 320, Height: 240},
		styles.BokehGrayscaleFastStyle{Width: 480, Height: 320},
		styles.BokehGrayscaleFastStyle{Width: 800, Height: 480},
		styles.BokehGrayscaleFastStyle{Width: 80, Height: 160},
		styles.BokehGrayscaleFastStyle{Width: 128, Height: 160},
		styles.BokehGrayscaleFastStyle{Width: 135, Height: 240},
		styles.BokehGrayscaleFastStyle{Width: 240, Height: 320},
		styles.BokehGrayscaleFastStyle{Width: 320, Height: 480},
		styles.BokehGrayscaleFastStyle{Width: 480, Height: 800},
		styles.BokehGrayscaleFastStyle{Width: 128, Height: 128},

		// ColorSlow (color e-paper)
		styles.BokehColorSlowStyle{Width: 122, Height: 250},
		styles.BokehColorSlowStyle{Width: 176, Height: 264},
		styles.BokehColorSlowStyle{Width: 200, Height: 200},
		styles.BokehColorSlowStyle{Width: 212, Height: 104},
		styles.BokehColorSlowStyle{Width: 296, Height: 128},
		styles.BokehColorSlowStyle{Width: 400, Height: 300},
		styles.BokehColorSlowStyle{Width: 480, Height: 800},
		styles.BokehColorSlowStyle{Width: 800, Height: 480},
		styles.BokehColorSlowStyle{Width: 104, Height: 212},
		styles.BokehColorSlowStyle{Width: 250, Height: 122},
		styles.BokehColorSlowStyle{Width: 128, Height: 296},
		styles.BokehColorSlowStyle{Width: 264, Height: 176},
		styles.BokehColorSlowStyle{Width: 300, Height: 400},

		// ColorFast (color TFT)
		styles.BokehColorFastStyle{Width: 16, Height: 8},
		styles.BokehColorFastStyle{Width: 8, Height: 16},
		styles.BokehColorFastStyle{Width: 160, Height: 80},
		styles.BokehColorFastStyle{Width: 160, Height: 128},
		styles.BokehColorFastStyle{Width: 240, Height: 135},
		styles.BokehColorFastStyle{Width: 240, Height: 240},
		styles.BokehColorFastStyle{Width: 320, Height: 240},
		styles.BokehColorFastStyle{Width: 480, Height: 320},
		styles.BokehColorFastStyle{Width: 800, Height: 480},
		styles.BokehColorFastStyle{Width: 80, Height: 160},
		styles.BokehColorFastStyle{Width: 128, Height: 160},
		styles.BokehColorFastStyle{Width: 135, Height: 240},
		styles.BokehColorFastStyle{Width: 240, Height: 320},
		styles.BokehColorFastStyle{Width: 320, Height: 480},
		styles.BokehColorFastStyle{Width: 480, Height: 800},
		styles.BokehColorFastStyle{Width: 128, Height: 128},
	)
	return r
}()

// resolvePanelType determines mono/eink flags from the best-fit style.
func resolvePanelType(hints textlayout.TextHints) (mono bool, eink bool) {
	s, _ := style.ResolveStyle(bokehRegistry, hints, "attract_bokeh", "")
	reqs := s.Requirements()
	switch reqs.Capability {
	case style.MonoFast:
		mono = true
	case style.MonoSlow:
		mono = true
		eink = true
	case style.GrayscaleSlow:
		eink = true
	case style.ColorSlow:
		eink = true
	}
	return
}

// resolveBestStyleName returns the name of the best-fit style.
func resolveBestStyleName(hints textlayout.TextHints) string {
	s, _ := style.ResolveStyle(bokehRegistry, hints, "attract_bokeh", "")
	return s.Name()
}
