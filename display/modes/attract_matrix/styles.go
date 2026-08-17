package attract_matrix

import (
	"github.com/databeast/cyberhud/display/modes/attract_matrix/source"
	matrixstyles "github.com/databeast/cyberhud/display/modes/attract_matrix/styles"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// matrixRegistry is the per-mode StyleRegistry for the matrix rain display mode.
var matrixRegistry = style.NewRegistry[source.MatrixSnapshot, source.Policy](
	// MonoSlow (e-paper mono)
	matrixstyles.MonoSlowStyle{Width: 122, Height: 250},
	matrixstyles.MonoSlowStyle{Width: 176, Height: 264},
	matrixstyles.MonoSlowStyle{Width: 200, Height: 200},
	matrixstyles.MonoSlowStyle{Width: 212, Height: 104},
	matrixstyles.MonoSlowStyle{Width: 296, Height: 128},
	matrixstyles.MonoSlowStyle{Width: 400, Height: 300},
	matrixstyles.MonoSlowStyle{Width: 480, Height: 800},
	matrixstyles.MonoSlowStyle{Width: 800, Height: 480},
	matrixstyles.MonoSlowStyle{Width: 104, Height: 212},
	matrixstyles.MonoSlowStyle{Width: 250, Height: 122},
	matrixstyles.MonoSlowStyle{Width: 128, Height: 296},
	matrixstyles.MonoSlowStyle{Width: 264, Height: 176},
	matrixstyles.MonoSlowStyle{Width: 300, Height: 400},

	// MonoFast (OLED mono)
	matrixstyles.MonoStyle{Width: 16, Height: 8},
	matrixstyles.MonoStyle{Width: 8, Height: 16},
	matrixstyles.MonoStyle{Width: 128, Height: 32},
	matrixstyles.MonoStyle{Width: 128, Height: 64},
	matrixstyles.MonoStyle{Width: 128, Height: 128},
	matrixstyles.MonoStyle{Width: 32, Height: 128},
	matrixstyles.MonoStyle{Width: 64, Height: 128},

	// GrayscaleSlow
	matrixstyles.GrayscaleSlowStyle{Width: 122, Height: 250},
	matrixstyles.GrayscaleSlowStyle{Width: 176, Height: 264},
	matrixstyles.GrayscaleSlowStyle{Width: 200, Height: 200},
	matrixstyles.GrayscaleSlowStyle{Width: 212, Height: 104},
	matrixstyles.GrayscaleSlowStyle{Width: 296, Height: 128},
	matrixstyles.GrayscaleSlowStyle{Width: 400, Height: 300},
	matrixstyles.GrayscaleSlowStyle{Width: 480, Height: 800},
	matrixstyles.GrayscaleSlowStyle{Width: 800, Height: 480},
	matrixstyles.GrayscaleSlowStyle{Width: 104, Height: 212},
	matrixstyles.GrayscaleSlowStyle{Width: 250, Height: 122},
	matrixstyles.GrayscaleSlowStyle{Width: 128, Height: 296},
	matrixstyles.GrayscaleSlowStyle{Width: 264, Height: 176},
	matrixstyles.GrayscaleSlowStyle{Width: 300, Height: 400},

	// GrayscaleFast
	matrixstyles.GrayscaleFastStyle{Width: 16, Height: 8},
	matrixstyles.GrayscaleFastStyle{Width: 8, Height: 16},
	matrixstyles.GrayscaleFastStyle{Width: 160, Height: 80},
	matrixstyles.GrayscaleFastStyle{Width: 160, Height: 128},
	matrixstyles.GrayscaleFastStyle{Width: 240, Height: 135},
	matrixstyles.GrayscaleFastStyle{Width: 240, Height: 240},
	matrixstyles.GrayscaleFastStyle{Width: 320, Height: 240},
	matrixstyles.GrayscaleFastStyle{Width: 480, Height: 320},
	matrixstyles.GrayscaleFastStyle{Width: 800, Height: 480},
	matrixstyles.GrayscaleFastStyle{Width: 80, Height: 160},
	matrixstyles.GrayscaleFastStyle{Width: 128, Height: 160},
	matrixstyles.GrayscaleFastStyle{Width: 135, Height: 240},
	matrixstyles.GrayscaleFastStyle{Width: 240, Height: 320},
	matrixstyles.GrayscaleFastStyle{Width: 320, Height: 480},
	matrixstyles.GrayscaleFastStyle{Width: 480, Height: 800},
	matrixstyles.GrayscaleFastStyle{Width: 128, Height: 128},

	// ColorSlow
	matrixstyles.ColorSlowStyle{Width: 122, Height: 250},
	matrixstyles.ColorSlowStyle{Width: 176, Height: 264},
	matrixstyles.ColorSlowStyle{Width: 200, Height: 200},
	matrixstyles.ColorSlowStyle{Width: 212, Height: 104},
	matrixstyles.ColorSlowStyle{Width: 296, Height: 128},
	matrixstyles.ColorSlowStyle{Width: 400, Height: 300},
	matrixstyles.ColorSlowStyle{Width: 480, Height: 800},
	matrixstyles.ColorSlowStyle{Width: 800, Height: 480},
	matrixstyles.ColorSlowStyle{Width: 104, Height: 212},
	matrixstyles.ColorSlowStyle{Width: 250, Height: 122},
	matrixstyles.ColorSlowStyle{Width: 128, Height: 296},
	matrixstyles.ColorSlowStyle{Width: 264, Height: 176},
	matrixstyles.ColorSlowStyle{Width: 300, Height: 400},

	// ColorFast
	matrixstyles.ColorStyle{Width: 16, Height: 8},
	matrixstyles.ColorStyle{Width: 8, Height: 16},
	matrixstyles.ColorStyle{Width: 160, Height: 80},
	matrixstyles.ColorStyle{Width: 160, Height: 128},
	matrixstyles.ColorStyle{Width: 240, Height: 135},
	matrixstyles.ColorStyle{Width: 240, Height: 240},
	matrixstyles.ColorStyle{Width: 320, Height: 240},
	matrixstyles.ColorStyle{Width: 480, Height: 320},
	matrixstyles.ColorStyle{Width: 800, Height: 480},
	matrixstyles.ColorStyle{Width: 80, Height: 160},
	matrixstyles.ColorStyle{Width: 128, Height: 160},
	matrixstyles.ColorStyle{Width: 135, Height: 240},
	matrixstyles.ColorStyle{Width: 240, Height: 320},
	matrixstyles.ColorStyle{Width: 320, Height: 480},
	matrixstyles.ColorStyle{Width: 480, Height: 800},
	matrixstyles.ColorStyle{Width: 128, Height: 128},
)

func registeredStyleNames() []string {
	registered := matrixRegistry.Enumerate()
	names := make([]string, len(registered))
	for i, s := range registered {
		names[i] = s.Name()
	}
	return names
}

func resolveBestStyleName(hints textlayout.TextHints) string {
	s, _ := style.ResolveStyle(matrixRegistry, hints, "attract_matrix", "")
	return s.Name()
}

func resolvePanelType(hints textlayout.TextHints) (mono bool, eink bool) {
	if cachedPanelTypeSet && cachedPanelHintsW == hints.PixelWidth && cachedPanelHintsH == hints.PixelHeight {
		return cachedMono, cachedEink
	}

	s, _ := style.ResolveStyle(matrixRegistry, hints, "attract_matrix", "")
	reqs := s.Requirements()
	mono = reqs.Capability == style.MonoFast
	eink = reqs.Capability == style.MonoSlow

	cachedPanelTypeSet = true
	cachedPanelHintsW = hints.PixelWidth
	cachedPanelHintsH = hints.PixelHeight
	cachedMono = mono
	cachedEink = eink
	return
}

var (
	cachedPanelTypeSet bool
	cachedPanelHintsW  int
	cachedPanelHintsH  int
	cachedMono         bool
	cachedEink         bool
)
