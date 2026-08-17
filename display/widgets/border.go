package widgets

import (
	"image"
	"image/color"
)

// DrawBorder draws a 1px border on the outer edges of img using fg.
// It sets the top row, bottom row, left column, and right column pixels.
func DrawBorder(img *image.RGBA, width, height int, fg color.RGBA) {
	for x := 0; x < width; x++ {
		img.SetRGBA(x, 0, fg)
		img.SetRGBA(x, height-1, fg)
	}
	for y := 0; y < height; y++ {
		img.SetRGBA(0, y, fg)
		img.SetRGBA(width-1, y, fg)
	}
}
