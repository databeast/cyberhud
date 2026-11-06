package attract_cube

import (
	"fmt"
	"image"

	"github.com/databeast/cyberhud/display/modes/attract_cube/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

type frame struct{ seconds float64 }

type cubeStyle struct {
	capability style.Capability
	width      int
	height     int
}

func (s cubeStyle) Name() string {
	return fmt.Sprintf("%s-%dx%d", s.capability, s.width, s.height)
}

func (s cubeStyle) Requirements() style.SurfaceRequirements {
	return style.SurfaceRequirements{MinWidth: s.width, MinHeight: s.height, Capability: s.capability}
}

func (s cubeStyle) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(s.Requirements(), hints)
}

func (s cubeStyle) Build(f frame, policy Policy, ctx style.StyleContext) style.ViewData {
	hints := ctx.Hints()
	w, h := hints.PixelWidth, hints.PixelHeight
	static := hints.Capability == textlayout.CapMonoSlow || hints.Capability == textlayout.CapGrayscaleSlow || hints.Capability == textlayout.CapColorSlow
	if static {
		f.seconds = 0
	}
	var img *image.RGBA
	switch hints.Capability {
	case textlayout.CapMonoFast, textlayout.CapMonoSlow:
		img = source.RenderMonoFrame(w, h, f.seconds)
	case textlayout.CapGrayscaleFast, textlayout.CapGrayscaleSlow:
		img = source.RenderGrayscaleFrame(w, h, f.seconds)
	case textlayout.CapColorFast, textlayout.CapColorSlow:
		_, base, err := source.ParseBaseColor(policy.BaseColor)
		if err != nil {
			return style.ViewData{Items: []string{"cube: invalid base color"}, Static: true}
		}
		img = source.RenderFrame(w, h, f.seconds, base)
	default:
		return style.ViewData{Items: []string{"cube: unsupported display"}, Static: true}
	}
	return style.ViewData{
		Static:  static,
		Sprites: []widgets.Sprite{{Image: img, Bounds: image.Rect(0, 0, w, h), Label: "cube"}},
	}
}

var cubeStyles = func() *style.StyleRegistry[frame, Policy] {
	slow := []image.Point{
		{122, 250}, {176, 264}, {200, 200}, {212, 104}, {296, 128}, {400, 300},
		{480, 800}, {800, 480}, {104, 212}, {250, 122}, {128, 296}, {264, 176}, {300, 400},
	}
	monoFast := []image.Point{
		{16, 8}, {8, 16}, {128, 32}, {128, 64}, {128, 128}, {32, 128}, {64, 128},
	}
	grayFast := []image.Point{
		{15, 7}, {16, 8}, {8, 16}, {128, 64}, {160, 80}, {160, 128}, {240, 135},
		{240, 240}, {320, 240}, {480, 320}, {800, 480}, {80, 160}, {128, 160},
		{135, 240}, {240, 320}, {320, 480}, {480, 800}, {128, 128}, {64, 128},
		{7, 15},
	}
	colorFast := []image.Point{
		{13, 9}, {9, 13}, {16, 8}, {8, 16}, {160, 80}, {160, 128}, {240, 135},
		{240, 240}, {320, 240}, {480, 320}, {800, 480}, {80, 160}, {128, 160},
		{135, 240}, {240, 320}, {320, 480}, {480, 800}, {128, 128},
	}
	var all []style.Style[frame, Policy]
	for _, family := range []struct {
		capability style.Capability
		sizes      []image.Point
	}{
		{style.MonoSlow, slow}, {style.MonoFast, monoFast},
		{style.GrayscaleSlow, slow}, {style.GrayscaleFast, grayFast},
		{style.ColorSlow, slow}, {style.ColorFast, colorFast},
	} {
		for _, size := range family.sizes {
			all = append(all, cubeStyle{capability: family.capability, width: size.X, height: size.Y})
		}
	}
	return style.NewRegistry(all...)
}()

// SnapshotStyles lists the surface-specific styles used by snapshot tooling.
func SnapshotStyles() []style.Style[frame, Policy] {
	return cubeStyles.Enumerate()
}
