package styles

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
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

	m := ctx.DefaultMetrics()
	cw := bridge.AvailableContentWidth()
	visibleRows := bridge.MaxVisibleRows(m)

	// fitRow sizes a row down from the default tier until its whole text fits the
	// width (never truncated — a clipped "myhost: 192.16…" is useless, and the
	// ellipsis would tofu on the bitmap fonts). Returns the text, its tier, and
	// that tier's row height.
	fitRow := func(s string) (string, tiercatalog.Tier, int) {
		tier, ok := ctx.FitTier(s, cw, tiercatalog.TierNormal)
		if !ok {
			if all := ctx.FontCatalog().Tiers(); len(all) > 0 {
				tier = all[0]
			}
		}
		return s, tier, ctx.Metrics(tier).RowHeight()
	}

	// Short strip: one dense line.
	if visibleRows <= 2 {
		line, tier, _ := fitRow(host + ": " + ip)
		var colors []color.Color
		if isColor {
			colors = []color.Color{defHostColor}
		}
		return style.ViewData{Items: []string{line}, Colors: colors, Tiers: []tiercatalog.Tier{tier}}
	}

	// Roomier panel: centred two-line block so the space reads as authored.
	hostText, hostTier, hostH := fitRow(host)
	ipText, ipTier, ipH := fitRow(ip)
	blockH := hostH + ipH
	availH := bridge.AvailableContentHeight()
	offsetY := style.TTYY(0)
	if availH > blockH {
		offsetY = style.TTYY((availH - blockH) / 2)
	}

	var colors []color.Color
	if isColor {
		colors = []color.Color{defHostColor, defLabelColor}
	}

	return style.ViewData{
		Items:      []string{hostText, ipText},
		Colors:     colors,
		Tiers:      []tiercatalog.Tier{hostTier, ipTier},
		OffsetY:    offsetY,
		RowHeights: []int{hostH, ipH},
	}
}
