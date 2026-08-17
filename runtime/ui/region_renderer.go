package ui

import (
	"image"
	"image/color"
	"log"
	"sync"
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/coordinator"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/surface/tierselect"
	"github.com/databeast/cyberhud/display/widgets"
)

// defaultFontFamily is the family preference used when a declared tier is absent
// from the region's catalog and must be resolved through tierselect.
const defaultFontFamily = "spleen"

// RegionRenderer implements region.Renderer. It renders a region's content by
// calling the Region's active ModeInstance directly.
type RegionRenderer struct {
	monochrome  bool
	warnings    []string
	inputMapper InputMapper

	// modeState is the shared coordinator.State that tracks the desired mode per panel.
	// When non-nil, syncMode detects external mode changes (from console commands)
	// and calls Region.SetMode to activate the new mode instance.
	modeState *coordinator.State
}

// RegionRendererOption configures a RegionRenderer at construction time.
type RegionRendererOption func(*RegionRenderer)

// WithRendererModeState injects the coordinator mode state for mode synchronization.
func WithRendererModeState(ms *coordinator.State) RegionRendererOption {
	return func(rr *RegionRenderer) { rr.modeState = ms }
}

// NewRegionRenderer creates a RegionRenderer with the given configuration.
func NewRegionRenderer(monochrome bool, warnings []string, inputMapper InputMapper, opts ...RegionRendererOption) *RegionRenderer {
	rr := &RegionRenderer{
		monochrome:  monochrome,
		warnings:    warnings,
		inputMapper: inputMapper,
	}
	for _, opt := range opts {
		opt(rr)
	}
	return rr
}

// Render implements region.Renderer. It renders the current frame for the given
// Region by calling the active ModeInstance's BuildView directly.
// If no instance is active (nil), the region is skipped unless syncMode can
// bootstrap it from coordinator.State.
func (rr *RegionRenderer) Render(r *region.Region) error {
	// Sync mode from coordinator.State — detect external mode changes and switch
	// to the new mode instance if needed.
	rr.syncMode(r)

	inst := r.Instance()
	if inst == nil {
		// Bootstrap: region has a mode string but no instance yet (initial allocation
		// only sets the mode field without constructing a ModeInstance). Trigger
		// SetMode to construct and activate the instance.
		if mode := r.CurrentMode(); mode != "" {
			if err := r.SetMode(mode); err != nil {
				return nil // can't construct, skip
			}
			inst = r.Instance()
		}
		if inst == nil {
			return nil // no mode active, skip
		}
	}

	surf := r.Surface()
	if surf == nil {
		return nil
	}

	surf.Clear(colBackground)

	// Build the frame.
	//
	// This used to be followed by a conditional second BuildView call, guarded on
	// the renderer's own row budget differing from the item count, with the stated
	// intent of "re-requesting with the correct maxVisible". It was vestigial: the
	// ModeInstance.BuildView signature takes no arguments, so the mode has no way to
	// observe a row budget and the second call could only return the same view as
	// the first. The row budget now travels the other way, from the style to the
	// renderer, as ViewData.VisibleCount.
	v := inst.BuildView()

	// Log style selection (deduplicated — only emits on change).
	displaymodes.LogStyleSelection(r.Name(), inst.ID(), v.StyleReport)

	// Resolve per-row tier intent into concrete font IDs. Modes declare Tiers;
	// this is the single place those become fonts, so no mode needs its own
	// resolveFont step.
	resolveTierFonts(&v, r.TextHints().Catalog)

	// Establish the baseline face before drawing, then let per-row IDs override it.
	//
	// This pairs with region.applyBaselineGlyphMetrics, which reports the same
	// catalog entry's metrics as hints.GlyphAdvance/RowHeight. A mode that measures
	// with those hints and emits no per-row font IDs must be drawn with the face
	// those metrics describe, or it measures with one font and is drawn with another.
	// Without this the surface would keep font.Default() (spleen-5x8) regardless of
	// what the region advertised.
	rr.setBaselineFont(surf, r.TextHints())

	if len(v.FontIDs) > 0 {
		rr.setFont(surf, v.FontIDs[0])
	}

	// Construct a LayoutCalculator for content positioning from the Region's
	// TextHints, using the same padding the style laid out against.
	//
	// The padding must match. This previously hardcoded PaddingPct: 0 while styles
	// obtain their calculator from StyleContext.Layout(pct) with a padding of their
	// choosing. Any style passing a non-zero padding centred its rows against a
	// narrower content width than the renderer then drew into, so its text sat off
	// centre by the inset. Clock happens to pass 0, which is why the disagreement
	// went unnoticed; ViewData.PaddingPct now carries the value so it cannot.
	bridge := layout.NewLayoutBridge(r.TextHints(), layout.BridgeConfig{
		PaddingPct: v.PaddingPct,
	})

	// Best-effort text fit: for views with no explicit tier or font declarations,
	// select the largest tier whose font allows each row's full text to fit within
	// the content width. Without this, bare-string rows always render at the
	// baseline (TierNormal), which on high-PPI small panels resolves to a
	// physically-sized font that may be far too wide for the pixel budget.
	resolveTextFitFonts(&v, r.TextHints().Catalog, bridge.AvailableContentWidth())

	// Static selects between a centred block and a scrolling cursor list. Both paths
	// honour the same ViewData layout contract; they differ only in windowing and
	// cursor highlighting.
	if v.Static {
		rr.renderStaticItems(surf, v, bridge)
	} else {
		rr.renderItems(surf, v, bridge)
	}

	// Sprites are positioned in absolute panel coordinates by the style, not relative
	// to the text content area. That is deliberate — sprite placement is intentionally
	// freeform, and widgets frequently want to sit outside the text insets (the clock's
	// border frame hugs the panel edge). Do not route Sprite bounds through the
	// LayoutCalculator.
	rr.renderSprites(surf, v.Sprites)

	return nil
}

// resolveTierFonts populates v.FontIDs from v.Tiers using the region's tier
// catalog.
//
// This is the single place in the display system where tier intent becomes a
// concrete font. Modes declare ViewData.Tiers and nothing else; putting the
// resolution here means it cannot be forgotten per-mode, which is exactly what
// happened before it existed — every mode was expected to do this privately, none
// did, and so every row rendered in the surface's default 5x8 face no matter what
// tier the catalog had chosen.
//
// Each tier resolves to its catalog entry's own FontID. This is the important
// detail: it is the same face whose metrics the style used to compute LineOffsets,
// RowHeights and OffsetY (via StyleContext.Entry, which reads the same entry). Any
// other resolution strategy — re-running a family-based search, say — can return a
// different face with a different advance, and the style's centring silently drifts.
// Layout and drawing must agree on one face.
//
// No-ops when the mode resolved fonts itself, when no tiers were declared, or when
// the region has no catalog. In the last case the surface keeps its current face,
// whose metrics Catalog.Entry also reports, so style and renderer still agree.
func resolveTierFonts(v *style.ViewData, catalog tiercatalog.Catalog) {
	if len(v.Tiers) == 0 || len(v.FontIDs) > 0 {
		return
	}
	if catalog.PixelWidth() == 0 {
		return
	}

	fontIDs := make([]string, len(v.Tiers))
	for i, tier := range v.Tiers {
		if entry := catalog.Entry(tier); entry.FontID != "" {
			fontIDs[i] = entry.FontID
			continue
		}
		// Catalog.Entry only yields an empty FontID when the font registry is
		// empty. TrySelect is used rather than tierselect.Select because this runs
		// once per row per frame: Select panics on an unresolvable request, and a
		// mode emitting an unrecognized tier string should degrade visibly rather
		// than terminate the render loop.
		face, reason := tierselect.TrySelect(catalog, tierselect.Request{
			Family: defaultFontFamily,
			Tier:   tier,
		})
		if face == nil {
			logTierFallbackOnce(reason)
			continue
		}
		if reason != "" {
			logTierFallbackOnce(reason)
		}
		fontIDs[i] = face.ID()
	}
	v.FontIDs = fontIDs
}

// resolveTextFitFonts is the renderer's best-effort fallback for ViewData rows
// that carry no explicit tier or font declaration. For each row it selects the
// largest tier whose glyph advance allows the complete text to fit within
// contentWidth, so the row is rendered at the biggest readable size that does
// not cause truncation.
//
// If no tier fits even at the smallest available face (text is genuinely too
// long for the region), the smallest face is chosen so that truncation clips
// as few characters as possible.
//
// No-op when v.Tiers or v.FontIDs is already set, when v.Items is empty, when
// the catalog has no pixel width, or when contentWidth is zero or negative.
func resolveTextFitFonts(v *style.ViewData, catalog tiercatalog.Catalog, contentWidth int) {
	if len(v.Tiers) > 0 || len(v.FontIDs) > 0 {
		return
	}
	if len(v.Items) == 0 || catalog.PixelWidth() == 0 || contentWidth <= 0 {
		return
	}

	// Collect catalog entries in descending tier order (largest font first).
	// Monotonicity enforcement means adjacent tiers can share the same face;
	// deduplicate so we do not redundantly re-check the same advance budget.
	tiers := catalog.Tiers() // ascending order
	type faceEntry struct {
		fontID  string
		advance int
	}
	seen := make(map[string]struct{}, len(tiers))
	// descending: iterate tiers in reverse
	candidates := make([]faceEntry, 0, len(tiers))
	for i := len(tiers) - 1; i >= 0; i-- {
		entry, ok := catalog.Get(tiers[i])
		if !ok || entry.GlyphAdvance <= 0 || entry.FontID == "" {
			continue
		}
		if _, dup := seen[entry.FontID]; dup {
			continue
		}
		seen[entry.FontID] = struct{}{}
		candidates = append(candidates, faceEntry{fontID: entry.FontID, advance: entry.GlyphAdvance})
	}
	if len(candidates) == 0 {
		return
	}

	fontIDs := make([]string, len(v.Items))
	for i, text := range v.Items {
		runeCount := utf8.RuneCountInString(text)
		chosen := candidates[len(candidates)-1].fontID // smallest as fallback
		for _, c := range candidates {
			if runeCount == 0 || c.advance*runeCount <= contentWidth {
				chosen = c.fontID
				break
			}
		}
		fontIDs[i] = chosen
	}
	v.FontIDs = fontIDs
}

// tierFallbackLogged deduplicates font-fallback warnings. Without this, a
// misconfigured tier would emit one log line per row per frame, which on a fast
// panel is tens of thousands of lines a minute and would bury everything else.
var tierFallbackLogged sync.Map

func logTierFallbackOnce(reason string) {
	if reason == "" {
		return
	}
	if _, seen := tierFallbackLogged.LoadOrStore(reason, struct{}{}); !seen {
		log.Printf("region-renderer: %s", reason)
	}
}

// setBaselineFont points the surface at the face whose metrics the region advertises
// as its baseline glyph metrics.
//
// No-op when the region has no catalog: the surface then keeps font.Default(), whose
// metrics tiercatalog.Catalog.Entry also reports in that situation, so the two remain
// in agreement.
func (rr *RegionRenderer) setBaselineFont(surf *surface.Surface, hints textlayout.TextHints) {
	if hints.Catalog.PixelWidth() == 0 {
		return
	}
	if entry := hints.Catalog.Entry(tiercatalog.TierNormal); entry.FontID != "" {
		rr.setFont(surf, entry.FontID)
	}
}

// setFont switches the surface to the given font ID, reporting a failure once.
//
// Surface.SetFontID returns whether the ID was registered, and every call site
// used to discard it with `_ =`. A silent failure leaves the previous face in
// place, so the row is measured with metrics from one font and drawn with another
// — which is exactly the mis-positioning that made this system's text land in the
// wrong place, arriving with no diagnostic at all. Font IDs now originate from the
// region's catalog so they should always resolve; this reports it if that
// assumption ever stops holding.
func (rr *RegionRenderer) setFont(surf *surface.Surface, id string) {
	if id == "" {
		return
	}
	if !surf.SetFontID(id) {
		logTierFallbackOnce("font ID " + id + " is not registered; surface kept its previous face")
	}
}

// --- layout helpers ---

// rowHeight reports a fallback row height from whatever face the surface currently
// holds.
//
// Callers should prefer the style's own per-row heights via staticRowHeight. This
// is order-dependent by nature: the "current" face is whichever row was drawn last,
// including the last row of the previous frame. It survives only for views that
// predate the layout contract and supply no RowHeights.
func (rr *RegionRenderer) rowHeight(surf *surface.Surface) int {
	if rh := surf.FontMetrics().RowHeight; rh > 0 {
		return rh
	}
	return defaultRowHeight
}

// maxVisible is the fallback row budget, used only when a view supplies no
// VisibleCount. It inherits rowHeight's order-dependence and cannot be correct for
// rows at differing tiers, since there is no single row height to divide by.
func (rr *RegionRenderer) maxVisible(surf *surface.Surface) int {
	rh := rr.rowHeight(surf)
	visible := surf.Bounds().Max.Y / rh
	if visible < 1 {
		return 1
	}
	return visible
}

// --- rendering helpers ---

// renderItems draws a scrolling, cursor-navigable list of rows.
//
// It honours the same layout contract as renderStaticItems — per-row heights,
// inter-row spacing, per-row horizontal offsets and the style's own row budget all
// come from ViewData rather than being re-derived here. Only the list-specific
// behaviour differs: rows are windowed from TopRow and the row at Cursor is drawn
// highlighted.
//
// Until this was unified, renderItems ignored LineOffsets, RowHeights, Spacing,
// VisibleCount and OffsetY entirely and positioned every row with one global row
// height. That was not merely latent: the wifi 240x240 and usb styles emit
// Static:false together with mixed tiers, their own computed offsets and their own
// visible count, so every one of those values was being discarded and those views
// rendered left-aligned with overlapping or gapped rows.
//
// OffsetY is treated as it is in the static path — an absolute top for the block
// when non-zero. For a list that scrolls, a style will normally leave it zero.
func (rr *RegionRenderer) renderItems(surf *surface.Surface, v style.ViewData, bridge layout.LayoutCalculator) {
	items := v.Items
	cursor, topRow := v.Cursor, v.TopRow

	m := surf.FontMetrics()
	contentLeft, _ := bridge.ContentOrigin()
	contentWidth := bridge.AvailableContentWidth()

	// Window the list, then bound it by the style's row budget when it supplied
	// one. The fallback measures with whichever face the surface currently holds,
	// which cannot be right for mixed-tier rows.
	start := 0
	if topRow > 0 && topRow < len(items) {
		start = topRow
	}
	visible := items[start:]
	if n := rr.staticVisibleCount(surf, v); n < len(visible) {
		visible = visible[:n]
	}

	_, contentY := bridge.ContentOrigin()
	y := surf.Bounds().Min.Y
	if v.OffsetY > 0 {
		y = contentY + v.OffsetY
	}

	for i, item := range visible {
		absIdx := start + i
		if v.FontIDs != nil && absIdx < len(v.FontIDs) {
			rr.setFont(surf, v.FontIDs[absIdx])
			m = surf.FontMetrics()
		}
		rh := rr.staticRowHeight(v, absIdx, m)
		textY := y + 1
		if m.GlyphHeight > 0 {
			textY = y + max(0, (rh-m.GlyphHeight)/2)
		}
		var bg, fg color.Color
		if absIdx == cursor {
			bg = colHighlight
			fg = rr.textOnBg(colHighlight)
		} else {
			bg = colBackground
			fg = colText
		}
		if v.Colors != nil && absIdx < len(v.Colors) && v.Colors[absIdx] != nil {
			fg = v.Colors[absIdx]
			if rr.monochrome && absIdx == cursor {
				fg = rr.textOnBg(bg)
			}
		}
		surf.DrawRect(image.Rect(contentLeft, y, contentLeft+contentWidth, y+rh), bg)

		baseX := contentLeft
		if v.LineOffsets != nil && absIdx < len(v.LineOffsets) {
			baseX += v.LineOffsets[absIdx]
		}
		surf.DrawText(baseX, textY, item, fg)

		y += rh
		if i < len(visible)-1 {
			y += v.Spacing
		}
	}
}

// renderStaticItems draws a centered, non-scrolling block of rows.
//
// Rows may each use a different font, so vertical advance is per-row rather than a
// single global row height. When the style supplied RowHeights and Spacing, those
// are used verbatim: the style already solved the fit against the panel, and
// re-deriving it here would disagree with the OffsetY it centered the block with.
// Otherwise each row advances by its own font's RowHeight.
func (rr *RegionRenderer) renderStaticItems(surf *surface.Surface, v style.ViewData, bridge layout.LayoutCalculator) {
	visible := v.Items
	if n := rr.staticVisibleCount(surf, v); n < len(visible) {
		visible = visible[:n]
	}
	m := surf.FontMetrics()

	// When OffsetY is provided (non-zero), it is relative to the content origin
	// (the same reference frame as bridge.CenterBlockY and the style's layout
	// arithmetic). Add the content's top inset so the block stays aligned within
	// the panel instead of drifting upward and clipping its first row.
	_, contentY := bridge.ContentOrigin()
	y := surf.Bounds().Min.Y
	if v.OffsetY > 0 {
		y = contentY + v.OffsetY
	}
	contentLeft, _ := bridge.ContentOrigin()
	contentWidth := bridge.AvailableContentWidth()

	for i, item := range visible {
		if v.FontIDs != nil && i < len(v.FontIDs) {
			_ = surf.SetFontID(v.FontIDs[i])
			m = surf.FontMetrics()
		}

		rh := rr.staticRowHeight(v, i, m)

		textY := y + 1
		if m.GlyphHeight > 0 {
			textY = y + max(0, (rh-m.GlyphHeight)/2)
		}
		surf.DrawRect(image.Rect(contentLeft, y, contentLeft+contentWidth, y+rh), colBackground)
		baseX := contentLeft
		if v.LineOffsets != nil && i < len(v.LineOffsets) {
			baseX += v.LineOffsets[i]
		}
		var fg color.Color = colText
		if v.Colors != nil && i < len(v.Colors) && v.Colors[i] != nil {
			fg = v.Colors[i]
		}
		surf.DrawText(baseX, textY, item, fg)

		y += rh
		if i < len(visible)-1 {
			y += v.Spacing
		}
	}
}

// staticRowHeight returns the pixel height of row i, preferring the height the
// style computed and falling back to the row's font metrics.
//
// Preferring the style's value is not redundant with reading the font metrics: the
// style's OffsetY was centred against these specific heights, so substituting a
// different height for any row shifts the block off centre by the difference.
func (rr *RegionRenderer) staticRowHeight(v style.ViewData, i int, m font.Metrics) int {
	if v.RowHeights != nil && i < len(v.RowHeights) && v.RowHeights[i] > 0 {
		return v.RowHeights[i]
	}
	if m.RowHeight > 0 {
		return m.RowHeight
	}
	return defaultRowHeight
}

// staticVisibleCount returns how many rows to draw.
//
// The style's own answer wins when it supplied one. It solved the fit against the
// real per-row font metrics via FitRows, which the renderer cannot reproduce: the
// fallback below divides the region height by a single row height, and with rows at
// different tiers there is no single row height that is correct. On the clock's
// 480x800 view, for instance, the rows are 66, 34 and 34 pixels tall.
//
// The fallback also measures against whatever face the surface currently holds,
// which is the previous frame's last row, making it order-dependent. It is retained
// only for views that predate VisibleCount and still emit zero.
func (rr *RegionRenderer) staticVisibleCount(surf *surface.Surface, v style.ViewData) int {
	if v.VisibleCount > 0 {
		return v.VisibleCount
	}
	return rr.maxVisible(surf)
}

func (rr *RegionRenderer) renderSprites(surf *surface.Surface, sprites []widgets.Sprite) {
	for _, sp := range sprites {
		if sp.Image == nil {
			continue
		}
		if sp.Bounds != (image.Rectangle{}) {
			surf.DrawImageScaled(sp.Image, sp.Bounds)
		} else {
			surf.DrawImage(sp.Image, sp.Position)
		}
	}
}

// textOnBg returns the appropriate text colour for drawing on bg.
// In monochrome mode it inverts to black when the background maps to an
// "on" pixel, ensuring text is always legible on a 1-bit OLED.
func (rr *RegionRenderer) textOnBg(bg color.Color) color.Color {
	if !rr.monochrome {
		return colText
	}
	r, g, b, _ := bg.RGBA()
	// Weighted-luma threshold must match the SH1106 driver's luminanceOn().
	luma := (299*r + 587*g + 114*b) / 1000
	if luma > 0x3000 {
		return colBlack // dark pixels readable on bright/white background
	}
	return colText // white pixels readable on dark background
}

// Ensure RegionRenderer satisfies region.Renderer at compile time.
var _ region.Renderer = (*RegionRenderer)(nil)
