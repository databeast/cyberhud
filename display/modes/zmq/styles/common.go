package styles

import (
	"strings"

	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

// itemsConfig tunes the shared message-feed builder for a given panel profile.
//
// The ZMQ mode is a pure live-message feed: the messages are the product, so no
// row is spent on chrome. There is deliberately no header — a label naming the
// socket/endpoint/topic would only echo back configuration the user already
// chose, which is reference, not information, and wastes the most prominent row
// (and, on large panels, most of the surface).
type itemsConfig struct {
	// tier is the font size the feed renders at. Small panels use a compact tier
	// to fit more log lines; large panels step up so the feed fills the surface
	// and stays readable at a glance rather than floating as a small block in a
	// void.
	tier tiercatalog.Tier
	// paddingPct is the edge inset the style lays out against (e-paper profiles
	// use a small inset so content clears the bezel; OLED/TFT use 0).
	paddingPct int
}

// buildItemsOnly is the shared message-feed layout used by every monochrome and
// grayscale ZMQ style. It renders the tail of the message buffer (newest at the
// bottom), left-aligned, filling from the top of the content area. Colour
// profiles use [buildColorItemsOnly], which layers a palette on top of this
// same structure.
func buildItemsOnly(snap content.ZMQData, _ content.Policy, ctx style.StyleContext, cfg itemsConfig) style.ViewData {
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
		Tiers:        tiers,
		LineOffsets:  offsets,
		RowHeights:   rowHeights,
		OffsetY:      offsetY,
		VisibleCount: visible,
		PaddingPct:   cfg.paddingPct,
		// A message feed is a non-interactive tail, not a selectable list.
		// Cursor -1 suppresses the renderer's row-highlight bar (a zero Cursor
		// would highlight the first row), and Static renders it as a plain block
		// rather than routing through the cursor/scroll path.
		Cursor: -1,
		Static: true,
	}
}

// resolveTier returns the configured tier, defaulting to TierNormal.
func resolveTier(t tiercatalog.Tier) tiercatalog.Tier {
	if strings.TrimSpace(string(t)) == "" {
		return tiercatalog.TierNormal
	}
	return t
}

// fitRowBudget returns how many rows of the given (tier) row height fit in the
// available content height. It uses the real tier row height rather than the
// layout bridge's hint-derived default, because the two can differ — the bridge
// defaults to a 10px row when hints carry no font, while styles lay out at the
// catalog tier's actual height. Computing the budget against the height rows are
// actually drawn at keeps the item count consistent with the anchor math.
func fitRowBudget(availableHeight, rowHeight int) int {
	if rowHeight <= 0 || availableHeight <= 0 {
		return 1
	}
	rows := availableHeight / rowHeight
	if rows < 1 {
		rows = 1
	}
	return rows
}

// topAnchor lays the message rows out flush from the top of the content area
// rather than vertically centred. A live feed fills downward as messages
// arrive, so top-anchoring reads as "a terminal that will fill up" instead of a
// small block stranded in the middle of the panel. Returns the vertical offset
// (zero: content origin) and the visible row count (all supplied rows, since
// the budget already fit them).
func topAnchor(bridge layout.LayoutCalculator, rowHeights []int) (offsetY, visible int) {
	visible = len(rowHeights)
	// Guard against a rounding overflow of one row past the content area.
	total := 0
	for _, h := range rowHeights {
		total += h
	}
	for visible > 1 && total > bridge.AvailableContentHeight() {
		total -= rowHeights[visible-1]
		visible--
	}
	return 0, visible
}

// assembleRows produces the message rows (the tail of the buffer, newest at the
// bottom) and their per-row left-align offsets, or a single empty-state row when
// the buffer is empty. Shared by the mono/grayscale and colour builders so every
// profile agrees on what rows exist.
func assembleRows(snap content.ZMQData, maxRows, maxChars, ox int) (items []string, offsets []int) {
	lines := snap.Lines
	if len(lines) == 0 {
		empty := emptyStateText(snap)
		if maxChars > 0 {
			empty = textlayout.Truncate(empty, maxChars)
		}
		items = append(items, empty)
	} else {
		if len(lines) > maxRows {
			lines = lines[len(lines)-maxRows:]
		}
		for _, raw := range lines {
			text := raw
			if maxChars > 0 {
				text = textlayout.Truncate(text, maxChars)
			}
			items = append(items, text)
		}
	}

	offsets = make([]int, len(items))
	for i := range offsets {
		offsets[i] = ox
	}
	return items, offsets
}

// emptyStateText returns the placeholder shown when the buffer holds no
// messages. It reflects whether the mode is configured yet — the one case where
// naming the lack of an endpoint is genuinely informational.
func emptyStateText(snap content.ZMQData) string {
	if strings.TrimSpace(snap.Endpoint) == "" {
		return "no endpoint set"
	}
	return "waiting for messages"
}
