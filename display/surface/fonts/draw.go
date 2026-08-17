package font

import (
	"image"
	"image/color"
)

// DrawGlyph renders a single character onto img at pixel position (x, y)
// using the given Face and foreground color. Pixels outside img bounds
// are clipped silently. Position (x, y) is the top-left corner of the glyph cell.
// If face is nil, DrawGlyph is a no-op.
//
// When fg.A < 255, alpha blending is performed against the existing pixel
// using standard "source over" compositing. When fg.A == 255 (fully opaque),
// the pixel is written directly for maximum performance.
func DrawGlyph(img *image.RGBA, face Face, ch rune, x, y int, fg color.RGBA) {
	if face == nil {
		return
	}
	bounds := img.Bounds()
	metrics := face.Metrics()
	blend := fg.A < 255

	for row := 0; row < metrics.GlyphHeight; row++ {
		py := y + row
		if py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}
		bits := face.GlyphRow(ch, row)
		if bits == 0 {
			continue
		}
		for col := 0; col < metrics.GlyphWidth; col++ {
			px := x + col
			if px < bounds.Min.X || px >= bounds.Max.X {
				continue
			}
			if bits&(1<<uint(31-col)) != 0 {
				if !blend {
					img.SetRGBA(px, py, fg)
				} else {
					// Source-over alpha blend.
					dst := img.RGBAAt(px, py)
					a := uint32(fg.A)
					invA := 255 - a
					r := (uint32(fg.R)*a + uint32(dst.R)*invA) / 255
					g := (uint32(fg.G)*a + uint32(dst.G)*invA) / 255
					b := (uint32(fg.B)*a + uint32(dst.B)*invA) / 255
					outA := a + (uint32(dst.A)*invA)/255
					img.SetRGBA(px, py, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(outA)})
				}
			}
		}
	}
}

// DrawText renders a string onto img starting at (x, y), advancing by
// GlyphAdvance per character. Stops when x exceeds maxX. Returns total
// pixel width advanced.
// If face is nil, DrawText returns 0.
func DrawText(img *image.RGBA, face Face, text string, x, y int, fg color.RGBA, maxX int) int {
	if face == nil {
		return 0
	}
	metrics := face.Metrics()
	xCursor := x
	for _, ch := range text {
		if xCursor >= maxX {
			break
		}
		DrawGlyph(img, face, ch, xCursor, y, fg)
		xCursor += metrics.GlyphAdvance
	}
	return xCursor - x
}
