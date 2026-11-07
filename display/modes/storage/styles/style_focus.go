package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/storage/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// buildFocusStyle renders a single-filesystem readout for space-constrained
// panels: the most-used mountpoint with its percentage, its free space, and a
// full-width usage-colored bar beneath. On tiny panels this beats cramming
// multiple mounts into illegible rows.
//
// Geometry uses the real tier font metrics (ctx.Entry) and pins the rows to that
// tier so the bar, positioned a row below the text, never collides with it.
func buildFocusStyle(snapshot source.StorageSnapshot, pol source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	if len(snapshot.Mounts) == 0 {
		return style.ViewData{Items: []string{"no storage"}}
	}

	m := snapshot.Mounts[0]
	fill := m.Percent()
	pct := int(fill*100 + 0.5)

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 6})

	entry := ctx.Entry(tiercatalog.TierNormal)
	rowHeight := entry.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}

	ox, oy := bridge.ContentOrigin()
	contentWidth := bridge.AvailableContentWidth()
	contentHeight := bridge.AvailableContentHeight()

	isColor := isColorPanel(ctx)
	fg := usageColor(fill, isColor)
	nativeFG := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	label := shortMount(m.Mountpoint)
	rows := []string{
		fmt.Sprintf("%s %d%%", label, pct),
		fmt.Sprintf("free %s", source.FormatBytes(m.FreeBytes, pol.Unit)),
	}

	vd := style.ViewData{PaddingPct: 6}
	for range rows {
		vd.Tiers = append(vd.Tiers, tiercatalog.TierNormal)
		vd.RowHeights = append(vd.RowHeights, rowHeight)
		if isColor {
			vd.Colors = append(vd.Colors, fg)
		} else {
			vd.Colors = append(vd.Colors, nativeFG)
		}
	}
	vd.Items = rows

	// Full-width bar one row below the text, only if it fits the content area.
	if contentWidth >= 6 {
		barHeight := rowHeight - 2
		if barHeight < 2 {
			barHeight = 2
		}
		barY := bridge.CellToPixel(len(rows), 0, layout.NewMetrics(entry.GlyphAdvance, rowHeight, entry.FontID)).Y + 1
		if barY+barHeight <= oy+contentHeight {
			barResult := progressbar.Render(progressbar.Config{
				Style:      progressbar.Linear,
				Value:      fill,
				Bounds:     image.Rect(ox, barY, ox+contentWidth, barY+barHeight),
				Foreground: fg,
				Background: barBackground(isColor),
			})
			if barResult != nil {
				barResult.Label = "storage-focus-bar"
				vd.Sprites = append(vd.Sprites, *barResult)
			}
		}
	}

	return vd
}
