package borderframe

import (
	"image"
	"math"
)

// applyOpacity multiplies all pixel alpha values by the effective opacity factor.
// Opacity semantics:
//   - nil: unset, full opacity, no modification (early return).
//   - *0.0: fully transparent, all alpha set to 0.
//   - *1.0: full opacity, no modification (early return).
//   - *v where 0 < v < 1: each pixel alpha = round(alpha * v).
//   - Values outside [0.0, 1.0] are clamped.
func applyOpacity(img *image.RGBA, opacity *float64) {
	if opacity == nil {
		return
	}

	// Clamp to [0.0, 1.0].
	v := *opacity
	if v < 0.0 {
		v = 0.0
	} else if v > 1.0 {
		v = 1.0
	}

	// Full opacity: no modification needed.
	if v == 1.0 {
		return
	}

	bounds := img.Bounds()

	// Fully transparent: zero all alpha bytes.
	if v == 0.0 {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rowStart := (y - bounds.Min.Y) * img.Stride
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				i := rowStart + (x-bounds.Min.X)*4
				img.Pix[i+3] = 0
			}
		}
		return
	}

	// Partial opacity: multiply each pixel's alpha by the factor.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rowStart := (y - bounds.Min.Y) * img.Stride
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := rowStart + (x-bounds.Min.X)*4
			a := img.Pix[i+3]
			if a == 0 {
				continue
			}
			img.Pix[i+3] = uint8(math.Round(float64(a) * v))
		}
	}
}
