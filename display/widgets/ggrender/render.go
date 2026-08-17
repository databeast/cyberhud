package ggrender

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/widgets"
)

// Config holds parameters for a gg-based reference widget.
type Config struct {
	Bounds image.Rectangle // Pixel region for rendering
	Label  string          // Sprite identification (max 128 chars)
	Color  color.RGBA      // Primary fill color
}

// Render produces a reference widget image. Returns nil for invalid bounds
// (zero/negative dimensions) or bounds smaller than MinBoundsWidth × MinBoundsHeight.
func Render(cfg Config) *widgets.Sprite {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()
	if w <= 0 || h <= 0 {
		return nil
	}
	if w < MinBoundsWidth || h < MinBoundsHeight {
		return nil
	}

	c := NewCanvas(w, h)
	if c == nil {
		return nil
	}

	// Draw a simple filled rectangle as the reference widget content
	c.FillRect(0, 0, float64(w), float64(h), cfg.Color)

	// Convert to Sprite with position from Bounds.Min
	return c.ToResult(cfg.Bounds.Min, cfg.Label)
}
