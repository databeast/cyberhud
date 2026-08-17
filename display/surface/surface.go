package surface

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/databeast/cyberhud/display/surface/fonts"
)

// Surface wraps an in-memory RGBA buffer with high-level drawing primitives.
// It is the primary abstraction for composing display content.
type Surface struct {
	fb      *image.RGBA
	font    font.Face
	logging bool
	drawLog []DrawCall
}

// New creates a Surface with the given bounds.
func New(bounds image.Rectangle) *Surface {
	return &Surface{
		fb:   image.NewRGBA(bounds),
		font: font.Default(),
	}
}

// NewFromSubImage creates a Surface backed by a sub-image of an existing RGBA buffer.
// The returned Surface uses local coordinates (origin at 0,0) but writes to shared memory.
// Drawing operations on the returned Surface are clipped to the sub-image bounds.
func NewFromSubImage(parent *image.RGBA, rect image.Rectangle) *Surface {
	// Intersect with parent bounds to avoid out-of-range access.
	rect = rect.Intersect(parent.Bounds())

	// Get a sub-image view that shares the same underlying pixel memory.
	sub := parent.SubImage(rect).(*image.RGBA)

	// Create a zero-origin RGBA that shares the same Pix slice and Stride,
	// translating the bounds to local coordinates (0,0 origin).
	localFB := &image.RGBA{
		Pix:    sub.Pix,
		Stride: sub.Stride,
		Rect:   image.Rect(0, 0, rect.Dx(), rect.Dy()),
	}

	return &Surface{
		fb:   localFB,
		font: font.Default(),
	}
}

// Clear fills the frame buffer with color c.
func (s *Surface) Clear(c color.Color) {
	if s.logging {
		s.drawLog = append(s.drawLog, DrawCall{Type: "clear", Rect: s.fb.Bounds(), Color: c})
	}
	draw.Draw(s.fb, s.fb.Bounds(), image.NewUniform(c), image.Point{}, draw.Src)
}

// DrawRect fills rect in the frame buffer.
func (s *Surface) DrawRect(rect image.Rectangle, c color.Color) {
	if s.logging {
		s.drawLog = append(s.drawLog, DrawCall{Type: "rect", Rect: rect, Color: c})
	}
	draw.Draw(s.fb, rect, image.NewUniform(c), image.Point{}, draw.Src)
}

// DrawText writes s starting at pixel (x,y).
// Uses direct pixel-slice access to avoid per-pixel interface dispatch overhead.
func (s *Surface) DrawText(x, y int, text string, fg color.Color) {
	if s.logging {
		s.drawLog = append(s.drawLog, DrawCall{Type: "text", X: x, Y: y, Text: text, Color: fg})
	}
	bounds := s.fb.Bounds()
	m := s.fontMetrics()
	if m.GlyphWidth <= 0 || m.GlyphHeight <= 0 || m.GlyphAdvance <= 0 {
		return
	}

	// Pre-convert fg to RGBA bytes once, avoiding repeated color.Color → RGBA
	// conversion per pixel.
	fr, fgg, fb, fa := fg.RGBA()
	fgR := uint8(fr >> 8)
	fgG := uint8(fgg >> 8)
	fgB := uint8(fb >> 8)
	fgA := uint8(fa >> 8)

	pix := s.fb.Pix
	stride := s.fb.Stride
	minX := bounds.Min.X
	maxX := bounds.Max.X
	minY := bounds.Min.Y
	maxY := bounds.Max.Y
	face := s.activeFont()

	for _, ch := range text {
		// Early-exit if the glyph is entirely off-screen to the right.
		if x >= maxX {
			break
		}
		// Skip glyphs entirely off-screen to the left.
		if x+m.GlyphWidth <= minX {
			x += m.GlyphAdvance
			continue
		}

		for row := 0; row < m.GlyphHeight; row++ {
			py := y + row
			if py < minY || py >= maxY {
				continue
			}
			mask := face.GlyphRow(ch, row)
			if mask == 0 {
				continue
			}
			rowOffset := (py-minY)*stride + (x-minX)*4
			for col := 0; col < m.GlyphWidth; col++ {
				bit := uint32(1 << uint(31-col))
				if mask&bit != 0 {
					px := x + col
					if px >= minX && px < maxX {
						offset := rowOffset + col*4
						pix[offset] = fgR
						pix[offset+1] = fgG
						pix[offset+2] = fgB
						pix[offset+3] = fgA
					}
				}
			}
		}
		x += m.GlyphAdvance
	}
}

// SetFontID configures surface text drawing to a registered font by id.
func (s *Surface) SetFontID(id string) bool {
	if face, ok := font.Get(id); ok {
		s.font = face
		return true
	}
	s.font = font.Default()
	return false
}

func (s *Surface) activeFont() font.Face {
	if s != nil && s.font != nil {
		return s.font
	}
	return font.Default()
}

func (s *Surface) fontMetrics() font.Metrics {
	if face := s.activeFont(); face != nil {
		return face.Metrics()
	}
	return font.Metrics{}
}

// DrawImage composites src onto the framebuffer at dst using source-over blending.
// Nil or zero-size images are no-ops. Pixels outside bounds are clipped.
func (s *Surface) DrawImage(src image.Image, dst image.Point) {
	if src == nil {
		return
	}
	srcBounds := src.Bounds()
	if srcBounds.Dx() <= 0 || srcBounds.Dy() <= 0 {
		return
	}
	// Compute the destination rectangle on the framebuffer.
	dstRect := image.Rectangle{
		Min: dst,
		Max: dst.Add(image.Pt(srcBounds.Dx(), srcBounds.Dy())),
	}
	// draw.Over handles clipping to s.fb bounds and negative offsets naturally.
	draw.Draw(s.fb, dstRect, src, srcBounds.Min, draw.Over)
}

// DrawImageScaled scales src to fill dstRect using nearest-neighbor interpolation,
// then composites onto the framebuffer. Nil/zero-size images or empty rects are no-ops.
func (s *Surface) DrawImageScaled(src image.Image, dstRect image.Rectangle) {
	if src == nil {
		return
	}
	srcBounds := src.Bounds()
	if srcBounds.Dx() <= 0 || srcBounds.Dy() <= 0 {
		return
	}
	if dstRect.Dx() <= 0 || dstRect.Dy() <= 0 {
		return
	}
	scaled := scaleNearestNeighbor(src, dstRect.Dx(), dstRect.Dy())
	draw.Draw(s.fb, dstRect, scaled, image.Point{}, draw.Over)
}

// FrameBuffer returns the underlying frame buffer.
func (s *Surface) FrameBuffer() *image.RGBA {
	return s.fb
}

// Bounds returns the surface's frame-buffer bounds.
func (s *Surface) Bounds() image.Rectangle {
	return s.fb.Bounds()
}

// FontMetrics returns the currently selected font's metrics.
func (s *Surface) FontMetrics() font.Metrics {
	return s.fontMetrics()
}
