package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// buildCompactStyle is the shared "compact" variant builder — the name tag.
//
// On a short strip (1-2 rows) it is a single dense "hostname: ip" line, the
// dignified label you would stick on the front of a rack unit. On a roomier
// panel the same one line stranded in the top-left reads as abandoned space, so
// compact instead presents a centred, accented two-line block (hostname above
// its primary address) that occupies the panel with intent. Emptiness, where it
// remains, is balanced around the block rather than dumped below it.
func buildCompactStyle(snapshot source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	_ = pol
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	isColor := isColorPanel(ctx)

	host := snapshot.Hostname
	ip := "(none)"
	if len(snapshot.IPs) > 0 {
		ip = snapshot.IPs[0]
	}

	maxChars := 0
	if bridge.GlyphAdvance() > 0 {
		maxChars = bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	}
	fit := func(s string) string {
		if maxChars > 0 {
			return textlayout.Truncate(s, maxChars)
		}
		return s
	}

	visibleRows := bridge.MaxVisibleRows()

	// Short strip: one dense line, as before.
	if visibleRows <= 2 {
		line := fit(host + ": " + ip)
		var colors []color.Color
		if isColor {
			colors = []color.Color{defHostColor}
		}
		return style.ViewData{Items: []string{line}, Colors: colors}
	}

	// Roomier panel: centred two-line block so the space reads as authored.
	rowHeight := bridge.RowHeight()
	availH := bridge.AvailableContentHeight()
	blockH := 2 * rowHeight
	offsetY := 0
	if availH > blockH {
		offsetY = (availH - blockH) / 2
	}

	items := []string{fit(host), fit(ip)}
	var colors []color.Color
	if isColor {
		colors = []color.Color{defHostColor, defLabelColor}
	}

	return style.ViewData{
		Items:      items,
		Colors:     colors,
		OffsetY:    offsetY,
		RowHeights: uniformRowHeights(2, rowHeight),
	}
}
