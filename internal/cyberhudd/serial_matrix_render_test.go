package cyberhudd

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/modes/serial"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/runtime/action"
	"github.com/databeast/cyberhud/runtime/ui"
)

type matrixRenderInstance struct {
	hints textlayout.TextHints
}

func (*matrixRenderInstance) ID() string                    { return "serial" }
func (*matrixRenderInstance) Activate()                     {}
func (*matrixRenderInstance) Deactivate()                   {}
func (*matrixRenderInstance) ActionHandler() action.Handler { return nil }
func (*matrixRenderInstance) RenderCacheKey() uint32        { return 0 }
func (i *matrixRenderInstance) BuildView() style.ViewData   { return serial.BuildView(i.hints) }

func TestMatrixSerialRendersThroughRegion(t *testing.T) {
	t.Cleanup(serial.ResetForTest)
	for _, tt := range []struct {
		product string
		w, h    int
		cap     style.Capability
	}{
		{"adafruit-13x9-is31fl3741", 13, 9, style.ColorFast},
		{"adafruit-15x7-charlieplex", 15, 7, style.GrayscaleFast},
		{"adafruit-4120-charlieplex", 16, 8, style.GrayscaleFast},
	} {
		t.Run(tt.product, func(t *testing.T) {
			serial.SetSnapshotForTest(serial.Snapshot{Lines: []string{"ABC"}}, serial.DefaultPolicy())
			bounds := image.Rect(0, 0, tt.w, tt.h)
			surf := surface.New(bounds)
			hints := textlayout.TextHints{PixelWidth: tt.w, PixelHeight: tt.h, Capability: int(tt.cap)}
			r := region.NewRegionWithScreens(tt.product, bounds, surf, []region.ScreenPosition{{
				Bounds: bounds, HintProvider: func() textlayout.TextHints { return hints },
			}}, tt.product, 0, 0)
			r.SetModeFactory(func(_ string, h textlayout.TextHints) (region.ModeInstance, bool) {
				return &matrixRenderInstance{hints: h}, true
			})
			if err := r.SetMode("serial"); err != nil {
				t.Fatal(err)
			}
			if err := ui.NewRegionRenderer(false, nil, nil).Render(r); err != nil {
				t.Fatal(err)
			}
			lit := 0
			for y := 0; y < tt.h; y++ {
				for x := 0; x < tt.w; x++ {
					c := surf.FrameBuffer().RGBAAt(x, y)
					if c.R != 0 || c.G != 0 || c.B != 0 {
						lit++
					}
				}
			}
			if lit == 0 {
				t.Fatal("serial character did not reach the framebuffer")
			}
		})
	}
}
