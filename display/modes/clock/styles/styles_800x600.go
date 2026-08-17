package styles

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/clock/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

// Clock style declarations for the 800x480 e-ink / large panels.
//
// Most entries are declarations over the core layouts in core.go. The
// MonoSlow800x480Style e-ink poster keeps a fully bespoke BuildFn showcasing
// per-display hand-tuning: TierHuge header with TierSmall secondary rows.

var GrayscaleSlow800x480Style = styleDef{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
}

// GrayscaleFast800x480Style is the grayscale fallback for the 800×480 color TFT panel.
var GrayscaleFast800x480Style = styleDef{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{Fast: true},
}

var ColorSlow800x480Style = styleDef{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
}

// ColorLarge800x480Style targets the 800×480 color TFT panel.
var ColorLarge800x480Style = styleDef{
	name: "color-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
}

// MonoSlow800x480Style targets the 800×480 landscape monochrome slow-refresh
// e-ink panel. Its bespoke poster layout uses TierHuge for the time header and
// TierSmall for date/weekday secondary rows, with the full 9-step polish
// pipeline: guard → borderframe → tier resolution → row assembly → adaptive
// fitting → truncation → centering → emit.
var MonoSlow800x480Style = styleDef{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: einkPoster800x480},
}

func einkPoster800x480(data source.ClockData, pol source.Policy, ctx style.StyleContext, d styleDef) style.ViewData {
	hints := ctx.Hints()

	// Step 1 (Guard): early return for zero/negative content dimensions.
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	// Step 2 (Borderframe): render decorative border sprite.
	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	}
	frameSprite := borderframe.Render(frameCfg)
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	// Step 3 (Tier resolution): TierHuge for header, TierSmall for secondary.
	hugeEntry := ctx.Entry(tiercatalog.TierHuge)
	smallEntry := ctx.Entry(tiercatalog.TierSmall)

	// Step 4 (Row assembly): header = time at TierHuge, secondary = date/weekday at TierSmall.
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	header := []rowSpec{
		{text: data.Time, tier: tiercatalog.TierHuge, entry: hugeEntry},
	}

	var secondary []rowSpec
	if data.Date != "" {
		secondary = append(secondary, rowSpec{text: data.Date, tier: tiercatalog.TierSmall, entry: smallEntry})
	}
	if data.Weekday != "" {
		secondary = append(secondary, rowSpec{text: data.Weekday, tier: tiercatalog.TierSmall, entry: smallEntry})
	}

	// Step 5 (GroupGap): half of TierHuge row height, minimum 8px.
	groupGap := hugeEntry.RowHeight / 2
	if groupGap < 8 {
		groupGap = 8
	}

	// Step 6 (Adaptive fitting): drop secondary rows from end until fit.
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

	// Step 7 (Truncation): per-tier maxChars.
	for i := range visibleRows {
		ga := visibleRows[i].entry.GlyphAdvance
		if ga > 0 {
			maxChars := bridge.AvailableContentWidth() / ga
			if maxChars > 0 {
				visibleRows[i].text = textlayout.Truncate(visibleRows[i].text, maxChars)
			}
		}
	}

	// Step 8 (Centering): horizontal per-tier, vertical center.
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

	// Step 9 (Emit): ViewData with Items, Tiers, LineOffsets, OffsetY, Static:true, Sprites.
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
