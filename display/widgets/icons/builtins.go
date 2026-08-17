package icons

import (
	"image"
	"image/color"
)

func init() {
	Register("wifi", makeAlpha8x8([]byte{
		0b00000000,
		0b01111110,
		0b00000000,
		0b00111100,
		0b00000000,
		0b00011000,
		0b00011000,
		0b00000000,
	}))

	Register("bluetooth", makeAlpha8x8([]byte{
		0b00010000,
		0b00010100,
		0b01010010,
		0b00111100,
		0b00010000,
		0b00111100,
		0b01010010,
		0b00010100,
	}))

	Register("battery", makeAlpha8x8([]byte{
		0b00111100,
		0b01111110,
		0b01000010,
		0b01000010,
		0b01000010,
		0b01111110,
		0b01111110,
		0b01111110,
	}))

	Register("error", makeAlpha8x8([]byte{
		0b00000000,
		0b01000010,
		0b00100100,
		0b00011000,
		0b00011000,
		0b00100100,
		0b01000010,
		0b00000000,
	}))

	Register("check", makeAlpha8x8([]byte{
		0b00000000,
		0b00000001,
		0b00000010,
		0b00000100,
		0b01001000,
		0b00110000,
		0b00100000,
		0b00000000,
	}))

	Register("circle-filled", makeAlpha8x8([]byte{
		0b00000000,
		0b00111100,
		0b01111110,
		0b01111110,
		0b01111110,
		0b01111110,
		0b00111100,
		0b00000000,
	}))

	Register("circle-hollow", makeAlpha8x8([]byte{
		0b00000000,
		0b00111100,
		0b01000010,
		0b01000010,
		0b01000010,
		0b01000010,
		0b00111100,
		0b00000000,
	}))
}

// makeAlpha8x8 creates an 8x8 *image.Alpha from a slice of 8 bytes,
// where each bit set to 1 produces an opaque (0xFF) alpha pixel and
// each bit set to 0 produces a transparent (0x00) alpha pixel.
// Bits are read MSB-first (bit 7 = x=0, bit 0 = x=7).
func makeAlpha8x8(rows []byte) *image.Alpha {
	img := image.NewAlpha(image.Rect(0, 0, 8, 8))
	for y, row := range rows {
		for x := 0; x < 8; x++ {
			if row&(1<<uint(7-x)) != 0 {
				img.SetAlpha(x, y, color.Alpha{A: 0xFF})
			}
		}
	}
	return img
}
