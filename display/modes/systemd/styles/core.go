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
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/gradient"
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

	entry := ctx.Entry(tiercatalog.TierNormal)
	if d.p.ContentMaxRows {
		if entry.RowHeight > 0 {
			maxRows := bridge.AvailableContentHeight() / entry.RowHeight
			if maxRows > 0 && len(items) > maxRows {
				items = items[:maxRows]
			}
		}
	} else {
		maxRows := bridge.MaxVisibleRows()
		if maxRows > 0 && len(items) > maxRows {
			items = items[:maxRows]
		}
	}

	var colors []color.Color
	var tiers []tiercatalog.Tier
	if d.p.ColorText {
		white := color.RGBA{255, 255, 255, 255}
		colors = make([]color.Color, len(items))
		tiers = make([]tiercatalog.Tier, len(items))
		for i := range items {
			colors[i] = white
			tiers[i] = tiercatalog.TierNormal
		}
	}

	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = bridge.CenterXWith(len(item), entry.GlyphAdvance)
	}
	rowHeights := make([]int, len(items))
	for i := range items {
		rowHeights[i] = entry.RowHeight
	}
	offsetY := bridge.CenterBlockY(rowHeights, 0)

	return style.ViewData{
		Items:       items,
		Colors:      colors,
		Tiers:       tiers,
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Sprites:     sprites,
		Static:      d.p.Static,
	}
}

func singleLabelBuild(snap source.Snapshot, ctx style.StyleContext, d def) style.ViewData {
	bridge := ctx.Layout(0)

	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Items: []string{"(no data)"}}
	}

	label := source.BootStatusLabel(snap)
	entry := ctx.Entry(tiercatalog.TierNormal)
	offsets := []int{bridge.CenterXWith(len(label), entry.GlyphAdvance)}
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
		maxChars := bridge.AvailableContentWidth() / bridge.GlyphAdvance()
		if maxChars > 0 && len(label) > maxChars {
			label = label[:maxChars]
		}
	}

	items := []string{label}
	entry := ctx.Entry(tiercatalog.TierNormal)
	offsets := []int{bridge.CenterXWith(len(label), entry.GlyphAdvance)}
	rowHeights := []int{entry.RowHeight}
	offsetY := bridge.CenterBlockY(rowHeights, 0)

	return style.ViewData{
		Items:       items,
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

	label := source.BootStatusLabel(snap)
	entry := ctx.Entry(tiercatalog.TierNormal)
	offsets := []int{bridge.CenterXWith(len(label), entry.GlyphAdvance)}
	offsetY := bridge.CenterBlockY([]int{entry.RowHeight}, 0)

	return style.ViewData{
		Items:       []string{label},
		LineOffsets: offsets,
		OffsetY:     offsetY,
		Sprites:     comp.Sprites(),
	}
}

func einkPoster800x480(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
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

	offsets := make([]int, len(visibleRows))
	for i, r := range visibleRows {
		offsets[i] = bridge.CenterXWith(len([]rune(r.text)), r.entry.GlyphAdvance)
	}

	var blockHeight int
	for _, r := range visibleRows {
		blockHeight += r.entry.RowHeight
	}
	if includeDetails && len(visibleRows) > len(header) {
		blockHeight += groupGap
	}
	offsetY := (bridge.AvailableContentHeight() - blockHeight) / 2
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
