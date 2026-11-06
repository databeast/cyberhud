package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// ZMQ colour palette for the message feed. Older messages sit in a calm
// near-white; the newest message — the one most likely to matter on a live
// feed — gets a bright cyan accent so the eye lands on it. The empty/idle
// placeholder is dimmed so it recedes rather than posing as data.
var (
	colorZMQData   = color.RGBA{R: 0xC8, G: 0xC8, B: 0xC8, A: 0xFF}
	colorZMQNewest = color.RGBA{R: 0x4F, G: 0xD6, B: 0xFF, A: 0xFF}
	colorZMQIdle   = color.RGBA{R: 0x70, G: 0x70, B: 0x70, A: 0xFF}
)

// buildColorItemsOnly renders the same top-anchored message feed as
// [buildItemsOnly], then tints the rows: the newest (bottom-most) message gets
// the cyan accent, older messages a calm grey, and the empty-state placeholder a
// dim grey. Only colour profiles use it; mono/grayscale panels drop the Colors
// slice and render in their native foreground.
func buildColorItemsOnly(snap content.ZMQData, _ content.Policy, ctx style.StyleContext, cfg itemsConfig) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: cfg.paddingPct})

	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	entry := ctx.Entry(resolveTier(cfg.tier))
	glyphAdvance := entry.GlyphAdvance
	rowHeight := entry.RowHeight

	maxRows := fitRowBudget(bridge.AvailableContentHeight(), rowHeight)
	maxChars := 0
	if glyphAdvance > 0 {
		maxChars = bridge.AvailableContentWidth() / glyphAdvance
	}

	ox, _ := bridge.ContentOrigin()

	items, offsets := assembleRows(snap, maxRows, maxChars, ox)
	colors := colorizeRows(snap, items)

	tier := resolveTier(cfg.tier)
	rowHeights := make([]int, len(items))
	tiers := make([]tiercatalog.Tier, len(items))
	for i := range items {
		rowHeights[i] = rowHeight
		tiers[i] = tier
	}
	offsetY, visible := topAnchor(bridge, rowHeights)

	return style.ViewData{
		Items:        items,
		Colors:       colors,
		Tiers:        tiers,
		LineOffsets:  offsets,
		RowHeights:   rowHeights,
		OffsetY:      offsetY,
		VisibleCount: visible,
		PaddingPct:   cfg.paddingPct,
		// See buildItemsOnly: a feed is a non-interactive static tail.
		Cursor: -1,
		Static: true,
	}
}

// colorizeRows tints each feed row: the empty-state placeholder is dimmed, the
// newest (last) message gets the cyan accent, and older messages are calm grey.
func colorizeRows(snap content.ZMQData, items []string) []color.Color {
	colors := make([]color.Color, len(items))
	if len(snap.Lines) == 0 {
		for i := range items {
			colors[i] = colorZMQIdle
		}
		return colors
	}
	for i := range items {
		if i == len(items)-1 {
			colors[i] = colorZMQNewest
		} else {
			colors[i] = colorZMQData
		}
	}
	return colors
}
