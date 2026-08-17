package svgrender

import (
	"image"

	"github.com/databeast/cyberhud/display/widgets"
)

// Canvas wraps an *image.RGBA buffer with dimension metadata, providing the
// drawable area for a single SVG render pass.
type Canvas struct {
	img    *image.RGBA
	width  int
	height int
}

// NewCanvas creates a Canvas with the given pixel dimensions.
// Returns nil if width or height <= 0.
// The canvas starts with a fully transparent (0,0,0,0) background.
func NewCanvas(width, height int) *Canvas {
	if width <= 0 || height <= 0 {
		return nil
	}
	return &Canvas{
		img:    image.NewRGBA(image.Rect(0, 0, width, height)),
		width:  width,
		height: height,
	}
}

// Image returns the underlying image as an image.Image interface.
func (c *Canvas) Image() image.Image {
	return c.img
}

// ToResult converts the Canvas content into a widgets.Sprite for compositing.
// Labels exceeding 128 runes are truncated.
func (c *Canvas) ToResult(position image.Point, label string) *widgets.Sprite {
	runes := []rune(label)
	if len(runes) > maxLabelLen {
		label = string(runes[:maxLabelLen])
	}
	return &widgets.Sprite{
		Image:    c.img,
		Position: position,
		Label:    label,
	}
}
