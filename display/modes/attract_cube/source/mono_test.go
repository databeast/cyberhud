package source

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"testing"
)

func TestMonoWireframe(t *testing.T) {
	for _, size := range []image.Point{{8, 16}, {128, 32}, {128, 64}, {200, 200}} {
		for _, seconds := range []float64{0, 4, 8} {
			img := RenderMonoFrame(size.X, size.Y, seconds)
			if img.Bounds() != image.Rect(0, 0, size.X, size.Y) {
				t.Fatalf("unexpected dimensions %v", img.Bounds())
			}
			white, black := 0, 0
			for y := 0; y < size.Y; y++ {
				for x := 0; x < size.X; x++ {
					switch pixel := img.RGBAAt(x, y); pixel {
					case color.RGBA{R: 255, G: 255, B: 255, A: 255}:
						white++
					case color.RGBA{A: 255}:
						black++
					default:
						t.Fatalf("non-binary pixel %v at (%d,%d)", pixel, x, y)
					}
				}
			}
			if white == 0 || black == 0 || (size.X >= 128 && white > size.X*size.Y/4) {
				t.Fatalf("unexpected wireframe coverage for %v at %g: %d white, %d black", size, seconds, white, black)
			}
			_, points := project(size.X, size.Y, seconds)
			for _, p := range points {
				x, y := int(math.Round(p.x)), int(math.Round(p.y))
				if img.RGBAAt(x, y).R != 255 {
					t.Fatalf("missing vertex at (%d,%d), size %v time %g", x, y, size, seconds)
				}
			}
		}
	}
	if bytes.Equal(RenderMonoFrame(128, 64, 0).Pix, RenderMonoFrame(128, 64, 4).Pix) {
		t.Fatal("wireframe did not animate")
	}
}
