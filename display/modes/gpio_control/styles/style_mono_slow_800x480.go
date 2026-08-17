package styles

import (
	"fmt"
	"image"

	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
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
	p:    Params{BuildFn: monoSlow800x480Build},
}

func monoSlow800x480Build(snap source.Data, _ source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()

	// Step 1: Guard
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	// Step 2: Borderframe
	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	}
	frameSprite := borderframe.Render(frameCfg)
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	// Step 3: Tier resolution
	largeEntry := ctx.Entry(tiercatalog.TierLarge)
	smallEntry := ctx.Entry(tiercatalog.TierSmall)

	// Step 4: Row assembly
	type rowSpec struct {
		text  string
		tier  tiercatalog.Tier
		entry tiercatalog.Entry
	}

	var header []rowSpec
	var secondary []rowSpec

	if len(snap.Pins) == 0 {
		header = []rowSpec{{text: "(no pins)", tier: tiercatalog.TierLarge, entry: largeEntry}}
	} else {
		header = []rowSpec{{text: "GPIO CTRL", tier: tiercatalog.TierLarge, entry: largeEntry}}

		for j := snap.TopRow; j < len(snap.Pins); j++ {
			p := snap.Pins[j]
			lvl := "LO"
			if p.Level {
				lvl = "HI"
			}
			pinText := fmt.Sprintf("GPIO%-2d %-3s %s", p.Number, p.Mode.String(), lvl)
			var prefix string
			if j == snap.Cursor {
				prefix = "> "
			} else {
				prefix = "  "
			}
			secondary = append(secondary, rowSpec{
				text:  prefix + pinText,
				tier:  tiercatalog.TierSmall,
				entry: smallEntry,
			})
		}
	}

	// Step 5: GroupGap
	groupGap := largeEntry.RowHeight / 2
	if groupGap < 8 {
		groupGap = 8
	}

	// Step 6: Adaptive fitting — drop secondary from end
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

	// Step 7: Truncation — per-tier maxChars
	for i := range visibleRows {
		ga := visibleRows[i].entry.GlyphAdvance
		if ga > 0 {
			maxChars := bridge.AvailableContentWidth() / ga
			if maxChars > 0 {
				visibleRows[i].text = textlayout.Truncate(visibleRows[i].text, maxChars)
			}
		}
	}

	// Step 8: Centering — horizontal per-tier, vertical center
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

	// Step 9: Emit ViewData
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
