package source

import (
	fonts "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"

	"github.com/databeast/cyberhud/display/style"
)

// staticSprite wraps a *widgets.Sprite to satisfy the widgets.Renderable interface,
// enabling the gradient background to be added to the Compositor.
type StaticSprite struct {
	s *widgets.Sprite
}

func NewStaticSprite(s *widgets.Sprite) *StaticSprite { return &StaticSprite{s: s} }

func (ss *StaticSprite) RenderFrame() *widgets.Sprite { return ss.s }

// Cached font metrics to avoid repeated face lookups within the same catalog.
var (
	cachedGlyphAdvance  int
	cachedRowHeight     int
	cachedFontResolved  bool
	cachedSpriteFace    fonts.Face
	cachedCatalogWidth  int
	cachedCatalogHeight int
)

// ResolveMatrixFont resolves font metrics using ctx.Face backed by the tier catalog.
// The matrix family includes a compact 10x10 variant for tiny panels and the
// original matrix-code face for larger panels; tier selection picks whichever
// registered face best fits the catalog's width budget.
// Returns glyph advance, row height, and the catalog-validated face for sprite rendering.
func ResolveMatrixFont(ctx style.StyleContext) (glyphAdvance, rowHeight int, face fonts.Face) {
	cat := ctx.FontCatalog()
	if cat.PixelWidth() == 0 {
		// No catalog available — use hardcoded defaults as degraded fallback.
		return 14, 14, nil
	}

	// Cache: if catalog dimensions haven't changed, return cached metrics.
	if cachedFontResolved && cachedCatalogWidth == cat.PixelWidth() && cachedCatalogHeight == cat.PixelHeight() {
		return cachedGlyphAdvance, cachedRowHeight, cachedSpriteFace
	}

	resolved := ctx.Face("matrix", tiercatalog.TierNormal)
	if resolved != nil {
		m := resolved.Metrics()
		cachedGlyphAdvance = m.GlyphAdvance
		cachedRowHeight = m.RowHeight
		cachedSpriteFace = resolved
	} else {
		// Degraded fallback if no font resolves: prefer the panel-provided
		// glyph metrics from hints over a hardcoded size.
		hints := ctx.Hints()
		cachedGlyphAdvance = hints.GlyphAdvance
		cachedRowHeight = hints.RowHeight
		if cachedGlyphAdvance <= 0 {
			cachedGlyphAdvance = 14
		}
		if cachedRowHeight <= 0 {
			cachedRowHeight = 14
		}
		cachedSpriteFace = nil
	}
	cachedFontResolved = true
	cachedCatalogWidth = cat.PixelWidth()
	cachedCatalogHeight = cat.PixelHeight()
	return cachedGlyphAdvance, cachedRowHeight, cachedSpriteFace
}

// ResolveSplashFont returns the largest catalog-validated font suitable for
// the "CYBERHUD" splash text. This intentionally uses "spleen" family at
// TierFullsize — the splash needs a large, dramatic font, not the narrow
// matrix-code font used for rain columns.
//
// This is catalog-validated (no bypass): the returned face satisfies
// GlyphAdvance ≤ maxAdvance for the panel. If TierFullsize doesn't fit
// the splash text width, renderSplash will skip rendering (existing guard).
func ResolveSplashFont(ctx style.StyleContext) fonts.Face {
	if ctx.FontCatalog().PixelWidth() == 0 {
		return nil
	}
	pixelWidth := ctx.Hints().PixelWidth

	// Try largest tier first, fall back to smaller tiers if needed.
	for _, tier := range []tiercatalog.Tier{
		tiercatalog.TierFullsize,
		tiercatalog.TierLarge,
		tiercatalog.TierNormal,
	} {
		face := ctx.Face("spleen", tier)
		if face == nil {
			continue
		}
		// Check if "CYBERHUD" (8 chars) fits the panel width with this face.
		textWidth := len(splashText) * face.Metrics().GlyphAdvance
		if textWidth <= pixelWidth {
			return face
		}
	}
	// Even TierNormal doesn't fit — return it anyway and let renderSplash's
	// existing width guard handle suppression.
	return ctx.Face("spleen", tiercatalog.TierNormal)
}
