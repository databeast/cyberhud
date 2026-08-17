package ggrender

import (
	"fmt"
	"image/color"
)

// Alignment specifies horizontal text alignment relative to a given anchor X coordinate.
type Alignment int

const (
	AlignLeft   Alignment = iota // Left-aligned (default)
	AlignCenter                  // Center-aligned
	AlignRight                   // Right-aligned
)

// DrawText draws a text string at the given position with the specified font, color, and alignment.
// Returns an error if f is nil (font not loaded).
func (c *Canvas) DrawText(text string, x, y float64, f *Font, col color.RGBA, align Alignment) error {
	if f == nil {
		return fmt.Errorf("ggrender: no font set")
	}
	c.ctx.SetFontFace(f.face)
	c.ctx.SetColor(col)

	switch align {
	case AlignCenter:
		c.ctx.DrawStringAnchored(text, x, y, 0.5, 0)
	case AlignRight:
		c.ctx.DrawStringAnchored(text, x, y, 1, 0)
	default: // AlignLeft
		c.ctx.DrawStringAnchored(text, x, y, 0, 0)
	}

	return nil
}

// MeasureText returns the pixel width and height of a text string without drawing.
// Returns an error if f is nil.
func (c *Canvas) MeasureText(text string, f *Font) (width, height float64, err error) {
	if f == nil {
		return 0, 0, fmt.Errorf("ggrender: no font set")
	}
	c.ctx.SetFontFace(f.face)
	w, h := c.ctx.MeasureString(text)
	return w, h, nil
}
