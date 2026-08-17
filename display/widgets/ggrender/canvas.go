package ggrender

import (
	"image"

	"github.com/databeast/cyberhud/display/widgets"
	"github.com/fogleman/gg"
)

// Canvas wraps a gg.Context with dimension metadata, providing the drawable
// area for a single widget render pass.
type Canvas struct {
	ctx    *gg.Context
	width  int
	height int
}

// NewCanvas creates a Canvas with the given pixel dimensions.
// Returns nil if width or height <= 0.
// The canvas starts with a fully transparent background.
func NewCanvas(width, height int) *Canvas {
	if width <= 0 || height <= 0 {
		return nil
	}
	ctx := gg.NewContext(width, height)
	// gg.NewContext already initializes with transparent (zero) pixels
	return &Canvas{
		ctx:    ctx,
		width:  width,
		height: height,
	}
}

// Image returns the underlying image.Image from the gg context.
func (c *Canvas) Image() image.Image {
	return c.ctx.Image()
}

// ToResult converts the Canvas content into a widgets.Sprite for compositing.
// Labels exceeding 128 runes are truncated.
func (c *Canvas) ToResult(position image.Point, label string) *widgets.Sprite {
	runes := []rune(label)
	if len(runes) > maxLabelLen {
		label = string(runes[:maxLabelLen])
	}
	return &widgets.Sprite{
		Image:    c.ctx.Image(),
		Position: position,
		Label:    label,
	}
}
