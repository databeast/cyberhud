package attract_cube

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestSurfaceStyles(t *testing.T) {
	styles := cubeStyles.Enumerate()
	counts := make(map[style.Capability]int)
	for _, s := range styles {
		req := s.Requirements()
		counts[req.Capability]++
		hints := textlayout.TextHints{PixelWidth: req.MinWidth, PixelHeight: req.MinHeight, Capability: int(req.Capability)}
		chosen, _ := style.ResolveStyle(cubeStyles, hints, "attract_cube", "")
		if chosen.Name() != s.Name() {
			t.Errorf("%s resolves to %s", s.Name(), chosen.Name())
		}
		view := chosen.Build(frame{seconds: 4}, DefaultPolicy(), style.NewStyleContext(hints))
		if len(view.Sprites) != 1 || view.Sprites[0].Image.Bounds() != image.Rect(0, 0, req.MinWidth, req.MinHeight) {
			t.Fatalf("%s: wrong sprite geometry", s.Name())
		}
		slow := req.Capability == style.MonoSlow || req.Capability == style.GrayscaleSlow || req.Capability == style.ColorSlow
		if view.Static != slow {
			t.Errorf("%s: Static = %v, want %v", s.Name(), view.Static, slow)
		}
	}
	for _, cap := range []style.Capability{style.MonoSlow, style.MonoFast, style.GrayscaleSlow, style.GrayscaleFast, style.ColorSlow, style.ColorFast} {
		if counts[cap] == 0 {
			t.Errorf("no styles for %s", cap)
		}
	}
}

func TestNonStandardDimensionsUsePanelCapabilities(t *testing.T) {
	for _, cap := range []int{textlayout.CapMonoSlow, textlayout.CapMonoFast, textlayout.CapGrayscaleSlow, textlayout.CapGrayscaleFast, textlayout.CapColorSlow, textlayout.CapColorFast} {
		i := &instance{}
		i.SetPanelHints(textlayout.TextHints{PixelWidth: 75, PixelHeight: 55, Capability: cap})
		i.Activate()
		view := i.BuildView()
		if len(view.Sprites) != 1 || view.Sprites[0].Image.Bounds() != image.Rect(0, 0, 75, 55) {
			t.Fatalf("capability %d: unexpected view %+v", cap, view)
		}
		wantStatic := cap == textlayout.CapMonoSlow || cap == textlayout.CapGrayscaleSlow || cap == textlayout.CapColorSlow
		if view.Static != wantStatic {
			t.Errorf("capability %d: Static = %v, want %v", cap, view.Static, wantStatic)
		}
	}
}
