package borderframe

import (
	"image"
	"image/draw"

	"github.com/databeast/cyberhud/display/widgets/icons"
)

// renderCornerAccents draws accent tiles at all four corner positions when
// CornerAccent is enabled. The accent tiles replace the standard corner tiles
// by drawing over them (draw.Over compositing). When CornerAccent is disabled,
// this function is a no-op and CornerFlash is ignored.
//
// Accent tiles are looked up using the naming convention:
//
//	"{prefix}/accent-{tl|tr|bl|br}"
func renderCornerAccents(img *image.RGBA, prefix string, cfg Config) {
	if !cfg.CornerAccent {
		return
	}

	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()
	if w < 16 || h < 16 {
		return
	}

	cols := w / tileSize
	rows := h / tileSize

	// Corner positions (same as standard corners).
	type cornerPos struct {
		name string
		x, y int
	}

	corners := [4]cornerPos{
		{name: prefix + "/accent-tl", x: 0, y: 0},
		{name: prefix + "/accent-tr", x: (cols - 1) * tileSize, y: 0},
		{name: prefix + "/accent-bl", x: 0, y: (rows - 1) * tileSize},
		{name: prefix + "/accent-br", x: (cols - 1) * tileSize, y: (rows - 1) * tileSize},
	}

	for _, c := range corners {
		tile, ok := icons.Get(c.name)
		if !ok {
			continue
		}
		draw.Draw(img, image.Rect(c.x, c.y, c.x+tileSize, c.y+tileSize), tile, image.Point{}, draw.Over)
	}
}

// applyCornerFlashAlpha multiplies the alpha of all pixels in the four 8×8
// corner tile areas by the given alpha factor. This is used by the render
// pipeline when corner flash is active to modulate corner accent visibility.
//
// The four corner areas are:
//   - Top-left:     (0, 0) to (8, 8)
//   - Top-right:    ((cols-1)*8, 0) to (cols*8, 8)
//   - Bottom-left:  (0, (rows-1)*8) to (8, rows*8)
//   - Bottom-right: ((cols-1)*8, (rows-1)*8) to (cols*8, rows*8)
func applyCornerFlashAlpha(img *image.RGBA, cols, rows int, alpha float64) {
	if alpha >= 1.0 {
		return // No modification needed.
	}
	if alpha < 0.0 {
		alpha = 0.0
	}

	// Define the four corner rectangles.
	corners := [4]image.Rectangle{
		image.Rect(0, 0, tileSize, tileSize),
		image.Rect((cols-1)*tileSize, 0, cols*tileSize, tileSize),
		image.Rect(0, (rows-1)*tileSize, tileSize, rows*tileSize),
		image.Rect((cols-1)*tileSize, (rows-1)*tileSize, cols*tileSize, rows*tileSize),
	}

	bounds := img.Bounds()
	for _, rect := range corners {
		// Intersect with image bounds for safety.
		r := rect.Intersect(bounds)
		for y := r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				off := img.PixOffset(x, y)
				a := float64(img.Pix[off+3]) * alpha
				// Round to nearest integer in [0, 255].
				if a < 0 {
					a = 0
				}
				if a > 255 {
					a = 255
				}
				img.Pix[off+3] = uint8(a + 0.5)
			}
		}
	}
}
