package borderframe

import (
	"image"
	"image/color"
)

// effectiveTint returns the tint color to use for compositing.
// If the input tint is zero-value (R=0, G=0, B=0, A=0), it returns opaque white.
func effectiveTint(tint color.RGBA) color.RGBA {
	if tint == (color.RGBA{}) {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	return tint
}

// applyTint applies foreground tinting and background fill to a rendered border image.
// For pixels where alpha > 0: RGB is replaced with tint RGB (preserving original alpha).
// For pixels where alpha == 0 and background is non-zero: pixel is filled with background RGBA.
func applyTint(img *image.RGBA, tint color.RGBA, background color.RGBA) {
	if img == nil {
		return
	}

	fg := effectiveTint(tint)
	hasBG := background != (color.RGBA{})

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rowStart := (y - bounds.Min.Y) * img.Stride
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := rowStart + (x-bounds.Min.X)*4

			a := img.Pix[i+3]
			if a > 0 {
				// Opaque pixel from tile: replace RGB with tint, preserve alpha.
				img.Pix[i+0] = fg.R
				img.Pix[i+1] = fg.G
				img.Pix[i+2] = fg.B
			} else if hasBG {
				// Transparent pixel within tile bounds: fill with background.
				img.Pix[i+0] = background.R
				img.Pix[i+1] = background.G
				img.Pix[i+2] = background.B
				img.Pix[i+3] = background.A
			}
		}
	}
}
