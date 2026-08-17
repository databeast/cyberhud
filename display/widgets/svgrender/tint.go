package svgrender

import (
	"image"
	"image/color"
)

// applyTint applies a channel-wise color multiply tint onto non-transparent
// pixels of the given image. For each pixel where alpha > 0, the RGB channels
// are multiplied by the tint color's corresponding channels (normalized to
// 0–255). The alpha channel is preserved unchanged.
//
// If col is the zero value (R=0, G=0, B=0, A=0), no tinting is performed.
func applyTint(img *image.RGBA, col color.RGBA) {
	if col == (color.RGBA{}) {
		return
	}

	tintR := uint16(col.R)
	tintG := uint16(col.G)
	tintB := uint16(col.B)

	pix := img.Pix
	for i := 0; i < len(pix); i += 4 {
		if pix[i+3] == 0 {
			continue
		}
		pix[i] = uint8(uint16(pix[i]) * tintR / 255)
		pix[i+1] = uint8(uint16(pix[i+1]) * tintG / 255)
		pix[i+2] = uint8(uint16(pix[i+2]) * tintB / 255)
	}
}
