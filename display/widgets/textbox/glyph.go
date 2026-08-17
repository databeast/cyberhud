package textbox

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/surface/fonts"
)

// drawGlyph renders a single character onto img at position (x, y) using the
// given font face and foreground color. Pixels outside the image bounds are
// clipped. The position (x, y) refers to the top-left corner of the glyph.
func drawGlyph(img *image.RGBA, face font.Face, ch rune, x, y int, fg color.RGBA) {
	bounds := img.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y

	metrics := face.Metrics()
	glyphWidth := metrics.GlyphWidth
	glyphHeight := metrics.GlyphHeight

	for row := 0; row < glyphHeight; row++ {
		py := y + row
		if py < 0 {
			continue
		}
		if py >= h {
			break
		}

		bits := face.GlyphRow(ch, row)
		for col := 0; col < glyphWidth; col++ {
			px := x + col
			if px >= w {
				break
			}
			if px < 0 {
				continue
			}
			// Bit 31-col being set means pixel at that column is "on".
			if bits&(1<<uint(31-col)) != 0 {
				img.SetRGBA(px, py, fg)
			}
		}
	}
}

// drawText renders a string onto img starting at position (x, y) using the
// given font face and foreground color. It stops rendering characters once the
// next glyph's starting X position would exceed maxX. Returns the total pixel
// width advanced (number of characters rendered × GlyphAdvance).
func drawText(img *image.RGBA, face font.Face, text string, x, y int, fg color.RGBA, maxX int) int {
	metrics := face.Metrics()
	glyphAdvance := metrics.GlyphAdvance

	xCursor := x
	for _, ch := range text {
		if xCursor >= maxX {
			break
		}
		drawGlyph(img, face, ch, xCursor, y, fg)
		xCursor += glyphAdvance
	}

	return xCursor - x
}
