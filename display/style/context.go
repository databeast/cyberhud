package style

import (
	"github.com/databeast/cyberhud/display/style/layout"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/surface/tierselect"
)

// StyleContext is the parameter for Style.Build(). It provides region
// constraints (TextHints), a pre-built LayoutBridge for backward compatibility,
// capability flags, and font catalog access. Styles that need to control their
// own layout construct a new LayoutBridge from Hints() with their own BridgeConfig.
type StyleContext struct {
	hints                    textlayout.TextHints
	catalog                  tiercatalog.Catalog
	capability               Capability
	supportsVerticalScroll   bool
	supportsHorizontalScroll bool
}

// Layout returns the encapsulated LayoutBridge for spatial queries.
// Styles that need custom padding should construct their own bridge from
// Hints() rather than relying on this pre-built one.
func (ctx StyleContext) Layout(paddingPercent int) layout.LayoutCalculator {
	return layout.NewLayoutBridge(ctx.hints, layout.BridgeConfig{PaddingPct: paddingPercent})
}

// Hints returns the region's TextHints, giving the style access to the raw
// panel constraints (dimensions, glyph metrics, capability). Styles use this
// to construct their own LayoutBridge with an appropriate BridgeConfig.
func (ctx StyleContext) Hints() textlayout.TextHints {
	return ctx.hints
}

// FontCatalog returns the encapsulated Catalog for font tier lookups.
//
// Prefer [StyleContext.Entry] for layout arithmetic. FontCatalog remains for
// styles that need catalog-level information such as MaxAdvance or Relaxed.
func (ctx StyleContext) FontCatalog() tiercatalog.Catalog {
	return ctx.catalog
}

// Entry returns the font metrics a style should lay out tier with, and cannot fail.
//
// This replaces the following block, which every tiered style used to carry:
//
//	entry, ok := ctx.FontCatalog().Get(tier)
//	if !ok {
//	    entry = tiercatalog.Entry{
//	        GlyphAdvance: hints.GlyphAdvance,
//	        RowHeight:    hints.RowHeight,
//	    }
//	}
//
// That fallback substituted textlayout's 6x10 defaults, which belong to no
// registered font, so the style measured with metrics the renderer would never
// draw with and its text landed wrong. See [tiercatalog.Catalog.Entry] for the
// full history.
//
// The returned Entry names its face in FontID. The renderer draws each row with
// exactly that face, so a style using Entry is guaranteed that the glyph advance
// it centred with is the advance actually rendered.
func (ctx StyleContext) Entry(tier tiercatalog.Tier) tiercatalog.Entry {
	return ctx.catalog.Entry(tier)
}

// Face returns the concrete font face for tier, or nil in degraded mode where no
// catalog was built for this region.
//
// Styles normally need [StyleContext.Entry] (metrics for arithmetic) rather than a
// face. Face is for styles that rasterize text themselves into a sprite instead of
// emitting rows for the renderer to draw.
//
// Never panics: unlike tierselect.Select, this routes through
// [tierselect.TrySelect], because Build is invoked during frame production.
func (ctx StyleContext) Face(family string, tier tiercatalog.Tier) font.Face {
	if ctx.catalog.PixelWidth() <= 0 {
		return nil
	}
	face, _ := tierselect.TrySelect(ctx.catalog, tierselect.Request{
		Family: family,
		Tier:   tier,
	})
	return face
}

// Cap returns the panel's hardware capability level.
func (ctx StyleContext) Cap() Capability {
	return ctx.capability
}

// VerticalScroll reports whether the panel supports vertical scrolling.
func (ctx StyleContext) VerticalScroll() bool {
	return ctx.supportsVerticalScroll
}

// HorizontalScroll reports whether the panel supports horizontal scrolling.
func (ctx StyleContext) HorizontalScroll() bool {
	return ctx.supportsHorizontalScroll
}

// NewStyleContext constructs a StyleContext from the given font Catalog, and TextHints.
func NewStyleContext(hints textlayout.TextHints) StyleContext {
	return StyleContext{
		hints:                    hints,
		catalog:                  hints.Catalog,
		capability:               Capability(hints.Capability),
		supportsVerticalScroll:   hints.SupportsVerticalScroll,
		supportsHorizontalScroll: hints.SupportsHorizontalScroll,
	}
}
