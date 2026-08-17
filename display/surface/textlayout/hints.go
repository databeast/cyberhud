package textlayout

import (
	"image"
	"strings"
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/surface/tierselect"
)

const (
	// Default font metrics for the built-in 5×7 bitmap font.
	GlyphWidth   = 5
	GlyphHeight  = 7
	GlyphAdvance = 6
	RowHeight    = 10
)

// Capability constants for TextHints.Capability.
// These mirror style.Capability values, defined here to avoid a circular import
// between textlayout and style. Drivers set TextHints.Capability using these.
const (
	CapMonoSlow      = 0 // 1-bit, slow refresh (e-paper mono)
	CapMonoFast      = 1 // 1-bit, fast refresh (OLED mono)
	CapGrayscaleSlow = 2 // Multi-level luminance, slow refresh (grayscale e-paper)
	CapGrayscaleFast = 3 // Multi-level luminance, fast refresh (grayscale LED matrix)
	CapColorSlow     = 4 // RGB, slow refresh (color e-paper)
	CapColorFast     = 5 // RGB, fast refresh (color TFT)
)

const (
	TickerDirectionVertical = "vertical"
	TickerDirectionNone     = "none"

	LineModeTruncate = "truncate"
	LineModeClip     = "clip"
)

// TextHints describes per-display text layout and scrolling guidance.
// It is screen-level metadata, independent of chipset.
type TextHints struct {
	PixelWidth  int
	PixelHeight int

	GlyphWidth   int
	GlyphHeight  int
	GlyphAdvance int
	RowHeight    int

	SupportsVerticalScroll   bool
	SupportsHorizontalScroll bool
	SupportsAutoScroll       bool
	PreferEventRefresh       bool

	// Capability encodes the panel's hardware capability level as a single ordered int.
	// Values correspond to style.Capability constants (MonoSlow=0 through ColorFast=5).
	// Used by EvaluateFitness for compatibility checks: a style requiring a higher
	// capability than the panel provides is Unsupported.
	Capability int

	DefaultTickerDirection string
	DefaultLineMode        string

	// Catalog provides tier-based font metrics for styles that need layout
	// calculations (GlyphWidth, GlyphHeight, GlyphAdvance, RowHeight per tier).
	// Styles query Catalog.Get(tier) instead of calling BestFontFor or
	// accessing fonts directly.
	Catalog tiercatalog.Catalog

	// PanelProduct is the normalized name of the active panel product
	// (e.g., "waveshare-1.3hat"). Empty string if unknown.
	PanelProduct string

	// ScreenName is the name of the specific screen within a multi-screen
	// panel product (e.g., "main"). Empty for single-screen panels.
	ScreenName string

	// Pixels per inch. Zero means "PPI unknown."
	PPI float64
}

// BoundsProvider is the minimal interface needed to resolve text hints from a display target.
type BoundsProvider interface {
	Bounds() image.Rectangle
}

// TextHintProvider is an optional interface display targets can implement
// to report their own text layout capabilities.
type TextHintProvider interface {
	TextHints() TextHints
}

// ResolveTextHints returns target-specific hints when available,
// falling back to defaults sized to the target bounds.
func ResolveTextHints(target BoundsProvider, hintFn ...func() TextHints) TextHints {
	if target == nil {
		return DefaultTextHints(image.Rect(0, 0, 240, 240))
	}
	if len(hintFn) > 0 && hintFn[0] != nil {
		return Normalize(hintFn[0](), target.Bounds())
	}
	return DefaultTextHints(target.Bounds())
}

// DefaultTextHints returns conservative defaults sized to bounds.
// Defaults match the SH1106 OLED: mono, fast refresh (Capability = MonoFast = 1).
func DefaultTextHints(bounds image.Rectangle) TextHints {
	h := TextHints{
		PixelWidth:               bounds.Dx(),
		PixelHeight:              bounds.Dy(),
		GlyphWidth:               GlyphWidth,
		GlyphHeight:              GlyphHeight,
		GlyphAdvance:             GlyphAdvance,
		RowHeight:                RowHeight,
		SupportsVerticalScroll:   true,
		SupportsHorizontalScroll: true,
		SupportsAutoScroll:       true,
		PreferEventRefresh:       false,
		Capability:               CapMonoFast, // MonoFast — matches style.MonoFast iota value
		DefaultTickerDirection:   TickerDirectionVertical,
		DefaultLineMode:          LineModeTruncate,
	}
	// Auto-build tier catalog from bounds dimensions.
	if bounds.Dx() > 0 && bounds.Dy() > 0 {
		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:  bounds.Dx(),
			PixelHeight: bounds.Dy(),
			MinChars:    10,
		})
		if err == nil {
			h.Catalog = cat
		}
	}
	return h
}

// Normalize fills zero-value fields in h from bounds.
// If h.Catalog is not already populated and the resolved pixel dimensions are
// non-zero, a tier catalog is built from those dimensions so that the returned
// hints are always ready for tier-based font selection, regardless of whether
// the caller used DefaultTextHints or a HintProvider.
func Normalize(h TextHints, bounds image.Rectangle) TextHints {
	h.PixelWidth = firstPositive(h.PixelWidth, bounds.Dx())
	h.PixelHeight = firstPositive(h.PixelHeight, bounds.Dy())
	h.GlyphWidth = firstPositive(h.GlyphWidth, GlyphWidth)
	h.GlyphHeight = firstPositive(h.GlyphHeight, GlyphHeight)
	h.GlyphAdvance = firstPositive(h.GlyphAdvance, GlyphAdvance)
	h.RowHeight = firstPositive(h.RowHeight, RowHeight)
	if h.DefaultTickerDirection == "" {
		h.DefaultTickerDirection = TickerDirectionVertical
	}
	if h.DefaultLineMode == "" {
		h.DefaultLineMode = LineModeTruncate
	}
	if h.Catalog.PixelWidth() <= 0 && h.PixelWidth > 0 && h.PixelHeight > 0 {
		cat, err := tiercatalog.Build(tiercatalog.Params{
			PixelWidth:           h.PixelWidth,
			PixelHeight:          h.PixelHeight,
			MinChars:             10,
			AllowRelaxedMinChars: true,
		})
		if err == nil {
			h.Catalog = cat
		}
	}
	return h
}

// WithFont returns a copy of h with glyph metrics updated from face.
// Pixel dimensions and scroll capabilities are preserved unchanged.
func WithFont(h TextHints, face font.Face) TextHints {
	if face == nil {
		return h
	}
	m := face.Metrics()
	h.GlyphWidth = m.GlyphWidth
	h.GlyphHeight = m.GlyphHeight
	h.GlyphAdvance = m.GlyphAdvance
	h.RowHeight = m.RowHeight
	return h
}

// MaxCharsPerRow returns how many monospace glyphs fit in the usable row width.
func MaxCharsPerRow(h TextHints, horizontalPadding int) int {
	usable := h.PixelWidth - horizontalPadding*2
	if h.GlyphAdvance <= 0 || usable <= 0 {
		return 0
	}
	return usable / h.GlyphAdvance
}

// MaxVisibleRows returns how many rows fit in the usable vertical space.
func MaxVisibleRows(h TextHints, verticalPadding int) int {
	usable := h.PixelHeight - verticalPadding*2
	if h.RowHeight <= 0 || usable <= 0 {
		return 0
	}
	return usable / h.RowHeight
}

func firstPositive(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

// Truncate returns s trimmed of surrounding whitespace. If the trimmed
// result exceeds max runes, it returns the first max-1 runes followed
// by a single ellipsis (U+2026).
func Truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max-1]) + "\u2026"
}

// Face returns a catalog-validated font.Face for the given family and tier.
// The face is guaranteed to satisfy GlyphAdvance ≤ maxAdvance for this panel.
// Returns nil when no catalog is available (degraded mode).
func (h TextHints) Face(family string, tier tiercatalog.Tier) font.Face {
	if h.Catalog.PixelWidth() <= 0 {
		return nil
	}
	return tierselect.Select(h.Catalog, tierselect.Request{
		Family: family,
		Tier:   tier,
	})
}

// AllFaces returns one catalog-validated face per tier (Small, Normal, Large, Fullsize)
// for the given family. Useful for diagnostic modes that enumerate available faces.
// Returns nil when no catalog is available.
func (h TextHints) AllFaces(family string) []font.Face {
	if h.Catalog.PixelWidth() <= 0 {
		return nil
	}
	tiers := h.Catalog.Tiers()
	faces := make([]font.Face, len(tiers))
	for i, tier := range tiers {
		faces[i] = tierselect.Select(h.Catalog, tierselect.Request{
			Family: family,
			Tier:   tier,
		})
	}
	return faces
}
