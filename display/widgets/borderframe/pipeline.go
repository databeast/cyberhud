package borderframe

import (
	"image"
	"image/color"
)

// applySegmentRevealMask sets alpha to 0 for all unrevealed tile positions.
// Uses perimeterTilePositions to get ordered positions, then revealStartIndex
// to determine the starting offset. Tiles from startIndex to startIndex+revealedCount
// (wrapping) are visible; all others are zeroed.
func applySegmentRevealMask(img *image.RGBA, cols, rows, revealedCount int, origin Corner) {
	positions := perimeterTilePositions(cols, rows)
	total := len(positions)
	if total == 0 || revealedCount >= total {
		return // All visible or no tiles.
	}

	startIdx := revealStartIndex(cols, rows, origin)

	// Build set of visible tile indices (wrapping).
	visible := make([]bool, total)
	for i := 0; i < revealedCount; i++ {
		idx := (startIdx + i) % total
		visible[idx] = true
	}

	// Zero out alpha for all non-visible tile positions.
	bounds := img.Bounds()
	for i, pos := range positions {
		if visible[i] {
			continue
		}
		// Zero the 8×8 tile at this position.
		for ty := 0; ty < tileSize; ty++ {
			py := pos.Y + ty
			if py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}
			for tx := 0; tx < tileSize; tx++ {
				px := pos.X + tx
				if px < bounds.Min.X || px >= bounds.Max.X {
					continue
				}
				off := img.PixOffset(px, py)
				img.Pix[off+3] = 0
			}
		}
	}
}

// applyGradientTint applies per-tile gradient coloring to the border image.
// For each tile in perimeter order, replaces the RGB of opaque pixels (alpha > 0)
// with the corresponding gradient color.
func applyGradientTint(img *image.RGBA, cols, rows int, colors []color.RGBA) {
	positions := perimeterTilePositions(cols, rows)
	if len(positions) == 0 || len(colors) == 0 {
		return
	}

	bounds := img.Bounds()
	for i, pos := range positions {
		if i >= len(colors) {
			break
		}
		c := colors[i]
		// Apply color to all opaque pixels in this tile.
		for ty := 0; ty < tileSize; ty++ {
			py := pos.Y + ty
			if py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}
			for tx := 0; tx < tileSize; tx++ {
				px := pos.X + tx
				if px < bounds.Min.X || px >= bounds.Max.X {
					continue
				}
				off := img.PixOffset(px, py)
				if img.Pix[off+3] > 0 {
					img.Pix[off+0] = c.R
					img.Pix[off+1] = c.G
					img.Pix[off+2] = c.B
				}
			}
		}
	}
}

// renderScanLine renders the scan line highlight onto the border image.
// For pixels within the scan line range (from scanPos - ScanLength to scanPos),
// adds alpha modulated by a linear gradient (full at leading edge, zero at trailing).
// Uses perimeterTilePositions to map pixel positions along the perimeter.
func renderScanLine(img *image.RGBA, cols, rows int, cfg Config, scanPos float64) {
	// No-op conditions.
	if cfg.ShowBorder != nil && !*cfg.ShowBorder {
		return
	}

	perimeterTiles := perimeterTileCount(cols, rows)
	perimeter := perimeterTiles * tileSize
	if perimeter < 32 {
		return
	}

	// Effective scan length.
	scanLength := cfg.ScanLength
	if scanLength == 0 {
		scanLength = 16
	}
	scanLength = clampInt(scanLength, 1, perimeter)

	// If scan length >= perimeter, render full-perimeter gradient (no advancement).
	if scanLength >= perimeter {
		scanLength = perimeter
	}

	positions := perimeterTilePositions(cols, rows)
	if len(positions) == 0 {
		return
	}

	// Determine tint color for the scan line.
	tint := effectiveTint(cfg.ColorTint)

	perimF := float64(perimeter)

	// For each pixel along the perimeter, check if it's within the scan line range.
	// The leading edge is at scanPos, trailing edge at scanPos - scanLength.
	bounds := img.Bounds()
	for tileIdx, pos := range positions {
		for ty := 0; ty < tileSize; ty++ {
			for tx := 0; tx < tileSize; tx++ {
				// Determine pixel's perimeter position.
				pixelPerimPos := float64(tileIdx*tileSize + ty*0 + tx)
				// Actually, we need to determine pixel position within the tile along
				// the perimeter direction. For simplicity, use tile-center based approach:
				// Each tile occupies tileSize pixels along the perimeter.
				_ = ty // pixels within a tile are at sub-tile positions

				// For the scan line, we treat each pixel column within a tile as one
				// perimeter pixel. The exact mapping depends on which edge we're on.
				// Simplified approach: pixel perimeter position = tileIdx*tileSize + sub-pixel offset.
				// For horizontal edges, offset is tx. For vertical edges, offset is ty.
				var subOffset int
				px := pos.X
				py := pos.Y

				// Determine edge direction for this tile.
				if pos.Y == 0 && pos.X > 0 && pos.X < (cols-1)*tileSize {
					// Top edge (left to right): sub-offset is tx.
					subOffset = tx
					px = pos.X + tx
					py = pos.Y + ty
				} else if pos.X == (cols-1)*tileSize && pos.Y > 0 && pos.Y < (rows-1)*tileSize {
					// Right edge (top to bottom): sub-offset is ty.
					subOffset = ty
					px = pos.X + tx
					py = pos.Y + ty
				} else if pos.Y == (rows-1)*tileSize && pos.X > 0 && pos.X < (cols-1)*tileSize {
					// Bottom edge (right to left): sub-offset is (tileSize-1-tx).
					subOffset = tileSize - 1 - tx
					px = pos.X + tx
					py = pos.Y + ty
				} else if pos.X == 0 && pos.Y > 0 && pos.Y < (rows-1)*tileSize {
					// Left edge (bottom to top): sub-offset is (tileSize-1-ty).
					subOffset = tileSize - 1 - ty
					px = pos.X + tx
					py = pos.Y + ty
				} else {
					// Corner tile: use tx as sub-offset for simplicity.
					subOffset = tx
					px = pos.X + tx
					py = pos.Y + ty
				}

				pixelPerimPos = float64(tileIdx*tileSize + subOffset)

				// Compute distance from leading edge (scanPos).
				// Distance wraps around the perimeter.
				dist := scanPos - pixelPerimPos
				if dist < 0 {
					dist += perimF
				}
				// dist is now how far behind the leading edge this pixel is.

				// Pixel is within scan line if dist < scanLength.
				if dist >= float64(scanLength) {
					continue
				}

				// Alpha factor: 1.0 at leading edge (dist=0), 0.0 at trailing edge (dist=scanLength).
				alphaFactor := 1.0 - dist/float64(scanLength)
				if alphaFactor <= 0 {
					continue
				}

				// Bounds check.
				if px < bounds.Min.X || px >= bounds.Max.X || py < bounds.Min.Y || py >= bounds.Max.Y {
					continue
				}

				// Apply scan line: add highlight (blend toward tint color with alpha).
				off := img.PixOffset(px, py)
				existingA := img.Pix[off+3]

				// Compute highlight alpha contribution.
				highlightA := uint8(float64(tint.A) * alphaFactor)
				if highlightA == 0 {
					continue
				}

				// Additive alpha blend: increase alpha toward full tint.
				newA := int(existingA) + int(highlightA)
				if newA > 255 {
					newA = 255
				}
				img.Pix[off+3] = uint8(newA)

				// Blend RGB toward tint color proportional to highlight contribution.
				if existingA == 0 {
					img.Pix[off+0] = tint.R
					img.Pix[off+1] = tint.G
					img.Pix[off+2] = tint.B
				}
			}
		}
	}
}

// effectiveGlowRadius returns the effective glow radius considering theme defaults.
// Clamps to [0, 32].
func effectiveGlowRadius(cfg Config) int {
	radius := cfg.GlowRadius
	if radius == 0 {
		// Check theme default.
		theme := LookupTheme(cfg.Theme)
		radius = theme.GlowRadius
	}
	return clampGlowRadius(radius)
}
