package styles

// Core layout assembly for systemd styles.
//
// Every per-resolution style in this package is declared as a def value in
// styles_defs.go: a name, the surface requirements it targets, and a Params
// block of hand-tweakable knobs. Edit Params for display-specific tweaks, or
// attach Params.BuildFn for a fully bespoke layout.

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/systemd/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/gradient"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
)

// Params holds the hand-tweakable knobs for a single per-resolution style.
// The zero value renders centered default boot-status rows.
type Params struct {
	// ColorText renders rows in white at TierNormal for color skeleton styles.
	ColorText bool

	// Static marks the ViewData static for e-ink/slow-refresh layouts that should
	// not be continuously refreshed.
	Static bool

	// SingleLabel renders only source.BootStatusLabel centered on the panel.
	SingleLabel bool

	// Summary renders the compact 128x32 mono summary, using "Boot OK" once boot
	// completes and truncating the loading target while boot is in progress.
	Summary bool

	// Gradient renders a boot-progress gradient background with the centered
	// boot label overlaid on top.
	Gradient bool

	// EinkSuppressed applies the e-ink suppression compositor.
	EinkSuppressed bool
	// ContentMaxRows uses content height / row height instead of bridge.MaxVisibleRows().
	ContentMaxRows bool

	// BuildFn, when set, replaces the core layout entirely for this style.
	BuildFn func(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData
}

// def is the declaration form of a systemd style: name + surface requirements
// + layout tweaks. It implements style.Style[source.Snapshot, source.Policy].
type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.Snapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(snap source.Snapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(snap, pol, ctx, d)
	}
	if d.p.Gradient {
		return gradientBuild(snap, pol, ctx, d)
	}
	if d.p.Summary {
		return summaryBuild(snap, ctx)
	}
	if d.p.SingleLabel {
		return singleLabelBuild(snap, ctx, d)
	}
	return rowsBuild(snap, ctx, d)
}

func rowsBuild(snap source.Snapshot, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(0)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(no data)"}, Static: d.p.Static}
	}

	var sprites []widgets.Sprite
	if d.p.EinkSuppressed {
		suppCtx := widgets.SuppressionContext{
			AvailableWidth:  bridge.AvailableContentWidth(),
			AvailableHeight: bridge.AvailableContentHeight(),
			IsEink:          true,
		}
		comp := widgets.NewCompositor(suppCtx, widgets.SuppressOnEink())
		sprites = comp.Sprites()
	}

	items := source.BuildDefaultItems(snap)

	// Every row is drawn at TierNormal — the boot status is a uniform block. Look
	// up that face once; its metrics measure the rows and name the drawn font.
	entry := ctx.Entry(tiercatalog.TierNormal)
	m := layout.NewMetrics(entry.GlyphAdvance, entry.RowHeight, entry.FontID)

	// ContentMaxRows is the e-ink pre-clamp: cap rows by raw content height so the
	// static poster never overflows. The remaining fit decision is FitRows' job via
	// LayoutRows below.
	if d.p.ContentMaxRows && entry.RowHeight > 0 {
		maxRows := bridge.AvailableContentHeight() / entry.RowHeight
		if maxRows > 0 && len(items) > maxRows {
			items = items[:maxRows]
		}
	}

	// Solve the whole block in one call, exactly as the clock does: per-row
	// centering, row heights, vertical offset, inter-row spacing and the visible
	// count all come from LayoutRows/FitRows. Hand the complete RowsResult to the
	// renderer — omitting RowHeights/Spacing/VisibleCount makes it re-derive a row
	// budget from the surface's current face, which collapses to zero rows on the
	// mono path and renders blank. (This is why the hand-rolled version drew
	// nothing on mono-fast.)
	rowInputs := make([]layout.RowInput, len(items))
	for i, item := range items {
		rowInputs[i] = layout.RowInput{TextLen: len(item), Font: m}
	}
	result := bridge.LayoutRows(rowInputs)

	// Bind every row to the measured tier so the renderer draws the face the rows
	// were centered against, instead of letting resolveTextFitFonts re-pick a wider
	// font (which leans the block right and clips the last glyph).
	tiers := make([]tiercatalog.Tier, len(items))
	for i := range items {
		tiers[i] = tiercatalog.TierNormal
	}

	// ColorText paints the rows white on color skeleton styles; independent of the
	// tier binding above.
	var colors []color.Color
	if d.p.ColorText {
		white := color.RGBA{255, 255, 255, 255}
		colors = make([]color.Color, len(items))
		for i := range items {
			colors[i] = white
		}
	}

	return style.ViewData{
		Items:        items,
		Colors:       colors,
		Tiers:        tiers,
		LineOffsets:  result.Offsets,
		OffsetY:      result.OffsetY,
		RowHeights:   result.RowHeights,
		Spacing:      result.Spacing,
		VisibleCount: result.VisibleCount,
		PaddingPct:   0,
		Sprites:      sprites,
		// Boot status is a non-interactive readout: render it as a centered static
		// block, never the scrolling cursor list. Routing through the cursor path
		// (Static:false) made the single row land at the default Cursor 0 and draw
		// highlighted — invisible on monochrome panels. It is categorically static
		// regardless of the per-def Static knob.
		Static: true,
		Cursor: -1,
	}
}

// labelFontFamily is the font family the single-label builders rasterize their
// centered boot-status text with. It matches the region's baseline family, so a
// panel with no catalog still resolves a sensible face.
const labelFontFamily = "spleen"

// centeredLabelSprite rasterizes text as a single sprite centered in the content
// area, using the face the style itself resolves.
//
// This is the correct-by-construction alternative to emitting a bare Items row.
// A bare row carries no tier, so the renderer's resolveTextFitFonts fallback
// re-picks the largest font that fits the width — a different, usually larger
// face than the one the style measured CenterX against. The style centers for a
// narrow face and the renderer draws a wide one, so the label leans right and
// clips (the "BOOT COMP" overflow). Rasterizing here binds measurement and
// drawing to one face: the sprite owns its glyphs, and the renderer composites
// it verbatim without re-resolving anything.
//
// Returns a zero Sprite (nil Image) when no face is available or the text does
// not fit; callers should fall back to an Items row in that degraded case.
func centeredLabelSprite(bridge layout.LayoutCalculator, face font.Face, text string, fg color.RGBA) (widgets.Sprite, bool) {
	if face == nil {
		return widgets.Sprite{}, false
	}
	cw := bridge.AvailableContentWidth()
	ch := bridge.AvailableContentHeight()
	if cw <= 0 || ch <= 0 {
		return widgets.Sprite{}, false
	}
	ox, oy := bridge.ContentOrigin()
	m := face.Metrics()
	rowH := m.GlyphHeight
	if rowH <= 0 {
		rowH = m.RowHeight
	}
	// Vertically center a single row band within the content height; textlabel
	// centers the glyphs within the band it is given.
	bandTop := oy + max(0, (ch-rowH)/2)
	sprite := textlabel.Render(textlabel.Config{
		Text:       text,
		Bounds:     image.Rect(ox, bandTop, ox+cw, bandTop+rowH),
		Font:       face,
		Alignment:  textlabel.Center,
		Foreground: fg,
	})
	if sprite == nil {
		return widgets.Sprite{}, false
	}
	return *sprite, true
}

// labelBoundsText truncates text to the glyphs that fit the content width for the
// given face, so a label can never be wider than the panel.
func labelBoundsText(bridge layout.LayoutCalculator, face font.Face, text string) string {
	if face == nil {
		return text
	}
	adv := face.Metrics().GlyphAdvance
	if adv <= 0 {
		return text
	}
	maxChars := bridge.AvailableContentWidth() / adv
	if maxChars > 0 && len(text) > maxChars {
		return text[:maxChars]
	}
	return text
}

func singleLabelBuild(snap source.Snapshot, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(0)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(no data)"}}
	}

	label := source.BootStatusLabel(snap)
	face := ctx.Face(labelFontFamily, tiercatalog.TierNormal)
	label = labelBoundsText(bridge, face, label)

	white := color.RGBA{255, 255, 255, 255}
	if sprite, ok := centeredLabelSprite(bridge, face, label, white); ok {
		return style.ViewData{
			Sprites: []widgets.Sprite{sprite},
			Static:  d.p.Static,
		}
	}

	// Degraded fallback: no face resolved, emit the row for the renderer to draw.
	entry := ctx.Entry(tiercatalog.TierNormal)
	offsets := []style.TTYX{bridge.CenterX(len(label), layout.NewMetrics(entry.GlyphAdvance, entry.RowHeight, entry.FontID))}
	offsetY := bridge.CenterBlockY([]int{entry.RowHeight}, 0)
	return style.ViewData{
		Items:       []string{label},
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      d.p.Static,
	}
}

func summaryBuild(snap source.Snapshot, ctx style.StyleContext) style.ViewData {
	bridge := ctx.Layout(0)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(no data)"}}
	}

	var label string
	if snap.BootComplete {
		label = "Boot OK"
	} else {
		label = source.BootStatusLabel(snap)
	}

	face := ctx.Face(labelFontFamily, tiercatalog.TierNormal)
	label = labelBoundsText(bridge, face, label)

	white := color.RGBA{255, 255, 255, 255}
	if sprite, ok := centeredLabelSprite(bridge, face, label, white); ok {
		return style.ViewData{Sprites: []widgets.Sprite{sprite}}
	}

	entry := ctx.Entry(tiercatalog.TierNormal)
	offsets := []style.TTYX{bridge.CenterX(len(label), layout.NewMetrics(entry.GlyphAdvance, entry.RowHeight, entry.FontID))}
	offsetY := bridge.CenterBlockY([]int{entry.RowHeight}, 0)
	return style.ViewData{
		Items:       []string{label},
		LineOffsets: offsets,
		OffsetY:     offsetY,
	}
}

func gradientBuild(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(0)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(no data)"}}
	}

	frac := bootFraction(snap)
	accent := gradientAccent(pol.ColorAccent, snap.BootComplete)

	suppCtx := widgets.SuppressionContext{
		AvailableWidth:  bridge.AvailableContentWidth(),
		AvailableHeight: bridge.AvailableContentHeight(),
	}
	comp := widgets.NewCompositor(suppCtx)
	comp.Add(gradient.New(gradient.Config{
		Style:  gradient.Linear,
		Angle:  0,
		Bounds: image.Rect(0, 0, bridge.AvailableContentWidth(), bridge.AvailableContentHeight()),
		Stops: []gradient.ColorStop{
			{Position: 0.0, Color: color.RGBA{0, 0, 0, 255}},
			{Position: frac, Color: accent},
		},
	}))

	// The gradient is a full-panel backdrop: sink it beneath the label so the text
	// reads on top. Without Z: -1 the opaque gradient composites over the label
	// (the default sprite depth draws above lower depths), burying it entirely.
	sprites := comp.Sprites()
	for i := range sprites {
		sprites[i].Z = -1
	}

	// The boot label is a sprite at the default depth (0), above the Z:-1 gradient.
	// Rasterizing it here (rather than emitting an Items row) binds its font to the
	// one the style measured, so it centers correctly over the gradient.
	label := source.BootStatusLabel(snap)
	face := ctx.Face(labelFontFamily, tiercatalog.TierNormal)
	label = labelBoundsText(bridge, face, label)
	white := color.RGBA{255, 255, 255, 255}
	if sprite, ok := centeredLabelSprite(bridge, face, label, white); ok {
		sprites = append(sprites, sprite)
		return style.ViewData{Sprites: sprites}
	}

	// Degraded fallback: emit the label as a row over the gradient backdrop.
	entry := ctx.Entry(tiercatalog.TierNormal)
	offsets := []style.TTYX{bridge.CenterX(len(label), layout.NewMetrics(entry.GlyphAdvance, entry.RowHeight, entry.FontID))}
	offsetY := bridge.CenterBlockY([]int{entry.RowHeight}, 0)
	return style.ViewData{
		Items:       []string{label},
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Sprites:     sprites,
	}
}

func monoSlowPoster800x480(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()

	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	}
	frameSprite := borderframe.Render(frameCfg)
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	smallEntry := ctx.Entry(tiercatalog.TierSmall)

	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	rawItems := source.BuildDefaultItems(snap)
	var header []rowSpec
	var details []rowSpec
	for i, item := range rawItems {
		if i == 0 {
			header = append(header, rowSpec{text: item, tier: tiercatalog.TierLarge, entry: largeEntry})
		} else {
			details = append(details, rowSpec{text: item, tier: tiercatalog.TierSmall, entry: smallEntry})
		}
	}

	groupGap := largeEntry.RowHeight / 2
	if groupGap < 8 {
		groupGap = 8
	}

	allRows := append(header, details...)
	includeDetails := len(details) > 0

	totalHeight := 0
	for _, r := range header {
		totalHeight += r.entry.RowHeight
	}
	for _, r := range details {
		totalHeight += r.entry.RowHeight
	}
	if len(details) > 0 {
		totalHeight += groupGap
	}

	visibleRows := allRows
	if totalHeight > bridge.AvailableContentHeight() {
		includeDetails = false
		visibleRows = header
		for i := len(details); i > 0; i-- {
			candidate := append(header, details[:i]...)
			h := 0
			for _, r := range candidate {
				h += r.entry.RowHeight
			}
			h += groupGap
			if h <= bridge.AvailableContentHeight() {
				visibleRows = candidate
				includeDetails = true
				break
			}
		}
	}

	for i := range visibleRows {
		ga := visibleRows[i].entry.GlyphAdvance
		if ga > 0 {
			maxChars := bridge.AvailableContentWidth() / ga
			if maxChars > 0 {
				visibleRows[i].text = textlayout.Truncate(visibleRows[i].text, maxChars)
			}
		}
	}

	offsets := make([]style.TTYX, len(visibleRows))
	for i, r := range visibleRows {
		offsets[i] = bridge.CenterX(len([]rune(r.text)), layout.NewMetrics(r.entry.GlyphAdvance, r.entry.RowHeight, r.entry.FontID))
	}

	var blockHeight int
	for _, r := range visibleRows {
		blockHeight += r.entry.RowHeight
	}
	if includeDetails && len(visibleRows) > len(header) {
		blockHeight += groupGap
	}
	offsetY := style.TTYY((bridge.AvailableContentHeight() - blockHeight) / 2)
	if offsetY < 0 {
		offsetY = 0
	}

	items := make([]string, len(visibleRows))
	tiers := make([]tiercatalog.Tier, len(visibleRows))
	for i, r := range visibleRows {
		items[i] = r.text
		tiers[i] = r.tier
	}

	return style.ViewData{
		PaddingPct: 2,
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Static:      true,
		Sprites:     sprites,
	}
}

// BootFraction exposes boot progress calculation for package-level tests.
func BootFraction(snap source.Snapshot) float64 { return bootFraction(snap) }
