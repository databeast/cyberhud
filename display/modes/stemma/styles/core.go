package styles

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

type layoutKind uint8

const (
	layoutSummary layoutKind = iota
	layoutList
	layoutPoster
)

// Params holds the hand-tweakable knobs for a STEMMA style declaration.
// Layout selects a shared STEMMA layout; BuildFn swaps in a bespoke layout for a
// single style when a mode needs something unique.
type Params struct {
	Layout       layoutKind
	HeaderText   string
	RowFormatter func(*source.Device, int) string
	UseBorder    bool
	BorderTheme  string
	BuildFn      func(data source.StemmaSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData
}

type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[source.StemmaSnapshot, source.Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data source.StemmaSnapshot, pol source.Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	switch d.p.Layout {
	case layoutSummary:
		return summaryBuild(data)
	case layoutPoster:
		return posterBuild(data, ctx, d.p)
	default:
		return listBuild(data, ctx, d.p)
	}
}

func summaryBuild(snapshot source.StemmaSnapshot) style.ViewData {
	if len(snapshot.Devices) == 0 {
		return style.ViewData{Items: []string{"(no devices found)"}, Static: true}
	}
	presentCount := 0
	for _, d := range snapshot.Devices {
		if d.Present {
			presentCount++
		}
	}
	return style.ViewData{Items: []string{summaryLine(presentCount, len(snapshot.Devices))}, Static: true}
}

func listBuild(snapshot source.StemmaSnapshot, ctx style.StyleContext, p Params) style.ViewData {
	if len(snapshot.Devices) == 0 {
		return style.ViewData{Items: []string{"(no devices found)"}, Static: true}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	maxRows := bridge.MaxVisibleRows()
	maxChars := 0
	if bridge.GlyphAdvance() > 0 {
		maxChars = bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	}
	iconCells := iconWidthInGlyphs(bridge.GlyphAdvance())

	visibleDevs := snapshot.Devices
	if maxRows > 0 && len(visibleDevs) > maxRows {
		visibleDevs = visibleDevs[:maxRows]
	}

	nameWidth := maxChars - iconCells
	if nameWidth < 1 {
		nameWidth = 1
	}
	items := make([]string, len(visibleDevs))
	for i, d := range visibleDevs {
		items[i] = textlayout.Truncate(formatDeviceRow(d, nameWidth, p.RowFormatter), nameWidth)
	}

	vd := style.ViewData{Items: items, Colors: BuildColors(visibleDevs, ColorPresent, ColorAbsent)}
	if snapshot.GetIcon != nil {
		ox, oy := bridge.ContentOrigin()
		rowH := bridge.RowHeight()
		suppCtx := widgets.SuppressionContext{
			IsEink:          false,
			AvailableWidth:  bridge.AvailableContentWidth(),
			AvailableHeight: bridge.AvailableContentHeight(),
		}
		comp := widgets.NewCompositor(suppCtx)
		for i, dev := range visibleDevs {
			comp.Add(&deviceIconRenderable{dev: dev, rowIndex: i, rowH: rowH, getIcon: snapshot.GetIcon, originX: ox, originY: oy})
		}
		vd.Sprites = comp.Sprites()
	}
	return vd
}

func posterBuild(snapshot source.StemmaSnapshot, ctx style.StyleContext, p Params) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	var sprites []widgets.Sprite
	if p.UseBorder {
		frameSprite := borderframe.Render(borderframe.Config{
			Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
			Theme:  p.BorderTheme,
		})
		if frameSprite != nil {
			sprites = append(sprites, *frameSprite)
		}
	}

	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	smallEntry := ctx.Entry(tiercatalog.TierSmall)
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}
	var header []rowSpec
	var secondary []rowSpec
	if len(snapshot.Devices) == 0 {
		header = []rowSpec{{text: "(no devices found)", tier: tiercatalog.TierLarge, entry: largeEntry}}
	} else {
		headerText := p.HeaderText
		if headerText == "" {
			headerText = "STEMMA QT"
		}
		header = []rowSpec{{text: headerText, tier: tiercatalog.TierLarge, entry: largeEntry}}
		for _, dev := range snapshot.Devices {
			secondary = append(secondary, rowSpec{text: formatDeviceRow(dev, 24, p.RowFormatter), tier: tiercatalog.TierSmall, entry: smallEntry})
		}
	}

	groupGap := largeEntry.RowHeight / 2
	if groupGap < 8 {
		groupGap = 8
	}
	allRows := append(header, secondary...)
	includeSecondary := len(secondary) > 0
	totalHeight := 0
	for _, r := range header {
		totalHeight += r.entry.RowHeight
	}
	for _, r := range secondary {
		totalHeight += r.entry.RowHeight
	}
	if includeSecondary {
		totalHeight += groupGap
	}

	visibleRows := allRows
	if totalHeight > bridge.AvailableContentHeight() {
		includeSecondary = false
		visibleRows = header
		for i := len(secondary); i > 0; i-- {
			candidate := append(header, secondary[:i]...)
			h := 0
			for _, r := range candidate {
				h += r.entry.RowHeight
			}
			h += groupGap
			if h <= bridge.AvailableContentHeight() {
				visibleRows = candidate
				includeSecondary = true
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
	if includeSecondary && len(visibleRows) > len(header) {
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
	return style.ViewData{Items: items, Tiers: tiers, LineOffsets: offsets, OffsetY: offsetY, Static: true, Sprites: sprites}
}
