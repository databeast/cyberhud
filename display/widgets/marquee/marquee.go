package marquee

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/surface/fonts"
)

// Direction specifies the scroll axis for a marquee strip.
type Direction int

const (
	Vertical   Direction = iota // Characters scroll top-to-bottom (matrix rain columns)
	Horizontal                  // Characters scroll left-to-right (ticker tape)
)

// CharSource provides characters for each cell position.
// Implementations can be deterministic (fixed string) or random (shuffled rune pool).
type CharSource interface {
	// CharAt returns the character for the given logical cell index.
	CharAt(index int) rune
}

// Config holds all parameters for a marquee strip.
type Config struct {
	Bounds    image.Rectangle // Pixel region this strip occupies
	Direction Direction       // Scroll direction
	Font      font.Face       // Bitmap font for rendering (nil → fonts.Default())
	Source    CharSource      // Character provider
	Colors    []color.RGBA    // Per-cell colors (index 0 = lead cell, gradient tail)
	Speed     float64         // Cells per second
	Phase     float64         // Initial scroll offset in cells (for staggering)
}

// FixedStringSource cycles through a fixed rune slice for each cell position.
// This is the general-purpose CharSource for scrolling text messages.
type FixedStringSource struct {
	Runes []rune
}

// CharAt returns the character at the given logical index, cycling through
// the rune slice using modular indexing. Returns ' ' if Runes is empty.
func (f *FixedStringSource) CharAt(index int) rune {
	if len(f.Runes) == 0 {
		return ' '
	}
	idx := index % len(f.Runes)
	if idx < 0 {
		idx += len(f.Runes)
	}
	return f.Runes[idx]
}
