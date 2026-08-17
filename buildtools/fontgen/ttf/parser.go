package ttf

import (
	"fmt"
	"image"
	"io"
	"log"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// CodepointRange defines an inclusive range of Unicode codepoints to extract.
type CodepointRange struct {
	Low  rune // Inclusive lower bound
	High rune // Inclusive upper bound
}

// ParseConfig controls TTF rasterization parameters.
type ParseConfig struct {
	Ranges       []CodepointRange // Codepoint ranges to extract
	TargetHeight int              // Target glyph height in pixels (used as ppem)
}

// GlyphData holds the rasterized bitmap for a single glyph.
type GlyphData struct {
	Codepoint rune
	Rows      []uint32 // One uint32 per pixel row; MSB-first bitmask (1=ink)
	Width     int      // Actual glyph pixel width (for monospace verification)
}

// Font holds all extracted glyphs from a TTF source.
type Font struct {
	Glyphs      map[rune]*GlyphData
	GlyphWidth  int // Max glyph width across all extracted glyphs
	GlyphHeight int // TargetHeight from config (consistent with ppem)
}

// Parse reads a TTF/OTF font from r and rasterizes glyphs at the configured
// pixel size for all codepoints within the configured ranges.
//
// The rasterization pipeline:
//  1. Load font via opentype.Parse
//  2. Create a font.Face at cfg.TargetHeight ppem (no hinting for pixel fonts)
//  3. For each codepoint in cfg.Ranges:
//     a. Retrieve the glyph advance and bounds
//     b. Rasterize at target size onto an image.Alpha buffer
//     c. Threshold at 50% alpha: pixels >= 128 become 1 bits in uint32 bitmask (MSB-first)
//     d. Store in the output map
//  4. Determine max width across all rasterized glyphs
//
// Returns error if: font data is invalid, no glyphs are found in the configured
// ranges, or targetHeight <= 0.
func Parse(r io.Reader, cfg ParseConfig) (*Font, error) {
	if cfg.TargetHeight <= 0 {
		return nil, fmt.Errorf("targetHeight must be > 0, got %d", cfg.TargetHeight)
	}
	if len(cfg.Ranges) == 0 {
		return nil, fmt.Errorf("no codepoint ranges specified")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading font data: %w", err)
	}

	otFont, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing font data: %w", err)
	}

	face, err := opentype.NewFace(otFont, &opentype.FaceOptions{
		Size:    float64(cfg.TargetHeight),
		DPI:     72, // 1 point = 1 pixel at 72 DPI
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, fmt.Errorf("creating font face: %w", err)
	}
	defer face.Close()

	result := &Font{
		Glyphs:      make(map[rune]*GlyphData),
		GlyphHeight: cfg.TargetHeight,
	}

	maxWidth := 0
	advances := make(map[int]int) // track advance widths for monospace verification

	for _, cr := range cfg.Ranges {
		for cp := cr.Low; cp <= cr.High; cp++ {
			gd := rasterizeGlyph(face, cp, cfg.TargetHeight)
			if gd == nil {
				continue
			}
			result.Glyphs[cp] = gd
			if gd.Width > maxWidth {
				maxWidth = gd.Width
			}
			advances[gd.Width]++
		}
	}

	if len(result.Glyphs) == 0 {
		return nil, fmt.Errorf("no glyphs found in configured ranges")
	}

	result.GlyphWidth = maxWidth

	// Monospace verification: warn if advance widths vary.
	if len(advances) > 1 {
		log.Printf("ttf: warning: glyph widths vary across extracted glyphs (%d distinct widths); font may not be monospace", len(advances))
	}

	return result, nil
}

// rasterizeGlyph renders a single glyph onto a temporary image and converts
// it to a uint32 bitmask array. Returns nil if the glyph is not present in
// the font or has no visible pixels.
func rasterizeGlyph(face font.Face, cp rune, targetHeight int) *GlyphData {
	// Check if glyph exists by getting its advance width.
	advance, ok := face.GlyphAdvance(cp)
	if !ok {
		return nil
	}

	advancePx := int(math.Ceil(fixedToFloat(advance)))
	if advancePx <= 0 {
		return nil
	}

	// Get glyph bounds to determine the rendering area.
	bounds, _, ok := face.GlyphBounds(cp)
	if !ok {
		return nil
	}

	// Calculate the image dimensions needed.
	// Use the advance width for horizontal extent and targetHeight for vertical.
	imgWidth := advancePx
	// Check if the glyph extends beyond the advance width.
	maxX := int(math.Ceil(fixedToFloat(bounds.Max.X)))
	if maxX > imgWidth {
		imgWidth = maxX
	}

	if imgWidth <= 0 || imgWidth > 32 {
		// Clamp to 32 bits max (uint32 bitmask constraint).
		if imgWidth > 32 {
			imgWidth = 32
		}
		if imgWidth <= 0 {
			return nil
		}
	}

	imgHeight := targetHeight

	// Create an RGBA image for drawing (font.Drawer works with draw.Image).
	rgba := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Determine baseline position.
	// The ascent tells us where the baseline sits from the top.
	metrics := face.Metrics()
	ascent := int(math.Ceil(fixedToFloat(metrics.Ascent)))

	// Draw the glyph using font.Drawer.
	drawer := &font.Drawer{
		Dst:  rgba,
		Src:  image.White,
		Face: face,
		Dot:  fixed.Point26_6{X: 0, Y: fixed.I(ascent)},
	}
	drawer.DrawString(string(cp))

	// Convert RGBA to uint32 bitmask rows with 50% threshold.
	// Since we drew white on black, the R channel gives us coverage.
	rows := make([]uint32, imgHeight)
	hasInk := false
	for y := 0; y < imgHeight; y++ {
		var bits uint32
		for x := 0; x < imgWidth; x++ {
			r, _, _, _ := rgba.At(x, y).RGBA()
			// r is in [0, 65535] range; threshold at 50% (32768)
			if r >= 32768 {
				bits |= 1 << uint(31-x) // MSB-first: bit 31 = leftmost pixel
				hasInk = true
			}
		}
		rows[y] = bits
	}

	if !hasInk {
		// Space-like character with no visible pixels — still valid.
		// Return it with empty rows so it occupies space.
		return &GlyphData{
			Codepoint: cp,
			Rows:      rows,
			Width:     advancePx,
		}
	}

	return &GlyphData{
		Codepoint: cp,
		Rows:      rows,
		Width:     advancePx,
	}
}

// fixedToFloat converts a fixed.Int26_6 to float64.
func fixedToFloat(f fixed.Int26_6) float64 {
	return float64(f) / 64.0
}
