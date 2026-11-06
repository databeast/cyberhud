package source

import (
	"image"
	"testing"
)

func TestGrayscaleFrameContrast(t *testing.T) {
	for _, size := range []image.Point{{128, 64}, {240, 240}, {320, 240}} {
		img := RenderGrayscaleFrame(size.X, size.Y, 4)
		levels := make(map[uint8]bool)
		for y := 0; y < size.Y; y++ {
			for x := 0; x < size.X; x++ {
				c := img.RGBAAt(x, y)
				if c.R != c.G || c.G != c.B || c.A != 255 {
					t.Fatalf("pixel (%d,%d) is not opaque gray: %v", x, y, c)
				}
				levels[c.R] = true
			}
		}
		bright := false
		for level := range levels {
			bright = bright || level >= 96
		}
		if len(levels) < 8 || !bright {
			t.Fatalf("grayscale frame %v has insufficient levels or contrast: %v", size, levels)
		}
	}
}
