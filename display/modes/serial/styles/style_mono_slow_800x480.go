package styles

import (
	"fmt"
	"image"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

var MonoSlow800x480Style = def{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMonoSlow800x480Style},
}

// It uses the 9-step polish pipeline: guard → borderframe → layout → tier resolution →
// row assembly → adaptive fitting → truncation → centering → emit.
func buildMonoSlow800x480Style(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	hints := ctx.Hints()

	// Step 1 (Guard): LayoutBridge with PaddingPct:2, early return for zero dims.
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	// Step 2 (Borderframe): Theme "sharp", e-ink safe.
	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	}
	frameSprite := borderframe.Render(frameCfg)
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	// Step 3 (Tier resolution): TierLarge for header, TierSmall for output lines.
	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	smallEntry := ctx.Entry(tiercatalog.TierSmall)

	// Step 4 (Row assembly): header + secondary lines with ANSI stripping.
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	var header []rowSpec
	var secondary []rowSpec

	if !snap.Connected {
		header = []rowSpec{
			{text: "Disconnected", tier: tiercatalog.TierLarge, entry: largeEntry},
		}
	} else {
		header = []rowSpec{
			{text: fmt.Sprintf("%s %d", snap.Port, snap.Baud), tier: tiercatalog.TierLarge, entry: largeEntry},
		}
		// Iterate snap.Lines (ignore ScrollOffset), strip ANSI codes.
		for _, line := range snap.Lines {
			stripped := string(source.StripAnsi([]byte(line)))
			secondary = append(secondary, rowSpec{text: stripped, tier: tiercatalog.TierSmall, entry: smallEntry})
		}
	}

	// GroupGap: half header tier row height, minimum 8px.
	groupGap := largeEntry.RowHeight / 2
	if groupGap < 8 {
		groupGap = 8
	}

	// Step 5 (Adaptive fitting): drop secondary rows from end when height insufficient.
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
		// Drop secondary rows from end until fit.
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

	// Step 6 (Truncation): per-tier maxChars.
	for i := range visibleRows {
		ga := visibleRows[i].entry.GlyphAdvance
		if ga > 0 {
			maxChars := bridge.AvailableContentWidth() / ga
			if maxChars > 0 {
				visibleRows[i].text = textlayout.Truncate(visibleRows[i].text, maxChars)
			}
		}
	}

	// Step 7 (Centering): CenterXWith per-tier, vertical centering.
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

	// Step 8 (Emit): ViewData with Items, Tiers, LineOffsets, OffsetY, Static:true, Sprites.
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
