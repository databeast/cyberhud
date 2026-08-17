package borderframe

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/databeast/cyberhud/display/widgets/icons"
)

// clampInnerOffset clamps the inner border offset to [1, 8].
// When offset is zero (Go zero-value / not set), defaults to 4.
func clampInnerOffset(offset int) int {
	if offset == 0 {
		return 4
	}
	if offset < 1 {
		return 1
	}
	if offset > 8 {
		return 8
	}
	return offset
}

// contentInset returns the total content inset per side in pixels.
// If InnerBorder is disabled or suppressed (content area < 8×8), returns 8 (primary border only).
// If InnerBorder is enabled and not suppressed, returns 8 + InnerOffset.
func contentInset(cfg Config) int {
	if !cfg.InnerBorder {
		return tileSize
	}

	offset := clampInnerOffset(cfg.InnerOffset)
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()

	// Content area after primary border (8px each side) and InnerOffset on each side.
	// Per Req 12.6: suppress if content area < 8×8.
	contentW := w - 2*(tileSize+offset)
	contentH := h - 2*(tileSize+offset)

	if contentW < tileSize || contentH < tileSize {
		return tileSize
	}

	return tileSize + offset
}

// renderInnerBorder creates an RGBA image containing the inner border tiles.
// Returns nil if InnerBorder is not enabled, or if the remaining content area
// after both primary border inset and InnerOffset on each side is less than 8×8 pixels.
//
// The inner border is drawn as a complete tile frame (corners + edges) whose outer
// perimeter starts at InnerOffset pixels from the primary border's inner edge.
// Primary border inner edge is at 8px from the outer bounds, so the inner border
// frame origin is at (tileSize + InnerOffset) pixels from each outer edge.
//
// The inner frame dimensions are:
//
//	width  = Bounds.Dx() - 2*(tileSize + InnerOffset)
//	height = Bounds.Dy() - 2*(tileSize + InnerOffset)
//
// Suppression: returns nil if content area < 8×8 (Req 12.6) or if the inner frame
// is too small to render a valid tile border (< 16×16).
func renderInnerBorder(cfg Config) *image.RGBA {
	if !cfg.InnerBorder {
		return nil
	}

	offset := clampInnerOffset(cfg.InnerOffset)
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()

	// Suppression check per Req 12.6: remaining content area after applying both
	// primary border inset (8px each side) AND InnerOffset on each side.
	contentW := w - 2*(tileSize+offset)
	contentH := h - 2*(tileSize+offset)

	if contentW < tileSize || contentH < tileSize {
		return nil
	}

	// Inner frame dimensions (same as content area calculation above).
	innerFrameW := contentW
	innerFrameH := contentH

	// Need at least 2 tiles in each direction for a valid tile frame (corners).
	if innerFrameW < 2*tileSize || innerFrameH < 2*tileSize {
		return nil
	}

	// Create the output image (same size as the full bounds for absolute positioning).
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Inner frame origin within the image.
	originX := tileSize + offset
	originY := tileSize + offset

	// Determine tile prefix for the inner border.
	prefix := resolveInnerPrefix(cfg)

	// Tile grid dimensions for the inner frame.
	innerCols := innerFrameW / tileSize
	innerRows := innerFrameH / tileSize

	// Draw corner tiles.
	drawTileAt(img, prefix+"/corner-tl", originX, originY)
	drawTileAt(img, prefix+"/corner-tr", originX+(innerCols-1)*tileSize, originY)
	drawTileAt(img, prefix+"/corner-bl", originX, originY+(innerRows-1)*tileSize)
	drawTileAt(img, prefix+"/corner-br", originX+(innerCols-1)*tileSize, originY+(innerRows-1)*tileSize)

	// Draw horizontal edge tiles (top and bottom rows, excluding corners).
	for col := 1; col < innerCols-1; col++ {
		x := originX + col*tileSize
		drawTileAt(img, prefix+"/h", x, originY)
		drawTileAt(img, prefix+"/h", x, originY+(innerRows-1)*tileSize)
	}

	// Draw vertical edge tiles (left and right columns, excluding corners).
	for row := 1; row < innerRows-1; row++ {
		y := originY + row*tileSize
		drawTileAt(img, prefix+"/v", originX, y)
		drawTileAt(img, prefix+"/v", originX+(innerCols-1)*tileSize, y)
	}

	// Apply inner border color tinting.
	innerColor := effectiveInnerColor(cfg.InnerColor)
	applyTintToRect(img, innerColor, image.Rect(originX, originY, originX+innerFrameW, originY+innerFrameH))

	return img
}

// resolveInnerPrefix determines the tile set prefix for the inner border.
// Uses InnerTileSet if non-empty, otherwise falls back to the primary border's prefix.
func resolveInnerPrefix(cfg Config) string {
	if cfg.InnerTileSet != "" {
		return cfg.InnerTileSet
	}
	return resolvePrefix(cfg)
}

// effectiveInnerColor returns the color to use for the inner border.
// Defaults to opaque white when InnerColor is zero-value.
func effectiveInnerColor(c color.RGBA) color.RGBA {
	if c == (color.RGBA{}) {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	return c
}

// drawTileAt composites a named icon tile onto dst at the given (x, y) position.
// Gracefully skips missing tiles without panicking.
func drawTileAt(dst *image.RGBA, name string, x, y int) {
	tile, ok := icons.Get(name)
	if !ok {
		return
	}
	draw.Draw(dst, image.Rect(x, y, x+tileSize, y+tileSize), tile, image.Point{}, draw.Over)
}

// applyTintToRect applies foreground color tinting to pixels within the specified
// rectangle that have alpha > 0. This is similar to applyTint but restricted to a rect.
func applyTintToRect(img *image.RGBA, tint color.RGBA, rect image.Rectangle) {
	if img == nil {
		return
	}

	bounds := img.Bounds()
	rect = rect.Intersect(bounds)

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			off := img.PixOffset(x, y)
			a := img.Pix[off+3]
			if a > 0 {
				img.Pix[off+0] = tint.R
				img.Pix[off+1] = tint.G
				img.Pix[off+2] = tint.B
			}
		}
	}
}
