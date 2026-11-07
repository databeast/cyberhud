package styles

import (
	"fmt"
	"image"
	"image/color"
	"unicode/utf8"

	"github.com/databeast/cyberhud/display/modes/storage/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// buildBarsStyle renders one row per filesystem: the mountpoint and a compact
// usage summary ("<used>/<total> NN%"), with a usage-colored progress bar in a
// fixed right-hand band. Mounts are already ordered most-used-first by the
// source, so the most urgent filesystem is on top and survives truncation on
// short panels.
//
// All geometry is computed against the real tier font metrics (via ctx.Entry),
// not the panel's default glyph advance, and every row is pinned to that tier so
// the width the text is measured at is the width the renderer draws it at. This
// is the layout contract on style.ViewData: a bar placed from a mismatched
// advance overlaps the text.
func buildBarsStyle(snapshot source.StorageSnapshot, pol source.Policy, ctx style.StyleContext, _ def) style.ViewData {
	if len(snapshot.Mounts) == 0 {
		return style.ViewData{Items: []string{"no storage data"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 4})

	entry := ctx.Entry(tiercatalog.TierNormal)
	glyphAdvance := entry.GlyphAdvance
	if glyphAdvance <= 0 {
		glyphAdvance = textlayout.GlyphAdvance
	}
	rowHeight := entry.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}

	rowMetrics := layout.NewMetrics(glyphAdvance, rowHeight, entry.FontID)
	ox, _ := bridge.ContentOrigin()
	contentWidth := bridge.AvailableContentWidth()

	// How many rows fit at this tier's row height.
	maxRows := 0
	if rowHeight > 0 {
		maxRows = bridge.AvailableContentHeight() / rowHeight
	}
	mounts := snapshot.Mounts
	if maxRows > 0 && len(mounts) > maxRows {
		mounts = mounts[:maxRows]
	}
	if len(mounts) == 0 {
		mounts = snapshot.Mounts[:1]
	}

	isColor := isColorPanel(ctx)
	barBG := barBackground(isColor)
	nativeFG := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Reserve a fixed right-hand band for the bar whenever the panel is wide
	// enough to show both a readable bar and at least a few characters of text
	// beside it. The band is a fraction of the width, so its left edge never
	// depends on the (variable) rendered text width.
	const minBarBand = 24 // px
	const gutter = 2      // px between text and bar
	barBandWidth := 0
	textBandWidth := contentWidth
	// Need room for ~6 chars of text plus a min bar plus the gutter.
	if contentWidth >= 6*glyphAdvance+minBarBand+gutter {
		barBandWidth = contentWidth * 2 / 5
		if barBandWidth < minBarBand {
			barBandWidth = minBarBand
		}
		textBandWidth = contentWidth - barBandWidth - gutter
	}

	textMaxChars := 0
	if glyphAdvance > 0 && textBandWidth > 0 {
		textMaxChars = textBandWidth / glyphAdvance
	}

	// Vertically center the block of rows at this row height.
	offsetY := bridge.CenterBlockY(uniformHeights(len(mounts), rowHeight), 0)

	vd := style.ViewData{
		OffsetY:    offsetY,
		PaddingPct: 4,
		
		Tiers:      make([]tiercatalog.Tier, 0, len(mounts)),
		RowHeights: make([]int, 0, len(mounts)),
	}

	for i, m := range mounts {
		fill := m.Percent()
		pct := int(fill*100 + 0.5)

		label := shortMount(m.Mountpoint)
		sizes := fmt.Sprintf("%s/%s", source.FormatBytes(m.UsedBytes, pol.Unit), source.FormatBytes(m.TotalBytes, pol.Unit))
		row := composeRow(label, sizes, pct, textMaxChars)

		vd.Items = append(vd.Items, row)
		vd.Tiers = append(vd.Tiers, tiercatalog.TierNormal)
		vd.RowHeights = append(vd.RowHeights, rowHeight)
		if isColor {
			vd.Colors = append(vd.Colors, usageColor(fill, isColor))
		} else {
			vd.Colors = append(vd.Colors, nativeFG)
		}

		if barBandWidth <= 0 {
			continue
		}
		barHeight := rowHeight - 2
		if barHeight < 2 {
			barHeight = 2
		}
		barX := ox + contentWidth - barBandWidth
		yPos := bridge.CellToPixel(i, 0, rowMetrics).Y + int(offsetY) + (rowHeight-barHeight)/2

		barResult := progressbar.Render(progressbar.Config{
			Style:      progressbar.Linear,
			Value:      fill,
			Bounds:     image.Rect(barX, yPos, barX+barBandWidth, yPos+barHeight),
			Foreground: usageColor(fill, isColor),
			Background: barBG,
		})
		if barResult != nil {
			barResult.Label = fmt.Sprintf("storage-bar-%d", i)
			vd.Sprites = append(vd.Sprites, *barResult)
		}
	}

	if len(vd.Items) == 0 {
		vd.Items = []string{"storage"}
	}
	return vd
}

// composeRow fits the fullest form of a mount row into maxChars, degrading in
// steps: "<label> <sizes> NN%" → "<label> NN%" → truncated "<label> NN%" →
// truncated label. maxChars <= 0 means "no width budget known" — return the
// full form.
func composeRow(label, sizes string, pct, maxChars int) string {
	full := fmt.Sprintf("%s %s %d%%", label, sizes, pct)
	if maxChars <= 0 || utf8.RuneCountInString(full) <= maxChars {
		return full
	}
	// Drop the sizes, keep label + percent.
	compact := fmt.Sprintf("%s %d%%", label, pct)
	if utf8.RuneCountInString(compact) <= maxChars {
		return compact
	}
	// Truncate, preserving the trailing percent which is the key fact.
	pctStr := fmt.Sprintf(" %d%%", pct)
	keep := maxChars - utf8.RuneCountInString(pctStr)
	if keep >= 1 {
		return textlayout.Truncate(label, keep) + pctStr
	}
	return textlayout.Truncate(label, maxChars)
}

// shortMount compacts a mountpoint for display. The root filesystem shows as
// "/"; deep paths keep the trailing component so each row stays narrow.
func shortMount(mountpoint string) string {
	if mountpoint == "/" || mountpoint == "" {
		return "/"
	}
	for i := len(mountpoint) - 1; i >= 0; i-- {
		if mountpoint[i] == '/' {
			if i == len(mountpoint)-1 {
				continue // trailing slash
			}
			return mountpoint[i+1:]
		}
	}
	return mountpoint
}

func uniformHeights(n, h int) []int {
	if n <= 0 {
		return nil
	}
	out := make([]int, n)
	for i := range out {
		out[i] = h
	}
	return out
}
