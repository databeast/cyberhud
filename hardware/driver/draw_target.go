package driver

import (
	"image"
	"image/draw"
)

// DrawTarget is the minimal surface required by a display output.
type DrawTarget interface {
	Bounds() image.Rectangle
	DrawImage(draw.Image) error
}
